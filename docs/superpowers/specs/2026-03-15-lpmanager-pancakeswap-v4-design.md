# LPManager PancakeSwap V4 Integration Design Spec

## Overview

Replace the mock LPManager implementation with real PancakeSwap V4 (Infinity) Concentrated Liquidity pool integration. LP positions are permanently locked (no withdraw functions exposed).

## Motivation

Current LPManager generates fake pool addresses via `keccak256` and holds tokens without creating real liquidity. Production deployment requires actual DEX integration for AWP/Alpha trading pairs.

## Design Decisions

| Decision | Choice | Reasoning |
|----------|--------|-----------|
| Pool type | Concentrated Liquidity (CL) | Standard for token pairs on PancakeSwap V4 |
| Fee tier | 1% (fee=10000, tickSpacing=200) | New token pairs with expected volatility |
| Price range | Full range (MIN_TICK to MAX_TICK) | Permanent lock = no rebalance, always has liquidity |
| LP lock | Never redeemable | LPManager exposes no decrease/burn functions |
| Permit2 flow | Per-call approve (SafeERC20.forceApprove) | Low frequency, safe for all token implementations |
| Testing | Fork BSC mainnet | Real PancakeSwap contracts, no mocks needed for LP tests |
| EVM version | Cancun | PancakeSwap V4 uses transient storage (TSTORE/TLOAD) |

## BSC Mainnet Contract Addresses

```
Vault:              0x238a358808379702088667322f80aC48bAd5e6c4
CLPoolManager:      0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b
CLPositionManager:  0x55f4c8abA71A1e923edC303eb4fEfF14608cC226
Permit2:            0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768
```

## Confirmed PancakeSwap V4 APIs (verified via BSC fork tests)

### CLPoolManager.initialize — NO hookData parameter
```solidity
// Selector: 0x8b0c1b22
function initialize(PoolKey calldata key, uint160 sqrtPriceX96) external returns (int24 tick);
```

### CLPositionManager.modifyLiquidities
```solidity
// Selector: 0xdd46508f
function modifyLiquidities(bytes calldata payload, uint256 deadline) external payable;
// payload = abi.encode(bytes actions, bytes[] params)
```

### CLPositionManager.nextTokenId — confirmed exists
```solidity
// Selector: 0x75794a3c
function nextTokenId() external view returns (uint256);
```

### PoolKey struct — Currency = plain address
```solidity
struct PoolKey {
    address currency0;     // lower address
    address currency1;     // higher address
    address hooks;         // address(0) for no hooks
    address poolManager;   // CLPoolManager address
    uint24 fee;            // 10000 = 1%
    bytes32 parameters;    // tickSpacing at bits [16:39]
}
```

### PoolId computation — raw memory hash, NOT abi.encode
```solidity
// PoolId = keccak256 of 6 raw 32-byte slots (0xc0 bytes)
bytes32 poolId;
assembly { poolId := keccak256(poolKey, 0xc0) }
```

### Action constants (confirmed from source)
```
CL_MINT_POSITION = 0x02
SETTLE_PAIR      = 0x0d
```

### CL_MINT_POSITION params encoding
```solidity
abi.encode(
    PoolKey poolKey,
    int24 tickLower,
    int24 tickUpper,
    uint256 liquidity,
    uint128 amount0Max,
    uint128 amount1Max,
    address owner,
    bytes hookData
)
```

### Parameters bytes32 encoding
```solidity
// tickSpacing at bits [16:39], OFFSET_TICK_SPACING = 16
bytes32 parameters = bytes32(uint256(uint24(tickSpacing)) << 16);
```

### Permit2.approve
```solidity
function approve(address token, address spender, uint160 amount, uint48 expiration) external;
```

## Contract Changes

### LPManager.sol — Rewrite

**Imports:**
```solidity
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
```

