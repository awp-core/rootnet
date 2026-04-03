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
import {IWorknetNFT} from "./interfaces/IWorknetNFT.sol";
import {ILPManager} from "./interfaces/ILPManager.sol";
import {IAWPRegistry} from "./interfaces/IAWPRegistry.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

/// @title AWPRegistry — Unified entry point for the AWP protocol (worknet management + staking management)
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

    /// @notice AWP token contract address (ERC20, 10B MAX_SUPPLY, configurable initialMint, remainder via AWPEmission)
    address public awpToken;
    /// @notice WorknetNFT contract address; each worknet corresponds to one NFT (tokenId = worknetId)
    address public worknetNFT;
    /// @notice AlphaToken factory contract address, used to deploy an independent Alpha token for each worknet
    address public alphaTokenFactory;
    /// @notice AWP emission contract address, the only contract holding AWP minting rights
    address public awpEmission;
    /// @notice LP manager contract address, responsible for creating AWP/Alpha trading pairs and managing liquidity
    address public lpManager;
    /// @notice Staking vault contract address — manages allocation bookkeeping
    address public stakingVault;
    /// @notice StakeNFT contract address — manages AWP staking positions as NFTs
    address public stakeNFT;
    /// @notice Treasury (Timelock) address — holds governance operation rights (onlyGuardian)
    address public treasury;
    /// @notice Guardian address — holds emergency pause rights (onlyGuardian)
    address public guardian;
    /// @notice Default worknet implementation address (for auto-deploying worknet contracts via proxy)
    address public defaultWorknetManagerImpl;
    /// @notice ABI-encoded DEX configuration passed to WorknetManager.initialize()
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
    //  Worknet data
    // ══════════════════════════════════════════════

    /// @notice worknetId => WorknetInfo mapping, stores the on-chain state of each worknet
    mapping(uint256 => WorknetInfo) public worknets;

    /// @dev Next local counter for worknet ID generation, auto-increments from 1.
    /// Global worknetId = (block.chainid << 64) | _nextLocalId
    uint256 private _nextLocalId;

    /// @dev Active worknet ID set
    EnumerableSet.UintSet private activeWorknetIds;

    /// @notice Maximum number of active worknets
    uint128 public constant MAX_ACTIVE_WORKNETS = 10000;

    /// @notice Initial price of the Alpha token when registering a worknet (denominated in AWP), default 0.01 AWP
    uint256 public initialAlphaPrice;

    /// @notice Alpha token mint amount per worknet registration (governance-settable, default 100M)
    uint256 public initialAlphaMint;

    /// @notice Worknet deregistration immunity period; Timelock cannot deregister the worknet during this window
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

    /// @dev EIP-712 type hash: ActivateWorknet(address user, uint256 worknetId, uint256 nonce, uint256 deadline)
    bytes32 private constant ACTIVATE_WORKNET_TYPEHASH =
        keccak256("ActivateWorknet(address user,uint256 worknetId,uint256 nonce,uint256 deadline)");

    /// @dev EIP-712 type hash for WorknetParams nested struct
    bytes32 private constant WORKNET_PARAMS_TYPEHASH =
        keccak256("WorknetParams(string name,string symbol,address worknetManager,bytes32 salt,uint128 minStake,string skillsURI)");

    /// @dev EIP-712 type hash: RegisterWorknet with nested WorknetParams (per EIP-712 §encodeType)
    bytes32 private constant REGISTER_WORKNET_TYPEHASH =
        keccak256("RegisterWorknet(address user,WorknetParams params,uint256 nonce,uint256 deadline)WorknetParams(string name,string symbol,address worknetManager,bytes32 salt,uint128 minStake,string skillsURI)");

    /// @dev EIP-712 type hash: GrantDelegate(address user, address delegate, uint256 nonce, uint256 deadline)
    bytes32 private constant GRANT_DELEGATE_TYPEHASH =
        keccak256("GrantDelegate(address user,address delegate,uint256 nonce,uint256 deadline)");

    /// @dev EIP-712 type hash: RevokeDelegate(address user, address delegate, uint256 nonce, uint256 deadline)
    bytes32 private constant REVOKE_DELEGATE_TYPEHASH =
        keccak256("RevokeDelegate(address user,address delegate,uint256 nonce,uint256 deadline)");

    /// @dev EIP-712 type hash: Unbind(address user, uint256 nonce, uint256 deadline)
    bytes32 private constant UNBIND_TYPEHASH =
        keccak256("Unbind(address user,uint256 nonce,uint256 deadline)");

    // ══════════════════════════════════════════════
    //  Permission modifiers
    // ══════════════════════════════════════════════

    /// @dev Only the Guardian may call
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
    /// @dev Caller is not the Guardian
    error NotGuardian();
    /// @dev Worknet name is invalid (empty or > 64 bytes)
    error InvalidWorknetName();
    /// @dev Worknet symbol is invalid (empty or > 16 bytes)
    error InvalidWorknetSymbol();
    /// @dev Computed LP AWP amount is zero
    error ZeroLPAmount();
    /// @dev Mint amount must be > 0
    error InvalidMintAmount();
    /// @dev Immunity period too short (minimum 7 days)
    error ImmunityTooShort();
    /// @dev Worknet contract address cannot be the zero address
    error WorknetManagerRequired();
    /// @dev Worknet status does not meet the precondition
    error InvalidWorknetStatus(uint256 worknetId, uint8 currentStatus);
    /// @dev Worknet is still within its immunity period and cannot be deregistered
    error ImmunityNotExpired();
    /// @dev Caller is not the worknet NFT owner
    error NotOwner();
    /// @dev Initial Alpha price is too low (minimum 1e12)
    error PriceTooLow();
    /// @dev EIP-712 signature has expired (block.timestamp > deadline)
    error ExpiredSignature();
    /// @dev EIP-712 signature verification failed (recovered signer does not match expected)
    error InvalidSignature();
    /// @dev Active worknet count has reached the maximum
    error MaxActiveWorknetsReached();
    /// @dev Initial Alpha price is too high
    error PriceTooHigh();
    /// @dev Zero address passed
    error ZeroAddress();
    /// @dev Self-bind is not allowed
    error SelfBind();
    /// @dev Binding would create a cycle in the tree
    error CycleDetected();
    /// @dev Binding chain exceeds maximum depth
    error ChainTooDeep();
    /// @dev Cannot revoke self as delegate
    error CannotRevokeSelf();
    /// @dev Name/symbol contains JSON-unsafe characters (" or \)
    error JsonUnsafeCharacter();
    /// @dev Address is already registered

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

        if (deployer_ == address(0) || treasury_ == address(0) || guardian_ == address(0)) revert ZeroAddress();
        _deployer = deployer_;
        treasury = treasury_;
        guardian = guardian_;

        // Default values written to proxy storage during initialize()
        _nextLocalId = 1;
        initialAlphaPrice = 1e15; // 0.001 AWP per Alpha
        initialAlphaMint = 100_000_000 * 1e18;
        immunityPeriod = 30 days;
    }

    /// @dev UUPS upgrade authorization — only Guardian may upgrade
    function _authorizeUpgrade(address) internal override onlyGuardian {}

    // ═══════════════════════════════════════════════
    //  Registry — module address registry
    // ═══════════════════════════════════════════════

    /// @notice One-time initialization of all external module addresses; only callable by deployer, and only once
    /// @dev After successful call _deployer is zeroed, permanently locked; cannot be called again.
    function initializeRegistry(
        address awpToken_,
        address worknetNFT_,
        address alphaTokenFactory_,
        address awpEmission_,
        address lpManager_,
        address stakingVault_,
        address stakeNFT_,
        address defaultWorknetManagerImpl_,
        bytes calldata dexConfig_
    ) external {
        // Only the deployer may call
        if (msg.sender != _deployer) revert NotDeployer();
        // Prevent re-initialization
        if (registryInitialized) revert AlreadyInitialized();
        // Validate critical addresses (cannot re-call after deployer is zeroed)
        if (awpToken_ == address(0) || worknetNFT_ == address(0) || alphaTokenFactory_ == address(0)
            || awpEmission_ == address(0) || lpManager_ == address(0) || stakingVault_ == address(0)
            || stakeNFT_ == address(0)) revert ZeroAddress();

        awpToken = awpToken_;
        worknetNFT = worknetNFT_;
        alphaTokenFactory = alphaTokenFactory_;
        awpEmission = awpEmission_;
        lpManager = lpManager_;
        stakingVault = stakingVault_;
        stakeNFT = stakeNFT_;
        defaultWorknetManagerImpl = defaultWorknetManagerImpl_;
        dexConfig = dexConfig_;

        // Link StakingVault → StakeNFT (one-time setter, resolves CREATE2 circular dependency)
        IStakingVault(stakingVault).setStakeNFT(stakeNFT);

        registryInitialized = true;
        // Permanently destroy deployer rights; this function can no longer be called
        _deployer = address(0);
    }

    /// @notice Retrieve all module addresses in a single call
    /// @return In order: awpToken, worknetNFT, alphaTokenFactory, awpEmission,
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
            worknetNFT,
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

    // ═══════════════════════════════════════════════
    //  Account V2: Binding (tree structure)
    // ═══════════════════════════════════════════════

    /// @notice Bind msg.sender to target in the tree structure
    /// @dev Anti-cycle: walk up from target, if we find msg.sender → cycle
    /// @param target Address to bind to
    function bind(address target) external nonReentrant whenNotPaused {
        if (target == address(0)) revert ZeroAddress();
        if (target == msg.sender) revert SelfBind();
        _checkCycle(msg.sender, target);
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
        if (target == address(0)) revert ZeroAddress();
        if (target == agent) revert SelfBind();
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
        if (addr == address(0)) revert ZeroAddress();
        if (recipient[msg.sender] == address(0) && boundTo[msg.sender] == address(0)) {
            registeredCount++;
            emit UserRegistered(msg.sender);
        }
        recipient[msg.sender] = addr;
        emit RecipientSet(msg.sender, addr);
    }

    /// @notice Gasless set recipient: relayer pays gas, user signs EIP-712
    function setRecipientFor(
        address user, address _recipient, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant whenNotPaused {
        if (_recipient == address(0)) revert ZeroAddress();
        _verifyDigest(user, keccak256(abi.encode(SET_RECIPIENT_TYPEHASH, user, _recipient, nonces[user]++, deadline)), deadline, v, r, s);

        if (recipient[user] == address(0) && boundTo[user] == address(0)) {
            registeredCount++;
            emit UserRegistered(user);
        }
        recipient[user] = _recipient;
        emit RecipientSet(user, _recipient);
    }

    /// @notice Resolve the reward recipient for an address by walking the binding tree to the root
    /// @param addr Address to resolve
    /// @return The reward recipient address
    function resolveRecipient(address addr) external view returns (address) {
        address cur = addr;
        uint256 depth;
        while (boundTo[cur] != address(0) && boundTo[cur] != cur) {
            cur = boundTo[cur];
            if (++depth >= 256) revert ChainTooDeep();
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
            uint256 depth;
            while (boundTo[cur] != address(0) && boundTo[cur] != cur) {
                cur = boundTo[cur];
                if (++depth >= 256) revert ChainTooDeep();
            }
            resolved[i] = recipient[cur] != address(0) ? recipient[cur] : cur;
            unchecked { ++i; }
        }
    }

    /// @notice Check if an address is registered (has a binding or a non-zero recipient)
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

    /// @notice Gasless grant delegate: relayer pays gas, user signs EIP-712
    function grantDelegateFor(
        address user, address delegate, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external whenNotPaused {
        _verifyDigest(user, keccak256(abi.encode(GRANT_DELEGATE_TYPEHASH, user, delegate, nonces[user]++, deadline)), deadline, v, r, s);
        delegates[user][delegate] = true;
        emit DelegateGranted(user, delegate);
    }

    /// @notice Gasless revoke delegate: relayer pays gas, user signs EIP-712
    function revokeDelegateFor(
        address user, address delegate, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external whenNotPaused {
        if (delegate == user) revert CannotRevokeSelf();
        _verifyDigest(user, keccak256(abi.encode(REVOKE_DELEGATE_TYPEHASH, user, delegate, nonces[user]++, deadline)), deadline, v, r, s);
        delegates[user][delegate] = false;
        emit DelegateRevoked(user, delegate);
    }

    /// @notice Gasless unbind: relayer pays gas, user signs EIP-712
    function unbindFor(
        address user, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant whenNotPaused {
        _verifyDigest(user, keccak256(abi.encode(UNBIND_TYPEHASH, user, nonces[user]++, deadline)), deadline, v, r, s);
        boundTo[user] = address(0);
        emit Unbound(user);
    }

    // ═══════════════════════════════════════════════
    //  Worknet Registration
    // ═══════════════════════════════════════════════

    /// @notice Register a new worknet: deploy Alpha token, create LP, mint NFT
    /// @param params Worknet parameters (name, symbol, worknetManager, salt, minStake)
    /// @return worknetId Newly created worknet ID
    function registerWorknet(WorknetParams calldata params) external nonReentrant whenNotPaused returns (uint256) {
        return _registerWorknet(msg.sender, params);
    }

    /// @notice Gasless worknet registration: relayer pays gas, user signs EIP-712 and pays AWP
    function registerWorknetFor(
        address user,
        WorknetParams calldata params,
        uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant whenNotPaused returns (uint256) {
        _verifyRegisterWorknetSignature(user, params, deadline, v, r, s);
        return _registerWorknet(user, params);
    }

    /// @notice Fully gasless worknet registration with EIP-2612 permit (no prior approve tx needed)
    function registerWorknetForWithPermit(
        address user,
        WorknetParams calldata params,
        uint256 deadline,
        uint8 permitV, bytes32 permitR, bytes32 permitS,
        uint8 registerV, bytes32 registerR, bytes32 registerS
    ) external nonReentrant whenNotPaused returns (uint256) {
        uint256 lpAWPAmount = initialAlphaMint * initialAlphaPrice / 1e18;
        // Wrap in try/catch: if permit was already consumed (e.g., front-run), allowance may already be set
        try IERC20Permit(awpToken).permit(user, address(this), lpAWPAmount, deadline, permitV, permitR, permitS) {} catch {}
        _verifyRegisterWorknetSignature(user, params, deadline, registerV, registerR, registerS);
        return _registerWorknet(user, params);
    }

    /// @dev Verify EIP-712 signature for registerWorknetFor
    function _verifyRegisterWorknetSignature(
        address user, WorknetParams calldata params, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) internal {
        bytes32 paramsStructHash = keccak256(abi.encode(
            WORKNET_PARAMS_TYPEHASH,
            keccak256(bytes(params.name)),
            keccak256(bytes(params.symbol)),
            params.worknetManager, params.salt, params.minStake,
            keccak256(bytes(params.skillsURI))
        ));
        _verifyDigest(user, keccak256(abi.encode(REGISTER_WORKNET_TYPEHASH, user, paramsStructHash, nonces[user]++, deadline)), deadline, v, r, s);
    }

    /// @dev Internal: shared logic for registerWorknet and registerWorknetFor
    function _registerWorknet(address user, WorknetParams calldata params) internal returns (uint256) {
        if (bytes(params.name).length == 0 || bytes(params.name).length > 64) revert InvalidWorknetName();
        if (bytes(params.symbol).length == 0 || bytes(params.symbol).length > 16) revert InvalidWorknetSymbol();
        // Reject characters that break on-chain JSON metadata in WorknetNFT.tokenURI
        _rejectJsonUnsafe(bytes(params.name));
        _rejectJsonUnsafe(bytes(params.symbol));
        if (bytes(params.skillsURI).length > 0) _rejectJsonUnsafe(bytes(params.skillsURI));
        bool autoDeployWorknet = params.worknetManager == address(0);
        if (autoDeployWorknet && defaultWorknetManagerImpl == address(0)) revert WorknetManagerRequired();

        uint256 lpAWPAmount = initialAlphaMint * initialAlphaPrice / 1e18;
        if (lpAWPAmount == 0) revert ZeroLPAmount();
        IERC20(awpToken).safeTransferFrom(user, lpManager, lpAWPAmount);

        uint256 worknetId = (block.chainid << 64) | _nextLocalId++;

        (address alphaToken, bytes32 poolId) = _deployAlphaAndLP(
            worknetId, params.name, params.symbol, lpAWPAmount, params.salt
        );

        address sc;
        if (autoDeployWorknet) {
            bytes memory initData = abi.encodeWithSignature(
                "initialize(address,address,address,bytes32,address,bytes)",
                address(this), alphaToken, awpToken, poolId, user, dexConfig
            );
            sc = address(new ERC1967Proxy(defaultWorknetManagerImpl, initData));
        } else {
            sc = params.worknetManager;
        }

        IAlphaToken(alphaToken).setWorknetMinter(sc);
        IWorknetNFT(worknetNFT).mint(user, worknetId, params.name, sc, alphaToken, params.minStake, params.skillsURI);
        worknets[worknetId] = WorknetInfo({
            lpPool: poolId,
            status: WorknetStatus.Pending,
            createdAt: uint64(block.timestamp),
            activatedAt: 0
        });

        _emitWorknetRegistered(worknetId, user, sc, alphaToken, poolId, lpAWPAmount, params);
        return worknetId;
    }

    /// @dev Deploy Alpha token + create LP (does NOT set worknet minter — caller must do that)
    function _deployAlphaAndLP(
        uint256 worknetId, string calldata name, string calldata symbol,
        uint256 lpAWPAmount, bytes32 salt
    ) internal returns (address alphaToken, bytes32 poolId) {
        alphaToken = IAlphaTokenFactory(alphaTokenFactory).deploy(worknetId, name, symbol, address(this), salt);
        IAlphaToken(alphaToken).mint(lpManager, initialAlphaMint);
        (poolId,) = ILPManager(lpManager).createPoolAndAddLiquidity(alphaToken, lpAWPAmount, initialAlphaMint);
    }

    /// @dev Emit worknet registration events
    function _emitWorknetRegistered(
        uint256 worknetId, address user, address sc, address alphaToken, bytes32 poolId, uint256 lpAWPAmount, WorknetParams calldata params
    ) internal {
        emit WorknetRegistered(
            worknetId, user,
            params.name, params.symbol,
            sc, alphaToken
        );
        emit LPCreated(worknetId, poolId, lpAWPAmount, initialAlphaMint);
    }

    // ═══════════════════════════════════════════════
    //  Worknet Lifecycle Management
    // ═══════════════════════════════════════════════

    /// @notice Activate a worknet: Pending → Active (only the NFT Owner may call)
    function activateWorknet(uint256 worknetId) external nonReentrant whenNotPaused {
        if (IWorknetNFT(worknetNFT).ownerOf(worknetId) != msg.sender) revert NotOwner();
        _activateWorknet(worknetId);
    }

    /// @notice Gasless activate worknet: relayer pays gas, NFT owner signs EIP-712
    function activateWorknetFor(
        address user, uint256 worknetId, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant whenNotPaused {
        _verifyDigest(user, keccak256(abi.encode(ACTIVATE_WORKNET_TYPEHASH, user, worknetId, nonces[user]++, deadline)), deadline, v, r, s);
        if (IWorknetNFT(worknetNFT).ownerOf(worknetId) != user) revert NotOwner();
        _activateWorknet(worknetId);
    }

    /// @dev Shared activation logic
    function _activateWorknet(uint256 worknetId) internal {
        WorknetInfo storage info = worknets[worknetId];
        if (info.status != WorknetStatus.Pending) revert InvalidWorknetStatus(worknetId, uint8(info.status));
        if (activeWorknetIds.length() >= MAX_ACTIVE_WORKNETS) revert MaxActiveWorknetsReached();
        info.status = WorknetStatus.Active;
        info.activatedAt = uint64(block.timestamp);
        activeWorknetIds.add(worknetId);
        emit WorknetActivated(worknetId);
    }

    /// @notice Pause a worknet: Active → Paused (only the NFT Owner may call)
    function pauseWorknet(uint256 worknetId) external nonReentrant whenNotPaused {
        if (IWorknetNFT(worknetNFT).ownerOf(worknetId) != msg.sender) revert NotOwner();
        WorknetInfo storage info = worknets[worknetId];
        if (info.status != WorknetStatus.Active) revert InvalidWorknetStatus(worknetId, uint8(info.status));

        info.status = WorknetStatus.Paused;
        activeWorknetIds.remove(worknetId);

        emit WorknetPaused(worknetId);
    }

    /// @notice Resume a worknet: Paused → Active (only the NFT Owner may call)
    function resumeWorknet(uint256 worknetId) external nonReentrant whenNotPaused {
        if (IWorknetNFT(worknetNFT).ownerOf(worknetId) != msg.sender) revert NotOwner();
        WorknetInfo storage info = worknets[worknetId];
        if (info.status != WorknetStatus.Paused) revert InvalidWorknetStatus(worknetId, uint8(info.status));

        if (activeWorknetIds.length() >= MAX_ACTIVE_WORKNETS) revert MaxActiveWorknetsReached();
        info.status = WorknetStatus.Active;
        activeWorknetIds.add(worknetId);

        emit WorknetResumed(worknetId);
    }

    /// @notice Ban a worknet: Active/Paused → Banned (Guardian only)
    function banWorknet(uint256 worknetId) external onlyGuardian {
        WorknetInfo storage info = worknets[worknetId];
        WorknetStatus status = info.status;
        if (status != WorknetStatus.Active && status != WorknetStatus.Paused && status != WorknetStatus.Pending)
            revert InvalidWorknetStatus(worknetId, uint8(status));

        IWorknetNFT.WorknetData memory sd = IWorknetNFT(worknetNFT).getWorknetData(worknetId);
        if (sd.worknetManager != address(0)) {
            IAlphaToken(sd.alphaToken).setMinterPaused(sd.worknetManager, true);
        }
        if (status == WorknetStatus.Active) {
            activeWorknetIds.remove(worknetId);
        }
        info.status = WorknetStatus.Banned;

        emit WorknetBanned(worknetId);
    }

    /// @notice Unban a worknet: Banned → Active (Guardian only)
    function unbanWorknet(uint256 worknetId) external onlyGuardian {
        WorknetInfo storage info = worknets[worknetId];
        if (info.status != WorknetStatus.Banned) revert InvalidWorknetStatus(worknetId, uint8(info.status));

        IWorknetNFT.WorknetData memory sd = IWorknetNFT(worknetNFT).getWorknetData(worknetId);
        if (sd.worknetManager != address(0)) {
            IAlphaToken(sd.alphaToken).setMinterPaused(sd.worknetManager, false);
        }
        if (activeWorknetIds.length() >= MAX_ACTIVE_WORKNETS) revert MaxActiveWorknetsReached();
        info.status = WorknetStatus.Active;
        activeWorknetIds.add(worknetId);

        emit WorknetUnbanned(worknetId);
    }

    /// @notice Deregister a worknet: permanently delete worknet data (Guardian only; immunity period must have elapsed)
    /// @dev Accepts Banned or Pending status. Pending worknets can be deregistered directly (never activated).
    function deregisterWorknet(uint256 worknetId) external onlyGuardian {
        WorknetInfo storage info = worknets[worknetId];
        WorknetStatus status = info.status;
        if (status != WorknetStatus.Banned && status != WorknetStatus.Pending)
            revert InvalidWorknetStatus(worknetId, uint8(status));
        uint256 immunityStart = info.activatedAt > 0 ? uint256(info.activatedAt) : uint256(info.createdAt);
        if (block.timestamp <= immunityStart + immunityPeriod) revert ImmunityNotExpired();

        // Pending worknets were never activated — no minter to pause, no activeWorknetIds entry
        // Banned worknets already had minter paused and removed from activeWorknetIds by banWorknet
        delete worknets[worknetId];
        IWorknetNFT(worknetNFT).burn(worknetId);

        emit WorknetDeregistered(worknetId);
    }

    // ═══════════════════════════════════════════════
    //  Worknet Parameters
    // ═══════════════════════════════════════════════

    /// @notice Set the initial Alpha price when registering a worknet (Guardian only)
    function setInitialAlphaPrice(uint256 price) external onlyGuardian {
        if (price < 1e12) revert PriceTooLow();
        if (price > 1e30) revert PriceTooHigh();
        initialAlphaPrice = price;
        emit InitialAlphaPriceUpdated(price);
    }

    /// @notice Update the initial Alpha token mint amount per worknet (Guardian only)
    function setInitialAlphaMint(uint256 amount) external onlyGuardian {
        if (amount == 0) revert InvalidMintAmount();
        initialAlphaMint = amount;
        emit InitialAlphaMintUpdated(amount);
    }

    /// @notice Update the guardian address (only Guardian may call — self-sovereign)
    /// @dev Guardian manages itself. If Guardian keys are lost, there is no on-chain recovery path.
    function setGuardian(address g) external onlyGuardian {
        if (g == address(0)) revert ZeroAddress();
        guardian = g;
        emit GuardianUpdated(g);
    }

    /// @notice Set the worknet deregistration immunity period (Guardian only)
    function setImmunityPeriod(uint256 p) external onlyGuardian {
        if (p < 7 days) revert ImmunityTooShort();
        immunityPeriod = p;
        emit ImmunityPeriodUpdated(p);
    }

    /// @notice Replace the AlphaToken factory (Guardian only)
    function setAlphaTokenFactory(address factory) external onlyGuardian {
        if (factory == address(0)) revert ZeroAddress();
        alphaTokenFactory = factory;
        emit AlphaTokenFactoryUpdated(factory);
    }

    /// @notice Set the default worknet implementation (Guardian only)
    function setWorknetManagerImpl(address impl) external onlyGuardian {
        if (impl == address(0)) revert ZeroAddress();
        defaultWorknetManagerImpl = impl;
        emit DefaultWorknetManagerImplUpdated(impl);
    }

    /// @notice Update DEX configuration for future auto-deployed WorknetManagers (Guardian only)
    function setDexConfig(bytes calldata dexConfig_) external onlyGuardian {
        dexConfig = dexConfig_;
        emit DexConfigUpdated();
    }

    /// @notice Set the base URI for WorknetNFT metadata (Guardian only)
    function setWorknetBaseURI(string calldata baseURI) external onlyGuardian {
        IWorknetNFT(worknetNFT).setBaseURI(baseURI);
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
        /// @dev Agent's stake amount on the specified worknet (uses root as staker)
        uint256 stake;
        /// @dev Reward recipient address for the root
        address rewardRecipient;
    }

    /// @notice Query complete information for a single agent on a specified worknet
    /// @param agent Agent address
    /// @param worknetId Worknet ID
    /// @return AgentInfo containing root, isValid, stake, rewardRecipient
    function getAgentInfo(address agent, uint256 worknetId) external view returns (AgentInfo memory) {
        // Walk the binding tree to find the root
        address root = _resolveRoot(agent);
        bool isValid = boundTo[agent] != address(0) || recipient[agent] != address(0);
        uint256 stake = IStakingVault(stakingVault).getAgentStake(root, agent, worknetId);
        address recip = recipient[root] != address(0) ? recipient[root] : root;
        return AgentInfo(root, isValid, stake, recip);
    }

    /// @notice Batch query information for multiple agents on a specified worknet
    function getAgentsInfo(address[] calldata agents, uint256 worknetId)
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
                ? IStakingVault(stakingVault).getAgentStake(root, agent, worknetId)
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

    /// @dev Reject strings containing " or \ (breaks on-chain JSON in WorknetNFT.tokenURI)
    function _rejectJsonUnsafe(bytes memory b) internal pure {
        for (uint256 i = 0; i < b.length;) {
            bytes1 c = b[i];
            // Reject: " (0x22), \ (0x5c), and control characters (0x00-0x1F) per RFC 8259
            if (c == 0x22 || c == 0x5c || c < 0x20) revert JsonUnsafeCharacter();
            unchecked { ++i; }
        }
    }

    /// @dev Anti-cycle check: walk up from target, revert if sender is found in the chain
    function _checkCycle(address sender, address target) internal view {
        address cur = target;
        uint256 depth;
        while (boundTo[cur] != address(0) && boundTo[cur] != cur) {
            if (boundTo[cur] == sender) revert CycleDetected();
            cur = boundTo[cur];
            if (++depth >= 256) revert ChainTooDeep();
        }
    }

    /// @dev Walk the binding chain to find the root
    function _resolveRoot(address addr) internal view returns (address) {
        address cur = addr;
        uint256 depth;
        while (boundTo[cur] != address(0) && boundTo[cur] != cur) {
            cur = boundTo[cur];
            if (++depth >= 256) revert ChainTooDeep();
        }
        return cur;
    }

    // ═══════════════════════════════════════════════
    //  View — general view functions
    // ═══════════════════════════════════════════════

    /// @notice Get worknet lifecycle state
    function getWorknet(uint256 worknetId) external view returns (WorknetInfo memory) {
        return worknets[worknetId];
    }

    /// @notice Get complete worknet info combining AWPRegistry state + WorknetNFT identity
    function getWorknetFull(uint256 worknetId) external view returns (WorknetFullInfo memory) {
        WorknetInfo storage info = worknets[worknetId];
        IWorknetNFT.WorknetData memory nftData = IWorknetNFT(worknetNFT).getWorknetData(worknetId);

        return WorknetFullInfo({
            worknetManager: nftData.worknetManager,
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

    /// @notice Get the current number of active worknets
    function getActiveWorknetCount() external view returns (uint256) {
        return activeWorknetIds.length();
    }

    /// @notice Get the active worknet ID at a given index
    function getActiveWorknetIdAt(uint256 index) external view returns (uint256) {
        return activeWorknetIds.at(index);
    }

    /// @notice Check whether a specified worknet is in the Active state
    function isWorknetActive(uint256 worknetId) external view returns (bool) {
        return worknets[worknetId].status == WorknetStatus.Active;
    }

    /// @notice Get the next worknet ID to be assigned (globally unique: chainId << 64 | localCounter)
    function nextWorknetId() external view returns (uint256) {
        return (block.chainid << 64) | _nextLocalId;
    }

    /// @notice Extract chainId from a global worknetId
    function extractChainId(uint256 worknetId) external pure returns (uint256) {
        return worknetId >> 64;
    }

    /// @notice Extract local counter from a global worknetId
    function extractLocalId(uint256 worknetId) external pure returns (uint256) {
        return worknetId & ((1 << 64) - 1);
    }

    // ═══════════════════════════════════════════════
    //  Pause — emergency pause
    // ═══════════════════════════════════════════════

    /// @notice Emergency pause the contract (only Guardian may call)
    function pause() external onlyGuardian {
        _pause();
    }

    /// @notice Unpause the contract (Guardian only)
    function unpause() external onlyGuardian {
        _unpause();
    }

    /// @notice Update the LP manager address (Guardian only, for DEX migration)
    function setLPManager(address lpManager_) external onlyGuardian {
        if (lpManager_ == address(0)) revert ZeroAddress();
        lpManager = lpManager_;
        emit LPManagerUpdated(lpManager_);
    }

    /// @notice Get active worknet IDs with pagination
    /// @param offset Starting index
    /// @param limit Maximum number of IDs to return
    function getActiveWorknetIds(uint256 offset, uint256 limit) external view returns (uint256[] memory) {
        uint256 total = activeWorknetIds.length();
        if (offset >= total) return new uint256[](0);
        uint256 end = offset + limit;
        if (end > total) end = total;
        uint256[] memory ids = new uint256[](end - offset);
        for (uint256 i = offset; i < end;) {
            ids[i - offset] = activeWorknetIds.at(i);
            unchecked { ++i; }
        }
        return ids;
    }
}
