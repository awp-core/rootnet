// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {EIP712Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {IStakeNFT} from "../interfaces/IStakeNFT.sol";

/// @dev Minimal interface for reading delegate auth from AWPRegistry
interface IAWPRegistryDelegates {
    function delegates(address staker, address delegate) external view returns (bool);
}

/// @title StakingVault — Allocation management with EIP-712 gasless support (UUPS proxy)
/// @notice Manages user stake allocations to (agent, subnetId) triples.
///         Deposit/withdraw have been moved to StakeNFT. All allocations are immediate.
///         Tracks which subnets each (user, agent) pair has allocations on, enabling
///         automatic complete freeze without requiring the caller to supply subnet IDs.
contract StakingVault is Initializable, UUPSUpgradeable, EIP712Upgradeable {
    using EnumerableSet for EnumerableSet.UintSet;

    /// @notice AWPRegistry contract address (storage, set at initialize for proxy pattern)
    address public awpRegistry;

    /// @notice Treasury/Timelock address (locally stored — upgrade auth does not depend on external calls)
    address public treasury;

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

    // ═══════════════════════════════════════════════
    // ── Gasless — EIP-712 ──
    // ═══════════════════════════════════════════════

    /// @notice Per-signer nonce for replay attack prevention
    mapping(address => uint256) public nonces;

    /// @dev EIP-712 type hash: Allocate(address staker, address agent, uint256 subnetId, uint256 amount, uint256 nonce, uint256 deadline)
    bytes32 private constant ALLOCATE_TYPEHASH =
        keccak256("Allocate(address staker,address agent,uint256 subnetId,uint256 amount,uint256 nonce,uint256 deadline)");

    /// @dev EIP-712 type hash: Deallocate(address staker, address agent, uint256 subnetId, uint256 amount, uint256 nonce, uint256 deadline)
    bytes32 private constant DEALLOCATE_TYPEHASH =
        keccak256("Deallocate(address staker,address agent,uint256 subnetId,uint256 amount,uint256 nonce,uint256 deadline)");

    // ── Errors ──

    error NotAWPRegistry();
    error NotAuthorized();
    error NotTimelock();
    error InsufficientUnallocated();
    error InsufficientAllocation();
    error InvalidAmount();
    error ZeroAddress();
    error AlreadySet();
    error ExpiredSignature();
    error InvalidSignature();

    // ── Events ──

    event Allocated(address indexed staker, address indexed agent, uint256 subnetId, uint256 amount, address operator);
    event Deallocated(address indexed staker, address indexed agent, uint256 subnetId, uint256 amount, address operator);
    event Reallocated(
        address indexed staker,
        address fromAgent,
        uint256 fromSubnetId,
        address toAgent,
        uint256 toSubnetId,
        uint256 amount,
        address operator
    );
    event AgentAllocationsFrozen(address indexed user, address indexed agent, uint256 totalFrozen);

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
    /// @param treasury_ Treasury/Timelock address (stored locally for upgrade auth)
    function initialize(address awpRegistry_, address treasury_) external initializer {
        __UUPSUpgradeable_init();
        __EIP712_init("StakingVault", "1");
        awpRegistry = awpRegistry_;
        treasury = treasury_;
    }

    /// @dev UUPS upgrade authorization — only Treasury/Timelock may upgrade
    function _authorizeUpgrade(address) internal view override {
        if (msg.sender != treasury) revert NotTimelock();
    }

    /// @notice Set StakeNFT address (one-time, resolves CREATE2 circular dependency)
    /// @param stakeNFT_ StakeNFT contract address
    function setStakeNFT(address stakeNFT_) external {
        if (msg.sender != awpRegistry) revert NotAWPRegistry();
        if (stakeNFT != address(0)) revert AlreadySet();
        if (stakeNFT_ == address(0)) revert ZeroAddress();
        stakeNFT = stakeNFT_;
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

    /// @notice Allocate staking to a specific (agent, subnetId), effective immediately
    /// @param staker Staker address (caller must be staker or delegate)
    /// @param agent Agent address
    /// @param subnetId Target subnet ID
    /// @param amount Allocation amount
    function allocate(address staker, address agent, uint256 subnetId, uint256 amount) external {
        if (!_isAuthorized(staker, msg.sender)) revert NotAuthorized();
        _allocate(staker, agent, subnetId, amount);
        emit Allocated(staker, agent, subnetId, amount, msg.sender);
    }

    /// @notice Gasless allocate: relayer pays gas, staker signs EIP-712
    function allocateFor(
        address staker, address agent, uint256 subnetId, uint256 amount, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external {
        _verifyDigest(staker, keccak256(abi.encode(ALLOCATE_TYPEHASH, staker, agent, subnetId, amount, nonces[staker]++, deadline)), deadline, v, r, s);
        _allocate(staker, agent, subnetId, amount);
        emit Allocated(staker, agent, subnetId, amount, msg.sender);
    }

    /// @notice Deallocate staking from a specific (agent, subnetId), effective immediately
    /// @param staker Staker address (caller must be staker or delegate)
    /// @param agent Agent address
    /// @param subnetId Target subnet ID
    /// @param amount Amount to deallocate
    function deallocate(address staker, address agent, uint256 subnetId, uint256 amount) external {
        if (!_isAuthorized(staker, msg.sender)) revert NotAuthorized();
        _deallocate(staker, agent, subnetId, amount);
        emit Deallocated(staker, agent, subnetId, amount, msg.sender);
    }

    /// @notice Gasless deallocate: relayer pays gas, staker signs EIP-712
    function deallocateFor(
        address staker, address agent, uint256 subnetId, uint256 amount, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external {
        _verifyDigest(staker, keccak256(abi.encode(DEALLOCATE_TYPEHASH, staker, agent, subnetId, amount, nonces[staker]++, deadline)), deadline, v, r, s);
        _deallocate(staker, agent, subnetId, amount);
        emit Deallocated(staker, agent, subnetId, amount, msg.sender);
    }

    /// @notice Reallocate staking: immediate atomic move from one (agent, subnet) to another
    /// @param staker Staker address (caller must be staker or delegate)
    /// @param fromAgent Source Agent address
    /// @param fromSubnetId Source subnet ID
    /// @param toAgent Target Agent address
    /// @param toSubnetId Target subnet ID
    /// @param amount Migration amount
    function reallocate(
        address staker,
        address fromAgent,
        uint256 fromSubnetId,
        address toAgent,
        uint256 toSubnetId,
        uint256 amount
    ) external {
        if (!_isAuthorized(staker, msg.sender)) revert NotAuthorized();
        if (amount == 0 || amount > type(uint128).max || fromSubnetId == 0 || toSubnetId == 0) revert InvalidAmount();

        uint128 amt128 = uint128(amount);
        uint128 current = _allocations[staker][fromAgent][fromSubnetId];
        if (current < amt128) revert InsufficientAllocation();

        // Subtract from source
        uint128 remaining = current - amt128;
        _allocations[staker][fromAgent][fromSubnetId] = remaining;
        if (remaining == 0) {
            _agentSubnets[staker][fromAgent].remove(fromSubnetId);
        }

        // Add to destination
        _allocations[staker][toAgent][toSubnetId] += amt128;
        _agentSubnets[staker][toAgent].add(toSubnetId);

        // Update subnet totals
        subnetTotalStake[fromSubnetId] -= amount;
        subnetTotalStake[toSubnetId] += amount;

        // userTotalAllocated unchanged (it's a move)

        emit Reallocated(staker, fromAgent, fromSubnetId, toAgent, toSubnetId, amount, msg.sender);
    }

    /// @notice Freeze all allocations of an Agent across all subnets it has been allocated to
    /// @dev Automatically enumerates subnets via the internal _agentSubnets set — no caller-supplied list needed.
    ///      Zeroes all allocations, updates subnet totals, and releases userTotalAllocated.
    /// @param user User (Principal) address
    /// @param agent Agent address to freeze
    function freezeAgentAllocations(address user, address agent) external onlyAWPRegistry {
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
    // ── Internal allocation helpers ──
    // ═══════════════════════════════════════════════

    /// @dev Internal allocate logic (shared by allocate and allocateFor)
    function _allocate(address staker, address agent, uint256 subnetId, uint256 amount) internal {
        if (amount == 0 || amount > type(uint128).max || subnetId == 0) revert InvalidAmount();

        // Check available balance via StakeNFT (O(1))
        uint256 staked = IStakeNFT(stakeNFT).getUserTotalStaked(staker);
        uint256 allocated = userTotalAllocated[staker];
        if (allocated + amount > staked) revert InsufficientUnallocated();

        _allocations[staker][agent][subnetId] += uint128(amount);
        userTotalAllocated[staker] += amount;
        subnetTotalStake[subnetId] += amount;
        // Track this subnet for the (staker, agent) pair
        _agentSubnets[staker][agent].add(subnetId);
    }

    /// @dev Internal deallocate logic (shared by deallocate and deallocateFor)
    function _deallocate(address staker, address agent, uint256 subnetId, uint256 amount) internal {
        if (subnetId == 0) revert InvalidAmount();
        if (amount == 0 || amount > type(uint128).max) revert InvalidAmount();

        uint128 amt128 = uint128(amount);
        uint128 current = _allocations[staker][agent][subnetId];
        if (current < amt128) revert InsufficientAllocation();

        uint128 remaining = current - amt128;
        _allocations[staker][agent][subnetId] = remaining;
        userTotalAllocated[staker] -= amount;
        subnetTotalStake[subnetId] -= amount;
        // Remove subnet from set when fully deallocated
        if (remaining == 0) {
            _agentSubnets[staker][agent].remove(subnetId);
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

    /// @dev Reserved storage gap for future upgrades (UUPS pattern)
    uint256[45] private __gap;
}
