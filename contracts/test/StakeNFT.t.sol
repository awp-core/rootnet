// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {IStakeNFT} from "../src/interfaces/IStakeNFT.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract StakeNFTTest is Test {
    AWPToken awp;
    StakeNFT stakeNFT;
    StakingVault vault;

    address deployer = address(0xD);
    address treasury = address(0xE);
    address awpRegistry; // will be address(this) for access control
    address user1 = address(0x1);
    address user2 = address(0x2);

    uint256 constant VOTE_WEIGHT_DIVISOR = 7 days;

    function setUp() public {
        awpRegistry = address(this);

        vm.startPrank(deployer);

        // Deploy AWPToken
        awp = new AWPToken("AWP Token", "AWP", deployer, 200_000_000 * 1e18);

        // Deploy StakingVault + StakeNFT (circular dependency)
        uint64 deployerNonce = vm.getNonce(deployer);
        address predictedVault = vm.computeCreateAddress(deployer, deployerNonce);
        address predictedStakeNFT = vm.computeCreateAddress(deployer, deployerNonce + 1);

        vault = StakingVault(address(new ERC1967Proxy(
            address(new StakingVault()), abi.encodeCall(StakingVault.initialize, (awpRegistry))
        )));
        stakeNFT = new StakeNFT(address(awp), address(vault), awpRegistry);
        vm.stopPrank();
        vault.setStakeNFT(address(stakeNFT));
        vm.startPrank(deployer);

        assertEq(address(vault), predictedVault);
        assertEq(address(stakeNFT), predictedStakeNFT);

        // Give users AWP
        awp.transfer(user1, 10_000_000 * 1e18);
        awp.transfer(user2, 10_000_000 * 1e18);

        vm.stopPrank();
    }

    // ── Deposit tests ──

    function test_deposit() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        assertEq(tokenId, 1);
        assertEq(stakeNFT.ownerOf(tokenId), user1);
        assertEq(stakeNFT.getUserTotalStaked(user1), 1000 * 1e18);

        (uint128 amount, uint64 lockEndTime, uint64 createdAt) = stakeNFT.positions(tokenId);
        assertEq(amount, 1000 * 1e18);
        assertEq(createdAt, uint64(block.timestamp));
        assertEq(lockEndTime, uint64(block.timestamp) + 52 weeks);
    }

    function test_deposit_revertsZeroAmount() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        vm.expectRevert(StakeNFT.InvalidAmount.selector);
        stakeNFT.deposit(0, 52 weeks);
        vm.stopPrank();
    }

    function test_deposit_revertsLockTooShort() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        vm.expectRevert(StakeNFT.LockTooShort.selector);
        stakeNFT.deposit(1000 * 1e18, 0);
        vm.stopPrank();
    }

    // ── Multiple positions ──

    function test_multiplePositions() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 3000 * 1e18);

        uint256 id1 = stakeNFT.deposit(1000 * 1e18, 10 weeks);
        uint256 id2 = stakeNFT.deposit(2000 * 1e18, 20 weeks);
        vm.stopPrank();

        assertEq(id1, 1);
        assertEq(id2, 2);
        assertEq(stakeNFT.getUserTotalStaked(user1), 3000 * 1e18);
        assertEq(stakeNFT.balanceOf(user1), 2);
    }

    // ── addToPosition tests ──

    function test_addToPosition_addAmount() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 2000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 52 weeks);

        stakeNFT.addToPosition(tokenId, 500 * 1e18, 0);
        vm.stopPrank();

        (uint128 amount,,) = stakeNFT.positions(tokenId);
        assertEq(amount, 1500 * 1e18);
        assertEq(stakeNFT.getUserTotalStaked(user1), 1500 * 1e18);
    }

    function test_addToPosition_extendLock() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 10 weeks);

        (, uint64 oldEnd,) = stakeNFT.positions(tokenId);
        stakeNFT.addToPosition(tokenId, 0, oldEnd + 10 weeks);
        vm.stopPrank();

        (, uint64 newEnd,) = stakeNFT.positions(tokenId);
        assertEq(newEnd, oldEnd + 10 weeks);
    }

    function test_addToPosition_cannotShortenLock() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        (, uint64 lockEnd,) = stakeNFT.positions(tokenId);

        vm.prank(user1);
        vm.expectRevert(StakeNFT.LockCannotShorten.selector);
        stakeNFT.addToPosition(tokenId, 0, lockEnd - 1);
    }

    function test_addToPosition_revertsNothingToUpdate() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 52 weeks);

        vm.expectRevert(StakeNFT.NothingToUpdate.selector);
        stakeNFT.addToPosition(tokenId, 0, 0);
        vm.stopPrank();
    }

    function test_addToPosition_revertsNotOwner() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        vm.prank(user2);
        vm.expectRevert(StakeNFT.NotTokenOwner.selector);
        stakeNFT.addToPosition(tokenId, 100 * 1e18, 0);
    }

    // ── Withdraw tests ──

    function test_withdraw_afterExpiry() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 1 days); // lock 1 day (MIN_LOCK_DURATION)
        vm.stopPrank();

        // Advance past lock
        vm.warp(block.timestamp + 1 days + 1);

        uint256 balBefore = awp.balanceOf(user1);
        vm.prank(user1);
        stakeNFT.withdraw(tokenId);
        uint256 balAfter = awp.balanceOf(user1);

        assertEq(balAfter - balBefore, 1000 * 1e18);
        assertEq(stakeNFT.getUserTotalStaked(user1), 0);
    }

    function test_withdraw_beforeExpiry_reverts() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        vm.prank(user1);
        vm.expectRevert(StakeNFT.LockNotExpired.selector);
        stakeNFT.withdraw(tokenId);
    }

    function test_withdraw_insufficientUnallocated_reverts() public {
        // Deposit and allocate via awpRegistry
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 1 days);
        vm.stopPrank();

        // Simulate allocation via awpRegistry
        vm.prank(awpRegistry);
        vault.allocate(user1, address(0x99), 1, 500 * 1e18);

        // Advance past lock
        vm.warp(block.timestamp + 1 days + 1);

        // Withdraw should fail because 500 AWP is allocated
        vm.prank(user1);
        vm.expectRevert(StakeNFT.InsufficientUnallocated.selector);
        stakeNFT.withdraw(tokenId);
    }

    // ── Transfer tests ──

    function test_transfer_updatesAccumulators() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 52 weeks);

        // Transfer to user2
        stakeNFT.transferFrom(user1, user2, tokenId);
        vm.stopPrank();

        assertEq(stakeNFT.getUserTotalStaked(user1), 0);
        assertEq(stakeNFT.getUserTotalStaked(user2), 1000 * 1e18);
        assertEq(stakeNFT.ownerOf(tokenId), user2);
    }

    function test_transfer_revertsIfUndercollateralized() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        // Allocate some stake via awpRegistry
        vm.prank(awpRegistry);
        vault.allocate(user1, address(0x99), 1, 500 * 1e18);

        // Transfer should fail: user1 has 1000 staked, 500 allocated, transfer would leave 0 staked
        vm.prank(user1);
        vm.expectRevert(StakeNFT.InsufficientUnallocated.selector);
        stakeNFT.transferFrom(user1, user2, tokenId);
    }

    // ── Voting power tests ──

    function test_votingPower_calculation() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        uint256 vp = stakeNFT.getVotingPower(tokenId);
        // remainingTime = 52 weeks, effective = min(52 weeks, 54 weeks) = 52 weeks
        // votingPower = 1000e18 * sqrt(52 weeks / 7 days) = 1000e18 * sqrt(52)
        uint256 expected = 1000 * 1e18 * Math.sqrt(52 weeks / VOTE_WEIGHT_DIVISOR);
        assertEq(vp, expected);
    }

    function test_votingPower_cappedAtSqrt54() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        // Lock for 100 weeks — remaining will exceed MAX_WEIGHT_DURATION
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 100 weeks);
        vm.stopPrank();

        uint256 vp = stakeNFT.getVotingPower(tokenId);
        // remainingTime = 100 weeks, effective = min(100 weeks, 54 weeks) = 54 weeks
        // votingPower = 1000e18 * sqrt(54 weeks / 7 days) = 1000e18 * sqrt(54)
        uint256 expected = 1000 * 1e18 * Math.sqrt(54 weeks / VOTE_WEIGHT_DIVISOR);
        assertEq(vp, expected);
    }

    function test_userVotingPower_multiplePositions() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 3000 * 1e18);
        uint256 id1 = stakeNFT.deposit(1000 * 1e18, 10 weeks);
        uint256 id2 = stakeNFT.deposit(2000 * 1e18, 20 weeks);
        vm.stopPrank();

        uint256[] memory tokenIds = new uint256[](2);
        tokenIds[0] = id1;
        tokenIds[1] = id2;
        uint256 total = stakeNFT.getUserVotingPower(user1, tokenIds);
        uint256 vp1 = stakeNFT.getVotingPower(id1);
        uint256 vp2 = stakeNFT.getVotingPower(id2);
        assertEq(total, vp1 + vp2);
    }

    function test_votingPower_zeroAfterExpiry() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 1 days);
        vm.stopPrank();

        // Advance past lock
        vm.warp(block.timestamp + 1 days + 1);

        uint256 vp = stakeNFT.getVotingPower(tokenId);
        assertEq(vp, 0);
    }

    // ── depositFor (only AWPRegistry) ──

    function test_depositFor_onlyAWPRegistry() public {
        vm.prank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);

        vm.prank(user2);
        vm.expectRevert(StakeNFT.NotAWPRegistry.selector);
        stakeNFT.depositFor(user1, 1000 * 1e18, 52 weeks);
    }

    function test_depositFor_success() public {
        vm.prank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);

        vm.prank(awpRegistry);
        uint256 tokenId = stakeNFT.depositFor(user1, 1000 * 1e18, 52 weeks);

        assertEq(stakeNFT.ownerOf(tokenId), user1);
        assertEq(stakeNFT.getUserTotalStaked(user1), 1000 * 1e18);
    }

    // ── remainingTime ──

    function test_remainingTime() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        uint256 tokenId = stakeNFT.deposit(1000 * 1e18, 10 weeks);
        vm.stopPrank();

        uint64 remaining = stakeNFT.remainingTime(tokenId);
        // lockEndTime = block.timestamp + 10 weeks, remaining = 10 weeks
        assertEq(remaining, 10 weeks);
    }

    // ── totalVotingPower ──

    function test_totalVotingPower() public {
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 1000 * 1e18);
        stakeNFT.deposit(1000 * 1e18, 52 weeks);
        vm.stopPrank();

        uint256 tvp = stakeNFT.totalVotingPower();
        // totalVotingPower = totalAWPInContract * SQRT_MAX_WEIGHT_FACTOR (= 7)
        assertEq(tvp, 1000 * 1e18 * 7);
    }
}
