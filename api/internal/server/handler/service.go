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

// ── Service layer error types ──

// svcErrKind distinguishes business error categories (mapped to HTTP status codes or RPC error codes)
type svcErrKind int

const (
	errInternal    svcErrKind = iota // 500 / rpcInternalError
	errBadInput                      // 400 / rpcInvalidParams
	errNotFound                      // 404 / rpcNotFound
	errUnavailable                   // 503 / rpcInternalError (service unavailable)
)

// svcError is the unified service layer error
type svcError struct {
	Kind    svcErrKind
	Message string
}

func (e *svcError) Error() string { return e.Message }

func newSvcErr(kind svcErrKind, msg string) *svcError {
	return &svcError{Kind: kind, Message: msg}
}

// ── Pagination constants and shared functions ──

// validSubnetStatuses valid subnet status values
var validSubnetStatuses = map[string]bool{"Pending": true, "Active": true, "Paused": true, "Banned": true}

// validProposalStatuses valid proposal status values
var validProposalStatuses = map[string]bool{"Active": true, "Canceled": true, "Defeated": true, "Succeeded": true, "Queued": true, "Expired": true, "Executed": true}

// computePageLimits unified pagination parameter calculation: page >= 1, limit 1..100, default limit=20
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

// ── Chain ID resolution ──

// resolveChainID extracts the chainId query parameter from an HTTP request, falling back to the default chain
func (h *Handler) resolveChainID(r *http.Request) int64 {
	if v := r.URL.Query().Get("chainId"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil && id > 0 {
			return id
		}
	}
	return h.defaultChainID()
}

// resolveRPCChainID extracts chainId from RPC parameters, falling back to the default chain
func (h *Handler) resolveRPCChainID(chainID int64) int64 {
	if chainID > 0 {
		return chainID
	}
	return h.defaultChainID()
}

// defaultChainID returns the default chain ID: first from loaded chains list if available, otherwise cfg.ChainID
func (h *Handler) defaultChainID() int64 {
	if h.chains != nil && len(h.chains) > 0 {
		return h.chains[0].ChainID
	}
	return h.cfg.ChainID
}

// getActiveChainIDs returns all configured chain IDs (multi-chain mode from chains yaml, single-chain mode falls back to cfg.ChainID)
func (h *Handler) getActiveChainIDs() []int64 {
	if h.chains != nil && len(h.chains) > 0 {
		ids := make([]int64, len(h.chains))
		for i, c := range h.chains {
			ids[i] = c.ChainID
		}
		return ids
	}
	if h.cfg.ChainID > 0 {
		return []int64{h.cfg.ChainID}
	}
	return nil
}

// ── Service methods ──

// svcGetRegistry returns contract address registry info (queries DB by chainID, falls back to cfg)
func (h *Handler) svcGetRegistry(ctx context.Context, chainID int64) registryResponse {
	// Try to read per-chain contract addresses from DB chains table
	if chainID > 0 {
		if chain, err := h.queries.GetChain(ctx, chainID); err == nil {
			return h.registryFromChainRow(chain)
		}
	}
	// Fall back to cfg configuration
	allocatorAddr := h.cfg.AWPAllocatorAddress
	return registryResponse{
		ChainID:             chainID,
		AWPRegistry:         h.cfg.AWPRegistryAddress,
		AWPToken:            h.cfg.AWPTokenAddress,
		AWPEmission:         h.cfg.AWPEmissionAddress,
		AWPAllocator:        allocatorAddr,
		VeAWP:               h.cfg.VeAWPAddress,
		AWPWorkNet:          h.cfg.AWPWorkNetAddress,
		LPManager:           h.cfg.LPManagerAddress,
		WorknetTokenFactory: h.cfg.WorknetTokenFactoryAddress,
		DAO:                 h.cfg.DAOAddress,
		Treasury:            h.cfg.TreasuryAddress,
		VeAWPHelper:         "0x0000561EDE5C1Ba0b81cE585964050bEAE730001",
		EIP712Domain: eip712DomainResponse{
			Name: "AWPRegistry", Version: "1",
			ChainID: chainID, VerifyingContract: h.cfg.AWPRegistryAddress,
		},
		AllocatorEIP712: eip712DomainResponse{
			Name: "AWPAllocator", Version: "1",
			ChainID: chainID, VerifyingContract: allocatorAddr,
		},
	}
}

