// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IAlphaTokenFactory — AlphaToken factory interface
/// @notice Deploys AlphaToken instances for each subnet via CREATE2
interface IAlphaTokenFactory {
    function awpRegistry() external view returns (address);
    function configured() external view returns (bool);

    /// @notice Set the AWPRegistry address and renounce ownership (one-time call)
    function setAddresses(address awpRegistry) external;
    /// @notice Deploy a new AlphaToken via CREATE2 (only AWPRegistry may call)
    /// @param salt CREATE2 salt; bytes32(0) skips vanity validation and uses subnetId as salt
    function deploy(uint256 subnetId, string memory name, string memory symbol, address admin, bytes32 salt)
        external
        returns (address);
}
