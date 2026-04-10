package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cortexia/rootnet/api/internal/config"
	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/cortexia/rootnet/api/internal/ratelimit"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

// Handler is the main handler holding all dependencies
type Handler struct {
	queries      *gen.Queries
	db           gen.DBTX              // raw DB access for queries not covered by sqlc
	rdb          *redis.Client
	cfg          *config.Config
	logger       *slog.Logger
	limiter      *ratelimit.Limiter
	chainReaders map[int64]ChainReader // chainId → on-chain reader (nonce, etc.)
	chains       []config.ChainConfig  // loaded chains (nil in single-chain mode)

	// JSON-RPC method table cache (initialized via sync.Once)
	rpcMethodsOnce    sync.Once
	rpcMethodTable    map[string]methodEntry
	rpcDiscoverResult []methodInfo
}

// ChainReader provides read-only access to on-chain state (optional dependency)
type ChainReader interface {
	GetNonce(addr string) (uint64, error)
	GetAllocatorNonce(addr string) (uint64, error)
	GetPermitNonce(addr string) (uint64, error)
	GetAWPTokenAddress() string
	ResolveRecipient(addr string) (string, error)
	BatchResolveRecipients(addrs []string) ([]string, error)
}

// NewHandler creates a new Handler instance
func NewHandler(queries *gen.Queries, db gen.DBTX, rdb *redis.Client, cfg *config.Config, logger *slog.Logger, limiter *ratelimit.Limiter) *Handler {
	return &Handler{
		queries: queries,
		db:      db,
		rdb:     rdb,
		cfg:     cfg,
		logger:  logger,
		limiter: limiter,
	}
}

// SetChains sets the loaded chain configurations for multi-chain name resolution
func (h *Handler) SetChains(chains []config.ChainConfig) {
	h.chains = chains
}

// SetChainReader adds a chain reader for the given chain ID
func (h *Handler) SetChainReader(chainID int64, cr ChainReader) {
	if h.chainReaders == nil {
		h.chainReaders = make(map[int64]ChainReader)
	}
	h.chainReaders[chainID] = cr
}

// GetChainReaders returns all registered chain readers (for sharing with RelayHandler)
func (h *Handler) GetChainReaders() map[int64]ChainReader {
	return h.chainReaders
}

// getChainReader returns the chain reader for a given chain ID, with single-chain fallback
func (h *Handler) getChainReader(chainID int64) ChainReader {
	if h.chainReaders == nil {
		return nil
	}
	if cr, ok := h.chainReaders[chainID]; ok {
		return cr
	}
	// Single-chain fallback: return the only reader when chainId is 0 or omitted
	if chainID == 0 && len(h.chainReaders) == 1 {
		for _, cr := range h.chainReaders {
			return cr
		}
	}
	return nil
}

// writeJSON writes data as JSON to the response
func (h *Handler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to write JSON response", "error", err)
	}
}

// writeError writes an error message as JSON to the response
func (h *Handler) writeError(w http.ResponseWriter, status int, msg string) {
	h.writeJSON(w, status, map[string]string{"error": msg})
}

// parsePageParams parses pagination parameters from the request; defaults are limit=20, offset=0
func (h *Handler) parsePageParams(r *http.Request) (limit, offset int) {
	page := 1
	lim := 20

	if v := r.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			lim = parsed
		}
	}
	if v := r.URL.Query().Get("page"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 1 {
			page = parsed
		}
	}

	l, o := computePageLimits(page, lim)
	return int(l), int(o)
}

// writeSvcError converts a svcError to an HTTP error response
func (h *Handler) writeSvcError(w http.ResponseWriter, err error) {
	var se *svcError
	if errors.As(err, &se) {
		switch se.Kind {
		case errNotFound:
			h.writeError(w, http.StatusNotFound, se.Message)
		case errBadInput:
			h.writeError(w, http.StatusBadRequest, se.Message)
		case errUnavailable:
			h.writeError(w, http.StatusServiceUnavailable, se.Message)
		default:
			h.writeError(w, http.StatusInternalServerError, se.Message)
		}
		return
	}
	h.writeError(w, http.StatusInternalServerError, "internal error")
}

// normalizeAddr converts an address to lowercase for consistency
func normalizeAddr(addr string) string {
	return strings.ToLower(addr)
}

// isValidAddress checks if a string is a valid Ethereum address (0x + 40 hex chars)
func isValidAddress(addr string) bool {
	if len(addr) != 42 || addr[:2] != "0x" {
		return false
	}
	for _, c := range addr[2:] {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// Version is set at build time via -ldflags "-X ...Version=xxx"
var Version = "dev"

// Health is the health-check endpoint (checks DB + Redis connectivity for load balancer use)
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	status := "ok"

	chainID := h.defaultChainID()
	if _, err := h.queries.GetUserCount(ctx, chainID); err != nil {
		status = "degraded"
	}
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		status = "degraded"
	}

	code := http.StatusOK
	if status != "ok" {
		code = http.StatusServiceUnavailable
	}
	h.writeJSON(w, code, map[string]string{"status": status, "version": Version})
}