// svcGetRegistryAll returns registry info for all configured chains
func (h *Handler) svcGetRegistryAll(ctx context.Context) []registryResponse {
	chainIDs := h.getActiveChainIDs()
	results := make([]registryResponse, 0, len(chainIDs))
	for _, id := range chainIDs {
		results = append(results, h.svcGetRegistry(ctx, id))
	}
	return results
}

// registryFromChainRow builds a registryResponse from a DB Chain record.
// Note: WorknetTokenFactory and Treasury always come from global config (h.cfg), not from the
// per-chain DB row. This is by design — these contracts share the same CREATE2 address on all chains.
func (h *Handler) registryFromChainRow(c gen.Chain) registryResponse {
	resolve := func(dbVal, cfgVal string) string {
		v := strings.TrimSpace(dbVal)
		if v != "" {
			return v
		}
		return cfgVal
	}
	registry := resolve(c.AwpRegistry, h.cfg.AWPRegistryAddress)
	allocator := resolve(c.AwpAllocator, h.cfg.AWPAllocatorAddress)
	return registryResponse{
		ChainID:             c.ChainID,
		AWPRegistry:         registry,
		AWPToken:            resolve(c.AwpToken, h.cfg.AWPTokenAddress),
		AWPEmission:         resolve(c.AwpEmission, h.cfg.AWPEmissionAddress),
		AWPAllocator:        allocator,
		VeAWP:               resolve(c.Veawp, h.cfg.VeAWPAddress),
		AWPWorkNet:          resolve(c.AwpWorknet, h.cfg.AWPWorkNetAddress),
		LPManager:           resolve(c.LpManager, h.cfg.LPManagerAddress),
		WorknetTokenFactory: h.cfg.WorknetTokenFactoryAddress,
		DAO:                 resolve(c.DaoAddress, h.cfg.DAOAddress),
		Treasury:            h.cfg.TreasuryAddress,
		VeAWPHelper:         "0x0000561EDE5C1Ba0b81cE585964050bEAE730001",
		EIP712Domain: eip712DomainResponse{
			Name: "AWPRegistry", Version: "1",
			ChainID: c.ChainID, VerifyingContract: registry,
		},
		AllocatorEIP712: eip712DomainResponse{
			Name: "AWPAllocator", Version: "1",
			ChainID: c.ChainID, VerifyingContract: allocator,
		},
	}
}

// svcGetChains returns the list of supported chains (reads from DB first, falls back to yaml config)
func (h *Handler) svcGetChains(ctx context.Context) any {
	// Try to read from DB
	if dbChains, err := h.queries.ListChains(ctx); err == nil && len(dbChains) > 0 {
		return dbChains
	}
	if h.chains != nil {
		return h.chains
	}
	return []map[string]interface{}{{"chainId": h.cfg.ChainID, "name": "Default"}}
}

// svcListUsers fetches a paginated user list
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

// svcGetUserCount fetches the total user count
func (h *Handler) svcGetUserCount(ctx context.Context, chainID int64) (int64, error) {
	count, err := h.queries.GetUserCount(ctx, chainID)
	if err != nil {
		h.logger.Error("failed to get user count", "error", err)
		return 0, newSvcErr(errInternal, "failed to get user count")
	}
	return count, nil
}

// svcGetUser fetches user details (including balance and bound agents)
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

