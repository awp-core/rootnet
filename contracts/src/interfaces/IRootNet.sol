// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IRootNet — RootNet contract interface
/// @notice Defines subnet status enums, data structures, and events
interface IRootNet {
    /// @notice Subnet lifecycle status
    enum SubnetStatus {
        Pending,  // Registered, awaiting activation
        Active,   // Active, participating in emission
        Paused,   // Paused (Owner-initiated)
        Banned    // Banned (DAO governance)
    }

    /// @notice Subnet lifecycle state (stored in RootNet; identity data in SubnetNFT)
    struct SubnetInfo {
        bytes32 lpPool;           // PancakeSwap V4 LP pool ID
        SubnetStatus status;      // Current lifecycle status
        uint64 createdAt;         // Creation timestamp (block.timestamp)
        uint64 activatedAt;       // Activation timestamp
    }

    /// @notice Full subnet view (combines RootNet state + SubnetNFT identity)
    struct SubnetFullInfo {
        address subnetManager;
        address alphaToken;
        bytes32 lpPool;
        SubnetStatus status;
        uint64 createdAt;
        uint64 activatedAt;
        string name;
        string skillsURI;
        uint128 minStake;
        address owner;
    }

    /// @notice Subnet registration parameters (contains strings; stored via events only, not on-chain)
    struct SubnetParams {
        string name;           // Alpha Token name
        string symbol;         // Alpha Token symbol
        string metadataURI;    // IPFS metadata URI
        address subnetManager; // Subnet contract address (0 = auto-deploy SubnetManager proxy)
        string coordinatorURL; // Coordinator service URL
        bytes32 salt;          // CREATE2 salt for vanity Alpha token address (0 = use subnetId)
        uint128 minStake;      // Minimum stake requirement for agents (0 = no minimum)
    }

    // ── User events ──
    event UserRegistered(address indexed user);
    event AgentBound(address indexed principal, address indexed agent, address oldPrincipal);
    /// @notice Emitted when an Agent voluntarily unbinds itself from its Principal
    event AgentUnbound(address indexed principal, address indexed agent);
    event AgentRemoved(address indexed user, address indexed agent, address operator);
    event DelegationUpdated(address indexed user, address indexed agent, bool isManager, address operator);
    event RewardRecipientUpdated(address indexed user, address recipient);

    // ── Staking events ──
    event Allocated(
        address indexed user, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator
    );
    event Deallocated(
        address indexed user, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator
    );
    event Reallocated(
        address indexed user,
        address fromAgent,
        uint256 fromSubnet,
        address toAgent,
        uint256 toSubnet,
        uint256 amount,
        address operator
    );

    // ── Subnet events ──
    event SubnetRegistered(
        uint256 indexed subnetId,
        address indexed owner,
        string name,
        string symbol,
        string metadataURI,
        address subnetManager,
        address alphaToken,
        string coordinatorURL
    );
    event LPCreated(uint256 indexed subnetId, bytes32 poolId, uint256 awpAmount, uint256 alphaAmount);
    /// @dev Off-chain notification only — no data stored on-chain. Indexer writes to DB.
    event MetadataUpdated(uint256 indexed subnetId, string metadataURI, string coordinatorURL);
    event SubnetActivated(uint256 indexed subnetId);
    event SubnetPaused(uint256 indexed subnetId);
    event SubnetResumed(uint256 indexed subnetId);
    event SubnetBanned(uint256 indexed subnetId);
    event SubnetUnbanned(uint256 indexed subnetId);
    event SubnetDeregistered(uint256 indexed subnetId);
}
