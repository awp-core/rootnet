// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AlphaTokenFactory} from "../src/token/AlphaTokenFactory.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {SubnetNFT} from "../src/core/SubnetNFT.sol";
import {MockLPManager} from "./helpers/MockLPManager.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {IAWPRegistry} from "../src/interfaces/IAWPRegistry.sol";
import {Treasury} from "../src/governance/Treasury.sol";

contract AWPRegistryTest is Test {
    AWPToken awp;
    AlphaTokenFactory factory;
    AWPEmission emission;
    StakingVault vault;
    StakeNFT stakeNFT;
    SubnetNFT nft;
    MockLPManager lp;
    AWPRegistry awpRegistry;
    Treasury treasury;

    address deployer = address(1);
    address guardian = address(2);
    address user1 = address(3);
    address agent1 = address(4);
    address subnetManager = address(5);
    address user2 = address(6);

    uint256 constant INITIAL_DAILY_EMISSION = 15_800_000 * 1e18;
    uint256 constant EPOCH_DURATION = 1 days;

    function setUp() public {
        vm.startPrank(deployer);

        // Deploy tokens
        awp = new AWPToken("AWP Token", "AWP", deployer);

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
        nft = new SubnetNFT("AWP Subnet", "AWPSUB", address(awpRegistry));
        lp = new MockLPManager(address(awpRegistry), address(awp));

        // Deploy AWPEmission (UUPS proxy)
        AWPEmission emissionImpl = new AWPEmission();
        bytes memory emissionInitData = abi.encodeCall(
            AWPEmission.initialize,
            (address(awp), address(treasury), INITIAL_DAILY_EMISSION, block.timestamp, EPOCH_DURATION)
        );
        ERC1967Proxy emissionProxy = new ERC1967Proxy(address(emissionImpl), emissionInitData);
        emission = AWPEmission(address(emissionProxy));

        // Set AWP minter
        awp.addMinter(address(emission));
        awp.renounceAdmin();

        // Configure factory
        factory.setAddresses(address(awpRegistry));

        // Deploy StakeNFT and StakingVault
        vault = new StakingVault(address(awpRegistry));
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
            address(0), // no default SubnetManager impl in unit tests
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
        awpRegistry.register();
        assertTrue(awpRegistry.isRegistered(user1));
    }

    function test_register_alreadyRegistered_reverts() public {
        vm.prank(user1);
        awpRegistry.register();
        vm.prank(user1);
        vm.expectRevert(AWPRegistry.AlreadyRegistered.selector);
        awpRegistry.register();
    }

    // ── Binding ──

    function test_bind() public {
        vm.prank(agent1);
        awpRegistry.bind(user1);
        assertEq(awpRegistry.boundTo(agent1), user1);
    }

    function test_bind_selfBind_reverts() public {
        vm.prank(user1);
        vm.expectRevert(AWPRegistry.InvalidAddress.selector);
        awpRegistry.bind(user1);
    }

    function test_bind_zeroAddress_reverts() public {
        vm.prank(agent1);
        vm.expectRevert(AWPRegistry.InvalidAddress.selector);
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
        // Setup: deposit via StakeNFT, register subnet
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        uint256 sid = _registerSubnet();

        vm.prank(user1);
        awpRegistry.activateSubnet(sid);

        // Allocate (staker = user1, agent = agent1)
        vm.prank(user1);
        awpRegistry.allocate(user1, agent1, sid, 500 * 1e18);
        assertEq(vault.getAgentStake(user1, agent1, sid), 500 * 1e18);

        // Deallocate
        vm.prank(user1);
        awpRegistry.deallocate(user1, agent1, sid, 200 * 1e18);
        assertEq(vault.getAgentStake(user1, agent1, sid), 300 * 1e18);
    }

    function test_allocate_delegate() public {
        // Setup: deposit
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        awpRegistry.grantDelegate(user2);
        vm.stopPrank();

        uint256 sid = _registerSubnet();

        vm.prank(user1);
        awpRegistry.activateSubnet(sid);

        // user2 allocates on behalf of user1
        vm.prank(user2);
        awpRegistry.allocate(user1, agent1, sid, 500 * 1e18);
        assertEq(vault.getAgentStake(user1, agent1, sid), 500 * 1e18);
    }

    function test_allocate_notAuthorized_reverts() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        uint256 sid = _registerSubnet();

        vm.prank(user1);
        awpRegistry.activateSubnet(sid);

        // user2 has no delegation from user1
        vm.prank(user2);
        vm.expectRevert(AWPRegistry.NotAuthorized.selector);
        awpRegistry.allocate(user1, agent1, sid, 500 * 1e18);
    }

    function test_reallocate_immediate() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        uint256 sid = _registerSubnet();

        vm.prank(user1);
        awpRegistry.activateSubnet(sid);

        address agent2 = address(10);

        // Allocate to agent1/subnet
        vm.prank(user1);
        awpRegistry.allocate(user1, agent1, sid, 500 * 1e18);

        // Reallocate to agent2/subnet — takes effect immediately
        vm.prank(user1);
        awpRegistry.reallocate(user1, agent1, sid, agent2, sid, 200 * 1e18);

        assertEq(vault.getAgentStake(user1, agent1, sid), 300 * 1e18);
        assertEq(vault.getAgentStake(user1, agent2, sid), 200 * 1e18);
    }

    // ── Subnet registration ──

    function test_registerSubnet() public {
        uint256 subnetId = _registerSubnet();
        assertEq(subnetId & ((1 << 64) - 1), 1);
        assertEq(subnetId >> 64, block.chainid);

        IAWPRegistry.SubnetInfo memory info = awpRegistry.getSubnet(subnetId);
        assertEq(awpRegistry.getSubnetFull(subnetId).subnetManager, subnetManager);
        assertTrue(info.status == IAWPRegistry.SubnetStatus.Pending);
    }

    function test_registerSubnetInvalidParams() public {
        vm.startPrank(user1);
        awp.approve(address(awpRegistry), 2_000_000 * 1e18);

        // Empty name
        vm.expectRevert(AWPRegistry.InvalidSubnetParams.selector);
        awpRegistry.registerSubnet(
            IAWPRegistry.SubnetParams({
                name: "",
                symbol: "TEST",
                subnetManager: subnetManager,
                salt: bytes32(0),
                minStake: 0,
                skillsURI: ""
            })
        );

        // Empty subnet contract address
        vm.expectRevert(AWPRegistry.SubnetManagerRequired.selector);
        awpRegistry.registerSubnet(
            IAWPRegistry.SubnetParams({
                name: "Test",
                symbol: "TEST",
                subnetManager: address(0),
                salt: bytes32(0),
                minStake: 0,
                skillsURI: ""
            })
        );
        vm.stopPrank();
    }

    // ── Subnet lifecycle ──

    function test_activateSubnet() public {
        uint256 subnetId = _registerSubnet();

        vm.prank(user1);
        awpRegistry.activateSubnet(subnetId);

        assertTrue(awpRegistry.isSubnetActive(subnetId));
        assertEq(awpRegistry.getActiveSubnetCount(), 1);
    }

    function test_pauseAndResumeSubnet() public {
        uint256 subnetId = _registerSubnet();
        vm.startPrank(user1);
        awpRegistry.activateSubnet(subnetId);
        awpRegistry.pauseSubnet(subnetId);
        assertFalse(awpRegistry.isSubnetActive(subnetId));

        awpRegistry.resumeSubnet(subnetId);
        assertTrue(awpRegistry.isSubnetActive(subnetId));
        vm.stopPrank();
    }

    function test_banAndUnbanSubnet() public {
        uint256 subnetId = _registerSubnet();
        vm.prank(user1);
        awpRegistry.activateSubnet(subnetId);

        vm.prank(address(treasury));
        awpRegistry.banSubnet(subnetId);

        IAWPRegistry.SubnetInfo memory info = awpRegistry.getSubnet(subnetId);
        assertTrue(info.status == IAWPRegistry.SubnetStatus.Banned);

        vm.prank(address(treasury));
        awpRegistry.unbanSubnet(subnetId);
        info = awpRegistry.getSubnet(subnetId);
        assertTrue(info.status == IAWPRegistry.SubnetStatus.Active);
    }

    function test_deregisterSubnet() public {
        uint256 subnetId = _registerSubnet();
        vm.prank(user1);
        awpRegistry.activateSubnet(subnetId);

        // Cannot deregister Active subnet (must ban first)
        vm.prank(address(treasury));
        vm.expectRevert(AWPRegistry.InvalidSubnetStatus.selector);
        awpRegistry.deregisterSubnet(subnetId);

        // Ban the subnet first
        vm.prank(address(treasury));
        awpRegistry.banSubnet(subnetId);

        // Cannot deregister during immunity period
        vm.prank(address(treasury));
        vm.expectRevert(AWPRegistry.ImmunityNotExpired.selector);
        awpRegistry.deregisterSubnet(subnetId);

        // After immunity period, deregister succeeds
        vm.warp(block.timestamp + 31 days);
        vm.prank(address(treasury));
        awpRegistry.deregisterSubnet(subnetId);

        assertEq(awpRegistry.getActiveSubnetCount(), 0);
    }

    // ── UUPS Upgrade tests ──

    function test_upgradeViaTimelock() public {
        AWPRegistry newImpl = new AWPRegistry();
        vm.prank(address(treasury));
        awpRegistry.upgradeToAndCall(address(newImpl), "");
    }

    function test_upgrade_revertsForNonTimelock() public {
        AWPRegistry newImpl = new AWPRegistry();
        vm.expectRevert(AWPRegistry.NotTimelock.selector);
        awpRegistry.upgradeToAndCall(address(newImpl), "");
    }

    function test_cannotInitializeImpl() public {
        AWPRegistry impl = new AWPRegistry();
        vm.expectRevert();
        impl.initialize(deployer, address(treasury), guardian);
    }

    // ── SubnetId encoding tests ──

    function test_subnetIdEncodesChainId() public {
        uint256 subnetId = _registerSubnet();
        uint256 expectedChainId = block.chainid;
        assertEq(subnetId >> 64, expectedChainId);
        assertEq(subnetId & ((1 << 64) - 1), 1);
    }

    function test_subnetIdIncrementsLocalCounter() public {
        uint256 id1 = _registerSubnet();
        vm.startPrank(user1);
        awp.approve(address(awpRegistry), 1_000_000 * 1e18);
        uint256 id2 = awpRegistry.registerSubnet(
            IAWPRegistry.SubnetParams("Sub2", "S2", subnetManager, bytes32(0), 0, "")
        );
        vm.stopPrank();
        assertEq(id1 >> 64, id2 >> 64);
        assertEq((id1 & ((1 << 64) - 1)) + 1, id2 & ((1 << 64) - 1));
    }

    function test_extractChainId() public view {
        uint256 subnetId = (uint256(8453) << 64) | 42;
        assertEq(awpRegistry.extractChainId(subnetId), 8453);
        assertEq(awpRegistry.extractLocalId(subnetId), 42);
    }

    // ── Permission tests ──

    function test_onlyTimelockFunctions() public {
        vm.expectRevert(AWPRegistry.NotTimelock.selector);
        awpRegistry.setInitialAlphaPrice(1e15);

        vm.expectRevert(AWPRegistry.NotTimelock.selector);
        awpRegistry.setGuardian(address(0));

        vm.expectRevert(AWPRegistry.NotTimelock.selector);
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
        awpRegistry.register();
        vm.prank(agent1);
        awpRegistry.bind(user1);

        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        uint256 sid = _registerSubnet();

        vm.prank(user1);
        awpRegistry.activateSubnet(sid);

        vm.prank(user1);
        awpRegistry.allocate(user1, agent1, sid, 500 * 1e18);

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

    function _registerSubnet() internal returns (uint256) {
        uint256 lpCost = 100_000_000 * 1e18 * 1e16 / 1e18;
        vm.startPrank(user1);
        awp.approve(address(awpRegistry), lpCost);
        uint256 subnetId = awpRegistry.registerSubnet(
            IAWPRegistry.SubnetParams({
                name: "TestSubnet",
                symbol: "TSUB",
                subnetManager: subnetManager,
                salt: bytes32(0),
                minStake: 0,
                skillsURI: ""
            })
        );
        vm.stopPrank();
        return subnetId;
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

    // ── Bind chain too long ──

    function test_bind_chainTooLong() public {
        // Build a chain of 100 addresses, then try to bind one more
        address[] memory addrs = new address[](101);
        for (uint256 i = 0; i < 101; i++) {
            addrs[i] = address(uint160(0x1000 + i));
        }
        // Chain: addrs[0] is root, addrs[1] → addrs[0], addrs[2] → addrs[1], ...
        for (uint256 i = 1; i < 101; i++) {
            vm.prank(addrs[i]);
            awpRegistry.bind(addrs[i-1]);
        }
        // addrs[100] is bound to addrs[99], chain depth is 100
        // Now try to bind a new address to addrs[100] — should revert ChainTooLong
        address newAddr = makeAddr("tooDeep");
        vm.prank(newAddr);
        vm.expectRevert(AWPRegistry.ChainTooLong.selector);
        awpRegistry.bind(addrs[100]);
    }

    // ── minStake enforcement ──


    // ── setImmunityPeriod tests ──

    function test_setImmunityPeriod() public {
        vm.prank(address(treasury));
        awpRegistry.setImmunityPeriod(60 days);
        assertEq(awpRegistry.immunityPeriod(), 60 days);
    }

    // ── nextSubnetId tests ──

    function test_nextSubnetId() public {
        uint256 before = awpRegistry.nextSubnetId();
        _registerSubnet();
        assertEq(awpRegistry.nextSubnetId(), before + 1);
    }
}
