// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import {ERC20BurnableUpgradeable} from
    "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20BurnableUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {IERC1363Receiver, IERC1363Spender} from "../interfaces/IERC1363Receiver.sol";

/// @title AlphaToken — Subnet token (CREATE2 deterministic deployment)
/// @notice Dual-minter + ERC1363 callbacks + burnable
/// @dev Deployed via AlphaTokenFactory using CREATE2 (full standalone contract, no proxy).
///      Does **not** inherit OwnableUpgradeable (admin is managed by this contract directly).
///      MAX_SUPPLY = 10 billion tokens per subnet.
///      Minter lifecycle:
///        1. During initialize(), admin is set as the initial minter (dual-minter phase)
///        2. setSubnetMinter(subnetManager) sets the subnet contract as the sole minter,
///           revokes admin's minting rights, and permanently sets mintersLocked = true
///        3. setMinterPaused() can pause/resume a specific minter (used for banning/unbanning subnets)
contract AlphaToken is Initializable, ERC20Upgradeable, ERC20BurnableUpgradeable {
    /// @notice Maximum supply per subnet token: 10 billion tokens (18 decimals)
    uint256 public constant MAX_SUPPLY = 10_000_000_000 * 1e18;

    /// @notice Subnet ID this token belongs to
    uint256 public subnetId;

    /// @notice Admin address (RootNet), responsible for minter management and pause control
    address public admin;

    /// @notice Minter whitelist
    mapping(address => bool) public minters;

    /// @notice Minter pause status (used to pause minting rights when banning a subnet)
    mapping(address => bool) public minterPaused;

    /// @notice Whether the minter list is permanently locked (set to true after setSubnetMinter)
    bool public mintersLocked;

    /// @notice Token creation timestamp (used for time-based minting cap)
    uint64 public createdAt;

    /// @notice Total supply snapshot at the moment setSubnetMinter locks the minter list
    /// @dev Pre-minted tokens (admin LP mint) are excluded from the time-based cap calculation
    uint256 public supplyAtLock;

    /// @notice Cumulative gross tokens minted since setSubnetMinter was called (not affected by burns)
    uint256 public grossMintedSinceLock;

    /// @notice Emitted when the subnet minter is permanently set
    event SubnetMinterSet(address indexed subnetManager);

    /// @dev Caller is not the admin
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
    /// @param subnetId_ Subnet ID
    /// @param admin_    Admin address (typically RootNet); also acts as minter in the initial phase
    /// @dev Can only be called once (enforced by the initializer modifier); admin_ is automatically added to the minter list
    function initialize(string memory name_, string memory symbol_, uint256 subnetId_, address admin_)
        external
        initializer
    {
        __ERC20_init(name_, symbol_);
        __ERC20Burnable_init();
        subnetId = subnetId_;
        admin = admin_;
        createdAt = uint64(block.timestamp);
        // Admin holds minting rights in the initial phase (first minter in the dual-minter setup)
        minters[admin_] = true;
    }

    /// @notice Set the subnet contract as the sole minter, revoke admin minting rights, and permanently lock
    /// @param subnetManager Subnet contract address (must not be the zero address)
    /// @dev Can only be called once; after the call mintersLocked = true and setSubnetMinter cannot be called again.
    ///      This ensures the subnet contract is the only minting source, preventing admin from minting arbitrarily.
    function setSubnetMinter(address subnetManager) external {
        if (msg.sender != admin) revert NotAdmin();
        // Ensure not yet locked
        if (mintersLocked) revert MintersLocked();
        // Subnet contract address must not be zero
        if (subnetManager == address(0)) revert NotMinter();
        // Authorize subnet contract as minter
        minters[subnetManager] = true;
        // Revoke admin's minting rights
        minters[admin] = false;
        // Snapshot supply and reset clock before locking
        supplyAtLock = totalSupply();
        createdAt = uint64(block.timestamp);
        grossMintedSinceLock = 0;
        // Permanently lock the minter list
        mintersLocked = true;
        emit SubnetMinterSet(subnetManager);
    }

    /// @notice Pause or resume a specific minter's minting rights (used for banning/unbanning subnets)
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
        // Time-based minting cap (only applies after minters are locked, i.e. subnet minter phase)
        // Admin minting (initial LP liquidity) is exempt from the time cap
        if (mintersLocked) {
            uint256 elapsed = block.timestamp - uint256(createdAt);
            if (elapsed == 0) elapsed = 1; // Allow same-block mint after setSubnetMinter
            uint256 maxMintable = (MAX_SUPPLY - supplyAtLock) * elapsed / 365 days;
            if (maxMintable > MAX_SUPPLY - supplyAtLock) maxMintable = MAX_SUPPLY - supplyAtLock;
            if (grossMintedSinceLock + amount > maxMintable) revert ExceedsMintableLimit();
            grossMintedSinceLock += amount;
        }
        _mint(to, amount);
    }

    /// @notice Get the remaining mintable headroom for subnet minters at the current timestamp
    /// @return Remaining amount that can be minted before hitting the time-based cap
    function currentMintableLimit() external view returns (uint256) {
        if (!mintersLocked) return 0;
        uint256 elapsed = block.timestamp - uint256(createdAt);
        if (elapsed == 0) elapsed = 1; // Allow same-block query after setSubnetMinter
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
