// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {EIP712Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

import {IAlphaToken} from "./interfaces/IAlphaToken.sol";
import {IAlphaTokenFactory} from "./interfaces/IAlphaTokenFactory.sol";
import {IStakingVault} from "./interfaces/IStakingVault.sol";
import {IStakeNFT} from "./interfaces/IStakeNFT.sol";
import {ISubnetNFT} from "./interfaces/ISubnetNFT.sol";
import {ILPManager} from "./interfaces/ILPManager.sol";
import {IAWPRegistry} from "./interfaces/IAWPRegistry.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

/// @title AWPRegistry — Unified entry point for the AWP protocol (subnet management + staking management)
/// @author AWP Team
/// @notice Account System V2: tree-based binding, optional registration, explicit staker parameter.
/// @dev Inheritance: Initializable + UUPSUpgradeable (UUPS proxy pattern), PausableUpgradeable (emergency pause),
///      ReentrancyGuardUpgradeable (reentrancy protection), EIP712Upgradeable (EIP-712 signing domain, domain name "AWPRegistry" v1).
contract AWPRegistry is Initializable, UUPSUpgradeable, PausableUpgradeable, ReentrancyGuardUpgradeable, EIP712Upgradeable, IAWPRegistry {
    using SafeERC20 for IERC20;
    using EnumerableSet for EnumerableSet.UintSet;

    // ══════════════════════════════════════════════
    //  Address registry — external module addresses injected via initializeRegistry after deployment
    // ══════════════════════════════════════════════

    /// @notice AWP token contract address (ERC20, 10B MAX_SUPPLY, 200M pre-minted, remainder via AWPEmission)
    address public awpToken;
    /// @notice SubnetNFT contract address; each subnet corresponds to one NFT (tokenId = subnetId)
    address public subnetNFT;
    /// @notice AlphaToken factory contract address, used to deploy an independent Alpha token for each subnet
    address public alphaTokenFactory;
    /// @notice AWP emission contract address, the only contract holding AWP minting rights
    address public awpEmission;
    /// @notice LP manager contract address, responsible for creating AWP/Alpha trading pairs and managing liquidity
    address public lpManager;
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
    /// @notice ABI-encoded DEX configuration passed to SubnetManager.initialize()
    ///         (clPoolManager, clPositionManager, clSwapRouter, permit2, poolFee, tickSpacing)
    bytes public dexConfig;

    /// @dev Deployer address; used only for initializeRegistry, zeroed immediately after the call
    address private _deployer;
    /// @notice Whether the registry has been initialized (can only be initialized once)
    bool public registryInitialized;

    // ══════════════════════════════════════════════
    //  Account System V2 — binding, recipient, delegation
    // ══════════════════════════════════════════════

    /// @notice Total number of registered users (incremented by register/registerFor, never decremented)
    uint256 public registeredCount;

    /// @notice Tree-based binding: addr → target (address(0) = root / unbound)
    mapping(address => address) public boundTo;
    /// @notice Reward recipient: addr → recipient (address(0) = unregistered, self-fallback)
    mapping(address => address) public recipient;
    /// @notice Delegation: staker → delegate → authorized
    mapping(address => mapping(address => bool)) public delegates;

    // ══════════════════════════════════════════════
    //  Subnet data
    // ══════════════════════════════════════════════

    /// @notice subnetId => SubnetInfo mapping, stores the on-chain state of each subnet
    mapping(uint256 => SubnetInfo) public subnets;

    /// @dev Next local counter for subnet ID generation, auto-increments from 1.
    /// Global subnetId = (block.chainid << 64) | _nextLocalId
    uint256 private _nextLocalId;

    /// @dev Active subnet ID set
    EnumerableSet.UintSet private activeSubnetIds;

    /// @notice Maximum number of active subnets
    uint128 public constant MAX_ACTIVE_SUBNETS = 10000;

    /// @notice Initial price of the Alpha token when registering a subnet (denominated in AWP), default 0.01 AWP
    uint256 public initialAlphaPrice;

    /// @notice Initial Alpha token mint amount per subnet: 100 million (100_000_000 * 1e18)
    /// @notice Alpha token mint amount per subnet registration (governance-settable)
    uint256 public initialAlphaMint;

    /// @notice Subnet deregistration immunity period; Timelock cannot deregister the subnet during this window
    uint256 public immunityPeriod;

    // ══════════════════════════════════════════════
    //  Gasless — EIP-712 signature related
    // ══════════════════════════════════════════════

    /// @notice Per-signer nonce for replay attack prevention (nonces[signer]++ when validated)
    mapping(address => uint256) public nonces;

    /// @dev Reserved storage gap for future upgrades (UUPS pattern)
    uint256[49] private __gap;

    /// @dev EIP-712 type hash: Bind(address agent, address target, uint256 nonce, uint256 deadline)
    bytes32 private constant BIND_TYPEHASH =
        keccak256("Bind(address agent,address target,uint256 nonce,uint256 deadline)");

    /// @dev EIP-712 type hash: SetRecipient(address user, address recipient, uint256 nonce, uint256 deadline)
    bytes32 private constant SET_RECIPIENT_TYPEHASH =
        keccak256("SetRecipient(address user,address recipient,uint256 nonce,uint256 deadline)");

    /// @dev EIP-712 type hash: ActivateSubnet(address user, uint256 subnetId, uint256 nonce, uint256 deadline)
    bytes32 private constant ACTIVATE_SUBNET_TYPEHASH =
        keccak256("ActivateSubnet(address user,uint256 subnetId,uint256 nonce,uint256 deadline)");

    /// @dev EIP-712 type hash for SubnetParams nested struct
    bytes32 private constant SUBNET_PARAMS_TYPEHASH =
        keccak256("SubnetParams(string name,string symbol,address subnetManager,bytes32 salt,uint128 minStake,string skillsURI)");

    /// @dev EIP-712 type hash: RegisterSubnet with nested SubnetParams (per EIP-712 §encodeType)
    bytes32 private constant REGISTER_SUBNET_TYPEHASH =
        keccak256("RegisterSubnet(address user,SubnetParams params,uint256 nonce,uint256 deadline)SubnetParams(string name,string symbol,address subnetManager,bytes32 salt,uint128 minStake,string skillsURI)");

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
    /// @dev Caller is not the Treasury (Timelock)
    error NotTimelock();
    /// @dev Caller is not the Guardian
    error NotGuardian();
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
    /// @dev Invalid address (zero address or self-bind)
    error InvalidAddress();
    /// @dev Binding would create a cycle in the tree
    error CycleDetected();
    /// @dev Binding chain is too long (safety limit)
    error ChainTooLong();
    /// @dev Cannot revoke self as delegate
    error CannotRevokeSelf();
    /// @dev Address is already registered
    error AlreadyRegistered();

    // ══════════════════════════════════════════════
    //  Constructor
    // ══════════════════════════════════════════════

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Initialize the registry (called once via proxy)
    /// @param deployer_ Deployer address (holds initializeRegistry rights)
    /// @param treasury_ Treasury (Timelock) address
    /// @param guardian_ Guardian address
    function initialize(
        address deployer_,
        address treasury_,
        address guardian_
    ) external initializer {
        __UUPSUpgradeable_init();
        __Pausable_init();
        __ReentrancyGuard_init();
        __EIP712_init("AWPRegistry", "1");

        _deployer = deployer_;
        treasury = treasury_;
        guardian = guardian_;

        // Inline initializer values (not executed for proxy storage)
        _nextLocalId = 1;
        initialAlphaPrice = 1e16;
        initialAlphaMint = 100_000_000 * 1e18;
        immunityPeriod = 30 days;
    }

    /// @dev UUPS upgrade authorization — only Timelock may upgrade
    function _authorizeUpgrade(address) internal override onlyTimelock {}

    // ═══════════════════════════════════════════════
    //  Registry — module address registry
    // ═══════════════════════════════════════════════

    /// @notice One-time initialization of all external module addresses; only callable by deployer, and only once
    /// @dev After successful call _deployer is zeroed, permanently locked; cannot be called again.
    function initializeRegistry(
        address awpToken_,
        address subnetNFT_,
        address alphaTokenFactory_,
        address awpEmission_,
        address lpManager_,
        address stakingVault_,
        address stakeNFT_,
        address defaultSubnetManagerImpl_,
        bytes calldata dexConfig_
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
        stakingVault = stakingVault_;
        stakeNFT = stakeNFT_;
        defaultSubnetManagerImpl = defaultSubnetManagerImpl_;
        dexConfig = dexConfig_;

        // Link StakingVault → StakeNFT (one-time setter, resolves CREATE2 circular dependency)
        IStakingVault(stakingVault).setStakeNFT(stakeNFT);

        registryInitialized = true;
        // Permanently destroy deployer rights; this function can no longer be called
        _deployer = address(0);
    }

    /// @notice Retrieve all module addresses in a single call
    /// @return In order: awpToken, subnetNFT, alphaTokenFactory, awpEmission,
    ///         lpManager, stakingVault, stakeNFT, treasury, guardian
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
            address
        )
    {
        return (
            awpToken,
            subnetNFT,
            alphaTokenFactory,
            awpEmission,
            lpManager,
            stakingVault,
            stakeNFT,
            treasury,
            guardian
        );
    }

    // ═══════════════════════════════════════════════
    //  Account V2: Registration
    // ═══════════════════════════════════════════════

    /// @notice Self-register: sets recipient to self and increments registeredCount
    function register() external nonReentrant whenNotPaused {
        if (recipient[msg.sender] != address(0)) revert AlreadyRegistered();
        recipient[msg.sender] = msg.sender;
        registeredCount++;
        emit UserRegistered(msg.sender);
    }

    // ═══════════════════════════════════════════════
    //  Account V2: Binding (tree structure)
    // ═══════════════════════════════════════════════

    /// @notice Bind msg.sender to target in the tree structure
    /// @dev Anti-cycle: walk up from target, if we find msg.sender → cycle
    /// @param target Address to bind to
    function bind(address target) external nonReentrant whenNotPaused {
        if (target == address(0)) revert InvalidAddress();
        if (target == msg.sender) revert InvalidAddress();
        _checkCycle(msg.sender, target);
        // First-time interaction: count as registered
        if (boundTo[msg.sender] == address(0) && recipient[msg.sender] == address(0)) {
            registeredCount++;
            emit UserRegistered(msg.sender);
        }
        boundTo[msg.sender] = target;
        emit Bound(msg.sender, target);
    }

    /// @notice Unbind msg.sender from the tree
    function unbind() external nonReentrant whenNotPaused {
        boundTo[msg.sender] = address(0);
        emit Unbound(msg.sender);
    }

    /// @notice Gasless bind: relayer pays gas, agent signs EIP-712
    /// @param agent Agent address (signer)
    /// @param target Address to bind to
    /// @param deadline Signature expiry time
    /// @param v EIP-712 signature v value
    /// @param r EIP-712 signature r value
    /// @param s EIP-712 signature s value
    function bindFor(address agent, address target, uint256 deadline, uint8 v, bytes32 r, bytes32 s)
        external
        nonReentrant
        whenNotPaused
    {
        if (target == address(0)) revert InvalidAddress();
        if (target == agent) revert InvalidAddress();
        _verifyDigest(agent, keccak256(abi.encode(BIND_TYPEHASH, agent, target, nonces[agent]++, deadline)), deadline, v, r, s);

        _checkCycle(agent, target);
        if (boundTo[agent] == address(0) && recipient[agent] == address(0)) {
            registeredCount++;
            emit UserRegistered(agent);
        }
        boundTo[agent] = target;
        emit Bound(agent, target);
    }

    // ═══════════════════════════════════════════════
    //  Account V2: Recipient
    // ═══════════════════════════════════════════════

    /// @notice Set reward recipient for msg.sender
    /// @param addr Recipient address
    function setRecipient(address addr) external nonReentrant whenNotPaused {
        if (addr == address(0)) revert InvalidAddress();
        bool firstTime = recipient[msg.sender] == address(0);
        recipient[msg.sender] = addr;
        if (firstTime) {
            registeredCount++;
            emit UserRegistered(msg.sender);
        }
        emit RecipientSet(msg.sender, addr);
    }

    /// @notice Gasless set recipient: relayer pays gas, user signs EIP-712
    function setRecipientFor(
        address user, address _recipient, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant whenNotPaused {
        if (_recipient == address(0)) revert InvalidAddress();
        _verifyDigest(user, keccak256(abi.encode(SET_RECIPIENT_TYPEHASH, user, _recipient, nonces[user]++, deadline)), deadline, v, r, s);

        // If this is the first time setting recipient (registration), increment counter
        bool firstTime = recipient[user] == address(0);
        recipient[user] = _recipient;
        if (firstTime) {
            registeredCount++;
            emit UserRegistered(user);
        }
        emit RecipientSet(user, _recipient);
    }

    /// @notice Resolve the reward recipient for an address by walking the binding tree to the root
    /// @param addr Address to resolve
    /// @return The reward recipient address
    function resolveRecipient(address addr) external view returns (address) {
        address cur = addr;
        uint256 depth = 0;
        while (boundTo[cur] != address(0) && boundTo[cur] != cur) {
            cur = boundTo[cur];
            if (++depth >= 100) break;
        }
        return recipient[cur] != address(0) ? recipient[cur] : cur;
    }

    /// @notice Batch resolve recipients for multiple addresses (view, no gas cost for callers)
    /// @param addrs Array of addresses to resolve
    /// @return resolved Array of resolved recipient addresses (same order as input)
    function batchResolveRecipients(address[] calldata addrs) external view returns (address[] memory resolved) {
        resolved = new address[](addrs.length);
        for (uint256 i = 0; i < addrs.length;) {
            address cur = addrs[i];
            uint256 depth = 0;
            while (boundTo[cur] != address(0) && boundTo[cur] != cur) {
                cur = boundTo[cur];
                if (++depth >= 100) break;
            }
            resolved[i] = recipient[cur] != address(0) ? recipient[cur] : cur;
            unchecked { ++i; }
        }
    }

    /// @notice Check if an address is registered (has a non-zero recipient)
    /// @param addr Address to check
    /// @return true if registered
    function isRegistered(address addr) external view returns (bool) {
        return boundTo[addr] != address(0) || recipient[addr] != address(0);
    }

    // ═══════════════════════════════════════════════
    //  Account V2: Delegation
    // ═══════════════════════════════════════════════

    /// @notice Grant delegate authorization to another address
    /// @param delegate Address to authorize
    function grantDelegate(address delegate) external whenNotPaused {
        delegates[msg.sender][delegate] = true;
        emit DelegateGranted(msg.sender, delegate);
    }

    /// @notice Revoke delegate authorization
    /// @param delegate Address to de-authorize
    function revokeDelegate(address delegate) external whenNotPaused {
        if (delegate == msg.sender) revert CannotRevokeSelf();
        delegates[msg.sender][delegate] = false;
        emit DelegateRevoked(msg.sender, delegate);
    }

    // ═══════════════════════════════════════════════
    //  Subnet Registration
    // ═══════════════════════════════════════════════

    /// @notice Register a new subnet: deploy Alpha token, create LP, mint NFT
    /// @param params Subnet parameters (name, symbol, subnetManager, salt, minStake)
    /// @return subnetId Newly created subnet ID
    function registerSubnet(SubnetParams calldata params) external nonReentrant whenNotPaused returns (uint256) {
        return _registerSubnet(msg.sender, params);
    }

    /// @notice Gasless subnet registration: relayer pays gas, user signs EIP-712 and pays AWP
    function registerSubnetFor(
        address user,
        SubnetParams calldata params,
        uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant whenNotPaused returns (uint256) {
        _verifyRegisterSubnetSignature(user, params, deadline, v, r, s);
        return _registerSubnet(user, params);
    }

    /// @notice Fully gasless subnet registration with EIP-2612 permit (no prior approve tx needed)
    function registerSubnetForWithPermit(
        address user,
        SubnetParams calldata params,
        uint256 deadline,
        uint8 permitV, bytes32 permitR, bytes32 permitS,
        uint8 registerV, bytes32 registerR, bytes32 registerS
    ) external nonReentrant whenNotPaused returns (uint256) {
        uint256 lpAWPAmount = initialAlphaMint * initialAlphaPrice / 1e18;
        IERC20Permit(awpToken).permit(user, address(this), lpAWPAmount, deadline, permitV, permitR, permitS);
        _verifyRegisterSubnetSignature(user, params, deadline, registerV, registerR, registerS);
        return _registerSubnet(user, params);
    }

    /// @dev Verify EIP-712 signature for registerSubnetFor
    function _verifyRegisterSubnetSignature(
        address user, SubnetParams calldata params, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) internal {
        bytes32 paramsStructHash = keccak256(abi.encode(
            SUBNET_PARAMS_TYPEHASH,
            keccak256(bytes(params.name)),
            keccak256(bytes(params.symbol)),
            params.subnetManager, params.salt, params.minStake,
            keccak256(bytes(params.skillsURI))
        ));
        _verifyDigest(user, keccak256(abi.encode(REGISTER_SUBNET_TYPEHASH, user, paramsStructHash, nonces[user]++, deadline)), deadline, v, r, s);
    }

    /// @dev Internal: shared logic for registerSubnet and registerSubnetFor
    function _registerSubnet(address user, SubnetParams calldata params) internal returns (uint256) {
        if (bytes(params.name).length == 0 || bytes(params.name).length > 64) revert InvalidSubnetParams();
        if (bytes(params.symbol).length == 0 || bytes(params.symbol).length > 16) revert InvalidSubnetParams();
        bool autoDeploySubnet = params.subnetManager == address(0);
        if (autoDeploySubnet && defaultSubnetManagerImpl == address(0)) revert SubnetManagerRequired();

        uint256 lpAWPAmount = initialAlphaMint * initialAlphaPrice / 1e18;
        IERC20(awpToken).safeTransferFrom(user, lpManager, lpAWPAmount);

        uint256 subnetId = (block.chainid << 64) | _nextLocalId++;

        (address alphaToken, bytes32 poolId) = _deployAlphaAndLP(
            subnetId, params.name, params.symbol, lpAWPAmount, params.salt
        );

        address sc;
        if (autoDeploySubnet) {
            bytes memory initData = abi.encodeWithSignature(
                "initialize(address,address,address,bytes32,address,bytes)",
                address(this), alphaToken, awpToken, poolId, user, dexConfig
            );
            sc = address(new ERC1967Proxy(defaultSubnetManagerImpl, initData));
        } else {
            sc = params.subnetManager;
        }

        IAlphaToken(alphaToken).setSubnetMinter(sc);
        ISubnetNFT(subnetNFT).mint(user, subnetId, params.name, sc, alphaToken, params.minStake, params.skillsURI);
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
        IAlphaToken(alphaToken).mint(lpManager, initialAlphaMint);
        (poolId,) = ILPManager(lpManager).createPoolAndAddLiquidity(alphaToken, lpAWPAmount, initialAlphaMint);
    }

    /// @dev Emit subnet registration events
    function _emitSubnetRegistered(
        uint256 subnetId, address user, address sc, address alphaToken, bytes32 poolId, uint256 lpAWPAmount, SubnetParams calldata params
    ) internal {
        emit SubnetRegistered(
            subnetId, user,
            params.name, params.symbol,
            sc, alphaToken
        );
        emit LPCreated(subnetId, poolId, lpAWPAmount, initialAlphaMint);
    }

    // ═══════════════════════════════════════════════
    //  Subnet Lifecycle Management
    // ═══════════════════════════════════════════════

    /// @notice Activate a subnet: Pending → Active (only the NFT Owner may call)
    function activateSubnet(uint256 subnetId) external nonReentrant whenNotPaused {
        if (ISubnetNFT(subnetNFT).ownerOf(subnetId) != msg.sender) revert NotOwner();
        _activateSubnet(subnetId);
    }

    /// @notice Gasless activate subnet: relayer pays gas, NFT owner signs EIP-712
    function activateSubnetFor(
        address user, uint256 subnetId, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant whenNotPaused {
        _verifyDigest(user, keccak256(abi.encode(ACTIVATE_SUBNET_TYPEHASH, user, subnetId, nonces[user]++, deadline)), deadline, v, r, s);
        if (ISubnetNFT(subnetNFT).ownerOf(subnetId) != user) revert NotOwner();
        _activateSubnet(subnetId);
    }

    /// @dev Shared activation logic
    function _activateSubnet(uint256 subnetId) internal {
        SubnetInfo storage info = subnets[subnetId];
        if (info.status != SubnetStatus.Pending) revert InvalidSubnetStatus();
        if (activeSubnetIds.length() >= MAX_ACTIVE_SUBNETS) revert MaxActiveSubnetsReached();
        info.status = SubnetStatus.Active;
        info.activatedAt = uint64(block.timestamp);
        activeSubnetIds.add(subnetId);
        emit SubnetActivated(subnetId);
    }

    /// @notice Pause a subnet: Active → Paused (only the NFT Owner may call)
    function pauseSubnet(uint256 subnetId) external nonReentrant whenNotPaused {
        if (ISubnetNFT(subnetNFT).ownerOf(subnetId) != msg.sender) revert NotOwner();
        SubnetInfo storage info = subnets[subnetId];
        if (info.status != SubnetStatus.Active) revert InvalidSubnetStatus();

        info.status = SubnetStatus.Paused;
        activeSubnetIds.remove(subnetId);

        emit SubnetPaused(subnetId);
    }

    /// @notice Resume a subnet: Paused → Active (only the NFT Owner may call)
    function resumeSubnet(uint256 subnetId) external nonReentrant whenNotPaused {
        if (ISubnetNFT(subnetNFT).ownerOf(subnetId) != msg.sender) revert NotOwner();
        SubnetInfo storage info = subnets[subnetId];
        if (info.status != SubnetStatus.Paused) revert InvalidSubnetStatus();

        if (activeSubnetIds.length() >= MAX_ACTIVE_SUBNETS) revert MaxActiveSubnetsReached();
        info.status = SubnetStatus.Active;
        activeSubnetIds.add(subnetId);

        emit SubnetResumed(subnetId);
    }

    /// @notice Ban a subnet: Active/Paused → Banned (only Timelock may call)
    function banSubnet(uint256 subnetId) external onlyTimelock {
        SubnetInfo storage info = subnets[subnetId];
        SubnetStatus status = info.status;
        if (status != SubnetStatus.Active && status != SubnetStatus.Paused) revert InvalidSubnetStatus();

        ISubnetNFT.SubnetData memory sd = ISubnetNFT(subnetNFT).getSubnetData(subnetId);
        if (sd.subnetManager != address(0)) {
            IAlphaToken(sd.alphaToken).setMinterPaused(sd.subnetManager, true);
        }
        if (status == SubnetStatus.Active) {
            activeSubnetIds.remove(subnetId);
        }
        info.status = SubnetStatus.Banned;

        emit SubnetBanned(subnetId);
    }

    /// @notice Unban a subnet: Banned → Active (only Timelock may call)
    function unbanSubnet(uint256 subnetId) external onlyTimelock {
        SubnetInfo storage info = subnets[subnetId];
        if (info.status != SubnetStatus.Banned) revert InvalidSubnetStatus();

        ISubnetNFT.SubnetData memory sd = ISubnetNFT(subnetNFT).getSubnetData(subnetId);
        if (sd.subnetManager != address(0)) {
            IAlphaToken(sd.alphaToken).setMinterPaused(sd.subnetManager, false);
        }
        if (activeSubnetIds.length() >= MAX_ACTIVE_SUBNETS) revert MaxActiveSubnetsReached();
        info.status = SubnetStatus.Active;
        activeSubnetIds.add(subnetId);

        emit SubnetUnbanned(subnetId);
    }

    /// @notice Deregister a subnet: permanently delete subnet data (only Timelock; immunity period must have elapsed)
    function deregisterSubnet(uint256 subnetId) external onlyTimelock {
        SubnetInfo storage info = subnets[subnetId];
        if (info.status != SubnetStatus.Banned) revert InvalidSubnetStatus();
        uint256 immunityStart = info.activatedAt > 0 ? uint256(info.activatedAt) : uint256(info.createdAt);
        if (block.timestamp <= immunityStart + immunityPeriod) revert ImmunityNotExpired();

        ISubnetNFT.SubnetData memory sd = ISubnetNFT(subnetNFT).getSubnetData(subnetId);
        if (sd.subnetManager != address(0)) {
            IAlphaToken(sd.alphaToken).setMinterPaused(sd.subnetManager, true);
        }
        activeSubnetIds.remove(subnetId);
        delete subnets[subnetId];
        ISubnetNFT(subnetNFT).burn(subnetId);

        emit SubnetDeregistered(subnetId);
    }

    // ═══════════════════════════════════════════════
    //  Subnet Parameters
    // ═══════════════════════════════════════════════

    /// @notice Set the initial Alpha price when registering a subnet (only Timelock may call)
    function setInitialAlphaPrice(uint256 price) external onlyTimelock {
        if (price < 1e12) revert PriceTooLow();
        if (price > 1e30) revert PriceTooHigh();
        initialAlphaPrice = price;
        emit InitialAlphaPriceUpdated(price);
    }

    /// @notice Update the initial Alpha token mint amount per subnet (only Timelock)
    function setInitialAlphaMint(uint256 amount) external onlyTimelock {
        if (amount == 0) revert InvalidSubnetParams();
        initialAlphaMint = amount;
        emit InitialAlphaMintUpdated(amount);
    }

    /// @notice Update the guardian address (only Guardian may call — self-sovereign)
    /// @dev Single-chain DAO cannot change guardian. If Guardian keys lost, recover via UUPS upgrade.
    function setGuardian(address g) external onlyGuardian {
        if (g == address(0)) revert InvalidAddress();
        guardian = g;
        emit GuardianUpdated(g);
    }

    /// @notice Set the subnet deregistration immunity period (only Timelock may call)
    function setImmunityPeriod(uint256 p) external onlyTimelock {
        if (p < 7 days) revert InvalidSubnetParams();
        immunityPeriod = p;
        emit ImmunityPeriodUpdated(p);
    }

    /// @notice Replace the AlphaToken factory (only Timelock may call)
    function setAlphaTokenFactory(address factory) external onlyTimelock {
        if (factory == address(0)) revert InvalidAddress();
        alphaTokenFactory = factory;
        emit AlphaTokenFactoryUpdated(factory);
    }

    /// @notice Set the default subnet implementation (only Timelock may call)
    function setSubnetManagerImpl(address impl) external onlyTimelock {
        if (impl == address(0)) revert InvalidAddress();
        defaultSubnetManagerImpl = impl;
        emit DefaultSubnetManagerImplUpdated(impl);
    }

    /// @notice Update DEX configuration for future auto-deployed SubnetManagers (only Timelock)
    function setDexConfig(bytes calldata dexConfig_) external onlyTimelock {
        dexConfig = dexConfig_;
        emit DexConfigUpdated();
    }

    // ═══════════════════════════════════════════════
    //  Queries
    // ═══════════════════════════════════════════════

    /// @dev Agent information aggregate struct, used for off-chain queries
    struct AgentInfo {
        /// @dev Resolved root address (walks binding chain)
        address root;
        /// @dev Whether the agent has a binding or the root has a recipient set
        bool isValid;
        /// @dev Agent's stake amount on the specified subnet (uses root as staker)
        uint256 stake;
        /// @dev Reward recipient address for the root
        address rewardRecipient;
    }

    /// @notice Query complete information for a single agent on a specified subnet
    /// @param agent Agent address
    /// @param subnetId Subnet ID
    /// @return AgentInfo containing root, isValid, stake, rewardRecipient
    function getAgentInfo(address agent, uint256 subnetId) external view returns (AgentInfo memory) {
        // Walk the binding tree to find the root
        address root = _resolveRoot(agent);
        bool isValid = boundTo[agent] != address(0) || recipient[agent] != address(0);
        uint256 stake = IStakingVault(stakingVault).getAgentStake(root, agent, subnetId);
        address recip = recipient[root] != address(0) ? recipient[root] : root;
        return AgentInfo(root, isValid, stake, recip);
    }

    /// @notice Batch query information for multiple agents on a specified subnet
    function getAgentsInfo(address[] calldata agents, uint256 subnetId)
        external
        view
        returns (AgentInfo[] memory)
    {
        AgentInfo[] memory infos = new AgentInfo[](agents.length);
        for (uint256 i = 0; i < agents.length;) {
            address agent = agents[i];
            address root = _resolveRoot(agent);
            bool isValid = boundTo[agent] != address(0) || recipient[agent] != address(0);
            uint256 stake = isValid
                ? IStakingVault(stakingVault).getAgentStake(root, agent, subnetId)
                : 0;
            address recip = isValid
                ? (recipient[root] != address(0) ? recipient[root] : root)
                : address(0);
            infos[i] = AgentInfo(root, isValid, stake, recip);
            unchecked { ++i; }
        }
        return infos;
    }

    /// @dev Check deadline, build EIP-712 digest, verify signer. Reverts on failure.
    function _verifyDigest(address expectedSigner, bytes32 structHash, uint256 deadline, uint8 v, bytes32 r, bytes32 s) internal view {
        if (block.timestamp > deadline) revert ExpiredSignature();
        bytes32 digest = _hashTypedDataV4(structHash);
        if (ECDSA.recover(digest, v, r, s) != expectedSigner) revert InvalidSignature();
    }

    /// @dev Anti-cycle check: walk up from target, revert if sender is found in the chain
    function _checkCycle(address sender, address target) internal view {
        address cur = target;
        uint256 depth = 0;
        while (boundTo[cur] != address(0) && boundTo[cur] != cur) {
            if (boundTo[cur] == sender) revert CycleDetected();
            cur = boundTo[cur];
            if (++depth >= 100) revert ChainTooLong();
        }
    }

    /// @dev Walk the binding chain to find the root
    function _resolveRoot(address addr) internal view returns (address) {
        address cur = addr;
        uint256 depth = 0;
        while (boundTo[cur] != address(0) && boundTo[cur] != cur) {
            cur = boundTo[cur];
            if (++depth >= 100) break;
        }
        return cur;
    }

    // ═══════════════════════════════════════════════
    //  View — general view functions
    // ═══════════════════════════════════════════════

    /// @notice Get subnet lifecycle state
    function getSubnet(uint256 subnetId) external view returns (SubnetInfo memory) {
        return subnets[subnetId];
    }

    /// @notice Get complete subnet info combining AWPRegistry state + SubnetNFT identity
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
    function getActiveSubnetCount() external view returns (uint256) {
        return activeSubnetIds.length();
    }

    /// @notice Get the active subnet ID at a given index
    function getActiveSubnetIdAt(uint256 index) external view returns (uint256) {
        return activeSubnetIds.at(index);
    }

    /// @notice Check whether a specified subnet is in the Active state
    function isSubnetActive(uint256 subnetId) external view returns (bool) {
        return subnets[subnetId].status == SubnetStatus.Active;
    }

    /// @notice Get the next subnet ID to be assigned (globally unique: chainId << 64 | localCounter)
    function nextSubnetId() external view returns (uint256) {
        return (block.chainid << 64) | _nextLocalId;
    }

    /// @notice Extract chainId from a global subnetId
    function extractChainId(uint256 subnetId) external pure returns (uint256) {
        return subnetId >> 64;
    }

    /// @notice Extract local counter from a global subnetId
    function extractLocalId(uint256 subnetId) external pure returns (uint256) {
        return subnetId & ((1 << 64) - 1);
    }

    // ═══════════════════════════════════════════════
    //  Pause — emergency pause
    // ═══════════════════════════════════════════════

    /// @notice Emergency pause the contract (only Guardian may call)
    function pause() external onlyGuardian {
        _pause();
    }

    /// @notice Unpause the contract (only Timelock may call)
    function unpause() external onlyTimelock {
        _unpause();
    }
}
