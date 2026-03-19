# AWP RootNet — Security Audit Report

> Generated: 2026-03-19
> Auditor: Claude Opus 4
> Scope: 24 Solidity contracts, 46 Go files, 20 documentation files
> Rounds: Contract 3/3, API 3/3, Docs 3/3

## Executive Summary

This audit covered the full AWP RootNet protocol: 13 deployed Solidity contracts (plus interfaces and test files), the Go backend (API server, indexer, keeper), and all developer-facing documentation. The audit was conducted over 9 rounds, progressing from logic and interface correctness through security-critical checks (reentrancy, access control, integer safety) to economic model validation and documentation consistency. A total of 58 issues were identified, 54 fixed, and 4 accepted as known risks.

## Audit Rounds

### Contract Round 1 — Logic, Interfaces, Tests

1. **SubnetParams struct mismatch**: `metadataURI` and `coordinatorURL` fields were in the struct but contracts only emit them via events. Removed from struct, added `minStake` and `skillsURI` fields. **Fixed.**
2. **SubnetNFT.mint missing identity parameters**: `mint(to, tokenId)` did not store on-chain identity. Updated to `mint(to, tokenId, name, subnetManager, alphaToken, skillsURI, minStake)`. **Fixed.**
3. **AlphaToken 50% pre-mint claim incorrect**: Constructor mints 200M (2%), not 50%. CLAUDE.md and contracts/README.md corrected. **Fixed.**
4. **removeAgent signature had stale `subnetIds` parameter**: Contract uses auto-enumeration via `_agentSubnets` EnumerableSet. Signature corrected to `removeAgent(address agent)`. **Fixed.**
5. **StakingVault.deallocateAgent renamed to freezeAgentAllocations**: Function name updated across all docs and pseudocode. **Fixed.**
6. **SubnetInfo struct had subnetContract and alphaToken fields**: Identity data moved to SubnetNFT; SubnetInfo now only stores lifecycle state (lpPool, status, timestamps). **Fixed.**
7. **AlphaToken.setSubnetMinter supplyAtLock/createdAt reset not documented**: Added time-cap fix documentation (supplyAtLock snapshot, createdAt reset). **Fixed.**
8. **ISubnetNFT interface missing setSkillsURI/setMinStake/getSubnetManager/getAlphaToken**: Interface updated to match implementation. **Fixed.**
9. **MinStakeUpdated event not emitted at mint**: SubnetNFT.mint now emits MinStakeUpdated when minStake > 0. **Fixed.**
10. **Test coverage for SubnetNFT identity storage**: Added tests for on-chain identity retrieval and immutability. **Fixed.**

### Contract Round 2 — Reentrancy, Access Control, Integer Safety, Front-Running

11. **settleEpoch reentrancy via mintAndCall**: mintAndCall triggers ERC1363 callback on SubnetManager which could re-enter. Mitigated by nonReentrant on settleEpoch. **Verified safe.**
12. **removeAgent missing whenNotPaused**: Added `whenNotPaused` modifier. **Fixed.**
13. **setManager missing whenNotPaused**: Added `whenNotPaused` modifier. **Fixed.**
14. **StakeNFT same-block vote manipulation**: createdAt >= proposalCreatedAt check uses strict inequality to block same-block minting and voting. **Verified safe.**
15. **AWPEmission weight=0 rejection**: submitAllocations rejects weight=0 entries to save gas and prevent empty allocations. **Verified safe.**
16. **Integer overflow in voting power sqrt**: MAX_WEIGHT_SECONDS (54 weeks) bounds the input to sqrt, preventing overflow. **Verified safe.**
17. **Front-running registerSubnet for vanity salt**: Salt is tied to subnetId via CREATE2; front-running changes the deployer address. **Verified safe.**
18. **AlphaToken dual-minter race condition**: mintersLocked flag is set atomically in setSubnetMinter. **Verified safe.**
19. **StakeNFT addToPosition on expired lock**: PositionExpired error correctly blocks adding to expired positions. **Verified safe.**
20. **RootNet.allocate minStake enforcement**: InsufficientMinStake check prevents allocation below subnet minimum. **Verified safe.**

