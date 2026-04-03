# CLAUDE.md

## Project
AWP — Agent Mining protocol on Base (Uniswap V4), BSC (PancakeSwap V4), Ethereum, and Arbitrum.

## Architecture Document
docs/architecture.md — Read the relevant section before starting any task.
.env - Read the bsc rpc endpoint

## Stack
- Contracts: Solidity 0.8.24, Foundry (evm_version = cancun), OpenZeppelin 5.x
- Backend: Go 1.26, Chi v5, sqlc + pgx/v5, PostgreSQL, Redis, go-ethereum
- Frontend: Next.js 14, Tailwind, wagmi/viem

## API Architecture
- Backend is read-only + on-chain data indexer + gasless relay; 7 relay endpoints (bind, set-recipient, register, allocate, deallocate, activate-worknet, register-worknet)
- Frontend sends transactions directly to chain via wagmi/viem
- Three independent Go processes: api (HTTP+WS) / indexer (event sync) / keeper (scheduled on-chain ops)
- Indexer → Redis Pub/Sub → API WebSocket
- Indexer uses optimistic indexing (processes up to chain tip, no confirmation delay) with block hash chain verification for reorg detection (max 64-block rollback)

## Core Architecture (11 contracts)
- AWPRegistry.sol = Unified entry: worknet management + account system (UUPS proxy). No deposit/withdraw — staking via StakeNFT. No epoch logic. EIP-712 domain name "AWPRegistry". No mandatory registration — every address is implicitly a root. `register()` is optional (= `setRecipient(msg.sender)`). Tree-based binding via `bind(target)`. `grantDelegate(delegate)` / `revokeDelegate(delegate)` for delegation. Gasless: bindFor, setRecipientFor, registerWorknetFor, registerWorknetForWithPermit. WorknetId globally unique: `(block.chainid << 64) | localCounter`. **Allocation functions moved to StakingVault.**
- StakeNFT.sol = ERC721 position NFT (not Enumerable). Deposit AWP with lock period. Each position = NFT with (amount, lockEndTime, createdAt). Transferable. O(1) balance tracking via _userTotalStaked. addToPosition blocked on expired locks (PositionExpired).
- StakingVault.sol = UUPS proxy with EIP-712 gasless support. `allocate(staker, agent, worknetId, amount)` — caller must be staker or delegate (reads AWPRegistry.delegates). `deallocate`, `reallocate` same auth. Gasless: `allocateFor`, `deallocateFor`. EIP-712 domain name "StakingVault". Cross-chain allocate: worknetId from any chain, no on-chain status check. worknetId=0 rejected. Auto-enumerates agent worknets via EnumerableSet.
- AWPEmission.sol = UUPS upgradeable proxy: generic address→weight distribution engine. Epoch authority: genesisTime + epochDuration (1 day). currentEpoch() = (block.timestamp - genesisTime) / epochDuration. Guardian (cross-chain multisig) submits epoch-versioned packed allocations (submitAllocations(address[] recipients, uint96[] weights, uint256 effectiveEpoch)). settleEpoch(limit) batch-mints AWP via mintAndCall (triggers WorknetManager.onTransferReceived). settledEpoch tracks settlement progress. 100% emission to recipients; Guardian includes treasury in recipients for DAO share. onlyGuardian manages all configuration.
- AWPDAO.sol = Inherits OZ Governor, GovernorSettings, GovernorTimelockControl. Overrides _getVotes and _countVote for StakeNFT-based voting (no delegate/checkpoint). No awpRegistry dependency. Voters submit tokenId[] arrays. Voting power = amount * sqrt(min(remainingTime, 54 weeks) / 7 days). Anti-manipulation: only NFTs with createdAt < proposalCreatedAt (strict: >= blocks same-block mint+vote). Per-tokenId double-vote prevention. totalVotingPower > 0 required for proposal creation. Two proposal types: proposeWithTokens (executable via Timelock) and signalPropose (vote-only). propose() is blocked.
- WorknetNFT.sol = ERC721 with on-chain identity storage. tokenId = worknetId. Stores immutable fields: name, worknetManager, alphaToken. Stores owner-updatable fields: skillsURI (via setSkillsURI), minStake (via setMinStake), metadataURI (via setMetadataURI, overrides tokenURI). Events: SkillsURIUpdated, MinStakeUpdated, MetadataURIUpdated. tokenURI 3-tier: per-token metadataURI → global baseURI → on-chain Base64 JSON. Lifecycle status managed by AWPRegistry, not WorknetNFT.
- WorknetManager.sol = Default worknet contract (deployed behind ERC1967Proxy via AWPRegistry when worknetManager=address(0)). UUPS upgradeable + AccessControlUpgradeable + ReentrancyGuardUpgradeable + IERC1363Receiver. Three roles: MERKLE_ROLE (submit Merkle roots), STRATEGY_ROLE (AWP handling), TRANSFER_ROLE (token transfers). Merkle claim mints Alpha to users. AWP strategy: Reserve / AddLiquidity / BuybackBurn. onTransferReceived auto-executes strategy on AWP receipt via mintAndCall. DEX addresses injected at init time via dexConfig (not hardcoded).
- LPManager = onlyAWPRegistry; compoundFees(alphaToken) reinvests accumulated LP fees (called by Keeper cron). StakeNFT = independent; AWPEmission = onlyGuardian (governance)

