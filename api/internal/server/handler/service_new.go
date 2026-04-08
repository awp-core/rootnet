package handler

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ── CRITICAL #1: users.getPortfolio ──

// portfolioResponse combines all user data in the response
type portfolioResponse struct {
	Address      string `json:"address"`
	IsRegistered bool   `json:"isRegistered"`
	BoundTo      string `json:"boundTo"`
	Recipient    string `json:"recipient"`

	Balance balanceResponse `json:"balance"`

	Positions   any `json:"positions"`
	Allocations any `json:"allocations"`
	Delegates   any `json:"delegates"`
}

// svcGetPortfolio returns full user profile: identity, balance, positions, allocations, delegates
func (h *Handler) svcGetPortfolio(ctx context.Context, chainID int64, address string) (*portfolioResponse, error) {
	resp := &portfolioResponse{Address: address}

	// 1. User identity
	if user, err := h.queries.GetUser(ctx, gen.GetUserParams{
		Address: address, ChainID: chainID,
	}); err == nil {
		resp.IsRegistered = user.RegisteredAt != 0 || user.BoundTo != "" || user.Recipient != ""
		resp.BoundTo = user.BoundTo
		resp.Recipient = user.Recipient
	} else if !errors.Is(err, pgx.ErrNoRows) {
		h.logger.Error("portfolio: failed to get user", "error", err, "address", address)
		return nil, newSvcErr(errInternal, "failed to get user")
	}

	// 2. Balance
	balance, err := h.svcGetBalance(ctx, chainID, address)
	if err != nil {
		return nil, err
	}
	resp.Balance = balance

	// 3. Positions
	positions, err := h.svcGetStakePositions(ctx, chainID, address)
	if err != nil {
		return nil, err
	}
	resp.Positions = positions

	// 4. Allocations (top 50)
	allocations, err := h.svcGetAllocations(ctx, chainID, address, 50, 0)
	if err != nil {
		return nil, err
	}
	resp.Allocations = allocations

	// 5. Delegates — query agents bound to this address from users table (approximation of delegates;
	//    the actual on-chain delegates mapping is not indexed in DB, only bound agents are returned)
	agents, err := h.queries.GetUsersByBoundTo(ctx, gen.GetUsersByBoundToParams{
		BoundTo: address, ChainID: chainID,
	})
	if err != nil {
		h.logger.Error("portfolio: failed to get delegates", "error", err, "address", address)
		resp.Delegates = []struct{}{}
	} else {
		resp.Delegates = agents
	}

	return resp, nil
}

// ── HIGH #4: subnets.listRanked ──

// rankedWorknetRow is a ranked worknet result row
type rankedWorknetRow struct {
	SubnetID    string `json:"worknetId"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Status      string `json:"status"`
	Owner       string `json:"owner"`
	TotalStake  string `json:"totalStake"`
	AgentCount  int64  `json:"agentCount"`
}

// svcListWorknetsRanked ranks worknets by total stake (raw SQL)
func (h *Handler) svcListWorknetsRanked(ctx context.Context, chainID int64, limit, offset int32) ([]rankedWorknetRow, error) {
	sql := `
		SELECT s.subnet_id, s.name, s.symbol, s.status, s.owner,
		       COALESCE(SUM(a.amount), 0) as total_stake,
		       COUNT(DISTINCT a.agent_address) as agent_count
		FROM subnets s
		LEFT JOIN stake_allocations a ON a.subnet_id = s.subnet_id AND a.chain_id = s.chain_id AND a.amount > 0 AND a.frozen = FALSE
		WHERE s.chain_id = $1 AND s.burned = FALSE
		GROUP BY s.subnet_id, s.name, s.symbol, s.status, s.owner
		ORDER BY total_stake DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := h.db.Query(ctx, sql, chainID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list ranked worknets", "error", err)
		return nil, newSvcErr(errInternal, "failed to list ranked worknets")
	}
	defer rows.Close()

	items := []rankedWorknetRow{}
	for rows.Next() {
		var r rankedWorknetRow
		var subnetID pgtype.Numeric
		var totalStake pgtype.Numeric
		if err := rows.Scan(&subnetID, &r.Name, &r.Symbol, &r.Status, &r.Owner, &totalStake, &r.AgentCount); err != nil {
			h.logger.Error("failed to scan ranked worknet row", "error", err)
			return nil, newSvcErr(errInternal, "failed to list ranked worknets")
		}
		if subnetID.Valid {
			r.SubnetID = numericString(subnetID)
		}
		if totalStake.Valid {
			r.TotalStake = numericString(totalStake)
		} else {
			r.TotalStake = "0"
		}
		items = append(items, r)
	}
	if err := rows.Err(); err != nil {
		h.logger.Error("failed to iterate ranked worknet rows", "error", err)
		return nil, newSvcErr(errInternal, "failed to list ranked worknets")
	}
	return items, nil
}

