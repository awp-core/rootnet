# CLAUDE.md

## Project
AWP RootNet — Agent Mining protocol on BSC.

## Architecture Document
docs/architecture.md — Read the relevant section before starting any task.
.env - Read the bsc rpc endpoint

## Stack
- Contracts: Solidity 0.8.24, Foundry (evm_version = cancun), OpenZeppelin 5.x
- Backend: Go 1.26, Chi v5, sqlc + pgx/v5, PostgreSQL, Redis, go-ethereum
- Frontend: Next.js 14, Tailwind, wagmi/viem

## API Architecture
- Backend is read-only + on-chain data indexer + gasless relay; three relay endpoints (register, bind, register-subnet)
- Frontend sends transactions directly to chain via wagmi/viem
- Three independent Go processes: api (HTTP+WS) / indexer (event sync) / keeper (scheduled on-chain ops)
- Indexer → Redis Pub/Sub → API WebSocket
- Indexer uses 15-block confirmation depth to avoid chain reorgs

## Core Architecture (13 contracts)
- RootNet.sol = Unified entry: subnet management + staking allocation. No deposit/withdraw — staking via StakeNFT. No epoch logic (epoch moved to AWPEmission). Gasless: registerFor, bindFor, registerSubnetFor, registerSubnetForWithPermit (fully gasless with ERC-2612 permit). Errors: NotTimelock, NotGuardian (distinct).
- StakeNFT.sol = ERC721 position NFT (not Enumerable). Users deposit AWP with lock period (timestamp-based, lockDuration in seconds). Each position = NFT with (amount, lockEndTime, createdAt). NFTs are transferable. O(1) balance tracking via _userTotalStaked accumulator. getUserVotingPower requires caller to pass tokenIds. addToPosition blocked on expired locks (PositionExpired error).
- StakingVault.sol = Pure allocation logic. Allocations are plain uint128. (user, agent, subnetId) triple; allocate/deallocate/reallocate all immediate. Auto-enumerates agent subnets via EnumerableSet — freezeAgentAllocations(user, agent) needs no caller-supplied list. allocate/reallocate reject subnetId=0.
- AWPEmission.sol = UUPS upgradeable proxy: generic address→weight distribution engine. Epoch authority: genesisTime + epochDuration (1 day). currentEpoch() = (block.timestamp - genesisTime) / epochDuration. Oracle multi-sig submits epoch-versioned packed allocations (submitAllocations(address[] recipients, uint96[] weights, bytes[] signatures, uint256 effectiveEpoch)). settleEpoch(limit) batch-mints AWP via mintAndCall (triggers SubnetManager.onTransferReceived). settledEpoch tracks settlement progress. emergencySetWeight(epoch, index, addr, weight) for Timelock override. onlyTimelock manages oracle config.
- AWPDAO.sol = Inherits OZ Governor, GovernorSettings, GovernorTimelockControl. Overrides _getVotes and _countVote for StakeNFT-based voting (no delegate/checkpoint). No rootNet dependency. Voters submit tokenId[] arrays. Voting power = amount * sqrt(min(remainingTime, 54 weeks) / 7 days). Anti-manipulation: only NFTs with createdAt < proposalCreatedAt (strict: >= blocks same-block mint+vote). Per-tokenId double-vote prevention. totalVotingPower > 0 required for proposal creation. Two proposal types: proposeWithTokens (executable via Timelock) and signalPropose (vote-only). propose() is blocked.
- SubnetNFT.sol = ERC721 with on-chain identity storage. tokenId = subnetId. Stores immutable fields: name, subnetManager, alphaToken. Stores owner-updatable fields: skillsURI (via setSkillsURI), minStake (via setMinStake). Events: SkillsURIUpdated, MinStakeUpdated. Lifecycle status managed by RootNet, not SubnetNFT.
- SubnetManager.sol = Default subnet contract (deployed behind ERC1967Proxy via RootNet when subnetManager=address(0)). Initializable + AccessControlUpgradeable + ReentrancyGuardUpgradeable + IERC1363Receiver. Three roles: MERKLE_ROLE (submit Merkle roots), STRATEGY_ROLE (AWP handling), TRANSFER_ROLE (token transfers). Merkle claim mints Alpha to users. AWP strategy: Reserve / AddLiquidity / BuybackBurn. onTransferReceived auto-executes strategy on AWP receipt via mintAndCall. PancakeSwap V4 BSC mainnet addresses hardcoded.
- AccessManager / LPManager = onlyRootNet; StakeNFT = independent; AWPEmission = onlyTimelock (governance)

