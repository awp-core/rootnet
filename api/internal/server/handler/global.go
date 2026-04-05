package handler

import (
	"encoding/json"
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

// BatchResolveRecipients resolves recipients for multiple addresses in one on-chain call
func (h *Handler) BatchResolveRecipients(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 65536)
	var req struct {
		Addresses []string `json:"addresses"`
		ChainID   int64    `json:"chainId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.Addresses) == 0 {
		h.writeJSON(w, http.StatusOK, []any{})
		return
	}
	if len(req.Addresses) > 500 {
		h.writeError(w, http.StatusBadRequest, "batch size exceeds limit (500)")
		return
	}
	for _, a := range req.Addresses {
		if !isValidAddress(a) {
			h.writeError(w, http.StatusBadRequest, "invalid address in batch: "+a)
			return
		}
	}

	chainID := req.ChainID
	if chainID == 0 {
		chainID = h.defaultChainID()
	}
	cr := h.getChainReader(chainID)
	if cr == nil {
		h.writeError(w, http.StatusServiceUnavailable, "chain reader not available")
		return
	}

	normalized := make([]string, len(req.Addresses))
	for i, a := range req.Addresses {
		normalized[i] = normalizeAddr(a)
	}

	resolved, err := cr.BatchResolveRecipients(normalized)
	if err != nil {
		h.logger.Error("batch resolve failed", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to resolve recipients")
		return
	}

	results := make([]map[string]string, len(normalized))
	for i := range normalized {
		results[i] = map[string]string{
			"address":           normalized[i],
			"resolvedRecipient": resolved[i],
		}
	}
	h.writeJSON(w, http.StatusOK, results)
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
