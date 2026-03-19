# LPManager PancakeSwap V4 Integration Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace mock LPManager with real PancakeSwap V4 CL pool integration; LP permanently locked.

**Architecture:** LPManager calls CLPoolManager.initialize() to create a 1% fee CL pool, then CLPositionManager.modifyLiquidities() to mint a full-range position via Permit2. Non-fork tests use MockLPManager (current mock logic extracted). Fork tests verify real PancakeSwap integration on BSC mainnet.

**Tech Stack:** Solidity 0.8.24, PancakeSwap V4 (Infinity), OpenZeppelin 5.x, Foundry fork testing

**Spec:** `docs/superpowers/specs/2026-03-15-lpmanager-pancakeswap-v4-design.md`

---

## Chunk 1: Core Contract Changes

### Task 1: Update IAWPRegistry.sol — lpPool type + LPCreated event

**Files:**
- Modify: `contracts/src/interfaces/IAWPRegistry.sol`

- [ ] **Step 1: Change SubnetInfo.lpPool from address to bytes32**

In the `SubnetInfo` struct, change:
```solidity
address lpPool;           // PancakeSwap V4 LP 池地址
```
to:
```solidity
bytes32 lpPool;           // PancakeSwap V4 PoolId (bytes32 hash)
```

- [ ] **Step 2: Change LPCreated event pool param from address to bytes32**

Change:
```solidity
event LPCreated(uint256 indexed subnetId, address pool, uint256 awpAmount, uint256 alphaAmount);
```
to:
```solidity
event LPCreated(uint256 indexed subnetId, bytes32 poolId, uint256 awpAmount, uint256 alphaAmount);
```

### Task 2: Update ILPManager.sol — return type

**Files:**
- Modify: `contracts/src/interfaces/ILPManager.sol`

- [ ] **Step 1: Change return type from address to bytes32**

Change `createPoolAndAddLiquidity` return type:
```solidity
function createPoolAndAddLiquidity(address alphaToken, uint256 awpAmount, uint256 alphaAmount)
    external returns (bytes32 poolId, uint256 lpTokenId);
```

### Task 3: Update AWPRegistry.sol — adapt to new return types

**Files:**
- Modify: `contracts/src/AWPRegistry.sol`

- [ ] **Step 1: Update _deployAlphaAndLP return type**

Change function signature from:
```solidity
function _deployAlphaAndLP(...) internal returns (address alphaToken, address pool)
```
to:
```solidity
function _deployAlphaAndLP(...) internal returns (address alphaToken, bytes32 poolId)
```

Update the LP call:
```solidity
(poolId,) = ILPManager(lpManager).createPoolAndAddLiquidity(alphaToken, lpAWPAmount, INITIAL_ALPHA_MINT);
```

- [ ] **Step 2: Update registerSubnet to store bytes32 lpPool**

In the `SubnetInfo` struct literal:
```solidity
lpPool: poolId,
```

- [ ] **Step 3: Update _emitSubnetRegistered**

Change parameter from `address pool` to `bytes32 poolId` and pass to `LPCreated` event.

### Task 4: Create MockLPManager for non-fork tests

**Files:**
- Create: `contracts/test/helpers/MockLPManager.sol`

- [ ] **Step 1: Extract current mock logic**