// ── HIGH #5: emission.getEpochDetail ──

// svcGetEpochDetail fetches distribution details for a given epoch (uses existing sqlc queries)
func (h *Handler) svcGetEpochDetail(ctx context.Context, chainID int64, epochID int64) (map[string]any, error) {
	// First fetch epoch basic info
	epoch, err := h.queries.GetEpoch(ctx, gen.GetEpochParams{ChainID: chainID, EpochID: epochID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newSvcErr(errNotFound, "epoch not found")
		}
		h.logger.Error("failed to get epoch", "error", err, "epochId", epochID)
		return nil, newSvcErr(errInternal, "failed to get epoch")
	}

	// Fetch distribution details
	distributions, err := h.queries.GetEpochDistributions(ctx, gen.GetEpochDistributionsParams{
		ChainID: chainID, EpochID: epochID,
	})
	if err != nil {
		h.logger.Error("failed to get epoch distributions", "error", err, "epochId", epochID)
		return nil, newSvcErr(errInternal, "failed to get epoch distributions")
	}

	dailyStr := "0"
	if epoch.DailyEmission.Valid {
		dailyStr = numericString(epoch.DailyEmission)
	}
	daoStr := "0"
	if epoch.DaoEmission.Valid {
		daoStr = numericString(epoch.DaoEmission)
	}

	return map[string]any{
		"epochId":       epoch.EpochID,
		"startTime":     epoch.StartTime,
		"dailyEmission": dailyStr,
		"daoEmission":   daoStr,
		"distributions": distributions,
	}, nil
}

// ── HIGH #6: subnets.listAgents ──

// worknetAgentRow is a worknet agent ranking result row
type worknetAgentRow struct {
	AgentAddress string `json:"agentAddress"`
	TotalStake   string `json:"totalStake"`
}

// svcListWorknetAgents returns the agent list within a worknet, ranked by stake
func (h *Handler) svcListWorknetAgents(ctx context.Context, chainID int64, subnetID pgtype.Numeric, limit, offset int32) ([]worknetAgentRow, error) {
	sql := `
		SELECT agent_address, SUM(amount) as total_stake
		FROM stake_allocations
		WHERE chain_id = $1 AND subnet_id = $2 AND amount > 0 AND frozen = FALSE
		GROUP BY agent_address
		ORDER BY total_stake DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := h.db.Query(ctx, sql, chainID, subnetID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list worknet agents", "error", err)
		return nil, newSvcErr(errInternal, "failed to list worknet agents")
	}
	defer rows.Close()

	items := []worknetAgentRow{}
	for rows.Next() {
		var r worknetAgentRow
		var totalStake pgtype.Numeric
		if err := rows.Scan(&r.AgentAddress, &totalStake); err != nil {
			h.logger.Error("failed to scan worknet agent row", "error", err)
			return nil, newSvcErr(errInternal, "failed to list worknet agents")
		}
		if totalStake.Valid {
			r.TotalStake = numericString(totalStake)
		} else {
			r.TotalStake = "0"
		}
		items = append(items, r)
	}
	if err := rows.Err(); err != nil {
		h.logger.Error("failed to iterate worknet agent rows", "error", err)
		return nil, newSvcErr(errInternal, "failed to list worknet agents")
	}
	return items, nil
}

// ── HIGH #8: subnets.search + subnets.getByOwner ──

// svcSearchWorknets fuzzy-searches worknets by name/symbol
func (h *Handler) svcSearchWorknets(ctx context.Context, chainID int64, query string, limit, offset int32) (any, error) {
	if len(query) < 1 || len(query) > 100 {
		return nil, newSvcErr(errBadInput, "query must be 1-100 characters")
	}
	// Sanitize query string to prevent SQL injection (pgx parameterized queries already prevent injection; this adds extra special character restrictions)
	pattern := "%" + strings.ReplaceAll(strings.ReplaceAll(query, "%", ""), "_", "") + "%"

	sql := `
		SELECT subnet_id, chain_id, owner, name, symbol, subnet_contract, skills_uri,
		       metadata_uri, min_stake, worknet_token, lp_pool, status, created_at,
		       activated_at, immunity_ends_at, burned
		FROM subnets
		WHERE chain_id = $1 AND burned = FALSE AND (name ILIKE $2 OR symbol ILIKE $2)
		ORDER BY subnet_id DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := h.db.Query(ctx, sql, chainID, pattern, limit, offset)
	if err != nil {
		h.logger.Error("failed to search worknets", "error", err)
		return nil, newSvcErr(errInternal, "failed to search worknets")
	}
	defer rows.Close()

	items := []gen.Subnet{}
	for rows.Next() {
		var s gen.Subnet
		if err := rows.Scan(
			&s.SubnetID, &s.ChainID, &s.Owner, &s.Name, &s.Symbol, &s.SubnetContract,
			&s.SkillsUri, &s.MetadataUri, &s.MinStake, &s.WorknetToken, &s.LpPool, &s.Status,
			&s.CreatedAt, &s.ActivatedAt, &s.ImmunityEndsAt, &s.Burned,
		); err != nil {
			h.logger.Error("failed to scan search result", "error", err)
			return nil, newSvcErr(errInternal, "failed to search worknets")
		}
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		h.logger.Error("failed to iterate search results", "error", err)
		return nil, newSvcErr(errInternal, "failed to search worknets")
	}
	return items, nil
}

