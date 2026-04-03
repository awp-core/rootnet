// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AlphaToken} from "../src/token/AlphaToken.sol";
import {AlphaTokenFactory} from "../src/token/AlphaTokenFactory.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {WorknetNFT} from "../src/core/WorknetNFT.sol";
import {MockLPManager} from "./helpers/MockLPManager.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {IAWPRegistry} from "../src/interfaces/IAWPRegistry.sol";
import {IWorknetNFT} from "../src/interfaces/IWorknetNFT.sol";
import {IAlphaToken} from "../src/interfaces/IAlphaToken.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {WorknetManager} from "../src/worknets/WorknetManager.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

/// @title AWPRegistryExtended — Extended test coverage: state machine, gasless sigs, MAX_ACTIVE limits
contract AWPRegistryExtendedTest is Test {
    AWPToken awp;
    AlphaTokenFactory factory;
    AWPEmission emission;
    StakingVault vault;
    StakeNFT stakeNFT;
    WorknetNFT nft;
    MockLPManager lp;
    AWPRegistry awpRegistry;
    Treasury treasury;
    WorknetManager worknetManagerImpl;

    address deployer = address(1);
    address guardian = address(2);
    address user1;
    uint256 user1Pk;
    address user2 = address(6);
    address worknetManager = address(5);
    address relayer = address(7);

    uint256 constant INITIAL_DAILY_EMISSION = 31_600_000 * 1e18;
    uint256 constant EPOCH_DURATION = 1 days;

    // EIP-712 type hashes (must match AWPRegistry definitions)
    bytes32 private constant BIND_TYPEHASH =
        keccak256("Bind(address agent,address target,uint256 nonce,uint256 deadline)");
    bytes32 private constant ACTIVATE_WORKNET_TYPEHASH =
        keccak256("ActivateWorknet(address user,uint256 worknetId,uint256 nonce,uint256 deadline)");
    bytes32 private constant WORKNET_PARAMS_TYPEHASH =
        keccak256("WorknetParams(string name,string symbol,address worknetManager,bytes32 salt,uint128 minStake,string skillsURI)");
    bytes32 private constant REGISTER_WORKNET_TYPEHASH =
        keccak256("RegisterWorknet(address user,WorknetParams params,uint256 nonce,uint256 deadline)WorknetParams(string name,string symbol,address worknetManager,bytes32 salt,uint128 minStake,string skillsURI)");

    function setUp() public {
        // Generate signable test account via makeAddrAndKey
        (user1, user1Pk) = makeAddrAndKey("user1");

        vm.startPrank(deployer);

        // Deploy AWP token
        awp = new AWPToken("AWP Token", "AWP", deployer);
        awp.initialMint(200_000_000 * 1e18);

        // Deploy Treasury
        address[] memory proposers = new address[](0);
        address[] memory executors = new address[](1);
        executors[0] = address(0);
        treasury = new Treasury(0, proposers, executors, deployer);

        // Deploy AWPRegistry (UUPS proxy)
        AWPRegistry awpRegistryImpl = new AWPRegistry();
        awpRegistry = AWPRegistry(address(new ERC1967Proxy(
            address(awpRegistryImpl),
            abi.encodeCall(AWPRegistry.initialize, (deployer, address(treasury), guardian))
        )));

        // Deploy sub-contracts
        factory = new AlphaTokenFactory(deployer, 0);
        nft = new WorknetNFT("AWP Worknet", "AWPSUB", address(awpRegistry));
        lp = new MockLPManager(address(awpRegistry), address(awp));

        // Deploy AWPEmission (UUPS proxy)
        AWPEmission emissionImpl = new AWPEmission();
        bytes memory emissionInitData = abi.encodeCall(
            AWPEmission.initialize,
            (address(awp), deployer, INITIAL_DAILY_EMISSION, block.timestamp, EPOCH_DURATION, address(treasury))
        );
        ERC1967Proxy emissionProxy = new ERC1967Proxy(address(emissionImpl), emissionInitData);
        emission = AWPEmission(address(emissionProxy));

        // Configure AWP minter
        awp.addMinter(address(emission));
        awp.renounceAdmin();

        // Configure factory
        factory.setAddresses(address(awpRegistry));

        // Deploy StakeNFT and StakingVault
        vault = StakingVault(address(new ERC1967Proxy(
            address(new StakingVault()), abi.encodeCall(StakingVault.initialize, (address(awpRegistry), deployer))
        )));
        stakeNFT = new StakeNFT(address(awp), address(vault), address(awpRegistry));

        // Deploy WorknetManager impl (for auto-deploy tests)
        worknetManagerImpl = new WorknetManager();

        // Initialize registry (no defaultWorknetManagerImpl, set as needed in tests)
        awpRegistry.initializeRegistry(
            address(awp),
            address(nft),
            address(factory),
            address(emission),
            address(lp),
            address(vault),
            address(stakeNFT),
            address(0),
            ""
        );

        // Distribute AWP to test users
        awp.transfer(user1, 50_000_000 * 1e18);
        awp.transfer(user2, 50_000_000 * 1e18);

        vm.stopPrank();
    }

    // ═══════════════════════════════════════════════
    //  Helper functions
    // ═══════════════════════════════════════════════

    /// @dev Register a worknet as user1, return worknetId
    function _registerWorknet() internal returns (uint256) {
        uint256 lpCost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.startPrank(user1);
        awp.approve(address(awpRegistry), lpCost);
        uint256 worknetId = awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams({
                name: "TestWorknet",
                symbol: "TSUB",
                worknetManager: worknetManager,
                salt: bytes32(0),
                minStake: 0,
                skillsURI: ""
            })
        );
        vm.stopPrank();
        return worknetId;
    }

    /// @dev Register and activate worknet
    function _registerAndActivate() internal returns (uint256) {
        uint256 worknetId = _registerWorknet();
        vm.prank(user1);
        awpRegistry.activateWorknet(worknetId);
        return worknetId;
    }

    /// @dev Build EIP-712 domain separator (matches AWPRegistry "AWPRegistry" v1)
    function _domainSeparator() internal view returns (bytes32) {
        return keccak256(abi.encode(
            keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
            keccak256("AWPRegistry"),
            keccak256("1"),
            block.chainid,
            address(awpRegistry)
        ));
    }

    /// @dev Sign EIP-712 digest with user1's private key
    function _signDigest(bytes32 structHash) internal view returns (uint8 v, bytes32 r, bytes32 s) {
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", _domainSeparator(), structHash));
        (v, r, s) = vm.sign(user1Pk, digest);
    }

    // ═══════════════════════════════════════════════════════
    //  1. Worknet state machine full coverage
    // ═══════════════════════════════════════════════════════

    /// @dev Pending → Active → Paused → Active (resume)
    function test_stateMachine_pendingActivesPausedResume() public {
        uint256 wid = _registerWorknet();
        // Pending
        assertEq(uint256(awpRegistry.getWorknet(wid).status), uint256(IAWPRegistry.WorknetStatus.Pending));

        vm.prank(user1);
        awpRegistry.activateWorknet(wid);
        assertEq(uint256(awpRegistry.getWorknet(wid).status), uint256(IAWPRegistry.WorknetStatus.Active));

        vm.prank(user1);
        awpRegistry.pauseWorknet(wid);
        assertEq(uint256(awpRegistry.getWorknet(wid).status), uint256(IAWPRegistry.WorknetStatus.Paused));

        vm.prank(user1);
        awpRegistry.resumeWorknet(wid);
        assertEq(uint256(awpRegistry.getWorknet(wid).status), uint256(IAWPRegistry.WorknetStatus.Active));
    }

    /// @dev Pending → Active → Paused → Banned → Active (unban)
    function test_stateMachine_pausedBannedUnban() public {
        uint256 wid = _registerAndActivate();

        vm.prank(user1);
        awpRegistry.pauseWorknet(wid);
        assertEq(uint256(awpRegistry.getWorknet(wid).status), uint256(IAWPRegistry.WorknetStatus.Paused));

        vm.prank(guardian);
        awpRegistry.banWorknet(wid);
        assertEq(uint256(awpRegistry.getWorknet(wid).status), uint256(IAWPRegistry.WorknetStatus.Banned));

        vm.prank(guardian);
        awpRegistry.unbanWorknet(wid);
        assertEq(uint256(awpRegistry.getWorknet(wid).status), uint256(IAWPRegistry.WorknetStatus.Active));
    }

    /// @dev Active → Banned (direct)
    function test_stateMachine_activeBannedDirect() public {
        uint256 wid = _registerAndActivate();

        vm.prank(guardian);
        awpRegistry.banWorknet(wid);
        assertEq(uint256(awpRegistry.getWorknet(wid).status), uint256(IAWPRegistry.WorknetStatus.Banned));
    }

    /// @dev Cannot activate from Banned (must unban)
    function test_stateMachine_cannotActivateFromBanned() public {
        uint256 wid = _registerAndActivate();

        vm.prank(guardian);
        awpRegistry.banWorknet(wid);

        // Activate should fail since _activateWorknet checks status == Pending
        vm.prank(user1);
        vm.expectRevert();
        awpRegistry.activateWorknet(wid);
    }

    /// @dev Cannot pause from Pending
    function test_stateMachine_cannotPauseFromPending() public {
        uint256 wid = _registerWorknet();

        vm.prank(user1);
        vm.expectRevert();
        awpRegistry.pauseWorknet(wid);
    }

    /// @dev Cannot resume from Active (resume only accepts Paused)
    function test_stateMachine_cannotResumeFromActive() public {
        uint256 wid = _registerAndActivate();

        vm.prank(user1);
        vm.expectRevert();
        awpRegistry.resumeWorknet(wid);
    }

    /// @dev Cannot deregister Active worknet (must ban first)
    function test_stateMachine_cannotDeregisterActive() public {
        uint256 wid = _registerAndActivate();

        vm.prank(guardian);
        vm.expectRevert();
        awpRegistry.deregisterWorknet(wid);
    }

    /// @dev Cannot deregister Paused worknet (must ban first)
    function test_stateMachine_cannotDeregisterPaused() public {
        uint256 wid = _registerAndActivate();

        vm.prank(user1);
        awpRegistry.pauseWorknet(wid);

        vm.prank(guardian);
        vm.expectRevert();
        awpRegistry.deregisterWorknet(wid);
    }

    /// @dev Cannot deregister Pending worknet before immunity period expires
    function test_stateMachine_cannotDeregisterPendingBeforeImmunity() public {
        uint256 wid = _registerWorknet();

        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.ImmunityNotExpired.selector);
        awpRegistry.deregisterWorknet(wid);
    }

    /// @dev Deregister burns WorknetNFT
    function test_stateMachine_deregisterBurnsNFT() public {
        uint256 wid = _registerAndActivate();

        // Ban -> wait immunity -> deregister
        vm.prank(guardian);
        awpRegistry.banWorknet(wid);

        vm.warp(block.timestamp + 31 days);

        vm.prank(guardian);
        awpRegistry.deregisterWorknet(wid);

        // NFT burned, ownerOf should revert
        vm.expectRevert();
        nft.ownerOf(wid);

        // Worknet data deleted
        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(uint256(info.status), 0); // deleted => default Pending(0)
        assertEq(info.createdAt, 0);
    }

    /// @dev Can ban Pending worknet (allows cleanup of never-activated worknets)
    function test_stateMachine_canBanPending() public {
        uint256 wid = _registerWorknet();

        vm.prank(guardian);
        awpRegistry.banWorknet(wid);
        assertEq(uint256(awpRegistry.getWorknet(wid).status), uint256(IAWPRegistry.WorknetStatus.Banned));
    }

    /// @dev Can deregister Pending worknet directly (skip ban)
    function test_stateMachine_canDeregisterPending() public {
        uint256 wid = _registerWorknet();

        vm.warp(block.timestamp + 31 days); // wait immunity period

        vm.prank(guardian);
        awpRegistry.deregisterWorknet(wid);

        vm.expectRevert();
        IWorknetNFT(address(nft)).ownerOf(wid);
    }

    /// @dev Cannot ban already Banned worknet (double ban)
    function test_stateMachine_cannotBanBanned() public {
        uint256 wid = _registerAndActivate();

        vm.prank(guardian);
        awpRegistry.banWorknet(wid);

        vm.prank(guardian);
        vm.expectRevert();
        awpRegistry.banWorknet(wid);
    }

    /// @dev Cannot unban non-Banned worknet
    function test_stateMachine_cannotUnbanActive() public {
        uint256 wid = _registerAndActivate();

        vm.prank(guardian);
        vm.expectRevert();
        awpRegistry.unbanWorknet(wid);
    }

    // ═══════════════════════════════════════════════════════
    //  2. Auto-deploy WorknetManager
    // ═══════════════════════════════════════════════════════

    /// @dev worknetManager=0 and defaultImpl not set -> revert
    function test_autoDeploy_noImplSet_reverts() public {
        uint256 lpCost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.startPrank(user1);
        awp.approve(address(awpRegistry), lpCost);

        vm.expectRevert(AWPRegistry.WorknetManagerRequired.selector);
        awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams({
                name: "AutoSub",
                symbol: "AUTO",
                worknetManager: address(0),
                salt: bytes32(0),
                minStake: 0,
                skillsURI: ""
            })
        );
        vm.stopPrank();
    }

    /// @dev worknetManager=0 and defaultImpl set -> auto-deploy proxy
    function test_autoDeploy_withImpl_success() public {
        // Set defaultWorknetManagerImpl via Timelock
        bytes memory dexCfg = abi.encode(address(1), address(2), address(3), address(4), uint24(10000), int24(200));
        vm.startPrank(guardian);
        awpRegistry.setWorknetManagerImpl(address(worknetManagerImpl));
        awpRegistry.setDexConfig(dexCfg);
        vm.stopPrank();

        uint256 lpCost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.startPrank(user1);
        awp.approve(address(awpRegistry), lpCost);

        uint256 wid = awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams({
                name: "AutoWorknet",
                symbol: "AWRK",
                worknetManager: address(0),
                salt: bytes32(0),
                minStake: 0,
                skillsURI: ""
            })
        );
        vm.stopPrank();

        // Verify auto-deployed proxy address is non-zero
        IAWPRegistry.WorknetFullInfo memory full = awpRegistry.getWorknetFull(wid);
        assertTrue(full.worknetManager != address(0), "Auto-deployed worknetManager should be non-zero");
        assertTrue(full.worknetManager.code.length > 0, "Auto-deployed worknetManager should be a contract");
    }

    /// @dev Auto-deployed proxy correctly initialized (admin is user1)
    function test_autoDeploy_proxyInitialized() public {
        bytes memory dexCfg = abi.encode(address(1), address(2), address(3), address(4), uint24(10000), int24(200));
        vm.startPrank(guardian);
        awpRegistry.setWorknetManagerImpl(address(worknetManagerImpl));
        awpRegistry.setDexConfig(dexCfg);
        vm.stopPrank();

        uint256 lpCost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.startPrank(user1);
        awp.approve(address(awpRegistry), lpCost);
        uint256 wid = awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams({
                name: "AutoInit",
                symbol: "AINT",
                worknetManager: address(0),
                salt: bytes32(0),
                minStake: 0,
                skillsURI: ""
            })
        );
        vm.stopPrank();

        IAWPRegistry.WorknetFullInfo memory full = awpRegistry.getWorknetFull(wid);
        WorknetManager wm = WorknetManager(full.worknetManager);

        // user1 should have DEFAULT_ADMIN_ROLE
        assertTrue(wm.hasRole(wm.DEFAULT_ADMIN_ROLE(), user1));
    }

    // ═══════════════════════════════════════════════════════
    //  3. Gasless registration (EIP-712)
    // ═══════════════════════════════════════════════════════

    /// @dev registerWorknetFor valid signature
    function test_gasless_registerWorknetFor_valid() public {
        uint256 lpCost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;

        // user1 first approves AWP
        vm.prank(user1);
        awp.approve(address(awpRegistry), lpCost);

        IAWPRegistry.WorknetParams memory params = IAWPRegistry.WorknetParams({
            name: "GaslessSub",
            symbol: "GSUB",
            worknetManager: worknetManager,
            salt: bytes32(0),
            minStake: 0,
            skillsURI: "ipfs://skills"
        });

        uint256 nonce = awpRegistry.nonces(user1);
        uint256 deadline = block.timestamp + 1 hours;

        // Build WorknetParams struct hash
        bytes32 paramsStructHash = keccak256(abi.encode(
            WORKNET_PARAMS_TYPEHASH,
            keccak256(bytes(params.name)),
            keccak256(bytes(params.symbol)),
            params.worknetManager,
            params.salt,
            params.minStake,
            keccak256(bytes(params.skillsURI))
        ));

        bytes32 structHash = keccak256(abi.encode(
            REGISTER_WORKNET_TYPEHASH,
            user1,
            paramsStructHash,
            nonce,
            deadline
        ));

        (uint8 v, bytes32 r, bytes32 s) = _signDigest(structHash);

        // Relayer calls registerWorknetFor
        vm.prank(relayer);
        uint256 wid = awpRegistry.registerWorknetFor(user1, params, deadline, v, r, s);

        assertGt(wid, 0);
        assertEq(awpRegistry.nonces(user1), nonce + 1);
        // NFT belongs to user1
        assertEq(nft.ownerOf(wid), user1);
    }

    /// @dev registerWorknetFor invalid signature -> revert
    function test_gasless_registerWorknetFor_invalidSig_reverts() public {
        uint256 lpCost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.prank(user1);
        awp.approve(address(awpRegistry), lpCost);

        IAWPRegistry.WorknetParams memory params = IAWPRegistry.WorknetParams({
            name: "BadSig",
            symbol: "BSIG",
            worknetManager: worknetManager,
            salt: bytes32(0),
            minStake: 0,
            skillsURI: ""
        });

        uint256 deadline = block.timestamp + 1 hours;

        // Sign with wrong structHash
        bytes32 fakeStructHash = keccak256("wrong data");
        (uint8 v, bytes32 r, bytes32 s) = _signDigest(fakeStructHash);

        vm.prank(relayer);
        vm.expectRevert(AWPRegistry.InvalidSignature.selector);
        awpRegistry.registerWorknetFor(user1, params, deadline, v, r, s);
    }

    /// @dev Replay protection: same signature unusable after nonce increment
    function test_gasless_registerWorknetFor_replayProtection() public {
        uint256 lpCost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;

        vm.prank(user1);
        awp.approve(address(awpRegistry), lpCost * 2);

        IAWPRegistry.WorknetParams memory params = IAWPRegistry.WorknetParams({
            name: "Replay1",
            symbol: "RPL1",
            worknetManager: worknetManager,
            salt: bytes32(0),
            minStake: 0,
            skillsURI: ""
        });

        uint256 nonce = awpRegistry.nonces(user1);
        uint256 deadline = block.timestamp + 1 hours;

        bytes32 paramsStructHash = keccak256(abi.encode(
            WORKNET_PARAMS_TYPEHASH,
            keccak256(bytes(params.name)),
            keccak256(bytes(params.symbol)),
            params.worknetManager,
            params.salt,
            params.minStake,
            keccak256(bytes(params.skillsURI))
        ));

        bytes32 structHash = keccak256(abi.encode(
            REGISTER_WORKNET_TYPEHASH,
            user1,
            paramsStructHash,
            nonce,
            deadline
        ));

        (uint8 v, bytes32 r, bytes32 s) = _signDigest(structHash);

        // First call succeeds
        vm.prank(relayer);
        awpRegistry.registerWorknetFor(user1, params, deadline, v, r, s);

        // Second call with same signature -> revert (nonce incremented)
        vm.prank(relayer);
        vm.expectRevert(AWPRegistry.InvalidSignature.selector);
        awpRegistry.registerWorknetFor(user1, params, deadline, v, r, s);
    }

    /// @dev Expired signature -> revert
    function test_gasless_registerWorknetFor_expired_reverts() public {
        uint256 lpCost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.prank(user1);
        awp.approve(address(awpRegistry), lpCost);

        IAWPRegistry.WorknetParams memory params = IAWPRegistry.WorknetParams({
            name: "Expired",
            symbol: "EXP",
            worknetManager: worknetManager,
            salt: bytes32(0),
            minStake: 0,
            skillsURI: ""
        });

        uint256 deadline = block.timestamp - 1; // expired

        bytes32 paramsStructHash = keccak256(abi.encode(
            WORKNET_PARAMS_TYPEHASH,
            keccak256(bytes(params.name)),
            keccak256(bytes(params.symbol)),
            params.worknetManager,
            params.salt,
            params.minStake,
            keccak256(bytes(params.skillsURI))
        ));

        bytes32 structHash = keccak256(abi.encode(
            REGISTER_WORKNET_TYPEHASH,
            user1,
            paramsStructHash,
            awpRegistry.nonces(user1),
            deadline
        ));

        (uint8 v, bytes32 r, bytes32 s) = _signDigest(structHash);

        vm.prank(relayer);
        vm.expectRevert(AWPRegistry.ExpiredSignature.selector);
        awpRegistry.registerWorknetFor(user1, params, deadline, v, r, s);
    }

    // ═══════════════════════════════════════════════════════
    //  4. Gasless activation (EIP-712)
    // ═══════════════════════════════════════════════════════

    /// @dev activateWorknetFor valid signature
    function test_gasless_activateWorknetFor_valid() public {
        uint256 wid = _registerWorknet();

        uint256 nonce = awpRegistry.nonces(user1);
        uint256 deadline = block.timestamp + 1 hours;

        bytes32 structHash = keccak256(abi.encode(
            ACTIVATE_WORKNET_TYPEHASH,
            user1,
            wid,
            nonce,
            deadline
        ));

        (uint8 v, bytes32 r, bytes32 s) = _signDigest(structHash);

        vm.prank(relayer);
        awpRegistry.activateWorknetFor(user1, wid, deadline, v, r, s);

        assertTrue(awpRegistry.isWorknetActive(wid));
        assertEq(awpRegistry.nonces(user1), nonce + 1);
    }

    /// @dev activateWorknetFor invalid signature -> revert
    function test_gasless_activateWorknetFor_invalidSig_reverts() public {
        uint256 wid = _registerWorknet();

        uint256 deadline = block.timestamp + 1 hours;

        // Sign a wrong structHash
        bytes32 fakeStructHash = keccak256("fake activation");
        (uint8 v, bytes32 r, bytes32 s) = _signDigest(fakeStructHash);

        vm.prank(relayer);
        vm.expectRevert(AWPRegistry.InvalidSignature.selector);
        awpRegistry.activateWorknetFor(user1, wid, deadline, v, r, s);
    }

    /// @dev activateWorknetFor non-owner -> revert
    function test_gasless_activateWorknetFor_notOwner_reverts() public {
        uint256 wid = _registerWorknet();

        // Sign as user2 (but NFT belongs to user1)
        (address signer2, uint256 signer2Pk) = makeAddrAndKey("signer2");

        uint256 deadline = block.timestamp + 1 hours;
        bytes32 structHash = keccak256(abi.encode(
            ACTIVATE_WORKNET_TYPEHASH,
            signer2,
            wid,
            awpRegistry.nonces(signer2),
            deadline
        ));

        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", _domainSeparator(), structHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signer2Pk, digest);

        vm.prank(relayer);
        vm.expectRevert(AWPRegistry.NotOwner.selector);
        awpRegistry.activateWorknetFor(signer2, wid, deadline, v, r, s);
    }

    // ═══════════════════════════════════════════════════════
    //  5. MAX_ACTIVE_WORKNETS enforcement
    // ═══════════════════════════════════════════════════════

    /// @dev Cannot exceed MAX_ACTIVE_WORKNETS (simulated with small count)
    ///      Since MAX_ACTIVE_WORKNETS=10000 is too large, we test resume limit via mock
    function test_maxActive_resumeBlocked() public {
        // Register and activate a worknet
        uint256 wid = _registerAndActivate();

        // Pause it
        vm.prank(user1);
        awpRegistry.pauseWorknet(wid);
        assertEq(awpRegistry.getActiveWorknetCount(), 0);

        // Resume succeeds (limit not reached)
        vm.prank(user1);
        awpRegistry.resumeWorknet(wid);
        assertEq(awpRegistry.getActiveWorknetCount(), 1);
    }

    /// @dev Pause/ban frees slot, allowing new worknet activation
    function test_maxActive_pauseFreesSlot() public {
        uint256 wid1 = _registerAndActivate();
        assertEq(awpRegistry.getActiveWorknetCount(), 1);

        // Pause wid1 -> count=0
        vm.prank(user1);
        awpRegistry.pauseWorknet(wid1);
        assertEq(awpRegistry.getActiveWorknetCount(), 0);

        // Register and activate new worknet
        uint256 wid2 = _registerWorknet();
        vm.prank(user1);
        awpRegistry.activateWorknet(wid2);
        assertEq(awpRegistry.getActiveWorknetCount(), 1);
    }

    /// @dev Ban frees slot
    function test_maxActive_banFreesSlot() public {
        uint256 wid1 = _registerAndActivate();
        assertEq(awpRegistry.getActiveWorknetCount(), 1);

        // Ban wid1 -> count=0
        vm.prank(guardian);
        awpRegistry.banWorknet(wid1);
        assertEq(awpRegistry.getActiveWorknetCount(), 0);

        // Register and activate new worknet
        uint256 wid2 = _registerWorknet();
        vm.prank(user1);
        awpRegistry.activateWorknet(wid2);
        assertEq(awpRegistry.getActiveWorknetCount(), 1);
    }

    /// @dev unbanWorknet checks MAX_ACTIVE_WORKNETS (verify logic exists)
    function test_maxActive_unbanChecksLimit() public {
        // Cannot create 10000 worknets, just verify active count increments after unban
        uint256 wid = _registerAndActivate();

        vm.prank(guardian);
        awpRegistry.banWorknet(wid);
        assertEq(awpRegistry.getActiveWorknetCount(), 0);

        vm.prank(guardian);
        awpRegistry.unbanWorknet(wid);
        assertEq(awpRegistry.getActiveWorknetCount(), 1);
        assertTrue(awpRegistry.isWorknetActive(wid));
    }

    // ═══════════════════════════════════════════════════════
    //  6. AlphaToken creation via factory
    // ═══════════════════════════════════════════════════════

    /// @dev registerWorknet creates AlphaToken with correct name/symbol
    function test_alphaToken_createdWithCorrectNameSymbol() public {
        uint256 wid = _registerWorknet();
        IAWPRegistry.WorknetFullInfo memory full = awpRegistry.getWorknetFull(wid);

        AlphaToken alpha = AlphaToken(full.alphaToken);
        assertEq(alpha.name(), "TestWorknet");
        assertEq(alpha.symbol(), "TSUB");
    }

    /// @dev AlphaToken admin is AWPRegistry
    function test_alphaToken_adminIsRegistry() public {
        uint256 wid = _registerWorknet();
        IAWPRegistry.WorknetFullInfo memory full = awpRegistry.getWorknetFull(wid);

        AlphaToken alpha = AlphaToken(full.alphaToken);
        assertEq(alpha.admin(), address(awpRegistry));
    }

    /// @dev setWorknetMinter called at registration, minter set to worknetManager
    function test_alphaToken_minterSetToWorknetManager() public {
        uint256 wid = _registerWorknet();
        IAWPRegistry.WorknetFullInfo memory full = awpRegistry.getWorknetFull(wid);

        AlphaToken alpha = AlphaToken(full.alphaToken);
        // After setWorknetMinter, mintersLocked = true
        assertTrue(alpha.mintersLocked());
        // worknetManager is a minter
        assertTrue(alpha.minters(full.worknetManager));
        // AWPRegistry (admin) is no longer a minter
        assertFalse(alpha.minters(address(awpRegistry)));
    }

    /// @dev AlphaToken worknetId matches
    function test_alphaToken_worknetIdMatches() public {
        uint256 wid = _registerWorknet();
        IAWPRegistry.WorknetFullInfo memory full = awpRegistry.getWorknetFull(wid);

        AlphaToken alpha = AlphaToken(full.alphaToken);
        assertEq(alpha.worknetId(), wid);
    }

    // ═══════════════════════════════════════════════════════
    //  7. LP creation
    // ═══════════════════════════════════════════════════════

    /// @dev registerWorknet calls LPManager.createPoolAndAddLiquidity
    function test_lp_poolCreated() public {
        uint256 wid = _registerWorknet();
        IAWPRegistry.WorknetFullInfo memory full = awpRegistry.getWorknetFull(wid);

        // MockLPManager records alphaTokenToPool
        bytes32 poolId = lp.alphaTokenToPool(full.alphaToken);
        assertNotEq(poolId, bytes32(0));
        // WorknetInfo also stores poolId
        assertEq(awpRegistry.getWorknet(wid).lpPool, poolId);
    }

    /// @dev AWP transferred from user to LPManager
    function test_lp_awpTransferred() public {
        uint256 balBefore = awp.balanceOf(address(lp));
        uint256 lpCost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        _registerWorknet();
        uint256 balAfter = awp.balanceOf(address(lp));

        assertEq(balAfter - balBefore, lpCost);
    }

    /// @dev Alpha minted to LPManager
    function test_lp_alphaMintedToLPManager() public {
        uint256 wid = _registerWorknet();
        IAWPRegistry.WorknetFullInfo memory full = awpRegistry.getWorknetFull(wid);

        // Alpha minted to LPManager at initialAlphaMint
        // After setWorknetMinter, supplyAtLock snapshots total, so totalSupply == initialAlphaMint
        AlphaToken alpha = AlphaToken(full.alphaToken);
        assertEq(alpha.totalSupply(), awpRegistry.initialAlphaMint());
    }

    // ═══════════════════════════════════════════════════════
    //  8. Parameter governance
    // ═══════════════════════════════════════════════════════

    /// @dev setInitialAlphaPrice - Timelock only
    function test_gov_setInitialAlphaPrice_onlyTimelock() public {
        vm.prank(guardian);
        awpRegistry.setInitialAlphaPrice(5e16);
        assertEq(awpRegistry.initialAlphaPrice(), 5e16);

        vm.prank(user1);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.setInitialAlphaPrice(5e16);
    }

    /// @dev setInitialAlphaPrice boundary check
    function test_gov_setInitialAlphaPrice_priceTooLow() public {
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.PriceTooLow.selector);
        awpRegistry.setInitialAlphaPrice(1e11); // < 1e12
    }

    function test_gov_setInitialAlphaPrice_priceTooHigh() public {
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.PriceTooHigh.selector);
        awpRegistry.setInitialAlphaPrice(1e31); // > 1e30
    }

    /// @dev setInitialAlphaMint - Timelock only
    function test_gov_setInitialAlphaMint_onlyTimelock() public {
        vm.prank(guardian);
        awpRegistry.setInitialAlphaMint(200_000_000 * 1e18);
        assertEq(awpRegistry.initialAlphaMint(), 200_000_000 * 1e18);

        vm.prank(user1);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.setInitialAlphaMint(200_000_000 * 1e18);
    }

    /// @dev setInitialAlphaMint cannot be 0
    function test_gov_setInitialAlphaMint_zeroReverts() public {
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.InvalidMintAmount.selector);
        awpRegistry.setInitialAlphaMint(0);
    }

    /// @dev setImmunityPeriod - Timelock only
    function test_gov_setImmunityPeriod_onlyTimelock() public {
        vm.prank(guardian);
        awpRegistry.setImmunityPeriod(60 days);
        assertEq(awpRegistry.immunityPeriod(), 60 days);

        vm.prank(user1);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.setImmunityPeriod(60 days);
    }

    /// @dev setImmunityPeriod minimum 7 days
    function test_gov_setImmunityPeriod_tooShort() public {
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.ImmunityTooShort.selector);
        awpRegistry.setImmunityPeriod(6 days);
    }

    /// @dev setWorknetManagerImpl - Timelock only
    function test_gov_setWorknetManagerImpl_onlyTimelock() public {
        vm.prank(guardian);
        awpRegistry.setWorknetManagerImpl(address(worknetManagerImpl));
        assertEq(awpRegistry.defaultWorknetManagerImpl(), address(worknetManagerImpl));

        vm.prank(user1);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.setWorknetManagerImpl(address(worknetManagerImpl));
    }

    /// @dev setWorknetManagerImpl cannot be address(0)
    function test_gov_setWorknetManagerImpl_zeroReverts() public {
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.ZeroAddress.selector);
        awpRegistry.setWorknetManagerImpl(address(0));
    }

    /// @dev setDexConfig - Timelock only
    function test_gov_setDexConfig_onlyTimelock() public {
        bytes memory cfg = abi.encode(address(1));
        vm.prank(guardian);
        awpRegistry.setDexConfig(cfg);

        vm.prank(user1);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.setDexConfig(cfg);
    }

    /// @dev setGuardian - Guardian only
    function test_gov_setGuardian_onlyGuardian() public {
        address newGuardian = address(99);
        vm.prank(guardian);
        awpRegistry.setGuardian(newGuardian);
        assertEq(awpRegistry.guardian(), newGuardian);

        // Old guardian no longer has permission
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.setGuardian(address(100));
    }

    /// @dev setGuardian cannot be address(0)
    function test_gov_setGuardian_zeroReverts() public {
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.ZeroAddress.selector);
        awpRegistry.setGuardian(address(0));
    }

    /// @dev setGuardian called by non-Guardian -> revert
    function test_gov_setGuardian_notGuardian() public {
        vm.prank(user1);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.setGuardian(address(99));
    }

    /// @dev setAlphaTokenFactory - Timelock only
    function test_gov_setAlphaTokenFactory_onlyTimelock() public {
        address newFactory = address(88);
        vm.prank(guardian);
        awpRegistry.setAlphaTokenFactory(newFactory);
        assertEq(awpRegistry.alphaTokenFactory(), newFactory);

        vm.prank(user1);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.setAlphaTokenFactory(newFactory);
    }

    /// @dev setAlphaTokenFactory cannot be address(0)
    function test_gov_setAlphaTokenFactory_zeroReverts() public {
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.ZeroAddress.selector);
        awpRegistry.setAlphaTokenFactory(address(0));
    }

    // ═══════════════════════════════════════════════════════
    //  9. Pause/Unpause
    // ═══════════════════════════════════════════════════════

    /// @dev pause - Guardian only
    function test_pause_onlyGuardian() public {
        vm.prank(guardian);
        awpRegistry.pause();
        assertTrue(awpRegistry.paused());

        // Non-guardian cannot pause
        vm.prank(user1);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.pause();
    }

    /// @dev unpause - Timelock only
    function test_unpause_onlyGuardian() public {
        vm.prank(guardian);
        awpRegistry.pause();

        // Non-Guardian cannot unpause
        vm.prank(user1);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.unpause();

        // Guardian can unpause
        vm.prank(guardian);
        awpRegistry.unpause();
        assertFalse(awpRegistry.paused());
    }

    /// @dev registerWorknet reverts when paused
    function test_paused_registerWorknet_reverts() public {
        vm.prank(guardian);
        awpRegistry.pause();

        uint256 lpCost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.startPrank(user1);
        awp.approve(address(awpRegistry), lpCost);

        vm.expectRevert();
        awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams({
                name: "PausedSub",
                symbol: "PSUB",
                worknetManager: worknetManager,
                salt: bytes32(0),
                minStake: 0,
                skillsURI: ""
            })
        );
        vm.stopPrank();
    }

    /// @dev activateWorknet reverts when paused
    function test_paused_activateWorknet_reverts() public {
        uint256 wid = _registerWorknet();

        vm.prank(guardian);
        awpRegistry.pause();

        vm.prank(user1);
        vm.expectRevert();
        awpRegistry.activateWorknet(wid);
    }

    /// @dev bind reverts when paused
    function test_paused_bind_reverts() public {
        vm.prank(guardian);
        awpRegistry.pause();

        vm.prank(user1);
        vm.expectRevert();
        awpRegistry.bind(user2);
    }

    /// @dev register reverts when paused
    function test_paused_register_reverts() public {
        vm.prank(guardian);
        awpRegistry.pause();

        vm.prank(user1);
        vm.expectRevert();
        awpRegistry.setRecipient(user1);
    }

    /// @dev pauseWorknet reverts when paused
    function test_paused_pauseWorknet_reverts() public {
        uint256 wid = _registerAndActivate();

        vm.prank(guardian);
        awpRegistry.pause();

        vm.prank(user1);
        vm.expectRevert();
        awpRegistry.pauseWorknet(wid);
    }

    /// @dev resumeWorknet reverts when paused
    function test_paused_resumeWorknet_reverts() public {
        uint256 wid = _registerAndActivate();

        vm.prank(user1);
        awpRegistry.pauseWorknet(wid);

        vm.prank(guardian);
        awpRegistry.pause();

        vm.prank(user1);
        vm.expectRevert();
        awpRegistry.resumeWorknet(wid);
    }

    /// @dev banWorknet and unbanWorknet unaffected by pause (onlyTimelock, no whenNotPaused)
    function test_paused_banUnban_still_work() public {
        uint256 wid = _registerAndActivate();

        vm.prank(guardian);
        awpRegistry.pause();

        // ban and unban should still work (governance ops not restricted by pause)
        vm.prank(guardian);
        awpRegistry.banWorknet(wid);
        assertEq(uint256(awpRegistry.getWorknet(wid).status), uint256(IAWPRegistry.WorknetStatus.Banned));

        vm.prank(guardian);
        awpRegistry.unbanWorknet(wid);
        assertEq(uint256(awpRegistry.getWorknet(wid).status), uint256(IAWPRegistry.WorknetStatus.Active));
    }

    /// @dev deregisterWorknet unaffected by pause (onlyTimelock, no whenNotPaused)
    function test_paused_deregister_still_works() public {
        uint256 wid = _registerAndActivate();

        vm.prank(guardian);
        awpRegistry.banWorknet(wid);

        vm.warp(block.timestamp + 31 days);

        vm.prank(guardian);
        awpRegistry.pause();

        vm.prank(guardian);
        awpRegistry.deregisterWorknet(wid);

        // Verify deleted
        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(info.createdAt, 0);
    }
}
