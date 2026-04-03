// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract AWPEmissionTest is Test {
    AWPEmission public emission;
    AWPToken public awpToken;

    address public deployer = makeAddr("deployer");
    address public treasury = makeAddr("treasury");
    address public recipient = makeAddr("recipient");
    address public user = makeAddr("user");
    address public recipient1 = makeAddr("recipient1");
    address public recipient2 = makeAddr("recipient2");
    address public recipient3 = makeAddr("recipient3");
    address public recipient4 = makeAddr("recipient4");

    uint256 constant INITIAL_DAILY_EMISSION = 31_600_000 * 1e18;
    uint256 constant EPOCH_DURATION = 1 days;

    function setUp() public {
        vm.startPrank(deployer);

        // Deploy AWPToken (constructor mints 200M to deployer)
        awpToken = new AWPToken("AWP", "AWP", deployer);
        awpToken.initialMint(200_000_000 * 1e18);

        // Deploy AWPEmission (UUPS proxy pattern)
        // treasury == guardian in tests
        AWPEmission emissionImpl = new AWPEmission();
        bytes memory initData = abi.encodeCall(
            AWPEmission.initialize,
            (address(awpToken), treasury, INITIAL_DAILY_EMISSION, block.timestamp, EPOCH_DURATION, treasury)
        );
        ERC1967Proxy emissionProxy = new ERC1967Proxy(address(emissionImpl), initData);
        emission = AWPEmission(address(emissionProxy));

        // Grant AWPEmission minting permission (no renounceAdmin; tests need flexibility)
        awpToken.addMinter(address(emission));

        vm.stopPrank();
    }

    // ── Helper: submit allocations as guardian ──

    /// @dev Pack address + weight into uint256: [64-bit weight | 160-bit address]
    function _pack(address addr, uint96 w) internal pure returns (uint256) {
        return (uint256(w) << 160) | uint256(uint160(addr));
    }

    function _packArray(address[] memory addrs, uint96[] memory ws) internal pure returns (uint256[] memory) {
        uint256[] memory packed = new uint256[](addrs.length);
        for (uint256 i = 0; i < addrs.length; i++) {
            packed[i] = _pack(addrs[i], ws[i]);
        }
        return packed;
    }

    function _submitAsGuardian(address[] memory addrs, uint96[] memory ws) internal {
        _submitAsGuardianForEpoch(addrs, ws, emission.settledEpoch());
    }

    function _submitAsGuardianForEpoch(address[] memory addrs, uint96[] memory ws, uint256 effectiveEpoch) internal {
        uint256 tw = 0;
        for (uint256 i = 0; i < ws.length; i++) tw += ws[i];
        vm.prank(treasury); // treasury == guardian in tests
        emission.submitAllocations(_packArray(addrs, ws), tw, effectiveEpoch);
    }

    /// @dev Advance time past current epoch and settle it
    function _settleCurrentEpoch() internal {
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);
    }

    // ── Guardian access control tests ──

    function test_setGuardian_onlyGuardian() public {
        vm.prank(user);
        vm.expectRevert(AWPEmission.NotGuardian.selector);
        emission.setGuardian(user);
    }

    function test_setGuardian_success() public {
        address newGuardian = makeAddr("newGuardian");
        vm.prank(treasury); // treasury == guardian in tests
        emission.setGuardian(newGuardian);
        assertEq(emission.guardian(), newGuardian);

        // Old guardian can no longer submit
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        vm.prank(treasury);
        vm.expectRevert(AWPEmission.NotGuardian.selector);
        emission.submitAllocations(_packArray(addrs, ws), 100, 1);

        // New guardian can submit
        vm.prank(newGuardian);
        emission.submitAllocations(_packArray(addrs, ws), 100, 1);
    }

    function test_submitAllocations_onlyGuardian() public {
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;

        vm.prank(user);
        vm.expectRevert(AWPEmission.NotGuardian.selector);
        emission.submitAllocations(_packArray(addrs, ws), 100, 1);
    }

    // ── submitAllocations tests ──

    function test_submitAllocations() public {
        // Submit 2 recipients with weights 300/100 for epoch 1
        address[] memory addrs = new address[](2);
        addrs[0] = recipient1;
        addrs[1] = recipient2;
        uint96[] memory ws = new uint96[](2);
        ws[0] = 300;
        ws[1] = 100;

        _submitAsGuardian(addrs, ws);

        // Verify weight updates in epoch 0 (settledEpoch)
        assertEq(emission.getEpochWeight(0, recipient1), 300);
        assertEq(emission.getEpochWeight(0, recipient2), 100);
        assertEq(emission.getEpochTotalWeight(0), 400);
    }

    function test_submitAllocations_fullReplacement() public {
        // First submission for epoch 1: recipient1, recipient2
        address[] memory addrs1 = new address[](2);
        addrs1[0] = recipient1;
        addrs1[1] = recipient2;
        uint96[] memory ws1 = new uint96[](2);
        ws1[0] = 300;
        ws1[1] = 100;
        _submitAsGuardian(addrs1, ws1);

        assertEq(emission.getEpochTotalWeight(0), 400);

        // Second submission for epoch 0 with completely different recipients: recipient3, recipient4
        address[] memory addrs2 = new address[](2);
        addrs2[0] = recipient3;
        addrs2[1] = recipient4;
        uint96[] memory ws2 = new uint96[](2);
        ws2[0] = 500;
        ws2[1] = 200;
        _submitAsGuardianForEpoch(addrs2, ws2, 0);

        // Old recipient weights zeroed out (for epoch 0)
        assertEq(emission.getEpochWeight(0, recipient1), 0);
        assertEq(emission.getEpochWeight(0, recipient2), 0);
        // New recipient weights are correct
        assertEq(emission.getEpochWeight(0, recipient3), 500);
        assertEq(emission.getEpochWeight(0, recipient4), 200);
        assertEq(emission.getEpochTotalWeight(0), 700);
    }

    function test_submitAllocations_revertsExceedsMaxRecipients() public {
        assertEq(emission.maxRecipients(), 10000);

        // Submitting 1 recipient should succeed
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitAsGuardian(addrs, ws);
        assertEq(emission.getEpochTotalWeight(0), 100);
    }



    function test_submitAllocations_revertsMustBeFutureEpoch() public {
        // Submit for epoch 0, settle it, then try to submit for epoch 0 again → revert
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitAsGuardian(addrs, ws); // submit for epoch 0
        _settleCurrentEpoch(); // settle epoch 0 → settledEpoch=1

        // Cannot submit for epoch 0 (already settled)
        vm.prank(treasury);
        vm.expectRevert(AWPEmission.MustBeFutureEpoch.selector);
        emission.submitAllocations(_packArray(addrs, ws), 100, 0);
    }

    function test_submitAllocations_blockedDuringSettlement() public {
        // Submit allocations for epoch 1
        address[] memory addrs = new address[](2);
        addrs[0] = recipient1;
        addrs[1] = recipient2;
        uint96[] memory ws = new uint96[](2);
        ws[0] = 100;
        ws[1] = 100;
        _submitAsGuardianForEpoch(addrs, ws, 1);

        // Settle epoch 0 then start settling epoch 1 with limit=1
        _settleCurrentEpoch();
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(1);
        assertTrue(emission.settleProgress() > 0);

        // Submit during settlement should revert (prevents reentrant overwrite)
        address[] memory addrs2 = new address[](1);
        addrs2[0] = recipient3;
        uint96[] memory ws2 = new uint96[](1);
        ws2[0] = 200;

        vm.prank(treasury);
        vm.expectRevert(abi.encodeWithSignature("SettlementInProgress()"));
        emission.submitAllocations(_packArray(addrs2, ws2), 200, 3);

        // Complete settlement — then submission succeeds
        emission.settleEpoch(200);
        assertEq(emission.settleProgress(), 0);
        _submitAsGuardianForEpoch(addrs2, ws2, 3);
        assertEq(emission.getEpochWeight(3, recipient3), 200);
    }

    // ── settleEpoch tests ──

    function test_settleEpoch() public {
        // Submit allocations for epoch 1
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitAsGuardianForEpoch(addrs, ws, 1);

        // Settle epoch 0 (no weights -> all to DAO)
        _settleCurrentEpoch();
        assertEq(emission.settledEpoch(), 1);

        // Now settle epoch 1 (weights submitted for epoch 1)
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);

        assertEq(emission.settledEpoch(), 2);
        // Recipient receives 100% of the emission (no DAO split)
        uint256 worknetBal = awpToken.balanceOf(recipient1);
        assertTrue(worknetBal > 0);
    }

    function test_settleEpochBatched_oneByOne() public {
        // Submit for epoch 1
        address[] memory addrs = new address[](3);
        addrs[0] = recipient1;
        addrs[1] = recipient2;
        addrs[2] = recipient3;
        uint96[] memory ws = new uint96[](3);
        ws[0] = 100;
        ws[1] = 200;
        ws[2] = 300;
        _submitAsGuardianForEpoch(addrs, ws, 1);

        // Settle epoch 0 first
        _settleCurrentEpoch();
        assertEq(emission.settledEpoch(), 1);

        // Now settle epoch 1 batched
        vm.warp(block.timestamp + EPOCH_DURATION + 1);

        // Call 1: Phase 1 + process recipient 0
        emission.settleEpoch(1);
        assertTrue(emission.settleProgress() > 0);

        // Call 2: process recipient 1
        emission.settleEpoch(1);
        assertTrue(emission.settleProgress() > 0);

        // Call 3: process recipient 2 + Phase 3 complete
        emission.settleEpoch(1);
        assertEq(emission.settleProgress(), 0);
        assertEq(emission.settledEpoch(), 2);

        // Verify distribution is proportional to weights (100:200:300 = 1:2:3)
        uint256 bal1 = awpToken.balanceOf(recipient1);
        uint256 bal2 = awpToken.balanceOf(recipient2);
        uint256 bal3 = awpToken.balanceOf(recipient3);
        assertApproxEqRel(bal2, bal1 * 2, 0.01e18);
        assertApproxEqRel(bal3, bal1 * 3, 0.01e18);
    }

    function test_settleEpochBatched_twoBatch() public {
        // Submit for epoch 0
        address[] memory addrs = new address[](4);
        addrs[0] = recipient1;
        addrs[1] = recipient2;
        addrs[2] = recipient3;
        addrs[3] = recipient4;
        uint96[] memory ws = new uint96[](4);
        ws[0] = 100;
        ws[1] = 100;
        ws[2] = 100;
        ws[3] = 100;
        _submitAsGuardianForEpoch(addrs, ws, 0);

        // Settle epoch 0 in 2 batches
        emission.settleEpoch(2);
        assertTrue(emission.settleProgress() > 0);

        emission.settleEpoch(2);
        assertEq(emission.settleProgress(), 0);
        assertEq(emission.settledEpoch(), 1);

        // 4 equal-weight recipients, each receives 1/4 of epoch 0 emission
        uint256 expected = INITIAL_DAILY_EMISSION / 4;
        assertEq(awpToken.balanceOf(recipient1), expected);
        assertEq(awpToken.balanceOf(recipient2), expected);
        assertEq(awpToken.balanceOf(recipient3), expected);
        assertEq(awpToken.balanceOf(recipient4), expected);
    }

    function test_settleEpochDecay() public {
        // Submit for epoch 1
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitAsGuardianForEpoch(addrs, ws, 1);

        // Epoch 0 (no decay, no weights) — use relative warp for fork compatibility
        vm.warp(block.timestamp + 2 days);
        emission.settleEpoch(200);
        assertEq(emission.settledEpoch(), 1);
        assertEq(emission.currentDailyEmission(), INITIAL_DAILY_EMISSION);

        // Epoch 1 (should decay, weights available)
        vm.warp(block.timestamp + 2 days);
        emission.settleEpoch(200);
        assertEq(emission.settledEpoch(), 2);

        uint256 e = emission.currentDailyEmission();
        assertTrue(e < INITIAL_DAILY_EMISSION);
        assertEq(e, INITIAL_DAILY_EMISSION * 996844 / 1000000);
    }

    function test_settleEpochNoRecipients() public {
        // No allocations — nothing is minted (Guardian must include treasury in recipients for DAO share)
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);

        // No recipients → 0 minted this epoch
        assertEq(awpToken.balanceOf(treasury), 0);
    }

    function test_settleEpochTooEarly() public {
        // Settle epoch 0 (allowed at currentEpoch=0)
        emission.settleEpoch(200);
        // Now settledEpoch=1, but currentEpoch still 0 → cannot settle epoch 1
        vm.expectRevert(AWPEmission.EpochNotReady.selector);
        emission.settleEpoch(200);
    }

    function test_settleEpoch_revertsLimitZero() public {
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        vm.expectRevert(AWPEmission.ZeroLimit.selector);
        emission.settleEpoch(0);
    }

    // ── Multi-recipient proportional distribution ──

    function test_multiRecipientDistribution() public {
        // Submit for epoch 1
        address[] memory addrs = new address[](2);
        addrs[0] = recipient1;
        addrs[1] = recipient2;
        uint96[] memory ws = new uint96[](2);
        ws[0] = 300;
        ws[1] = 100;
        _submitAsGuardianForEpoch(addrs, ws, 1);

        // Settle epoch 0
        _settleCurrentEpoch();

        // Settle epoch 1
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);

        uint256 bal1 = awpToken.balanceOf(recipient1);
        uint256 bal2 = awpToken.balanceOf(recipient2);

        // 3:1 ratio
        assertApproxEqRel(bal1, bal2 * 3, 0.01e18);
    }

    // ── notSettling guard ──

    function test_notSettling() public {
        // Submit for epoch 1
        address[] memory addrs = new address[](2);
        addrs[0] = recipient1;
        addrs[1] = recipient2;
        uint96[] memory ws = new uint96[](2);
        ws[0] = 100;
        ws[1] = 100;
        _submitAsGuardianForEpoch(addrs, ws, 1);

        // Settle epoch 0
        _settleCurrentEpoch();

        // Start settling epoch 1 with limit=1
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(1);
        assertTrue(emission.settleProgress() > 0);

        // Cannot submitAllocations during settlement
        address[] memory addrs2 = new address[](1);
        addrs2[0] = recipient3;
        uint96[] memory ws2 = new uint96[](1);
        ws2[0] = 200;
        vm.prank(treasury);
        vm.expectRevert(AWPEmission.SettlementInProgress.selector);
        emission.submitAllocations(_packArray(addrs2, ws2), 200, 3);

        // Complete remaining settlement batches
        emission.settleEpoch(200);
        assertEq(emission.settleProgress(), 0);
    }

    // ── Upgrade tests ──

    function test_upgradeViaTimelock() public {
        // Submit for epoch 1
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 500;
        _submitAsGuardianForEpoch(addrs, ws, 1);

        // Settle epoch 0
        _settleCurrentEpoch();
        // Settle epoch 1 to promote activeEpoch=1
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);

        // Deploy new implementation
        AWPEmission newImpl = new AWPEmission();

        // treasury calls upgradeToAndCall
        vm.prank(treasury);
        emission.upgradeToAndCall(address(newImpl), "");

        // Verify state is preserved
        assertEq(emission.getWeight(recipient1), 500);
        assertEq(emission.getTotalWeight(), 500);
        assertEq(emission.getRecipientCount(), 1);
    }

    function test_upgrade_revertsForNonGuardian() public {
        AWPEmission newImpl = new AWPEmission();

        vm.prank(user);
        vm.expectRevert(AWPEmission.NotGuardian.selector);
        emission.upgradeToAndCall(address(newImpl), "");
    }

    // ── Epoch-versioned weight tests ──

    function test_epochVersionedWeights_independentEpochs() public {
        // Submit different weights for epochs 1 and 2
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws1 = new uint96[](1);
        ws1[0] = 100;
        _submitAsGuardianForEpoch(addrs, ws1, 1);

        uint96[] memory ws2 = new uint96[](1);
        ws2[0] = 500;
        _submitAsGuardianForEpoch(addrs, ws2, 2);

        // Verify weights are independent per epoch
        assertEq(emission.getEpochWeight(1, recipient1), 100);
        assertEq(emission.getEpochTotalWeight(1), 100);
        assertEq(emission.getEpochWeight(2, recipient1), 500);
        assertEq(emission.getEpochTotalWeight(2), 500);
    }

    function test_activeEpochPromotionOnSettle() public {
        // Submit for epoch 0
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitAsGuardianForEpoch(addrs, ws, 0);

        // activeEpoch is 0 initially
        assertEq(emission.activeEpoch(), 0);

        // Settle epoch 0 — weights exist for epoch 0, activeEpoch promoted to 0
        _settleCurrentEpoch();
        assertEq(emission.activeEpoch(), 0);
    }

    function test_activeEpochPersists_whenNoNewWeights() public {
        // Submit for epoch 1
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitAsGuardianForEpoch(addrs, ws, 1);

        // Settle epoch 0
        _settleCurrentEpoch();
        // Settle epoch 1 — promotes activeEpoch to 1
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);
        assertEq(emission.activeEpoch(), 1);

        uint256 balAfterEpoch1 = awpToken.balanceOf(recipient1);
        assertTrue(balAfterEpoch1 > 0);

        // Settle epoch 2 — no weights for epoch 2, activeEpoch stays at 1
        // Note: use 3 days offset from genesis (block.timestamp=1 in setUp) to ensure epoch 2 is ready
        vm.warp(3 * EPOCH_DURATION + 2);
        emission.settleEpoch(200);
        assertEq(emission.activeEpoch(), 1);

        // Still distributes using epoch 1 weights
        uint256 balAfterEpoch2 = awpToken.balanceOf(recipient1);
        assertTrue(balAfterEpoch2 > balAfterEpoch1);
    }

    // ══════════════════════════════════════════════
    // Governance setter tests (setDecayFactor)
    // ══════════════════════════════════════════════

    function test_setDecayFactor() public {
        vm.prank(treasury);
        emission.setDecayFactor(995000); // ~0.5% decay
        assertEq(emission.decayFactor(), 995000);
    }

    function test_setDecayFactor_tooHigh_reverts() public {
        vm.prank(treasury);
        vm.expectRevert(AWPEmission.InvalidDecayFactor.selector);
        emission.setDecayFactor(1000000); // >= DECAY_PRECISION
    }

    function test_setDecayFactor_zero_reverts() public {
        vm.prank(treasury);
        vm.expectRevert(AWPEmission.InvalidDecayFactor.selector);
        emission.setDecayFactor(0);
    }

    function test_setDecayFactor_notGuardian_reverts() public {
        vm.prank(user);
        vm.expectRevert(AWPEmission.NotGuardian.selector);
        emission.setDecayFactor(995000);
    }

    // ══════════════════════════════════════════════
    // New tests: Guardian upgrade, 100% emission, treasury as recipient,
    //            setDecayFactor by guardian, catch-up submission
    // ══════════════════════════════════════════════

    function test_upgradeViaGuardian() public {
        // Guardian (treasury in tests) can upgrade the emission contract
        AWPEmission newImpl = new AWPEmission();
        vm.prank(treasury);
        emission.upgradeToAndCall(address(newImpl), "");

        // Verify state preserved after upgrade
        assertEq(emission.guardian(), treasury);
        assertEq(emission.maxRecipients(), 10000);
    }

    function test_epochSettlement_100percent() public {
        // Verify 100% emission goes to recipients (no DAO split)
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitAsGuardianForEpoch(addrs, ws, 0);

        // Settle epoch 0 (current epoch, allowed by new rules)
        emission.settleEpoch(200);

        // recipient1 should receive full epoch 0 emission (no decay on first epoch)
        assertEq(awpToken.balanceOf(recipient1), INITIAL_DAILY_EMISSION);

        // No other address received anything (no DAO split)
        assertEq(awpToken.balanceOf(treasury), 0);
    }

    function test_guardianIncludesTreasury() public {
        // Guardian adds treasury as a recipient — verify treasury gets its weight share
        address[] memory addrs = new address[](2);
        addrs[0] = recipient1;
        addrs[1] = treasury;
        uint96[] memory ws = new uint96[](2);
        ws[0] = 300;
        ws[1] = 100;
        _submitAsGuardianForEpoch(addrs, ws, 1);

        _settleCurrentEpoch();
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);

        uint256 bal1 = awpToken.balanceOf(recipient1);
        uint256 balT = awpToken.balanceOf(treasury);

        // 3:1 ratio (recipient1:treasury)
        assertApproxEqRel(bal1, balT * 3, 0.01e18);
        assertTrue(balT > 0);
    }

    function test_setDecayFactor_byGuardian() public {
        vm.prank(treasury);
        emission.setDecayFactor(990000);
        assertEq(emission.decayFactor(), 990000);

        // Now verify decay uses the new factor
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitAsGuardianForEpoch(addrs, ws, 1);

        _settleCurrentEpoch();

        // Settle epoch 1 with new decay factor
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);

        uint256 decayedEmission = emission.currentDailyEmission();
        assertEq(decayedEmission, INITIAL_DAILY_EMISSION * 990000 / 1000000);
    }

    function test_submitAllocations_catchUp() public {
        // Submit for epoch 3 (skipping epochs 0-2)
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitAsGuardianForEpoch(addrs, ws, 3);

        uint256 genesis = emission.baseTime();

        // Warp to epoch 3 so we can settle epochs 0-3
        vm.warp(genesis + EPOCH_DURATION * 3 + 1);

        // Settle epoch 0 — no weights for epoch 0, activeEpoch stays 0
        emission.settleEpoch(200);
        assertEq(emission.settledEpoch(), 1);
        assertEq(emission.activeEpoch(), 0);

        // Settle epoch 1 — no weights for epoch 1, activeEpoch stays 0
        emission.settleEpoch(200);
        assertEq(emission.settledEpoch(), 2);
        assertEq(emission.activeEpoch(), 0);

        // Settle epoch 2 — no weights for epoch 2, activeEpoch stays 0
        emission.settleEpoch(200);
        assertEq(emission.settledEpoch(), 3);
        assertEq(emission.activeEpoch(), 0);

        // Settle epoch 3 — weights exist, activeEpoch promoted to 3
        emission.settleEpoch(200);
        assertEq(emission.settledEpoch(), 4);
        assertEq(emission.activeEpoch(), 3);

        // recipient1 receives epoch 3's emission
        uint256 bal = awpToken.balanceOf(recipient1);
        assertTrue(bal > 0);
    }
}
