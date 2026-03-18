// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IERC1363Receiver — ERC1363 transfer callback interface
/// @notice Receiving contracts implement this to respond to transferAndCall
interface IERC1363Receiver {
    /// @notice Callback invoked after a token transfer
    /// @param operator Address that initiated the transfer
    /// @param from Token source address
    /// @param amount Transfer amount
    /// @param data Additional data
    /// @return Must return IERC1363Receiver.onTransferReceived.selector
    function onTransferReceived(address operator, address from, uint256 amount, bytes calldata data)
        external
        returns (bytes4);
}

/// @title IERC1363Spender — ERC1363 approval callback interface
/// @notice Approved contracts implement this to respond to approveAndCall
interface IERC1363Spender {
    /// @notice Callback invoked after a token approval
    /// @param owner Approver address
    /// @param amount Approved amount
    /// @param data Additional data
    /// @return Must return IERC1363Spender.onApprovalReceived.selector
    function onApprovalReceived(address owner, uint256 amount, bytes calldata data) external returns (bytes4);
}
