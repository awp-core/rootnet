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
import {SubnetNFT} from "../src/core/SubnetNFT.sol";
import {MockLPManager} from "./helpers/MockLPManager.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {IAWPRegistry} from "../src/interfaces/IAWPRegistry.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {TimelockController} from "@openzeppelin/contracts/governance/TimelockController.sol";
import {IGovernor} from "@openzeppelin/contracts/governance/IGovernor.sol";
import {EmissionSigningHelper} from "./helpers/EmissionSigningHelper.sol";

/// @title E2E — End-to-end tests based on the architecture document (Account System V2)
contract E2ETest is EmissionSigningHelper {
    AWPToken awp;
    AlphaTokenFactory factory;
    AWPEmission emission;
    StakingVault vault;
    StakeNFT stakeNFT;
    SubnetNFT nft;
    MockLPManager lp;
    AWPRegistry rootNet;
    Treasury treasury;
    AWPDAO dao;

    address deployer = address(0xD);
    address guardian = address(0xE);
    address airdropAddr = address(0xF4);

    // Mock users
    address alice = address(0x1001);
    address bob = address(0x1002);
    address charlie = address(0x1003);

    // Mock Agents
    address agentA = address(0x2001);
    address agentB = address(0x2002);
    address agentC = address(0x2003);

    // Mock subnet contracts
    address subnetC1 = address(0x3001);
    address subnetC2 = address(0x3002);
    address subnetC3 = address(0x3003);

    uint256 constant INITIAL_DAILY = 15_800_000 * 1e18;
    uint256 constant EPOCH = 1 days;
    uint256 constant LP_COST = 1_000_000 * 1e18; // 100M Alpha * 0.01 AWP

    // Oracle private keys
    uint256 constant ORACLE_PK1 = 0xA1;
    uint256 constant ORACLE_PK2 = 0xA2;
    uint256 constant ORACLE_PK3 = 0xA3;

    function setUp() public {
        _deploy();
    }

    function _deploy() internal {
        vm.startPrank(deployer);
        awp = new AWPToken("AWP Token", "AWP", deployer);
        factory = new AlphaTokenFactory(deployer, 0);

        address[] memory p = new address[](0);
        address[] memory e = new address[](1);
        e[0] = address(0);
        treasury = new Treasury(1, p, e, deployer);

        rootNet = new AWPRegistry(deployer, address(treasury), guardian);
        nft = new SubnetNFT("AWP Subnet", "AWPSUB", address(rootNet));
        lp = new MockLPManager(address(rootNet), address(awp));

        // Deploy AWPEmission (UUPS proxy)
        AWPEmission emissionImpl = new AWPEmission();
        bytes memory emissionInitData = abi.encodeCall(
            AWPEmission.initialize,
            (address(awp), address(treasury), INITIAL_DAILY, block.timestamp, EPOCH)
        );
        ERC1967Proxy emissionProxy = new ERC1967Proxy(address(emissionImpl), emissionInitData);
        emission = AWPEmission(address(emissionProxy));

        awp.addMinter(address(emission));
        awp.renounceAdmin();
        factory.setAddresses(address(rootNet));

        // Deploy StakingVault + StakeNFT
        vault = new StakingVault(address(rootNet));
        stakeNFT = new StakeNFT(address(awp), address(vault), address(rootNet));

        // Deploy AWPDAO
        dao = new AWPDAO(
            address(stakeNFT),
            address(awp),
            TimelockController(payable(address(treasury))),
            1,    // votingDelay
            100,  // votingPeriod
            4     // quorum 4%
        );
        treasury.grantRole(treasury.PROPOSER_ROLE(), address(dao));
        treasury.grantRole(treasury.CANCELLER_ROLE(), address(dao));
        treasury.renounceRole(treasury.DEFAULT_ADMIN_ROLE(), deployer);

        // Initialize registry (no accessManager)
        rootNet.initializeRegistry(
            address(awp), address(nft), address(factory), address(emission),
            address(lp), address(vault), address(stakeNFT), address(0), ""
        );

        vm.stopPrank();

        // Configure oracles
        address[] memory oracleList = new address[](3);
        oracleList[0] = vm.addr(ORACLE_PK1);
        oracleList[1] = vm.addr(ORACLE_PK2);
        oracleList[2] = vm.addr(ORACLE_PK3);
        vm.prank(address(treasury));
        emission.setOracleConfig(oracleList, 2);

        // Distribute tokens from deployer
        vm.startPrank(deployer);
        awp.transfer(address(treasury), 50_000_000 * 1e18);
        awp.transfer(airdropAddr, 150_000_000 * 1e18);
        vm.stopPrank();

        // Distribute AWP to test users
        vm.startPrank(airdropAddr);
        awp.transfer(alice, 50_000_000 * 1e18);
        awp.transfer(bob, 50_000_000 * 1e18);
        awp.transfer(charlie, 50_000_000 * 1e18);
        vm.stopPrank();
    }

    // ── Helper functions ──

    function _registerUser(address user) internal {
        vm.prank(user);
        rootNet.register();
    }

    function _bindAgent(address agent, address target) internal {
        vm.prank(agent);
        rootNet.bind(target);
    }

    function _registerSubnet(address owner, address sc) internal returns (uint256) {
        vm.startPrank(owner);
        awp.approve(address(rootNet), LP_COST);
        uint256 id = rootNet.registerSubnet(
            IAWPRegistry.SubnetParams("Subnet", "SUB", sc, bytes32(0), 0, "")
        );
        vm.stopPrank();
        return id;
    }

    /// @dev Deposit AWP via StakeNFT and allocate via RootNet (explicit staker)
    function _depositAndAllocate(address staker, address agent, uint256 subnetId, uint256 deposit, uint256 alloc)
        internal
    {
        vm.startPrank(staker);
        awp.approve(address(stakeNFT), deposit);
        stakeNFT.deposit(deposit, 52 weeks);
        rootNet.allocate(staker, agent, subnetId, alloc);
        vm.stopPrank();
    }

    function _settleEpoch() internal {
        vm.warp(block.timestamp + EPOCH + 1);
        emission.settleEpoch(200);
    }

    function _sortByAddress(address[] memory addrs, uint96[] memory ws) internal pure {
        uint256 n = addrs.length;
        for (uint256 i = 0; i < n; i++) {
            for (uint256 j = i + 1; j < n; j++) {
                if (uint160(addrs[i]) > uint160(addrs[j])) {
                    (addrs[i], addrs[j]) = (addrs[j], addrs[i]);
                    (ws[i], ws[j]) = (ws[j], ws[i]);
                }
            }
        }
    }

    function _submitWeights(address[] memory recipients, uint96[] memory weights) internal {
        _sortByAddress(recipients, weights);
        uint256 nonce = emission.allocationNonce();
        uint256 effectiveEpoch = emission.settledEpoch() + 1;
        bytes[] memory sigs = new bytes[](2);
        sigs[0] = _signAllocations(ORACLE_PK1, recipients, weights, nonce, effectiveEpoch, address(emission));
        sigs[1] = _signAllocations(ORACLE_PK2, recipients, weights, nonce, effectiveEpoch, address(emission));
        emission.submitAllocations(recipients, weights, sigs, effectiveEpoch);
    }

    function _submitWeight(address _recipient, uint96 weight) internal {
        address[] memory addrs = new address[](1);
        addrs[0] = _recipient;
        uint96[] memory ws = new uint96[](1);
        ws[0] = weight;
        _submitWeights(addrs, ws);
    }

    // ════════════════════════════════════════════
    //  E2E 1: Subnet NFT transfer -> ownership change
    // ════════════════════════════════════════════

    function test_e2e_subnetNFTTransfer() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);

        vm.prank(alice);
        rootNet.activateSubnet(sid);

        vm.prank(alice);
        nft.transferFrom(alice, bob, sid);
        assertEq(nft.ownerOf(sid), bob);

        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotOwner.selector);
        rootNet.pauseSubnet(sid);

        vm.prank(bob);
        rootNet.pauseSubnet(sid);
        assertFalse(rootNet.isSubnetActive(sid));

        vm.prank(bob);
        rootNet.resumeSubnet(sid);
        assertTrue(rootNet.isSubnetActive(sid));
    }

    // ════════════════════════════════════════════
    //  E2E 2: Reallocate is immediate
    // ════════════════════════════════════════════

    function test_e2e_reallocateImmediate() public {
        _registerUser(alice);
        _bindAgent(agentA, alice);
        _bindAgent(agentB, alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        _depositAndAllocate(alice, agentA, sid, 10_000 * 1e18, 8_000 * 1e18);

        // Reallocate 5000 from agentA -> agentB (immediate)
        vm.prank(alice);
        rootNet.reallocate(alice, agentA, sid, agentB, sid, 5_000 * 1e18);

        assertEq(vault.getAgentStake(alice, agentA, sid), 3_000 * 1e18);
        assertEq(vault.getAgentStake(alice, agentB, sid), 5_000 * 1e18);
        assertEq(vault.userTotalAllocated(alice), 8_000 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 3: Full subnet lifecycle
    // ════════════════════════════════════════════

    function test_e2e_subnetFullLifecycle() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);

        vm.prank(alice);
        rootNet.activateSubnet(sid);
        assertEq(uint256(rootNet.getSubnet(sid).status), uint256(IAWPRegistry.SubnetStatus.Active));

        vm.prank(alice);
        rootNet.pauseSubnet(sid);
        assertEq(uint256(rootNet.getSubnet(sid).status), uint256(IAWPRegistry.SubnetStatus.Paused));

        vm.prank(alice);
        rootNet.resumeSubnet(sid);
        assertEq(uint256(rootNet.getSubnet(sid).status), uint256(IAWPRegistry.SubnetStatus.Active));

        vm.prank(address(treasury));
        rootNet.banSubnet(sid);
        assertEq(uint256(rootNet.getSubnet(sid).status), uint256(IAWPRegistry.SubnetStatus.Banned));

        AlphaToken alpha = AlphaToken(rootNet.getSubnetFull(sid).alphaToken);
        assertTrue(alpha.minterPaused(subnetC1));

        vm.prank(address(treasury));
        rootNet.unbanSubnet(sid);
        assertEq(uint256(rootNet.getSubnet(sid).status), uint256(IAWPRegistry.SubnetStatus.Active));
        assertFalse(alpha.minterPaused(subnetC1));

        vm.warp(block.timestamp + 31 days);
        vm.prank(address(treasury));
        rootNet.deregisterSubnet(sid);

        vm.expectRevert();
        nft.ownerOf(sid);
        vm.expectRevert();
        rootNet.getSubnetFull(sid);
    }

    // ════════════════════════════════════════════
    //  E2E 4: DAO governance with NFT-based voting
    // ════════════════════════════════════════════

    function test_e2e_daoGovernanceWeight() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        uint256 stakeAmount = 40_000_000 * 1e18;
        vm.startPrank(alice);
        awp.approve(address(stakeNFT), stakeAmount);
        uint256 tokenId = stakeNFT.deposit(stakeAmount, 52 weeks);
        vm.stopPrank();

        _submitWeight(subnetC1, uint96(100));
        _settleEpoch();
        _settleEpoch();

        vm.roll(block.number + 1);

        address[] memory targets = new address[](1);
        targets[0] = address(emission);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(emission.emergencySetWeight, (uint256(1), uint256(0), subnetC1, uint96(500)));
        bytes32 descHash = keccak256("Set weight for subnet 1");

        uint256[] memory propTokenIds = new uint256[](1);
        propTokenIds[0] = tokenId;
        vm.prank(alice);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, "Set weight for subnet 1", propTokenIds);

        vm.roll(block.number + 2);

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = tokenId;
        bytes memory params = abi.encode(tokenIds);
        vm.prank(alice);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", params);

        vm.roll(block.number + 101);

        dao.queue(targets, values, calldatas, descHash);
        vm.warp(block.timestamp + 2);
        dao.execute(targets, values, calldatas, descHash);

        assertEq(emission.getTotalWeight(), 500);
    }

    // ════════════════════════════════════════════
    //  E2E 5: Multi-epoch emission decay
    // ════════════════════════════════════════════

    function test_e2e_emissionDecayMultiEpoch() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        _submitWeight(subnetC1, uint96(100));

        uint256 prevEmission = INITIAL_DAILY;

        for (uint256 i = 0; i < 10; i++) {
            if (i > 0) {
                _submitWeight(subnetC1, uint96(100));
            }
            _settleEpoch();
            uint256 currentEmission = emission.currentDailyEmission();
            if (i > 0) {
                assertEq(currentEmission, prevEmission * 996844 / 1000000);
            }
            prevEmission = currentEmission;
        }

        assertEq(emission.currentEpoch(), 10);
    }

    // ════════════════════════════════════════════
    //  E2E 6: Multi-user multi-subnet concurrency
    // ════════════════════════════════════════════

    function test_e2e_multiUserMultiSubnet() public {
        _registerUser(alice);
        _registerUser(bob);
        _registerUser(charlie);

        uint256 sid1 = _registerSubnet(alice, subnetC1);
        uint256 sid2 = _registerSubnet(alice, subnetC2);
        vm.startPrank(alice);
        rootNet.activateSubnet(sid1);
        rootNet.activateSubnet(sid2);
        vm.stopPrank();

        {
            address[] memory addrs = new address[](2);
            addrs[0] = subnetC1;
            addrs[1] = subnetC2;
            uint96[] memory ws = new uint96[](2);
            ws[0] = 200;
            ws[1] = 100;
            _submitWeights(addrs, ws);
        }

        _depositAndAllocate(alice, agentA, sid1, 100_000 * 1e18, 80_000 * 1e18);
        _depositAndAllocate(bob, agentB, sid1, 50_000 * 1e18, 50_000 * 1e18);
        _depositAndAllocate(charlie, agentC, sid2, 30_000 * 1e18, 30_000 * 1e18);

        assertEq(vault.getSubnetTotalStake(sid1), 130_000 * 1e18);
        assertEq(vault.getSubnetTotalStake(sid2), 30_000 * 1e18);

        _settleEpoch();
        _settleEpoch();

        uint256 bal1 = awp.balanceOf(subnetC1);
        uint256 bal2 = awp.balanceOf(subnetC2);
        assertApproxEqRel(bal1, bal2 * 2, 0.01e18);
    }

    // ════════════════════════════════════════════
    //  E2E 7: Binding tree (replaces mutual exclusion)
    // ════════════════════════════════════════════

    function test_e2e_bindingTree() public {
        // Build chain: agentA -> alice, agentB -> agentA
        _bindAgent(agentA, alice);
        _bindAgent(agentB, agentA);

        assertEq(rootNet.boundTo(agentA), alice);
        assertEq(rootNet.boundTo(agentB), agentA);

        // resolveRecipient should walk to alice
        vm.prank(alice);
        rootNet.setRecipient(charlie);
        assertEq(rootNet.resolveRecipient(agentB), charlie);
    }

    // ════════════════════════════════════════════
    //  E2E 8: Emission clamp near MAX_SUPPLY
    // ════════════════════════════════════════════

    function test_e2e_emissionClampNearMaxSupply() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);
        _submitWeight(subnetC1, uint96(100));

        _settleEpoch();
        _settleEpoch();
        uint256 afterEpoch1 = awp.totalSupply();
        assertTrue(afterEpoch1 > awp.INITIAL_MINT());

        uint256 remaining = awp.MAX_SUPPLY() - awp.totalSupply();
        assertTrue(remaining > 0);
    }

    // ════════════════════════════════════════════
    //  E2E 9: notSettling guard
    // ════════════════════════════════════════════

    function test_e2e_notSettlingGuard() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        uint256 sid2 = _registerSubnet(alice, subnetC2);
        vm.startPrank(alice);
        rootNet.activateSubnet(sid);
        rootNet.activateSubnet(sid2);
        vm.stopPrank();

        {
            address[] memory addrs = new address[](2);
            addrs[0] = subnetC1;
            addrs[1] = subnetC2;
            uint96[] memory ws = new uint96[](2);
            ws[0] = 100;
            ws[1] = 100;
            _submitWeights(addrs, ws);
        }

        _settleEpoch();

        vm.warp(block.timestamp + EPOCH + 1);
        emission.settleEpoch(1);
        assertTrue(emission.settleProgress() > 0);

        emission.settleEpoch(200);
        assertEq(emission.settleProgress(), 0);

        uint256 sid3 = _registerSubnet(alice, subnetC3);
        vm.prank(alice);
        rootNet.activateSubnet(sid3);
        assertTrue(rootNet.isSubnetActive(sid3));
    }

    // ════════════════════════════════════════════
    //  E2E 10: Deallocate releases allocations
    // ════════════════════════════════════════════

    function test_e2e_deallocateReleasesAllocations() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        _depositAndAllocate(alice, agentA, sid, 10_000 * 1e18, 5_000 * 1e18);

        // Deallocate all
        vm.prank(alice);
        rootNet.deallocate(alice, agentA, sid, 5_000 * 1e18);

        assertEq(vault.getAgentStake(alice, agentA, sid), 0);
        assertEq(vault.userTotalAllocated(alice), 0);
    }

    // ════════════════════════════════════════════
    //  E2E 11: Deposit + allocate flow
    // ════════════════════════════════════════════

    function test_e2e_depositAndAllocate() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        vm.startPrank(bob);
        awp.approve(address(stakeNFT), 20_000 * 1e18);
        stakeNFT.deposit(20_000 * 1e18, 52 weeks);
        rootNet.allocate(bob, agentB, sid, 15_000 * 1e18);
        vm.stopPrank();

        assertEq(stakeNFT.getUserTotalStaked(bob), 20_000 * 1e18);
        assertEq(vault.getAgentStake(bob, agentB, sid), 15_000 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 12: Reward Recipient
    // ════════════════════════════════════════════

    function test_e2e_rewardRecipient() public {
        _registerUser(alice);
        _bindAgent(agentA, alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        assertEq(rootNet.resolveRecipient(alice), alice);

        vm.prank(alice);
        rootNet.setRecipient(bob);
        assertEq(rootNet.resolveRecipient(alice), bob);

        _depositAndAllocate(alice, agentA, sid, 1_000 * 1e18, 500 * 1e18);
        AWPRegistry.AgentInfo memory info = rootNet.getAgentInfo(agentA, sid);
        assertEq(info.rewardRecipient, bob);
    }

    // ════════════════════════════════════════════
    //  E2E 13: Multi-epoch with stake changes
    // ════════════════════════════════════════════

    function test_e2e_multiEpochWithStakeChanges() public {
        _registerUser(alice);
        _registerUser(bob);

        uint256 sid1 = _registerSubnet(alice, subnetC1);
        uint256 sid2 = _registerSubnet(alice, subnetC2);
        vm.startPrank(alice);
        rootNet.activateSubnet(sid1);
        rootNet.activateSubnet(sid2);
        vm.stopPrank();

        {
            address[] memory addrs = new address[](2);
            addrs[0] = subnetC1;
            addrs[1] = subnetC2;
            uint96[] memory ws = new uint96[](2);
            ws[0] = 100;
            ws[1] = 100;
            _submitWeights(addrs, ws);
        }

        _depositAndAllocate(alice, agentA, sid1, 10_000 * 1e18, 8_000 * 1e18);
        _depositAndAllocate(bob, agentB, sid2, 5_000 * 1e18, 5_000 * 1e18);

        _settleEpoch();
        _settleEpoch();
        uint256 sc1Bal1 = awp.balanceOf(subnetC1);
        uint256 sc2Bal1 = awp.balanceOf(subnetC2);
        assertEq(sc1Bal1, sc2Bal1);

        // Alice deallocates
        vm.prank(alice);
        rootNet.deallocate(alice, agentA, sid1, 3_000 * 1e18);

        _submitWeights(
            _toArray2(subnetC1, subnetC2),
            _toUint96Array2(100, 100)
        );
        _settleEpoch();

        // Bob deallocates
        vm.prank(bob);
        rootNet.deallocate(bob, agentB, sid2, 5_000 * 1e18);

        _submitWeights(
            _toArray2(subnetC1, subnetC2),
            _toUint96Array2(100, 100)
        );
        _settleEpoch();
        assertEq(emission.currentEpoch(), 4);

        assertTrue(awp.balanceOf(subnetC1) > sc1Bal1);
        assertTrue(awp.balanceOf(subnetC2) > sc2Bal1);
    }

    // ════════════════════════════════════════════
    //  E2E 14: Agent unbind
    // ════════════════════════════════════════════

    function test_e2e_agentUnbind() public {
        _bindAgent(agentA, alice);
        assertEq(rootNet.boundTo(agentA), alice);

        vm.prank(agentA);
        rootNet.unbind();
        assertEq(rootNet.boundTo(agentA), address(0));

        // agentA can re-bind
        _bindAgent(agentA, alice);
        assertEq(rootNet.boundTo(agentA), alice);
    }

    // ════════════════════════════════════════════
    //  E2E 15: Alpha Token subnet contract minting
    // ════════════════════════════════════════════

    function test_e2e_alphaTokenMinting() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);

        AlphaToken alpha = AlphaToken(rootNet.getSubnetFull(sid).alphaToken);
        assertFalse(alpha.minters(address(rootNet)));
        assertTrue(alpha.mintersLocked());

        vm.warp(block.timestamp + 10 days);
        vm.prank(subnetC1);
        alpha.mint(alice, 1_000_000 * 1e18);
        assertEq(alpha.balanceOf(alice), 1_000_000 * 1e18);

        vm.prank(alice);
        rootNet.activateSubnet(sid);
        vm.prank(address(treasury));
        rootNet.banSubnet(sid);

        vm.prank(subnetC1);
        vm.expectRevert(AlphaToken.MinterPaused.selector);
        alpha.mint(alice, 100);

        vm.prank(address(treasury));
        rootNet.unbanSubnet(sid);
        vm.warp(block.timestamp + 1 days);
        vm.prank(subnetC1);
        alpha.mint(alice, 500_000 * 1e18);
        assertEq(alpha.balanceOf(alice), 1_500_000 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 16: Guardian emergency pause
    // ════════════════════════════════════════════

    function test_e2e_emergencyPause() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);
        _depositAndAllocate(alice, agentA, sid, 10_000 * 1e18, 5_000 * 1e18);

        vm.prank(guardian);
        rootNet.pause();

        // User operations blocked
        vm.startPrank(alice);
        vm.expectRevert();
        rootNet.register();
        vm.expectRevert();
        rootNet.allocate(alice, agentA, sid, 100);
        vm.expectRevert();
        rootNet.deallocate(alice, agentA, sid, 100);
        vm.expectRevert();
        rootNet.activateSubnet(sid);
        vm.stopPrank();

        // Emission unaffected
        vm.warp(block.timestamp + EPOCH + 1);
        emission.settleEpoch(200);
        assertEq(emission.settledEpoch(), 1);

        // Unpause
        vm.prank(address(treasury));
        rootNet.unpause();

        vm.prank(alice);
        rootNet.deallocate(alice, agentA, sid, 1_000 * 1e18);
        assertEq(vault.getAgentStake(alice, agentA, sid), 4_000 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 17: Batch Agent query
    // ════════════════════════════════════════════

    function test_e2e_batchAgentQuery() public {
        _registerUser(alice);
        _registerUser(bob);
        _bindAgent(agentA, alice);
        _bindAgent(agentB, bob);

        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        _depositAndAllocate(alice, agentA, sid, 5_000 * 1e18, 3_000 * 1e18);
        _depositAndAllocate(bob, agentB, sid, 2_000 * 1e18, 2_000 * 1e18);

        vm.prank(alice);
        rootNet.setRecipient(charlie);

        address[] memory agents = new address[](3);
        agents[0] = agentA;
        agents[1] = agentB;
        agents[2] = address(0x9999);

        AWPRegistry.AgentInfo[] memory infos = rootNet.getAgentsInfo(agents, sid);

        assertEq(infos.length, 3);
        assertEq(infos[0].root, alice);
        assertTrue(infos[0].isValid);
        assertEq(infos[0].stake, 3_000 * 1e18);
        assertEq(infos[0].rewardRecipient, charlie);

        assertEq(infos[1].root, bob);
        assertTrue(infos[1].isValid);
        assertEq(infos[1].stake, 2_000 * 1e18);

        assertEq(infos[2].root, address(0x9999));
        assertFalse(infos[2].isValid);
        assertEq(infos[2].stake, 0);
    }

    // ════════════════════════════════════════════
    //  E2E 18: Gasless bind
    // ════════════════════════════════════════════

    function test_e2e_gaslessAgentBind() public {
        (address agent, uint256 agentPk) = makeAddrAndKey("gaslessAgent");

        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = rootNet.nonces(agent);

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Bind(address agent,address target,uint256 nonce,uint256 deadline)"),
            agent, alice, nonce, deadline
        ));
        bytes32 digest = _getDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(agentPk, digest);

        address relayer = address(0x9999);
        vm.prank(relayer);
        rootNet.bindFor(agent, alice, deadline, v, r, s);

        assertEq(rootNet.boundTo(agent), alice);
    }

    // ════════════════════════════════════════════
    //  E2E 19: Gasless expired signature reverts
    // ════════════════════════════════════════════

    function test_e2e_gaslessExpiredSignature() public {
        (address agent, uint256 agentPk) = makeAddrAndKey("expiredAgent");

        uint256 deadline = block.timestamp - 1;
        bytes32 structHash = keccak256(abi.encode(
            keccak256("Bind(address agent,address target,uint256 nonce,uint256 deadline)"),
            agent, alice, 0, deadline
        ));
        bytes32 digest = _getDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(agentPk, digest);

        vm.expectRevert(AWPRegistry.ExpiredSignature.selector);
        rootNet.bindFor(agent, alice, deadline, v, r, s);
    }

    // ════════════════════════════════════════════
    //  E2E 20: Multi-user same subnet reallocate
    // ════════════════════════════════════════════

    function test_e2e_multiUserSameSubnetReallocate() public {
        _registerUser(alice);
        _registerUser(bob);

        uint256 sid1 = _registerSubnet(alice, subnetC1);
        uint256 sid2 = _registerSubnet(alice, subnetC2);
        vm.startPrank(alice);
        rootNet.activateSubnet(sid1);
        rootNet.activateSubnet(sid2);
        vm.stopPrank();

        _depositAndAllocate(alice, agentA, sid1, 1_000 * 1e18, 500 * 1e18);
        _depositAndAllocate(bob, agentB, sid1, 1_000 * 1e18, 300 * 1e18);

        // Reallocate
        vm.prank(alice);
        rootNet.reallocate(alice, agentA, sid1, agentA, sid2, 200 * 1e18);
        vm.prank(bob);
        rootNet.reallocate(bob, agentB, sid1, agentB, sid2, 100 * 1e18);

        assertEq(vault.getAgentStake(alice, agentA, sid1), 300 * 1e18);
        assertEq(vault.getAgentStake(alice, agentA, sid2), 200 * 1e18);
        assertEq(vault.getAgentStake(bob, agentB, sid1), 200 * 1e18);
        assertEq(vault.getAgentStake(bob, agentB, sid2), 100 * 1e18);

        assertEq(vault.getSubnetTotalStake(sid1), 500 * 1e18);
        assertEq(vault.getSubnetTotalStake(sid2), 300 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 21: Batched settlement verification
    // ════════════════════════════════════════════

    function test_e2e_batchSettleMultiCall() public {
        _registerUser(alice);

        address[5] memory scs = [
            address(0x3001),
            address(0x3002),
            address(0x3003),
            address(0x3004),
            address(0x3005)
        ];

        vm.startPrank(alice);
        awp.approve(address(rootNet), LP_COST * 5);
        for (uint256 i = 0; i < 5; i++) {
            rootNet.registerSubnet(
                IAWPRegistry.SubnetParams("S", "S", scs[i], bytes32(0), 0, "")
            );
            rootNet.activateSubnet(i + 1);
        }
        vm.stopPrank();

        {
            address[] memory addrs = new address[](5);
            uint96[] memory ws = new uint96[](5);
            for (uint256 i = 0; i < 5; i++) {
                addrs[i] = scs[i];
                ws[i] = uint96((i + 1) * 100);
            }
            _submitWeights(addrs, ws);
        }

        _settleEpoch();
        _settleEpoch();

        for (uint256 i = 0; i < 5; i++) {
            assertTrue(awp.balanceOf(scs[i]) > 0);
        }
    }

    // ════════════════════════════════════════════
    //  E2E 22: Delegate operations
    // ════════════════════════════════════════════

    function test_e2e_delegateOperations() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        _depositAndAllocate(alice, agentA, sid, 10_000 * 1e18, 5_000 * 1e18);

        // Grant delegate to bob
        vm.prank(alice);
        rootNet.grantDelegate(bob);

        // Bob can deallocate on behalf of alice
        vm.prank(bob);
        rootNet.deallocate(alice, agentA, sid, 2_000 * 1e18);
        assertEq(vault.getAgentStake(alice, agentA, sid), 3_000 * 1e18);

        // Revoke delegation
        vm.prank(alice);
        rootNet.revokeDelegate(bob);

        // Bob can no longer operate
        vm.prank(bob);
        vm.expectRevert(AWPRegistry.NotAuthorized.selector);
        rootNet.deallocate(alice, agentA, sid, 1_000 * 1e18);
    }

    // ── EIP-712 helpers ──

    function _getDigest(bytes32 structHash) internal view returns (bytes32) {
        bytes32 domainSeparator = keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
                keccak256("AWPRegistry"),
                keccak256("1"),
                block.chainid,
                address(rootNet)
            )
        );
        return keccak256(abi.encodePacked("\x19\x01", domainSeparator, structHash));
    }

    // ── Helper arrays ──

    function _toArray2(address a, address b) internal pure returns (address[] memory) {
        address[] memory arr = new address[](2);
        arr[0] = a;
        arr[1] = b;
        return arr;
    }

    function _toUint96Array2(uint96 a, uint96 b) internal pure returns (uint96[] memory) {
        uint96[] memory arr = new uint96[](2);
        arr[0] = a;
        arr[1] = b;
        return arr;
    }
}
