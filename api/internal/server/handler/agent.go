package handler

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/cortexia/rootnet/api/internal/ratelimit"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// agentInfoItem holds a single agent's info including binding data and stake
type agentInfoItem struct {
	Address  string `json:"address"`
	BoundTo  string `json:"boundTo"`
	Stake    string `json:"stake"`
}

// batchAgentInfoRequest is the request body for batch agent info queries
type batchAgentInfoRequest struct {
	Agents   []string    `json:"agents"`
	SubnetID json.Number `json:"subnetId"`
}

// GetAgentsByOwner returns all agents (addresses) bound to a given owner
func (h *Handler) GetAgentsByOwner(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "owner")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid owner address")
		return
	}
	owner := normalizeAddr(raw)

	agents, err := h.queries.GetUsersByBoundTo(r.Context(), gen.GetUsersByBoundToParams{
		BoundTo: owner,
		ChainID: h.cfg.ChainID,
	})
	if err != nil {
		h.logger.Error("failed to get agents by owner", "error", err, "owner", owner)
		h.writeError(w, http.StatusInternalServerError, "failed to get agents")
		return
	}

	h.writeJSON(w, http.StatusOK, agents)
}

// GetAgentDetail returns details for a single agent (user record with binding info)
func (h *Handler) GetAgentDetail(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "agent")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid agent address")
		return
	}
	agentAddr := normalizeAddr(raw)

	user, err := h.queries.GetUser(r.Context(), gen.GetUserParams{
		Address: agentAddr,
		ChainID: h.cfg.ChainID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.writeError(w, http.StatusNotFound, "agent not found")
			return
		}
		h.logger.Error("failed to get agent detail", "error", err, "agent", agentAddr)
		h.writeError(w, http.StatusInternalServerError, "failed to get agent detail")
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

// LookupAgent looks up the owner (boundTo) of an agent by agent address
func (h *Handler) LookupAgent(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "agent")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid agent address")
		return
	}
	agentAddr := normalizeAddr(raw)

	user, err := h.queries.GetUser(r.Context(), gen.GetUserParams{
		Address: agentAddr,
		ChainID: h.cfg.ChainID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.writeError(w, http.StatusNotFound, "agent not found")
			return
		}
		h.logger.Error("failed to lookup agent owner", "error", err, "agent", agentAddr)
		h.writeError(w, http.StatusInternalServerError, "failed to lookup agent owner")
		return
	}

	if user.BoundTo == "" {
		h.writeError(w, http.StatusNotFound, "agent not bound")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"ownerAddress": user.BoundTo})
}

// BatchAgentInfo returns batch agent info along with each agent's stake in the specified subnet
func (h *Handler) BatchAgentInfo(w http.ResponseWriter, r *http.Request) {
	ip := ratelimit.GetClientIP(r)
	if exceeded, err := h.limiter.CheckAndIncrement(r.Context(), "batch_agent_info", ip); exceeded {
		h.writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
		return
	} else if err != nil {
		h.logger.Error("batch_agent_info rate limit error", "error", err)
	}

	r.Body = http.MaxBytesReader(w, r.Body, 65536) // max ~100 addresses
	var req batchAgentInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}

	if len(req.Agents) == 0 {
		h.writeJSON(w, http.StatusOK, []agentInfoItem{})
		return
	}

	if len(req.Agents) > 100 {
		h.writeError(w, http.StatusBadRequest, "batch size exceeds limit (100)")
		return
	}

	subnetNum, err := parseSubnetIDString(req.SubnetID.String())
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "subnetId must be a positive integer")
		return
	}

	ctx := r.Context()

	// Validate and normalize all addresses upfront
	validAddrs := make([]string, 0, len(req.Agents))
	for _, addr := range req.Agents {
		if !isValidAddress(addr) {
			continue
		}
		validAddrs = append(validAddrs, normalizeAddr(addr))
	}

	if len(validAddrs) == 0 {
		h.writeJSON(w, http.StatusOK, []agentInfoItem{})
		return
	}

	// Batch query: fetch all users in one DB call
	users, err := h.queries.GetUsersBatch(ctx, gen.GetUsersBatchParams{
		ChainID:   h.cfg.ChainID,
		Addresses: validAddrs,
	})
	if err != nil {
		h.logger.Error("failed to batch get users", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get agent info")
		return
	}
	userMap := make(map[string]gen.GetUsersBatchRow, len(users))
	for _, u := range users {
		userMap[u.Address] = u
	}

	// Batch query: fetch all stakes in one DB call
	stakes, err := h.queries.GetAgentSubnetStakesBatch(ctx, gen.GetAgentSubnetStakesBatchParams{
		ChainID:  h.cfg.ChainID,
		Agents:   validAddrs,
		SubnetID: subnetNum,
	})
	if err != nil {
		h.logger.Error("failed to batch get stakes", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get agent info")
		return
	}
	stakeMap := make(map[string]string, len(stakes))
	for _, s := range stakes {
		if s.Total.Valid {
			stakeMap[s.AgentAddress] = s.Total.Int.String()
		}
	}

	// Assemble results preserving input order
	results := make([]agentInfoItem, 0, len(validAddrs))
	for _, addr := range validAddrs {
		user, ok := userMap[addr]
		if !ok {
			continue
		}
		item := agentInfoItem{
			Address: user.Address,
			BoundTo: user.BoundTo,
			Stake:   "0",
		}
		if s, ok := stakeMap[addr]; ok {
			item.Stake = s
		}
		results = append(results, item)
	}

	h.writeJSON(w, http.StatusOK, results)
}

// parseSubnetID parses the subnetId URL parameter into pgtype.Numeric.
// subnetId can exceed int64 range (e.g. (8453<<64)|1 = 77 bits), so we parse as big.Int.
func parseSubnetID(r *http.Request) (pgtype.Numeric, error) {
	raw := chi.URLParam(r, "subnetId")
	if raw == "" {
		return pgtype.Numeric{}, errors.New("missing subnetId parameter")
	}
	return parseSubnetIDString(raw)
}

// parseSubnetIDString converts a decimal string to a pgtype.Numeric, validating > 0.
func parseSubnetIDString(s string) (pgtype.Numeric, error) {
	id, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return pgtype.Numeric{}, errors.New("subnetId must be an integer")
	}
	if id.Sign() <= 0 {
		return pgtype.Numeric{}, errors.New("subnetId must be a positive integer")
	}
	return pgtype.Numeric{Int: id, Exp: 0, Valid: true}, nil
}
