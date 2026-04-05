package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/go-chi/chi/v5"
)

// startedAt records process start time, used for calculating uptime
var startedAt = time.Now()

// requireAdmin verifies Bearer token identity. Returns true if auth passed, false if error response was written.
func (h *Handler) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	if h.cfg.AdminToken == "" {
		h.writeError(w, http.StatusServiceUnavailable, "admin API not configured")
		return false
	}
	auth := r.Header.Get("Authorization")
	if auth == "" || auth != "Bearer "+h.cfg.AdminToken {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return false
	}
	return true
}

// ════════════════════════════════════════════════════════════
// Chain Management
// ════════════════════════════════════════════════════════════

// adminChainResponse is admin chain info including rpc_url (hidden in public endpoint)
type adminChainResponse struct {
	ChainID      int64  `json:"chainId"`
	Name         string `json:"name"`
	RpcUrl       string `json:"rpcUrl"`
	Dex          string `json:"dex"`
	Explorer     string `json:"explorer"`
	Status       string `json:"status"`
	AwpRegistry  string `json:"awpRegistry"`
	AwpToken     string `json:"awpToken"`
	AwpEmission  string `json:"awpEmission"`
	AwpAllocator string `json:"awpAllocator"`
	VeAWP        string `json:"veAWP"`
	AwpWorkNet   string `json:"awpWorkNet"`
	DaoAddress   string `json:"daoAddress"`
	LpManager    string `json:"lpManager"`
	PoolManager  string `json:"poolManager"`
	DeployBlock  int64  `json:"deployBlock"`
	CreatedAt    string `json:"createdAt,omitempty"`
}

// toAdminChainResponse converts a DB Chain to admin response (includes rpc_url)
func toAdminChainResponse(c gen.Chain) adminChainResponse {
	var createdAt string
	if c.CreatedAt.Valid {
		createdAt = c.CreatedAt.Time.Format(time.RFC3339)
	}
	return adminChainResponse{
		ChainID:      c.ChainID,
		Name:         c.Name,
		RpcUrl:       c.RpcUrl,
		Dex:          c.Dex,
		Explorer:     c.Explorer,
		Status:       c.Status,
		AwpRegistry:  c.AwpRegistry,
		AwpToken:     c.AwpToken,
		AwpEmission:  c.AwpEmission,
		AwpAllocator: c.AwpAllocator,
		VeAWP:        c.Veawp,
		AwpWorkNet:   c.AwpWorknet,
		DaoAddress:   c.DaoAddress,
		LpManager:    c.LpManager,
		PoolManager:  c.PoolManager,
		DeployBlock:  c.DeployBlock,
		CreatedAt:    createdAt,
	}
}

