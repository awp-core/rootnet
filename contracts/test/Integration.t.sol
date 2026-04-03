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
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @title Integration — Full deployment + registration + staking + emission flow
contract IntegrationTest is Test {
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
    address airdrop = address(0xF4);

    address owner1 = address(0x101);
    address agent1 = address(0x201);
    address agent2 = address(0x202);
    address worknetManager1 = address(0x301);
    address worknetManager2 = address(0x302);

    uint256 constant INITIAL_DAILY_EMISSION = 31_600_000 * 1e18;
    uint256 constant EPOCH_DURATION = 1 days;

    function setUp() public {
        _fullDeployment();
    }

    /// @notice Full deployment flow
    function _fullDeployment() internal {
        vm.startPrank(deployer);

        // Step 1: AWPToken
        awp = new AWPToken("AWP Token", "AWP", deployer);
        awp.initialMint(200_000_000 * 1e18);

        // Step 2: AlphaTokenFactory
        factory = new AlphaTokenFactory(deployer, 0);

        // Step 4: Treasury
        address[] memory proposers = new address[](0);
        address[] memory executors = new address[](1);
        executors[0] = address(0);
        treasury = new Treasury(0, proposers, executors, deployer);

        // Step 5: AWPRegistry
        AWPRegistry awpRegistryImpl = new AWPRegistry();
        awpRegistry = AWPRegistry(address(new ERC1967Proxy(
            address(awpRegistryImpl),
            abi.encodeCall(AWPRegistry.initialize, (deployer, address(treasury), guardian))
        )));

        // Step 6-7: Sub-contracts
        nft = new WorknetNFT("AWP Worknet", "AWPSUB", address(awpRegistry));
        lp = new MockLPManager(address(awpRegistry), address(awp));

        // Step 8: AWPEmission (UUPS proxy) — deployer == guardian in tests
        AWPEmission emissionImpl = new AWPEmission();
        bytes memory emissionInitData = abi.encodeCall(
            AWPEmission.initialize,
            (address(awp), deployer, INITIAL_DAILY_EMISSION, block.timestamp, EPOCH_DURATION, address(treasury))
        );
        ERC1967Proxy emissionProxy = new ERC1967Proxy(address(emissionImpl), emissionInitData);
        emission = AWPEmission(address(emissionProxy));

        // Step 9: Add minter
        awp.addMinter(address(emission));

        // Step 10: Permanently lock the minter list
        awp.renounceAdmin();

        // Step 11: Configure factory
        factory.setAddresses(address(awpRegistry));

        // Step 12: Deploy StakingVault + StakeNFT
        vault = StakingVault(address(new ERC1967Proxy(
            address(new StakingVault()), abi.encodeCall(StakingVault.initialize, (address(awpRegistry), deployer))
        )));
        stakeNFT = new StakeNFT(address(awp), address(vault), address(awpRegistry));

        // Step 13: AWPDAO
        dao = new AWPDAO(
            address(stakeNFT),
            address(awp),
            TimelockController(payable(address(treasury))),
            1,      // votingDelay
            50400,  // votingPeriod
            4       // quorum 4%
        );

        // Step 14: Grant Treasury roles
        treasury.grantRole(treasury.PROPOSER_ROLE(), address(dao));
        treasury.grantRole(treasury.CANCELLER_ROLE(), address(dao));

        // Step 15: Transfer Treasury admin to guardian, deployer renounces
        treasury.grantRole(treasury.DEFAULT_ADMIN_ROLE(), guardian);
        treasury.renounceRole(treasury.DEFAULT_ADMIN_ROLE(), deployer);

        // Step 16: Initialize registry (no accessManager)
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

        // Distribute tokens
        awp.transfer(address(treasury), 90_000_000 * 1e18);
        awp.transfer(airdrop, 100_000_000 * 1e18);

        vm.stopPrank();

        // Verify post-deployment state
        assertEq(awp.balanceOf(deployer), 10_000_000 * 1e18);
        assertEq(awp.balanceOf(address(treasury)), 90_000_000 * 1e18);
        assertTrue(awp.minters(address(emission)));
        assertFalse(awp.minters(deployer));
        assertTrue(awpRegistry.registryInitialized());
    }

    function _settleOneEpoch() internal {
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);
    }

    // ── Guardian submission helpers ──

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
        vm.prank(deployer); // deployer == guardian in tests
        emission.submitAllocations(_packArray(recipients, weights), tw, effectiveEpoch);
    }

    function _submitWeight(address _recipient, uint96 weight) internal {
        address[] memory addrs = new address[](1);
        addrs[0] = _recipient;
        uint96[] memory ws = new uint96[](1);
        ws[0] = weight;
        _submitWeights(addrs, ws);
    }

    /// @notice Full flow: register -> bind -> register worknet -> stake -> allocate -> emission
    function test_fullFlow() public {
        // Give owner1 some AWP
        vm.prank(airdrop);
        awp.transfer(owner1, 2_000_000 * 1e18);

        // 1. Register user
        vm.prank(owner1);
        awpRegistry.setRecipient(owner1);
        assertTrue(awpRegistry.isRegistered(owner1));

        // 2. Bind Agent
        vm.prank(agent1);
        awpRegistry.bind(owner1);
        assertEq(awpRegistry.boundTo(agent1), owner1);

        // 3. Register worknet
        vm.startPrank(owner1);
        awp.approve(address(awpRegistry), 1_000_000 * 1e18);
        uint256 worknetId = awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams({
                name: "TestWorknet",
                symbol: "TSUB",
                worknetManager: worknetManager1,
                salt: bytes32(0),
                minStake: 0,
                skillsURI: ""
            })
        );
        vm.stopPrank();
        assertEq(worknetId & ((1 << 64) - 1), 1);
        assertEq(worknetId >> 64, block.chainid);

        // Verify Alpha Token
        AlphaToken alpha = AlphaToken(awpRegistry.getWorknetFull(worknetId).alphaToken);
        assertTrue(alpha.mintersLocked());
        assertTrue(alpha.minters(worknetManager1));

        // 4. Activate worknet
        vm.prank(owner1);
        awpRegistry.activateWorknet(worknetId);
        assertTrue(awpRegistry.isWorknetActive(worknetId));

        // 5. Set governance weight for epoch 1
        _submitWeight(worknetManager1, uint96(1000));

        // 6. Stake via StakeNFT
        vm.startPrank(owner1);
        awp.approve(address(stakeNFT), 500_000 * 1e18);
        stakeNFT.deposit(500_000 * 1e18, 52 weeks);

        // 7. Allocate to agent1/worknet1 (explicit staker)
        vault.allocate(owner1, agent1, worknetId, 300_000 * 1e18);
        vm.stopPrank();

        assertEq(vault.getAgentStake(owner1, agent1, worknetId), 300_000 * 1e18);
        assertEq(vault.getWorknetTotalStake(worknetId), 300_000 * 1e18);

        // 8. Query Agent info
        AWPRegistry.AgentInfo memory agentInfo = awpRegistry.getAgentInfo(agent1, worknetId);
        assertEq(agentInfo.root, owner1);
        assertTrue(agentInfo.isValid);
        assertEq(agentInfo.stake, 300_000 * 1e18);

        // 9. Emission — settle epoch 0
        _settleOneEpoch();

        // 10. Settle epoch 1
        _settleOneEpoch();

        uint256 worknetBal = awp.balanceOf(worknetManager1);
        // Epoch 0 (no decay, 100% to recipients) + Epoch 1 (decayed, 100%)
        uint256 epoch0Worknet = INITIAL_DAILY_EMISSION;
        uint256 decayedEmission = INITIAL_DAILY_EMISSION * 996844 / 1000000;
        uint256 epoch1Worknet = decayedEmission;
        assertEq(worknetBal, epoch0Worknet + epoch1Worknet);

        // 11. Third epoch
        _submitWeight(worknetManager1, uint96(1000));
        _settleOneEpoch();
    }

    /// @notice Test multi-worknet emission distribution
    function test_multiWorknetEmission() public {
        vm.prank(airdrop);
        awp.transfer(owner1, 3_000_000 * 1e18);

        vm.prank(owner1);
        awpRegistry.setRecipient(owner1);

        vm.startPrank(owner1);
        awp.approve(address(awpRegistry), 2_000_000 * 1e18);

        uint256 worknet1 = awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams("Sub1", "S1", worknetManager1, bytes32(0), 0, "")
        );
        uint256 worknet2 = awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams("Sub2", "S2", worknetManager2, bytes32(0), 0, "")
        );

        awpRegistry.activateWorknet(worknet1);
        awpRegistry.activateWorknet(worknet2);
        vm.stopPrank();

        {
            address[] memory addrs = new address[](2);
            addrs[0] = worknetManager1;
            addrs[1] = worknetManager2;
            uint96[] memory ws = new uint96[](2);
            ws[0] = 300;
            ws[1] = 100;
            _submitWeights(addrs, ws);
        }

        _settleOneEpoch();
        _settleOneEpoch();

        uint256 bal1 = awp.balanceOf(worknetManager1);
        uint256 bal2 = awp.balanceOf(worknetManager2);
        assertApproxEqRel(bal1, bal2 * 3, 0.01e18);
    }

    /// @notice Test stake via StakeNFT + withdraw flow
    function test_stakeWithdrawFlow() public {
        vm.prank(airdrop);
        awp.transfer(owner1, 1_000_000 * 1e18);

        vm.startPrank(owner1);
        awpRegistry.setRecipient(owner1);

        awp.approve(address(stakeNFT), 500_000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(500_000 * 1e18, 1 days);
        vm.stopPrank();

        assertEq(stakeNFT.getUserTotalStaked(owner1), 500_000 * 1e18);

        vm.prank(owner1);
        vm.expectRevert(StakeNFT.LockNotExpired.selector);
        stakeNFT.withdraw(tokenId);

        _settleOneEpoch();
        _settleOneEpoch();

        vm.prank(owner1);
        stakeNFT.withdraw(tokenId);

        assertEq(awp.balanceOf(owner1), 1_000_000 * 1e18);
        assertEq(stakeNFT.getUserTotalStaked(owner1), 0);
    }

    /// @notice Test deallocate flow (replaces agent removal freeze test)
    function test_deallocateFreesStake() public {
        vm.prank(airdrop);
        awp.transfer(owner1, 2_000_000 * 1e18);

        vm.prank(owner1);
        awpRegistry.setRecipient(owner1);

        vm.startPrank(owner1);
        awp.approve(address(awpRegistry), 1_000_000 * 1e18);
        uint256 worknetId = awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams("Sub", "SUB", worknetManager1, bytes32(0), 0, "")
        );
        awpRegistry.activateWorknet(worknetId);
        awp.approve(address(stakeNFT), 500_000 * 1e18);
        stakeNFT.deposit(500_000 * 1e18, 52 weeks);
        vault.allocate(owner1, agent1, worknetId, 300_000 * 1e18);

        // Deallocate all
        vault.deallocate(owner1, agent1, worknetId, 300_000 * 1e18);
        vm.stopPrank();

        assertEq(vault.getAgentStake(owner1, agent1, worknetId), 0);
        assertEq(vault.userTotalAllocated(owner1), 0);
    }

    /// @notice Test delegate operations
    function test_delegateOperations() public {
        vm.prank(airdrop);
        awp.transfer(owner1, 2_000_000 * 1e18);

        vm.prank(owner1);
        awpRegistry.setRecipient(owner1);

        // Register worknet
        vm.startPrank(owner1);
        awp.approve(address(awpRegistry), 1_000_000 * 1e18);
        uint256 worknetId = awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams("Sub", "SUB", worknetManager1, bytes32(0), 0, "")
        );
        awpRegistry.activateWorknet(worknetId);
        awp.approve(address(stakeNFT), 500_000 * 1e18);
        stakeNFT.deposit(500_000 * 1e18, 52 weeks);
        // Grant delegate to agent1
        awpRegistry.grantDelegate(agent1);
        vm.stopPrank();

        // agent1 allocates on behalf of owner1
        vm.prank(agent1);
        vault.allocate(owner1, agent2, worknetId, 100_000 * 1e18);
        assertEq(vault.getAgentStake(owner1, agent2, worknetId), 100_000 * 1e18);

        // agent1 deallocates on behalf of owner1
        vm.prank(agent1);
        vault.deallocate(owner1, agent2, worknetId, 50_000 * 1e18);
        assertEq(vault.getAgentStake(owner1, agent2, worknetId), 50_000 * 1e18);
    }

    /// @notice Test batch emission (many worknets)
    function test_batchSettle() public {
        vm.prank(airdrop);
        awp.transfer(owner1, 10_000_000 * 1e18);

        vm.prank(owner1);
        awpRegistry.setRecipient(owner1);

        vm.startPrank(owner1);
        awp.approve(address(awpRegistry), 10_000_000 * 1e18);

        for (uint256 i = 0; i < 3; i++) {
            address sc = address(uint160(0x400 + i));
            uint256 sid = awpRegistry.registerWorknet(IAWPRegistry.WorknetParams("Sub", "SUB", sc, bytes32(0), 0, ""));
            awpRegistry.activateWorknet(sid);
        }
        vm.stopPrank();

        {
            address[] memory addrs = new address[](3);
            addrs[0] = address(uint160(0x400));
            addrs[1] = address(uint160(0x401));
            addrs[2] = address(uint160(0x402));
            uint96[] memory ws = new uint96[](3);
            ws[0] = 100;
            ws[1] = 200;
            ws[2] = 300;
            _submitWeights(addrs, ws);
        }

        _settleOneEpoch();
        _settleOneEpoch();

        // Epoch 0 (weights promoted from epoch 1, 100%) + Epoch 1 (decayed, 100%)
        uint256 epoch0Pool = INITIAL_DAILY_EMISSION;
        uint256 decayedEmission = INITIAL_DAILY_EMISSION * 996844 / 1000000;
        uint256 epoch1Pool = decayedEmission;
        uint256 totalPool = epoch0Pool + epoch1Pool;
        uint256 expected1 = totalPool * 100 / 600;
        uint256 expected2 = totalPool * 200 / 600;
        uint256 expected3 = totalPool * 300 / 600;

        // Allow 1 wei rounding tolerance from integer division across two epochs
        assertApproxEqAbs(awp.balanceOf(address(uint160(0x400))), expected1, 1);
        assertApproxEqAbs(awp.balanceOf(address(uint160(0x401))), expected2, 1);
        assertApproxEqAbs(awp.balanceOf(address(uint160(0x402))), expected3, 1);
    }

    /// @notice Test pause protection
    function test_pauseProtection() public {
        vm.prank(guardian);
        awpRegistry.pause();

        vm.prank(owner1);
        vm.expectRevert();
        awpRegistry.setRecipient(owner1);

        vm.prank(guardian);
        awpRegistry.unpause();

        vm.prank(owner1);
        awpRegistry.setRecipient(owner1);
    }

    // ══════════════════════════════════════════════
    // New integration tests: Guardian/Timelock access control, treasury as recipient
    // ══════════════════════════════════════════════

    /// @notice Timelock cannot call AWPEmission.setGuardian (it's onlyGuardian)
    function test_guardianCannotBeChangedByTimelock() public {
        // deployer is the guardian in tests, not treasury/Timelock
        // Timelock should NOT be able to change guardian
        vm.prank(address(treasury));
        vm.expectRevert(AWPEmission.NotGuardian.selector);
        emission.setGuardian(address(0x999));
    }

    /// @notice Timelock cannot upgrade AWPEmission (it's onlyGuardian)
    function test_timelockCannotUpgradeEmission() public {
        AWPEmission newImpl = new AWPEmission();
        vm.prank(address(treasury));
        vm.expectRevert(AWPEmission.NotGuardian.selector);
        emission.upgradeToAndCall(address(newImpl), "");
    }

    /// @notice Full flow with treasury as a recipient in emission
    function test_fullFlowWithTreasuryAsRecipient() public {
        vm.prank(airdrop);
        awp.transfer(owner1, 2_000_000 * 1e18);

        // Register user and worknet
        vm.prank(owner1);
        awpRegistry.setRecipient(owner1);

        vm.startPrank(owner1);
        awp.approve(address(awpRegistry), 1_000_000 * 1e18);
        uint256 worknetId = awpRegistry.registerWorknet(
            IAWPRegistry.WorknetParams("TestWorknet", "TSUB", worknetManager1, bytes32(0), 0, "")
        );
        awpRegistry.activateWorknet(worknetId);
        vm.stopPrank();

        // Guardian includes treasury as a recipient: worknetManager1=700, treasury=300
        {
            address[] memory addrs = new address[](2);
            addrs[0] = worknetManager1;
            addrs[1] = address(treasury);
            uint96[] memory ws = new uint96[](2);
            ws[0] = 700;
            ws[1] = 300;
            uint256 tw = 0;
            for (uint256 i = 0; i < ws.length; i++) tw += ws[i];
            uint256 effectiveEpoch = emission.settledEpoch();
            vm.prank(deployer); // deployer == guardian in tests
            emission.submitAllocations(_packArray(addrs, ws), tw, effectiveEpoch);
        }

        uint256 treasuryBefore = awp.balanceOf(address(treasury));

        // Settle epoch 0 + epoch 1
        _settleOneEpoch();
        _settleOneEpoch();

        uint256 treasuryAfter = awp.balanceOf(address(treasury));
        uint256 worknetBal = awp.balanceOf(worknetManager1);

        // Treasury received emission share
        uint256 treasuryGain = treasuryAfter - treasuryBefore;
        assertTrue(treasuryGain > 0);

        // 7:3 ratio (worknetManager1:treasury)
        assertApproxEqRel(worknetBal, treasuryGain * 7 / 3, 0.01e18);
    }
}
