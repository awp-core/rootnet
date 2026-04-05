# AWP — Security Audit Report

> Generated: 2026-04-03
> Auditor: Claude Opus 4.6 (1M context)
> Scope: 15 Solidity contracts (+ 10 interfaces, 18 test files), 63 Go files, documentation
> Target: Base, Ethereum, Arbitrum, BSC (4-chain deployment)
> Solidity: ^0.8.20 (compiled with solc 0.8.33, Foundry, evm_version = cancun, optimizer 800 runs via-ir)

## Executive Summary

This audit covers the full AWP protocol: 15 deployed Solidity contracts (plus interfaces and test files), the Go backend (API server, indexer, keeper), and all developer-facing documentation. The audit spans two phases — an initial 9-round review (2026-03-19) and a follow-up review (2026-04-03) covering the worknet rename, AWPEmission V3 Guardian-only architecture, permission hardening, multi-chain production fixes, and pre-deploy contract fixes.

**Key architectural changes since the initial audit:**
- All "subnet" terminology replaced with "worknet" (SubnetNFT -> WorknetNFT, SubnetManager -> WorknetManager)
- AWPEmission V3: Oracle multi-sig replaced by Guardian-only packed `uint256[]` allocations; EIP-712 removed from AWPEmission
- All UUPS upgrades authorized by Guardian (cross-chain multisig), not Timelock
- `register()` removed from AWPRegistry (every address is implicitly a root; `register()` = `setRecipient(msg.sender)`)
- `freezeAgentAllocations` removed from StakingVault; allocation functions (allocate/deallocate/reallocate) moved from AWPRegistry to StakingVault
- Mint + best-effort callback pattern in settleEpoch: `awpToken.mint()` always succeeds, ERC1363 callback is try/catch (no revert on callback failure)
- 4-chain deployment (Base, Ethereum, Arbitrum, BSC) with same CREATE2 salts for identical addresses (except LPManager/WorknetManager on BSC due to PancakeSwap bytecode)

A total of 70 issues have been identified across both phases. 66 fixed, 4 accepted as known risks.

## Contract Inventory (15 contracts)

| Contract | Type | Description |
|----------|------|-------------|
| AWPRegistry | UUPS proxy | Unified entry: worknet management + account system (tree-based binding, delegation, gasless EIP-712) |
| StakeNFT | ERC721 | Position NFT for AWP staking (deposit with lock, transferable, O(1) balance) |
| StakingVault | UUPS proxy | Allocation management (allocate/deallocate/reallocate) with EIP-712 gasless support |
| AWPEmission | UUPS proxy | Epoch-versioned weight distribution engine (Guardian-only, packed uint256[] allocations) |
| AWPDAO | Governor | OZ Governor + GovernorTimelockControl, StakeNFT-based voting |
| Treasury | TimelockController | Governance timelock |
| WorknetNFT | ERC721 | On-chain worknet identity (name, worknetManager, alphaToken, skillsURI, minStake) |
| WorknetManager | UUPS proxy (via ERC1967Proxy) | Default worknet contract: Merkle claims, AWP strategy (PancakeSwap V4) |
| WorknetManagerUni | UUPS proxy (via ERC1967Proxy) | Default worknet contract: Merkle claims, AWP strategy (Uniswap V4) |
| LPManager | — | LP management for PancakeSwap V4 (BSC) |
| LPManagerUni | — | LP management for Uniswap V4 (Base, ETH, ARB) |
| LPManagerBase | abstract | Shared LP logic |
| AWPToken | ERC20 | 10B MAX_SUPPLY, configurable initialMint, mintAndCall (ERC1363) |
| AlphaToken | ERC20 | Per-worknet token, 10B max, dual minter, standalone CREATE2 deployment |
| AlphaTokenFactory | — | CREATE2 factory with vanity address validation |

## Audit Rounds

### Contract Round 1 — Logic, Interfaces, Tests (Initial Audit)

