package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
			newRedis,
		),
		fx.Invoke(startKeepers),
	).Run()
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
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

// newKeeperForChain creates a Keeper for a single chain using the given RPC URL and chain ID
func newKeeperForChain(rpcURL string, chainID *big.Int, rdb *redis.Client, cfg *config.Config, logger *slog.Logger) (*chain.Keeper, error) {
	if cfg.KeeperPrivateKey == "" {
		return nil, fmt.Errorf("KEEPER_PRIVATE_KEY is required")
	}
	key, err := crypto.HexToECDSA(cfg.KeeperPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid keeper private key: %w", err)
	}

	dialCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := ethclient.DialContext(dialCtx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial RPC: %w", err)
	}

	awpEmissionAddr := common.HexToAddress(cfg.AWPEmissionAddress)
	awpTokenAddr := common.HexToAddress(cfg.AWPTokenAddress)
	return chain.NewKeeper(client, awpEmissionAddr, awpTokenAddr, key, chainID, rdb, logger)
}

// startKeepers starts one or more keeper goroutines depending on configuration.
// Multi-chain mode: CHAINS_FILE env set => spawn one goroutine per chain.
// Single-chain mode: fallback using CHAIN_ID + RPC_URL from env.
func startKeepers(lc fx.Lifecycle, rdb *redis.Client, cfg *config.Config, logger *slog.Logger) error {
	if cfg.ChainsFile != "" {
		return startMultiChain(lc, rdb, cfg, logger)
	}
	return startSingleChain(lc, rdb, cfg, logger)
}

// startSingleChain uses env-configured CHAIN_ID and RPC_URL (backward compatible)
func startSingleChain(lc fx.Lifecycle, rdb *redis.Client, cfg *config.Config, logger *slog.Logger) error {
	// 获取链上chain ID
	dialCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := ethclient.DialContext(dialCtx, cfg.RPCURL)
	if err != nil {
		return fmt.Errorf("failed to dial RPC: %w", err)
	}

	chainCtx, chainCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer chainCancel()
	chainID, err := client.ChainID(chainCtx)
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to get chain ID: %w", err)
	}
	client.Close() // newKeeperForChain会创建自己的连接

	k, err := newKeeperForChain(cfg.RPCURL, chainID, rdb, cfg, logger)
	if err != nil {
		return fmt.Errorf("keeper: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Keeper starting (single-chain)", "chainId", chainID.Int64())
			k.Start(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			k.Stop()
			return nil
		},
	})
	return nil
}

// startMultiChain loads chains.yaml and spawns one keeper goroutine per chain
func startMultiChain(lc fx.Lifecycle, rdb *redis.Client, cfg *config.Config, logger *slog.Logger) error {
	chains, err := config.LoadChains(cfg.ChainsFile)
	if err != nil {
		return fmt.Errorf("load chains: %w", err)
	}

	type keeperEntry struct {
		keeper  *chain.Keeper
		chainID int64
		name    string
	}

	// 创建所有keeper（在启动前完成，失败则阻止启动）
	entries := make([]keeperEntry, 0, len(chains))
	for _, ch := range chains {
		k, kerr := newKeeperForChain(ch.RPCURL, big.NewInt(ch.ChainID), rdb, cfg, logger)
		if kerr != nil {
			return fmt.Errorf("keeper for %s (chainId=%d): %w", ch.Name, ch.ChainID, kerr)
		}
		entries = append(entries, keeperEntry{keeper: k, chainID: ch.ChainID, name: ch.Name})
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Keeper starting (multi-chain)", "chains", len(entries))

			var wg sync.WaitGroup
			for _, e := range entries {
				wg.Add(1)
				go func(entry keeperEntry) {
					defer wg.Done()
					logger.Info("keeper goroutine started", "chain", entry.name, "chainId", entry.chainID)
					entry.keeper.Start(ctx)
				}(e)
			}
			// WaitGroup不在OnStart中等待——goroutine在后台运行直到OnStop
			return nil
		},
		OnStop: func(ctx context.Context) error {
			for _, e := range entries {
				e.keeper.Stop()
			}
			return nil
		},
	})
	return nil
}
