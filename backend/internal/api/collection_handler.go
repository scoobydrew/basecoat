package api

import (
	"net/http"
	"time"

	"github.com/drews/basecoat/internal/auth"
	"github.com/drews/basecoat/internal/claude"
	"github.com/drews/basecoat/internal/db"
	"github.com/drews/basecoat/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type collectionHandler struct {
	collections db.CollectionRepository
	miniatures  db.MiniatureRepository
	claude      *claude.Client
}

func (h *collectionHandler) list(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	cols, err := h.collections.ListByUser(claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if cols == nil {
		cols = []models.Collection{}
	}
	writeJSON(w, http.StatusOK, cols)
}

func (h *collectionHandler) create(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())

	var req struct {
		Name  string `json:"name"`
		Game  string `json:"game"`
		Set   string `json:"set"`
		Notes string `json:"notes"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Game == "" {
		writeError(w, http.StatusBadRequest, "name and game are required")
		return
	}

	col := &models.Collection{
		ID:        uuid.NewString(),
		UserID:    claims.UserID,
		Name:      req.Name,
		Game:      req.Game,
		Notes:     req.Notes,
		CreatedAt: time.Now(),
	}
	if err := h.collections.Create(col); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	// If a set was provided, ask Claude to seed the miniature list.
	var minis []models.Miniature
	if req.Set != "" && h.claude != nil {
		suggestions, err := h.claude.LookupMinis(r.Context(), req.Game, req.Set)
		if err == nil {
			for _, s := range suggestions {
				qty := s.Quantity
				if qty < 1 {
					qty = 1
				}
				m := &models.Miniature{
					ID:           uuid.NewString(),
					CollectionID: col.ID,
					UserID:       claims.UserID,
					Name:         s.Name,
					UnitType:     s.UnitType,
					Quantity:     qty,
					Status:       models.StatusUnpainted,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				if err := h.miniatures.Create(m); err == nil {
					minis = append(minis, *m)
				}
			}
		}
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"collection": col,
		"miniatures": minis,
	})
}

func (h *collectionHandler) get(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	col, err := h.collections.GetByID(chi.URLParam(r, "id"))
	if err != nil || col == nil {
		writeError(w, http.StatusNotFound, "collection not found")
		return
	}
	if col.UserID != claims.UserID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	writeJSON(w, http.StatusOK, col)
}

func (h *collectionHandler) update(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	col, err := h.collections.GetByID(chi.URLParam(r, "id"))
	if err != nil || col == nil {
		writeError(w, http.StatusNotFound, "collection not found")
		return
	}
	if col.UserID != claims.UserID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	var req struct {
		Name  string `json:"name"`
		Game  string `json:"game"`
		Notes string `json:"notes"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name != "" {
		col.Name = req.Name
	}
	if req.Game != "" {
		col.Game = req.Game
	}
	col.Notes = req.Notes

	if err := h.collections.Update(col); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, col)
}

func (h *collectionHandler) delete(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	col, err := h.collections.GetByID(chi.URLParam(r, "id"))
	if err != nil || col == nil {
		writeError(w, http.StatusNotFound, "collection not found")
		return
	}
	if col.UserID != claims.UserID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	if err := h.collections.Delete(col.ID, claims.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