**Inline interfaces (minimal — only what we call):**
```solidity
interface ICLPoolManager {
    function initialize(PoolKey calldata key, uint160 sqrtPriceX96)
        external returns (int24 tick);
}

interface ICLPositionManager {
    function modifyLiquidities(bytes calldata payload, uint256 deadline) external payable;
    function nextTokenId() external view returns (uint256);
}

interface IPermit2 {
    function approve(address token, address spender, uint160 amount, uint48 expiration) external;
}
```

**PoolKey struct:**
```solidity
struct PoolKey {
    address currency0;
    address currency1;
    address hooks;
    address poolManager;
    uint24 fee;
    bytes32 parameters;
}
```

**Storage:**
```solidity
address public immutable rootNet;
address public immutable clPoolManager;
address public immutable clPositionManager;
address public immutable permit2;
IERC20 public immutable awpToken;

uint24 public constant POOL_FEE = 10000;       // 1%
int24 public constant TICK_SPACING = 200;
int24 public constant MIN_TICK = -887200;       // Aligned to tickSpacing=200
int24 public constant MAX_TICK = 887200;

mapping(address => bytes32) public alphaTokenToPoolId;  // PoolId (bytes32, not address)
mapping(address => uint256) public alphaTokenToTokenId;  // Position token ID
```

Note: `alphaTokenToPool` changes from `mapping(address => address)` to `mapping(address => bytes32)` because PancakeSwap V4 PoolId is a bytes32 hash, not a contract address.

**Constructor (5 params):**
```solidity
constructor(
    address rootNet_,
    address clPoolManager_,
    address clPositionManager_,
    address permit2_,
    address awpToken_
)
```

**ILPManager interface change:**
Return type of `createPoolAndAddLiquidity` changes:
```solidity
// Old:
function createPoolAndAddLiquidity(address alphaToken, uint256 awpAmount, uint256 alphaAmount)
    external returns (address pool, uint256 lpTokenId);

// New:
function createPoolAndAddLiquidity(address alphaToken, uint256 awpAmount, uint256 alphaAmount)
    external returns (bytes32 poolId, uint256 lpTokenId);
```

This requires updating `RootNet.registerSubnet` which stores the return value in `SubnetInfo.lpPool`. The `lpPool` field type in `IRootNet.SubnetInfo` changes from `address` to `bytes32`.

**`createPoolAndAddLiquidity` implementation:**

```solidity
function createPoolAndAddLiquidity(address alphaToken, uint256 awpAmount, uint256 alphaAmount)
    external onlyRootNet returns (bytes32 poolId, uint256 lpTokenId)
{
    if (alphaTokenToPoolId[alphaToken] != bytes32(0)) revert PoolAlreadyExists();

    // 1. Sort tokens
    address awp = address(awpToken);
    (address c0, address c1) = awp < alphaToken ? (awp, alphaToken) : (alphaToken, awp);
    (uint256 amt0, uint256 amt1) = awp < alphaToken ? (awpAmount, alphaAmount) : (alphaAmount, awpAmount);

    // 2. Construct PoolKey
    PoolKey memory poolKey = PoolKey({
        currency0: c0,
        currency1: c1,
        hooks: address(0),
        poolManager: clPoolManager,
        fee: POOL_FEE,
        parameters: bytes32(uint256(uint24(TICK_SPACING)) << 16)
    });

    // 3. Calculate sqrtPriceX96
    // price = amt1 / amt0; sqrtPriceX96 = sqrt(price) * 2^96
    // Safe computation: sqrt(amt1) * 2^96 / sqrt(amt0)
    uint160 sqrtPriceX96 = uint160(
        (Math.sqrt(amt1) << 96) / Math.sqrt(amt0)
    );

    // 4. Initialize pool (NO hookData param)
    ICLPoolManager(clPoolManager).initialize(poolKey, sqrtPriceX96);

    // 5. Approve tokens via Permit2 (forceApprove for safety)
    SafeERC20.forceApprove(IERC20(c0), permit2, amt0);
    IPermit2(permit2).approve(c0, clPositionManager, uint160(amt0), type(uint48).max);
    SafeERC20.forceApprove(IERC20(c1), permit2, amt1);
    IPermit2(permit2).approve(c1, clPositionManager, uint160(amt1), type(uint48).max);

    // 6. Calculate liquidity for full-range position
    // For full range: L = amt0 * sqrtPriceX96 / 2^96 (simplified)
    // Use: LiquidityAmounts.getLiquidityForAmounts equivalent
    uint256 liquidity = _getLiquidityForAmounts(sqrtPriceX96, amt0, amt1);

    // 7. Record token ID before mint
    lpTokenId = ICLPositionManager(clPositionManager).nextTokenId();

    // 8. Encode and execute mint
    bytes memory actions = abi.encodePacked(uint8(0x02), uint8(0x0d)); // CL_MINT_POSITION, SETTLE_PAIR
    bytes[] memory params = new bytes[](2);
    params[0] = abi.encode(
        poolKey,
        MIN_TICK,
        MAX_TICK,
        liquidity,
        uint128(amt0),     // amount0Max (slippage protection)
        uint128(amt1),     // amount1Max
        address(this),     // owner = LPManager (permanently locked)
        bytes("")          // hookData
    );
    params[1] = abi.encode(c0, c1);

    ICLPositionManager(clPositionManager).modifyLiquidities(
        abi.encode(actions, params),
        block.timestamp
    );

    // 9. Compute and store PoolId (raw memory hash per PancakeSwap V4)
    poolId = _computePoolId(poolKey);
    alphaTokenToPoolId[alphaToken] = poolId;
    alphaTokenToTokenId[alphaToken] = lpTokenId;
}
```

