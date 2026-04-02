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

// newRelayHandler creates a multi-chain RelayHandler (optional: returns nil if RELAYER_PRIVATE_KEY not set)
func newRelayHandler(lc fx.Lifecycle, cfg *config.Config, rdb *redis.Client, limiter *ratelimit.Limiter, logger *slog.Logger) (*handler.RelayHandler, error) {
	if cfg.RelayerPrivateKey == "" {
		logger.Info("RELAYER_PRIVATE_KEY not set, relay endpoints disabled")
		return nil, nil
	}

	key, err := crypto.HexToECDSA(cfg.RelayerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid RELAYER_PRIVATE_KEY: %w", err)
	}

	relayers := make(map[int64]*chain.Relayer)
	var clients []*ethclient.Client

	// 收集要连接的 RPC URL 列表（单链 or 多链）
	type chainRPC struct {
		rpcURL       string
		awpRegistry  string
		stakingVault string
	}
	var rpcs []chainRPC

	if cfg.ChainsFile != "" {
		chains, err := config.LoadChains(cfg.ChainsFile)
		if err != nil {
			return nil, fmt.Errorf("load chains for relay: %w", err)
		}
		for _, ch := range chains {
			rpcs = append(rpcs, chainRPC{
				rpcURL:       ch.RPCURL,
				awpRegistry:  config.ResolveAddress(ch.AWPRegistry, cfg.AWPRegistryAddress),
				stakingVault: config.ResolveAddress(ch.StakingVault, cfg.StakingVaultAddress),
			})
		}
	} else if cfg.RPCURL != "" {
		rpcs = append(rpcs, chainRPC{
			rpcURL:       cfg.RPCURL,
			awpRegistry:  cfg.AWPRegistryAddress,
			stakingVault: cfg.StakingVaultAddress,
		})
	}

	for _, rpc := range rpcs {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, dialErr := ethclient.DialContext(ctx, rpc.rpcURL)
		cancel()
		if dialErr != nil {
			logger.Warn("relay: failed to dial RPC, skipping", "rpc", rpc.rpcURL, "error", dialErr)
			continue
		}

		ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
		chainID, cidErr := client.ChainID(ctx2)
		cancel2()
		if cidErr != nil {
			client.Close()
			logger.Warn("relay: failed to get chainID, skipping", "error", cidErr)
			continue
		}

		awpRegistryAddr := common.HexToAddress(rpc.awpRegistry)
		stakingVaultAddr := common.HexToAddress(rpc.stakingVault)
		rl, rlErr := chain.NewRelayer(client, awpRegistryAddr, stakingVaultAddr, key, chainID, rdb, logger)
		if rlErr != nil {
			client.Close()
			logger.Warn("relay: failed to create relayer, skipping", "chainId", chainID, "error", rlErr)
			continue
		}

		relayers[chainID.Int64()] = rl
		clients = append(clients, client)
		logger.Info("relay enabled for chain", "chainId", chainID.Int64())
	}

	if len(relayers) == 0 {
		logger.Warn("no relay chains configured")
		return nil, nil
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			for _, c := range clients {
				c.Close()
			}
			return nil
		},
	})

	return handler.NewRelayHandler(relayers, limiter, logger), nil
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
	// 默认 chainID：多链时用第一条链，单链时用 cfg.ChainID
	vanityChainID := cfg.ChainID
	if cfg.ChainsFile != "" {
		if chains, err := config.LoadChains(cfg.ChainsFile); err == nil && len(chains) > 0 {
			vanityChainID = chains[0].ChainID
		}
	}
	return handler.NewVanityHandler(cfg.AlphaFactoryAddress, cfg.AlphaInitCodeHash, rule, vanityChainID, queries, limiter, logger)
}

// wireChainReader creates a lightweight chain client for on-chain reads (nonce, etc.)
func wireChainReader(lc fx.Lifecycle, h *handler.Handler, hub *ws.Hub, queries *gen.Queries, cfg *config.Config, logger *slog.Logger) {
	// 注入 WebSocket 分配查询接口（chainID 从事件中获取，fallback 用第一条配置链）
	defaultCID := cfg.ChainID
	if cfg.ChainsFile != "" {
		if chains, err := config.LoadChains(cfg.ChainsFile); err == nil && len(chains) > 0 {
			defaultCID = chains[0].ChainID
		}
	}
	hub.SetAllocationQuerier(handler.NewWSAllocQuerier(queries), defaultCID)

	var clients []*chain.Client

	// 多链模式：为每条链创建独立 chain reader
	if cfg.ChainsFile != "" {
		chains, loadErr := config.LoadChains(cfg.ChainsFile)
		if loadErr == nil && len(chains) > 0 {
			h.SetChains(chains)
			for _, ch := range chains {
				addrs := map[string]string{
					"AWPRegistry":  config.ResolveAddress(ch.AWPRegistry, cfg.AWPRegistryAddress),
					"AWPToken":     config.ResolveAddress(ch.AWPToken, cfg.AWPTokenAddress),
					"AWPEmission":  config.ResolveAddress(ch.AWPEmission, cfg.AWPEmissionAddress),
					"StakingVault": config.ResolveAddress(ch.StakingVault, cfg.StakingVaultAddress),
					"SubnetNFT":    config.ResolveAddress(ch.SubnetNFT, cfg.SubnetNFTAddress),
					"AWPDAO":       config.ResolveAddress(ch.DAOAddress, cfg.DAOAddress),
					"StakeNFT":     config.ResolveAddress(ch.StakeNFT, cfg.StakeNFTAddress),
				}
				dialCtx, dialCancel := context.WithTimeout(context.Background(), 10*time.Second)
				client, dialErr := chain.NewClient(dialCtx, ch.RPCURL, addrs)
				dialCancel()
				if dialErr != nil {
					logger.Warn("chain reader unavailable for chain", "chainId", ch.ChainID, "error", dialErr)
					continue
				}
				h.SetChainReader(ch.ChainID, client)
				clients = append(clients, client)
				logger.Info("chain reader connected", "chainId", ch.ChainID)
			}
		} else if loadErr != nil {
			logger.Warn("failed to load chains config", "error", loadErr)
		}
	} else if cfg.RPCURL != "" {
		// 单链模式
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
		h.SetChainReader(cfg.ChainID, client)
		clients = append(clients, client)
		logger.Info("chain reader connected (single-chain)", "chainId", cfg.ChainID)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			for _, c := range clients {
				c.Close()
			}
			return nil
		},
	})
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
