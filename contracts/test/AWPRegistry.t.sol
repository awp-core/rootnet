// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AlphaTokenFactory} from "../src/token/AlphaTokenFactory.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {WorknetNFT} from "../src/core/WorknetNFT.sol";
import {MockLPManager} from "./helpers/MockLPManager.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {IAWPRegistry} from "../src/interfaces/IAWPRegistry.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {WorknetManager} from "../src/worknets/WorknetManager.sol";

contract AWPRegistryTest is Test {
    AWPToken awp;
    AlphaTokenFactory factory;
    AWPEmission emission;
    StakingVault vault;
    StakeNFT stakeNFT;
    WorknetNFT nft;
    MockLPManager lp;
    AWPRegistry awpRegistry;
    Treasury treasury;

    address deployer = address(1);
    address guardian = address(2);
    address user1 = address(3);
    address agent1 = address(4);
    address worknetManager = address(5);
    address user2 = address(6);

    uint256 constant INITIAL_DAILY_EMISSION = 31_600_000 * 1e18;
    uint256 constant EPOCH_DURATION = 1 days;

    function setUp() public {
        vm.startPrank(deployer);

        // Deploy tokens
        awp = new AWPToken("AWP Token", "AWP", deployer);
        awp.initialMint(200_000_000 * 1e18);

        // Deploy Treasury
        address[] memory proposers = new address[](0);
        address[] memory executors = new address[](1);
        executors[0] = address(0);
        treasury = new Treasury(0, proposers, executors, deployer);

        // Deploy AWPRegistry
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

        // Set AWP minter
        awp.addMinter(address(emission));
        awp.renounceAdmin();

        // Configure factory
        factory.setAddresses(address(awpRegistry));

        // Deploy StakeNFT and StakingVault
        vault = StakingVault(address(new ERC1967Proxy(
            address(new StakingVault()), abi.encodeCall(StakingVault.initialize, (address(awpRegistry), deployer))
        )));
        stakeNFT = new StakeNFT(address(awp), address(vault), address(awpRegistry));

        // Initialize registry (no accessManager parameter)
        awpRegistry.initializeRegistry(
            address(awp),
            address(nft),
            address(factory),
            address(emission),
            address(lp),
            address(vault),
            address(stakeNFT),
            address(0), // no default WorknetManager impl in unit tests
            "" // no dexConfig in unit tests
        );

        // Give users some AWP
        awp.transfer(user1, 10_000_000 * 1e18);
        awp.transfer(user2, 10_000_000 * 1e18);

        vm.stopPrank();
    }

    // ── Registry tests ──

    function test_initializeRegistryOnce() public {
        vm.expectRevert(AWPRegistry.NotDeployer.selector);
        awpRegistry.initializeRegistry(
            address(awp), address(nft), address(factory), address(emission),
            address(lp), address(vault), address(stakeNFT), address(0), ""
        );
    }

    // ── Registration ──

    function test_register() public {
        vm.prank(user1);
        awpRegistry.setRecipient(user1);
        assertTrue(awpRegistry.isRegistered(user1));
    }


    // ── Binding ──

    function test_bind() public {
        vm.prank(agent1);
        awpRegistry.bind(user1);
        assertEq(awpRegistry.boundTo(agent1), user1);
    }

    function test_bind_selfBind_reverts() public {
        vm.prank(user1);
        vm.expectRevert(AWPRegistry.SelfBind.selector);
        awpRegistry.bind(user1);
    }

    function test_bind_zeroAddress_reverts() public {
        vm.prank(agent1);
        vm.expectRevert(AWPRegistry.ZeroAddress.selector);
        awpRegistry.bind(address(0));
    }

    function test_bind_rebind() public {
        vm.prank(agent1);
        awpRegistry.bind(user1);
        assertEq(awpRegistry.boundTo(agent1), user1);

        // rebind to user2
        vm.prank(agent1);
        awpRegistry.bind(user2);
        assertEq(awpRegistry.boundTo(agent1), user2);
    }

    function test_bind_antiCycle() public {
        // A -> B -> C, then C tries to bind to A => cycle
        vm.prank(address(0x10));
        awpRegistry.bind(address(0x11));
        vm.prank(address(0x11));
        awpRegistry.bind(address(0x12));

        vm.prank(address(0x12));
        vm.expectRevert(AWPRegistry.CycleDetected.selector);
        awpRegistry.bind(address(0x10));
    }

    // ── Unbind ──

    function test_unbind() public {
        vm.prank(agent1);
        awpRegistry.bind(user1);
        assertEq(awpRegistry.boundTo(agent1), user1);

        vm.prank(agent1);
        awpRegistry.unbind();
        assertEq(awpRegistry.boundTo(agent1), address(0));
    }

    // ── Recipient ──

    function test_setRecipient() public {
        vm.prank(user1);
        awpRegistry.setRecipient(user2);
        assertEq(awpRegistry.recipient(user1), user2);
    }

    function test_resolveRecipient() public {
        // Set up: agent1 -> user1, user1 has recipient = user2
        vm.prank(user1);
        awpRegistry.setRecipient(user2);
        vm.prank(agent1);
        awpRegistry.bind(user1);

        // resolveRecipient(agent1) should walk to user1 and return user2
        assertEq(awpRegistry.resolveRecipient(agent1), user2);
    }

    function test_resolveRecipient_unregistered() public {
        // No binding, no recipient => returns the address itself
        assertEq(awpRegistry.resolveRecipient(user1), user1);
    }

    // ── Delegation ──

    function test_grantAndRevokeDelegate() public {
        vm.prank(user1);
        awpRegistry.grantDelegate(user2);
        assertTrue(awpRegistry.delegates(user1, user2));

        vm.prank(user1);
        awpRegistry.revokeDelegate(user2);
        assertFalse(awpRegistry.delegates(user1, user2));
    }

    function test_revokeDelegate_self_reverts() public {
        vm.prank(user1);
        vm.expectRevert(AWPRegistry.CannotRevokeSelf.selector);
        awpRegistry.revokeDelegate(user1);
    }

    // ── Staking (allocation with explicit staker) ──

    function test_allocateAndDeallocate() public {
        // Setup: deposit via StakeNFT, register worknet
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        uint256 sid = _registerWorknet();

        vm.prank(user1);
        awpRegistry.activateWorknet(sid);

        // Allocate (staker = user1, agent = agent1)
        vm.prank(user1);
        vault.allocate(user1, agent1, sid, 500 * 1e18);
        assertEq(vault.getAgentStake(user1, agent1, sid), 500 * 1e18);

        // Deallocate
        vm.prank(user1);
        vault.deallocate(user1, agent1, sid, 200 * 1e18);
        assertEq(vault.getAgentStake(user1, agent1, sid), 300 * 1e18);
    }

    function test_allocate_delegate() public {
        // Setup: deposit
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        awpRegistry.grantDelegate(user2);
        vm.stopPrank();

        uint256 sid = _registerWorknet();

        vm.prank(user1);
        awpRegistry.activateWorknet(sid);

        // user2 allocates on behalf of user1
        vm.prank(user2);
        vault.allocate(user1, agent1, sid, 500 * 1e18);
        assertEq(vault.getAgentStake(user1, agent1, sid), 500 * 1e18);
    }

    function test_allocate_notAuthorized_reverts() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        uint256 sid = _registerWorknet();

        vm.prank(user1);
        awpRegistry.activateWorknet(sid);

        // user2 has no delegation from user1
        vm.prank(user2);
        vm.expectRevert(StakingVault.NotAuthorized.selector);
        vault.allocate(user1, agent1, sid, 500 * 1e18);
    }

    function test_reallocate_immediate() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        uint256 sid = _registerWorknet();

        vm.prank(user1);
        awpRegistry.activateWorknet(sid);

        address agent2 = address(10);

        // Allocate to agent1/worknet
        vm.prank(user1);
        vault.allocate(user1, agent1, sid, 500 * 1e18);

        // Reallocate to agent2/worknet — takes effect immediately
        vm.prank(user1);
        vault.reallocate(user1, agent1, sid, agent2, sid, 200 * 1e18);

        assertEq(vault.getAgentStake(user1, agent1, sid), 300 * 1e18);
        assertEq(vault.getAgentStake(user1, agent2, sid), 200 * 1e18);
    }

    // ── Worknet registration ──

    function test_registerWorknet() public {
        uint256 worknetId = _registerWorknet();
        assertEq(worknetId & ((1 << 64) - 1), 1);
        assertEq(worknetId >> 64, block.chainid);

        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(worknetId);
        assertEq(awpRegistry.getWorknetFull(worknetId).worknetManager, worknetManager);
        assertTrue(info.status == IAWPRegistry.WorknetStatus.Pending);
    }

    function test_registerWorknetInvalidParams() public {
        vm.startPrank(user1);
        awp.approve(address(awpRegistry), 2_000_000 * 1e18);

        // Empty name
        vm.expectRevert(AWPRegistry.InvalidWorknetName.selector);
        awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams({
                name: "",
                symbol: "TEST",
                worknetManager: worknetManager,
                salt: bytes32(0),
                minStake: 0,
                skillsURI: ""
            })
        );

        // Empty worknet contract address
        vm.expectRevert(AWPRegistry.WorknetManagerRequired.selector);
        awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams({
                name: "Test",
                symbol: "TEST",
                worknetManager: address(0),
                salt: bytes32(0),
                minStake: 0,
                skillsURI: ""
            })
        );
        vm.stopPrank();
    }

    // ── Worknet lifecycle ──

    function test_activateWorknet() public {
        uint256 worknetId = _registerWorknet();

        vm.prank(user1);
        awpRegistry.activateWorknet(worknetId);

        assertTrue(awpRegistry.isWorknetActive(worknetId));
        assertEq(awpRegistry.getActiveWorknetCount(), 1);
    }

    function test_pauseAndResumeWorknet() public {
        uint256 worknetId = _registerWorknet();
        vm.startPrank(user1);
        awpRegistry.activateWorknet(worknetId);
        awpRegistry.pauseWorknet(worknetId);
        assertFalse(awpRegistry.isWorknetActive(worknetId));

        awpRegistry.resumeWorknet(worknetId);
        assertTrue(awpRegistry.isWorknetActive(worknetId));
        vm.stopPrank();
    }

    function test_banAndUnbanWorknet() public {
        uint256 worknetId = _registerWorknet();
        vm.prank(user1);
        awpRegistry.activateWorknet(worknetId);

        vm.prank(guardian);
        awpRegistry.banWorknet(worknetId);

        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(worknetId);
        assertTrue(info.status == IAWPRegistry.WorknetStatus.Banned);

        vm.prank(guardian);
        awpRegistry.unbanWorknet(worknetId);
        info = awpRegistry.getWorknet(worknetId);
        assertTrue(info.status == IAWPRegistry.WorknetStatus.Active);
    }

    function test_deregisterWorknet() public {
        uint256 worknetId = _registerWorknet();
        vm.prank(user1);
        awpRegistry.activateWorknet(worknetId);

        // Cannot deregister Active worknet (must ban first)
        vm.prank(guardian);
        vm.expectRevert();
        awpRegistry.deregisterWorknet(worknetId);

        // Ban the worknet first
        vm.prank(guardian);
        awpRegistry.banWorknet(worknetId);

        // Cannot deregister during immunity period
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.ImmunityNotExpired.selector);
        awpRegistry.deregisterWorknet(worknetId);

        // After immunity period, deregister succeeds
        vm.warp(block.timestamp + 31 days);
        vm.prank(guardian);
        awpRegistry.deregisterWorknet(worknetId);

        assertEq(awpRegistry.getActiveWorknetCount(), 0);
    }

    // ── UUPS Upgrade tests ──

    function test_upgradeViaGuardian() public {
        AWPRegistry newImpl = new AWPRegistry();
        vm.prank(guardian);
        awpRegistry.upgradeToAndCall(address(newImpl), "");
    }

    function test_upgrade_revertsForNonGuardian() public {
        AWPRegistry newImpl = new AWPRegistry();
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.upgradeToAndCall(address(newImpl), "");
    }

    function test_cannotInitializeImpl() public {
        AWPRegistry impl = new AWPRegistry();
        vm.expectRevert();
        impl.initialize(deployer, address(treasury), guardian);
    }

    // ── WorknetId encoding tests ──

    function test_worknetIdEncodesChainId() public {
        uint256 worknetId = _registerWorknet();
        uint256 expectedChainId = block.chainid;
        assertEq(worknetId >> 64, expectedChainId);
        assertEq(worknetId & ((1 << 64) - 1), 1);
    }

    function test_worknetIdIncrementsLocalCounter() public {
        uint256 id1 = _registerWorknet();
        vm.startPrank(user1);
        awp.approve(address(awpRegistry), 1_000_000 * 1e18);
        uint256 id2 = awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams("Sub2", "S2", worknetManager, bytes32(0), 0, "")
        );
        vm.stopPrank();
        assertEq(id1 >> 64, id2 >> 64);
        assertEq((id1 & ((1 << 64) - 1)) + 1, id2 & ((1 << 64) - 1));
    }

    function test_extractChainId() public view {
        uint256 worknetId = (uint256(8453) << 64) | 42;
        assertEq(awpRegistry.extractChainId(worknetId), 8453);
        assertEq(awpRegistry.extractLocalId(worknetId), 42);
    }

    // ── WorknetManager UUPS Upgrade tests ──

    function test_worknetManagerUpgradeByAdmin() public {
        // Deploy WorknetManager impl + proxy directly (no AWPRegistry auto-deploy, no DEX needed)
        WorknetManager smImpl = new WorknetManager();
        bytes memory dexCfg = abi.encode(address(1), address(2), address(3), address(4), uint24(10000), int24(200));
        bytes memory initData = abi.encodeCall(WorknetManager.initialize, (address(awpRegistry), address(awp), address(awp), bytes32(0), user1, dexCfg));
        address worknetProxy = address(new ERC1967Proxy(address(smImpl), initData));

        // user1 is DEFAULT_ADMIN_ROLE — can upgrade
        WorknetManager newImpl = new WorknetManager();
        vm.prank(user1);
        WorknetManager(worknetProxy).upgradeToAndCall(address(newImpl), "");
    }

    function test_worknetManagerUpgradeRevertsForNonAdmin() public {
        WorknetManager smImpl = new WorknetManager();
        bytes memory dexCfg = abi.encode(address(1), address(2), address(3), address(4), uint24(10000), int24(200));
        bytes memory initData = abi.encodeCall(WorknetManager.initialize, (address(awpRegistry), address(awp), address(awp), bytes32(0), user1, dexCfg));
        address worknetProxy = address(new ERC1967Proxy(address(smImpl), initData));

        // user2 has no role — upgrade should revert
        WorknetManager newImpl = new WorknetManager();
        vm.prank(user2);
        vm.expectRevert();
        WorknetManager(worknetProxy).upgradeToAndCall(address(newImpl), "");
    }

    // ── Cross-chain allocate tests ──

    function test_allocateToCrossChainWorknet() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 10_000 * 1e18);
        stakeNFT.deposit(10_000 * 1e18, 52 weeks);
        vm.stopPrank();

        // Allocate to a "foreign" worknetId (Arbitrum chain, local ID 5)
        uint256 foreignWorknetId = (uint256(42161) << 64) | 5;
        vm.prank(user1);
        vault.allocate(user1, user1, foreignWorknetId, 5_000 * 1e18);

        assertEq(vault.getAgentStake(user1, user1, foreignWorknetId), 5_000 * 1e18);
    }

    function test_reallocateToCrossChainWorknet() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 10_000 * 1e18);
        stakeNFT.deposit(10_000 * 1e18, 52 weeks);
        vm.stopPrank();

        uint256 localWorknetId = (block.chainid << 64) | 999;
        uint256 foreignWorknetId = (uint256(42161) << 64) | 5;

        vm.startPrank(user1);
        vault.allocate(user1, user1, localWorknetId, 5_000 * 1e18);
        vault.reallocate(user1, user1, localWorknetId, user1, foreignWorknetId, 2_000 * 1e18);
        vm.stopPrank();

        assertEq(vault.getAgentStake(user1, user1, localWorknetId), 3_000 * 1e18);
        assertEq(vault.getAgentStake(user1, user1, foreignWorknetId), 2_000 * 1e18);
    }

    // ── Permission tests ──

    function test_onlyTimelockFunctions() public {
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.setInitialAlphaPrice(1e15);

        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.setGuardian(address(0));

        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.unpause();
    }

    function test_onlyGuardianPause() public {
        vm.prank(guardian);
        awpRegistry.pause();
        assertTrue(awpRegistry.paused());

        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.pause();
    }

    // ── Queries ──

    function test_getAgentInfo() public {
        vm.prank(user1);
        awpRegistry.setRecipient(user1);
        vm.prank(agent1);
        awpRegistry.bind(user1);

        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        uint256 sid = _registerWorknet();

        vm.prank(user1);
        awpRegistry.activateWorknet(sid);

        vm.prank(user1);
        vault.allocate(user1, agent1, sid, 500 * 1e18);

        AWPRegistry.AgentInfo memory info = awpRegistry.getAgentInfo(agent1, sid);
        assertEq(info.root, user1);
        assertTrue(info.isValid);
        assertEq(info.stake, 500 * 1e18);
        assertEq(info.rewardRecipient, user1);
    }

    function test_getRegistry() public view {
        (address a, address b, address c, address d, address e, address f, address g,,) =
            awpRegistry.getRegistry();
        assertEq(a, address(awp));
        assertEq(b, address(nft));
        assertEq(c, address(factory));
        assertEq(d, address(emission));
        assertEq(e, address(lp));
        assertEq(f, address(vault));
        assertEq(g, address(stakeNFT));
    }

    // ── Helper functions ──

    function _registerWorknet() internal returns (uint256) {
        uint256 lpCost = 100_000_000 * 1e18 * 1e16 / 1e18;
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

    // ── Unbind idempotent ──

    function test_unbind_idempotent() public {
        // unbind when already unbound — should not revert
        vm.prank(user1);
        awpRegistry.unbind();
        assertEq(awpRegistry.boundTo(user1), address(0));
    }

    // ── resolveRecipient deep chain ──

    function test_resolveRecipient_deepChain() public {
        // Build chain: agent3 → agent2 → agent1 → user1 (root with recipient)
        address agent2 = makeAddr("agent2");
        address agent3 = makeAddr("agent3");

        vm.prank(user1);
        awpRegistry.setRecipient(address(0xBEEF));

        vm.prank(agent1);
        awpRegistry.bind(user1);

        vm.prank(agent2);
        awpRegistry.bind(agent1);

        vm.prank(agent3);
        awpRegistry.bind(agent2);

        // resolveRecipient(agent3) should walk up to user1 and return 0xBEEF
        assertEq(awpRegistry.resolveRecipient(agent3), address(0xBEEF));
    }

    // ── Bind chain depth (no artificial limit, uses gas as natural bound) ──

    // ── minStake enforcement ──


    // ── setImmunityPeriod tests ──

    function test_setImmunityPeriod() public {
        vm.prank(guardian);
        awpRegistry.setImmunityPeriod(60 days);
        assertEq(awpRegistry.immunityPeriod(), 60 days);
    }

    // ── nextWorknetId tests ──

    function test_nextWorknetId() public {
        uint256 before = awpRegistry.nextWorknetId();
        _registerWorknet();
        assertEq(awpRegistry.nextWorknetId(), before + 1);
    }
}
