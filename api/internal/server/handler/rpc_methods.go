package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

// ── 通用参数结构 ──

type addressParams struct {
	Address string `json:"address"`
}

type subnetParams struct {
	SubnetID string `json:"subnetId"`
}

type pageParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// parsePage 解析分页参数，返回 (limit, offset)
func parsePage(p pageParams) (int32, int32) {
	limit := 20
	if p.Limit > 0 {
		limit = p.Limit
	}
	if limit > 100 {
		limit = 100
	}
	offset := 0
	if p.Page > 1 && p.Page <= 10000 {
		offset = (p.Page - 1) * limit
	}
	return int32(limit), int32(offset)
}

// parseSubnetNum 将字符串 subnetId 解析为 pgtype.Numeric
func parseSubnetNum(s string) (pgtype.Numeric, *RPCErr) {
	if s == "" {
		return pgtype.Numeric{}, &RPCErr{Code: rpcInvalidParams, Message: "subnetId is required"}
	}
	id, ok := new(big.Int).SetString(s, 10)
	if !ok || id.Sign() <= 0 {
		return pgtype.Numeric{}, &RPCErr{Code: rpcInvalidParams, Message: "subnetId must be a positive integer"}
	}
	return pgtype.Numeric{Int: id, Exp: 0, Valid: true}, nil
}

// requireAddress 验证并规范化地址
func requireAddress(addr string) (string, *RPCErr) {
	if !isValidAddress(addr) {
		return "", &RPCErr{Code: rpcInvalidParams, Message: "invalid address: must be 0x + 40 hex chars"}
	}
	return normalizeAddr(addr), nil
}

// internalErr 返回内部错误
func internalErr(msg string) *RPCErr {
	return &RPCErr{Code: rpcInternalError, Message: msg}
}

// ═══════════════════════════════════════════════
// ── registry ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcRegistryGet(_ context.Context, _ json.RawMessage) (any, *RPCErr) {
	return registryResponse{
		ChainID:           h.cfg.ChainID,
		AWPRegistry:       h.cfg.AWPRegistryAddress,
		AWPToken:          h.cfg.AWPTokenAddress,
		AWPEmission:       h.cfg.AWPEmissionAddress,
		StakingVault:      h.cfg.StakingVaultAddress,
		StakeNFT:          h.cfg.StakeNFTAddress,
		SubnetNFT:         h.cfg.SubnetNFTAddress,
		LPManager:         h.cfg.LPManagerAddress,
		AlphaTokenFactory: h.cfg.AlphaFactoryAddress,
		DAO:               h.cfg.DAOAddress,
		Treasury:          h.cfg.TreasuryAddress,
		EIP712Domain: eip712DomainResponse{
			Name: "AWPRegistry", Version: "1",
			ChainID: h.cfg.ChainID, VerifyingContract: h.cfg.AWPRegistryAddress,
		},
		StakingVaultEIP712: eip712DomainResponse{
			Name: "StakingVault", Version: "1",
			ChainID: h.cfg.ChainID, VerifyingContract: h.cfg.StakingVaultAddress,
		},
	}, nil
}

// ═══════════════════════════════════════════════
// ── health ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcHealthCheck(_ context.Context, _ json.RawMessage) (any, *RPCErr) {
	return map[string]string{"status": "ok"}, nil
}

func (h *Handler) rpcHealthDetailed(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	return h.buildDetailedHealth(ctx), nil
}

// ═══════════════════════════════════════════════
// ── chains ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcChainsList(_ context.Context, _ json.RawMessage) (any, *RPCErr) {
	if h.chains == nil {
		return []map[string]interface{}{{"chainId": h.cfg.ChainID, "name": "Default"}}, nil
	}
	return h.chains, nil
}

// ═══════════════════════════════════════════════
// ── users ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcUsersList(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p pageParams
	_ = json.Unmarshal(raw, &p)
	limit, offset := parsePage(p)

	users, err := h.queries.ListUsers(ctx, gen.ListUsersParams{
		ChainID: h.cfg.ChainID, Limit: limit, Offset: offset,
	})
	if err != nil {
		return nil, internalErr("failed to list users")
	}
	return users, nil
}

