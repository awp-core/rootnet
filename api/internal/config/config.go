package config

import (
	"github.com/caarlos0/env/v11"
)

// Config holds application configuration parsed from environment variables
type Config struct {
	// Database
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/awp?sslmode=disable"`

	// Redis
	RedisURL string `env:"REDIS_URL" envDefault:"redis://localhost:6379/0"`

	// HTTP server
	HTTPAddr string `env:"HTTP_ADDR" envDefault:":8080"`

	// Chain
	ChainID         int64  `env:"CHAIN_ID" envDefault:"56"` // BSC Mainnet
	RPCURL          string `env:"RPC_URL" envDefault:"https://bsc-testnet-rpc.publicnode.com"`
	RootNetAddress  string `env:"ROOTNET_ADDRESS"`
	SubnetNFTAddress string `env:"SUBNETNFT_ADDRESS"`
	DAOAddress      string `env:"DAO_ADDRESS"`
	AWPTokenAddress string `env:"AWP_TOKEN_ADDRESS"`

	// Keeper
	KeeperPrivateKey string `env:"KEEPER_PRIVATE_KEY"`

	// Relayer (gasless transaction relay)
	RelayerPrivateKey string `env:"RELAYER_PRIVATE_KEY"`
	// Rate limits are configured via Redis (HSET ratelimit:config relay "100:3600")
	// Defaults are compiled into the ratelimit.Limiter package

	// Contract address registry (all 11 protocol contracts)
	StakingVaultAddress  string `env:"STAKING_VAULT_ADDRESS"`
	AWPEmissionAddress   string `env:"AWP_EMISSION_ADDRESS"`
	TreasuryAddress      string `env:"TREASURY_ADDRESS"`
	StakeNFTAddress      string `env:"STAKE_NFT_ADDRESS"`
	AccessManagerAddress string `env:"ACCESS_MANAGER_ADDRESS"`
	LPManagerAddress     string `env:"LP_MANAGER_ADDRESS"`
	AlphaFactoryAddress  string `env:"ALPHA_FACTORY_ADDRESS"`

	// Vanity address mining config
	AlphaInitCodeHash string `env:"ALPHA_INITCODE_HASH"` // keccak256(AlphaToken.creationCode), hex
	VanityRule        string `env:"VANITY_RULE"`          // AlphaTokenFactory.vanityRule(), uint64 hex (e.g. "0x1001FFFF0C0A0F0E")

	// Indexer start block (deploy block); used only on first run when sync_states is empty
	DeployBlock int64 `env:"DEPLOY_BLOCK" envDefault:"0"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
