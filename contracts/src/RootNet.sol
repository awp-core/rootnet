// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Pausable} from "@openzeppelin/contracts/utils/Pausable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {EIP712} from "@openzeppelin/contracts/utils/cryptography/EIP712.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

import {IAlphaToken} from "./interfaces/IAlphaToken.sol";
import {IAlphaTokenFactory} from "./interfaces/IAlphaTokenFactory.sol";
import {IAccessManager} from "./interfaces/IAccessManager.sol";
import {IStakingVault} from "./interfaces/IStakingVault.sol";
import {IStakeNFT} from "./interfaces/IStakeNFT.sol";
import {ISubnetNFT} from "./interfaces/ISubnetNFT.sol";
import {ILPManager} from "./interfaces/ILPManager.sol";
import {IRootNet} from "./interfaces/IRootNet.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

/// @title RootNet — Unified entry point for the AWP protocol (subnet management + staking management)
/// @author AWP Team
/// @notice This contract is the core control layer of the entire AWP RootNet Agent Mining protocol.
///         All user interactions (registration, staking, subnet management) go through this contract.
///         Emission logic has been migrated to the AWPEmission contract.
/// @dev Inheritance: IRootNet (interface defining enums/structs/events), Pausable (emergency pause),
///      ReentrancyGuard (reentrancy protection), EIP712 (EIP-712 signing domain, domain name "AWPRootNet" v1).
///      10 external module addresses are injected once via initializeRegistry; the deployer is then zeroed and cannot call it again.
contract RootNet is IRootNet, Pausable, ReentrancyGuard, EIP712 {
    using SafeERC20 for IERC20;
    using EnumerableSet for EnumerableSet.UintSet;

    // ══════════════════════════════════════════════
    //  Address registry — external module addresses injected via initializeRegistry after deployment
    // ══════════════════════════════════════════════

    /// @notice AWP token contract address (ERC20, 10B MAX_SUPPLY, 50% minted in constructor, 50% minted on-demand)
    address public awpToken;
    /// @notice SubnetNFT contract address; each subnet corresponds to one NFT (tokenId = subnetId)
    address public subnetNFT;
    /// @notice AlphaToken factory contract address, used to deploy an independent Alpha token for each subnet
    address public alphaTokenFactory;
    /// @notice AWP emission contract address, the only contract holding AWP minting rights
    address public awpEmission;
    /// @notice LP manager contract address, responsible for creating AWP/Alpha trading pairs and managing liquidity
    address public lpManager;
    /// @notice Access control manager — manages user registration / Agent registration / Manager privileges
    address public accessManager;
    /// @notice Staking vault contract address — manages allocation bookkeeping
    address public stakingVault;
    /// @notice StakeNFT contract address — manages AWP staking positions as NFTs
    address public stakeNFT;
    /// @notice Treasury (Timelock) address — holds governance operation rights (onlyTimelock)
    address public treasury;
    /// @notice Guardian address — holds emergency pause rights (onlyGuardian)
    address public guardian;
    /// @notice Default subnet implementation address (for auto-deploying subnet contracts via proxy)
    address public defaultSubnetManagerImpl;

    /// @dev Deployer address; used only for initializeRegistry, zeroed immediately after the call
    address private _deployer;
    /// @notice Whether the registry has been initialized (can only be initialized once)
    bool public registryInitialized;

    // ══════════════════════════════════════════════
    //  Subnet data
    // ══════════════════════════════════════════════

    /// @notice subnetId => SubnetInfo mapping, stores the on-chain state of each subnet
    /// @dev Strings are not stored on-chain (name/symbol/metadataURI/coordinatorURL); they are recorded via events and written to DB by the Indexer
    mapping(uint256 => SubnetInfo) public subnets;

    /// @dev Next subnet ID to be assigned, auto-increments from 1 (tokenId = subnetId)
    uint256 private _nextSubnetId = 1;

    /// @dev Active subnet ID set (migrated here from AWPEmission V3)
    EnumerableSet.UintSet private activeSubnetIds;

    /// @notice Maximum number of active subnets
    uint128 public constant MAX_ACTIVE_SUBNETS = 10000;

    /// @notice Initial price of the Alpha token when registering a subnet (denominated in AWP), default 0.01 AWP
    uint256 public initialAlphaPrice = 1e16; // 0.01 AWP

    /// @notice Initial Alpha token mint amount per subnet: 100 million (100_000_000 * 1e18)
    uint256 public constant INITIAL_ALPHA_MINT = 100_000_000 * 1e18;

    /// @notice Subnet deregistration immunity period; Timelock cannot deregister the subnet during this window
    uint256 public immunityPeriod = 30 days;

    // ══════════════════════════════════════════════
    //  Gasless registration — EIP-712 signature related
    // ══════════════════════════════════════════════

    /// @notice Per-signer nonce for replay attack prevention (nonces[signer]++ when validated)
    mapping(address => uint256) public nonces;

    /// @dev EIP-712 type hash: Register(address user, uint256 nonce, uint256 deadline)
    bytes32 private constant REGISTER_TYPEHASH =
        keccak256("Register(address user,uint256 nonce,uint256 deadline)");

    /// @dev EIP-712 type hash: Bind(address agent, address principal, uint256 nonce, uint256 deadline)
    bytes32 private constant BIND_TYPEHASH =
        keccak256("Bind(address agent,address principal,uint256 nonce,uint256 deadline)");

    /// @dev EIP-712 type hash: RegisterSubnet(address user, bytes32 paramsHash, uint256 nonce, uint256 deadline)
    bytes32 private constant REGISTER_SUBNET_TYPEHASH =
        keccak256("RegisterSubnet(address user,bytes32 paramsHash,uint256 nonce,uint256 deadline)");

    // ══════════════════════════════════════════════
    //  Permission modifiers
    // ══════════════════════════════════════════════

    /// @dev Only the Treasury (Timelock) may call, used for governance operations
    modifier onlyTimelock() {
        if (msg.sender != treasury) revert NotTimelock();
        _;
    }

    /// @dev Only the Guardian may call, used for emergency pausing
    modifier onlyGuardian() {
        if (msg.sender != guardian) revert NotGuardian();
        _;
    }

    // ══════════════════════════════════════════════
    //  Custom errors
    // ══════════════════════════════════════════════

    /// @dev Non-deployer called initializeRegistry
    error NotDeployer();
    /// @dev Registry is already initialized; cannot call again
    error AlreadyInitialized();
    /// @dev Unknown address
    error UnknownAddress();
    /// @dev Caller is not the Treasury (Timelock)
    error NotTimelock();
    /// @dev Caller is not the Guardian
    error NotGuardian();
    /// @dev Caller is not a Manager
    error NotManager();
    /// @dev Caller is not registered as a user
    error NotRegistered();
    /// @dev Subnet parameters are invalid (name/symbol length out of bounds)
    error InvalidSubnetParams();
    /// @dev Subnet contract address cannot be the zero address
    error SubnetManagerRequired();
    /// @dev Subnet status does not meet the precondition for the operation
    error InvalidSubnetStatus();
    /// @dev Subnet is still within its immunity period and cannot be deregistered
    error ImmunityNotExpired();
    /// @dev Caller is not the subnet NFT owner
    error NotOwner();
    /// @dev The specified Agent is not a valid Agent for this user
    error InvalidAgent();
    /// @dev Initial Alpha price is too low (minimum 1e12)
    error PriceTooLow();
    /// @dev EIP-712 signature has expired (block.timestamp > deadline)
    error ExpiredSignature();
    /// @dev EIP-712 signature verification failed (recovered signer does not match expected)
    error InvalidSignature();
    /// @dev Active subnet count has reached the maximum
    error MaxActiveSubnetsReached();
    /// @dev Initial Alpha price is too high
    error PriceTooHigh();
    /// @dev Allocation does not meet the subnet's minimum stake requirement
    error InsufficientMinStake();

    // ══════════════════════════════════════════════
    //  Constructor
    // ══════════════════════════════════════════════

    /// @notice Deploy the RootNet contract
    /// @dev EIP-712 domain name is "AWPRootNet", version "1".
    ///      deployer_ is used only for the subsequent initializeRegistry call; zeroed immediately after.
    /// @param deployer_ Deployer address (holds initializeRegistry rights, self-destructs after call)
    /// @param treasury_ Treasury (Timelock) address, holds governance rights
    /// @param guardian_ Guardian address, holds emergency pause rights
    constructor(
        address deployer_,
        address treasury_,
        address guardian_
    ) EIP712("AWPRootNet", "1") {
        _deployer = deployer_;
        treasury = treasury_;
        guardian = guardian_;
    }

    // ═══════════════════════════════════════════════
    //  Registry — module address registry
    // ═══════════════════════════════════════════════

    /// @notice One-time initialization of all external module addresses; only callable by deployer, and only once
    /// @dev After successful call _deployer is zeroed, permanently locked; cannot be called again.
    ///      Deployment flow: deploy all modules → deploy RootNet → initializeRegistry → deployer zeroed.
    /// @param awpToken_ AWP token contract address
    /// @param subnetNFT_ SubnetNFT contract address
    /// @param alphaTokenFactory_ AlphaToken factory contract address
    /// @param awpEmission_ AWP emission contract address
    /// @param lpManager_ LP manager contract address
    /// @param accessManager_ Access control manager contract address
    /// @param stakingVault_ Staking vault contract address
    /// @param stakeNFT_ StakeNFT contract address
    function initializeRegistry(
        address awpToken_,
        address subnetNFT_,
        address alphaTokenFactory_,
        address awpEmission_,
        address lpManager_,
        address accessManager_,
        address stakingVault_,
        address stakeNFT_
    ) external {
        // Only the deployer may call
        if (msg.sender != _deployer) revert NotDeployer();
        // Prevent re-initialization
        if (registryInitialized) revert AlreadyInitialized();

        awpToken = awpToken_;
        subnetNFT = subnetNFT_;
        alphaTokenFactory = alphaTokenFactory_;
        awpEmission = awpEmission_;
        lpManager = lpManager_;
        accessManager = accessManager_;
        stakingVault = stakingVault_;
        stakeNFT = stakeNFT_;

        // Link StakingVault → StakeNFT (one-time setter, resolves CREATE2 circular dependency)
        IStakingVault(stakingVault).setStakeNFT(stakeNFT);

        registryInitialized = true;
        // Permanently destroy deployer rights; this function can no longer be called
        _deployer = address(0);
    }

    /// @notice Retrieve all module addresses in a single call
    /// @return In order: awpToken, subnetNFT, alphaTokenFactory, awpEmission,
    ///         lpManager, accessManager, stakingVault, stakeNFT, treasury, guardian
    function getRegistry()
        external
        view
        returns (
            address,
            address,
            address,
            address,
            address,
            address,
            address,
            address,
            address,
            address
        )
    {
        return (
            awpToken,
            subnetNFT,
            alphaTokenFactory,
            awpEmission,
            lpManager,
            accessManager,
            stakingVault,
            stakeNFT,
            treasury,
            guardian
        );
    }

    // ═══════════════════════════════════════════════
    //  Identity Resolution
    // ═══════════════════════════════════════════════

    /// @notice Resolve the caller's identity: returns msg.sender if registered Principal, returns its Principal if a Manager Agent
    /// @dev Uses AccessManager.resolveCallerRole for a single external call resolving the three-part identity (Principal / owner / isManager).
    ///      - msg.sender is a registered Principal → returns msg.sender (Principal operates directly)
    ///      - msg.sender is a Manager Agent → returns its Principal address (Manager acts on behalf)
    ///      - Otherwise → revert
    /// @return Resolved Principal address
    function _resolvePrincipalOrDelegated() internal view returns (address) {
        (address owner, bool isUser, bool isManager_) =
            IAccessManager(accessManager).resolveCallerRole(msg.sender);
        // Caller is itself a registered Principal; return directly
        if (isUser) return msg.sender;
        // owner is zero address means msg.sender is neither a Principal nor an Agent
        if (owner == address(0)) revert UnknownAddress();
        // Is an Agent but not a Manager; no permission to operate on behalf
        if (!isManager_) revert NotManager();
        return owner;
    }

    /// @notice Resolve the caller's identity: only registered users are allowed, Manager proxy is not permitted
    /// @dev Used for sensitive operations such as setRewardRecipient
    /// @return Caller address (must be a registered user)
    function _resolveOwnerOnly() internal view returns (address) {
        (,bool isUser,) = IAccessManager(accessManager).resolveCallerRole(msg.sender);
        if (!isUser) revert NotRegistered();
        return msg.sender;
    }

    /// @notice Validate that the specified agent is a valid Agent belonging to the specified user
    /// @param user User address
    /// @param agent Agent address
    function _validateAgent(address user, address agent) internal view {
        if (!IAccessManager(accessManager).isAgent(user, agent)) revert InvalidAgent();
    }

    // ═══════════════════════════════════════════════
    //  User Registration
    // ═══════════════════════════════════════════════

    /// @notice Self-register as a Principal (msg.sender registers directly)
    /// @dev Calls AccessManager.register; enforces mutual exclusion (cannot be both Principal and Agent)
    function register() external nonReentrant whenNotPaused {
        IAccessManager(accessManager).register(msg.sender);
        emit UserRegistered(msg.sender);
    }

    /// @notice Register as a Principal with optional reward recipient and/or initial stake
    /// @dev Combines register + setRewardRecipient + StakeNFT.depositFor in a single transaction.
    ///      Any parameter can be omitted (zero value = skip that step).
    ///      User must pre-approve AWP to the StakeNFT contract if depositAmount > 0.
    /// @param recipient Custom reward recipient address (address(0) = use own address / default)
    /// @param depositAmount AWP amount to stake via StakeNFT (0 = skip staking)
    /// @param lockDuration Lock duration in seconds for the StakeNFT position (0 = skip staking)
    function register(address recipient, uint256 depositAmount, uint64 lockDuration)
        external
        nonReentrant
        whenNotPaused
    {
        // Register if not already registered
        if (!IAccessManager(accessManager).isRegisteredUser(msg.sender)) {
            IAccessManager(accessManager).register(msg.sender);
            emit UserRegistered(msg.sender);
        }
        // Set custom reward recipient if provided
        if (recipient != address(0)) {
            IAccessManager(accessManager).setRewardRecipient(msg.sender, recipient);
            emit RewardRecipientUpdated(msg.sender, recipient);
        }
        // Stake if both amount and duration are non-zero
        if (depositAmount > 0 && lockDuration > 0) {
            IStakeNFT(stakeNFT).depositFor(msg.sender, depositAmount, lockDuration);
        }
    }

    /// @notice Gasless user registration: anyone can pay gas to submit a user-signed registration request
    /// @dev EIP-712 signature verification flow:
    ///      1. Check that the signature has not expired (block.timestamp <= deadline)
    ///      2. Construct structHash = keccak256(abi.encode(REGISTER_TYPEHASH, user, nonces[user]++, deadline))
    ///         nonces[user]++ reads the current value first then increments, preventing replay attacks
    ///      3. Generate EIP-712 digest via _hashTypedDataV4
    ///      4. ECDSA.recover recovers the signer; verify signer == user
    /// @param user User address to register (must be the signer)
    /// @param deadline Signature expiry time (unix timestamp)
    /// @param v EIP-712 signature v value
    /// @param r EIP-712 signature r value
    /// @param s EIP-712 signature s value
    function registerFor(address user, uint256 deadline, uint8 v, bytes32 r, bytes32 s)
        external
        nonReentrant
        whenNotPaused
    {
        // Check whether the signature has expired
        if (block.timestamp > deadline) revert ExpiredSignature();
        // Construct EIP-712 struct hash; nonce is read then incremented to prevent replay
        bytes32 structHash = keccak256(abi.encode(REGISTER_TYPEHASH, user, nonces[user]++, deadline));
        // Generate the full EIP-712 digest (including domain separator)
        bytes32 digest = _hashTypedDataV4(structHash);
        // Recover the signer from the signature
        address signer = ECDSA.recover(digest, v, r, s);
        // Verify the recovered signer matches the supplied user address
        if (signer != user) revert InvalidSignature();

        IAccessManager(accessManager).register(user);
        emit UserRegistered(user);
    }

    /// @notice One-stop register + stake (via StakeNFT) + allocate: reduces the number of user interactions
    /// @dev Executes in order: register (if not yet registered) → depositFor via StakeNFT (if depositAmount > 0)
    ///      → allocate (if parameters are complete). Each step is optional.
    ///      User must have pre-approved AWP to StakeNFT contract.
    /// @param depositAmount AWP amount to deposit into StakeNFT (0 to skip deposit)
    /// @param lockDuration Lock duration in seconds for the StakeNFT position
    /// @param agent Agent address to allocate to (address(0) to skip allocation)
    /// @param subnetId Subnet ID to allocate to (0 to skip allocation)
    /// @param allocateAmount Staking amount to allocate (0 to skip allocation)
    function registerAndStake(
        uint256 depositAmount,
        uint64 lockDuration,
        address agent,
        uint256 subnetId,
        uint256 allocateAmount
    ) external nonReentrant whenNotPaused {
        // Auto-register the user if not yet registered
        if (!IAccessManager(accessManager).isRegisteredUser(msg.sender)) {
            IAccessManager(accessManager).register(msg.sender);
            emit UserRegistered(msg.sender);
        }
        // Deposit via StakeNFT: create a new position for the user
        if (depositAmount > 0 && lockDuration > 0) {
            IStakeNFT(stakeNFT).depositFor(msg.sender, depositAmount, lockDuration);
        }
        // Allocate: assign staked AWP to the (agent, subnetId) triple
        if (allocateAmount > 0 && agent != address(0) && subnetId > 0) {
            if (subnets[subnetId].status != SubnetStatus.Active) revert InvalidSubnetStatus();
            _validateAgent(msg.sender, agent);
            IStakingVault(stakingVault).allocate(msg.sender, agent, subnetId, allocateAmount);
            emit Allocated(msg.sender, agent, subnetId, allocateAmount, msg.sender);
        }
    }

    // ═══════════════════════════════════════════════
    //  Agent Binding
    // ═══════════════════════════════════════════════

    /// @notice Bind msg.sender as an Agent to the specified Principal
    /// @dev Calls AccessManager.bind; auto-registers Principal if not yet registered.
    ///      Supports rebind: if Agent was previously bound to another Principal,
    ///      its allocations are frozen automatically (StakingVault enumerates subnets).
    ///      No-op (silent return) if Agent is already bound to the same Principal.
    /// @param principal Principal address this Agent belongs to
    function bind(address principal) external nonReentrant whenNotPaused {
        if (principal == address(0)) revert UnknownAddress();
        address oldPrincipal = IAccessManager(accessManager).bind(msg.sender, principal);
        // Same-principal rebind: AccessManager made no changes, nothing to do
        if (oldPrincipal == principal) return;
        if (oldPrincipal != address(0)) {
            // True rebind: freeze all allocations from the Agent under its old Principal
            IStakingVault(stakingVault).freezeAgentAllocations(oldPrincipal, msg.sender);
        }
        emit AgentBound(principal, msg.sender, oldPrincipal);
    }

    /// @notice Gasless Agent bind: anyone can pay gas to submit an agent-signed bind request
    /// @dev The agent signature proves possession of the private key for that address (mining devices do not need to hold BNB to bind).
    ///      EIP-712 signature verification flow:
    ///      1. Check signature expiry
    ///      2. structHash = keccak256(abi.encode(BIND_TYPEHASH, agent, principal, nonces[agent]++, deadline))
    ///         Note: the nonce uses the agent address's nonce (the signer is the agent)
    ///      3. ECDSA.recover recovers the signer; verify signer == agent
    /// @param agent Agent address (signer, mining device address)
    /// @param principal Principal address this Agent belongs to
    /// @param deadline Signature expiry time
    /// @param v EIP-712 signature v value
    /// @param r EIP-712 signature r value
    /// @param s EIP-712 signature s value
    function bindFor(address agent, address principal, uint256 deadline, uint8 v, bytes32 r, bytes32 s)
        external
        nonReentrant
        whenNotPaused
    {
        if (principal == address(0)) revert UnknownAddress();
        // Check whether the signature has expired
        if (block.timestamp > deadline) revert ExpiredSignature();
        // Construct EIP-712 struct hash using the agent's nonce (the signer is the agent)
        bytes32 structHash = keccak256(abi.encode(BIND_TYPEHASH, agent, principal, nonces[agent]++, deadline));
        bytes32 digest = _hashTypedDataV4(structHash);
        address signer = ECDSA.recover(digest, v, r, s);
        // Verify the signer is the agent itself
        if (signer != agent) revert InvalidSignature();

        address oldPrincipal = IAccessManager(accessManager).bind(agent, principal);
        if (oldPrincipal == principal) return;
        if (oldPrincipal != address(0)) {
            IStakingVault(stakingVault).freezeAgentAllocations(oldPrincipal, agent);
        }
        emit AgentBound(principal, agent, oldPrincipal);
    }

    /// @notice Agent voluntarily unbinds itself from its Principal
    /// @dev The Agent returns to unregistered status; all its allocations are frozen automatically
    ///      (StakingVault enumerates subnets — no caller-supplied list required).
    ///      After unbinding the address may bind to any Principal in the future.
    ///      Intentionally NOT gated by whenNotPaused — Agents must always be able to voluntarily detach.
    function unbind() external nonReentrant {
        address oldPrincipal = IAccessManager(accessManager).unbind(msg.sender);
        IStakingVault(stakingVault).freezeAgentAllocations(oldPrincipal, msg.sender);
        emit AgentUnbound(oldPrincipal, msg.sender);
    }

    // ═══════════════════════════════════════════════
    //  Agent Management
    // ═══════════════════════════════════════════════

    /// @notice Remove an Agent: validate ownership, freeze all its allocations, then remove from AccessManager
    /// @dev Validates agent ownership first, then freezes via StakingVault (auto-enumerates subnets).
    ///      Supports Principal or delegated Agent calling (identity resolved via _resolvePrincipalOrDelegated).
    /// @param agent Agent address to remove
    function removeAgent(address agent) external nonReentrant {
        address user = _resolvePrincipalOrDelegated();
        // Validate agent belongs to user before mutating StakingVault state
        if (!IAccessManager(accessManager).isAgent(user, agent) || agent == user) revert UnknownAddress();
        // Freeze all allocations this Agent has (StakingVault enumerates subnets automatically)
        IStakingVault(stakingVault).freezeAgentAllocations(user, agent);
        // Remove the Agent binding from AccessManager
        IAccessManager(accessManager).removeAgent(user, agent, msg.sender);
        emit AgentRemoved(user, agent, msg.sender);
    }

    /// @notice Set delegation privileges for an Agent (delegated Agents may perform allocate/deallocate/reallocate on behalf of the Principal)
    /// @dev Supports Principal or existing delegated Agent calling
    /// @param agent Target Agent address
    /// @param _isManager true = grant delegation, false = revoke
    function setDelegation(address agent, bool _isManager) external {
        address user = _resolvePrincipalOrDelegated();
        IAccessManager(accessManager).setManager(user, agent, _isManager, msg.sender);
        emit DelegationUpdated(user, agent, _isManager, msg.sender);
    }

    /// @notice Set the reward recipient address (only the Owner may call; Manager proxy is not allowed)
    /// @param recipient New reward recipient address
    function setRewardRecipient(address recipient) external {
        address user = _resolveOwnerOnly();
        IAccessManager(accessManager).setRewardRecipient(user, recipient);
        emit RewardRecipientUpdated(user, recipient);
    }

    // ═══════════════════════════════════════════════
    //  Staking: Allocation
    // ═══════════════════════════════════════════════

    /// @notice Allocate deposited AWP to a (user, agent, subnetId) triple
    /// @dev Supports Owner or Manager calling. Allocation takes effect immediately.
    /// @param agent Target Agent address (must be a valid Agent of the user)
    /// @param subnetId Target subnet ID
    /// @param amount AWP amount to allocate
    function allocate(address agent, uint256 subnetId, uint256 amount) external nonReentrant whenNotPaused {
        address user = _resolvePrincipalOrDelegated();
        if (subnets[subnetId].status != SubnetStatus.Active) revert InvalidSubnetStatus();
        _validateAgent(user, agent);
        IStakingVault(stakingVault).allocate(user, agent, subnetId, amount);
        // Enforce subnet minStake requirement after allocation
        uint128 minStake = ISubnetNFT(subnetNFT).getMinStake(subnetId);
        if (minStake > 0) {
            uint256 totalStake = IStakingVault(stakingVault).getAgentStake(user, agent, subnetId);
            if (totalStake < minStake) revert InsufficientMinStake();
        }
        emit Allocated(user, agent, subnetId, amount, msg.sender);
    }

    /// @notice Deallocate: release staking from a (user, agent, subnetId) triple
    /// @dev Supports Owner or Manager calling. Released AWP returns to the user's unallocated balance.
    ///      StakingVault internally validates that the fromAgent allocation is sufficient.
    /// @param agent Source Agent address (must be a valid Agent of the user)
    /// @param subnetId Source subnet ID
    /// @param amount AWP amount to deallocate
    function deallocate(address agent, uint256 subnetId, uint256 amount) external nonReentrant whenNotPaused {
        address user = _resolvePrincipalOrDelegated();
        _validateAgent(user, agent);
        IStakingVault(stakingVault).deallocate(user, agent, subnetId, amount);
        emit Deallocated(user, agent, subnetId, amount, msg.sender);
    }

    /// @notice Reallocate: move staking from one (agent, subnet) triple to another (immediate)
    /// @dev Supports Owner or Manager calling. Reallocation takes effect immediately (atomic).
    ///      Both agents must be validated as valid Agents of the user.
    /// @param fromAgent Source Agent address
    /// @param fromSubnetId Source subnet ID
    /// @param toAgent Target Agent address
    /// @param toSubnetId Target subnet ID
    /// @param amount AWP amount to reallocate
    function reallocate(
        address fromAgent,
        uint256 fromSubnetId,
        address toAgent,
        uint256 toSubnetId,
        uint256 amount
    ) external nonReentrant whenNotPaused {
        address user = _resolvePrincipalOrDelegated();
        if (subnets[toSubnetId].status != SubnetStatus.Active) revert InvalidSubnetStatus();
        _validateAgent(user, fromAgent);
        _validateAgent(user, toAgent);
        IStakingVault(stakingVault).reallocate(
            user, fromAgent, fromSubnetId, toAgent, toSubnetId, amount
        );
        emit Reallocated(user, fromAgent, fromSubnetId, toAgent, toSubnetId, amount, msg.sender);
    }

    // ═══════════════════════════════════════════════
    //  Subnet Registration
    // ═══════════════════════════════════════════════

    /// @notice Register a new subnet: deploy Alpha token, create LP, mint NFT
    /// @dev Full registration flow (8 steps):
    ///      1. Calculate AWP amount required for LP creation: INITIAL_ALPHA_MINT * initialAlphaPrice / 1e18
    ///      2. Transfer AWP directly from the user to LPManager (RootNet does not intermediate)
    ///      3. Mint SubnetNFT (tokenId = _nextSubnetId++) to msg.sender
    ///      4. Deploy Alpha token via factory (admin = address(this), i.e. RootNet)
    ///      5. RootNet mints INITIAL_ALPHA_MINT Alpha tokens to LPManager
    ///      6. LPManager creates the AWP/Alpha trading pair and injects initial liquidity
    ///      7. Set subnetManager as the sole minter of the Alpha token (setSubnetMinter permanently locked)
    ///      8. Store SubnetInfo (strings not stored on-chain; name/symbol/metadataURI/coordinatorURL recorded via events)
    /// @param params Subnet parameters (name, symbol, metadataURI, subnetManager, coordinatorURL)
    /// @return subnetId Newly created subnet ID
    function registerSubnet(SubnetParams calldata params) external nonReentrant whenNotPaused returns (uint256) {
        return _registerSubnet(msg.sender, params);
    }

    /// @notice Gasless subnet registration: relayer pays gas, user signs EIP-712 and pays AWP
    /// @dev User must have pre-approved AWP to RootNet (or use registerSubnetForWithPermit for fully gasless).
    /// @param user User address (signer, pays AWP, receives NFT + admin)
    /// @param params Subnet parameters
    /// @param deadline Signature expiry time
    function registerSubnetFor(
        address user,
        SubnetParams calldata params,
        uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant whenNotPaused returns (uint256) {
        if (block.timestamp > deadline) revert ExpiredSignature();
        _verifyRegisterSubnetSignature(user, params, deadline, v, r, s);
        return _registerSubnet(user, params);
    }

    /// @notice Fully gasless subnet registration with EIP-2612 permit (no prior approve tx needed)
    /// @dev User signs two off-chain messages: (1) ERC20 permit for AWP, (2) EIP-712 registerSubnet.
    ///      Relayer submits both in one tx. User pays zero gas.
    /// @param user User address (signer, pays AWP, receives NFT + admin)
    /// @param params Subnet parameters
    /// @param deadline Shared deadline for both permit and registerSubnet signatures
    /// @param permitV ERC-2612 permit signature v
    /// @param permitR ERC-2612 permit signature r
    /// @param permitS ERC-2612 permit signature s
    /// @param registerV EIP-712 registerSubnet signature v
    /// @param registerR EIP-712 registerSubnet signature r
    /// @param registerS EIP-712 registerSubnet signature s
    function registerSubnetForWithPermit(
        address user,
        SubnetParams calldata params,
        uint256 deadline,
        uint8 permitV, bytes32 permitR, bytes32 permitS,
        uint8 registerV, bytes32 registerR, bytes32 registerS
    ) external nonReentrant whenNotPaused returns (uint256) {
        if (block.timestamp > deadline) revert ExpiredSignature();
        // Step 1: Execute ERC-2612 permit (user approves AWP to RootNet without a prior tx)
        uint256 lpAWPAmount = INITIAL_ALPHA_MINT * initialAlphaPrice / 1e18;
        IERC20Permit(awpToken).permit(user, address(this), lpAWPAmount, deadline, permitV, permitR, permitS);
        // Step 2: Verify registerSubnet EIP-712 signature
        _verifyRegisterSubnetSignature(user, params, deadline, registerV, registerR, registerS);
        // Step 3: Execute registration (safeTransferFrom uses the just-permitted allowance)
        return _registerSubnet(user, params);
    }

    /// @dev Verify EIP-712 signature for registerSubnetFor
    function _verifyRegisterSubnetSignature(
        address user, SubnetParams calldata params, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) internal {
        bytes32 paramsHash = keccak256(abi.encode(
            params.name, params.symbol, params.metadataURI,
            params.subnetManager, params.coordinatorURL, params.salt, params.minStake
        ));
        bytes32 structHash = keccak256(abi.encode(
            REGISTER_SUBNET_TYPEHASH, user, paramsHash, nonces[user]++, deadline
        ));
        bytes32 digest = _hashTypedDataV4(structHash);
        address signer = ECDSA.recover(digest, v, r, s);
        if (signer != user) revert InvalidSignature();
    }

    /// @dev Internal: shared logic for registerSubnet and registerSubnetFor
    function _registerSubnet(address user, SubnetParams calldata params) internal returns (uint256) {
        if (bytes(params.name).length == 0 || bytes(params.name).length > 64) revert InvalidSubnetParams();
        if (bytes(params.symbol).length == 0 || bytes(params.symbol).length > 16) revert InvalidSubnetParams();
        bool autoDeploySubnet = params.subnetManager == address(0);
        if (autoDeploySubnet && defaultSubnetManagerImpl == address(0)) revert SubnetManagerRequired();

        uint256 lpAWPAmount = INITIAL_ALPHA_MINT * initialAlphaPrice / 1e18;
        IERC20(awpToken).safeTransferFrom(user, lpManager, lpAWPAmount);

        uint256 subnetId = _nextSubnetId++;

        (address alphaToken, bytes32 poolId) = _deployAlphaAndLP(
            subnetId, params.name, params.symbol, lpAWPAmount, params.salt
        );

        address sc;
        if (autoDeploySubnet) {
            bytes memory initData = abi.encodeWithSignature(
                "initialize(address,address,bytes32,address)",
                alphaToken, awpToken, poolId, user
            );
            sc = address(new ERC1967Proxy(defaultSubnetManagerImpl, initData));
        } else {
            sc = params.subnetManager;
        }

        IAlphaToken(alphaToken).setSubnetMinter(sc);
        ISubnetNFT(subnetNFT).mint(user, subnetId, params.name, sc, alphaToken, params.minStake);
        subnets[subnetId] = SubnetInfo({
            lpPool: poolId,
            status: SubnetStatus.Pending,
            createdAt: uint64(block.timestamp),
            activatedAt: 0
        });

        _emitSubnetRegistered(subnetId, user, sc, alphaToken, poolId, lpAWPAmount, params);
        return subnetId;
    }

    /// @dev Deploy Alpha token + create LP (does NOT set subnet minter — caller must do that)
    function _deployAlphaAndLP(
        uint256 subnetId, string calldata name, string calldata symbol,
        uint256 lpAWPAmount, bytes32 salt
    ) internal returns (address alphaToken, bytes32 poolId) {
        alphaToken = IAlphaTokenFactory(alphaTokenFactory).deploy(subnetId, name, symbol, address(this), salt);
        IAlphaToken(alphaToken).mint(lpManager, INITIAL_ALPHA_MINT);
        (poolId,) = ILPManager(lpManager).createPoolAndAddLiquidity(alphaToken, lpAWPAmount, INITIAL_ALPHA_MINT);
    }

    /// @dev Emit subnet registration events
    function _emitSubnetRegistered(
        uint256 subnetId, address user, address sc, address alphaToken, bytes32 poolId, uint256 lpAWPAmount, SubnetParams calldata params
    ) internal {
        emit SubnetRegistered(
            subnetId, user,
            params.name, params.symbol, params.metadataURI,
            sc, alphaToken,
            params.coordinatorURL
        );
        emit LPCreated(subnetId, poolId, lpAWPAmount, INITIAL_ALPHA_MINT);
    }

    // ═══════════════════════════════════════════════
    //  Subnet Lifecycle Management
    //  State transitions: Pending → Active ⇌ Paused → Active, Active/Paused → Banned → Active
    //  Deregistration: any state → deleted (requires immunity period to have elapsed)
    // ═══════════════════════════════════════════════

    /// @notice Activate a subnet: Pending → Active (only the NFT Owner may call)
    /// @param subnetId Subnet ID to activate
    function activateSubnet(uint256 subnetId) external whenNotPaused {
        if (ISubnetNFT(subnetNFT).ownerOf(subnetId) != msg.sender) revert NotOwner();
        SubnetInfo storage info = subnets[subnetId];
        if (info.status != SubnetStatus.Pending) revert InvalidSubnetStatus();

        if (activeSubnetIds.length() >= MAX_ACTIVE_SUBNETS) revert MaxActiveSubnetsReached();
        info.status = SubnetStatus.Active;
        info.activatedAt = uint64(block.timestamp);
        activeSubnetIds.add(subnetId);

        emit SubnetActivated(subnetId);
    }

    /// @notice Pause a subnet: Active → Paused (only the NFT Owner may call)
    /// @param subnetId Subnet ID to pause
    function pauseSubnet(uint256 subnetId) external {
        if (ISubnetNFT(subnetNFT).ownerOf(subnetId) != msg.sender) revert NotOwner();
        SubnetInfo storage info = subnets[subnetId];
        if (info.status != SubnetStatus.Active) revert InvalidSubnetStatus();

        info.status = SubnetStatus.Paused;
        activeSubnetIds.remove(subnetId);

        emit SubnetPaused(subnetId);
    }

    /// @notice Resume a subnet: Paused → Active (only the NFT Owner may call)
    /// @param subnetId Subnet ID to resume
    function resumeSubnet(uint256 subnetId) external whenNotPaused {
        if (ISubnetNFT(subnetNFT).ownerOf(subnetId) != msg.sender) revert NotOwner();
        SubnetInfo storage info = subnets[subnetId];
        if (info.status != SubnetStatus.Paused) revert InvalidSubnetStatus();

        if (activeSubnetIds.length() >= MAX_ACTIVE_SUBNETS) revert MaxActiveSubnetsReached();
        info.status = SubnetStatus.Active;
        activeSubnetIds.add(subnetId);

        emit SubnetResumed(subnetId);
    }

    /// @notice Ban a subnet: Active/Paused → Banned (only Timelock may call)
    /// @param subnetId Subnet ID to ban
    function banSubnet(uint256 subnetId) external onlyTimelock {
        SubnetInfo storage info = subnets[subnetId];
        SubnetStatus status = info.status;
        if (status != SubnetStatus.Active && status != SubnetStatus.Paused) revert InvalidSubnetStatus();

        address sc = ISubnetNFT(subnetNFT).getSubnetManager(subnetId);
        if (sc != address(0)) {
            IAlphaToken(ISubnetNFT(subnetNFT).getAlphaToken(subnetId)).setMinterPaused(sc, true);
        }
        if (status == SubnetStatus.Active) {
            activeSubnetIds.remove(subnetId);
        }
        info.status = SubnetStatus.Banned;

        emit SubnetBanned(subnetId);
    }

    /// @notice Unban a subnet: Banned → Active (only Timelock may call)
    /// @param subnetId Subnet ID to unban
    function unbanSubnet(uint256 subnetId) external onlyTimelock {
        SubnetInfo storage info = subnets[subnetId];
        if (info.status != SubnetStatus.Banned) revert InvalidSubnetStatus();

        address sc = ISubnetNFT(subnetNFT).getSubnetManager(subnetId);
        if (sc != address(0)) {
            IAlphaToken(ISubnetNFT(subnetNFT).getAlphaToken(subnetId)).setMinterPaused(sc, false);
        }
        if (activeSubnetIds.length() >= MAX_ACTIVE_SUBNETS) revert MaxActiveSubnetsReached();
        info.status = SubnetStatus.Active;
        activeSubnetIds.add(subnetId);

        emit SubnetUnbanned(subnetId);
    }

    /// @notice Deregister a subnet: permanently delete subnet data (only Timelock may call; immunity period must have elapsed)
    /// @param subnetId Subnet ID to deregister
    function deregisterSubnet(uint256 subnetId) external onlyTimelock {
        SubnetInfo storage info = subnets[subnetId];
        if (info.createdAt == 0) revert InvalidSubnetStatus();
        uint256 immunityStart = info.activatedAt > 0 ? uint256(info.activatedAt) : uint256(info.createdAt);
        if (block.timestamp <= immunityStart + immunityPeriod) revert ImmunityNotExpired();

        address sc = ISubnetNFT(subnetNFT).getSubnetManager(subnetId);
        if (sc != address(0)) {
            IAlphaToken(ISubnetNFT(subnetNFT).getAlphaToken(subnetId)).setMinterPaused(sc, true);
        }
        activeSubnetIds.remove(subnetId);
        delete subnets[subnetId];
        ISubnetNFT(subnetNFT).burn(subnetId);

        emit SubnetDeregistered(subnetId);
    }

    // ═══════════════════════════════════════════════
    //  Subnet Parameters
    // ═══════════════════════════════════════════════

    /// @notice Notify off-chain services that subnet metadata has changed (only NFT Owner may call)
    /// @dev This function does NOT store any data on-chain. It only emits a MetadataUpdated event
    ///      which the Indexer listens to and writes to the database. On-chain metadata is served
    ///      via SubnetNFT.tokenURI (baseURI + tokenId). This function exists solely as an on-chain
    ///      notification mechanism for off-chain consumers (Indexer, API, frontends).
    /// @param subnetId Subnet ID
    /// @param metadataURI New metadata URI (e.g. IPFS hash for subnet description/avatar)
    /// @param coordinatorURL New Coordinator service URL
    function updateMetadata(uint256 subnetId, string calldata metadataURI, string calldata coordinatorURL) external {
        if (ISubnetNFT(subnetNFT).ownerOf(subnetId) != msg.sender) revert NotOwner();
        emit MetadataUpdated(subnetId, metadataURI, coordinatorURL);
    }

    /// @notice Set the initial Alpha price when registering a subnet (only Timelock may call)
    /// @dev Minimum price limit 1e12 (prevents precision loss that would cause lpAWPAmount to be 0)
    /// @param price New initial price (AWP wei / Alpha)
    function setInitialAlphaPrice(uint256 price) external onlyTimelock {
        if (price < 1e12) revert PriceTooLow();
        if (price > 1e30) revert PriceTooHigh();
        initialAlphaPrice = price;
    }

    /// @notice Update the guardian address (only Timelock may call)
    /// @param g New guardian address
    function setGuardian(address g) external onlyTimelock {
        if (g == address(0)) revert UnknownAddress();
        guardian = g;
    }

    /// @notice Set the subnet deregistration immunity period (only Timelock may call)
    /// @param p New immunity period duration (seconds, minimum 7 days)
    function setImmunityPeriod(uint256 p) external onlyTimelock {
        if (p < 7 days) revert InvalidSubnetParams();
        immunityPeriod = p;
    }

    /// @notice Replace the AlphaToken factory (only Timelock may call)
    /// @dev New subnets will use the new factory; existing subnets are unaffected.
    /// @param factory New AlphaTokenFactory address
    function setAlphaTokenFactory(address factory) external onlyTimelock {
        if (factory == address(0)) revert UnknownAddress();
        alphaTokenFactory = factory;
    }

    /// @notice Set the default subnet implementation (only Timelock may call)
    /// @dev When subnetManager is address(0) in registerSubnet, an ERC1967Proxy
    ///      pointing to this implementation is deployed automatically.
    ///      Setting address(0) disables auto-deployment (subnetManager becomes required).
    /// @param impl SubnetManager implementation address
    function setSubnetManagerImpl(address impl) external onlyTimelock {
        defaultSubnetManagerImpl = impl;
    }

    // ═══════════════════════════════════════════════
    //  Subnet Queries (view functions)
    // ═══════════════════════════════════════════════

    /// @dev Agent information aggregate struct, used for off-chain queries
    struct AgentInfo {
        /// @dev Agent's owner address
        address owner;
        /// @dev Whether the Agent is valid (has not been removed)
        bool isValid;
        /// @dev Agent's stake amount on the specified subnet
        uint256 stake;
        /// @dev Reward recipient address set by the Owner
        address rewardRecipient;
    }

    /// @notice Query complete information for a single Agent on a specified subnet
    /// @param agent Agent address
    /// @param subnetId Subnet ID
    /// @return AgentInfo containing owner, isValid, stake, rewardRecipient
    function getAgentInfo(address agent, uint256 subnetId) external view returns (AgentInfo memory) {
        address owner = IAccessManager(accessManager).getOwner(agent);
        bool isValid = IAccessManager(accessManager).isAgent(owner, agent);
        uint256 stake = IStakingVault(stakingVault).getAgentStake(owner, agent, subnetId);
        address rewardRecipient = IAccessManager(accessManager).getRewardRecipient(owner);
        return AgentInfo(owner, isValid, stake, rewardRecipient);
    }

    /// @notice Batch query information for multiple Agents on a specified subnet
    /// @param agents Array of Agent addresses
    /// @param subnetId Subnet ID
    /// @return infos Array of AgentInfo
    function getAgentsInfo(address[] calldata agents, uint256 subnetId)
        external
        view
        returns (AgentInfo[] memory)
    {
        AgentInfo[] memory infos = new AgentInfo[](agents.length);
        for (uint256 i = 0; i < agents.length;) {
            address agent = agents[i];
            address owner = IAccessManager(accessManager).getOwner(agent);
            bool isValid = owner != address(0);
            uint256 stake = isValid
                ? IStakingVault(stakingVault).getAgentStake(owner, agent, subnetId)
                : 0;
            address recipient = isValid
                ? IAccessManager(accessManager).getRewardRecipient(owner)
                : address(0);
            infos[i] = AgentInfo(owner, isValid, stake, recipient);
            unchecked { ++i; }
        }
        return infos;
    }

    // ═══════════════════════════════════════════════
    //  View — general view functions
    // ═══════════════════════════════════════════════

    /// @notice Get subnet lifecycle state (status, lpPool, timestamps)
    /// @param subnetId Subnet ID
    /// @return SubnetInfo struct (RootNet state only; use getSubnetFull for complete data)
    function getSubnet(uint256 subnetId) external view returns (SubnetInfo memory) {
        return subnets[subnetId];
    }

    /// @notice Get complete subnet info combining RootNet state + SubnetNFT identity
    /// @param subnetId Subnet ID
    /// @return Full subnet data including name, addresses, skills, minStake, owner, status
    function getSubnetFull(uint256 subnetId) external view returns (SubnetFullInfo memory) {
        SubnetInfo storage info = subnets[subnetId];
        ISubnetNFT.SubnetData memory nftData = ISubnetNFT(subnetNFT).getSubnetData(subnetId);

        return SubnetFullInfo({
            subnetManager: nftData.subnetManager,
            alphaToken: nftData.alphaToken,
            lpPool: info.lpPool,
            status: info.status,
            createdAt: info.createdAt,
            activatedAt: info.activatedAt,
            name: nftData.name,
            skillsURI: nftData.skillsURI,
            minStake: nftData.minStake,
            owner: nftData.owner
        });
    }

    /// @notice Get the current number of active subnets
    /// @return Active subnet count
    function getActiveSubnetCount() external view returns (uint256) {
        return activeSubnetIds.length();
    }

    /// @notice Get the active subnet ID at a given index
    /// @param index Index position
    /// @return Subnet ID
    function getActiveSubnetIdAt(uint256 index) external view returns (uint256) {
        return activeSubnetIds.at(index);
    }

    /// @notice Check whether a specified subnet is in the Active state
    /// @param subnetId Subnet ID
    /// @return Whether it is Active
    function isSubnetActive(uint256 subnetId) external view returns (bool) {
        return subnets[subnetId].status == SubnetStatus.Active;
    }

    /// @notice Get the next subnet ID to be assigned
    /// @return Next subnet ID
    function nextSubnetId() external view returns (uint256) {
        return _nextSubnetId;
    }

    // ═══════════════════════════════════════════════
    //  Pause — emergency pause
    // ═══════════════════════════════════════════════

    /// @notice Emergency pause the contract (only Guardian may call)
    /// @dev After pausing, all functions with the whenNotPaused modifier will be uncallable
    function pause() external onlyGuardian {
        _pause();
    }

    /// @notice Unpause the contract (only Timelock may call)
    /// @dev Pause is triggered by the Guardian; resumption requires Timelock governance, ensuring dual approval
    function unpause() external onlyTimelock {
        _unpause();
    }
}