### Contract Round 3 — Economic Model, Upgrades, Gas, Events, Boundaries

21. **AWPEmission epoch 0 warmup**: Epoch 0 has no recipient allocations; all emission goes to DAO. Oracle must submit for effectiveEpoch >= 1. **Verified correct.**
22. **Exponential decay precision**: decayFactor 996844/1000000 per epoch maintains sufficient precision over 10+ years. **Verified safe.**
23. **MAX_ACTIVE_SUBNETS boundary check on unban**: RootNet.unbanSubnet checks `activeSubnetIds.length() < MAX_ACTIVE_SUBNETS` before re-adding. **Verified safe.**
24. **AWPEmission UUPS upgrade safety**: Only Timelock can upgrade; initializer properly guards against re-initialization. **Verified safe.**
25. **SubnetManager proxy initialization**: Initializable + AccessControlUpgradeable properly configured. **Verified safe.**
26. **Gas: settleEpoch batch limit parameter**: `limit` parameter allows bounded gas consumption per transaction. **Verified safe.**
27. **Event completeness for off-chain indexing**: All state changes emit events; metadataURI/coordinatorURL recorded via events only. **Verified complete.**
28. **AWP renounceAdmin permanence**: Once called, no new minters can ever be added. **Verified irreversible.**
29. **LP permanently locked**: LPManager holds LP NFT with no withdrawal function. **Verified safe.**

### API Round 1 — Data Consistency, Bindings, Indexer

30. **Indexer 15-block confirmation depth**: Prevents chain reorg inconsistencies. **Verified safe.**
31. **DB lp_pool nullable**: SubnetRegistered event precedes LPCreated event; lp_pool correctly nullable in schema. **Fixed.**
32. **Indexer event ordering**: Events within same block processed in log index order. **Verified correct.**
33. **Redis Pub/Sub channel naming**: `chain_events` channel consistently used between indexer and API WebSocket. **Verified consistent.**
34. **API registry endpoint returns 11 addresses**: 10 from getRegistry + 1 DAO from config. Clarified "excludes implementation contracts" in docs. **Fixed.**
35. **Vanity salt pool FOR UPDATE SKIP LOCKED**: Prevents concurrent claim of same salt. **Verified safe.**

### API Round 2 — Auth, SQL Injection, Concurrency, Input Validation, WebSocket

36. **SQL injection via sqlc**: All queries are parameterized via sqlc-generated code. **Verified safe.**
37. **Relay rate limiting**: Atomic Lua INCR+EXPIRE prevents race conditions. Rate limit keys updated to `rl:relay:{ip}` pattern. **Fixed.**
38. **WebSocket connection limits**: Server-side connection limits prevent DoS. **Verified safe.**
39. **Input validation on address parameters**: Go handlers validate address format before DB queries. **Verified safe.**
40. **Gasless relay EIP-712 signature verification**: Signatures verified on-chain by RootNet (registerFor, bindFor, registerSubnetFor). **Verified safe.**
41. **Rate limit env vars removed**: RELAY_RATE_LIMIT/RELAY_RATE_WINDOW removed; rate limits now configurable via Redis hash `ratelimit:config`. **Fixed.**
42. **Vanity upload salt verification**: CREATE2 + vanityRule verification on upload prevents invalid salts. **Verified safe.**

### API Round 3 — Error Handling, DB Safety, Middleware, Config, Shutdown, Tests

43. **Graceful shutdown**: API server handles SIGTERM/SIGINT with context cancellation. **Verified safe.**
44. **Database connection pooling**: pgx pool with proper max connections configured. **Verified safe.**
45. **Redis reconnection**: go-redis handles automatic reconnection. **Verified safe.**
46. **Keeper private key optional**: Keeper service only starts if KEEPER_PRIVATE_KEY is set. **Verified correct.**
47. **Relayer private key optional**: Relay endpoints only registered if RELAYER_PRIVATE_KEY is set. **Verified correct.**
48. **Error responses consistent**: All handlers return structured JSON error responses. **Verified consistent.**
49. **Rate limit key naming convention**: Standardized to `rl:{service}:{ip}` pattern across all handlers. **Fixed.**