func (h *Handler) rpcUsersCount(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	count, err := h.queries.GetUserCount(ctx, h.cfg.ChainID)
	if err != nil {
		return nil, internalErr("failed to get user count")
	}
	return map[string]int64{"count": count}, nil
}

func (h *Handler) rpcUsersGet(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p addressParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}

	user, err := h.queries.GetUser(ctx, gen.GetUserParams{Address: address, ChainID: h.cfg.ChainID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &RPCErr{Code: rpcNotFound, Message: "user not found"}
		}
		return nil, internalErr("failed to get user")
	}

	resp := userDetailResponse{User: user, Agents: []gen.GetUsersByBoundToRow{}}
	if balance, err := h.queries.GetUserBalance(ctx, gen.GetUserBalanceParams{
		UserAddress: address, ChainID: h.cfg.ChainID,
	}); err == nil {
		resp.Balance = &balance
	}
	if agents, err := h.queries.GetUsersByBoundTo(ctx, gen.GetUsersByBoundToParams{
		BoundTo: address, ChainID: h.cfg.ChainID,
	}); err == nil {
		resp.Agents = agents
	}

	return resp, nil
}

// ═══════════════════════════════════════════════
// ── address ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcAddressCheck(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p addressParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}

	resp := checkAddressResponse{}
	if user, err := h.queries.GetUser(ctx, gen.GetUserParams{
		Address: address, ChainID: h.cfg.ChainID,
	}); err == nil {
		resp.IsRegistered = user.RegisteredAt != 0 || user.BoundTo != "" || user.Recipient != ""
		resp.BoundTo = user.BoundTo
		resp.Recipient = user.Recipient
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, internalErr("failed to check address")
	}

	return resp, nil
}

// ═══════════════════════════════════════════════
// ── nonce ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcNonceGet(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	// Rate limit（与 REST 端点一致）
	if ip := rpcClientIP(ctx); ip != "" {
		if exceeded, _ := h.limiter.CheckAndIncrement(ctx, "nonce", ip); exceeded {
			return nil, &RPCErr{Code: rpcInvalidRequest, Message: "rate limit exceeded"}
		}
	}
	var p addressParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}
	if h.chain == nil {
		return nil, &RPCErr{Code: rpcInternalError, Message: "chain reader not available"}
	}
	nonce, err := h.chain.GetNonce(address)
	if err != nil {
		return nil, internalErr("failed to read nonce")
	}
	return map[string]uint64{"nonce": nonce}, nil
}

func (h *Handler) rpcNonceGetStaking(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	if ip := rpcClientIP(ctx); ip != "" {
		if exceeded, _ := h.limiter.CheckAndIncrement(ctx, "nonce", ip); exceeded {
			return nil, &RPCErr{Code: rpcInvalidRequest, Message: "rate limit exceeded"}
		}
	}
	var p addressParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}
	if h.chain == nil {
		return nil, &RPCErr{Code: rpcInternalError, Message: "chain reader not available"}
	}
	nonce, err := h.chain.GetStakingNonce(address)
	if err != nil {
		return nil, internalErr("failed to read staking nonce")
	}
	return map[string]uint64{"nonce": nonce}, nil
}

// ═══════════════════════════════════════════════
// ── agents ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcAgentsGetByOwner(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Owner string `json:"owner"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	owner, rpcErr := requireAddress(p.Owner)
	if rpcErr != nil {
		return nil, rpcErr
	}

	agents, err := h.queries.GetUsersByBoundTo(ctx, gen.GetUsersByBoundToParams{
		BoundTo: owner, ChainID: h.cfg.ChainID,
	})
	if err != nil {
		return nil, internalErr("failed to get agents")
	}
	return agents, nil
}

