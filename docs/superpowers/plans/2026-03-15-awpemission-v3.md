# AWPEmission V3 Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor AWPEmission to a generic address→weight distribution engine, removing all subnet concepts and simplifying settlement with a call-parameter batch limit.

**Architecture:** AWPEmission V3 is a UUPS proxy that manages `(address → weight)` pairs via oracle multi-sig full-replacement submissions. Settlement uses `settleEpoch(limit)` with progress tracked by a single counter. RootNet self-maintains `activeSubnetIds` and no longer calls AWPEmission for lifecycle events.

**Tech Stack:** Solidity 0.8.24, OpenZeppelin 5.x (Upgradeable), Foundry, Go 1.26, abigen

**Spec:** `docs/superpowers/specs/2026-03-15-awpemission-v3-design.md`

---

## Chunk 1: Core Contract Changes

### Task 1: Rewrite IAWPEmission Interface

**Files:**
- Rewrite: `contracts/src/interfaces/IAWPEmission.sol`

- [ ] **Step 1: Write V3 interface**

Replace the entire file. Key changes from V2:
- `submitAllocations` params change from `(uint256[] subnetIds, uint128[] weights, bytes[] sigs)` to `(address[] recipients, uint128[] weights, bytes[] sigs)`
- `settleEpoch()` → `settleEpoch(uint256 limit)`
- Remove: `registerSubnet`, `activateSubnet`, `deactivateSubnet`, `reactivateSubnet`, `removeSubnet`, `setBatchSize`, `setMaxActiveSubnets`, `getActiveSubnetCount`, `getActiveSubnetIdAt`, `epochSettling`
- Add: `settleProgress`, `getRecipientCount`, `maxRecipients`
- Events: `SubnetAWPDistributed` second param changes from `uint256 indexed subnetId` to `address indexed recipient`; `AllocationsSubmitted` changes from `uint256[] subnetIds` to `address[] recipients`

Use the exact interface from the spec (section 3).

### Task 2: Rewrite AWPEmission.sol

**Files:**
- Rewrite: `contracts/src/token/AWPEmission.sol`

- [ ] **Step 1: Write the full V3 contract**

Complete rewrite from V2. Inheritance unchanged: `Initializable, UUPSUpgradeable, ReentrancyGuardUpgradeable, EIP712Upgradeable, IAWPEmission`.

**Key differences from V2:**

Storage (spec section 1, slots 0-19 + gap):
- Remove: `SubnetWeight` struct, `subnetWeights` mapping, `activeSubnetIds` EnumerableSet, `maxActiveSubnets`, `batchSize`, `epochSettling`, `settleIndex`, `_epochTotalWeight`, `_epochActiveCount`, `_epochSubnetPool`, `_epochSubnetMinted`
- Add: `recipients` (address[]), `weights` (mapping(address=>uint128)), `totalWeight`, `settleProgress`, `_snapshotLen`, `_snapshotWeight`, `_snapshotPool`, `_epochMinted`, `maxRecipients` (default 10000)
- Remove: `EnumerableSet` import (no longer needed)

Functions removed:
- `registerSubnet`, `activateSubnet`, `deactivateSubnet`, `reactivateSubnet`, `removeSubnet`
- `onlyRootNet` modifier
- `setBatchSize`, `setMaxActiveSubnets`

Functions modified:
- `submitAllocations(address[] recipients_, uint128[] weights_, bytes[] sigs)`:
  - Full replacement: clear old weights (loop delete), write new ones
  - Require `settleProgress == 0`
  - Require `recipients_.length <= maxRecipients`
  - Check no address(0), no duplicates
  - Duplicate check: during the write phase, use a sentinel — write `weights[addr] = w + 1` temporarily (shifted by 1), so even weight=0 produces a non-zero storage value. After write loop, shift back. Or simpler: just require weight > 0 for all entries — submitting weight=0 is pointless in a full-replacement model (just exclude the address). This aligns with the spec's intent: "no stale entries."
  - Require `weights_[i] > 0` for all entries. Duplicate check: `weights[addr] != 0` after clear means duplicate.
  - ALLOCATION_TYPEHASH changes: `"SubmitAllocations(address[] recipients,uint128[] weights,uint256 nonce)"`
  - EIP-712 struct hash: `keccak256(abi.encode(TYPEHASH, keccak256(abi.encodePacked(recipients_)), keccak256(abi.encodePacked(weights_)), allocationNonce))`

