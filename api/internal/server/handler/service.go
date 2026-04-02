package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

// ── 服务层错误类型 ──

// svcErrKind 区分业务错误类别（映射到 HTTP 状态码或 RPC 错误码）
type svcErrKind int

const (
	errInternal    svcErrKind = iota // 500 / rpcInternalError
	errBadInput                      // 400 / rpcInvalidParams
	errNotFound                      // 404 / rpcNotFound
	errUnavailable                   // 503 / rpcInternalError (service unavailable)
)

// svcError 是服务层统一错误
type svcError struct {
	Kind    svcErrKind
	Message string
}

func (e *svcError) Error() string { return e.Message }

func newSvcErr(kind svcErrKind, msg string) *svcError {
	return &svcError{Kind: kind, Message: msg}
}

// ── 分页常量与共享函数 ──

// validSubnetStatuses 子网状态有效值
var validSubnetStatuses = map[string]bool{"Pending": true, "Active": true, "Paused": true, "Banned": true}

// validProposalStatuses 提案状态有效值
var validProposalStatuses = map[string]bool{"Active": true, "Canceled": true, "Defeated": true, "Succeeded": true, "Queued": true, "Expired": true, "Executed": true}

// computePageLimits 统一分页参数计算：page >= 1, limit 1..100, 默认 limit=20
func computePageLimits(page, limit int) (int32, int32) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := 0
	if page > 1 && page <= 10000 {
		offset = (page - 1) * limit
	}
	return int32(limit), int32(offset)
}

// ── Chain ID 解析 ──

// resolveChainID 从 HTTP 请求中提取 chainId 查询参数，回退到默认链
func (h *Handler) resolveChainID(r *http.Request) int64 {
	if v := r.URL.Query().Get("chainId"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil && id > 0 {
			return id
		}
	}
	return h.defaultChainID()
}

// resolveRPCChainID 从 RPC 参数中提取 chainId，回退到默认链
func (h *Handler) resolveRPCChainID(chainID int64) int64 {
	if chainID > 0 {
		return chainID
	}
	return h.defaultChainID()
}

// defaultChainID 返回默认链 ID：优先从已加载的 chains 列表取第一个，否则用 cfg.ChainID
func (h *Handler) defaultChainID() int64 {
	if h.chains != nil && len(h.chains) > 0 {
		return h.chains[0].ChainID
	}
	return h.cfg.ChainID
}

// ── 服务方法 ──

// svcGetRegistry 返回合约地址注册信息（按 chainID 查询 DB，回退到 cfg）
func (h *Handler) svcGetRegistry(chainID int64) registryResponse {
	// 尝试从 DB chains 表读取 per-chain 合约地址
	if chainID > 0 {
		if chain, err := h.queries.GetChain(context.Background(), chainID); err == nil {
			return h.registryFromChainRow(chain)
		}
	}
	// 回退到 cfg 配置
	return registryResponse{
		ChainID:           chainID,
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
			ChainID: chainID, VerifyingContract: h.cfg.AWPRegistryAddress,
		},
		StakingVaultEIP712: eip712DomainResponse{
			Name: "StakingVault", Version: "1",
			ChainID: chainID, VerifyingContract: h.cfg.StakingVaultAddress,
		},
	}
}