// buildDetailedHealth gathers per-chain health, Redis connectivity, and DB connectivity.
// Shared by both HealthDetailed (REST) and rpcHealthDetailed (JSON-RPC).
func (h *Handler) buildDetailedHealth(ctx context.Context) map[string]interface{} {
	health := map[string]interface{}{
		"status": "ok",
	}

	// Determine which chain IDs to check
	chainIDs := []int64{h.defaultChainID()}
	if h.chains != nil {
		chainIDs = make([]int64, len(h.chains))
		for i, c := range h.chains {
			chainIDs[i] = c.ChainID
		}
	}

	// Per-chain health: indexer sync + keeper cache
	chainsHealth := make([]map[string]interface{}, 0, len(chainIDs))
	for _, cid := range chainIDs {
		if cid == 0 {
			continue
		}
		ch := map[string]interface{}{"chainId": cid}

		if syncState, err := h.queries.GetSyncState(ctx, gen.GetSyncStateParams{
			ChainID:      cid,
			ContractName: "indexer",
		}); err == nil {
			ch["indexerLastBlock"] = syncState.LastBlock
		}

		emissionKey := fmt.Sprintf("emission_current:%d", cid)
		// go-redis TTL returns: -2s = key doesn't exist, -1s = no expiry (persistent), >0 = TTL remaining.
		// A persistent key (ttl == -1s) is still alive, so we check ttl >= -1s (i.e., not -2s).
		if ttl, err := h.rdb.TTL(ctx, emissionKey).Result(); err == nil && ttl >= -1*time.Second {
			ch["keeperCacheAlive"] = true
		} else {
			ch["keeperCacheAlive"] = false
		}

		// Relayer native token balance (updated by keeper every 25s)
		balanceKey := fmt.Sprintf("relayer_balance:%d", cid)
		if bal, err := h.rdb.Get(ctx, balanceKey).Result(); err == nil {
			ch["relayerBalance"] = bal
		}

		chainsHealth = append(chainsHealth, ch)
	}
	health["chains"] = chainsHealth

	// Redis connectivity
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		health["redis"] = "error"
		health["status"] = "degraded"
	} else {
		health["redis"] = "ok"
	}

	// Database connectivity
	if _, err := h.queries.GetUserCount(ctx, h.defaultChainID()); err != nil {
		health["database"] = "error"
		health["status"] = "degraded"
	} else {
		health["database"] = "ok"
	}

	return health
}

// HealthDetailed returns detailed system health including indexer lag, keeper status, chain info
func (h *Handler) HealthDetailed(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, h.buildDetailedHealth(r.Context()))
}

// eip712DomainResponse provides all info needed to construct EIP-712 signatures for gasless relay
type eip712DomainResponse struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	ChainID           int64  `json:"chainId"`
	VerifyingContract string `json:"verifyingContract"`
}

// registryResponse is the response type for the contract address registry
type registryResponse struct {
	ChainID             int64                `json:"chainId"`
	AWPRegistry         string               `json:"awpRegistry"`
	AWPToken            string               `json:"awpToken"`
	AWPEmission         string               `json:"awpEmission"`
	AWPAllocator        string               `json:"awpAllocator"`
	VeAWP               string               `json:"veAWP"`
	AWPWorkNet          string               `json:"awpWorkNet"`
	LPManager           string               `json:"lpManager"`
	WorknetTokenFactory string               `json:"worknetTokenFactory"`
	DAO                 string               `json:"dao"`
	Treasury            string               `json:"treasury"`
	VeAWPHelper         string               `json:"veAWPHelper"`
	EIP712Domain        eip712DomainResponse `json:"eip712Domain"`
	AllocatorEIP712     eip712DomainResponse `json:"allocatorEip712Domain"`
}

// GetRegistry returns contract address registry. Without ?chainId, returns all chains.
func (h *Handler) GetRegistry(w http.ResponseWriter, r *http.Request) {
	if v := r.URL.Query().Get("chainId"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil && id > 0 {
			h.writeJSON(w, http.StatusOK, h.svcGetRegistry(r.Context(), id))
			return
		}
	}
	// No chainId — return all chains
	h.writeJSON(w, http.StatusOK, h.svcGetRegistryAll(r.Context()))
}