- `settleEpoch(uint256 limit)`:
  - `require(limit > 0)` — InvalidParameter
  - Phase 1 (settleProgress == 0): snapshot `_snapshotLen`, `_snapshotWeight`, `_snapshotPool`, set `settleProgress = 1`
  - Phase 2: `start = settleProgress - 1`, `end = min(start + limit, _snapshotLen)`, iterate and mint
  - Guard: `if share > 0` before `_mintTo` call
  - Phase 3: `end >= _snapshotLen` → mint DAO share, reset, epoch++
  - Use local `epoch` var for EpochSettled event (not `currentEpoch - 1`)

- `emergencySetWeight(address addr, uint128 weight)`:
  - Require `settleProgress == 0`
  - Loop `recipients` to verify `addr` exists (revert `RecipientNotFound`)
  - Update `weights[addr]`, adjust `totalWeight`

- `setOracleConfig`: add `require(settleProgress == 0)`

- `setEpochDuration`: unchanged (kept `require(d > 0)`)

- `initialize()`: add `maxRecipients = 10000` initialization

Errors: Replace `CurrentlySettling` with `SettlementInProgress`, remove `NotRootNet`, `MaxActiveSubnetsReached`, `SubnetNotRegistered`. Add `DuplicateRecipient`, `RecipientNotFound`.

### Task 3: Update RootNet.sol

**Files:**
- Modify: `contracts/src/RootNet.sol`

- [ ] **Step 1: Restore activeSubnetIds and remove AWPEmission lifecycle calls**

Changes:
1. Re-add `EnumerableSet` import and `using EnumerableSet for EnumerableSet.UintSet`
2. Add `EnumerableSet.UintSet private activeSubnetIds` state variable
3. Add `uint128 public constant MAX_ACTIVE_SUBNETS = 10000`

4. `registerSubnet`: Remove `IAWPEmission(awpEmission).registerSubnet(subnetId, params.subnetContract)` call

5. `activateSubnet`: Replace `IAWPEmission(awpEmission).activateSubnet(subnetId)` with:
   ```solidity
   if (activeSubnetIds.length() >= MAX_ACTIVE_SUBNETS) revert MaxActiveSubnetsReached();
   activeSubnetIds.add(subnetId);
   ```
   Re-add `MaxActiveSubnetsReached` error.

6. `pauseSubnet`: Replace `IAWPEmission(awpEmission).deactivateSubnet(subnetId)` with `activeSubnetIds.remove(subnetId)`

7. `resumeSubnet`: Replace `IAWPEmission(awpEmission).reactivateSubnet(subnetId)` with:
   ```solidity
   if (activeSubnetIds.length() >= MAX_ACTIVE_SUBNETS) revert MaxActiveSubnetsReached();
   activeSubnetIds.add(subnetId);
   ```

8. `banSubnet`: Replace `IAWPEmission(awpEmission).deactivateSubnet(subnetId)` with `activeSubnetIds.remove(subnetId)` (only when status was Active)

9. `unbanSubnet`: Replace `IAWPEmission(awpEmission).reactivateSubnet(subnetId)` with:
   ```solidity
   if (activeSubnetIds.length() >= MAX_ACTIVE_SUBNETS) revert MaxActiveSubnetsReached();
   activeSubnetIds.add(subnetId);
   ```

10. `deregisterSubnet`: Replace `IAWPEmission(awpEmission).removeSubnet(subnetId)` with:
    ```solidity
    if (activeSubnetIds.contains(subnetId)) {
        activeSubnetIds.remove(subnetId);
    }
    ```

11. `getActiveSubnetCount()`: Change from `IAWPEmission(awpEmission).getActiveSubnetCount()` to `activeSubnetIds.length()`

12. Add `getActiveSubnetIdAt(uint256 index)`: return `activeSubnetIds.at(index)`

### Task 4: Update Deploy Scripts and Signing Helper

**Files:**
- Modify: `contracts/script/Deploy.s.sol`
- Modify: `contracts/script/TestDeploy.s.sol`
- Modify: `contracts/test/helpers/EmissionSigningHelper.sol`

- [ ] **Step 1: Update EmissionSigningHelper typehash**

Change `ALLOCATION_TYPEHASH` from:
```
"SubmitAllocations(uint256[] subnetIds,uint128[] weights,uint256 nonce)"
```
to:
```
"SubmitAllocations(address[] recipients,uint128[] weights,uint256 nonce)"
```