// registryFromChainRow 从 DB Chain 记录构建 registryResponse
func (h *Handler) registryFromChainRow(c gen.Chain) registryResponse {
	resolve := func(dbVal, cfgVal string) string {
		v := strings.TrimSpace(dbVal)
		if v != "" {
			return v
		}
		return cfgVal
	}
	registry := resolve(c.AwpRegistry, h.cfg.AWPRegistryAddress)
	vault := resolve(c.StakingVault, h.cfg.StakingVaultAddress)
	return registryResponse{
		ChainID:           c.ChainID,
		AWPRegistry:       registry,
		AWPToken:          resolve(c.AwpToken, h.cfg.AWPTokenAddress),
		AWPEmission:       resolve(c.AwpEmission, h.cfg.AWPEmissionAddress),
		StakingVault:      vault,
		StakeNFT:          resolve(c.StakeNft, h.cfg.StakeNFTAddress),
		SubnetNFT:         resolve(c.SubnetNft, h.cfg.SubnetNFTAddress),
		LPManager:         resolve(c.LpManager, h.cfg.LPManagerAddress),
		AlphaTokenFactory: h.cfg.AlphaFactoryAddress,
		DAO:               resolve(c.DaoAddress, h.cfg.DAOAddress),
		Treasury:          h.cfg.TreasuryAddress,
		EIP712Domain: eip712DomainResponse{
			Name: "AWPRegistry", Version: "1",
			ChainID: c.ChainID, VerifyingContract: registry,
		},
		StakingVaultEIP712: eip712DomainResponse{
			Name: "StakingVault", Version: "1",
			ChainID: c.ChainID, VerifyingContract: vault,
		},
	}
}

// svcGetChains 返回支持的链列表（优先从 DB 读取，回退到 yaml 配置）
func (h *Handler) svcGetChains(ctx context.Context) any {
	// 尝试从 DB 读取
	if dbChains, err := h.queries.ListChains(ctx); err == nil && len(dbChains) > 0 {
		return dbChains
	}
	if h.chains != nil {
		return h.chains
	}
	return []map[string]interface{}{{"chainId": h.cfg.ChainID, "name": "Default"}}
}

// svcListUsers 分页获取用户列表
func (h *Handler) svcListUsers(ctx context.Context, chainID int64, limit, offset int32) (any, error) {
	users, err := h.queries.ListUsers(ctx, gen.ListUsersParams{
		ChainID: chainID, Limit: limit, Offset: offset,
	})
	if err != nil {
		h.logger.Error("failed to list users", "error", err)
		return nil, newSvcErr(errInternal, "failed to list users")
	}
	return users, nil
}

// svcGetUserCount 获取用户总数
func (h *Handler) svcGetUserCount(ctx context.Context, chainID int64) (int64, error) {
	count, err := h.queries.GetUserCount(ctx, chainID)
	if err != nil {
		h.logger.Error("failed to get user count", "error", err)
		return 0, newSvcErr(errInternal, "failed to get user count")
	}
	return count, nil
}

// svcGetUser 获取用户详情（含余额和绑定的 agent）
func (h *Handler) svcGetUser(ctx context.Context, chainID int64, address string) (*userDetailResponse, error) {
	user, err := h.queries.GetUser(ctx, gen.GetUserParams{Address: address, ChainID: chainID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newSvcErr(errNotFound, "user not found")
		}
		h.logger.Error("failed to get user", "error", err, "address", address)
		return nil, newSvcErr(errInternal, "failed to get user")
	}

	resp := &userDetailResponse{User: user, Agents: []gen.GetUsersByBoundToRow{}}
	if balance, err := h.queries.GetUserBalance(ctx, gen.GetUserBalanceParams{
		UserAddress: address, ChainID: chainID,
	}); err == nil {
		resp.Balance = &balance
	}
	if agents, err := h.queries.GetUsersByBoundTo(ctx, gen.GetUsersByBoundToParams{
		BoundTo: address, ChainID: chainID,
	}); err == nil {
		resp.Agents = agents
	}
	return resp, nil
}

// svcCheckAddress 检查地址注册状态、绑定和收款地址
func (h *Handler) svcCheckAddress(ctx context.Context, chainID int64, address string) (checkAddressResponse, error) {
	resp := checkAddressResponse{}
	if user, err := h.queries.GetUser(ctx, gen.GetUserParams{
		Address: address, ChainID: chainID,
	}); err == nil {
		resp.IsRegistered = user.RegisteredAt != 0 || user.BoundTo != "" || user.Recipient != ""
		resp.BoundTo = user.BoundTo
		resp.Recipient = user.Recipient
	} else if !errors.Is(err, pgx.ErrNoRows) {
		h.logger.Error("failed to check address", "error", err, "address", address)
		return resp, newSvcErr(errInternal, "failed to check address")
	}
	return resp, nil
}