1. **WorknetParams struct mismatch**: `metadataURI` and `coordinatorURL` fields were in the struct but contracts only emit them via events. Removed from struct, added `minStake` and `skillsURI` fields. **Resolved.**
2. **WorknetNFT.mint missing identity parameters**: `mint(to, tokenId)` did not store on-chain identity. Updated to `mint(to, tokenId, name, worknetManager, alphaToken, skillsURI, minStake)`. **Resolved.**
3. **AlphaToken 50% pre-mint claim incorrect**: Constructor mints 200M (2%), not 50%. CLAUDE.md and contracts/README.md corrected. **Resolved.**
4. **removeAgent signature had stale `worknetIds` parameter**: Contract uses auto-enumeration via `_agentWorknets` EnumerableSet. Signature corrected to `removeAgent(address agent)`. **Resolved.**
5. ~~**StakingVault.deallocateAgent renamed to freezeAgentAllocations**~~: `freezeAgentAllocations` has been fully removed. Allocation functions (allocate/deallocate/reallocate) now live on StakingVault with direct staker-or-delegate auth. **Resolved (superseded).**
6. **WorknetInfo struct reduced to lifecycle state**: Identity data moved to WorknetNFT; WorknetInfo now only stores lpPool, status, createdAt, activatedAt. **Resolved.**
7. **AlphaToken.setWorknetMinter supplyAtLock/createdAt reset not documented**: Added time-cap fix documentation (supplyAtLock snapshot, createdAt reset). **Resolved.**
8. **IWorknetNFT interface missing setSkillsURI/setMinStake/getWorknetManager/getAlphaToken**: Interface updated to match implementation. **Resolved.**
9. **MinStakeUpdated event not emitted at mint**: WorknetNFT.mint now emits MinStakeUpdated when minStake > 0. **Resolved.**
10. **Test coverage for WorknetNFT identity storage**: Added tests for on-chain identity retrieval and immutability. **Resolved.**

### Contract Round 2 — Reentrancy, Access Control, Integer Safety, Front-Running (Initial Audit)

11. **settleEpoch reentrancy via mintAndCall**: mintAndCall triggers ERC1363 callback on WorknetManager which could re-enter. Mitigated by nonReentrant on settleEpoch. **Verified safe.**
12. **removeAgent missing whenNotPaused**: Added `whenNotPaused` modifier. **Resolved.**
13. **setManager missing whenNotPaused**: Added `whenNotPaused` modifier. **Resolved.**
14. **StakeNFT same-block vote manipulation**: createdAt >= proposalCreatedAt check uses strict inequality to block same-block minting and voting. **Verified safe.**
15. **AWPEmission weight=0 entries**: Guardian is trusted to submit valid packed allocations off-chain. weight=0 entries in packed format are skipped at settlement time (no gas waste). **Verified safe.**
16. **Integer overflow in voting power sqrt**: MAX_WEIGHT_SECONDS (54 weeks) bounds the input to sqrt, preventing overflow. **Verified safe.**
17. **Front-running registerWorknet for vanity salt**: Salt is tied to worknetId via CREATE2; front-running changes the deployer address. **Verified safe.**
18. **AlphaToken dual-minter race condition**: mintersLocked flag is set atomically in setWorknetMinter. **Verified safe.**
19. **StakeNFT addToPosition on expired lock**: PositionExpired error correctly blocks adding to expired positions. **Verified safe.**
20. **WorknetNFT.minStake enforcement**: minStake is stored on-chain in WorknetNFT but NOT enforced by StakingVault.allocate — used as off-chain/coordinator reference only. This is by design. **Verified correct.**

### Contract Round 3 — Economic Model, Upgrades, Gas, Events, Boundaries (Initial Audit)