Change `_signAllocations` first parameter from `uint256[] memory subnetIds` to `address[] memory recipients`.

Change `keccak256(abi.encodePacked(subnetIds))` to `keccak256(abi.encodePacked(recipients))`.

- [ ] **Step 2: Update Deploy.s.sol if needed**

The `initialize()` params are unchanged. The only change: remove any `setBatchSize` or `setMaxActiveSubnets` calls if present. Check and clean.

- [ ] **Step 3: Verify compilation**

Run: `cd /home/ubuntu/code/Cortexia/contracts && /home/ubuntu/.foundry/bin/forge build 2>&1`
Source files and scripts should compile. Tests will fail (expected).

- [ ] **Step 4: Commit**

```bash
git add contracts/src/ contracts/script/ contracts/test/helpers/
git commit -m "feat: AWPEmission V3 — generic distribution engine, remove subnet coupling, simplify settlement"
```

## Chunk 2: Test Updates

### Task 5: Rewrite AWPEmission.t.sol

**Files:**
- Rewrite: `contracts/test/AWPEmission.t.sol`

- [ ] **Step 1: Write complete test file**

setUp() deploys via proxy (same as V2). Oracle setup same as V2.

Key differences from V2 tests:
- `_signAllocations` now takes `address[] memory recipients` instead of `uint256[] memory subnetIds`
- `submitAllocations` takes `address[]` recipients directly (no registerSubnet needed)
- `settleEpoch(200)` with limit parameter
- No lifecycle tests (registerSubnet/activateSubnet etc. removed from AWPEmission)
- `emergencySetWeight` takes `address` not `uint256 subnetId`

Required tests:

**submitAllocations (full replacement):**
- `test_submitAllocations`: submit 2 recipients with weights 300/100, verify storage
- `test_submitAllocations_fullReplacement`: submit once, then submit again with different recipients, verify old weights cleared
- `test_submitAllocations_revertsBeforeOracleConfig`
- `test_submitAllocations_revertsBelowThreshold`
- `test_submitAllocations_revertsDuplicateSigner`
- `test_submitAllocations_revertsUnknownOracle`
- `test_submitAllocations_revertsExceedsMaxRecipients`
- `test_submitAllocations_revertsDuplicateRecipient`
- `test_submitAllocations_nonceIncrement`
- `test_submitAllocations_revertsNonceReplay`: submit once (nonce=0), then replay same signatures (nonce=0 again), expect revert because nonce is now 1
- `test_submitAllocations_revertsDuringSettlement`: start settlement with limit=1, try submit, expect revert

**settleEpoch:**
- `test_settleEpoch`: single call processes all recipients
- `test_settleEpochBatched_oneByOne`: limit=1 with 3 recipients, call 3+1 times
- `test_settleEpochBatched_twoBatch`: limit=2 with 4 recipients, call 2+1 times
- `test_settleEpochDecay`: verify no decay epoch 0, decay epoch 1
- `test_settleEpochNoRecipients`: all emission to DAO
- `test_settleEpochTooEarly`
- `test_settleEpoch_revertsLimitZero`

**emergencySetWeight:**
- `test_emergencySetWeight`
- `test_emergencySetWeight_revertsForNonTimelock`
- `test_emergencySetWeight_revertsRecipientNotFound`
- `test_emergencySetWeight_revertsDuringSettlement`

**Oracle config, upgrade (same as V2):**
- `test_setOracleConfig`, revert cases
- `test_upgradeViaTimelock`, `test_upgrade_revertsForNonTimelock`

- [ ] **Step 2: Run tests**

Run: `/home/ubuntu/.foundry/bin/forge test --match-contract AWPEmissionTest -v`

### Task 6: Update RootNet.t.sol

**Files:**
- Modify: `contracts/test/RootNet.t.sol`

- [ ] **Step 1: Verify lifecycle tests work with restored activeSubnetIds**

The setUp already deploys AWPEmission via proxy. After removing AWPEmission lifecycle calls from RootNet, the existing lifecycle tests (test_activateSubnet, test_pauseAndResumeSubnet, test_banAndUnbanSubnet, test_deregisterSubnet) should still work since they only check RootNet state. Verify `getActiveSubnetCount` returns correct values.

The `test_reallocate` test calls `emission.settleEpoch()` — update to `emission.settleEpoch(200)`.

- [ ] **Step 2: Run tests**

Run: `/home/ubuntu/.foundry/bin/forge test --match-contract RootNetTest -v`

