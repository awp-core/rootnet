package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/cortexia/rootnet/api/internal/chain"
	"github.com/cortexia/rootnet/api/internal/config"
)

func main() {
	fx.New(
		fx.Provide(
			newLogger,
			config.Load,
			newDBPool,
			newRedis,
		),
		fx.Invoke(startIndexers),
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
	poolCfg.MaxConns = 5
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

// newChainClient creates a chain.Client for a single chain from env-configured addresses
func newChainClient(cfg *config.Config) (*chain.Client, error) {
	addrs := map[string]string{
		"AWPRegistry":  cfg.AWPRegistryAddress,
		"AWPToken":     cfg.AWPTokenAddress,
		"AWPEmission":  cfg.AWPEmissionAddress,
		"StakingVault": cfg.StakingVaultAddress,
		"SubnetNFT":    cfg.SubnetNFTAddress,
		"AWPDAO":       cfg.DAOAddress,
		"StakeNFT":     cfg.StakeNFTAddress,
	}
	dialCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return chain.NewClient(dialCtx, cfg.RPCURL, addrs)
}

// startIndexers starts one or more indexer goroutines depending on configuration.
// Multi-chain mode: CHAINS_FILE env set => spawn one goroutine per chain.
// Single-chain mode: fallback using CHAIN_ID + RPC_URL from env.
func startIndexers(lc fx.Lifecycle, pool *pgxpool.Pool, rdb *redis.Client, cfg *config.Config, logger *slog.Logger) error {
	if cfg.ChainsFile != "" {
		return startMultiChain(lc, pool, rdb, cfg, logger)
	}
	return startSingleChain(lc, pool, rdb, cfg, logger)
}

// startSingleChain uses env-configured CHAIN_ID and RPC_URL (backward compatible)
func startSingleChain(lc fx.Lifecycle, pool *pgxpool.Pool, rdb *redis.Client, cfg *config.Config, logger *slog.Logger) error {
	client, err := newChainClient(cfg)
	if err != nil {
		return fmt.Errorf("chain client: %w", err)
	}
	idx, err := chain.NewIndexer(client, pool, rdb, cfg.ChainID, cfg.DeployBlock)
	if err != nil {
		return fmt.Errorf("indexer: %w", err)
	}

	var cancel context.CancelFunc
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("indexer starting (single-chain)", "chainId", cfg.ChainID)
			runCtx, c := context.WithCancel(context.Background())
			cancel = c
			go func() {
				if err := idx.Run(runCtx); err != nil {
					logger.Error("indexer error", "chainId", cfg.ChainID, "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if cancel != nil {
				cancel()
			}
			return nil
		},
	})
	return nil
}

// startMultiChain loads chains.yaml and spawns one indexer goroutine per chain
func startMultiChain(lc fx.Lifecycle, pool *pgxpool.Pool, rdb *redis.Client, cfg *config.Config, logger *slog.Logger) error {
	chains, err := config.LoadChains(cfg.ChainsFile)
	if err != nil {
		return fmt.Errorf("load chains: %w", err)
	}

	type indexerEntry struct {
		idx     *chain.Indexer
		chainID int64
		name    string
	}

	// 创建所有chain client和indexer（在启动前完成，失败则阻止启动）
	entries := make([]indexerEntry, 0, len(chains))
	for _, ch := range chains {
		// 为每条链创建 chain client（per-chain 地址覆盖优先于全局 env）
		addrs := map[string]string{
			"AWPRegistry":  config.ResolveAddress(ch.AWPRegistry, cfg.AWPRegistryAddress),
			"AWPToken":     config.ResolveAddress(ch.AWPToken, cfg.AWPTokenAddress),
			"AWPEmission":  config.ResolveAddress(ch.AWPEmission, cfg.AWPEmissionAddress),
			"StakingVault": config.ResolveAddress(ch.StakingVault, cfg.StakingVaultAddress),
			"SubnetNFT":    config.ResolveAddress(ch.SubnetNFT, cfg.SubnetNFTAddress),
			"AWPDAO":       config.ResolveAddress(ch.DAOAddress, cfg.DAOAddress),
			"StakeNFT":     config.ResolveAddress(ch.StakeNFT, cfg.StakeNFTAddress),
		}
		dialCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, clientErr := chain.NewClient(dialCtx, ch.RPCURL, addrs)
		cancel()
		if clientErr != nil {
			return fmt.Errorf("chain client for %s (chainId=%d): %w", ch.Name, ch.ChainID, clientErr)
		}
		deployBlock := ch.DeployBlock
		if deployBlock == 0 {
			deployBlock = cfg.DeployBlock // fallback to global
		}
		idx, idxErr := chain.NewIndexer(client, pool, rdb, ch.ChainID, deployBlock)
		if idxErr != nil {
			return fmt.Errorf("indexer for %s (chainId=%d): %w", ch.Name, ch.ChainID, idxErr)
		}
		entries = append(entries, indexerEntry{idx: idx, chainID: ch.ChainID, name: ch.Name})
	}

	var cancelAll context.CancelFunc
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("indexer starting (multi-chain)", "chains", len(entries))
			runCtx, c := context.WithCancel(context.Background())
			cancelAll = c

			var wg sync.WaitGroup
			for _, e := range entries {
				wg.Add(1)
				go func(entry indexerEntry) {
					defer wg.Done()
					logger.Info("indexer goroutine started", "chain", entry.name, "chainId", entry.chainID)
					if err := entry.idx.Run(runCtx); err != nil {
						logger.Error("indexer error", "chain", entry.name, "chainId", entry.chainID, "error", err)
					}
				}(e)
			}
			// WaitGroup不在OnStart中等待——goroutine在后台运行直到OnStop
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if cancelAll != nil {
				cancelAll()
			}
			return nil
		},
	})
	return nil
}
