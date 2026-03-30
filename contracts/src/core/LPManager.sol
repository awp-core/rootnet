// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {LPManagerBase, IPermit2} from "./LPManagerBase.sol";

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
    function getSlot0(bytes32 id) external view returns (uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee);
}

/// @dev PancakeSwap V4 CLPositionManager interface
interface ICLPositionManager {
    function modifyLiquidities(bytes calldata payload, uint256 deadline) external payable;
    function nextTokenId() external view returns (uint256);
}

/// @title LPManager — PancakeSwap V4 CL liquidity management
/// @notice Only AWPRegistry may call; LP is permanently locked in this contract and cannot be withdrawn
/// @dev Integrates PancakeSwap V4 Concentrated Liquidity:
///   - Full-range liquidity (MIN_TICK ~ MAX_TICK)
///   - Token approvals via Permit2
///   - Pool ID computed as keccak256(PoolKey)
contract LPManager is LPManagerBase {
    using SafeERC20 for IERC20;

    // ══════════════════════════════════════════════
    //  PancakeSwap-specific immutable storage
    // ══════════════════════════════════════════════

    /// @notice PancakeSwap V4 CLPoolManager address
    address public immutable clPoolManager;
    /// @notice PancakeSwap V4 CLPositionManager address
    address public immutable clPositionManager;

    /// @notice Constructor
    /// @param awpRegistry_ AWPRegistry contract address
    /// @param clPoolManager_ PancakeSwap V4 CLPoolManager address
    /// @param clPositionManager_ PancakeSwap V4 CLPositionManager address
    /// @param permit2_ Permit2 contract address
    /// @param awpToken_ AWP token contract address
    constructor(address awpRegistry_, address clPoolManager_, address clPositionManager_, address permit2_, address awpToken_)
        LPManagerBase(awpRegistry_, permit2_, awpToken_)
    {
        clPoolManager = clPoolManager_;
        clPositionManager = clPositionManager_;
    }

    /// @dev 初始化 PancakeSwap V4 CL pool
    function _initializePool(address c0, address c1, uint160 sqrtPriceX96) internal override {
        PoolKey memory poolKey = _buildPoolKey(c0, c1);
        ICLPoolManager(clPoolManager).initialize(poolKey, sqrtPriceX96);
    }

    /// @dev Permit2 approval + PancakeSwap V4 LP position minting
    function _approveAndMint(
        address c0, address c1,
        uint256 amt0, uint256 amt1, uint160 sqrtPriceX96
    ) internal override returns (uint256 lpTokenId) {
        // 通过 Permit2 approve token 到 CLPositionManager
        IERC20(c0).forceApprove(permit2, amt0);
        IPermit2(permit2).approve(c0, clPositionManager, uint160(amt0), uint48(block.timestamp + 600));
        IERC20(c1).forceApprove(permit2, amt1);
        IPermit2(permit2).approve(c1, clPositionManager, uint160(amt1), uint48(block.timestamp + 600));

        // 使用 PancakeSwap 官方 library 计算 full-range liquidity
        uint128 liquidity = LiquidityAmounts.getLiquidityForAmounts(
            sqrtPriceX96, MIN_SQRT_RATIO, MAX_SQRT_RATIO, amt0, amt1
        );

        // mint 前记录 nextTokenId
        lpTokenId = ICLPositionManager(clPositionManager).nextTokenId();

        // 编码 mint 操作: CL_MINT_POSITION(0x02) + SETTLE_PAIR(0x0d)
        PoolKey memory poolKey = _buildPoolKey(c0, c1);
        bytes memory actions = abi.encodePacked(uint8(0x02), uint8(0x0d));
        bytes[] memory params = new bytes[](2);
        params[0] = abi.encode(poolKey, MIN_TICK, MAX_TICK, liquidity, uint128(amt0), uint128(amt1), address(this), bytes(""));
        params[1] = abi.encode(c0, c1);
        ICLPositionManager(clPositionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);
    }

    /// @dev 计算 PancakeSwap V4 pool ID: keccak256(PoolKey)
    function _computePoolId(address c0, address c1) internal view override returns (bytes32) {
        PoolKey memory poolKey = _buildPoolKey(c0, c1);
        return keccak256(abi.encode(poolKey));
    }

    /// @dev Compound accumulated fees back into the LP position (PancakeSwap V4)
    ///      Step 1: DECREASE_LIQUIDITY(0x01) with liquidityDelta=0 → collects accrued fees to LPManager
    ///      Step 2: Re-add collected fees as liquidity via INCREASE_LIQUIDITY(0x00)
    function _getCurrentSqrtPrice(address c0, address c1) internal view override returns (uint160) {
        PoolKey memory poolKey = _buildPoolKey(c0, c1);
        bytes32 pid = keccak256(abi.encode(poolKey));
        (uint160 sqrtPriceX96,,,) = ICLPoolManager(clPoolManager).getSlot0(pid);
        return sqrtPriceX96;
    }

    function _compoundFees(uint256 tokenId, address c0, address c1, uint160 sqrtPriceX96) internal override {
        // Step 1: Collect fees — DECREASE_LIQUIDITY(0x01) with 0 delta + TAKE_PAIR(0x11)
        {
            bytes memory actions = abi.encodePacked(uint8(0x01), uint8(0x11));
            bytes[] memory params = new bytes[](2);
            params[0] = abi.encode(tokenId, uint256(0), uint128(0), uint128(0), bytes(""));
            params[1] = abi.encode(c0, c1);
            ICLPositionManager(clPositionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);
        }

        // Read collected fee balances
        uint256 bal0 = IERC20(c0).balanceOf(address(this));
        uint256 bal1 = IERC20(c1).balanceOf(address(this));
        if (bal0 == 0 && bal1 == 0) return; // no fees to compound

        // Step 2: Re-add fees as liquidity — approve + INCREASE_LIQUIDITY(0x00) + SETTLE_PAIR(0x0d)
        if (bal0 > 0) {
            IERC20(c0).forceApprove(permit2, bal0);
            IPermit2(permit2).approve(c0, clPositionManager, uint160(bal0), uint48(block.timestamp + 600));
        }
        if (bal1 > 0) {
            IERC20(c1).forceApprove(permit2, bal1);
            IPermit2(permit2).approve(c1, clPositionManager, uint160(bal1), uint48(block.timestamp + 600));
        }

        uint128 liquidity = LiquidityAmounts.getLiquidityForAmounts(
            sqrtPriceX96, MIN_SQRT_RATIO, MAX_SQRT_RATIO, bal0, bal1
        );
        if (liquidity == 0) return;

        PoolKey memory poolKey = _buildPoolKey(c0, c1);
        bytes memory actions = abi.encodePacked(uint8(0x00), uint8(0x0d));
        bytes[] memory params = new bytes[](2);
        params[0] = abi.encode(poolKey, MIN_TICK, MAX_TICK, liquidity, uint128(bal0), uint128(bal1), address(this), bytes(""));
        params[1] = abi.encode(c0, c1);
        ICLPositionManager(clPositionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);
    }

    /// @dev 构建 PancakeSwap V4 PoolKey
    function _buildPoolKey(address c0, address c1) internal view returns (PoolKey memory) {
        return PoolKey({
            currency0: c0,
            currency1: c1,
            hooks: address(0),
            poolManager: clPoolManager,
            fee: POOL_FEE,
            parameters: bytes32(uint256(int256(TICK_SPACING)) << 16)
        });
    }
}
