// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {WorknetManager, ICLPositionManager, IPermit2} from "./WorknetManager.sol";
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

/// @dev Uniswap V4 PoolManager — unlock + swap + settle + take
interface IUniPoolManager {
    struct SwapParams {
        bool zeroForOne;
        int256 amountSpecified;
        uint160 sqrtPriceLimitX96;
    }
    function unlock(bytes calldata data) external returns (bytes memory);
    function swap(UniPoolKey memory key, SwapParams memory params, bytes calldata hookData) external returns (int256);
    function sync(address currency) external;
    function settle() external payable returns (uint256);
    function take(address currency, address to, uint256 amount) external;
}

/// @title WorknetManagerUni — WorknetManager variant for Uniswap V4 (Base, Ethereum, etc.)
/// @notice Overrides initialize to construct Uniswap V4 PoolKey, and overrides all DEX
///         interaction functions to use the Uni V4 PoolKey format and StateView for reads.
///         dexConfig must encode 7 fields: (poolManager, positionManager, swapRouter, permit2, poolFee, tickSpacing, stateView)
/// @dev Deploy this implementation for chains using Uniswap V4 (not PancakeSwap V4).
///      _buybackAndBurn uses PoolManager.unlock + swap callback (PositionManager does NOT handle swap actions).
contract WorknetManagerUni is WorknetManager {
    using SafeERC20 for IERC20;

    // ── Uniswap V4 specific storage ──
    // NOTE: These are placed after WorknetManager's __gap[38], so they do NOT consume gap slots.
    // WorknetManagerUni is a separate implementation — not an upgrade of WorknetManager.
    // If WorknetManager is independently upgraded and expands into __gap, a new WorknetManagerUni
    // implementation must be deployed with matching storage layout.
    address public stateView;
    UniPoolKey public uniPoolKey;

    error NotPoolManager();
    error SlippageExceeded();

    /// @notice Override initialize to construct Uniswap V4 PoolKey and decode stateView
    /// @dev dexConfig_ must be abi.encode(poolManager, positionManager, swapRouter, permit2, poolFee, tickSpacing, stateView)
    function initialize(
        address awpRegistry_, address alphaToken_, address awpToken_, bytes32 poolId_, address admin_,
        bytes calldata dexConfig_
    ) external override initializer {
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

        // Shared init: AccessControl, ReentrancyGuard, storage, role grants
        _initializeBase(
            awpRegistry_, alphaToken_, awpToken_, poolId_, admin_,
            clPoolManager_, clPositionManager_, clSwapRouter_, permit2_, poolFee_, tickSpacing_
        );

        // Uni-specific storage
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
    }

    /// @dev Override: read slot0 from StateView (Uniswap V4 PoolManager does not expose getSlot0 directly)
    function _getSlot0() internal view override returns (uint160 sqrtPriceX96, int24 tick) {
        (sqrtPriceX96, tick,,) = IStateView(stateView).getSlot0(poolId);
    }

    /// @dev Override: encode with Uniswap V4 PoolKey (5 fields, different struct from PancakeSwap)
    function _addSingleSidedLiquidity(uint256 amount) internal override {
        (, int24 currentTick) = _getSlot0();

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

    /// @dev Override: use PoolManager.unlock + swap callback (Uniswap V4 PositionManager does NOT handle swap actions)
    function _buybackAndBurn(uint256 amount) internal override {
        bool zeroForOne = address(awpToken) < address(alphaToken);

        // Read current pool price via virtual _getSlot0() for slippage protection
        (uint160 sqrtPriceX96,) = _getSlot0();
        uint256 expectedOut;
        if (zeroForOne) {
            expectedOut = FullMath.mulDiv(FullMath.mulDiv(amount, sqrtPriceX96, 1 << 96), sqrtPriceX96, 1 << 96);
        } else {
            expectedOut = FullMath.mulDiv(FullMath.mulDiv(amount, 1 << 96, sqrtPriceX96), 1 << 96, sqrtPriceX96);
        }
        uint128 minOut = uint128(expectedOut * (10000 - slippageBps) / 10000);

        uint256 before = alphaToken.balanceOf(address(this));

        // Unlock PoolManager → triggers unlockCallback where swap + settle + take happen
        IUniPoolManager(clPoolManager).unlock(abi.encode(amount, zeroForOne, minOut));

        uint256 received = alphaToken.balanceOf(address(this)) - before;
        if (received > 0) alphaToken.burn(received);
        emit BuybackBurned(amount, received);
    }

    /// @dev Called by PoolManager during unlock — executes swap, settles input, takes output
    function unlockCallback(bytes calldata data) external returns (bytes memory) {
        if (msg.sender != clPoolManager) revert NotPoolManager();

        (uint256 amount, bool zeroForOne, uint128 minOut) = abi.decode(data, (uint256, bool, uint128));

        UniPoolKey memory pk = uniPoolKey;

        // Execute swap via PoolManager (exact input)
        int256 delta = IUniPoolManager(clPoolManager).swap(
            pk,
            IUniPoolManager.SwapParams({
                zeroForOne: zeroForOne,
                amountSpecified: -int256(amount), // negative = exact input
                sqrtPriceLimitX96: zeroForOne
                    ? TickMath.getSqrtRatioAtTick(TickMath.MIN_TICK) + 1
                    : TickMath.getSqrtRatioAtTick(TickMath.MAX_TICK) - 1
            }),
            bytes("")
        );

        // Decode BalanceDelta: upper 128 bits = amount0, lower 128 bits = amount1
        int128 delta0 = int128(delta >> 128);
        int128 delta1 = int128(delta);

        // Determine input/output amounts from deltas
        // For zeroForOne: delta0 < 0 (we owe token0), delta1 > 0 (we receive token1)
        // For !zeroForOne: delta0 > 0 (we receive token0), delta1 < 0 (we owe token1)
        address inputCurrency = zeroForOne ? pk.currency0 : pk.currency1;
        address outputCurrency = zeroForOne ? pk.currency1 : pk.currency0;
        uint256 inputAmt = uint256(uint128(zeroForOne ? -delta0 : -delta1));
        uint256 outputAmt = uint256(uint128(zeroForOne ? delta1 : delta0));

        // Slippage check
        if (outputAmt < minOut) revert SlippageExceeded();

        // Settle input: sync balance → transfer tokens → settle diff
        IUniPoolManager(clPoolManager).sync(inputCurrency);
        IERC20(inputCurrency).safeTransfer(clPoolManager, inputAmt);
        IUniPoolManager(clPoolManager).settle();

        // Take output: PoolManager transfers tokens to us
        IUniPoolManager(clPoolManager).take(outputCurrency, address(this), outputAmt);

        return bytes("");
    }
}
