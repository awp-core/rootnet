// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {AWPToken} from "../src/token/AWPToken.sol";

contract StakingVaultTest is Test {
    StakingVault public vault;
    StakeNFT public stakeNFT;
    AWPToken public awp;

    address public deployer = address(this);
    address public user1 = makeAddr("user1");
    address public user2 = makeAddr("user2");
    address public agent1 = makeAddr("agent1");
    address public agent2 = makeAddr("agent2");
    address public nonAWPRegistry = makeAddr("nonAWPRegistry");

    uint256 public constant SUBNET_1 = 1;
    uint256 public constant SUBNET_2 = 2;
    uint256 public constant DEPOSIT_AMOUNT = 1000 ether;
    uint256 public constant EPOCH_DURATION = 7 days;
    uint256 public genesisTime;

    /// @dev This test contract acts as awpRegistry; StakeNFT calls awpRegistry.currentEpoch()
    function currentEpoch() external view returns (uint256) {
        return (block.timestamp - genesisTime) / EPOCH_DURATION;
    }

    function setUp() public {
        genesisTime = block.timestamp;

        // Deploy AWPToken (deployer gets INITIAL_MINT)
        awp = new AWPToken("AWP", "AWP", deployer);

        // Deploy StakingVault + StakeNFT (circular dependency)
        // This test contract (address(this)) acts as awpRegistry
        uint64 nonce = vm.getNonce(deployer);
        address predictedVault = vm.computeCreateAddress(deployer, nonce);
        address predictedStakeNFT = vm.computeCreateAddress(deployer, nonce + 1);

        vault = new StakingVault(address(this));
        stakeNFT = new StakeNFT(address(awp), address(vault), address(this));
        vault.setStakeNFT(address(stakeNFT));

        // Give user1 AWP and have them deposit into StakeNFT
        awp.transfer(user1, 10_000 ether);
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 10_000 ether);
        stakeNFT.deposit(DEPOSIT_AMOUNT, 52 weeks);
        vm.stopPrank();
    }

    // ══════════════════════════════════════════════
    // Allocate tests
    // ══════════════════════════════════════════════

    function test_allocate_basic() public {
        vault.allocate(user1, agent1, SUBNET_1, 300 ether);

        assertEq(vault.getAgentStake(user1, agent1, SUBNET_1), 300 ether);
        assertEq(vault.userTotalAllocated(user1), 300 ether);
        assertEq(vault.subnetTotalStake(SUBNET_1), 300 ether);
    }

    function test_allocate_moreThanUnallocated_reverts() public {
        vault.allocate(user1, agent1, SUBNET_1, 800 ether);

        // Only 200 unallocated, allocating 300 should fail
        vm.expectRevert(StakingVault.InsufficientUnallocated.selector);
        vault.allocate(user1, agent2, SUBNET_2, 300 ether);
    }

    function test_allocate_zeroAmount_reverts() public {
        vm.expectRevert(StakingVault.InvalidAmount.selector);
        vault.allocate(user1, agent1, SUBNET_1, 0);
    }

    // ══════════════════════════════════════════════
    // Deallocate tests
    // ══════════════════════════════════════════════

    function test_deallocate_basic() public {
        vault.allocate(user1, agent1, SUBNET_1, 500 ether);
        vault.deallocate(user1, agent1, SUBNET_1, 200 ether);

        assertEq(vault.getAgentStake(user1, agent1, SUBNET_1), 300 ether);
        assertEq(vault.userTotalAllocated(user1), 300 ether);
        assertEq(vault.subnetTotalStake(SUBNET_1), 300 ether);
    }

    function test_deallocate_full_zerosStake() public {
        vault.allocate(user1, agent1, SUBNET_1, 500 ether);
        vault.deallocate(user1, agent1, SUBNET_1, 500 ether);

        assertEq(vault.getAgentStake(user1, agent1, SUBNET_1), 0);
    }

    function test_deallocate_moreThanAllocated_reverts() public {
        vault.allocate(user1, agent1, SUBNET_1, 200 ether);

        vm.expectRevert(StakingVault.InsufficientAllocation.selector);
        vault.deallocate(user1, agent1, SUBNET_1, 300 ether);
    }

    // ══════════════════════════════════════════════
    // Reallocate tests (immediate, no dual-slot)
    // ══════════════════════════════════════════════

    function test_reallocate_immediate() public {
        vault.allocate(user1, agent1, SUBNET_1, 500 ether);

        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 200 ether);

        // Immediate effect
        assertEq(vault.getAgentStake(user1, agent1, SUBNET_1), 300 ether);
        assertEq(vault.getAgentStake(user1, agent2, SUBNET_2), 200 ether);

        // Subnet totals
        assertEq(vault.subnetTotalStake(SUBNET_1), 300 ether);
        assertEq(vault.subnetTotalStake(SUBNET_2), 200 ether);

        // userTotalAllocated unchanged
        assertEq(vault.userTotalAllocated(user1), 500 ether);
    }

    function test_reallocate_multipleAccumulate() public {
        vault.allocate(user1, agent1, SUBNET_1, 500 ether);

        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 100 ether);
        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 150 ether);

        assertEq(vault.getAgentStake(user1, agent1, SUBNET_1), 250 ether);
        assertEq(vault.getAgentStake(user1, agent2, SUBNET_2), 250 ether);
    }

    function test_reallocate_insufficientAllocation_reverts() public {
        vault.allocate(user1, agent1, SUBNET_1, 100 ether);

        vm.expectRevert(StakingVault.InsufficientAllocation.selector);
        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 200 ether);
    }

    function test_reallocate_zeroAmount_reverts() public {
        vault.allocate(user1, agent1, SUBNET_1, 100 ether);

        vm.expectRevert(StakingVault.InvalidAmount.selector);
        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 0);
    }

    // ══════════════════════════════════════════════
    // Freeze Agent allocations tests
    // ══════════════════════════════════════════════

    function test_freezeAgentAllocations_immediateRelease() public {
        vault.allocate(user1, agent1, SUBNET_1, 300 ether);
        vault.allocate(user1, agent1, SUBNET_2, 200 ether);

        vault.freezeAgentAllocations(user1, agent1);

        // Allocations zeroed
        assertEq(vault.getAgentStake(user1, agent1, SUBNET_1), 0);
        assertEq(vault.getAgentStake(user1, agent1, SUBNET_2), 0);

        // Subnet totals reduced
        assertEq(vault.subnetTotalStake(SUBNET_1), 0);
        assertEq(vault.subnetTotalStake(SUBNET_2), 0);

        // userTotalAllocated released
        assertEq(vault.userTotalAllocated(user1), 0);
    }

    function test_freezeAgentAllocations_agentSubnetsCleared() public {
        vault.allocate(user1, agent1, SUBNET_1, 300 ether);
        vault.allocate(user1, agent1, SUBNET_2, 200 ether);

        assertEq(vault.getAgentSubnets(user1, agent1).length, 2);

        vault.freezeAgentAllocations(user1, agent1);

        // Set must be fully cleared after freeze
        assertEq(vault.getAgentSubnets(user1, agent1).length, 0);
    }

    function test_deallocate_full_clearsAgentSubnets() public {
        vault.allocate(user1, agent1, SUBNET_1, 500 ether);
        assertEq(vault.getAgentSubnets(user1, agent1).length, 1);

        vault.deallocate(user1, agent1, SUBNET_1, 500 ether);

        // Subnet removed from set after full deallocation
        assertEq(vault.getAgentSubnets(user1, agent1).length, 0);
    }

    function test_freezeAfterReallocate_setsConsistent() public {
        vault.allocate(user1, agent1, SUBNET_1, 500 ether);

        // Reallocate everything from agent1/SUBNET_1 to agent2/SUBNET_2
        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 500 ether);

        // agent1 should have no subnets left
        assertEq(vault.getAgentSubnets(user1, agent1).length, 0);
        // agent2 should have SUBNET_2
        assertEq(vault.getAgentSubnets(user1, agent2).length, 1);

        // Freeze agent1 — should be a no-op (no allocations)
        vault.freezeAgentAllocations(user1, agent1);
        assertEq(vault.userTotalAllocated(user1), 500 ether);

        // Freeze agent2 — should clear everything
        vault.freezeAgentAllocations(user1, agent2);
        assertEq(vault.getAgentStake(user1, agent2, SUBNET_2), 0);
        assertEq(vault.getAgentSubnets(user1, agent2).length, 0);
        assertEq(vault.userTotalAllocated(user1), 0);
    }

    // ══════════════════════════════════════════════
    // onlyAWPRegistry access control tests
    // ══════════════════════════════════════════════

    function test_onlyAWPRegistry_allocate() public {
        vm.prank(nonAWPRegistry);
        vm.expectRevert(StakingVault.NotAWPRegistry.selector);
        vault.allocate(user1, agent1, SUBNET_1, 100 ether);
    }

    function test_onlyAWPRegistry_deallocate() public {
        vm.prank(nonAWPRegistry);
        vm.expectRevert(StakingVault.NotAWPRegistry.selector);
        vault.deallocate(user1, agent1, SUBNET_1, 100 ether);
    }

    function test_onlyAWPRegistry_reallocate() public {
        vm.prank(nonAWPRegistry);
        vm.expectRevert(StakingVault.NotAWPRegistry.selector);
        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 100 ether);
    }

    function test_onlyAWPRegistry_freezeAgentAllocations() public {
        vm.prank(nonAWPRegistry);
        vm.expectRevert(StakingVault.NotAWPRegistry.selector);
        vault.freezeAgentAllocations(user1, agent1);
    }
}