// svcGetBalance 获取用户 AWP 质押余额（总质押/已分配/可用）
func (h *Handler) svcGetBalance(ctx context.Context, chainID int64, address string) (balanceResponse, error) {
	totalStakedNum, err := h.queries.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{
		ChainID: chainID, Owner: address,
	})
	if err != nil {
		h.logger.Error("failed to get user total staked", "error", err, "address", address)
		return balanceResponse{}, newSvcErr(errInternal, "failed to get user balance")
	}

	totalStaked := "0"
	if totalStakedNum.Valid {
		totalStaked = totalStakedNum.Int.String()
	}

	totalAllocated := "0"
	balance, err := h.queries.GetUserBalance(ctx, gen.GetUserBalanceParams{
		UserAddress: address, ChainID: chainID,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			h.logger.Error("failed to get user balance", "error", err, "address", address)
			return balanceResponse{}, newSvcErr(errInternal, "failed to get user balance")
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
		TotalStaked: totalStaked, TotalAllocated: totalAllocated, Unallocated: unallocated,
	}, nil
}

// svcGetAllocations 分页获取用户的质押分配列表
func (h *Handler) svcGetAllocations(ctx context.Context, chainID int64, address string, limit, offset int32) (any, error) {
	allocations, err := h.queries.GetAllocationsByUser(ctx, gen.GetAllocationsByUserParams{
		ChainID: chainID, UserAddress: address, Limit: limit, Offset: offset,
	})
	if err != nil {
		h.logger.Error("failed to get stake allocations", "error", err, "address", address)
		return nil, newSvcErr(errInternal, "failed to get stake allocations")
	}
	return allocations, nil
}

// svcGetFrozen 获取用户冻结的质押分配
func (h *Handler) svcGetFrozen(ctx context.Context, chainID int64, address string) (any, error) {
	frozen, err := h.queries.GetFrozenByUser(ctx, gen.GetFrozenByUserParams{
		ChainID: chainID, UserAddress: address,
	})
	if err != nil {
		h.logger.Error("failed to get frozen allocations", "error", err, "address", address)
		return nil, newSvcErr(errInternal, "failed to get frozen allocations")
	}
	return frozen, nil
}

// svcGetStakePositions 获取用户的 StakeNFT 持仓列表
func (h *Handler) svcGetStakePositions(ctx context.Context, chainID int64, address string) (any, error) {
	positions, err := h.queries.GetUserStakePositions(ctx, gen.GetUserStakePositionsParams{
		ChainID: chainID, Owner: address,
	})
	if err != nil {
		h.logger.Error("failed to get stake positions", "error", err, "address", address)
		return nil, newSvcErr(errInternal, "failed to get stake positions")
	}
	return positions, nil
}

// svcGetAgentSubnetStake 获取 agent 在某子网的质押量（跨链聚合，不需要 chainID）
func (h *Handler) svcGetAgentSubnetStake(ctx context.Context, agent string, subnetID pgtype.Numeric) (string, error) {
	stake, err := h.queries.GetAgentSubnetStakeGlobal(ctx, gen.GetAgentSubnetStakeGlobalParams{
		AgentAddress: agent, SubnetID: subnetID,
	})
	if err != nil {
		h.logger.Error("failed to get agent subnet stake", "error", err)
		return "", newSvcErr(errInternal, "failed to get agent subnet stake")
	}
	if stake.Valid {
		return stake.Int.String(), nil
	}
	return "0", nil
}

// svcGetAgentSubnets 获取 agent 参与的所有子网及质押量（跨链聚合）
func (h *Handler) svcGetAgentSubnets(ctx context.Context, agent string) (any, error) {
	subnets, err := h.queries.GetAgentSubnetsGlobal(ctx, agent)
	if err != nil {
		h.logger.Error("failed to get agent subnets", "error", err, "agent", agent)
		return nil, newSvcErr(errInternal, "failed to get agent subnets")
	}
	return subnets, nil
}

