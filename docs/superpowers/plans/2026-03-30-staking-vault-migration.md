# Staking Allocation Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move allocate/deallocate/reallocate (+ gasless variants) from AWPRegistry to StakingVault. Reduces AWPRegistry bytecode by ~3-4KB, shortens call chain, improves separation of concerns.

**Architecture:** StakingVault becomes a UUPS proxy with EIP-712 support. It reads delegate authorization from AWPRegistry via cross-contract call. AWPRegistry removes 5 allocation functions + 2 typehash constants. Relay and indexer update to target StakingVault.

**Tech Stack:** Solidity 0.8.24, Foundry, OpenZeppelin Upgradeable 5.x

---

### Task 1: StakingVault → UUPS Upgradeable with EIP-712

**Files:**
- Modify: `contracts/src/core/StakingVault.sol`
- Modify: `contracts/src/interfaces/IStakingVault.sol`

- [ ] **Step 1: Add UUPS + EIP-712 imports and inheritance**

Change StakingVault from a plain contract to UUPS upgradeable:

```solidity
// Add imports:
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {EIP712Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

// Change declaration:
contract StakingVault is Initializable, UUPSUpgradeable, EIP712Upgradeable {
```

- [ ] **Step 2: Replace constructor with initialize**

Remove the existing constructor. Add:

```solidity
constructor() { _disableInitializers(); }

function initialize(address awpRegistry_, address stakeNFT_) external initializer {
    __UUPSUpgradeable_init();
    __EIP712_init("StakingVault", "1");
    awpRegistry = awpRegistry_;
    stakeNFT = IStakeNFT(stakeNFT_);
}
```

Change `awpRegistry` and `stakeNFT` from `immutable` to regular storage variables. Remove `setStakeNFT` function (both set in initialize now).

- [ ] **Step 3: Add nonces, typehashes, delegate check, _authorizeUpgrade**

```solidity
mapping(address => uint256) public nonces;

bytes32 private constant ALLOCATE_TYPEHASH = keccak256("Allocate(address staker,address agent,uint256 subnetId,uint256 amount,uint256 nonce,uint256 deadline)");
bytes32 private constant DEALLOCATE_TYPEHASH = keccak256("Deallocate(address staker,address agent,uint256 subnetId,uint256 amount,uint256 nonce,uint256 deadline)");

error ExpiredSignature();
error InvalidSignature();
error NotAuthorized();

// Read delegate authorization from AWPRegistry
function _isAuthorized(address staker, address caller) internal view returns (bool) {
    return caller == staker || IAWPRegistry(awpRegistry).delegates(staker, caller);
}

function _verifyDigest(address user, bytes32 structHash, uint256 deadline, uint8 v, bytes32 r, bytes32 s) internal view {
    if (block.timestamp > deadline) revert ExpiredSignature();
    bytes32 digest = _hashTypedDataV4(structHash);
    if (ECDSA.recover(digest, v, r, s) != user) revert InvalidSignature();
}

function _authorizeUpgrade(address) internal view override {
    if (msg.sender != IAWPRegistry(awpRegistry).treasury()) revert NotAuthorized();
}

uint256[47] private __gap;
```

- [ ] **Step 4: Make allocate/deallocate/reallocate public with auth**

Change `onlyAWPRegistry` modifier to `_isAuthorized` check:

```solidity
function allocate(address staker, address agent, uint256 subnetId, uint256 amount) external {
    if (!_isAuthorized(staker, msg.sender)) revert NotAuthorized();
    // ... existing allocation logic
    emit Allocated(staker, agent, subnetId, amount, msg.sender);
}
```

Same for `deallocate` and `reallocate`. Add events:

```solidity
event Allocated(address indexed staker, address indexed agent, uint256 subnetId, uint256 amount, address operator);
event Deallocated(address indexed staker, address indexed agent, uint256 subnetId, uint256 amount, address operator);
event Reallocated(address indexed staker, address fromAgent, uint256 fromSubnetId, address toAgent, uint256 toSubnetId, uint256 amount, address operator);
```

- [ ] **Step 5: Add gasless allocateFor / deallocateFor**

```solidity
function allocateFor(
    address staker, address agent, uint256 subnetId, uint256 amount, uint256 deadline,
    uint8 v, bytes32 r, bytes32 s
) external {
    _verifyDigest(staker, keccak256(abi.encode(ALLOCATE_TYPEHASH, staker, agent, subnetId, amount, nonces[staker]++, deadline)), deadline, v, r, s);
    _allocate(staker, agent, subnetId, amount);
    emit Allocated(staker, agent, subnetId, amount, msg.sender);
}

function deallocateFor(
    address staker, address agent, uint256 subnetId, uint256 amount, uint256 deadline,
    uint8 v, bytes32 r, bytes32 s
) external {
    _verifyDigest(staker, keccak256(abi.encode(DEALLOCATE_TYPEHASH, staker, agent, subnetId, amount, nonces[staker]++, deadline)), deadline, v, r, s);
    _deallocate(staker, agent, subnetId, amount);
    emit Deallocated(staker, agent, subnetId, amount, msg.sender);
}
```

Extract `_allocate` and `_deallocate` internal functions from the existing logic.

- [ ] **Step 6: Update IStakingVault interface**

Remove `onlyAWPRegistry` comments. Add `allocateFor`, `deallocateFor`, events, `nonces`, `initialize`.

- [ ] **Step 7: Build and test**

```bash
cd /home/ubuntu/code/Cortexia/contracts && forge build && forge test
```

- [ ] **Step 8: Commit**

---

### Task 2: Remove Allocation Functions from AWPRegistry