func (h *Handler) rpcAgentsGetDetail(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Agent string `json:"agent"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	agent, rpcErr := requireAddress(p.Agent)
	if rpcErr != nil {
		return nil, rpcErr
	}

	user, err := h.queries.GetUser(ctx, gen.GetUserParams{Address: agent, ChainID: h.cfg.ChainID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &RPCErr{Code: rpcNotFound, Message: "agent not found"}
		}
		return nil, internalErr("failed to get agent detail")
	}
	return user, nil
}

func (h *Handler) rpcAgentsLookup(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Agent string `json:"agent"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	agent, rpcErr := requireAddress(p.Agent)
	if rpcErr != nil {
		return nil, rpcErr
	}

	user, err := h.queries.GetUser(ctx, gen.GetUserParams{Address: agent, ChainID: h.cfg.ChainID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &RPCErr{Code: rpcNotFound, Message: "agent not found"}
		}
		return nil, internalErr("failed to lookup agent")
	}
	if user.BoundTo == "" {
		return nil, &RPCErr{Code: rpcNotFound, Message: "agent not bound"}
	}
	return map[string]string{"ownerAddress": user.BoundTo}, nil
}

func (h *Handler) rpcAgentsBatchInfo(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	if ip := rpcClientIP(ctx); ip != "" {
		if exceeded, _ := h.limiter.CheckAndIncrement(ctx, "batch_agent_info", ip); exceeded {
			return nil, &RPCErr{Code: rpcInvalidRequest, Message: "rate limit exceeded"}
		}
	}
	var p struct {
		Agents   []string `json:"agents"`
		SubnetID string   `json:"subnetId"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	if len(p.Agents) == 0 {
		return []agentInfoItem{}, nil
	}
	if len(p.Agents) > 100 {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "batch size exceeds limit (100)"}
	}

	subnetNum, rpcErr := parseSubnetNum(p.SubnetID)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Validate and normalize all addresses upfront
	validAddrs := make([]string, 0, len(p.Agents))
	for _, addr := range p.Agents {
		if !isValidAddress(addr) {
			continue
		}
		validAddrs = append(validAddrs, normalizeAddr(addr))
	}
	if len(validAddrs) == 0 {
		return []agentInfoItem{}, nil
	}

	// Batch query: fetch all users in one DB call
	users, err := h.queries.GetUsersBatch(ctx, gen.GetUsersBatchParams{
		ChainID: h.cfg.ChainID, Addresses: validAddrs,
	})
	if err != nil {
		return nil, internalErr("failed to get agent info")
	}
	userMap := make(map[string]gen.GetUsersBatchRow, len(users))
	for _, u := range users {
		userMap[u.Address] = u
	}

	// Batch query: fetch all stakes in one DB call
	stakes, err := h.queries.GetAgentSubnetStakesBatch(ctx, gen.GetAgentSubnetStakesBatchParams{
		ChainID: h.cfg.ChainID, Agents: validAddrs, SubnetID: subnetNum,
	})
	if err != nil {
		return nil, internalErr("failed to get agent info")
	}
	stakeMap := make(map[string]string, len(stakes))
	for _, s := range stakes {
		if s.Total.Valid {
			stakeMap[s.AgentAddress] = s.Total.Int.String()
		}
	}

	// Assemble results preserving input order
	results := make([]agentInfoItem, 0, len(validAddrs))
	for _, addr := range validAddrs {
		user, ok := userMap[addr]
		if !ok {
			continue
		}
		item := agentInfoItem{Address: user.Address, BoundTo: user.BoundTo, Stake: "0"}
		if s, ok := stakeMap[addr]; ok {
			item.Stake = s
		}
		results = append(results, item)
	}
	return results, nil
}

// ═══════════════════════════════════════════════
// ── staking ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcStakingGetBalance(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p addressParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}

	totalStakedNum, err := h.queries.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{
		ChainID: h.cfg.ChainID, Owner: address,
	})
	if err != nil {
		return nil, internalErr("failed to get user balance")
	}

	totalStaked := "0"
	if totalStakedNum.Valid {
		totalStaked = totalStakedNum.Int.String()
	}

	totalAllocated := "0"
	balance, err := h.queries.GetUserBalance(ctx, gen.GetUserBalanceParams{
		UserAddress: address, ChainID: h.cfg.ChainID,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, internalErr("failed to get user balance")
		}
	} else if balance.TotalAllocated.Valid {
		totalAllocated = balance.TotalAllocated.Int.String()
	}

	unallocated := "0"
	if totalStakedNum.Valid {
		stakedBig := totalStakedNum.Int
		allocBig := new(big.Int)
		if balance.TotalAllocated.Valid {
			allocBig.Set(balance.TotalAllocated.Int)
		}
		unallocated = new(big.Int).Sub(stakedBig, allocBig).String()
	}

	return balanceResponse{
		TotalStaked:    totalStaked,
		TotalAllocated: totalAllocated,
		Unallocated:    unallocated,
	}, nil
}

func (h *Handler) rpcStakingGetPositions(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p addressParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}

	positions, err := h.queries.GetUserStakePositions(ctx, gen.GetUserStakePositionsParams{
		ChainID: h.cfg.ChainID, Owner: address,
	})
	if err != nil {
		return nil, internalErr("failed to get stake positions")
	}
	return positions, nil
}

func (h *Handler) rpcStakingGetAllocations(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Address string `json:"address"`
		pageParams
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}
	limit, offset := parsePage(p.pageParams)

	allocations, err := h.queries.GetAllocationsByUser(ctx, gen.GetAllocationsByUserParams{
		ChainID: h.cfg.ChainID, UserAddress: address, Limit: limit, Offset: offset,
	})
	if err != nil {
		return nil, internalErr("failed to get allocations")
	}
	return allocations, nil
}

func (h *Handler) rpcStakingGetFrozen(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p addressParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}

	frozen, err := h.queries.GetFrozenByUser(ctx, gen.GetFrozenByUserParams{
		ChainID: h.cfg.ChainID, UserAddress: address,
	})
	if err != nil {
		return nil, internalErr("failed to get frozen allocations")
	}
	return frozen, nil
}

func (h *Handler) rpcStakingGetPending(_ context.Context, _ json.RawMessage) (any, *RPCErr) {
	// 双槽模式下 reallocate 即时生效，无待处理记录
	return []struct{}{}, nil
}

func (h *Handler) rpcStakingGetAgentSubnetStake(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Agent    string `json:"agent"`
		SubnetID string `json:"subnetId"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	agent, rpcErr := requireAddress(p.Agent)
	if rpcErr != nil {
		return nil, rpcErr
	}
	subnetNum, rpcErr := parseSubnetNum(p.SubnetID)
	if rpcErr != nil {
		return nil, rpcErr
	}

	stake, err := h.queries.GetAgentSubnetStakeGlobal(ctx, gen.GetAgentSubnetStakeGlobalParams{
		AgentAddress: agent, SubnetID: subnetNum,
	})
	if err != nil {
		return nil, internalErr("failed to get agent subnet stake")
	}
	amount := "0"
	if stake.Valid {
		amount = stake.Int.String()
	}
	return map[string]string{"amount": amount}, nil
}

func (h *Handler) rpcStakingGetAgentSubnets(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Agent string `json:"agent"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	agent, rpcErr := requireAddress(p.Agent)
	if rpcErr != nil {
		return nil, rpcErr
	}

	subnets, err := h.queries.GetAgentSubnetsGlobal(ctx, agent)
	if err != nil {
		return nil, internalErr("failed to get agent subnets")
	}
	return subnets, nil
}

func (h *Handler) rpcStakingGetSubnetTotalStake(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p subnetParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	subnetNum, rpcErr := parseSubnetNum(p.SubnetID)
	if rpcErr != nil {
		return nil, rpcErr
	}

	total, err := h.queries.GetSubnetTotalStake(ctx, subnetNum)
	if err != nil {
		return nil, internalErr("failed to get subnet total stake")
	}
	amount := "0"
	if total.Valid {
		amount = total.Int.String()
	}
	return map[string]string{"total": amount}, nil
}

// ═══════════════════════════════════════════════
// ── subnets ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcSubnetsList(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Status string `json:"status"`
		pageParams
	}
	_ = json.Unmarshal(raw, &p)
	limit, offset := parsePage(p.pageParams)

	if p.Status != "" {
		validStatuses := map[string]bool{"Pending": true, "Active": true, "Paused": true, "Banned": true}
		if !validStatuses[p.Status] {
			return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid status: must be Pending/Active/Paused/Banned"}
		}
		subnets, err := h.queries.ListSubnetsByStatus(ctx, gen.ListSubnetsByStatusParams{
			ChainID: h.cfg.ChainID, Status: p.Status, Limit: limit, Offset: offset,
		})
		if err != nil {
			return nil, internalErr("failed to list subnets")
		}
		return subnets, nil
	}

	subnets, err := h.queries.ListSubnets(ctx, gen.ListSubnetsParams{
		ChainID: h.cfg.ChainID, Limit: limit, Offset: offset,
	})
	if err != nil {
		return nil, internalErr("failed to list subnets")
	}
	return subnets, nil
}

func (h *Handler) rpcSubnetsGet(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p subnetParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	subnetNum, rpcErr := parseSubnetNum(p.SubnetID)
	if rpcErr != nil {
		return nil, rpcErr
	}

	subnet, err := h.queries.GetSubnet(ctx, subnetNum)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &RPCErr{Code: rpcNotFound, Message: "subnet not found"}
		}
		return nil, internalErr("failed to get subnet")
	}
	return subnet, nil
}

func (h *Handler) rpcSubnetsGetSkills(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p subnetParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	subnetNum, rpcErr := parseSubnetNum(p.SubnetID)
	if rpcErr != nil {
		return nil, rpcErr
	}

	skillsURI, err := h.queries.GetSubnetSkills(ctx, subnetNum)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &RPCErr{Code: rpcNotFound, Message: "subnet not found"}
		}
		return nil, internalErr("failed to get subnet skills")
	}
	var uri string
	if skillsURI.Valid {
		uri = skillsURI.String
	}
	return map[string]interface{}{"subnetId": subnetNum, "skillsURI": uri}, nil
}

func (h *Handler) rpcSubnetsGetEarnings(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		SubnetID string `json:"subnetId"`
		pageParams
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	subnetNum, rpcErr := parseSubnetNum(p.SubnetID)
	if rpcErr != nil {
		return nil, rpcErr
	}
	limit, offset := parsePage(p.pageParams)

	earnings, err := h.queries.GetSubnetEarningsByID(ctx, gen.GetSubnetEarningsByIDParams{
		SubnetID: subnetNum, Limit: limit, Offset: offset,
	})
	if err != nil {
		return nil, internalErr("failed to get subnet earnings")
	}
	return earnings, nil
}

func (h *Handler) rpcSubnetsGetAgentInfo(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		SubnetID string `json:"subnetId"`
		Agent    string `json:"agent"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	subnetNum, rpcErr := parseSubnetNum(p.SubnetID)
	if rpcErr != nil {
		return nil, rpcErr
	}
	agent, rpcErr := requireAddress(p.Agent)
	if rpcErr != nil {
		return nil, rpcErr
	}

	total, err := h.queries.GetAgentSubnetStakeGlobal(ctx, gen.GetAgentSubnetStakeGlobalParams{
		AgentAddress: agent, SubnetID: subnetNum,
	})
	stakeStr := "0"
	if err == nil && total.Valid {
		stakeStr = total.Int.String()
	}
	return map[string]any{"agent": agent, "subnetId": subnetNum, "stake": stakeStr}, nil
}

// ═══════════════════════════════════════════════
// ── emission ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcEmissionGetCurrent(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	val, err := h.rdb.Get(ctx, fmt.Sprintf("emission_current:%d", h.cfg.ChainID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, &RPCErr{Code: rpcInternalError, Message: "emission data not yet available"}
		}
		return nil, internalErr("failed to get emission data")
	}
	var data any
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, internalErr("emission data format error")
	}
	return data, nil
}

