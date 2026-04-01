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

	"github.com/cortexia/rootnet/api/internal/config"
	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/cortexia/rootnet/api/internal/ratelimit"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

// Handler is the main handler holding all dependencies
type Handler struct {
	queries *gen.Queries
	rdb     *redis.Client
	cfg     *config.Config
	logger  *slog.Logger
	limiter *ratelimit.Limiter
	chain   ChainReader           // optional: for on-chain reads (nonce, etc.)
	chains  []config.ChainConfig // loaded chains (nil in single-chain mode)

	// JSON-RPC 方法表缓存（sync.Once 初始化）
	rpcMethodsOnce    sync.Once
	rpcMethodTable    map[string]methodEntry
	rpcDiscoverResult []methodInfo
}

// ChainReader provides read-only access to on-chain state (optional dependency)
type ChainReader interface {
	GetNonce(addr string) (uint64, error)
	GetStakingNonce(addr string) (uint64, error)
}

// NewHandler creates a new Handler instance
func NewHandler(queries *gen.Queries, rdb *redis.Client, cfg *config.Config, logger *slog.Logger, limiter *ratelimit.Limiter) *Handler {
	return &Handler{
		queries: queries,
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

// SetChainReader sets the optional chain reader for on-chain queries
func (h *Handler) SetChainReader(cr ChainReader) {
	h.chain = cr
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
	limit = 20
	offset = 0

	if v := r.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > 100 {
		limit = 100
	}

	if v := r.URL.Query().Get("page"); v != "" {
		if page, err := strconv.Atoi(v); err == nil && page >= 1 && page <= 10000 {
			offset = (page - 1) * limit
		}
	}

	return limit, offset
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

// Health is the health-check endpoint
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// buildDetailedHealth gathers per-chain health, Redis connectivity, and DB connectivity.
// Shared by both HealthDetailed (REST) and rpcHealthDetailed (JSON-RPC).
func (h *Handler) buildDetailedHealth(ctx context.Context) map[string]interface{} {
	health := map[string]interface{}{
		"status": "ok",
	}

	// Determine which chain IDs to check
	chainIDs := []int64{h.cfg.ChainID}
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
		if ttl, err := h.rdb.TTL(ctx, emissionKey).Result(); err == nil && ttl > 0 {
			ch["keeperCacheAlive"] = true
		} else {
			ch["keeperCacheAlive"] = false
		}

		chainsHealth = append(chainsHealth, ch)
	}
	health["chains"] = chainsHealth

	// Redis 连通性
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		health["redis"] = "error"
		health["status"] = "degraded"
	} else {
		health["redis"] = "ok"
	}

	// 数据库连通性
	if _, err := h.queries.GetUserCount(ctx, h.cfg.ChainID); err != nil {
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
	ChainID           int64               `json:"chainId"`
	AWPRegistry       string              `json:"awpRegistry"`
	AWPToken          string              `json:"awpToken"`
	AWPEmission       string              `json:"awpEmission"`
	StakingVault      string              `json:"stakingVault"`
	StakeNFT          string              `json:"stakeNFT"`
	SubnetNFT         string              `json:"subnetNFT"`
	LPManager         string              `json:"lpManager"`
	AlphaTokenFactory string              `json:"alphaTokenFactory"`
	DAO               string              `json:"dao"`
	Treasury          string              `json:"treasury"`
	EIP712Domain      eip712DomainResponse `json:"eip712Domain"`
	StakingVaultEIP712 eip712DomainResponse `json:"stakingVaultEip712Domain"`
}

// GetRegistry returns the contract address registry with chain ID
func (h *Handler) GetRegistry(w http.ResponseWriter, r *http.Request) {
	resp := registryResponse{
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
			Name:              "AWPRegistry",
			Version:           "1",
			ChainID:           h.cfg.ChainID,
			VerifyingContract: h.cfg.AWPRegistryAddress,
		},
		StakingVaultEIP712: eip712DomainResponse{
			Name:              "StakingVault",
			Version:           "1",
			ChainID:           h.cfg.ChainID,
			VerifyingContract: h.cfg.StakingVaultAddress,
		},
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// checkAddressResponse is the response type for an address lookup check (V2)
type checkAddressResponse struct {
	IsRegistered bool   `json:"isRegistered"`
	BoundTo      string `json:"boundTo"`
	Recipient    string `json:"recipient"`
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
	if h.chain == nil {
		h.writeError(w, http.StatusServiceUnavailable, "chain reader not available")
		return
	}
	nonce, err := h.chain.GetNonce(address)
	if err != nil {
		h.logger.Error("failed to read nonce", "error", err, "address", address)
		h.writeError(w, http.StatusInternalServerError, "failed to read nonce")
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]uint64{"nonce": nonce})
}

// GetStakingNonce returns the EIP-712 nonce from StakingVault (for gasless allocate/deallocate)
func (h *Handler) GetStakingNonce(w http.ResponseWriter, r *http.Request) {
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
	if h.chain == nil {
		h.writeError(w, http.StatusServiceUnavailable, "chain reader not available")
		return
	}
	nonce, err := h.chain.GetStakingNonce(address)
	if err != nil {
		h.logger.Error("failed to read staking nonce", "error", err, "address", address)
		h.writeError(w, http.StatusInternalServerError, "failed to read staking nonce")
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]uint64{"nonce": nonce})
}

// GetChains returns the list of supported chains
func (h *Handler) GetChains(w http.ResponseWriter, r *http.Request) {
	if h.chains == nil {
		// 单链模式 — 仅返回当前配置的链
		h.writeJSON(w, http.StatusOK, []map[string]interface{}{
			{"chainId": h.cfg.ChainID, "name": "Default"},
		})
		return
	}
	h.writeJSON(w, http.StatusOK, h.chains) // ChainConfig 有 json 标签
}

// CheckAddress checks whether an address is registered and returns binding/recipient info (V2)
func (h *Handler) CheckAddress(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "address")
	if !isValidAddress(raw) {
		h.writeError(w, http.StatusBadRequest, "invalid address parameter")
		return
	}
	address := normalizeAddr(raw)

	ctx := r.Context()
	resp := checkAddressResponse{}

	// Check whether this address exists in the users table
	if user, err := h.queries.GetUser(ctx, gen.GetUserParams{
		Address: address,
		ChainID: h.cfg.ChainID,
	}); err == nil {
		resp.IsRegistered = user.RegisteredAt != 0 || user.BoundTo != "" || user.Recipient != ""
		resp.BoundTo = user.BoundTo
		resp.Recipient = user.Recipient
	} else if !errors.Is(err, pgx.ErrNoRows) {
		h.logger.Error("failed to check address", "error", err, "address", address)
		h.writeError(w, http.StatusInternalServerError, "failed to check address")
		return
	}

	h.writeJSON(w, http.StatusOK, resp)
}
