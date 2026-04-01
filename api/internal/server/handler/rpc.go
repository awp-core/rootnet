package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/cortexia/rootnet/api/internal/ratelimit"
)

// JSON-RPC 2.0 错误码
const (
	rpcParseError     = -32700
	rpcInvalidRequest = -32600
	rpcMethodNotFound = -32601
	rpcInvalidParams  = -32602
	rpcInternalError  = -32603
	rpcNotFound       = -32001 // 应用层：资源不存在
)

// rpcCtxKey 用于在 context 中存储请求信息
type rpcCtxKey struct{}

// rpcCtxVal 存储 RPC 请求的元信息（如客户端 IP）
type rpcCtxVal struct {
	ClientIP string
}

// RPCRequest JSON-RPC 2.0 请求
type RPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      any             `json:"id"`
}

// RPCResponse JSON-RPC 2.0 响应
type RPCResponse struct {
	JSONRPC string   `json:"jsonrpc"`
	Result  any      `json:"result,omitempty"`
	Error   *RPCErr  `json:"error,omitempty"`
	ID      any      `json:"id"`
}

// RPCErr JSON-RPC 2.0 错误对象
type RPCErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// methodFunc 方法处理函数签名
type methodFunc func(ctx context.Context, params json.RawMessage) (any, *RPCErr)

// paramInfo 参数元数据（用于 rpc.discover）
type paramInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// methodInfo 方法元数据（用于 rpc.discover）
type methodInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Params      []paramInfo `json:"params"`
}

// methodEntry 方法注册表条目
type methodEntry struct {
	fn   methodFunc
	info methodInfo
}

// HandleRPC 处理 /v2 的 JSON-RPC 2.0 请求
// GET → 返回 rpc.discover（API 文档）
// POST → 处理 JSON-RPC 请求（支持单请求和批量请求）
func (h *Handler) HandleRPC(w http.ResponseWriter, r *http.Request) {
	// GET /v2 → 返回 rpc.discover 结果（方便浏览器查看可用方法）
	if r.Method == http.MethodGet {
		h.initRPCMethods()
		h.writeJSON(w, http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"result":  map[string]any{"methods": h.rpcDiscoverResult},
			"id":      nil,
		})
		return
	}
	if r.Method != http.MethodPost {
		h.writeJSON(w, http.StatusMethodNotAllowed, RPCResponse{
			JSONRPC: "2.0",
			Error:   &RPCErr{Code: rpcInvalidRequest, Message: "method not allowed, use POST"},
		})
		return
	}

	// 限制请求体大小
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB

	// 注入客户端 IP 到 context（用于 rate limiting）
	ctx := context.WithValue(r.Context(), rpcCtxKey{}, &rpcCtxVal{
		ClientIP: ratelimit.GetClientIP(r),
	})

	var raw json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		h.writeJSON(w, http.StatusOK, RPCResponse{
			JSONRPC: "2.0",
			Error:   &RPCErr{Code: rpcParseError, Message: "parse error"},
		})
		return
	}

	// 检测是否为批量请求
	if len(raw) > 0 && raw[0] == '[' {
		var reqs []RPCRequest
		if err := json.Unmarshal(raw, &reqs); err != nil {
			h.writeJSON(w, http.StatusOK, RPCResponse{
				JSONRPC: "2.0",
				Error:   &RPCErr{Code: rpcParseError, Message: "parse error"},
			})
			return
		}
		if len(reqs) == 0 {
			h.writeJSON(w, http.StatusOK, RPCResponse{
				JSONRPC: "2.0",
				Error:   &RPCErr{Code: rpcInvalidRequest, Message: "empty batch"},
			})
			return
		}
		if len(reqs) > 20 {
			h.writeJSON(w, http.StatusOK, RPCResponse{
				JSONRPC: "2.0",
				Error:   &RPCErr{Code: rpcInvalidRequest, Message: "batch size exceeds limit (20)"},
			})
			return
		}
		// 批量请求并发执行
		responses := make([]RPCResponse, len(reqs))
		var wg sync.WaitGroup
		for i, req := range reqs {
			wg.Add(1)
			go func(idx int, r RPCRequest) {
				defer wg.Done()
				responses[idx] = h.dispatchRPC(ctx, r)
			}(i, req)
		}
		wg.Wait()
		h.writeJSON(w, http.StatusOK, responses)
		return
	}

	var req RPCRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		h.writeJSON(w, http.StatusOK, RPCResponse{
			JSONRPC: "2.0",
			Error:   &RPCErr{Code: rpcParseError, Message: "parse error"},
		})
		return
	}

	resp := h.dispatchRPC(ctx, req)
	h.writeJSON(w, http.StatusOK, resp)
}

