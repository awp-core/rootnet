// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import {ERC20BurnableUpgradeable} from
    "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20BurnableUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {IERC1363Receiver, IERC1363Spender} from "../interfaces/IERC1363Receiver.sol";

/// @title AlphaToken — Worknet token (CREATE2 deterministic deployment)
/// @notice Dual-minter + ERC1363 callbacks + burnable
/// @dev Deployed via AlphaTokenFactory using CREATE2 (full standalone contract, no proxy).
///      Does **not** inherit OwnableUpgradeable (admin is managed by this contract directly).
///      MAX_SUPPLY = 10 billion tokens per worknet.
///      Minter lifecycle:
///        1. During initialize(), admin is set as the initial minter (dual-minter phase)
///        2. setWorknetMinter(worknetManager) sets the worknet contract as the sole minter,
///           revokes admin's minting rights, and permanently sets mintersLocked = true
///        3. setMinterPaused() can pause/resume a specific minter (used for banning/unbanning worknets)
contract AlphaToken is Initializable, ERC20Upgradeable, ERC20BurnableUpgradeable {
    /// @notice Maximum supply per worknet token: 10 billion tokens (18 decimals)
    uint256 public constant MAX_SUPPLY = 10_000_000_000 * 1e18;

    /// @notice Worknet ID this token belongs to
    uint256 public worknetId;

    /// @notice Admin address (AWPRegistry), responsible for minter management and pause control
    /// @dev Packed with mintersLocked (1 byte) + createdAt (8 bytes) into one 32-byte slot
    address public admin;
    bool public mintersLocked;
    uint64 public createdAt;

    /// @notice Minter whitelist
    mapping(address => bool) public minters;

    /// @notice Minter pause status (used to pause minting rights when banning a worknet)
    mapping(address => bool) public minterPaused;

    /// @notice Total supply snapshot at the moment setWorknetMinter locks the minter list
    /// @dev Pre-minted tokens (admin LP mint) are excluded from the time-based cap calculation
    uint256 public supplyAtLock;

    /// @notice Cumulative gross tokens minted since setWorknetMinter was called (not affected by burns)
    uint256 public grossMintedSinceLock;

    /// @notice Emitted when the worknet minter is permanently set
    event WorknetMinterSet(address indexed worknetManager);

    /// @dev Caller is not the admin
    error ZeroAddress();
    error NotAdmin();
    /// @dev Minter is invalid or unauthorized
    error NotMinter();
    /// @dev Minter has been paused
    error MinterPaused();
    /// @dev Minter list is permanently locked and cannot be modified
    error MintersLocked();
    /// @dev Minting would cause total supply to exceed MAX_SUPPLY
    error ExceedsMaxSupply();
    /// @dev Minting would exceed time-based cap (MAX_SUPPLY * elapsed / 365 days)
    error ExceedsMintableLimit();
    /// @dev ERC1363 callback returned an incorrect value
    error InvalidCallback();

    /// @dev No-op constructor; the initializer modifier on initialize() prevents double-init.
    ///      Empty constructor is required for CREATE2 deployment via AlphaTokenFactory.
    constructor() {}

    /// @notice Initialize an AlphaToken clone instance (called by AlphaTokenFactory.deploy())
    /// @param name_     Token name
    /// @param symbol_   Token symbol
    /// @param worknetId_ Worknet ID
    /// @param admin_    Admin address (typically AWPRegistry); also acts as minter in the initial phase
    /// @dev Can only be called once (enforced by the initializer modifier); admin_ is automatically added to the minter list
    function initialize(string memory name_, string memory symbol_, uint256 worknetId_, address admin_)
        external
        initializer
    {
        __ERC20_init(name_, symbol_);
        __ERC20Burnable_init();
        worknetId = worknetId_;
        admin = admin_;
        createdAt = uint64(block.timestamp);
        // Admin holds minting rights in the initial phase (first minter in the dual-minter setup)
        minters[admin_] = true;
    }

    /// @notice Set the worknet contract as the sole minter, revoke admin minting rights, and permanently lock
    /// @param worknetManager Worknet contract address (must not be the zero address)
    /// @dev Can only be called once; after the call mintersLocked = true and setWorknetMinter cannot be called again.
    ///      This ensures the worknet contract is the only minting source, preventing admin from minting arbitrarily.
    function setWorknetMinter(address worknetManager) external {
        if (msg.sender != admin) revert NotAdmin();
        // Ensure not yet locked
        if (mintersLocked) revert MintersLocked();
        // Worknet contract address must not be zero
        if (worknetManager == address(0)) revert ZeroAddress();
        // Authorize worknet contract as minter
        minters[worknetManager] = true;
        // Revoke admin's minting rights
        minters[admin] = false;
        // Snapshot supply and reset clock before locking
        supplyAtLock = totalSupply();
        createdAt = uint64(block.timestamp);
        grossMintedSinceLock = 0;
        // Permanently lock the minter list
        mintersLocked = true;
        emit WorknetMinterSet(worknetManager);
    }

    /// @notice Pause or resume a specific minter's minting rights (used for banning/unbanning worknets)
    /// @param minter  Target minter address
    /// @param paused  true = pause, false = resume
    /// @dev Only admin may call; does not modify the minters mapping, only controls the extra check in mint()
    function setMinterPaused(address minter, bool paused) external {
        if (msg.sender != admin) revert NotAdmin();
        minterPaused[minter] = paused;
    }

    /// @notice Mint tokens
    /// @param to     Recipient address
    /// @param amount Amount to mint (18 decimals)
    /// @dev Checks in order: is an authorized minter → is not paused → does not exceed MAX_SUPPLY
    function mint(address to, uint256 amount) external {
        // Check that caller is an authorized minter
        if (!minters[msg.sender]) revert NotMinter();
        // Check that the minter is not paused (banned state)
        if (minterPaused[msg.sender]) revert MinterPaused();
        // Ensure minting does not exceed the supply cap
        uint256 supply = totalSupply();
        if (supply + amount > MAX_SUPPLY) revert ExceedsMaxSupply();
        // Time-based minting cap (only applies after minters are locked, i.e. worknet minter phase)
        // Admin minting (initial LP liquidity) is exempt from the time cap
        if (mintersLocked) {
            uint256 elapsed = block.timestamp - uint256(createdAt);
            if (elapsed == 0) elapsed = 1; // Allow same-block mint after setWorknetMinter
            uint256 maxMintable = (MAX_SUPPLY - supplyAtLock) * elapsed / 365 days;
            if (maxMintable > MAX_SUPPLY - supplyAtLock) maxMintable = MAX_SUPPLY - supplyAtLock;
            if (grossMintedSinceLock + amount > maxMintable) revert ExceedsMintableLimit();
            grossMintedSinceLock += amount;
        }
        _mint(to, amount);
    }

    /// @notice Get the remaining mintable headroom for worknet minters at the current timestamp
    /// @return Remaining amount that can be minted before hitting the time-based cap
    function currentMintableLimit() external view returns (uint256) {
        if (!mintersLocked) return 0;
        uint256 elapsed = block.timestamp - uint256(createdAt);
        if (elapsed == 0) elapsed = 1; // Allow same-block query after setWorknetMinter
        uint256 maxMintable = (MAX_SUPPLY - supplyAtLock) * elapsed / 365 days;
        if (maxMintable > MAX_SUPPLY - supplyAtLock) maxMintable = MAX_SUPPLY - supplyAtLock;
        return grossMintedSinceLock >= maxMintable ? 0 : maxMintable - grossMintedSinceLock;
    }

    // ── ERC1363 callbacks ──

    /// @notice Transfer tokens and notify the receiving contract (ERC1363 transferAndCall)
    /// @param to     Recipient address
    /// @param amount Transfer amount
    /// @param data   Additional data passed to the receiver's onTransferReceived callback
    /// @return Always returns true (reverts on failure)
    /// @dev If the recipient is a contract, calls onTransferReceived and verifies the return value
    function transferAndCall(address to, uint256 amount, bytes calldata data) external returns (bool) {
        transfer(to, amount);
        // If the target address is a contract (code.length > 0), execute the ERC1363 callback
        if (to.code.length > 0) {
            bytes4 retval = IERC1363Receiver(to).onTransferReceived(msg.sender, msg.sender, amount, data);
            // Verify that the receiver returns the correct selector
            if (retval != IERC1363Receiver.onTransferReceived.selector) revert InvalidCallback();
        }
        return true;
    }

    /// @notice Approve and notify the spender contract (ERC1363 approveAndCall)
    /// @param spender Approved address
    /// @param amount  Approved amount
    /// @param data    Additional data passed to the spender's onApprovalReceived callback
    /// @return Always returns true (reverts on failure)
    /// @dev If the spender is a contract, calls onApprovalReceived and verifies the return value
    function approveAndCall(address spender, uint256 amount, bytes calldata data) external returns (bool) {
        approve(spender, amount);
        // If the target address is a contract, execute the ERC1363 approval callback
        if (spender.code.length > 0) {
            bytes4 retval = IERC1363Spender(spender).onApprovalReceived(msg.sender, amount, data);
            // Verify that the spender returns the correct selector
            if (retval != IERC1363Spender.onApprovalReceived.selector) revert InvalidCallback();
        }
        return true;
    }
}
