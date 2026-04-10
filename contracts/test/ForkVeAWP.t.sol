// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {veAWP} from "../src/core/veAWP.sol";
import {AWPAllocator} from "../src/core/AWPAllocator.sol";
import {IveAWP} from "../src/interfaces/IveAWP.sol";
import {WorknetToken} from "../src/token/WorknetToken.sol";

/// @title ForkVeAWP — Fork tests for veAWP staking, withdrawal, and edge cases on live Base chain
contract ForkVeAWP is Test {
    // ── Live Base mainnet addresses ──
    address constant AWP_TOKEN = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant VEAWP = 0x0000b534C63D78212f1BDCc315165852793A00A8;
    address constant ALLOCATOR = 0x0000D6BB5e040E35081b3AaF59DD71b21C9800AA;
    address constant REGISTRY = 0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A;

    veAWP ve;
    IERC20 awp;
    AWPAllocator allocator;

    address alice;
    address bob;
    uint256 constant STAKE_AMOUNT = 10_000 ether; // 10,000 AWP
    uint64 constant LOCK_90_DAYS = 90 days;
    uint64 constant LOCK_1_DAY = 1 days;

    function setUp() public {
        // Fork Base mainnet
        vm.createSelectFork(vm.envString("BASE_RPC_URL"));

        ve = veAWP(VEAWP);
        awp = IERC20(AWP_TOKEN);
        allocator = AWPAllocator(ALLOCATOR);

        alice = makeAddr("alice");
        bob = makeAddr("bob");

        // Fund alice and bob with AWP (mint via deal cheatcode)
        deal(AWP_TOKEN, alice, 100_000 ether);
        deal(AWP_TOKEN, bob, 100_000 ether);
    }

    // ═══════════════════════════════════════════
    //  Basic deposit / withdraw
    // ═══════════════════════════════════════════

    function test_deposit_basic() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_90_DAYS);
        vm.stopPrank();

        assertEq(ve.ownerOf(tokenId), alice, "NFT owner");
        (uint128 amount, uint64 lockEnd, uint64 createdAt) = ve.positions(tokenId);
        assertEq(amount, uint128(STAKE_AMOUNT), "position amount");
        assertGt(lockEnd, block.timestamp, "lock not expired");
        assertEq(createdAt, uint64(block.timestamp), "createdAt");
        assertEq(ve.getUserTotalStaked(alice), STAKE_AMOUNT, "userTotalStaked");
    }

    function test_withdraw_after_lock_expires() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        vm.stopPrank();

        uint256 balBefore = awp.balanceOf(alice);

        // Warp past lock
        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        vm.prank(alice);
        ve.withdraw(tokenId);

        assertEq(awp.balanceOf(alice), balBefore + STAKE_AMOUNT, "AWP returned");
        // NFT should be burned
        vm.expectRevert();
        ve.ownerOf(tokenId);
        assertEq(ve.getUserTotalStaked(alice), 0, "userTotalStaked cleared");
    }

    function test_withdraw_reverts_before_lock_expires() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_90_DAYS);

        vm.expectRevert(veAWP.LockNotExpired.selector);
        ve.withdraw(tokenId);
        vm.stopPrank();
    }

    // ═══════════════════════════════════════════
    //  Partial withdraw
    // ═══════════════════════════════════════════

    function test_partialWithdraw() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        vm.stopPrank();

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        uint128 withdrawAmt = 3_000 ether;
        vm.prank(alice);
        ve.partialWithdraw(tokenId, withdrawAmt);

        (uint128 remaining,,) = ve.positions(tokenId);
        assertEq(remaining, uint128(STAKE_AMOUNT) - withdrawAmt, "remaining amount");
        assertEq(ve.getUserTotalStaked(alice), STAKE_AMOUNT - withdrawAmt, "userTotalStaked reduced");
        // NFT still exists
        assertEq(ve.ownerOf(tokenId), alice, "NFT still owned");
    }

    function test_partialWithdraw_reverts_if_equals_full_amount() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        vm.stopPrank();

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        vm.prank(alice);
        vm.expectRevert(veAWP.PartialWithdrawExceedsBalance.selector);
        ve.partialWithdraw(tokenId, uint128(STAKE_AMOUNT)); // amount == pos.amount → revert
    }

    function test_partialWithdraw_reverts_before_lock() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_90_DAYS);

        vm.expectRevert(veAWP.LockNotExpired.selector);
        ve.partialWithdraw(tokenId, 1_000 ether);
        vm.stopPrank();
    }

    // ═══════════════════════════════════════════
    //  Batch withdraw
    // ═══════════════════════════════════════════

    function test_batchWithdraw() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT * 3);
        uint256 id1 = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        uint256 id2 = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        uint256 id3 = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        vm.stopPrank();

        assertEq(ve.getUserTotalStaked(alice), STAKE_AMOUNT * 3, "3x staked");

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        uint256 balBefore = awp.balanceOf(alice);
        uint256[] memory ids = new uint256[](3);
        ids[0] = id1; ids[1] = id2; ids[2] = id3;

        vm.prank(alice);
        ve.batchWithdraw(ids);

        assertEq(awp.balanceOf(alice), balBefore + STAKE_AMOUNT * 3, "all AWP returned");
        assertEq(ve.getUserTotalStaked(alice), 0, "staked cleared");
    }

    function test_batchWithdraw_reverts_if_one_locked() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT * 2);
        uint256 id1 = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        uint256 id2 = ve.deposit(STAKE_AMOUNT, LOCK_90_DAYS); // still locked
        vm.stopPrank();

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        uint256[] memory ids = new uint256[](2);
        ids[0] = id1; ids[1] = id2;

        vm.prank(alice);
        vm.expectRevert(veAWP.LockNotExpired.selector);
        ve.batchWithdraw(ids);
    }

    // ═══════════════════════════════════════════
    //  Allocation coverage check
    // ═══════════════════════════════════════════

    function test_withdraw_blocked_by_allocation() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);

        // Allocate half to a worknet
        allocator.allocate(alice, alice, 845300000002, STAKE_AMOUNT / 2);
        vm.stopPrank();

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        // Full withdraw should fail (5000 AWP allocated, trying to withdraw 10000)
        vm.prank(alice);
        vm.expectRevert(veAWP.InsufficientUnallocated.selector);
        ve.withdraw(tokenId);
    }

    function test_partialWithdraw_respects_allocation() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);

        // Allocate 8000 AWP
        allocator.allocate(alice, alice, 845300000002, 8_000 ether);
        vm.stopPrank();

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        // Can withdraw at most 2000 AWP (10000 - 8000 allocated)
        vm.prank(alice);
        ve.partialWithdraw(tokenId, 2_000 ether); // OK

        // Trying to withdraw 1 more should fail
        vm.prank(alice);
        vm.expectRevert(veAWP.InsufficientUnallocated.selector);
        ve.partialWithdraw(tokenId, 1);
    }

    function test_withdraw_after_deallocate() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        allocator.allocate(alice, alice, 845300000002, STAKE_AMOUNT);
        vm.stopPrank();

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        // Can't withdraw while allocated
        vm.prank(alice);
        vm.expectRevert(veAWP.InsufficientUnallocated.selector);
        ve.withdraw(tokenId);

        // Deallocate all
        vm.prank(alice);
        allocator.deallocateAll(alice, alice, 845300000002);

        // Now withdraw should work
        vm.prank(alice);
        ve.withdraw(tokenId);
        assertEq(ve.getUserTotalStaked(alice), 0);
    }

    // ═══════════════════════════════════════════
    //  Transfer and allocation coverage
    // ═══════════════════════════════════════════

    function test_transfer_blocked_if_insufficient_unallocated() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        allocator.allocate(alice, alice, 845300000002, STAKE_AMOUNT);
        vm.stopPrank();

        // Transfer to bob should fail (alice has 10000 allocated, transferring 10000 staked)
        vm.prank(alice);
        vm.expectRevert(veAWP.InsufficientUnallocated.selector);
        ve.transferFrom(alice, bob, tokenId);
    }

    function test_transfer_succeeds_if_unallocated() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        vm.stopPrank();

        // No allocation — transfer should work
        vm.prank(alice);
        ve.transferFrom(alice, bob, tokenId);

        assertEq(ve.ownerOf(tokenId), bob);
        assertEq(ve.getUserTotalStaked(alice), 0);
        assertEq(ve.getUserTotalStaked(bob), STAKE_AMOUNT);
    }

    // ═══════════════════════════════════════════
    //  addToPosition
    // ═══════════════════════════════════════════

    function test_addToPosition_extends_lock_and_adds_amount() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT * 2);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);

        uint64 newLockEnd = uint64(block.timestamp + LOCK_90_DAYS);
        ve.addToPosition(tokenId, STAKE_AMOUNT, newLockEnd);
        vm.stopPrank();

        (uint128 amount, uint64 lockEnd, uint64 createdAt) = ve.positions(tokenId);
        assertEq(amount, uint128(STAKE_AMOUNT * 2), "doubled amount");
        assertEq(lockEnd, newLockEnd, "extended lock");
        assertEq(createdAt, uint64(block.timestamp), "createdAt reset");
        assertEq(ve.getUserTotalStaked(alice), STAKE_AMOUNT * 2);
    }

    function test_addToPosition_blocked_on_expired_lock() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT * 2);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        vm.stopPrank();

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        // Adding amount to expired position without extending lock should fail
        vm.prank(alice);
        vm.expectRevert(veAWP.PositionExpired.selector);
        ve.addToPosition(tokenId, STAKE_AMOUNT, 0); // newLockEndTime=0 means no extension
    }

    function test_addToPosition_extend_expired_then_add() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT * 2);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        vm.stopPrank();

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        // Extend lock to future THEN add amount in same call — should work
        uint64 newLockEnd = uint64(block.timestamp + LOCK_90_DAYS);
        vm.prank(alice);
        ve.addToPosition(tokenId, STAKE_AMOUNT, newLockEnd);

        (uint128 amount, uint64 lockEnd,) = ve.positions(tokenId);
        assertEq(amount, uint128(STAKE_AMOUNT * 2), "added to expired+extended");
        assertEq(lockEnd, newLockEnd);
    }

    function test_addToPosition_cannot_shorten_lock() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_90_DAYS);

        (, uint64 currentLockEnd,) = ve.positions(tokenId);

        vm.expectRevert(veAWP.LockCannotShorten.selector);
        ve.addToPosition(tokenId, 0, currentLockEnd - 1);
        vm.stopPrank();
    }

    // ═══════════════════════════════════════════
    //  Voting power
    // ═══════════════════════════════════════════

    function test_votingPower_decreases_over_time() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_90_DAYS);
        vm.stopPrank();

        uint256 vpStart = ve.getVotingPower(tokenId);
        assertGt(vpStart, 0, "initial VP > 0");

        // After 45 days, VP should be lower
        vm.warp(block.timestamp + 45 days);
        uint256 vpMid = ve.getVotingPower(tokenId);
        assertLt(vpMid, vpStart, "VP decreased");

        // After lock expires, VP = 0
        vm.warp(block.timestamp + 90 days);
        uint256 vpEnd = ve.getVotingPower(tokenId);
        assertEq(vpEnd, 0, "VP = 0 after expiry");
    }

    function test_votingPower_capped_at_54_weeks() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, uint64(365 days)); // 1 year > 54 weeks
        vm.stopPrank();

        uint256 vpLong = ve.getVotingPower(tokenId);

        vm.startPrank(bob);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId2 = ve.deposit(STAKE_AMOUNT, uint64(54 * 7 days)); // exactly 54 weeks
        vm.stopPrank();

        uint256 vpCapped = ve.getVotingPower(tokenId2);

        // Both should have same voting power (capped at 54 weeks)
        assertEq(vpLong, vpCapped, "VP capped at 54 weeks");
    }

    // ═══════════════════════════════════════════
    //  Edge cases
    // ═══════════════════════════════════════════

    function test_deposit_minimum_lock() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, 1 ether);
        uint256 tokenId = ve.deposit(1 ether, 1 days); // minimum lock
        vm.stopPrank();

        (uint128 amount,,) = ve.positions(tokenId);
        assertEq(amount, 1 ether);
    }

    function test_deposit_reverts_zero_amount() public {
        vm.prank(alice);
        vm.expectRevert(veAWP.InvalidAmount.selector);
        ve.deposit(0, LOCK_90_DAYS);
    }

    function test_deposit_reverts_lock_too_short() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        vm.expectRevert(veAWP.LockTooShort.selector);
        ve.deposit(STAKE_AMOUNT, 1 hours); // less than MIN_LOCK_DURATION (1 day)
        vm.stopPrank();
    }

    function test_withdraw_by_non_owner_reverts() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        vm.stopPrank();

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        vm.prank(bob);
        vm.expectRevert(veAWP.NotTokenOwner.selector);
        ve.withdraw(tokenId);
    }

    function test_multiple_positions_independent() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT * 2);
        uint256 id1 = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        uint256 id2 = ve.deposit(STAKE_AMOUNT, LOCK_90_DAYS);
        vm.stopPrank();

        assertEq(ve.getUserTotalStaked(alice), STAKE_AMOUNT * 2);

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        // Can withdraw id1 (expired) but not id2 (still locked)
        vm.prank(alice);
        ve.withdraw(id1);
        assertEq(ve.getUserTotalStaked(alice), STAKE_AMOUNT);

        vm.prank(alice);
        vm.expectRevert(veAWP.LockNotExpired.selector);
        ve.withdraw(id2);
    }

    function test_totalStaked_consistency() public {
        uint256 initialTotal = ve.totalStaked();

        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        vm.stopPrank();

        assertEq(ve.totalStaked(), initialTotal + STAKE_AMOUNT, "total increased");

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        vm.prank(alice);
        ve.withdraw(tokenId);

        assertEq(ve.totalStaked(), initialTotal, "total restored");
    }

    // ═══════════════════════════════════════════
    //  Transfer after withdrawal by new owner
    // ═══════════════════════════════════════════

    function test_transfer_then_new_owner_withdraws() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_1_DAY);
        ve.transferFrom(alice, bob, tokenId);
        vm.stopPrank();

        vm.warp(block.timestamp + LOCK_1_DAY + 1);

        uint256 bobBalBefore = awp.balanceOf(bob);
        vm.prank(bob);
        ve.withdraw(tokenId);

        assertEq(awp.balanceOf(bob), bobBalBefore + STAKE_AMOUNT, "bob received AWP");
        assertEq(ve.getUserTotalStaked(bob), 0);
    }

    // ═══════════════════════════════════════════
    //  On-chain metadata (tokenURI)
    // ═══════════════════════════════════════════

    function test_tokenURI_returns_valid_json() public {
        vm.startPrank(alice);
        awp.approve(VEAWP, STAKE_AMOUNT);
        uint256 tokenId = ve.deposit(STAKE_AMOUNT, LOCK_90_DAYS);
        vm.stopPrank();

        string memory uri = ve.tokenURI(tokenId);
        // Should start with data:application/json;base64,
        assertGt(bytes(uri).length, 40, "URI not empty");
    }
}

