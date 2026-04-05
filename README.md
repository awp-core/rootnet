# Agent Working Protocol

<p align="center">
  <a href="https://awp.pro/">
    <img src="assets/banner.png" alt="AWP - Agent Work Protocol" width="800">
  </a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Base-0052FF?style=flat&logo=coinbase&logoColor=white" alt="Base">
  <img src="https://img.shields.io/badge/Ethereum-3C3C3D?style=flat&logo=ethereum&logoColor=white" alt="Ethereum">
  <img src="https://img.shields.io/badge/Arbitrum-28A0F0?style=flat&logo=arbitrum&logoColor=white" alt="Arbitrum">
  <img src="https://img.shields.io/badge/BNB_Chain-F0B90B?style=flat&logo=bnbchain&logoColor=white" alt="BNB Chain">
  <img src="https://img.shields.io/badge/Uniswap_V4-FF007A?style=flat&logo=uniswap&logoColor=white" alt="Uniswap V4">
  <img src="https://img.shields.io/badge/PancakeSwap_V4-1FC7D4?style=flat" alt="PancakeSwap V4">
  <img src="https://img.shields.io/badge/Solidity_0.8.24-363636?style=flat&logo=solidity&logoColor=white" alt="Solidity">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/License-MIT-97CA00?style=flat" alt="MIT">
</p>

> **Mainnet.** AWP is deployed on Base, Ethereum, Arbitrum, and BSC. Protocol parameters may change via governance.

## Abstract

AWP is a decentralized **Agent Working** protocol deployed on Base (Uniswap V4), Ethereum (Uniswap V4), Arbitrum (Uniswap V4), and BNB Smart Chain (PancakeSwap V4). The protocol establishes a permissionless marketplace where autonomous AI agent networks (*worknets*) compete for protocol-level emission rewards through stake-weighted Guardian consensus. Each worknet deploys an independent economy backed by a dedicated ERC-20 token (WorknetToken), with initial liquidity bootstrapped via Concentrated Liquidity DEX at registration time.

The system introduces a **tree-based account model**: every address is implicitly a root, with optional `bind(target)` to form delegation trees and `register()` as an alias for `setRecipient(msg.sender)`. Users deposit AWP tokens into non-fungible position NFTs (veAWP) with time-locked commitments, then allocate stake across (agent, worknet) triples via explicit `allocate(staker, agent, worknetId, amount)`. An exponentially-decaying emission schedule distributes newly-minted AWP to worknet managers proportional to Guardian-assigned governance weights, with 100% to recipients; Guardian includes treasury in recipients for DAO share, governed by NFT-weighted quadratic voting.

Key design contributions include: (1) a gasless relay layer enabling device-bound agents to participate without holding native gas tokens; (2) an ERC-1363 `mintAndCall` emission pathway that auto-triggers configurable AWP handling strategies (reserve, single-sided liquidity provision, or buyback-and-burn) at the worknet manager level; (3) a tiered CREATE2 vanity address system with pre-mined salt pools for deterministic cross-chain deployment; and (4) a modular worknet architecture where a default `WorknetManager` proxy contract provides Merkle-based reward distribution, multi-role access control, and DEX integration out of the box, while advanced operators may deploy custom manager contracts.

The protocol consists of 11 Solidity contracts (Foundry, Solidity 0.8.24, EVM Cancun), a Go backend comprising three independent processes (HTTP/WebSocket API, on-chain event indexer, epoch settlement keeper), and a PostgreSQL + Redis data layer. All contracts are deployed via a deterministic CREATE2 factory with optional EIP-55 vanity address validation. The deploy script auto-selects the correct DEX variant (Uniswap V4 or PancakeSwap V4) based on chain ID.

> **Note:** The AWP Emission mechanism (AWPEmission contract, Guardian consensus, epoch settlement) is under active design and has not been finalized. All emission-related descriptions in this document are preliminary and subject to change.

## Architecture

```
User
 ├── AWPRegistry ─── bind / allocate / worknet lifecycle / delegation
 │    ├── veAWP ── ERC721 position NFTs (deposit AWP + lock)
 │    ├── AWPAllocator ── allocation bookkeeping (auto-enumerates agent worknets)
 │    ├── AWPWorkNet ── worknet ownership NFTs
 │    └── LPManager / LPManagerUni ── DEX V4 CL liquidity (PancakeSwap / Uniswap)
 │
 ├── AWPEmission (UUPS proxy) ── epoch settlement + AWP minting
 │    └── Guardian multi-sig → submitAllocations → settleEpoch
 │
 ├── AWPDAO ── NFT-based voting (executable + signal proposals)
 │    └── Treasury (TimelockController) ── governance execution
 │
 └── WorknetTokenFactory ── CREATE2 per-worknet tokens with vanity addresses
```

