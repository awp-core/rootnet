package handler

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// ── 通用参数结构 ──

type addressParams struct {
	Address string `json:"address"`
	ChainID int64  `json:"chainId"`
}

type subnetParams struct {
	SubnetID string `json:"subnetId"`
	ChainID  int64  `json:"chainId"`
}

type pageParams struct {
	Page    int   `json:"page"`
	Limit   int   `json:"limit"`
	ChainID int64 `json:"chainId"`
}

// parsePage 解析分页参数，返回 (limit, offset)；委托给 computePageLimits
func parsePage(p pageParams) (int32, int32) {
	return computePageLimits(p.Page, p.Limit)
}

// parseSubnetNum 将字符串 subnetId 解析为 pgtype.Numeric，返回 RPCErr
func parseSubnetNum(s string) (pgtype.Numeric, *RPCErr) {
	num, err := parseSubnetIDString(s)
	if err != nil {
		return pgtype.Numeric{}, &RPCErr{Code: rpcInvalidParams, Message: err.Error()}
	}
	return num, nil
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

// svcToRPC 将 svcError 转为 RPCErr
func svcToRPC(err error) *RPCErr {
	if se, ok := err.(*svcError); ok {
		switch se.Kind {
		case errNotFound:
			return &RPCErr{Code: rpcNotFound, Message: se.Message}
		case errBadInput:
			return &RPCErr{Code: rpcInvalidParams, Message: se.Message}
		default:
			return &RPCErr{Code: rpcInternalError, Message: se.Message}
		}
	}
	return internalErr("internal error")
}

// ═══════════════════════════════════════════════
// ── registry ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcRegistryGet(_ context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		ChainID int64 `json:"chainId"`
	}
	_ = json.Unmarshal(raw, &p)
	chainID := h.resolveRPCChainID(p.ChainID)
	return h.svcGetRegistry(chainID), nil
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

func (h *Handler) rpcChainsList(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	return h.svcGetChains(ctx), nil
}