Create MockLPManager that implements ILPManager with the new bytes32 return type:
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract MockLPManager {
    address public immutable rootNet;
    IERC20 public immutable awpToken;
    mapping(address => bytes32) public alphaTokenToPool;

    error NotAWPRegistry();
    error PoolAlreadyExists();

    modifier onlyAWPRegistry() {
        if (msg.sender != rootNet) revert NotAWPRegistry();
        _;
    }

    constructor(address rootNet_, address awpToken_) {
        rootNet = rootNet_;
        awpToken = IERC20(awpToken_);
    }

    function createPoolAndAddLiquidity(address alphaToken, uint256, uint256)
        external onlyAWPRegistry returns (bytes32 poolId, uint256 lpTokenId)
    {
        if (alphaTokenToPool[alphaToken] != bytes32(0)) revert PoolAlreadyExists();
        poolId = keccak256(abi.encodePacked(alphaToken, address(awpToken)));
        alphaTokenToPool[alphaToken] = poolId;
        lpTokenId = 0;
    }
}
```

### Task 5: Rewrite LPManager.sol — real PancakeSwap V4

**Files:**
- Rewrite: `contracts/src/core/LPManager.sol`

- [ ] **Step 1: Write the full implementation**

Complete rewrite per spec. Key elements:
- Inline interfaces: `ICLPoolManager` (initialize without hookData), `ICLPositionManager` (modifyLiquidities, nextTokenId), `IPermit2` (approve), `PoolKey` struct
- Constructor: 5 params (rootNet, clPoolManager, clPositionManager, permit2, awpToken)
- Constants: POOL_FEE=10000, TICK_SPACING=200, MIN_TICK=-887200, MAX_TICK=887200, MIN_SQRT_RATIO=4295128739, MAX_SQRT_RATIO=1461446703485210103287273052203988822378723970342
- Storage: `alphaTokenToPoolId` (bytes32), `alphaTokenToTokenId` (uint256)
- `createPoolAndAddLiquidity`: sort tokens, construct PoolKey, compute sqrtPriceX96 via `Math.sqrt(amt1) << 96 / Math.sqrt(amt0)`, initialize pool, forceApprove + Permit2.approve, compute liquidity via `_getLiquidityForAmounts`, encode CL_MINT_POSITION(0x02) + SETTLE_PAIR(0x0d) payload, call modifyLiquidities, compute PoolId via assembly keccak256
- `_getLiquidityForAmounts`: two-step Math.mulDiv for L0 (avoid overflow), standard L1 formula, return min
- `_computePoolId`: assembly keccak256(key, 0xc0)
- Import `Math` from OpenZeppelin for sqrt and mulDiv

### Task 6: Update Deploy Scripts

**Files:**
- Modify: `contracts/script/Deploy.s.sol`
- Modify: `contracts/script/TestDeploy.s.sol`

- [ ] **Step 1: Update Deploy.s.sol**

Change LPManager constructor from 4 params to 5:
```solidity
LPManager lp = new LPManager(
    address(rootNet),
    0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b,  // CLPoolManager
    0x55f4c8abA71A1e923edC303eb4fEfF14608cC226,  // CLPositionManager
    0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768,  // Permit2
    address(awp)
);
```

- [ ] **Step 2: Update TestDeploy.s.sol**

Same 5-param constructor with real BSC addresses.

### Task 7: Update foundry.toml

**Files:**
- Modify: `contracts/foundry.toml`

- [ ] **Step 1: Change evm_version to cancun**

Change `evm_version = "paris"` to `evm_version = "cancun"` (PancakeSwap V4 uses transient storage).

- [ ] **Step 2: Verify source compilation**

Run: `cd /home/ubuntu/code/Cortexia/contracts && /home/ubuntu/.foundry/bin/forge build --skip test 2>&1`
Expected: Source files compile. Tests will fail (expected — not updated yet).

- [ ] **Step 3: Commit**

```bash
git add contracts/src/ contracts/script/ contracts/test/helpers/MockLPManager.sol contracts/foundry.toml
git commit -m "feat: LPManager PancakeSwap V4 integration + MockLPManager for tests"
```

## Chunk 2: Test Updates

### Task 8: Update non-fork tests to use MockLPManager

**Files:**
- Modify: `contracts/test/AWPRegistry.t.sol`
- Modify: `contracts/test/E2E.t.sol`
- Modify: `contracts/test/Integration.t.sol`

- [ ] **Step 1: Update AWPRegistry.t.sol**

Replace `LPManager` import with `MockLPManager`. Change setUp:
```solidity
import {MockLPManager} from "./helpers/MockLPManager.sol";
// ...
lp = new MockLPManager(address(rootNet), address(awp));
```

Update any assertions that reference `lpPool` as `address` to `bytes32`.

Note: AWPRegistry.initializeRegistry still accepts `address(lp)` since `lpManager` is stored as `address`.

- [ ] **Step 2: Update E2E.t.sol**

Same pattern: import MockLPManager, change constructor in `_deploy()`.
Update `LPCreated` event assertions if any exist (pool param is now bytes32).

- [ ] **Step 3: Update Integration.t.sol**

Same pattern in `_fullDeployment()`.

- [ ] **Step 4: Run non-fork tests**

Run: `cd /home/ubuntu/code/Cortexia/contracts && /home/ubuntu/.foundry/bin/forge test --no-match-contract "Fork|Research" 2>&1`
Expected: All 269 tests pass.

### Task 9: Rewrite LPManager.t.sol as fork test

**Files:**
- Rewrite: `contracts/test/LPManager.t.sol`

- [ ] **Step 1: Write fork test file**

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AlphaToken} from "../src/token/AlphaToken.sol";
import {LPManager} from "../src/core/LPManager.sol";

contract LPManagerForkTest is Test {
    // BSC mainnet PancakeSwap V4 addresses
    address constant CL_POOL_MANAGER = 0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b;
    address constant CL_POSITION_MANAGER = 0x55f4c8abA71A1e923edC303eb4fEfF14608cC226;
    address constant PERMIT2 = 0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768;

    AWPToken awp;
    AlphaToken alphaImpl;
    LPManager lp;
    address rootNet;
    address deployer;

    function setUp() public {
        deployer = makeAddr("deployer");
        rootNet = makeAddr("rootNet");

        vm.startPrank(deployer);
        awp = new AWPToken("AWP", "AWP", deployer);

        // Deploy Alpha token (use implementation directly, not clone)
        alphaImpl = new AlphaToken();
        alphaImpl.initialize("Alpha", "ALPHA", rootNet);

        // Mint Alpha to rootNet for LP
        vm.startPrank(rootNet);
        alphaImpl.mint(rootNet, 100_000_000 * 1e18);

        // Deploy real LPManager
        lp = new LPManager(rootNet, CL_POOL_MANAGER, CL_POSITION_MANAGER, PERMIT2, address(awp));

        // Transfer tokens to LPManager (simulating AWPRegistry.registerSubnet flow)
        uint256 awpAmount = 1_000_000 * 1e18;
        uint256 alphaAmount = 100_000_000 * 1e18;

        vm.startPrank(deployer);
        awp.transfer(address(lp), awpAmount);
        vm.startPrank(rootNet);
        alphaImpl.transfer(address(lp), alphaAmount);
        vm.stopPrank();
    }
}
```

