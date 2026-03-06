package api

import (
	"net/http"

	"github.com/drews/basecoat/internal/auth"
	"github.com/drews/basecoat/internal/db"
)

type dashboardHandler struct {
	miniatures db.MiniatureRepository
}

func (h *dashboardHandler) get(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	stats, err := h.miniatures.GetDashboardStats(claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, stats)
}
