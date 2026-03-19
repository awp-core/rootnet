# AWP — API Reference

## Table of Contents

1. [Smart Contract API](#1-smart-contract-api)
   - [AWPRegistry](#11-rootnet)
   - [AWPEmission](#12-awpemission)
   - [StakingVault](#13-stakingvault)
   - [StakeNFT](#13b-stakenft)
   - [~~AccessManager~~](#14-accessmanager-removed)
   - [AWPToken](#15-awptoken)
   - [AlphaToken](#16-alphatoken)
   - [LPManager](#17-lpmanager)
   - [SubnetNFT](#18-subnetnft)
   - [AWPDAO](#19-awpdao)
2. [REST API](#2-rest-api)
   - [System](#21-system)
   - [Users](#22-users)
   - [Agents](#23-agents)
   - [Staking](#24-staking)
   - [Subnets](#25-subnets)
   - [Emission](#26-emission)
   - [Tokens](#27-tokens)
   - [Governance](#28-governance)
   - [WebSocket](#29-websocket)
   - [Relay](#210-relay-gasless-transactions)
3. [Data Structures](#3-data-structures)
4. [Events](#4-events)
5. [Error Codes](#5-error-codes)
6. [Constants](#6-constants)

---

## 1. Smart Contract API

### 1.1 AWPRegistry

> Unified entry point for subnet management and staking. All user-facing write operations go through AWPRegistry.

#### Account System (V2)

> EIP-712 domain name: "AWPRegistry". No mandatory registration — every address is implicitly a root.

| Function | Access | Description |
|----------|--------|-------------|
| `register()` | Anyone | Optional; equivalent to `setRecipient(msg.sender)` |
| `bind(address target)` | Anyone | Tree-based binding with anti-cycle check |
| `bindFor(address user, address target, uint256 deadline, uint8 v, bytes32 r, bytes32 s)` | Anyone | Gasless bind via EIP-712 |
| `setRecipient(address recipient)` | Anyone | Set reward recipient address |
| `setRecipientFor(address user, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s)` | Anyone | Gasless set recipient via EIP-712 |
| `grantDelegate(address delegate)` | Anyone | Grant delegation to an address |
| `revokeDelegate(address delegate)` | Anyone | Revoke delegation from an address |
| `resolveRecipient(address addr)` | View | Walk boundTo chain to root, return reward recipient |
| `isRegistered(address addr)` | View | `boundTo[addr] != 0 \|\| recipient[addr] != 0` |

#### Staking (Allocation Only)

| Function | Access | Description |
|----------|--------|-------------|
| `allocate(address staker, address agent, uint256 subnetId, uint256 amount)` | Staker / Delegate | Allocate stake to (agent, subnet) triple; staker is explicit parameter |
| `deallocate(address staker, address agent, uint256 subnetId, uint256 amount)` | Staker / Delegate | Release stake allocation |
| `reallocate(address staker, address fromAgent, uint256 fromSubnetId, address toAgent, uint256 toSubnetId, uint256 amount)` | Staker / Delegate | Move stake between triples (immediate) |

> **Note:** Deposit/withdraw is handled by StakeNFT directly. AWPRegistry only manages allocations.

#### Subnet Lifecycle

| Function | Access | Description |
|----------|--------|-------------|
| `registerSubnet(SubnetParams params)` → `uint256` | Anyone | Register new subnet (CREATE2-deploys Alpha + LP). `params.salt=0` uses subnetId; non-zero enables vanity address. |
| `activateSubnet(uint256 subnetId)` | NFT Owner | Pending → Active |
| `pauseSubnet(uint256 subnetId)` | NFT Owner | Active → Paused |
| `resumeSubnet(uint256 subnetId)` | NFT Owner | Paused → Active |
| `banSubnet(uint256 subnetId)` | Timelock | Active/Paused → Banned |
| `unbanSubnet(uint256 subnetId)` | Timelock | Banned → Active |
| `deregisterSubnet(uint256 subnetId)` | Timelock | Delete subnet (after immunity period) |

#### Governance Parameters

| Function | Access | Description |
|----------|--------|-------------|
| `setInitialAlphaPrice(uint256 price)` | Timelock | Set LP creation price (min 1e12) |
| `setGuardian(address g)` | Timelock | Update guardian address |
| `setImmunityPeriod(uint256 p)` | Timelock | Set deregister immunity period |
| `setSubnetManagerImpl(address impl)` | Timelock | Set/update default SubnetManager impl |

#### View Functions

| Function | Returns | Description |
|----------|---------|-------------|
| ~~`currentEpoch()`~~ | — | Removed from AWPRegistry. Epoch logic now lives in AWPEmission. |
| `getSubnet(uint256 subnetId)` | `SubnetInfo` | Full subnet on-chain data |
| `getActiveSubnetCount()` | `uint256` | Number of active subnets |
| `getActiveSubnetIdAt(uint256 index)` | `uint256` | Active subnet ID by index |
| `isSubnetActive(uint256 subnetId)` | `bool` | Whether subnet is Active |
| `nextSubnetId()` | `uint256` | Next subnet ID to be assigned |
| `getAgentInfo(address agent, uint256 subnetId)` | `AgentInfo` | Agent stake + boundTo + recipient |
| `getAgentsInfo(address[] agents, uint256 subnetId)` | `AgentInfo[]` | Batch agent info query |
| `getRegistry()` | 9 addresses | All module contract addresses (awpToken, subnetNFT, alphaTokenFactory, awpEmission, lpManager, stakingVault, stakeNFT, treasury, guardian) |
| `resolveRecipient(address addr)` | `address` | Walk boundTo chain to root |
| `isRegistered(address addr)` | `bool` | `boundTo[addr] != 0 \|\| recipient[addr] != 0` |
| `boundTo(address addr)` | `address` | Direct binding target |
| `recipient(address addr)` | `address` | Set recipient address |
| `delegates(address user, address delegate)` | `bool` | Whether delegate is authorized |

#### Emergency

| Function | Access | Description |
|----------|--------|-------------|
| `pause()` | Guardian | Emergency pause all operations |
| `unpause()` | Timelock | Resume operations |

---

### 1.2 AWPEmission

> UUPS upgradeable emission engine. Generic address→weight distribution. Oracle multi-sig for weight submission.

#### Oracle Weight Submission

| Function | Access | Description |
|----------|--------|-------------|
| `submitAllocations(address[] recipients, uint96[] weights, bytes[] signatures, uint256 effectiveEpoch)` | Anyone (oracle-signed) | Full-replacement weight submission with EIP-712 multi-sig |

#### Epoch Settlement

| Function | Access | Description |
|----------|--------|-------------|
| `settleEpoch(uint256 limit)` | Anyone | Process up to `limit` recipients per call. Batched 3-phase design. |

#### Governance

| Function | Access | Description |
|----------|--------|-------------|
| `emergencySetWeight(uint256 epoch, uint256 index, address addr, uint96 weight)` | Timelock | Override a single recipient's weight |
| `setOracleConfig(address[] oracles, uint256 threshold)` | Timelock | Configure oracle set and multi-sig threshold |
| `upgradeToAndCall(address newImpl, bytes data)` | Timelock | UUPS upgrade |

#### View Functions

| Function | Returns | Description |
|----------|---------|-------------|
| `settledEpoch()` | `uint256` | Number of epochs settled |
| `activeEpoch()` | `uint256` | Most recently promoted weight epoch |
| `currentDailyEmission()` | `uint256` | Current epoch emission (wei) |
| `settleProgress()` | `uint256` | 0=idle, >0=in progress |
| `epochEmissionLocked()` | `uint256` | Locked emission for current epoch |
| `oracleThreshold()` | `uint256` | Required signature count |
| `allocationNonce()` | `uint256` | Replay protection counter |
| `maxRecipients()` | `uint256` | Maximum recipients allowed |
| `awpRegistry()` | `address` | AWPRegistry contract address |
| `getOracleCount()` | `uint256` | Number of registered oracles |
| `getRecipientCount()` | `uint256` | Number of active-epoch recipients |
| `getRecipient(uint256 index)` | `address` | Recipient by index |
| `getWeight(address addr)` | `uint96` | Weight for address (O(n) scan) |
| `getTotalWeight()` | `uint256` | Total weight in active epoch |
| `getEpochRecipientCount(uint256 epoch)` | `uint256` | Recipient count for a specific epoch |
| `getEpochWeight(uint256 epoch, address addr)` | `uint96` | Weight for address in a specific epoch |
| `getEpochTotalWeight(uint256 epoch)` | `uint256` | Total weight in a specific epoch |
| `oracles(uint256 index)` | `address` | Oracle by index |

---

### 1.3 StakingVault

> Pure allocation logic. No deposit/withdraw/cooldown/STP. Only callable by AWPRegistry.

| Function | Description |
|----------|-------------|
| `allocate(staker, agent, subnetId, amount)` | Allocate stake (onlyAWPRegistry) |
| `deallocate(staker, agent, subnetId, amount)` | Release allocation (onlyAWPRegistry) |
| `reallocate(staker, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)` | Move allocation (onlyAWPRegistry) |
| `freezeAgentAllocations(staker, agent)` | Freeze agent allocations — auto-enumerates subnets (onlyAWPRegistry) |

**View functions:** `userTotalAllocated`, `getAgentStake`, `subnetTotalStake`, `getSubnetTotalStake`, `getAgentSubnets(user, agent) → uint256[]`

---

### 1.3b StakeNFT

> ERC721 position NFT. Users deposit AWP with lock period (timestamp-based). Each position = NFT with (amount, lockEndTime, createdAt). Transferable.

#### Deposit / Withdraw

| Function | Access | Description |
|----------|--------|-------------|
| `deposit(uint256 amount, uint64 lockDuration)` → `uint256 tokenId` | Anyone | Deposit AWP + mint position NFT (lockDuration in seconds) |
| `depositFor(address user, uint256 amount, uint64 lockDuration)` → `uint256 tokenId` | onlyAWPRegistry | Deposit AWP for another address |
| `addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime)` | NFT Owner | Add more AWP to existing position |
| `withdraw(uint256 tokenId)` | NFT Owner | Withdraw after lock expires (burns NFT) |

#### View Functions

| Function | Returns | Description |
|----------|---------|-------------|
| `positions(uint256 tokenId)` | `(uint128 amount, uint64 lockEndTime, uint64 createdAt)` | Position data |
| `getUserTotalStaked(address user)` | `uint256` | O(1) total staked balance |
| `getVotingPower(uint256 tokenId)` | `uint256` | Voting power: amount * sqrt(min(remainingTime, 54 weeks) / 7 days) |
| `getUserVotingPower(address user, uint256[] tokenIds)` | `uint256` | Total voting power for user's NFTs |
| `totalVotingPower()` | `uint256` | Total voting power across all positions |
| `getPositionForVoting(uint256 tokenId)` | `(address owner, uint128 amount, uint64 lockEndTime, uint64 createdAt, uint64 remainingSeconds, uint256 votingPower)` | Position data for voting calculations |
| `remainingTime(uint256 tokenId)` | `uint64` | Remaining lock time in seconds for position |

#### Events

| Event | Parameters |
|-------|-----------|
| `Deposited` | `address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndTime` |
| `PositionIncreased` | `uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndTime` |
| `Withdrawn` | `address indexed user, uint256 indexed tokenId, uint256 amount` |

---

### 1.4 ~~AccessManager~~ (REMOVED)

> AccessManager contract has been removed in Account System V2. All account management (binding, delegation, recipient) is now handled directly by AWPRegistry.

---

### 1.5 AWPToken

> ERC20 + ERC1363 + Votes. 10B MAX_SUPPLY. 200M (2%) minted in constructor; remainder via AWPEmission.

| Function | Access | Description |
|----------|--------|-------------|
| `mint(address to, uint256 amount)` | Minters only | Mint AWP |
| `addMinter(address minter)` | Admin only | Add minter (before renounce) |
| `renounceAdmin()` | Admin only | Permanently lock minter list |
| `burn(uint256 amount)` | Anyone | Burn own tokens |
| `transferAndCall(to, amount, data)` | Anyone | ERC1363 transfer + callback |
| `delegate(address delegatee)` | Anyone | Delegate voting power |

---

### 1.6 AlphaToken

> Standalone ERC20 deployed via CREATE2. 10B MAX_SUPPLY per subnet. Dual minter: admin (AWPRegistry) + subnetMinter. No proxy pattern — no `_disableInitializers()` needed.

| Function | Access | Description |
|----------|--------|-------------|
| `mint(address to, uint256 amount)` | Minters | Mint Alpha (up to 10B, with time-based cap after lock) |
| `setSubnetMinter(address sc)` | Admin | Set subnet as sole minter (one-time, permanent); snapshots `supplyAtLock` and resets `createdAt` |
| `setMinterPaused(address minter, bool paused)` | Admin | Pause/unpause minting (used for ban) |
| `currentMintableLimit()` | View | Current max mintable (since lock) based on elapsed time |

> **Time-cap design:** After `setSubnetMinter`, `supplyAtLock` snapshots the pre-activation supply (excluding admin LP mint) and `createdAt` is reset to `block.timestamp`. Subnet minters can therefore mint immediately after activation — there is no 4-day lockout. The annual cap is `MAX_SUPPLY * elapsed / 365 days` measured from lock time.

### 1.6b AlphaTokenFactory

> Deploys AlphaToken instances via CREATE2. No Clones/EIP-1167 proxy — each token is a standalone contract. Vanity address rules configured at factory deployment (immutable).

| Function | Access | Description |
|----------|--------|-------------|
| `constructor(deployer, vanityRule)` | — | Deploy factory with packed vanity rule (0 = no validation) |
| `setAddresses(awpRegistry)` | Owner | Link to AWPRegistry and renounce ownership (one-time) |
| `deploy(subnetId, name, symbol, admin, salt)` | AWPRegistry | CREATE2-deploy AlphaToken; salt=0 uses subnetId |
| `predictDeployAddress(bytes32 salt)` | View | Predict address for a given salt (standard CREATE2 formula) |

**Vanity rule encoding** (`uint64`, 8 positions packed):

| Byte position | Addresses positions checked |
|---|---|
| bytes [7..4] | First 4 hex chars of address (prefix) |
| bytes [3..0] | Last 4 hex chars of address (suffix) |

Per-position value: `0-9` = digit, `10-15` = lowercase `a-f` (EIP-55 must stay lower), `16-21` = uppercase `A-F` (EIP-55 must be upper), `>=22` = wildcard.

Example: `"A1????cafe"` → `vanityRule = 0x1001FFFF0C0A0F0E`

---

### 1.7 LPManager

> PancakeSwap V4 CL pool creation. Full-range liquidity, permanently locked.

| Function | Access | Description |
|----------|--------|-------------|
| `createPoolAndAddLiquidity(address alphaToken, uint256 awpAmount, uint256 alphaAmount)` → `(bytes32 poolId, uint256 lpTokenId)` | AWPRegistry | Create CL pool + mint full-range LP |

---

### 1.8 SubnetNFT

> ERC721. tokenId = subnetId. Ownership determines subnet control.

| Function | Access | Description |
|----------|--------|-------------|
| `mint(address to, uint256 tokenId, string name_, address subnetManager_, address alphaToken_, uint128 minStake_, string skillsURI_)` | AWPRegistry | Mint on subnet registration |
| `burn(uint256 tokenId)` | AWPRegistry | Burn on deregister |
| Standard ERC721 | Anyone | `transferFrom`, `approve`, `ownerOf`, etc. |

---

### 1.9 AWPDAO

> Custom NFT-based voting. No delegate/checkpoint. Voters submit StakeNFT tokenId arrays.

#### Proposal Creation

| Function | Access | Description |
|----------|--------|-------------|
| `proposeWithTokens(address[] targets, uint256[] values, bytes[] calldatas, string description, uint256[] tokenIds)` | NFT Owner | Create proposal with token-based proposer threshold |

#### Voting

| Function | Access | Description |
|----------|--------|-------------|
| `castVoteWithReasonAndParams(uint256 proposalId, uint8 support, string reason, bytes params)` | NFT Owner | Vote with position NFTs. `params` = `abi.encode(tokenIds)`. Power = amount * sqrt(min(remainingTime, 54 weeks) / 7 days). Note: `castVote()` is blocked. |

#### View Functions

| Function | Returns | Description |
|----------|---------|-------------|
| `hasVotedWithToken(uint256 proposalId, uint256 tokenId)` | `bool` | Whether tokenId has been used to vote |
| `proposalCreatedAt(uint256 proposalId)` | `uint256` | Timestamp when proposal was created |
| `proposalVotes(uint256 proposalId)` | `(uint256 againstVotes, uint256 forVotes, uint256 abstainVotes)` | Current vote tallies |
| `quorum(uint256 timepoint)` | `uint256` | Required quorum for a given timepoint |
| `proposalThreshold()` | `uint256` | Minimum voting power to create a proposal |

> **Anti-manipulation:** Only NFTs with createdAt < proposalCreatedAt can vote (timestamp-based). Per-tokenId double-vote prevention. MAX_WEIGHT_SECONDS = 54 weeks.

---

## 2. REST API

> Base URL: `https://tapi.awp.sh/api`

### 2.1 System

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check → `{"status": "ok"}` |
| GET | `/registry` | All protocol contract addresses (awpRegistry, awpToken, awpEmission, stakingVault, stakeNFT, subnetNFT, lpManager, alphaTokenFactory, dao, treasury). Per-subnet addresses (subnet_contract, alpha_token) are in `/subnets/{id}`. |

### 2.2 Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users?page=1&limit=20` | Paginated user list |
| GET | `/users/count` | Total registered users |
| GET | `/users/{address}` | User detail (balance, agents, reward recipient) |
| GET | `/address/{address}/check` | Check registration status |

### 2.3 Address Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/address/{address}/check` | Check registration status: `{isRegistered, boundTo, recipient}` |

### 2.4 Staking

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/staking/user/{address}/balance` | totalStaked (from stake_positions) + totalAllocated + unallocated |
| GET | `/staking/user/{address}/positions` | User's StakeNFT positions (tokenId, amount, lockEndTime, createdAt) |
| GET | `/staking/user/{address}/allocations?page=1&limit=20` | User's allocations |
| GET | `/staking/user/{address}/frozen` | Frozen allocations |
| GET | `/staking/agent/{agent}/subnet/{subnetId}` | Agent's total stake on subnet |
| GET | `/staking/agent/{agent}/subnets` | Agent's stakes across all subnets |
| GET | `/staking/user/{address}/pending` | Pending operations (returns `[]`) |
| GET | `/staking/subnet/{subnetId}/total` | Subnet total staked |

### 2.5 Subnets

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/subnets?status=Active&page=1&limit=20` | List subnets (filterable by status) |
| GET | `/subnets/{subnetId}` | Subnet detail |
| GET | `/subnets/{subnetId}/earnings?page=1&limit=20` | AWP emission history |
| GET | `/subnets/{subnetId}/skills` | Subnet skills file URI |
| GET | `/subnets/{subnetId}/agents/{agent}` | Agent info on subnet |

> Subnet response includes: `subnet_id`, `owner`, `name`, `symbol`, `subnet_contract`, `skills_uri`, `alpha_token`, `lp_pool`, `status`, `created_at`, `activated_at`, `min_stake`, `immunity_ends_at` (nullable), `burned` (boolean).

### 2.6 Emission [DRAFT]

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/emission/current` | Current epoch, daily emission, total weight (from Redis) |
| GET | `/emission/schedule` | 30/90/365 day emission projections |
| GET | `/emission/epochs?page=1&limit=20` | Epoch settlement history |

### 2.7 Tokens

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tokens/awp` | AWP total supply + max supply (from Redis) |
| GET | `/tokens/alpha/{subnetId}` | Alpha token info (name, symbol, address) |
| GET | `/tokens/alpha/{subnetId}/price` | Alpha price in AWP (from Redis) |

### 2.8 Governance

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/governance/proposals?status=Active&page=1&limit=20` | List proposals |
| GET | `/governance/proposals/{proposalId}` | Proposal detail |
| GET | `/governance/treasury` | Treasury address |

### 2.10 Relay (Gasless Transactions)

> Rate limit: 100 requests per IP per 1 hour (shared across all relay endpoints)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/relay/bind` | Gasless tree-based bind via EIP-712 signature |
| POST | `/relay/set-recipient` | Gasless set recipient via EIP-712 signature |
| POST | `/relay/register-subnet` | Fully gasless subnet registration via ERC-2612 permit + EIP-712 |

**POST /relay/bind request:**
```json
{"user": "0x...", "target": "0x...", "deadline": 1742400000, "signature": "0x...130 hex chars (65 bytes)"}
```

**POST /relay/set-recipient request:**
```json
{"user": "0x...", "recipient": "0x...", "deadline": 1742400000, "signature": "0x...130 hex chars (65 bytes)"}
```

**POST /relay/register-subnet request:**
```json
{
  "user": "0x...", "name": "EVO Alpha", "symbol": "EVO",
  "subnetManager": "0x0000...0000", "salt": "0x...",
  "minStake": "0", "skillsUri": "https://example.com/skills.md",
  "deadline": 1742400000,
  "permitSignature": "0x...65 bytes (ERC-2612 AWP permit)",
  "registerSignature": "0x...65 bytes (EIP-712 registerSubnet)"
}
```

**Response (all three):**
```json
{"txHash": "0x..."}
```

**Errors:**
| Code | Meaning |
|------|---------|
| 400 | Invalid params, expired deadline, bad signature format |
| 429 | Rate limit exceeded |
| 500 | Relay transaction failed |

### 2.11 Vanity Address (Salt Mining)

> Requires `ALPHA_FACTORY_ADDRESS`, `ALPHA_INITCODE_HASH`, and `VANITY_RULE` configured. Uses Foundry `cast create2` for high-speed parallel mining. Pattern is determined by the factory's on-chain `vanityRule` — no request parameters needed.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/vanity/compute-salt` | Compute a CREATE2 salt matching the factory's vanity rule |

**Request:** empty body or `{}`

**Response:**
```json
{"salt": "0x530c11...", "address": "0xA1b275...cafe", "elapsed": "6.998s"}
```

The returned `salt` is passed as `SubnetParams.salt` in `registerSubnet()`.

| Error | Meaning |
|-------|---------|
| 408 | Search timed out (120s) |
| 500 | Mining engine error |

### 2.9 WebSocket

**Endpoint:** `WS /ws/live`

```javascript
const ws = new WebSocket('wss://tapi.awp.sh/ws/live');
ws.send(JSON.stringify({ subscribe: ["RecipientAWPDistributed", "EpochSettled"] }));
ws.onmessage = (e) => console.log(JSON.parse(e.data));
```

**Event format:**
```json
{"type": "RecipientAWPDistributed", "blockNumber": 12345, "txHash": "0x...", "data": {...}}
```

---

## 3. Data Structures

### Solidity

```solidity
enum SubnetStatus { Pending, Active, Paused, Banned }

struct SubnetInfo {
    bytes32 lpPool;          // PancakeSwap V4 PoolId
    SubnetStatus status;
    uint64 createdAt;
    uint64 activatedAt;
}

struct SubnetParams {
    string name;             // 1-64 bytes
    string symbol;           // 1-16 bytes
    address subnetManager;   // address(0) = auto-deploy SubnetManager proxy
    bytes32 salt;            // CREATE2 salt for Alpha token address; bytes32(0) = use subnetId as salt
    uint128 minStake;        // Minimum stake requirement for agents (0 = no minimum)
    string skillsURI;        // Skills file URI (IPFS/HTTPS)
}

struct AgentInfo {
    address boundTo;
    bool isValid;
    uint256 stake;
    address recipient;
}
```

---

## 4. Events

### AWPRegistry Events

| Event | Parameters |
|-------|-----------|
| `Bound` | `address indexed user, address indexed target, address oldTarget` |
| `RecipientUpdated` | `address indexed user, address recipient` |
| `DelegateGranted` | `address indexed user, address indexed delegate` |
| `DelegateRevoked` | `address indexed user, address indexed delegate` |
| `Allocated` | `address indexed user, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator` |
| `Deallocated` | `address indexed user, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator` |
| `Reallocated` | `address indexed user, address fromAgent, uint256 fromSubnet, address toAgent, uint256 toSubnet, uint256 amount, address operator` |
| `SubnetRegistered` | `uint256 indexed subnetId, address indexed owner, string name, string symbol, address subnetManager, address alphaToken` |
| `LPCreated` | `uint256 indexed subnetId, bytes32 poolId, uint256 awpAmount, uint256 alphaAmount` |
| `SubnetActivated` | `uint256 indexed subnetId` |
| `SubnetPaused` | `uint256 indexed subnetId` |
| `SubnetResumed` | `uint256 indexed subnetId` |
| `SubnetBanned` | `uint256 indexed subnetId` |
| `SubnetUnbanned` | `uint256 indexed subnetId` |
| `SubnetDeregistered` | `uint256 indexed subnetId` |
| `GuardianUpdated` | `address indexed oldGuardian, address indexed newGuardian` |
| `InitialAlphaPriceUpdated` | `uint256 newPrice` |
| `ImmunityPeriodUpdated` | `uint256 newPeriod` |
| `AlphaTokenFactoryUpdated` | `address indexed newFactory` |
| `DefaultSubnetManagerImplUpdated` | `address indexed newImpl` |

### StakingVault Events

| Event | Parameters |
|-------|-----------|
| `AgentAllocationsFrozen` | `address indexed user, address indexed agent, uint256 totalFrozen` |

### StakeNFT Events

| Event | Parameters |
|-------|-----------|
| `Deposited` | `address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndTime` |
| `PositionIncreased` | `uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndTime` |
| `Withdrawn` | `address indexed user, uint256 indexed tokenId, uint256 amount` |

### AWPEmission Events

| Event | Parameters |
|-------|-----------|
| `AllocationsSubmitted` | `uint256 indexed nonce, address[] recipients, uint96[] weights` |
| `OracleConfigUpdated` | `address[] oracles, uint256 threshold` |
| `GovernanceWeightUpdated` | `address indexed addr, uint96 weight` |
| `RecipientAWPDistributed` | `uint256 indexed epoch, address indexed recipient, uint256 awpAmount` |
| `DAOMatchDistributed` | `uint256 indexed epoch, uint256 amount` |
| `EpochSettled` | `uint256 indexed epoch, uint256 totalEmission, uint256 recipientCount` |

---

## 5. Error Codes

### AWPRegistry

| Error | Trigger |
|-------|---------|
| `NotDeployer()` | Non-deployer calls initializeRegistry |
| `AlreadyInitialized()` | Registry already initialized |
| `UnknownAddress()` | Invalid Timelock/Guardian caller or unknown updateAddress key |
| `NotDelegate()` | Caller is not a delegate of the staker |
| `CycleDetected()` | Binding would create a cycle in the tree |
| `InvalidSubnetParams()` | name/symbol length invalid |
| `SubnetManagerRequired()` | subnetManager is zero address |
| `NotOwner()` | Non-NFT holder calling lifecycle function |
| `InvalidSubnetStatus()` | Status precondition not met |
| `MaxActiveSubnetsReached()` | Active count >= 10,000 |
| `ImmunityNotExpired()` | Deregister during immunity period |
| `InvalidAgent()` | Agent is not valid |
| `PriceTooLow()` | initialAlphaPrice < 1e12 |
| `PriceTooHigh()` | initialAlphaPrice exceeds maximum |
| `InsufficientMinStake()` | Allocation results in agent stake below subnet minStake |
| `ExpiredSignature()` | Gasless signature expired |
| `InvalidSignature()` | Gasless signature invalid |

### AWPEmission

| Error | Trigger |
|-------|---------|
| `NotTimelock()` | Non-Timelock caller |
| `InvalidRecipient()` | Zero address recipient |
| `InvalidAmount()` | Zero amount |
| `EpochNotReady()` | All epochs up to current time-based epoch have been settled |
| `MiningComplete()` | AWP fully minted |
| `SettlementInProgress()` | Cannot modify during settlement |
| `OracleNotConfigured()` | No oracle set |
| `InvalidOracleConfig()` | Bad threshold or oracle list |
| `InvalidSignatureCount()` | Below threshold or exceeds oracle count |
| `DuplicateOracle()` | Duplicate in oracle list |
| `UnknownOracle()` | Signer not registered oracle |
| `DuplicateSigner()` | Same oracle signed twice |
| `DuplicateRecipient()` | Duplicate in recipient list |
| `ArrayLengthMismatch()` | recipients/weights length differ |
| `InvalidParameter()` | Zero limit, zero epoch duration, etc. |

### StakingVault

| Error | Trigger |
|-------|---------|
| `InsufficientUnallocated()` | Allocate > unallocated balance |
| `InsufficientAllocation()` | Deallocate > available allocation |
| `InvalidAmount()` | Zero amount |

### StakeNFT

| Error | Trigger |
|-------|---------|
| `InvalidAmount()` | Zero deposit amount |
| `LockTooShort()` | Lock period too short |
| `LockNotExpired()` | Withdraw before lock end time |
| `NotTokenOwner()` | Caller does not own the tokenId |
| `InsufficientUnallocated()` | Withdraw exceeds unallocated balance |
| `NothingToUpdate()` | No changes to apply to position |
| `LockCannotShorten()` | New lock end time is earlier than current |
| `LockMustExceedCurrentTime()` | Lock end time must be in the future |
| `NotRootNet()` | Caller is not the AWPRegistry contract |

### API HTTP Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 400 | Bad request |
| 404 | Not found |
| 429 | Rate limit exceeded (relay endpoints) |
| 500 | Internal error |

---

## 6. Constants

| Constant | Value | Location |
|----------|-------|----------|
| AWP MAX_SUPPLY | 10,000,000,000 × 10^18 | AWPToken |
| Alpha MAX_SUPPLY | 10,000,000,000 × 10^18 | AlphaToken (per subnet) |
| INITIAL_ALPHA_MINT | 100,000,000 × 10^18 | AWPRegistry |
| INITIAL_DAILY_EMISSION | 15,800,000 × 10^18 | Deploy.s.sol |
| EPOCH_DURATION | 86,400 (1 day) | AWPEmission (initialized via Deploy.s.sol) |
| DECAY_FACTOR | 996844 / 1000000 | AWPEmission |
| EMISSION_SPLIT_BPS | 5000 (50%) | AWPEmission |
| MAX_ACTIVE_SUBNETS | 10,000 | AWPRegistry |
| maxRecipients | 10,000 | AWPEmission |
| MAX_WEIGHT_SECONDS | 54 weeks (32,659,200s) | StakeNFT / AWPDAO |
| POOL_FEE | 10,000 (1%) | LPManager |
| TICK_SPACING | 200 | LPManager |
| Default immunity | 30 days | AWPRegistry |
| TIMELOCK_DELAY | 172,800 (2 days) | Deploy.s.sol |
