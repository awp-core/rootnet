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
		// Filter by status
		subnets, err := h.queries.ListSubnetsByStatus(ctx, gen.ListSubnetsByStatusParams{
			Status: status,
			Limit:  int32(limit),
			Offset: int32(offset),
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
		Limit:  int32(limit),
		Offset: int32(offset),
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

// GetSubnetEarnings returns a paginated AWP earnings history for a subnet
func (h *Handler) GetSubnetEarnings(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	// Look up the subnet_contract address by subnet_id
	subnet, err := h.queries.GetSubnet(ctx, subnetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.writeError(w, http.StatusNotFound, "subnet not found")
			return
		}
		h.logger.Error("failed to get subnet", "error", err, "subnetId", subnetID)
		h.writeError(w, http.StatusInternalServerError, "failed to get subnet")
		return
	}

	limit, offset := h.parsePageParams(r)

	earnings, err := h.queries.GetRecipientEarnings(ctx, gen.GetRecipientEarningsParams{
		Recipient: subnet.SubnetContract,
		Limit:     int32(limit),
		Offset:    int32(offset),
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
	agentAddr := normalizeAddr(chi.URLParam(r, "agent"))
	if agentAddr == "" {
		h.writeError(w, http.StatusBadRequest, "missing agent parameter")
		return
	}
	ctx := r.Context()

	total, err := h.queries.GetAgentSubnetStake(ctx, gen.GetAgentSubnetStakeParams{
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
	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"subnetId":  subnetID,
		"skillsURI": skillsURI,
	})
}
