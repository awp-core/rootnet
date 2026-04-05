package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/cortexia/rootnet/api/internal/ratelimit"
)

// JSON-RPC 2.0 error codes
const (
	rpcParseError     = -32700
	rpcInvalidRequest = -32600
	rpcMethodNotFound = -32601
	rpcInvalidParams  = -32602
	rpcInternalError  = -32603
	rpcNotFound       = -32001 // Application layer: resource not found
)

// rpcCtxKey is used to store request info in context
type rpcCtxKey struct{}

// rpcCtxVal stores RPC request metadata (e.g., client IP)
type rpcCtxVal struct {
	ClientIP string
}

// RPCRequest JSON-RPC 2.0 request
type RPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      any             `json:"id"`
}

// RPCResponse JSON-RPC 2.0 response
type RPCResponse struct {
	JSONRPC string   `json:"jsonrpc"`
	Result  any      `json:"result,omitempty"`
	Error   *RPCErr  `json:"error,omitempty"`
	ID      any      `json:"id"`
}

// RPCErr JSON-RPC 2.0 error object
type RPCErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// methodFunc is the method handler function signature
type methodFunc func(ctx context.Context, params json.RawMessage) (any, *RPCErr)

// paramInfo is parameter metadata (used for rpc.discover)
type paramInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// methodInfo is method metadata (used for rpc.discover)
type methodInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Params      []paramInfo `json:"params"`
}

// methodEntry is a method registry entry
type methodEntry struct {
	fn   methodFunc
	info methodInfo
}

