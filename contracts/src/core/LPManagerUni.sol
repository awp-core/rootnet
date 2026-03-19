// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {FullMath} from "infinity-core/src/pool-cl/libraries/FullMath.sol";
import {FixedPoint96} from "infinity-core/src/pool-cl/libraries/FixedPoint96.sol";

/// @dev Uniswap V4 PoolKey struct (5 fields, different from PancakeSwap V4's 6 fields)
struct UniPoolKey {
    address currency0;
    address currency1;
    uint24 fee;
    int24 tickSpacing;
    address hooks;
}

/// @dev Uniswap V4 PoolManager interface
interface IUniPoolManager {
    function initialize(UniPoolKey calldata key, uint160 sqrtPriceX96) external returns (int24 tick);
}

/// @dev Uniswap V4 PositionManager interface
interface IUniPositionManager {
    function modifyLiquidities(bytes calldata payload, uint256 deadline) external payable;
    function nextTokenId() external view returns (uint256);
}

/// @dev Permit2 interface for token approvals
interface IPermit2Uni {
    function approve(address token, address spender, uint160 amount, uint48 expiration) external;
}

/// @title LPManagerUni — Uniswap V4 CL liquidity management (Base, Ethereum, etc.)
/// @notice Identical logic to LPManager but uses Uniswap V4 PoolKey struct and interfaces.
/// @dev Deploy this for chains using Uniswap V4 (not PancakeSwap V4).
contract LPManagerUni {
    using SafeERC20 for IERC20;

    address public immutable awpRegistry;
    address public immutable poolManager;
    address public immutable positionManager;
    address public immutable permit2;
    IERC20 public immutable awpToken;

    uint24 public constant POOL_FEE = 10000;
    int24 public constant TICK_SPACING = 200;
    int24 public constant MIN_TICK = -887200;
    int24 public constant MAX_TICK = 887200;
    uint160 public constant MIN_SQRT_RATIO = 4295128739;
    uint160 public constant MAX_SQRT_RATIO = 1461446703485210103287273052203988822378723970342;

    mapping(address => bytes32) public alphaTokenToPoolId;
    mapping(address => uint256) public alphaTokenToTokenId;

    error NotAWPRegistry();
    error PoolAlreadyExists();

    modifier onlyAWPRegistry() {
        if (msg.sender != awpRegistry) revert NotAWPRegistry();
        _;
    }

    constructor(address awpRegistry_, address poolManager_, address positionManager_, address permit2_, address awpToken_) {
        awpRegistry = awpRegistry_;
        poolManager = poolManager_;
        positionManager = positionManager_;
        permit2 = permit2_;
        awpToken = IERC20(awpToken_);
    }

    function createPoolAndAddLiquidity(address alphaToken, uint256 awpAmount, uint256 alphaAmount)
        external
        onlyAWPRegistry
        returns (bytes32 poolId, uint256 lpTokenId)
    {
        if (alphaTokenToPoolId[alphaToken] != bytes32(0)) revert PoolAlreadyExists();

        address awp = address(awpToken);
        (address c0, address c1) = awp < alphaToken ? (awp, alphaToken) : (alphaToken, awp);
        (uint256 amt0, uint256 amt1) = awp < alphaToken ? (awpAmount, alphaAmount) : (alphaAmount, awpAmount);

        // Uniswap V4 PoolKey: (currency0, currency1, fee, tickSpacing, hooks)
        UniPoolKey memory poolKey = UniPoolKey({
            currency0: c0,
            currency1: c1,
            fee: POOL_FEE,
            tickSpacing: TICK_SPACING,
            hooks: address(0)
        });

        uint256 ratioX192 = FullMath.mulDiv(amt1, FixedPoint96.Q96 * FixedPoint96.Q96, amt0);
        uint160 sqrtPriceX96 = uint160(Math.sqrt(ratioX192));

        // Initialize pool via Uniswap V4 PoolManager
        IUniPoolManager(poolManager).initialize(poolKey, sqrtPriceX96);

        lpTokenId = _approveAndMint(poolKey, c0, c1, amt0, amt1, sqrtPriceX96);
        poolId = keccak256(abi.encode(poolKey));

        alphaTokenToPoolId[alphaToken] = poolId;
        alphaTokenToTokenId[alphaToken] = lpTokenId;
    }

    function _approveAndMint(
        UniPoolKey memory poolKey, address c0, address c1,
        uint256 amt0, uint256 amt1, uint160 sqrtPriceX96
    ) internal returns (uint256 lpTokenId) {
        IERC20(c0).forceApprove(permit2, amt0);
        IPermit2Uni(permit2).approve(c0, positionManager, uint160(amt0), uint48(block.timestamp + 600));
        IERC20(c1).forceApprove(permit2, amt1);
        IPermit2Uni(permit2).approve(c1, positionManager, uint160(amt1), uint48(block.timestamp + 600));

        uint128 liquidity = LiquidityAmounts.getLiquidityForAmounts(
            sqrtPriceX96, MIN_SQRT_RATIO, MAX_SQRT_RATIO, amt0, amt1
        );

        lpTokenId = IUniPositionManager(positionManager).nextTokenId();

        // MINT_POSITION(0x02) + SETTLE_PAIR(0x0d)
        bytes memory actions = abi.encodePacked(uint8(0x02), uint8(0x0d));
        bytes[] memory params = new bytes[](2);
        params[0] = abi.encode(poolKey, MIN_TICK, MAX_TICK, liquidity, uint128(amt0), uint128(amt1), address(this), bytes(""));
        params[1] = abi.encode(c0, c1);
        IUniPositionManager(positionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);
    }
}
