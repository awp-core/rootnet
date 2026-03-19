// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @title IAlphaToken — Subnet Alpha token interface
/// @notice Clonable upgrade pattern; one instance per subnet; MAX_SUPPLY = 10B
interface IAlphaToken is IERC20 {
    function MAX_SUPPLY() external view returns (uint256);
    function subnetId() external view returns (uint256);
    function admin() external view returns (address);
    function minters(address) external view returns (bool);
    function minterPaused(address) external view returns (bool);
    function mintersLocked() external view returns (bool);

    /// @notice Initialize (called after factory clone, cannot be called again)
    function initialize(string memory name, string memory symbol, uint256 subnetId, address admin) external;
    /// @notice Mint (authorized minters only, not in paused state)
    function mint(address to, uint256 amount) external;
    /// @notice Set subnet contract as sole minter and permanently lock
    function setSubnetMinter(address subnetManager) external;
    /// @notice Pause/resume a specific minter (used for subnet banning)
    function setMinterPaused(address minter, bool paused) external;

    function burn(uint256 amount) external;
    function currentMintableLimit() external view returns (uint256);
    function supplyAtLock() external view returns (uint256);
    function grossMintedSinceLock() external view returns (uint256);
    function createdAt() external view returns (uint64);

    function transferAndCall(address to, uint256 amount, bytes calldata data) external returns (bool);
    function approveAndCall(address spender, uint256 amount, bytes calldata data) external returns (bool);
}
