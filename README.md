# AWP

[![BSC Mainnet](https://img.shields.io/badge/BSC-Mainnet-yellow)](https://bscscan.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-green)](LICENSE)

> **Testnet.** AWP is currently in testnet on BSC mainnet. AWP mainnet deployment (BSC + Base) is planned. Protocol parameters may change before the official mainnet launch.

## Abstract

AWP is a decentralized **Agent Working** protocol deployed on BNB Smart Chain (BSC). The protocol establishes a permissionless marketplace where autonomous AI agent networks (*subnets*) compete for protocol-level emission rewards through stake-weighted oracle consensus. Each subnet deploys an independent economy backed by a dedicated ERC-20 token (Alpha), with initial liquidity bootstrapped via PancakeSwap V4 Concentrated Liquidity at registration time.

The system introduces a **Principal–Agent staking model**: Principals deposit AWP tokens into non-fungible position NFTs with time-locked commitments, then allocate stake across (agent, subnet) triples. An exponentially-decaying emission schedule distributes newly-minted AWP to subnet managers proportional to oracle-assigned governance weights, with a 50/50 split between subnet recipients and a DAO treasury governed by NFT-weighted quadratic voting.

Key design contributions include: (1) a gasless relay layer enabling device-bound agents to participate without holding native gas tokens; (2) an ERC-1363 `mintAndCall` emission pathway that auto-triggers configurable AWP handling strategies (reserve, single-sided liquidity provision, or buyback-and-burn) at the subnet manager level; (3) a tiered CREATE2 vanity address system with pre-mined salt pools for deterministic cross-chain deployment; and (4) a modular subnet architecture where a default `SubnetManager` proxy contract provides Merkle-based reward distribution, multi-role access control, and PancakeSwap V4 integration out of the box, while advanced operators may deploy custom manager contracts.

The protocol consists of 13 Solidity contracts (Foundry, Solidity 0.8.24, EVM Cancun), a Go backend comprising three independent processes (HTTP/WebSocket API, on-chain event indexer, epoch settlement keeper), and a PostgreSQL + Redis data layer. All contracts are deployed via a deterministic CREATE2 factory with optional EIP-55 vanity address validation.

> **Note:** The AWP Emission mechanism (AWPEmission contract, oracle consensus, epoch settlement) is under active design and has not been finalized. All emission-related descriptions in this document are preliminary and subject to change.

## Architecture

```
User
 ├── AWPRegistry ─── register / join / allocate / subnet lifecycle
 │    ├── StakeNFT ── ERC721 position NFTs (deposit AWP + lock)
 │    ├── StakingVault ── allocation bookkeeping (auto-enumerates agent subnets)
 │    ├── AccessManager ── Principal/Agent identity + delegation
 │    ├── SubnetNFT ── subnet ownership NFTs
 │    └── LPManager ── PancakeSwap V4 CL liquidity
 │
 ├── AWPEmission (UUPS proxy) ── epoch settlement + AWP minting
 │    └── Oracle multi-sig → submitAllocations → settleEpoch
 │
 ├── AWPDAO ── NFT-based voting (executable + signal proposals)
 │    └── Treasury (TimelockController) ── governance execution
 │
 └── AlphaTokenFactory ── CREATE2 per-subnet tokens with vanity addresses
```

**13 contracts**, 3 Go backend processes (API / Indexer / Keeper), PostgreSQL, Redis.

**Backend API** provides:
- Read-only REST API + WebSocket real-time events
- Gasless relay endpoints (`/api/relay/register`, `/api/relay/bind`) — EIP-712 signed, relayer pays gas
- Vanity salt mining (`/api/vanity/compute-salt`) — uses Foundry `cast create2` for high-speed parallel mining

## Key Design

- **Principal/Agent**: Principals register, Agents join via `bind(principal)`. Supports rebind (auto-freezes old allocations), unbind, gasless `bindFor` with EIP-712 signatures.
- **Staking**: deposit AWP into StakeNFT (ERC721 positions with lock period). Allocate to (agent, subnet) triples via StakingVault. Auto-enumeration of agent subnets — no caller-supplied subnet list needed for freeze.
- **Epoch**: time-based on AWPEmission (`(block.timestamp - genesisTime) / epochDuration`, 1 day).
- **Emission**: exponential decay. 50% to subnets, 50% to DAO. Batch settlement via `settleEpoch(limit)`.
- **Voting**: quadratic voting with time-weighted staking positions. Two proposal types: executable (Timelock) and signal (vote-only).
- **Subnets**: registration deploys Alpha token (CREATE2 vanity address) + PancakeSwap V4 LP. Time-based mint cap on Alpha.

## Deployment

### Prerequisites

- [Foundry](https://getfoundry.sh/) (`forge`, `cast`)
- Go 1.26+
- PostgreSQL 15+
- Redis 7+
- `jq`, `psql`

### Step 1: Configure

```bash
# Copy and fill the deployment config
cp contracts/.env.example contracts/.env
vim contracts/.env
```

Required fields in `contracts/.env`:
```
ETH_RPC_URL=...           # BSC RPC endpoint
DEPLOYER_PRIVATE_KEY=...  # Deployer wallet private key
GUARDIAN=...              # Emergency pause guardian address
LIQUIDITY_POOL=...       # LP wallet
AIRDROP=...              # Airdrop wallet
POOL_MANAGER=...         # PancakeSwap V4 CLPoolManager
POSITION_MANAGER=...     # PancakeSwap V4 CLPositionManager
PERMIT2=...              # Permit2 address
VANITY_RULE=0            # Alpha token vanity rule (0 = disabled)
```

### Step 2: Deploy Contracts

```bash
# Dry-run first (no transactions)
./scripts/deploy.sh --dry-run

# Full deployment → generates api/.env with all contract addresses
./scripts/deploy.sh
```

This deploys all 13 contracts via deterministic CREATE2, generates `api/.env` with:
- All contract addresses
- Deploy block (indexer start)
- AlphaToken initCodeHash + vanity rule (for mining API)

### Step 3: Verify Contracts (optional)

```bash
# Set BscScan API key in contracts/.env
echo 'ETHERSCAN_API_KEY=your_key' >> contracts/.env

# Verify all contracts on BscScan
./scripts/deploy.sh --verify-only
```

### Step 4: Deploy API Services

```bash
# Fill private keys for keeper/relayer
cp api/.env.local.example api/.env.local
vim api/.env.local  # KEEPER_PRIVATE_KEY, RELAYER_PRIVATE_KEY

# One-click: build + init DB + start all services
./scripts/deploy-api.sh
```

This builds 3 Go binaries, creates PostgreSQL tables, and starts:
- **api** — HTTP + WebSocket server (REST API, relay, vanity)
- **indexer** — On-chain event scanner → PostgreSQL + Redis Pub/Sub
- **keeper** — Epoch settlement + cache updates (every 30s)

### Managing Services

```bash
./scripts/deploy-api.sh --status     # Check service status + health
./scripts/deploy-api.sh --stop       # Stop all services
./scripts/deploy-api.sh --start      # Restart services
./scripts/deploy-api.sh --build-only # Rebuild binaries only
./scripts/deploy-api.sh --db-only    # Reinitialize database only
```

### Configuration Files

| File | Purpose | Git tracked |
|------|---------|:-----------:|
| `contracts/.env` | Deployment config (keys, addresses, salts) | No |
| `contracts/.env.example` | Template for contracts/.env | Yes |
| `api/.env` | Auto-generated by deploy.sh (contract addresses) | No |
| `api/.env.local` | Manual overrides (private keys) | No |
| `api/.env.local.example` | Template for api/.env.local | Yes |

## Development

```bash
# Build contracts
cd contracts && forge build

# Run tests (282 pass, 2 require BSC fork RPC)
forge test

# Build Go backend
cd api && go build ./...

# Run Go tests
go test ./...

# Regenerate Go bindings after contract changes
jq '.abi' out/AWPRegistry.sol/AWPRegistry.json > /tmp/AWPRegistry.abi
jq '.bytecode.object' out/AWPRegistry.sol/AWPRegistry.json | tr -d '"' > /tmp/AWPRegistry.bin
abigen --abi /tmp/AWPRegistry.abi --bin /tmp/AWPRegistry.bin --pkg bindings --type AWPRegistry --out api/internal/chain/bindings/awp_registry.go
```

## Project Structure

```
contracts/
  src/
    AWPRegistry.sol             # Unified entry point
    token/
      AWPToken.sol              # ERC20+Votes, 10B supply
      AWPEmission.sol           # UUPS proxy emission engine
      AlphaToken.sol            # Per-subnet ERC20 (CREATE2)
      AlphaTokenFactory.sol     # CREATE2 deployer + vanity rules
    core/
      StakeNFT.sol              # ERC721 staking positions
      StakingVault.sol          # Allocation bookkeeping
      AccessManager.sol         # Principal/Agent registration
      SubnetNFT.sol             # Subnet ownership
      LPManager.sol             # PancakeSwap V4 CL
    governance/
      AWPDAO.sol                # NFT-based voting
      Treasury.sol              # TimelockController
    subnets/
      SubnetManager.sol         # Default subnet contract (proxy)
  test/                         # 284 tests
  script/                       # Deploy.s.sol, Predict.s.sol, InitCodeHashes.s.sol

api/
  cmd/api/                      # HTTP + WebSocket server
  cmd/indexer/                  # On-chain event scanner
  cmd/keeper/                   # Epoch settlement + cache updates
  internal/
    chain/                      # Contract bindings, indexer, keeper, relayer, vanity
    db/                         # PostgreSQL schema + sqlc queries
    server/                     # Chi router + handlers + WebSocket hub
    config/                     # Environment config

docs/
  architecture.md               # Full technical architecture
  api-reference.md              # Contract + REST + WebSocket API reference
  deployment-guide.md           # Deploy + operations guide
  subnet-developer-guide.md     # For subnet developers

skills-dev/
  contract-api.md               # Contract API quick reference
  rest-api.md                   # REST API reference
  examples.md                   # Code examples (viem)
  agent-skill-guide.md          # Agent skill discovery + install
  abi/                          # Contract ABI JSON files
  config.md                     # Constants + addresses + env vars

scripts/
  deploy.sh                     # Contract deployment + verification + api/.env generation
  deploy-api.sh                 # API service build + DB init + process management
```

## API Endpoints

| Group | Endpoints | Description |
|-------|-----------|-------------|
| System | `GET /api/health`, `/api/registry` | Health check, all 11 contract addresses (excludes implementation contracts) |
| Users | `GET /api/users/*` | User list, detail, registration check |
| Agents | `GET /api/agents/*`, `POST /api/agents/batch-info` | Agent lookup, batch query |
| Staking | `GET /api/staking/*` | Balances, positions, allocations, subnet totals |
| Subnets | `GET /api/subnets/*` | Subnet list/detail/skills/earnings (includes `subnet_contract` + `alpha_token`) |
| Emission | `GET /api/emission/*` | Current epoch, schedule, history |
| Tokens | `GET /api/tokens/*` | AWP info, Alpha token info/price |
| Governance | `GET /api/governance/*` | Proposals, treasury |
| Relay | `POST /api/relay/register`, `/bind` | Gasless EIP-712 transactions |
| Vanity | `POST /api/vanity/compute-salt` | CREATE2 salt mining (factory vanity rule) |
| WebSocket | `WS /ws/live` | Real-time on-chain events |

## Stack

| Layer | Technology |
|-------|-----------|
| Contracts | Solidity 0.8.24, Foundry, OpenZeppelin 5.x |
| Backend | Go 1.26, Chi v5, sqlc + pgx/v5, PostgreSQL, Redis |
| Frontend | Next.js 14, Tailwind, wagmi/viem |
| Chain | BSC (EVM Cancun), PancakeSwap V4 |

## Documentation

- [Architecture](docs/architecture.md) — Full technical design
- [API Reference](docs/api-reference.md) — Contract + REST + WebSocket
- [Deployment Guide](docs/deployment-guide.md) — Deploy + operations
- [Subnet Developer Guide](docs/subnet-developer-guide.md) — For subnet builders
- [Agent Skill Guide](skills-dev/agent-skill-guide.md) — Skill discovery + install

## License

[MIT](LICENSE)
