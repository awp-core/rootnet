// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IWorknetTokenFactory — WorknetToken factory interface
/// @notice Deploys WorknetToken instances for each worknet via CREATE2.
///         Token constructor reads params via pendingName/Symbol/WorknetId callbacks.
interface IWorknetTokenFactory {
    function awpRegistry() external view returns (address);
    function configured() external view returns (bool);

    function setAddresses(address awpRegistry) external;
    function deploy(uint256 worknetId, string memory name, string memory symbol, bytes32 salt)
        external returns (address);
    function predictDeployAddress(bytes32 salt) external view returns (address);

    // Pending deploy params (read by WorknetToken constructor during CREATE2)
    function pendingName() external view returns (string memory);
    function pendingSymbol() external view returns (string memory);
    function pendingWorknetId() external view returns (uint256);
}
