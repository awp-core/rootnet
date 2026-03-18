// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";
import {ERC20Votes} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Votes.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {Nonces} from "@openzeppelin/contracts/utils/Nonces.sol";
import {IERC1363Receiver, IERC1363Spender} from "../interfaces/IERC1363Receiver.sol";

/// @title AWPToken — AWP main token
/// @notice ERC20Votes + minter model + ERC1363 callbacks
/// @dev Total supply capped at 10 billion (MAX_SUPPLY = 10B × 1e18).
///      Constructor pre-mints 50% (INITIAL_MINT = 5B) to the deployer; remaining 50% minted on-demand by AWPEmission.
///      Minter management flow: admin → addMinter(awpEmission) → renounceAdmin(); minter list is then permanently immutable.
///      Inherits ERC20Votes for on-chain governance voting; inherits ERC20Burnable for token burns.
contract AWPToken is ERC20, ERC20Permit, ERC20Votes, ERC20Burnable {
    /// @notice AWP maximum supply: 10 billion tokens (18 decimals)
    uint256 public constant MAX_SUPPLY = 10_000_000_000 * 1e18;

    /// @notice Constructor pre-mint: 5 billion tokens transferred to the deployer for distribution
    uint256 public constant INITIAL_MINT = 5_000_000_000 * 1e18;

    /// @notice Minter whitelist; only authorized addresses may call mint()
    mapping(address => bool) public minters;

    /// @notice Admin address; set to address(0) after renounceAdmin(), permanently locked
    address public admin;

    /// @dev Caller is not the admin
    error NotAdmin();
    /// @dev Minter is invalid or unauthorized
    error NotMinter();
    /// @dev Minting would cause total supply to exceed MAX_SUPPLY
    error ExceedsMaxSupply();
    /// @dev ERC1363 callback returned an incorrect value
    error InvalidCallback();

    /// @notice Deploy the AWP token
    /// @param name_     Token name (e.g. "AWP Token")
    /// @param symbol_   Token symbol (e.g. "AWP")
    /// @param deployer_ Deployer address — receives admin rights and the INITIAL_MINT pre-mint
    /// @dev deployer_ is **not** a minter; AWPEmission must be explicitly authorized via addMinter
    constructor(string memory name_, string memory symbol_, address deployer_)
        ERC20(name_, symbol_)
        ERC20Permit(name_)
    {
        // Set deployer as admin
        admin = deployer_;
        // Pre-mint 50% of supply to deployer for subsequent distribution (Treasury / DAO / LP, etc.)
        _mint(deployer_, INITIAL_MINT);
    }

    /// @notice Add a minter (admin only)
    /// @param minter Address to authorize as minter (typically the AWPEmission contract)
    /// @dev May be called multiple times before renounceAdmin(); list is permanently locked after renounce
    function addMinter(address minter) external {
        if (msg.sender != admin) revert NotAdmin();
        // Zero address may not be a minter
        if (minter == address(0)) revert NotMinter();
        minters[minter] = true;
    }

    /// @notice Admin renounces control (minter list is permanently locked after this call)
    /// @dev Sets admin to address(0); addMinter() and renounceAdmin() can no longer be called
    function renounceAdmin() external {
        if (msg.sender != admin) revert NotAdmin();
        admin = address(0);
    }

    /// @notice Mint tokens (authorized minters only)
    /// @param to     Recipient address
    /// @param amount Amount to mint (18 decimals)
    /// @dev Checks total supply does not exceed MAX_SUPPLY before minting
    function mint(address to, uint256 amount) external {
        if (!minters[msg.sender]) revert NotMinter();
        if (totalSupply() + amount > MAX_SUPPLY) revert ExceedsMaxSupply();
        _mint(to, amount);
    }

    /// @notice Mint tokens and notify the receiving contract (ERC1363-style callback)
    /// @param to     Recipient address
    /// @param amount Amount to mint
    /// @param data   Additional data passed to the receiver's onTransferReceived callback
    /// @dev Combines mint + ERC1363 callback in one call. If recipient is a contract,
    ///      calls onTransferReceived. Used by AWPEmission to auto-trigger SubnetManager strategies.
    function mintAndCall(address to, uint256 amount, bytes calldata data) external {
        if (!minters[msg.sender]) revert NotMinter();
        if (totalSupply() + amount > MAX_SUPPLY) revert ExceedsMaxSupply();
        _mint(to, amount);
        if (to.code.length > 0) {
            bytes4 retval = IERC1363Receiver(to).onTransferReceived(msg.sender, msg.sender, amount, data);
            if (retval != IERC1363Receiver.onTransferReceived.selector) revert InvalidCallback();
        }
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
            // Revert if the receiver does not return the correct selector
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

    // ── Required overrides (resolve diamond inheritance between ERC20 and ERC20Votes) ──

    /// @dev Override _update to update both ERC20 balances and ERC20Votes voting power
    /// @param from  Sender address
    /// @param to    Recipient address
    /// @param value Transfer amount
    function _update(address from, address to, uint256 value) internal override(ERC20, ERC20Votes) {
        super._update(from, to, value);
    }

    /// @dev Override nonces to resolve conflict between ERC20Permit and Nonces
    /// @param owner Address to query
    /// @return Current nonce value
    function nonces(address owner) public view override(ERC20Permit, Nonces) returns (uint256) {
        return super.nonces(owner);
    }
}
