// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IAWPAllocator — Allocation management interface (UUPS proxy with EIP-712 gasless support)
/// @notice Pure bookkeeping — holds no tokens. Deposit/withdraw are in veAWP.
interface IAWPAllocator {
    // ── Write ──
    function allocate(address staker, address agent, uint256 worknetId, uint256 amount) external;
    function deallocate(address staker, address agent, uint256 worknetId, uint256 amount) external;
    function deallocateAll(address staker, address agent, uint256 worknetId) external;
    function reallocate(
        address staker, address fromAgent, uint256 fromWorknetId,
        address toAgent, uint256 toWorknetId, uint256 amount
    ) external;

    // ── Batch ──
    function batchAllocate(address staker, address[] calldata agents, uint256[] calldata worknetIds, uint256[] calldata amounts) external;
    function batchDeallocate(address staker, address[] calldata agents, uint256[] calldata worknetIds, uint256[] calldata amounts) external;

    // ── Gasless (EIP-712) ──
    function allocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) external;
    function deallocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) external;

    // ── Events ──
    event Allocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator);
    event Deallocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator);
    event Reallocated(address indexed staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount, address operator);

    // ── Query ──
    function nonces(address user) external view returns (uint256);
    function userTotalAllocated(address user) external view returns (uint256);
    function worknetTotalStake(uint256 worknetId) external view returns (uint256);
    function getAgentStake(address user, address agent, uint256 worknetId) external view returns (uint256);
    function getAgentWorknets(address user, address agent) external view returns (uint256[] memory);
    function getWorknetTotalStake(uint256 worknetId) external view returns (uint256);
}
