# AWP Registry — Smart Contracts

Solidity 0.8.24, Foundry, OpenZeppelin 5.x. Targets Base, Ethereum, Arbitrum, BSC (EVM Cancun).

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

## Contracts (11 deployed per chain)

> The protocol deploys 11 contracts per chain. `WorknetManager` / `WorknetManagerUni` and `LPManager` / `LPManagerUni` are chain-specific variants (PancakeSwap V4 for BSC, Uniswap V4 for Base/Ethereum/Arbitrum) — only one variant is deployed per chain.

| Contract | Description |
|----------|-------------|
| `AWPRegistry` | Unified entry: worknet lifecycle + account system (UUPS proxy). onlyGuardian admin. |
| `AWPToken` | ERC20 + ERC1363 + Permit. 10B max, configurable initial mint per chain. |
| `AWPEmission` | UUPS proxy. Guardian-only weight submission, batch epoch settlement, exponential decay. |
| `WorknetToken` | Per-worknet ERC20 via CREATE2. No constructor args (callback pattern). Time-based mint cap. |
| `WorknetTokenFactory` | CREATE2 deployer with configurable EIP-55 vanity address rules. Universal salt. |
| `veAWP` | ERC721 position NFTs. Deposit AWP with lock period. Voting power = amount * sqrt(min(remaining, 54w)). |
| `AWPAllocator` | UUPS proxy + EIP-712. Allocation bookkeeping. (staker, agent, worknetId) triples. |
| `AWPWorkNet` | ERC721. tokenId = worknetId. On-chain identity storage. Ownership = worknet control. |
| `WorknetManager` | Default worknet contract (PancakeSwap V4). Merkle claim, AWP strategy, deployed behind ERC1967Proxy. |
| `WorknetManagerUni` | Uniswap V4 variant of WorknetManager for Base/Ethereum/Arbitrum. |
| `LPManager` | PancakeSwap V4 Concentrated Liquidity. Full-range, permanently locked. |
| `LPManagerUni` | Uniswap V4 variant with stateView for reading pool slot0. |
| `AWPDAO` | veAWP-based voting. Executable proposals (Timelock) + signal proposals (vote-only). |
| `Treasury` | OZ TimelockController. |

## Dependencies

```bash
forge install  # OpenZeppelin, forge-std, infinity-periphery (PancakeSwap V4)
```
