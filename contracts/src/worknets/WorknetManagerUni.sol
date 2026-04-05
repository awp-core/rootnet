// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {WorknetManagerBase, IWorknetToken, IERC20} from "./WorknetManagerBase.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {TickMath} from "infinity-core/src/pool-cl/libraries/TickMath.sol";
import {FullMath} from "infinity-core/src/pool-cl/libraries/FullMath.sol";

struct UniPoolKey {
    address currency0;
    address currency1;
    uint24 fee;
    int24 tickSpacing;
    address hooks;
}

interface IStateView {
    function getSlot0(bytes32 id) external view returns (uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee);
}

interface IUniPoolManager {
    struct SwapParams { bool zeroForOne; int256 amountSpecified; uint160 sqrtPriceLimitX96; }
    function unlock(bytes calldata data) external returns (bytes memory);
    function swap(UniPoolKey memory key, SwapParams memory params, bytes calldata hookData) external returns (int256);
    function sync(address currency) external;
    function settle() external payable returns (uint256);
    function take(address currency, address to, uint256 amount) external;
}

interface IUniPositionManager {
    function modifyLiquidities(bytes calldata payload, uint256 deadline) external payable;
    function nextTokenId() external view returns (uint256);
}

interface IPermit2 {
    function approve(address token, address spender, uint160 amount, uint48 expiration) external;
}

