// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {FullMath} from "infinity-core/src/pool-cl/libraries/FullMath.sol";
import {FixedPoint96} from "infinity-core/src/pool-cl/libraries/FixedPoint96.sol";

/// @dev Permit2 interface for token approvals
interface IPermit2 {
    function approve(address token, address spender, uint160 amount, uint48 expiration) external;
}

/// @title LPManagerBase — Shared LP management logic for PancakeSwap V4 and Uniswap V4
/// @notice Only AWPRegistry may call; LP is permanently locked and cannot be withdrawn
/// @dev Subclasses implement DEX-specific pool initialization, minting, and pool ID computation
abstract contract LPManagerBase {
    using SafeERC20 for IERC20;

    // ══════════════════════════════════════════════
    //  Immutable storage
    // ══════════════════════════════════════════════

    /// @notice AWPRegistry contract address
    address public immutable awpRegistry;
    /// @notice Permit2 contract address
    address public immutable permit2;
    /// @notice AWP token contract address
    IERC20 public immutable awpToken;

    // ══════════════════════════════════════════════
    //  Constants
    // ══════════════════════════════════════════════

    /// @dev Pool fee: 1% (10000 = 1%)
    uint24 public constant POOL_FEE = 10000;
    /// @dev Tick spacing
    int24 public constant TICK_SPACING = 200;
    /// @dev Full-range minimum tick (rounded to TICK_SPACING)
    int24 public constant MIN_TICK = -887200;
    /// @dev Full-range maximum tick (rounded to TICK_SPACING)
    int24 public constant MAX_TICK = 887200;
    /// @dev Precomputed TickMath.getSqrtRatioAtTick(MIN_TICK)
    uint160 public constant MIN_SQRT_RATIO = 4295128739;
    /// @dev Precomputed TickMath.getSqrtRatioAtTick(MAX_TICK)
    uint160 public constant MAX_SQRT_RATIO = 1461446703485210103287273052203988822378723970342;

    // ══════════════════════════════════════════════
    //  Mappings
    // ══════════════════════════════════════════════

    /// @notice Alpha token address → pool ID
    mapping(address => bytes32) public alphaTokenToPoolId;
    /// @notice Alpha token address → LP NFT tokenId
    mapping(address => uint256) public alphaTokenToTokenId;

    // ══════════════════════════════════════════════
    //  Errors
    // ══════════════════════════════════════════════

    error NotAWPRegistry();
    /// @dev A LP pool already exists for this Alpha token
    error PoolAlreadyExists();
    /// @dev Amount exceeds Permit2 uint160 limit
    error AmountExceedsPermit2Limit();

    /// @dev Only the AWPRegistry contract may call
    modifier onlyAWPRegistry() {
        if (msg.sender != awpRegistry) revert NotAWPRegistry();
        _;
    }

    /// @notice Constructor
    /// @param awpRegistry_ AWPRegistry contract address
    /// @param permit2_ Permit2 contract address
    /// @param awpToken_ AWP token contract address
    constructor(address awpRegistry_, address permit2_, address awpToken_) {
        awpRegistry = awpRegistry_;
        permit2 = permit2_;
        awpToken = IERC20(awpToken_);
    }

    /// @notice Create an LP pool and add full-range liquidity (called once during worknet registration)
    /// @dev Full flow:
    ///   1. Check pool doesn't exist
    ///   2. Sort tokens (V4 requires currency0 < currency1)
    ///   3. Compute sqrtPriceX96
    ///   4. Call _initializePool (DEX-specific)
    ///   5. Call _approveAndMint (DEX-specific)
    ///   6. Compute poolId via _computePoolId (DEX-specific)
    ///   7. Store mappings
    /// @param alphaToken Alpha token address
    /// @param awpAmount AWP amount
    /// @param alphaAmount Alpha amount
    /// @return poolId Pool ID (bytes32)
    /// @return lpTokenId LP NFT ID
    function createPoolAndAddLiquidity(address alphaToken, uint256 awpAmount, uint256 alphaAmount)
        external
        onlyAWPRegistry
        returns (bytes32 poolId, uint256 lpTokenId)
    {
        // Each Alpha token can only have one LP pool
        if (alphaTokenToPoolId[alphaToken] != bytes32(0)) revert PoolAlreadyExists();

        // Sort tokens: V4 requires currency0 < currency1
        address awp = address(awpToken);
        (address c0, address c1) = awp < alphaToken ? (awp, alphaToken) : (alphaToken, awp);
        (uint256 amt0, uint256 amt1) = awp < alphaToken ? (awpAmount, alphaAmount) : (alphaAmount, awpAmount);

        // Compute initial price: sqrtPriceX96 = sqrt(amt1 * 2^192 / amt0)
        uint256 ratioX192 = FullMath.mulDiv(amt1, FixedPoint96.Q96 * FixedPoint96.Q96, amt0);
        uint160 sqrtPriceX96 = uint160(Math.sqrt(ratioX192));

        // Permit2 uses uint160 amounts; verify no truncation
        if (amt0 > type(uint160).max || amt1 > type(uint160).max) revert AmountExceedsPermit2Limit();

        // DEX-specific: initialize pool
        _initializePool(c0, c1, sqrtPriceX96);

        // DEX-specific: approve + mint LP position
        lpTokenId = _approveAndMint(c0, c1, amt0, amt1, sqrtPriceX96);

        // DEX-specific: compute pool ID
        poolId = _computePoolId(c0, c1);

        // Store mappings
        alphaTokenToPoolId[alphaToken] = poolId;
        alphaTokenToTokenId[alphaToken] = lpTokenId;
    }

    // ══════════════════════════════════════════════
    //  Fee Compounding (auto-compound LP fees back into liquidity)
    // ══════════════════════════════════════════════

    error NoPool();

    event FeesCompounded(address indexed alphaToken, uint256 tokenId);

    /// @notice Compound accumulated LP fees back into liquidity for a given Alpha token's pool.
    ///         Anyone can call — no access restriction needed since fees belong to the locked LP position.
    /// @param alphaToken The Alpha token whose LP position fees should be compounded
    function compoundFees(address alphaToken) external {
        uint256 tokenId = alphaTokenToTokenId[alphaToken];
        if (tokenId == 0) revert NoPool();

        address awp = address(awpToken);
        (address c0, address c1) = awp < alphaToken ? (awp, alphaToken) : (alphaToken, awp);

        // Read current pool price for accurate liquidity computation
        uint160 sqrtPriceX96 = _getCurrentSqrtPrice(c0, c1);

        _compoundFees(tokenId, c0, c1, sqrtPriceX96);

        emit FeesCompounded(alphaToken, tokenId);
    }

    /// @notice Check whether a pool has fees worth compounding (keeper query)
    /// @param alphaToken The Alpha token address
    /// @return hasPool Whether a pool exists
    /// @return tokenId The LP NFT token ID (0 if no pool)
    function needsCompounding(address alphaToken) external view returns (bool hasPool, uint256 tokenId) {
        tokenId = alphaTokenToTokenId[alphaToken];
        hasPool = tokenId != 0;
    }

    /// @dev DEX-specific: read current pool sqrtPriceX96
    function _getCurrentSqrtPrice(address c0, address c1) internal virtual view returns (uint160);

    /// @dev DEX-specific: compound fees back into liquidity for a position
    function _compoundFees(uint256 tokenId, address c0, address c1, uint160 sqrtPriceX96) internal virtual;

    /// @dev DEX-specific: initialize the pool
    function _initializePool(address c0, address c1, uint160 sqrtPriceX96) internal virtual;

    /// @dev DEX-specific: Permit2 approval + LP position minting
    function _approveAndMint(
        address c0, address c1,
        uint256 amt0, uint256 amt1, uint160 sqrtPriceX96
    ) internal virtual returns (uint256 lpTokenId);

    /// @dev DEX-specific: compute the pool ID
    function _computePoolId(address c0, address c1) internal virtual view returns (bytes32);
}
