// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {ReentrancyGuardTransient} from "@openzeppelin/contracts/utils/ReentrancyGuardTransient.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {EIP712Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

import {IWorknetToken} from "./interfaces/IWorknetToken.sol";
import {IWorknetTokenFactory} from "./interfaces/IWorknetTokenFactory.sol";
import {IAWPAllocator} from "./interfaces/IAWPAllocator.sol";
import {IAWPWorkNet} from "./interfaces/IAWPWorkNet.sol";
import {ILPManager} from "./interfaces/ILPManager.sol";
import {IAWPRegistry} from "./interfaces/IAWPRegistry.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

/// @title AWPRegistry — Unified entry point for the AWP protocol (worknet management + staking management)
/// @author AWP Team
/// @notice Account System V2: tree-based binding, optional registration, explicit staker parameter.
/// @dev Inheritance: UUPSUpgradeable (UUPS proxy pattern), PausableUpgradeable (emergency pause),
///      ReentrancyGuardTransient (reentrancy protection), EIP712Upgradeable (EIP-712 signing domain, domain name "AWPRegistry" v1).
///      8 immutable addresses baked into impl bytecode. Storage-mutable: guardian, defaultWorknetManagerImpl, initialAlphaPrice/Mint.
contract AWPRegistry is UUPSUpgradeable, PausableUpgradeable, ReentrancyGuardTransient, EIP712Upgradeable, IAWPRegistry {
    using SafeERC20 for IERC20;
    using EnumerableSet for EnumerableSet.UintSet;

    // ══════════════════════════════════════════════
    //  Immutables (baked into impl bytecode — zero SLOAD cost)
    //  All proxy addresses are immutable: proxy address never changes, impl changes via UUPS upgrade.
    // ══════════════════════════════════════════════

    /// @notice AWP token address (same on all chains)
    address public immutable awpToken;
    /// @notice AWPWorkNet proxy address (UUPS upgradeable)
    address public immutable awpWorkNet;
    /// @notice WorknetToken factory proxy address (UUPS upgradeable)
    address public immutable worknetTokenFactory;
    /// @notice AWP emission proxy address (same on all chains)
    address public immutable awpEmission;
    /// @notice LP manager proxy address (UUPS upgradeable)
    address public immutable lpManager;
    /// @notice AWPAllocator proxy address (same on all chains)
    address public immutable awpAllocator;
    /// @notice veAWP proxy address (UUPS upgradeable)
    address public immutable veAWP;
    /// @notice Treasury (Timelock) address (same on all chains)
    address public immutable treasury;

    // ══════════════════════════════════════════════
    //  Storage — (slots 0-11 must match deployed proxy layout)
    // ══════════════════════════════════════════════

    /// @dev Freed slots 0-7: were module addresses, now all immutable
    uint256 private __freed_slot0;
    uint256 private __freed_slot1;
    uint256 private __freed_slot2;
    uint256 private __freed_slot3;
    uint256 private __freed_slot4;
    uint256 private __freed_slot5;
    uint256 private __freed_slot6;
    uint256 private __freed_slot7;
    /// @notice Guardian address — emergency pause + parameter management
    address public guardian;
    /// @notice Default WorknetManager implementation (for auto-deploying worknet proxies)
    address public defaultWorknetManagerImpl;
    /// @dev Freed slot 10: was dexConfig
    bytes private __freed_slot10;

    /// @dev Freed slot 11: was _deployer (packed with registryInitialized)
    uint256 private __freed_slot11;

    // ══════════════════════════════════════════════
    //  Account System V2 — binding, recipient, delegation
    // ══════════════════════════════════════════════

    /// @notice Approximate registration count (incremented on first bind/setRecipient, never decremented).
    /// @dev Not deduplicated: unbind→rebind increments again. Use off-chain indexer for exact unique count.
    ///      Intentional trade-off: O(1) on-chain vs exact count would require an additional mapping (21k gas per new user).
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
    /// Global worknetId = block.chainid * CHAIN_ID_MULTIPLIER + _nextLocalId
    uint256 private _nextLocalId;

    /// @dev Active worknet ID set
    EnumerableSet.UintSet private activeWorknetIds;

    /// @notice Maximum number of active worknets
    uint128 public constant MAX_ACTIVE_WORKNETS = 10000;

    /// @notice Multiplier for worknetId generation: worknetId = chainId * CHAIN_ID_MULTIPLIER + localId
    uint256 public constant CHAIN_ID_MULTIPLIER = 100_000_000;

    /// @notice Initial price of the Alpha token when registering a worknet (denominated in AWP wei, default 1e15 = 0.001 AWP)
    uint256 public initialAlphaPrice;

    /// @notice Alpha token mint amount per worknet registration (governance-settable, default 1B = 1e27 wei)
    uint256 public initialAlphaMint;

    /// @dev Deprecated: was immunity period. Kept for UUPS storage layout compatibility.
    uint256 private __deprecated_immunityPeriod;

    // ══════════════════════════════════════════════
    //  Gasless — EIP-712 signature related
    // ══════════════════════════════════════════════

    /// @notice Per-signer nonce for replay attack prevention (nonces[signer]++ when validated)
    mapping(address => uint256) public nonces;

    // ══════════════════════════════════════════════
    //  Worknet escrow (AWP held until activation)
    // ══════════════════════════════════════════════

    /// @notice Escrowed AWP and Alpha mint amount for pending worknets (locked at registration time)
    struct EscrowInfo {
        uint128 lpAWPAmount;   // AWP escrowed (in token units, not wei — fits uint128)
        uint128 alphaMint;     // Alpha mint amount snapshot (in token units)
    }
    /// @notice worknetId => escrow info (only exists for Pending worknets)
    mapping(uint256 => EscrowInfo) public worknetEscrow;

    /// @notice Pending worknet registration params (stored until activation or cancellation)
    struct PendingParams {
        string name;
        string symbol;
        address worknetManager;    // address(0) = auto-deploy
        bytes32 salt;              // CREATE2 salt for Alpha token
        uint128 minStake;
        string skillsURI;
        address owner;
    }
    /// @notice worknetId => pending registration params (deleted on activate/cancel/reject)
    mapping(uint256 => PendingParams) public pendingWorknets;

    /// @dev Reserved storage gap for future upgrades (UUPS pattern)
    uint256[47] private __gap;

    /// @dev EIP-712 type hash: Bind(address agent, address target, uint256 nonce, uint256 deadline)
    bytes32 private constant BIND_TYPEHASH =
        keccak256("Bind(address agent,address target,uint256 nonce,uint256 deadline)");

    /// @dev EIP-712 type hash: SetRecipient(address user, address recipient, uint256 nonce, uint256 deadline)
    bytes32 private constant SET_RECIPIENT_TYPEHASH =
        keccak256("SetRecipient(address user,address recipient,uint256 nonce,uint256 deadline)");

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

    /// @dev Pre-computed selector for WorknetManager.initialize(address,bytes32,address)
    bytes4 private constant WM_INIT_SELECTOR = bytes4(keccak256("initialize(address,bytes32,address)"));

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
    /// @dev Worknet contract address cannot be the zero address
    error WorknetManagerRequired();
    /// @dev Worknet status does not meet the precondition
    error InvalidWorknetStatus(uint256 worknetId, uint8 currentStatus);
    /// @dev Escrow amount overflows uint128
    error AmountOverflow();
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

    // ══════════════════════════════════════════════
    //  Constructor
    // ══════════════════════════════════════════════

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor(
        address awpToken_,
        address awpWorkNet_,
        address worknetTokenFactory_,
        address awpEmission_,
        address lpManager_,
        address awpAllocator_,
        address veAWP_,
        address treasury_
    ) {
        awpToken = awpToken_;
        awpWorkNet = awpWorkNet_;
        worknetTokenFactory = worknetTokenFactory_;
        awpEmission = awpEmission_;
        lpManager = lpManager_;
        awpAllocator = awpAllocator_;
        veAWP = veAWP_;
        treasury = treasury_;
        _disableInitializers();
    }

    /// @notice Initialize the registry (called once via proxy)
    /// @param deployer_ Unused (kept for ABI compatibility with existing proxy)
    /// @param treasury_ Unused (treasury is now immutable from constructor)
    /// @param guardian_ Guardian address
    function initialize(
        address deployer_,
        address treasury_,
        address guardian_
    ) external initializer {
        __Pausable_init();
        // ReentrancyGuardTransient — no init needed (uses TSTORE)
        __EIP712_init("AWPRegistry", "1");

        if (guardian_ == address(0)) revert ZeroAddress();
        guardian = guardian_;

        _nextLocalId = 1;
        initialAlphaPrice = 1e15;
        initialAlphaMint = 1_000_000_000 * 1e18;
    }

    /// @dev UUPS upgrade authorization — only Guardian may upgrade
    function _authorizeUpgrade(address) internal override onlyGuardian {}

    /// @notice Retrieve all module addresses in a single call
    /// @return In order: awpToken, awpWorkNet, worknetTokenFactory, awpEmission,
    ///         lpManager, awpAllocator, veAWP, treasury, guardian
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
            awpWorkNet,
            worknetTokenFactory,
            awpEmission,
            lpManager,
            awpAllocator,
            veAWP,
            treasury,
            guardian
        );
    }

    // ═══════════════════════════════════════════════
    //  Account V2: Binding (tree structure)
    // ═══════════════════════════════════════════════

    /// @notice Bind msg.sender to target
    function bind(address target) external nonReentrant whenNotPaused {
        _bind(msg.sender, target);
    }

    /// @notice Gasless bind: relayer pays gas, agent signs EIP-712
    function bindFor(address agent, address target, uint256 deadline, uint8 v, bytes32 r, bytes32 s)
        external nonReentrant whenNotPaused
    {
        _verifyDigest(agent, keccak256(abi.encode(BIND_TYPEHASH, agent, target, nonces[agent]++, deadline)), deadline, v, r, s);
        _bind(agent, target);
    }

    /// @notice Unbind msg.sender from the tree
    function unbind() external nonReentrant whenNotPaused {
        _unbind(msg.sender);
    }

    /// @notice Gasless unbind
    function unbindFor(address user, uint256 deadline, uint8 v, bytes32 r, bytes32 s)
        external nonReentrant whenNotPaused
    {
        _verifyDigest(user, keccak256(abi.encode(UNBIND_TYPEHASH, user, nonces[user]++, deadline)), deadline, v, r, s);
        _unbind(user);
    }

    function _bind(address agent, address target) internal {
        if (target == address(0)) revert ZeroAddress();
        if (target == agent) revert SelfBind();
        _checkCycle(agent, target);
        if (boundTo[agent] == address(0) && recipient[agent] == address(0)) {
            registeredCount++;
            emit UserRegistered(agent);
        }
        boundTo[agent] = target;
        emit Bound(agent, target);
    }

    function _unbind(address user) internal {
        boundTo[user] = address(0);
        emit Unbound(user);
    }

    // ═══════════════════════════════════════════════
    //  Account V2: Recipient
    // ═══════════════════════════════════════════════

    /// @notice Set reward recipient. Pass address(0) to clear (reverts to self as default).
    function setRecipient(address addr) external nonReentrant whenNotPaused {
        _setRecipient(msg.sender, addr);
    }

    /// @notice Gasless set recipient. Pass address(0) to clear.
    function setRecipientFor(
        address user, address _recipient, uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant whenNotPaused {
        _verifyDigest(user, keccak256(abi.encode(SET_RECIPIENT_TYPEHASH, user, _recipient, nonces[user]++, deadline)), deadline, v, r, s);
        _setRecipient(user, _recipient);
    }

    function _setRecipient(address user, address addr) internal {
        if (addr == address(0)) {
            recipient[user] = address(0);
            emit RecipientSet(user, address(0));
            return;
        }
        if (recipient[user] == address(0) && boundTo[user] == address(0)) {
            registeredCount++;
            emit UserRegistered(user);
        }
        recipient[user] = addr;
        emit RecipientSet(user, addr);
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

    /// @notice Grant delegate authorization
    function grantDelegate(address delegate) external whenNotPaused {
        _grantDelegate(msg.sender, delegate);
    }

    /// @notice Gasless grant delegate
    function grantDelegateFor(address user, address delegate, uint256 deadline, uint8 v, bytes32 r, bytes32 s)
        external whenNotPaused
    {
        _verifyDigest(user, keccak256(abi.encode(GRANT_DELEGATE_TYPEHASH, user, delegate, nonces[user]++, deadline)), deadline, v, r, s);
        _grantDelegate(user, delegate);
    }

    /// @notice Revoke delegate authorization
    function revokeDelegate(address delegate) external whenNotPaused {
        _revokeDelegate(msg.sender, delegate);
    }

    /// @notice Gasless revoke delegate
    function revokeDelegateFor(address user, address delegate, uint256 deadline, uint8 v, bytes32 r, bytes32 s)
        external whenNotPaused
    {
        _verifyDigest(user, keccak256(abi.encode(REVOKE_DELEGATE_TYPEHASH, user, delegate, nonces[user]++, deadline)), deadline, v, r, s);
        _revokeDelegate(user, delegate);
    }

    function _grantDelegate(address user, address delegate) internal {
        delegates[user][delegate] = true;
        emit DelegateGranted(user, delegate);
    }

    function _revokeDelegate(address user, address delegate) internal {
        if (delegate == user) revert CannotRevokeSelf();
        delegates[user][delegate] = false;
        emit DelegateRevoked(user, delegate);
    }

    // ═══════════════════════════════════════════════
    //  Worknet Registration
    // ═══════════════════════════════════════════════

    /// @notice Register a new worknet: escrow AWP and store params (Alpha deploy + NFT mint deferred to activateWorknet)
    /// @param params Worknet parameters (name, symbol, worknetManager, salt, minStake, skillsURI)
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

    /// @dev Internal: shared logic for registerWorknet and registerWorknetFor.
    ///      Escrows AWP and stores registration params. No Alpha deploy, no NFT mint — deferred to activateWorknet.
    function _registerWorknet(address user, WorknetParams calldata params) internal returns (uint256) {
        if (bytes(params.name).length == 0 || bytes(params.name).length > 64) revert InvalidWorknetName();
        if (bytes(params.symbol).length == 0 || bytes(params.symbol).length > 16) revert InvalidWorknetSymbol();
        _rejectJsonUnsafe(bytes(params.name));
        _rejectJsonUnsafe(bytes(params.symbol));
        if (bytes(params.skillsURI).length > 0) _rejectJsonUnsafe(bytes(params.skillsURI));
        if (params.worknetManager == address(0) && defaultWorknetManagerImpl == address(0)) revert WorknetManagerRequired();

        // Snapshot current price/mint, validate before external calls
        uint256 alphaMint = initialAlphaMint;
        uint256 lpAWPAmount = alphaMint * initialAlphaPrice / 1e18;
        if (lpAWPAmount == 0) revert ZeroLPAmount();
        if (lpAWPAmount > type(uint128).max || alphaMint > type(uint128).max) revert AmountOverflow();

        // Escrow AWP
        IERC20(awpToken).safeTransferFrom(user, address(this), lpAWPAmount);

        uint256 worknetId = block.chainid * CHAIN_ID_MULTIPLIER + _nextLocalId++;

        // Store state — no Alpha deploy, no NFT mint (deferred to activateWorknet)
        worknets[worknetId] = WorknetInfo({
            lpPool: bytes32(0),
            status: WorknetStatus.Pending,
            createdAt: uint64(block.timestamp),
            activatedAt: 0
        });
        worknetEscrow[worknetId] = EscrowInfo(uint128(lpAWPAmount), uint128(alphaMint));
        pendingWorknets[worknetId] = PendingParams({
            name: params.name,
            symbol: params.symbol,
            worknetManager: params.worknetManager,
            salt: params.salt,
            minStake: params.minStake,
            skillsURI: params.skillsURI,
            owner: user
        });

        emit WorknetRegistered(worknetId, user, params.name, params.symbol);
        return worknetId;
    }

    // ═══════════════════════════════════════════════
    //  Worknet Lifecycle Management
    // ═══════════════════════════════════════════════

    /// @notice Activate a worknet: Pending → Active (Guardian approval required)
    function activateWorknet(uint256 worknetId) external nonReentrant whenNotPaused onlyGuardian {
        _activateWorknet(worknetId);
    }

    /// @dev Shared activation logic: deploy Alpha → create LP → deploy WorknetManager → mint NFT
    function _activateWorknet(uint256 worknetId) internal {
        WorknetInfo storage info = worknets[worknetId];
        if (info.status != WorknetStatus.Pending) revert InvalidWorknetStatus(worknetId, uint8(info.status));
        if (activeWorknetIds.length() >= MAX_ACTIVE_WORKNETS) revert MaxActiveWorknetsReached();

        EscrowInfo memory escrow = worknetEscrow[worknetId];
        PendingParams memory pp = pendingWorknets[worknetId];

        // 1. Deploy worknet token
        address worknetToken = IWorknetTokenFactory(worknetTokenFactory).deploy(
            worknetId, pp.name, pp.symbol, pp.salt
        );

        // 2. Transfer escrowed AWP to LPManager, mint worknet tokens, create LP pool
        IERC20(awpToken).safeTransfer(lpManager, escrow.lpAWPAmount);
        IWorknetToken(worknetToken).mint(lpManager, escrow.alphaMint);
        (bytes32 poolId,) = ILPManager(lpManager).createPoolAndAddLiquidity(worknetToken, escrow.lpAWPAmount, escrow.alphaMint);

        // 3. Deploy WorknetManager proxy if needed (CREATE2 with worknetId as salt)
        address sc = pp.worknetManager;
        if (sc == address(0)) {
            bytes memory initData = abi.encodeWithSelector(WM_INIT_SELECTOR, worknetToken, poolId, pp.owner);
            sc = address(new ERC1967Proxy{salt: bytes32(worknetId)}(defaultWorknetManagerImpl, initData));
        }

        // 4. Lock worknet token minter to WorknetManager
        IWorknetToken(worknetToken).setMinter(sc);

        // 5. Mint NFT with complete identity (all fields known at this point)
        IAWPWorkNet(awpWorkNet).mint(
            pp.owner, worknetId,
            pp.name, pp.symbol,
            sc, worknetToken, poolId,
            pp.minStake, pp.skillsURI
        );

        // 6. Update state
        info.lpPool = poolId;
        info.status = WorknetStatus.Active;
        info.activatedAt = uint64(block.timestamp);
        activeWorknetIds.add(worknetId);
        delete worknetEscrow[worknetId];
        delete pendingWorknets[worknetId];

        emit LPCreated(worknetId, poolId, escrow.lpAWPAmount, escrow.alphaMint);
        emit WorknetActivated(worknetId);
    }

    /// @notice Pause a worknet: Active → Paused (only the NFT Owner may call)
    function pauseWorknet(uint256 worknetId) external nonReentrant whenNotPaused {
        if (IAWPWorkNet(awpWorkNet).ownerOf(worknetId) != msg.sender) revert NotOwner();
        WorknetInfo storage info = worknets[worknetId];
        if (info.status != WorknetStatus.Active) revert InvalidWorknetStatus(worknetId, uint8(info.status));

        info.status = WorknetStatus.Paused;
        activeWorknetIds.remove(worknetId);

        emit WorknetPaused(worknetId);
    }

    /// @notice Resume a worknet: Paused → Active (only the NFT Owner may call)
    function resumeWorknet(uint256 worknetId) external nonReentrant whenNotPaused {
        if (IAWPWorkNet(awpWorkNet).ownerOf(worknetId) != msg.sender) revert NotOwner();
        WorknetInfo storage info = worknets[worknetId];
        if (info.status != WorknetStatus.Paused) revert InvalidWorknetStatus(worknetId, uint8(info.status));

        if (activeWorknetIds.length() >= MAX_ACTIVE_WORKNETS) revert MaxActiveWorknetsReached();
        info.status = WorknetStatus.Active;
        activeWorknetIds.add(worknetId);

        emit WorknetResumed(worknetId);
    }

    /// @notice Cancel a pending worknet: Pending → deleted (original registrant only, full AWP refund)
    function cancelWorknet(uint256 worknetId) external nonReentrant whenNotPaused {
        if (pendingWorknets[worknetId].owner != msg.sender) revert NotOwner();
        _cancelPending(worknetId, msg.sender);
        emit WorknetCancelled(worknetId);
    }

    /// @notice Reject a pending worknet: Pending → deleted (Guardian only, full AWP refund to registrant)
    function rejectWorknet(uint256 worknetId) external nonReentrant onlyGuardian {
        address owner = pendingWorknets[worknetId].owner;
        if (owner == address(0)) revert InvalidWorknetStatus(worknetId, 0);
        _cancelPending(worknetId, owner);
        emit WorknetRejected(worknetId);
    }

    /// @dev Shared cancel/reject logic: refund escrowed AWP, delete state (no NFT to burn — NFT not yet minted)
    function _cancelPending(uint256 worknetId, address refundTo) internal {
        WorknetInfo storage info = worknets[worknetId];
        if (info.status != WorknetStatus.Pending) revert InvalidWorknetStatus(worknetId, uint8(info.status));

        uint256 refund = worknetEscrow[worknetId].lpAWPAmount;
        delete worknetEscrow[worknetId];
        delete pendingWorknets[worknetId];
        delete worknets[worknetId];

        if (refund > 0) {
            IERC20(awpToken).safeTransfer(refundTo, refund);
        }
    }

    /// @notice Ban a worknet: Active/Paused → Banned (Guardian only)
    function banWorknet(uint256 worknetId) external nonReentrant onlyGuardian {
        WorknetInfo storage info = worknets[worknetId];
        WorknetStatus status = info.status;
        if (status != WorknetStatus.Active && status != WorknetStatus.Paused)
            revert InvalidWorknetStatus(worknetId, uint8(status));

        if (status == WorknetStatus.Active) {
            activeWorknetIds.remove(worknetId);
        }
        info.status = WorknetStatus.Banned;

        emit WorknetBanned(worknetId);
    }

    /// @notice Unban a worknet: Banned → Active (Guardian only)
    function unbanWorknet(uint256 worknetId) external nonReentrant onlyGuardian {
        WorknetInfo storage info = worknets[worknetId];
        if (info.status != WorknetStatus.Banned) revert InvalidWorknetStatus(worknetId, uint8(info.status));

        if (activeWorknetIds.length() >= MAX_ACTIVE_WORKNETS) revert MaxActiveWorknetsReached();
        info.status = WorknetStatus.Active;
        activeWorknetIds.add(worknetId);

        emit WorknetUnbanned(worknetId);
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

    /// @notice Set the default worknet implementation (Guardian only)
    function setWorknetManagerImpl(address impl) external onlyGuardian {
        if (impl == address(0)) revert ZeroAddress();
        defaultWorknetManagerImpl = impl;
        emit DefaultWorknetManagerImplUpdated(impl);
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
        uint256 stake = IAWPAllocator(awpAllocator).getAgentStake(root, agent, worknetId);
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
                ? IAWPAllocator(awpAllocator).getAgentStake(root, agent, worknetId)
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

    /// @dev Reject strings containing " or \ (breaks on-chain JSON in AWPWorkNet.tokenURI)
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

    /// @notice Get complete worknet info. For Pending worknets returns PendingParams; for Active+ returns NFT data.
    function getWorknetFull(uint256 worknetId) external view returns (WorknetFullInfo memory) {
        WorknetInfo storage info = worknets[worknetId];
        if (info.status == WorknetStatus.None) revert InvalidWorknetStatus(worknetId, 0);

        if (info.status == WorknetStatus.Pending) {
            // NFT not yet minted — return data from PendingParams
            PendingParams storage pp = pendingWorknets[worknetId];
            return WorknetFullInfo({
                worknetManager: pp.worknetManager,
                worknetToken: address(0),
                lpPool: bytes32(0),
                status: info.status,
                createdAt: info.createdAt,
                activatedAt: 0,
                name: pp.name,
                symbol: pp.symbol,
                skillsURI: pp.skillsURI,
                minStake: pp.minStake,
                owner: pp.owner
            });
        }

        // Active/Paused/Banned — read from NFT
        IAWPWorkNet.WorknetData memory nftData = IAWPWorkNet(awpWorkNet).getWorknetData(worknetId);
        return WorknetFullInfo({
            worknetManager: nftData.worknetManager,
            worknetToken: nftData.worknetToken,
            lpPool: info.lpPool,
            status: info.status,
            createdAt: info.createdAt,
            activatedAt: info.activatedAt,
            name: nftData.name,
            symbol: nftData.symbol,
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

    /// @notice Get the next worknet ID to be assigned (globally unique: chainId * 1e8 + localCounter)
    function nextWorknetId() external view returns (uint256) {
        return block.chainid * CHAIN_ID_MULTIPLIER + _nextLocalId;
    }

    /// @notice Extract chainId from a global worknetId
    function extractChainId(uint256 worknetId) external pure returns (uint256) {
        return worknetId / CHAIN_ID_MULTIPLIER;
    }

    /// @notice Extract local counter from a global worknetId
    function extractLocalId(uint256 worknetId) external pure returns (uint256) {
        return worknetId % CHAIN_ID_MULTIPLIER;
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

    // ═══════════════════════════════════════════════
    //  Token rescue
    // ═══════════════════════════════════════════════

    error CannotRescueEscrowedToken();

    /// @notice Rescue accidentally sent ERC20 tokens (Guardian only). Cannot rescue escrowed AWP.
    function rescueToken(address token, address to, uint256 amount) external onlyGuardian {
        if (token == awpToken) revert CannotRescueEscrowedToken();
        IERC20(token).safeTransfer(to, amount);
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