// rpcClientIP 从 context 中提取客户端 IP
func rpcClientIP(ctx context.Context) string {
	if v, ok := ctx.Value(rpcCtxKey{}).(*rpcCtxVal); ok {
		return v.ClientIP
	}
	return ""
}

// initRPCMethods 初始化方法注册表（在首次调用时缓存）
func (h *Handler) initRPCMethods() {
	h.rpcMethodsOnce.Do(func() {
		h.rpcMethodTable = h.rpcMethods()
		h.rpcDiscoverResult = func() []methodInfo {
			infos := make([]methodInfo, 0, len(h.rpcMethodTable))
			for _, m := range h.rpcMethodTable {
				infos = append(infos, m.info)
			}
			return infos
		}()
	})
}

// dispatchRPC 分发单个 JSON-RPC 请求
func (h *Handler) dispatchRPC(ctx context.Context, req RPCRequest) RPCResponse {
	if req.JSONRPC != "2.0" {
		return RPCResponse{JSONRPC: "2.0", Error: &RPCErr{Code: rpcInvalidRequest, Message: "jsonrpc must be \"2.0\""}, ID: req.ID}
	}

	h.initRPCMethods()

	if req.Method == "rpc.discover" {
		return RPCResponse{JSONRPC: "2.0", Result: map[string]any{"methods": h.rpcDiscoverResult}, ID: req.ID}
	}

	entry, ok := h.rpcMethodTable[req.Method]
	if !ok {
		return RPCResponse{JSONRPC: "2.0", Error: &RPCErr{Code: rpcMethodNotFound, Message: "method not found: " + req.Method}, ID: req.ID}
	}

	result, rpcErr := entry.fn(ctx, req.Params)
	if rpcErr != nil {
		return RPCResponse{JSONRPC: "2.0", Error: rpcErr, ID: req.ID}
	}
	return RPCResponse{JSONRPC: "2.0", Result: result, ID: req.ID}
}