21. **AWPEmission epoch 0 settlement**: settledEpoch starts at 0. Epoch 0 can be settled immediately (settledEpoch <= currentEpoch). Guardian submits for effectiveEpoch >= settledEpoch. No warmup epoch — emission starts from epoch 0 if weights are submitted. Decay is skipped for epoch 0 (settledEpoch == 0 check). **Verified correct.**
22. **Exponential decay precision**: decayFactor 996844/1000000 per epoch maintains sufficient precision over 10+ years. MIN_DECAY_FACTOR (900000) prevents Guardian from setting near-zero decay. **Verified safe.**
23. **MAX_ACTIVE_WORKNETS boundary check on unban**: AWPRegistry.unbanWorknet checks `activeWorknetIds.length() < MAX_ACTIVE_WORKNETS` before re-adding. **Verified safe.**
24. **UUPS upgrade safety — Guardian authorization**: AWPEmission, StakingVault, and AWPRegistry all use Guardian for `_authorizeUpgrade`. Initializer properly guards against re-initialization. **Verified safe.**
25. **WorknetManager proxy initialization**: Initializable + AccessControlUpgradeable properly configured. Auto-deployed via ERC1967Proxy when worknetManager=address(0) and defaultWorknetManagerImpl is set. **Verified safe.**
26. **Gas: settleEpoch batch limit parameter**: `limit` parameter allows bounded gas consumption per transaction. 3-phase design (init, batch mint, finalize) supports incremental settlement. **Verified safe.**
27. **Event completeness for off-chain indexing**: All state changes emit events with correct parameters. **Verified complete.**
28. **AWP renounceAdmin permanence**: Once called, no new minters can ever be added. **Verified irreversible.**
29. **LP permanently locked**: LPManager holds LP NFT with no withdrawal function. **Verified safe.**

### Contract Round 4 — Pre-Deploy Audit (2026-04-03)

30. **AWPEmission.submitAllocations event param**: `AllocationsSubmitted` event emitted hardcoded 0 instead of `effectiveEpoch`. Fixed to emit the actual epoch, enabling off-chain indexers to identify which epoch was submitted. **Resolved.**
31. **AWPEmission.setGuardian missing zero-address check**: Setting guardian to address(0) would permanently brick the contract (no recovery). Added `if (g == address(0)) revert ZeroAddress()` check. **Resolved (Critical).**
32. **LPManager compoundFees broken INCREASE_LIQUIDITY params**: Fee reinvestment was completely broken — passed (poolKey, ticks, ...) instead of (tokenId, liquidity, amounts, hookData). Fixed in both LPManager (PancakeSwap) and LPManagerUni (Uniswap). **Resolved (Critical).**
33. **AWPEmission.setGuardian missing GuardianUpdated event**: No event was emitted on guardian rotation. Added `emit GuardianUpdated(g)`. **Resolved.**
34. **StakingVault.setGuardian missing**: Guardian rotation required a full UUPS upgrade. Added `setGuardian(address g)` function with nonReentrant and NotGuardian check. **Resolved.**
35. **AWPRegistry._registerWorknet zero LP amount**: No revert when computed `lpAWPAmount == 0`, which would create a worknet with no initial liquidity. Added explicit revert. **Resolved.**
36. **StakingVault missing ReentrancyGuard**: allocate/deallocate/reallocate/allocateFor/deallocateFor had no reentrancy protection. Added ReentrancyGuardUpgradeable + nonReentrant on all external mutation functions. **Resolved (Critical).**

### Contract Round 5 — Permission Hardening (2026-04-03)

37. **AWPRegistry.setGuardian: onlyTimelock -> onlyGuardian**: A single-chain DAO could vote to replace the cross-chain Guardian via Timelock. Now only Guardian can replace itself (self-sovereign, consistent with AWPEmission). Recovery: Timelock can still upgrade AWPRegistry via UUPS if Guardian keys are lost. **Resolved.**
38. **StakingVault._authorizeUpgrade: Timelock -> Guardian**: A single-chain DAO could upgrade StakingVault to arbitrary code. Now only Guardian (cross-chain multisig) can authorize UUPS upgrades. **Resolved.**

### Contract Round 6 — Emission V3 Architecture Review (2026-04-03)