// svcGetSubnetTotalStake 获取子网总质押量（跨链聚合）
func (h *Handler) svcGetSubnetTotalStake(ctx context.Context, subnetID pgtype.Numeric) (string, error) {
	total, err := h.queries.GetSubnetTotalStake(ctx, subnetID)
	if err != nil {
		h.logger.Error("failed to get subnet total stake", "error", err, "subnetId", subnetID)
		return "", newSvcErr(errInternal, "failed to get subnet total stake")
	}
	if total.Valid {
		return total.Int.String(), nil
	}
	return "0", nil
}

// svcListSubnets 分页获取子网列表（可按状态筛选）
// svcListSubnets 获取子网列表。chainID=0 时返回所有链的子网（跨链），否则返回指定链的子网
func (h *Handler) svcListSubnets(ctx context.Context, chainID int64, status string, limit, offset int32) (any, error) {
	if status != "" {
		if !validSubnetStatuses[status] {
			return nil, newSvcErr(errBadInput, "invalid status filter: must be one of Pending, Active, Paused, Banned")
		}
		if chainID == 0 {
			subnets, err := h.queries.ListAllSubnetsByStatus(ctx, gen.ListAllSubnetsByStatusParams{
				Status: status, Limit: limit, Offset: offset,
			})
			if err != nil {
				return nil, newSvcErr(errInternal, "failed to list subnets")
			}
			return subnets, nil
		}
		subnets, err := h.queries.ListSubnetsByStatus(ctx, gen.ListSubnetsByStatusParams{
			ChainID: chainID, Status: status, Limit: limit, Offset: offset,
		})
		if err != nil {
			return nil, newSvcErr(errInternal, "failed to list subnets")
		}
		return subnets, nil
	}

	if chainID == 0 {
		subnets, err := h.queries.ListAllSubnets(ctx, gen.ListAllSubnetsParams{
			Limit: limit, Offset: offset,
		})
		if err != nil {
			return nil, newSvcErr(errInternal, "failed to list subnets")
		}
		return subnets, nil
	}

	subnets, err := h.queries.ListSubnets(ctx, gen.ListSubnetsParams{
		ChainID: chainID, Limit: limit, Offset: offset,
	})
	if err != nil {
		return nil, newSvcErr(errInternal, "failed to list subnets")
	}
	return subnets, nil
}

// svcGetSubnet 获取子网详情（按 subnetID 查询，不需要 chainID）
func (h *Handler) svcGetSubnet(ctx context.Context, subnetID pgtype.Numeric) (any, error) {
	subnet, err := h.queries.GetSubnet(ctx, subnetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newSvcErr(errNotFound, "subnet not found")
		}
		h.logger.Error("failed to get subnet", "error", err, "subnetId", subnetID)
		return nil, newSvcErr(errInternal, "failed to get subnet")
	}
	return subnet, nil
}

// svcGetSubnetSkills 获取子网 skills URI
func (h *Handler) svcGetSubnetSkills(ctx context.Context, subnetID pgtype.Numeric) (map[string]interface{}, error) {
	skillsURI, err := h.queries.GetSubnetSkills(ctx, subnetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newSvcErr(errNotFound, "subnet not found")
		}
		h.logger.Error("failed to get subnet skills", "error", err, "subnetId", subnetID)
		return nil, newSvcErr(errInternal, "failed to get subnet skills")
	}
	var uri string
	if skillsURI.Valid {
		uri = skillsURI.String
	}
	return map[string]interface{}{"subnetId": subnetID, "skillsURI": uri}, nil
}

// svcGetSubnetEarnings 获取子网 AWP 收益历史（分页）
func (h *Handler) svcGetSubnetEarnings(ctx context.Context, subnetID pgtype.Numeric, limit, offset int32) (any, error) {
	earnings, err := h.queries.GetSubnetEarningsByID(ctx, gen.GetSubnetEarningsByIDParams{
		SubnetID: subnetID, Limit: limit, Offset: offset,
	})
	if err != nil {
		h.logger.Error("failed to get subnet earnings", "error", err, "subnetId", subnetID)
		return nil, newSvcErr(errInternal, "failed to get subnet earnings")
	}
	return earnings, nil
}

