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

type paintHandler struct {
	paints db.PaintRepository
}

func (h *paintHandler) list(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	paints, err := h.paints.ListByUser(claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if paints == nil {
		paints = []models.Paint{}
	}
	writeJSON(w, http.StatusOK, paints)
}

func (h *paintHandler) create(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())

	var req struct {
		Brand     string `json:"brand"`
		Name      string `json:"name"`
		Color     string `json:"color"`
		PaintType string `json:"type"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Brand == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, "brand and name are required")
		return
	}

	p := &models.Paint{
		ID:        uuid.NewString(),
		UserID:    claims.UserID,
		Brand:     req.Brand,
		Name:      req.Name,
		Color:     req.Color,
		PaintType: req.PaintType,
		CreatedAt: time.Now(),
	}
	if err := h.paints.Create(p); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *paintHandler) update(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	p, err := h.paints.GetByID(chi.URLParam(r, "id"))
	if err != nil || p == nil {
		writeError(w, http.StatusNotFound, "paint not found")
		return
	}
	if p.UserID != claims.UserID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	var req struct {
		Brand     string `json:"brand"`
		Name      string `json:"name"`
		Color     string `json:"color"`
		PaintType string `json:"type"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Brand != "" {
		p.Brand = req.Brand
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	p.Color = req.Color
	p.PaintType = req.PaintType

	if err := h.paints.Update(p); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *paintHandler) delete(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	p, err := h.paints.GetByID(chi.URLParam(r, "id"))
	if err != nil || p == nil {
		writeError(w, http.StatusNotFound, "paint not found")
		return
	}
	if p.UserID != claims.UserID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	if err := h.paints.Delete(p.ID, claims.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
