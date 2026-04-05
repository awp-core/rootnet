// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {DeployHelper} from "./helpers/DeployHelper.sol";
import {AWPAllocator} from "../src/core/AWPAllocator.sol";

contract AWPAllocatorTest is DeployHelper {
    uint256 constant STAKE_AMOUNT = 10_000e18;
    uint256 tokenId;

    function setUp() public {
        _deployAll();
        // Alice stakes AWP
        vm.startPrank(alice);
        awp.approve(address(veAwp), STAKE_AMOUNT);
        tokenId = veAwp.deposit(STAKE_AMOUNT, 30 days);
        vm.stopPrank();
    }

    // ═══════════════════════════════════════════════
    //  Allocate
    // ═══════════════════════════════════════════════

    function test_allocate() public {
        vm.prank(alice);
        awpAllocator.allocate(alice, bob, 845300000001, 1000e18);

        assertEq(awpAllocator.getAgentStake(alice, bob, 845300000001), 1000e18);
        assertEq(awpAllocator.userTotalAllocated(alice), 1000e18);
        assertEq(awpAllocator.worknetTotalStake(845300000001), 1000e18);
    }

    function test_allocate_exceedsStake_reverts() public {
        vm.prank(alice);
        vm.expectRevert(AWPAllocator.InsufficientUnallocated.selector);
        awpAllocator.allocate(alice, bob, 845300000001, STAKE_AMOUNT + 1);
    }

    function test_allocate_zeroWorknetId_reverts() public {
        vm.prank(alice);
        vm.expectRevert(AWPAllocator.ZeroWorknetId.selector);
        awpAllocator.allocate(alice, bob, 0, 1000e18);
    }

    function test_allocate_zeroAmount_reverts() public {
        vm.prank(alice);
        vm.expectRevert(AWPAllocator.ZeroAmount.selector);
        awpAllocator.allocate(alice, bob, 845300000001, 0);
    }

    // ═══════════════════════════════════════════════
    //  Deallocate
    // ═══════════════════════════════════════════════

    function test_deallocate() public {
        vm.prank(alice);
        awpAllocator.allocate(alice, bob, 845300000001, 1000e18);

        vm.prank(alice);
        awpAllocator.deallocate(alice, bob, 845300000001, 400e18);

        assertEq(awpAllocator.getAgentStake(alice, bob, 845300000001), 600e18);
        assertEq(awpAllocator.userTotalAllocated(alice), 600e18);
    }

    function test_deallocate_exceedsAllocation_reverts() public {
        vm.prank(alice);
        awpAllocator.allocate(alice, bob, 845300000001, 1000e18);

        vm.prank(alice);
        vm.expectRevert(AWPAllocator.InsufficientAllocation.selector);
        awpAllocator.deallocate(alice, bob, 845300000001, 1001e18);
    }

    // ═══════════════════════════════════════════════
    //  DeallocateAll
    // ═══════════════════════════════════════════════

    function test_deallocateAll() public {
        vm.prank(alice);
        awpAllocator.allocate(alice, bob, 845300000001, 1000e18);

        vm.prank(alice);
        awpAllocator.deallocateAll(alice, bob, 845300000001);

        assertEq(awpAllocator.getAgentStake(alice, bob, 845300000001), 0);
        assertEq(awpAllocator.userTotalAllocated(alice), 0);
    }

    // ═══════════════════════════════════════════════
    //  Reallocate
    // ═══════════════════════════════════════════════

    function test_reallocate() public {
        vm.prank(alice);
        awpAllocator.allocate(alice, bob, 845300000001, 1000e18);

        vm.prank(alice);
        awpAllocator.reallocate(alice, bob, 845300000001, alice, 845300000002, 400e18);

        assertEq(awpAllocator.getAgentStake(alice, bob, 845300000001), 600e18);
        assertEq(awpAllocator.getAgentStake(alice, alice, 845300000002), 400e18);
        assertEq(awpAllocator.userTotalAllocated(alice), 1000e18); // unchanged
    }

    // ═══════════════════════════════════════════════
    //  Batch
    // ═══════════════════════════════════════════════

    function test_batchAllocate() public {
        address[] memory agents = new address[](2);
        agents[0] = bob; agents[1] = alice;
        uint256[] memory wids = new uint256[](2);
        wids[0] = 845300000001; wids[1] = 845300000002;
        uint256[] memory amounts = new uint256[](2);
        amounts[0] = 1000e18; amounts[1] = 2000e18;

        vm.prank(alice);
        awpAllocator.batchAllocate(alice, agents, wids, amounts);

        assertEq(awpAllocator.getAgentStake(alice, bob, 845300000001), 1000e18);
        assertEq(awpAllocator.getAgentStake(alice, alice, 845300000002), 2000e18);
        assertEq(awpAllocator.userTotalAllocated(alice), 3000e18);
    }

    function test_batchDeallocate_zeroMeansAll() public {
        vm.prank(alice);
        awpAllocator.allocate(alice, bob, 845300000001, 1000e18);
        vm.prank(alice);
        awpAllocator.allocate(alice, bob, 845300000002, 2000e18);

        address[] memory agents = new address[](2);
        agents[0] = bob; agents[1] = bob;
        uint256[] memory wids = new uint256[](2);
        wids[0] = 845300000001; wids[1] = 845300000002;
        uint256[] memory amounts = new uint256[](2);
        amounts[0] = 0; amounts[1] = 0; // 0 = deallocate all

        vm.prank(alice);
        awpAllocator.batchDeallocate(alice, agents, wids, amounts);

        assertEq(awpAllocator.userTotalAllocated(alice), 0);
    }

    // ═══════════════════════════════════════════════
    //  Delegation
    // ═══════════════════════════════════════════════

    function test_delegate_canAllocate() public {
        // Alice grants delegate to bob
        vm.prank(alice);
        awpRegistry.grantDelegate(bob);

        // Bob allocates on behalf of alice
        vm.prank(bob);
        awpAllocator.allocate(alice, bob, 845300000001, 1000e18);

        assertEq(awpAllocator.getAgentStake(alice, bob, 845300000001), 1000e18);
    }

    function test_nonDelegate_reverts() public {
        vm.prank(bob);
        vm.expectRevert(AWPAllocator.NotAuthorized.selector);
        awpAllocator.allocate(alice, bob, 845300000001, 1000e18);
    }

    // ═══════════════════════════════════════════════
    //  AgentWorknets
    // ═══════════════════════════════════════════════

    function test_agentWorknets_tracking() public {
        vm.startPrank(alice);
        awpAllocator.allocate(alice, bob, 845300000001, 500e18);
        awpAllocator.allocate(alice, bob, 845300000002, 500e18);
        vm.stopPrank();

        uint256[] memory wids = awpAllocator.getAgentWorknets(alice, bob);
        assertEq(wids.length, 2);

        vm.prank(alice);
        awpAllocator.deallocateAll(alice, bob, 845300000001);

        wids = awpAllocator.getAgentWorknets(alice, bob);
        assertEq(wids.length, 1);
    }

    // ═══════════════════════════════════════════════
    //  Guardian
    // ═══════════════════════════════════════════════

    function test_setGuardian() public {
        vm.prank(guardian);
        awpAllocator.setGuardian(alice);
        assertEq(awpAllocator.guardian(), alice);
    }

    function test_setGuardian_notGuardian_reverts() public {
        vm.prank(alice);
        vm.expectRevert(AWPAllocator.NotGuardian.selector);
        awpAllocator.setGuardian(bob);
    }
}