**11 contracts**, 3 Go backend processes (API / Indexer / Keeper), PostgreSQL, Redis.

**Chain-agnostic DEX support:**
- **Base / Ethereum / Arbitrum** — Uniswap V4: `LPManagerUni` + `WorknetManagerUni` (auto-selected for chainId != 56/97)
- **BSC** — PancakeSwap V4: `LPManager` + `WorknetManager` (auto-selected for chainId 56/97)

**Backend API** provides:
- Read-only REST API + WebSocket real-time events
- Gasless relay endpoints (`/api/relay/bind`, `/api/relay/set-recipient`, `/api/relay/register-worknet`, etc.) — EIP-712 signed, relayer pays gas
- Vanity salt mining (`/api/vanity/compute-salt`) — uses Foundry `cast create2` for high-speed parallel mining


## Key Design

- **Account System V2**: No mandatory registration — every address is implicitly a root. `register()` is optional (= `setRecipient(msg.sender)`). Tree-based binding via `bind(target)` with anti-cycle check. No address mutual exclusion. `grantDelegate(delegate)` / `revokeDelegate(delegate)` for delegation. `resolveRecipient(addr)` walks boundTo chain to root.
- **Staking**: deposit AWP into veAWP (ERC721 positions with lock period). `allocate(staker, agent, worknetId, amount)` — staker is explicit parameter. Auto-enumeration of agent worknets — no caller-supplied worknet list needed for freeze.
- **Epoch**: time-based on AWPEmission (`(block.timestamp - genesisTime) / epochDuration`, 1 day).
- **Emission**: exponential decay. 100% to recipients; Guardian includes treasury in recipients for DAO share. Batch settlement via `settleEpoch(limit)`.
- **Voting**: quadratic voting with time-weighted staking positions. Two proposal types: executable (Timelock) and signal (vote-only).
- **Worknets**: registration deploys WorknetToken (CREATE2 vanity address) + DEX V4 LP. Time-based mint cap on WorknetToken. Auto-deploys WorknetManager proxy if no custom manager provided.
- **Chain-agnostic**: Deploy script auto-selects Uniswap V4 or PancakeSwap V4 contracts based on chain ID. PoolKey struct differences (5 fields vs 6 fields) handled transparently.

## Multi-Chain

AWP deploys identical contracts on multiple EVM chains using CREATE2 (same deployer + salts = same addresses).

**Supported chains:** Base, Ethereum, Arbitrum, BSC (configured in `chains.yaml`)

**WorknetId encoding:** `(chainId << 64) | localCounter` — globally unique across all chains. Use `extractChainId(worknetId)` / `extractLocalId(worknetId)` to decode.

**Cross-chain allocate:** Users stake AWP on one chain and can allocate to worknets on ANY chain. The AWPAllocator only checks local balance; worknet validity is verified off-chain by the Guardian/indexer.

**Per-chain independence:**
- Each chain has its own AWPToken, AWPEmission, AWPDAO, and Treasury
- Emission quotas are coordinated by the Guardian across chains
- DAO voting power is aggregated off-chain from all chains' veAWP positions

**Deploy to a new chain:**
```bash
# List available chains
./scripts/deploy-multichain.sh --list

# Deploy to a specific chain
./scripts/deploy-multichain.sh base

# Deploy to all chains
./scripts/deploy-multichain.sh --all
```

## Live Mainnet

