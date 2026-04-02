package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ListProposals returns a paginated list of governance proposals with optional status filter
func (h *Handler) ListProposals(w http.ResponseWriter, r *http.Request) {
	chainID := h.resolveChainID(r)
	limit, offset := h.parsePageParams(r)
	status := r.URL.Query().Get("status")

	result, err := h.svcListProposals(r.Context(), chainID, status, int32(limit), int32(offset))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetProposal returns details for a single governance proposal
func (h *Handler) GetProposal(w http.ResponseWriter, r *http.Request) {
	chainID := h.resolveChainID(r)
	proposalID := chi.URLParam(r, "proposalId")
	if proposalID == "" {
		h.writeError(w, http.StatusBadRequest, "missing proposalId parameter")
		return
	}
	result, err := h.svcGetProposal(r.Context(), chainID, proposalID)
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// GetTreasury returns the treasury contract address (same on all chains via CREATE2)
func (h *Handler) GetTreasury(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{
		"treasuryAddress": h.cfg.TreasuryAddress,
	})
}
