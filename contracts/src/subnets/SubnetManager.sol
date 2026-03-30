// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {MerkleProof} from "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {IERC1363Receiver} from "../interfaces/IERC1363Receiver.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {TickMath} from "infinity-core/src/pool-cl/libraries/TickMath.sol";
import {FullMath} from "infinity-core/src/pool-cl/libraries/FullMath.sol";

interface IAlphaToken {
    function mint(address to, uint256 amount) external;
    function burn(uint256 amount) external;
    function balanceOf(address account) external view returns (uint256);
}

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

/// @title SubnetManager — Reference subnet contract (proxy-compatible)
/// @notice Deployed behind ERC1967Proxy by AWPRegistry when subnetManager is not provided.
///         Implementation is shared; each proxy gets its own storage via initialize().
/// @dev Three roles via OZ AccessControl:
///   - MERKLE_ROLE:   Submit Merkle roots → claim mints Alpha to users
///   - STRATEGY_ROLE: Choose AWP handling strategy + execute
///   - TRANSFER_ROLE: Transfer any token held by this contract
contract SubnetManager is Initializable, AccessControlUpgradeable, ReentrancyGuardUpgradeable, IERC1363Receiver {
    using SafeERC20 for IERC20;

    bytes32 public constant MERKLE_ROLE = keccak256("MERKLE_ROLE");
    bytes32 public constant STRATEGY_ROLE = keccak256("STRATEGY_ROLE");
    bytes32 public constant TRANSFER_ROLE = keccak256("TRANSFER_ROLE");

    enum AWPStrategy { Reserve, AddLiquidity, BuybackBurn }

    // ── DEX addresses (set via initialize, chain-agnostic) ──
    address public clPoolManager;
    address public clPositionManager;
    address public clSwapRouter;
    address public permit2;
    uint24 public poolFee;
    int24 public tickSpacing;

    // ── Action codes ──
    uint8 internal constant ACT_CL_MINT_POSITION = 0x02;
    uint8 internal constant ACT_CL_SWAP_EXACT_IN_SINGLE = 0x06;
    uint8 internal constant ACT_SETTLE_ALL = 0x0c;
    uint8 internal constant ACT_SETTLE_PAIR = 0x0d;
    uint8 internal constant ACT_TAKE_ALL = 0x0f;

    // ── Storage (set via initialize) ──
    IAlphaToken public alphaToken;
    IERC20 public awpToken;
    bytes32 public poolId;
    PoolKey public poolKey;
    AWPStrategy public currentStrategy;

    mapping(uint32 => bytes32) public merkleRoots;
    mapping(uint32 => mapping(address => bool)) public claimed;

    /// @dev Reserved storage gap for future upgrades
    uint256[38] private __gap;

    event MerkleRootSet(uint32 indexed epoch, bytes32 merkleRoot);
    event Claimed(uint32 indexed epoch, address indexed account, uint256 amount);
    event StrategyUpdated(AWPStrategy indexed strategy);
    event AWPProcessed(AWPStrategy indexed strategy, uint256 amount);
    event LiquidityAdded(uint256 tokenId, uint256 awpAmount);
    event BuybackBurned(uint256 awpSpent, uint256 alphaBurned);
    event TokenTransferred(address indexed token, address indexed to, uint256 amount);

    error AlreadyClaimed();
    error InvalidProof();
    error RootAlreadySet();
    error NoRootForEpoch();
    error ZeroAmount();

    /// @dev Prevent direct construction; must use proxy
    constructor() {
        _disableInitializers();
    }

    /// @notice Initialize (called once via proxy constructor)
    /// @param alphaToken_ Alpha token address
    /// @param awpToken_ AWP token address
    /// @param poolId_ LP pool ID (bytes32)
    /// @param admin_ Admin address (receives DEFAULT_ADMIN_ROLE)
    /// @param dexConfig_ ABI-encoded DEX configuration:
    ///        (address clPoolManager, address clPositionManager, address clSwapRouter, address permit2, uint24 poolFee, int24 tickSpacing)
    function initialize(
        address alphaToken_, address awpToken_, bytes32 poolId_, address admin_,
        bytes calldata dexConfig_
    ) external virtual initializer {
        // Decode DEX addresses and pool parameters (chain-agnostic)
        (
            address clPoolManager_,
            address clPositionManager_,
            address clSwapRouter_,
            address permit2_,
            uint24 poolFee_,
            int24 tickSpacing_
        ) = abi.decode(dexConfig_, (address, address, address, address, uint24, int24));

        _initializeBase(
            alphaToken_, awpToken_, poolId_, admin_,
            clPoolManager_, clPositionManager_, clSwapRouter_, permit2_, poolFee_, tickSpacing_
        );

        // Construct PancakeSwap V4 PoolKey (6 fields)
        (address c0, address c1) = awpToken_ < alphaToken_
            ? (awpToken_, alphaToken_)
            : (alphaToken_, awpToken_);
        poolKey = PoolKey({
            currency0: c0,
            currency1: c1,
            hooks: address(0),
            poolManager: clPoolManager_,
            fee: poolFee_,
            parameters: bytes32(uint256(int256(tickSpacing_)) << 16)
        });
    }

    /// @dev 共享初始化逻辑：AccessControl、ReentrancyGuard、存储、角色授予
    function _initializeBase(
        address alphaToken_, address awpToken_, bytes32 poolId_, address admin_,
        address clPoolManager_, address clPositionManager_, address clSwapRouter_,
        address permit2_, uint24 poolFee_, int24 tickSpacing_
    ) internal {
        __AccessControl_init();
        __ReentrancyGuard_init();

        alphaToken = IAlphaToken(alphaToken_);
        awpToken = IERC20(awpToken_);
        poolId = poolId_;
        clPoolManager = clPoolManager_;
        clPositionManager = clPositionManager_;
        clSwapRouter = clSwapRouter_;
        permit2 = permit2_;
        poolFee = poolFee_;
        tickSpacing = tickSpacing_;

        _grantRole(DEFAULT_ADMIN_ROLE, admin_);
        _grantRole(MERKLE_ROLE, admin_);
        _grantRole(STRATEGY_ROLE, admin_);
        _grantRole(TRANSFER_ROLE, admin_);
    }

    // ═══════════════════════════════════════════════
    //  Merkle Distribution (MERKLE_ROLE)
    // ═══════════════════════════════════════════════

    function setMerkleRoot(uint32 epoch, bytes32 root) external onlyRole(MERKLE_ROLE) {
        if (merkleRoots[epoch] != bytes32(0)) revert RootAlreadySet();
        if (root == bytes32(0)) revert ZeroAmount();
        merkleRoots[epoch] = root;
        emit MerkleRootSet(epoch, root);
    }

    function claim(uint32 epoch, uint256 amount, bytes32[] calldata proof) external nonReentrant {
        if (merkleRoots[epoch] == bytes32(0)) revert NoRootForEpoch();
        if (claimed[epoch][msg.sender]) revert AlreadyClaimed();

        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(msg.sender, amount))));
        if (!MerkleProof.verify(proof, merkleRoots[epoch], leaf)) revert InvalidProof();

        claimed[epoch][msg.sender] = true;
        alphaToken.mint(msg.sender, amount);
        emit Claimed(epoch, msg.sender, amount);
    }

    function isClaimed(uint32 epoch, address account) external view returns (bool) {
        return claimed[epoch][account];
    }

    // ═══════════════════════════════════════════════
    //  AWP Strategy (STRATEGY_ROLE)
    // ═══════════════════════════════════════════════

    function setStrategy(AWPStrategy strategy) external onlyRole(STRATEGY_ROLE) {
        currentStrategy = strategy;
        emit StrategyUpdated(strategy);
    }

    function executeStrategy(uint256 amount) external nonReentrant onlyRole(STRATEGY_ROLE) {
        if (amount == 0) revert ZeroAmount();
        AWPStrategy strategy = currentStrategy;
        if (strategy == AWPStrategy.AddLiquidity) {
            _addSingleSidedLiquidity(amount);
        } else if (strategy == AWPStrategy.BuybackBurn) {
            _buybackAndBurn(amount);
        }
        // Reserve is a no-op (AWP stays in contract); skip event to avoid misleading indexers
        if (strategy != AWPStrategy.Reserve) {
            emit AWPProcessed(strategy, amount);
        }
    }

    // ═══════════════════════════════════════════════
    //  ERC1363 Receiver — auto-execute strategy on AWP transferAndCall
    // ═══════════════════════════════════════════════

    /// @notice Called by AWPToken.transferAndCall when AWP is sent to this contract
    /// @dev Automatically executes the current AWP strategy on the received amount.
    ///      Only responds to AWP token transfers; other tokens are accepted silently.
    function onTransferReceived(address, address, uint256 amount, bytes calldata)
        external
        override
        nonReentrant
        returns (bytes4)
    {
        if (msg.sender == address(awpToken) && amount > 0) {
            AWPStrategy strategy = currentStrategy;
            if (strategy == AWPStrategy.AddLiquidity) {
                _addSingleSidedLiquidity(amount);
                emit AWPProcessed(strategy, amount);
            } else if (strategy == AWPStrategy.BuybackBurn) {
                _buybackAndBurn(amount);
                emit AWPProcessed(strategy, amount);
            }
            // Reserve: AWP stays in contract, no action
        }
        return IERC1363Receiver.onTransferReceived.selector;
    }

    // ═══════════════════════════════════════════════
    //  Token Transfer (TRANSFER_ROLE)
    // ═══════════════════════════════════════════════

    function transferToken(address token, address to, uint256 amount) external onlyRole(TRANSFER_ROLE) {
        IERC20(token).safeTransfer(to, amount);
        emit TokenTransferred(token, to, amount);
    }

    // ═══════════════════════════════════════════════
    //  Internal: Pool Slot0 Read (virtual — overridden by SubnetManagerUni)
    // ═══════════════════════════════════════════════

    /// @dev 读取池子当前 sqrtPriceX96 和 tick；子合约可覆写以使用不同的数据源
    function _getSlot0() internal view virtual returns (uint160 sqrtPriceX96, int24 tick) {
        (sqrtPriceX96, tick,,) = ICLPoolManager(clPoolManager).getSlot0(poolId);
    }

    // ═══════════════════════════════════════════════
    //  Internal: PancakeSwap V4 — Add Single-Sided Liquidity
    // ═══════════════════════════════════════════════

    function _addSingleSidedLiquidity(uint256 amount) internal virtual {
        (, int24 currentTick) = _getSlot0();

        // Floor-align currentTick to tickSpacing
        int24 ts = tickSpacing;
        int24 aligned = (currentTick / ts) * ts;
        if (aligned > currentTick) aligned -= ts;

        // Compute min/max tick aligned to tickSpacing
        int24 minTick = (-887272 / ts) * ts;
        int24 maxTick = (887272 / ts) * ts;

        bool awpIs0 = address(awpToken) < address(alphaToken);

        int24 tickLower;
        int24 tickUpper;
        if (awpIs0) {
            // AWP is token0: single-sided deposit requires range ABOVE current price
            tickLower = aligned + ts;
            tickUpper = maxTick;
        } else {
            // AWP is token1: single-sided deposit requires range BELOW current price
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
        params[0] = abi.encode(
            poolKey, tickLower, tickUpper, liquidity,
            awpIs0 ? uint128(amount) : uint128(0),
            awpIs0 ? uint128(0) : uint128(amount),
            address(this), bytes("")
        );
        params[1] = abi.encode(poolKey.currency0, poolKey.currency1);

        ICLPositionManager(clPositionManager).modifyLiquidities(abi.encode(actions, params), block.timestamp);
        emit LiquidityAdded(tokenId, amount);
    }

    // ═══════════════════════════════════════════════
    //  Internal: PancakeSwap V4 — Buyback + Burn
    // ═══════════════════════════════════════════════

    function _buybackAndBurn(uint256 amount) internal virtual {
        IERC20(address(awpToken)).forceApprove(permit2, amount);
        IPermit2(permit2).approve(address(awpToken), clSwapRouter, uint160(amount), uint48(block.timestamp + 600));

        bool zeroForOne = address(awpToken) < address(alphaToken);

        // Read current pool price for slippage protection
        (uint160 sqrtPriceX96,) = _getSlot0();
        uint256 expectedOut;
        if (zeroForOne) {
            // Selling token0 for token1: expectedOut = amount * sqrtPrice^2 / 2^192
            // Split into two mulDiv to avoid sqrtPrice^2 overflow
            expectedOut = FullMath.mulDiv(FullMath.mulDiv(amount, sqrtPriceX96, 1 << 96), sqrtPriceX96, 1 << 96);
        } else {
            // Selling token1 for token0: expectedOut = amount * 2^192 / sqrtPrice^2
            // Split: (amount * 2^96 / sqrtPrice) * 2^96 / sqrtPrice
            expectedOut = FullMath.mulDiv(FullMath.mulDiv(amount, 1 << 96, sqrtPriceX96), 1 << 96, sqrtPriceX96);
        }
        uint128 minOut = uint128(expectedOut * 95 / 100); // 5% slippage tolerance

        bytes memory actions = abi.encodePacked(ACT_CL_SWAP_EXACT_IN_SINGLE, ACT_SETTLE_ALL, ACT_TAKE_ALL);
        bytes[] memory params = new bytes[](3);
        params[0] = abi.encode(poolKey, zeroForOne, uint128(amount), minOut, bytes(""));
        params[1] = abi.encode(address(awpToken), amount);
        params[2] = abi.encode(address(alphaToken), 0);

        uint256 before = alphaToken.balanceOf(address(this));
        ICLSwapRouter(clSwapRouter).executeActions(abi.encode(actions, params));
        uint256 received = alphaToken.balanceOf(address(this)) - before;

        if (received > 0) alphaToken.burn(received);
        emit BuybackBurned(amount, received);
    }

}