// Test that WorknetToken mint works after burn (the fix for supply-lock underflow)
contract ForkWorknetTokenBurn is Test {
    address constant AWP_TOKEN = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant REGISTRY = 0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A;

    function setUp() public {
        vm.createSelectFork(vm.envString("BASE_RPC_URL"));
    }

    function test_mint_after_burn_does_not_revert() public {
        // Deploy a fresh WorknetToken using the fixed code
        MockTokenFactory factory = new MockTokenFactory();
        address token = factory.deployForTest("Test", "TST", 1);

        WorknetToken wt = WorknetToken(token);

        // Open mint phase: mint some tokens
        wt.mint(address(this), 1_000_000 ether);
        assertEq(wt.totalSupply(), 1_000_000 ether);

        // Lock minter
        wt.setMinter(address(this));
        uint256 lock = wt.supplyAtLock();
        assertEq(lock, 1_000_000 ether);

        // Burn tokens below supplyAtLock
        wt.burn(500_000 ether);
        assertLt(wt.totalSupply(), lock, "supply < supplyAtLock after burn");

        // Warp 1 day for time cap
        vm.warp(block.timestamp + 1 days);

        // Mint should NOT revert (the fix handles supply < lock gracefully)
        wt.mint(address(this), 1 ether);
        assertEq(wt.totalSupply(), 500_001 ether);
    }
}

// Minimal factory mock for testing (WorknetToken reads params from msg.sender via callback)
contract MockTokenFactory {
    string public pendingName;
    string public pendingSymbol;
    uint256 public pendingWorknetId;

    function deployForTest(string memory name, string memory symbol, uint256 wid) external returns (address) {
        pendingName = name;
        pendingSymbol = symbol;
        pendingWorknetId = wid;
        return address(new WorknetToken());
    }
}
