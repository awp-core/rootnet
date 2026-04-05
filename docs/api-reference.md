# AWP — API Reference

## Table of Contents

1. [Smart Contract API](#1-smart-contract-api)
   - [AWPRegistry](#11-awpregistry)
   - [AWPEmission](#12-awpemission)
   - [StakingVault](#13-stakingvault)
   - [StakeNFT](#13b-stakenft)
   - [AWPToken](#14-awptoken)
   - [AlphaToken](#15-alphatoken)
   - [AlphaTokenFactory](#15b-alphatokenfactory)
   - [WorknetNFT](#16-worknetnft)
   - [LPManager](#17-lpmanager)
   - [AWPDAO](#18-awpdao)
2. [REST API](#2-rest-api)
   - [System](#21-system)
   - [Users](#22-users)
   - [Address / Nonce](#23-address--nonce)
   - [Agents](#24-agents)
   - [Staking](#25-staking)
   - [Worknets](#26-worknets)
   - [Emission](#27-emission)
   - [Tokens](#28-tokens)
   - [Governance](#29-governance)
   - [Relay](#210-relay-gasless-transactions)
   - [Vanity](#211-vanity-address-salt-mining)
   - [Announcements](#212-announcements)
   - [Admin](#213-admin)
   - [WebSocket](#214-websocket)
   - [JSON-RPC 2.0](#215-json-rpc-20)
3. [Data Structures](#3-data-structures)
4. [Events](#4-events)
5. [Error Codes](#5-error-codes)
6. [Constants](#6-constants)
7. [Multi-Chain](#7-multi-chain)

---

## 1. Smart Contract API

### 1.1 AWPRegistry

> Unified entry point for worknet management and the account system (UUPS proxy). All worknet lifecycle operations go through AWPRegistry. Allocation/deallocation is handled by StakingVault.

#### Account System (V2)

> EIP-712 domain name: "AWPRegistry". No mandatory registration -- every address is implicitly a root.

| Function | Access | Description |
|----------|--------|-------------|
| `bind(address target)` | Anyone | Tree-based binding with anti-cycle check |
| `unbind()` | Anyone | Remove binding (become a root again) |
| `setRecipient(address addr)` | Anyone | Set reward recipient address |
| `grantDelegate(address delegate)` | Anyone | Grant delegation to an address |
| `revokeDelegate(address delegate)` | Anyone | Revoke delegation from an address |
| `resolveRecipient(address addr)` | View | Walk boundTo chain to root, return reward recipient |
| `batchResolveRecipients(address[] addrs)` | View | Batch resolve recipients |

#### Gasless Account Operations (EIP-712)

| Function | Access | Description |
|----------|--------|-------------|
| `bindFor(address agent, address target, uint256 deadline, uint8 v, bytes32 r, bytes32 s)` | Anyone | Gasless bind via EIP-712 |
| `setRecipientFor(address user, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s)` | Anyone | Gasless set recipient via EIP-712 |
| `grantDelegateFor(address user, address delegate, uint256 deadline, uint8 v, bytes32 r, bytes32 s)` | Anyone | Gasless grant delegate via EIP-712 |
| `revokeDelegateFor(address user, address delegate, uint256 deadline, uint8 v, bytes32 r, bytes32 s)` | Anyone | Gasless revoke delegate via EIP-712 |
| `unbindFor(address user, uint256 deadline, uint8 v, bytes32 r, bytes32 s)` | Anyone | Gasless unbind via EIP-712 |

#### Worknet Lifecycle

| Function | Access | Description |
|----------|--------|-------------|
| `registerWorknet(WorknetParams params)` -> `uint256` | Anyone | Register new worknet (CREATE2-deploys Alpha + LP). `params.salt=0` uses worknetId; non-zero enables vanity address. |
| `registerWorknetFor(...)` | Anyone | Gasless EIP-712 worknet registration (requires prior AWP approve) |
| `registerWorknetForWithPermit(...)` | Anyone | Fully gasless -- ERC-2612 permit + EIP-712 in one tx (zero gas for user) |
| `activateWorknet(uint256 worknetId)` | NFT Owner | Pending -> Active |
| `activateWorknetFor(uint256 worknetId, uint256 deadline, uint8 v, bytes32 r, bytes32 s)` | Anyone | Gasless activation via EIP-712 |
| `pauseWorknet(uint256 worknetId)` | NFT Owner | Active -> Paused |
| `resumeWorknet(uint256 worknetId)` | NFT Owner | Paused -> Active |
| `banWorknet(uint256 worknetId)` | Timelock | Active/Paused -> Banned |
| `unbanWorknet(uint256 worknetId)` | Timelock | Banned -> Active (checks MAX_ACTIVE_WORKNETS) |
| `deregisterWorknet(uint256 worknetId)` | Timelock | Delete worknet (after immunity period) |

#### Governance Parameters

| Function | Access | Description |
|----------|--------|-------------|
| `setInitialAlphaPrice(uint256 price)` | Timelock | Set LP creation price (min 1e12) |
| `setInitialAlphaMint(uint256 amount)` | Timelock | Set initial Alpha mint amount for LP |
| `setGuardian(address g)` | Guardian | Update guardian address |
| `setImmunityPeriod(uint256 p)` | Timelock | Set deregister immunity period |
| `setWorknetManagerImpl(address impl)` | Timelock | Set/update default WorknetManager impl |
| `setAlphaTokenFactory(address factory)` | Timelock | Replace AlphaTokenFactory |
| `setLPManager(address lpManager_)` | Timelock | Replace LPManager |
| `setDexConfig(bytes dexConfig_)` | Timelock | Update DEX configuration |

#### View Functions

| Function | Returns | Description |
|----------|---------|-------------|
| `getWorknet(uint256 worknetId)` | `WorknetInfo` | Worknet lifecycle state |
| `getActiveWorknetCount()` | `uint256` | Number of active worknets |
| `getActiveWorknetIdAt(uint256 index)` | `uint256` | Active worknet ID by index |
| `isWorknetActive(uint256 worknetId)` | `bool` | Whether worknet is Active |
| `nextWorknetId()` | `uint256` | Next worknet ID to be assigned |
| `getRegistry()` | 9 addresses | All module contract addresses (awpToken, worknetNFT, alphaTokenFactory, awpEmission, lpManager, stakingVault, stakeNFT, treasury, guardian) |
| `resolveRecipient(address addr)` | `address` | Walk boundTo chain to root |
| `batchResolveRecipients(address[] addrs)` | `address[]` | Batch resolve |
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

> UUPS upgradeable emission engine. Generic address->weight distribution. Guardian (cross-chain multisig) submits epoch-versioned packed allocations.

#### Allocation Management (Guardian only)

| Function | Access | Description |
|----------|--------|-------------|
| `submitAllocations(uint256[] packed, uint256 totalWeight, uint256 effectiveEpoch)` | Guardian | Full-replacement weight submission. Packed format: `[32-bit zero \| 64-bit weight \| 160-bit address]` |
| `appendAllocations(uint256[] packed, uint256 effectiveEpoch)` | Guardian | Add new recipients to existing allocations |
| `modifyAllocations(uint256[] patches, uint256 newTotalWeight, uint256 effectiveEpoch)` | Guardian | Modify weights for existing recipients |

#### Epoch Settlement

| Function | Access | Description |
|----------|--------|-------------|
| `settleEpoch(uint256 limit)` | Anyone | Process up to `limit` recipients per call. Uses mint + best-effort ERC1363 callback. |

#### Guardian Configuration

| Function | Access | Description |
|----------|--------|-------------|
| `setDecayFactor(uint256 newDecayFactor)` | Guardian | Update emission decay factor |
| `setMaxRecipients(uint256 newMax)` | Guardian | Update max recipients |
| `setEpochDuration(uint256 newDuration)` | Guardian | Update epoch duration |
| `setTreasury(address t)` | Guardian | Update treasury address |
| `setGuardian(address g)` | Guardian | Transfer guardian role |
| `pauseEpochUntil(uint64 resumeTime)` | Guardian | Pause emission until a future time |

#### View Functions

| Function | Returns | Description |
|----------|---------|-------------|
| `currentEpoch()` | `uint256` | Current time-based epoch: `(block.timestamp - baseTime) / epochDuration + baseEpoch` |
| `settledEpoch()` | `uint256` | Next epoch to settle |
| `activeEpoch()` | `uint256` | Most recently promoted weight epoch |
| `currentDailyEmission()` | `uint256` | Current epoch emission (wei) |
| `settleProgress()` | `uint256` | 0=idle, >0=in progress |
| `epochEmissionLocked()` | `uint256` | Locked emission for current settlement |
| `maxRecipients()` | `uint256` | Maximum recipients allowed |
| `decayFactor()` | `uint256` | Current decay factor |
| `epochDuration()` | `uint256` | Epoch length in seconds |
| `baseTime()` | `uint256` | Genesis or rebased time |
| `baseEpoch()` | `uint256` | Epoch number at baseTime |
| `guardian()` | `address` | Guardian address |
| `treasury()` | `address` | Treasury address |
| `pausedUntil()` | `uint64` | Pause resume timestamp (0 = not paused) |
| `frozenEpoch()` | `uint64` | Epoch frozen at pause time |
| `getRecipientCount()` | `uint256` | Number of active-epoch recipients |
| `getRecipient(uint256 index)` | `address` | Recipient by index |
| `getWeight(address addr)` | `uint96` | Weight for address |
| `getTotalWeight()` | `uint256` | Total weight in active epoch |
| `getEpochRecipientCount(uint256 epoch)` | `uint256` | Recipient count for a specific epoch |
| `getEpochWeight(uint256 epoch, address addr)` | `uint96` | Weight for address in a specific epoch |
| `getEpochTotalWeight(uint256 epoch)` | `uint256` | Total weight in a specific epoch |

---

### 1.3 StakingVault

> UUPS proxy with EIP-712 gasless support. Manages allocations only -- deposit/withdraw is handled by StakeNFT. Auth: caller must be staker or delegate (reads AWPRegistry.delegates cross-contract).

> EIP-712 domain name: "StakingVault"

#### Write Functions

| Function | Access | Description |
|----------|--------|-------------|
| `allocate(address staker, address agent, uint256 worknetId, uint256 amount)` | Staker / Delegate | Allocate stake to (staker, agent, worknetId) triple. worknetId=0 rejected. |
| `deallocate(address staker, address agent, uint256 worknetId, uint256 amount)` | Staker / Delegate | Release stake allocation |
| `reallocate(address staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount)` | Staker / Delegate | Move stake between triples (immediate) |

#### Gasless (EIP-712)

| Function | Access | Description |
|----------|--------|-------------|
| `allocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s)` | Anyone | Gasless allocate via EIP-712 signature |
| `deallocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s)` | Anyone | Gasless deallocate via EIP-712 signature |

#### View Functions

| Function | Returns | Description |
|----------|---------|-------------|
| `nonces(address user)` | `uint256` | EIP-712 nonce for user |
| `userTotalAllocated(address user)` | `uint256` | Total allocated by user |
| `worknetTotalStake(uint256 worknetId)` | `uint256` | Total stake on worknet |
| `getAgentStake(address user, address agent, uint256 worknetId)` | `uint256` | Stake for (user, agent, worknet) triple |
| `getAgentWorknets(address user, address agent)` | `uint256[]` | All worknetIds for a (user, agent) pair |
| `getWorknetTotalStake(uint256 worknetId)` | `uint256` | Total stake on worknet (alias) |

---

### 1.3b StakeNFT

> ERC721 position NFT (not Enumerable). Users deposit AWP with lock period (timestamp-based). Each position = NFT with (amount, lockEndTime, createdAt). Transferable.

#### Deposit / Withdraw

| Function | Access | Description |
|----------|--------|-------------|
| `deposit(uint256 amount, uint64 lockDuration)` -> `uint256 tokenId` | Anyone | Deposit AWP + mint position NFT (lockDuration in seconds) |
| `depositFor(address user, uint256 amount, uint64 lockDuration)` -> `uint256 tokenId` | onlyAWPRegistry | Deposit AWP for another address |
| `addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime)` | NFT Owner | Add more AWP to existing position (blocked if lock expired -- PositionExpired) |
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
| `remainingTime(uint256 tokenId)` | `uint64` | Remaining lock time in seconds |

---

### 1.4 AWPToken

> ERC20 + ERC1363 + Votes. 10B MAX_SUPPLY. Initial mint configurable per chain (INITIAL_MINT constructor param, immutable); remainder via AWPEmission.

| Function | Access | Description |
|----------|--------|-------------|
| `mint(address to, uint256 amount)` | Minters only | Mint AWP |
| `mintAndCall(address to, uint256 amount, bytes data)` | Minters only | Mint + ERC1363 callback on recipient |
| `addMinter(address minter)` | Admin only | Add minter (before renounce) |
| `renounceAdmin()` | Admin only | Permanently lock minter list |
| `burn(uint256 amount)` | Anyone | Burn own tokens |
| `transferAndCall(to, amount, data)` | Anyone | ERC1363 transfer + callback |
| `delegate(address delegatee)` | Anyone | Delegate voting power |

---

### 1.5 AlphaToken

> Standalone ERC20 deployed via CREATE2. 10B MAX_SUPPLY per worknet. Dual minter: admin (AWPRegistry) + worknetMinter. No proxy pattern.

| Function | Access | Description |
|----------|--------|-------------|
| `initialize(string name, string symbol, uint256 worknetId, address admin)` | Factory | Initialize after deployment (one-time) |
| `mint(address to, uint256 amount)` | Minters | Mint Alpha (up to 10B, with time-based cap after lock) |
| `setWorknetMinter(address worknetManager)` | Admin | Set worknet as sole minter (one-time, permanent); snapshots `supplyAtLock` and resets `createdAt` |
| `setMinterPaused(address minter, bool paused)` | Admin | Pause/unpause minting (used for ban) |
| `burn(uint256 amount)` | Anyone | Burn own tokens |
| `transferAndCall(address to, uint256 amount, bytes data)` | Anyone | ERC1363 transfer + callback |
| `approveAndCall(address spender, uint256 amount, bytes data)` | Anyone | ERC1363 approve + callback |
| `currentMintableLimit()` | View | Current max mintable (since lock) based on elapsed time |

> **Time-cap design:** After `setWorknetMinter`, `supplyAtLock` snapshots the pre-activation supply and `createdAt` is reset to `block.timestamp`. Worknet minters can mint immediately after activation. The annual cap is `MAX_SUPPLY * elapsed / 365 days` measured from lock time.

### 1.5b AlphaTokenFactory

> Deploys AlphaToken instances via CREATE2. No Clones/EIP-1167 proxy -- each token is a standalone contract. Vanity address rules configured at factory deployment (immutable).

| Function | Access | Description |
|----------|--------|-------------|
| `constructor(deployer, vanityRule)` | -- | Deploy factory with packed vanity rule (0 = no validation) |
| `setAddresses(awpRegistry)` | Owner | Link to AWPRegistry and renounce ownership (one-time) |
| `deploy(worknetId, name, symbol, admin, salt)` | AWPRegistry | CREATE2-deploy AlphaToken; salt=0 uses worknetId |
| `predictDeployAddress(bytes32 salt)` | View | Predict address for a given salt (standard CREATE2 formula) |

**Vanity rule encoding** (`uint64`, 8 positions packed):

| Byte position | Address positions checked |
|---|---|
| bytes [7..4] | First 4 hex chars of address (prefix) |
| bytes [3..0] | Last 4 hex chars of address (suffix) |

Per-position value: `0-9` = digit, `10-15` = lowercase `a-f` (EIP-55 must stay lower), `16-21` = uppercase `A-F` (EIP-55 must be upper), `>=22` = wildcard.

Example: `"A1????cafe"` -> `vanityRule = 0x1001FFFF0C0A0F0E`

---

### 1.6 WorknetNFT

> ERC721 with on-chain identity storage. tokenId = worknetId. Ownership determines worknet control.

| Function | Access | Description |
|----------|--------|-------------|
| `mint(address to, uint256 tokenId, string name_, address worknetManager_, address alphaToken_, uint128 minStake_, string skillsURI_)` | AWPRegistry | Mint on worknet registration |
| `burn(uint256 tokenId)` | AWPRegistry | Burn on deregister |
| `setSkillsURI(uint256 tokenId, string skillsURI_)` | NFT Owner | Update skills URI |
| `setMinStake(uint256 tokenId, uint128 minStake_)` | NFT Owner | Update minimum stake |
| `setBaseURI(string uri)` | Owner | Set base URI for tokenURI fallback |
| Standard ERC721 | Anyone | `transferFrom`, `approve`, `ownerOf`, etc. |

#### View Functions

| Function | Returns | Description |
|----------|---------|-------------|
| `getWorknetManager(uint256 tokenId)` | `address` | Worknet manager contract |
| `getAlphaToken(uint256 tokenId)` | `address` | Alpha token contract |
| `getMinStake(uint256 tokenId)` | `uint128` | Minimum stake requirement |
| `getWorknetData(uint256 tokenId)` | `WorknetData` | Full on-chain identity (name, worknetManager, alphaToken, skillsURI, minStake, owner) |

> **tokenURI resolution:** 3-tier: per-token metadataURI -> global baseURI -> on-chain Base64 JSON.

---

### 1.7 LPManager

> DEX pool creation (Uniswap V4 CL on Base/ETH/ARB, PancakeSwap V4 CL on BSC). Full-range liquidity, permanently locked.

| Function | Access | Description |
|----------|--------|-------------|
| `createPoolAndAddLiquidity(address alphaToken, uint256 awpAmount, uint256 alphaAmount)` -> `(bytes32 poolId, uint256 lpTokenId)` | AWPRegistry | Create CL pool + mint full-range LP |
| `compoundFees(address alphaToken)` | AWPRegistry | Reinvest accumulated LP fees (called by Keeper cron) |

---

### 1.8 AWPDAO

> Custom NFT-based voting. No delegate/checkpoint. Voters submit StakeNFT tokenId arrays. No awpRegistry dependency.

#### Proposal Creation

| Function | Access | Description |
|----------|--------|-------------|
| `proposeWithTokens(address[] targets, uint256[] values, bytes[] calldatas, string description, uint256[] tokenIds)` | NFT Owner | Create executable proposal (via Timelock). `totalVotingPower > 0` required. |
| `signalPropose(string description, uint256[] tokenIds)` | NFT Owner | Create vote-only (signal) proposal |

> `propose()` is blocked -- use `proposeWithTokens` or `signalPropose`.

#### Voting

| Function | Access | Description |
|----------|--------|-------------|
| `castVoteWithReasonAndParams(uint256 proposalId, uint8 support, string reason, bytes params)` | NFT Owner | Vote with position NFTs. `params` = `abi.encode(tokenIds)`. Power = amount * sqrt(min(remainingTime, 54 weeks) / 7 days). |

> `castVote()` is blocked -- use `castVoteWithReasonAndParams`.

#### View Functions

| Function | Returns | Description |
|----------|---------|-------------|
| `hasVotedWithToken(uint256 proposalId, uint256 tokenId)` | `bool` | Whether tokenId has been used to vote |
| `proposalCreatedAt(uint256 proposalId)` | `uint256` | Timestamp when proposal was created |
| `proposalVotes(uint256 proposalId)` | `(uint256 againstVotes, uint256 forVotes, uint256 abstainVotes)` | Current vote tallies |
| `quorum(uint256 timepoint)` | `uint256` | Required quorum for a given timepoint |
| `proposalThreshold()` | `uint256` | Minimum voting power to create a proposal |

> **Anti-manipulation:** Only NFTs with createdAt < proposalCreatedAt can vote (strict: >= blocks same-block mint+vote). Per-tokenId double-vote prevention. MAX_WEIGHT_SECONDS = 54 weeks.

---

## 2. REST API

> Base URL: `https://tapi.awp.sh/api`
>
> All endpoints accept an optional `?chain_id=` query parameter to target a specific chain (default: primary chain). See [Multi-Chain](#7-multi-chain).

### 2.1 System

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check -> `{"status": "ok"}` |
| GET | `/health/detailed` | Detailed health with DB, Redis, chain connectivity |
| GET | `/registry` | All protocol contract addresses plus `chainId` and `eip712Domain` (AWPRegistry) and `stakingVaultEip712Domain` (StakingVault). Excludes implementation contracts. |
| GET | `/chains` | List all configured chains with their contract addresses |
| GET | `/stats` | Global protocol statistics (cross-chain aggregated) |

> **EIP-712 Domains:** `/api/registry` returns two EIP-712 domains: `eip712Domain` (AWPRegistry, name "AWPRegistry" -- for bind/setRecipient/registerWorknet) and `stakingVaultEip712Domain` (StakingVault, name "StakingVault" -- for allocate/deallocate). Use the correct domain for each operation.

### 2.2 Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users?page=1&limit=20` | Paginated user list |
| GET | `/users/count` | Total registered users |
| GET | `/users/global` | Cross-chain aggregated user list |
| GET | `/users/{address}` | User detail (balance, bound_to, recipient) |
| GET | `/users/{address}/portfolio` | User portfolio (positions, allocations, rewards) |
| GET | `/users/{address}/delegates` | User's granted delegates |

### 2.3 Address / Nonce

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/address/{address}/check` | Check registration status |
| GET | `/address/{address}/resolve-recipient` | Resolve reward recipient for address |
| POST | `/address/batch-resolve-recipients` | Batch resolve recipients (JSON array of addresses) |
| GET | `/nonce/{address}` | AWPRegistry EIP-712 nonce (for bind/setRecipient/registerWorknet relay signatures) |
| GET | `/staking-nonce/{address}` | StakingVault EIP-712 nonce (for allocate/deallocate relay signatures) |

### 2.4 Agents

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/agents/by-owner/{owner}` | List all agents owned by an address |
| GET | `/agents/by-owner/{owner}/{agent}` | Get specific agent detail for an owner |
| GET | `/agents/lookup/{agent}` | Lookup agent by agent address |
| POST | `/agents/batch-info` | Batch query agent info (accepts JSON array of agent addresses) |

### 2.5 Staking

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/staking/user/{address}/balance` | totalStaked (from stake_positions) + totalAllocated + unallocated |
| GET | `/staking/user/{address}/balance/global` | Cross-chain aggregated balance |
| GET | `/staking/user/{address}/positions` | User's StakeNFT positions (tokenId, amount, lockEndTime, createdAt) |
| GET | `/staking/user/{address}/positions/global` | Cross-chain aggregated positions |
| GET | `/staking/user/{address}/allocations?page=1&limit=20` | User's allocations |
| GET | `/staking/user/{address}/pending` | Pending operations (returns `[]`) |
| GET | `/staking/user/{address}/frozen` | Frozen allocations |
| GET | `/staking/agent/{agent}/subnet/{worknetId}` | Agent's total stake on worknet |
| GET | `/staking/agent/{agent}/subnets` | Agent's stakes across all worknets |
| GET | `/staking/subnet/{worknetId}/total` | Worknet total staked |

### 2.6 Worknets

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/subnets` | List worknets (filterable by `?status=Active&page=1&limit=20`) |
| GET | `/subnets/ranked` | List worknets ranked by total stake |
| GET | `/subnets/search?q=keyword` | Search worknets by name/symbol |
| GET | `/subnets/by-owner/{owner}` | List worknets owned by an address |
| GET | `/subnets/{worknetId}` | Worknet detail |
| GET | `/subnets/{worknetId}/skills` | Worknet skills file URI |
| GET | `/subnets/{worknetId}/earnings?page=1&limit=20` | AWP emission history |
| GET | `/subnets/{worknetId}/agents` | List agents on worknet |
| GET | `/subnets/{worknetId}/agents/{agent}` | Agent info on worknet |

> Worknet response includes: `worknet_id`, `owner`, `name`, `symbol`, `worknet_manager`, `skills_uri`, `alpha_token`, `lp_pool`, `status`, `created_at`, `activated_at`, `min_stake`, `immunity_ends_at` (nullable), `burned` (boolean).

### 2.7 Emission

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/emission/current` | Current epoch, daily emission, total weight (from Redis) |
| GET | `/emission/schedule` | 30/90/365 day emission projections (single chain) |
| GET | `/emission/global-schedule` | Cross-chain aggregated emission schedule |
| GET | `/emission/epochs?page=1&limit=20` | Epoch settlement history |
| GET | `/emission/epochs/{epochId}` | Epoch detail (recipients, weights, amounts) |

### 2.8 Tokens

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tokens/awp` | AWP total supply + max supply (from Redis) |
| GET | `/tokens/awp/global` | Cross-chain aggregated AWP info |
| GET | `/tokens/alpha/{worknetId}` | Alpha token info (name, symbol, address) |
| GET | `/tokens/alpha/{worknetId}/price` | Alpha price in AWP (from Redis) |

### 2.9 Governance

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/governance/proposals?status=Active&page=1&limit=20` | List proposals |
| GET | `/governance/proposals/global` | Cross-chain aggregated proposals |
| GET | `/governance/proposals/{proposalId}` | Proposal detail |
| GET | `/governance/treasury` | Treasury address |

### 2.10 Relay (Gasless Transactions)

> Rate limit: configurable per endpoint via Redis (default 100 requests per IP per hour, shared across all relay endpoints). Hot-updatable via admin API.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/relay/bind` | Gasless tree-based bind via EIP-712 signature |
| POST | `/relay/unbind` | Gasless unbind via EIP-712 signature |
| POST | `/relay/set-recipient` | Gasless set recipient via EIP-712 signature |
| POST | `/relay/register` | Gasless user registration via EIP-712 signature |
| POST | `/relay/grant-delegate` | Gasless grant delegate via EIP-712 signature |
| POST | `/relay/revoke-delegate` | Gasless revoke delegate via EIP-712 signature |
| POST | `/relay/allocate` | Gasless stake allocation via EIP-712 signature (StakingVault domain) |
| POST | `/relay/deallocate` | Gasless stake deallocation via EIP-712 signature (StakingVault domain) |
| POST | `/relay/activate-subnet` | Gasless worknet activation via EIP-712 signature |
| POST | `/relay/register-subnet` | Fully gasless worknet registration via ERC-2612 permit + EIP-712 |
| GET | `/relay/status/{txHash}` | Check relay transaction status |

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
  "worknetManager": "0x0000...0000", "salt": "0x...",
  "minStake": "0", "skillsUri": "https://example.com/skills.md",
  "deadline": 1742400000,
  "permitSignature": "0x...65 bytes (ERC-2612 AWP permit)",
  "registerSignature": "0x...65 bytes (EIP-712 registerWorknet)"
}
```

**POST /relay/allocate request:**
```json
{"staker": "0x...", "agent": "0x...", "worknetId": "123", "amount": "1000000000000000000", "deadline": 1742400000, "signature": "0x...130 hex chars (65 bytes)"}
```

**Response (all relay endpoints):**
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

> Requires `ALPHA_FACTORY_ADDRESS`, `ALPHA_INITCODE_HASH`, and `VANITY_RULE` configured. Uses Foundry `cast create2` for high-speed parallel mining.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/vanity/mining-params` | Get current factory address, initcode hash, and vanity rule |
| POST | `/vanity/upload-salts` | Upload pre-computed salts to the DB pool (rate limited: 5/hr/IP) |
| GET | `/vanity/salts` | List available salts in the pool |
| GET | `/vanity/salts/count` | Count available salts |
| POST | `/vanity/compute-salt` | Compute a CREATE2 salt matching the factory's vanity rule (rate limited: 20/hr/IP) |

**POST /vanity/compute-salt response:**
```json
{"salt": "0x530c11...", "address": "0xA1b275...cafe", "elapsed": "6.998s"}
```

The returned `salt` is passed as `WorknetParams.salt` in `registerWorknet()`.

| Error | Meaning |
|-------|---------|
| 408 | Search timed out (120s) |
| 500 | Mining engine error |

### 2.12 Announcements

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/announcements?chain_id=&category=&limit=20` | List active announcements |
| GET | `/announcements/{id}` | Get announcement by ID |
| GET | `/announcements/llm-context` | Get announcements formatted for LLM context |

Admin announcement endpoints are under `/api/admin/announcements` (see [Admin](#213-admin)).

### 2.13 Admin

> All admin endpoints require `Authorization: Bearer <token>` header.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/admin/chains` | Add a new chain configuration |
| DELETE | `/admin/chains/{chainId}` | Remove a chain configuration |
| GET | `/admin/chains` | List all chain configurations |
| PUT | `/admin/ratelimit` | Hot-update rate limit configuration |
| GET | `/admin/ratelimit` | Get current rate limit configuration |
| GET | `/admin/system` | System info (uptime, DB pool, Redis, chains) |
| POST | `/admin/announcements` | Create announcement |
| PUT | `/admin/announcements/{id}` | Update announcement |
| DELETE | `/admin/announcements/{id}` | Delete announcement |

### 2.14 WebSocket

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

### 2.15 JSON-RPC 2.0

**Endpoint:** `GET /v2` (method discovery), `POST /v2` (execute RPC call)

Standard JSON-RPC 2.0 format. All REST endpoints are mirrored as RPC methods (e.g., `users.get`, `subnets.listRanked`, `staking.getPositionsGlobal`).

---

## 3. Data Structures

### Solidity

```solidity
enum WorknetStatus { Pending, Active, Paused, Banned }

struct WorknetInfo {
    bytes32 lpPool;          // DEX LP pool ID
    WorknetStatus status;
    uint64 createdAt;
    uint64 activatedAt;
}

struct WorknetFullInfo {
    address worknetManager;
    address alphaToken;
    bytes32 lpPool;
    WorknetStatus status;
    uint64 createdAt;
    uint64 activatedAt;
    string name;
    string skillsURI;
    uint128 minStake;
    address owner;
}

struct WorknetParams {
    string name;             // 1-64 bytes
    string symbol;           // 1-16 bytes
    address worknetManager;  // address(0) = auto-deploy WorknetManager proxy
    bytes32 salt;            // CREATE2 salt for Alpha token address; bytes32(0) = use worknetId as salt
    uint128 minStake;        // Minimum stake requirement for agents (0 = no minimum)
    string skillsURI;        // Skills file URI (IPFS/HTTPS)
}

// WorknetNFT on-chain identity
struct WorknetData {
    string name;
    address worknetManager;
    address alphaToken;
    string skillsURI;
    uint128 minStake;
    address owner;
}
```

---

## 4. Events

### AWPRegistry Events

| Event | Parameters |
|-------|-----------|
| `UserRegistered` | `address indexed user` |
| `Bound` | `address indexed addr, address indexed target` |
| `Unbound` | `address indexed addr` |
| `RecipientSet` | `address indexed addr, address recipient` |
| `DelegateGranted` | `address indexed staker, address indexed delegate` |
| `DelegateRevoked` | `address indexed staker, address indexed delegate` |
| `WorknetRegistered` | `uint256 indexed worknetId, address indexed owner, string name, string symbol, address worknetManager, address alphaToken` |
| `LPCreated` | `uint256 indexed worknetId, bytes32 poolId, uint256 awpAmount, uint256 alphaAmount` |
| `WorknetActivated` | `uint256 indexed worknetId` |
| `WorknetPaused` | `uint256 indexed worknetId` |
| `WorknetResumed` | `uint256 indexed worknetId` |
| `WorknetBanned` | `uint256 indexed worknetId` |
| `WorknetUnbanned` | `uint256 indexed worknetId` |
| `WorknetDeregistered` | `uint256 indexed worknetId` |
| `GuardianUpdated` | `address indexed newGuardian` |
| `InitialAlphaPriceUpdated` | `uint256 newPrice` |
| `InitialAlphaMintUpdated` | `uint256 amount` |
| `ImmunityPeriodUpdated` | `uint256 newPeriod` |
| `AlphaTokenFactoryUpdated` | `address indexed newFactory` |
| `DefaultWorknetManagerImplUpdated` | `address indexed newImpl` |
| `DexConfigUpdated` | (no params) |
| `LPManagerUpdated` | `address indexed newLPManager` |

### StakingVault Events

| Event | Parameters |
|-------|-----------|
| `Allocated` | `address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator` |
| `Deallocated` | `address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator` |
| `Reallocated` | `address indexed staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount, address operator` |

### StakeNFT Events

| Event | Parameters |
|-------|-----------|
| `Deposited` | `address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndTime` |
| `PositionIncreased` | `uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndTime` |
| `Withdrawn` | `address indexed user, uint256 indexed tokenId, uint256 amount` |

### AWPEmission Events

| Event | Parameters |
|-------|-----------|
| `AllocationsSubmitted` | `uint256 indexed effectiveEpoch, uint256[] packed, uint256 totalWeight` |
| `AllocationsAppended` | `uint256 indexed effectiveEpoch, uint256[] packed` |
| `AllocationsModified` | `uint256 indexed effectiveEpoch, uint256[] patches, uint256 newTotalWeight` |
| `RecipientAWPDistributed` | `uint256 indexed epoch, address indexed recipient, uint256 awpAmount` |
| `EpochSettled` | `uint256 indexed epoch, uint256 totalEmission, uint256 recipientCount` |
| `EpochDurationUpdated` | `uint256 oldDuration, uint256 newDuration` |
| `EpochPausedUntil` | `uint64 resumeTime, uint64 frozenEpoch` |
| `GuardianUpdated` | `address indexed newGuardian` |
| `TreasuryUpdated` | `address indexed newTreasury` |
| `MaxRecipientsUpdated` | `uint256 newMax` |
| `DecayFactorUpdated` | `uint256 newDecayFactor` |

### WorknetNFT Events

| Event | Parameters |
|-------|-----------|
| `SkillsURIUpdated` | `uint256 indexed tokenId, string skillsURI` |
| `MinStakeUpdated` | `uint256 indexed tokenId, uint128 minStake` |
| `MetadataURIUpdated` | `uint256 indexed tokenId, string metadataURI` |

---

## 5. Error Codes

### AWPRegistry

| Error | Trigger |
|-------|---------|
| `NotDeployer()` | Non-deployer calls initializeRegistry |
| `AlreadyInitialized()` | Registry already initialized |
| `ChainTooLong()` | Binding chain exceeds maximum depth (256) |
| `CannotRevokeSelf()` | Cannot revoke self as delegate |
| `InvalidAddress()` | Zero or invalid address provided |
| `NotAuthorized()` | Caller is not authorized (not staker or delegate) |
| `CycleDetected()` | Binding would create a cycle in the tree |
| `InvalidWorknetParams()` | name/symbol length or character invalid |
| `WorknetManagerRequired()` | worknetManager is zero and auto-deploy not available |
| `NotOwner()` | Non-NFT holder calling lifecycle function |
| `InvalidWorknetStatus()` | Status precondition not met |
| `MaxActiveWorknetsReached()` | Active count >= 10,000 |
| `ImmunityNotExpired()` | Deregister during immunity period |
| `PriceTooLow()` | initialAlphaPrice < 1e12 |
| `PriceTooHigh()` | initialAlphaPrice exceeds maximum |
| `ExpiredSignature()` | Gasless signature expired |
| `InvalidSignature()` | Gasless signature invalid |
| `SelfBindNotAllowed()` | Cannot bind to self |

### AWPEmission

| Error | Trigger |
|-------|---------|
| `NotGuardian()` | Non-Guardian caller |
| `InvalidRecipient()` | Zero address recipient |
| `EpochNotReady()` | All epochs up to current time-based epoch have been settled |
| `MiningComplete()` | AWP fully minted (MAX_SUPPLY reached) |
| `SettlementInProgress()` | Cannot modify during settlement |
| `ZeroWeightNotAllowed()` | Weight=0 entry in submission |
| `TooManyRecipients()` | Recipient count exceeds maxRecipients |
| `InvalidDecayFactor()` | Bad decay factor value |
| `LimitCannotBeZero()` | Zero limit in settleEpoch |
| `GenesisTimeNotReached()` | Block time before genesis |
| `EpochDurationCannotBeZero()` | Zero epoch duration |
| `MustBeFutureEpoch()` | effectiveEpoch is in the past |

### StakingVault

| Error | Trigger |
|-------|---------|
| `InsufficientAllocation()` | Deallocate > available allocation |
| `AmountCannotBeZero()` | Zero amount |
| `AmountExceedsUint128()` | Amount overflows uint128 |
| `WorknetIdCannotBeZero()` | worknetId=0 |
| `AllocationOverflow()` | Allocation would overflow uint128 |
| `AlreadySet()` | StakeNFT already configured |
| `ExpiredSignature()` | Gasless signature expired |
| `InvalidSignature()` | Gasless signature invalid |

### StakeNFT

| Error | Trigger |
|-------|---------|
| `InvalidAmount()` | Zero deposit amount |
| `LockTooShort()` | Lock period too short |
| `LockNotExpired()` | Withdraw before lock end time |
| `NotTokenOwner()` | Caller does not own the tokenId |
| `InsufficientUnallocated()` | Withdraw exceeds unallocated balance |
| `NothingToUpdate()` | No changes to apply to position |
| `PositionExpired()` | addToPosition blocked on expired lock |
| `LockCannotShorten()` | New lock end time is earlier than current |
| `LockMustExceedCurrentTime()` | Lock end time must be in the future |
| `NotAWPRegistry()` | Caller is not the AWPRegistry contract |

### API HTTP Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 400 | Bad request |
| 404 | Not found |
| 408 | Timeout (vanity computation) |
| 429 | Rate limit exceeded (relay/vanity endpoints) |
| 500 | Internal error |
| 503 | Service unavailable (keeper cache not ready) |

---

## 6. Constants

| Constant | Value | Location |
|----------|-------|----------|
| AWP MAX_SUPPLY | 10,000,000,000 * 10^18 | AWPToken |
| Alpha MAX_SUPPLY | 10,000,000,000 * 10^18 | AlphaToken (per worknet) |
| INITIAL_DAILY_EMISSION | 15,800,000 * 10^18 | Deploy.s.sol |
| EPOCH_DURATION | 86,400 (1 day) | AWPEmission (initialized via Deploy.s.sol) |
| DECAY_FACTOR | 996844 / 1000000 | AWPEmission |
| MAX_ACTIVE_WORKNETS | 10,000 | AWPRegistry |
| maxRecipients | 10,000 | AWPEmission |
| MAX_WEIGHT_SECONDS | 54 weeks (32,659,200s) | StakeNFT / AWPDAO |
| POOL_FEE | 10,000 (1%) | LPManager |
| TICK_SPACING | 200 | LPManager |
| Default immunity | 30 days | AWPRegistry |
| TIMELOCK_DELAY | 172,800 (2 days) | Deploy.s.sol |
| MAX_CHAIN_LENGTH | 256 | AWPRegistry (binding tree depth) |

---

## 7. Multi-Chain

AWP is deployed on 4 chains with identical contract addresses (via CREATE2):

| Chain | Chain ID | DEX |
|-------|----------|-----|
| Base | 8453 | Uniswap V4 |
| Ethereum | 1 | Uniswap V4 |
| Arbitrum | 42161 | Uniswap V4 |
| BSC | 56 | PancakeSwap V4 |

> LPManager and WorknetManager bytecode differs on BSC (PancakeSwap vs Uniswap), so these two contracts have different addresses on BSC.

**Key multi-chain design:**

- **WorknetId:** `(block.chainid << 64) | localCounter` -- globally unique across all chains.
- **Allocation is local:** Users allocate on their staking chain to any chain's worknet (cross-chain worknetId accepted, no on-chain status check).
- **Emission:** Per-chain AWPEmission. Guardian coordinates quotas across chains.
- **DAO:** Per-chain AWPDAO + Treasury. Off-chain aggregated voting.
- **AWPToken:** Per-chain independent mint (INITIAL_MINT configurable per chain).
- **API:** All REST endpoints accept `?chain_id=` to target a specific chain. Global/aggregated endpoints (suffixed `/global`) aggregate across all chains.
