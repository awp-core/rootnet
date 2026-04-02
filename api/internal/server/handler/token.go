package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetAWPInfo retrieves AWP token info from the Redis cache
func (h *Handler) GetAWPInfo(w http.ResponseWriter, r *http.Request) {
	chainID := h.resolveChainID(r)
	data, err := h.svcGetAWPInfo(r.Context(), chainID)
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, data)
}

// GetAlphaInfo retrieves subnet Alpha token info from the database
func (h *Handler) GetAlphaInfo(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, svcErr := h.svcGetAlphaInfo(r.Context(), subnetID)
	if svcErr != nil {
		h.writeSvcError(w, svcErr)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetAlphaPrice retrieves the Alpha token price from the Redis cache.
// Reads directly from alpha_price:{subnetId} — no DB lookup needed.
func (h *Handler) GetAlphaPrice(w http.ResponseWriter, r *http.Request) {
	subnetIDRaw := chi.URLParam(r, "subnetId")
	if subnetIDRaw == "" {
		h.writeError(w, http.StatusBadRequest, "missing subnetId parameter")
		return
	}
	if _, err := parseSubnetIDString(subnetIDRaw); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	data, svcErr := h.svcGetAlphaPrice(r.Context(), subnetIDRaw)
	if svcErr != nil {
		h.writeSvcError(w, svcErr)
		return
	}
	h.writeJSON(w, http.StatusOK, data)
}