// svcGetSubnetAgentInfo 获取 agent 在子网的质押信息
func (h *Handler) svcGetSubnetAgentInfo(ctx context.Context, agent string, subnetID pgtype.Numeric) (map[string]any, error) {
	amount, err := h.svcGetAgentSubnetStake(ctx, agent, subnetID)
	if err != nil {
		return nil, err
	}
	return map[string]any{"agent": agent, "subnetId": subnetID, "stake": amount}, nil
}

// svcGetEmissionSchedule 获取排放预测（30/90/365 天）
func (h *Handler) svcGetEmissionSchedule(ctx context.Context, chainID int64) (map[string]any, error) {
	currentDaily := new(big.Int).Set(initialDailyEmission)
	latestEpoch, err := h.queries.GetLatestEpoch(ctx, chainID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			h.logger.Error("failed to get latest epoch", "error", err)
			return nil, newSvcErr(errInternal, "failed to get emission data")
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

// svcListEpochs 分页获取 epoch 列表
func (h *Handler) svcListEpochs(ctx context.Context, chainID int64, limit, offset int32) (any, error) {
	epochs, err := h.queries.ListEpochs(ctx, gen.ListEpochsParams{
		ChainID: chainID, Limit: limit, Offset: offset,
	})
	if err != nil {
		h.logger.Error("failed to list epochs", "error", err)
		return nil, newSvcErr(errInternal, "failed to list epochs")
	}
	return epochs, nil
}

// svcReadRedisJSON 读取 Redis 缓存中的 JSON 值
func (h *Handler) svcReadRedisJSON(ctx context.Context, key string) (any, error) {
	val, err := h.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 调用者决定如何处理 nil（空对象 vs 错误）
		}
		h.logger.Error("failed to read Redis key", "error", err, "key", key)
		return nil, newSvcErr(errInternal, "failed to read cache")
	}
	var data any
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		h.logger.Error("failed to parse Redis JSON", "error", err, "key", key)
		return nil, newSvcErr(errInternal, "cache data format error")
	}
	return data, nil
}

// svcGetAgentsByOwner 获取绑定到指定 owner 的所有 agent
func (h *Handler) svcGetAgentsByOwner(ctx context.Context, chainID int64, owner string) (any, error) {
	agents, err := h.queries.GetUsersByBoundTo(ctx, gen.GetUsersByBoundToParams{
		BoundTo: owner, ChainID: chainID,
	})
	if err != nil {
		h.logger.Error("failed to get agents by owner", "error", err, "owner", owner)
		return nil, newSvcErr(errInternal, "failed to get agents")
	}
	return agents, nil
}

// svcGetAgentDetail 获取 agent 详情（user 记录）
func (h *Handler) svcGetAgentDetail(ctx context.Context, chainID int64, agent string) (any, error) {
	user, err := h.queries.GetUser(ctx, gen.GetUserParams{Address: agent, ChainID: chainID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newSvcErr(errNotFound, "agent not found")
		}
		h.logger.Error("failed to get agent detail", "error", err, "agent", agent)
		return nil, newSvcErr(errInternal, "failed to get agent detail")
	}
	return user, nil
}

// svcLookupAgent 查找 agent 的 owner（boundTo）
func (h *Handler) svcLookupAgent(ctx context.Context, chainID int64, agent string) (string, error) {
	user, err := h.queries.GetUser(ctx, gen.GetUserParams{Address: agent, ChainID: chainID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", newSvcErr(errNotFound, "agent not found")
		}
		h.logger.Error("failed to lookup agent owner", "error", err, "agent", agent)
		return "", newSvcErr(errInternal, "failed to lookup agent owner")
	}
	if user.BoundTo == "" {
		return "", newSvcErr(errNotFound, "agent not bound")
	}
	return user.BoundTo, nil
}

