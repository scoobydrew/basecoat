package api

import (
	"net/http"
	"time"

	"github.com/drews/basecoat/internal/auth"
	"github.com/drews/basecoat/internal/db"
	"github.com/drews/basecoat/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type collectionHandler struct {
	collections db.CollectionRepository
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
		Notes string `json:"notes"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	col := &models.Collection{
		ID:        uuid.NewString(),
		UserID:    claims.UserID,
		Name:      req.Name,
		Notes:     req.Notes,
		CreatedAt: time.Now(),
	}
	if err := h.collections.Create(col); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, col)
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
		Notes string `json:"notes"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name != "" {
		col.Name = req.Name
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