**Files:**
- Modify: `contracts/src/AWPRegistry.sol`

- [ ] **Step 1: Remove these functions from AWPRegistry**

Delete entirely:
- `allocate(staker, agent, subnetId, amount)`
- `deallocate(staker, agent, subnetId, amount)`
- `allocateFor(...)`
- `deallocateFor(...)`
- `reallocate(...)`

Delete constants:
- `ALLOCATE_TYPEHASH`
- `DEALLOCATE_TYPEHASH`

Delete events (moved to StakingVault):
- `Allocated`
- `Deallocated`
- `Reallocated`

Keep: `delegates` mapping, `grantDelegate`, `revokeDelegate` (account system stays).

- [ ] **Step 2: Add public `treasury()` getter if not already public**

StakingVault's `_authorizeUpgrade` needs to read `AWPRegistry.treasury`. Verify it's already `address public treasury` (it is).

Also ensure `delegates(address,address)` is public (it is: `mapping(address => mapping(address => bool)) public delegates`).

- [ ] **Step 3: Build — check bytecode size**

```bash
forge build --sizes 2>&1 | grep "AWPRegistry"
```

Expected: ~19-20KB (down from 23.1KB).

- [ ] **Step 4: Commit**

---

### Task 3: Update Deploy.s.sol

**Files:**
- Modify: `contracts/script/Deploy.s.sol`

- [ ] **Step 1: Deploy StakingVault as impl + ERC1967Proxy**

Change from direct CREATE2 to proxy pattern (like AWPRegistry and AWPEmission):

```solidity
// Step 8: StakingVault (UUPS proxy)
bytes32 saltVaultImpl = _readSalt("SALT_STAKING_VAULT_IMPL");
StakingVault vaultImpl = StakingVault(_create2(saltVaultImpl, abi.encodePacked(type(StakingVault).creationCode)));
bytes memory vaultInitData = abi.encodeCall(StakingVault.initialize, (address(awpRegistry), address(stakeNft)));
StakingVault vault = StakingVault(_create2(saltVault, abi.encodePacked(type(ERC1967Proxy).creationCode, abi.encode(address(vaultImpl), vaultInitData))));
```

Wait — StakeNFT is deployed AFTER StakingVault in the current sequence. Need to reorder or use `initializeRegistry` to set stakeNFT later.

Actually, the simpler approach: keep the 2-step initialization. Deploy StakingVault proxy with `initialize(awpRegistry, address(0))`, then call `setStakeNFT(stakeNFT)` after StakeNFT is deployed. But we removed `setStakeNFT`...

Better: change `initialize` to only take `awpRegistry_`, and keep `setStakeNFT` as a one-time setter called from `initializeRegistry`.

- [ ] **Step 2: Update initializeRegistry**

Remove the `IStakingVault(stakingVault).setStakeNFT(stakeNFT)` call from AWPRegistry's `initializeRegistry` — StakingVault now has its own `setStakeNFT` still callable by AWPRegistry.

- [ ] **Step 3: Update InitCodeHashes.s.sol**

Add `SALT_STAKING_VAULT_IMPL` and update tier assignments.

- [ ] **Step 4: Build and test**

- [ ] **Step 5: Commit**

---

### Task 4: Update Tests

**Files:**
- Modify: `contracts/test/AWPRegistry.t.sol`
- Modify: `contracts/test/Integration.t.sol`
- Modify: `contracts/test/E2E.t.sol`
- Modify: `contracts/test/StakingVault.t.sol`

- [ ] **Step 1: Update test setUp — deploy StakingVault via proxy**

All test files that deploy StakingVault need `new StakingVault()` + ERC1967Proxy + initialize.

- [ ] **Step 2: Update allocation tests**

Tests that call `awpRegistry.allocate(...)` must change to `vault.allocate(...)`. Same for deallocate, reallocate.

- [ ] **Step 3: Add gasless allocateFor/deallocateFor tests in StakingVault.t.sol**

Test EIP-712 signature flow for allocateFor.

- [ ] **Step 4: Run full test suite**

```bash
forge test
```

- [ ] **Step 5: Commit**

---

### Task 5: Update API (Relay + Indexer)

**Files:**
- Modify: `api/internal/chain/relayer.go` — add StakingVault binding
- Modify: `api/internal/server/handler/relay.go` — allocate/deallocate relay
- Modify: `api/internal/chain/indexer.go` — listen for events from StakingVault
- Modify: `api/internal/chain/client.go` — add StakingVault binding
- Modify: `api/Makefile` — generate StakingVault bindings

- [ ] **Step 1: Regenerate Go bindings for StakingVault**

```bash
cd api && make generate-bindings
```

- [ ] **Step 2: Add StakingVault to relayer**

Relayer needs a `stakingVault` binding. `RelayAllocate` and `RelayDeallocate` call StakingVault instead of AWPRegistry.

- [ ] **Step 3: Update indexer to listen for allocation events from StakingVault**

Add StakingVault address to the FilterLogs address list. Parse `Allocated`/`Deallocated`/`Reallocated` from StakingVault binding instead of AWPRegistry binding.

- [ ] **Step 4: Build and test**

```bash
go build ./... && go test ./internal/server/handler/ -count=1 -timeout 60s
```

- [ ] **Step 5: Commit**

---

## Task Summary

| Task | Description | Estimate |
|------|-------------|----------|
| 1 | StakingVault UUPS + EIP-712 + public allocate | 15 min |
| 2 | Remove allocation from AWPRegistry | 5 min |
| 3 | Update Deploy.s.sol | 5 min |
| 4 | Update tests | 10 min |
| 5 | Update API (relay + indexer) | 10 min |
