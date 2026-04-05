// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {IERC1363Receiver, IERC1363Spender} from "../interfaces/IERC1363Receiver.sol";

/// @title AWPToken — AWP main token
/// @author AWP Protocol Team
/// @notice ERC20Permit + minter model + ERC1363 callbacks
///
/// @dev Total supply capped at 10 billion (MAX_SUPPLY = 10B × 1e18).
///      Constructor pre-mints INITIAL_MINT to the deployer for distribution; remainder minted on-demand by AWPEmission.
///      Minter management flow: admin → addMinter(awpEmission) → renounceAdmin(); minter list is then permanently immutable.
///      Inherits ERC20Permit for gasless approvals (ERC-2612); inherits ERC20Burnable for token burns.
///      NOTE: ERC20Votes intentionally omitted — AWPDAO uses veAWP position-based voting, not ERC20 checkpoints.
///
/// @dev ──────────────────────────────────────────────────────────────────────
///      Security Design Notes
///      ──────────────────────────────────────────────────────────────────────
///
///      DEPLOYMENT STATE
///      ~~~~~~~~~~~~~~~~
///      Admin has been permanently renounced (admin == address(0)).
///      A single minter (AWPEmission contract) has been registered and its
///      ownership is also renounced / locked, making the minter set immutable.
///      initialMint() was never invoked; all token supply originates from the
///      AWPEmission minter via scheduled emission. Since admin == address(0),
///      initialMint() is now permanently uncallable.
///
///      1. ADMIN & MINTER IMMUTABILITY — admin has been renounced on-chain,
///         so addMinter(), initialMint(), and renounceAdmin() are all permanently
///         dead code. The minter whitelist is frozen with exactly one entry
///         (AWPEmission). No removeMinter() exists by design: removing the sole
///         minter would halt all emissions, and adding it back would be impossible
///         after admin renunciation. The single-minter-then-renounce pattern
///         ensures the token's mint authority is fully deterministic.
///
///      2. NO REENTRANCY GUARD — transferAndCall, mintAndCall, and approveAndCall
///         follow the Checks-Effects-Interactions (CEI) pattern: all state changes
///         (via OpenZeppelin's ERC20._transfer / _mint / _approve) complete before
///         any external callback is invoked. Even if the callback re-enters this
///         contract, balances and allowances are already consistent. A ReentrancyGuard
///         modifier was evaluated but omitted because CEI compliance makes it
///         redundant for this contract's own state invariants. Integrators that
///         perform multi-step operations inside callbacks should apply their own
///         reentrancy protection.
///
///      3. ERC1363 SUBSET — This contract implements transferAndCall and
///         approveAndCall but intentionally omits transferFromAndCall. Only the
///         direct-transfer callback pattern is required by the AWP protocol.
///         The contract does not register ERC-1363 via ERC-165 supportsInterface,
///         so no standard-compliance claim is made; these are simply convenience
///         methods following the ERC-1363 callback convention.
///
///      4. mintAndCall CALLBACK CONVENTION — onTransferReceived(operator, from,
///         amount, data) receives msg.sender (the minter) for both `operator`
///         and `from`. In a mint context no prior owner exists; passing the minter
///         address (rather than address(0)) lets the receiver identify which
///         authorized minter initiated the operation. In current deployment the
///         sole minter (AWPEmission) uses mint() exclusively; mintAndCall() is
///         reserved for future emission contracts and is not actively called.
///
///      5. EVENT COVERAGE — addMinter() and renounceAdmin() do not emit custom
///         events. These functions are part of a one-time, irreversible setup
///         sequence (deploy → addMinter → renounceAdmin, typically 2 transactions).
///         Admin has already been renounced on-chain, so these functions are now
///         permanently unreachable. Off-chain indexing of the setup can rely on
///         transaction receipts from the deployment block range.
///
///      6. INPUT VALIDATION IN ADMIN FUNCTIONS — addMinter() does not include
///         an explicit address(0) check. This is acceptable because: (a) only the
///         admin could call it and admin has been permanently renounced, making the
///         function unreachable; (b) even hypothetically, minting from address(0) is
///         impossible since no transaction can originate from address(0) to pass the
///         minters[msg.sender] check in mint(); (c) OpenZeppelin's ERC20._mint()
///         independently reverts on mint-to-zero-address, providing a secondary guard.
///
///      ──────────────────────────────────────────────────────────────────────
contract AWPToken is ERC20, ERC20Permit, ERC20Burnable {
    /// @notice AWP maximum supply: 10 billion tokens (18 decimals)
    uint256 public constant MAX_SUPPLY = 10_000_000_000 * 1e18;

    /// @notice Minter whitelist; only authorized addresses may call mint()
    /// @dev Frozen after admin renunciation — contains exactly one entry (AWPEmission)
    mapping(address => bool) public minters;

    /// @notice Admin address; permanently set to address(0) after renounceAdmin()
    /// @dev Admin has been renounced on-chain; all admin-gated functions are unreachable
    address public admin;

    /// @notice Whether initialMint has been called (one-time only)
    /// @dev Never invoked in current deployment; permanently uncallable (admin == address(0))
    bool public initialMinted;

    /// @dev Caller is not the admin
    error NotAdmin();
    /// @dev Minter is invalid or unauthorized
    error NotMinter();
    /// @dev Minting would cause total supply to exceed MAX_SUPPLY
    error ExceedsMaxSupply();
    /// @dev ERC1363 callback returned an incorrect value
    error InvalidCallback();
    /// @dev initialMint already called
    error AlreadyMinted();

    /// @notice Deploy the AWP token
    /// @param name_     Token name (e.g. "AWP Token")
    /// @param symbol_   Token symbol (e.g. "AWP")
    /// @param deployer_ Deployer address — receives admin rights
    /// @dev Constructor does NOT mint — keeps CREATE2 bytecode identical across chains.
    ///      Call initialMint(amount) after deployment to pre-mint tokens (amount configurable per chain, 0 = skip).
    constructor(string memory name_, string memory symbol_, address deployer_)
        ERC20(name_, symbol_)
        ERC20Permit(name_)
    {
        admin = deployer_;
    }

    // ── Admin Functions (setup-only; permanently unreachable after renounceAdmin) ──

    /// @notice One-time pre-mint to admin for distribution (admin only, callable once)
    /// @param amount Amount to mint (0 = no-op, allows chains that need no pre-mint)
    /// @dev Permanently unreachable: admin has been renounced on-chain
    function initialMint(uint256 amount) external {
        if (msg.sender != admin) revert NotAdmin();
        if (initialMinted) revert AlreadyMinted();
        initialMinted = true;
        if (amount > 0) {
            if (amount > MAX_SUPPLY) revert ExceedsMaxSupply();
            _mint(admin, amount);
        }
    }

    /// @notice Add a minter (admin only)
    /// @param minter Address to authorize as minter (typically the AWPEmission contract)
    /// @dev Permanently unreachable: admin has been renounced on-chain.
    ///      No removeMinter() by design — see Security Design Notes §1.
    function addMinter(address minter) external {
        if (msg.sender != admin) revert NotAdmin();
        minters[minter] = true;
    }

    /// @notice Permanently renounce admin privileges — no new minters can ever be added
    /// @dev Irreversible. Already executed on-chain; this function is now permanently unreachable.
    function renounceAdmin() external {
        if (msg.sender != admin) revert NotAdmin();
        admin = address(0);
    }

    // ── Minting Functions ──

    /// @notice Mint new AWP tokens (only authorized minters)
    /// @param to     Recipient address
    /// @param amount Amount to mint (18 decimals)
    /// @dev Supply cap enforced: totalSupply() + amount must not exceed MAX_SUPPLY.
    ///      Currently callable only by the sole registered minter (AWPEmission).
    function mint(address to, uint256 amount) external {
        if (!minters[msg.sender]) revert NotMinter();
        if (totalSupply() + amount > MAX_SUPPLY) revert ExceedsMaxSupply();
        _mint(to, amount);
    }

    /// @notice Mint and notify the receiver via ERC1363 callback
    /// @param to     Recipient address
    /// @param amount Amount to mint (18 decimals)
    /// @param data   Additional data passed to the receiver's onTransferReceived callback
    /// @dev CEI-compliant: _mint() completes all state changes before the external callback.
    ///      Callback signature: onTransferReceived(operator=minter, from=minter, amount, data) —
    ///      see Security Design Notes §4 for the from-parameter convention.
    ///      Not actively invoked in current deployment; the sole minter uses mint() exclusively.
    function mintAndCall(address to, uint256 amount, bytes calldata data) external {
        if (!minters[msg.sender]) revert NotMinter();
        if (totalSupply() + amount > MAX_SUPPLY) revert ExceedsMaxSupply();
        _mint(to, amount);
        if (to.code.length > 0) {
            bytes4 retval = IERC1363Receiver(to).onTransferReceived(msg.sender, msg.sender, amount, data);
            if (retval != IERC1363Receiver.onTransferReceived.selector) revert InvalidCallback();
        }
    }

    // ── ERC1363 Callbacks (intentional subset — see Security Design Notes §3) ──

    /// @notice Transfer tokens and notify the receiving contract (ERC1363 transferAndCall)
    /// @param to     Recipient address
    /// @param amount Transfer amount
    /// @param data   Additional data passed to the receiver's onTransferReceived callback
    /// @return Always returns true (reverts on failure)
    /// @dev CEI-compliant: transfer() completes all state changes before the external callback.
    ///      If the recipient is a contract, calls onTransferReceived and verifies the return value.
    ///      Callback receives operator=msg.sender, from=msg.sender (direct transfer, no delegation).
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
    /// @dev CEI-compliant: approve() completes all state changes before the external callback.
    ///      If the spender is a contract, calls onApprovalReceived and verifies the return value.
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
