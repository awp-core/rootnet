// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {WorknetManager, ICLPositionManager, IPermit2, IWorknetToken} from "./WorknetManager.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
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

/// @title WorknetManagerUni — WorknetManager variant for Uniswap V4
/// @dev DEX addresses are immutable (constructor). stateView is also immutable.
contract WorknetManagerUni is WorknetManager {
    using SafeERC20 for IERC20;

    // ── Uniswap V4 specific (immutable) ──
    address public immutable stateView;

    // ── Per-worknet storage (after WorknetManager's __gap) ──
    UniPoolKey public uniPoolKey;

    /// @dev Reserved storage gap for future WorknetManagerUni upgrades
    uint256[46] private __uniGap;

    error NotPoolManager();
    error SlippageExceeded();

    /// @param permit2_ Permit2
    /// @param clPoolManager_ Uniswap V4 PoolManager
    /// @param clPositionManager_ Uniswap V4 PositionManager
    /// @param clSwapRouter_ Not used for Uni (swap via unlock callback), pass address(0)
    /// @param stateView_ Uniswap V4 StateView for reading pool state
    constructor(
        address permit2_, address clPoolManager_, address clPositionManager_, address clSwapRouter_, address stateView_
    ) WorknetManager(permit2_, clPoolManager_, clPositionManager_, clSwapRouter_) {
        stateView = stateView_;
    }

    /// @notice Initialize per-worknet state + Uniswap V4 PoolKey
    function initialize(address worknetToken_, bytes32 poolId_, address admin_) external override initializer {
        __UUPSUpgradeable_init();
        __AccessControl_init();
        __ReentrancyGuard_init();

        worknetToken = IWorknetToken(worknetToken_);
        poolId = poolId_;

        _grantRole(DEFAULT_ADMIN_ROLE, admin_);
        _grantRole(MERKLE_ROLE, admin_);
        _grantRole(STRATEGY_ROLE, admin_);
        _grantRole(TRANSFER_ROLE, admin_);

        slippageBps = 500;

        // Construct Uniswap V4 PoolKey
        (address c0, address c1) = awpToken < worknetToken_
            ? (awpToken, worknetToken_)
            : (worknetToken_, awpToken);
        uniPoolKey = UniPoolKey({
            currency0: c0, currency1: c1, fee: POOL_FEE, tickSpacing: TICK_SPACING, hooks: address(0)
        });
    }

    function _getSlot0() internal view override returns (uint160 sqrtPriceX96, int24 tick) {
        (sqrtPriceX96, tick,,) = IStateView(stateView).getSlot0(poolId);
    }

    function _addSingleSidedLiquidity(uint256 amount) internal override {
        (, int24 currentTick) = _getSlot0();

        int24 ts = TICK_SPACING;
        int24 aligned = (currentTick / ts) * ts;
        if (aligned > currentTick) aligned -= ts;

        int24 minTick = (-887272 / ts) * ts;
        int24 maxTick = (887272 / ts) * ts;

        bool awpIs0 = awpToken < address(worknetToken);

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

        IERC20(awpToken).forceApprove(permit2, amount);
        IPermit2(permit2).approve(awpToken, clPositionManager, uint160(amount), uint48(block.timestamp + 600));

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

    function _buybackAndBurn(uint256 amount) internal override {
        bool zeroForOne = awpToken < address(worknetToken);

        (uint160 sqrtPriceX96,) = _getSlot0();
        uint256 expectedOut;
        if (zeroForOne) {
            expectedOut = FullMath.mulDiv(FullMath.mulDiv(amount, sqrtPriceX96, 1 << 96), sqrtPriceX96, 1 << 96);
        } else {
            expectedOut = FullMath.mulDiv(FullMath.mulDiv(amount, 1 << 96, sqrtPriceX96), 1 << 96, sqrtPriceX96);
        }
        uint128 minOut = uint128(expectedOut * (10000 - slippageBps) / 10000);

        uint256 before = worknetToken.balanceOf(address(this));
        IUniPoolManager(clPoolManager).unlock(abi.encode(amount, zeroForOne, minOut));

        uint256 received = worknetToken.balanceOf(address(this)) - before;
        if (received > 0) worknetToken.burn(received);
        emit BuybackBurned(amount, received);
    }

    function unlockCallback(bytes calldata data) external returns (bytes memory) {
        if (msg.sender != clPoolManager) revert NotPoolManager();

        (uint256 amount, bool zeroForOne, uint128 minOut) = abi.decode(data, (uint256, bool, uint128));

        UniPoolKey memory pk = uniPoolKey;

        int256 delta = IUniPoolManager(clPoolManager).swap(
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

        IUniPoolManager(clPoolManager).sync(inputCurrency);
        IERC20(inputCurrency).safeTransfer(clPoolManager, inputAmt);
        IUniPoolManager(clPoolManager).settle();

        IUniPoolManager(clPoolManager).take(outputCurrency, address(this), outputAmt);

        return bytes("");
    }
}