// svcCheckAddress checks address registration status, binding, and recipient address.
// When chainID is 0 (caller did not specify), returns registration info across all chains.
func (h *Handler) svcCheckAddress(ctx context.Context, chainID int64, address string) (checkAddressResponse, error) {
	resp := checkAddressResponse{}

	if chainID > 0 {
		// Specific chain requested
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
	} else {
		// No chain specified — return all chains where user is registered
		rows, err := h.queries.GetUserAllChains(ctx, address)
		if err != nil {
			h.logger.Error("failed to check address", "error", err, "address", address)
			return resp, newSvcErr(errInternal, "failed to check address")
		}
		if len(rows) > 0 {
			resp.IsRegistered = true
			for _, r := range rows {
				resp.Chains = append(resp.Chains, chainRegistration{
					ChainID:      r.ChainID,
					IsRegistered: true,
					BoundTo:      r.BoundTo,
					Recipient:    r.Recipient,
				})
			}
		}
	}
	return resp, nil
}

// svcGetBalance fetches the user AWP staking balance (total staked / allocated / available)
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
		totalStaked = numericString(totalStakedNum)
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
		totalAllocated = numericString(balance.TotalAllocated)
	}

	unallocated := "0"
	if totalStakedNum.Valid {
		stakedBig := totalStakedNum.Int
		allocBig := new(big.Int)
		if balance.TotalAllocated.Valid {
			allocBig.Set(balance.TotalAllocated.Int)
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

// svcGetAllocations fetches a paginated list of the user's stake allocations
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

// svcGetFrozen fetches the user's frozen stake allocations
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

// svcGetStakePositions fetches the user's veAWP position list
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

// numericString converts a pgtype.Numeric to its full decimal string representation.
// pgtype.Numeric stores (Int, Exp) where value = Int * 10^Exp.
// Int.String() alone returns just the mantissa, ignoring the exponent.
func numericString(n pgtype.Numeric) string {
	if !n.Valid || n.Int == nil {
		return "0"
	}
	if n.Exp == 0 {
		return n.Int.String()
	}
	if n.Exp > 0 {
		return new(big.Int).Mul(n.Int, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n.Exp)), nil)).String()
	}
	// Negative exponent: integer part only (truncate fractional digits)
	// e.g., Int=12345, Exp=-2 → 123 (drop .45)
	// e.g., Int=5, Exp=-3 → 0 (0.005 truncated to 0)
	// Token amounts in this codebase always use Exp=0, so this path is a safety fallback.
	s := n.Int.String()
	neg := ""
	if s[0] == '-' {
		neg = "-"
		s = s[1:]
	}
	absExp := int(-n.Exp)
	if absExp >= len(s) {
		return "0"
	}
	return neg + s[:len(s)-absExp]
}

// svcGetAgentSubnetStake fetches the agent's stake in a subnet (cross-chain aggregation, chainID not needed)
func (h *Handler) svcGetAgentSubnetStake(ctx context.Context, agent string, subnetID pgtype.Numeric) (string, error) {
	stake, err := h.queries.GetAgentSubnetStakeGlobal(ctx, gen.GetAgentSubnetStakeGlobalParams{
		AgentAddress: agent, SubnetID: subnetID,
	})
	if err != nil {
		h.logger.Error("failed to get agent subnet stake", "error", err)
		return "", newSvcErr(errInternal, "failed to get agent subnet stake")
	}
	if stake.Valid {
		return numericString(stake), nil
	}
	return "0", nil
}

// svcGetAgentSubnets fetches all subnets the agent participates in with stake amounts (cross-chain aggregation)
func (h *Handler) svcGetAgentSubnets(ctx context.Context, agent string) (any, error) {
	subnets, err := h.queries.GetAgentSubnetsGlobal(ctx, agent)
	if err != nil {
		h.logger.Error("failed to get agent subnets", "error", err, "agent", agent)
		return nil, newSvcErr(errInternal, "failed to get agent subnets")
	}
	return subnets, nil
}

// svcGetSubnetTotalStake fetches the subnet total stake (cross-chain aggregation)
func (h *Handler) svcGetSubnetTotalStake(ctx context.Context, subnetID pgtype.Numeric) (string, error) {
	total, err := h.queries.GetSubnetTotalStake(ctx, subnetID)
	if err != nil {
		h.logger.Error("failed to get subnet total stake", "error", err, "worknetId", subnetID)
		return "", newSvcErr(errInternal, "failed to get subnet total stake")
	}
	if total.Valid {
		return numericString(total), nil
	}
	return "0", nil
}

