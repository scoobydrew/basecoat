package api

import (
	"log"
	"net/http"
	"time"

	"github.com/drews/basecoat/internal/auth"
	"github.com/drews/basecoat/internal/db"
	"github.com/drews/basecoat/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type gameHandler struct {
	games       db.GameRepository
	collections db.CollectionRepository
	catalog     db.CatalogRepository
}

func (h *gameHandler) list(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	collectionID := chi.URLParam(r, "collectionID")

	col, err := h.collections.GetByID(collectionID)
	if err != nil || col == nil || col.UserID != claims.UserID {
		writeError(w, http.StatusNotFound, "collection not found")
		return
	}

	games, err := h.games.ListByCollection(collectionID, claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if games == nil {
		games = []models.Game{}
	}
	writeJSON(w, http.StatusOK, games)
}

func (h *gameHandler) create(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	collectionID := chi.URLParam(r, "collectionID")

	col, err := h.collections.GetByID(collectionID)
	if err != nil || col == nil || col.UserID != claims.UserID {
		writeError(w, http.StatusNotFound, "collection not found")
		return
	}

	var req struct {
		Name      string `json:"name"`
		Publisher string `json:"publisher"`
		Year      *int   `json:"year"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	g := &models.Game{
		ID:           uuid.NewString(),
		CollectionID: collectionID,
		UserID:       claims.UserID,
		Name:         req.Name,
		CreatedAt:    time.Now(),
	}

	// Always link to the shared catalog — find existing or create new.
	if h.catalog != nil {
		catalogGameID, err := h.findOrCreateCatalogGame(req.Name, req.Publisher, req.Year)
		if err != nil {
			log.Printf("catalog game find/create: %v", err)
		} else {
			g.CatalogGameID = catalogGameID
		}
	}

	if err := h.games.Create(g); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	// Re-fetch to populate publisher/year from the catalog JOIN.
	created, err := h.games.GetByID(g.ID)
	if err != nil || created == nil {
		writeJSON(w, http.StatusCreated, g)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *gameHandler) findOrCreateCatalogGame(name, publisher string, year *int) (string, error) {
	cg, err := h.catalog.FindGame(name, publisher)
	if err != nil {
		return "", err
	}
	if cg != nil {
		return cg.ID, nil
	}
	newCG := &models.CatalogGame{
		ID:        uuid.NewString(),
		Name:      name,
		Publisher: publisher,
		Year:      year,
		CreatedAt: time.Now(),
	}
	if err := h.catalog.CreateGame(newCG); err != nil {
		return "", err
	}
	return newCG.ID, nil
}

func (h *gameHandler) get(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	g, err := h.games.GetByID(chi.URLParam(r, "id"))
	if err != nil || g == nil || g.UserID != claims.UserID {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}
	writeJSON(w, http.StatusOK, g)
}

func (h *gameHandler) update(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	g, err := h.games.GetByID(chi.URLParam(r, "id"))
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
	if req.Name != "" {
		g.Name = req.Name
	}
	if err := h.games.Update(g); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, g)
}

func (h *gameHandler) delete(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	g, err := h.games.GetByID(chi.URLParam(r, "id"))
	if err != nil || g == nil || g.UserID != claims.UserID {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}
	if err := h.games.Delete(g.ID, claims.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