// svcBatchAgentInfo 批量查询 agent 信息及质押
func (h *Handler) svcBatchAgentInfo(ctx context.Context, chainID int64, agents []string, subnetID pgtype.Numeric) ([]agentInfoItem, error) {
	// Validate and normalize
	validAddrs := make([]string, 0, len(agents))
	for _, addr := range agents {
		if !isValidAddress(addr) {
			continue
		}
		validAddrs = append(validAddrs, normalizeAddr(addr))
	}
	if len(validAddrs) == 0 {
		return []agentInfoItem{}, nil
	}

	users, err := h.queries.GetUsersBatch(ctx, gen.GetUsersBatchParams{
		ChainID: chainID, Addresses: validAddrs,
	})
	if err != nil {
		h.logger.Error("failed to batch get users", "error", err)
		return nil, newSvcErr(errInternal, "failed to get agent info")
	}
	userMap := make(map[string]gen.GetUsersBatchRow, len(users))
	for _, u := range users {
		userMap[u.Address] = u
	}

	stakes, err := h.queries.GetAgentSubnetStakesBatch(ctx, gen.GetAgentSubnetStakesBatchParams{
		ChainID: chainID, Agents: validAddrs, SubnetID: subnetID,
	})
	if err != nil {
		h.logger.Error("failed to batch get stakes", "error", err)
		return nil, newSvcErr(errInternal, "failed to get agent info")
	}
	stakeMap := make(map[string]string, len(stakes))
	for _, s := range stakes {
		if s.Total.Valid {
			stakeMap[s.AgentAddress] = s.Total.Int.String()
		}
	}

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

// ── 跨链聚合服务方法 ──

// svcGetGlobalStats 获取全局协议统计（跨链聚合）
func (h *Handler) svcGetGlobalStats(ctx context.Context) (map[string]any, error) {
	totalSubnets, err := h.queries.CountAllSubnets(ctx)
	if err != nil {
		h.logger.Error("failed to count all subnets", "error", err)
		return nil, newSvcErr(errInternal, "failed to get global stats")
	}

	totalUsers, err := h.queries.CountAllDistinctUsers(ctx)
	if err != nil {
		h.logger.Error("failed to count all users", "error", err)
		return nil, newSvcErr(errInternal, "failed to get global stats")
	}

	totalStaked, err := h.queries.SumAllStaked(ctx)
	if err != nil {
		h.logger.Error("failed to sum all staked", "error", err)
		return nil, newSvcErr(errInternal, "failed to get global stats")
	}

	totalAllocated, err := h.queries.SumAllAllocated(ctx)
	if err != nil {
		h.logger.Error("failed to sum all allocated", "error", err)
		return nil, newSvcErr(errInternal, "failed to get global stats")
	}

	stakedStr := "0"
	if totalStaked.Valid {
		stakedStr = totalStaked.Int.String()
	}
	allocatedStr := "0"
	if totalAllocated.Valid {
		allocatedStr = totalAllocated.Int.String()
	}

	chainCount := 1
	if h.chains != nil && len(h.chains) > 0 {
		chainCount = len(h.chains)
	}

	return map[string]any{
		"totalSubnets":   totalSubnets,
		"totalUsers":     totalUsers,
		"totalStaked":    stakedStr,
		"totalAllocated": allocatedStr,
		"chains":         chainCount,
	}, nil
}

// svcGetUserBalanceGlobal 获取用户跨链聚合质押余额
func (h *Handler) svcGetUserBalanceGlobal(ctx context.Context, address string) (balanceResponse, error) {
	totalStakedNum, err := h.queries.GetUserTotalStakedGlobal(ctx, address)
	if err != nil {
		h.logger.Error("failed to get user total staked global", "error", err, "address", address)
		return balanceResponse{}, newSvcErr(errInternal, "failed to get user balance")
	}

	totalStaked := "0"
	if totalStakedNum.Valid {
		totalStaked = totalStakedNum.Int.String()
	}

	totalAllocatedNum, err := h.queries.GetUserBalanceGlobal(ctx, address)
	if err != nil {
		h.logger.Error("failed to get user balance global", "error", err, "address", address)
		return balanceResponse{}, newSvcErr(errInternal, "failed to get user balance")
	}

	totalAllocated := "0"
	if totalAllocatedNum.Valid {
		totalAllocated = totalAllocatedNum.Int.String()
	}

	unallocated := "0"
	if totalStakedNum.Valid {
		stakedBig := totalStakedNum.Int
		allocBig := new(big.Int)
		if totalAllocatedNum.Valid {
			allocBig.Set(totalAllocatedNum.Int)
		}
		diff := new(big.Int).Sub(stakedBig, allocBig)
		if diff.Sign() < 0 {
			diff.SetInt64(0)
		}
		unallocated = diff.String()
	}

	return balanceResponse{
		TotalStaked: totalStaked, TotalAllocated: totalAllocated, Unallocated: unallocated,
	}, nil
}

// svcGetGlobalEmissionSchedule 获取全链排放汇总
func (h *Handler) svcGetGlobalEmissionSchedule(ctx context.Context) (map[string]any, error) {
	chainIDs := []int64{h.defaultChainID()}
	if h.chains != nil && len(h.chains) > 0 {
		chainIDs = make([]int64, len(h.chains))
		for i, c := range h.chains {
			chainIDs[i] = c.ChainID
		}
	}

	globalDaily := new(big.Int)
	perChain := make([]map[string]any, 0, len(chainIDs))

	for _, cid := range chainIDs {
		key := fmt.Sprintf("emission_current:%d", cid)
		data, err := h.svcReadRedisJSON(ctx, key)
		if err != nil || data == nil {
			perChain = append(perChain, map[string]any{"chainId": cid, "dailyEmission": "0", "available": false})
			continue
		}
		// 从 Redis JSON 中提取 dailyEmission
		dataMap, ok := data.(map[string]any)
		if !ok {
			perChain = append(perChain, map[string]any{"chainId": cid, "dailyEmission": "0", "available": false})
			continue
		}
		dailyStr := "0"
		if v, ok := dataMap["dailyEmission"]; ok {
			dailyStr = fmt.Sprintf("%v", v)
		}
		daily, ok := new(big.Int).SetString(dailyStr, 10)
		if !ok {
			daily = new(big.Int)
		}
		globalDaily.Add(globalDaily, daily)
		perChain = append(perChain, map[string]any{"chainId": cid, "dailyEmission": dailyStr, "available": true})
	}

	// 计算全局预测
	periods := []int{30, 90, 365}
	projections := make([]emissionProjection, 0, len(periods))
	for _, days := range periods {
		total := new(big.Int)
		daily := new(big.Int).Set(globalDaily)
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
		"globalDailyEmission": globalDaily.String(),
		"chains":              perChain,
		"projections":         projections,
	}, nil
}

// svcListUsersGlobal 跨链去重用户列表
func (h *Handler) svcListUsersGlobal(ctx context.Context, limit, offset int32) (map[string]any, error) {
	users, err := h.queries.ListAllUsers(ctx, gen.ListAllUsersParams{
		Limit: limit, Offset: offset,
	})
	if err != nil {
		h.logger.Error("failed to list all users", "error", err)
		return nil, newSvcErr(errInternal, "failed to list users")
	}
	total, err := h.queries.CountAllDistinctUsers(ctx)
	if err != nil {
		h.logger.Error("failed to count all distinct users", "error", err)
		return nil, newSvcErr(errInternal, "failed to count users")
	}
	return map[string]any{"users": users, "total": total}, nil
}

// svcListAllProposals 跨链提案列表
func (h *Handler) svcListAllProposals(ctx context.Context, status string, limit, offset int32) (any, error) {
	if status != "" {
		if !validProposalStatuses[status] {
			return nil, newSvcErr(errBadInput, "invalid status filter: must be one of Active, Canceled, Defeated, Succeeded, Queued, Expired, Executed")
		}
		proposals, err := h.queries.ListAllProposalsByStatus(ctx, gen.ListAllProposalsByStatusParams{
			Status: status, Limit: limit, Offset: offset,
		})
		if err != nil {
			h.logger.Error("failed to list all proposals by status", "error", err, "status", status)
			return nil, newSvcErr(errInternal, "failed to list proposals")
		}
		return proposals, nil
	}

	proposals, err := h.queries.ListAllProposals(ctx, gen.ListAllProposalsParams{
		Limit: limit, Offset: offset,
	})
	if err != nil {
		h.logger.Error("failed to list all proposals", "error", err)
		return nil, newSvcErr(errInternal, "failed to list proposals")
	}
	return proposals, nil
}

// svcGetAlphaInfo 获取子网 Alpha 代币信息
func (h *Handler) svcGetAlphaInfo(ctx context.Context, subnetID pgtype.Numeric) (map[string]any, error) {
	subnet, err := h.queries.GetSubnet(ctx, subnetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newSvcErr(errNotFound, "subnet not found")
		}
		h.logger.Error("failed to get subnet info", "error", err, "subnetId", subnetID)
		return nil, newSvcErr(errInternal, "failed to get subnet info")
	}
	return map[string]any{
		"subnetId": subnet.SubnetID, "name": subnet.Name,
		"symbol": subnet.Symbol, "alphaToken": subnet.AlphaToken,
	}, nil
}