/// @title WorknetManagerUni — Uniswap V4 worknet manager
contract WorknetManagerUni is WorknetManagerBase {
    using SafeERC20 for IERC20;

    address public immutable poolManager;
    address public immutable positionManager;
    address public immutable stateView;
    address public immutable permit2;

    uint24 public constant POOL_FEE = 10000;
    int24 public constant TICK_SPACING = 200;

    uint8 internal constant ACT_MINT_POSITION = 0x02;
    uint8 internal constant ACT_INCREASE_LIQUIDITY = 0x00;   // [B1]
    uint8 internal constant ACT_SETTLE_PAIR = 0x0d;

    UniPoolKey public uniPoolKey;

    uint256[40] private __gap;

    error NotPoolManager();
    error SlippageExceeded();

    constructor(address permit2_, address poolManager_, address positionManager_, address stateView_) {
        permit2 = permit2_;
        poolManager = poolManager_;
        positionManager = positionManager_;
        stateView = stateView_;
        _disableInitializers();
    }

    function initialize(address worknetToken_, bytes32 poolId_, address admin_) external initializer {
        __WorknetManagerBase_init(worknetToken_, poolId_, admin_);

        (address c0, address c1) = awpToken < worknetToken_
            ? (awpToken, worknetToken_)
            : (worknetToken_, awpToken);
        uniPoolKey = UniPoolKey({
            currency0: c0, currency1: c1, fee: POOL_FEE, tickSpacing: TICK_SPACING, hooks: address(0)
        });
    }

    // ── [B1] Add liquidity: reuse existing position if tick range matches ──

    function _addSingleSidedLiquidity(uint256 amount) internal override {
        (, int24 currentTick) = _getSlot0();

        int24 ts = TICK_SPACING;
        int24 aligned = (currentTick / ts) * ts;
        if (aligned > currentTick) aligned -= ts;

        bool awpIs0 = awpToken < address(worknetToken);

        int24 tickLower;
        int24 tickUpper;
        if (awpIs0) {
            tickLower = aligned + ts;
            tickUpper = (887272 / ts) * ts;
        } else {
            tickUpper = aligned < currentTick ? aligned : aligned - ts;
            tickLower = (-887272 / ts) * ts;
        }

        uint160 sqrtLower = TickMath.getSqrtRatioAtTick(tickLower);
        uint160 sqrtUpper = TickMath.getSqrtRatioAtTick(tickUpper);

        uint128 liquidity = awpIs0
            ? LiquidityAmounts.getLiquidityForAmount0(sqrtLower, sqrtUpper, amount)
            : LiquidityAmounts.getLiquidityForAmount1(sqrtLower, sqrtUpper, amount);

        if (liquidity == 0) return;  // [S2] amount too small for any liquidity

        IERC20(awpToken).forceApprove(permit2, amount);
        IPermit2(permit2).approve(awpToken, positionManager, uint160(amount), uint48(block.timestamp));

        bool reuse = _lastLpTokenId != 0 && tickLower == _lastTickLower && tickUpper == _lastTickUpper;

        UniPoolKey memory pk = uniPoolKey;
        bytes memory actions;
        bytes[] memory params = new bytes[](2);

        if (reuse) {
            actions = abi.encodePacked(ACT_INCREASE_LIQUIDITY, ACT_SETTLE_PAIR);
            params[0] = abi.encode(
                _lastLpTokenId, liquidity,
                awpIs0 ? uint128(amount) : uint128(0),
                awpIs0 ? uint128(0) : uint128(amount),
                bytes("")
            );
        } else {
            uint256 tokenId = IUniPositionManager(positionManager).nextTokenId();
            actions = abi.encodePacked(ACT_MINT_POSITION, ACT_SETTLE_PAIR);
            params[0] = abi.encode(
                pk, tickLower, tickUpper, liquidity,
                awpIs0 ? uint128(amount) : uint128(0),
                awpIs0 ? uint128(0) : uint128(amount),
                address(this), bytes("")
            );
            _lastLpTokenId = tokenId;
            _lastTickLower = tickLower;
            _lastTickUpper = tickUpper;
        }
        params[1] = abi.encode(pk.currency0, pk.currency1);

        IUniPositionManager(positionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);

        if (reuse) {
            emit LiquidityIncreased(_lastLpTokenId, amount);
        } else {
            emit LiquidityAdded(_lastLpTokenId, amount);
        }
    }

    // ── [B2] Buyback: accept minAmountOut for MEV protection ──

    function _buybackAndBurn(uint256 amount, uint256 minAmountOut) internal override {
        bool zeroForOne = awpToken < address(worknetToken);

        uint128 minOut;
        if (minAmountOut > 0) {
            minOut = uint128(minAmountOut);
        } else {
            (uint160 sqrtPriceX96,) = _getSlot0();
            uint256 expectedOut;
            if (zeroForOne) {
                expectedOut = FullMath.mulDiv(FullMath.mulDiv(amount, sqrtPriceX96, 1 << 96), sqrtPriceX96, 1 << 96);
            } else {
                expectedOut = FullMath.mulDiv(FullMath.mulDiv(amount, 1 << 96, sqrtPriceX96), 1 << 96, sqrtPriceX96);
            }
            minOut = uint128(expectedOut * (10000 - slippageBps) / 10000);
        }

        uint256 before = worknetToken.balanceOf(address(this));
        IUniPoolManager(poolManager).unlock(abi.encode(amount, zeroForOne, minOut));

        uint256 received = worknetToken.balanceOf(address(this)) - before;
        if (received > 0) worknetToken.burn(received);
        emit BuybackBurned(amount, received);
    }

    function unlockCallback(bytes calldata data) external returns (bytes memory) {
        if (msg.sender != poolManager) revert NotPoolManager();

        (uint256 amount, bool zeroForOne, uint128 minOut) = abi.decode(data, (uint256, bool, uint128));

        UniPoolKey memory pk = uniPoolKey;

        int256 delta = IUniPoolManager(poolManager).swap(
            pk,
            IUniPoolManager.SwapParams({
                zeroForOne: zeroForOne,
                amountSpecified: -int256(amount),
                sqrtPriceLimitX96: zeroForOne
                    ? TickMath.getSqrtRatioAtTick(TickMath.MIN_TICK) + 1
                    : TickMath.getSqrtRatioAtTick(TickMath.MAX_TICK) - 1
            }),
            bytes("")
        );

        int128 delta0 = int128(delta >> 128);
        int128 delta1 = int128(delta);

        address inputCurrency = zeroForOne ? pk.currency0 : pk.currency1;
        address outputCurrency = zeroForOne ? pk.currency1 : pk.currency0;
        uint256 inputAmt = uint256(uint128(zeroForOne ? -delta0 : -delta1));
        uint256 outputAmt = uint256(uint128(zeroForOne ? delta1 : delta0));

        if (outputAmt < minOut) revert SlippageExceeded();

        IUniPoolManager(poolManager).sync(inputCurrency);
        IERC20(inputCurrency).safeTransfer(poolManager, inputAmt);
        IUniPoolManager(poolManager).settle();

        IUniPoolManager(poolManager).take(outputCurrency, address(this), outputAmt);

        return bytes("");
    }

    function _getSlot0() internal view returns (uint160 sqrtPriceX96, int24 tick) {
        (sqrtPriceX96, tick,,) = IStateView(stateView).getSlot0(poolId);
    }
}
