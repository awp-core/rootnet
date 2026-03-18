// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {FullMath} from "infinity-core/src/pool-cl/libraries/FullMath.sol";
import {FixedPoint96} from "infinity-core/src/pool-cl/libraries/FixedPoint96.sol";

/// @dev PancakeSwap V4 CL PoolKey struct (verified on-chain via fork test)
struct PoolKey {
    address currency0;
    address currency1;
    address hooks;
    address poolManager;
    uint24 fee;
    bytes32 parameters;
}

/// @dev PancakeSwap V4 CLPoolManager interface (selector: 0x8b0c1b22, no hookData parameter)
interface ICLPoolManager {
    function initialize(PoolKey calldata key, uint160 sqrtPriceX96) external returns (int24 tick);
}

/// @dev PancakeSwap V4 CLPositionManager interface
interface ICLPositionManager {
    function modifyLiquidities(bytes calldata payload, uint256 deadline) external payable;
    function nextTokenId() external view returns (uint256);
}

/// @dev Permit2 interface for token approvals
interface IPermit2 {
    function approve(address token, address spender, uint160 amount, uint48 expiration) external;
}

/// @title LPManager — PancakeSwap V4 CL liquidity management
/// @notice Only RootNet may call; LP is permanently locked in this contract and cannot be withdrawn
/// @dev Integrates PancakeSwap V4 Concentrated Liquidity:
///   - Full-range liquidity (MIN_TICK ~ MAX_TICK)
///   - Token approvals via Permit2
///   - Pool ID computed as keccak256(PoolKey)
contract LPManager {
    using SafeERC20 for IERC20;

    // ══════════════════════════════════════════════
    //  Immutable storage
    // ══════════════════════════════════════════════

    /// @notice RootNet contract address
    address public immutable rootNet;
    /// @notice PancakeSwap V4 CLPoolManager address
    address public immutable clPoolManager;
    /// @notice PancakeSwap V4 CLPositionManager address
    address public immutable clPositionManager;
    /// @notice Permit2 contract address
    address public immutable permit2;
    /// @notice AWP token contract address
    IERC20 public immutable awpToken;

    // ══════════════════════════════════════════════
    //  Constants
    // ══════════════════════════════════════════════

    /// @dev Pool fee: 1% (10000 = 1%)
    uint24 public constant POOL_FEE = 10000;
    /// @dev Tick spacing
    int24 public constant TICK_SPACING = 200;
    /// @dev Full-range minimum tick (rounded to TICK_SPACING)
    int24 public constant MIN_TICK = -887200;
    /// @dev Full-range maximum tick (rounded to TICK_SPACING)
    int24 public constant MAX_TICK = 887200;
    /// @dev Precomputed TickMath.getSqrtRatioAtTick(MIN_TICK)
    uint160 public constant MIN_SQRT_RATIO = 4295128739;
    /// @dev Precomputed TickMath.getSqrtRatioAtTick(MAX_TICK)
    uint160 public constant MAX_SQRT_RATIO = 1461446703485210103287273052203988822378723970342;

    // ══════════════════════════════════════════════
    //  Mappings
    // ══════════════════════════════════════════════

    /// @notice Alpha token address → pool ID
    mapping(address => bytes32) public alphaTokenToPoolId;
    /// @notice Alpha token address → LP NFT tokenId
    mapping(address => uint256) public alphaTokenToTokenId;

    // ══════════════════════════════════════════════
    //  Errors
    // ══════════════════════════════════════════════

    error NotRootNet();
    /// @dev A LP pool already exists for this Alpha token
    error PoolAlreadyExists();

    /// @dev Only the RootNet contract may call
    modifier onlyRootNet() {
        if (msg.sender != rootNet) revert NotRootNet();
        _;
    }

    /// @notice Constructor
    /// @param rootNet_ RootNet contract address
    /// @param clPoolManager_ PancakeSwap V4 CLPoolManager address
    /// @param clPositionManager_ PancakeSwap V4 CLPositionManager address
    /// @param permit2_ Permit2 contract address
    /// @param awpToken_ AWP token contract address
    constructor(address rootNet_, address clPoolManager_, address clPositionManager_, address permit2_, address awpToken_) {
        rootNet = rootNet_;
        clPoolManager = clPoolManager_;
        clPositionManager = clPositionManager_;
        permit2 = permit2_;
        awpToken = IERC20(awpToken_);
    }

    /// @notice Create an LP pool and add full-range liquidity (called once during subnet registration)
    /// @dev Full flow:
    ///   1. Sort token addresses (PancakeSwap V4 requires currency0 < currency1)
    ///   2. Construct PoolKey and initialize the pool
    ///   3. Approve tokens to CLPositionManager via Permit2
    ///   4. Compute liquidity and mint the LP position
    ///   5. Record pool ID and LP NFT tokenId
    /// @param alphaToken Alpha token address
    /// @param awpAmount AWP amount
    /// @param alphaAmount Alpha amount
    /// @return poolId Pool ID (bytes32)
    /// @return lpTokenId LP NFT ID
    function createPoolAndAddLiquidity(address alphaToken, uint256 awpAmount, uint256 alphaAmount)
        external
        onlyRootNet
        returns (bytes32 poolId, uint256 lpTokenId)
    {
        // Each Alpha token may only have one LP pool
        if (alphaTokenToPoolId[alphaToken] != bytes32(0)) revert PoolAlreadyExists();

        // Sort tokens: PancakeSwap V4 requires currency0 < currency1
        address awp = address(awpToken);
        (address c0, address c1) = awp < alphaToken ? (awp, alphaToken) : (alphaToken, awp);
        (uint256 amt0, uint256 amt1) = awp < alphaToken ? (awpAmount, alphaAmount) : (alphaAmount, awpAmount);

        // Construct PoolKey
        PoolKey memory poolKey = PoolKey({
            currency0: c0,
            currency1: c1,
            hooks: address(0),
            poolManager: clPoolManager,
            fee: POOL_FEE,
            parameters: bytes32(uint256(int256(TICK_SPACING)) << 16)
        });

        // Compute initial price: sqrtPriceX96 = sqrt(amt1 * 2^192 / amt0)
        // Using FullMath.mulDiv for overflow-safe (amt1 * 2^192 / amt0), then sqrt
        uint256 ratioX192 = FullMath.mulDiv(amt1, FixedPoint96.Q96 * FixedPoint96.Q96, amt0);
        uint160 sqrtPriceX96 = uint160(Math.sqrt(ratioX192));

        // Initialize pool (no hookData parameter)
        ICLPoolManager(clPoolManager).initialize(poolKey, sqrtPriceX96);

        // Approve + mint LP position (split into an internal function to avoid stack too deep)
        lpTokenId = _approveAndMint(poolKey, c0, c1, amt0, amt1, sqrtPriceX96);

        // Compute pool ID: keccak256(PoolKey); PoolKey occupies 6 slots = 0xc0 bytes
        poolId = _computePoolId(poolKey);

        // Store mappings
        alphaTokenToPoolId[alphaToken] = poolId;
        alphaTokenToTokenId[alphaToken] = lpTokenId;
    }

    /// @dev Internal function: Permit2 approval + LP position minting (split from createPoolAndAddLiquidity to avoid stack too deep)
    function _approveAndMint(
        PoolKey memory poolKey, address c0, address c1,
        uint256 amt0, uint256 amt1, uint160 sqrtPriceX96
    ) internal returns (uint256 lpTokenId) {
        // Approve tokens to CLPositionManager via Permit2
        IERC20(c0).forceApprove(permit2, amt0);
        IPermit2(permit2).approve(c0, clPositionManager, uint160(amt0), uint48(block.timestamp + 600));
        IERC20(c1).forceApprove(permit2, amt1);
        IPermit2(permit2).approve(c1, clPositionManager, uint160(amt1), uint48(block.timestamp + 600));

        // Compute full-range liquidity using PancakeSwap's official library
        uint128 liquidity = LiquidityAmounts.getLiquidityForAmounts(
            sqrtPriceX96,
            MIN_SQRT_RATIO,
            MAX_SQRT_RATIO,
            amt0,
            amt1
        );

        // Record nextTokenId before minting
        lpTokenId = ICLPositionManager(clPositionManager).nextTokenId();

        // Encode mint operation: CL_MINT_POSITION(0x02) + SETTLE_PAIR(0x0d)
        bytes memory actions = abi.encodePacked(uint8(0x02), uint8(0x0d));
        bytes[] memory params = new bytes[](2);
        params[0] = abi.encode(poolKey, MIN_TICK, MAX_TICK, liquidity, uint128(amt0), uint128(amt1), address(this), bytes(""));
        params[1] = abi.encode(c0, c1);
        ICLPositionManager(clPositionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);
    }

    /// @dev Compute the pool ID (keccak256 hash) of a PoolKey
    function _computePoolId(PoolKey memory key) internal pure returns (bytes32) {
        return keccak256(abi.encode(key));
    }
}
