# AWPEmission V3 Design Spec

## Overview

Refactor AWPEmission from a subnet-aware emission scheduler to a generic distribution engine. Remove all subnet concepts from AWPEmission. Simplify settlement by eliminating state-based batch control in favor of a call parameter.

## Motivation

V2 AWPEmission has two design problems:
1. **Subnet coupling**: AWPEmission stores `SubnetWeight` structs keyed by `subnetId`, maintains `activeSubnetIds`, and exposes 5 lifecycle sync functions (`registerSubnet`, `activateSubnet`, `deactivateSubnet`, `reactivateSubnet`, `removeSubnet`). It should only care about `(address -> weight)`.
2. **Over-engineered settlement**: `batchSize` as state variable, `epochSettling` bool lock, `notSettling` modifier, `settleIndex` — all unnecessary. Batch size can be a call parameter, and settlement progress tracking needs only an index counter.

## Design Decisions

| Decision | Choice | Reasoning |
|----------|--------|-----------|
| Core abstraction | `address -> weight` (no subnet concept) | AWPEmission is a distribution engine, not subnet manager |
| Allocation mode | Full replacement per submission | Simple, no orphaned entries, oracle submits complete picture |
| Batch control | `settleEpoch(uint256 limit)` parameter | No state variable needed; caller decides gas budget |
| Settlement locking | `require(settleProgress == 0)` on submitAllocations | Simpler than dedicated bool + modifier; prevents array mutation during iteration |
| `setMaxActiveSubnets` | Removed | `settleEpoch(limit)` makes gas control caller-side |
| `setBatchSize` | Removed | Replaced by `limit` parameter |
| RootNet lifecycle notifications | Removed from AWPEmission | Oracle excludes banned/paused subnets; DAO uses `emergencySetWeight` |
| `emergencySetWeight` | Kept | DAO emergency override; accepted that next oracle submission overwrites it |
| `getActiveSubnetCount` | Moved to RootNet self-maintained | RootNet restores its own `activeSubnetIds` EnumerableSet |

## Contract Changes

### 1. AWPEmission V3

**Inheritance (unchanged from V2):**
```solidity
contract AWPEmission is
    Initializable,
    UUPSUpgradeable,
    ReentrancyGuardUpgradeable,
    EIP712Upgradeable,
    IAWPEmission
```

**Storage layout:**

NOTE: This is a greenfield deployment (V3 deploys as a new proxy, not upgrading a live V2).
Future V4 upgrades MUST preserve this slot order and only append within __gap.

```
// ── Inherited OZ storage (ERC-7201 namespaced, not in sequential slots) ──

// ── Sequential storage (slot 0+) ──
address public rootNet;                          // slot 0 (kept for reference, unused after init)
IAWPToken public awpToken;                       // slot 1
address public treasury;                         // slot 2

uint256 public epochDuration;                    // slot 3
uint256 public currentEpoch;                     // slot 4
uint256 public lastSettleTime;                   // slot 5
uint256 public currentDailyEmission;             // slot 6
// constants occupy no slots

address[] public recipients;                     // slot 7 (dynamic array length)
mapping(address => uint128) public weights;      // slot 8 (mapping root)
uint256 public totalWeight;                      // slot 9

uint256 public settleProgress;                   // slot 10: 0=idle, >0 = processed up to (settleProgress-1)
uint256 public epochEmissionLocked;              // slot 11
uint256 private _snapshotLen;                    // slot 12
uint256 private _snapshotWeight;                 // slot 13
uint256 private _snapshotPool;                   // slot 14
uint256 private _epochMinted;                    // slot 15

address[] public oracles;                        // slot 16 (dynamic array length)
uint256 public oracleThreshold;                  // slot 17
uint256 public allocationNonce;                  // slot 18

uint256 public maxRecipients;                    // slot 19 (default 10000)

uint256[40] private __gap;                       // slots 20-59 (reserve for future V4)
```

