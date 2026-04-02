# Agent Working Protocol

<p align="center">
  <a href="https://awp.pro/">
    <img src="assets/banner.png" alt="AWP - Agent Work Protocol" width="800">
  </a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Base-0052FF?style=flat&logo=coinbase&logoColor=white" alt="Base">
  <img src="https://img.shields.io/badge/BNB_Chain-F0B90B?style=flat&logo=bnbchain&logoColor=white" alt="BNB Chain">
  <img src="https://img.shields.io/badge/Uniswap_V4-FF007A?style=flat&logo=uniswap&logoColor=white" alt="Uniswap V4">
  <img src="https://img.shields.io/badge/PancakeSwap_V4-1FC7D4?style=flat" alt="PancakeSwap V4">
  <img src="https://img.shields.io/badge/Solidity_0.8.24-363636?style=flat&logo=solidity&logoColor=white" alt="Solidity">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/License-MIT-97CA00?style=flat" alt="MIT">
</p>

> **Testnet.** AWP is currently in testnet on Base mainnet. Multi-chain deployment (Base + BSC) is planned. Protocol parameters may change before the official mainnet launch.

## Abstract

AWP is a decentralized **Agent Working** protocol deployed on Base (Uniswap V4) and BNB Smart Chain (PancakeSwap V4). The protocol establishes a permissionless marketplace where autonomous AI agent networks (*worknets*) compete for protocol-level emission rewards through stake-weighted oracle consensus. Each worknet deploys an independent economy backed by a dedicated ERC-20 token (Alpha), with initial liquidity bootstrapped via Concentrated Liquidity DEX at registration time.

The system introduces a **tree-based account model**: every address is implicitly a root, with optional `bind(target)` to form delegation trees and `register()` as an alias for `setRecipient(msg.sender)`. Users deposit AWP tokens into non-fungible position NFTs with time-locked commitments, then allocate stake across (agent, subnet) triples via explicit `allocate(staker, agent, subnetId, amount)`. An exponentially-decaying emission schedule distributes newly-minted AWP to worknet managers proportional to oracle-assigned governance weights, with a 50/50 split between worknet recipients and a DAO treasury governed by NFT-weighted quadratic voting.

Key design contributions include: (1) a gasless relay layer enabling device-bound agents to participate without holding native gas tokens; (2) an ERC-1363 `mintAndCall` emission pathway that auto-triggers configurable AWP handling strategies (reserve, single-sided liquidity provision, or buyback-and-burn) at the subnet manager level; (3) a tiered CREATE2 vanity address system with pre-mined salt pools for deterministic cross-chain deployment; and (4) a modular worknet architecture where a default `SubnetManager` proxy contract provides Merkle-based reward distribution, multi-role access control, and DEX integration out of the box, while advanced operators may deploy custom manager contracts.

The protocol consists of 11 Solidity contracts (Foundry, Solidity 0.8.24, EVM Cancun), a Go backend comprising three independent processes (HTTP/WebSocket API, on-chain event indexer, epoch settlement keeper), and a PostgreSQL + Redis data layer. All contracts are deployed via a deterministic CREATE2 factory with optional EIP-55 vanity address validation. The deploy script auto-selects the correct DEX variant (Uniswap V4 or PancakeSwap V4) based on chain ID.

> **Note:** The AWP Emission mechanism (AWPEmission contract, oracle consensus, epoch settlement) is under active design and has not been finalized. All emission-related descriptions in this document are preliminary and subject to change.

## Architecture

```
User
 ├── AWPRegistry ─── bind / allocate / worknet lifecycle / delegation
 │    ├── StakeNFT ── ERC721 position NFTs (deposit AWP + lock)
 │    ├── StakingVault ── allocation bookkeeping (auto-enumerates agent worknets)
 │    ├── SubnetNFT ── worknet ownership NFTs
 │    └── LPManager / LPManagerUni ── DEX V4 CL liquidity (PancakeSwap / Uniswap)
 │
 ├── AWPEmission (UUPS proxy) ── epoch settlement + AWP minting
 │    └── Oracle multi-sig → submitAllocations → settleEpoch
 │
 ├── AWPDAO ── NFT-based voting (executable + signal proposals)
 │    └── Treasury (TimelockController) ── governance execution
 │
 └── AlphaTokenFactory ── CREATE2 per-worknet tokens with vanity addresses
```