// AdminListChains lists all chains (including inactive), rpc_url visible
func (h *Handler) AdminListChains(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	chains, err := h.queries.ListAllChains(r.Context())
	if err != nil {
		h.logger.Error("admin: failed to list chains", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to list chains")
		return
	}
	resp := make([]adminChainResponse, len(chains))
	for i, c := range chains {
		resp[i] = toAdminChainResponse(c)
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// adminAddChainRequest is the request body for adding a chain
type adminAddChainRequest struct {
	ChainID      int64  `json:"chainId"`
	Name         string `json:"name"`
	RpcUrl       string `json:"rpcUrl"`
	Dex          string `json:"dex"`
	Explorer     string `json:"explorer"`
	AwpRegistry  string `json:"awpRegistry"`
	AwpToken     string `json:"awpToken"`
	AwpEmission  string `json:"awpEmission"`
	AwpAllocator string `json:"awpAllocator"`
	VeAWP        string `json:"veAWP"`
	AwpWorkNet   string `json:"awpWorkNet"`
	DaoAddress   string `json:"daoAddress"`
	LpManager    string `json:"lpManager"`
	PoolManager  string `json:"poolManager"`
	DeployBlock  int64  `json:"deployBlock"`
}

// AdminAddChain adds a new chain
func (h *Handler) AdminAddChain(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 65536)
	var req adminAddChainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ChainID <= 0 {
		h.writeError(w, http.StatusBadRequest, "chainId must be positive")
		return
	}
	if req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	err := h.queries.InsertChain(r.Context(), gen.InsertChainParams{
		ChainID:      req.ChainID,
		Name:         req.Name,
		RpcUrl:       req.RpcUrl,
		Dex:          req.Dex,
		Explorer:     req.Explorer,
		AwpRegistry:  req.AwpRegistry,
		AwpToken:     req.AwpToken,
		AwpEmission:  req.AwpEmission,
		AwpAllocator: req.AwpAllocator,
		Veawp:        req.VeAWP,
		AwpWorknet:   req.AwpWorkNet,
		DaoAddress:   req.DaoAddress,
		LpManager:    req.LpManager,
		PoolManager:  req.PoolManager,
		DeployBlock:  req.DeployBlock,
	})
	if err != nil {
		h.logger.Error("admin: failed to insert chain", "error", err, "chainId", req.ChainID)
		h.writeError(w, http.StatusConflict, "failed to insert chain (duplicate chainId?)")
		return
	}
	h.writeJSON(w, http.StatusCreated, map[string]any{
		"chainId": req.ChainID,
		"status":  "added",
		"note":    "Chain metadata added to DB. Indexer and keeper require process restart with updated chains.yaml to begin syncing this chain.",
	})
}

// AdminDeleteChain deletes a chain
func (h *Handler) AdminDeleteChain(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	raw := chi.URLParam(r, "chainId")
	chainID, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || chainID <= 0 {
		h.writeError(w, http.StatusBadRequest, "invalid chainId")
		return
	}
	if err := h.queries.DeleteChain(r.Context(), chainID); err != nil {
		h.logger.Error("admin: failed to delete chain", "error", err, "chainId", chainID)
		h.writeError(w, http.StatusInternalServerError, "failed to delete chain")
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]any{
		"chainId": chainID,
		"status":  "deactivated",
		"note":    "Chain marked inactive in DB. Indexer and keeper require process restart to stop syncing this chain.",
	})
}

// ════════════════════════════════════════════════════════════
// Rate Limit Management
// ════════════════════════════════════════════════════════════

// adminRateLimitRequest is the request for updating rate limit configuration
type adminRateLimitRequest struct {
	Key    string `json:"key"`
	Limit  int    `json:"limit"`
	Window int    `json:"window"`
}

// AdminUpdateRateLimit updates rate limit configuration (writes to Redis HSET ratelimit:config)
func (h *Handler) AdminUpdateRateLimit(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var req adminRateLimitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Key == "" {
		h.writeError(w, http.StatusBadRequest, "key is required")
		return
	}
	if req.Limit <= 0 || req.Window <= 0 {
		h.writeError(w, http.StatusBadRequest, "limit and window must be positive")
		return
	}

	value := fmt.Sprintf("%d:%d", req.Limit, req.Window)
	if err := h.rdb.HSet(r.Context(), "ratelimit:config", req.Key, value).Err(); err != nil {
		h.logger.Error("admin: failed to update ratelimit", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to update ratelimit")
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// AdminGetRateLimit fetches all rate limit configurations
func (h *Handler) AdminGetRateLimit(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	result, err := h.rdb.HGetAll(r.Context(), "ratelimit:config").Result()
	if err != nil {
		h.logger.Error("admin: failed to get ratelimit config", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get ratelimit config")
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// ════════════════════════════════════════════════════════════
// System Info
// ════════════════════════════════════════════════════════════

// AdminSystemInfo returns extended system info (Go version, uptime, goroutines, memory, Redis, DB connection pool, etc.)
func (h *Handler) AdminSystemInfo(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	ctx := r.Context()

	// Memory statistics
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	info := map[string]any{
		"goVersion":      runtime.Version(),
		"uptime":         time.Since(startedAt).String(),
		"uptimeSeconds":  int64(time.Since(startedAt).Seconds()),
		"goroutines":     runtime.NumGoroutine(),
		"numCPU":         runtime.NumCPU(),
		"memAllocMB":     mem.Alloc / 1024 / 1024,
		"memSysMB":       mem.Sys / 1024 / 1024,
		"memTotalAllocMB": mem.TotalAlloc / 1024 / 1024,
		"numGC":          mem.NumGC,
	}

	// Redis info
	redisInfo := map[string]string{}
	if pong, err := h.rdb.Ping(ctx).Result(); err == nil {
		redisInfo["status"] = pong
	} else {
		redisInfo["status"] = "error"
		redisInfo["error"] = err.Error()
	}
	if dbSize, err := h.rdb.DBSize(ctx).Result(); err == nil {
		redisInfo["keys"] = fmt.Sprintf("%d", dbSize)
	}
	info["redis"] = redisInfo

	// DB connection pool (via simple connectivity check)
	dbInfo := map[string]string{}
	chainID := h.defaultChainID()
	if cnt, err := h.queries.GetUserCount(ctx, chainID); err == nil {
		dbInfo["status"] = "ok"
		dbInfo["userCount"] = fmt.Sprintf("%d", cnt)
	} else {
		dbInfo["status"] = "error"
		dbInfo["error"] = err.Error()
	}
	info["database"] = dbInfo

	// Connected chains
	connectedChains := make([]int64, 0, len(h.chainReaders))
	for cid := range h.chainReaders {
		connectedChains = append(connectedChains, cid)
	}
	info["connectedChains"] = connectedChains

	h.writeJSON(w, http.StatusOK, info)
}
