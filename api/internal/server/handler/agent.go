package handler

import (
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/cortexia/rootnet/api/internal/ratelimit"
	"github.com/go-chi/chi/v5"
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
	chainID := h.resolveChainID(r)
	result, err := h.svcGetAgentsByOwner(r.Context(), chainID, normalizeAddr(raw))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetAgentDetail returns details for a single agent (user record with binding info)
func (h *Handler) GetAgentDetail(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "agent")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid agent address")
		return
	}
	chainID := h.resolveChainID(r)
	result, err := h.svcGetAgentDetail(r.Context(), chainID, normalizeAddr(raw))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// LookupAgent looks up the owner (boundTo) of an agent by agent address
func (h *Handler) LookupAgent(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "agent")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid agent address")
		return
	}
	chainID := h.resolveChainID(r)
	owner, err := h.svcLookupAgent(r.Context(), chainID, normalizeAddr(raw))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]string{"ownerAddress": owner})
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

	r.Body = http.MaxBytesReader(w, r.Body, 65536)
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

	chainID := h.resolveChainID(r)
	results, svcErr := h.svcBatchAgentInfo(r.Context(), chainID, req.Agents, subnetNum)
	if svcErr != nil {
		h.writeSvcError(w, svcErr)
		return
	}
	h.writeJSON(w, http.StatusOK, results)
}

// parseSubnetID parses the subnetId URL parameter into pgtype.Numeric.
// subnetId can exceed int64 range (e.g. (8453<<64)|1 = 77 bits), so we parse as big.Int.
func parseSubnetID(r *http.Request) (pgtype.Numeric, error) {
	raw := chi.URLParam(r, "subnetId")
	if raw == "" {
		return pgtype.Numeric{}, errMissingSubnetID
	}
	return parseSubnetIDString(raw)
}

// errMissingSubnetID is a sentinel error for missing subnetId parameter
var errMissingSubnetID = &svcError{Kind: errBadInput, Message: "missing subnetId parameter"}

// parseSubnetIDString converts a decimal string to a pgtype.Numeric, validating > 0.
func parseSubnetIDString(s string) (pgtype.Numeric, error) {
	id, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return pgtype.Numeric{}, &svcError{Kind: errBadInput, Message: "subnetId must be an integer"}
	}
	if id.Sign() <= 0 {
		return pgtype.Numeric{}, &svcError{Kind: errBadInput, Message: "subnetId must be a positive integer"}
	}
	return pgtype.Numeric{Int: id, Exp: 0, Valid: true}, nil
}
