package handler

import (
	"errors"
	"math/big"
	"net/http"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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
	address := normalizeAddr(raw)

	ctx := r.Context()

	// Get total staked from stake_positions
	totalStakedNum, err := h.queries.GetUserTotalStaked(ctx, address)
	if err != nil {
		h.logger.Error("failed to get user total staked", "error", err, "address", address)
		h.writeError(w, http.StatusInternalServerError, "failed to get user balance")
		return
	}

	totalStaked := "0"
	if totalStakedNum.Valid {
		totalStaked = totalStakedNum.Int.String()
	}

	// Get total allocated from user_balances
	totalAllocated := "0"
	balance, err := h.queries.GetUserBalance(ctx, address)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			h.logger.Error("failed to get user balance", "error", err, "address", address)
			h.writeError(w, http.StatusInternalServerError, "failed to get user balance")
			return
		}
		// No balance record; totalAllocated stays 0
	} else {
		if balance.TotalAllocated.Valid {
			totalAllocated = balance.TotalAllocated.Int.String()
		}
	}

	// Compute unallocated = totalStaked - totalAllocated
	unallocated := "0"
	if totalStakedNum.Valid {
		stakedBig := totalStakedNum.Int
		allocBig := new(big.Int)
		if balance.TotalAllocated.Valid {
			allocBig.Set(balance.TotalAllocated.Int)
		}
		diff := new(big.Int).Sub(stakedBig, allocBig)
		unallocated = diff.String()
	}

	resp := balanceResponse{
		TotalStaked:    totalStaked,
		TotalAllocated: totalAllocated,
		Unallocated:    unallocated,
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
	address := normalizeAddr(rawAddr)

	positions, err := h.queries.GetUserStakePositions(r.Context(), address)
	if err != nil {
		h.logger.Error("failed to get stake positions", "error", err, "address", address)
		h.writeError(w, http.StatusInternalServerError, "failed to get stake positions")
		return
	}

	h.writeJSON(w, http.StatusOK, positions)
}

// GetAllocations returns a paginated list of stake allocations for a user
func (h *Handler) GetAllocations(w http.ResponseWriter, r *http.Request) {
	rawAddr := chi.URLParam(r, "address")
	if !isValidAddress(rawAddr) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	address := normalizeAddr(rawAddr)

	limit, offset := h.parsePageParams(r)

	allocations, err := h.queries.GetAllocationsByUser(r.Context(), gen.GetAllocationsByUserParams{
		UserAddress: address,
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		h.logger.Error("failed to get stake allocations", "error", err, "address", address)
		h.writeError(w, http.StatusInternalServerError, "failed to get stake allocations")
		return
	}

	h.writeJSON(w, http.StatusOK, allocations)
}

// GetPending returns pending reallocations; in dual-slot mode reallocate takes effect immediately so this is always empty
func (h *Handler) GetPending(w http.ResponseWriter, r *http.Request) {
	// In dual-slot mode reallocate takes effect immediately; no pending records to query
	h.writeJSON(w, http.StatusOK, []struct{}{})
}

// GetFrozen returns a user's frozen stake allocations
func (h *Handler) GetFrozen(w http.ResponseWriter, r *http.Request) {
	rawAddr := chi.URLParam(r, "address")
	if !isValidAddress(rawAddr) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	address := normalizeAddr(rawAddr)

	frozen, err := h.queries.GetFrozenByUser(r.Context(), address)
	if err != nil {
		h.logger.Error("failed to get frozen allocations", "error", err, "address", address)
		h.writeError(w, http.StatusInternalServerError, "failed to get frozen allocations")
		return
	}

	h.writeJSON(w, http.StatusOK, frozen)
}

// GetAgentSubnetStake returns the stake amount for an agent in a given subnet
func (h *Handler) GetAgentSubnetStake(w http.ResponseWriter, r *http.Request) {
	rawAgent := chi.URLParam(r, "agent")
	if !isValidAddress(rawAgent) {
		h.writeError(w, http.StatusBadRequest, "invalid agent address")
		return
	}
	agentAddr := normalizeAddr(rawAgent)

	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	stake, err := h.queries.GetAgentSubnetStake(r.Context(), gen.GetAgentSubnetStakeParams{
		AgentAddress: agentAddr,
		SubnetID:     subnetID,
	})
	if err != nil {
		h.logger.Error("failed to get agent subnet stake", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get agent subnet stake")
		return
	}

	amount := "0"
	if stake.Valid {
		amount = stake.Int.String()
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
	agentAddr := normalizeAddr(rawAgent)

	subnets, err := h.queries.GetAgentSubnets(r.Context(), agentAddr)
	if err != nil {
		h.logger.Error("failed to get agent subnets", "error", err, "agent", agentAddr)
		h.writeError(w, http.StatusInternalServerError, "failed to get agent subnets")
		return
	}

	h.writeJSON(w, http.StatusOK, subnets)
}

// GetSubnetTotalStake returns the total stake amount for a subnet
func (h *Handler) GetSubnetTotalStake(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	total, err := h.queries.GetSubnetTotalStake(r.Context(), subnetID)
	if err != nil {
		h.logger.Error("failed to get subnet total stake", "error", err, "subnetId", subnetID)
		h.writeError(w, http.StatusInternalServerError, "failed to get subnet total stake")
		return
	}

	amount := "0"
	if total.Valid {
		amount = total.Int.String()
	}
	h.writeJSON(w, http.StatusOK, map[string]string{"total": amount})
}
