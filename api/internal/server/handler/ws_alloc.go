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
// 使用全局查询（无 chain_id 过滤）— subnetId 全局唯一，跨链聚合是正确的
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
		return stake.Int.String(), nil
	}
	return "0", nil
}
