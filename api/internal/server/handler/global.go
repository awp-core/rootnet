package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetGlobalStats returns aggregated protocol statistics across all chains
func (h *Handler) GetGlobalStats(w http.ResponseWriter, r *http.Request) {
	result, err := h.svcGetGlobalStats(r.Context())
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetUserBalanceGlobal returns a user's staking balance aggregated across all chains
func (h *Handler) GetUserBalanceGlobal(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "address")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	resp, err := h.svcGetUserBalanceGlobal(r.Context(), normalizeAddr(raw))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// GetGlobalEmissionSchedule returns emission schedule aggregated across all chains
func (h *Handler) GetGlobalEmissionSchedule(w http.ResponseWriter, r *http.Request) {
	result, err := h.svcGetGlobalEmissionSchedule(r.Context())
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// ListUsersGlobal returns a deduplicated cross-chain user list
func (h *Handler) ListUsersGlobal(w http.ResponseWriter, r *http.Request) {
	limit, offset := h.parsePageParams(r)
	result, err := h.svcListUsersGlobal(r.Context(), int32(limit), int32(offset))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// ListAllProposals returns governance proposals across all chains
func (h *Handler) ListAllProposals(w http.ResponseWriter, r *http.Request) {
	limit, offset := h.parsePageParams(r)
	status := r.URL.Query().Get("status")
	result, err := h.svcListAllProposals(r.Context(), status, int32(limit), int32(offset))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}
