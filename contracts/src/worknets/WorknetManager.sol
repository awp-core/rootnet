// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {MerkleProof} from "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {IERC1363Receiver} from "../interfaces/IERC1363Receiver.sol";
import {IAWPRegistry} from "../interfaces/IAWPRegistry.sol";
import {LiquidityAmounts} from "infinity-periphery/src/pool-cl/libraries/LiquidityAmounts.sol";
import {TickMath} from "infinity-core/src/pool-cl/libraries/TickMath.sol";
import {FullMath} from "infinity-core/src/pool-cl/libraries/FullMath.sol";

interface IWorknetToken {
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

/// @title WorknetManager — Reference worknet contract (proxy-compatible)
/// @dev Protocol addresses (awpRegistry, awpToken) are constants. DEX addresses are immutable per impl.
///      Per-worknet state (worknetToken, poolId) is storage set via initialize.
contract WorknetManager is Initializable, UUPSUpgradeable, AccessControlUpgradeable, ReentrancyGuardUpgradeable, IERC1363Receiver {
    using SafeERC20 for IERC20;

    bytes32 public constant MERKLE_ROLE = keccak256("MERKLE_ROLE");
    bytes32 public constant STRATEGY_ROLE = keccak256("STRATEGY_ROLE");
    bytes32 public constant TRANSFER_ROLE = keccak256("TRANSFER_ROLE");

    enum AWPStrategy { Reserve, AddLiquidity, BuybackBurn }

    // ── Protocol addresses (constant — same on all chains) ──
    address public constant awpRegistry = 0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A;
    address public constant awpToken = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;

    // ── DEX addresses (immutable — set in impl constructor, differ per chain) ──
    address public immutable clPoolManager;
    address public immutable clPositionManager;
    address public immutable clSwapRouter;
    address public immutable permit2;

    // ── Constants ──
    uint24 public constant POOL_FEE = 10000;
    int24 public constant TICK_SPACING = 200;

    uint8 internal constant ACT_CL_MINT_POSITION = 0x02;
    uint8 internal constant ACT_CL_SWAP_EXACT_IN_SINGLE = 0x06;
    uint8 internal constant ACT_SETTLE_ALL = 0x0c;
    uint8 internal constant ACT_SETTLE_PAIR = 0x0d;
    uint8 internal constant ACT_TAKE_ALL = 0x0f;

    // ── Per-worknet storage (set via initialize) ──
    IWorknetToken public worknetToken;
    bytes32 public poolId;
    PoolKey public poolKey;
    AWPStrategy public currentStrategy;

    mapping(uint32 => bytes32) public merkleRoots;
    mapping(uint32 => mapping(address => bool)) public claimed;

    uint256 public slippageBps;
    bool public strategyPaused;
    uint256 public minStrategyAmount;

    /// @dev Reserved storage gap for future upgrades
    uint256[35] private __gap;

    event MerkleRootSet(uint32 indexed epoch, bytes32 merkleRoot);
    event Claimed(uint32 indexed epoch, address indexed account, uint256 amount);
    event StrategyUpdated(AWPStrategy indexed strategy);
    event AWPProcessed(AWPStrategy indexed strategy, uint256 amount);
    event LiquidityAdded(uint256 tokenId, uint256 awpAmount);
    event BuybackBurned(uint256 awpSpent, uint256 alphaBurned);
    event TokenTransferred(address indexed token, address indexed to, uint256 amount);
    event SlippageUpdated(uint256 bps);
    event StrategyPausedChanged(bool paused);

    error StrategyIsPaused();
    error InvalidSlippage();
    error ArrayLengthMismatch();
    error AlreadyClaimed();
    error InvalidProof();
    error RootAlreadySet();
    error NoRootForEpoch();
    error ZeroAmount();
    error ZeroRoot();

    function _authorizeUpgrade(address) internal override onlyRole(DEFAULT_ADMIN_ROLE) {}

    /// @param permit2_ Permit2 address
    /// @param clPoolManager_ DEX PoolManager
    /// @param clPositionManager_ DEX PositionManager
    /// @param clSwapRouter_ DEX SwapRouter
    constructor(address permit2_, address clPoolManager_, address clPositionManager_, address clSwapRouter_) {
        permit2 = permit2_;
        clPoolManager = clPoolManager_;
        clPositionManager = clPositionManager_;
        clSwapRouter = clSwapRouter_;
        _disableInitializers();
    }

    /// @notice Initialize per-worknet state (called once via proxy constructor)
    function initialize(address worknetToken_, bytes32 poolId_, address admin_) external virtual initializer {
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

        // Construct PancakeSwap V4 PoolKey
        (address c0, address c1) = awpToken < worknetToken_
            ? (awpToken, worknetToken_)
            : (worknetToken_, awpToken);
        poolKey = PoolKey({
            currency0: c0, currency1: c1, hooks: address(0), poolManager: clPoolManager,
            fee: POOL_FEE, parameters: bytes32(uint256(int256(TICK_SPACING)) << 16)
        });
    }

    // ═══════════════════════════════════════════════
    //  Merkle Distribution (MERKLE_ROLE)
    // ═══════════════════════════════════════════════

    function setMerkleRoot(uint32 epoch, bytes32 root) external onlyRole(MERKLE_ROLE) {
        if (merkleRoots[epoch] != bytes32(0)) revert RootAlreadySet();
        if (root == bytes32(0)) revert ZeroRoot();
        merkleRoots[epoch] = root;
        emit MerkleRootSet(epoch, root);
    }

    function claim(uint32 epoch, uint256 amount, bytes32[] calldata proof) external nonReentrant {
        if (merkleRoots[epoch] == bytes32(0)) revert NoRootForEpoch();
        if (claimed[epoch][msg.sender]) revert AlreadyClaimed();

        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(msg.sender, amount))));
        if (!MerkleProof.verify(proof, merkleRoots[epoch], leaf)) revert InvalidProof();

        claimed[epoch][msg.sender] = true;

        address to = IAWPRegistry(awpRegistry).resolveRecipient(msg.sender);
        worknetToken.mint(to, amount);
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
        if (strategyPaused) revert StrategyIsPaused();
        if (amount == 0) revert ZeroAmount();
        AWPStrategy strategy = currentStrategy;
        if (strategy == AWPStrategy.AddLiquidity) {
            _addSingleSidedLiquidity(amount);
        } else if (strategy == AWPStrategy.BuybackBurn) {
            _buybackAndBurn(amount);
        }
        if (strategy != AWPStrategy.Reserve) {
            emit AWPProcessed(strategy, amount);
        }
    }

    // ═══════════════════════════════════════════════
    //  ERC1363 Receiver
    // ═══════════════════════════════════════════════

    function onTransferReceived(address, address, uint256 amount, bytes calldata)
        external override nonReentrant returns (bytes4)
    {
        if (msg.sender == awpToken && amount > 0 && !strategyPaused && amount >= minStrategyAmount) {
            AWPStrategy strategy = currentStrategy;
            if (strategy == AWPStrategy.AddLiquidity) {
                _addSingleSidedLiquidity(amount);
                emit AWPProcessed(strategy, amount);
            } else if (strategy == AWPStrategy.BuybackBurn) {
                _buybackAndBurn(amount);
                emit AWPProcessed(strategy, amount);
            }
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

    function batchTransferToken(address token, address[] calldata recipients, uint256[] calldata amounts)
        external onlyRole(TRANSFER_ROLE)
    {
        if (recipients.length != amounts.length) revert ArrayLengthMismatch();
        for (uint256 i = 0; i < recipients.length;) {
            IERC20(token).safeTransfer(recipients[i], amounts[i]);
            emit TokenTransferred(token, recipients[i], amounts[i]);
            unchecked { ++i; }
        }
    }

    // ═══════════════════════════════════════════════
    //  Configuration
    // ═══════════════════════════════════════════════

    function setSlippageTolerance(uint256 bps) external onlyRole(STRATEGY_ROLE) {
        if (bps == 0 || bps > 5000) revert InvalidSlippage();
        slippageBps = bps;
        emit SlippageUpdated(bps);
    }

    function setStrategyPaused(bool paused) external onlyRole(DEFAULT_ADMIN_ROLE) {
        strategyPaused = paused;
        emit StrategyPausedChanged(paused);
    }

    function setMinStrategyAmount(uint256 amount) external onlyRole(STRATEGY_ROLE) {
        minStrategyAmount = amount;
    }

    // ═══════════════════════════════════════════════
    //  Internal: Pool Slot0 Read
    // ═══════════════════════════════════════════════

    function _getSlot0() internal view virtual returns (uint160 sqrtPriceX96, int24 tick) {
        (sqrtPriceX96, tick,,) = ICLPoolManager(clPoolManager).getSlot0(poolId);
    }

    // ═══════════════════════════════════════════════
    //  Internal: Add Single-Sided Liquidity
    // ═══════════════════════════════════════════════

    function _addSingleSidedLiquidity(uint256 amount) internal virtual {
        PoolKey memory pk = poolKey;
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

    // ═══════════════════════════════════════════════
    //  Internal: Buyback + Burn
    // ═══════════════════════════════════════════════

    function _buybackAndBurn(uint256 amount) internal virtual {
        PoolKey memory pk = poolKey;
        IERC20(awpToken).forceApprove(permit2, amount);
        IPermit2(permit2).approve(awpToken, clSwapRouter, uint160(amount), uint48(block.timestamp + 600));

        bool zeroForOne = awpToken < address(worknetToken);

        (uint160 sqrtPriceX96,) = _getSlot0();
        uint256 expectedOut;
        if (zeroForOne) {
            expectedOut = FullMath.mulDiv(FullMath.mulDiv(amount, sqrtPriceX96, 1 << 96), sqrtPriceX96, 1 << 96);
        } else {
            expectedOut = FullMath.mulDiv(FullMath.mulDiv(amount, 1 << 96, sqrtPriceX96), 1 << 96, sqrtPriceX96);
        }
        uint128 minOut = uint128(expectedOut * (10000 - slippageBps) / 10000);

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
}