// rpcMethods 返回方法注册表
func (h *Handler) rpcMethods() map[string]methodEntry {
	return map[string]methodEntry{
		// ── registry ──
		"registry.get": {fn: h.rpcRegistryGet, info: methodInfo{
			Name: "registry.get", Description: "获取所有合约地址和 EIP-712 域信息", Params: []paramInfo{},
		}},

		// ── health ──
		"health.check": {fn: h.rpcHealthCheck, info: methodInfo{
			Name: "health.check", Description: "健康检查", Params: []paramInfo{},
		}},
		"health.detailed": {fn: h.rpcHealthDetailed, info: methodInfo{
			Name: "health.detailed", Description: "详细健康状态（含各链 indexer/keeper 状态）", Params: []paramInfo{},
		}},

		// ── chains ──
		"chains.list": {fn: h.rpcChainsList, info: methodInfo{
			Name: "chains.list", Description: "获取支持的链列表", Params: []paramInfo{},
		}},

		// ── users ──
		"users.list": {fn: h.rpcUsersList, info: methodInfo{
			Name: "users.list", Description: "分页获取用户列表",
			Params: []paramInfo{
				{Name: "page", Type: "integer", Required: false, Description: "页码（默认 1）"},
				{Name: "limit", Type: "integer", Required: false, Description: "每页数量（默认 20，最大 100）"},
			},
		}},
		"users.count": {fn: h.rpcUsersCount, info: methodInfo{
			Name: "users.count", Description: "获取用户总数", Params: []paramInfo{},
		}},
		"users.get": {fn: h.rpcUsersGet, info: methodInfo{
			Name: "users.get", Description: "获取用户详情（含余额和绑定的 agent）",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "用户地址 (0x...)"},
			},
		}},

		// ── address ──
		"address.check": {fn: h.rpcAddressCheck, info: methodInfo{
			Name: "address.check", Description: "检查地址注册状态、绑定和收款地址",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "地址 (0x...)"},
			},
		}},

		// ── nonce ──
		"nonce.get": {fn: h.rpcNonceGet, info: methodInfo{
			Name: "nonce.get", Description: "获取 AWPRegistry EIP-712 nonce",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "地址 (0x...)"},
			},
		}},
		"nonce.getStaking": {fn: h.rpcNonceGetStaking, info: methodInfo{
			Name: "nonce.getStaking", Description: "获取 StakingVault EIP-712 nonce",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "地址 (0x...)"},
			},
		}},

		// ── agents ──
		"agents.getByOwner": {fn: h.rpcAgentsGetByOwner, info: methodInfo{
			Name: "agents.getByOwner", Description: "获取绑定到指定 owner 的所有 agent",
			Params: []paramInfo{
				{Name: "owner", Type: "string", Required: true, Description: "owner 地址 (0x...)"},
			},
		}},
		"agents.getDetail": {fn: h.rpcAgentsGetDetail, info: methodInfo{
			Name: "agents.getDetail", Description: "获取 agent 详情",
			Params: []paramInfo{
				{Name: "agent", Type: "string", Required: true, Description: "agent 地址 (0x...)"},
			},
		}},
		"agents.lookup": {fn: h.rpcAgentsLookup, info: methodInfo{
			Name: "agents.lookup", Description: "查找 agent 的 owner",
			Params: []paramInfo{
				{Name: "agent", Type: "string", Required: true, Description: "agent 地址 (0x...)"},
			},
		}},
		"agents.batchInfo": {fn: h.rpcAgentsBatchInfo, info: methodInfo{
			Name: "agents.batchInfo", Description: "批量查询 agent 信息及在指定子网的质押量",
			Params: []paramInfo{
				{Name: "agents", Type: "array<string>", Required: true, Description: "agent 地址列表（最多 100 个）"},
				{Name: "subnetId", Type: "string", Required: true, Description: "子网 ID"},
			},
		}},

		// ── staking ──
		"staking.getBalance": {fn: h.rpcStakingGetBalance, info: methodInfo{
			Name: "staking.getBalance", Description: "获取用户 AWP 质押余额（总质押/已分配/可用）",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "用户地址 (0x...)"},
			},
		}},
		"staking.getPositions": {fn: h.rpcStakingGetPositions, info: methodInfo{
			Name: "staking.getPositions", Description: "获取用户的 StakeNFT 持仓列表",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "用户地址 (0x...)"},
			},
		}},
		"staking.getAllocations": {fn: h.rpcStakingGetAllocations, info: methodInfo{
			Name: "staking.getAllocations", Description: "获取用户的质押分配列表（分页）",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "用户地址 (0x...)"},
				{Name: "page", Type: "integer", Required: false, Description: "页码（默认 1）"},
				{Name: "limit", Type: "integer", Required: false, Description: "每页数量（默认 20，最大 100）"},
			},
		}},
		"staking.getFrozen": {fn: h.rpcStakingGetFrozen, info: methodInfo{
			Name: "staking.getFrozen", Description: "获取用户冻结的质押分配",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "用户地址 (0x...)"},
			},
		}},
		"staking.getPending": {fn: h.rpcStakingGetPending, info: methodInfo{
			Name: "staking.getPending", Description: "获取待处理的分配变更（当前总是空数组）",
			Params: []paramInfo{},
		}},
		"staking.getAgentSubnetStake": {fn: h.rpcStakingGetAgentSubnetStake, info: methodInfo{
			Name: "staking.getAgentSubnetStake", Description: "获取 agent 在某子网的质押量",
			Params: []paramInfo{
				{Name: "agent", Type: "string", Required: true, Description: "agent 地址 (0x...)"},
				{Name: "subnetId", Type: "string", Required: true, Description: "子网 ID"},
			},
		}},
		"staking.getAgentSubnets": {fn: h.rpcStakingGetAgentSubnets, info: methodInfo{
			Name: "staking.getAgentSubnets", Description: "获取 agent 参与的所有子网及质押量",
			Params: []paramInfo{
				{Name: "agent", Type: "string", Required: true, Description: "agent 地址 (0x...)"},
			},
		}},
		"staking.getSubnetTotalStake": {fn: h.rpcStakingGetSubnetTotalStake, info: methodInfo{
			Name: "staking.getSubnetTotalStake", Description: "获取子网总质押量",
			Params: []paramInfo{
				{Name: "subnetId", Type: "string", Required: true, Description: "子网 ID"},
			},
		}},

		// ── subnets ──
		"subnets.list": {fn: h.rpcSubnetsList, info: methodInfo{
			Name: "subnets.list", Description: "分页获取子网列表（可按状态筛选）",
			Params: []paramInfo{
				{Name: "status", Type: "string", Required: false, Description: "状态筛选: Pending/Active/Paused/Banned"},
				{Name: "page", Type: "integer", Required: false, Description: "页码（默认 1）"},
				{Name: "limit", Type: "integer", Required: false, Description: "每页数量（默认 20，最大 100）"},
			},
		}},
		"subnets.get": {fn: h.rpcSubnetsGet, info: methodInfo{
			Name: "subnets.get", Description: "获取子网详情",
			Params: []paramInfo{
				{Name: "subnetId", Type: "string", Required: true, Description: "子网 ID"},
			},
		}},
		"subnets.getSkills": {fn: h.rpcSubnetsGetSkills, info: methodInfo{
			Name: "subnets.getSkills", Description: "获取子网 skills URI",
			Params: []paramInfo{
				{Name: "subnetId", Type: "string", Required: true, Description: "子网 ID"},
			},
		}},
		"subnets.getEarnings": {fn: h.rpcSubnetsGetEarnings, info: methodInfo{
			Name: "subnets.getEarnings", Description: "获取子网 AWP 收益历史（分页）",
			Params: []paramInfo{
				{Name: "subnetId", Type: "string", Required: true, Description: "子网 ID"},
				{Name: "page", Type: "integer", Required: false, Description: "页码（默认 1）"},
				{Name: "limit", Type: "integer", Required: false, Description: "每页数量（默认 20，最大 100）"},
			},
		}},
		"subnets.getAgentInfo": {fn: h.rpcSubnetsGetAgentInfo, info: methodInfo{
			Name: "subnets.getAgentInfo", Description: "获取 agent 在子网的质押信息",
			Params: []paramInfo{
				{Name: "subnetId", Type: "string", Required: true, Description: "子网 ID"},
				{Name: "agent", Type: "string", Required: true, Description: "agent 地址 (0x...)"},
			},
		}},

		// ── emission ──
		"emission.getCurrent": {fn: h.rpcEmissionGetCurrent, info: methodInfo{
			Name: "emission.getCurrent", Description: "获取当前排放数据", Params: []paramInfo{},
		}},
		"emission.getSchedule": {fn: h.rpcEmissionGetSchedule, info: methodInfo{
			Name: "emission.getSchedule", Description: "获取排放预测（30/90/365 天）", Params: []paramInfo{},
		}},
		"emission.listEpochs": {fn: h.rpcEmissionListEpochs, info: methodInfo{
			Name: "emission.listEpochs", Description: "分页获取 epoch 列表",
			Params: []paramInfo{
				{Name: "page", Type: "integer", Required: false, Description: "页码（默认 1）"},
				{Name: "limit", Type: "integer", Required: false, Description: "每页数量（默认 20，最大 100）"},
			},
		}},

		// ── tokens ──
		"tokens.getAWP": {fn: h.rpcTokensGetAWP, info: methodInfo{
			Name: "tokens.getAWP", Description: "获取 AWP 代币信息", Params: []paramInfo{},
		}},
		"tokens.getAlphaInfo": {fn: h.rpcTokensGetAlphaInfo, info: methodInfo{
			Name: "tokens.getAlphaInfo", Description: "获取子网 Alpha 代币信息",
			Params: []paramInfo{
				{Name: "subnetId", Type: "string", Required: true, Description: "子网 ID"},
			},
		}},
		"tokens.getAlphaPrice": {fn: h.rpcTokensGetAlphaPrice, info: methodInfo{
			Name: "tokens.getAlphaPrice", Description: "获取子网 Alpha 代币价格",
			Params: []paramInfo{
				{Name: "subnetId", Type: "string", Required: true, Description: "子网 ID"},
			},
		}},

		// ── governance ──
		"governance.listProposals": {fn: h.rpcGovernanceListProposals, info: methodInfo{
			Name: "governance.listProposals", Description: "分页获取治理提案（可按状态筛选）",
			Params: []paramInfo{
				{Name: "status", Type: "string", Required: false, Description: "状态筛选: Active/Canceled/Defeated/Succeeded/Queued/Expired/Executed"},
				{Name: "page", Type: "integer", Required: false, Description: "页码（默认 1）"},
				{Name: "limit", Type: "integer", Required: false, Description: "每页数量（默认 20，最大 100）"},
			},
		}},
		"governance.getProposal": {fn: h.rpcGovernanceGetProposal, info: methodInfo{
			Name: "governance.getProposal", Description: "获取单个治理提案详情",
			Params: []paramInfo{
				{Name: "proposalId", Type: "string", Required: true, Description: "提案 ID"},
			},
		}},
		"governance.getTreasury": {fn: h.rpcGovernanceTreasury, info: methodInfo{
			Name: "governance.getTreasury", Description: "获取 Treasury 合约地址", Params: []paramInfo{},
		}},
	}
}