**Internal helpers:**

```solidity
/// @dev Compute PoolId — raw keccak256 of 6 ABI-packed 32-byte slots (0xc0 bytes)
function _computePoolId(PoolKey memory key) internal pure returns (bytes32 id) {
    assembly {
        id := keccak256(key, 0xc0)
    }
}

/// @dev Calculate liquidity for full-range position given amounts and sqrtPrice
/// For full range (MIN_TICK to MAX_TICK):
///   sqrtRatioA = TickMath.MIN_SQRT_RATIO ≈ 4295128739
///   sqrtRatioB = TickMath.MAX_SQRT_RATIO ≈ 1461446703485210103287273052203988822378723970342
///   L_from_amt0 = amt0 * sqrtPrice * sqrtRatioB / ((sqrtRatioB - sqrtPrice) * 2^96)
///   L_from_amt1 = amt1 * 2^96 / (sqrtPrice - sqrtRatioA)
///   L = min(L_from_amt0, L_from_amt1)
/// Simplified for full range: use amt0 * sqrtPrice / 2^96 as approximation
/// since sqrtRatioB >> sqrtPrice for typical prices
function _getLiquidityForAmounts(uint160 sqrtPriceX96, uint256 amt0, uint256 amt1)
    internal pure returns (uint256)
{
    // Use the standard Uniswap/PancakeSwap LiquidityAmounts formula
    uint160 sqrtRatioAX96 = 4295128739;  // MIN_SQRT_RATIO
    uint160 sqrtRatioBX96 = 1461446703485210103287273052203988822378723970342;  // MAX_SQRT_RATIO

    // 两步 mulDiv 避免 sqrtPriceX96 * sqrtRatioBX96 溢出 uint256
    uint256 intermediate = Math.mulDiv(amt0, uint256(sqrtPriceX96), 1 << 96);
    uint256 L0 = Math.mulDiv(intermediate, uint256(sqrtRatioBX96),
                             uint256(sqrtRatioBX96 - sqrtPriceX96));
    uint256 L1 = Math.mulDiv(amt1, 1 << 96, uint256(sqrtPriceX96) - uint256(sqrtRatioAX96));

    return L0 < L1 ? L0 : L1;
}
```

### IRootNet.sol — SubnetInfo + LPCreated event type change