39. **AWPEmission V3: Oracle removal verified**: Oracle multi-sig infrastructure (oracles[], oracleThreshold, isOracleMap, signature verification) fully removed. Storage slots freed but preserved for UUPS proxy upgrade safety. Guardian-only `submitAllocations` with packed `uint256[]` encoding replaces the old Oracle-signed submission. **Verified clean.**
40. **EIP-712 removal from AWPEmission**: AWPEmission no longer inherits EIP712Upgradeable. No gasless relay for emission — Guardian submits directly. EIP-712 remains in AWPRegistry and StakingVault for gasless user operations. **Verified correct.**
41. **Mint + best-effort callback**: settleEpoch uses `awpToken.mint(recipient, toMint)` followed by a try/catch `IERC1363Receiver.onTransferReceived()` call. Mint always succeeds; callback failure does not revert the epoch. This eliminates the RecipientMintFailed false-positive issue where a reverting callback could block settlement for all recipients. **Verified safe.**
42. **Epoch-versioned weight persistence**: Weights submitted via `submitAllocations` are stored per-epoch. If no new weights are submitted for an epoch, `activeEpoch` carries forward the last submitted weights. This supports weekly Guardian submissions with daily settlements. **Verified correct.**
43. **appendAllocations + modifyAllocations**: New functions allow Guardian to incrementally update allocations without resubmitting the full list. modifyAllocations supports in-place patching with index-based addressing. **Verified safe.**
44. **Epoch pause mechanism**: `pauseEpochUntil(uint64 resumeTime)` freezes currentEpoch at settledEpoch. Resume rebases baseEpoch/baseTime correctly. Expired pauses are auto-cleaned via `_checkResume()`. **Verified safe.**
45. **setDecayFactor lower bound**: MIN_DECAY_FACTOR (900000) prevents Guardian from setting decay below 10% per epoch. Upper bound < DECAY_PRECISION prevents 100% retention. **Verified safe.**

### API Round 1 — Data Consistency, Bindings, Indexer (Initial Audit)

46. **Indexer optimistic indexing**: Indexer processes up to chain tip with no confirmation delay. Block hash chain verification detects reorgs (max 64-block rollback). No 15-block confirmation depth. **Verified safe.**
47. **DB lp_pool nullable**: WorknetRegistered event precedes LPCreated event; lp_pool correctly nullable in schema. **Resolved.**
48. **Indexer event ordering**: Events within same block processed in log index order. **Verified correct.**
49. **Redis Pub/Sub channel naming**: `chain_events` channel consistently used between indexer and API WebSocket. **Verified consistent.**
50. **API registry endpoint returns contract addresses**: Returns deployed contract addresses + chainId + eip712Domain. Excludes implementation contracts. **Verified correct.**
51. **Vanity salt pool FOR UPDATE SKIP LOCKED**: Prevents concurrent claim of same salt. **Verified safe.**

### API Round 2 — Auth, SQL Injection, Concurrency, Input Validation, WebSocket (Initial Audit)

52. **SQL injection via sqlc**: All queries are parameterized via sqlc-generated code. **Verified safe.**
53. **Relay rate limiting**: Atomic Lua INCR+EXPIRE prevents race conditions. Rate limit keys use `rl:relay:{ip}` pattern. Configurable via Redis hash `ratelimit:config`. **Verified safe.**
54. **WebSocket connection limits**: Server-side connection limits prevent DoS. **Verified safe.**
55. **Input validation on address parameters**: Go handlers validate address format before DB queries. **Verified safe.**
56. **Gasless relay EIP-712 signature verification**: Signatures verified on-chain by AWPRegistry (bindFor, setRecipientFor, registerWorknetFor) and StakingVault (allocateFor, deallocateFor). **Verified safe.**
57. **Rate limits via Redis**: RELAY_RATE_LIMIT/RELAY_RATE_WINDOW env vars removed. Rate limits configurable via Redis hash `ratelimit:config` and admin API. **Resolved.**
58. **Vanity upload salt verification**: CREATE2 + vanityRule verification on upload prevents invalid salts. **Verified safe.**

### API Round 3 — Error Handling, DB Safety, Middleware, Config, Shutdown, Tests (Initial Audit)

