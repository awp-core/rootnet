// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IStakingVault — Staking vault interface (simplified: pure allocation management)
/// @notice Deposit/withdraw moved to StakeNFT. StakingVault manages allocations only.
///         Allocations are immediate — no pending/dual-slot mechanism.
interface IStakingVault {
    // ── Setup (AWPRegistry only, one-time) ──
    function setStakeNFT(address stakeNFT_) external;

    // ── Write (AWPRegistry only) ──
    function allocate(address user, address agent, uint256 subnetId, uint256 amount) external;
    function deallocate(address user, address agent, uint256 subnetId, uint256 amount) external;
    function reallocate(
        address user, address fromAgent, uint256 fromSubnetId,
        address toAgent, uint256 toSubnetId, uint256 amount
    ) external;
    /// @notice Freeze all allocations of agent — no subnet list needed
    function freezeAgentAllocations(address user, address agent) external;

    // ── Query ──
    function userTotalAllocated(address user) external view returns (uint256);
    function subnetTotalStake(uint256 subnetId) external view returns (uint256);
    function getAgentStake(address user, address agent, uint256 subnetId) external view returns (uint256);
    function getAgentSubnets(address user, address agent) external view returns (uint256[] memory);
    function getSubnetTotalStake(uint256 subnetId) external view returns (uint256);
}
