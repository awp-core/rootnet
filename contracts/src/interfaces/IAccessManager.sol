// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IAccessManager — Principal registration + Agent permission management interface
/// @notice All write functions are callable only by RootNet
interface IAccessManager {
    function rootNet() external view returns (address);
    function isRegistered(address) external view returns (bool);
    function registeredAt(address) external view returns (uint64);
    function totalUsers() external view returns (uint256);

    // ── Write ──
    function register(address user) external;
    /// @notice Bind an agent to a principal; auto-registers principal; returns old principal (unchanged = same principal, no-op)
    function bind(address agent, address principal) external returns (address oldPrincipal);
    /// @notice Unbind an agent from its principal, returning it to unregistered status
    function unbind(address agent) external returns (address oldPrincipal);
    function removeAgent(address user, address agent, address operator) external;
    function setManager(address user, address agent, bool _isManager, address operator) external;
    function setRewardRecipient(address user, address recipient) external;

    // ── Query ──
    function getOwner(address addr) external view returns (address);
    function isRegisteredUser(address addr) external view returns (bool);
    function isRegisteredAgent(address addr) external view returns (bool);
    function isKnownAddress(address addr) external view returns (bool);
    function getAgents(address user) external view returns (address[] memory);
    /// @notice Check whether agent belongs to user (a user is also considered their own valid agent)
    function isAgent(address user, address agent) external view returns (bool);
    function isManagerAgent(address agent) external view returns (bool);
    function getRewardRecipient(address user) external view returns (address);
    function getTotalUsers() external view returns (uint256);
    /// @notice Resolve address role in a single call (replaces 3 external calls)
    function resolveCallerRole(address addr) external view returns (address owner, bool isUser, bool isManager_);
}
