// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {WorknetManagerBase, IWorknetToken, IERC20} from "./WorknetManagerBase.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {TickMath} from "infinity-core/src/pool-cl/libraries/TickMath.sol";
import {FullMath} from "infinity-core/src/pool-cl/libraries/FullMath.sol";

struct PoolKey {
    address currency0;
    address currency1;
    address hooks;
    address poolManager;
    uint24 fee;
    bytes32 parameters;
}

interface ICLPoolManager {
    function getSlot0(bytes32 id) external view returns (uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee);
}

interface ICLPositionManager {
    function modifyLiquidities(bytes calldata payload, uint256 deadline) external payable;
    function nextTokenId() external view returns (uint256);
}

interface ICLSwapRouter {
    function executeActions(bytes calldata data) external payable;
}

interface IPermit2 {
    function approve(address token, address spender, uint160 amount, uint48 expiration) external;
}

/// @title WorknetManager — PancakeSwap V4 worknet manager
contract WorknetManager is WorknetManagerBase {
    using SafeERC20 for IERC20;

    error SlippageExceeded();
    error AmountExceedsPermit2Limit();

    address public immutable clPoolManager;
    address public immutable clPositionManager;
    address public immutable clSwapRouter;
    address public immutable permit2;

    uint24 public constant POOL_FEE = 10000;
    int24 public constant TICK_SPACING = 200;

    uint8 internal constant ACT_CL_MINT_POSITION = 0x02;
    uint8 internal constant ACT_CL_INCREASE_LIQUIDITY = 0x00;   // [B1]
    uint8 internal constant ACT_CL_SWAP_EXACT_IN_SINGLE = 0x06;
    uint8 internal constant ACT_SETTLE_ALL = 0x0c;
    uint8 internal constant ACT_SETTLE_PAIR = 0x0d;
    uint8 internal constant ACT_TAKE_ALL = 0x0f;

    PoolKey public poolKey;

    uint256[40] private __gap;

    constructor(address permit2_, address clPoolManager_, address clPositionManager_, address clSwapRouter_) {
        permit2 = permit2_;
        clPoolManager = clPoolManager_;
        clPositionManager = clPositionManager_;
        clSwapRouter = clSwapRouter_;
        _disableInitializers();
    }

    function initialize(address worknetToken_, bytes32 poolId_, address admin_) external initializer {
        __WorknetManagerBase_init(worknetToken_, poolId_, admin_);

        (address c0, address c1) = awpToken < worknetToken_
            ? (awpToken, worknetToken_)
            : (worknetToken_, awpToken);
        poolKey = PoolKey({
            currency0: c0, currency1: c1, hooks: address(0), poolManager: clPoolManager,
            fee: POOL_FEE, parameters: bytes32(uint256(int256(TICK_SPACING)) << 16)
        });
    }

    // ── [B1] Add liquidity: reuse existing position if tick range matches ──

    function _addSingleSidedLiquidity(uint256 amount) internal override {
        PoolKey memory pk = poolKey;
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

        if (amount > type(uint160).max) revert AmountExceedsPermit2Limit();
        IERC20(awpToken).forceApprove(permit2, amount);
        IPermit2(permit2).approve(awpToken, clPositionManager, uint160(amount), uint48(block.timestamp));

        // Reuse existing position if tick range matches
        bool reuse = _lastLpTokenId != 0 && tickLower == _lastTickLower && tickUpper == _lastTickUpper;

        bytes memory actions;
        bytes[] memory params = new bytes[](2);

        if (reuse) {
            actions = abi.encodePacked(ACT_CL_INCREASE_LIQUIDITY, ACT_SETTLE_PAIR);
            params[0] = abi.encode(
                _lastLpTokenId, liquidity,
                awpIs0 ? uint128(amount) : uint128(0),
                awpIs0 ? uint128(0) : uint128(amount),
                bytes("")
            );
        } else {
            uint256 tokenId = ICLPositionManager(clPositionManager).nextTokenId();
            actions = abi.encodePacked(ACT_CL_MINT_POSITION, ACT_SETTLE_PAIR);
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

        ICLPositionManager(clPositionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);

        if (reuse) {
            emit LiquidityIncreased(_lastLpTokenId, amount);
        } else {
            emit LiquidityAdded(_lastLpTokenId, amount);
        }
    }

    // ── [B2] Buyback: accept minAmountOut for MEV protection ──

    function _buybackAndBurn(uint256 amount, uint256 minAmountOut) internal override {
        PoolKey memory pk = poolKey;
        if (amount > type(uint160).max) revert AmountExceedsPermit2Limit();
        IERC20(awpToken).forceApprove(permit2, amount);
        IPermit2(permit2).approve(awpToken, clSwapRouter, uint160(amount), uint48(block.timestamp));

        bool zeroForOne = awpToken < address(worknetToken);

        // If no explicit minOut, compute from slippage + spot price
        uint128 minOut;
        if (minAmountOut > 0) {
            if (minAmountOut > type(uint128).max) revert SlippageExceeded();
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

        bytes memory actions = abi.encodePacked(ACT_CL_SWAP_EXACT_IN_SINGLE, ACT_SETTLE_ALL, ACT_TAKE_ALL);
        bytes[] memory params = new bytes[](3);
        params[0] = abi.encode(pk, zeroForOne, uint128(amount), minOut, bytes(""));
        params[1] = abi.encode(awpToken, amount);
        params[2] = abi.encode(address(worknetToken), 0);

        uint256 before = worknetToken.balanceOf(address(this));
        ICLSwapRouter(clSwapRouter).executeActions(abi.encode(actions, params));
        uint256 received = worknetToken.balanceOf(address(this)) - before;

        if (received > 0) worknetToken.burn(received);
        emit BuybackBurned(amount, received);
    }

    function _getSlot0() internal view returns (uint160 sqrtPriceX96, int24 tick) {
        (sqrtPriceX96, tick,,) = ICLPoolManager(clPoolManager).getSlot0(poolId);
    }
}
