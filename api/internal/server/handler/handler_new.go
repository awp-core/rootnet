package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetPortfolio GET /api/users/{address}/portfolio — returns the full user profile
func (h *Handler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "address")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	chainID := h.resolveChainID(r)
	resp, err := h.svcGetPortfolio(r.Context(), chainID, normalizeAddr(raw))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// GetDelegates GET /api/users/{address}/delegates — returns agents bound to the user
func (h *Handler) GetDelegates(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "address")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	chainID := h.resolveChainID(r)
	result, err := h.svcGetDelegates(r.Context(), chainID, normalizeAddr(raw))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// ListWorknetsRanked GET /api/subnets/ranked — list worknets ranked by stake
func (h *Handler) ListWorknetsRanked(w http.ResponseWriter, r *http.Request) {
	chainID := h.resolveChainID(r)
	limit, offset := h.parsePageParams(r)
	result, err := h.svcListWorknetsRanked(r.Context(), chainID, int32(limit), int32(offset))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// ListWorknetAgents GET /api/subnets/{worknetId}/agents — worknet agent list
func (h *Handler) ListWorknetAgents(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	chainID := h.resolveChainID(r)
	limit, offset := h.parsePageParams(r)
	result, svcErr := h.svcListWorknetAgents(r.Context(), chainID, subnetID, int32(limit), int32(offset))
	if svcErr != nil {
		h.writeSvcError(w, svcErr)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// SearchWorknets GET /api/subnets/search?q=xxx — worknet search
func (h *Handler) SearchWorknets(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		h.writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}
	chainID := h.resolveChainID(r)
	limit, offset := h.parsePageParams(r)
	result, err := h.svcSearchWorknets(r.Context(), chainID, query, int32(limit), int32(offset))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetWorknetsByOwner GET /api/subnets/by-owner/{owner} — query worknets by owner
func (h *Handler) GetWorknetsByOwner(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "owner")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid owner address")
		return
	}
	chainID := h.resolveChainID(r)
	limit, offset := h.parsePageParams(r)
	result, err := h.svcGetWorknetsByOwner(r.Context(), chainID, normalizeAddr(raw), int32(limit), int32(offset))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetEpochDetail GET /api/emission/epochs/{epochId} — epoch detail
func (h *Handler) GetEpochDetail(w http.ResponseWriter, r *http.Request) {
	epochIDStr := chi.URLParam(r, "epochId")
	epochID, parseErr := strconv.ParseInt(epochIDStr, 10, 64)
	if parseErr != nil || epochID < 0 {
		h.writeError(w, http.StatusBadRequest, "invalid epochId")
		return
	}
	chainID := h.resolveChainID(r)
	result, err := h.svcGetEpochDetail(r.Context(), chainID, epochID)
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetStakePositionsGlobal GET /api/staking/user/{address}/positions/global — cross-chain positions
func (h *Handler) GetStakePositionsGlobal(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "address")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	result, err := h.svcGetPositionsGlobal(r.Context(), normalizeAddr(raw))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}