// svcGetWorknetsByOwner queries worknets by owner
func (h *Handler) svcGetWorknetsByOwner(ctx context.Context, chainID int64, owner string, limit, offset int32) (any, error) {
	sql := `
		SELECT subnet_id, chain_id, owner, name, symbol, subnet_contract, skills_uri,
		       metadata_uri, min_stake, worknet_token, lp_pool, status, created_at,
		       activated_at, immunity_ends_at, burned
		FROM subnets
		WHERE chain_id = $1 AND owner = $2 AND burned = FALSE
		ORDER BY subnet_id DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := h.db.Query(ctx, sql, chainID, owner, limit, offset)
	if err != nil {
		h.logger.Error("failed to get worknets by owner", "error", err)
		return nil, newSvcErr(errInternal, "failed to get worknets by owner")
	}
	defer rows.Close()

	items := []gen.Subnet{}
	for rows.Next() {
		var s gen.Subnet
		if err := rows.Scan(
			&s.SubnetID, &s.ChainID, &s.Owner, &s.Name, &s.Symbol, &s.SubnetContract,
			&s.SkillsUri, &s.MetadataUri, &s.MinStake, &s.WorknetToken, &s.LpPool, &s.Status,
			&s.CreatedAt, &s.ActivatedAt, &s.ImmunityEndsAt, &s.Burned,
		); err != nil {
			h.logger.Error("failed to scan owner worknet", "error", err)
			return nil, newSvcErr(errInternal, "failed to get worknets by owner")
		}
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		h.logger.Error("failed to iterate owner worknets", "error", err)
		return nil, newSvcErr(errInternal, "failed to get worknets by owner")
	}
	return items, nil
}

// ── HIGH #9: staking.getPositionsGlobal ──

// svcGetPositionsGlobal fetches user veAWP positions across all chains
func (h *Handler) svcGetPositionsGlobal(ctx context.Context, address string) (map[string]any, error) {
	sql := `
		SELECT chain_id, token_id, owner, amount, lock_end_time, created_at, burned
		FROM stake_positions
		WHERE owner = $1 AND burned = FALSE
		ORDER BY chain_id, token_id
		LIMIT 500
	`
	rows, err := h.db.Query(ctx, sql, address)
	if err != nil {
		h.logger.Error("failed to get positions global", "error", err, "address", address)
		return nil, newSvcErr(errInternal, "failed to get positions")
	}
	defer rows.Close()

	items := []gen.StakePosition{}
	totalStaked := new(big.Int)
	for rows.Next() {
		var p gen.StakePosition
		if err := rows.Scan(&p.ChainID, &p.TokenID, &p.Owner, &p.Amount, &p.LockEndTime, &p.CreatedAt, &p.Burned); err != nil {
			h.logger.Error("failed to scan position", "error", err)
			return nil, newSvcErr(errInternal, "failed to get positions")
		}
		if p.Amount.Valid {
			totalStaked.Add(totalStaked, p.Amount.Int)
		}
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		h.logger.Error("failed to iterate positions", "error", err)
		return nil, newSvcErr(errInternal, "failed to get positions")
	}

	return map[string]any{
		"positions":   items,
		"totalStaked": totalStaked.String(),
		"count":       len(items),
	}, nil
}

// ── CRITICAL #3: users.getDelegates (DB-based, from users table bound_to) ──

// svcGetDelegates returns addresses that have bound_to = address (i.e., users who bound themselves
// to this address via bind()). This is NOT the same as on-chain delegates (grantDelegate/revokeDelegate).
// The indexer publishes DelegateGranted/DelegateRevoked events to Redis Pub/Sub for real-time
// WebSocket notifications, but delegate state is NOT persisted to the DB.
// TODO: persist DelegateGranted/DelegateRevoked events to a delegates table for accurate on-chain delegate queries.
func (h *Handler) svcGetDelegates(ctx context.Context, chainID int64, address string) (any, error) {
	agents, err := h.queries.GetUsersByBoundTo(ctx, gen.GetUsersByBoundToParams{
		BoundTo: address, ChainID: chainID,
	})
	if err != nil {
		h.logger.Error("failed to get delegates", "error", err, "address", address)
		return nil, newSvcErr(errInternal, "failed to get delegates")
	}
	return agents, nil
}
