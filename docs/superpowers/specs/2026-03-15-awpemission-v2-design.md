> **OBSOLETE — Superseded by V3 (2026-03-15-awpemission-v3-design.md)**
> This document describes the V2 design which was never deployed. See V3 spec for the current architecture.

# AWPEmission V2 Design Spec

## Overview

Refactor AWPEmission to a UUPS upgradeable proxy with multi-oracle consensus for batch weight submission. Remove all emission delegation from RootNet. DAO governance controls oracle configuration and contract upgrades, not individual weight settings.

## Motivation

Current design has three problems:
1. **RootNet as middleman**: DAO votes route through RootNet for no reason, adding gas and coupling.
2. **Individual weight voting**: 10000 subnets = 10000 DAO proposals. Unscalable.
3. **Immutable emission contract**: AWPToken's minter list is permanently locked after `renounceAdmin()`. Cannot deploy a new AWPEmission without losing mint permission.

## Architecture

```
Off-chain voting system (complex scoring logic)
        |
        v
N oracle nodes (independently verify voting results)
        | off-chain consensus + aggregated signatures
        v
AWPEmission V2 (UUPS Proxy)
|-- submitAllocations(ids[], weights[], sigs[])  <-- oracle submission
|-- settleEpoch()                                 <-- permissionless
|-- register/activate/deactivate/...              <-- RootNet lifecycle
|-- setOracleConfig(oracles[], threshold)          <-- DAO governance
|-- upgradeTo(newImpl)                             <-- DAO governance
```

## Design Decisions

| Decision | Choice | Reasoning |
|----------|--------|-----------|
| Trust model | Multi-oracle consensus | Eliminates committee, DAO only governs oracle config |
| Oracle consensus | Off-chain aggregated signatures, single on-chain tx | Lowest gas (1 tx), oracles coordinate off-chain |
| N/M threshold | DAO-governed, initially 3/5 | Flexible, can adjust as network grows |
| Weight submission frequency | On-demand, long-lived | Weights persist until replaced by new batch |
| RootNet delegation functions | Remove all 4 | DAO calls AWPEmission directly via Treasury |
| Upgradeability | UUPS proxy | Required because AWPToken minter list is permanently locked |

## Contract Changes

### 1. AWPEmission V2 (UUPS Upgradeable)

**Inheritance:**
```solidity
contract AWPEmission is
    Initializable,
    UUPSUpgradeable,
    ReentrancyGuardUpgradeable,
    EIP712Upgradeable,
    IAWPEmission
```

**Constructor → Initializer:**
```solidity
/// @custom:oz-upgrades-unsafe-allow constructor
constructor() { _disableInitializers(); }

function initialize(
    address rootNet_,
    address awpToken_,
    address treasury_,
    uint256 initialDailyEmission_,
    uint256 epochDuration_
) external initializer {
    __UUPSUpgradeable_init();
    __ReentrancyGuard_init();
    __EIP712_init("AWPEmission", "2");
    rootNet = rootNet_;
    awpToken = IAWPToken(awpToken_);
    treasury = treasury_;
    currentDailyEmission = initialDailyEmission_;
    epochDuration = epochDuration_;
    lastSettleTime = block.timestamp;
}
```

**Storage changes:**
- `rootNet`: `immutable` → regular state variable (proxy incompatible with immutable)
- `awpToken`: `immutable` → regular state variable
- New: `address[] public oracles` — registered oracle addresses
- New: `uint256 public oracleThreshold` — minimum valid signatures required (M)
- New: `uint256 public allocationNonce` — replay protection counter
- New: `uint256[47] private __gap` — storage gap for future upgrades (3 new vars + 47 gap = 50-slot block)

**New functions:**

```solidity
/// @notice Batch-submit subnet weights with aggregated oracle signatures
/// @param subnetIds Array of subnet IDs to update
/// @param weights Array of new weights (parallel to subnetIds)
/// @param signatures Array of EIP-712 signatures from oracles
function submitAllocations(
    uint256[] calldata subnetIds,
    uint128[] calldata weights,
    bytes[] calldata signatures
) external notSettling;
```

Verification logic:
1. Require `oracleThreshold > 0` (guard against uninitialized oracle config)
2. Construct EIP-712 struct hash: `keccak256(abi.encode(ALLOCATION_TYPEHASH, keccak256(abi.encodePacked(subnetIds)), keccak256(abi.encodePacked(weights)), allocationNonce))`
3. Compute digest via `_hashTypedDataV4(structHash)` (includes EIP-712 domain separator)
4. For each signature: `ECDSA.recover(digest, sig)`, verify signer is a registered oracle
5. Duplicate signer check: maintain `address[signatures.length] memory seen` array, O(N^2) check against prior signers. Revert if duplicate found.
6. Require valid unique signers >= `oracleThreshold`
7. Batch update: for each `(subnetId, weight)` pair, require `subnetWeights[subnetId].recipient != address(0)` (subnet must be registered), update weight and adjust `totalWeight` if active
8. Increment `allocationNonce`
9. Emit `AllocationsSubmitted(nonce, subnetIds, weights)`