```solidity
struct SubnetInfo {
    address subnetContract;
    address alphaToken;
    bytes32 lpPool;          // Changed from address to bytes32 (PancakeSwap V4 PoolId)
    SubnetStatus status;
    uint64 createdAt;
    uint64 activatedAt;
}

// LPCreated event: pool param changes from address to bytes32
event LPCreated(uint256 indexed subnetId, bytes32 poolId, uint256 awpAmount, uint256 alphaAmount);
```

### RootNet.registerSubnet — update lpPool assignment

```solidity
// Old:
(address alphaToken, address pool) = _deployAlphaAndLP(...);
subnets[subnetId] = SubnetInfo({
    ...
    lpPool: pool,
    ...
});

// New:
(address alphaToken, bytes32 poolId) = _deployAlphaAndLP(...);
subnets[subnetId] = SubnetInfo({
    ...
    lpPool: poolId,
    ...
});
```

`_deployAlphaAndLP` return type changes from `(address, address)` to `(address, bytes32)`.

### Deploy.s.sol and TestDeploy.s.sol

```solidity
// Production:
lp = new LPManager(
    address(rootNet),
    0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b,  // CLPoolManager
    0x55f4c8abA71A1e923edC303eb4fEfF14608cC226,  // CLPositionManager
    0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768,  // Permit2
    address(awp)
);
```

### foundry.toml — EVM version

```toml
evm_version = "cancun"
```

Required because PancakeSwap V4 uses EIP-1153 transient storage. Already set to 0.8.24.

## Test Strategy

### MockLPManager for non-fork tests

Create `test/helpers/MockLPManager.sol` — extract current mock behavior (keccak256 fake pool, no real DEX calls). Implements `ILPManager` interface with the new `bytes32` return type. All 269 existing tests use MockLPManager to avoid fork dependency.

### Fork tests for real integration

Create `test/LPManager.t.sol` — run with `forge test --match-contract LPManagerForkTest --fork-url $ETH_RPC_URL --evm-version cancun`

Tests:
- `test_createPoolAndAddLiquidity`: Deploy tokens, fund LPManager, create pool. Verify pool exists via CLPoolManager.getSlot0(poolId), verify position created (nextTokenId incremented), verify tokens transferred from LPManager.
- `test_createPool_revertsDoubleCreate`: Same Alpha twice → PoolAlreadyExists.
- `test_tokenOrdering`: AWP address < Alpha and AWP address > Alpha both work.
- `test_swapAfterPoolCreation`: Create pool, then swap via CLPoolManager to verify liquidity is real and tradeable.
- `test_lpPermanentlyLocked`: Verify LPManager has no decrease/burn functions.

### Existing tests update

E2E.t.sol, Integration.t.sol, RootNet.t.sol:
- Replace `new LPManager(rootNet, address(0), address(0), awp)` with `new MockLPManager(rootNet, awp)`
- Update any assertions on `lpPool` from `address` to `bytes32` type
- LPCreated event parameter changes if lpPool type changes

## Files Changed

| File | Action |
|------|--------|
| `contracts/src/core/LPManager.sol` | Rewrite (real PancakeSwap V4) |
| `contracts/src/interfaces/ILPManager.sol` | Update return type to bytes32 |
| `contracts/src/interfaces/IRootNet.sol` | SubnetInfo.lpPool: address → bytes32 |
| `contracts/src/RootNet.sol` | Update _deployAlphaAndLP return type |
| `contracts/test/helpers/MockLPManager.sol` | Create (extract current mock) |
| `contracts/test/LPManager.t.sol` | Rewrite (fork tests) |
| `contracts/test/E2E.t.sol` | Use MockLPManager |
| `contracts/test/Integration.t.sol` | Use MockLPManager |
| `contracts/test/RootNet.t.sol` | Use MockLPManager |
| `contracts/test/PancakeSwapV4Research.t.sol` | Delete (research complete) |
| `contracts/script/Deploy.s.sol` | Update LPManager constructor |
| `contracts/script/TestDeploy.s.sol` | Update LPManager constructor |
| `contracts/foundry.toml` | Ensure evm_version = "cancun" |