// chainRegistration represents registration info on a single chain
type chainRegistration struct {
	ChainID      int64  `json:"chainId"`
	IsRegistered bool   `json:"isRegistered"`
	BoundTo      string `json:"boundTo,omitempty"`
	Recipient    string `json:"recipient,omitempty"`
}

// checkAddressResponse is the response type for an address lookup check (V2)
type checkAddressResponse struct {
	IsRegistered bool                `json:"isRegistered"`
	BoundTo      string              `json:"boundTo,omitempty"`
	Recipient    string              `json:"recipient,omitempty"`
	Chains       []chainRegistration `json:"chains,omitempty"`
}

// GetNonce returns the EIP-712 nonce for an address (used for gasless relay signature construction)
func (h *Handler) GetNonce(w http.ResponseWriter, r *http.Request) {
	// Rate limit nonce lookups to prevent abuse as an oracle
	ip := ratelimit.GetClientIP(r)
	if exceeded, err := h.limiter.CheckAndIncrement(r.Context(), "nonce", ip); exceeded {
		h.writeError(w, http.StatusTooManyRequests, h.limiter.FormatError(r.Context(), "nonce"))
		return
	} else if err != nil {
		h.logger.Error("nonce rate limit error", "error", err)
	}

	raw := chi.URLParam(r, "address")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	address := normalizeAddr(raw)
	chainID := h.resolveChainID(r)
	cr := h.getChainReader(chainID)
	if cr == nil {
		h.writeError(w, http.StatusServiceUnavailable, "chain reader not available for chainId")
		return
	}
	nonce, err := cr.GetNonce(address)
	if err != nil {
		h.logger.Error("failed to read nonce", "error", err, "address", address, "chainId", chainID)
		h.writeError(w, http.StatusInternalServerError, "failed to read nonce")
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]uint64{"nonce": nonce})
}

// GetAllocatorNonce returns the EIP-712 nonce from AWPAllocator (for gasless allocate/deallocate)
func (h *Handler) GetAllocatorNonce(w http.ResponseWriter, r *http.Request) {
	ip := ratelimit.GetClientIP(r)
	if exceeded, err := h.limiter.CheckAndIncrement(r.Context(), "nonce", ip); exceeded {
		h.writeError(w, http.StatusTooManyRequests, h.limiter.FormatError(r.Context(), "nonce"))
		return
	} else if err != nil {
		h.logger.Error("nonce rate limit error", "error", err)
	}

	raw := chi.URLParam(r, "address")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	address := normalizeAddr(raw)
	chainID := h.resolveChainID(r)
	cr := h.getChainReader(chainID)
	if cr == nil {
		h.writeError(w, http.StatusServiceUnavailable, "chain reader not available for chainId")
		return
	}
	nonce, err := cr.GetAllocatorNonce(address)
	if err != nil {
		h.logger.Error("failed to read allocator nonce", "error", err, "address", address, "chainId", chainID)
		h.writeError(w, http.StatusInternalServerError, "failed to read allocator nonce")
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]uint64{"nonce": nonce})
}

// GetChains returns the list of supported chains
func (h *Handler) GetChains(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, h.svcGetChains(r.Context()))
}

// ResolveRecipient returns the resolved recipient for an address (walks the on-chain bind chain to root)
func (h *Handler) ResolveRecipient(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "address")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid address")
		return
	}
	chainID := h.resolveChainID(r)
	cr := h.getChainReader(chainID)
	if cr == nil {
		h.writeError(w, http.StatusServiceUnavailable, "chain reader not available")
		return
	}
	resolved, err := cr.ResolveRecipient(normalizeAddr(raw))
	if err != nil {
		h.logger.Error("failed to resolve recipient", "error", err, "address", raw)
		h.writeError(w, http.StatusInternalServerError, "failed to resolve recipient")
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]string{"address": normalizeAddr(raw), "resolvedRecipient": strings.ToLower(resolved)})
}

// CheckAddress checks whether an address is registered and returns binding/recipient info (V2)
func (h *Handler) CheckAddress(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "address")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid address parameter")
		return
	}
	chainID := h.resolveChainID(r)
	resp, err := h.svcCheckAddress(r.Context(), chainID, normalizeAddr(raw))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// GetDB returns the raw DB connection for announcement routes
func (h *Handler) GetDB() gen.DBTX {
	return h.db
}

// AdminAuthMiddleware is a chi middleware that checks Bearer token auth
func (h *Handler) AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.cfg.AdminToken == "" {
			h.writeError(w, http.StatusServiceUnavailable, "admin API not configured")
			return
		}
		auth := r.Header.Get("Authorization")
		if auth == "" || auth != "Bearer "+h.cfg.AdminToken {
			h.writeError(w, http.StatusUnauthorized, "invalid or missing admin token")
			return
		}
		next.ServeHTTP(w, r)
	})
}
