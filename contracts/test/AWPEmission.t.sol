// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {DeployHelper} from "./helpers/DeployHelper.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";

contract AWPEmissionTest is DeployHelper {
    function setUp() public {
        _deployAll();
        // Warp to genesis time so currentEpoch works
        vm.warp(GENESIS_TIME + 1);
    }

    function test_currentEpoch() public view {
        assertEq(awpEmission.currentEpoch(), 0);
    }

    function test_currentEpoch_afterOneDay() public {
        vm.warp(GENESIS_TIME + 1 days + 1);
        assertEq(awpEmission.currentEpoch(), 1);
    }

    function test_submitAllocations() public {
        uint256[] memory packed = new uint256[](1);
        packed[0] = (uint256(100) << 160) | uint256(uint160(alice));

        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 0);

        assertEq(awpEmission.getEpochRecipientCount(0), 1);
        assertEq(awpEmission.getEpochTotalWeight(0), 100);
    }

    function test_settleEpoch() public {
        // Submit weights for epoch 0
        uint256[] memory packed = new uint256[](1);
        packed[0] = (uint256(100) << 160) | uint256(uint160(alice));

        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 0);

        // Settle
        awpEmission.settleEpoch(100);

        assertEq(awpEmission.settledEpoch(), 1);
    }

    function test_settleEpoch_tooEarly_reverts() public {
        // Don't warp — epoch 0 at genesis+1, settledEpoch=0, should be able to settle epoch 0
        // But epoch 1 won't be ready
        uint256[] memory packed = new uint256[](1);
        packed[0] = (uint256(100) << 160) | uint256(uint160(alice));

        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 0);

        // Settle epoch 0 (should work)
        awpEmission.settleEpoch(100);

        // Try settle epoch 1 (not ready)
        vm.expectRevert(AWPEmission.EpochNotReady.selector);
        awpEmission.settleEpoch(100);
    }

    function test_pauseEpochUntil_blocksSettle() public {
        uint256[] memory packed = new uint256[](1);
        packed[0] = (uint256(100) << 160) | uint256(uint160(alice));
        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 0);

        // Pause
        vm.prank(guardian);
        awpEmission.pauseEpochUntil(uint64(block.timestamp + 1 days));

        // Try settle — should fail (paused)
        vm.expectRevert(AWPEmission.EpochNotReady.selector);
        awpEmission.settleEpoch(100);
    }

    function test_pauseEpochUntil_pastTimestamp_autoResumes() public {
        uint256[] memory packed = new uint256[](1);
        packed[0] = (uint256(100) << 160) | uint256(uint160(alice));
        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 0);

        // Pause with past timestamp
        vm.prank(guardian);
        awpEmission.pauseEpochUntil(uint64(block.timestamp - 1));

        // Should be auto-resumed — settle works
        awpEmission.settleEpoch(100);
        assertEq(awpEmission.settledEpoch(), 1);
    }

    function test_setDecayFactor() public {
        vm.prank(guardian);
        awpEmission.setDecayFactor(999000);
        assertEq(awpEmission.decayFactor(), 999000);
    }

    function test_setDecayFactor_invalidRange_reverts() public {
        vm.prank(guardian);
        vm.expectRevert(AWPEmission.InvalidDecayFactor.selector);
        awpEmission.setDecayFactor(1_000_000); // >= DECAY_PRECISION
    }

    function test_setGuardian() public {
        vm.prank(guardian);
        awpEmission.setGuardian(alice);
        assertEq(awpEmission.guardian(), alice);
    }

    function test_immutables() public view {
        assertEq(address(awpEmission.awpToken()), address(awp));
        assertEq(awpEmission.cachedMaxSupply(), awp.MAX_SUPPLY());
    }
}