## Tokens
- AWP: 10B MAX_SUPPLY; 200M (2%) minted in constructor; 98% AWPEmission mint on demand. mintAndCall(to, amount, data) triggers ERC1363 callback on recipient.
- Alpha: 10B max per subnet, dual minter; standalone CREATE2 deployment (not proxy). `supplyAtLock` snapshot + `createdAt` reset at `setSubnetMinter` — subnet minters can mint immediately after activation.

## Gas Optimization Design
- settleEpoch(limit) (AWPEmission) iterates recipients[] in bounded batches via mintAndCall
- SubnetInfo on RootNet stores only lifecycle state (lpPool, status, createdAt, activatedAt) — identity data in SubnetNFT
- Allocations use plain uint128 mapping, no struct wrapper
- RootNet manages activeSubnetIds locally with MAX_ACTIVE_SUBNETS = 10000 constant; AWPEmission has maxRecipients = 10000

> **AWP Emission: DRAFT — mechanism not finalized. Descriptions below are preliminary.**

## Emission
- Epoch 0 is a warmup epoch: no recipient allocations (all emission goes to DAO). Oracle must submit for effectiveEpoch >= 1 before epoch 0 is settled; weights take effect starting epoch 1.
- Exponential decay: currentEmission *= 996844 / 1000000
- Per epoch: 50% mintAndCall to subnets (by governanceWeight), daoShare = total - subnet minted
- Recipients omitted from submitAllocations have their share go to DAO. weight=0 entries are rejected to save gas. addr==address(0) entries are rejected by emergencySetWeight.

## AlphaTokenFactory
- Uses CREATE2 full deployment (no Clones / EIP-1167 proxy). Each AlphaToken is a standalone contract with no delegatecall overhead.
- Constructor: `(deployer, vanityRule)`. Replaceable via RootNet.setAlphaTokenFactory (Timelock).
- Vanity address system: 8 positions (4 prefix + 4 suffix) packed into uint64. Per-position: 0-9=digit, 10-15=lowercase a-f, 16-21=uppercase A-F (EIP-55), >=22=wildcard. Set at factory deployment (immutable). vanityRule=0 skips all validation.
- Example `"A1????cafe"`: 0x1001FFFF0C0A0F0E
- `deploy(subnetId, name, symbol, admin, salt)` — salt=bytes32(0) uses subnetId as salt; salt!=0 is user-provided

## Subnet Registration
- If subnetManager == address(0) and defaultSubnetManagerImpl is set, auto-deploys SubnetManager proxy via ERC1967Proxy
- AWP transferFrom(user → LPManager); Alpha mint(LPManager)
- setSubnetMinter(sc) permanently locks minter to subnet manager
- SubnetParams: name, symbol, subnetManager, salt, minStake, skillsURI
- registerSubnetFor: gasless EIP-712 subnet registration (requires prior AWP approve)
- registerSubnetForWithPermit: fully gasless — ERC-2612 permit + EIP-712 in one tx (user signs two messages, zero gas)
- SubnetNFT.mint stores identity (name, subnetManager, alphaToken) + initial minStake

## Staking
- StakeNFT: deposit AWP with lockDuration (seconds) → mint position NFT (amount, lockEndTime, createdAt). Transferable. addToPosition increases amount (blocked if lock expired — PositionExpired). withdraw after lock expires.
- StakingVault: pure allocation logic. Auto-enumerates agent subnets via _agentSubnets EnumerableSet. freezeAgentAllocations(user, agent) needs no subnet list.
- removeAgent(agent) — no subnetIds param; StakingVault auto-enumerates
- (user, agent, subnetId) triple; allocate/deallocate/reallocate all immediate; subnetId=0 rejected
- Epoch duration: 1 day (daily epochs, AWPEmission only)

