package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// batchAgentInfoRequest is the request body for batch agent info queries
type batchAgentInfoRequest struct {
	Agents   []string `json:"agents"`
	SubnetID int64    `json:"subnetId"`
}

// agentInfoItem holds a single agent's info including agent details and stake data
type agentInfoItem struct {
	Agent gen.Agent `json:"agent"`
	Stake string    `json:"stake"`
}

// GetAgentsByOwner returns all agents belonging to a given owner
func (h *Handler) GetAgentsByOwner(w http.ResponseWriter, r *http.Request) {
	owner := normalizeAddr(chi.URLParam(r, "owner"))
	if owner == "" {
		h.writeError(w, http.StatusBadRequest, "missing owner parameter")
		return
	}

	agents, err := h.queries.GetAgentsByOwner(r.Context(), owner)
	if err != nil {
		h.logger.Error("failed to get agents by owner", "error", err, "owner", owner)
		h.writeError(w, http.StatusInternalServerError, "failed to get agents")
		return
	}

	h.writeJSON(w, http.StatusOK, agents)
}

// GetAgentDetail returns details for a single agent
func (h *Handler) GetAgentDetail(w http.ResponseWriter, r *http.Request) {
	agentAddr := normalizeAddr(chi.URLParam(r, "agent"))
	if agentAddr == "" {
		h.writeError(w, http.StatusBadRequest, "missing agent parameter")
		return
	}

	agent, err := h.queries.GetAgent(r.Context(), agentAddr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.writeError(w, http.StatusNotFound, "agent not found")
			return
		}
		h.logger.Error("failed to get agent detail", "error", err, "agent", agentAddr)
		h.writeError(w, http.StatusInternalServerError, "failed to get agent detail")
		return
	}

	h.writeJSON(w, http.StatusOK, agent)
}

// LookupAgent looks up the owner of an agent by agent address
func (h *Handler) LookupAgent(w http.ResponseWriter, r *http.Request) {
	agentAddr := normalizeAddr(chi.URLParam(r, "agent"))
	if agentAddr == "" {
		h.writeError(w, http.StatusBadRequest, "missing agent parameter")
		return
	}

	owner, err := h.queries.LookupAgentOwner(r.Context(), agentAddr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.writeError(w, http.StatusNotFound, "agent not found or already removed")
			return
		}
		h.logger.Error("failed to lookup agent owner", "error", err, "agent", agentAddr)
		h.writeError(w, http.StatusInternalServerError, "failed to lookup agent owner")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"ownerAddress": owner})
}

// BatchAgentInfo returns batch agent info along with each agent's stake in the specified subnet
func (h *Handler) BatchAgentInfo(w http.ResponseWriter, r *http.Request) {
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
		agent, err := h.queries.GetAgent(ctx, addr)
		if err != nil {
			continue
		}

		item := agentInfoItem{
			Agent: agent,
			Stake: "0",
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
	id, err := strconv.Atoi(raw)
	if err != nil {
		return 0, errors.New("subnetId must be an integer")
	}
	return int64(id), nil
}
