package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
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
			newChainClient,
			newIndexer,
		),
		fx.Invoke(startIndexer),
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

func newIndexer(client *chain.Client, pool *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) (*chain.Indexer, error) {
	return chain.NewIndexer(client, pool, rdb, cfg.DeployBlock)
}

func startIndexer(lc fx.Lifecycle, idx *chain.Indexer, logger *slog.Logger) {
	var idxCancel context.CancelFunc

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Indexer starting")
			runCtx, cancel := context.WithCancel(context.Background())
			idxCancel = cancel
			go func() {
				if err := idx.Run(runCtx); err != nil {
					logger.Error("Indexer error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if idxCancel != nil {
				idxCancel()
			}
			return nil
		},
	})
}