Tests to implement:
- `test_createPoolAndAddLiquidity`: Call from rootNet, verify poolId != 0, tokenId > 0, tokens transferred out of LPManager
- `test_createPool_revertsDoubleCreate`: Create same Alpha pool twice → PoolAlreadyExists
- `test_createPool_revertsNonAWPRegistry`: User calls → NotAWPRegistry
- `test_tokenOrdering`: Verify works regardless of AWP vs Alpha address ordering
- `test_swapAfterPoolCreation`: After creating pool, verify swap works via CLPoolManager

Each test calls `lp.createPoolAndAddLiquidity(address(alphaImpl), awpAmount, alphaAmount)` from rootNet prank.

- [ ] **Step 2: Run fork tests**

Run: `cd /home/ubuntu/code/Cortexia/contracts && /home/ubuntu/.foundry/bin/forge test --match-contract LPManagerForkTest --fork-url https://delicate-delicate-fire.bsc.quiknode.pro/afb2d8691739eff71081fc07e5816747eac49db4 -vvv 2>&1`
Expected: All tests pass.

- [ ] **Step 3: Delete research test**

Remove: `contracts/test/PancakeSwapV4Research.t.sol`

- [ ] **Step 4: Commit**

```bash
git add contracts/test/
git commit -m "test: MockLPManager for unit tests, fork tests for PancakeSwap V4 integration"
```

## Chunk 3: Go Backend + Cleanup

### Task 10: Regenerate Go bindings

**Files:**
- Regenerate: `api/internal/chain/bindings/root_net.go`

- [ ] **Step 1: Regenerate AWPRegistry binding**

The SubnetInfo struct changed (lpPool: address → bytes32), LPCreated event changed. Regenerate:
```bash
cd /home/ubuntu/code/Cortexia/contracts
/home/ubuntu/.foundry/bin/forge inspect AWPRegistry abi --json > /tmp/rootnet_abi.json
/home/ubuntu/go/bin/abigen --abi /tmp/rootnet_abi.json --pkg bindings --type AWPRegistry --out /home/ubuntu/code/Cortexia/api/internal/chain/bindings/root_net.go
```

- [ ] **Step 2: Update indexer if LPCreated parsing needs changes**

Read indexer.go, find `ParseLPCreated` usage. The `Pool` field changes from `common.Address` to `[32]byte`. Update accordingly:
```go
// Old: evt.Pool.Hex()
// New: common.Bytes2Hex(evt.PoolId[:]) or hex.EncodeToString(evt.PoolId[:])
```

Also update the DB query parameter — `UpdateSubnetLP` likely expects a string for `lp_pool`. Convert bytes32 to hex string.

- [ ] **Step 3: Verify Go build**

Run: `cd /home/ubuntu/code/Cortexia/api && go build ./...`
Expected: Clean build.

- [ ] **Step 4: Commit**

```bash
git add api/
git commit -m "feat: regenerate bindings for LPManager PancakeSwap V4, update indexer"
```

### Task 11: Final Verification

- [ ] **Step 1: Non-fork Solidity tests**

Run: `cd /home/ubuntu/code/Cortexia/contracts && /home/ubuntu/.foundry/bin/forge test --no-match-contract "Fork|Research" 2>&1`
Expected: All tests pass.

- [ ] **Step 2: Fork tests**

Run: `cd /home/ubuntu/code/Cortexia/contracts && /home/ubuntu/.foundry/bin/forge test --match-contract LPManagerForkTest --fork-url https://delicate-delicate-fire.bsc.quiknode.pro/afb2d8691739eff71081fc07e5816747eac49db4 -vvv 2>&1`
Expected: All fork tests pass.

- [ ] **Step 3: Go build**

Run: `cd /home/ubuntu/code/Cortexia/api && go build ./...`
Expected: Clean build.
