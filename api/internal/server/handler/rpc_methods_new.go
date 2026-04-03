package handler

import (
	"context"
	"encoding/json"
)

// ═══════════════════════════════════════════════
// ── users.getPortfolio ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcUsersGetPortfolio(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p addressParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	resp, err := h.svcGetPortfolio(ctx, chainID, address)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return resp, nil
}

// ═══════════════════════════════════════════════
// ── users.getDelegates ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcUsersGetDelegates(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p addressParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	result, err := h.svcGetDelegates(ctx, chainID, address)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

// ═══════════════════════════════════════════════
// ── subnets.listRanked ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcSubnetsListRanked(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p pageParams
	_ = json.Unmarshal(raw, &p)
	chainID := h.resolveRPCChainID(p.ChainID)
	limit, offset := parsePage(p)
	result, err := h.svcListWorknetsRanked(ctx, chainID, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

// ═══════════════════════════════════════════════
// ── emission.getEpochDetail ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcEmissionGetEpochDetail(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		EpochID int64 `json:"epochId"`
		ChainID int64 `json:"chainId"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	if p.EpochID < 0 {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "epochId must be >= 0"}
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	result, err := h.svcGetEpochDetail(ctx, chainID, p.EpochID)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

// ═══════════════════════════════════════════════
// ── subnets.listAgents ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcSubnetsListAgents(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		SubnetID string `json:"worknetId"`
		ChainID  int64  `json:"chainId"`
		pageParams
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	subnetNum, rpcErr := parseSubnetNum(p.SubnetID)
	if rpcErr != nil {
		return nil, rpcErr
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	limit, offset := parsePage(p.pageParams)
	result, err := h.svcListWorknetAgents(ctx, chainID, subnetNum, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

// ═══════════════════════════════════════════════
// ── subnets.search ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcSubnetsSearch(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Query   string `json:"query"`
		ChainID int64  `json:"chainId"`
		pageParams
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	if p.Query == "" {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "query is required"}
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	limit, offset := parsePage(p.pageParams)
	result, err := h.svcSearchWorknets(ctx, chainID, p.Query, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

// ═══════════════════════════════════════════════
// ── subnets.getByOwner ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcSubnetsGetByOwner(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p struct {
		Owner   string `json:"owner"`
		ChainID int64  `json:"chainId"`
		pageParams
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	owner, rpcErr := requireAddress(p.Owner)
	if rpcErr != nil {
		return nil, rpcErr
	}
	chainID := h.resolveRPCChainID(p.ChainID)
	limit, offset := parsePage(p.pageParams)
	result, err := h.svcGetWorknetsByOwner(ctx, chainID, owner, limit, offset)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}

// ═══════════════════════════════════════════════
// ── staking.getPositionsGlobal ──
// ═══════════════════════════════════════════════

func (h *Handler) rpcStakingGetPositionsGlobal(ctx context.Context, raw json.RawMessage) (any, *RPCErr) {
	var p addressParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCErr{Code: rpcInvalidParams, Message: "invalid params"}
	}
	address, rpcErr := requireAddress(p.Address)
	if rpcErr != nil {
		return nil, rpcErr
	}
	result, err := h.svcGetPositionsGlobal(ctx, address)
	if err != nil {
		return nil, svcToRPC(err)
	}
	return result, nil
}