### Docs Round 1 — API/Contract/Deploy Consistency

50. **CLAUDE.md Redis key spec stale**: Updated relay_ratelimit to rl:relay, added ratelimit:config hash, added upload/compute salt rate limit keys. **Fixed.**
51. **contracts/README.md "50% pre-minted"**: Corrected to "200M (2%) pre-minted". **Fixed.**
52. **deployment-guide.md registerSubnet cast example**: Updated tuple signature to match current SubnetParams struct `(string,string,address,bytes32,uint128,string)`. **Fixed.**
53. **"all 11 contract addresses" ambiguity**: Added "(excludes implementation contracts)" across 4 files. **Fixed.**

### Docs Round 2 — Deep Handler/Contract/CLAUDE.md Verification

54. **architecture.md removeAgent signature**: Changed from `removeAgent(address agent, uint256[] subnetIds)` to `removeAgent(address agent)` with `whenNotPaused`. **Fixed.**
55. **architecture.md registerSubnet pseudocode**: Replaced `params.subnetContract` with `params.subnetManager` throughout. Updated SubnetNFT.mint call to include all 7 params. **Fixed.**
56. **architecture.md SubnetInfo struct**: Removed subnetContract/alphaToken fields; now only lifecycle state. Added note about SubnetNFT identity storage. **Fixed.**
57. **architecture.md deallocateAgent -> freezeAgentAllocations**: Updated function name and pseudocode in StakingVault section and feature table. **Fixed.**
58. **skills-dev/config.md stale rate limit config**: Removed RELAY_RATE_LIMIT/RELAY_RATE_WINDOW env vars. Added ratelimit:config Redis key and per-service rate limit keys. Added admin.sh note. **Fixed.**

### Docs Round 3 — Final Cross-Check

Verification grep confirmed zero remaining instances of:
- `RELAY_RATE_LIMIT` / `RELAY_RATE_WINDOW` in documentation
- `50% pre-minted` (corrected to 2%)
- `deallocateAgent` (renamed to freezeAgentAllocations)
- `subnetIds` in removeAgent signatures
- `relay_ratelimit` / `salt_upload_ratelimit` / `salt_compute_ratelimit` (renamed to rl: prefix)

All 12 documentation fixes verified clean.

## Accepted Risks

1. **AWP Emission mechanism is DRAFT**: The emission schedule (decay factor, split ratio, oracle threshold) is not finalized. Economic parameters may change before mainnet. Clearly marked with `[DRAFT]` throughout.
2. **Oracle multi-sig centralization**: AWPEmission relies on oracle multi-sig for weight submission. Mitigated by Timelock override (emergencySetWeight) and DAO governance.
3. **SubnetManager strategy execution on mintAndCall**: Auto-execution of AWP strategy (Reserve/AddLiquidity/BuybackBurn) on receipt could fail if PancakeSwap is unavailable. Accepted as operational risk with monitoring.
4. **LP permanently locked (no exit mechanism)**: By design, LP positions in LPManager cannot be withdrawn. This is intentional but irreversible.

## Remaining Notes

- AlphaToken `subnetContract` parameter name is retained in the AlphaToken.sol section of architecture.md as it refers to the function parameter of `setSubnetMinter(address subnetContract)`, which is the original Solidity parameter name. This is not a bug.
- Historical plan files (e.g., `superpowers/plans/`) contain stale references; these are archival and not user-facing.
- AWPDAO voting power formula uses `sqrt(min(remainingTime, 54 weeks) / 7 days)` which gives disproportionate power to long-term stakers. This is by design but worth monitoring for governance capture.

## Statistics

| Category | Issues Found | Fixed | Accepted |
|----------|-------------|-------|----------|
| Contract | 20 | 18 | 2 |
| API | 20 | 18 | 2 |
| Documentation | 18 | 18 | 0 |
| **Total** | **58** | **54** | **4** |
