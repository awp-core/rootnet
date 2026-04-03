package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ListSubnets returns a paginated list of subnets with optional status filter.
// Without ?chainId → returns subnets from ALL chains. With ?chainId=8453 → single chain.
func (h *Handler) ListSubnets(w http.ResponseWriter, r *http.Request) {
	// Subnet list defaults to cross-chain: returns subnets from all chains when chainId is not specified
	var chainID int64
	if v := r.URL.Query().Get("chainId"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil && id > 0 {
			chainID = id
		}
	}
	limit, offset := h.parsePageParams(r)
	status := r.URL.Query().Get("status")

	result, err := h.svcListSubnets(r.Context(), chainID, status, int32(limit), int32(offset))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetSubnet returns details for a single subnet
func (h *Handler) GetSubnet(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, svcErr := h.svcGetSubnet(r.Context(), subnetID)
	if svcErr != nil {
		h.writeSvcError(w, svcErr)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetSubnetEarnings returns a paginated AWP earnings history for a subnet (single JOIN query)
func (h *Handler) GetSubnetEarnings(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	limit, offset := h.parsePageParams(r)
	result, svcErr := h.svcGetSubnetEarnings(r.Context(), subnetID, int32(limit), int32(offset))
	if svcErr != nil {
		h.writeSvcError(w, svcErr)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
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
	result, svcErr := h.svcGetSubnetAgentInfo(r.Context(), normalizeAddr(rawAgent), subnetID)
	if svcErr != nil {
		h.writeSvcError(w, svcErr)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetSubnetSkills returns the skills URI for a subnet
func (h *Handler) GetSubnetSkills(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, svcErr := h.svcGetSubnetSkills(r.Context(), subnetID)
	if svcErr != nil {
		h.writeSvcError(w, svcErr)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}
