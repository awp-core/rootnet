# Multi-Chain AWP Architecture Design

**Date:** 2026-03-29
**Status:** Approved
**Scope:** Core architecture change — multi-chain deployment with off-chain aggregation

---

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Target chains | Ethereum L1 + EVM L2 (Base, Arbitrum, BSC, etc.) | Broad EVM coverage; CREATE2 gives same addresses across all chains |
| Cross-chain communication | Off-chain aggregation (Indexer) | Matches existing oracle/indexer/relay architecture; no bridge dependency |
| SubnetId uniqueness | Contract-layer: `(chainId << 64) \| localCounter` | Allocate happens on-chain with global IDs; must be verifiable on-chain |
| DAO | Per-chain independent DAO + off-chain aggregated voting | Each chain's Treasury manages its own AWPRegistry; oracle submits aggregated votes |
| AWP Token | Per-chain independent mint + emission | No bridge dependency; each chain fully autonomous; oracle coordinates emission quotas |
| Allocate location | On staker's chain; indexer aggregates globally | StakingVault must verify staker balance locally; subnetId is just a number |

---

## Architecture Overview

```
                    ┌─────────────────────────┐
                    │   Unified API + Indexer  │
                    │   (Go, single DB+Redis)  │
                    │                          │
                    │  Aggregates ALL chains:  │
                    │  - subnets (global view) │
                    │  - allocations (global)  │
                    │  - voting power (global) │
                    └──────┬──────┬──────┬─────┘
                           │      │      │
              ┌────────────┘      │      └────────────┐
              ▼                   ▼                    ▼
    ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
    │  Ethereum L1     │ │  Base            │ │  Arbitrum        │
    │  (chainId 1)     │ │  (chainId 8453)  │ │  (chainId 42161) │
    │                  │ │                  │ │                  │
    │  Full 11-contract│ │  Full 11-contract│ │  Full 11-contract│
    │  deployment      │ │  deployment      │ │  deployment      │
    │  (same addresses)│ │  (same addresses)│ │  (same addresses)│
    └──────────────────┘ └──────────────────┘ └──────────────────┘
```

Each chain runs the full protocol independently. The Indexer aggregates data from all chains into a single database, providing a unified global view.

---

## Contract Changes

### 1. SubnetId — Global Uniqueness via Contract-Layer Encoding

**Current:** `_nextSubnetId++` (local counter, starts at 1)

**New:** `globalSubnetId = (block.chainid << 64) | _nextLocalId++`

```solidity
// AWPRegistry.sol
uint256 private _nextLocalId; // set to 1 in initialize()

function _registerSubnet(...) internal returns (uint256) {
    uint256 subnetId = (block.chainid << 64) | _nextLocalId++;
    // ... rest of registration logic uses subnetId
}
```

**Properties:**
- Upper 192 bits: chainId (fits any EVM chainId)
- Lower 64 bits: local counter (supports 18 quintillion subnets per chain)
- Globally unique across all chains
- Human-readable: can extract chainId and localId from any subnetId
- SubnetNFT tokenId = globalSubnetId

**Helper view functions:**
```solidity
function extractChainId(uint256 subnetId) public pure returns (uint256) {
    return subnetId >> 64;
}
function extractLocalId(uint256 subnetId) public pure returns (uint256) {
    return subnetId & ((1 << 64) - 1);
}
```

### 2. StakingVault — Remove Subnet Status Check for Cross-Chain Allocate

**Current:** AWPRegistry.allocate checks `subnets[subnetId].status == Active` before forwarding to StakingVault.

**New:** Remove the on-chain subnet status check entirely. StakingVault only verifies:
- `userTotalAllocated + amount <= userTotalStaked` (balance check)
- `subnetId != 0` (zero-ID rejection)

Subnet validity (active, exists, correct chain) is verified off-chain by the indexer/oracle. This is necessary because a subnet registered on Arbitrum cannot be verified on-chain by Base's AWPRegistry.

### 3. AWPToken — Per-Chain Independent

Each chain deploys an independent AWPToken with:
- Same MAX_SUPPLY (10B) — but each chain's actual mint is capped by its emission quota
- Constructor mint amount is configurable per chain (not hardcoded 200M)
- AWPEmission on each chain has its own `currentDailyEmission` controlled by oracle

**Emission coordination:**
- Oracle reads all chains' emission state
- Submits per-chain weights that sum to the protocol-wide target emission
- Each chain's AWPEmission independently settles its allocated portion
- Total cross-chain emission stays within protocol bounds via oracle enforcement

### 4. AWPDAO — Per-Chain with Off-Chain Aggregated Voting

