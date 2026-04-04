// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {FullMath} from "infinity-core/src/pool-cl/libraries/FullMath.sol";
import {FixedPoint96} from "infinity-core/src/pool-cl/libraries/FixedPoint96.sol";

/// @dev Permit2 interface for token approvals
interface IPermit2 {
    function approve(address token, address spender, uint160 amount, uint48 expiration) external;
}

/// @dev Minimal interface to read guardian from AWPRegistry
interface IRegistryGuardian {
    function guardian() external view returns (address);
}

/// @title LPManagerBase — UUPS upgradeable LP management for PancakeSwap V4 and Uniswap V4
/// @notice Only AWPRegistry may call pool creation; LP is permanently locked and cannot be withdrawn.
/// @dev Protocol addresses (awpRegistry, awpToken) are hardcoded constants (same on all chains).
///      DEX-specific addresses (poolManager, positionManager, permit2) are immutable in each impl.
///      Proxy address is identical on all chains; Guardian upgrades proxy to chain-specific impl.
abstract contract LPManagerBase is Initializable, UUPSUpgradeable {
    using SafeERC20 for IERC20;

    // ══════════════════════════════════════════════
    //  Protocol addresses (hardcoded — same on all chains)
    // ══════════════════════════════════════════════

    address public constant awpRegistry = 0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A;
    address public constant awpToken = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;

    // ══════════════════════════════════════════════
    //  DEX addresses (immutable per impl — set in subclass constructor)
    // ══════════════════════════════════════════════

    /// @notice Permit2 contract address (differs on BSC)
    address public immutable permit2;

    // ══════════════════════════════════════════════
    //  Constants
    // ══════════════════════════════════════════════

    uint24 public constant POOL_FEE = 10000;
    int24 public constant TICK_SPACING = 200;
    int24 public constant MIN_TICK = -887200;
    int24 public constant MAX_TICK = 887200;
    uint160 public constant MIN_SQRT_RATIO = 4295128739;
    uint160 public constant MAX_SQRT_RATIO = 1461446703485210103287273052203988822378723970342;

    // ══════════════════════════════════════════════
    //  Storage (proxy state — survives impl upgrades)
    // ══════════════════════════════════════════════

    mapping(address => bytes32) public worknetTokenToPoolId;
    mapping(address => uint256) public worknetTokenToTokenId;

    /// @dev Reserved storage gap for upgrades
    uint256[48] private __gap;

    // ══════════════════════════════════════════════
    //  Errors & Events
    // ══════════════════════════════════════════════

    error NotAWPRegistry();
    error NotGuardian();
    error PoolAlreadyExists();
    error AmountExceedsPermit2Limit();
    error NoPool();

    event FeesCompounded(address indexed worknetToken, uint256 tokenId);

    // ══════════════════════════════════════════════
    //  Modifiers
    // ══════════════════════════════════════════════

    modifier onlyAWPRegistry() {
        if (msg.sender != awpRegistry) revert NotAWPRegistry();
        _;
    }

    // ══════════════════════════════════════════════
    //  Constructor + Initialization
    // ══════════════════════════════════════════════

    /// @param permit2_ Permit2 address (immutable per impl, differs on BSC)
    constructor(address permit2_) {
        permit2 = permit2_;
        _disableInitializers();
    }

    /// @notice Initialize (called once via proxy, no params — protocol addresses hardcoded)
    function initialize() external initializer {
        __UUPSUpgradeable_init();
    }

    /// @dev UUPS upgrade authorization — reads guardian from AWPRegistry
    function _authorizeUpgrade(address) internal view override {
        if (msg.sender != IRegistryGuardian(awpRegistry).guardian()) revert NotGuardian();
    }

    // ══════════════════════════════════════════════
    //  Pool Creation
    // ══════════════════════════════════════════════

    function createPoolAndAddLiquidity(address worknetToken, uint256 awpAmount, uint256 alphaAmount)
        external
        onlyAWPRegistry
        returns (bytes32 poolId, uint256 lpTokenId)
    {
        if (worknetTokenToPoolId[worknetToken] != bytes32(0)) revert PoolAlreadyExists();

        address awp = awpToken;
        (address c0, address c1) = awp < worknetToken ? (awp, worknetToken) : (worknetToken, awp);
        (uint256 amt0, uint256 amt1) = awp < worknetToken ? (awpAmount, alphaAmount) : (alphaAmount, awpAmount);

        uint256 ratioX192 = FullMath.mulDiv(amt1, FixedPoint96.Q96 * FixedPoint96.Q96, amt0);
        uint160 sqrtPriceX96 = uint160(Math.sqrt(ratioX192));

        if (amt0 > type(uint160).max || amt1 > type(uint160).max) revert AmountExceedsPermit2Limit();

        _initializePool(c0, c1, sqrtPriceX96);
        lpTokenId = _approveAndMint(c0, c1, amt0, amt1, sqrtPriceX96);
        poolId = _computePoolId(c0, c1);

        worknetTokenToPoolId[worknetToken] = poolId;
        worknetTokenToTokenId[worknetToken] = lpTokenId;
    }

    // ══════════════════════════════════════════════
    //  Fee Compounding
    // ══════════════════════════════════════════════

    function compoundFees(address worknetToken) external {
        uint256 tokenId = worknetTokenToTokenId[worknetToken];
        if (tokenId == 0) revert NoPool();

        address awp = awpToken;
        (address c0, address c1) = awp < worknetToken ? (awp, worknetToken) : (worknetToken, awp);

        uint160 sqrtPriceX96 = _getCurrentSqrtPrice(c0, c1);
        _compoundFees(tokenId, c0, c1, sqrtPriceX96);

        emit FeesCompounded(worknetToken, tokenId);
    }

    function needsCompounding(address worknetToken) external view returns (bool hasPool, uint256 tokenId) {
        tokenId = worknetTokenToTokenId[worknetToken];
        hasPool = tokenId != 0;
    }

    // ══════════════════════════════════════════════
    //  DEX-specific virtual functions
    // ══════════════════════════════════════════════

    function _getCurrentSqrtPrice(address c0, address c1) internal virtual view returns (uint160);
    function _compoundFees(uint256 tokenId, address c0, address c1, uint160 sqrtPriceX96) internal virtual;
    function _initializePool(address c0, address c1, uint160 sqrtPriceX96) internal virtual;
    function _approveAndMint(address c0, address c1, uint256 amt0, uint256 amt1, uint160 sqrtPriceX96) internal virtual returns (uint256 lpTokenId);
    function _computePoolId(address c0, address c1) internal virtual view returns (bytes32);
}
