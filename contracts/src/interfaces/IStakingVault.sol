// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IStakingVault — Staking vault interface (UUPS proxy with EIP-712 gasless support)
/// @notice Deposit/withdraw moved to StakeNFT. StakingVault manages allocations only.
///         Allocations are immediate — no pending/dual-slot mechanism.
interface IStakingVault {
    // ── Setup ──
    function initialize(address awpRegistry_, address guardian_) external;
    function setStakeNFT(address stakeNFT_) external;

    // ── Write (public with delegate auth) ──
    function allocate(address staker, address agent, uint256 subnetId, uint256 amount) external;
    function deallocate(address staker, address agent, uint256 subnetId, uint256 amount) external;
    function reallocate(
        address staker, address fromAgent, uint256 fromSubnetId,
        address toAgent, uint256 toSubnetId, uint256 amount
    ) external;

    // ── Gasless (EIP-712) ──
    function allocateFor(
        address staker, address agent, uint256 subnetId, uint256 amount, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external;
    function deallocateFor(
        address staker, address agent, uint256 subnetId, uint256 amount, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external;

    // ── Write (AWPRegistry only) ──
    /// @notice Freeze all allocations of agent — no subnet list needed
    function freezeAgentAllocations(address user, address agent) external;

    // ── Events ──
    event Allocated(address indexed staker, address indexed agent, uint256 subnetId, uint256 amount, address operator);
    event Deallocated(address indexed staker, address indexed agent, uint256 subnetId, uint256 amount, address operator);
    event AgentAllocationsFrozen(address indexed staker, address indexed agent, uint256 totalFrozen);
    event Reallocated(
        address indexed staker,
        address fromAgent,
        uint256 fromSubnetId,
        address toAgent,
        uint256 toSubnetId,
        uint256 amount,
        address operator
    );

    // ── Query ──
    function nonces(address user) external view returns (uint256);
    function userTotalAllocated(address user) external view returns (uint256);
    function subnetTotalStake(uint256 subnetId) external view returns (uint256);
    function getAgentStake(address user, address agent, uint256 subnetId) external view returns (uint256);
    function getAgentSubnets(address user, address agent) external view returns (uint256[] memory);
    function getSubnetTotalStake(uint256 subnetId) external view returns (uint256);
}