func (h *Handler) rpcEmissionGetSchedule(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	currentDaily := new(big.Int).Set(initialDailyEmission)
	latestEpoch, err := h.queries.GetLatestEpoch(ctx, h.cfg.ChainID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, internalErr("failed to get emission data")
		}
	} else if latestEpoch.DailyEmission.Valid {
		currentDaily = latestEpoch.DailyEmission.Int
	}

	periods := []int{30, 90, 365}
	projections := make([]emissionProjection, 0, len(periods))
	for _, days := range periods {
		total := new(big.Int)
		daily := new(big.Int).Set(currentDaily)
		for d := 0; d < days; d++ {
			total.Add(total, daily)
			daily.Mul(daily, decayFactor)
			daily.Div(daily, decayPrecision)
		}
		projections = append(projections, emissionProjection{
			Days: days, TotalEmission: total.String(), FinalDailyRate: daily.String(),
		})
	}

	return map[string]any{
		"currentDailyEmission": currentDaily.String(),
		"projections":          projections,
	}, nil
}

func (h *Handler) rpcEmissionListEpochs(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p pageParams
	_ = json.Unmarshal(raw, &p)
	limit, offset := parsePage(p)

	epochs, err := h.queries.ListEpochs(ctx, gen.ListEpochsParams{
		ChainID: h.cfg.ChainID, Limit: limit, Offset: offset,
	})
	if err != nil {
		return nil, internalErr("failed to list epochs")
	}
	return epochs, nil
}

