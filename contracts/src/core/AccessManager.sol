// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/// @title AccessManager — user registration + Agent permission management
/// @notice Only RootNet may call write functions; manages user/Agent identity registration, binding relationships, and role permissions
/// @dev Core design:
///   - User and Agent addresses are mutually exclusive (the same address cannot be both a User and an Agent)
///   - No on-chain cap on the number of Agents per user
///   - resolveCallerRole returns (owner, isUser, isManager) in a single call, reducing external call overhead
contract AccessManager {
    using EnumerableSet for EnumerableSet.AddressSet;

    /// @notice RootNet contract address (immutable, set at deployment)
    address public immutable rootNet;

    // ── User registration state ──

    /// @notice Whether an address is registered as a user
    mapping(address => bool) public isRegistered;
    /// @notice Block timestamp when the user registered
    mapping(address => uint64) public registeredAt;
    /// @notice Total number of registered users
    uint256 public totalUsers;

    // ── Agent management state ──

    /// @notice Agent address → owner user address (address(0) means not registered as an Agent)
    mapping(address => address) public agentOwner;
    /// @dev User address → set of Agents belonging to that user (EnumerableSet supports enumeration)
    mapping(address => EnumerableSet.AddressSet) private userAgents;
    /// @notice Whether an Agent has Manager privileges (Managers may perform advanced operations on behalf of users)
    mapping(address => bool) public isManager;
    /// @notice User-defined reward recipient address (address(0) means use the user's own address)
    mapping(address => address) public rewardRecipients;

    // ── Error definitions ──

    error NotRootNet();
    error AlreadyRegistered();
    /// @dev Target address is already registered as an Agent and cannot also be a Principal
    error AddressIsAgent();
    /// @dev Target address is already registered as a Principal and cannot also be an Agent
    error AddressIsPrincipal();
    /// @dev An address may not register itself as its own Agent
    error AgentIsSelf();
    error NotAgentOwner();
    /// @dev The address is not bound as an Agent
    error NotBound();
    /// @dev An Agent may not remove itself
    error CannotRemoveSelf();
    /// @dev An Agent may not revoke its own Manager privileges
    error CannotRevokeSelf();
    error InvalidRecipient();

    /// @dev Only the RootNet contract may call
    modifier onlyRootNet() {
        if (msg.sender != rootNet) revert NotRootNet();
        _;
    }

    /// @notice Constructor — sets the RootNet contract address
    /// @param rootNet_ RootNet contract address
    constructor(address rootNet_) {
        rootNet = rootNet_;
    }

    // ═══════════════════════════════════════════════
    // ── User registration ──
    // ═══════════════════════════════════════════════

    /// @notice Register a new user
    /// @dev Mutual exclusion check: if the address is already an Agent it cannot be registered as a user
    /// @param user Address to register
    function register(address user) external onlyRootNet {
        if (isRegistered[user]) revert AlreadyRegistered();
        // Mutual exclusion: an Agent address cannot simultaneously be a user
        if (agentOwner[user] != address(0)) revert AddressIsAgent();
        isRegistered[user] = true;
        registeredAt[user] = uint64(block.timestamp);
        totalUsers++;
    }

    // ═══════════════════════════════════════════════
    // ── Agent binding ──
    // ═══════════════════════════════════════════════

    /// @notice Bind an Agent to a Principal; supports rebind (Agent moves to a new Principal)
    /// @dev Auto-registers Principal if not yet registered.
    ///      Returns the old Principal address (address(0) if this is a first-time bind).
    ///      Returns the current Principal unchanged (no state change) if already bound to the same Principal.
    ///      On rebind, caller (RootNet) is responsible for freezing the Agent's old allocations.
    /// @param agent Agent address
    /// @param principal Principal address that will own this Agent
    /// @return oldPrincipal Previous Principal address, or address(0) if not previously bound
    function bind(address agent, address principal) external onlyRootNet returns (address oldPrincipal) {
        // Mutual exclusion: a Principal address cannot simultaneously be an Agent
        if (isRegistered[agent]) revert AddressIsPrincipal();
        // An address may not register itself as its own Agent
        if (agent == principal) revert AgentIsSelf();
        // Mutual exclusion: an Agent address cannot be a Principal
        if (agentOwner[principal] != address(0)) revert AddressIsAgent();
        // Auto-register Principal if not yet registered
        if (!isRegistered[principal]) {
            isRegistered[principal] = true;
            registeredAt[principal] = uint64(block.timestamp);
            totalUsers++;
        }
        oldPrincipal = agentOwner[agent];
        // Self-rebind: already bound to this Principal — no state change needed
        if (oldPrincipal == principal) return oldPrincipal;
        if (oldPrincipal != address(0)) {
            // Rebind: remove from old Principal's set
            userAgents[oldPrincipal].remove(agent);
        }
        agentOwner[agent] = principal;
        userAgents[principal].add(agent);
    }

    // ═══════════════════════════════════════════════
    // ── Agent unbinding ──
    // ═══════════════════════════════════════════════

    /// @notice Unbind an Agent from its Principal, returning it to unregistered status
    /// @dev Clears agentOwner, isManager, and removes from the Principal's set.
    ///      After unbinding the Agent address is no longer known to the system and can
    ///      be re-bound (to the same or a different Principal) in the future.
    /// @param agent Agent address to unbind
    /// @return oldPrincipal The Principal this Agent was bound to
    function unbind(address agent) external onlyRootNet returns (address oldPrincipal) {
        oldPrincipal = agentOwner[agent];
        if (oldPrincipal == address(0)) revert NotBound();
        delete agentOwner[agent];
        delete isManager[agent];
        userAgents[oldPrincipal].remove(agent);
    }

    // ═══════════════════════════════════════════════
    // ── Agent management ──
    // ═══════════════════════════════════════════════

    /// @notice Remove a user's Agent (also clears Manager privileges)
    /// @dev An Agent may not remove itself (prevents an Agent from detaching without user authorization)
    /// @param user User address that owns the Agent
    /// @param agent Agent address to remove
    /// @param operator Address initiating the operation (used to block Agent self-removal)
    function removeAgent(address user, address agent, address operator) external onlyRootNet {
        if (agentOwner[agent] != user) revert NotAgentOwner();
        // An Agent may not remove itself
        if (agent == operator) revert CannotRemoveSelf();
        delete agentOwner[agent];
        delete isManager[agent];
        userAgents[user].remove(agent);
    }

    /// @notice Grant or revoke an Agent's Manager privileges
    /// @dev An Agent may not revoke its own Manager privileges
    /// @param user User address that owns the Agent
    /// @param agent Target Agent address
    /// @param _isManager Whether to grant Manager privileges
    /// @param operator Address initiating the operation (used to block Agent self-revocation)
    function setManager(address user, address agent, bool _isManager, address operator) external onlyRootNet {
        if (agentOwner[agent] != user) revert NotAgentOwner();
        // An Agent may not revoke its own Manager privileges
        if (!_isManager && agent == operator) revert CannotRevokeSelf();
        isManager[agent] = _isManager;
    }

    /// @notice Set a user's custom reward recipient address
    /// @param user User address
    /// @param recipient Reward recipient address (must not be zero address)
    function setRewardRecipient(address user, address recipient) external onlyRootNet {
        if (recipient == address(0)) revert InvalidRecipient();
        rewardRecipients[user] = recipient;
    }

    // ═══════════════════════════════════════════════
    // ── Query functions ──
    // ═══════════════════════════════════════════════

    /// @notice Get the owner of an address: returns its owner if Agent, returns itself if registered user, otherwise returns address(0)
    /// @param addr Address to query
    /// @return Owner address
    function getOwner(address addr) external view returns (address) {
        address owner = agentOwner[addr];
        // Check if it is an Agent first
        if (owner != address(0)) return owner;
        // Then check if it is a registered user (user itself is the owner)
        if (isRegistered[addr]) return addr;
        return address(0);
    }

    /// @notice Check whether an address is a registered user
    /// @param addr Address to query
    /// @return Whether it is a registered user
    function isRegisteredUser(address addr) external view returns (bool) {
        return isRegistered[addr];
    }

    /// @notice Check whether an address is a registered Agent
    /// @param addr Address to query
    /// @return Whether it is a registered Agent
    function isRegisteredAgent(address addr) external view returns (bool) {
        return agentOwner[addr] != address(0);
    }

    /// @notice Check whether an address is known in the system (user or Agent)
    /// @param addr Address to query
    /// @return Whether it is a known address
    function isKnownAddress(address addr) external view returns (bool) {
        return isRegistered[addr] || agentOwner[addr] != address(0);
    }

    /// @notice Get all Agent addresses belonging to a user
    /// @param user User address
    /// @return Array of Agent addresses
    function getAgents(address user) external view returns (address[] memory) {
        return userAgents[user].values();
    }

    /// @notice Check whether agent belongs to user (a user is also considered their own Agent)
    /// @param user User address
    /// @param agent Agent address
    /// @return Whether agent is user's Agent
    function isAgent(address user, address agent) external view returns (bool) {
        return agentOwner[agent] == user || (agent == user && isRegistered[user]);
    }

    /// @notice Check whether an Agent has Manager privileges
    /// @param agent Agent address
    /// @return Whether it is a Manager
    function isManagerAgent(address agent) external view returns (bool) {
        return isManager[agent];
    }

    /// @notice Get the reward recipient address for a user (returns the user's own address if not set)
    /// @param user User address
    /// @return Reward recipient address
    function getRewardRecipient(address user) external view returns (address) {
        address r = rewardRecipients[user];
        return r != address(0) ? r : user;
    }

    /// @notice Get the total number of registered users
    /// @return Total user count
    function getTotalUsers() external view returns (uint256) {
        return totalUsers;
    }

    /// @notice Resolve the role of an address in a single call (replaces three separate calls to isRegisteredUser + getOwner + isManagerAgent)
    /// @dev Return value semantics:
    ///   - owner: returns addr itself if it is a user, returns its owner if it is an Agent, otherwise address(0)
    ///   - isUser: whether addr is a registered user
    ///   - isManager_: whether addr has Manager privileges
    /// @param addr Address to resolve
    /// @return owner Owner address
    /// @return isUser Whether it is a registered user
    /// @return isManager_ Whether it is a Manager
    function resolveCallerRole(address addr) external view returns (address owner, bool isUser, bool isManager_) {
        isUser = isRegistered[addr];
        if (isUser) {
            // User itself is the owner
            owner = addr;
        } else {
            // Not a user — try to look up owner as an Agent
            owner = agentOwner[addr];
        }
        isManager_ = isManager[addr];
    }
}
