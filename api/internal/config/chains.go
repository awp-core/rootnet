package config

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

// ChainConfig holds per-chain deployment configuration loaded from chains.yaml
type ChainConfig struct {
	ChainID         int64  `yaml:"chainId" json:"chainId"`
	Name            string `yaml:"name" json:"name"`
	RPCURL          string `yaml:"rpcUrl" json:"-"` // never expose RPC URL in API
	DEX             string `yaml:"dex" json:"dex"`
	InitialMint     int64  `yaml:"initialMint" json:"-"`
	DeployBlock     int64  `yaml:"deployBlock" json:"-"`
	Explorer        string `yaml:"explorer" json:"explorer"`
	PoolManager     string `yaml:"poolManager" json:"-"`
	PositionManager string `yaml:"positionManager" json:"-"`
	Permit2         string `yaml:"permit2" json:"-"`
	SwapRouter      string `yaml:"swapRouter" json:"-"`
	StateView       string `yaml:"stateView" json:"-"`

	// 可选的 per-chain 合约地址覆盖（为空时使用全局 env 配置）
	AWPRegistry  string `yaml:"awpRegistry" json:"-"`
	AWPToken     string `yaml:"awpToken" json:"-"`
	AWPEmission  string `yaml:"awpEmission" json:"-"`
	StakingVault string `yaml:"stakingVault" json:"-"`
	StakeNFT     string `yaml:"stakeNFT" json:"-"`
	SubnetNFT    string `yaml:"subnetNFT" json:"-"`
	DAOAddress   string `yaml:"dao" json:"-"`
	LPManager    string `yaml:"lpManager" json:"-"`
}

// ResolveAddress 返回 per-chain 地址覆盖，如果为空则返回全局默认值
func ResolveAddress(perChain, global string) string {
	if perChain != "" {
		return perChain
	}
	return global
}

// ChainsFile is the top-level structure of chains.yaml
type ChainsFile struct {
	Chains map[string]ChainConfig `yaml:"chains"`
}

// LoadChains reads chains.yaml and resolves env vars in rpcUrl fields.
// Returns nil if path is empty (single-chain mode).
func LoadChains(path string) ([]ChainConfig, error) {
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read chains config: %w", err)
	}
	var file ChainsFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse chains config: %w", err)
	}
	chains := make([]ChainConfig, 0, len(file.Chains))
	for _, cfg := range file.Chains {
		cfg.RPCURL = os.ExpandEnv(cfg.RPCURL)
		if cfg.RPCURL == "" || cfg.ChainID == 0 {
			continue // skip misconfigured chains
		}
		chains = append(chains, cfg)
	}
	if len(chains) == 0 {
		return nil, fmt.Errorf("no valid chains in %s", path)
	}
	sort.Slice(chains, func(i, j int) bool {
		return chains[i].ChainID < chains[j].ChainID
	})
	return chains, nil
}
