package handler

import (
	"errors"
	"net/http"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// userDetailResponse is the response type for user details including balance and bound agents (V2)
type userDetailResponse struct {
	User    gen.GetUserRow              `json:"user"`
	Balance *gen.GetUserBalanceRow      `json:"balance,omitempty"`
	Agents  []gen.GetUsersByBoundToRow  `json:"agents"`
}

// ListUsers returns a paginated list of users
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, offset := h.parsePageParams(r)

	users, err := h.queries.ListUsers(r.Context(), gen.ListUsersParams{
		ChainID: h.cfg.ChainID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		h.logger.Error("failed to list users", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	h.writeJSON(w, http.StatusOK, users)
}

// GetUserCount returns the total number of users
func (h *Handler) GetUserCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.queries.GetUserCount(r.Context(), h.cfg.ChainID)
	if err != nil {
		h.logger.Error("failed to get user count", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get user count")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]int64{"count": count})
}

// GetUser returns details for a single user including balance and bound agents (V2)
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "address")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	address := normalizeAddr(raw)

	ctx := r.Context()

	// Fetch basic user info
	user, err := h.queries.GetUser(ctx, gen.GetUserParams{
		Address: address,
		ChainID: h.cfg.ChainID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.writeError(w, http.StatusNotFound, "user not found")
			return
		}
		h.logger.Error("failed to get user", "error", err, "address", address)
		h.writeError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	resp := userDetailResponse{
		User:   user,
		Agents: []gen.GetUsersByBoundToRow{},
	}

	// Fetch user balance
	if balance, err := h.queries.GetUserBalance(ctx, gen.GetUserBalanceParams{
		UserAddress: address,
		ChainID:     h.cfg.ChainID,
	}); err == nil {
		resp.Balance = &balance
	}

	// Fetch the user's bound agents (addresses where bound_to = this user)
	if agents, err := h.queries.GetUsersByBoundTo(ctx, gen.GetUsersByBoundToParams{
		BoundTo: address,
		ChainID: h.cfg.ChainID,
	}); err == nil {
		resp.Agents = agents
	}

	h.writeJSON(w, http.StatusOK, resp)
}
