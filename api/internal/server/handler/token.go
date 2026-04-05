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

// GetAWPInfoGlobal returns AWP token info aggregated across all chains
func (h *Handler) GetAWPInfoGlobal(w http.ResponseWriter, r *http.Request) {
	data, err := h.svcGetAWPInfoGlobal(r.Context())
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, data)
}

// GetWorknetTokenInfo retrieves worknet token info from the database
func (h *Handler) GetWorknetTokenInfo(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, svcErr := h.svcGetWorknetTokenInfo(r.Context(), subnetID)
	if svcErr != nil {
		h.writeSvcError(w, svcErr)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetWorknetTokenPrice retrieves the worknet token price from the Redis cache.
// Reads directly from worknet_token_price:{worknetId} — no DB lookup needed.
func (h *Handler) GetWorknetTokenPrice(w http.ResponseWriter, r *http.Request) {
	subnetIDRaw := chi.URLParam(r, "worknetId")
	if subnetIDRaw == "" {
		h.writeError(w, http.StatusBadRequest, "missing worknetId parameter")
		return
	}
	if _, err := parseSubnetIDString(subnetIDRaw); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	data, svcErr := h.svcGetWorknetTokenPrice(r.Context(), subnetIDRaw)
	if svcErr != nil {
		h.writeSvcError(w, svcErr)
		return
	}
	h.writeJSON(w, http.StatusOK, data)
}