**Removed from V2:**
- `struct SubnetWeight` — replaced by `recipients[]` + `weights` mapping
- `mapping(uint256 => SubnetWeight) public subnetWeights` — gone
- `EnumerableSet.UintSet private activeSubnetIds` — moved to RootNet
- `uint128 public maxActiveSubnets` — deleted
- `uint128 public batchSize` — replaced by call parameter
- `bool public epochSettling` — replaced by `settleProgress > 0` check
- `uint256 public settleIndex` — replaced by `settleProgress`
- `_epochTotalWeight`, `_epochActiveCount`, `_epochSubnetPool`, `_epochSubnetMinted` — replaced by `_snapshotWeight`, `_snapshotLen`, `_snapshotPool`, `_epochMinted`
- All lifecycle functions: `registerSubnet`, `activateSubnet`, `deactivateSubnet`, `reactivateSubnet`, `removeSubnet`
- `onlyRootNet` modifier — no longer needed
- `NotRootNet` error — no longer needed
- `MaxActiveSubnetsReached` error — no longer needed
- `SubnetNotRegistered` error — replaced by address(0) check in submitAllocations
- `CurrentlySettling` error — replaced by `SettlementInProgress`

**New/modified functions:**

```solidity
/// @notice Full-replacement allocation submission with oracle multi-sig
function submitAllocations(
    address[] calldata recipients_,
    uint128[] calldata weights_,
    bytes[] calldata signatures
) external
```
Requirements:
- `settleProgress == 0` (not mid-settlement)
- Oracle signature verification (same as V2)
- `recipients_.length == weights_.length`
- `recipients_.length <= maxRecipients` (default 10000, bounds O(N) loops)
- No `address(0)` in recipients
- No duplicate addresses in recipients
Logic:
1. Verify oracle signatures on `(recipients_, weights_, allocationNonce)`
2. Clear old allocation: for each `recipients[i]`, `delete weights[recipients[i]]`; then `delete recipients`; `totalWeight = 0`
3. Write new allocation: `recipients = recipients_`; for each `(addr, w)`: `weights[addr] = w`; `totalWeight += w`
4. `allocationNonce++`
5. Emit `AllocationsSubmitted(nonce, recipients_, weights_)`

```solidity
/// @notice Settle current epoch, processing up to `limit` recipients per call
function settleEpoch(uint256 limit) external nonReentrant
```
Requirements: `limit > 0` (revert `InvalidParameter`)

Logic:
```
require(limit > 0)
if settleProgress == 0:
    require(block.timestamp >= lastSettleTime + epochDuration)
    if currentEpoch > 0: currentDailyEmission *= DECAY_FACTOR / DECAY_PRECISION
    epochEmissionLocked = min(currentDailyEmission, awpRemaining)
    require(epochEmissionLocked > 0)  // MiningComplete
    _snapshotLen = recipients.length
    _snapshotWeight = totalWeight
    _snapshotPool = epochEmissionLocked * EMISSION_SPLIT_BPS / 10000
    _epochMinted = 0
    settleProgress = 1  // marks Phase 1 done, start processing from index 0

start = settleProgress - 1
end = min(start + limit, _snapshotLen)

if _snapshotWeight > 0:
    for i in start..end:
        addr = recipients[i]
        w = weights[addr]
        if w > 0 && addr != address(0):
            share = _snapshotPool * w / _snapshotWeight
            if share > 0:
                minted = _mintTo(addr, share)
                _epochMinted += minted
                emit SubnetAWPDistributed(currentEpoch, addr, minted)

settleProgress = end + 1

if end >= _snapshotLen:
    daoShare = epochEmissionLocked > _epochMinted ? epochEmissionLocked - _epochMinted : 0
    if daoShare > 0: _mintTo(treasury, daoShare)
    emit DAOMatchDistributed(currentEpoch, daoShare)
    settleProgress = 0
    currentEpoch++
    lastSettleTime = block.timestamp
    emit EpochSettled(currentEpoch - 1, epochEmissionLocked, _snapshotLen)
```

Note on reading `recipients[i]` and `weights[addr]` during settlement: since `submitAllocations` requires `settleProgress == 0`, the recipients array and weights mapping cannot change mid-settlement. This is the same guarantee as V2's `notSettling` but simpler.

Note on `EpochSettled` event: uses `currentEpoch - 1` after increment. Per the V2 review, a local variable `epoch` captured before increment is cleaner — will be used in implementation.

```solidity
/// @notice DAO emergency weight override
function emergencySetWeight(address addr, uint128 weight) external onlyTimelock
```
- Modifies `weights[addr]` and adjusts `totalWeight`
- Requires `addr` exists in `recipients[]` (loop check)
- `settleProgress == 0` required (cannot modify during settlement)

