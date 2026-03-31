package handler

import (
	"context"
	"math/big"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

// WSAllocQuerier 实现 ws.AllocationQuerier 接口，供 WebSocket Hub 查询分配余额
type WSAllocQuerier struct {
	queries *gen.Queries
}

// NewWSAllocQuerier 创建 WSAllocQuerier
func NewWSAllocQuerier(queries *gen.Queries) *WSAllocQuerier {
	return &WSAllocQuerier{queries: queries}
}

// GetAgentSubnetStakeWS 查询 (agent, subnetId) 当前分配量
func (q *WSAllocQuerier) GetAgentSubnetStakeWS(ctx context.Context, chainID int64, agent string, subnetID string) (string, error) {
	id, ok := new(big.Int).SetString(subnetID, 10)
	if !ok {
		return "0", nil
	}
	subnetNum := pgtype.Numeric{Int: id, Exp: 0, Valid: true}

	stake, err := q.queries.GetAgentSubnetStake(ctx, gen.GetAgentSubnetStakeParams{
		ChainID:      chainID,
		AgentAddress: agent,
		SubnetID:     subnetNum,
	})
	if err != nil {
		return "0", err
	}
	if stake.Valid {
		return stake.Int.String(), nil
	}
	return "0", nil
}
