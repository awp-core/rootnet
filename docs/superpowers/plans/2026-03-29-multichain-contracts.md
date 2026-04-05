# Multi-Chain Contracts Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make AWP contracts multi-chain ready: global subnetId encoding, cross-chain allocate support, configurable initial mint, and multi-chain deploy config.

**Architecture:** AWPRegistry encodes subnetId as `(block.chainid << 64) | localCounter`. Allocate/reallocate remove on-chain subnet status check (cross-chain subnets can't be verified locally). AWPToken constructor mint becomes configurable. Deploy uses `chains.yaml` for per-chain config.

**Tech Stack:** Solidity 0.8.24, Foundry, OpenZeppelin 5.x

---

### Task 1: AWPRegistry — Global SubnetId Encoding

**Files:**
- Modify: `contracts/src/AWPRegistry.sol:88,233,629,954`
- Modify: `contracts/test/AWPRegistry.t.sol`

- [ ] **Step 1: Write failing test for global subnetId**

Add to `contracts/test/AWPRegistry.t.sol` after the UUPS tests:

```solidity
// ── SubnetId encoding tests ──

function test_subnetIdEncodesChainId() public {
    uint256 subnetId = _registerSubnet();
    // subnetId should be (block.chainid << 64) | 1
    uint256 expectedChainId = block.chainid;
    assertEq(subnetId >> 64, expectedChainId);
    assertEq(subnetId & ((1 << 64) - 1), 1);
}

function test_subnetIdIncrementsLocalCounter() public {
    uint256 id1 = _registerSubnet();
    // Register a second subnet
    vm.startPrank(user1);
    awp.approve(address(awpRegistry), 1_000_000 * 1e18);
    uint256 id2 = awpRegistry.registerSubnet(
        IAWPRegistry.SubnetParams("Sub2", "S2", address(0), bytes32(0), 0, "")
    );
    vm.stopPrank();

    assertEq(id1 >> 64, id2 >> 64); // same chainId
    assertEq((id1 & ((1 << 64) - 1)) + 1, id2 & ((1 << 64) - 1)); // localId increments
}

function test_extractChainId() public view {
    uint256 subnetId = (uint256(8453) << 64) | 42;
    assertEq(awpRegistry.extractChainId(subnetId), 8453);
    assertEq(awpRegistry.extractLocalId(subnetId), 42);
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd contracts && forge test --match-test "test_subnetIdEncodesChainId|test_subnetIdIncrementsLocalCounter|test_extractChainId" -v`
Expected: FAIL (functions don't exist / subnetId doesn't encode chainId)

- [ ] **Step 3: Implement global subnetId in AWPRegistry.sol**

In `contracts/src/AWPRegistry.sol`, make these changes:

1. Rename storage variable (line 88):
```solidity
// OLD:
uint256 private _nextSubnetId;
// NEW:
uint256 private _nextLocalId;
```

2. Update `initialize()` (line 233):
```solidity
// OLD:
_nextSubnetId = 1;
// NEW:
_nextLocalId = 1;
```

3. Update `_registerSubnet()` (line 629):
```solidity
// OLD:
uint256 subnetId = _nextSubnetId++;
// NEW:
uint256 subnetId = (block.chainid << 64) | _nextLocalId++;
```

4. Update `nextSubnetId()` view function (around line 954):
```solidity
// OLD:
function nextSubnetId() external view returns (uint256) {
    return _nextSubnetId;
}
// NEW:
function nextSubnetId() external view returns (uint256) {
    return (block.chainid << 64) | _nextLocalId;
}

/// @notice Extract chainId from a global subnetId
function extractChainId(uint256 subnetId) external pure returns (uint256) {
    return subnetId >> 64;
}

/// @notice Extract local counter from a global subnetId
function extractLocalId(uint256 subnetId) external pure returns (uint256) {
    return subnetId & ((1 << 64) - 1);
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `forge test --match-test "test_subnetId" -v`
Expected: All 3 new tests PASS

- [ ] **Step 5: Run full test suite**

Run: `forge test`
Expected: 235/237 pass (2 BSC fork expected failures). Existing tests still pass because `block.chainid` in Foundry test defaults to `31337`, so subnetIds are `(31337 << 64) | 1`, `(31337 << 64) | 2`, etc. Tests that compare subnetId directly (like `assertEq(subnetId, 1)`) will break — fix them in Step 6.

- [ ] **Step 6: Fix existing tests that hardcode subnetId=1**

Search all test files for `assertEq(subnetId, 1)` or `assertEq(id, 1)` or `subnetId == 1` and replace with the encoded form. For example in `Integration.t.sol`:

```solidity
// OLD:
assertEq(subnetId, 1);
// NEW:
assertEq(subnetId & ((1 << 64) - 1), 1); // local ID is 1
```

Also fix any test that calls endpoints with literal `1` as subnetId — replace with the variable returned from `_registerSubnet()`.

- [ ] **Step 7: Commit**

```bash
git add contracts/src/AWPRegistry.sol contracts/test/
git commit -m "feat: global subnetId encoding — (chainId << 64) | localCounter

SubnetId is now globally unique across all chains. Upper bits encode
block.chainid, lower 64 bits are a per-chain local counter.
Added extractChainId() and extractLocalId() helper view functions."
```

---

### Task 2: AWPRegistry — Remove Allocate Subnet Status Check

**Files:**
- Modify: `contracts/src/AWPRegistry.sol:487,517,554`
- Add test: `contracts/test/AWPRegistry.t.sol`

- [ ] **Step 1: Write test for cross-chain allocate (allocate to non-local subnetId)**

Add to `contracts/test/AWPRegistry.t.sol`:

```solidity
function test_allocateToCrossChainSubnet() public {
    // Stake AWP first
    vm.startPrank(user1);
    awp.approve(address(stakeNFT), 10_000 * 1e18);
    stakeNFT.deposit(10_000 * 1e18, 52 weeks);

    // Allocate to a "foreign" subnetId (Arbitrum chain, local ID 5)
    uint256 foreignSubnetId = (uint256(42161) << 64) | 5;
    awpRegistry.allocate(user1, user1, foreignSubnetId, 5_000 * 1e18);
    vm.stopPrank();

    // Verify allocation was recorded
    assertEq(vault.getAgentStake(user1, user1, foreignSubnetId), 5_000 * 1e18);
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `forge test --match-test test_allocateToCrossChainSubnet -v`
Expected: FAIL with `InvalidSubnetStatus()` (foreign subnet doesn't exist locally)

- [ ] **Step 3: Remove subnet status checks from allocate/allocateFor/reallocate**

In `contracts/src/AWPRegistry.sol`:

Line 487 — `allocate()`: DELETE this line:
```solidity
if (subnets[subnetId].status != SubnetStatus.Active) revert InvalidSubnetStatus();
```

Line 517 — `allocateFor()`: DELETE this line:
```solidity
if (subnets[subnetId].status != SubnetStatus.Active) revert InvalidSubnetStatus();
```

Line 554 — `reallocate()`: DELETE this line:
```solidity
if (subnets[toSubnetId].status != SubnetStatus.Active) revert InvalidSubnetStatus();
```

- [ ] **Step 4: Run test to verify it passes**

Run: `forge test --match-test test_allocateToCrossChainSubnet -v`
Expected: PASS

- [ ] **Step 5: Run full test suite**

Run: `forge test`
Expected: 238/240 pass (3 new tests + existing). Some existing tests may have tested the `InvalidSubnetStatus` revert on allocate — those tests need updating or removal since that check no longer exists.

- [ ] **Step 6: Commit**

```bash
git add contracts/src/AWPRegistry.sol contracts/test/AWPRegistry.t.sol
git commit -m "feat: remove on-chain subnet status check from allocate/reallocate

Cross-chain allocate requires allocating to subnets that don't exist
on the staker's chain. Subnet validity is now verified off-chain by
the indexer/oracle. StakingVault still enforces balance limits."
```

---

### Task 3: AWPToken — Configurable Initial Mint

**Files:**
- Modify: `contracts/src/token/AWPToken.sol:21,43-51`
- Modify: `contracts/script/Deploy.s.sol`
- Add test: `contracts/test/AWPToken.t.sol` (or modify existing)

- [ ] **Step 1: Modify AWPToken constructor to accept initialMint parameter**

In `contracts/src/token/AWPToken.sol`:

```solidity
// OLD (line 21):
uint256 public constant INITIAL_MINT = 200_000_000 * 1e18;

// NEW:
uint256 public immutable INITIAL_MINT;
```

```solidity
// OLD constructor (line 43-51):
constructor(string memory name_, string memory symbol_, address deployer_)
    ERC20(name_, symbol_)
    ERC20Permit(name_)
{
    admin = deployer_;
    _mint(deployer_, INITIAL_MINT);
}

// NEW:
constructor(string memory name_, string memory symbol_, address deployer_, uint256 initialMint_)
    ERC20(name_, symbol_)
    ERC20Permit(name_)
{
    admin = deployer_;
    INITIAL_MINT = initialMint_;
    if (initialMint_ > 0) {
        _mint(deployer_, initialMint_);
    }
}
```

- [ ] **Step 2: Update Deploy.s.sol to pass initialMint**

Read `INITIAL_MINT` from env var (default 200M for backward compat):

```solidity
uint256 initialMint = vm.envOr("INITIAL_MINT", uint256(200_000_000)) * 1e18;
```

Update the AWPToken CREATE2 deployment to pass the 4th parameter:

```solidity
AWPToken awp = AWPToken(_create2(
    saltAWPToken,
    abi.encodePacked(type(AWPToken).creationCode, abi.encode("AWP Token", "AWP", deployer, initialMint))
));
```

- [ ] **Step 3: Update all test setUp functions that deploy AWPToken**

Every test file that does `new AWPToken(...)` needs the 4th param. Search and update:

```solidity
// OLD:
awp = new AWPToken("AWP Token", "AWP", deployer);
// NEW:
awp = new AWPToken("AWP Token", "AWP", deployer, 200_000_000 * 1e18);
```

- [ ] **Step 4: Update InitCodeHashes.s.sol**

The AWPToken hash now includes the `initialMint` param:

```solidity
uint256 initialMint = vm.envOr("INITIAL_MINT", uint256(200_000_000)) * 1e18;
_logHash("AWPToken", abi.encodePacked(type(AWPToken).creationCode, abi.encode("AWP Token", "AWP", deployer, initialMint)));
```

- [ ] **Step 5: Run full test suite**

Run: `forge test`
Expected: All pass (same counts as before)

- [ ] **Step 6: Commit**

```bash
git add contracts/src/token/AWPToken.sol contracts/script/ contracts/test/
git commit -m "feat: AWPToken configurable initial mint

Constructor now accepts initialMint parameter instead of hardcoded
200M. Enables per-chain deployment with different initial distributions.
INITIAL_MINT changed from constant to immutable."
```

---

### Task 4: Multi-Chain Deploy Configuration (chains.yaml)

**Files:**
- Create: `chains.yaml`
- Create: `scripts/deploy-multichain.sh`
- Modify: `contracts/.env.example`

- [ ] **Step 1: Create chains.yaml template**

Create `chains.yaml`:

```yaml
# AWP Multi-Chain Deployment Configuration
# Each chain gets a full independent protocol deployment with identical contract addresses.
# Deploy: ./scripts/deploy-multichain.sh [chainName]

chains:
  base:
    chainId: 8453
    name: Base
    rpcUrl: ${BASE_RPC_URL}
    dex: uniswap_v4
    initialMint: 50000000  # 50M AWP
    explorer: https://basescan.org
    blockTime: 2
    poolManager: "0x498581ff718922c3f8e6a244956af099b2652b2b"
    positionManager: "0x7c5f5a4bbd8fd63184577525326123b519429bdc"
    permit2: "0x000000000022D473030F116dDEE9F6B43aC78BA3"
    swapRouter: "0x6ff5693b99212da76ad316178a184ab56d299b43"
    stateView: "0xa3c0c9b65bad0b08107aa264b0f3db444b867a71"

  arbitrum:
    chainId: 42161
    name: Arbitrum One
    rpcUrl: ${ARB_RPC_URL}
    dex: uniswap_v4
    initialMint: 50000000
    explorer: https://arbiscan.io
    blockTime: 0.25
    poolManager: "0x..."
    positionManager: "0x..."
    permit2: "0x000000000022D473030F116dDEE9F6B43aC78BA3"
    swapRouter: "0x..."
    stateView: "0x..."

  ethereum:
    chainId: 1
    name: Ethereum
    rpcUrl: ${ETH_RPC_URL}
    dex: uniswap_v4
    initialMint: 100000000  # 100M AWP
    explorer: https://etherscan.io
    blockTime: 12
    poolManager: "0x..."
    positionManager: "0x..."
    permit2: "0x000000000022D473030F116dDEE9F6B43aC78BA3"
    swapRouter: "0x..."
    stateView: "0x..."

  bsc:
    chainId: 56
    name: BNB Smart Chain
    rpcUrl: ${BSC_RPC_URL}
    dex: pancakeswap_v4
    initialMint: 0  # no initial mint on BSC (emission only)
    explorer: https://bscscan.com
    blockTime: 3
    poolManager: "0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b"
    positionManager: "0x55f4c8abA71A1e923edC303eb4fEfF14608cC226"
    permit2: "0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768"
    swapRouter: "0x1b81D678ffb9C0263b24A97847620C99d213eB14"
    stateView: ""
```

- [ ] **Step 2: Create deploy-multichain.sh**

Create `scripts/deploy-multichain.sh`:

```bash
#!/bin/bash
set -euo pipefail

# Usage: ./scripts/deploy-multichain.sh <chainName>
# Example: ./scripts/deploy-multichain.sh base
#          ./scripts/deploy-multichain.sh --all

CHAIN_NAME="${1:-}"
CHAINS_FILE="chains.yaml"

if [[ -z "$CHAIN_NAME" ]]; then
    echo "Usage: $0 <chainName|--all>"
    echo "Available chains:"
    python3 -c "import yaml; [print(f'  {k}') for k in yaml.safe_load(open('$CHAINS_FILE'))['chains']]"
    exit 1
fi

deploy_chain() {
    local name="$1"
    echo "═══════════════════════════════════════"
    echo "  Deploying to: $name"
    echo "═══════════════════════════════════════"

    # Read chain config from YAML
    local cfg
    cfg=$(python3 -c "
import yaml, json, os
chains = yaml.safe_load(open('$CHAINS_FILE'))['chains']
c = chains['$name']
# Resolve env vars in rpcUrl
rpc = os.path.expandvars(c['rpcUrl'])
c['rpcUrl'] = rpc
print(json.dumps(c))
")

    # Export as env vars for deploy.sh
    export ETH_RPC_URL=$(echo "$cfg" | jq -r '.rpcUrl')
    export POOL_MANAGER=$(echo "$cfg" | jq -r '.poolManager')
    export POSITION_MANAGER=$(echo "$cfg" | jq -r '.positionManager')
    export PERMIT2=$(echo "$cfg" | jq -r '.permit2')
    export CL_SWAP_ROUTER=$(echo "$cfg" | jq -r '.swapRouter')
    export STATE_VIEW=$(echo "$cfg" | jq -r '.stateView // empty')
    export INITIAL_MINT=$(echo "$cfg" | jq -r '.initialMint')

    # Run the standard deploy script (reuses existing salt.json for identical addresses)
    ./scripts/deploy.sh --skip-mine

    echo "✓ $name deployment complete"
}

if [[ "$CHAIN_NAME" == "--all" ]]; then
    for name in $(python3 -c "import yaml; [print(k) for k in yaml.safe_load(open('$CHAINS_FILE'))['chains']]"); do
        deploy_chain "$name"
    done
else
    deploy_chain "$CHAIN_NAME"
fi
```

- [ ] **Step 3: Make executable and commit**

```bash
chmod +x scripts/deploy-multichain.sh
git add chains.yaml scripts/deploy-multichain.sh
git commit -m "feat: multi-chain deployment config and script

chains.yaml defines per-chain config (RPC, DEX addresses, initial mint).
deploy-multichain.sh wraps deploy.sh to deploy to any configured chain
with --skip-mine (reuses shared salt.json for identical addresses)."
```

---

### Task 5: Update CLAUDE.md and Documentation

**Files:**
- Modify: `CLAUDE.md`
- Modify: `README.md`

- [ ] **Step 1: Update CLAUDE.md**

Add to the Core Architecture section:
- SubnetId encoding: `(block.chainid << 64) | localCounter`
- Cross-chain allocate: no on-chain subnet status check
- Multi-chain deployment: `chains.yaml` config, identical addresses via CREATE2
- AWPToken: configurable `initialMint` per chain

Update the Deployment Sequence to mention `chains.yaml` and per-chain config.

- [ ] **Step 2: Update README.md**

Add a "Multi-Chain" section explaining:
- Supported chains (from chains.yaml)
- How to deploy to a new chain
- SubnetId structure
- Cross-chain allocate semantics

- [ ] **Step 3: Commit**

```bash
git add CLAUDE.md README.md
git commit -m "docs: update for multi-chain architecture"
```

---

## Task Summary

| Task | Description | Files | Estimate |
|------|-------------|-------|----------|
| 1 | Global subnetId encoding | AWPRegistry.sol + tests | 5 min |
| 2 | Remove allocate subnet check | AWPRegistry.sol + tests | 3 min |
| 3 | Configurable AWPToken mint | AWPToken.sol + Deploy.s.sol + tests | 5 min |
| 4 | Multi-chain deploy config | chains.yaml + deploy script | 5 min |
| 5 | Documentation updates | CLAUDE.md + README.md | 3 min |