```solidity
/// @notice Set oracle configuration
function setOracleConfig(address[] calldata oracles_, uint256 threshold_) external onlyTimelock
```
- `settleProgress == 0` required
- Same validations as V2 (threshold > 0, threshold <= length, no address(0), no duplicates)

```solidity
function setEpochDuration(uint256 d) external onlyTimelock
```
- `require(d > 0)`

```solidity
function _authorizeUpgrade(address) internal override onlyTimelock
```

```solidity
function getRecipientCount() external view returns (uint256)
function getOracleCount() external view returns (uint256)
```

**EIP-712:**
```solidity
bytes32 private constant ALLOCATION_TYPEHASH = keccak256(
    "SubmitAllocations(address[] recipients,uint128[] weights,uint256 nonce)"
);
```
Note: typehash changed from V2 (`subnetIds` → `recipients`, `uint256[]` → `address[]`).

**Errors:**
```solidity
error NotTimelock();
error InvalidRecipient();
error InvalidAmount();
error EpochNotReady();
error MiningComplete();
error SettlementInProgress();
error OracleNotConfigured();
error InvalidOracleConfig();
error InvalidSignatureCount();
error DuplicateOracle();
error UnknownOracle();
error DuplicateSigner();
error DuplicateRecipient();
error ArrayLengthMismatch();
error InvalidParameter();
error RecipientNotFound();
```

### 2. RootNet Changes

**Restore `activeSubnetIds` to RootNet:**
```solidity
EnumerableSet.UintSet private activeSubnetIds;
```

**Modify lifecycle functions** — remove AWPEmission calls, manage local set:
- `activateSubnet`: `activeSubnetIds.add(subnetId)` (was calling `IAWPEmission.activateSubnet`)
- `pauseSubnet`: `activeSubnetIds.remove(subnetId)` (was calling `IAWPEmission.deactivateSubnet`)
- `resumeSubnet`: `require(activeSubnetIds.length() < MAX_ACTIVE_SUBNETS)` then `activeSubnetIds.add(subnetId)` (capacity check + add)
- `banSubnet`: `activeSubnetIds.remove(subnetId)` (was calling `IAWPEmission.deactivateSubnet`)
- `unbanSubnet`: `require(activeSubnetIds.length() < MAX_ACTIVE_SUBNETS)` then `activeSubnetIds.add(subnetId)` (capacity check + add)
- `deregisterSubnet`: `activeSubnetIds.remove(subnetId)` (was calling `IAWPEmission.removeSubnet`)
- `registerSubnet`: remove `IAWPEmission.registerSubnet()` call

**Restore view functions:**
```solidity
function getActiveSubnetCount() external view returns (uint256) {
    return activeSubnetIds.length();
}
function getActiveSubnetIdAt(uint256 index) external view returns (uint256) {
    return activeSubnetIds.at(index);
}
```

**Keep:**
- `currentEpoch()` — still delegates to `IAWPEmission(awpEmission).currentEpoch()`

**Re-add imports:**
- `EnumerableSet` import restored to RootNet

### 3. IAWPEmission Interface

```solidity
interface IAWPEmission {
    function settleEpoch(uint256 limit) external;
    function submitAllocations(address[] calldata recipients, uint128[] calldata weights, bytes[] calldata signatures) external;
    function emergencySetWeight(address addr, uint128 weight) external;
    function setOracleConfig(address[] calldata oracles_, uint256 threshold_) external;
    function setEpochDuration(uint256 d) external;

    function currentEpoch() external view returns (uint256);
    function totalWeight() external view returns (uint256);
    function epochDuration() external view returns (uint256);
    function lastSettleTime() external view returns (uint256);
    function currentDailyEmission() external view returns (uint256);
    function settleProgress() external view returns (uint256);
    function oracleThreshold() external view returns (uint256);
    function allocationNonce() external view returns (uint256);
    function getOracleCount() external view returns (uint256);
    function getRecipientCount() external view returns (uint256);

    event AllocationsSubmitted(uint256 indexed nonce, address[] recipients, uint128[] weights);
    event OracleConfigUpdated(address[] oracles, uint256 threshold);
    event GovernanceWeightUpdated(address indexed addr, uint128 weight);
    event SubnetAWPDistributed(uint256 indexed epoch, address indexed recipient, uint256 awpAmount);
    event DAOMatchDistributed(uint256 indexed epoch, uint256 amount);
    event EpochSettled(uint256 indexed epoch, uint256 totalEmission, uint256 recipientCount);
}
```