## Tokens
- AWP: 10B MAX_SUPPLY; initial mint configurable per chain (INITIAL_MINT constructor param, immutable); emission via AWPEmission. mintAndCall(to, amount, data) triggers ERC1363 callback on recipient.
- Alpha: 10B max per worknet, dual minter; standalone CREATE2 deployment (not proxy). `supplyAtLock` snapshot + `createdAt` reset at `setWorknetMinter` — worknet minters can mint immediately after activation.

## Gas Optimization Design
- settleEpoch(limit) (AWPEmission) iterates recipients[] in bounded batches via mintAndCall
- WorknetInfo on AWPRegistry stores only lifecycle state (lpPool, status, createdAt, activatedAt) — identity data in WorknetNFT
- Account system V2: no mandatory registration, tree-based binding, no address mutual exclusion
- Allocations use plain uint128 mapping, no struct wrapper
- AWPRegistry manages activeWorknetIds locally with MAX_ACTIVE_WORKNETS = 10000 constant; AWPEmission has maxRecipients = 10000

## Emission
- Guardian-only submission: cross-chain multisig submits weights directly via submitAllocations (no Oracle signatures, no Timelock dependency).
- settledEpoch starts at 0 (= next epoch to settle). Epoch 0 can be settled immediately (settledEpoch <= currentEpoch). Guardian submits for effectiveEpoch >= settledEpoch. No warmup epoch — emission starts from epoch 0 if weights are submitted.
- Exponential decay: currentEmission *= 996844 / 1000000
- Per epoch: 100% mintAndCall to recipients (by weight). Guardian includes treasury address in recipients for DAO share.
- Removed: oracle multi-sig, emissionSplitBps, emergencySetWeight, daoShare. weight=0 entries are rejected to save gas. addr==address(0) entries are rejected.

## AlphaTokenFactory
- Uses CREATE2 full deployment (no Clones / EIP-1167 proxy). Each AlphaToken is a standalone contract with no delegatecall overhead.
- Constructor: `(deployer, vanityRule)`. Replaceable via AWPRegistry.setAlphaTokenFactory (Timelock).
- Vanity address system: 8 positions (4 prefix + 4 suffix) packed into uint64. Per-position: 0-9=digit, 10-15=lowercase a-f, 16-21=uppercase A-F (EIP-55), >=22=wildcard. Set at factory deployment (immutable). vanityRule=0 skips all validation.
- Example `"A1????cafe"`: 0x1001FFFF0C0A0F0E
- `deploy(worknetId, name, symbol, admin, salt)` — salt=bytes32(0) uses worknetId as salt; salt!=0 is user-provided

## Worknet Registration
- If worknetManager == address(0) and defaultWorknetManagerImpl is set, auto-deploys WorknetManager proxy via ERC1967Proxy
- AWP transferFrom(user → LPManager); Alpha mint(LPManager)
- setWorknetMinter(sc) permanently locks minter to worknet manager
- WorknetParams: name, symbol, worknetManager, salt, minStake, skillsURI
- registerWorknetFor: gasless EIP-712 worknet registration (requires prior AWP approve)
- registerWorknetForWithPermit: fully gasless — ERC-2612 permit + EIP-712 in one tx (user signs two messages, zero gas)
- WorknetNFT.mint stores identity (name, worknetManager, alphaToken) + initial minStake

## Staking
- StakeNFT: deposit AWP with lockDuration (seconds) → mint position NFT (amount, lockEndTime, createdAt). Transferable. addToPosition increases amount (blocked if lock expired — PositionExpired). withdraw after lock expires.
- StakingVault (UUPS proxy, EIP-712 domain "StakingVault"): allocation management + gasless support. Users call `allocate(staker, agent, worknetId, amount)` directly on StakingVault (not AWPRegistry). Auth: staker or delegate (reads AWPRegistry.delegates cross-contract). Gasless: `allocateFor`, `deallocateFor`. worknetId from ANY chain (globally unique). No on-chain worknet status check. Balance check is local only.
- (staker, agent, worknetId) triple; allocate/deallocate/reallocate all immediate; worknetId=0 rejected
- Epoch duration: 1 day (daily epochs, AWPEmission only)

## Deployment Sequence
AWPToken(constructor: name, symbol, deployer, initialMint) → AlphaTokenFactory(deployer, vanityRule)
→ Treasury → AWPRegistry impl → ERC1967Proxy(impl, initialize(deployer, treasury, guardian)) → WorknetNFT → LPManager
→ AWPEmission impl → ERC1967Proxy(impl, initData with genesisTime_ and epochDuration_)
→ StakingVault impl → ERC1967Proxy(impl, initialize(awpRegistry, treasury)) → StakeNFT(awpToken, stakingVault, awpRegistry)
→ WorknetManager impl
→ AWPDAO (6 params, no awpRegistry_)
→ grantRole(dao) + renounce + addMinter(awpEmissionProxy) + renounceAdmin + AlphaTokenFactory.setAddresses
→ AWPRegistry.initializeRegistry (deployer calls, then zeroed; 9 params: 8 addresses + dexConfig bytes)
→ AWP transfer distribution

