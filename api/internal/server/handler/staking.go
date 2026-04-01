package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// balanceResponse is the response type for a user's staking balance
type balanceResponse struct {
	TotalStaked    string `json:"totalStaked"`
	TotalAllocated string `json:"totalAllocated"`
	Unallocated    string `json:"unallocated"`
}

// GetBalance returns a user's staking balance derived from stake_positions and user_balances
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "address")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	resp, err := h.svcGetBalance(r.Context(), normalizeAddr(raw))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// GetStakePositions returns a user's active stake NFT positions
func (h *Handler) GetStakePositions(w http.ResponseWriter, r *http.Request) {
	rawAddr := chi.URLParam(r, "address")
	if !isValidAddress(rawAddr) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	result, err := h.svcGetStakePositions(r.Context(), normalizeAddr(rawAddr))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetAllocations returns a paginated list of stake allocations for a user
func (h *Handler) GetAllocations(w http.ResponseWriter, r *http.Request) {
	rawAddr := chi.URLParam(r, "address")
	if !isValidAddress(rawAddr) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	limit, offset := h.parsePageParams(r)
	result, err := h.svcGetAllocations(r.Context(), normalizeAddr(rawAddr), int32(limit), int32(offset))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetPending returns pending reallocations; in dual-slot mode reallocate takes effect immediately so this is always empty
func (h *Handler) GetPending(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, []struct{}{})
}

// GetFrozen returns a user's frozen stake allocations
func (h *Handler) GetFrozen(w http.ResponseWriter, r *http.Request) {
	rawAddr := chi.URLParam(r, "address")
	if !isValidAddress(rawAddr) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	result, err := h.svcGetFrozen(r.Context(), normalizeAddr(rawAddr))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetAgentSubnetStake returns the stake amount for an agent in a given subnet
func (h *Handler) GetAgentSubnetStake(w http.ResponseWriter, r *http.Request) {
	rawAgent := chi.URLParam(r, "agent")
	if !isValidAddress(rawAgent) {
		h.writeError(w, http.StatusBadRequest, "invalid agent address")
		return
	}
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	amount, svcErr := h.svcGetAgentSubnetStake(r.Context(), normalizeAddr(rawAgent), subnetID)
	if svcErr != nil {
		h.writeSvcError(w, svcErr)
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]string{"amount": amount})
}

// GetAgentSubnets returns all subnets an agent participates in along with their stake amounts
func (h *Handler) GetAgentSubnets(w http.ResponseWriter, r *http.Request) {
	rawAgent := chi.URLParam(r, "agent")
	if !isValidAddress(rawAgent) {
		h.writeError(w, http.StatusBadRequest, "invalid agent address")
		return
	}
	result, err := h.svcGetAgentSubnets(r.Context(), normalizeAddr(rawAgent))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetSubnetTotalStake returns the total stake amount for a subnet
func (h *Handler) GetSubnetTotalStake(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	amount, svcErr := h.svcGetSubnetTotalStake(r.Context(), subnetID)
	if svcErr != nil {
		h.writeSvcError(w, svcErr)
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]string{"total": amount})
}
