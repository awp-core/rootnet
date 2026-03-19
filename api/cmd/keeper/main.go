package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
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
			newEthClient,
			newKeeper,
		),
		fx.Invoke(startKeeper),
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

func newEthClient(lc fx.Lifecycle, cfg *config.Config) (*ethclient.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := ethclient.DialContext(ctx, cfg.RPCURL)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			client.Close()
			return nil
		},
	})
	return client, nil
}

func newKeeper(client *ethclient.Client, rdb *redis.Client, cfg *config.Config, logger *slog.Logger) (*chain.Keeper, error) {
	if cfg.KeeperPrivateKey == "" {
		return nil, fmt.Errorf("KEEPER_PRIVATE_KEY is required")
	}
	key, err := crypto.HexToECDSA(cfg.KeeperPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid keeper private key: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	awpEmissionAddr := common.HexToAddress(cfg.AWPEmissionAddress)
	awpTokenAddr := common.HexToAddress(cfg.AWPTokenAddress)
	return chain.NewKeeper(client, awpEmissionAddr, awpTokenAddr, key, chainID, rdb, logger)
}

func startKeeper(lc fx.Lifecycle, k *chain.Keeper, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Keeper starting")
			k.Start(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			k.Stop()
			return nil
		},
	})
}
