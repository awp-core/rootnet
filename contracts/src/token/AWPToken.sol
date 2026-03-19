// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {IERC1363Receiver, IERC1363Spender} from "../interfaces/IERC1363Receiver.sol";

/// @title AWPToken — AWP main token
/// @notice ERC20Permit + minter model + ERC1363 callbacks
/// @dev Total supply capped at 10 billion (MAX_SUPPLY = 10B × 1e18).
///      Constructor pre-mints INITIAL_MINT to the deployer for distribution; remainder minted on-demand by AWPEmission.
///      Minter management flow: admin → addMinter(awpEmission) → renounceAdmin(); minter list is then permanently immutable.
///      Inherits ERC20Permit for gasless approvals (ERC-2612); inherits ERC20Burnable for token burns.
///      NOTE: ERC20Votes intentionally omitted — AWPDAO uses StakeNFT position-based voting, not ERC20 checkpoints.
contract AWPToken is ERC20, ERC20Permit, ERC20Burnable {
    /// @notice AWP maximum supply: 10 billion tokens (18 decimals)
    uint256 public constant MAX_SUPPLY = 10_000_000_000 * 1e18;

    /// @notice Constructor pre-mint: 200 million tokens transferred to the deployer for distribution
    uint256 public constant INITIAL_MINT = 200_000_000 * 1e18;

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
        // Pre-mint 200M (2% of max supply) to deployer for subsequent distribution
        _mint(deployer_, INITIAL_MINT);
    }

    /// @notice Add a minter (admin only)
    /// @param minter Address to authorize as minter (typically the AWPEmission contract)
    function addMinter(address minter) external {
        if (msg.sender != admin) revert NotAdmin();
        minters[minter] = true;
    }

    /// @notice Permanently renounce admin privileges — no new minters can ever be added
    function renounceAdmin() external {
        if (msg.sender != admin) revert NotAdmin();
        admin = address(0);
    }

    /// @notice Mint new AWP tokens (only authorized minters)
    /// @param to     Recipient address
    /// @param amount Amount to mint (18 decimals)
    function mint(address to, uint256 amount) external {
        if (!minters[msg.sender]) revert NotMinter();
        if (totalSupply() + amount > MAX_SUPPLY) revert ExceedsMaxSupply();
        _mint(to, amount);
    }

    /// @notice Mint and notify the receiver via ERC1363 callback (used by AWPEmission.settleEpoch)
    /// @param to     Recipient address
    /// @param amount Amount to mint (18 decimals)
    /// @param data   Additional data passed to the receiver's onTransferReceived callback
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
}
