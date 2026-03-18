// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @title IAWPToken — AWP main token interface
/// @notice Extends IERC20 with minter model + ERC1363 callbacks + burn
interface IAWPToken is IERC20 {
    function MAX_SUPPLY() external view returns (uint256);
    function admin() external view returns (address);
    function minters(address) external view returns (bool);

    /// @notice Mint tokens (authorized minters only)
    function mint(address to, uint256 amount) external;
    /// @notice Mint tokens and trigger ERC1363 onTransferReceived callback on recipient
    function mintAndCall(address to, uint256 amount, bytes calldata data) external;
    /// @notice Burn tokens held by the caller
    function burn(uint256 amount) external;
    /// @notice Burn tokens from another address (requires allowance)
    function burnFrom(address account, uint256 amount) external;
    /// @notice Add a minter (admin only)
    function addMinter(address minter) external;
    /// @notice Permanently renounce admin rights, locking the minter list
    function renounceAdmin() external;

    /// @notice ERC1363 transfer callback
    function transferAndCall(address to, uint256 amount, bytes calldata data) external returns (bool);
    /// @notice ERC1363 approval callback
    function approveAndCall(address spender, uint256 amount, bytes calldata data) external returns (bool);
}
