package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cortexia/rootnet/api/internal/chain"
	"github.com/go-chi/chi/v5"
	"github.com/cortexia/rootnet/api/internal/config"
	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/cortexia/rootnet/api/internal/ratelimit"
	"github.com/cortexia/rootnet/api/internal/server"
	"github.com/cortexia/rootnet/api/internal/server/handler"
	"github.com/cortexia/rootnet/api/internal/server/ws"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			newLogger,
			config.Load,
			newDBPool,
			newRedis,
			newQueries,
			newLimiter,
			handler.NewHandler,
			ws.NewHub,
			newRelayHandler,
			newVanityHandler,
			newRouterParams,
			server.NewRouter,
		),
		fx.Invoke(wireChainReader),
		fx.Invoke(startServer),
	).Run()
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func newDBPool(lc fx.Lifecycle, cfg *config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	poolCfg.MaxConns = 20
	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})
	return pool, nil
}

func newRedis(lc fx.Lifecycle, cfg *config.Config) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_URL: %w", err)
	}
	rdb := redis.NewClient(opt)
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return rdb.Close()
		},
	})
	return rdb, nil
}

func newQueries(pool *pgxpool.Pool) *gen.Queries {
	return gen.New(pool)
}

func newLimiter(rdb *redis.Client, logger *slog.Logger) *ratelimit.Limiter {
	return ratelimit.NewLimiter(rdb, logger)
}

// newRelayHandler creates RelayHandler (optional: returns nil if RELAYER_PRIVATE_KEY not set)
// If key is configured but invalid, returns error to fail fast
func newRelayHandler(lc fx.Lifecycle, cfg *config.Config, limiter *ratelimit.Limiter, logger *slog.Logger) (*handler.RelayHandler, error) {
	if cfg.RelayerPrivateKey == "" {
		logger.Info("RELAYER_PRIVATE_KEY not set, relay endpoints disabled")
		return nil, nil
	}

	key, err := crypto.HexToECDSA(cfg.RelayerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid RELAYER_PRIVATE_KEY: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("relay ethclient dial: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("relay get chainID: %w", err)
	}

	awpRegistryAddr := common.HexToAddress(cfg.AWPRegistryAddress)
	relayer, err := chain.NewRelayer(client, awpRegistryAddr, key, chainID, logger)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("create relayer: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			client.Close()
			return nil
		},
	})

	logger.Info("relay endpoints enabled", "chainID", chainID.String())
	return handler.NewRelayHandler(relayer, limiter, logger), nil
}

// newVanityHandler creates VanityHandler (optional: returns nil if not configured)
func newVanityHandler(cfg *config.Config, queries *gen.Queries, limiter *ratelimit.Limiter, logger *slog.Logger) *handler.VanityHandler {
	if cfg.AlphaFactoryAddress == "" || cfg.AlphaInitCodeHash == "" || cfg.VanityRule == "" {
		logger.Info("ALPHA_FACTORY_ADDRESS, ALPHA_INITCODE_HASH or VANITY_RULE not set, vanity endpoint disabled")
		return nil
	}

	// Validate initCodeHash format
	hashHex := strings.TrimPrefix(cfg.AlphaInitCodeHash, "0x")
	if _, err := hex.DecodeString(hashHex); err != nil || len(hashHex) != 64 {
		logger.Error("invalid ALPHA_INITCODE_HASH", "value", cfg.AlphaInitCodeHash)
		return nil
	}

	// Decode vanityRule
	rule, err := chain.DecodeVanityRule(cfg.VanityRule)
	if err != nil {
		logger.Error("invalid VANITY_RULE", "value", cfg.VanityRule, "error", err)
		return nil
	}
	if rule.IsEmpty() {
		logger.Info("VANITY_RULE is all wildcards, vanity endpoint disabled (no constraint)")
		return nil
	}

	logger.Info("vanity compute-salt endpoint enabled", "factory", cfg.AlphaFactoryAddress, "vanityRule", cfg.VanityRule)
	return handler.NewVanityHandler(cfg.AlphaFactoryAddress, cfg.AlphaInitCodeHash, rule, cfg.ChainID, queries, limiter, logger)
}

// wireChainReader creates a lightweight chain client for on-chain reads (nonce, etc.)
func wireChainReader(lc fx.Lifecycle, h *handler.Handler, cfg *config.Config, logger *slog.Logger) {
	addrs := map[string]string{
		"AWPRegistry":  cfg.AWPRegistryAddress,
		"AWPToken":     cfg.AWPTokenAddress,
		"AWPEmission":  cfg.AWPEmissionAddress,
		"StakingVault": cfg.StakingVaultAddress,
		"SubnetNFT":    cfg.SubnetNFTAddress,
		"AWPDAO":       cfg.DAOAddress,
		"StakeNFT":     cfg.StakeNFTAddress,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := chain.NewClient(ctx, cfg.RPCURL, addrs)
	if err != nil {
		logger.Warn("chain reader unavailable, /api/nonce endpoint disabled", "error", err)
		return
	}
	h.SetChainReader(client)
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			client.Close()
			return nil
		},
	})
	logger.Info("chain reader connected for on-chain queries")
}

// newRouterParams assembles RouterParams
func newRouterParams(cfg *config.Config, h *handler.Handler, hub *ws.Hub, rh *handler.RelayHandler, vh *handler.VanityHandler) server.RouterParams {
	return server.RouterParams{
		Config:        cfg,
		Handler:       h,
		Hub:           hub,
		RelayHandler:  rh,
		VanityHandler: vh,
	}
}

func startServer(lc fx.Lifecycle, router chi.Router, cfg *config.Config, hub *ws.Hub, logger *slog.Logger) {
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      150 * time.Second, // allows vanity compute-salt 120s timeout
	}

	var hubCancel context.CancelFunc

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			hubCtx, cancel := context.WithCancel(context.Background())
			hubCancel = cancel
			go hub.Run(hubCtx)

			logger.Info("HTTP server starting", "addr", cfg.HTTPAddr)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if hubCancel != nil {
				hubCancel()
			}
			return srv.Shutdown(ctx)
		},
	})
}
