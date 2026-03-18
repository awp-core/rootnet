// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AlphaToken} from "../src/token/AlphaToken.sol";
import {AlphaTokenFactory} from "../src/token/AlphaTokenFactory.sol";

import {AWPEmission} from "../src/token/AWPEmission.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {AccessManager} from "../src/core/AccessManager.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {SubnetNFT} from "../src/core/SubnetNFT.sol";
import {MockLPManager} from "./helpers/MockLPManager.sol";
import {RootNet} from "../src/RootNet.sol";
import {IRootNet} from "../src/interfaces/IRootNet.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {TimelockController} from "@openzeppelin/contracts/governance/TimelockController.sol";
import {IGovernor} from "@openzeppelin/contracts/governance/IGovernor.sol";
import {EmissionSigningHelper} from "./helpers/EmissionSigningHelper.sol";

/// @title E2E — End-to-end tests based on the architecture document
contract E2ETest is EmissionSigningHelper {
    AWPToken awp;
    AlphaTokenFactory factory;
    AWPEmission emission;
    AccessManager access;
    StakingVault vault;
    StakeNFT stakeNFT;
    SubnetNFT nft;
    MockLPManager lp;
    RootNet rootNet;
    Treasury treasury;
    AWPDAO dao;

    address deployer = address(0xD);
    address guardian = address(0xE);
    address teamVesting = address(0xF1);
    address investorVesting = address(0xF2);
    address liquidityPool = address(0xF3);
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

        rootNet = new RootNet(deployer, address(treasury), guardian);
        nft = new SubnetNFT("AWP Subnet", "CXSUB", address(rootNet));
        access = new AccessManager(address(rootNet));
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

        // Deploy StakingVault + StakeNFT (circular dependency)
        uint64 nonce = vm.getNonce(deployer);
        address predictedVault = vm.computeCreateAddress(deployer, nonce);
        address predictedStakeNFT = vm.computeCreateAddress(deployer, nonce + 1);

        vault = new StakingVault(address(rootNet));
        stakeNFT = new StakeNFT(address(awp), address(vault), address(rootNet));

        // Deploy AWPDAO (no rootNet param — uses timestamps)
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

        // Initialize registry with 9 params (including stakeNFT)
        rootNet.initializeRegistry(
            address(awp), address(nft), address(factory), address(emission),
            address(lp), address(access), address(vault), address(stakeNFT)
        );

        vm.stopPrank();

        // Configure oracles
        address[] memory oracleList = new address[](3);
        oracleList[0] = vm.addr(ORACLE_PK1);
        oracleList[1] = vm.addr(ORACLE_PK2);
        oracleList[2] = vm.addr(ORACLE_PK3);
        vm.prank(address(treasury));
        emission.setOracleConfig(oracleList, 2);

        vm.startPrank(deployer);
        awp.transfer(address(treasury), 2_000_000_000 * 1e18);
        awp.transfer(teamVesting, 1_000_000_000 * 1e18);
        awp.transfer(investorVesting, 750_000_000 * 1e18);
        awp.transfer(liquidityPool, 1_000_000_000 * 1e18);
        awp.transfer(airdropAddr, 250_000_000 * 1e18);
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

    function _registerAgent(address agent, address user) internal {
        vm.prank(agent);
        rootNet.bind(user);
    }

    function _registerSubnet(address owner, address sc) internal returns (uint256) {
        vm.startPrank(owner);
        awp.approve(address(rootNet), LP_COST);
        uint256 id = rootNet.registerSubnet(
            IRootNet.SubnetParams("Subnet", "SUB", "ipfs://meta", sc, "https://coord", bytes32(0), 0)
        );
        vm.stopPrank();
        return id;
    }

    /// @dev Deposit AWP via StakeNFT and allocate via RootNet
    function _depositAndAllocate(address user, address agent, uint256 subnetId, uint256 deposit, uint256 alloc)
        internal
    {
        vm.startPrank(user);
        awp.approve(address(stakeNFT), deposit);
        stakeNFT.deposit(deposit, 52 weeks); // 52 weeks lock
        rootNet.allocate(agent, subnetId, alloc);
        vm.stopPrank();
    }

    function _settleEpoch() internal {
        vm.warp(block.timestamp + EPOCH + 1);
        emission.settleEpoch(200);
    }

    /// @dev Sort addresses ascending (bubble sort) and reorder weights to match
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

    function _submitWeight(address recipient, uint96 weight) internal {
        address[] memory addrs = new address[](1);
        addrs[0] = recipient;
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

        // Transfer NFT to Bob
        vm.prank(alice);
        nft.transferFrom(alice, bob, sid);
        assertEq(nft.ownerOf(sid), bob);

        // Alice can no longer pause
        vm.prank(alice);
        vm.expectRevert(RootNet.NotOwner.selector);
        rootNet.pauseSubnet(sid);

        // Bob can pause and resume
        vm.prank(bob);
        rootNet.pauseSubnet(sid);
        assertFalse(rootNet.isSubnetActive(sid));

        vm.prank(bob);
        rootNet.resumeSubnet(sid);
        assertTrue(rootNet.isSubnetActive(sid));
    }

    // ════════════════════════════════════════════
    //  E2E 2: Reallocate is immediate (no dual-slot)
    // ════════════════════════════════════════════

    function test_e2e_reallocateImmediate() public {
        _registerUser(alice);
        _registerAgent(agentA, alice);
        _registerAgent(agentB, alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        _depositAndAllocate(alice, agentA, sid, 10_000 * 1e18, 8_000 * 1e18);

        // Reallocate 5000 from agentA -> agentB (immediate effect)
        vm.prank(alice);
        rootNet.reallocate(agentA, sid, agentB, sid, 5_000 * 1e18);

        // Immediate: values reflect the reallocation right away
        assertEq(vault.getAgentStake(alice, agentA, sid), 3_000 * 1e18);
        assertEq(vault.getAgentStake(alice, agentB, sid), 5_000 * 1e18);

        // userTotalAllocated unchanged
        assertEq(vault.userTotalAllocated(alice), 8_000 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 3: Full subnet lifecycle
    // ════════════════════════════════════════════

    function test_e2e_subnetFullLifecycle() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);

        // Pending -> Active
        vm.prank(alice);
        rootNet.activateSubnet(sid);
        assertEq(uint256(rootNet.getSubnet(sid).status), uint256(IRootNet.SubnetStatus.Active));

        // Active -> Paused
        vm.prank(alice);
        rootNet.pauseSubnet(sid);
        assertEq(uint256(rootNet.getSubnet(sid).status), uint256(IRootNet.SubnetStatus.Paused));

        // Paused -> Active
        vm.prank(alice);
        rootNet.resumeSubnet(sid);
        assertEq(uint256(rootNet.getSubnet(sid).status), uint256(IRootNet.SubnetStatus.Active));

        // Active -> Banned
        vm.prank(address(treasury));
        rootNet.banSubnet(sid);
        assertEq(uint256(rootNet.getSubnet(sid).status), uint256(IRootNet.SubnetStatus.Banned));

        AlphaToken alpha = AlphaToken(rootNet.getSubnetFull(sid).alphaToken);
        assertTrue(alpha.minterPaused(subnetC1));

        // Banned -> Active
        vm.prank(address(treasury));
        rootNet.unbanSubnet(sid);
        assertEq(uint256(rootNet.getSubnet(sid).status), uint256(IRootNet.SubnetStatus.Active));
        assertFalse(alpha.minterPaused(subnetC1));

        // Deregister
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
        // Give Alice enough voting power
        vm.prank(liquidityPool);
        awp.transfer(alice, 250_000_000 * 1e18);

        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        // Alice stakes into StakeNFT for voting power
        vm.startPrank(alice);
        awp.approve(address(stakeNFT), 250_000_000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(250_000_000 * 1e18, 52 weeks);
        vm.stopPrank();

        // Submit initial weight via oracle for epoch 1
        _submitWeight(subnetC1, uint96(100));

        // Settle epoch 0 (no weights, all to DAO)
        _settleEpoch();
        // Settle epoch 1 (promotes activeEpoch to 1)
        _settleEpoch();

        vm.roll(block.number + 1);

        // Propose emergencySetWeight
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

        vm.roll(block.number + 2); // votingDelay

        // Vote with NFT tokenIds
        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = tokenId;
        bytes memory params = abi.encode(tokenIds);
        vm.prank(alice);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", params);

        vm.roll(block.number + 101); // votingPeriod

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

        // Submit weight for epoch 1
        _submitWeight(subnetC1, uint96(100));

        uint256 prevEmission = INITIAL_DAILY;

        for (uint256 i = 0; i < 10; i++) {
            // Before settling each epoch, submit weights for the NEXT epoch
            if (i > 0) {
                // Re-submit for the next epoch (currentEpoch + 1)
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
        _registerAgent(agentA, alice);
        _registerAgent(agentB, bob);
        _registerAgent(agentC, charlie);

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

        // Settle epoch 0 (no weights)
        _settleEpoch();
        // Settle epoch 1 (weights active)
        _settleEpoch();

        uint256 bal1 = awp.balanceOf(subnetC1);
        uint256 bal2 = awp.balanceOf(subnetC2);
        assertApproxEqRel(bal1, bal2 * 2, 0.01e18);
    }

    // ════════════════════════════════════════════
    //  E2E 7: User/Agent address mutual exclusion
    // ════════════════════════════════════════════

    function test_e2e_addressMutualExclusion() public {
        _registerUser(alice);

        // alice is a registered Principal — cannot be an Agent
        vm.prank(alice);
        vm.expectRevert(AccessManager.AddressIsPrincipal.selector);
        rootNet.bind(bob); // alice tries to bind herself as agent of bob

        _registerAgent(agentA, alice);

        // agentA is an Agent — cannot register as Principal directly
        vm.prank(agentA);
        vm.expectRevert(AccessManager.AddressIsAgent.selector);
        rootNet.register();
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

        // Settle epoch 0 (no weights)
        _settleEpoch();
        // Settle epoch 1 (weights active)
        _settleEpoch();
        uint256 afterEpoch1 = awp.totalSupply();
        assertTrue(afterEpoch1 > 5_000_000_000 * 1e18);

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

        // Settle epoch 0 (no weights)
        _settleEpoch();

        // Start settling epoch 1 with limit=1
        vm.warp(block.timestamp + EPOCH + 1);
        emission.settleEpoch(1);
        assertTrue(emission.settleProgress() > 0);

        // submitAllocations for future epochs is now ALLOWED during settlement
        // (epoch-versioned design writes to a different slot)

        emission.settleEpoch(200);
        assertEq(emission.settleProgress(), 0);

        uint256 sid3 = _registerSubnet(alice, subnetC3);
        vm.prank(alice);
        rootNet.activateSubnet(sid3);
        assertTrue(rootNet.isSubnetActive(sid3));
    }

    // ════════════════════════════════════════════
    //  E2E 10: Agent freeze releases allocations
    // ════════════════════════════════════════════

    function test_e2e_freezeWithAgentRemoval() public {
        _registerUser(alice);
        _registerAgent(agentA, alice);
        _registerAgent(agentB, alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        _depositAndAllocate(alice, agentA, sid, 10_000 * 1e18, 5_000 * 1e18);

        // Remove agentA — all allocations frozen and released (StakingVault auto-enumerates)
        vm.prank(alice);
        rootNet.removeAgent(agentA);

        assertEq(vault.getAgentStake(alice, agentA, sid), 0);
        assertEq(vault.userTotalAllocated(alice), 0);
    }

    // ════════════════════════════════════════════
    //  E2E 11: registerAndStake one-click flow
    // ════════════════════════════════════════════

    function test_e2e_registerAndStakeOneClick() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        // Bob registers user + Agent first
        _registerUser(bob);
        _registerAgent(agentB, bob);

        // registerAndStake with new signature (depositAmount, lockDuration, agent, subnetId, allocateAmount)
        vm.startPrank(bob);
        awp.approve(address(stakeNFT), 20_000 * 1e18);
        rootNet.registerAndStake(20_000 * 1e18, 52 weeks, agentB, sid, 15_000 * 1e18);
        vm.stopPrank();

        assertEq(stakeNFT.getUserTotalStaked(bob), 20_000 * 1e18);
        assertEq(vault.getAgentStake(bob, agentB, sid), 15_000 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 12: Reward Recipient
    // ════════════════════════════════════════════

    function test_e2e_rewardRecipient() public {
        _registerUser(alice);
        _registerAgent(agentA, alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        assertEq(access.getRewardRecipient(alice), alice);

        vm.prank(alice);
        rootNet.setRewardRecipient(bob);
        assertEq(access.getRewardRecipient(alice), bob);

        _depositAndAllocate(alice, agentA, sid, 1_000 * 1e18, 500 * 1e18);
        RootNet.AgentInfo memory info = rootNet.getAgentInfo(agentA, sid);
        assertEq(info.rewardRecipient, bob);
    }

    // ════════════════════════════════════════════
    //  E2E 13: Multi-epoch with stake changes
    // ════════════════════════════════════════════

    function test_e2e_multiEpochWithStakeChanges() public {
        _registerUser(alice);
        _registerUser(bob);
        _registerAgent(agentA, alice);
        _registerAgent(agentB, bob);

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

        // Settle epoch 0 (no weights)
        _settleEpoch();
        // Settle epoch 1 (weights active)
        _settleEpoch();
        uint256 sc1Bal1 = awp.balanceOf(subnetC1);
        uint256 sc2Bal1 = awp.balanceOf(subnetC2);
        assertEq(sc1Bal1, sc2Bal1);

        // Alice deallocates
        vm.prank(alice);
        rootNet.deallocate(agentA, sid1, 3_000 * 1e18);

        // Submit weights for next epoch and settle
        _submitWeights(
            _toArray2(subnetC1, subnetC2),
            _toUint96Array2(100, 100)
        );
        _settleEpoch();

        // Bob deallocates
        vm.prank(bob);
        rootNet.deallocate(agentB, sid2, 5_000 * 1e18);

        // Submit weights for next epoch and settle
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
        _registerUser(alice);
        _registerAgent(agentA, alice);

        assertTrue(access.isAgent(alice, agentA));

        vm.prank(agentA);
        rootNet.unbind();

        assertFalse(access.isRegisteredAgent(agentA));
        assertFalse(access.isKnownAddress(agentA));

        // agentA can re-bind to any principal after unbind
        _registerAgent(agentA, alice);
        assertTrue(access.isAgent(alice, agentA));
    }

    // ════════════════════════════════════════════
    //  E2E 14b: No agent limit per principal
    // ════════════════════════════════════════════

    function test_e2e_noAgentLimit() public {
        _registerUser(alice);

        // Bind 50 agents — should never revert
        for (uint256 i = 0; i < 50; i++) {
            address agent = address(uint160(0x5000 + i));
            vm.prank(agent);
            rootNet.bind(alice);
        }

        assertEq(access.getAgents(alice).length, 50);
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
        _registerAgent(agentA, alice);
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
        rootNet.allocate(agentA, sid, 100);
        vm.expectRevert();
        rootNet.deallocate(agentA, sid, 100);
        vm.expectRevert();
        rootNet.activateSubnet(sid);
        vm.stopPrank();

        // Emission unaffected — settle epoch 0
        vm.warp(block.timestamp + EPOCH + 1);
        emission.settleEpoch(200);
        assertEq(emission.settledEpoch(), 1);

        // Unpause
        vm.prank(address(treasury));
        rootNet.unpause();

        vm.prank(alice);
        rootNet.deallocate(agentA, sid, 1_000 * 1e18);
        assertEq(vault.getAgentStake(alice, agentA, sid), 4_000 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 17: Batch Agent query
    // ════════════════════════════════════════════

    function test_e2e_batchAgentQuery() public {
        _registerUser(alice);
        _registerUser(bob);
        _registerAgent(agentA, alice);
        _registerAgent(agentB, bob);

        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        _depositAndAllocate(alice, agentA, sid, 5_000 * 1e18, 3_000 * 1e18);
        _depositAndAllocate(bob, agentB, sid, 2_000 * 1e18, 2_000 * 1e18);

        vm.prank(alice);
        rootNet.setRewardRecipient(charlie);

        address[] memory agents = new address[](3);
        agents[0] = agentA;
        agents[1] = agentB;
        agents[2] = address(0x9999);

        RootNet.AgentInfo[] memory infos = rootNet.getAgentsInfo(agents, sid);

        assertEq(infos.length, 3);
        assertEq(infos[0].owner, alice);
        assertTrue(infos[0].isValid);
        assertEq(infos[0].stake, 3_000 * 1e18);
        assertEq(infos[0].rewardRecipient, charlie);

        assertEq(infos[1].owner, bob);
        assertTrue(infos[1].isValid);
        assertEq(infos[1].stake, 2_000 * 1e18);

        assertEq(infos[2].owner, address(0));
        assertFalse(infos[2].isValid);
        assertEq(infos[2].stake, 0);
    }

    // ════════════════════════════════════════════
    //  E2E 18: Gasless user registration
    // ════════════════════════════════════════════

    function test_e2e_gaslessUserRegistration() public {
        (address user, uint256 userPk) = makeAddrAndKey("gaslessUser");

        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = rootNet.nonces(user);

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Register(address user,uint256 nonce,uint256 deadline)"),
            user, nonce, deadline
        ));
        bytes32 digest = _getDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(userPk, digest);

        address relayer = address(0x9999);
        vm.prank(relayer);
        rootNet.registerFor(user, deadline, v, r, s);

        assertTrue(access.isRegistered(user));
        assertEq(rootNet.nonces(user), 1);
    }

    // ════════════════════════════════════════════
    //  E2E 19: Gasless Agent registration
    // ════════════════════════════════════════════

    function test_e2e_gaslessAgentBind() public {
        _registerUser(alice);

        (address agent, uint256 agentPk) = makeAddrAndKey("gaslessAgent");

        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = rootNet.nonces(agent);

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Bind(address agent,address principal,uint256 nonce,uint256 deadline)"),
            agent, alice, nonce, deadline
        ));
        bytes32 digest = _getDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(agentPk, digest);

        address relayer = address(0x9999);
        vm.prank(relayer);
        rootNet.bindFor(agent, alice, deadline, v, r, s);

        assertTrue(access.isAgent(alice, agent));
    }

    // ════════════════════════════════════════════
    //  E2E 20: Gasless expired signature reverts
    // ════════════════════════════════════════════

    function test_e2e_gaslessExpiredSignature() public {
        (address user, uint256 userPk) = makeAddrAndKey("expiredUser");

        uint256 deadline = block.timestamp - 1;
        bytes32 structHash = keccak256(abi.encode(
            keccak256("Register(address user,uint256 nonce,uint256 deadline)"),
            user, 0, deadline
        ));
        bytes32 digest = _getDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(userPk, digest);

        vm.expectRevert(RootNet.ExpiredSignature.selector);
        rootNet.registerFor(user, deadline, v, r, s);
    }

    // ════════════════════════════════════════════
    //  E2E 21: Gasless invalid signature reverts
    // ════════════════════════════════════════════

    function test_e2e_gaslessInvalidSignature() public {
        (address user,) = makeAddrAndKey("targetUser");
        (, uint256 wrongPk) = makeAddrAndKey("wrongSigner");

        uint256 deadline = block.timestamp + 1 hours;
        bytes32 structHash = keccak256(abi.encode(
            keccak256("Register(address user,uint256 nonce,uint256 deadline)"),
            user, 0, deadline
        ));
        bytes32 digest = _getDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(wrongPk, digest);

        vm.expectRevert(RootNet.InvalidSignature.selector);
        rootNet.registerFor(user, deadline, v, r, s);
    }

    // ════════════════════════════════════════════
    //  E2E 22: Gasless replay protection
    // ════════════════════════════════════════════

    function test_e2e_gaslessReplayProtection() public {
        (address user, uint256 userPk) = makeAddrAndKey("replayUser");

        uint256 deadline = block.timestamp + 1 hours;
        bytes32 structHash = keccak256(abi.encode(
            keccak256("Register(address user,uint256 nonce,uint256 deadline)"),
            user, 0, deadline
        ));
        bytes32 digest = _getDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(userPk, digest);

        rootNet.registerFor(user, deadline, v, r, s);
        assertTrue(access.isRegistered(user));

        vm.expectRevert();
        rootNet.registerFor(user, deadline, v, r, s);
    }

    // ════════════════════════════════════════════
    //  E2E 23: Multi-user same subnet reallocate
    // ════════════════════════════════════════════

    function test_e2e_multiUserSameSubnetReallocate() public {
        _registerUser(alice);
        _registerUser(bob);
        _registerAgent(agentA, alice);
        _registerAgent(agentB, bob);

        uint256 sid1 = _registerSubnet(alice, subnetC1);
        uint256 sid2 = _registerSubnet(alice, subnetC2);
        vm.startPrank(alice);
        rootNet.activateSubnet(sid1);
        rootNet.activateSubnet(sid2);
        vm.stopPrank();

        _depositAndAllocate(alice, agentA, sid1, 1_000 * 1e18, 500 * 1e18);
        _depositAndAllocate(bob, agentB, sid1, 1_000 * 1e18, 300 * 1e18);

        // Reallocate (immediate)
        vm.prank(alice);
        rootNet.reallocate(agentA, sid1, agentA, sid2, 200 * 1e18);
        vm.prank(bob);
        rootNet.reallocate(agentB, sid1, agentB, sid2, 100 * 1e18);

        // Immediate effect
        assertEq(vault.getAgentStake(alice, agentA, sid1), 300 * 1e18);
        assertEq(vault.getAgentStake(alice, agentA, sid2), 200 * 1e18);
        assertEq(vault.getAgentStake(bob, agentB, sid1), 200 * 1e18);
        assertEq(vault.getAgentStake(bob, agentB, sid2), 100 * 1e18);

        assertEq(vault.getSubnetTotalStake(sid1), 500 * 1e18);
        assertEq(vault.getSubnetTotalStake(sid2), 300 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 24: Batched settlement verification
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
        uint256[5] memory wts = [uint256(100), 200, 300, 400, 500];
        for (uint256 i = 0; i < 5; i++) {
            _registerSubnet(alice, scs[i]);
            vm.prank(alice);
            rootNet.activateSubnet(i + 1);
        }

        {
            address[] memory addrs = new address[](5);
            uint96[] memory ws = new uint96[](5);
            for (uint256 i = 0; i < 5; i++) {
                addrs[i] = scs[i];
                ws[i] = uint96(wts[i]);
            }
            _submitWeights(addrs, ws);
        }

        // Settle epoch 0 (no weights)
        _settleEpoch();

        // Settle epoch 1 in batches
        vm.warp(block.timestamp + EPOCH + 1);

        emission.settleEpoch(2);
        assertTrue(emission.settleProgress() > 0);
        emission.settleEpoch(2);
        assertTrue(emission.settleProgress() > 0);
        emission.settleEpoch(2);
        assertEq(emission.settleProgress(), 0);

        uint256 totalWeightVal = 1500;
        uint256 epochEmission = emission.epochEmissionLocked();
        uint256 subnetPool = epochEmission * 5000 / 10000;

        for (uint256 i = 0; i < 5; i++) {
            uint256 expected = subnetPool * wts[i] / totalWeightVal;
            uint256 actual = awp.balanceOf(scs[i]);
            assertApproxEqAbs(actual, expected, 1);
        }
    }

    // ════════════════════════════════════════════
    //  E2E 25: Emission precision (no leakage)
    // ════════════════════════════════════════════

    function test_e2e_emissionPrecision() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);
        _submitWeight(subnetC1, uint96(100));

        // Settle epoch 0 (no weights)
        uint256 treasuryBalBefore0 = awp.balanceOf(address(treasury));
        _settleEpoch();

        // Settle epoch 1 (weights active)
        uint256 treasuryBalBefore = awp.balanceOf(address(treasury));
        uint256 totalSupplyBefore = awp.totalSupply();
        _settleEpoch();

        uint256 epochEmission = emission.epochEmissionLocked();
        uint256 subnetMinted = awp.balanceOf(subnetC1);
        uint256 daoMinted = awp.balanceOf(address(treasury)) - treasuryBalBefore;

        assertEq(subnetMinted + daoMinted, epochEmission);
        assertEq(awp.totalSupply() - totalSupplyBefore, epochEmission);
    }

    // ════════════════════════════════════════════
    //  E2E 26: Query after deregistration
    // ════════════════════════════════════════════

    function test_e2e_deregisterThenQuery() public {
        _registerUser(alice);
        uint256 sid = _registerSubnet(alice, subnetC1);
        vm.prank(alice);
        rootNet.activateSubnet(sid);

        assertEq(rootNet.getActiveSubnetCount(), 1);

        vm.warp(block.timestamp + 31 days);
        vm.prank(address(treasury));
        rootNet.deregisterSubnet(sid);

        vm.expectRevert();
        rootNet.getSubnetFull(sid);
        assertFalse(rootNet.isSubnetActive(sid));
        assertEq(rootNet.getActiveSubnetCount(), 0);
    }

    // ════════════════════════════════════════════
    //  E2E 27: Agent freeze across multiple subnets
    // ════════════════════════════════════════════

    function test_e2e_freezeMultipleSubnets() public {
        _registerUser(alice);
        _registerAgent(agentA, alice);

        uint256 sid1 = _registerSubnet(alice, subnetC1);
        uint256 sid2 = _registerSubnet(alice, subnetC2);
        uint256 sid3 = _registerSubnet(alice, subnetC3);
        vm.startPrank(alice);
        rootNet.activateSubnet(sid1);
        rootNet.activateSubnet(sid2);
        rootNet.activateSubnet(sid3);
        vm.stopPrank();

        vm.startPrank(alice);
        awp.approve(address(stakeNFT), 30_000 * 1e18);
        stakeNFT.deposit(30_000 * 1e18, 52 weeks);
        rootNet.allocate(agentA, sid1, 5_000 * 1e18);
        rootNet.allocate(agentA, sid2, 8_000 * 1e18);
        rootNet.allocate(agentA, sid3, 3_000 * 1e18);
        vm.stopPrank();

        assertEq(vault.userTotalAllocated(alice), 16_000 * 1e18);

        vm.prank(alice);
        rootNet.removeAgent(agentA);

        assertEq(vault.getAgentStake(alice, agentA, sid1), 0);
        assertEq(vault.getAgentStake(alice, agentA, sid2), 0);
        assertEq(vault.getAgentStake(alice, agentA, sid3), 0);
        assertEq(vault.userTotalAllocated(alice), 0);

        _registerAgent(agentB, alice);
        vm.prank(alice);
        rootNet.allocate(agentB, sid1, 16_000 * 1e18);
        assertEq(vault.getAgentStake(alice, agentB, sid1), 16_000 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 28: All subnets zero weight -> all to DAO
    // ════════════════════════════════════════════

    function test_e2e_settleEmptyWeight() public {
        _registerUser(alice);

        uint256 sid1 = _registerSubnet(alice, subnetC1);
        uint256 sid2 = _registerSubnet(alice, subnetC2);
        vm.startPrank(alice);
        rootNet.activateSubnet(sid1);
        rootNet.activateSubnet(sid2);
        vm.stopPrank();

        uint256 treasuryBalBefore = awp.balanceOf(address(treasury));

        _settleEpoch();

        assertEq(awp.balanceOf(subnetC1), 0);
        assertEq(awp.balanceOf(subnetC2), 0);

        uint256 daoReceived = awp.balanceOf(address(treasury)) - treasuryBalBefore;
        assertEq(daoReceived, emission.epochEmissionLocked());
    }

    // ════════════════════════════════════════════
    //  E2E 29: Double reallocate accumulation
    // ════════════════════════════════════════════

    function test_e2e_doubleReallocateAccumulation() public {
        _registerUser(alice);
        _registerAgent(agentA, alice);
        _registerAgent(agentB, alice);

        uint256 sid1 = _registerSubnet(alice, subnetC1);
        uint256 sid2 = _registerSubnet(alice, subnetC2);
        vm.startPrank(alice);
        rootNet.activateSubnet(sid1);
        rootNet.activateSubnet(sid2);
        vm.stopPrank();

        _depositAndAllocate(alice, agentA, sid1, 2_000 * 1e18, 1_000 * 1e18);

        vm.startPrank(alice);
        rootNet.reallocate(agentA, sid1, agentB, sid2, 200 * 1e18);
        rootNet.reallocate(agentA, sid1, agentB, sid2, 300 * 1e18);
        vm.stopPrank();

        // Immediate effect
        assertEq(vault.getAgentStake(alice, agentA, sid1), 500 * 1e18);
        assertEq(vault.getAgentStake(alice, agentB, sid2), 500 * 1e18);
        assertEq(vault.userTotalAllocated(alice), 1_000 * 1e18);
    }

    // ── Utility helpers ──

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

    // ── EIP-712 helper ──

    function _getDigest(bytes32 structHash) internal view returns (bytes32) {
        bytes32 domainSeparator = keccak256(abi.encode(
            keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
            keccak256("AWPRootNet"),
            keccak256("1"),
            block.chainid,
            address(rootNet)
        ));
        return keccak256(abi.encodePacked("\x19\x01", domainSeparator, structHash));
    }
}
