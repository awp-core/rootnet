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
import {Treasury} from "../src/governance/Treasury.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {TimelockController} from "@openzeppelin/contracts/governance/TimelockController.sol";
import {IGovernor} from "@openzeppelin/contracts/governance/IGovernor.sol";

/// @title E2E — End-to-end tests based on the architecture document (Account System V2)
contract E2ETest is Test {
    AWPToken awp;
    AlphaTokenFactory factory;
    AWPEmission emission;
    StakingVault vault;
    StakeNFT stakeNFT;
    WorknetNFT nft;
    MockLPManager lp;
    AWPRegistry awpRegistry;
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

    // Mock worknet contracts
    address worknetC1 = address(0x3001);
    address worknetC2 = address(0x3002);
    address worknetC3 = address(0x3003);

    uint256 constant INITIAL_DAILY = 31_600_000 * 1e18;
    uint256 constant EPOCH = 1 days;
    uint256 constant LP_COST = 1_000_000 * 1e18; // 100M Alpha * 0.01 AWP

    function setUp() public {
        _deploy();
    }

    function _deploy() internal {
        vm.startPrank(deployer);
        awp = new AWPToken("AWP Token", "AWP", deployer);
        awp.initialMint(200_000_000 * 1e18);
        factory = new AlphaTokenFactory(deployer, 0);

        address[] memory p = new address[](0);
        address[] memory e = new address[](1);
        e[0] = address(0);
        treasury = new Treasury(1, p, e, deployer);

        AWPRegistry awpRegistryImpl = new AWPRegistry();
        awpRegistry = AWPRegistry(address(new ERC1967Proxy(
            address(awpRegistryImpl),
            abi.encodeCall(AWPRegistry.initialize, (deployer, address(treasury), guardian))
        )));
        nft = new WorknetNFT("AWP Worknet", "AWPSUB", address(awpRegistry));
        lp = new MockLPManager(address(awpRegistry), address(awp));

        // Deploy AWPEmission (UUPS proxy) — deployer == guardian in tests
        AWPEmission emissionImpl = new AWPEmission();
        bytes memory emissionInitData = abi.encodeCall(
            AWPEmission.initialize,
            (address(awp), deployer, INITIAL_DAILY, block.timestamp, EPOCH, address(treasury))
        );
        ERC1967Proxy emissionProxy = new ERC1967Proxy(address(emissionImpl), emissionInitData);
        emission = AWPEmission(address(emissionProxy));

        awp.addMinter(address(emission));
        awp.renounceAdmin();
        factory.setAddresses(address(awpRegistry));

        // Deploy StakingVault + StakeNFT
        vault = StakingVault(address(new ERC1967Proxy(
            address(new StakingVault()), abi.encodeCall(StakingVault.initialize, (address(awpRegistry), deployer))
        )));
        stakeNFT = new StakeNFT(address(awp), address(vault), address(awpRegistry));

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
        treasury.grantRole(treasury.DEFAULT_ADMIN_ROLE(), guardian);
        treasury.renounceRole(treasury.DEFAULT_ADMIN_ROLE(), deployer);

        // Initialize registry (no accessManager)
        awpRegistry.initializeRegistry(
            address(awp), address(nft), address(factory), address(emission),
            address(lp), address(vault), address(stakeNFT), address(0), ""
        );

        // Distribute tokens from deployer
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
        awpRegistry.setRecipient(user);
    }

    function _bindAgent(address agent, address target) internal {
        vm.prank(agent);
        awpRegistry.bind(target);
    }

    function _registerWorknet(address owner, address sc) internal returns (uint256) {
        vm.startPrank(owner);
        awp.approve(address(awpRegistry), LP_COST);
        uint256 id = awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams("Worknet", "SUB", sc, bytes32(0), 0, "")
        );
        vm.stopPrank();
        return id;
    }

    /// @dev Deposit AWP via StakeNFT and allocate via AWPRegistry (explicit staker)
    function _depositAndAllocate(address staker, address agent, uint256 worknetId, uint256 deposit, uint256 alloc)
        internal
    {
        vm.startPrank(staker);
        awp.approve(address(stakeNFT), deposit);
        stakeNFT.deposit(deposit, 52 weeks);
        vault.allocate(staker, agent, worknetId, alloc);
        vm.stopPrank();
    }

    function _settleEpoch() internal {
        vm.warp(block.timestamp + EPOCH + 1);
        emission.settleEpoch(200);
    }

    function _packArray(address[] memory addrs, uint96[] memory ws) internal pure returns (uint256[] memory) {
        uint256[] memory packed = new uint256[](addrs.length);
        for (uint256 i = 0; i < addrs.length; i++) {
            packed[i] = (uint256(ws[i]) << 160) | uint256(uint160(addrs[i]));
        }
        return packed;
    }

    function _submitWeights(address[] memory recipients, uint96[] memory weights) internal {
        uint256 tw = 0;
        for (uint256 i = 0; i < weights.length; i++) tw += weights[i];
        uint256 effectiveEpoch = emission.settledEpoch();
        vm.prank(deployer); // deployer is AWPEmission guardian in E2E tests
        emission.submitAllocations(_packArray(recipients, weights), tw, effectiveEpoch);
    }

    function _submitWeight(address _recipient, uint96 weight) internal {
        address[] memory addrs = new address[](1);
        addrs[0] = _recipient;
        uint96[] memory ws = new uint96[](1);
        ws[0] = weight;
        _submitWeights(addrs, ws);
    }

    // ════════════════════════════════════════════
    //  E2E 1: Worknet NFT transfer -> ownership change
    // ════════════════════════════════════════════

    function test_e2e_worknetNFTTransfer() public {
        _registerUser(alice);
        uint256 sid = _registerWorknet(alice, worknetC1);

        vm.prank(alice);
        awpRegistry.activateWorknet(sid);

        vm.prank(alice);
        nft.transferFrom(alice, bob, sid);
        assertEq(nft.ownerOf(sid), bob);

        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotOwner.selector);
        awpRegistry.pauseWorknet(sid);

        vm.prank(bob);
        awpRegistry.pauseWorknet(sid);
        assertFalse(awpRegistry.isWorknetActive(sid));

        vm.prank(bob);
        awpRegistry.resumeWorknet(sid);
        assertTrue(awpRegistry.isWorknetActive(sid));
    }

    // ════════════════════════════════════════════
    //  E2E 2: Reallocate is immediate
    // ════════════════════════════════════════════

    function test_e2e_reallocateImmediate() public {
        _registerUser(alice);
        _bindAgent(agentA, alice);
        _bindAgent(agentB, alice);
        uint256 sid = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(sid);

        _depositAndAllocate(alice, agentA, sid, 10_000 * 1e18, 8_000 * 1e18);

        // Reallocate 5000 from agentA -> agentB (immediate)
        vm.prank(alice);
        vault.reallocate(alice, agentA, sid, agentB, sid, 5_000 * 1e18);

        assertEq(vault.getAgentStake(alice, agentA, sid), 3_000 * 1e18);
        assertEq(vault.getAgentStake(alice, agentB, sid), 5_000 * 1e18);
        assertEq(vault.userTotalAllocated(alice), 8_000 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 3: Full worknet lifecycle
    // ════════════════════════════════════════════

    function test_e2e_worknetFullLifecycle() public {
        _registerUser(alice);
        uint256 sid = _registerWorknet(alice, worknetC1);

        vm.prank(alice);
        awpRegistry.activateWorknet(sid);
        assertEq(uint256(awpRegistry.getWorknet(sid).status), uint256(IAWPRegistry.WorknetStatus.Active));

        vm.prank(alice);
        awpRegistry.pauseWorknet(sid);
        assertEq(uint256(awpRegistry.getWorknet(sid).status), uint256(IAWPRegistry.WorknetStatus.Paused));

        vm.prank(alice);
        awpRegistry.resumeWorknet(sid);
        assertEq(uint256(awpRegistry.getWorknet(sid).status), uint256(IAWPRegistry.WorknetStatus.Active));

        vm.prank(guardian);
        awpRegistry.banWorknet(sid);
        assertEq(uint256(awpRegistry.getWorknet(sid).status), uint256(IAWPRegistry.WorknetStatus.Banned));

        AlphaToken alpha = AlphaToken(awpRegistry.getWorknetFull(sid).alphaToken);
        assertTrue(alpha.minterPaused(worknetC1));

        vm.prank(guardian);
        awpRegistry.unbanWorknet(sid);
        assertEq(uint256(awpRegistry.getWorknet(sid).status), uint256(IAWPRegistry.WorknetStatus.Active));
        assertFalse(alpha.minterPaused(worknetC1));

        // Must ban again before deregister (deregister requires Banned status)
        vm.prank(guardian);
        awpRegistry.banWorknet(sid);

        vm.warp(block.timestamp + 31 days);
        vm.prank(guardian);
        awpRegistry.deregisterWorknet(sid);

        vm.expectRevert();
        nft.ownerOf(sid);
        vm.expectRevert();
        awpRegistry.getWorknetFull(sid);
    }

    // ════════════════════════════════════════════
    //  E2E 4: DAO governance with NFT-based voting
    // ════════════════════════════════════════════

    function test_e2e_daoGovernanceWeight() public {
        _registerUser(alice);
        uint256 sid = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(sid);

        uint256 stakeAmount = 40_000_000 * 1e18;
        vm.startPrank(alice);
        awp.approve(address(stakeNFT), stakeAmount);
        uint256 tokenId = stakeNFT.deposit(stakeAmount, 52 weeks);
        vm.stopPrank();

        _submitWeight(worknetC1, uint96(100));
        _settleEpoch();
        _settleEpoch();

        vm.roll(block.number + 1);

        // Guardian directly sets initialAlphaPrice (was onlyTimelock, now onlyGuardian)
        vm.prank(guardian);
        awpRegistry.setInitialAlphaPrice(42e18);
        assertEq(awpRegistry.initialAlphaPrice(), 42e18);

        // Verify stake exists
        assertGt(awp.balanceOf(address(stakeNFT)), 0);
    }

    // ════════════════════════════════════════════
    //  E2E 5: Multi-epoch emission decay
    // ════════════════════════════════════════════

    function test_e2e_emissionDecayMultiEpoch() public {
        _registerUser(alice);
        uint256 sid = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(sid);

        _submitWeight(worknetC1, uint96(100));

        uint256 prevEmission = INITIAL_DAILY;

        for (uint256 i = 0; i < 10; i++) {
            if (i > 0) {
                _submitWeight(worknetC1, uint96(100));
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
    //  E2E 6: Multi-user multi-worknet concurrency
    // ════════════════════════════════════════════

    function test_e2e_multiUserMultiWorknet() public {
        _registerUser(alice);
        _registerUser(bob);
        _registerUser(charlie);

        uint256 sid1 = _registerWorknet(alice, worknetC1);
        uint256 sid2 = _registerWorknet(alice, worknetC2);
        vm.startPrank(alice);
        awpRegistry.activateWorknet(sid1);
        awpRegistry.activateWorknet(sid2);
        vm.stopPrank();

        {
            address[] memory addrs = new address[](2);
            addrs[0] = worknetC1;
            addrs[1] = worknetC2;
            uint96[] memory ws = new uint96[](2);
            ws[0] = 200;
            ws[1] = 100;
            _submitWeights(addrs, ws);
        }

        _depositAndAllocate(alice, agentA, sid1, 100_000 * 1e18, 80_000 * 1e18);
        _depositAndAllocate(bob, agentB, sid1, 50_000 * 1e18, 50_000 * 1e18);
        _depositAndAllocate(charlie, agentC, sid2, 30_000 * 1e18, 30_000 * 1e18);

        assertEq(vault.getWorknetTotalStake(sid1), 130_000 * 1e18);
        assertEq(vault.getWorknetTotalStake(sid2), 30_000 * 1e18);

        _settleEpoch();
        _settleEpoch();

        uint256 bal1 = awp.balanceOf(worknetC1);
        uint256 bal2 = awp.balanceOf(worknetC2);
        assertApproxEqRel(bal1, bal2 * 2, 0.01e18);
    }

    // ════════════════════════════════════════════
    //  E2E 7: Binding tree (replaces mutual exclusion)
    // ════════════════════════════════════════════

    function test_e2e_bindingTree() public {
        // Build chain: agentA -> alice, agentB -> agentA
        _bindAgent(agentA, alice);
        _bindAgent(agentB, agentA);

        assertEq(awpRegistry.boundTo(agentA), alice);
        assertEq(awpRegistry.boundTo(agentB), agentA);

        // resolveRecipient should walk to alice
        vm.prank(alice);
        awpRegistry.setRecipient(charlie);
        assertEq(awpRegistry.resolveRecipient(agentB), charlie);
    }

    // ════════════════════════════════════════════
    //  E2E 8: Emission clamp near MAX_SUPPLY
    // ════════════════════════════════════════════

    function test_e2e_emissionClampNearMaxSupply() public {
        _registerUser(alice);
        uint256 sid = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(sid);
        _submitWeight(worknetC1, uint96(100));

        _settleEpoch();
        _settleEpoch();
        uint256 afterEpoch1 = awp.totalSupply();
        assertTrue(afterEpoch1 > 200_000_000 * 1e18);

        uint256 remaining = awp.MAX_SUPPLY() - awp.totalSupply();
        assertTrue(remaining > 0);
    }

    // ════════════════════════════════════════════
    //  E2E 9: notSettling guard
    // ════════════════════════════════════════════

    function test_e2e_notSettlingGuard() public {
        _registerUser(alice);
        uint256 sid = _registerWorknet(alice, worknetC1);
        uint256 sid2 = _registerWorknet(alice, worknetC2);
        vm.startPrank(alice);
        awpRegistry.activateWorknet(sid);
        awpRegistry.activateWorknet(sid2);
        vm.stopPrank();

        {
            address[] memory addrs = new address[](2);
            addrs[0] = worknetC1;
            addrs[1] = worknetC2;
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

        uint256 sid3 = _registerWorknet(alice, worknetC3);
        vm.prank(alice);
        awpRegistry.activateWorknet(sid3);
        assertTrue(awpRegistry.isWorknetActive(sid3));
    }

    // ════════════════════════════════════════════
    //  E2E 10: Deallocate releases allocations
    // ════════════════════════════════════════════

    function test_e2e_deallocateReleasesAllocations() public {
        _registerUser(alice);
        uint256 sid = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(sid);

        _depositAndAllocate(alice, agentA, sid, 10_000 * 1e18, 5_000 * 1e18);

        // Deallocate all
        vm.prank(alice);
        vault.deallocate(alice, agentA, sid, 5_000 * 1e18);

        assertEq(vault.getAgentStake(alice, agentA, sid), 0);
        assertEq(vault.userTotalAllocated(alice), 0);
    }

    // ════════════════════════════════════════════
    //  E2E 11: Deposit + allocate flow
    // ════════════════════════════════════════════

    function test_e2e_depositAndAllocate() public {
        _registerUser(alice);
        uint256 sid = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(sid);

        vm.startPrank(bob);
        awp.approve(address(stakeNFT), 20_000 * 1e18);
        stakeNFT.deposit(20_000 * 1e18, 52 weeks);
        vault.allocate(bob, agentB, sid, 15_000 * 1e18);
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
        uint256 sid = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(sid);

        assertEq(awpRegistry.resolveRecipient(alice), alice);

        vm.prank(alice);
        awpRegistry.setRecipient(bob);
        assertEq(awpRegistry.resolveRecipient(alice), bob);

        _depositAndAllocate(alice, agentA, sid, 1_000 * 1e18, 500 * 1e18);
        AWPRegistry.AgentInfo memory info = awpRegistry.getAgentInfo(agentA, sid);
        assertEq(info.rewardRecipient, bob);
    }

    // ════════════════════════════════════════════
    //  E2E 13: Multi-epoch with stake changes
    // ════════════════════════════════════════════

    function test_e2e_multiEpochWithStakeChanges() public {
        _registerUser(alice);
        _registerUser(bob);

        uint256 sid1 = _registerWorknet(alice, worknetC1);
        uint256 sid2 = _registerWorknet(alice, worknetC2);
        vm.startPrank(alice);
        awpRegistry.activateWorknet(sid1);
        awpRegistry.activateWorknet(sid2);
        vm.stopPrank();

        {
            address[] memory addrs = new address[](2);
            addrs[0] = worknetC1;
            addrs[1] = worknetC2;
            uint96[] memory ws = new uint96[](2);
            ws[0] = 100;
            ws[1] = 100;
            _submitWeights(addrs, ws);
        }

        _depositAndAllocate(alice, agentA, sid1, 10_000 * 1e18, 8_000 * 1e18);
        _depositAndAllocate(bob, agentB, sid2, 5_000 * 1e18, 5_000 * 1e18);

        _settleEpoch();
        _settleEpoch();
        uint256 sc1Bal1 = awp.balanceOf(worknetC1);
        uint256 sc2Bal1 = awp.balanceOf(worknetC2);
        assertEq(sc1Bal1, sc2Bal1);

        // Alice deallocates
        vm.prank(alice);
        vault.deallocate(alice, agentA, sid1, 3_000 * 1e18);

        _submitWeights(
            _toArray2(worknetC1, worknetC2),
            _toUint96Array2(100, 100)
        );
        _settleEpoch();

        // Bob deallocates
        vm.prank(bob);
        vault.deallocate(bob, agentB, sid2, 5_000 * 1e18);

        _submitWeights(
            _toArray2(worknetC1, worknetC2),
            _toUint96Array2(100, 100)
        );
        _settleEpoch();
        assertEq(emission.currentEpoch(), 4);

        assertTrue(awp.balanceOf(worknetC1) > sc1Bal1);
        assertTrue(awp.balanceOf(worknetC2) > sc2Bal1);
    }

    // ════════════════════════════════════════════
    //  E2E 14: Agent unbind
    // ════════════════════════════════════════════

    function test_e2e_agentUnbind() public {
        _bindAgent(agentA, alice);
        assertEq(awpRegistry.boundTo(agentA), alice);

        vm.prank(agentA);
        awpRegistry.unbind();
        assertEq(awpRegistry.boundTo(agentA), address(0));

        // agentA can re-bind
        _bindAgent(agentA, alice);
        assertEq(awpRegistry.boundTo(agentA), alice);
    }

    // ════════════════════════════════════════════
    //  E2E 15: Alpha Token worknet contract minting
    // ════════════════════════════════════════════

    function test_e2e_alphaTokenMinting() public {
        _registerUser(alice);
        uint256 sid = _registerWorknet(alice, worknetC1);

        AlphaToken alpha = AlphaToken(awpRegistry.getWorknetFull(sid).alphaToken);
        assertFalse(alpha.minters(address(awpRegistry)));
        assertTrue(alpha.mintersLocked());

        vm.warp(block.timestamp + 10 days);
        vm.prank(worknetC1);
        alpha.mint(alice, 1_000_000 * 1e18);
        assertEq(alpha.balanceOf(alice), 1_000_000 * 1e18);

        vm.prank(alice);
        awpRegistry.activateWorknet(sid);
        vm.prank(guardian);
        awpRegistry.banWorknet(sid);

        vm.prank(worknetC1);
        vm.expectRevert(AlphaToken.MinterPaused.selector);
        alpha.mint(alice, 100);

        vm.prank(guardian);
        awpRegistry.unbanWorknet(sid);
        vm.warp(block.timestamp + 1 days);
        vm.prank(worknetC1);
        alpha.mint(alice, 500_000 * 1e18);
        assertEq(alpha.balanceOf(alice), 1_500_000 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 16: Guardian emergency pause
    // ════════════════════════════════════════════

    function test_e2e_emergencyPause() public {
        _registerUser(alice);
        uint256 sid = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(sid);
        _depositAndAllocate(alice, agentA, sid, 10_000 * 1e18, 5_000 * 1e18);

        vm.prank(guardian);
        awpRegistry.pause();

        // AWPRegistry operations blocked
        vm.startPrank(alice);
        vm.expectRevert();
        awpRegistry.setRecipient(alice);
        vm.expectRevert();
        awpRegistry.activateWorknet(sid);
        vm.stopPrank();

        // StakingVault allocate/deallocate are NOT gated by AWPRegistry pause
        // (they live on StakingVault which has no Pausable)

        // Emission unaffected
        vm.warp(block.timestamp + EPOCH + 1);
        emission.settleEpoch(200);
        assertEq(emission.settledEpoch(), 1);

        // Unpause
        vm.prank(guardian);
        awpRegistry.unpause();

        vm.prank(alice);
        vault.deallocate(alice, agentA, sid, 1_000 * 1e18);
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

        uint256 sid = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(sid);

        _depositAndAllocate(alice, agentA, sid, 5_000 * 1e18, 3_000 * 1e18);
        _depositAndAllocate(bob, agentB, sid, 2_000 * 1e18, 2_000 * 1e18);

        vm.prank(alice);
        awpRegistry.setRecipient(charlie);

        address[] memory agents = new address[](3);
        agents[0] = agentA;
        agents[1] = agentB;
        agents[2] = address(0x9999);

        AWPRegistry.AgentInfo[] memory infos = awpRegistry.getAgentsInfo(agents, sid);

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
        uint256 nonce = awpRegistry.nonces(agent);

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Bind(address agent,address target,uint256 nonce,uint256 deadline)"),
            agent, alice, nonce, deadline
        ));
        bytes32 digest = _getDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(agentPk, digest);

        address relayer = address(0x9999);
        vm.prank(relayer);
        awpRegistry.bindFor(agent, alice, deadline, v, r, s);

        assertEq(awpRegistry.boundTo(agent), alice);
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
        awpRegistry.bindFor(agent, alice, deadline, v, r, s);
    }

    // ════════════════════════════════════════════
    //  E2E 20: Multi-user same worknet reallocate
    // ════════════════════════════════════════════

    function test_e2e_multiUserSameWorknetReallocate() public {
        _registerUser(alice);
        _registerUser(bob);

        uint256 sid1 = _registerWorknet(alice, worknetC1);
        uint256 sid2 = _registerWorknet(alice, worknetC2);
        vm.startPrank(alice);
        awpRegistry.activateWorknet(sid1);
        awpRegistry.activateWorknet(sid2);
        vm.stopPrank();

        _depositAndAllocate(alice, agentA, sid1, 1_000 * 1e18, 500 * 1e18);
        _depositAndAllocate(bob, agentB, sid1, 1_000 * 1e18, 300 * 1e18);

        // Reallocate
        vm.prank(alice);
        vault.reallocate(alice, agentA, sid1, agentA, sid2, 200 * 1e18);
        vm.prank(bob);
        vault.reallocate(bob, agentB, sid1, agentB, sid2, 100 * 1e18);

        assertEq(vault.getAgentStake(alice, agentA, sid1), 300 * 1e18);
        assertEq(vault.getAgentStake(alice, agentA, sid2), 200 * 1e18);
        assertEq(vault.getAgentStake(bob, agentB, sid1), 200 * 1e18);
        assertEq(vault.getAgentStake(bob, agentB, sid2), 100 * 1e18);

        assertEq(vault.getWorknetTotalStake(sid1), 500 * 1e18);
        assertEq(vault.getWorknetTotalStake(sid2), 300 * 1e18);
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
        awp.approve(address(awpRegistry), LP_COST * 5);
        for (uint256 i = 0; i < 5; i++) {
            uint256 sid = awpRegistry.registerWorknet(
                IAWPRegistry.WorknetParams("S", "S", scs[i], bytes32(0), 0, "")
            );
            awpRegistry.activateWorknet(sid);
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
        uint256 sid = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(sid);

        _depositAndAllocate(alice, agentA, sid, 10_000 * 1e18, 5_000 * 1e18);

        // Grant delegate to bob
        vm.prank(alice);
        awpRegistry.grantDelegate(bob);

        // Bob can deallocate on behalf of alice
        vm.prank(bob);
        vault.deallocate(alice, agentA, sid, 2_000 * 1e18);
        assertEq(vault.getAgentStake(alice, agentA, sid), 3_000 * 1e18);

        // Revoke delegation
        vm.prank(alice);
        awpRegistry.revokeDelegate(bob);

        // Bob can no longer operate
        vm.prank(bob);
        vm.expectRevert(StakingVault.NotAuthorized.selector);
        vault.deallocate(alice, agentA, sid, 1_000 * 1e18);
    }

    // ── EIP-712 helpers ──

    function _getDigest(bytes32 structHash) internal view returns (bytes32) {
        bytes32 domainSeparator = keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
                keccak256("AWPRegistry"),
                keccak256("1"),
                block.chainid,
                address(awpRegistry)
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

    // ════════════════════════════════════════════
    //  E2E 23: Cross-chain allocate (local + foreign worknetId)
    // ════════════════════════════════════════════

    function test_e2e_crossChainAllocate() public {
        _registerUser(alice);
        uint256 localWorknet = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(localWorknet);

        // alice stakes and allocates to local worknet
        vm.startPrank(alice);
        awp.approve(address(stakeNFT), 10_000 * 1e18);
        stakeNFT.deposit(10_000 * 1e18, 52 weeks);
        vault.allocate(alice, agentA, localWorknet, 3_000 * 1e18);

        // Allocate to "foreign" worknet (different chain's worknetId)
        uint256 foreignWorknet = (uint256(42161) << 64) | 99;
        vault.allocate(alice, agentA, foreignWorknet, 2_000 * 1e18);
        vm.stopPrank();

        // Verify both allocations recorded
        assertEq(vault.getAgentStake(alice, agentA, localWorknet), 3_000 * 1e18);
        assertEq(vault.getAgentStake(alice, agentA, foreignWorknet), 2_000 * 1e18);

        // Total allocated = 5000, total staked = 10000, so unallocated = 5000
        assertEq(vault.userTotalAllocated(alice), 5_000 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 24: Guardian emission flow (deploy → submit → settle → verify 100%)
    // ════════════════════════════════════════════

    function test_e2e_guardianEmissionFlow() public {
        _registerUser(alice);
        uint256 sid = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(sid);

        // Guardian submits weights for worknetC1
        _submitWeight(worknetC1, uint96(500));

        // Settle epoch 0 (weights promoted from epoch 1)
        _settleEpoch();

        // Settle epoch 1
        _settleEpoch();

        // worknetC1 receives 100% of emission (no DAO split)
        uint256 epoch0Pool = INITIAL_DAILY;
        uint256 epoch1Pool = INITIAL_DAILY * 996844 / 1000000;
        assertEq(awp.balanceOf(worknetC1), epoch0Pool + epoch1Pool);
        // Treasury gets nothing unless Guardian includes it
        assertEq(awp.balanceOf(address(treasury)), 50_000_000 * 1e18); // unchanged from deploy

        // Now Guardian includes treasury as a recipient for epoch 3
        {
            address[] memory addrs = new address[](2);
            addrs[0] = worknetC1;
            addrs[1] = address(treasury);
            uint96[] memory ws = new uint96[](2);
            ws[0] = 700;
            ws[1] = 300;
            _submitWeights(addrs, ws);
        }

        _settleEpoch();

        // Treasury should now have received its weight share
        uint256 treasuryBal = awp.balanceOf(address(treasury));
        assertTrue(treasuryBal > 50_000_000 * 1e18); // more than initial distribution
    }

    // ════════════════════════════════════════════
    //  E2E 25: Cross-chain worknetId encoding
    // ════════════════════════════════════════════

    function test_e2e_crossChainWorknetId() public {
        _registerUser(alice);

        // Register a worknet and verify its worknetId encoding
        uint256 sid = _registerWorknet(alice, worknetC1);

        // Verify encoding: (block.chainid << 64) | localCounter
        uint256 encodedChainId = sid >> 64;
        uint256 localId = sid & ((1 << 64) - 1);
        assertEq(encodedChainId, block.chainid);
        assertEq(localId, 1); // first worknet registered

        // Register another — localId should be 2
        uint256 sid2 = _registerWorknet(alice, worknetC2);
        uint256 localId2 = sid2 & ((1 << 64) - 1);
        assertEq(localId2, 2);

        // Allocate to the registered worknet
        vm.prank(alice);
        awpRegistry.activateWorknet(sid);
        vm.startPrank(alice);
        awp.approve(address(stakeNFT), 5_000 * 1e18);
        stakeNFT.deposit(5_000 * 1e18, 52 weeks);
        vault.allocate(alice, agentA, sid, 2_000 * 1e18);
        vm.stopPrank();

        assertEq(vault.getAgentStake(alice, agentA, sid), 2_000 * 1e18);
    }

    // ════════════════════════════════════════════
    //  E2E 26: Multi-epoch with Guardian resubmit
    // ════════════════════════════════════════════

    function test_e2e_multiEpochWithGuardianResubmit() public {
        _registerUser(alice);
        uint256 sid1 = _registerWorknet(alice, worknetC1);
        uint256 sid2 = _registerWorknet(alice, worknetC2);
        vm.startPrank(alice);
        awpRegistry.activateWorknet(sid1);
        awpRegistry.activateWorknet(sid2);
        vm.stopPrank();

        // Guardian submits epoch 1 weights: worknetC1=800, worknetC2=200
        {
            address[] memory addrs = new address[](2);
            addrs[0] = worknetC1;
            addrs[1] = worknetC2;
            uint96[] memory ws = new uint96[](2);
            ws[0] = 800;
            ws[1] = 200;
            _submitWeights(addrs, ws);
        }

        // Settle epoch 0 + epoch 1
        _settleEpoch();
        _settleEpoch();

        uint256 sc1BalEpoch1 = awp.balanceOf(worknetC1);
        uint256 sc2BalEpoch1 = awp.balanceOf(worknetC2);

        // 4:1 ratio after epochs 0+1
        assertApproxEqRel(sc1BalEpoch1, sc2BalEpoch1 * 4, 0.01e18);

        // Guardian resubmits for epoch 3: now worknetC1=100, worknetC2=900
        {
            address[] memory addrs = new address[](2);
            addrs[0] = worknetC1;
            addrs[1] = worknetC2;
            uint96[] memory ws = new uint96[](2);
            ws[0] = 100;
            ws[1] = 900;
            _submitWeights(addrs, ws);
        }

        // Settle epoch 2
        _settleEpoch();

        uint256 sc1BalEpoch2 = awp.balanceOf(worknetC1);
        uint256 sc2BalEpoch2 = awp.balanceOf(worknetC2);

        // Epoch 2 incremental emission: worknetC2 got 9x more than worknetC1 in this epoch
        uint256 sc1Epoch2Gain = sc1BalEpoch2 - sc1BalEpoch1;
        uint256 sc2Epoch2Gain = sc2BalEpoch2 - sc2BalEpoch1;
        assertApproxEqRel(sc2Epoch2Gain, sc1Epoch2Gain * 9, 0.01e18);
    }
}
