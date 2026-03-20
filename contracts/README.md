# AWP Registry — Smart Contracts

Solidity 0.8.24, Foundry, OpenZeppelin 5.x. Targets BSC and Base (EVM Cancun).

## Build

```bash
forge build
```

## Test

```bash
# All tests (284 pass, 2 BSC fork tests require RPC)
forge test

# Single file
forge test --match-contract AWPEmissionTest -vv

# With gas report
forge test --gas-report
```

## Deploy

```bash
# Configure contracts/.env (see docs/deployment-guide.md)
# Then run one-click deploy:
../scripts/deploy.sh

# Or manual:
forge script script/Deploy.s.sol --rpc-url $ETH_RPC_URL --broadcast
```

## Contracts (14)

| Contract | Description |
|----------|-------------|
| `AWPRegistry` | Unified entry: subnet lifecycle + staking allocation. Epoch authority (immutable genesisTime + epochDuration). |
| `AWPToken` | ERC20 + ERC1363 + Permit. 10B max, 200M (2%) pre-minted. |
| `AWPEmission` | **[DRAFT]** UUPS proxy. Oracle multi-sig weight submission, batch epoch settlement, exponential decay. |
| `AlphaToken` | Per-subnet ERC20 via CREATE2. Time-based mint cap. |
| `AlphaTokenFactory` | CREATE2 deployer with configurable EIP-55 vanity address rules. |
| `StakeNFT` | ERC721 position NFTs. Deposit AWP with lock period. Voting power = amount * sqrt(min(remaining, 54)). |
| `StakingVault` | Pure allocation logic. (user, agent, subnetId) triples. Auto-enumerates agent subnets. |
| `SubnetNFT` | ERC721. tokenId = subnetId. Ownership = subnet control. |
| `SubnetManager` | Default subnet contract (PancakeSwap V4). Merkle claim, AWP strategy, deployed behind ERC1967Proxy. |
| `SubnetManagerUni` | Uniswap V4 variant of SubnetManager for Base deployment. |
| `LPManager` | PancakeSwap V4 Concentrated Liquidity. Full-range, permanently locked. |
| `LPManagerUni` | Uniswap V4 variant of LPManager for Base deployment. |
| `AWPDAO` | Custom NFT-based voting. Executable proposals (Timelock) + signal proposals (vote-only). |
| `Treasury` | OZ TimelockController. |

## Dependencies

```bash
forge install  # OpenZeppelin, forge-std, infinity-periphery (PancakeSwap V4)
```
