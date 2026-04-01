package handler

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
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

func (h *Handler) rpcRegistryGet(_ context.Context, _ json.RawMessage) (any, *RPCErr) {
	return h.svcGetRegistry(), nil
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
	return h.svcGetChains(), nil
}

// ═══════════════════════════════════════════════
// ── users ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcUsersList(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p pageParams
	_ = json.Unmarshal(raw, &p)
	limit, offset := parsePage(p)
	result, err := h.svcListUsers(ctx, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcUsersCount(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	count, err := h.svcGetUserCount(ctx)
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
	resp, err := h.svcGetUser(ctx, address)
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
	resp, err := h.svcCheckAddress(ctx, address)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return resp, nil
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
	result, err := h.svcGetAgentsByOwner(ctx, owner)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
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
	result, err := h.svcGetAgentDetail(ctx, agent)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
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
	owner, err := h.svcLookupAgent(ctx, agent)
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
	results, err := h.svcBatchAgentInfo(ctx, p.Agents, subnetNum)
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
	resp, err := h.svcGetBalance(ctx, address)
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
	result, err := h.svcGetStakePositions(ctx, address)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
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
	result, err := h.svcGetAllocations(ctx, address, limit, offset)
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
	result, err := h.svcGetFrozen(ctx, address)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcStakingGetPending(_ context.Context, _ json.RawMessage) (any, *RPCErr) {
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
	amount, err := h.svcGetAgentSubnetStake(ctx, agent, subnetNum)
	if err != nil {
		return nil, svcToRPC(err)
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
	result, err := h.svcGetAgentSubnets(ctx, agent)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
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
		Status string `json:"status"`
		pageParams
	}
	_ = json.Unmarshal(raw, &p)
	limit, offset := parsePage(p.pageParams)
	result, err := h.svcListSubnets(ctx, p.Status, limit, offset)
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

func (h *Handler) rpcEmissionGetCurrent(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	data, err := h.svcGetCurrentEmission(ctx)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return data, nil
}

func (h *Handler) rpcEmissionGetSchedule(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	result, err := h.svcGetEmissionSchedule(ctx)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcEmissionListEpochs(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p pageParams
	_ = json.Unmarshal(raw, &p)
	limit, offset := parsePage(p)
	result, err := h.svcListEpochs(ctx, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

// ═══════════════════════════════════════════════
// ── tokens ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcTokensGetAWP(ctx context.Context, _ json.RawMessage) (any, *RPCErr) {
	data, err := h.svcGetAWPInfo(ctx)
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
		Status string `json:"status"`
		pageParams
	}
	_ = json.Unmarshal(raw, &p)
	limit, offset := parsePage(p.pageParams)
	result, err := h.svcListProposals(ctx, p.Status, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
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
	result, err := h.svcGetProposal(ctx, p.ProposalID)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

func (h *Handler) rpcGovernanceTreasury(_ context.Context, _ json.RawMessage) (any, *RPCErr) {
	return map[string]string{"treasuryAddress": h.cfg.TreasuryAddress}, nil
}
