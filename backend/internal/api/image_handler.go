package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/drews/basecoat/internal/auth"
	"github.com/drews/basecoat/internal/db"
	"github.com/drews/basecoat/internal/models"
	"github.com/drews/basecoat/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const maxUploadSize = 10 << 20 // 10 MB

type imageHandler struct {
	images     db.ImageRepository
	miniatures db.MiniatureRepository
	storage    storage.Storage
}

func (h *imageHandler) upload(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	miniID := chi.URLParam(r, "id")

	mini, err := h.miniatures.GetByID(miniID)
	if err != nil || mini == nil || mini.UserID != claims.UserID {
		writeError(w, http.StatusNotFound, "miniature not found")
		return
	}

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, "file too large or invalid form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		writeError(w, http.StatusBadRequest, "image file is required")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		writeError(w, http.StatusBadRequest, "file must be an image")
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	stage := r.FormValue("stage")
	caption := r.FormValue("caption")

	imgID := uuid.NewString()
	key := fmt.Sprintf("users/%s/minis/%s/%s%s", claims.UserID, miniID, imgID, ext)

	if err := h.storage.Put(key, file, contentType); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to store image")
		return
	}

	img := &models.MiniatureImage{
		ID:          imgID,
		MiniatureID: miniID,
		UserID:      claims.UserID,
		Stage:       models.PaintingStatus(stage),
		StoragePath: key,
		Caption:     caption,
		CreatedAt:   time.Now(),
	}
	img.URL = h.storage.URL(key)

	if err := h.images.Create(img); err != nil {
		_ = h.storage.Delete(key)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, img)
}

func (h *imageHandler) delete(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	img, err := h.images.GetByID(chi.URLParam(r, "imageID"))
	if err != nil || img == nil {
		writeError(w, http.StatusNotFound, "image not found")
		return
	}
	if img.UserID != claims.UserID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	_ = h.storage.Delete(img.StoragePath)
	if err := h.images.Delete(img.ID, claims.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
