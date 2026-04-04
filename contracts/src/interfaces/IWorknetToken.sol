// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @title IWorknetToken — Worknet token interface
/// @notice CREATE2 deployed, one per worknet; MAX_SUPPLY = 10B; ERC1363 for callbacks
interface IWorknetToken is IERC20 {
    function MAX_SUPPLY() external view returns (uint256);
    function worknetId() external view returns (uint256);
    function minter() external view returns (address);
    function initialized() external view returns (bool);
    function createdAt() external view returns (uint64);
    function supplyAtLock() external view returns (uint256);

    function setMinter(address newMinter) external;
    function mint(address to, uint256 amount) external;
    function burn(uint256 amount) external;
    function currentMintableLimit() external view returns (uint256);
}
