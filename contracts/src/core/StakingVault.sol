// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {IStakeNFT} from "../interfaces/IStakeNFT.sol";

/// @title StakingVault — Pure allocation management (simplified)
/// @notice Manages user stake allocations to (agent, subnetId) triples.
///         Deposit/withdraw have been moved to StakeNFT. All allocations are immediate.
///         Tracks which subnets each (user, agent) pair has allocations on, enabling
///         automatic complete freeze without requiring the caller to supply subnet IDs.
contract StakingVault {
    using EnumerableSet for EnumerableSet.UintSet;

    /// @notice RootNet contract address (immutable)
    address public immutable rootNet;

    /// @notice StakeNFT contract address (for balance checks, set once via setStakeNFT)
    address public stakeNFT;

    // ═══════════════════════════════════════════════
    // ── Allocation storage ──
    // ═══════════════════════════════════════════════

    /// @dev Allocation mapping: user => agent => subnetId => amount (plain uint128)
    mapping(address => mapping(address => mapping(uint256 => uint128))) private _allocations;

    /// @dev Tracks all subnet IDs where (user, agent) has a non-zero allocation.
    ///      Maintained in sync with _allocations: added on first allocate, removed when zeroed.
    mapping(address => mapping(address => EnumerableSet.UintSet)) private _agentSubnets;

    /// @notice Total allocation per user (increased on allocate, decreased on deallocate/freeze)
    mapping(address => uint256) public userTotalAllocated;

    /// @notice Total stake per subnet
    mapping(uint256 => uint256) public subnetTotalStake;

    // ── Errors ──

    error NotRootNet();
    error InsufficientUnallocated();
    error InsufficientAllocation();
    error InvalidAmount();

    event AgentAllocationsFrozen(address indexed user, address indexed agent, uint256 totalFrozen);

    /// @dev Only the RootNet contract may call
    modifier onlyRootNet() {
        if (msg.sender != rootNet) revert NotRootNet();
        _;
    }

    error ZeroAddress();
    error AlreadySet();

    /// @notice Constructor
    /// @param rootNet_ RootNet contract address
    constructor(address rootNet_) {
        rootNet = rootNet_;
    }

    /// @notice Set StakeNFT address (one-time, resolves CREATE2 circular dependency)
    /// @param stakeNFT_ StakeNFT contract address
    function setStakeNFT(address stakeNFT_) external {
        if (msg.sender != rootNet) revert NotRootNet();
        if (stakeNFT != address(0)) revert AlreadySet();
        if (stakeNFT_ == address(0)) revert ZeroAddress();
        stakeNFT = stakeNFT_;
    }

    // ═══════════════════════════════════════════════
    // ── Allocation (all immediate) ──
    // ═══════════════════════════════════════════════

    /// @notice Allocate staking to a specific (agent, subnetId), effective immediately
    /// @param user User address
    /// @param agent Agent address
    /// @param subnetId Target subnet ID
    /// @param amount Allocation amount
    function allocate(address user, address agent, uint256 subnetId, uint256 amount)
        external
        onlyRootNet
    {
        if (amount == 0 || amount > type(uint128).max || subnetId == 0) revert InvalidAmount();

        // Check available balance via StakeNFT (O(1))
        uint256 staked = IStakeNFT(stakeNFT).getUserTotalStaked(user);
        uint256 allocated = userTotalAllocated[user];
        if (allocated + amount > staked) revert InsufficientUnallocated();

        _allocations[user][agent][subnetId] += uint128(amount);
        userTotalAllocated[user] += amount;
        subnetTotalStake[subnetId] += amount;
        // Track this subnet for the (user, agent) pair
        _agentSubnets[user][agent].add(subnetId);
    }

    /// @notice Deallocate staking from a specific (agent, subnetId), effective immediately
    /// @param user User address
    /// @param agent Agent address
    /// @param subnetId Target subnet ID
    /// @param amount Amount to deallocate
    function deallocate(address user, address agent, uint256 subnetId, uint256 amount)
        external
        onlyRootNet
    {
        if (amount == 0 || amount > type(uint128).max) revert InvalidAmount();

        uint128 amt128 = uint128(amount);
        uint128 current = _allocations[user][agent][subnetId];
        if (current < amt128) revert InsufficientAllocation();

        uint128 remaining = current - amt128;
        _allocations[user][agent][subnetId] = remaining;
        userTotalAllocated[user] -= amount;
        subnetTotalStake[subnetId] -= amount;
        // Remove subnet from set when fully deallocated
        if (remaining == 0) {
            _agentSubnets[user][agent].remove(subnetId);
        }
    }

    /// @notice Reallocate staking: immediate atomic move from one (agent, subnet) to another
    /// @param user User address
    /// @param fromAgent Source Agent address
    /// @param fromSubnetId Source subnet ID
    /// @param toAgent Target Agent address
    /// @param toSubnetId Target subnet ID
    /// @param amount Migration amount
    function reallocate(
        address user,
        address fromAgent,
        uint256 fromSubnetId,
        address toAgent,
        uint256 toSubnetId,
        uint256 amount
    ) external onlyRootNet {
        if (amount == 0 || amount > type(uint128).max || fromSubnetId == 0 || toSubnetId == 0) revert InvalidAmount();

        uint128 amt128 = uint128(amount);
        uint128 current = _allocations[user][fromAgent][fromSubnetId];
        if (current < amt128) revert InsufficientAllocation();

        // Subtract from source
        uint128 remaining = current - amt128;
        _allocations[user][fromAgent][fromSubnetId] = remaining;
        if (remaining == 0) {
            _agentSubnets[user][fromAgent].remove(fromSubnetId);
        }

        // Add to destination
        _allocations[user][toAgent][toSubnetId] += amt128;
        _agentSubnets[user][toAgent].add(toSubnetId);

        // Update subnet totals
        subnetTotalStake[fromSubnetId] -= amount;
        subnetTotalStake[toSubnetId] += amount;

        // userTotalAllocated unchanged (it's a move)
    }

    /// @notice Freeze all allocations of an Agent across all subnets it has been allocated to
    /// @dev Automatically enumerates subnets via the internal _agentSubnets set — no caller-supplied list needed.
    ///      Zeroes all allocations, updates subnet totals, and releases userTotalAllocated.
    /// @param user User (Principal) address
    /// @param agent Agent address to freeze
    function freezeAgentAllocations(address user, address agent) external onlyRootNet {
        EnumerableSet.UintSet storage subnets = _agentSubnets[user][agent];
        uint256 count = subnets.length();
        if (count == 0) return;

        uint256 totalFrozen = 0;
        // Single reverse-iterating loop: zero allocation + remove from set atomically.
        // Reverse iteration is safe because EnumerableSet.remove swaps with the last
        // element — removing index i never invalidates indices < i.
        for (uint256 i = count; i > 0;) {
            unchecked { --i; }
            uint256 sid = subnets.at(i);
            uint256 amt = _allocations[user][agent][sid];
            if (amt > 0) {
                _allocations[user][agent][sid] = 0;
                subnetTotalStake[sid] -= amt;
                totalFrozen += amt;
            }
            subnets.remove(sid);
        }

        userTotalAllocated[user] -= totalFrozen;

        if (totalFrozen > 0) {
            emit AgentAllocationsFrozen(user, agent, totalFrozen);
        }
    }

    // ═══════════════════════════════════════════════
    // ── Query functions ──
    // ═══════════════════════════════════════════════

    /// @notice Get the current effective allocation for a specific (user, agent, subnetId)
    /// @param user User address
    /// @param agent Agent address
    /// @param subnetId Subnet ID
    /// @return Current allocation amount
    function getAgentStake(address user, address agent, uint256 subnetId)
        external
        view
        returns (uint256)
    {
        return _allocations[user][agent][subnetId];
    }

    /// @notice Get all subnet IDs where (user, agent) has active allocations
    /// @param user User address
    /// @param agent Agent address
    /// @return Array of subnet IDs with non-zero allocations
    function getAgentSubnets(address user, address agent) external view returns (uint256[] memory) {
        return _agentSubnets[user][agent].values();
    }

    /// @notice Get the total stake for a subnet
    /// @param subnetId Subnet ID
    /// @return Total stake for the subnet
    function getSubnetTotalStake(uint256 subnetId) external view returns (uint256) {
        return subnetTotalStake[subnetId];
    }
}
