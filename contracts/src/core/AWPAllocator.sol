// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {EIP712Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {IveAWP} from "../interfaces/IveAWP.sol";

/// @dev Minimal interface for reading delegate auth from AWPRegistry
interface IAWPRegistryDelegates {
    function delegates(address staker, address delegate) external view returns (bool);
}

/// @title AWPAllocator — Allocation management with EIP-712 gasless support (UUPS proxy)
/// @notice Manages user stake allocations to (agent, worknetId) triples. Pure bookkeeping — holds no tokens.
///         Deposit/withdraw are in veAWP. All allocations are immediate.
contract AWPAllocator is UUPSUpgradeable, EIP712Upgradeable {
    using EnumerableSet for EnumerableSet.UintSet;

    // ── Immutables ──

    /// @notice AWPRegistry contract address (for delegate auth checks)
    address public immutable awpRegistry;

    /// @notice veAWP contract address (for staked balance checks)
    address public immutable veAWP;

    // ── Storage ──

    /// @notice Guardian address (cross-chain multisig — upgrade auth)
    address public guardian;

    /// @dev Allocation mapping: user => agent => worknetId => amount
    mapping(address => mapping(address => mapping(uint256 => uint256))) private _allocations;

    /// @dev Tracks all worknet IDs where (user, agent) has a non-zero allocation
    mapping(address => mapping(address => EnumerableSet.UintSet)) private _agentWorknets;

    /// @notice Total allocation per user
    mapping(address => uint256) public userTotalAllocated;

    /// @notice Total stake per worknet
    mapping(uint256 => uint256) public worknetTotalStake;

    /// @notice Per-signer nonce for EIP-712 replay attack prevention
    mapping(address => uint256) public nonces;

    /// @dev Reserved storage gap for future upgrades
    uint256[43] private __gap;

    // ── Constants ──

    bytes32 private constant ALLOCATE_TYPEHASH =
        keccak256("Allocate(address staker,address agent,uint256 worknetId,uint256 amount,uint256 nonce,uint256 deadline)");

    bytes32 private constant DEALLOCATE_TYPEHASH =
        keccak256("Deallocate(address staker,address agent,uint256 worknetId,uint256 amount,uint256 nonce,uint256 deadline)");

    // ── Errors ──

    error NotAuthorized();
    error NotGuardian();
    error InsufficientUnallocated();
    error InsufficientAllocation();
    error ZeroAmount();
    error ZeroWorknetId();
    error ZeroAddress();
    error ArrayLengthMismatch();
    error ExpiredSignature();
    error InvalidSignature();

    // ── Events ──

    event Allocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator);
    event Deallocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator);
    event Reallocated(
        address indexed staker,
        address fromAgent, uint256 fromWorknetId,
        address toAgent, uint256 toWorknetId,
        uint256 amount, address operator
    );
    event GuardianUpdated(address indexed newGuardian);

    // ═══════════════════════════════════════════════
    // ── Constructor / Initialize ──
    // ═══════════════════════════════════════════════

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor(address awpRegistry_, address veAWP_) {
        awpRegistry = awpRegistry_;
        veAWP = veAWP_;
        _disableInitializers();
    }

    /// @notice Initialize the allocator (called once via proxy)
    /// @param awpRegistry_ Unused (kept for ABI compatibility; actual value is immutable from constructor)
    /// @param guardian_ Guardian address
    function initialize(address awpRegistry_, address guardian_) external initializer {
        __EIP712_init("AWPAllocator", "1");
        if (guardian_ == address(0)) revert ZeroAddress();
        guardian = guardian_;
    }

    /// @dev UUPS upgrade authorization — only Guardian may upgrade
    function _authorizeUpgrade(address) internal view override {
        if (msg.sender != guardian) revert NotGuardian();
    }

    /// @notice Update guardian address (only current Guardian may call)
    function setGuardian(address g) external {
        if (msg.sender != guardian) revert NotGuardian();
        if (g == address(0)) revert ZeroAddress();
        guardian = g;
        emit GuardianUpdated(g);
    }

    // ═══════════════════════════════════════════════
    // ── Single allocation ──
    // ═══════════════════════════════════════════════

    /// @notice Allocate staking to a specific (agent, worknetId)
    function allocate(address staker, address agent, uint256 worknetId, uint256 amount) external {
        _requireAuth(staker);
        _allocate(staker, agent, worknetId, amount);
        emit Allocated(staker, agent, worknetId, amount, msg.sender);
    }

    /// @notice Gasless allocate: relayer pays gas, staker signs EIP-712
    function allocateFor(
        address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external {
        _verifyDigest(staker, keccak256(abi.encode(ALLOCATE_TYPEHASH, staker, agent, worknetId, amount, nonces[staker]++, deadline)), deadline, v, r, s);
        _allocate(staker, agent, worknetId, amount);
        emit Allocated(staker, agent, worknetId, amount, msg.sender);
    }

    /// @notice Deallocate staking from a specific (agent, worknetId)
    function deallocate(address staker, address agent, uint256 worknetId, uint256 amount) external {
        _requireAuth(staker);
        _deallocate(staker, agent, worknetId, amount);
        emit Deallocated(staker, agent, worknetId, amount, msg.sender);
    }

    /// @notice Deallocate all from a specific (agent, worknetId) — convenience, no need to query amount first
    function deallocateAll(address staker, address agent, uint256 worknetId) external {
        _requireAuth(staker);
        uint256 amount = _allocations[staker][agent][worknetId];
        if (amount == 0) revert ZeroAmount();
        _deallocate(staker, agent, worknetId, amount);
        emit Deallocated(staker, agent, worknetId, amount, msg.sender);
    }

    /// @notice Gasless deallocate: relayer pays gas, staker signs EIP-712
    function deallocateFor(
        address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external {
        _verifyDigest(staker, keccak256(abi.encode(DEALLOCATE_TYPEHASH, staker, agent, worknetId, amount, nonces[staker]++, deadline)), deadline, v, r, s);
        _deallocate(staker, agent, worknetId, amount);
        emit Deallocated(staker, agent, worknetId, amount, msg.sender);
    }

    /// @notice Reallocate: atomic move from one (agent, worknet) to another
    function reallocate(
        address staker,
        address fromAgent, uint256 fromWorknetId,
        address toAgent, uint256 toWorknetId,
        uint256 amount
    ) external {
        _requireAuth(staker);
        if (amount == 0) revert ZeroAmount();
        if (fromWorknetId == 0 || toWorknetId == 0) revert ZeroWorknetId();

        // Subtract from source
        uint256 current = _allocations[staker][fromAgent][fromWorknetId];
        if (current < amount) revert InsufficientAllocation();
        uint256 remaining = current - amount;
        _allocations[staker][fromAgent][fromWorknetId] = remaining;
        if (remaining == 0) {
            _agentWorknets[staker][fromAgent].remove(fromWorknetId);
        }

        // Add to destination
        _allocations[staker][toAgent][toWorknetId] += amount;
        _agentWorknets[staker][toAgent].add(toWorknetId);

        // Update worknet totals (userTotalAllocated unchanged — it's a move)
        worknetTotalStake[fromWorknetId] -= amount;
        worknetTotalStake[toWorknetId] += amount;

        emit Reallocated(staker, fromAgent, fromWorknetId, toAgent, toWorknetId, amount, msg.sender);
    }

    // ═══════════════════════════════════════════════
    // ── Batch operations ──
    // ═══════════════════════════════════════════════

    /// @notice Allocate to multiple (agent, worknetId) pairs in a single transaction
    /// @param staker Staker address (caller must be staker or delegate)
    /// @param agents Array of agent addresses
    /// @param worknetIds Array of worknet IDs
    /// @param amounts Array of allocation amounts
    function batchAllocate(
        address staker,
        address[] calldata agents,
        uint256[] calldata worknetIds,
        uint256[] calldata amounts
    ) external {
        _requireAuth(staker);
        uint256 len = agents.length;
        if (len != worknetIds.length || len != amounts.length) revert ArrayLengthMismatch(); // length mismatch
        for (uint256 i = 0; i < len;) {
            _allocate(staker, agents[i], worknetIds[i], amounts[i]);
            emit Allocated(staker, agents[i], worknetIds[i], amounts[i], msg.sender);
            unchecked { ++i; }
        }
    }

    /// @notice Deallocate from multiple (agent, worknetId) pairs in a single transaction
    /// @param staker Staker address (caller must be staker or delegate)
    /// @param agents Array of agent addresses
    /// @param worknetIds Array of worknet IDs
    /// @param amounts Array of deallocation amounts (0 = deallocate all for that entry)
    function batchDeallocate(
        address staker,
        address[] calldata agents,
        uint256[] calldata worknetIds,
        uint256[] calldata amounts
    ) external {
        _requireAuth(staker);
        uint256 len = agents.length;
        if (len != worknetIds.length || len != amounts.length) revert ArrayLengthMismatch();
        for (uint256 i = 0; i < len;) {
            uint256 amt = amounts[i];
            if (amt == 0) {
                // deallocate all
                amt = _allocations[staker][agents[i]][worknetIds[i]];
                if (amt == 0) { unchecked { ++i; } continue; } // skip zero allocations
            }
            _deallocate(staker, agents[i], worknetIds[i], amt);
            emit Deallocated(staker, agents[i], worknetIds[i], amt, msg.sender);
            unchecked { ++i; }
        }
    }

    // ═══════════════════════════════════════════════
    // ── Internal ──
    // ═══════════════════════════════════════════════

    function _requireAuth(address staker) internal view {
        if (msg.sender != staker && !IAWPRegistryDelegates(awpRegistry).delegates(staker, msg.sender)) {
            revert NotAuthorized();
        }
    }

    function _verifyDigest(address user, bytes32 structHash, uint256 deadline, uint8 v, bytes32 r, bytes32 s) internal view {
        if (block.timestamp > deadline) revert ExpiredSignature();
        bytes32 digest = _hashTypedDataV4(structHash);
        if (ECDSA.recover(digest, v, r, s) != user) revert InvalidSignature();
    }

    function _allocate(address staker, address agent, uint256 worknetId, uint256 amount) internal {
        if (amount == 0) revert ZeroAmount();
        if (worknetId == 0) revert ZeroWorknetId();

        // Check available balance via veAWP (O(1))
        uint256 allocated = userTotalAllocated[staker];
        if (allocated + amount > IveAWP(veAWP).getUserTotalStaked(staker)) revert InsufficientUnallocated();

        _allocations[staker][agent][worknetId] += amount;
        userTotalAllocated[staker] = allocated + amount;
        worknetTotalStake[worknetId] += amount;
        _agentWorknets[staker][agent].add(worknetId);
    }

    function _deallocate(address staker, address agent, uint256 worknetId, uint256 amount) internal {
        if (worknetId == 0) revert ZeroWorknetId();
        if (amount == 0) revert ZeroAmount();

        uint256 current = _allocations[staker][agent][worknetId];
        if (current < amount) revert InsufficientAllocation();

        uint256 remaining = current - amount;
        _allocations[staker][agent][worknetId] = remaining;
        userTotalAllocated[staker] -= amount;
        worknetTotalStake[worknetId] -= amount;
        if (remaining == 0) {
            _agentWorknets[staker][agent].remove(worknetId);
        }
    }

    // ═══════════════════════════════════════════════
    // ── Query ──
    // ═══════════════════════════════════════════════

    function getAgentStake(address user, address agent, uint256 worknetId) external view returns (uint256) {
        return _allocations[user][agent][worknetId];
    }

    function getAgentWorknets(address user, address agent) external view returns (uint256[] memory) {
        return _agentWorknets[user][agent].values();
    }

    function getWorknetTotalStake(uint256 worknetId) external view returns (uint256) {
        return worknetTotalStake[worknetId];
    }
}
