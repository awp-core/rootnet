package handler

import (
	"net/http"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/go-chi/chi/v5"
)

// userDetailResponse is the response type for user details including balance and bound agents (V2)
type userDetailResponse struct {
	User    gen.GetUserRow              `json:"user"`
	Balance *gen.GetUserBalanceRow      `json:"balance,omitempty"`
	Agents  []gen.GetUsersByBoundToRow  `json:"agents"`
}

// ListUsers returns a paginated list of users
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	chainID := h.resolveChainID(r)
	limit, offset := h.parsePageParams(r)
	result, err := h.svcListUsers(r.Context(), chainID, int32(limit), int32(offset))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetUserCount returns the total number of users
func (h *Handler) GetUserCount(w http.ResponseWriter, r *http.Request) {
	chainID := h.resolveChainID(r)
	count, err := h.svcGetUserCount(r.Context(), chainID)
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]int64{"count": count})
}

// GetUser returns details for a single user including balance and bound agents (V2)
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "address")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	chainID := h.resolveChainID(r)
	resp, err := h.svcGetUser(r.Context(), chainID, normalizeAddr(raw))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, resp)
}
