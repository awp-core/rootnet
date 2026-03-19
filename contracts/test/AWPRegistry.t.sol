// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AlphaTokenFactory} from "../src/token/AlphaTokenFactory.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {AccessManager} from "../src/core/AccessManager.sol";
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
    AccessManager access;
    StakingVault vault;
    StakeNFT stakeNFT;
    SubnetNFT nft;
    MockLPManager lp;
    AWPRegistry rootNet;
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

        // Deploy AWPRegistry (no epochDuration param)
        rootNet = new AWPRegistry(deployer, address(treasury), guardian);

        // Deploy sub-contracts
        factory = new AlphaTokenFactory(deployer, 0);
        nft = new SubnetNFT("AWP Subnet", "AWPSUB", address(rootNet));
        access = new AccessManager(address(rootNet));
        lp = new MockLPManager(address(rootNet), address(awp));

        // Deploy AWPEmission (UUPS proxy) — now has its own epoch timing
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
        factory.setAddresses(address(rootNet));

        // Deploy StakeNFT and StakingVault (circular dependency resolution)
        uint64 nonce = vm.getNonce(deployer);
        address predictedVault = vm.computeCreateAddress(deployer, nonce);
        address predictedStakeNFT = vm.computeCreateAddress(deployer, nonce + 1);

        vault = new StakingVault(address(rootNet));
        stakeNFT = new StakeNFT(address(awp), address(vault), address(rootNet));

        // Initialize registry
        rootNet.initializeRegistry(
            address(awp),
            address(nft),
            address(factory),
            address(emission),
            address(lp),
            address(access),
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
        rootNet.initializeRegistry(
            address(awp), address(nft), address(factory), address(emission),
            address(lp), address(access), address(vault), address(stakeNFT), address(0), ""
        );
    }

    // ── User registration ──

    function test_register() public {
        vm.prank(user1);
        rootNet.register();
        assertTrue(access.isRegistered(user1));
    }

    function test_registerAndStake() public {
        vm.startPrank(user1);
        // Approve StakeNFT for deposit
        awp.approve(address(stakeNFT), 1000 * 1e18);
        // Register + deposit via StakeNFT, without allocating (lockDuration in seconds)
        rootNet.registerAndStake(1000 * 1e18, 52 weeks, address(0), 0, 0);
        vm.stopPrank();

        assertTrue(access.isRegistered(user1));
        assertEq(stakeNFT.getUserTotalStaked(user1), 1000 * 1e18);
    }

    // ── Agent binding ──

    function test_bind() public {
        // bind auto-registers the principal
        vm.prank(agent1);
        rootNet.bind(user1);

        assertTrue(access.isAgent(user1, agent1));
        assertTrue(access.isRegistered(user1));
    }

    function test_bind_rebind() public {
        vm.prank(user1);
        rootNet.register();

        vm.prank(agent1);
        rootNet.bind(user1);

        // rebind to user2
        vm.prank(agent1);
        rootNet.bind(user2);

        assertFalse(access.isAgent(user1, agent1));
        assertTrue(access.isAgent(user2, agent1));
    }

    // ── Agent management ──

    function test_removeAgent() public {
        vm.prank(user1);
        rootNet.register();
        vm.prank(agent1);
        rootNet.bind(user1);

        vm.prank(user1);
        rootNet.removeAgent(agent1);

        assertFalse(access.isAgent(user1, agent1));
    }

    function test_setDelegation() public {
        vm.prank(user1);
        rootNet.register();
        vm.prank(agent1);
        rootNet.bind(user1);

        vm.prank(user1);
        rootNet.setDelegation(agent1, true);

        assertTrue(access.isManagerAgent(agent1));
    }

    function test_unbind() public {
        vm.prank(user1);
        rootNet.register();
        vm.prank(agent1);
        rootNet.bind(user1);

        assertTrue(access.isAgent(user1, agent1));

        vm.prank(agent1);
        rootNet.unbind();

        assertFalse(access.isRegisteredAgent(agent1));
        assertFalse(access.isKnownAddress(agent1));
    }

    function test_registerWithOptions_recipientAndStake() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        rootNet.register(user2, 1000 * 1e18, 52 weeks);
        vm.stopPrank();

        assertTrue(access.isRegistered(user1));
        assertEq(access.getRewardRecipient(user1), user2);
        assertEq(stakeNFT.getUserTotalStaked(user1), 1000 * 1e18);
    }

    function test_registerWithOptions_emptyParams() public {
        vm.prank(user1);
        rootNet.register(address(0), 0, 0);

        assertTrue(access.isRegistered(user1));
        assertEq(access.getRewardRecipient(user1), user1); // default = self
        assertEq(stakeNFT.getUserTotalStaked(user1), 0);
    }

    function test_registerWithOptions_idempotent() public {
        // Calling register(options) on an already-registered user should not revert
        vm.prank(user1);
        rootNet.register();

        vm.prank(user1);
        rootNet.register(user2, 0, 0); // just update recipient

        assertEq(access.getRewardRecipient(user1), user2);
    }

    function test_setRewardRecipient() public {
        vm.prank(user1);
        rootNet.register();

        vm.prank(user1);
        rootNet.setRewardRecipient(user2);

        assertEq(access.getRewardRecipient(user1), user2);
    }

    // ── Staking (via StakeNFT) ──

    function test_allocateAndDeallocate() public {
        // Setup: register, deposit via StakeNFT, register agent, register subnet
        vm.startPrank(user1);
        rootNet.register();
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        vm.prank(agent1);
        rootNet.bind(user1);

        _registerSubnet();

        vm.prank(user1);
        rootNet.activateSubnet(1);

        // Allocate
        vm.prank(user1);
        rootNet.allocate(agent1, 1, 500 * 1e18);
        assertEq(vault.getAgentStake(user1, agent1, 1), 500 * 1e18);

        // Deallocate
        vm.prank(user1);
        rootNet.deallocate(agent1, 1, 200 * 1e18);
        assertEq(vault.getAgentStake(user1, agent1, 1), 300 * 1e18);
    }

    function test_reallocate_immediate() public {
        vm.startPrank(user1);
        rootNet.register();
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        vm.prank(agent1);
        rootNet.bind(user1);

        _registerSubnet(); // subnetId=1

        vm.prank(user1);
        rootNet.activateSubnet(1);

        address agent2 = address(10);
        vm.prank(agent2);
        rootNet.bind(user1);

        // Allocate to agent1/subnet1
        vm.prank(user1);
        rootNet.allocate(agent1, 1, 500 * 1e18);

        // Reallocate to agent2/subnet1 — takes effect immediately
        vm.prank(user1);
        rootNet.reallocate(agent1, 1, agent2, 1, 200 * 1e18);

        // Immediate effect, no epoch advance needed
        assertEq(vault.getAgentStake(user1, agent1, 1), 300 * 1e18);
        assertEq(vault.getAgentStake(user1, agent2, 1), 200 * 1e18);
    }

    // ── Subnet registration ──

    function test_registerSubnet() public {
        uint256 subnetId = _registerSubnet();
        assertEq(subnetId, 1);

        IAWPRegistry.SubnetInfo memory info = rootNet.getSubnet(1);
        assertEq(rootNet.getSubnetFull(subnetId).subnetManager, subnetManager);
        assertTrue(info.status == IAWPRegistry.SubnetStatus.Pending);
    }

    function test_registerSubnetInvalidParams() public {
        vm.startPrank(user1);
        awp.approve(address(rootNet), 2_000_000 * 1e18);

        // Empty name
        vm.expectRevert(AWPRegistry.InvalidSubnetParams.selector);
        rootNet.registerSubnet(
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
        rootNet.registerSubnet(
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
        rootNet.activateSubnet(subnetId);

        assertTrue(rootNet.isSubnetActive(subnetId));
        assertEq(rootNet.getActiveSubnetCount(), 1);
    }

    function test_pauseAndResumeSubnet() public {
        uint256 subnetId = _registerSubnet();
        vm.startPrank(user1);
        rootNet.activateSubnet(subnetId);
        rootNet.pauseSubnet(subnetId);
        assertFalse(rootNet.isSubnetActive(subnetId));

        rootNet.resumeSubnet(subnetId);
        assertTrue(rootNet.isSubnetActive(subnetId));
        vm.stopPrank();
    }

    function test_banAndUnbanSubnet() public {
        uint256 subnetId = _registerSubnet();
        vm.prank(user1);
        rootNet.activateSubnet(subnetId);

        vm.prank(address(treasury));
        rootNet.banSubnet(subnetId);

        IAWPRegistry.SubnetInfo memory info = rootNet.getSubnet(subnetId);
        assertTrue(info.status == IAWPRegistry.SubnetStatus.Banned);

        vm.prank(address(treasury));
        rootNet.unbanSubnet(subnetId);
        info = rootNet.getSubnet(subnetId);
        assertTrue(info.status == IAWPRegistry.SubnetStatus.Active);
    }

    function test_deregisterSubnet() public {
        uint256 subnetId = _registerSubnet();
        vm.prank(user1);
        rootNet.activateSubnet(subnetId);

        vm.prank(address(treasury));
        vm.expectRevert(AWPRegistry.ImmunityNotExpired.selector);
        rootNet.deregisterSubnet(subnetId);

        vm.warp(block.timestamp + 31 days);
        vm.prank(address(treasury));
        rootNet.deregisterSubnet(subnetId);

        assertEq(rootNet.getActiveSubnetCount(), 0);
    }

    // ── Permission tests ──

    function test_onlyTimelockFunctions() public {
        vm.expectRevert(AWPRegistry.NotTimelock.selector);
        rootNet.setInitialAlphaPrice(1e15);

        vm.expectRevert(AWPRegistry.NotTimelock.selector);
        rootNet.setGuardian(address(0));

        vm.expectRevert(AWPRegistry.NotTimelock.selector);
        rootNet.unpause();
    }

    function test_onlyGuardianPause() public {
        vm.prank(guardian);
        rootNet.pause();
        assertTrue(rootNet.paused());

        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        rootNet.pause();
    }

    // ── Queries ──

    function test_getAgentInfo() public {
        vm.prank(user1);
        rootNet.register();
        vm.prank(agent1);
        rootNet.bind(user1);

        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        _registerSubnet();

        vm.prank(user1);
        rootNet.activateSubnet(1);

        vm.prank(user1);
        rootNet.allocate(agent1, 1, 500 * 1e18);

        AWPRegistry.AgentInfo memory info = rootNet.getAgentInfo(agent1, 1);
        assertEq(info.owner, user1);
        assertTrue(info.isValid);
        assertEq(info.stake, 500 * 1e18);
        assertEq(info.rewardRecipient, user1);
    }

    function test_getRegistry() public view {
        (address a, address b, address c, address d, address e, address f, address g, address h,,) =
            rootNet.getRegistry();
        assertEq(a, address(awp));
        assertEq(b, address(nft));
        assertEq(c, address(factory));
        assertEq(d, address(emission));
        assertEq(e, address(lp));
        assertEq(f, address(access));
        assertEq(g, address(vault));
        assertEq(h, address(stakeNFT));
    }

    // ── Helper functions ──

    function _registerSubnet() internal returns (uint256) {
        uint256 lpCost = 100_000_000 * 1e18 * 1e16 / 1e18;
        vm.startPrank(user1);
        awp.approve(address(rootNet), lpCost);
        uint256 subnetId = rootNet.registerSubnet(
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

    // ── minStake enforcement ──

    function test_allocate_insufficientMinStake() public {
        // Register subnet with minStake = 1000 * 1e18
        uint256 lpCost = 100_000_000 * 1e18 * 1e16 / 1e18;
        vm.startPrank(user1);
        awp.approve(address(rootNet), lpCost);
        uint256 subnetId = rootNet.registerSubnet(
            IAWPRegistry.SubnetParams({
                name: "MinStakeSubnet",
                symbol: "MSUB",
                subnetManager: subnetManager,
                salt: bytes32(0),
                minStake: 1000 * 1e18,
                skillsURI: ""
            })
        );
        rootNet.activateSubnet(subnetId);
        vm.stopPrank();

        // Setup: register user, bind agent, deposit AWP
        vm.startPrank(user1);
        rootNet.register();
        awp.approve(address(stakeNFT), 5000 * 1e18);
        stakeNFT.deposit(5000 * 1e18, 52 weeks);
        vm.stopPrank();

        vm.prank(agent1);
        rootNet.bind(user1);

        // Allocate below minStake — should revert
        vm.prank(user1);
        vm.expectRevert(AWPRegistry.InsufficientMinStake.selector);
        rootNet.allocate(agent1, subnetId, 500 * 1e18);

        // Allocate at minStake — should succeed
        vm.prank(user1);
        rootNet.allocate(agent1, subnetId, 1000 * 1e18);
        assertEq(vault.getAgentStake(user1, agent1, subnetId), 1000 * 1e18);
    }

    // ── setImmunityPeriod tests ──

    function test_setImmunityPeriod() public {
        vm.prank(address(treasury));
        rootNet.setImmunityPeriod(60 days);
        assertEq(rootNet.immunityPeriod(), 60 days);
    }

    // ── nextSubnetId tests ──

    function test_nextSubnetId() public {
        uint256 before = rootNet.nextSubnetId();
        _registerSubnet();
        assertEq(rootNet.nextSubnetId(), before + 1);
    }
}
