package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/cortexia/rootnet/api/internal/config"
	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/cortexia/rootnet/api/internal/ratelimit"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

// Handler is the main handler holding all dependencies
type Handler struct {
	queries *gen.Queries
	rdb     *redis.Client
	cfg     *config.Config
	logger  *slog.Logger
	limiter *ratelimit.Limiter
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
		if page, err := strconv.Atoi(v); err == nil && page >= 1 {
			offset = (page - 1) * limit
		}
	}

	return limit, offset
}

// normalizeAddr converts an address to lowercase for consistency
func normalizeAddr(addr string) string {
	return strings.ToLower(addr)
}

// Health is the health-check endpoint
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// registryResponse is the response type for the contract address registry
type registryResponse struct {
	ChainID           int64  `json:"chainId"`
	AWPRegistry       string `json:"awpRegistry"`
	AWPToken          string `json:"awpToken"`
	AWPEmission       string `json:"awpEmission"`
	StakingVault      string `json:"stakingVault"`
	StakeNFT          string `json:"stakeNFT"`
	SubnetNFT         string `json:"subnetNFT"`
	AccessManager     string `json:"accessManager"`
	LPManager         string `json:"lpManager"`
	AlphaTokenFactory string `json:"alphaTokenFactory"`
	DAO               string `json:"dao"`
	Treasury          string `json:"treasury"`
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
		AccessManager:     h.cfg.AccessManagerAddress,
		LPManager:         h.cfg.LPManagerAddress,
		AlphaTokenFactory: h.cfg.AlphaFactoryAddress,
		DAO:               h.cfg.DAOAddress,
		Treasury:          h.cfg.TreasuryAddress,
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// checkAddressResponse is the response type for an address lookup check
type checkAddressResponse struct {
	IsRegisteredUser  bool   `json:"isRegisteredUser"`
	IsRegisteredAgent bool   `json:"isRegisteredAgent"`
	OwnerAddress      string `json:"ownerAddress,omitempty"`
	IsManager         bool   `json:"isManager"`
}

// CheckAddress checks whether an address is a registered user or agent
func (h *Handler) CheckAddress(w http.ResponseWriter, r *http.Request) {
	address := normalizeAddr(chi.URLParam(r, "address"))
	if address == "" || len(address) != 42 {
		h.writeError(w, http.StatusBadRequest, "invalid address parameter")
		return
	}

	ctx := r.Context()
	resp := checkAddressResponse{}

	// Check whether this is a registered user
	if _, err := h.queries.GetUser(ctx, address); err == nil {
		resp.IsRegisteredUser = true
	}

	// Check whether this is a registered agent
	if agent, err := h.queries.GetAgent(ctx, address); err == nil {
		resp.IsRegisteredAgent = true
		resp.OwnerAddress = agent.OwnerAddress
		resp.IsManager = agent.IsManager
	}

	h.writeJSON(w, http.StatusOK, resp)
}
