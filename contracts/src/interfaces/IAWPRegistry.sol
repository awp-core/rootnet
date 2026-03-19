// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IAWPRegistry — AWPRegistry contract interface
/// @notice Defines subnet status enums, data structures, and events
interface IAWPRegistry {
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

    /// @notice Subnet registration parameters
    struct SubnetParams {
        string name;           // Alpha Token name
        string symbol;         // Alpha Token symbol
        address subnetManager; // Subnet contract address (0 = auto-deploy SubnetManager proxy)
        bytes32 salt;          // CREATE2 salt for vanity Alpha token address (0 = use subnetId)
        uint128 minStake;      // Minimum stake requirement for agents (0 = no minimum)
        string skillsURI;      // Skills file URI (set at registration, updatable later via SubnetNFT)
    }

    // ── Account V2 events ──
    event UserRegistered(address indexed user);
    event Bound(address indexed addr, address indexed target);
    event Unbound(address indexed addr);
    event RecipientSet(address indexed addr, address recipient);
    event DelegateGranted(address indexed staker, address indexed delegate);
    event DelegateRevoked(address indexed staker, address indexed delegate);

    // ── Staking events ──
    event Allocated(
        address indexed staker, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator
    );
    event Deallocated(
        address indexed staker, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator
    );
    event Reallocated(
        address indexed staker,
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
        address subnetManager,
        address alphaToken
    );
    event LPCreated(uint256 indexed subnetId, bytes32 poolId, uint256 awpAmount, uint256 alphaAmount);
    event SubnetActivated(uint256 indexed subnetId);
    event SubnetPaused(uint256 indexed subnetId);
    event SubnetResumed(uint256 indexed subnetId);
    event SubnetBanned(uint256 indexed subnetId);
    event SubnetUnbanned(uint256 indexed subnetId);
    event SubnetDeregistered(uint256 indexed subnetId);

    // ── Governance parameter events ──
    event GuardianUpdated(address indexed newGuardian);
    event InitialAlphaPriceUpdated(uint256 newPrice);
    event ImmunityPeriodUpdated(uint256 newPeriod);
    event AlphaTokenFactoryUpdated(address indexed newFactory);
    event DefaultSubnetManagerImplUpdated(address indexed newImpl);
}