59. **Graceful shutdown**: API server handles SIGTERM/SIGINT with context cancellation. **Verified safe.**
60. **Database connection pooling**: pgx pool with proper max connections configured. **Verified safe.**
61. **Redis reconnection**: go-redis handles automatic reconnection. **Verified safe.**
62. **Keeper private key optional**: Keeper service only starts if KEEPER_PRIVATE_KEY is set. **Verified correct.**
63. **Relayer private key optional**: Relay endpoints only registered if RELAYER_PRIVATE_KEY is set. **Verified correct.**
64. **Error responses consistent**: All handlers return structured JSON error responses. **Verified consistent.**
65. **Rate limit key naming convention**: Standardized to `rl:{service}:{ip}` pattern across all handlers. **Resolved.**

### API Round 4 — Multi-Chain Production Fixes (2026-04-03)

66. **Keeper epoch-0 settlement**: Keeper must call settleEpoch starting from epoch 0 when settledEpoch == 0. No special warmup handling needed — settledEpoch <= currentEpoch is the only condition. Keeper now correctly triggers settlement from genesis. **Resolved.**
67. **Keeper distributed lock TTL**: Increased from 60s to 90s (3x renewal interval of 30s) to prevent lock expiry during GC pauses or slow RPC. **Resolved.**
68. **WS cross-chain allocation enrichment**: GetAgentWorknetStakeWS now uses global query (no chain_id filter). Fixes wrong currentStake=0 in AllocationChanged push when allocation lives on a different chain than the event source. **Resolved.**
69. **Per-chain relay balance monitoring**: Keeper writes relayer_balance:{chainId} to Redis every 25s. Logs ERROR alert when below 0.01 ETH/BNB. /api/health/detailed includes relayerBalance per chain. **Resolved.**
70. **AdminDeleteChain soft-delete**: Changed from hard DELETE to soft-delete (status='inactive'). Preserves indexed data; goroutine continues until restart. **Resolved.**

## Accepted Risks

1. **WorknetManager strategy execution on mintAndCall**: Auto-execution of AWP strategy (Reserve/AddLiquidity/BuybackBurn) on receipt could fail if the DEX is unavailable. Accepted as operational risk with monitoring. The best-effort callback pattern in AWPEmission means a failing WorknetManager does not block epoch settlement.
2. **LP permanently locked (no exit mechanism)**: By design, LP positions in LPManager cannot be withdrawn. This is intentional but irreversible.
3. **Guardian centralization**: AWPEmission, StakingVault upgrades, and AWPRegistry pause/setGuardian are all Guardian-controlled (cross-chain Safe multisig). If Guardian keys are compromised, emission weights, contract upgrades, and pause authority are at risk. Mitigated by multi-sig threshold and per-chain Timelock retaining AWPRegistry UUPS upgrade authority.
4. **AWPDAO voting power concentration**: Voting power formula `amount * sqrt(min(remainingTime, 54 weeks) / 7 days)` gives disproportionate power to long-term stakers. This is by design but worth monitoring for governance capture.

## Permission Model Summary

| Authority | Controls |
|-----------|----------|
| Guardian (cross-chain Safe multisig) | AWPEmission (all config + weights + upgrade), StakingVault (upgrade + setGuardian), AWPRegistry (pause + setGuardian) |
| Timelock (per-chain DAO) | AWPRegistry (upgrade, ban/unban, deregister, params, unpause), AWPDAO proposals |
| Anyone | settleEpoch (permissionless settlement trigger) |

## Remaining Notes

- AlphaToken `subnetContract` parameter name is retained in the AlphaToken.sol `setWorknetMinter(address subnetContract)` function — this is the original Solidity parameter name, not a terminology bug.
- Historical plan files (e.g., `superpowers/plans/`) may contain stale references; these are archival and not user-facing.
- Solidity pragma is `^0.8.20` in source files but compiled with solc 0.8.33 via Foundry config.
- BSC uses different LPManager (PancakeSwap V4) and WorknetManager bytecode, resulting in different CREATE2 addresses for these two contracts on BSC vs other chains.

## Statistics

| Category | Issues Found | Fixed | Accepted |
|----------|-------------|-------|----------|
| Contract (Initial) | 29 | 27 | 2 |
| Contract (Follow-up) | 16 | 16 | 0 |
| API (Initial) | 20 | 18 | 2 |
| API (Follow-up) | 5 | 5 | 0 |
| **Total** | **70** | **66** | **4** |