### Task 7: Update E2E.t.sol

**Files:**
- Modify: `contracts/test/E2E.t.sol`

- [ ] **Step 1: Update signing helper and submitAllocations calls**

- `_signAllocations` now takes `address[]` (via updated EmissionSigningHelper)
- `_submitWeights` and `_submitWeight` helpers: change to use `address[]` recipients
- All `_submitWeight(sid, weight)` → `_submitWeight(subnetContractAddr, weight)`
- Must look up the actual subnet contract address from the test context, not subnetId
- `emission.settleEpoch()` → `emission.settleEpoch(200)` throughout
- `_settleEpoch()` helper: update to call `emission.settleEpoch(200)`
- `emission.emergencySetWeight(sid, weight)` → `emission.emergencySetWeight(subnetContractAddr, weight)`

- [ ] **Step 2: Run tests**

Run: `/home/ubuntu/.foundry/bin/forge test --match-contract E2ETest -v`

### Task 8: Update Integration.t.sol

**Files:**
- Modify: `contracts/test/Integration.t.sol`

- [ ] **Step 1: Same changes as E2E**

Update signing helpers, submitAllocations calls, settleEpoch(200), emergencySetWeight with addresses.

- [ ] **Step 2: Run full test suite**

Run: `/home/ubuntu/.foundry/bin/forge test`
All tests must pass.

- [ ] **Step 3: Commit**

```bash
git add contracts/test/
git commit -m "test: update all tests for AWPEmission V3 address-based distribution"
```

## Chunk 3: Go Backend + Docs

### Task 9: Regenerate Bindings and Update Keeper

**Files:**
- Regenerate: `api/internal/chain/bindings/a_w_p_emission.go`
- Regenerate: `api/internal/chain/bindings/root_net.go`
- Modify: `api/internal/chain/keeper.go`

- [ ] **Step 1: Regenerate bindings**

```bash
cd /home/ubuntu/code/Cortexia/contracts
/home/ubuntu/.foundry/bin/forge inspect AWPEmission abi --json > /tmp/awp_emission_abi.json
/home/ubuntu/go/bin/abigen --abi /tmp/awp_emission_abi.json --pkg bindings --type AWPEmission --out /home/ubuntu/code/Cortexia/api/internal/chain/bindings/a_w_p_emission.go
/home/ubuntu/.foundry/bin/forge inspect RootNet abi --json > /tmp/rootnet_abi.json
/home/ubuntu/go/bin/abigen --abi /tmp/rootnet_abi.json --pkg bindings --type RootNet --out /home/ubuntu/code/Cortexia/api/internal/chain/bindings/root_net.go
```

- [ ] **Step 2: Update keeper.go**

Change `trySettleEpoch`:
- Replace `k.awpEmission.EpochSettling(nil)` with `k.awpEmission.SettleProgress(nil)`
- Check `settleProgress.Cmp(big.NewInt(0)) > 0` instead of `settling == true`

Change `sendSettleEpoch`:
- `k.awpEmission.SettleEpoch(txAuth)` → `k.awpEmission.SettleEpoch(txAuth, big.NewInt(200))`

- [ ] **Step 3: Verify Go build**

Run: `cd /home/ubuntu/code/Cortexia/api && go build ./...`

### Task 10: Update Indexer

**Files:**
- Modify: `api/internal/chain/indexer.go`

- [ ] **Step 1: Update event parsing**

The `AllocationsSubmitted` event data changes: `SubnetIds` field becomes `Recipients` (type `[]common.Address`). The `SubnetAWPDistributed` event second indexed param changes from `SubnetId` to `Recipient` (type `common.Address`).

Update the parsing code to match the new binding struct field names. The Go binding auto-generates field names from event parameter names.

- [ ] **Step 2: Verify Go build**

Run: `cd /home/ubuntu/code/Cortexia/api && go build ./...`

- [ ] **Step 3: Commit**

```bash
git add api/
git commit -m "feat: update keeper and indexer for AWPEmission V3 address-based settlement"
```

### Task 11: Final Verification

- [ ] **Step 1: Full Solidity test suite**

Run: `cd /home/ubuntu/code/Cortexia/contracts && /home/ubuntu/.foundry/bin/forge test`
Expected: All tests pass.

- [ ] **Step 2: Full Go build**

Run: `cd /home/ubuntu/code/Cortexia/api && go build ./...`
Expected: Clean build.
