# Deployment & Operations Guide — AWP

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [Contract Deployment (BSC)](#2-contract-deployment-bsc)
3. [Database Setup](#3-database-setup)
4. [Backend Services](#4-backend-services)
5. [Configuration Reference](#5-configuration-reference)
6. [Post-Deployment Setup](#6-post-deployment-setup)
7. [Monitoring & Maintenance](#7-monitoring--maintenance)

---

## 1. Prerequisites

### System Requirements

- **Go** 1.26+
- **Foundry** (forge, cast)
- **PostgreSQL** 15+
- **Redis** 7+
- **Node.js** 18+ (for frontend, optional)
- **abigen** (go-ethereum code generator)
- **sqlc** (SQL code generator)

### Tool Installation

```bash
# Foundry
curl -L https://foundry.paradigm.xyz | bash
foundryup

# Go tools
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# Solidity dependencies
cd contracts
forge install
```

### Required Accounts & Keys

| Key | Purpose | Security Level |
|-----|---------|---------------|
| Deployer private key | One-time contract deployment, then discarded | Hot wallet (use only for deploy) |
| Keeper private key | Signs `settleEpoch` transactions every 30s | Hot wallet on server |
| Guardian address | Emergency pause (no private key needed on server) | Cold wallet / multisig |

### BSC PancakeSwap V4 Addresses (Mainnet)

```
CLPoolManager:      0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b
CLPositionManager:  0x55f4c8abA71A1e923edC303eb4fEfF14608cC226
Permit2:            0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768
```

---

## 2. Contract Deployment (BSC)

### 2.1 Prepare Environment

Create `.env` in the `contracts/` directory:

```bash
# RPC endpoint
ETH_RPC_URL=https://bsc-mainnet.example.com

# Deployer key (remove after deployment!)
DEPLOYER_PRIVATE_KEY=0x...

# Guardian address (cold wallet / multisig)
GUARDIAN=0x...

# Token distribution addresses
LIQUIDITY_POOL=0x...
AIRDROP=0x...

# PancakeSwap V4 (BSC Mainnet)
POOL_MANAGER=0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b
POSITION_MANAGER=0x55f4c8abA71A1e923edC303eb4fEfF14608cC226
PERMIT2=0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768
```

### 2.2 Deploy

```bash
cd contracts

# Dry run (simulation)
forge script script/Deploy.s.sol --rpc-url $ETH_RPC_URL --broadcast --dry-run

# Actual deployment
forge script script/Deploy.s.sol --rpc-url $ETH_RPC_URL --broadcast --verify
```

> **AlphaTokenFactory vanity rule:** The factory is constructed with a `vanityRule` (uint64) that enforces a pattern on all Alpha token addresses. Set `vanityRule=0` to disable validation (default / testnet). For mainnet, configure the desired pattern before deployment — it is **immutable** after construction. See section 4.3 of `docs/architecture.md` for the encoding format.

### 2.3 Deployment Steps (22 steps, automated)

| Step | Contract | Notes |
|------|----------|-------|
| 1 | AWPToken | 200M minted to deployer |
| 2 | AlphaToken (impl not needed) | No longer deployed separately — CREATE2 full deployment |
| 3 | AlphaTokenFactory | CREATE2 deployer with optional vanity rule |
| 4 | Treasury | TimelockController (2-day delay) |
| 5 | AWPDAO | Custom NFT-based voting (6 params, no awpRegistry dependency) |
| 6 | Grant roles | DAO gets PROPOSER + CANCELLER on Treasury |
| 7 | Renounce admin | Treasury admin permanently locked |
| 8 | AWPRegistry | Unified entry (deployer, treasury, guardian) |
| 9 | SubnetNFT | ERC721, only AWPRegistry can mint/burn |
| 10 | StakingVault | Pure allocation logic |
| 10b | StakeNFT | ERC721 position NFT (awpToken, stakingVault, awpRegistry) |
| 12 | LPManager | PancakeSwap V4 CL integration |
| 13 | AWPEmission | UUPS proxy (impl + ERC1967Proxy + initialize(awpToken, treasury, initialDailyEmission, genesisTime_, epochDuration_=86400)) |
| 14 | Add minter | AWPEmission added as sole AWP minter |
| 15 | Renounce admin | AWP minter list permanently locked |
| 16 | Configure factory | `factory.setAddresses(awpRegistry)` — links to AWPRegistry and renounces ownership |
| 17 | Initialize registry | All module addresses injected into AWPRegistry |
| 18-22 | Distribute AWP | Treasury 90M, LP 10M, Airdrop 100M |

### 2.4 Verify Deployment

```bash
# Verify all contracts on BSCScan
forge verify-contract <address> <ContractName> --chain bsc

# Check key invariants
cast call <AWPToken> "admin()" --rpc-url $ETH_RPC_URL
# Should return 0x0000000000000000000000000000000000000000

cast call <AWPToken> "minters(address)" <AWPEmission> --rpc-url $ETH_RPC_URL
# Should return true

cast call <AWPRegistry> "registryInitialized()" --rpc-url $ETH_RPC_URL
# Should return true
```

### 2.5 Record Deployed Addresses

After deployment, save all addresses for backend configuration:

```bash
# Output format from Deploy.s.sol console.log:
# Step 1: AWPToken deployed at 0x...
# Step 8: AWPRegistry at 0x...
# Step 9: SubnetNFT at 0x...
# Step 13: AWPEmission proxy at 0x...
# ...
```

---

## 3. Database Setup

### 3.1 Create Database

```bash
createdb awp
```

### 3.2 Apply Schema

```bash
psql awp < api/internal/db/schema.sql
```

### 3.3 Schema Overview

| Table | Purpose |
|-------|---------|
| `users` | User addresses with `bound_to` and `recipient` columns |
| `subnets` | Subnet metadata and status |
| `stake_allocations` | (user, agent, subnet) stake allocations |
| `user_balances` | User allocation totals (total_allocated only) |
| `stake_positions` | StakeNFT positions (tokenId, owner, amount, lockEndTime, createdAt) |
| `epochs` | Epoch settlement records |
| `recipient_awp_distributions` | Per-recipient AWP emission history |
| `proposals` | DAO governance proposals |
| `sync_states` | Indexer block sync progress |

### 3.4 Regenerate Go Code (after schema changes)

```bash
cd api
sqlc generate
```

---

## 4. Backend Services

Three independent Go processes:

| Process | Purpose | Frequency |
|---------|---------|-----------|
| **api** | HTTP + WebSocket server | Always running |
| **indexer** | Scans chain events → PostgreSQL + Redis Pub/Sub | Polls every 3s |
| **keeper** | Calls `settleEpoch`, updates Redis caches | Cron: 30s (settle), 5m (prices) |

### 4.1 Build

```bash
cd api

# Build all three binaries
make build
# Produces: bin/api, bin/indexer, bin/keeper

# Or build individually
go build -o bin/api ./cmd/api
go build -o bin/indexer ./cmd/indexer
go build -o bin/keeper ./cmd/keeper
```

### 4.2 Regenerate Bindings (after contract changes)

```bash
cd api
make generate
# Runs: forge build → abigen (all contracts) → sqlc generate
```

### 4.3 Configure Environment

Create `.env` for the backend (or export variables):

```bash
# Database
DATABASE_URL=postgres://user:pass@localhost:5432/awp?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379/0

# HTTP server (api only)
HTTP_ADDR=:8080

# BSC RPC
RPC_URL=https://bsc-mainnet.example.com

# Contract addresses (from deployment output)
AWP_REGISTRY_ADDRESS=0x...
AWP_TOKEN_ADDRESS=0x...
AWP_EMISSION_ADDRESS=0x...
STAKING_VAULT_ADDRESS=0x...
STAKENFT_ADDRESS=0x...
SUBNETNFT_ADDRESS=0x...
DAO_ADDRESS=0x...
TREASURY_ADDRESS=0x...

# Keeper only
KEEPER_PRIVATE_KEY=abcdef1234...  # No 0x prefix
```

### 4.4 Start Services

```bash
# Start all three (use systemd, supervisor, or screen in production)

# 1. API server
./bin/api

# 2. Indexer
./bin/indexer

# 3. Keeper
./bin/keeper
```

### 4.5 Systemd Service Examples

#### API Service

```ini
# /etc/systemd/system/awp-api.service
[Unit]
Description=AWP API Server
After=postgresql.service redis.service

[Service]
Type=simple
User=awp
WorkingDirectory=/opt/awp
ExecStart=/opt/awp/bin/api
EnvironmentFile=/opt/awp/.env
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

#### Indexer Service

```ini
# /etc/systemd/system/awp-indexer.service
[Unit]
Description=AWP Event Indexer
After=postgresql.service redis.service

[Service]
Type=simple
User=awp
WorkingDirectory=/opt/awp
ExecStart=/opt/awp/bin/indexer
EnvironmentFile=/opt/awp/.env
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

#### Keeper Service

```ini
# /etc/systemd/system/awp-keeper.service
[Unit]
Description=AWP Keeper (Epoch Settlement)
After=network.target

[Service]
Type=simple
User=awp
WorkingDirectory=/opt/awp
ExecStart=/opt/awp/bin/keeper
EnvironmentFile=/opt/awp/.env
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable awp-api awp-indexer awp-keeper
sudo systemctl start awp-api awp-indexer awp-keeper
```

---

## 5. Configuration Reference

### 5.1 Environment Variables

| Variable | Required By | Default | Description |
|----------|------------|---------|-------------|
| `DATABASE_URL` | api, indexer | `postgres://postgres:postgres@localhost:5432/awp?sslmode=disable` | PostgreSQL connection string |
| `REDIS_URL` | api, indexer, keeper | `redis://localhost:6379/0` | Redis connection string |
| `HTTP_ADDR` | api | `:8080` | HTTP listen address |
| `RPC_URL` | indexer, keeper | `https://bsc-testnet-rpc.publicnode.com` | BSC JSON-RPC endpoint |
| `AWP_REGISTRY_ADDRESS` | indexer | — | AWPRegistry contract address |
| `AWP_TOKEN_ADDRESS` | keeper | — | AWP token address |
| `AWP_EMISSION_ADDRESS` | indexer, keeper | — | AWPEmission proxy address |
| `STAKING_VAULT_ADDRESS` | indexer | — | StakingVault address |
| `STAKENFT_ADDRESS` | indexer | — | StakeNFT address |
| `SUBNETNFT_ADDRESS` | indexer | — | SubnetNFT address |
| `DAO_ADDRESS` | indexer | — | AWPDAO address |
| `TREASURY_ADDRESS` | api | — | Treasury (Timelock) address |
| `KEEPER_PRIVATE_KEY` | keeper | — | Hex private key for signing settle transactions |

### 5.2 Contract Constants

| Constant | Value | Location |
|----------|-------|----------|
| `INITIAL_DAILY_EMISSION` | 15,800,000 AWP | Deploy.s.sol |
| `EPOCH_DURATION` | 1 day (86400s) | AWPEmission (initialized via Deploy.s.sol) |
| `TIMELOCK_DELAY` | 2 days (172800s) | Deploy.s.sol |
| `POOL_FEE` | 1% (10000) | LPManager.sol |
| `TICK_SPACING` | 200 | LPManager.sol |
| `MAX_ACTIVE_SUBNETS` | 10,000 | AWPRegistry.sol |
| `maxRecipients` | 10,000 | AWPEmission.sol |
| `DECAY_FACTOR` | 996844 / 1000000 | AWPEmission.sol |
| `EMISSION_SPLIT_BPS` | 5000 (50%) | AWPEmission.sol |

### 5.3 Redis Key Spec

| Key | Type | TTL | Updated By |
|-----|------|-----|------------|
| `emission_current` | JSON | 30s | Keeper (every 5m) |
| `awp_info` | JSON | 1m | Keeper (every 5m) |
| `alpha_price:{subnetId}` | JSON | 10m | External |
| `chain_events` | Pub/Sub channel | — | Indexer |

---

## 6. Post-Deployment Setup

### 6.1 Configure Oracle Network

After deployment, AWPEmission has no oracle configured. Weights cannot be submitted until oracles are set up.

**Via DAO Governance Proposal:**

```solidity
// Create a proposal to call setOracleConfig on AWPEmission
targets = [awpEmissionProxy];
calldatas = [abi.encodeCall(
    AWPEmission.setOracleConfig,
    (oracleAddresses, threshold)  // e.g., 5 oracles, threshold 3
)];
```

**For Testing (direct Timelock call):**

```bash
cast send <AWPEmission> "setOracleConfig(address[],uint256)" \
  "[0xOracle1,0xOracle2,0xOracle3]" 2 \
  --private-key $TIMELOCK_KEY \
  --rpc-url $RPC_URL
```

### 6.2 Verify Indexer Sync

```bash
# Check sync progress
psql awp -c "SELECT * FROM sync_states;"
# Should show: indexer | <recent_block_number>

# Verify events are being indexed
psql awp -c "SELECT COUNT(*) FROM users;"
```

### 6.3 Verify Keeper Operation

```bash
# Check if epoch is settling
cast call <AWPEmission> "settleProgress()" --rpc-url $RPC_URL
# 0 = idle, >0 = in progress

# Check current epoch (on AWPEmission)
cast call <AWPEmission> "currentEpoch()" --rpc-url $RPC_URL

# Check settlement progress
cast call <AWPEmission> "settledEpoch()" --rpc-url $RPC_URL

# Verify epoch duration (on AWPEmission)
cast call <AWPEmission> "epochDuration()" --rpc-url $RPC_URL

# Check Redis cache
redis-cli GET emission_current
redis-cli GET awp_info
```

### 6.4 First Subnet Registration

```bash
# 1. Approve AWP for LP (no mandatory registration needed)
cast send <AWPToken> "approve(address,uint256)" <AWPRegistry> 1000000000000000000000000 \
  --private-key $USER_KEY --rpc-url $RPC_URL

# 2. Register subnet (salt=0x00..00 uses subnetId as CREATE2 salt)
cast send <AWPRegistry> "registerSubnet((string,string,address,bytes32,uint128,string))" \
  "(\"My Subnet\",\"MSUB\",0x0000000000000000000000000000000000000000,0x0000000000000000000000000000000000000000000000000000000000000000,0,\"ipfs://QmSkills...\")" \
  --private-key $USER_KEY --rpc-url $RPC_URL

# 3. Activate
cast send <AWPRegistry> "activateSubnet(uint256)" 1 \
  --private-key $USER_KEY --rpc-url $RPC_URL
```

---

## 7. Monitoring & Maintenance

### 7.1 Health Checks

```bash
# API health
curl http://localhost:8080/api/health
# {"status": "ok"}

# Indexer: check sync lag
LATEST=$(cast block-number --rpc-url $RPC_URL)
SYNCED=$(psql awp -t -c "SELECT last_block FROM sync_states WHERE contract_name='indexer'")
echo "Lag: $(($LATEST - $SYNCED)) blocks"

# Keeper: check epoch currency
cast call <AWPEmission> "currentEpoch()" --rpc-url $RPC_URL
cast call <AWPEmission> "settledEpoch()" --rpc-url $RPC_URL
```

### 7.2 Common Operations

**Emergency Pause (Guardian):**
```bash
cast send <AWPRegistry> "pause()" --private-key $GUARDIAN_KEY --rpc-url $RPC_URL
```

**Resume (via DAO/Timelock):**
```bash
cast send <AWPRegistry> "unpause()" --private-key $TIMELOCK_KEY --rpc-url $RPC_URL
```

**Upgrade AWPEmission (via DAO Proposal):**
```bash
# 1. Deploy new implementation
forge create AWPEmission --rpc-url $RPC_URL --private-key $DEPLOYER_KEY

# 2. Create DAO proposal to upgrade
# targets = [awpEmissionProxy]
# calldatas = [abi.encodeCall(UUPSUpgradeable.upgradeToAndCall, (newImpl, ""))]
```

### 7.3 Log Monitoring

All three services output structured JSON logs:

```bash
# API
journalctl -u awp-api -f | jq '.msg'

# Indexer
journalctl -u awp-indexer -f | jq 'select(.msg == "scan complete")'

# Keeper
journalctl -u awp-keeper -f | jq 'select(.msg | contains("SettleEpoch"))'
```

### 7.4 Backup

```bash
# Database backup
pg_dump awp > cortexia_backup_$(date +%Y%m%d).sql

# No need to backup Redis (caches are rebuilt by Keeper)
# No need to backup indexer state (re-indexable from chain)
```
