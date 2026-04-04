// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {DeployHelper} from "./helpers/DeployHelper.sol";
import {IAWPRegistry} from "../src/interfaces/IAWPRegistry.sol";

/// @title E2E — End-to-end integration tests
contract E2ETest is DeployHelper {
    function setUp() public {
        _deployAll();
    }

    /// @dev Full workflow: register → activate → stake → allocate → emit → settle
    function test_fullWorkflow() public {
        // 1. Register worknet
        uint256 wid = _registerWorknet(alice, "E2E Worknet", "E2E");
        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Pending));

        // 2. Activate (Guardian)
        _activateWorknet(wid);
        info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Active));
        assertEq(awpWorkNet.ownerOf(wid), alice);

        // 3. Bob stakes AWP
        vm.startPrank(bob);
        awp.approve(address(veAwp), 5000e18);
        uint256 tokenId = veAwp.deposit(5000e18, 30 days);
        vm.stopPrank();

        assertEq(veAwp.getUserTotalStaked(bob), 5000e18);
        assertEq(veAwp.totalStaked(), 5000e18);

        // 4. Bob allocates to agent (himself) on the worknet
        vm.prank(bob);
        awpAllocator.allocate(bob, bob, wid, 3000e18);

        assertEq(awpAllocator.getAgentStake(bob, bob, wid), 3000e18);
        assertEq(awpAllocator.worknetTotalStake(wid), 3000e18);

        // 5. Submit emission weights
        vm.warp(GENESIS_TIME + 1);
        uint256[] memory packed = new uint256[](1);
        // weight=100, recipient=alice (worknet manager)
        packed[0] = (uint256(100) << 160) | uint256(uint160(alice));

        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 0);

        // 6. Settle epoch 0
        awpEmission.settleEpoch(100);
        assertEq(awpEmission.settledEpoch(), 1);

        // 7. Bob deallocates
        vm.prank(bob);
        awpAllocator.deallocateAll(bob, bob, wid);
        assertEq(awpAllocator.userTotalAllocated(bob), 0);

        // 8. Warp past lock and withdraw
        vm.warp(block.timestamp + 31 days);
        vm.prank(bob);
        veAwp.withdraw(tokenId);
        assertEq(veAwp.getUserTotalStaked(bob), 0);
    }

    /// @dev Register → cancel → verify refund
    function test_registerAndCancel() public {
        uint256 balBefore = awp.balanceOf(alice);
        uint256 wid = _registerWorknet(alice);

        vm.prank(alice);
        awpRegistry.cancelWorknet(wid);

        assertEq(awp.balanceOf(alice), balBefore); // full refund
    }

    /// @dev Register → reject by guardian → verify refund
    function test_registerAndReject() public {
        uint256 balBefore = awp.balanceOf(alice);
        uint256 wid = _registerWorknet(alice);

        vm.prank(guardian);
        awpRegistry.rejectWorknet(wid);

        assertEq(awp.balanceOf(alice), balBefore);
    }

    /// @dev Multiple worknets on same chain
    function test_multipleWorknets() public {
        uint256 wid1 = _registerWorknet(alice, "Net1", "N1");
        uint256 wid2 = _registerWorknet(bob, "Net2", "N2");

        _activateWorknet(wid1);
        _activateWorknet(wid2);

        assertTrue(awpRegistry.isWorknetActive(wid1));
        assertTrue(awpRegistry.isWorknetActive(wid2));
        assertEq(awpRegistry.getActiveWorknetCount(), 2);
    }

    /// @dev Stake → allocate → partial withdraw → verify allocation coverage
    function test_stakeAllocatePartialWithdraw() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 1 days);
        vm.stopPrank();

        // Allocate 600
        vm.prank(alice);
        awpAllocator.allocate(alice, bob, 845300000001, 600e18);

        // Warp past lock
        vm.warp(block.timestamp + 2 days);

        // Try partial withdraw 500 — should fail (600 allocated, only 400 available)
        vm.prank(alice);
        vm.expectRevert(); // InsufficientUnallocated
        veAwp.partialWithdraw(tokenId, 500e18);

        // Partial withdraw 300 — should work (600 allocated, 700 remaining after)
        vm.prank(alice);
        veAwp.partialWithdraw(tokenId, 300e18);

        assertEq(veAwp.getUserTotalStaked(alice), 700e18);
    }

    /// @dev Delegation: delegate allocates on behalf of staker
    function test_delegateAllocates() public {
        // Alice stakes
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        veAwp.deposit(1000e18, 30 days);
        // Alice grants Bob as delegate
        awpRegistry.grantDelegate(bob);
        vm.stopPrank();

        // Bob allocates alice's stake
        vm.prank(bob);
        awpAllocator.allocate(alice, bob, 845300000001, 500e18);

        assertEq(awpAllocator.getAgentStake(alice, bob, 845300000001), 500e18);
    }
}
