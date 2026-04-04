// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {LPManagerBase, IPermit2} from "./LPManagerBase.sol";

struct UniPoolKey {
    address currency0;
    address currency1;
    uint24 fee;
    int24 tickSpacing;
    address hooks;
}

interface IUniPoolManager {
    function initialize(UniPoolKey calldata key, uint160 sqrtPriceX96) external returns (int24 tick);
    function getSlot0(bytes32 id) external view returns (uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee);
}

interface IUniPositionManager {
    function modifyLiquidities(bytes calldata payload, uint256 deadline) external payable;
    function nextTokenId() external view returns (uint256);
}

/// @title LPManagerUni — UUPS upgradeable Uniswap V4 CL liquidity management
/// @dev DEX addresses are immutable in the impl bytecode. Proxy upgrades to chain-specific impl.
contract LPManagerUni is LPManagerBase {
    using SafeERC20 for IERC20;

    address public immutable poolManager;
    address public immutable positionManager;

    /// @param permit2_ Permit2 address
    /// @param poolManager_ Uniswap V4 PoolManager
    /// @param positionManager_ Uniswap V4 PositionManager
    constructor(address permit2_, address poolManager_, address positionManager_)
        LPManagerBase(permit2_)
    {
        poolManager = poolManager_;
        positionManager = positionManager_;
    }

    function _initializePool(address c0, address c1, uint160 sqrtPriceX96) internal override {
        IUniPoolManager(poolManager).initialize(_buildPoolKey(c0, c1), sqrtPriceX96);
    }

    function _approveAndMint(
        address c0, address c1,
        uint256 amt0, uint256 amt1, uint160 sqrtPriceX96
    ) internal override returns (uint256 lpTokenId) {
        address _permit2 = permit2;
        address _posMgr = positionManager;

        IERC20(c0).forceApprove(_permit2, amt0);
        IPermit2(_permit2).approve(c0, _posMgr, uint160(amt0), uint48(block.timestamp + 600));
        IERC20(c1).forceApprove(_permit2, amt1);
        IPermit2(_permit2).approve(c1, _posMgr, uint160(amt1), uint48(block.timestamp + 600));

        uint128 liquidity = LiquidityAmounts.getLiquidityForAmounts(
            sqrtPriceX96, MIN_SQRT_RATIO, MAX_SQRT_RATIO, amt0, amt1
        );

        lpTokenId = IUniPositionManager(_posMgr).nextTokenId();

        UniPoolKey memory poolKey = _buildPoolKey(c0, c1);
        bytes memory actions = abi.encodePacked(uint8(0x02), uint8(0x0d));
        bytes[] memory params = new bytes[](2);
        params[0] = abi.encode(poolKey, MIN_TICK, MAX_TICK, liquidity, uint128(amt0), uint128(amt1), address(this), bytes(""));
        params[1] = abi.encode(c0, c1);
        IUniPositionManager(_posMgr).modifyLiquidities(abi.encode(actions, params), block.timestamp);
    }

    function _computePoolId(address c0, address c1) internal view override returns (bytes32) {
        return keccak256(abi.encode(_buildPoolKey(c0, c1)));
    }

    function _getCurrentSqrtPrice(address c0, address c1) internal view override returns (uint160) {
        bytes32 pid = keccak256(abi.encode(_buildPoolKey(c0, c1)));
        (uint160 sqrtPriceX96,,,) = IUniPoolManager(poolManager).getSlot0(pid);
        return sqrtPriceX96;
    }

    function _compoundFees(uint256 tokenId, address c0, address c1, uint160 sqrtPriceX96) internal override {
        address _permit2 = permit2;
        address _posMgr = positionManager;

        uint256 pre0 = IERC20(c0).balanceOf(address(this));
        uint256 pre1 = IERC20(c1).balanceOf(address(this));

        {
            bytes memory actions = abi.encodePacked(uint8(0x01), uint8(0x11));
            bytes[] memory params = new bytes[](2);
            params[0] = abi.encode(tokenId, uint256(0), uint128(0), uint128(0), bytes(""));
            params[1] = abi.encode(c0, c1, address(this));
            IUniPositionManager(_posMgr).modifyLiquidities(abi.encode(actions, params), block.timestamp);
        }

        uint256 bal0 = IERC20(c0).balanceOf(address(this)) - pre0;
        uint256 bal1 = IERC20(c1).balanceOf(address(this)) - pre1;
        if (bal0 == 0 && bal1 == 0) return;

        if (bal0 > 0) {
            IERC20(c0).forceApprove(_permit2, bal0);
            IPermit2(_permit2).approve(c0, _posMgr, uint160(bal0), uint48(block.timestamp + 600));
        }
        if (bal1 > 0) {
            IERC20(c1).forceApprove(_permit2, bal1);
            IPermit2(_permit2).approve(c1, _posMgr, uint160(bal1), uint48(block.timestamp + 600));
        }

        uint128 liquidity = LiquidityAmounts.getLiquidityForAmounts(
            sqrtPriceX96, MIN_SQRT_RATIO, MAX_SQRT_RATIO, bal0, bal1
        );
        if (liquidity == 0) return;

        bytes memory actions = abi.encodePacked(uint8(0x00), uint8(0x0d));
        bytes[] memory params = new bytes[](2);
        params[0] = abi.encode(tokenId, liquidity, uint128(bal0), uint128(bal1), bytes(""));
        params[1] = abi.encode(c0, c1);
        IUniPositionManager(_posMgr).modifyLiquidities(abi.encode(actions, params), block.timestamp);
    }

    function _buildPoolKey(address c0, address c1) internal pure returns (UniPoolKey memory) {
        return UniPoolKey({ currency0: c0, currency1: c1, fee: POOL_FEE, tickSpacing: TICK_SPACING, hooks: address(0) });
    }
}