## Deployment Sequence
AWPToken(constructor mint 200M) → AlphaTokenFactory(deployer, vanityRule)
→ Treasury → RootNet(deployer, treasury, guardian) → SubnetNFT → AccessManager → LPManager
→ AWPEmission impl → ERC1967Proxy(impl, initData with genesisTime_ and epochDuration_)
→ StakingVault → StakeNFT(awpToken, stakingVault, rootNet)
→ SubnetManager impl
→ AWPDAO (6 params, no rootNet_)
→ grantRole(dao) + renounce + addMinter(awpEmissionProxy) + renounceAdmin + AlphaTokenFactory.setAddresses
→ RootNet.initializeRegistry (deployer calls, then zeroed; 9 addresses)
→ AWP transfer distribution

## API Endpoints
- REST: /api/registry (all 11 contract addresses, excludes implementation contracts), /api/users/*, /api/agents/*, /api/staking/*, /api/subnets/*, /api/emission/*, /api/tokens/*, /api/governance/*
- Relay: POST /api/relay/register, /api/relay/bind, /api/relay/register-subnet (gasless EIP-712, rate-limited 100/IP/1h, configurable via Redis)
- Vanity: GET /api/vanity/mining-params, POST /api/vanity/upload-salts, GET /api/vanity/salts, GET /api/vanity/salts/count, POST /api/vanity/compute-salt (DB pool first, cast fallback)
- WebSocket: WS /ws/live

## Redis Key Spec
- alpha_price:{subnetId} → JSON, TTL=10m
- awp_info → JSON, TTL=1m
- emission_current → JSON, TTL=30s
- ratelimit:config → hash, persistent, hot-updatable rate limit configs (admin.sh)
- rl:relay:{ip} → counter, TTL=1h (atomic Lua INCR+EXPIRE)
- rl:upload_salts:{ip} → counter, TTL=1h (5 uploads/hr/IP)
- rl:compute_salt:{ip} → counter, TTL=1h (20 compute/hr/IP)
- channel: chain_events → Pub/Sub

## Bug Prevention
- tokenId uses _nextSubnetId++
- User/Agent address mutual exclusion (Principal/Agent model with bind/unbind)
- removeAgent(agent) freezes stake via StakingVault auto-enumeration (no caller-supplied subnetIds)
- StakeNFT: only NFTs with createdAt < proposalCreatedAt can vote (strict: >= blocks same-block)
- StakeNFT: addToPosition blocked on expired locks (PositionExpired)
- StakingVault: allocate/reallocate reject subnetId=0
- RootNet.allocate enforces SubnetNFT.minStake — allocation must result in total agent stake >= minStake (InsufficientMinStake)
- AWPDAO: totalVotingPower > 0 required for proposals
- deregisterSubnet: users must manually deallocate from deregistered subnets (deallocate has no status check); frontend should alert on SubnetDeregistered
- subnetManager == address(0) auto-deploys SubnetManager proxy if defaultSubnetManagerImpl is set
- setSubnetMinter permanently locked; ban uses minterPaused
- AWPEmission weights submitted by oracle multi-sig; epoch-versioned packed allocations
- RootNet.unbanSubnet checks MAX_ACTIVE_SUBNETS before re-adding
- AWP: deployer is never a minter; renounceAdmin permanently locks
- settleEpoch has nonReentrant
- DB lp_pool is nullable (SubnetRegistered precedes LPCreated)
- DB vanity_salts: salt pool with CREATE2 + vanityRule verification on upload; FOR UPDATE SKIP LOCKED on claim
- SubnetNFT stores identity on-chain (name, subnetManager, alphaToken, skillsURI, minStake)
- Permit2 BSC mainnet: 0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768
- admin.sh: Hot-update rate limits, manage salt pool, view system status (scripts/admin.sh)