// ═══════════════════════════════════════════════
// ── users ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcUsersList(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p pageParams
	_ = json.Unmarshal(raw, &p)
	chainID := h.resolveRPCChainID(p.ChainID)
	limit, offset := parsePage(p)
	result, err := h.svcListUsers(ctx, chainID, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcUsersCount(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		ChainID int64 `json:"chainId"`
	}
	_ = json.Unmarshal(raw, &p)
	chainID := h.resolveRPCChainID(p.ChainID)
	count, err := h.svcGetUserCount(ctx, chainID)
	if err != nil {
		return nil, svcToRPC(err)
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
	chainID := h.resolveRPCChainID(p.ChainID)
	resp, err := h.svcGetUser(ctx, chainID, address)
	if err != nil {
		return nil, svcToRPC(err)
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
	chainID := h.resolveRPCChainID(p.ChainID)
	resp, err := h.svcCheckAddress(ctx, chainID, address)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return resp, nil
}

func (h *Handler) rpcResolveRecipient(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p addressParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	cr := h.getChainReader(chainID)
	if cr == nil {
		return nil, &RPCErr{Code: rpcInternalError, Message: "chain reader not available"}
	}
	resolved, err := cr.ResolveRecipient(address)
	if err != nil {
		return nil, internalErr("failed to resolve recipient")
	}
	return map[string]string{"address": address, "resolvedRecipient": strings.ToLower(resolved)}, nil
}

func (h *Handler) rpcBatchResolveRecipients(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Addresses []string `json:"addresses"`
		ChainID   int64    `json:"chainId"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	if len(p.Addresses) == 0 {
		return []any{}, nil
	}
	if len(p.Addresses) > 500 {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "batch size exceeds limit (500)"}
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	cr := h.getChainReader(chainID)
	if cr == nil {
		return nil, &RPCErr{Code: rpcInternalError, Message: "chain reader not available"}
	}
	normalized := make([]string, len(p.Addresses))
	for i, a := range p.Addresses {
		addr, rpcErr := requireAddress(a)
		if rpcErr != nil {
			return nil, rpcErr
		}
		normalized[i] = addr
	}
	resolved, err := cr.BatchResolveRecipients(normalized)
	if err != nil {
		return nil, internalErr("failed to resolve recipients")
	}
	results := make([]map[string]string, len(normalized))
	for i := range normalized {
		results[i] = map[string]string{"address": normalized[i], "resolvedRecipient": resolved[i]}
	}
	return results, nil
}

// ═══════════════════════════════════════════════
// ── nonce ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcNonceGet(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
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
	chainID := h.resolveRPCChainID(p.ChainID)
	cr := h.getChainReader(chainID)
	if cr == nil {
		return nil, &RPCErr{Code: rpcInternalError, Message: "chain reader not available for chainId"}
	}
	nonce, err := cr.GetNonce(address)
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
	chainID := h.resolveRPCChainID(p.ChainID)
	cr := h.getChainReader(chainID)
	if cr == nil {
		return nil, &RPCErr{Code: rpcInternalError, Message: "chain reader not available for chainId"}
	}
	nonce, err := cr.GetStakingNonce(address)
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
		Owner   string `json:"owner"`
		ChainID int64  `json:"chainId"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	owner, rpcErr := requireAddress(p.Owner)
	if rpcErr != nil {
		return nil, rpcErr
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	result, err := h.svcGetAgentsByOwner(ctx, chainID, owner)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcAgentsGetDetail(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Agent   string `json:"agent"`
		ChainID int64  `json:"chainId"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	agent, rpcErr := requireAddress(p.Agent)
	if rpcErr != nil {
		return nil, rpcErr
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	result, err := h.svcGetAgentDetail(ctx, chainID, agent)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcAgentsLookup(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Agent   string `json:"agent"`
		ChainID int64  `json:"chainId"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	agent, rpcErr := requireAddress(p.Agent)
	if rpcErr != nil {
		return nil, rpcErr
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	owner, err := h.svcLookupAgent(ctx, chainID, agent)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return map[string]string{"ownerAddress": owner}, nil
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
		ChainID  int64    `json:"chainId"`
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
	chainID := h.resolveRPCChainID(p.ChainID)
	results, err := h.svcBatchAgentInfo(ctx, chainID, p.Agents, subnetNum)
	if err != nil {
		return nil, svcToRPC(err)
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
	chainID := h.resolveRPCChainID(p.ChainID)
	resp, err := h.svcGetBalance(ctx, chainID, address)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return resp, nil
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
	chainID := h.resolveRPCChainID(p.ChainID)
	result, err := h.svcGetStakePositions(ctx, chainID, address)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcStakingGetAllocations(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Address string `json:"address"`
		ChainID int64  `json:"chainId"`
		pageParams
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	limit, offset := parsePage(p.pageParams)
	result, err := h.svcGetAllocations(ctx, chainID, address, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
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
	chainID := h.resolveRPCChainID(p.ChainID)
	result, err := h.svcGetFrozen(ctx, chainID, address)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcStakingGetPending(_ context.Context, _ json.RawMessage) (any, *RPCErr) {
	return []struct{}{}, nil
}

// 跨链查询：不需要 chainId
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
	amount, err := h.svcGetAgentSubnetStake(ctx, agent, subnetNum)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return map[string]string{"amount": amount}, nil
}

// 跨链查询：不需要 chainId
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
	result, err := h.svcGetAgentSubnets(ctx, agent)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

// 跨链查询：不需要 chainId
func (h *Handler) rpcStakingGetSubnetTotalStake(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p subnetParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	subnetNum, rpcErr := parseSubnetNum(p.SubnetID)
	if rpcErr != nil {
		return nil, rpcErr
	}
	amount, err := h.svcGetSubnetTotalStake(ctx, subnetNum)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return map[string]string{"total": amount}, nil
}

// ═══════════════════════════════════════════════
// ── subnets ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcSubnetsList(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Status  string `json:"status"`
		ChainID int64  `json:"chainId"` // 0 = all chains (cross-chain)
		pageParams
	}
	_ = json.Unmarshal(raw, &p)
	limit, offset := parsePage(p.pageParams)
	// chainId=0（默认）返回所有链的子网
	result, err := h.svcListSubnets(ctx, p.ChainID, p.Status, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
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
	result, err := h.svcGetSubnet(ctx, subnetNum)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
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
	result, err := h.svcGetSubnetSkills(ctx, subnetNum)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
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
	result, err := h.svcGetSubnetEarnings(ctx, subnetNum, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
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
	result, err := h.svcGetSubnetAgentInfo(ctx, agent, subnetNum)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

// ═══════════════════════════════════════════════
// ── emission ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcEmissionGetCurrent(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		ChainID int64 `json:"chainId"`
	}
	_ = json.Unmarshal(raw, &p)
	chainID := h.resolveRPCChainID(p.ChainID)
	data, err := h.svcGetCurrentEmission(ctx, chainID)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return data, nil
}

func (h *Handler) rpcEmissionGetSchedule(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		ChainID int64 `json:"chainId"`
	}
	_ = json.Unmarshal(raw, &p)
	chainID := h.resolveRPCChainID(p.ChainID)
	result, err := h.svcGetEmissionSchedule(ctx, chainID)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcEmissionListEpochs(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p pageParams
	_ = json.Unmarshal(raw, &p)
	chainID := h.resolveRPCChainID(p.ChainID)
	limit, offset := parsePage(p)
	result, err := h.svcListEpochs(ctx, chainID, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

// ═══════════════════════════════════════════════
// ── tokens ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcTokensGetAWP(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		ChainID int64 `json:"chainId"`
	}
	_ = json.Unmarshal(raw, &p)
	chainID := h.resolveRPCChainID(p.ChainID)
	data, err := h.svcGetAWPInfo(ctx, chainID)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return data, nil
}

func (h *Handler) rpcTokensGetAWPGlobal(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	data, err := h.svcGetAWPInfoGlobal(ctx)
	if err != nil {
		return nil, svcToRPC(err)
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
	result, err := h.svcGetAlphaInfo(ctx, subnetNum)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcTokensGetAlphaPrice(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p subnetParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	if _, rpcErr := parseSubnetNum(p.SubnetID); rpcErr != nil {
		return nil, rpcErr
	}
	data, err := h.svcGetAlphaPrice(ctx, p.SubnetID)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return data, nil
}

// ═══════════════════════════════════════════════
// ── governance ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcGovernanceListProposals(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Status  string `json:"status"`
		ChainID int64  `json:"chainId"`
		pageParams
	}
	_ = json.Unmarshal(raw, &p)
	chainID := h.resolveRPCChainID(p.ChainID)
	limit, offset := parsePage(p.pageParams)
	result, err := h.svcListProposals(ctx, chainID, p.Status, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcGovernanceGetProposal(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		ProposalID string `json:"proposalId"`
		ChainID    int64  `json:"chainId"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	if p.ProposalID == "" {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "proposalId is required"}
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	result, err := h.svcGetProposal(ctx, chainID, p.ProposalID)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcGovernanceTreasury(_ context.Context, _ json.RawMessage) (any, *RPCErr) {
	return map[string]string{"treasuryAddress": h.cfg.TreasuryAddress}, nil
}

// ═══════════════════════════════════════════════
// ── cross-chain global ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcStatsGlobal(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	result, err := h.svcGetGlobalStats(ctx)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcStakingGetUserBalanceGlobal(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p addressParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}
	resp, err := h.svcGetUserBalanceGlobal(ctx, address)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return resp, nil
}

func (h *Handler) rpcEmissionGetGlobalSchedule(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	result, err := h.svcGetGlobalEmissionSchedule(ctx)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcUsersListGlobal(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p pageParams
	_ = json.Unmarshal(raw, &p)
	limit, offset := parsePage(p)
	result, err := h.svcListUsersGlobal(ctx, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcGovernanceListAllProposals(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Status string `json:"status"`
		pageParams
	}
	_ = json.Unmarshal(raw, &p)
	limit, offset := parsePage(p.pageParams)
	result, err := h.svcListAllProposals(ctx, p.Status, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}
