// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import {ReentrancyGuardTransient} from "@openzeppelin/contracts/utils/ReentrancyGuardTransient.sol";

interface IveAWPDeposit {
    function deposit(uint256 amount, uint64 lockDuration) external returns (uint256 tokenId);
}

/// @title VeAWPHelper — Gasless staking helper for veAWP
/// @notice User signs one ERC-2612 permit off-chain, relayer calls depositFor().
/// @dev Stateless, non-upgradeable. Max-approves veAWP at construction.
contract VeAWPHelper is ReentrancyGuardTransient {
    address public immutable awpToken;
    address public immutable veAWP;

    error ZeroAmount();
    error ZeroAddress();
    error InvalidUser();
    error TransferFailed();
    error ApproveFailed();

    constructor(address awpToken_, address veAWP_) {
        if (awpToken_ == address(0) || veAWP_ == address(0)) revert ZeroAddress();
        awpToken = awpToken_;
        veAWP = veAWP_;
        if (!IERC20(awpToken_).approve(veAWP_, type(uint256).max)) revert ApproveFailed();
    }

    /// @notice Gasless deposit: user signs ERC-2612 permit, relayer pays gas.
    ///         Flow: permit(user→helper) → transferFrom(user→helper) → deposit → NFT(helper→user)
    /// @param user Permit signer, receives the veAWP NFT
    /// @param amount AWP amount to stake (must match permit value)
    /// @param lockDuration Lock period in seconds (min 1 day)
    /// @param deadline Permit expiry timestamp
    /// @param v ERC-2612 permit v
    /// @param r ERC-2612 permit r
    /// @param s ERC-2612 permit s
    /// @return tokenId Minted veAWP position NFT token ID
    function depositFor(
        address user,
        uint256 amount,
        uint64 lockDuration,
        uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant returns (uint256 tokenId) {
        if (amount == 0) revert ZeroAmount();
        if (user == address(0)) revert ZeroAddress();
        if (user == address(this)) revert InvalidUser();

        // 1. Execute ERC-2612 permit (user approves helper). Silently ignore failure (user may have pre-approved).
        _permit(user, amount, deadline, v, r, s);

        // 2. Pull AWP from user (helper already max-approved veAWP in constructor)
        if (!IERC20(awpToken).transferFrom(user, address(this), amount)) revert TransferFailed();

        // 3. Deposit into veAWP (NFT minted to helper)
        tokenId = IveAWPDeposit(veAWP).deposit(amount, lockDuration);

        // 4. Transfer NFT to user
        IERC721(veAWP).transferFrom(address(this), user, tokenId);
    }

    // NOTE: addToPositionFor is intentionally NOT provided.
    // veAWP._update hook checks allocation coverage on every NFT transfer.
    // The transfer-add-transfer pattern (user→helper→addToPosition→helper→user)
    // reverts at the first transfer for any user with active allocations,
    // making it non-functional for the primary use case.
    // Users who need to add to positions must call veAWP.addToPosition directly.

    /// @dev Low-level permit call (~200 gas cheaper than try/catch). Silently ignores failure.
    function _permit(address user, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) internal {
        address token = awpToken;
        assembly {
            let m := mload(0x40)
            mstore(m, 0xd505accf00000000000000000000000000000000000000000000000000000000)
            mstore(add(m, 0x04), user)
            mstore(add(m, 0x24), address())
            mstore(add(m, 0x44), amount)
            mstore(add(m, 0x64), deadline)
            mstore(add(m, 0x84), v)
            mstore(add(m, 0xa4), r)
            mstore(add(m, 0xc4), s)
            pop(call(gas(), token, 0, m, 0xe4, 0, 0))
        }
    }
}
