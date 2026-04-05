// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {DeployHelper} from "./helpers/DeployHelper.sol";
import {veAWP} from "../src/core/veAWP.sol";
import {IveAWP} from "../src/interfaces/IveAWP.sol";

contract veAWPTest is DeployHelper {
    function setUp() public {
        _deployAll();
    }

    // ═══════════════════════════════════════════════
    //  Deposit
    // ═══════════════════════════════════════════════

    function test_deposit() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 30 days);
        vm.stopPrank();

        assertEq(veAwp.ownerOf(tokenId), alice);
        assertEq(veAwp.getUserTotalStaked(alice), 1000e18);
        assertEq(veAwp.totalStaked(), 1000e18);

        (uint128 amount, uint64 lockEndTime, uint64 createdAt) = veAwp.positions(tokenId);
        assertEq(amount, 1000e18);
        assertEq(lockEndTime, uint64(block.timestamp) + 30 days);
        assertEq(createdAt, uint64(block.timestamp));
    }

    function test_deposit_zeroAmount_reverts() public {
        vm.prank(alice);
        vm.expectRevert(veAWP.InvalidAmount.selector);
        veAwp.deposit(0, 30 days);
    }

    function test_deposit_lockTooShort_reverts() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        vm.expectRevert(veAWP.LockTooShort.selector);
        veAwp.deposit(1000e18, 1 hours); // < 1 day
        vm.stopPrank();
    }

    // ═══════════════════════════════════════════════
    //  AddToPosition
    // ═══════════════════════════════════════════════

    function test_addToPosition_amount() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 2000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 30 days);

        veAwp.addToPosition(tokenId, 500e18, 0);
        vm.stopPrank();

        (uint128 amount,,) = veAwp.positions(tokenId);
        assertEq(amount, 1500e18);
        assertEq(veAwp.getUserTotalStaked(alice), 1500e18);
        assertEq(veAwp.totalStaked(), 1500e18);
    }

    function test_addToPosition_extendLock() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 30 days);

        uint64 newEnd = uint64(block.timestamp) + 60 days;
        veAwp.addToPosition(tokenId, 0, newEnd);
        vm.stopPrank();

        (, uint64 lockEndTime,) = veAwp.positions(tokenId);
        assertEq(lockEndTime, newEnd);
    }

    function test_addToPosition_expiredPosition_extendAndAdd() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 2000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 1 days);
        vm.stopPrank();

        // Fast forward past lock
        vm.warp(block.timestamp + 2 days);

        // Extend lock first, then add — should work in one tx
        uint64 newEnd = uint64(block.timestamp) + 30 days;
        vm.startPrank(alice);
        awp.approve(address(veAwp), 500e18);
        veAwp.addToPosition(tokenId, 500e18, newEnd);
        vm.stopPrank();

        (uint128 amount, uint64 lockEndTime,) = veAwp.positions(tokenId);
        assertEq(amount, 1500e18);
        assertEq(lockEndTime, newEnd);
    }

    function test_addToPosition_expiredPosition_addOnly_reverts() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 2000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 1 days);
        vm.stopPrank();

        vm.warp(block.timestamp + 2 days);

        vm.startPrank(alice);
        vm.expectRevert(veAWP.PositionExpired.selector);
        veAwp.addToPosition(tokenId, 500e18, 0); // no lock extension
        vm.stopPrank();
    }

    // ═══════════════════════════════════════════════
    //  Withdraw
    // ═══════════════════════════════════════════════

    function test_withdraw() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 1 days);
        vm.stopPrank();

        vm.warp(block.timestamp + 2 days);

        uint256 balBefore = awp.balanceOf(alice);
        vm.prank(alice);
        veAwp.withdraw(tokenId);

        assertEq(awp.balanceOf(alice) - balBefore, 1000e18);
        assertEq(veAwp.getUserTotalStaked(alice), 0);
        assertEq(veAwp.totalStaked(), 0);
    }

    function test_withdraw_lockNotExpired_reverts() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 30 days);
        vm.stopPrank();

        vm.prank(alice);
        vm.expectRevert(veAWP.LockNotExpired.selector);
        veAwp.withdraw(tokenId);
    }

    // ═══════════════════════════════════════════════
    //  Partial Withdraw
    // ═══════════════════════════════════════════════

    function test_partialWithdraw() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 1 days);
        vm.stopPrank();

        vm.warp(block.timestamp + 2 days);

        vm.prank(alice);
        veAwp.partialWithdraw(tokenId, 400e18);

        (uint128 amount,,) = veAwp.positions(tokenId);
        assertEq(amount, 600e18);
        assertEq(veAwp.getUserTotalStaked(alice), 600e18);
        assertEq(veAwp.totalStaked(), 600e18);
    }

    function test_partialWithdraw_exceedsBalance_reverts() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 1 days);
        vm.stopPrank();

        vm.warp(block.timestamp + 2 days);

        vm.prank(alice);
        vm.expectRevert(veAWP.PartialWithdrawExceedsBalance.selector);
        veAwp.partialWithdraw(tokenId, 1000e18); // must be < amount
    }

    // ═══════════════════════════════════════════════
    //  Batch Withdraw
    // ═══════════════════════════════════════════════

    function test_batchWithdraw() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 3000e18);
        uint256 t1 = veAwp.deposit(1000e18, 1 days);
        uint256 t2 = veAwp.deposit(1000e18, 1 days);
        uint256 t3 = veAwp.deposit(1000e18, 1 days);
        vm.stopPrank();

        vm.warp(block.timestamp + 2 days);

        uint256 balBefore = awp.balanceOf(alice);
        uint256[] memory ids = new uint256[](3);
        ids[0] = t1; ids[1] = t2; ids[2] = t3;

        vm.prank(alice);
        veAwp.batchWithdraw(ids);

        assertEq(awp.balanceOf(alice) - balBefore, 3000e18);
        assertEq(veAwp.totalStaked(), 0);
    }

    // ═══════════════════════════════════════════════
    //  Transfer
    // ═══════════════════════════════════════════════

    function test_transfer_updatesAccumulators() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 30 days);

        veAwp.transferFrom(alice, bob, tokenId);
        vm.stopPrank();

        assertEq(veAwp.getUserTotalStaked(alice), 0);
        assertEq(veAwp.getUserTotalStaked(bob), 1000e18);
        assertEq(veAwp.totalStaked(), 1000e18); // unchanged
    }

    // ═══════════════════════════════════════════════
    //  Voting Power
    // ═══════════════════════════════════════════════

    function test_votingPower() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 54 weeks);
        vm.stopPrank();

        uint256 vp = veAwp.getVotingPower(tokenId);
        // amount * sqrt(54 weeks / 7 days) = 1000e18 * sqrt(54) = 1000e18 * 7
        assertEq(vp, 1000e18 * 7);
    }

    function test_totalVotingPower_basedOnTotalStaked() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        veAwp.deposit(1000e18, 30 days);
        vm.stopPrank();

        assertEq(veAwp.totalVotingPower(), 1000e18); // = totalStaked, no multiplier
    }

    // ═══════════════════════════════════════════════
    //  TokenURI
    // ═══════════════════════════════════════════════

    function test_tokenURI_returns_json() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1000e18);
        uint256 tokenId = veAwp.deposit(1000e18, 30 days);
        vm.stopPrank();

        string memory uri = veAwp.tokenURI(tokenId);
        // Should start with data:application/json;base64,
        assertTrue(bytes(uri).length > 40);
    }

    // ═══════════════════════════════════════════════
    //  Guardian
    // ═══════════════════════════════════════════════

    function test_setGuardian() public {
        vm.prank(guardian);
        veAwp.setGuardian(alice);
        assertEq(veAwp.guardian(), alice);
    }

    function test_setGuardian_notGuardian_reverts() public {
        vm.prank(alice);
        vm.expectRevert(veAWP.NotGuardian.selector);
        veAwp.setGuardian(bob);
    }

    function test_rescueToken_cannotRescueAWP() public {
        vm.prank(guardian);
        vm.expectRevert(veAWP.CannotRescueStakedToken.selector);
        veAwp.rescueToken(address(awp), guardian, 1);
    }
}