// svcListSubnets fetches a paginated subnet list (optionally filtered by status).
// When chainID=0, returns subnets from all chains (cross-chain); otherwise returns subnets for the specified chain
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

// svcGetSubnet fetches subnet details (queried by subnetID, chainID not needed)
func (h *Handler) svcGetSubnet(ctx context.Context, subnetID pgtype.Numeric) (any, error) {
	subnet, err := h.queries.GetSubnet(ctx, subnetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newSvcErr(errNotFound, "subnet not found")
		}
		h.logger.Error("failed to get subnet", "error", err, "worknetId", subnetID)
		return nil, newSvcErr(errInternal, "failed to get subnet")
	}
	return subnet, nil
}

// svcGetSubnetSkills fetches the subnet skills URI
func (h *Handler) svcGetSubnetSkills(ctx context.Context, subnetID pgtype.Numeric) (map[string]interface{}, error) {
	skillsURI, err := h.queries.GetSubnetSkills(ctx, subnetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newSvcErr(errNotFound, "subnet not found")
		}
		h.logger.Error("failed to get subnet skills", "error", err, "worknetId", subnetID)
		return nil, newSvcErr(errInternal, "failed to get subnet skills")
	}
	var uri string
	if skillsURI.Valid {
		uri = skillsURI.String
	}
	return map[string]interface{}{"worknetId": subnetID, "skillsURI": uri}, nil
}

// svcGetSubnetEarnings fetches the subnet AWP earnings history (paginated)
func (h *Handler) svcGetSubnetEarnings(ctx context.Context, subnetID pgtype.Numeric, limit, offset int32) (any, error) {
	earnings, err := h.queries.GetSubnetEarningsByID(ctx, gen.GetSubnetEarningsByIDParams{
		SubnetID: subnetID, Limit: limit, Offset: offset,
	})
	if err != nil {
		h.logger.Error("failed to get subnet earnings", "error", err, "worknetId", subnetID)
		return nil, newSvcErr(errInternal, "failed to get subnet earnings")
	}
	return earnings, nil
}

// svcGetSubnetAgentInfo fetches the agent's staking info in a subnet
func (h *Handler) svcGetSubnetAgentInfo(ctx context.Context, agent string, subnetID pgtype.Numeric) (map[string]any, error) {
	amount, err := h.svcGetAgentSubnetStake(ctx, agent, subnetID)
	if err != nil {
		return nil, err
	}
	return map[string]any{"agent": agent, "worknetId": subnetID, "stake": amount}, nil
}

// svcGetEmissionSchedule fetches emission projections (30/90/365 days)
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

// svcListEpochs fetches a paginated epoch list
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

// svcReadRedisJSON reads a JSON value from Redis cache
func (h *Handler) svcReadRedisJSON(ctx context.Context, key string) (any, error) {
	val, err := h.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // caller decides how to handle nil (empty object vs error)
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

// svcGetAgentsByOwner fetches all agents bound to the specified owner
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

// svcGetAgentDetail fetches agent details (user record)
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

// svcLookupAgent looks up the agent's owner (boundTo)
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

// svcBatchAgentInfo batch-queries agent info and stake
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
			stakeMap[s.AgentAddress] = numericString(s.Total)
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

// ── Cross-chain aggregation service methods ──

// svcGetGlobalStats fetches global protocol statistics (cross-chain aggregation)
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
		stakedStr = numericString(totalStaked)
	}
	allocatedStr := "0"
	if totalAllocated.Valid {
		allocatedStr = numericString(totalAllocated)
	}

	chainCount := 1
	if h.chains != nil && len(h.chains) > 0 {
		chainCount = len(h.chains)
	}

	// Per-chain sync freshness — indexer sync status for each chain
	chainIDs := h.getActiveChainIDs()
	syncStates := []map[string]any{}
	for _, cid := range chainIDs {
		state, err := h.queries.GetSyncState(ctx, gen.GetSyncStateParams{ChainID: cid, ContractName: "indexer"})
		entry := map[string]any{"chainId": cid}
		if err == nil {
			entry["lastBlock"] = state.LastBlock
		}
		syncStates = append(syncStates, entry)
	}

	return map[string]any{
		"totalSubnets":   totalSubnets,
		"totalUsers":     totalUsers,
		"totalStaked":    stakedStr,
		"totalAllocated": allocatedStr,
		"chains":         chainCount,
		"chainSyncStates": syncStates,
	}, nil
}

