// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {SubnetManager, IAlphaToken, ICLPositionManager, IPermit2} from "./SubnetManager.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {TickMath} from "infinity-core/src/pool-cl/libraries/TickMath.sol";
import {FullMath} from "infinity-core/src/pool-cl/libraries/FullMath.sol";

/// @dev Uniswap V4 PoolKey struct (5 fields, field order: currency0, currency1, fee, tickSpacing, hooks)
struct UniPoolKey {
    address currency0;
    address currency1;
    uint24 fee;
    int24 tickSpacing;
    address hooks;
}

/// @dev Uniswap V4 StateView — read pool state (getSlot0)
interface IStateView {
    function getSlot0(bytes32 id) external view returns (uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee);
}

/// @title SubnetManagerUni — SubnetManager variant for Uniswap V4 (Base, Ethereum, etc.)
/// @notice Overrides initialize to construct Uniswap V4 PoolKey, and overrides all DEX
///         interaction functions to use the Uni V4 PoolKey format and StateView for reads.
///         dexConfig must encode 7 fields: (poolManager, positionManager, swapRouter, permit2, poolFee, tickSpacing, stateView)
/// @dev Deploy this implementation for chains using Uniswap V4 (not PancakeSwap V4).
contract SubnetManagerUni is SubnetManager {
    using SafeERC20 for IERC20;

    // ── Uniswap V4 specific storage (uses gap space) ──
    address public stateView;
    UniPoolKey public uniPoolKey;

    /// @notice Override initialize to construct Uniswap V4 PoolKey and decode stateView
    /// @dev dexConfig_ must be abi.encode(poolManager, positionManager, swapRouter, permit2, poolFee, tickSpacing, stateView)
    function initialize(
        address alphaToken_, address awpToken_, bytes32 poolId_, address admin_,
        bytes calldata dexConfig_
    ) external override initializer {
        __AccessControl_init();
        __ReentrancyGuard_init();

        alphaToken = IAlphaToken(alphaToken_);
        awpToken = IERC20(awpToken_);
        poolId = poolId_;

        // Decode DEX addresses — 7 fields for Uniswap V4 (extra: stateView)
        (
            address clPoolManager_,
            address clPositionManager_,
            address clSwapRouter_,
            address permit2_,
            uint24 poolFee_,
            int24 tickSpacing_,
            address stateView_
        ) = abi.decode(dexConfig_, (address, address, address, address, uint24, int24, address));

        clPoolManager = clPoolManager_;
        clPositionManager = clPositionManager_;
        clSwapRouter = clSwapRouter_;
        permit2 = permit2_;
        poolFee = poolFee_;
        tickSpacing = tickSpacing_;
        stateView = stateView_;

        // Construct Uniswap V4 PoolKey (5 fields, different order from PancakeSwap)
        (address c0, address c1) = awpToken_ < alphaToken_
            ? (awpToken_, alphaToken_)
            : (alphaToken_, awpToken_);
        uniPoolKey = UniPoolKey({
            currency0: c0,
            currency1: c1,
            fee: poolFee_,
            tickSpacing: tickSpacing_,
            hooks: address(0)
        });

        _grantRole(DEFAULT_ADMIN_ROLE, admin_);
    }

    /// @dev Override: use Uni V4 PoolKey and StateView for single-sided liquidity
    function _addSingleSidedLiquidity(uint256 amount) internal override {
        (, int24 currentTick,,) = IStateView(stateView).getSlot0(poolId);

        int24 ts = tickSpacing;
        int24 aligned = (currentTick / ts) * ts;
        if (aligned > currentTick) aligned -= ts;

        int24 minTick = (-887272 / ts) * ts;
        int24 maxTick = (887272 / ts) * ts;

        bool awpIs0 = address(awpToken) < address(alphaToken);

        int24 tickLower;
        int24 tickUpper;
        if (awpIs0) {
            tickLower = aligned + ts;
            tickUpper = maxTick;
        } else {
            tickUpper = aligned < currentTick ? aligned : aligned - ts;
            tickLower = minTick;
        }

        uint160 sqrtLower = TickMath.getSqrtRatioAtTick(tickLower);
        uint160 sqrtUpper = TickMath.getSqrtRatioAtTick(tickUpper);

        uint128 liquidity = awpIs0
            ? LiquidityAmounts.getLiquidityForAmount0(sqrtLower, sqrtUpper, amount)
            : LiquidityAmounts.getLiquidityForAmount1(sqrtLower, sqrtUpper, amount);

        IERC20(address(awpToken)).forceApprove(permit2, amount);
        IPermit2(permit2).approve(address(awpToken), clPositionManager, uint160(amount), uint48(block.timestamp + 600));

        uint256 tokenId = ICLPositionManager(clPositionManager).nextTokenId();
        bytes memory actions = abi.encodePacked(ACT_CL_MINT_POSITION, ACT_SETTLE_PAIR);
        bytes[] memory params = new bytes[](2);

        // Encode with Uniswap V4 PoolKey
        UniPoolKey memory pk = uniPoolKey;
        params[0] = abi.encode(
            pk, tickLower, tickUpper, liquidity,
            awpIs0 ? uint128(amount) : uint128(0),
            awpIs0 ? uint128(0) : uint128(amount),
            address(this), bytes("")
        );
        params[1] = abi.encode(pk.currency0, pk.currency1);

        ICLPositionManager(clPositionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);
        emit LiquidityAdded(tokenId, amount);
    }

    /// @dev Override: use Uni V4 PoolKey, StateView, and route swap through PositionManager
    function _buybackAndBurn(uint256 amount) internal override {
        IERC20(address(awpToken)).forceApprove(permit2, amount);
        IPermit2(permit2).approve(address(awpToken), clPositionManager, uint160(amount), uint48(block.timestamp + 600));

        bool zeroForOne = address(awpToken) < address(alphaToken);

        // Read current pool price from StateView
        (uint160 sqrtPriceX96,,,) = IStateView(stateView).getSlot0(poolId);
        uint256 expectedOut;
        if (zeroForOne) {
            expectedOut = FullMath.mulDiv(amount, uint256(sqrtPriceX96) * uint256(sqrtPriceX96), 1 << 192);
        } else {
            expectedOut = FullMath.mulDiv(amount, 1 << 192, uint256(sqrtPriceX96) * uint256(sqrtPriceX96));
        }
        uint128 minOut = uint128(expectedOut * 95 / 100);

        bytes memory actions = abi.encodePacked(ACT_CL_SWAP_EXACT_IN_SINGLE, ACT_SETTLE_ALL, ACT_TAKE_ALL);
        bytes[] memory params = new bytes[](3);

        // Encode with Uniswap V4 PoolKey
        UniPoolKey memory pk = uniPoolKey;
        params[0] = abi.encode(pk, zeroForOne, uint128(amount), minOut, bytes(""));
        params[1] = abi.encode(address(awpToken), amount);
        params[2] = abi.encode(address(alphaToken), 0);

        uint256 before = alphaToken.balanceOf(address(this));
        // Route through PositionManager (Uniswap V4)
        ICLPositionManager(clPositionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);
        uint256 received = alphaToken.balanceOf(address(this)) - before;

        if (received > 0) alphaToken.burn(received);
        emit BuybackBurned(amount, received);
    }
}
