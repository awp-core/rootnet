package handler

import (
	"errors"
	"net/http"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// ListProposals returns a paginated list of governance proposals with optional status filter
func (h *Handler) ListProposals(w http.ResponseWriter, r *http.Request) {
	limit, offset := h.parsePageParams(r)
	status := r.URL.Query().Get("status")

	ctx := r.Context()

	if status != "" {
		// Validate status filter
		validStatuses := map[string]bool{"Active": true, "Canceled": true, "Defeated": true, "Succeeded": true, "Queued": true, "Expired": true, "Executed": true}
		if !validStatuses[status] {
			h.writeError(w, http.StatusBadRequest, "invalid status filter: must be one of Active, Canceled, Defeated, Succeeded, Queued, Expired, Executed")
			return
		}
		// Filter by status
		proposals, err := h.queries.ListProposalsByStatus(ctx, gen.ListProposalsByStatusParams{
			Status: status,
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			h.logger.Error("failed to list proposals by status", "error", err, "status", status)
			h.writeError(w, http.StatusInternalServerError, "failed to list proposals")
			return
		}
		h.writeJSON(w, http.StatusOK, proposals)
		return
	}

	// List all proposals
	proposals, err := h.queries.ListProposals(ctx, gen.ListProposalsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		h.logger.Error("failed to list proposals", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to list proposals")
		return
	}

	h.writeJSON(w, http.StatusOK, proposals)
}

// GetProposal returns details for a single governance proposal
func (h *Handler) GetProposal(w http.ResponseWriter, r *http.Request) {
	proposalID := chi.URLParam(r, "proposalId")
	if proposalID == "" {
		h.writeError(w, http.StatusBadRequest, "missing proposalId parameter")
		return
	}

	proposal, err := h.queries.GetProposal(r.Context(), proposalID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.writeError(w, http.StatusNotFound, "proposal not found")
			return
		}
		h.logger.Error("failed to get proposal", "error", err, "proposalId", proposalID)
		h.writeError(w, http.StatusInternalServerError, "failed to get proposal")
		return
	}

	h.writeJSON(w, http.StatusOK, proposal)
}

// GetTreasury returns the treasury contract address
func (h *Handler) GetTreasury(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{
		"treasuryAddress": h.cfg.TreasuryAddress,
	})
}