// svcGetUserBalanceGlobal fetches the user's cross-chain aggregated staking balance
func (h *Handler) svcGetUserBalanceGlobal(ctx context.Context, address string) (balanceResponse, error) {
	totalStakedNum, err := h.queries.GetUserTotalStakedGlobal(ctx, address)
	if err != nil {
		h.logger.Error("failed to get user total staked global", "error", err, "address", address)
		return balanceResponse{}, newSvcErr(errInternal, "failed to get user balance")
	}

	totalStaked := "0"
	if totalStakedNum.Valid {
		totalStaked = numericString(totalStakedNum)
	}

	totalAllocatedNum, err := h.queries.GetUserBalanceGlobal(ctx, address)
	if err != nil {
		h.logger.Error("failed to get user balance global", "error", err, "address", address)
		return balanceResponse{}, newSvcErr(errInternal, "failed to get user balance")
	}

	totalAllocated := "0"
	if totalAllocatedNum.Valid {
		totalAllocated = numericString(totalAllocatedNum)
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

// svcGetGlobalEmissionSchedule fetches global emission schedule aggregated across all chains
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
		// Extract dailyEmission from Redis JSON
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

	// Calculate global projections
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

// svcListUsersGlobal returns a cross-chain deduplicated user list
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

// svcListAllProposals returns a cross-chain proposal list
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

// svcGetWorknetTokenInfo fetches worknet token info
func (h *Handler) svcGetWorknetTokenInfo(ctx context.Context, subnetID pgtype.Numeric) (map[string]any, error) {
	subnet, err := h.queries.GetSubnet(ctx, subnetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newSvcErr(errNotFound, "subnet not found")
		}
		h.logger.Error("failed to get subnet info", "error", err, "worknetId", subnetID)
		return nil, newSvcErr(errInternal, "failed to get subnet info")
	}
	return map[string]any{
		"worknetId": subnet.SubnetID, "name": subnet.Name,
		"symbol": subnet.Symbol, "worknetToken": subnet.WorknetToken,
	}, nil
}

// svcGetWorknetTokenPrice fetches the worknet token price
func (h *Handler) svcGetWorknetTokenPrice(ctx context.Context, subnetIDRaw string) (any, error) {
	data, err := h.svcReadRedisJSON(ctx, fmt.Sprintf("worknet_token_price:%s", subnetIDRaw))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return map[string]any{}, nil
	}
	return data, nil
}

// svcGetAWPInfo fetches AWP token info
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

// svcGetAWPInfoGlobal aggregates AWP token info across chains (each chain is independent, no bridge)
func (h *Handler) svcGetAWPInfoGlobal(ctx context.Context) (map[string]any, error) {
	chainIDs := h.getActiveChainIDs()
	perChain := []map[string]any{}
	for _, cid := range chainIDs {
		data, err := h.svcReadRedisJSON(ctx, fmt.Sprintf("awp_info:%d", cid))
		if err != nil {
			perChain = append(perChain, map[string]any{"chainId": cid, "error": "unavailable"})
			continue
		}
		entry := map[string]any{"chainId": cid, "data": data}
		perChain = append(perChain, entry)
	}
	return map[string]any{
		"chains":      perChain,
		"independent": true,
		"note":        "AWP tokens on each chain are independent (no bridge). Total supply = sum of per-chain supplies.",
	}, nil
}

// svcGetCurrentEmission fetches current emission data
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

// svcListProposals fetches a paginated list of governance proposals (optionally filtered by status)
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

// svcGetProposal fetches a single governance proposal's details
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