## Multi-Chain
- Deploy config: chains.yaml (Base, Ethereum, Arbitrum, BSC)
- Deploy script: scripts/deploy-multichain.sh <chainName|--all|--list>
- Same CREATE2 salts → identical addresses on all chains. GENESIS_TIME must be set explicitly (no block.timestamp fallback). AWPToken.initialMint() called post-deploy with per-chain amount. LPManager/WorknetManager bytecode differs by DEX (Uniswap vs PancakeSwap) so BSC addresses differ for these two contracts — use per-chain overrides in chains.yaml
- WorknetId: (block.chainid << 64) | localCounter — globally unique
- Allocate is local: user allocates on their staking chain to any chain's worknet
- Emission: per-chain AWPEmission, Guardian coordinates quotas
- DAO: per-chain AWPDAO + Treasury, off-chain aggregated voting
- AWPToken: per-chain independent mint (INITIAL_MINT configurable)

## API Endpoints
- REST: /api/registry (all contract addresses + chainId + eip712Domain, excludes implementation contracts). eip712Domain object includes name, version, chainId, verifyingContract. /api/users/*, /api/staking/*, /api/worknets/*, /api/emission/*, /api/tokens/*, /api/governance/*
- Nonce: GET /api/nonce/{address} — AWPRegistry nonce (bind/setRecipient/registerWorknet). GET /api/staking-nonce/{address} — StakingVault nonce (allocate/deallocate)
- Relay: POST /api/relay/bind, /api/relay/set-recipient, /api/relay/register, /api/relay/allocate, /api/relay/deallocate, /api/relay/activate-worknet, /api/relay/register-worknet (gasless EIP-712, rate-limited 100/IP/1h, configurable via Redis)
- Vanity: GET /api/vanity/mining-params, POST /api/vanity/upload-salts, GET /api/vanity/salts, GET /api/vanity/salts/count, POST /api/vanity/compute-salt (DB pool first, cast fallback)
- WebSocket: WS /ws/live

## Redis Key Spec
- alpha_price:{worknetId} → JSON, TTL=10m
- awp_info:{chainId} → JSON, TTL=1m
- emission_current:{chainId} → JSON, TTL=30s
- ratelimit:config → hash, persistent, hot-updatable rate limit configs (admin.sh)
- rl:relay:{ip} → counter, TTL=1h (atomic Lua INCR+EXPIRE)
- rl:upload_salts:{ip} → counter, TTL=1h (5 uploads/hr/IP)
- rl:compute_salt:{ip} → counter, TTL=1h (20 compute/hr/IP)
- channel: chain_events → Pub/Sub

## Bug Prevention
- tokenId uses _nextWorknetId++
- Tree-based binding with anti-cycle check (no address mutual exclusion)
- `bind(target)` walks chain to detect cycles before binding
- `resolveRecipient(addr)` walks boundTo chain to root for reward distribution
- StakeNFT: only NFTs with createdAt < proposalCreatedAt can vote (strict: >= blocks same-block)
- StakeNFT: addToPosition blocked on expired locks (PositionExpired)
- StakingVault: allocate/deallocate/reallocate all reject worknetId=0
- WorknetNFT.minStake is stored on-chain but NOT enforced by AWPRegistry.allocate (used as off-chain/coordinator reference only)
- AWPDAO: totalVotingPower > 0 required for proposals
- deregisterWorknet: users must manually deallocate from deregistered worknets (deallocate has no status check); frontend should alert on WorknetDeregistered
- worknetManager == address(0) auto-deploys WorknetManager proxy if defaultWorknetManagerImpl is set
- setWorknetMinter permanently locked; ban uses minterPaused
- AWPEmission weights submitted by Guardian multi-sig; epoch-versioned packed allocations
- AWPRegistry.unbanWorknet checks MAX_ACTIVE_WORKNETS before re-adding
- AWP: deployer is never a minter; renounceAdmin permanently locks
- settleEpoch has nonReentrant
- DB lp_pool is nullable (WorknetRegistered precedes LPCreated)
- DB vanity_salts: salt pool with CREATE2 + vanityRule verification on upload; FOR UPDATE SKIP LOCKED on claim
- DB: `agents` table removed; `users` table has `bound_to` and `recipient` columns
- WorknetNFT stores identity on-chain (name, worknetManager, alphaToken, skillsURI, minStake)
- Permit2 BSC mainnet: 0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768
- admin.sh: Hot-update rate limits, manage salt pool, view system status (scripts/admin.sh)
- Vanity address rules: contracts/.env defines VANITY_PREFIX_*/VANITY_SUFFIX_* for key contracts (AWPToken: 0000/00a1, AWPRegistry: 0000/00b1). deploy.sh MUST mine salts via scripts/vanity/mine.sh before deployment — never deploy with --skip-mine on first deploy. Changing contract constructors/init params invalidates initCodeHash → must re-mine salts