// HandleRPC handles JSON-RPC 2.0 requests on /v2
// GET -> returns rpc.discover (API documentation)
// POST -> processes JSON-RPC requests (supports single and batch requests)
func (h *Handler) HandleRPC(w http.ResponseWriter, r *http.Request) {
	// GET /v2 -> return rpc.discover result (allows browsing available methods)
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

	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB

	// Inject client IP into context (used for rate limiting)
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

	// Detect whether this is a batch request
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
		// Execute batch requests concurrently
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

// rpcClientIP extracts the client IP from context
func rpcClientIP(ctx context.Context) string {
	if v, ok := ctx.Value(rpcCtxKey{}).(*rpcCtxVal); ok {
		return v.ClientIP
	}
	return ""
}

// initRPCMethods initializes the method registry (cached on first call)
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

// dispatchRPC dispatches a single JSON-RPC request
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

// rpcMethods returns the method registry
func (h *Handler) rpcMethods() map[string]methodEntry {
	return map[string]methodEntry{
		// ── stats ──
		"stats.global": {fn: h.rpcStatsGlobal, info: methodInfo{
			Name: "stats.global", Description: "Get global protocol statistics across all chains", Params: []paramInfo{},
		}},

		// ── registry ──
		"registry.get": {fn: h.rpcRegistryGet, info: methodInfo{
			Name: "registry.get", Description: "Get all contract addresses and EIP-712 domain info", Params: []paramInfo{},
		}},

		// ── health ──
		"health.check": {fn: h.rpcHealthCheck, info: methodInfo{
			Name: "health.check", Description: "Health check", Params: []paramInfo{},
		}},
		"health.detailed": {fn: h.rpcHealthDetailed, info: methodInfo{
			Name: "health.detailed", Description: "Detailed health status (per-chain indexer/keeper)", Params: []paramInfo{},
		}},

		// ── chains ──
		"chains.list": {fn: h.rpcChainsList, info: methodInfo{
			Name: "chains.list", Description: "List supported chains", Params: []paramInfo{},
		}},

		// ── users ──
		"users.list": {fn: h.rpcUsersList, info: methodInfo{
			Name: "users.list", Description: "List users (paginated)",
			Params: []paramInfo{
				{Name: "page", Type: "integer", Required: false, Description: "Page number (default 1)"},
				{Name: "limit", Type: "integer", Required: false, Description: "Items per page (default 20, max 100)"},
			},
		}},
		"users.listGlobal": {fn: h.rpcUsersListGlobal, info: methodInfo{
			Name: "users.listGlobal", Description: "List users across all chains (deduplicated, paginated)",
			Params: []paramInfo{
				{Name: "page", Type: "integer", Required: false, Description: "Page number (default 1)"},
				{Name: "limit", Type: "integer", Required: false, Description: "Items per page (default 20, max 100)"},
			},
		}},
		"users.count": {fn: h.rpcUsersCount, info: methodInfo{
			Name: "users.count", Description: "Get total user count", Params: []paramInfo{},
		}},
		"users.get": {fn: h.rpcUsersGet, info: methodInfo{
			Name: "users.get", Description: "Get user details (balance + bound agents)",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "User address (0x...)"},
			},
		}},

		// ── address ──
		"address.check": {fn: h.rpcAddressCheck, info: methodInfo{
			Name: "address.check", Description: "Check address registration, binding, and recipient",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "Address (0x...)"},
			},
		}},

		"address.batchResolveRecipients": {fn: h.rpcBatchResolveRecipients, info: methodInfo{
			Name: "address.batchResolveRecipients", Description: "Batch resolve effective recipients for multiple addresses (on-chain)",
			Params: []paramInfo{
				{Name: "addresses", Type: "array<string>", Required: true, Description: "Address list (max 500)"},
				{Name: "chainId", Type: "integer", Required: false, Description: "Chain ID for on-chain read"},
			},
		}},
		"address.resolveRecipient": {fn: h.rpcResolveRecipient, info: methodInfo{
			Name: "address.resolveRecipient", Description: "Resolve the effective recipient for an address (walks bind chain to root)",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "Address (0x...)"},
				{Name: "chainId", Type: "integer", Required: false, Description: "Chain ID (required for on-chain read)"},
			},
		}},

		// ── nonce ──
		"nonce.get": {fn: h.rpcNonceGet, info: methodInfo{
			Name: "nonce.get", Description: "Get AWPRegistry EIP-712 nonce",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "Address (0x...)"},
			},
		}},
		"nonce.getStaking": {fn: h.rpcNonceGetStaking, info: methodInfo{
			Name: "nonce.getStaking", Description: "Get AWPAllocator EIP-712 nonce",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "Address (0x...)"},
			},
		}},

		// ── agents ──
		"agents.getByOwner": {fn: h.rpcAgentsGetByOwner, info: methodInfo{
			Name: "agents.getByOwner", Description: "Get all agents bound to an owner",
			Params: []paramInfo{
				{Name: "owner", Type: "string", Required: true, Description: "Owner address (0x...)"},
			},
		}},
		"agents.getDetail": {fn: h.rpcAgentsGetDetail, info: methodInfo{
			Name: "agents.getDetail", Description: "Get agent details",
			Params: []paramInfo{
				{Name: "agent", Type: "string", Required: true, Description: "Agent address (0x...)"},
			},
		}},
		"agents.lookup": {fn: h.rpcAgentsLookup, info: methodInfo{
			Name: "agents.lookup", Description: "Look up agent owner",
			Params: []paramInfo{
				{Name: "agent", Type: "string", Required: true, Description: "Agent address (0x...)"},
			},
		}},
		"agents.batchInfo": {fn: h.rpcAgentsBatchInfo, info: methodInfo{
			Name: "agents.batchInfo", Description: "Batch query agent info and stake in a subnet",
			Params: []paramInfo{
				{Name: "agents", Type: "array<string>", Required: true, Description: "Agent address list (max 100)"},
				{Name: "worknetId", Type: "string", Required: true, Description: "Subnet ID"},
			},
		}},

		// ── staking ──
		"staking.getBalance": {fn: h.rpcStakingGetBalance, info: methodInfo{
			Name: "staking.getBalance", Description: "Get user AWP staking balance (staked/allocated/available)",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "User address (0x...)"},
			},
		}},
		"staking.getUserBalanceGlobal": {fn: h.rpcStakingGetUserBalanceGlobal, info: methodInfo{
			Name: "staking.getUserBalanceGlobal", Description: "Get user AWP staking balance aggregated across all chains",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "User address (0x...)"},
			},
		}},
		"staking.getPositions": {fn: h.rpcStakingGetPositions, info: methodInfo{
			Name: "staking.getPositions", Description: "Get user veAWP positions",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "User address (0x...)"},
			},
		}},
		"staking.getAllocations": {fn: h.rpcStakingGetAllocations, info: methodInfo{
			Name: "staking.getAllocations", Description: "Get user stake allocations (paginated)",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "User address (0x...)"},
				{Name: "page", Type: "integer", Required: false, Description: "Page number (default 1)"},
				{Name: "limit", Type: "integer", Required: false, Description: "Items per page (default 20, max 100)"},
			},
		}},
		"staking.getFrozen": {fn: h.rpcStakingGetFrozen, info: methodInfo{
			Name: "staking.getFrozen", Description: "Get user frozen allocations",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "User address (0x...)"},
			},
		}},
		"staking.getPending": {fn: h.rpcStakingGetPending, info: methodInfo{
			Name: "staking.getPending", Description: "Get pending allocation changes (always empty)",
			Params: []paramInfo{},
		}},
		"staking.getAgentSubnetStake": {fn: h.rpcStakingGetAgentSubnetStake, info: methodInfo{
			Name: "staking.getAgentSubnetStake", Description: "Get agent stake in a subnet",
			Params: []paramInfo{
				{Name: "agent", Type: "string", Required: true, Description: "Agent address (0x...)"},
				{Name: "worknetId", Type: "string", Required: true, Description: "Subnet ID"},
			},
		}},
		"staking.getAgentSubnets": {fn: h.rpcStakingGetAgentSubnets, info: methodInfo{
			Name: "staking.getAgentSubnets", Description: "Get all subnets an agent participates in",
			Params: []paramInfo{
				{Name: "agent", Type: "string", Required: true, Description: "Agent address (0x...)"},
			},
		}},
		"staking.getSubnetTotalStake": {fn: h.rpcStakingGetSubnetTotalStake, info: methodInfo{
			Name: "staking.getSubnetTotalStake", Description: "Get subnet total stake",
			Params: []paramInfo{
				{Name: "worknetId", Type: "string", Required: true, Description: "Subnet ID"},
			},
		}},

		// ── subnets ──
		"subnets.list": {fn: h.rpcSubnetsList, info: methodInfo{
			Name: "subnets.list", Description: "List subnets (paginated, optional status filter)",
			Params: []paramInfo{
				{Name: "status", Type: "string", Required: false, Description: "Status filter: Pending/Active/Paused/Banned"},
				{Name: "page", Type: "integer", Required: false, Description: "Page number (default 1)"},
				{Name: "limit", Type: "integer", Required: false, Description: "Items per page (default 20, max 100)"},
			},
		}},
		"subnets.get": {fn: h.rpcSubnetsGet, info: methodInfo{
			Name: "subnets.get", Description: "Get subnet details",
			Params: []paramInfo{
				{Name: "worknetId", Type: "string", Required: true, Description: "Subnet ID"},
			},
		}},
		"subnets.getSkills": {fn: h.rpcSubnetsGetSkills, info: methodInfo{
			Name: "subnets.getSkills", Description: "Get subnet skills URI",
			Params: []paramInfo{
				{Name: "worknetId", Type: "string", Required: true, Description: "Subnet ID"},
			},
		}},
		"subnets.getEarnings": {fn: h.rpcSubnetsGetEarnings, info: methodInfo{
			Name: "subnets.getEarnings", Description: "Get subnet AWP earnings history (paginated)",
			Params: []paramInfo{
				{Name: "worknetId", Type: "string", Required: true, Description: "Subnet ID"},
				{Name: "page", Type: "integer", Required: false, Description: "Page number (default 1)"},
				{Name: "limit", Type: "integer", Required: false, Description: "Items per page (default 20, max 100)"},
			},
		}},
		"subnets.getAgentInfo": {fn: h.rpcSubnetsGetAgentInfo, info: methodInfo{
			Name: "subnets.getAgentInfo", Description: "Get agent staking info in a subnet",
			Params: []paramInfo{
				{Name: "worknetId", Type: "string", Required: true, Description: "Subnet ID"},
				{Name: "agent", Type: "string", Required: true, Description: "Agent address (0x...)"},
			},
		}},

		// ── emission ──
		"emission.getCurrent": {fn: h.rpcEmissionGetCurrent, info: methodInfo{
			Name: "emission.getCurrent", Description: "Get current emission data", Params: []paramInfo{},
		}},
		"emission.getSchedule": {fn: h.rpcEmissionGetSchedule, info: methodInfo{
			Name: "emission.getSchedule", Description: "Get emission projections (30/90/365 days)", Params: []paramInfo{},
		}},
		"emission.getGlobalSchedule": {fn: h.rpcEmissionGetGlobalSchedule, info: methodInfo{
			Name: "emission.getGlobalSchedule", Description: "Get emission schedule aggregated across all chains", Params: []paramInfo{},
		}},
		"emission.listEpochs": {fn: h.rpcEmissionListEpochs, info: methodInfo{
			Name: "emission.listEpochs", Description: "List epochs (paginated)",
			Params: []paramInfo{
				{Name: "page", Type: "integer", Required: false, Description: "Page number (default 1)"},
				{Name: "limit", Type: "integer", Required: false, Description: "Items per page (default 20, max 100)"},
			},
		}},

		// ── tokens ──
		"tokens.getAWP": {fn: h.rpcTokensGetAWP, info: methodInfo{
			Name: "tokens.getAWP", Description: "Get AWP token info", Params: []paramInfo{},
		}},
		"tokens.getAWPGlobal": {fn: h.rpcTokensGetAWPGlobal, info: methodInfo{
			Name: "tokens.getAWPGlobal", Description: "Get AWP token info aggregated across all chains", Params: []paramInfo{},
		}},
		"tokens.getAlphaInfo": {fn: h.rpcTokensGetAlphaInfo, info: methodInfo{
			Name: "tokens.getAlphaInfo", Description: "Get subnet Alpha token info",
			Params: []paramInfo{
				{Name: "worknetId", Type: "string", Required: true, Description: "Subnet ID"},
			},
		}},
		"tokens.getAlphaPrice": {fn: h.rpcTokensGetAlphaPrice, info: methodInfo{
			Name: "tokens.getAlphaPrice", Description: "Get subnet Alpha token price",
			Params: []paramInfo{
				{Name: "worknetId", Type: "string", Required: true, Description: "Subnet ID"},
			},
		}},

		// ── governance ──
		"governance.listProposals": {fn: h.rpcGovernanceListProposals, info: methodInfo{
			Name: "governance.listProposals", Description: "List governance proposals (paginated, optional status filter)",
			Params: []paramInfo{
				{Name: "status", Type: "string", Required: false, Description: "Status filter: Active/Canceled/Defeated/Succeeded/Queued/Expired/Executed"},
				{Name: "page", Type: "integer", Required: false, Description: "Page number (default 1)"},
				{Name: "limit", Type: "integer", Required: false, Description: "Items per page (default 20, max 100)"},
			},
		}},
		"governance.listAllProposals": {fn: h.rpcGovernanceListAllProposals, info: methodInfo{
			Name: "governance.listAllProposals", Description: "List governance proposals across all chains (paginated, optional status filter)",
			Params: []paramInfo{
				{Name: "status", Type: "string", Required: false, Description: "Status filter: Active/Canceled/Defeated/Succeeded/Queued/Expired/Executed"},
				{Name: "page", Type: "integer", Required: false, Description: "Page number (default 1)"},
				{Name: "limit", Type: "integer", Required: false, Description: "Items per page (default 20, max 100)"},
			},
		}},
		"governance.getProposal": {fn: h.rpcGovernanceGetProposal, info: methodInfo{
			Name: "governance.getProposal", Description: "Get governance proposal details",
			Params: []paramInfo{
				{Name: "proposalId", Type: "string", Required: true, Description: "Proposal ID"},
			},
		}},
		"governance.getTreasury": {fn: h.rpcGovernanceTreasury, info: methodInfo{
			Name: "governance.getTreasury", Description: "Get Treasury contract address", Params: []paramInfo{},
		}},

		// ── users (new) ──
		"users.getPortfolio": {fn: h.rpcUsersGetPortfolio, info: methodInfo{
			Name: "users.getPortfolio", Description: "Get full user portfolio (identity, balance, positions, allocations, delegates)",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "User address (0x...)"},
				{Name: "chainId", Type: "integer", Required: false, Description: "Chain ID (default: primary chain)"},
			},
		}},
		"users.getDelegates": {fn: h.rpcUsersGetDelegates, info: methodInfo{
			Name: "users.getDelegates", Description: "Get agents bound to a user (delegate approximation from bind tree)",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "User address (0x...)"},
				{Name: "chainId", Type: "integer", Required: false, Description: "Chain ID (default: primary chain)"},
			},
		}},

		// ── subnets (new) ──
		"subnets.listRanked": {fn: h.rpcSubnetsListRanked, info: methodInfo{
			Name: "subnets.listRanked", Description: "List worknets ranked by total stake",
			Params: []paramInfo{
				{Name: "chainId", Type: "integer", Required: false, Description: "Chain ID (default: primary chain)"},
				{Name: "page", Type: "integer", Required: false, Description: "Page number (default 1)"},
				{Name: "limit", Type: "integer", Required: false, Description: "Items per page (default 20, max 100)"},
			},
		}},
		"subnets.listAgents": {fn: h.rpcSubnetsListAgents, info: methodInfo{
			Name: "subnets.listAgents", Description: "List agents in a worknet ranked by stake",
			Params: []paramInfo{
				{Name: "worknetId", Type: "string", Required: true, Description: "Subnet ID"},
				{Name: "chainId", Type: "integer", Required: false, Description: "Chain ID (default: primary chain)"},
				{Name: "page", Type: "integer", Required: false, Description: "Page number (default 1)"},
				{Name: "limit", Type: "integer", Required: false, Description: "Items per page (default 20, max 100)"},
			},
		}},
		"subnets.search": {fn: h.rpcSubnetsSearch, info: methodInfo{
			Name: "subnets.search", Description: "Search worknets by name or symbol (ILIKE)",
			Params: []paramInfo{
				{Name: "query", Type: "string", Required: true, Description: "Search query (1-100 chars)"},
				{Name: "chainId", Type: "integer", Required: false, Description: "Chain ID (default: primary chain)"},
				{Name: "page", Type: "integer", Required: false, Description: "Page number (default 1)"},
				{Name: "limit", Type: "integer", Required: false, Description: "Items per page (default 20, max 100)"},
			},
		}},
		"subnets.getByOwner": {fn: h.rpcSubnetsGetByOwner, info: methodInfo{
			Name: "subnets.getByOwner", Description: "Get worknets owned by an address",
			Params: []paramInfo{
				{Name: "owner", Type: "string", Required: true, Description: "Owner address (0x...)"},
				{Name: "chainId", Type: "integer", Required: false, Description: "Chain ID (default: primary chain)"},
				{Name: "page", Type: "integer", Required: false, Description: "Page number (default 1)"},
				{Name: "limit", Type: "integer", Required: false, Description: "Items per page (default 20, max 100)"},
			},
		}},

		// ── emission (new) ──
		"emission.getEpochDetail": {fn: h.rpcEmissionGetEpochDetail, info: methodInfo{
			Name: "emission.getEpochDetail", Description: "Get epoch detail with per-recipient distributions",
			Params: []paramInfo{
				{Name: "epochId", Type: "integer", Required: true, Description: "Epoch ID"},
				{Name: "chainId", Type: "integer", Required: false, Description: "Chain ID (default: primary chain)"},
			},
		}},

		// ── staking (new) ──
		"staking.getPositionsGlobal": {fn: h.rpcStakingGetPositionsGlobal, info: methodInfo{
			Name: "staking.getPositionsGlobal", Description: "Get user veAWP positions across all chains",
			Params: []paramInfo{
				{Name: "address", Type: "string", Required: true, Description: "User address (0x...)"},
			},
		}},
	}
}