Each chain has:
- `AWPDAO` — Governor contract
- `Treasury` — TimelockController (manages that chain's AWPRegistry upgrades)

**Voting flow:**
1. Indexer aggregates StakeNFT positions from ALL chains
2. Computes global voting power per address (sum of all chains' positions)
3. Voting happens off-chain (Snapshot-style) using aggregated power
4. Oracle submits vote results to each chain's AWPDAO for on-chain execution
5. Each chain's Treasury executes governance actions for its own contracts

### 5. Deploy Script — Per-Chain Configuration

```yaml
# chains.yaml — multi-chain deployment configuration
chains:
  - chainId: 1
    name: Ethereum
    rpcUrl: ${ETH_RPC_URL}
    dex: uniswap_v4
    initialMint: 100000000  # 100M AWP constructor mint
    emissionQuota: 40       # 40% of total daily emission
    poolManager: "0x..."
    positionManager: "0x..."
    permit2: "0x..."
    swapRouter: "0x..."
    stateView: "0x..."

  - chainId: 8453
    name: Base
    rpcUrl: ${BASE_RPC_URL}
    dex: uniswap_v4
    initialMint: 50000000   # 50M AWP
    emissionQuota: 30       # 30%
    poolManager: "0x498581ff718922c3f8e6a244956af099b2652b2b"
    positionManager: "0x7c5f5a4bbd8fd63184577525326123b519429bdc"
    permit2: "0x000000000022D473030F116dDEE9F6B43aC78BA3"
    swapRouter: "0x6ff5693b99212da76ad316178a184ab56d299b43"
    stateView: "0xa3c0c9b65bad0b08107aa264b0f3db444b867a71"

  - chainId: 42161
    name: Arbitrum
    rpcUrl: ${ARB_RPC_URL}
    dex: uniswap_v4
    initialMint: 50000000   # 50M AWP
    emissionQuota: 30       # 30%
    poolManager: "0x..."
    # ...
```

**Deployment ensures identical addresses:**
- Same deployer wallet (same private key) on all chains
- Same CREATE2 salts (from salt.json)
- Same contract bytecode
- Only DEX addresses and chain-specific config differ (passed via dexConfig, not in bytecode)

---

## API / Indexer Changes

### Multi-Chain Indexer

**Current:** Single-chain indexer scanning one RPC.

**New:** Multi-chain indexer with one goroutine per chain, all writing to the same DB.

```
┌─────────────┐
│  Indexer     │
│             │
│  chain_1    │──→ goroutine: poll Base events
│  chain_2    │──→ goroutine: poll Arbitrum events
│  chain_3    │──→ goroutine: poll Ethereum events
│             │
│  shared DB  │──→ PostgreSQL (all chains in one DB)
│  shared     │──→ Redis (unified Pub/Sub)
└─────────────┘
```

### Database Schema Changes

All chain-scoped tables add `chain_id BIGINT NOT NULL`:

```sql
-- subnets: subnet_id is already globally unique (chainId << 64 | localId)
-- but chain_id column enables efficient per-chain queries
ALTER TABLE subnets ADD COLUMN chain_id BIGINT NOT NULL;
ALTER TABLE users ADD COLUMN chain_id BIGINT NOT NULL;
ALTER TABLE stake_allocations ADD COLUMN chain_id BIGINT NOT NULL;
ALTER TABLE stake_positions ADD COLUMN chain_id BIGINT NOT NULL;
ALTER TABLE user_balances ADD COLUMN chain_id BIGINT NOT NULL;
ALTER TABLE epochs ADD COLUMN chain_id BIGINT NOT NULL;
ALTER TABLE vanity_salts ADD COLUMN chain_id BIGINT NOT NULL;
ALTER TABLE sync_states — already keyed by contract_name, add chain_id
ALTER TABLE indexed_blocks ADD COLUMN chain_id BIGINT NOT NULL;

-- Users exist on multiple chains (same address, different positions)
-- Primary keys become composite: (chain_id, address) for users, etc.
```

### New API Endpoints

```
GET /api/chains                     — list all supported chains
GET /api/subnets?chainId=8453       — filter subnets by chain (optional)
GET /api/subnets/{subnetId}         — subnet detail (includes chain_id, chain_name)
GET /api/staking/user/{addr}/global — aggregated balance across all chains
GET /api/staking/user/{addr}/allocations?chainId=... — per-chain or global
GET /api/registry?chainId=8453      — per-chain contract addresses
```

### Subnet Response (new fields)

```json
{
  "subnet_id": 155865067315552256001,
  "local_id": 1,
  "chain_id": 8453,
  "chain_name": "Base",
  "owner": "0x...",
  "name": "Benchmark Subnet",
  "status": "Active",
  "alpha_token": "0x...",
  ...
}
```

---

## Key Invariants

1. **SubnetId globally unique:** `(chainId << 64) | localId` — enforced at contract level
2. **Allocate is local:** User calls allocate on their staking chain; target subnetId can be from any chain
3. **Balance check is local:** `userTotalAllocated <= userTotalStaked` checked on staker's chain only
4. **Subnet status is off-chain:** Oracle/indexer validates subnet is active before including in emission
5. **Emission is per-chain:** Each chain's AWPEmission operates independently; oracle coordinates quotas
6. **DAO is per-chain:** Each chain's Treasury governs its own contracts; voting power aggregated off-chain
7. **Addresses are identical:** CREATE2 + same deployer + same salts = same addresses on all chains

---

## Migration Path (from current single-chain)

1. **Phase 1:** Contract changes (subnetId encoding, remove allocate subnet check)
2. **Phase 2:** Deploy to additional chains using same salts
3. **Phase 3:** Indexer multi-chain support (goroutine per chain, chain_id in DB)
4. **Phase 4:** API changes (chain filter, global aggregation endpoints)
5. **Phase 5:** Oracle multi-chain emission coordination
6. **Phase 6:** DAO off-chain voting aggregation

---

## Risk Analysis

| Risk | Mitigation |
|------|-----------|
| Oracle is single point of failure for emission coordination | Oracle is already trusted for single-chain; multi-sig oracle set |
| SubnetId encoding change breaks existing data | Migration: re-encode existing subnets with `(currentChainId << 64) \| existingId` |
| Same address on all chains means one key compromise affects all | Already true for deployer; use hardware wallet + multisig |
| Indexer DB grows large with multiple chains | Partition tables by chain_id; archive old epochs |
| One chain's RPC failure affects global view | Indexer goroutines are independent; API returns partial data with warning |