// ═══════════════════════════════════════════════
// ── tokens ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcTokensGetAWP(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	val, err := h.rdb.Get(ctx, fmt.Sprintf("awp_info:%d", h.cfg.ChainID)).Result()
	if err != nil {
		if err == redis.Nil {
			return map[string]any{}, nil
		}
		return nil, internalErr("failed to get AWP info")
	}
	var data any
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, internalErr("AWP data format error")
	}
	return data, nil
}

func (h *Handler) rpcTokensGetAlphaInfo(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p subnetParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	subnetNum, rpcErr := parseSubnetNum(p.SubnetID)
	if rpcErr != nil {
		return nil, rpcErr
	}

	subnet, err := h.queries.GetSubnet(ctx, subnetNum)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &RPCErr{Code: rpcNotFound, Message: "subnet not found"}
		}
		return nil, internalErr("failed to get subnet info")
	}
	return map[string]any{
		"subnetId": subnet.SubnetID, "name": subnet.Name,
		"symbol": subnet.Symbol, "alphaToken": subnet.AlphaToken,
	}, nil
}

func (h *Handler) rpcTokensGetAlphaPrice(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p subnetParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	// Validate subnetId is a positive integer
	if _, rpcErr := parseSubnetNum(p.SubnetID); rpcErr != nil {
		return nil, rpcErr
	}

	key := fmt.Sprintf("alpha_price:%s", p.SubnetID)
	val, err := h.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return map[string]any{}, nil
		}
		return nil, internalErr("failed to get alpha price")
	}
	var data any
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, internalErr("price data format error")
	}
	return data, nil
}