| | |
|---|---|
| **API** | `https://api.awp.sh/api` |
| **WebSocket** | `wss://api.awp.sh/ws/live` |
| **Chains** | Base (8453), Ethereum (1), Arbitrum (42161), BSC (56) |
| **Explorer** | [basescan.org](https://basescan.org) |

### Deployed Contracts (Base Mainnet)

| Contract | Address |
|----------|---------|
| AWPToken | [`0x0000A105...00A1`](https://basescan.org/address/0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1) |
| AWPRegistry | [`0x0000F34E...001A`](https://basescan.org/address/0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A) |
| Treasury | [`0x82562023...759e`](https://basescan.org/address/0x82562023a053025F3201785160CaE6051efD759e) |
| WorknetTokenFactory | [`0x0000D499...FAC7`](https://basescan.org/address/0x0000D4996BDBb99c772e3fA9f0e94AB52AAFFAC7) |
| AWPWorkNet | [`0x00000bfb...00A7`](https://basescan.org/address/0x00000bfbdEf8533E5F3228c9C846522D906100A7) |
| LPManager | [`0x00001961...00A2`](https://basescan.org/address/0x00001961b9AcCD86b72DE19Be24FaD6f7c5b00A2) |
| AWPEmission (proxy) | [`0x3C9cB73f...EaA9`](https://basescan.org/address/0x3C9cB73f8B81083882c5308Cce4F31f93600EaA9) |
| AWPAllocator | [`0x0000D6BB...00AA`](https://basescan.org/address/0x0000D6BB5e040E35081b3AaF59DD71b21C9800AA) |
| veAWP | [`0x0000b534...00A8`](https://basescan.org/address/0x0000b534C63D78212f1BDCc315165852793A00A8) |
| WorknetManager (impl) | [`0x000011EE...00A9`](https://basescan.org/address/0x000011EE4117c52dC0Eb146cBC844cb155B200A9) |
| AWPDAO | [`0x00006879...0DA0`](https://basescan.org/address/0x00006879f79f3Da189b5D0fF6e58ad0127Cc0DA0) |
| Guardian | [`0x000002bE...00A3`](https://basescan.org/address/0x000002bEfa6A1C99A710862Feb6dB50525dF00A3) |

### Quick Start

```bash
# Query contract addresses + chain ID
curl https://api.awp.sh/api/registry

# List worknets
curl https://api.awp.sh/api/worknets

# Get emission info
curl https://api.awp.sh/api/emission/current

# Compute a vanity salt for WorknetToken
curl -X POST https://api.awp.sh/api/vanity/compute-salt
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
4. Auto-selects `LPManagerUni` + `WorknetManagerUni` for non-BSC chains
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
      WorknetToken.sol          # Per-worknet ERC20 (CREATE2)
      WorknetTokenFactory.sol   # CREATE2 deployer + vanity rules
    core/
      veAWP.sol                 # ERC721 staking positions
      AWPAllocator.sol          # Allocation bookkeeping
      AWPWorkNet.sol            # Worknet ownership
      LPManager.sol             # PancakeSwap V4 CL (BSC)
      LPManagerBase.sol         # Shared LP logic
      LPManagerUni.sol          # Uniswap V4 CL (Base, Ethereum, Arbitrum)
    governance/
      AWPDAO.sol                # NFT-based voting
      Treasury.sol              # TimelockController
    worknets/
      WorknetManager.sol        # Default worknet contract — PancakeSwap V4 (BSC)
      WorknetManagerBase.sol    # Shared worknet manager logic
      WorknetManagerUni.sol     # Default worknet contract — Uniswap V4 (Base, Ethereum, Arbitrum)
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
| Worknets | `GET /api/worknets/*` | Worknet list/detail/skills/earnings (includes `worknet_contract` + `worknet_token`) |
| Emission | `GET /api/emission/*` | Current epoch, schedule, history |
| Tokens | `GET /api/tokens/*` | AWP info, WorknetToken info/price |
| Governance | `GET /api/governance/*` | Proposals, treasury |
| Relay | `POST /api/relay/*` | Gasless EIP-712: bind, set-recipient, register, allocate, deallocate, activate-worknet, register-worknet |
| Vanity | `GET/POST /api/vanity/*` | Mining params, salt pool, compute-salt, upload-salts |
| WebSocket | `WS /ws/live` | Real-time on-chain events |

## Stack

| Layer | Technology |
|-------|-----------|
| Contracts | Solidity 0.8.24, Foundry, OpenZeppelin 5.x |
| Backend | Go 1.26, Chi v5, sqlc + pgx/v5, PostgreSQL, Redis |
| Frontend | Next.js 14, Tailwind, wagmi/viem |
| Chain | Base (Uniswap V4), Ethereum (Uniswap V4), Arbitrum (Uniswap V4), BSC (PancakeSwap V4) |

## Documentation

- [Architecture](docs/architecture.md) — Full technical design
- [API Reference](docs/api-reference.md) — Contract + REST + WebSocket
- [Deployment Guide](docs/deployment-guide.md) — Deploy + operations
- [Worknet Developer Guide](docs/subnet-developer-guide.md) — For worknet builders
- [Agent Skill Guide](skills-dev/agent-skill-guide.md) — Skill discovery + install

## License

[MIT](LICENSE)
