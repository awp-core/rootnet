// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IAWPRegistry — AWPRegistry contract interface
/// @notice Defines worknet status enums, data structures, and events
interface IAWPRegistry {
    /// @notice Worknet lifecycle status (None=0 reserved for deleted/nonexistent)
    enum WorknetStatus {
        None,     // 0 — default after delete, not a valid state
        Pending,  // 1 — Registered, awaiting activation
        Active,   // 2 — Active, participating in emission
        Paused,   // 3 — Paused (Owner-initiated)
        Banned    // 4 — Banned (Guardian)
    }

    /// @notice Worknet lifecycle state (stored in AWPRegistry; identity data in AWPWorkNet)
    struct WorknetInfo {
        bytes32 lpPool;           // PancakeSwap V4 LP pool ID
        WorknetStatus status;      // Current lifecycle status
        uint64 createdAt;         // Creation timestamp (block.timestamp)
        uint64 activatedAt;       // Activation timestamp
    }

    /// @notice Full worknet view (combines AWPRegistry state + AWPWorkNet identity)
    struct WorknetFullInfo {
        address worknetManager;
        address worknetToken;
        bytes32 lpPool;
        WorknetStatus status;
        uint64 createdAt;
        uint64 activatedAt;
        string name;
        string symbol;
        string skillsURI;
        uint128 minStake;
        address owner;
    }

    /// @notice Worknet registration parameters
    struct WorknetParams {
        string name;           // Alpha Token name
        string symbol;         // Alpha Token symbol
        address worknetManager; // Worknet contract address (0 = auto-deploy WorknetManager proxy)
        bytes32 salt;          // CREATE2 salt for vanity Alpha token address (0 = use worknetId)
        uint128 minStake;      // Minimum stake requirement for agents (0 = no minimum)
        string skillsURI;      // Skills file URI (set at registration, updatable later via AWPWorkNet)
    }

    // ── Account V2 events ──
    event UserRegistered(address indexed user);
    event Bound(address indexed addr, address indexed target);
    event Unbound(address indexed addr);
    event RecipientSet(address indexed addr, address recipient);
    event DelegateGranted(address indexed staker, address indexed delegate);
    event DelegateRevoked(address indexed staker, address indexed delegate);

    // ── Worknet events ──
    event WorknetRegistered(
        uint256 indexed worknetId,
        address indexed owner,
        string name,
        string symbol
    );
    event LPCreated(uint256 indexed worknetId, bytes32 poolId, uint256 awpAmount, uint256 worknetTokenAmount);
    event WorknetActivated(uint256 indexed worknetId);
    event WorknetPaused(uint256 indexed worknetId);
    event WorknetResumed(uint256 indexed worknetId);
    event WorknetCancelled(uint256 indexed worknetId);
    event WorknetRejected(uint256 indexed worknetId);
    event WorknetBanned(uint256 indexed worknetId);
    event WorknetUnbanned(uint256 indexed worknetId);

    // ── Governance parameter events ──
    event GuardianUpdated(address indexed newGuardian);
    event InitialAlphaPriceUpdated(uint256 newPrice);
    event InitialAlphaMintUpdated(uint256 amount);
    event DefaultWorknetManagerImplUpdated(address indexed newImpl);
    // DexConfigUpdated, LPManagerUpdated, veAWPUpdated, AWPWorkNetUpdated removed — all immutable proxies

    // ── View functions ──
    function resolveRecipient(address addr) external view returns (address);
    function batchResolveRecipients(address[] calldata addrs) external view returns (address[] memory);
}
