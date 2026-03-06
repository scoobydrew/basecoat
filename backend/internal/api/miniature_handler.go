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

type miniatureHandler struct {
	miniatures db.MiniatureRepository
	paints     db.MiniaturePaintRepository
	images     db.ImageRepository
}

func (h *miniatureHandler) list(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	collectionID := chi.URLParam(r, "collectionID")

	minis, err := h.miniatures.ListByCollection(collectionID, claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if minis == nil {
		minis = []models.Miniature{}
	}
	writeJSON(w, http.StatusOK, minis)
}

func (h *miniatureHandler) create(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	collectionID := chi.URLParam(r, "collectionID")

	var req struct {
		Name     string `json:"name"`
		UnitType string `json:"unit_type"`
		Quantity int    `json:"quantity"`
		Notes    string `json:"notes"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.Quantity < 1 {
		req.Quantity = 1
	}

	now := time.Now()
	m := &models.Miniature{
		ID:           uuid.NewString(),
		CollectionID: collectionID,
		UserID:       claims.UserID,
		Name:         req.Name,
		UnitType:     req.UnitType,
		Quantity:     req.Quantity,
		Status:       models.StatusUnpainted,
		Notes:        req.Notes,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := h.miniatures.Create(m); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, m)
}

func (h *miniatureHandler) get(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	m, err := h.miniatures.GetByID(chi.URLParam(r, "id"))
	if err != nil || m == nil {
		writeError(w, http.StatusNotFound, "miniature not found")
		return
	}
	if m.UserID != claims.UserID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	// Hydrate paints and images
	m.Paints, _ = h.paints.ListByMiniature(m.ID)
	m.Images, _ = h.images.ListByMiniature(m.ID)

	writeJSON(w, http.StatusOK, m)
}

func (h *miniatureHandler) update(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	m, err := h.miniatures.GetByID(chi.URLParam(r, "id"))
	if err != nil || m == nil {
		writeError(w, http.StatusNotFound, "miniature not found")
		return
	}
	if m.UserID != claims.UserID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	var req struct {
		Name     *string                `json:"name"`
		UnitType *string                `json:"unit_type"`
		Quantity *int                   `json:"quantity"`
		Status   *models.PaintingStatus `json:"status"`
		Notes    *string                `json:"notes"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.UnitType != nil {
		m.UnitType = *req.UnitType
	}
	if req.Quantity != nil && *req.Quantity >= 1 {
		m.Quantity = *req.Quantity
	}
	if req.Status != nil {
		m.Status = *req.Status
	}
	if req.Notes != nil {
		m.Notes = *req.Notes
	}
	m.UpdatedAt = time.Now()

	if err := h.miniatures.Update(m); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (h *miniatureHandler) delete(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	m, err := h.miniatures.GetByID(chi.URLParam(r, "id"))
	if err != nil || m == nil {
		writeError(w, http.StatusNotFound, "miniature not found")
		return
	}
	if m.UserID != claims.UserID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	if err := h.miniatures.Delete(m.ID, claims.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *miniatureHandler) addPaint(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	m, err := h.miniatures.GetByID(chi.URLParam(r, "id"))
	if err != nil || m == nil || m.UserID != claims.UserID {
		writeError(w, http.StatusNotFound, "miniature not found")
		return
	}

	var req struct {
		PaintID string `json:"paint_id"`
		Purpose string `json:"purpose"`
		Notes   string `json:"notes"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.PaintID == "" {
		writeError(w, http.StatusBadRequest, "paint_id is required")
		return
	}

	mp := &models.MiniaturePaint{
		ID:          uuid.NewString(),
		MiniatureID: m.ID,
		PaintID:     req.PaintID,
		Purpose:     req.Purpose,
		Notes:       req.Notes,
		CreatedAt:   time.Now(),
	}
	if err := h.paints.Add(mp); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, mp)
}

func (h *miniatureHandler) removePaint(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if err := h.paints.Remove(chi.URLParam(r, "paintLinkID"), claims.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
