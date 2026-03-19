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
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {EmissionSigningHelper} from "./helpers/EmissionSigningHelper.sol";

/// @title Integration — Full deployment + registration + staking + emission flow
contract IntegrationTest is EmissionSigningHelper {
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
    address airdrop = address(0xF4);

    address owner1 = address(0x101);
    address agent1 = address(0x201);
    address agent2 = address(0x202);
    address subnetManager1 = address(0x301);
    address subnetManager2 = address(0x302);

    uint256 constant INITIAL_DAILY_EMISSION = 15_800_000 * 1e18;
    uint256 constant EPOCH_DURATION = 1 days;

    // Oracle private keys
    uint256 constant ORACLE_PK1 = 0xA1;
    uint256 constant ORACLE_PK2 = 0xA2;
    uint256 constant ORACLE_PK3 = 0xA3;

    function setUp() public {
        _fullDeployment();
    }

    /// @notice Full deployment flow
    function _fullDeployment() internal {
        vm.startPrank(deployer);

        // Step 1: AWPToken
        awp = new AWPToken("AWP Token", "AWP", deployer);
        assertEq(awp.totalSupply(), 200_000_000 * 1e18);
        assertEq(awp.balanceOf(deployer), 200_000_000 * 1e18);

        // Step 2: AlphaTokenFactory (CREATE2-based)
        factory = new AlphaTokenFactory(deployer, 0);

        // Step 4: Treasury (minDelay=0 to simplify tests)
        address[] memory proposers = new address[](0);
        address[] memory executors = new address[](1);
        executors[0] = address(0);
        treasury = new Treasury(0, proposers, executors, deployer);

        // Step 5: RootNet (no epochDuration — epochs now owned by AWPEmission)
        rootNet = new RootNet(deployer, address(treasury), guardian);

        // Step 6-7: Sub-contracts
        nft = new SubnetNFT("AWP Subnet", "AWPSUB", address(rootNet));
        access = new AccessManager(address(rootNet));
        lp = new MockLPManager(address(rootNet), address(awp));

        // Step 8: AWPEmission (UUPS proxy)
        AWPEmission emissionImpl = new AWPEmission();
        bytes memory emissionInitData = abi.encodeCall(
            AWPEmission.initialize,
            (address(awp), address(treasury), INITIAL_DAILY_EMISSION, block.timestamp, EPOCH_DURATION)
        );
        ERC1967Proxy emissionProxy = new ERC1967Proxy(address(emissionImpl), emissionInitData);
        emission = AWPEmission(address(emissionProxy));

        // Step 9: Add minter
        awp.addMinter(address(emission));

        // Step 10: Permanently lock the minter list
        awp.renounceAdmin();
        assertEq(awp.admin(), address(0));

        // Step 11: Configure factory
        factory.setAddresses(address(rootNet));

        // Step 12: Deploy StakingVault + StakeNFT (circular dependency)
        uint64 nonce = vm.getNonce(deployer);
        address predictedVault = vm.computeCreateAddress(deployer, nonce);
        address predictedStakeNFT = vm.computeCreateAddress(deployer, nonce + 1);

        vault = new StakingVault(address(rootNet));
        stakeNFT = new StakeNFT(address(awp), address(vault), address(rootNet));

        // Step 13: AWPDAO (no rootNet param — uses timestamps)
        dao = new AWPDAO(
            address(stakeNFT),
            address(awp),
            TimelockController(payable(address(treasury))),
            1,      // votingDelay
            50400,  // votingPeriod (~1 week)
            4       // quorum 4%
        );

        // Step 14: Grant Treasury roles
        treasury.grantRole(treasury.PROPOSER_ROLE(), address(dao));
        treasury.grantRole(treasury.CANCELLER_ROLE(), address(dao));

        // Step 15: Renounce Treasury admin role
        treasury.renounceRole(treasury.DEFAULT_ADMIN_ROLE(), deployer);

        // Step 16: Initialize registry
        rootNet.initializeRegistry(
            address(awp),
            address(nft),
            address(factory),
            address(emission),
            address(lp),
            address(access),
            address(vault),
            address(stakeNFT),
            address(0), // no default SubnetManager impl in integration tests
            "" // no dexConfig in integration tests
        );

        // Distribute tokens: Treasury 90M, Airdrop 100M (10M remains with deployer for tests)
        awp.transfer(address(treasury), 90_000_000 * 1e18);
        awp.transfer(airdrop, 100_000_000 * 1e18);

        vm.stopPrank();

        // Configure oracles
        address[] memory oracleList = new address[](3);
        oracleList[0] = vm.addr(ORACLE_PK1);
        oracleList[1] = vm.addr(ORACLE_PK2);
        oracleList[2] = vm.addr(ORACLE_PK3);
        vm.prank(address(treasury));
        emission.setOracleConfig(oracleList, 2);

        // Verify post-deployment state
        assertEq(awp.balanceOf(deployer), 10_000_000 * 1e18);
        assertEq(awp.balanceOf(address(treasury)), 90_000_000 * 1e18);
        assertTrue(awp.minters(address(emission)));
        assertFalse(awp.minters(deployer));
        assertTrue(rootNet.registryInitialized());
    }

    /// @dev Settle one epoch; isolated in its own function to prevent the via_ir
    ///      optimizer from caching block.timestamp across vm.warp boundaries.
    function _settleOneEpoch() internal {
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);
    }

    // ── Oracle signing helpers ──

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

    /// @notice Full flow: register user -> register agent -> register subnet -> stake via StakeNFT -> allocate -> emission
    function test_fullFlow() public {
        // Give owner1 some AWP
        vm.prank(airdrop);
        awp.transfer(owner1, 2_000_000 * 1e18);

        // 1. Register user
        vm.prank(owner1);
        rootNet.register();
        assertTrue(access.isRegistered(owner1));

        // 2. Register Agent
        vm.prank(agent1);
        rootNet.bind(owner1);
        assertTrue(access.isAgent(owner1, agent1));

        // 3. Register subnet
        vm.startPrank(owner1);
        awp.approve(address(rootNet), 1_000_000 * 1e18);
        uint256 subnetId = rootNet.registerSubnet(
            IRootNet.SubnetParams({
                name: "TestSubnet",
                symbol: "TSUB",
                subnetManager: subnetManager1,
                salt: bytes32(0),
                minStake: 0,
                skillsURI: ""
            })
        );
        vm.stopPrank();
        assertEq(subnetId, 1);

        // Verify Alpha Token
        IRootNet.SubnetInfo memory info = rootNet.getSubnet(subnetId);
        assertFalse(rootNet.getSubnetFull(subnetId).alphaToken == address(0));
        AlphaToken alpha = AlphaToken(rootNet.getSubnetFull(subnetId).alphaToken);
        assertTrue(alpha.mintersLocked());
        assertFalse(alpha.minters(address(rootNet)));
        assertTrue(alpha.minters(subnetManager1));

        // 4. Activate subnet
        vm.prank(owner1);
        rootNet.activateSubnet(subnetId);
        assertTrue(rootNet.isSubnetActive(subnetId));

        // 5. Set governance weight for epoch 1
        _submitWeight(subnetManager1, uint96(1000));

        // 6. Stake via StakeNFT
        vm.startPrank(owner1);
        awp.approve(address(stakeNFT), 500_000 * 1e18);
        stakeNFT.deposit(500_000 * 1e18, 52 weeks);
        assertEq(stakeNFT.getUserTotalStaked(owner1), 500_000 * 1e18);

        // 7. Allocate to agent1/subnet1
        rootNet.allocate(agent1, subnetId, 300_000 * 1e18);
        vm.stopPrank();

        assertEq(vault.getAgentStake(owner1, agent1, subnetId), 300_000 * 1e18);
        assertEq(vault.getSubnetTotalStake(subnetId), 300_000 * 1e18);

        // 8. Query Agent info
        RootNet.AgentInfo memory agentInfo = rootNet.getAgentInfo(agent1, subnetId);
        assertEq(agentInfo.owner, owner1);
        assertTrue(agentInfo.isValid);
        assertEq(agentInfo.stake, 300_000 * 1e18);

        // 9. Emission — settle epoch 0 (no weights for epoch 0)
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);
        assertEq(emission.currentEpoch(), 1);

        // 10. Settle epoch 1 (weights active)
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);
        assertEq(emission.currentEpoch(), 2);

        // Subnet contract receives 50% of emission (epoch 1, with decay)
        uint256 subnetBal = awp.balanceOf(subnetManager1);
        uint256 decayedEmission = INITIAL_DAILY_EMISSION * 996844 / 1000000;
        assertEq(subnetBal, decayedEmission / 2);

        // DAO receives 50% of emission
        uint256 treasuryBal = awp.balanceOf(address(treasury));
        assertTrue(treasuryBal > 90_000_000 * 1e18);

        // 11. Third epoch (verify continued decay)
        _submitWeight(subnetManager1, uint96(1000));
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);
        assertEq(emission.currentEpoch(), 3);

        uint256 subnetBal2 = awp.balanceOf(subnetManager1);
        uint256 epoch2SubnetEmission = subnetBal2 - subnetBal;
        assertTrue(epoch2SubnetEmission < decayedEmission / 2);
    }

    /// @notice Test multi-subnet emission distribution
    function test_multiSubnetEmission() public {
        vm.prank(airdrop);
        awp.transfer(owner1, 3_000_000 * 1e18);

        vm.prank(owner1);
        rootNet.register();

        vm.startPrank(owner1);
        awp.approve(address(rootNet), 2_000_000 * 1e18);

        uint256 subnet1 = rootNet.registerSubnet(
            IRootNet.SubnetParams("Sub1", "S1", subnetManager1, bytes32(0), 0, "")
        );
        uint256 subnet2 = rootNet.registerSubnet(
            IRootNet.SubnetParams("Sub2", "S2", subnetManager2, bytes32(0), 0, "")
        );

        rootNet.activateSubnet(subnet1);
        rootNet.activateSubnet(subnet2);
        vm.stopPrank();

        {
            address[] memory addrs = new address[](2);
            addrs[0] = subnetManager1;
            addrs[1] = subnetManager2;
            uint96[] memory ws = new uint96[](2);
            ws[0] = 300;
            ws[1] = 100;
            _submitWeights(addrs, ws);
        }

        // Settle epoch 0 (no weights — all to DAO)
        _settleOneEpoch();

        // Settle epoch 1 (weights active)
        _settleOneEpoch();

        uint256 bal1 = awp.balanceOf(subnetManager1);
        uint256 bal2 = awp.balanceOf(subnetManager2);

        assertApproxEqRel(bal1, bal2 * 3, 0.01e18);
    }

    /// @notice Test stake via StakeNFT + withdraw flow
    function test_stakeWithdrawFlow() public {
        vm.prank(airdrop);
        awp.transfer(owner1, 1_000_000 * 1e18);

        vm.startPrank(owner1);
        rootNet.register();

        // Deposit via StakeNFT with short lock
        awp.approve(address(stakeNFT), 500_000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(500_000 * 1e18, 1 days); // 1 day lock
        vm.stopPrank();

        assertEq(stakeNFT.getUserTotalStaked(owner1), 500_000 * 1e18);

        // Cannot withdraw before lock expires
        vm.prank(owner1);
        vm.expectRevert(StakeNFT.LockNotExpired.selector);
        stakeNFT.withdraw(tokenId);

        // Advance past 1-day timestamp lock
        _settleOneEpoch();
        _settleOneEpoch();

        vm.prank(owner1);
        stakeNFT.withdraw(tokenId);

        assertEq(awp.balanceOf(owner1), 1_000_000 * 1e18); // full balance restored
        assertEq(stakeNFT.getUserTotalStaked(owner1), 0);
    }

    /// @notice Test freeze release flow after Agent removal
    function test_agentRemovalAndFreeze() public {
        vm.prank(airdrop);
        awp.transfer(owner1, 2_000_000 * 1e18);

        vm.prank(owner1);
        rootNet.register();
        vm.prank(agent1);
        rootNet.bind(owner1);

        // Register subnet + deposit via StakeNFT + allocate
        vm.startPrank(owner1);
        awp.approve(address(rootNet), 1_000_000 * 1e18);
        uint256 subnetId = rootNet.registerSubnet(
            IRootNet.SubnetParams("Sub", "SUB", subnetManager1, bytes32(0), 0, "")
        );
        rootNet.activateSubnet(subnetId);
        awp.approve(address(stakeNFT), 500_000 * 1e18);
        stakeNFT.deposit(500_000 * 1e18, 52 weeks);
        rootNet.allocate(agent1, subnetId, 300_000 * 1e18);

        // Remove Agent (StakingVault auto-enumerates subnets)
        rootNet.removeAgent(agent1);
        vm.stopPrank();

        // Allocation zeroed out
        assertEq(vault.getAgentStake(owner1, agent1, subnetId), 0);
        // userTotalAllocated released
        assertEq(vault.userTotalAllocated(owner1), 0);
    }

    /// @notice Test registerAndStake one-click operation
    function test_registerAndStake() public {
        vm.prank(airdrop);
        awp.transfer(owner1, 2_000_000 * 1e18);

        // Register agent and subnet first
        vm.prank(owner1);
        rootNet.register();
        vm.prank(agent1);
        rootNet.bind(owner1);

        vm.startPrank(owner1);
        awp.approve(address(rootNet), 1_000_000 * 1e18);
        uint256 subnetId = rootNet.registerSubnet(
            IRootNet.SubnetParams("Sub", "SUB", subnetManager1, bytes32(0), 0, "")
        );
        rootNet.activateSubnet(subnetId);
        vm.stopPrank();

        // New user does registerAndStake
        address owner2 = address(0x102);
        vm.prank(airdrop);
        awp.transfer(owner2, 100_000 * 1e18);

        vm.prank(owner2);
        rootNet.register();
        address agent3 = address(0x203);
        vm.prank(agent3);
        rootNet.bind(owner2);

        vm.startPrank(owner2);
        awp.approve(address(stakeNFT), 50_000 * 1e18);
        // registerAndStake: depositAmount, lockDuration, agent, subnetId, allocateAmount
        rootNet.registerAndStake(50_000 * 1e18, 52 weeks, agent3, subnetId, 30_000 * 1e18);
        vm.stopPrank();

        assertTrue(access.isRegistered(owner2));
        assertEq(stakeNFT.getUserTotalStaked(owner2), 50_000 * 1e18);
        assertEq(vault.getAgentStake(owner2, agent3, subnetId), 30_000 * 1e18);
    }

    /// @notice Test batch emission (many subnets)
    function test_batchSettle() public {
        vm.prank(airdrop);
        awp.transfer(owner1, 10_000_000 * 1e18);

        vm.prank(owner1);
        rootNet.register();

        vm.startPrank(owner1);
        awp.approve(address(rootNet), 10_000_000 * 1e18);

        for (uint256 i = 0; i < 3; i++) {
            address sc = address(uint160(0x400 + i));
            rootNet.registerSubnet(IRootNet.SubnetParams("Sub", "SUB", sc, bytes32(0), 0, ""));
            rootNet.activateSubnet(i + 1);
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

        // Settle epoch 0 (no weights — all to DAO)
        _settleOneEpoch();

        // Settle epoch 1 (weights active)
        _settleOneEpoch();

        uint256 decayedEmission = INITIAL_DAILY_EMISSION * 996844 / 1000000;
        uint256 subnetPool = decayedEmission / 2;
        uint256 expected1 = subnetPool * 100 / 600;
        uint256 expected2 = subnetPool * 200 / 600;
        uint256 expected3 = subnetPool * 300 / 600;

        assertEq(awp.balanceOf(address(uint160(0x400))), expected1);
        assertEq(awp.balanceOf(address(uint160(0x401))), expected2);
        assertEq(awp.balanceOf(address(uint160(0x402))), expected3);
    }

    /// @notice Test pause protection
    function test_pauseProtection() public {
        vm.prank(guardian);
        rootNet.pause();

        vm.prank(owner1);
        vm.expectRevert();
        rootNet.register();

        vm.prank(address(treasury));
        rootNet.unpause();

        vm.prank(owner1);
        rootNet.register();
    }

    /// @notice Test Manager Agent proxy operations
    function test_managerAgentOperations() public {
        vm.prank(airdrop);
        awp.transfer(owner1, 2_000_000 * 1e18);

        vm.prank(owner1);
        rootNet.register();
        vm.prank(agent1);
        rootNet.bind(owner1);
        vm.prank(agent2);
        rootNet.bind(owner1);

        vm.prank(owner1);
        rootNet.setDelegation(agent1, true);

        // Register subnet
        vm.startPrank(owner1);
        awp.approve(address(rootNet), 1_000_000 * 1e18);
        uint256 subnetId = rootNet.registerSubnet(
            IRootNet.SubnetParams("Sub", "SUB", subnetManager1, bytes32(0), 0, "")
        );
        rootNet.activateSubnet(subnetId);
        // Deposit via StakeNFT
        awp.approve(address(stakeNFT), 500_000 * 1e18);
        stakeNFT.deposit(500_000 * 1e18, 52 weeks);
        vm.stopPrank();

        // Manager Agent proxies allocation
        vm.prank(agent1);
        rootNet.allocate(agent2, subnetId, 100_000 * 1e18);

        assertEq(vault.getAgentStake(owner1, agent2, subnetId), 100_000 * 1e18);

        // Manager Agent proxies deallocation
        vm.prank(agent1);
        rootNet.deallocate(agent2, subnetId, 50_000 * 1e18);

        assertEq(vault.getAgentStake(owner1, agent2, subnetId), 50_000 * 1e18);
    }
}
