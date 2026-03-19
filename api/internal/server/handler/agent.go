package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/cortexia/rootnet/api/internal/ratelimit"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// agentInfoItem holds a single agent's info including binding data and stake
type agentInfoItem struct {
	Address  string `json:"address"`
	BoundTo  string `json:"boundTo"`
	Stake    string `json:"stake"`
}

// batchAgentInfoRequest is the request body for batch agent info queries
type batchAgentInfoRequest struct {
	Agents   []string `json:"agents"`
	SubnetID int64    `json:"subnetId"`
}

// GetAgentsByOwner returns all agents (addresses) bound to a given owner
func (h *Handler) GetAgentsByOwner(w http.ResponseWriter, r *http.Request) {
	owner := normalizeAddr(chi.URLParam(r, "owner"))
	if owner == "" {
		h.writeError(w, http.StatusBadRequest, "missing owner parameter")
		return
	}

	agents, err := h.queries.GetUsersByBoundTo(r.Context(), owner)
	if err != nil {
		h.logger.Error("failed to get agents by owner", "error", err, "owner", owner)
		h.writeError(w, http.StatusInternalServerError, "failed to get agents")
		return
	}

	h.writeJSON(w, http.StatusOK, agents)
}

// GetAgentDetail returns details for a single agent (user record with binding info)
func (h *Handler) GetAgentDetail(w http.ResponseWriter, r *http.Request) {
	agentAddr := normalizeAddr(chi.URLParam(r, "agent"))
	if agentAddr == "" {
		h.writeError(w, http.StatusBadRequest, "missing agent parameter")
		return
	}

	user, err := h.queries.GetUser(r.Context(), agentAddr)
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
	agentAddr := normalizeAddr(chi.URLParam(r, "agent"))
	if agentAddr == "" {
		h.writeError(w, http.StatusBadRequest, "missing agent parameter")
		return
	}

	user, err := h.queries.GetUser(r.Context(), agentAddr)
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
	if exceeded, _ := h.limiter.CheckAndIncrement(r.Context(), "batch_agent_info", ip); exceeded {
		h.writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
		return
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

	ctx := r.Context()
	results := make([]agentInfoItem, 0, len(req.Agents))

	for _, addr := range req.Agents {
		addr = normalizeAddr(addr)
		user, err := h.queries.GetUser(ctx, addr)
		if err != nil {
			continue
		}

		item := agentInfoItem{
			Address: user.Address,
			BoundTo: user.BoundTo,
			Stake:   "0",
		}

		// Fetch the agent's stake in the specified subnet
		stake, err := h.queries.GetAgentSubnetStake(ctx, gen.GetAgentSubnetStakeParams{
			AgentAddress: addr,
			SubnetID:     req.SubnetID,
		})
		if err == nil && stake.Valid {
			item.Stake = stake.Int.String()
		}

		results = append(results, item)
	}

	h.writeJSON(w, http.StatusOK, results)
}

// parseSubnetID parses the subnetId URL parameter
func parseSubnetID(r *http.Request) (int64, error) {
	raw := chi.URLParam(r, "subnetId")
	if raw == "" {
		return 0, errors.New("missing subnetId parameter")
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, errors.New("subnetId must be an integer")
	}
	return id, nil
}