// svcGetAlphaPrice 获取子网 Alpha 代币价格
func (h *Handler) svcGetAlphaPrice(ctx context.Context, subnetIDRaw string) (any, error) {
	data, err := h.svcReadRedisJSON(ctx, fmt.Sprintf("alpha_price:%s", subnetIDRaw))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return map[string]any{}, nil
	}
	return data, nil
}

// svcGetAWPInfo 获取 AWP 代币信息
func (h *Handler) svcGetAWPInfo(ctx context.Context, chainID int64) (any, error) {
	data, err := h.svcReadRedisJSON(ctx, fmt.Sprintf("awp_info:%d", chainID))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return map[string]any{}, nil
	}
	return data, nil
}

// svcGetCurrentEmission 获取当前排放数据
func (h *Handler) svcGetCurrentEmission(ctx context.Context, chainID int64) (any, error) {
	data, err := h.svcReadRedisJSON(ctx, fmt.Sprintf("emission_current:%d", chainID))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, newSvcErr(errUnavailable, "emission data not yet available (keeper may not be running)")
	}
	return data, nil
}

// svcListProposals 分页获取治理提案（可按状态筛选）
func (h *Handler) svcListProposals(ctx context.Context, chainID int64, status string, limit, offset int32) (any, error) {
	if status != "" {
		if !validProposalStatuses[status] {
			return nil, newSvcErr(errBadInput, "invalid status filter: must be one of Active, Canceled, Defeated, Succeeded, Queued, Expired, Executed")
		}
		proposals, err := h.queries.ListProposalsByStatus(ctx, gen.ListProposalsByStatusParams{
			ChainID: chainID, Status: status, Limit: limit, Offset: offset,
		})
		if err != nil {
			h.logger.Error("failed to list proposals by status", "error", err, "status", status)
			return nil, newSvcErr(errInternal, "failed to list proposals")
		}
		return proposals, nil
	}

	proposals, err := h.queries.ListProposals(ctx, gen.ListProposalsParams{
		ChainID: chainID, Limit: limit, Offset: offset,
	})
	if err != nil {
		h.logger.Error("failed to list proposals", "error", err)
		return nil, newSvcErr(errInternal, "failed to list proposals")
	}
	return proposals, nil
}

// svcGetProposal 获取单个治理提案详情
func (h *Handler) svcGetProposal(ctx context.Context, chainID int64, proposalID string) (any, error) {
	proposal, err := h.queries.GetProposal(ctx, gen.GetProposalParams{
		ChainID: chainID, ProposalID: proposalID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newSvcErr(errNotFound, "proposal not found")
		}
		h.logger.Error("failed to get proposal", "error", err, "proposalId", proposalID)
		return nil, newSvcErr(errInternal, "failed to get proposal")
	}
	return proposal, nil
}