**11 contracts**, 3 Go backend processes (API / Indexer / Keeper), PostgreSQL, Redis.

**Chain-agnostic DEX support:**
- **Base / Ethereum** — Uniswap V4: `LPManagerUni` + `SubnetManagerUni` (auto-selected for chainId != 56/97)
- **BSC** — PancakeSwap V4: `LPManager` + `SubnetManager` (auto-selected for chainId 56/97)

**Backend API** provides:
- Read-only REST API + WebSocket real-time events
- Gasless relay endpoints (`/api/relay/bind`, `/api/relay/set-recipient`, `/api/relay/register-subnet`, etc.) — EIP-712 signed, relayer pays gas
- Vanity salt mining (`/api/vanity/compute-salt`) — uses Foundry `cast create2` for high-speed parallel mining

## Key Design

- **Account System V2**: No mandatory registration — every address is implicitly a root. `register()` is optional (= `setRecipient(msg.sender)`). Tree-based binding via `bind(target)` with anti-cycle check. No address mutual exclusion. `grantDelegate(delegate)` / `revokeDelegate(delegate)` for delegation. `resolveRecipient(addr)` walks boundTo chain to root.
- **Staking**: deposit AWP into StakeNFT (ERC721 positions with lock period). `allocate(staker, agent, subnetId, amount)` — staker is explicit parameter. Auto-enumeration of agent worknets — no caller-supplied subnet list needed for freeze.
- **Epoch**: time-based on AWPEmission (`(block.timestamp - genesisTime) / epochDuration`, 1 day).
- **Emission**: exponential decay. 50% to worknets, 50% to DAO. Batch settlement via `settleEpoch(limit)`.
- **Voting**: quadratic voting with time-weighted staking positions. Two proposal types: executable (Timelock) and signal (vote-only).
- **Worknets**: registration deploys Alpha token (CREATE2 vanity address) + DEX V4 LP. Time-based mint cap on Alpha. Auto-deploys SubnetManager proxy if no custom manager provided.
- **Chain-agnostic**: Deploy script auto-selects Uniswap V4 or PancakeSwap V4 contracts based on chain ID. PoolKey struct differences (5 fields vs 6 fields) handled transparently.

## Live Testnet

