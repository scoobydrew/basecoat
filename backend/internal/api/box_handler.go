package api

import (
	"log"
	"net/http"
	"time"

	"github.com/drews/basecoat/internal/auth"
	"github.com/drews/basecoat/internal/claude"
	"github.com/drews/basecoat/internal/db"
	"github.com/drews/basecoat/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type boxHandler struct {
	boxes      db.BoxRepository
	games      db.GameRepository
	miniatures db.MiniatureRepository
	catalog    db.CatalogRepository
	claude     *claude.Client
}

func (h *boxHandler) list(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	gameID := chi.URLParam(r, "gameID")

	g, err := h.games.GetByID(gameID)
	if err != nil || g == nil || g.UserID != claims.UserID {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	boxes, err := h.boxes.ListByGame(gameID, claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if boxes == nil {
		boxes = []models.Box{}
	}
	writeJSON(w, http.StatusOK, boxes)
}

func (h *boxHandler) create(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	gameID := chi.URLParam(r, "gameID")

	g, err := h.games.GetByID(gameID)
	if err != nil || g == nil || g.UserID != claims.UserID {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	b := &models.Box{
		ID:        uuid.NewString(),
		GameID:    gameID,
		UserID:    claims.UserID,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	// Check catalog first — if found, skip Claude entirely.
	var suggestions []claude.MiniSuggestion
	var claudeErr string
	source := "none"

	if h.catalog != nil && g.CatalogGameID != "" {
		if cb, err := h.catalog.FindBox(g.CatalogGameID, req.Name); err != nil {
			log.Printf("catalog box lookup: %v", err)
		} else if cb != nil {
			b.CatalogBoxID = cb.ID
			minis, err := h.catalog.ListBoxMiniatures(cb.ID)
			if err != nil {
				log.Printf("catalog minis lookup: %v", err)
			} else {
				for _, m := range minis {
					suggestions = append(suggestions, claude.MiniSuggestion{
						Name:     m.Name,
						UnitType: m.UnitType,
						Quantity: m.Quantity,
					})
				}
				source = "catalog"
			}
		}
	}

	// Fall back to Claude if catalog had nothing.
	if source == "none" && h.claude != nil {
		s, err := h.claude.LookupMinis(r.Context(), g.Name, b.Name, claude.GameMeta{
			Publisher: g.Publisher,
			Year:      g.Year,
		})
		if err != nil {
			claudeErr = err.Error()
		} else {
			suggestions = s
			source = "claude"
		}
	}

	if err := h.boxes.Create(b); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"box":          b,
		"suggestions":  suggestions,
		"claude_error": claudeErr,
		"source":       source,
	})
}

// confirm batch-creates miniatures for a box and, if the box isn't already in
// the catalog, contributes it as a new catalog entry.
func (h *boxHandler) confirm(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())

	b, err := h.boxes.GetByID(chi.URLParam(r, "id"))
	if err != nil || b == nil || b.UserID != claims.UserID {
		writeError(w, http.StatusNotFound, "box not found")
		return
	}

	var req struct {
		Miniatures []struct {
			Name     string `json:"name"`
			UnitType string `json:"unit_type"`
			Quantity int    `json:"quantity"`
		} `json:"miniatures"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create user miniatures.
	now := time.Now()
	var created []models.Miniature
	for _, row := range req.Miniatures {
		if row.Name == "" {
			continue
		}
		qty := row.Quantity
		if qty < 1 {
			qty = 1
		}
		m := &models.Miniature{
			ID:        uuid.NewString(),
			BoxID:     b.ID,
			UserID:    claims.UserID,
			Name:      row.Name,
			UnitType:  row.UnitType,
			Quantity:  qty,
			Status:    models.StatusUnpainted,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := h.miniatures.Create(m); err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		created = append(created, *m)
	}

	// Contribute to catalog if this box isn't already there.
	if h.catalog != nil && b.CatalogBoxID == "" && len(req.Miniatures) > 0 {
		g, err := h.games.GetByID(b.GameID)
		if err != nil || g == nil {
			log.Printf("confirm: could not load game for catalog contribution: %v", err)
		} else {
			h.contributeTocatalog(b, g, req.Miniatures)
		}
	}

	writeJSON(w, http.StatusCreated, map[string]any{"miniatures": created})
}

// contributeToCatalog saves the box and its minis to the shared catalog.
// The game's catalog entry is guaranteed to exist before this is called.
// Errors are logged but not returned to the client — the user's data was already saved.
func (h *boxHandler) contributeToCatalog(b *models.Box, g *models.Game, rows []struct {
	Name     string `json:"name"`
	UnitType string `json:"unit_type"`
	Quantity int    `json:"quantity"`
}) {
	if g.CatalogGameID == "" {
		log.Printf("catalog: skipping box contribution — game %s has no catalog entry", g.ID)
		return
	}

	now := time.Now()

	cb := &models.CatalogBox{
		ID:            uuid.NewString(),
		CatalogGameID: g.CatalogGameID,
		Name:          b.Name,
		CreatedAt:     now,
	}
	if err := h.catalog.CreateBox(cb); err != nil {
		log.Printf("catalog: create box: %v", err)
		return
	}
	if err := h.boxes.SetCatalogBoxID(b.ID, cb.ID); err != nil {
		log.Printf("catalog: link box: %v", err)
	}

	for _, row := range rows {
		if row.Name == "" {
			continue
		}
		cm := &models.CatalogMiniature{
			ID:           uuid.NewString(),
			CatalogBoxID: cb.ID,
			Name:         row.Name,
			UnitType:     row.UnitType,
			Quantity:     max(row.Quantity, 1),
			CreatedAt:    now,
		}
		if err := h.catalog.CreateMiniature(cm); err != nil {
			log.Printf("catalog: create miniature %q: %v", row.Name, err)
		}
	}
}

func (h *boxHandler) get(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	b, err := h.boxes.GetByID(chi.URLParam(r, "id"))
	if err != nil || b == nil || b.UserID != claims.UserID {
		writeError(w, http.StatusNotFound, "box not found")
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (h *boxHandler) update(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	b, err := h.boxes.GetByID(chi.URLParam(r, "id"))
	if err != nil || b == nil || b.UserID != claims.UserID {
		writeError(w, http.StatusNotFound, "box not found")
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name != "" {
		b.Name = req.Name
	}
	if err := h.boxes.Update(b); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (h *boxHandler) delete(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	b, err := h.boxes.GetByID(chi.URLParam(r, "id"))
	if err != nil || b == nil || b.UserID != claims.UserID {
		writeError(w, http.StatusNotFound, "box not found")
		return
	}
	if err := h.boxes.Delete(b.ID, claims.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