```solidity
/// @notice Set oracle configuration (DAO governance)
/// @param oracles_ New oracle address list
/// @param threshold_ Minimum signatures required
function setOracleConfig(
    address[] calldata oracles_,
    uint256 threshold_
) external onlyTimelock;
```

Validations:
- `threshold_ > 0`
- `threshold_ <= oracles_.length`
- No `address(0)` in `oracles_` (ECDSA.recover returns address(0) for invalid signatures, which would bypass threshold)
- No duplicate addresses in `oracles_`

```solidity
/// @notice UUPS upgrade authorization — only Treasury
function _authorizeUpgrade(address) internal override onlyTimelock;
```

**Modified modifier:**
```solidity
/// @dev onlyTimelock — only Treasury can call (RootNet no longer accepted)
modifier onlyTimelock() {
    if (msg.sender != treasury) revert NotTimelock();
    _;
}
```

**Kept functions (onlyTimelock, DAO calls directly):**
- `setEpochDuration(uint256 d)` — unchanged
- `setBatchSize(uint128 b)` — unchanged
- `setMaxActiveSubnets(uint128 m)` — unchanged

**Removed functions:**
- `setWeight()` — replaced by `submitAllocations()` for normal ops. Kept as `emergencySetWeight()` with `onlyTimelock` for DAO emergency override.

**EIP-712 Domain:**
```solidity
// Domain: name="AWPEmission", version="2"
bytes32 private constant ALLOCATION_TYPEHASH = keccak256(
    "SubmitAllocations(uint256[] subnetIds,uint128[] weights,uint256 nonce)"
);
```

### 2. RootNet Changes

**Remove these functions entirely:**
- `setGovernanceWeight(uint256 subnetId, uint128 w)`
- `setEpochDuration(uint256 d)`
- `setBatchSize(uint128 b)`
- `setMaxActiveSubnets(uint128 m)`

**Keep unchanged:**
- All subnet lifecycle notifications to AWPEmission (`registerSubnet`, `activateSubnet`, `deactivateSubnet`, `reactivateSubnet`, `removeSubnet`)
- `currentEpoch()` view function (still delegates to AWPEmission)
- `getActiveSubnetCount()` view function (still delegates to AWPEmission)

**Remove `currentEpoch()` passthrough? No.** Keep it for backward compatibility — staking functions use it internally, and external callers may depend on it.

### 3. IAWPEmission Interface Updates

**Add:**
```solidity
function submitAllocations(
    uint256[] calldata subnetIds,
    uint128[] calldata weights,
    bytes[] calldata signatures
) external;

function emergencySetWeight(uint256 subnetId, uint128 weight) external;

function setOracleConfig(address[] calldata oracles, uint256 threshold) external;

function oracles(uint256 index) external view returns (address);
function oracleThreshold() external view returns (uint256);
function allocationNonce() external view returns (uint256);
function getOracleCount() external view returns (uint256);

event AllocationsSubmitted(uint256 indexed nonce, uint256[] subnetIds, uint128[] weights);
event OracleConfigUpdated(address[] oracles, uint256 threshold);
```

**Remove:**
- `setWeight()` — replaced by `emergencySetWeight()` and `submitAllocations()`

### 4. Deployment Changes

**New deployment sequence (Deploy.s.sol):**
```
Step 13a: Deploy AWPEmission implementation
Step 13b: Deploy ERC1967Proxy(impl, initData)
Step 14:  awp.addMinter(proxy address)  // proxy address is permanent
Step 15:  awp.renounceAdmin()           // locks minter to proxy forever
Step 16:  AWPEmission(proxy).setOracleConfig(initialOracles, 3)
```

After deployment, future upgrades:
```
1. Deploy new implementation contract
2. DAO proposal: AWPEmission(proxy).upgradeTo(newImpl)
3. Treasury executes after timelock delay
```

**TestDeploy.s.sol:**
Same pattern with simplified parameters.

### 5. Keeper Changes

**keeper.go:**
- Remove `rootNet` field entirely (already done in optimization pass)
- No changes needed — Keeper already reads from AWPEmission and calls `AWPEmission.SettleEpoch()`

### 6. Indexer Changes

**indexer.go:**
- Add parsing for new events: `AllocationsSubmitted`, `OracleConfigUpdated`
- These events come from AWPEmission address (already in filter list)

## Data Flow

### Normal weight update (via oracles):
```
1. Off-chain voting system produces allocation list
2. Oracle nodes independently verify results
3. Oracles coordinate off-chain, produce aggregated signatures
4. Any oracle submits: AWPEmission.submitAllocations(ids, weights, sigs)
5. Contract verifies >= M valid oracle signatures
6. Weights updated in batch, totalWeight adjusted
7. Next settleEpoch() uses new weights
```

### Emergency weight override (via DAO):
```
1. DAO proposal: AWPEmission.emergencySetWeight(subnetId, weight)
2. Treasury executes after timelock delay
3. Single subnet weight updated
```

