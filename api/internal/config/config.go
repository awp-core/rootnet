package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config holds application configuration parsed from environment variables
type Config struct {
	// Database
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/awp?sslmode=disable"`

	// Redis
	RedisURL string `env:"REDIS_URL" envDefault:"redis://localhost:6379/0"`

	// HTTP server
	HTTPAddr    string `env:"HTTP_ADDR" envDefault:":8080"`
	TrustProxy  bool   `env:"TRUST_PROXY" envDefault:"false"` // Set to true ONLY when behind a trusted reverse proxy (nginx/caddy)

	// Multi-chain
	ChainsFile string `env:"CHAINS_FILE" envDefault:""` // Path to chains.yaml; empty = single-chain mode

	// Chain
	ChainID         int64  `env:"CHAIN_ID" envDefault:"0"`
	RPCURL          string `env:"RPC_URL" envDefault:""`
	AWPRegistryAddress string `env:"AWP_REGISTRY_ADDRESS"`
	AWPWorkNetAddress string `env:"AWP_WORKNET_ADDRESS"`
	DAOAddress      string `env:"DAO_ADDRESS"`
	AWPTokenAddress string `env:"AWP_TOKEN_ADDRESS"`

	// Keeper
	KeeperPrivateKey string `env:"KEEPER_PRIVATE_KEY"`
	KeeperSkipSettle bool   `env:"KEEPER_SKIP_SETTLE" envDefault:"false"`

	// Relayer (gasless transaction relay)
	RelayerPrivateKey string `env:"RELAYER_PRIVATE_KEY"`
	// Rate limits are configured via Redis (HSET ratelimit:config relay "100:3600")
	// Defaults are compiled into the ratelimit.Limiter package

	// Contract address registry (protocol contracts)
	AWPAllocatorAddress        string `env:"AWP_ALLOCATOR_ADDRESS"`
	AWPEmissionAddress         string `env:"AWP_EMISSION_ADDRESS"`
	TreasuryAddress            string `env:"TREASURY_ADDRESS"`
	VeAWPAddress               string `env:"VEAWP_ADDRESS"`
	LPManagerAddress           string `env:"LP_MANAGER_ADDRESS"`
	PoolManagerAddress         string `env:"POOL_MANAGER" envDefault:""` // DEX V4 PoolManager (for worknet token price reads)
	WorknetTokenFactoryAddress string `env:"WORKNET_TOKEN_FACTORY_ADDRESS"`

	// Vanity address mining config
	WorknetTokenBytecodeHash string `env:"WORKNET_TOKEN_BYTECODE_HASH"` // keccak256(WorknetToken.creationCode), hex
	VanityRule               string `env:"VANITY_RULE"`                 // WorknetTokenFactory.vanityRule(), uint64 hex (e.g. "0x1001FFFF0C0A0F0E")

	// Indexer start block (deploy block); used only on first run when sync_states is empty
	DeployBlock int64 `env:"DEPLOY_BLOCK" envDefault:"0"`

	// Admin API
	AdminToken string `env:"ADMIN_TOKEN" envDefault:""`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	// Validate: either multi-chain (ChainsFile) or single-chain (ChainID+RPC_URL)
	if cfg.ChainsFile == "" && (cfg.ChainID == 0 || cfg.RPCURL == "") {
		return nil, fmt.Errorf("either CHAINS_FILE or both CHAIN_ID and RPC_URL must be set")
	}
	return cfg, nil
}
