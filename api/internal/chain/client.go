package chain

import (
	"context"
	"fmt"

	"github.com/cortexia/rootnet/api/internal/chain/bindings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client wraps ethclient.Client and provides convenient access to contract binding instances
type Client struct {
	Eth *ethclient.Client

	// Contract addresses
	AWPRegistryAddr  common.Address
	AWPTokenAddr     common.Address
	AWPEmissionAddr  common.Address
	StakingVaultAddr common.Address
	SubnetNFTAddr    common.Address
	AlphaTokenAddr   common.Address
	AWPDAOAddr       common.Address
	StakeNFTAddr     common.Address

	// Contract binding instances
	AWPRegistry  *bindings.AWPRegistry
	AWPToken     *bindings.AWPToken
	AWPEmission  *bindings.AWPEmission
	StakingVault *bindings.StakingVault
	SubnetNFT    *bindings.SubnetNFT
	AlphaToken   *bindings.AlphaToken
	AWPDAO       *bindings.AWPDAO
	StakeNFT     *bindings.StakeNFT
}

// NewClient creates a chain client, connects to RPC, and initializes all contract binding instances.
// The addresses map uses contract names (e.g. "AWPRegistry", "AWPToken") as keys and hex addresses as values.
func NewClient(ctx context.Context, rpcURL string, addresses map[string]string) (*Client, error) {
	eth, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	c := &Client{Eth: eth}

	// Parse contract addresses
	parseAddr := func(name string) (common.Address, error) {
		raw, ok := addresses[name]
		if !ok {
			return common.Address{}, fmt.Errorf("missing contract address: %s", name)
		}
		if !common.IsHexAddress(raw) {
			return common.Address{}, fmt.Errorf("invalid contract address %s: %s", name, raw)
		}
		return common.HexToAddress(raw), nil
	}

	// AWPRegistry (required)
	c.AWPRegistryAddr, err = parseAddr("AWPRegistry")
	if err != nil {
		return nil, err
	}
	c.AWPRegistry, err = bindings.NewAWPRegistry(c.AWPRegistryAddr, eth)
	if err != nil {
		return nil, fmt.Errorf("failed to bind AWPRegistry: %w", err)
	}

	// AWPToken
	c.AWPTokenAddr, err = parseAddr("AWPToken")
	if err != nil {
		return nil, err
	}
	c.AWPToken, err = bindings.NewAWPToken(c.AWPTokenAddr, eth)
	if err != nil {
		return nil, fmt.Errorf("failed to bind AWPToken: %w", err)
	}

	// AWPEmission
	c.AWPEmissionAddr, err = parseAddr("AWPEmission")
	if err != nil {
		return nil, err
	}
	c.AWPEmission, err = bindings.NewAWPEmission(c.AWPEmissionAddr, eth)
	if err != nil {
		return nil, fmt.Errorf("failed to bind AWPEmission: %w", err)
	}

	// StakingVault
	c.StakingVaultAddr, err = parseAddr("StakingVault")
	if err != nil {
		return nil, err
	}
	c.StakingVault, err = bindings.NewStakingVault(c.StakingVaultAddr, eth)
	if err != nil {
		return nil, fmt.Errorf("failed to bind StakingVault: %w", err)
	}

	// SubnetNFT
	c.SubnetNFTAddr, err = parseAddr("SubnetNFT")
	if err != nil {
		return nil, err
	}
	c.SubnetNFT, err = bindings.NewSubnetNFT(c.SubnetNFTAddr, eth)
	if err != nil {
		return nil, fmt.Errorf("failed to bind SubnetNFT: %w", err)
	}

	// AlphaToken (optional)
	if raw, ok := addresses["AlphaToken"]; ok && common.IsHexAddress(raw) {
		c.AlphaTokenAddr = common.HexToAddress(raw)
		c.AlphaToken, err = bindings.NewAlphaToken(c.AlphaTokenAddr, eth)
		if err != nil {
			return nil, fmt.Errorf("failed to bind AlphaToken: %w", err)
		}
	}

	// AWPDAO
	c.AWPDAOAddr, err = parseAddr("AWPDAO")
	if err != nil {
		return nil, err
	}
	c.AWPDAO, err = bindings.NewAWPDAO(c.AWPDAOAddr, eth)
	if err != nil {
		return nil, fmt.Errorf("failed to bind AWPDAO: %w", err)
	}

	// StakeNFT
	c.StakeNFTAddr, err = parseAddr("StakeNFT")
	if err != nil {
		return nil, err
	}
	c.StakeNFT, err = bindings.NewStakeNFT(c.StakeNFTAddr, eth)
	if err != nil {
		return nil, fmt.Errorf("failed to bind StakeNFT: %w", err)
	}

	return c, nil
}

// BlockNumber returns the latest block number
func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	return c.Eth.BlockNumber(ctx)
}

// Close closes the underlying RPC connection
func (c *Client) Close() {
	c.Eth.Close()
}
