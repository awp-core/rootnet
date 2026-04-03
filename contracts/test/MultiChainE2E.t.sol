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
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {TimelockController} from "@openzeppelin/contracts/governance/TimelockController.sol";

/// @title MultiChainE2EBase — Shared deployment logic and helper functions
/// @dev Chain-specific test contracts inherit this base, returning chain ID via _targetChainId()
abstract contract MultiChainE2EBase is Test {
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

    // Test users
    address alice = address(0x1001);
    address bob = address(0x1002);
    address charlie = address(0x1003);

    // Test agents
    address agentA = address(0x2001);
    address agentB = address(0x2002);

    // Test worknet contract addresses
    address worknetC1 = address(0x3001);
    address worknetC2 = address(0x3002);

    uint256 constant INITIAL_DAILY = 31_600_000 * 1e18;
    uint256 constant EPOCH = 1 days;
    uint256 constant LP_COST = 1_000_000 * 1e18; // 100M Alpha * 0.01 AWP

    /// @dev Subclass must implement, returns target chain ID
    function _targetChainId() internal pure virtual returns (uint256);

    /// @dev Subclass may override to indicate PancakeSwap (BSC) vs Uniswap (other chains)
    function _isPancakeSwap() internal pure virtual returns (bool) {
        return false;
    }

    function setUp() public {
        vm.chainId(_targetChainId());
        _deploy();
    }

    function _deploy() internal {
        vm.startPrank(deployer);

        // Deploy AWPToken
        awp = new AWPToken("AWP Token", "AWP", deployer);
        awp.initialMint(200_000_000 * 1e18);

        // Deploy AlphaTokenFactory (vanityRule=0 skips validation)
        factory = new AlphaTokenFactory(deployer, 0);

        // Deploy Treasury (TimelockController)
        address[] memory proposers = new address[](0);
        address[] memory executors = new address[](1);
        executors[0] = address(0);
        treasury = new Treasury(1, proposers, executors, deployer);

        // Deploy AWPRegistry (UUPS proxy)
        AWPRegistry awpRegistryImpl = new AWPRegistry();
        awpRegistry = AWPRegistry(address(new ERC1967Proxy(
            address(awpRegistryImpl),
            abi.encodeCall(AWPRegistry.initialize, (deployer, address(treasury), guardian))
        )));

        // Deploy WorknetNFT
        nft = new WorknetNFT("AWP Worknet", "AWPSUB", address(awpRegistry));

        // MockLPManager (no DEX dependency)
        lp = new MockLPManager(address(awpRegistry), address(awp));

        // Deploy AWPEmission (UUPS proxy) — deployer acts as guardian in tests
        AWPEmission emissionImpl = new AWPEmission();
        bytes memory emissionInitData = abi.encodeCall(
            AWPEmission.initialize,
            (address(awp), deployer, INITIAL_DAILY, block.timestamp, EPOCH, address(treasury))
        );
        ERC1967Proxy emissionProxy = new ERC1967Proxy(address(emissionImpl), emissionInitData);
        emission = AWPEmission(address(emissionProxy));

        // Set AWP minter and renounce admin
        awp.addMinter(address(emission));
        awp.renounceAdmin();

        // Set factory addresses
        factory.setAddresses(address(awpRegistry));

        // Deploy StakingVault (UUPS proxy) + StakeNFT
        vault = StakingVault(address(new ERC1967Proxy(
            address(new StakingVault()),
            abi.encodeCall(StakingVault.initialize, (address(awpRegistry), deployer))
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

        // Configure Treasury roles
        treasury.grantRole(treasury.PROPOSER_ROLE(), address(dao));
        treasury.grantRole(treasury.CANCELLER_ROLE(), address(dao));
        treasury.grantRole(treasury.DEFAULT_ADMIN_ROLE(), guardian);
        treasury.renounceRole(treasury.DEFAULT_ADMIN_ROLE(), deployer);

        // Initialize Registry (inject all module addresses, also calls vault.setStakeNFT)
        awpRegistry.initializeRegistry(
            address(awp),
            address(nft),
            address(factory),
            address(emission),
            address(lp),
            address(vault),
            address(stakeNFT),
            address(0), // defaultWorknetManagerImpl — no auto-deploy in tests
            ""          // dexConfig — not needed with MockLPManager
        );

        // Distribute tokens
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

    // ══════════════════════════════════════════════
    //  Helper functions
    // ══════════════════════════════════════════════

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

    function _depositAndAllocate(
        address staker,
        address agent,
        uint256 worknetId,
        uint256 deposit,
        uint256 alloc
    ) internal {
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

    // ══════════════════════════════════════════════
    //  Test 1: Verify chain ID is set correctly
    // ══════════════════════════════════════════════

    function test_chainIdIsCorrect() public view {
        assertEq(block.chainid, _targetChainId(), "block.chainid mismatch");
    }

    // ══════════════════════════════════════════════
    //  Test 2: WorknetId encoding (chainId << 64 | localCounter)
    // ══════════════════════════════════════════════

    function test_worknetIdEncoding() public {
        _registerUser(alice);
        uint256 worknetId = _registerWorknet(alice, worknetC1);

        // First worknet: localCounter = 1
        uint256 expectedId = (_targetChainId() << 64) | 1;
        assertEq(worknetId, expectedId, "worknetId encoding: (chainId << 64) | 1");

        // Verify chainId can be extracted
        uint256 extractedChainId = worknetId >> 64;
        assertEq(extractedChainId, _targetChainId(), "extracted chainId mismatch");

        // Verify localCounter
        uint256 localCounter = worknetId & ((1 << 64) - 1);
        assertEq(localCounter, 1, "localCounter should be 1");

        // Second worknet: localCounter = 2
        uint256 worknetId2 = _registerWorknet(alice, worknetC2);
        uint256 expectedId2 = (_targetChainId() << 64) | 2;
        assertEq(worknetId2, expectedId2, "second worknetId encoding: (chainId << 64) | 2");
    }

    // ══════════════════════════════════════════════
    //  Test 3: Full lifecycle — register, bind, register worknet, activate, stake, allocate, emit
    // ══════════════════════════════════════════════

    function test_fullLifecycle() public {
        // 1. Register user and bind agent
        _registerUser(alice);
        _bindAgent(agentA, alice);

        assertEq(awpRegistry.boundTo(agentA), alice, "agentA should be bound to alice");
        assertEq(awpRegistry.recipient(alice), alice, "alice recipient should be self after register");

        // 2. Register worknet, verify worknetId encoding
        uint256 worknetId = _registerWorknet(alice, worknetC1);
        uint256 expectedId = (_targetChainId() << 64) | 1;
        assertEq(worknetId, expectedId, "worknetId should match (chainId << 64) | 1");

        // 3. Activate worknet
        vm.prank(alice);
        awpRegistry.activateWorknet(worknetId);
        assertTrue(awpRegistry.isWorknetActive(worknetId), "worknet should be active");

        // 4. Stake and allocate
        uint256 stakeAmount = 10_000_000 * 1e18;
        uint256 allocAmount = 8_000_000 * 1e18;
        _depositAndAllocate(alice, agentA, worknetId, stakeAmount, allocAmount);

        assertEq(vault.getAgentStake(alice, agentA, worknetId), allocAmount, "allocation mismatch");
        assertEq(vault.getWorknetTotalStake(worknetId), allocAmount, "worknet total stake mismatch");

        // 5. Submit weights and settle 2 epochs
        _submitWeight(worknetC1, uint96(100));
        _settleEpoch(); // settle epoch 0

        uint256 balAfterEpoch0 = awp.balanceOf(worknetC1);

        _submitWeight(worknetC1, uint96(100));
        _settleEpoch(); // settle epoch 1

        // 6. Verify emission amounts
        uint256 worknetBalance = awp.balanceOf(worknetC1);
        assertTrue(worknetBalance > 0, "worknet should have received emission");

        // Total emission = epoch 0 + epoch 1 (both emit, epoch 0 uses initial weights)
        // epoch 0 uses INITIAL_DAILY, epoch 1 uses decayed amount
        uint256 epoch1Emission = worknetBalance - balAfterEpoch0;
        uint256 expectedEpoch1 = INITIAL_DAILY * 996844 / 1000000;
        assertApproxEqRel(epoch1Emission, expectedEpoch1, 0.01e18, "epoch 1 emission amount mismatch");
    }

    // ══════════════════════════════════════════════
    //  Test 4: Cross-chain worknetId allocation
    // ══════════════════════════════════════════════

    function test_crossChainWorknetIdAllocation() public {
        _registerUser(alice);
        _bindAgent(agentA, alice);

        // Register local worknet
        uint256 localWorknetId = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(localWorknetId);

        // Construct foreign-chain worknetId (different chainId)
        uint256 foreignChainId;
        if (_targetChainId() == 8453) {
            foreignChainId = 1; // Base -> Ethereum
        } else if (_targetChainId() == 1) {
            foreignChainId = 56; // Ethereum -> BSC
        } else if (_targetChainId() == 56) {
            foreignChainId = 42161; // BSC -> Arbitrum
        } else {
            foreignChainId = 8453; // Arbitrum -> Base
        }

        // Foreign worknetId (assuming localCounter=1)
        uint256 foreignWorknetId = (foreignChainId << 64) | 1;

        // Stake and allocate to foreign worknetId (StakingVault doesn't check on-chain status)
        uint256 stakeAmount = 5_000_000 * 1e18;
        _depositAndAllocate(alice, agentA, foreignWorknetId, stakeAmount, 3_000_000 * 1e18);

        // Verify cross-chain allocation succeeded
        assertEq(
            vault.getAgentStake(alice, agentA, foreignWorknetId),
            3_000_000 * 1e18,
            "cross-chain allocation should succeed"
        );
        assertEq(
            vault.getWorknetTotalStake(foreignWorknetId),
            3_000_000 * 1e18,
            "foreign worknet total stake mismatch"
        );

        // Also allocate to local worknet
        vm.prank(alice);
        vault.allocate(alice, agentA, localWorknetId, 2_000_000 * 1e18);
        assertEq(
            vault.getAgentStake(alice, agentA, localWorknetId),
            2_000_000 * 1e18,
            "local allocation should succeed"
        );

        // User total allocated = 3M + 2M = 5M
        assertEq(vault.userTotalAllocated(alice), 5_000_000 * 1e18, "user total allocated mismatch");
    }

    // ══════════════════════════════════════════════
    //  Test 5: Multi-epoch emission decay verification
    // ══════════════════════════════════════════════

    function test_emissionDecayMultiEpoch() public {
        _registerUser(alice);
        uint256 worknetId = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(worknetId);

        _submitWeight(worknetC1, uint96(100));

        uint256 prevEmission = INITIAL_DAILY;
        for (uint256 i = 0; i < 5; i++) {
            if (i > 0) {
                _submitWeight(worknetC1, uint96(100));
            }
            _settleEpoch();
            uint256 currentEmission = emission.currentDailyEmission();
            if (i > 0) {
                // Per-epoch decay: newEmission = prevEmission * 996844 / 1000000
                assertEq(
                    currentEmission,
                    prevEmission * 996844 / 1000000,
                    "emission decay mismatch at epoch"
                );
            }
            prevEmission = currentEmission;
        }
    }

    // ══════════════════════════════════════════════
    //  Test 6: Multi-user multi-worknet emission distribution
    // ══════════════════════════════════════════════

    function test_multiUserMultiWorknetEmission() public {
        _registerUser(alice);
        _registerUser(bob);

        // Register two worknets
        uint256 worknetId1 = _registerWorknet(alice, worknetC1);
        uint256 worknetId2 = _registerWorknet(alice, worknetC2);
        vm.startPrank(alice);
        awpRegistry.activateWorknet(worknetId1);
        awpRegistry.activateWorknet(worknetId2);
        vm.stopPrank();

        // Verify worknetId encoding
        assertEq(worknetId1, (_targetChainId() << 64) | 1, "worknetId1 encoding");
        assertEq(worknetId2, (_targetChainId() << 64) | 2, "worknetId2 encoding");

        // Submit weights: worknetC1=200, worknetC2=100 (2:1 ratio)
        {
            address[] memory addrs = new address[](2);
            addrs[0] = worknetC1;
            addrs[1] = worknetC2;
            uint96[] memory ws = new uint96[](2);
            ws[0] = 200;
            ws[1] = 100;
            _submitWeights(addrs, ws);
        }

        // Multi-user staking
        _depositAndAllocate(alice, agentA, worknetId1, 10_000_000 * 1e18, 8_000_000 * 1e18);
        _depositAndAllocate(bob, agentB, worknetId2, 5_000_000 * 1e18, 5_000_000 * 1e18);

        // Settle 2 epochs (epoch 0 is warmup)
        _settleEpoch();
        {
            address[] memory addrs = new address[](2);
            addrs[0] = worknetC1;
            addrs[1] = worknetC2;
            uint96[] memory ws = new uint96[](2);
            ws[0] = 200;
            ws[1] = 100;
            _submitWeights(addrs, ws);
        }
        _settleEpoch();

        // Verify emission ratio worknetC1 : worknetC2 ~= 2:1
        uint256 bal1 = awp.balanceOf(worknetC1);
        uint256 bal2 = awp.balanceOf(worknetC2);
        assertTrue(bal1 > 0 && bal2 > 0, "both worknets should have emission");
        assertApproxEqRel(bal1, bal2 * 2, 0.01e18, "emission ratio should be ~2:1");
    }

    // ══════════════════════════════════════════════
    //  Test 7: Worknet status lifecycle (Pending -> Active -> Paused -> Active)
    // ══════════════════════════════════════════════

    function test_worknetStatusLifecycle() public {
        _registerUser(alice);
        uint256 worknetId = _registerWorknet(alice, worknetC1);

        // Pending status
        assertEq(
            uint256(awpRegistry.getWorknet(worknetId).status),
            uint256(IAWPRegistry.WorknetStatus.Pending),
            "should be Pending"
        );

        // Activate
        vm.prank(alice);
        awpRegistry.activateWorknet(worknetId);
        assertEq(
            uint256(awpRegistry.getWorknet(worknetId).status),
            uint256(IAWPRegistry.WorknetStatus.Active),
            "should be Active"
        );

        // Pause
        vm.prank(alice);
        awpRegistry.pauseWorknet(worknetId);
        assertEq(
            uint256(awpRegistry.getWorknet(worknetId).status),
            uint256(IAWPRegistry.WorknetStatus.Paused),
            "should be Paused"
        );

        // Resume
        vm.prank(alice);
        awpRegistry.resumeWorknet(worknetId);
        assertTrue(awpRegistry.isWorknetActive(worknetId), "should be active again");
    }

    // ══════════════════════════════════════════════
    //  Test 8: Delegate-authorized allocation
    // ══════════════════════════════════════════════

    function test_delegateAllocation() public {
        _registerUser(alice);
        _registerUser(bob);

        uint256 worknetId = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(worknetId);

        // Alice stakes
        uint256 stakeAmount = 10_000_000 * 1e18;
        vm.startPrank(alice);
        awp.approve(address(stakeNFT), stakeAmount);
        stakeNFT.deposit(stakeAmount, 52 weeks);
        vm.stopPrank();

        // Alice grants Bob as delegate
        vm.prank(alice);
        awpRegistry.grantDelegate(bob);

        // Bob allocates on behalf of Alice
        vm.prank(bob);
        vault.allocate(alice, agentA, worknetId, 5_000_000 * 1e18);

        assertEq(
            vault.getAgentStake(alice, agentA, worknetId),
            5_000_000 * 1e18,
            "delegate allocation should succeed"
        );

        // Revoke delegate
        vm.prank(alice);
        awpRegistry.revokeDelegate(bob);

        // Bob can no longer allocate
        vm.prank(bob);
        vm.expectRevert(StakingVault.NotAuthorized.selector);
        vault.allocate(alice, agentA, worknetId, 1_000_000 * 1e18);
    }

    // ══════════════════════════════════════════════
    //  Test 9: Binding tree + resolveRecipient
    // ══════════════════════════════════════════════

    function test_bindingTreeResolveRecipient() public {
        // Build binding chain: agentA -> alice, agentB -> agentA
        _bindAgent(agentA, alice);
        _bindAgent(agentB, agentA);

        assertEq(awpRegistry.boundTo(agentA), alice, "agentA bound to alice");
        assertEq(awpRegistry.boundTo(agentB), agentA, "agentB bound to agentA");

        // Alice sets recipient to charlie
        vm.prank(alice);
        awpRegistry.setRecipient(charlie);

        // resolveRecipient should walk binding chain to alice, then return charlie
        assertEq(awpRegistry.resolveRecipient(agentB), charlie, "resolve should walk to charlie");
        assertEq(awpRegistry.resolveRecipient(agentA), charlie, "resolve agentA should get charlie");
    }

    // ══════════════════════════════════════════════
    //  Test 10: Reallocate — immediate reallocation
    // ══════════════════════════════════════════════

    function test_reallocateImmediate() public {
        _registerUser(alice);
        _bindAgent(agentA, alice);
        _bindAgent(agentB, alice);

        uint256 worknetId1 = _registerWorknet(alice, worknetC1);
        uint256 worknetId2 = _registerWorknet(alice, worknetC2);
        vm.startPrank(alice);
        awpRegistry.activateWorknet(worknetId1);
        awpRegistry.activateWorknet(worknetId2);
        vm.stopPrank();

        _depositAndAllocate(alice, agentA, worknetId1, 10_000_000 * 1e18, 8_000_000 * 1e18);

        // Cross-agent and cross-worknet reallocation
        vm.prank(alice);
        vault.reallocate(alice, agentA, worknetId1, agentB, worknetId2, 3_000_000 * 1e18);

        assertEq(vault.getAgentStake(alice, agentA, worknetId1), 5_000_000 * 1e18, "from reduced");
        assertEq(vault.getAgentStake(alice, agentB, worknetId2), 3_000_000 * 1e18, "to increased");
        assertEq(vault.userTotalAllocated(alice), 8_000_000 * 1e18, "total unchanged");
    }

    // ══════════════════════════════════════════════
    //  Test 11: subnetId=0 rejected
    // ══════════════════════════════════════════════

    function test_zeroWorknetIdRejected() public {
        _registerUser(alice);

        uint256 stakeAmount = 1_000_000 * 1e18;
        vm.startPrank(alice);
        awp.approve(address(stakeNFT), stakeAmount);
        stakeNFT.deposit(stakeAmount, 52 weeks);

        // worknetId=0 should be rejected
        vm.expectRevert(StakingVault.ZeroWorknetId.selector);
        vault.allocate(alice, agentA, 0, 500_000 * 1e18);
        vm.stopPrank();
    }

    // ══════════════════════════════════════════════
    //  Test 12: Verify DEX type identifier (chain-specific)
    // ══════════════════════════════════════════════

    function test_dexTypeIdentifier() public view {
        // Verify each chain test contract correctly identifies DEX type
        if (_targetChainId() == 56) {
            assertTrue(_isPancakeSwap(), "BSC should use PancakeSwap");
        } else {
            assertFalse(_isPancakeSwap(), "non-BSC chains should use Uniswap");
        }
    }
}

// ══════════════════════════════════════════════════════
//  Base (chainId=8453) — Uniswap V4 (LPManagerUni + WorknetManagerUni)
// ══════════════════════════════════════════════════════

contract BaseE2ETest is MultiChainE2EBase {
    function _targetChainId() internal pure override returns (uint256) {
        return 8453;
    }

    /// @dev Base chain-specific: worknetId high-bits encoding
    function test_base_worknetIdHighBits() public {
        _registerUser(alice);
        uint256 worknetId = _registerWorknet(alice, worknetC1);

        // 8453 << 64 = 0x2105_0000_0000_0000_0000
        uint256 highBits = worknetId >> 64;
        assertEq(highBits, 8453, "Base chainId in high bits");
    }

    /// @dev Base chain-specific: cross-chain allocation to Ethereum
    function test_base_crossChainToEthereum() public {
        _registerUser(alice);
        _bindAgent(agentA, alice);

        uint256 stakeAmount = 5_000_000 * 1e18;
        vm.startPrank(alice);
        awp.approve(address(stakeNFT), stakeAmount);
        stakeNFT.deposit(stakeAmount, 52 weeks);

        // Allocate to Ethereum worknetId
        uint256 ethWorknetId = (uint256(1) << 64) | 42;
        vault.allocate(alice, agentA, ethWorknetId, 2_000_000 * 1e18);
        vm.stopPrank();

        assertEq(vault.getAgentStake(alice, agentA, ethWorknetId), 2_000_000 * 1e18);
    }
}

// ══════════════════════════════════════════════════════
//  Ethereum (chainId=1) — Uniswap V4 (LPManagerUni + WorknetManagerUni)
// ══════════════════════════════════════════════════════

contract EthereumE2ETest is MultiChainE2EBase {
    function _targetChainId() internal pure override returns (uint256) {
        return 1;
    }

    /// @dev Ethereum chain-specific: worknetId high-bits encoding
    function test_ethereum_worknetIdHighBits() public {
        _registerUser(alice);
        uint256 worknetId = _registerWorknet(alice, worknetC1);

        // 1 << 64
        uint256 highBits = worknetId >> 64;
        assertEq(highBits, 1, "Ethereum chainId in high bits");

        // localCounter = 1
        uint256 localId = worknetId & type(uint64).max;
        assertEq(localId, 1, "first worknet localId");
    }

    /// @dev Ethereum chain-specific: cross-chain allocation to Arbitrum and BSC
    function test_ethereum_crossChainMultiple() public {
        _registerUser(alice);
        _bindAgent(agentA, alice);

        uint256 stakeAmount = 10_000_000 * 1e18;
        vm.startPrank(alice);
        awp.approve(address(stakeNFT), stakeAmount);
        stakeNFT.deposit(stakeAmount, 52 weeks);

        // Allocate to Arbitrum
        uint256 arbWorknetId = (uint256(42161) << 64) | 1;
        vault.allocate(alice, agentA, arbWorknetId, 3_000_000 * 1e18);

        // Allocate to BSC
        uint256 bscWorknetId = (uint256(56) << 64) | 5;
        vault.allocate(alice, agentA, bscWorknetId, 2_000_000 * 1e18);
        vm.stopPrank();

        assertEq(vault.getAgentStake(alice, agentA, arbWorknetId), 3_000_000 * 1e18, "arb alloc");
        assertEq(vault.getAgentStake(alice, agentA, bscWorknetId), 2_000_000 * 1e18, "bsc alloc");
        assertEq(vault.userTotalAllocated(alice), 5_000_000 * 1e18, "total alloc");
    }
}

// ══════════════════════════════════════════════════════
//  BSC (chainId=56) — PancakeSwap V4 (LPManager + WorknetManager)
// ══════════════════════════════════════════════════════

contract BSCE2ETest is MultiChainE2EBase {
    function _targetChainId() internal pure override returns (uint256) {
        return 56;
    }

    function _isPancakeSwap() internal pure override returns (bool) {
        return true;
    }

    /// @dev BSC chain-specific: worknetId high-bits encoding
    function test_bsc_worknetIdHighBits() public {
        _registerUser(alice);
        uint256 worknetId = _registerWorknet(alice, worknetC1);

        // 56 << 64
        uint256 highBits = worknetId >> 64;
        assertEq(highBits, 56, "BSC chainId in high bits");
    }

    /// @dev BSC chain-specific: PancakeSwap identifier check
    function test_bsc_pancakeSwapIdentifier() public view {
        assertTrue(_isPancakeSwap(), "BSC must use PancakeSwap");
    }

    /// @dev BSC chain-specific: cross-chain allocation to Base
    function test_bsc_crossChainToBase() public {
        _registerUser(alice);
        _bindAgent(agentA, alice);

        uint256 stakeAmount = 5_000_000 * 1e18;
        vm.startPrank(alice);
        awp.approve(address(stakeNFT), stakeAmount);
        stakeNFT.deposit(stakeAmount, 52 weeks);

        // Allocate to Base worknetId
        uint256 baseWorknetId = (uint256(8453) << 64) | 10;
        vault.allocate(alice, agentA, baseWorknetId, 4_000_000 * 1e18);
        vm.stopPrank();

        assertEq(vault.getAgentStake(alice, agentA, baseWorknetId), 4_000_000 * 1e18);
        assertEq(vault.getWorknetTotalStake(baseWorknetId), 4_000_000 * 1e18);
    }

    /// @dev BSC: Full 3-epoch emission verification
    function test_bsc_threeEpochEmission() public {
        _registerUser(alice);
        uint256 worknetId = _registerWorknet(alice, worknetC1);
        vm.prank(alice);
        awpRegistry.activateWorknet(worknetId);

        _depositAndAllocate(alice, agentA, worknetId, 10_000_000 * 1e18, 5_000_000 * 1e18);

        // Submit weights and settle 3 epochs
        _submitWeight(worknetC1, uint96(100));
        _settleEpoch(); // epoch 0

        uint256 balAfterEpoch0 = awp.balanceOf(worknetC1);

        _submitWeight(worknetC1, uint96(100));
        _settleEpoch(); // epoch 1

        uint256 balAfterEpoch1 = awp.balanceOf(worknetC1);
        uint256 epoch1Emission = balAfterEpoch1 - balAfterEpoch0;

        _submitWeight(worknetC1, uint96(100));
        _settleEpoch(); // epoch 2

        uint256 balAfterEpoch2 = awp.balanceOf(worknetC1);
        uint256 epoch2Emission = balAfterEpoch2 - balAfterEpoch1;

        // epoch 2 amount should be < epoch 1 (decay)
        assertTrue(epoch2Emission < epoch1Emission, "epoch 2 emission should be less due to decay");

        // Verify exact decay ratio
        assertApproxEqRel(
            epoch2Emission,
            epoch1Emission * 996844 / 1000000,
            0.01e18,
            "decay ratio mismatch"
        );
    }
}

// ══════════════════════════════════════════════════════
//  Arbitrum (chainId=42161) — Uniswap V4 (LPManagerUni + WorknetManagerUni)
// ══════════════════════════════════════════════════════

contract ArbitrumE2ETest is MultiChainE2EBase {
    function _targetChainId() internal pure override returns (uint256) {
        return 42161;
    }

    /// @dev Arbitrum chain-specific: worknetId high-bits encoding
    function test_arbitrum_worknetIdHighBits() public {
        _registerUser(alice);
        uint256 worknetId = _registerWorknet(alice, worknetC1);

        // 42161 << 64
        uint256 highBits = worknetId >> 64;
        assertEq(highBits, 42161, "Arbitrum chainId in high bits");
    }

    /// @dev Arbitrum chain-specific: register multiple worknets, verify incrementing localCounter
    function test_arbitrum_multipleWorknetsIncrementCounter() public {
        _registerUser(alice);

        uint256 id1 = _registerWorknet(alice, worknetC1);
        uint256 id2 = _registerWorknet(alice, worknetC2);

        uint256 base = uint256(42161) << 64;
        assertEq(id1, base | 1, "first worknet");
        assertEq(id2, base | 2, "second worknet");

        // Verify nextWorknetId
        uint256 next = awpRegistry.nextWorknetId();
        assertEq(next, base | 3, "next should be 3");
    }

    /// @dev Arbitrum chain-specific: cross-chain allocation to all four chains
    function test_arbitrum_crossChainAllFourChains() public {
        _registerUser(alice);
        _bindAgent(agentA, alice);

        uint256 stakeAmount = 20_000_000 * 1e18;
        vm.startPrank(alice);
        awp.approve(address(stakeNFT), stakeAmount);
        stakeNFT.deposit(stakeAmount, 52 weeks);

        // Allocate to worknetIds on all four chains
        uint256 baseWid = (uint256(8453) << 64) | 1;
        uint256 ethWid = (uint256(1) << 64) | 1;
        uint256 bscWid = (uint256(56) << 64) | 1;
        uint256 arbWid = (uint256(42161) << 64) | 1;

        vault.allocate(alice, agentA, baseWid, 5_000_000 * 1e18);
        vault.allocate(alice, agentA, ethWid, 5_000_000 * 1e18);
        vault.allocate(alice, agentA, bscWid, 5_000_000 * 1e18);
        vault.allocate(alice, agentA, arbWid, 5_000_000 * 1e18);
        vm.stopPrank();

        // Verify allocation per chain
        assertEq(vault.getAgentStake(alice, agentA, baseWid), 5_000_000 * 1e18, "base");
        assertEq(vault.getAgentStake(alice, agentA, ethWid), 5_000_000 * 1e18, "eth");
        assertEq(vault.getAgentStake(alice, agentA, bscWid), 5_000_000 * 1e18, "bsc");
        assertEq(vault.getAgentStake(alice, agentA, arbWid), 5_000_000 * 1e18, "arb");
        assertEq(vault.userTotalAllocated(alice), 20_000_000 * 1e18, "total = 20M");
    }

    /// @dev Arbitrum: verify extractChainId helper
    function test_arbitrum_extractChainId() public {
        _registerUser(alice);
        uint256 worknetId = _registerWorknet(alice, worknetC1);

        uint256 extracted = awpRegistry.extractChainId(worknetId);
        assertEq(extracted, 42161, "extractChainId should return 42161");

        // Verify extractChainId for foreign worknetId
        uint256 baseWorknetId = (uint256(8453) << 64) | 99;
        assertEq(awpRegistry.extractChainId(baseWorknetId), 8453, "extractChainId for Base");
    }
}
