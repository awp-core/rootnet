package handler

import (
	"context"
	"math/big"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

// WSAllocQuerier implements the ws.AllocationQuerier interface for WebSocket Hub allocation balance queries
type WSAllocQuerier struct {
	queries *gen.Queries
}

// NewWSAllocQuerier creates a WSAllocQuerier
func NewWSAllocQuerier(queries *gen.Queries) *WSAllocQuerier {
	return &WSAllocQuerier{queries: queries}
}

// GetAgentSubnetStakeWS queries the current allocation for (agent, worknetId)
// Uses global query (no chain_id filter) — worknetId is globally unique, cross-chain aggregation is correct
func (q *WSAllocQuerier) GetAgentSubnetStakeWS(ctx context.Context, _ int64, agent string, subnetID string) (string, error) {
	id, ok := new(big.Int).SetString(subnetID, 10)
	if !ok || id.Sign() <= 0 {
		return "0", nil
	}
	subnetNum := pgtype.Numeric{Int: id, Exp: 0, Valid: true}

	stake, err := q.queries.GetAgentSubnetStakeGlobal(ctx, gen.GetAgentSubnetStakeGlobalParams{
		AgentAddress: agent,
		SubnetID:     subnetNum,
	})
	if err != nil {
		return "0", err
	}
	if stake.Valid {
		return numericString(stake), nil
	}
	return "0", nil
}
