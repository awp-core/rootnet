// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {LPManagerBase, IPermit2} from "./LPManagerBase.sol";

struct PoolKey {
    address currency0;
    address currency1;
    address hooks;
    address poolManager;
    uint24 fee;
    bytes32 parameters;
}

interface ICLPoolManager {
    function initialize(PoolKey calldata key, uint160 sqrtPriceX96) external returns (int24 tick);
    function getSlot0(bytes32 id) external view returns (uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee);
}

interface ICLPositionManager {
    function modifyLiquidities(bytes calldata payload, uint256 deadline) external payable;
    function nextTokenId() external view returns (uint256);
}

/// @title LPManager — UUPS upgradeable PancakeSwap V4 CL liquidity management (BSC)
/// @dev DEX addresses are immutable in the impl bytecode. Proxy upgrades to this impl on BSC.
contract LPManager is LPManagerBase {
    using SafeERC20 for IERC20;

    address public immutable clPoolManager;
    address public immutable clPositionManager;

    /// @param permit2_ Permit2 address (BSC)
    /// @param clPoolManager_ PancakeSwap V4 CLPoolManager
    /// @param clPositionManager_ PancakeSwap V4 CLPositionManager
    constructor(address permit2_, address clPoolManager_, address clPositionManager_)
        LPManagerBase(permit2_)
    {
        clPoolManager = clPoolManager_;
        clPositionManager = clPositionManager_;
    }

    function _initializePool(address c0, address c1, uint160 sqrtPriceX96) internal override {
        ICLPoolManager(clPoolManager).initialize(_buildPoolKey(c0, c1), sqrtPriceX96);
    }

    function _approveAndMint(
        address c0, address c1,
        uint256 amt0, uint256 amt1, uint160 sqrtPriceX96
    ) internal override returns (uint256 lpTokenId) {
        address _permit2 = permit2;
        address _posMgr = clPositionManager;

        IERC20(c0).forceApprove(_permit2, amt0);
        IPermit2(_permit2).approve(c0, _posMgr, uint160(amt0), uint48(block.timestamp + 600));
        IERC20(c1).forceApprove(_permit2, amt1);
        IPermit2(_permit2).approve(c1, _posMgr, uint160(amt1), uint48(block.timestamp + 600));

        uint128 liquidity = LiquidityAmounts.getLiquidityForAmounts(
            sqrtPriceX96, MIN_SQRT_RATIO, MAX_SQRT_RATIO, amt0, amt1
        );

        lpTokenId = ICLPositionManager(_posMgr).nextTokenId();

        PoolKey memory poolKey = _buildPoolKey(c0, c1);
        bytes memory actions = abi.encodePacked(uint8(0x02), uint8(0x0d));
        bytes[] memory params = new bytes[](2);
        params[0] = abi.encode(poolKey, MIN_TICK, MAX_TICK, liquidity, uint128(amt0), uint128(amt1), address(this), bytes(""));
        params[1] = abi.encode(c0, c1);
        ICLPositionManager(_posMgr).modifyLiquidities(abi.encode(actions, params), block.timestamp);
    }

    function _computePoolId(address c0, address c1) internal view override returns (bytes32) {
        return keccak256(abi.encode(_buildPoolKey(c0, c1)));
    }

    function _getCurrentSqrtPrice(address c0, address c1) internal view override returns (uint160) {
        bytes32 pid = keccak256(abi.encode(_buildPoolKey(c0, c1)));
        (uint160 sqrtPriceX96,,,) = ICLPoolManager(clPoolManager).getSlot0(pid);
        return sqrtPriceX96;
    }

    function _compoundFees(uint256 tokenId, address c0, address c1, uint160 sqrtPriceX96) internal override {
        address _permit2 = permit2;
        address _posMgr = clPositionManager;

        uint256 pre0 = IERC20(c0).balanceOf(address(this));
        uint256 pre1 = IERC20(c1).balanceOf(address(this));

        {
            bytes memory actions = abi.encodePacked(uint8(0x01), uint8(0x11));
            bytes[] memory params = new bytes[](2);
            params[0] = abi.encode(tokenId, uint256(0), uint128(0), uint128(0), bytes(""));
            params[1] = abi.encode(c0, c1, address(this));
            ICLPositionManager(_posMgr).modifyLiquidities(abi.encode(actions, params), block.timestamp);
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
        ICLPositionManager(_posMgr).modifyLiquidities(abi.encode(actions, params), block.timestamp);
    }

    function _buildPoolKey(address c0, address c1) internal view returns (PoolKey memory) {
        return PoolKey({
            currency0: c0, currency1: c1, hooks: address(0), poolManager: clPoolManager,
            fee: POOL_FEE, parameters: bytes32(uint256(int256(TICK_SPACING)) << 16)
        });
    }
}
