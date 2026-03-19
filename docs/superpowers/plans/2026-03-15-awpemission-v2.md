# AWPEmission V2 Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor AWPEmission to UUPS upgradeable proxy with multi-oracle batch weight submission, removing all emission delegation from AWPRegistry.

**Architecture:** AWPEmission becomes a UUPS proxy with EIP-712 signature verification for oracle-submitted weight batches. AWPRegistry retains only subnet lifecycle notifications. DAO governs oracle config and upgrades via Treasury (Timelock).

**Tech Stack:** Solidity 0.8.20, OpenZeppelin 5.x (Upgradeable), Foundry, Go 1.26, abigen

**Spec:** `docs/superpowers/specs/2026-03-15-awpemission-v2-design.md`

---

## Chunk 1: Core Contract Changes

### Task 1: Update IAWPEmission Interface

**Files:**
- Modify: `contracts/src/interfaces/IAWPEmission.sol`

- [ ] **Step 1: Rewrite interface**

Replace the entire file with the V2 interface. Key changes:
- Remove `setWeight()` — replaced by `emergencySetWeight()` and `submitAllocations()`
- Add oracle functions: `submitAllocations`, `setOracleConfig`, `emergencySetWeight`
- Add oracle view functions: `oracles`, `oracleThreshold`, `allocationNonce`, `getOracleCount`
- Add events: `AllocationsSubmitted`, `OracleConfigUpdated`
- Keep all existing view functions and lifecycle functions unchanged

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IAWPEmission {
    // ── 结算 ──
    function settleEpoch() external;

    // ── 预言机批量提交 ──
    function submitAllocations(
        uint256[] calldata subnetIds,
        uint128[] calldata weights,
        bytes[] calldata signatures
    ) external;

    // ── DAO 治理（onlyTimelock） ──
    function emergencySetWeight(uint256 subnetId, uint128 weight) external;
    function setOracleConfig(address[] calldata oracles_, uint256 threshold_) external;
    function setEpochDuration(uint256 d) external;
    function setBatchSize(uint128 b) external;
    function setMaxActiveSubnets(uint128 m) external;

    // ── AWPRegistry 生命周期同步（onlyAWPRegistry） ──
    function registerSubnet(uint256 subnetId, address recipient) external;
    function activateSubnet(uint256 subnetId) external;
    function deactivateSubnet(uint256 subnetId) external;
    function reactivateSubnet(uint256 subnetId) external;
    function removeSubnet(uint256 subnetId) external;

    // ── View ──
    function rootNet() external view returns (address);
    function treasury() external view returns (address);
    function currentEpoch() external view returns (uint256);
    function epochSettling() external view returns (bool);
    function currentDailyEmission() external view returns (uint256);
    function totalWeight() external view returns (uint256);
    function epochDuration() external view returns (uint256);
    function lastSettleTime() external view returns (uint256);
    function getActiveSubnetCount() external view returns (uint256);
    function getActiveSubnetIdAt(uint256 index) external view returns (uint256);
    function oracles(uint256 index) external view returns (address);
    function oracleThreshold() external view returns (uint256);
    function allocationNonce() external view returns (uint256);
    function getOracleCount() external view returns (uint256);

    // ── 事件 ──
    event GovernanceWeightUpdated(uint256 indexed subnetId, uint256 weight);
    event SubnetAWPDistributed(uint256 indexed epoch, uint256 indexed subnetId, uint256 awpAmount);
    event DAOMatchDistributed(uint256 indexed epoch, uint256 amount);
    event EpochSettled(uint256 indexed epoch, uint256 totalEmission, uint256 subnetCount);
    event AllocationsSubmitted(uint256 indexed nonce, uint256[] subnetIds, uint128[] weights);
    event OracleConfigUpdated(address[] oracles, uint256 threshold);
}
```

- [ ] **Step 2: Verify compilation**

Run: `/home/ubuntu/.foundry/bin/forge build`
Expected: Compilation errors in AWPEmission.sol (expected — we haven't rewritten it yet) and test files, but no errors in IAWPEmission.sol itself.

### Task 2: Rewrite AWPEmission.sol as UUPS Proxy

**Files:**
- Rewrite: `contracts/src/token/AWPEmission.sol`

- [ ] **Step 1: Write the full V2 contract**

Complete rewrite. Changes from V1:
- Inherit `Initializable`, `UUPSUpgradeable`, `ReentrancyGuardUpgradeable`, `EIP712Upgradeable`, `IAWPEmission`
- Constructor only calls `_disableInitializers()`
- `initialize()` replaces constructor logic, calls all OZ `__*_init()` functions including `__EIP712_init("AWPEmission", "2")`
- `rootNet` and `awpToken` become regular state variables (not immutable — proxy incompatible)
- `onlyTimelock` modifier only checks `treasury` (no longer accepts `rootNet`)
- `setWeight()` renamed to `emergencySetWeight()`
- New `submitAllocations()` with EIP-712 signature verification
- New `setOracleConfig()` with zero-address and duplicate checks
- New state: `oracles[]`, `oracleThreshold`, `allocationNonce`, `__gap[47]`
- New errors: `OracleNotConfigured`, `InvalidOracleConfig`, `InvalidSignatureCount`, `DuplicateOracle`, `UnknownOracle`, `DuplicateSigner`, `SubnetNotRegistered`
- `_authorizeUpgrade()` restricted to treasury
- `ALLOCATION_TYPEHASH` constant for EIP-712
- `_isOracle()` internal helper using loop (oracle list is small)

Key implementation details for `submitAllocations`:
1. `require(oracleThreshold > 0)` — guard uninitialized config
2. Build struct hash: `keccak256(abi.encode(ALLOCATION_TYPEHASH, keccak256(abi.encodePacked(subnetIds)), keccak256(abi.encodePacked(weights)), allocationNonce))`
3. Compute digest: `_hashTypedDataV4(structHash)`
4. For each sig: `ECDSA.recover(digest, sig)` → check `_isOracle(signer)` → check not in `seen[]` array
5. `require(validCount >= oracleThreshold)`
6. Batch update weights with `recipient != address(0)` guard
7. Increment `allocationNonce`, emit event

Storage layout must exactly match V1 order for existing variables, then append new V2 variables:
```
rootNet, awpToken, treasury, epochDuration, currentEpoch, lastSettleTime,
currentDailyEmission, subnetWeights(mapping), activeSubnetIds(set),
totalWeight, maxActiveSubnets+batchSize(packed), settleIndex,
epochEmissionLocked, _epochSubnetPool, _epochSubnetMinted,
_epochTotalWeight, _epochActiveCount, epochSettling,
oracles, oracleThreshold, allocationNonce, __gap[47]
```

- [ ] **Step 2: Verify compilation**

Run: `/home/ubuntu/.foundry/bin/forge build`
Expected: AWPEmission.sol compiles. Test files and deploy scripts will still fail (expected).

### Task 3: Remove Delegation Functions from AWPRegistry

**Files:**
- Modify: `contracts/src/AWPRegistry.sol`

- [ ] **Step 1: Remove 4 delegation functions**

Delete these functions from AWPRegistry.sol:
- `setGovernanceWeight(uint256 subnetId, uint128 w)` (line 744-749)
- `setMaxActiveSubnets(uint128 m)` (line 765-769)
- `setEpochDuration(uint256 d)` (line 783-787)
- `setBatchSize(uint128 b)` (line 789-793)

Also remove `setGovernanceWeight` and `setMaxActiveSubnets` from `test_onlyTimelockFunctions` test expectations.

- [ ] **Step 2: Verify compilation of AWPRegistry.sol**

Run: `/home/ubuntu/.foundry/bin/forge build`
Expected: AWPRegistry.sol compiles. Tests that call removed functions will fail (expected).

### Task 4: Update Deploy Scripts

**Files:**
- Modify: `contracts/script/Deploy.s.sol`
- Modify: `contracts/script/TestDeploy.s.sol`

- [ ] **Step 1: Update Deploy.s.sol**

Change Step 13 from direct deployment to proxy pattern:
```solidity
// Step 13a: AWPEmission implementation
AWPEmission emissionImpl = new AWPEmission();
// Step 13b: Deploy proxy with initialize call
bytes memory initData = abi.encodeCall(
    AWPEmission.initialize,
    (address(rootNet), address(awp), address(treasury), INITIAL_DAILY_EMISSION, EPOCH_DURATION)
);
ERC1967Proxy emissionProxy = new ERC1967Proxy(address(emissionImpl), initData);
AWPEmission emission = AWPEmission(address(emissionProxy));
```

Add import for `ERC1967Proxy`:
```solidity
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
```

After `rootNet.initializeRegistry(...)`, add oracle config:
```solidity
// Step 18: Configure oracles (placeholder — real oracles set by DAO governance)
// In production, this would be called with actual oracle addresses
```

Note: `awp.addMinter(address(emission))` requires no source change — the `emission` variable is reassigned to point to the proxy address in step 13b, so `addMinter` automatically receives the proxy address.

- [ ] **Step 2: Update TestDeploy.s.sol**

Same proxy pattern but simplified. Use the same `ERC1967Proxy` approach.

- [ ] **Step 3: Verify deploy scripts compile**

Run: `/home/ubuntu/.foundry/bin/forge build`
Expected: All source files and scripts compile.

- [ ] **Step 4: Commit**

```bash
git add contracts/src/ contracts/script/
git commit -m "feat: AWPEmission V2 — UUPS proxy + oracle consensus + remove AWPRegistry delegation"
```

## Chunk 2: Test Updates

### Task 5: Rewrite AWPEmission.t.sol

**Files:**
- Rewrite: `contracts/test/AWPEmission.t.sol`

- [ ] **Step 1: Write complete test file**

Test setup deploys via proxy pattern:
```solidity
AWPEmission emissionImpl = new AWPEmission();
bytes memory initData = abi.encodeCall(AWPEmission.initialize, (rootNet, address(awpToken), treasury, INITIAL_DAILY_EMISSION, EPOCH_DURATION));
ERC1967Proxy proxy = new ERC1967Proxy(address(emissionImpl), initData);
emission = AWPEmission(address(proxy));
```

Required tests (each is a separate `function test_*`):

**Oracle config tests:**
- `test_setOracleConfig` — treasury sets 5 oracles with threshold 3, verify storage
- `test_setOracleConfig_revertsForNonTimelock` — user calls, expect `NotTimelock`
- `test_setOracleConfig_revertsForZeroThreshold` — threshold=0, expect `InvalidOracleConfig`
- `test_setOracleConfig_revertsForThresholdExceedsOracles` — threshold > oracles.length
- `test_setOracleConfig_revertsForZeroAddress` — oracle list contains address(0)
- `test_setOracleConfig_revertsForDuplicateOracle` — duplicate address in list

**submitAllocations tests:**
- `test_submitAllocations` — valid 3/5 oracle signatures, 2 subnets, verify weights updated and totalWeight adjusted
- `test_submitAllocations_revertsBeforeOracleConfig` — call before setOracleConfig, expect `OracleNotConfigured`
- `test_submitAllocations_revertsBelowThreshold` — only 2/3 valid sigs, expect `InvalidSignatureCount`
- `test_submitAllocations_revertsDuplicateSigner` — same oracle signs twice, expect `DuplicateSigner`
- `test_submitAllocations_revertsUnknownOracle` — sig from non-oracle address, expect `UnknownOracle`
- `test_submitAllocations_revertsUnregisteredSubnet` — subnetId with no recipient, expect `SubnetNotRegistered`
- `test_submitAllocations_nonceIncrement` — submit twice, verify nonce increments and first sigs can't replay

**emergencySetWeight tests:**
- `test_emergencySetWeight` — treasury sets weight, verify storage
- `test_emergencySetWeight_revertsForNonTimelock` — user calls, expect `NotTimelock`

**Existing emission tests (adapted for proxy):**
- `test_settleEpoch` — same as V1 but through proxy
- `test_settleEpochDecay` — verify decay after epoch 0 (no decay) and epoch 1 (decayed)
- `test_settleEpochNoActiveSubnets` — all emission to DAO
- `test_settleEpochTooEarly` — expect `EpochNotReady`
- `test_multiSubnetDistribution` — 2 subnets, 3:1 weight ratio
- `test_notSettling` — batchSize=0, verify submitAllocations blocked during settlement

**Lifecycle tests (same as V1):**
- `test_registerSubnet`, `test_activateSubnet`, `test_deactivateSubnet`, etc.

**Upgrade test:**
- `test_upgradeViaTimelock` — treasury calls `upgradeToAndCall(newImpl, "")`
- `test_upgrade_revertsForNonTimelock` — user calls, expect `NotTimelock`

Helper: `_signAllocations(uint256 pk, uint256[] subnetIds, uint128[] weights, uint256 nonce)` returns `bytes memory` — uses `vm.sign` with EIP-712 digest.

- [ ] **Step 2: Run tests**

Run: `/home/ubuntu/.foundry/bin/forge test --match-contract AWPEmissionTest -v`
Expected: All tests pass.

### Task 6: Update AWPRegistry.t.sol

**Files:**
- Modify: `contracts/test/AWPRegistry.t.sol`

- [ ] **Step 1: Update setUp for proxy deployment**

Change AWPEmission deployment from `new AWPEmission(...)` to proxy pattern.

- [ ] **Step 2: Remove tests for deleted functions**

Delete:
- `test_setMaxActiveSubnets` (tests `rootNet.setMaxActiveSubnets` which no longer exists)
- `test_setEpochDuration` (tests `rootNet.setEpochDuration` which no longer exists)

Update `test_onlyTimelockFunctions`:
- Remove `rootNet.setGovernanceWeight(1, 100)` expectation (function removed)
- Keep other timelock function checks (`setInitialAlphaPrice`, `setGuardian`, `unpause`)

- [ ] **Step 3: Update test_reallocate**

`rootNet.settleEpoch()` was already changed to `emission.settleEpoch()` in prior refactoring. Verify it still references `emission.settleEpoch()`. The `setGovernanceWeight` call is no longer needed since V2 uses oracle submission — but reallocate test doesn't set weights, it just needs a subnet to be active. No weight needed for reallocate (just tests stake movement).

- [ ] **Step 4: Run tests**

Run: `/home/ubuntu/.foundry/bin/forge test --match-contract AWPRegistryTest -v`
Expected: All tests pass.

- [ ] **Step 5: Commit**

```bash
git add contracts/test/AWPEmission.t.sol contracts/test/AWPRegistry.t.sol
git commit -m "test: rewrite AWPEmission tests for V2 proxy + oracle, update AWPRegistry tests"
```

### Task 7: Update E2E.t.sol

**Files:**
- Modify: `contracts/test/E2E.t.sol`

- [ ] **Step 1: Update _deploy() for proxy pattern**

Change AWPEmission deployment to proxy in `_deploy()`.

- [ ] **Step 2: Add oracle setup helper**

```solidity
uint256 oracle1Pk = 0xA1;
uint256 oracle2Pk = 0xA2;
uint256 oracle3Pk = 0xA3;
address oracle1 = vm.addr(oracle1Pk);
address oracle2 = vm.addr(oracle2Pk);
address oracle3 = vm.addr(oracle3Pk);
```

In `_deploy()`, after `rootNet.initializeRegistry(...)`, add:
```solidity
address[] memory oracleList = new address[](3);
oracleList[0] = oracle1; oracleList[1] = oracle2; oracleList[2] = oracle3;
emission.setOracleConfig(oracleList, 2); // 2/3 threshold
```

Add helper `_submitWeights(uint256[] memory ids, uint128[] memory weights)` that signs with oracle1 and oracle2, then calls `emission.submitAllocations(...)`.

- [ ] **Step 3: Update test_e2e_daoGovernanceWeight**

Change from DAO proposal targeting `rootNet.setGovernanceWeight` to:
- DAO proposal targeting `emission.emergencySetWeight(sid, 500)`
- Update `targets[0] = address(emission)` and `calldatas[0] = abi.encodeCall(emission.emergencySetWeight, (sid, uint128(500)))`

- [ ] **Step 4: Replace all `rootNet.setGovernanceWeight` calls with oracle submission**

Throughout E2E tests, replace:
```solidity
vm.prank(address(treasury));
rootNet.setGovernanceWeight(sid, uint128(100));
```
with oracle-signed `submitAllocations` call using the helper.

- [ ] **Step 5: Update test_e2e_notSettlingGuard**

Replace `emission.setBatchSize(0)` and `emission.setBatchSize(200)` with direct treasury calls to `emission`:
```solidity
vm.prank(address(treasury));
emission.setBatchSize(0);
```
(This already works since `onlyTimelock` checks `msg.sender == treasury`.)

The `notSettling` guard test should also verify `submitAllocations` is blocked during settlement.

- [ ] **Step 6: Add test_e2e_upgradeViaDAOProposal**

DAO proposes `emission.upgradeToAndCall(newImpl, "")`, goes through vote + queue + execute. Verify upgrade succeeds and emission state is preserved (currentEpoch, totalWeight unchanged). Deploy a trivial V2 impl (same code, just a new deployment) for the test.

- [ ] **Step 7: Run tests**

Run: `/home/ubuntu/.foundry/bin/forge test --match-contract E2ETest -v`
Expected: All tests pass.

### Task 8: Update Integration.t.sol

**Files:**
- Modify: `contracts/test/Integration.t.sol`

- [ ] **Step 1: Update setUp for proxy deployment + oracle setup**

Same proxy pattern as E2E. Add oracle private keys and setup.

- [ ] **Step 2: Replace rootNet.setGovernanceWeight calls**

Replace all `rootNet.setGovernanceWeight(subnetId, uint128(N))` with oracle-signed `submitAllocations`.

- [ ] **Step 3: Run tests**

Run: `/home/ubuntu/.foundry/bin/forge test --match-contract IntegrationTest -v`
Expected: All tests pass.

- [ ] **Step 4: Run full test suite**

Run: `/home/ubuntu/.foundry/bin/forge test`
Expected: All tests pass (256+ tests).

- [ ] **Step 5: Commit**

```bash
git add contracts/test/E2E.t.sol contracts/test/Integration.t.sol
git commit -m "test: update E2E and Integration tests for AWPEmission V2 oracle flow"
```

## Chunk 3: Go Backend + Bindings + Docs

### Task 9: Regenerate Go Bindings

**Files:**
- Regenerate: `api/internal/chain/bindings/a_w_p_emission.go`
- Regenerate: `api/internal/chain/bindings/root_net.go`

- [ ] **Step 1: Regenerate AWPEmission binding**

```bash
cd /home/ubuntu/code/Cortexia/contracts
/home/ubuntu/.foundry/bin/forge inspect AWPEmission abi --json > /tmp/awp_emission_abi.json
/home/ubuntu/go/bin/abigen --abi /tmp/awp_emission_abi.json --pkg bindings --type AWPEmission --out /home/ubuntu/code/Cortexia/api/internal/chain/bindings/a_w_p_emission.go
```

- [ ] **Step 2: Regenerate AWPRegistry binding**

```bash
/home/ubuntu/.foundry/bin/forge inspect AWPRegistry abi --json > /tmp/rootnet_abi.json
/home/ubuntu/go/bin/abigen --abi /tmp/rootnet_abi.json --pkg bindings --type AWPRegistry --out /home/ubuntu/code/Cortexia/api/internal/chain/bindings/root_net.go
```

- [ ] **Step 3: Verify Go build**

Run: `cd /home/ubuntu/code/Cortexia/api && go build ./...`
Expected: Build succeeds.

### Task 10: Update Indexer for New Events

**Files:**
- Modify: `api/internal/chain/indexer.go`

- [ ] **Step 1: Add AllocationsSubmitted event parsing**

After the `EpochSettled` parsing block, add:

```go
// AllocationsSubmitted (from AWPEmission)
if evt, err := awpEmission.ParseAllocationsSubmitted(lg); err == nil {
    return []redisEvent{makeEvent("AllocationsSubmitted", lg, map[string]interface{}{
        "nonce":     evt.Nonce.String(),
        "subnetIds": bigIntSliceToStrings(evt.SubnetIds),
        "weights":   uint128SliceToStrings(evt.Weights),
    })}, nil
}
```

- [ ] **Step 2: Add OracleConfigUpdated event parsing**

```go
// OracleConfigUpdated (from AWPEmission)
if evt, err := awpEmission.ParseOracleConfigUpdated(lg); err == nil {
    return []redisEvent{makeEvent("OracleConfigUpdated", lg, map[string]interface{}{
        "oracles":   addressSliceToStrings(evt.Oracles),
        "threshold": evt.Threshold.String(),
    })}, nil
}
```

Add helper functions for slice conversions if not already present.

- [ ] **Step 3: Verify Go build**

Run: `cd /home/ubuntu/code/Cortexia/api && go build ./...`
Expected: Build succeeds.

- [ ] **Step 4: Commit**

```bash
git add api/
git commit -m "feat: regenerate bindings for AWPEmission V2, add new event parsing"
```

### Task 11: Update Documentation

**Files:**
- Modify: `CLAUDE.md`
- Modify: `docs/architecture.md`

- [ ] **Step 1: Update CLAUDE.md**

Update the Core Architecture section:
- AWPEmission description: Add "UUPS upgradeable, multi-oracle weight submission"
- Remove references to `setGovernanceWeight` on AWPRegistry
- Update deployment sequence to show proxy deployment

- [ ] **Step 2: Update docs/architecture.md**

Update AWPEmission section (4.4) to describe V2:
- UUPS proxy pattern
- `submitAllocations` with oracle consensus
- `emergencySetWeight` for DAO override
- Oracle config management

Update AWPRegistry section:
- Remove `setGovernanceWeight`, `setEpochDuration`, `setBatchSize`, `setMaxActiveSubnets`

Update deployment steps:
- Step 13 split into 13a (impl) + 13b (proxy)

- [ ] **Step 3: Commit**

```bash
git add CLAUDE.md docs/
git commit -m "docs: update architecture for AWPEmission V2 oracle consensus"
```

### Task 12: Final Verification

- [ ] **Step 1: Full Solidity build**

Run: `/home/ubuntu/.foundry/bin/forge build`
Expected: Clean compilation.

- [ ] **Step 2: Full Solidity test suite**

Run: `/home/ubuntu/.foundry/bin/forge test`
Expected: All tests pass.

- [ ] **Step 3: Full Go build**

Run: `cd /home/ubuntu/code/Cortexia/api && go build ./...`
Expected: Clean build.

- [ ] **Step 4: Go tests**

Run: `cd /home/ubuntu/code/Cortexia/api && go test ./...`
Expected: No new failures (pre-existing handler_test failures are OK).