// ═══════════════════════════════════════════════
// ── governance ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcGovernanceListProposals(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Status string `json:"status"`
		pageParams
	}
	_ = json.Unmarshal(raw, &p)
	limit, offset := parsePage(p.pageParams)

	if p.Status != "" {
		validStatuses := map[string]bool{"Active": true, "Canceled": true, "Defeated": true, "Succeeded": true, "Queued": true, "Expired": true, "Executed": true}
		if !validStatuses[p.Status] {
			return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid status: must be Active/Canceled/Defeated/Succeeded/Queued/Expired/Executed"}
		}
		proposals, err := h.queries.ListProposalsByStatus(ctx, gen.ListProposalsByStatusParams{
			ChainID: h.cfg.ChainID, Status: p.Status, Limit: limit, Offset: offset,
		})
		if err != nil {
			return nil, internalErr("failed to list proposals")
		}
		return proposals, nil
	}

	proposals, err := h.queries.ListProposals(ctx, gen.ListProposalsParams{
		ChainID: h.cfg.ChainID, Limit: limit, Offset: offset,
	})
	if err != nil {
		return nil, internalErr("failed to list proposals")
	}
	return proposals, nil
}

func (h *Handler) rpcGovernanceGetProposal(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		ProposalID string `json:"proposalId"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	if p.ProposalID == "" {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "proposalId is required"}
	}

	proposal, err := h.queries.GetProposal(ctx, gen.GetProposalParams{
		ChainID: h.cfg.ChainID, ProposalID: p.ProposalID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &RPCErr{Code: rpcNotFound, Message: "proposal not found"}
		}
		return nil, internalErr("failed to get proposal")
	}
	return proposal, nil
}

func (h *Handler) rpcGovernanceTreasury(_ context.Context, _ json.RawMessage) (any, *RPCErr) {
	return map[string]string{"treasuryAddress": h.cfg.TreasuryAddress}, nil
}