### Contract upgrade:
```
1. Deploy new AWPEmission implementation
2. DAO proposal: AWPEmission.upgradeTo(newImpl)
3. Treasury executes after timelock delay
4. Proxy delegates to new implementation
5. All state preserved, address unchanged, minter permission intact
```

## Security Considerations

1. **Replay protection**: `allocationNonce` increments on each submission, included in EIP-712 digest
2. **Oracle compromise**: Requires M/N oracles compromised simultaneously. DAO can replace oracles via `setOracleConfig`
3. **Upgrade safety**: `_authorizeUpgrade` restricted to Treasury (Timelock). Storage layout must be append-only (gaps reserved)
4. **Settlement protection**: `notSettling` modifier on `submitAllocations` prevents weight changes during epoch settlement
5. **Minter invariant**: Proxy address never changes, so AWPToken minter permission is preserved across upgrades
6. **Deployment window**: `submitAllocations` requires `oracleThreshold > 0`, preventing submissions before `setOracleConfig` is called
7. **Duplicate signature prevention**: O(N^2) check against `seen` array in `submitAllocations` ensures a single oracle cannot meet threshold alone
8. **Unregistered subnet guard**: `submitAllocations` requires `subnetWeights[subnetId].recipient != address(0)` for each entry, preventing weight pollution for non-existent subnets

## Storage Layout

```
// Slot 0-N: inherited from Initializable, UUPSUpgradeable, ReentrancyGuard
// Existing AWPEmission V1 state (preserved in exact order):
address public rootNet;
IAWPToken public awpToken;
address public treasury;
uint256 public epochDuration;
uint256 public currentEpoch;
uint256 public lastSettleTime;
uint256 public currentDailyEmission;
mapping(uint256 => SubnetWeight) public subnetWeights;
EnumerableSet.UintSet private activeSubnetIds;
uint256 public totalWeight;
uint128 public maxActiveSubnets;
uint128 public batchSize;
uint256 public settleIndex;
uint256 public epochEmissionLocked;
uint256 private _epochSubnetPool;
uint256 private _epochSubnetMinted;
uint256 private _epochTotalWeight;
uint256 private _epochActiveCount;
bool public epochSettling;

// New V2 state (appended after existing):
address[] public oracles;
uint256 public oracleThreshold;
uint256 public allocationNonce;

// Storage gap for future upgrades (3 new V2 vars + 47 gap = 50-slot V2 extension block)
// Future V3: reduce __gap by the number of new variables added
uint256[47] private __gap;
```

## Test Plan

1. **AWPEmission.t.sol** — rewrite for proxy deployment pattern:
   - Deploy impl + proxy + initialize
   - Test submitAllocations with valid/invalid signatures
   - Test oracle threshold enforcement
   - Test nonce replay protection
   - Test setOracleConfig access control
   - Test emergencySetWeight via treasury
   - Test settleEpoch still works (3-phase, unchanged)
   - Test upgrade via treasury

2. **RootNet.t.sol** — remove emission delegation tests:
   - Remove test_setEpochDuration, test_setMaxActiveSubnets
   - Remove setGovernanceWeight from test_onlyTimelockFunctions
   - Keep lifecycle notification tests

3. **E2E.t.sol** — update DAO governance tests:
   - DAO proposal targets AWPEmission directly (not RootNet)
   - Oracle-based weight submission flow
   - Upgrade flow via DAO proposal

4. **Integration.t.sol** — update full flow:
   - Oracle submits weights → settleEpoch → verify distribution

## Migration Notes

This is a **greenfield deployment** (not upgrading a live contract). The UUPS pattern is for future upgradeability after initial deployment. No migration of existing state is needed.

If deployed on a chain where AWPEmission V1 is already live with locked minters:
1. Cannot deploy V2 as new contract (minter list locked)
2. Would need AWPToken governance to add new minter (impossible after renounceAdmin)
3. Only option: deploy new AWPToken + new AWPEmission V2 together

## Files Changed

| File | Action |
|------|--------|
| `contracts/src/token/AWPEmission.sol` | Rewrite (UUPS + oracle consensus) |
| `contracts/src/interfaces/IAWPEmission.sol` | Update (new functions/events) |
| `contracts/src/RootNet.sol` | Remove 4 delegation functions |
| `contracts/src/interfaces/IRootNet.sol` | No change |
| `contracts/script/Deploy.s.sol` | Proxy deployment pattern |
| `contracts/script/TestDeploy.s.sol` | Proxy deployment pattern |
| `contracts/test/AWPEmission.t.sol` | Rewrite |
| `contracts/test/RootNet.t.sol` | Remove emission delegation tests |
| `contracts/test/E2E.t.sol` | Update DAO governance flow |
| `contracts/test/Integration.t.sol` | Update emission flow |
| `api/internal/chain/indexer.go` | Add new event parsing |
| `api/internal/chain/bindings/` | Regenerate |
| `CLAUDE.md` | Update architecture description |
| `docs/architecture.md` | Update AWPEmission section |
