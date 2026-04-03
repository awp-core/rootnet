// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {EIP712Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {IStakeNFT} from "../interfaces/IStakeNFT.sol";

/// @dev Minimal interface for reading delegate auth from AWPRegistry
interface IAWPRegistryDelegates {
    function delegates(address staker, address delegate) external view returns (bool);
}

/// @title StakingVault — Allocation management with EIP-712 gasless support (UUPS proxy)
/// @notice Manages user stake allocations to (agent, worknetId) triples.
///         Deposit/withdraw have been moved to StakeNFT. All allocations are immediate.
///         Tracks which worknets each (user, agent) pair has allocations on, enabling
///         automatic complete freeze without requiring the caller to supply worknet IDs.
contract StakingVault is Initializable, UUPSUpgradeable, ReentrancyGuardUpgradeable, EIP712Upgradeable {
    using EnumerableSet for EnumerableSet.UintSet;

    /// @notice AWPRegistry contract address (storage, set at initialize for proxy pattern)
    address public awpRegistry;

    /// @notice Guardian address (cross-chain multisig — upgrade auth, locally stored)
    address public guardian;

    /// @notice StakeNFT contract address (for balance checks, set once via setStakeNFT)
    address public stakeNFT;

    // ═══════════════════════════════════════════════
    // ── Allocation storage ──
    // ═══════════════════════════════════════════════

    /// @dev Allocation mapping: user => agent => worknetId => amount (plain uint128)
    mapping(address => mapping(address => mapping(uint256 => uint128))) private _allocations;

    /// @dev Tracks all worknet IDs where (user, agent) has a non-zero allocation.
    ///      Maintained in sync with _allocations: added on first allocate, removed when zeroed.
    mapping(address => mapping(address => EnumerableSet.UintSet)) private _agentWorknets;

    /// @notice Total allocation per user (increased on allocate, decreased on deallocate/freeze)
    mapping(address => uint256) public userTotalAllocated;

    /// @notice Total stake per worknet
    mapping(uint256 => uint256) public worknetTotalStake;

    // ═══════════════════════════════════════════════
    // ── Gasless — EIP-712 ──
    // ═══════════════════════════════════════════════

    /// @notice Per-signer nonce for replay attack prevention
    mapping(address => uint256) public nonces;

    /// @dev EIP-712 type hash: Allocate(address staker, address agent, uint256 worknetId, uint256 amount, uint256 nonce, uint256 deadline)
    bytes32 private constant ALLOCATE_TYPEHASH =
        keccak256("Allocate(address staker,address agent,uint256 worknetId,uint256 amount,uint256 nonce,uint256 deadline)");

    /// @dev EIP-712 type hash: Deallocate(address staker, address agent, uint256 worknetId, uint256 amount, uint256 nonce, uint256 deadline)
    bytes32 private constant DEALLOCATE_TYPEHASH =
        keccak256("Deallocate(address staker,address agent,uint256 worknetId,uint256 amount,uint256 nonce,uint256 deadline)");

    // ── Errors ──

    error NotAWPRegistry();
    error NotAuthorized();
    error NotGuardian();
    error InsufficientUnallocated();
    error InsufficientAllocation();
    error ZeroAmount();
    error AmountExceedsUint128();
    error ZeroWorknetId();
    error AllocationOverflow();
    error ZeroAddress();
    error AlreadySet();
    error ExpiredSignature();
    error InvalidSignature();

    // ── Events ──

    event Allocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator);
    event Deallocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator);
    event Reallocated(
        address indexed staker,
        address fromAgent,
        uint256 fromWorknetId,
        address toAgent,
        uint256 toWorknetId,
        uint256 amount,
        address operator
    );
    event GuardianUpdated(address indexed newGuardian);
    event StakeNFTSet(address indexed stakeNFT);

    /// @dev Only the AWPRegistry contract may call
    modifier onlyAWPRegistry() {
        if (msg.sender != awpRegistry) revert NotAWPRegistry();
        _;
    }

    // ═══════════════════════════════════════════════
    // ── Constructor / Initialize ──
    // ═══════════════════════════════════════════════

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Initialize the vault (called once via proxy)
    /// @param awpRegistry_ AWPRegistry contract address
    /// @param guardian_ Guardian address (cross-chain multisig — controls upgrades)
    function initialize(address awpRegistry_, address guardian_) external initializer {
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
        __EIP712_init("StakingVault", "1");
        if (awpRegistry_ == address(0) || guardian_ == address(0)) revert ZeroAddress();
        awpRegistry = awpRegistry_;
        guardian = guardian_;
    }

    /// @dev UUPS upgrade authorization — only Guardian may upgrade
    function _authorizeUpgrade(address) internal view override {
        if (msg.sender != guardian) revert NotGuardian();
    }

    /// @notice Update guardian address (only current Guardian may call)
    function setGuardian(address g) external nonReentrant {
        if (msg.sender != guardian) revert NotGuardian();
        if (g == address(0)) revert ZeroAddress();
        guardian = g;
        emit GuardianUpdated(g);
    }

    /// @notice Set StakeNFT address (one-time, resolves CREATE2 circular dependency)
    /// @param stakeNFT_ StakeNFT contract address
    function setStakeNFT(address stakeNFT_) external nonReentrant {
        if (msg.sender != awpRegistry) revert NotAWPRegistry();
        if (stakeNFT != address(0)) revert AlreadySet();
        if (stakeNFT_ == address(0)) revert ZeroAddress();
        stakeNFT = stakeNFT_;
        emit StakeNFTSet(stakeNFT_);
    }

    // ═══════════════════════════════════════════════
    // ── Auth helpers ──
    // ═══════════════════════════════════════════════

    /// @dev Check if caller is authorized to act on behalf of staker
    function _isAuthorized(address staker, address caller) internal view returns (bool) {
        return caller == staker || IAWPRegistryDelegates(awpRegistry).delegates(staker, caller);
    }

    /// @dev Verify EIP-712 digest + deadline. Reverts on failure.
    function _verifyDigest(address user, bytes32 structHash, uint256 deadline, uint8 v, bytes32 r, bytes32 s) internal view {
        if (block.timestamp > deadline) revert ExpiredSignature();
        bytes32 digest = _hashTypedDataV4(structHash);
        if (ECDSA.recover(digest, v, r, s) != user) revert InvalidSignature();
    }

    // ═══════════════════════════════════════════════
    // ── Allocation (all immediate) ──
    // ═══════════════════════════════════════════════

    /// @notice Allocate staking to a specific (agent, worknetId), effective immediately
    /// @param staker Staker address (caller must be staker or delegate)
    /// @param agent Agent address
    /// @param worknetId Target worknet ID
    /// @param amount Allocation amount
    function allocate(address staker, address agent, uint256 worknetId, uint256 amount) external nonReentrant {
        if (!_isAuthorized(staker, msg.sender)) revert NotAuthorized();
        _allocate(staker, agent, worknetId, amount);
        emit Allocated(staker, agent, worknetId, amount, msg.sender);
    }

    /// @notice Gasless allocate: relayer pays gas, staker signs EIP-712
    function allocateFor(
        address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant {
        _verifyDigest(staker, keccak256(abi.encode(ALLOCATE_TYPEHASH, staker, agent, worknetId, amount, nonces[staker]++, deadline)), deadline, v, r, s);
        _allocate(staker, agent, worknetId, amount);
        emit Allocated(staker, agent, worknetId, amount, msg.sender);
    }

    /// @notice Deallocate staking from a specific (agent, worknetId), effective immediately
    /// @param staker Staker address (caller must be staker or delegate)
    /// @param agent Agent address
    /// @param worknetId Target worknet ID
    /// @param amount Amount to deallocate
    function deallocate(address staker, address agent, uint256 worknetId, uint256 amount) external nonReentrant {
        if (!_isAuthorized(staker, msg.sender)) revert NotAuthorized();
        _deallocate(staker, agent, worknetId, amount);
        emit Deallocated(staker, agent, worknetId, amount, msg.sender);
    }

    /// @notice Gasless deallocate: relayer pays gas, staker signs EIP-712
    function deallocateFor(
        address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant {
        _verifyDigest(staker, keccak256(abi.encode(DEALLOCATE_TYPEHASH, staker, agent, worknetId, amount, nonces[staker]++, deadline)), deadline, v, r, s);
        _deallocate(staker, agent, worknetId, amount);
        emit Deallocated(staker, agent, worknetId, amount, msg.sender);
    }

    /// @notice Reallocate staking: immediate atomic move from one (agent, worknet) to another
    /// @param staker Staker address (caller must be staker or delegate)
    /// @param fromAgent Source Agent address
    /// @param fromWorknetId Source worknet ID
    /// @param toAgent Target Agent address
    /// @param toWorknetId Target worknet ID
    /// @param amount Migration amount
    function reallocate(
        address staker,
        address fromAgent,
        uint256 fromWorknetId,
        address toAgent,
        uint256 toWorknetId,
        uint256 amount
    ) external nonReentrant {
        if (!_isAuthorized(staker, msg.sender)) revert NotAuthorized();
        if (amount == 0) revert ZeroAmount();
        if (amount > type(uint128).max) revert AmountExceedsUint128();
        if (fromWorknetId == 0 || toWorknetId == 0) revert ZeroWorknetId();

        uint128 amt128 = uint128(amount);
        uint128 current = _allocations[staker][fromAgent][fromWorknetId];
        if (current < amt128) revert InsufficientAllocation();

        // Subtract from source
        uint128 remaining = current - amt128;
        _allocations[staker][fromAgent][fromWorknetId] = remaining;
        if (remaining == 0) {
            _agentWorknets[staker][fromAgent].remove(fromWorknetId);
        }

        // Add to destination — check uint128 overflow
        uint128 destCurrent = _allocations[staker][toAgent][toWorknetId];
        if (uint256(destCurrent) + amount > type(uint128).max) revert AllocationOverflow();
        _allocations[staker][toAgent][toWorknetId] = destCurrent + amt128;
        _agentWorknets[staker][toAgent].add(toWorknetId);

        // Update worknet totals
        worknetTotalStake[fromWorknetId] -= amount;
        worknetTotalStake[toWorknetId] += amount;

        // userTotalAllocated unchanged (it's a move)

        emit Reallocated(staker, fromAgent, fromWorknetId, toAgent, toWorknetId, amount, msg.sender);
    }


    // ═══════════════════════════════════════════════
    // ── Internal allocation helpers ──
    // ═══════════════════════════════════════════════

    /// @dev Internal allocate logic (shared by allocate and allocateFor)
    function _allocate(address staker, address agent, uint256 worknetId, uint256 amount) internal {
        if (amount == 0) revert ZeroAmount();
        if (amount > type(uint128).max) revert AmountExceedsUint128();
        if (worknetId == 0) revert ZeroWorknetId();

        // Check available balance via StakeNFT (O(1))
        uint256 staked = IStakeNFT(stakeNFT).getUserTotalStaked(staker);
        uint256 allocated = userTotalAllocated[staker];
        if (allocated + amount > staked) revert InsufficientUnallocated();

        // Check uint128 overflow before adding
        uint128 current = _allocations[staker][agent][worknetId];
        if (uint256(current) + amount > type(uint128).max) revert AllocationOverflow();
        _allocations[staker][agent][worknetId] = current + uint128(amount);
        userTotalAllocated[staker] += amount;
        worknetTotalStake[worknetId] += amount;
        // Track this worknet for the (staker, agent) pair
        _agentWorknets[staker][agent].add(worknetId);
    }

    /// @dev Internal deallocate logic (shared by deallocate and deallocateFor)
    function _deallocate(address staker, address agent, uint256 worknetId, uint256 amount) internal {
        if (worknetId == 0) revert ZeroWorknetId();
        if (amount == 0) revert ZeroAmount();
        if (amount > type(uint128).max) revert AmountExceedsUint128();

        uint128 amt128 = uint128(amount);
        uint128 current = _allocations[staker][agent][worknetId];
        if (current < amt128) revert InsufficientAllocation();

        uint128 remaining = current - amt128;
        _allocations[staker][agent][worknetId] = remaining;
        userTotalAllocated[staker] -= amount;
        worknetTotalStake[worknetId] -= amount;
        // Remove worknet from set when fully deallocated
        if (remaining == 0) {
            _agentWorknets[staker][agent].remove(worknetId);
        }
    }

    // ═══════════════════════════════════════════════
    // ── Query functions ──
    // ═══════════════════════════════════════════════

    /// @notice Get the current effective allocation for a specific (user, agent, worknetId)
    /// @param user User address
    /// @param agent Agent address
    /// @param worknetId Worknet ID
    /// @return Current allocation amount
    function getAgentStake(address user, address agent, uint256 worknetId)
        external
        view
        returns (uint256)
    {
        return _allocations[user][agent][worknetId];
    }

    /// @notice Get all worknet IDs where (user, agent) has active allocations
    /// @param user User address
    /// @param agent Agent address
    /// @return Array of worknet IDs with non-zero allocations
    function getAgentWorknets(address user, address agent) external view returns (uint256[] memory) {
        return _agentWorknets[user][agent].values();
    }

    /// @notice Get the total stake for a worknet
    /// @param worknetId Worknet ID
    /// @return Total stake for the worknet
    function getWorknetTotalStake(uint256 worknetId) external view returns (uint256) {
        return worknetTotalStake[worknetId];
    }

    /// @dev Reserved storage gap for future upgrades (UUPS pattern)
    uint256[45] private __gap;
}
