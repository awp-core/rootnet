// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {MerkleProof} from "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuardTransient} from "@openzeppelin/contracts/utils/ReentrancyGuardTransient.sol";
import {IERC1363Receiver} from "@openzeppelin/contracts/interfaces/IERC1363Receiver.sol";
import {IAWPRegistry} from "../interfaces/IAWPRegistry.sol";
import {IWorknetToken} from "../interfaces/IWorknetToken.sol";

/// @title WorknetManagerBase — Abstract base for worknet management (DEX-agnostic)
/// @dev Merkle distribution, AWP strategy dispatch, token transfers, access control.
///      Subclasses implement DEX-specific liquidity and swap operations.
///      Protocol addresses (awpRegistry, awpToken) are constants (same on all chains).
abstract contract WorknetManagerBase is Initializable, UUPSUpgradeable, AccessControlUpgradeable, ReentrancyGuardTransient, IERC1363Receiver {
    using SafeERC20 for IERC20;

    bytes32 public constant MERKLE_ROLE = keccak256("MERKLE_ROLE");
    bytes32 public constant STRATEGY_ROLE = keccak256("STRATEGY_ROLE");
    bytes32 public constant TRANSFER_ROLE = keccak256("TRANSFER_ROLE");

    enum AWPStrategy { Reserve, AddLiquidity, BuybackBurn }

    // ── Protocol addresses (constant — same on all chains) ──
    address public constant awpRegistry = 0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A;
    address public constant awpToken = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;

    // ── Per-worknet storage (packed: slot 0 = 24 bytes) ──  [O5]
    IWorknetToken public worknetToken;       // 20 bytes
    AWPStrategy public currentStrategy;      // 1 byte
    bool public strategyPaused;              // 1 byte
    uint16 public slippageBps;               // 2 bytes  (max 5000, fits uint16)

    bytes32 public poolId;                   // slot 1
    uint256 public minStrategyAmount;        // slot 2

    mapping(uint32 => bytes32) public merkleRoots;
    mapping(uint32 => mapping(address => bool)) public claimed;

    // ── LP position cache (avoids minting new NFT every call) ── [B1]
    uint256 internal _lastLpTokenId;
    int24 internal _lastTickLower;
    int24 internal _lastTickUpper;

    /// @dev Reserved storage gap for future base upgrades
    uint256[33] private __baseGap;

    // ── Events ──
    event MerkleRootSet(uint32 indexed epoch, bytes32 merkleRoot);
    event Claimed(uint32 indexed epoch, address indexed account, uint256 amount);
    event StrategyUpdated(AWPStrategy indexed strategy);
    event AWPProcessed(AWPStrategy indexed strategy, uint256 amount);
    event LiquidityAdded(uint256 tokenId, uint256 awpAmount);
    event LiquidityIncreased(uint256 tokenId, uint256 awpAmount);
    event BuybackBurned(uint256 awpSpent, uint256 alphaBurned);
    event TokenTransferred(address indexed token, address indexed to, uint256 amount);
    event SlippageUpdated(uint16 bps);
    event StrategyPausedChanged(bool paused);

    // ── Errors ──
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

    /// @dev Shared initialization (called by subclass initialize)
    function __WorknetManagerBase_init(address worknetToken_, bytes32 poolId_, address admin_) internal onlyInitializing {
        __AccessControl_init();
        // ReentrancyGuardTransient uses TSTORE — no init needed

        worknetToken = IWorknetToken(worknetToken_);
        poolId = poolId_;

        _grantRole(DEFAULT_ADMIN_ROLE, admin_);
        _grantRole(MERKLE_ROLE, admin_);
        _grantRole(STRATEGY_ROLE, admin_);
        _grantRole(TRANSFER_ROLE, admin_);

        slippageBps = 500;
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
        bytes32 root = merkleRoots[epoch];                          // [G3] single SLOAD
        if (root == bytes32(0)) revert NoRootForEpoch();
        if (claimed[epoch][msg.sender]) revert AlreadyClaimed();

        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(msg.sender, amount))));
        if (!MerkleProof.verify(proof, root, leaf)) revert InvalidProof();

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

    /// @notice Execute strategy manually with optional MEV protection  [B2]
    /// @param amount AWP amount to use
    /// @param minAmountOut For BuybackBurn: minimum output (0 = use slippageBps). Prevents sandwich attacks.
    function executeStrategy(uint256 amount, uint256 minAmountOut) external nonReentrant onlyRole(STRATEGY_ROLE) {
        if (strategyPaused) revert StrategyIsPaused();
        if (amount == 0) revert ZeroAmount();
        AWPStrategy strategy = currentStrategy;
        if (strategy == AWPStrategy.AddLiquidity) {
            _addSingleSidedLiquidity(amount);
        } else if (strategy == AWPStrategy.BuybackBurn) {
            _buybackAndBurn(amount, minAmountOut);
        }
        emit AWPProcessed(strategy, amount);   // [B3] emit for all strategies including Reserve
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
            } else if (strategy == AWPStrategy.BuybackBurn) {
                _buybackAndBurn(amount, 0); // auto: use slippage, no explicit minOut
            }
            emit AWPProcessed(strategy, amount);   // [B3] emit for all strategies including Reserve
        }
        return IERC1363Receiver.onTransferReceived.selector;
    }

    // ═══════════════════════════════════════════════
    //  Token Transfer (TRANSFER_ROLE)
    // ═══════════════════════════════════════════════

    function transferToken(address token, address to, uint256 amount) external onlyRole(TRANSFER_ROLE) {
        address resolved = IAWPRegistry(awpRegistry).resolveRecipient(to);
        IERC20(token).safeTransfer(resolved, amount);
        emit TokenTransferred(token, resolved, amount);
    }

    function batchTransferToken(address token, address[] calldata recipients, uint256[] calldata amounts)
        external onlyRole(TRANSFER_ROLE)
    {
        if (recipients.length != amounts.length) revert ArrayLengthMismatch();
        address[] memory resolved = IAWPRegistry(awpRegistry).batchResolveRecipients(recipients);
        for (uint256 i = 0; i < resolved.length;) {
            IERC20(token).safeTransfer(resolved[i], amounts[i]);
            emit TokenTransferred(token, resolved[i], amounts[i]);
            unchecked { ++i; }
        }
    }

    // ═══════════════════════════════════════════════
    //  Configuration
    // ═══════════════════════════════════════════════

    function setSlippageTolerance(uint16 bps) external onlyRole(STRATEGY_ROLE) {   // [O5] uint16
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
    //  DEX-specific virtual functions
    // ═══════════════════════════════════════════════

    function _addSingleSidedLiquidity(uint256 amount) internal virtual;
    /// @param minAmountOut 0 = compute from slippageBps; >0 = use directly (MEV protection) [B2]
    function _buybackAndBurn(uint256 amount, uint256 minAmountOut) internal virtual;
}
