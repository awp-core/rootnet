// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {LPManagerBase, IPermit2} from "./LPManagerBase.sol";

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
    function getSlot0(bytes32 id) external view returns (uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee);
}

/// @dev Uniswap V4 PositionManager interface
interface IUniPositionManager {
    function modifyLiquidities(bytes calldata payload, uint256 deadline) external payable;
    function nextTokenId() external view returns (uint256);
}

/// @title LPManagerUni — Uniswap V4 CL liquidity management (Base, Ethereum, etc.)
/// @notice Only AWPRegistry may call; LP is permanently locked and cannot be withdrawn
/// @dev Deploy this for chains using Uniswap V4 (not PancakeSwap V4).
contract LPManagerUni is LPManagerBase {
    using SafeERC20 for IERC20;

    // ══════════════════════════════════════════════
    //  Uniswap-specific immutable storage
    // ══════════════════════════════════════════════

    /// @notice Uniswap V4 PoolManager address
    address public immutable poolManager;
    /// @notice Uniswap V4 PositionManager address
    address public immutable positionManager;

    /// @notice Constructor
    /// @param awpRegistry_ AWPRegistry contract address
    /// @param poolManager_ Uniswap V4 PoolManager address
    /// @param positionManager_ Uniswap V4 PositionManager address
    /// @param permit2_ Permit2 contract address
    /// @param awpToken_ AWP token contract address
    constructor(address awpRegistry_, address poolManager_, address positionManager_, address permit2_, address awpToken_)
        LPManagerBase(awpRegistry_, permit2_, awpToken_)
    {
        poolManager = poolManager_;
        positionManager = positionManager_;
    }

    /// @dev Initialize Uniswap V4 pool
    function _initializePool(address c0, address c1, uint160 sqrtPriceX96) internal override {
        UniPoolKey memory poolKey = _buildPoolKey(c0, c1);
        IUniPoolManager(poolManager).initialize(poolKey, sqrtPriceX96);
    }

    /// @dev Permit2 approval + Uniswap V4 LP position minting
    function _approveAndMint(
        address c0, address c1,
        uint256 amt0, uint256 amt1, uint160 sqrtPriceX96
    ) internal override returns (uint256 lpTokenId) {
        // Approve tokens to PositionManager via Permit2
        IERC20(c0).forceApprove(permit2, amt0);
        IPermit2(permit2).approve(c0, positionManager, uint160(amt0), uint48(block.timestamp + 600));
        IERC20(c1).forceApprove(permit2, amt1);
        IPermit2(permit2).approve(c1, positionManager, uint160(amt1), uint48(block.timestamp + 600));

        // Compute full-range liquidity
        uint128 liquidity = LiquidityAmounts.getLiquidityForAmounts(
            sqrtPriceX96, MIN_SQRT_RATIO, MAX_SQRT_RATIO, amt0, amt1
        );

        // Record nextTokenId before mint
        lpTokenId = IUniPositionManager(positionManager).nextTokenId();

        // Encode mint action: MINT_POSITION(0x02) + SETTLE_PAIR(0x0d)
        UniPoolKey memory poolKey = _buildPoolKey(c0, c1);
        bytes memory actions = abi.encodePacked(uint8(0x02), uint8(0x0d));
        bytes[] memory params = new bytes[](2);
        params[0] = abi.encode(poolKey, MIN_TICK, MAX_TICK, liquidity, uint128(amt0), uint128(amt1), address(this), bytes(""));
        params[1] = abi.encode(c0, c1);
        IUniPositionManager(positionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);
    }

    /// @dev Compute Uniswap V4 pool ID: keccak256(UniPoolKey)
    function _computePoolId(address c0, address c1) internal view override returns (bytes32) {
        UniPoolKey memory poolKey = _buildPoolKey(c0, c1);
        return keccak256(abi.encode(poolKey));
    }

    /// @dev Compound accumulated fees back into the LP position (Uniswap V4)
    ///      Step 1: DECREASE_LIQUIDITY(0x01) with 0 delta → collect fees to LPManager
    ///      Step 2: INCREASE_LIQUIDITY(0x00) with collected amounts → reinvest as liquidity
    function _getCurrentSqrtPrice(address c0, address c1) internal view override returns (uint160) {
        UniPoolKey memory poolKey = _buildPoolKey(c0, c1);
        bytes32 pid = keccak256(abi.encode(poolKey));
        (uint160 sqrtPriceX96,,,) = IUniPoolManager(poolManager).getSlot0(pid);
        return sqrtPriceX96;
    }

    function _compoundFees(uint256 tokenId, address c0, address c1, uint160 sqrtPriceX96) internal override {
        // Step 1: Collect fees
        {
            bytes memory actions = abi.encodePacked(uint8(0x01), uint8(0x11));
            bytes[] memory params = new bytes[](2);
            params[0] = abi.encode(tokenId, uint256(0), uint128(0), uint128(0), bytes(""));
            params[1] = abi.encode(c0, c1);
            IUniPositionManager(positionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);
        }

        uint256 bal0 = IERC20(c0).balanceOf(address(this));
        uint256 bal1 = IERC20(c1).balanceOf(address(this));
        if (bal0 == 0 && bal1 == 0) return;

        // Step 2: Reinvest fees as liquidity
        if (bal0 > 0) {
            IERC20(c0).forceApprove(permit2, bal0);
            IPermit2(permit2).approve(c0, positionManager, uint160(bal0), uint48(block.timestamp + 600));
        }
        if (bal1 > 0) {
            IERC20(c1).forceApprove(permit2, bal1);
            IPermit2(permit2).approve(c1, positionManager, uint160(bal1), uint48(block.timestamp + 600));
        }

        uint128 liquidity = LiquidityAmounts.getLiquidityForAmounts(
            sqrtPriceX96, MIN_SQRT_RATIO, MAX_SQRT_RATIO, bal0, bal1
        );
        if (liquidity == 0) return;

        // INCREASE_LIQUIDITY (0x00): params = (tokenId, liquidity, amount0Max, amount1Max, hookData)
        bytes memory actions = abi.encodePacked(uint8(0x00), uint8(0x0d));
        bytes[] memory params = new bytes[](2);
        params[0] = abi.encode(tokenId, liquidity, uint128(bal0), uint128(bal1), bytes(""));
        params[1] = abi.encode(c0, c1);
        IUniPositionManager(positionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);
    }

    /// @dev Build Uniswap V4 UniPoolKey
    function _buildPoolKey(address c0, address c1) internal view returns (UniPoolKey memory) {
        return UniPoolKey({
            currency0: c0,
            currency1: c1,
            fee: POOL_FEE,
            tickSpacing: TICK_SPACING,
            hooks: address(0)
        });
    }
}