| | |
|---|---|
| **API** | `https://tapi.awp.sh/api` |
| **WebSocket** | `wss://tapi.awp.sh/ws/live` |
| **Chain** | Base Mainnet (chain ID 8453) |
| **Explorer** | [basescan.org](https://basescan.org) |

### Deployed Contracts (Base Mainnet)

| Contract | Address |
|----------|---------|
| AWPToken | [`0x0000d0e3...00a1`](https://basescan.org/address/0x0000d0e38e9c6ba147b0098bb42007b942ef00a1) |
| AWPRegistry | [`0x00003a7f...00b1`](https://basescan.org/address/0x00003a7fa04c3af3adba2dc3c6622277501400b1) |
| Treasury | [`0x9ee82684...d994`](https://basescan.org/address/0x9ee82684e4214edb405d930001e9058d1913d994) |
| AlphaTokenFactory | [`0x3ebe3168...5164`](https://basescan.org/address/0x3ebe3168c898f4b05ebf0c0d17f4739e111e5164) |
| SubnetNFT | [`0x0f86ec2f...1cda`](https://basescan.org/address/0x0f86ec2f2fbf234b00b18e66e7c5e00518091cda) |
| LPManagerUni | [`0x2703d681...ebae`](https://basescan.org/address/0x2703d681ff3f7c4dc9eeed6f3ebaba3e82f8ebae) |
| AWPEmission (proxy) | [`0xd31b6fed...f1cf`](https://basescan.org/address/0xd31b6fedf7e568091b7fcf3cb5aac86c3a0ef1cf) |
| StakingVault | [`0x0367e9c2...6fca`](https://basescan.org/address/0x0367e9c2f79ab35dc65e6876405a747882296fca) |
| StakeNFT | [`0x4f7e8d44...ce39`](https://basescan.org/address/0x4f7e8d4487c0c514b72ed0e35ed707cb8acdce39) |
| SubnetManagerUni (impl) | [`0x56788237...e574`](https://basescan.org/address/0x567882378dcc11ec0d763fc5ca6c862487bbe574) |
| AWPDAO | [`0x71712119...939a`](https://basescan.org/address/0x7171211da849a2c569643fb1e8f5399ddd71939a) |

### Quick Start

```bash
# Query contract addresses + chain ID
curl https://tapi.awp.sh/api/registry

# List worknets
curl https://tapi.awp.sh/api/subnets

# Get emission info
curl https://tapi.awp.sh/api/emission/current

# Compute a vanity salt for Alpha token
curl -X POST https://tapi.awp.sh/api/vanity/compute-salt
```

## Deployment

### Prerequisites

- [Foundry](https://getfoundry.sh/) (`forge`, `cast`)
- Go 1.26+
- PostgreSQL 15+
- Redis 7+
- `jq`, `psql`

### Step 1: Configure

```bash
# Copy the template for your target chain
cp contracts/.env.example contracts/.env      # BSC
cp contracts/.env.base.example contracts/.env # Base

vim contracts/.env
```

Required fields in `contracts/.env`:
```
ETH_RPC_URL=...           # Chain RPC endpoint
DEPLOYER_PRIVATE_KEY=...  # Deployer wallet private key
GUARDIAN=...              # Emergency pause guardian address
LIQUIDITY_POOL=...       # LP wallet
AIRDROP=...              # Airdrop wallet
POOL_MANAGER=...         # DEX V4 PoolManager
POSITION_MANAGER=...     # DEX V4 PositionManager
PERMIT2=...              # Permit2 address
CL_SWAP_ROUTER=...       # DEX V4 SwapRouter / UniversalRouter
STATE_VIEW=...           # Uniswap V4 StateView (Base only, optional on BSC)
VANITY_RULE=0            # Alpha token vanity rule (0 = disabled)
```

### Step 2: Deploy Contracts

```bash
# Dry-run first (no transactions)
./scripts/deploy.sh --dry-run

# Full deployment: mine vanity salts → deploy → generate api/.env
./scripts/deploy.sh

# Skip salt mining (reuse existing SALT_* from .env)
./scripts/deploy.sh --skip-mine
```

The deploy script:
1. Builds contracts
2. Mines vanity salts (tiered, based on `VANITY_PREFIX_*` / `VANITY_SUFFIX_*` env vars)
3. Deploys all 11 contracts via deterministic CREATE2
4. Auto-selects `LPManagerUni` + `SubnetManagerUni` for non-BSC chains
5. Generates `api/.env` with all contract addresses

### Step 3: Verify Contracts (optional)

```bash
echo 'ETHERSCAN_API_KEY=your_key' >> contracts/.env
./scripts/deploy.sh --verify-only
```

### Step 4: Deploy API Services

```bash
cd api && make build

# Initialize database
psql -U postgres -d awp -f internal/db/schema.sql

# Start all 3 processes (source api/.env first)
set -a && source .env && set +a
./bin/api &       # HTTP + WebSocket (:8080 default)
./bin/indexer &   # On-chain event scanner
./bin/keeper &    # Epoch settlement + cache updates
```

The three processes:
- **api** — HTTP + WebSocket server (REST API, relay, vanity)
- **indexer** — On-chain event scanner → PostgreSQL + Redis Pub/Sub
- **keeper** — Epoch settlement + token price cache updates (every 25-30s)

### Configuration Files

| File | Purpose | Git tracked |
|------|---------|:-----------:|
| `contracts/.env` | Deployment config (keys, addresses, salts) | No |
| `contracts/.env.example` | Template (BSC) | Yes |
| `contracts/.env.base.example` | Template (Base) | Yes |
| `api/.env` | Auto-generated by deploy.sh (contract addresses) | No |

## Development

```bash
# Build contracts
cd contracts && forge build

# Run tests (232 pass, 2 require BSC fork RPC)
forge test

# Build Go backend
cd api && make build

# Run Go tests
make test

# Regenerate Go bindings after contract changes
make generate-bindings
```

## Project Structure

```
contracts/
  src/
    AWPRegistry.sol             # Unified entry point
    token/
      AWPToken.sol              # ERC20, 10B supply
      AWPEmission.sol           # UUPS proxy emission engine
      AlphaToken.sol            # Per-worknet ERC20 (CREATE2)
      AlphaTokenFactory.sol     # CREATE2 deployer + vanity rules
    core/
      StakeNFT.sol              # ERC721 staking positions
      StakingVault.sol          # Allocation bookkeeping
      SubnetNFT.sol             # Worknet ownership
      LPManager.sol             # PancakeSwap V4 CL (BSC)
      LPManagerUni.sol          # Uniswap V4 CL (Base, Ethereum)
    governance/
      AWPDAO.sol                # NFT-based voting
      Treasury.sol              # TimelockController
    subnets/
      SubnetManager.sol         # Default worknet contract — PancakeSwap V4 (BSC)
      SubnetManagerUni.sol      # Default worknet contract — Uniswap V4 (Base, Ethereum)
  test/                         # 234 tests
  script/                       # Deploy.s.sol, InitCodeHashes.s.sol

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
  subnet-developer-guide.md     # For worknet developers

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
  admin.sh                      # Hot-update rate limits, manage salt pool, view system status
```

## API Endpoints

| Group | Endpoints | Description |
|-------|-----------|-------------|
| System | `GET /api/health`, `/api/registry` | Health check, all contract addresses + chainId |
| Users | `GET /api/users/*` | User list, detail, count |
| Address | `GET /api/address/{address}/check` | Registration status (`isRegistered`, `boundTo`, `recipient`) |
| Staking | `GET /api/staking/*` | Balances, positions, allocations, worknet totals |
| Worknets | `GET /api/subnets/*` | Worknet list/detail/skills/earnings (includes `subnet_contract` + `alpha_token`) |
| Emission | `GET /api/emission/*` | Current epoch, schedule, history |
| Tokens | `GET /api/tokens/*` | AWP info, Alpha token info/price |
| Governance | `GET /api/governance/*` | Proposals, treasury |
| Relay | `POST /api/relay/*` | Gasless EIP-712: bind, set-recipient, register, allocate, deallocate, activate-subnet, register-subnet |
| Vanity | `GET/POST /api/vanity/*` | Mining params, salt pool, compute-salt, upload-salts |
| WebSocket | `WS /ws/live` | Real-time on-chain events |

## Stack

| Layer | Technology |
|-------|-----------|
| Contracts | Solidity 0.8.24, Foundry, OpenZeppelin 5.x |
| Backend | Go 1.26, Chi v5, sqlc + pgx/v5, PostgreSQL, Redis |
| Frontend | Next.js 14, Tailwind, wagmi/viem |
| Chain | Base (Uniswap V4), BSC (PancakeSwap V4) |

## Documentation

- [Architecture](docs/architecture.md) — Full technical design
- [API Reference](docs/api-reference.md) — Contract + REST + WebSocket
- [Deployment Guide](docs/deployment-guide.md) — Deploy + operations
- [Worknet Developer Guide](docs/subnet-developer-guide.md) — For worknet builders
- [Agent Skill Guide](skills-dev/agent-skill-guide.md) — Skill discovery + install

## License

[MIT](LICENSE)
