package handler

import (
	"errors"
	"net/http"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// ListSubnets returns a paginated list of subnets with optional status filter
func (h *Handler) ListSubnets(w http.ResponseWriter, r *http.Request) {
	limit, offset := h.parsePageParams(r)
	status := r.URL.Query().Get("status")

	ctx := r.Context()

	if status != "" {
		// Validate status filter
		validStatuses := map[string]bool{"Pending": true, "Active": true, "Paused": true, "Banned": true}
		if !validStatuses[status] {
			h.writeError(w, http.StatusBadRequest, "invalid status filter: must be one of Pending, Active, Paused, Banned")
			return
		}
		// Filter by status
		subnets, err := h.queries.ListSubnetsByStatus(ctx, gen.ListSubnetsByStatusParams{
			ChainID: h.cfg.ChainID,
			Status:  status,
			Limit:   int32(limit),
			Offset:  int32(offset),
		})
		if err != nil {
			h.logger.Error("failed to list subnets by status", "error", err, "status", status)
			h.writeError(w, http.StatusInternalServerError, "failed to list subnets")
			return
		}
		h.writeJSON(w, http.StatusOK, subnets)
		return
	}

	// List all subnets
	subnets, err := h.queries.ListSubnets(ctx, gen.ListSubnetsParams{
		ChainID: h.cfg.ChainID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		h.logger.Error("failed to list subnets", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to list subnets")
		return
	}

	h.writeJSON(w, http.StatusOK, subnets)
}

// GetSubnet returns details for a single subnet
func (h *Handler) GetSubnet(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	subnet, err := h.queries.GetSubnet(r.Context(), subnetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.writeError(w, http.StatusNotFound, "subnet not found")
			return
		}
		h.logger.Error("failed to get subnet", "error", err, "subnetId", subnetID)
		h.writeError(w, http.StatusInternalServerError, "failed to get subnet")
		return
	}

	h.writeJSON(w, http.StatusOK, subnet)
}

// GetSubnetEarnings returns a paginated AWP earnings history for a subnet (single JOIN query)
func (h *Handler) GetSubnetEarnings(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	limit, offset := h.parsePageParams(r)

	earnings, err := h.queries.GetSubnetEarningsByID(r.Context(), gen.GetSubnetEarningsByIDParams{
		SubnetID: subnetID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		h.logger.Error("failed to get subnet earnings", "error", err, "subnetId", subnetID)
		h.writeError(w, http.StatusInternalServerError, "failed to get subnet earnings")
		return
	}

	h.writeJSON(w, http.StatusOK, earnings)
}

// GetSubnetAgentInfo returns the staking info for an agent in a given subnet
func (h *Handler) GetSubnetAgentInfo(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	rawAgent := chi.URLParam(r, "agent")
	if !isValidAddress(rawAgent) {
		h.writeError(w, http.StatusBadRequest, "invalid agent address")
		return
	}
	agentAddr := normalizeAddr(rawAgent)
	ctx := r.Context()

	total, err := h.queries.GetAgentSubnetStakeGlobal(ctx, gen.GetAgentSubnetStakeGlobalParams{
		AgentAddress: agentAddr,
		SubnetID:     subnetID,
	})
	stakeStr := "0"
	if err == nil && total.Valid {
		stakeStr = total.Int.String()
	}

	h.writeJSON(w, http.StatusOK, map[string]any{
		"agent":    agentAddr,
		"subnetId": subnetID,
		"stake":    stakeStr,
	})
}

// GetSubnetSkills returns the skills URI for a subnet
func (h *Handler) GetSubnetSkills(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	skillsURI, err := h.queries.GetSubnetSkills(r.Context(), subnetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.writeError(w, http.StatusNotFound, "subnet not found")
			return
		}
		h.logger.Error("failed to get subnet skills", "error", err, "subnetId", subnetID)
		h.writeError(w, http.StatusInternalServerError, "failed to get subnet skills")
		return
	}
	var uri string
	if skillsURI.Valid {
		uri = skillsURI.String
	}
	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"subnetId":  subnetID,
		"skillsURI": uri,
	})
}