**Removed from V2 interface:**
- `registerSubnet`, `activateSubnet`, `deactivateSubnet`, `reactivateSubnet`, `removeSubnet`
- `rootNet()` view (still exists on contract but not needed in interface — RootNet doesn't call it)
- `setBatchSize`, `setMaxActiveSubnets`
- `getActiveSubnetCount`, `getActiveSubnetIdAt` (moved to RootNet)
- `epochSettling` (replaced by `settleProgress`)

### 4. Deployment Changes

Deploy.s.sol and TestDeploy.s.sol:
- AWPEmission `initialize()` unchanged (same 5 params)
- Remove any `setMaxActiveSubnets` or `setBatchSize` calls post-deploy
- Oracle config (`setOracleConfig`) remains as post-deploy Timelock action

### 5. Keeper Changes

keeper.go:
- `sendSettleEpoch` must pass a `limit` argument: `awpEmission.SettleEpoch(txAuth, big.NewInt(200))`
- `trySettleEpoch` reads `settleProgress` instead of `epochSettling` to check if mid-settlement

### 6. Indexer Changes

indexer.go:
- Remove parsing for `GovernanceWeightUpdated` with `subnetId` — now uses `address`
- Update `SubnetAWPDistributed` parsing to use `address` instead of `subnetId`
- `AllocationsSubmitted` event data changes from `subnetIds` to `recipients` (address[])

## Security Considerations

1. **Settlement integrity**: `submitAllocations` requires `settleProgress == 0`, preventing array mutation during settlement
2. **Oracle consensus**: Same multi-sig verification as V2
3. **Replay protection**: `allocationNonce` incremented per submission
4. **Emergency override**: `emergencySetWeight` only modifies existing recipients, also requires `settleProgress == 0`
5. **Full replacement safety**: Old weights are deleted before new ones written — no stale entries
6. **Duplicate recipient guard**: `submitAllocations` checks for duplicate addresses to prevent double-counting in `totalWeight`

## Test Plan

1. **AWPEmission.t.sol** — rewrite:
   - Proxy deployment (unchanged)
   - `submitAllocations` with addresses (not subnetIds): valid sigs, full replacement verification, duplicate check
   - `settleEpoch(limit)` with varying limits: process all in one call, process in 2 batches, process one-by-one
   - `settleProgress` blocking: submitAllocations reverts during settlement, emergencySetWeight reverts during settlement
   - Emission decay, DAO share, EpochSettled event
   - Oracle config, nonce replay, upgrade

2. **RootNet.t.sol** — update:
   - `getActiveSubnetCount` now reads from RootNet's own set
   - Lifecycle functions no longer call AWPEmission

3. **E2E.t.sol** — update:
   - Oracle submits addresses (not subnetIds)
   - `settleEpoch` called with limit parameter

4. **Integration.t.sol** — update:
   - Same as E2E

## Files Changed

| File | Action |
|------|--------|
| `contracts/src/token/AWPEmission.sol` | Rewrite (remove subnet, simplify settlement) |
| `contracts/src/interfaces/IAWPEmission.sol` | Rewrite |
| `contracts/src/RootNet.sol` | Restore activeSubnetIds, remove AWPEmission lifecycle calls |
| `contracts/test/helpers/EmissionSigningHelper.sol` | Update typehash (subnetIds→recipients) |
| `contracts/test/AWPEmission.t.sol` | Rewrite |
| `contracts/test/RootNet.t.sol` | Update |
| `contracts/test/E2E.t.sol` | Update |
| `contracts/test/Integration.t.sol` | Update |
| `contracts/script/Deploy.s.sol` | Minor update |
| `contracts/script/TestDeploy.s.sol` | Minor update |
| `api/internal/chain/keeper.go` | Add limit param to SettleEpoch call |
| `api/internal/chain/indexer.go` | Update event parsing |
| `api/internal/chain/bindings/` | Regenerate |
