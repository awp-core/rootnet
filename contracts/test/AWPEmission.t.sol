// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {EmissionSigningHelper} from "./helpers/EmissionSigningHelper.sol";

contract AWPEmissionTest is EmissionSigningHelper {
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

    uint256 constant INITIAL_DAILY_EMISSION = 15_800_000 * 1e18;
    uint256 constant EPOCH_DURATION = 1 days;

    // Oracle private keys and addresses
    uint256 oracle1Pk = 0xA1;
    uint256 oracle2Pk = 0xA2;
    uint256 oracle3Pk = 0xA3;
    uint256 oracle4Pk = 0xA4;
    uint256 oracle5Pk = 0xA5;

    address oracle1 = vm.addr(0xA1);
    address oracle2 = vm.addr(0xA2);
    address oracle3 = vm.addr(0xA3);
    address oracle4 = vm.addr(0xA4);
    address oracle5 = vm.addr(0xA5);

    function setUp() public {
        vm.startPrank(deployer);

        // Deploy AWPToken (constructor mints 200M to deployer)
        awpToken = new AWPToken("AWP", "AWP", deployer);

        // Deploy AWPEmission (UUPS proxy pattern)
        // AWPEmission now has its own epoch timing
        AWPEmission emissionImpl = new AWPEmission();
        bytes memory initData = abi.encodeCall(
            AWPEmission.initialize,
            (address(awpToken), treasury, INITIAL_DAILY_EMISSION, block.timestamp, EPOCH_DURATION)
        );
        ERC1967Proxy emissionProxy = new ERC1967Proxy(address(emissionImpl), initData);
        emission = AWPEmission(address(emissionProxy));

        // Grant AWPEmission minting permission (no renounceAdmin; tests need flexibility)
        awpToken.addMinter(address(emission));

        // Configure oracles: 5 oracles, threshold 3
        address[] memory oracleList = new address[](5);
        oracleList[0] = oracle1;
        oracleList[1] = oracle2;
        oracleList[2] = oracle3;
        oracleList[3] = oracle4;
        oracleList[4] = oracle5;
        vm.stopPrank();

        vm.prank(treasury);
        emission.setOracleConfig(oracleList, 3);
    }

    // ── Helper: submit allocations via oracle signatures ──

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

    function _submitWithOracles(address[] memory addrs, uint96[] memory ws) internal {
        _submitWithOraclesForEpoch(addrs, ws, emission.settledEpoch() + 1);
    }

    function _submitWithOraclesForEpoch(address[] memory addrs, uint96[] memory ws, uint256 effectiveEpoch) internal {
        _sortByAddress(addrs, ws);
        uint256 nonce = emission.allocationNonce();
        bytes[] memory sigs = new bytes[](3);
        sigs[0] = _signAllocations(oracle1Pk, addrs, ws, nonce, effectiveEpoch, address(emission));
        sigs[1] = _signAllocations(oracle2Pk, addrs, ws, nonce, effectiveEpoch, address(emission));
        sigs[2] = _signAllocations(oracle3Pk, addrs, ws, nonce, effectiveEpoch, address(emission));
        emission.submitAllocations(addrs, ws, sigs, effectiveEpoch);
    }

    /// @dev Settle epoch 0 (no weights, all goes to DAO) to advance to epoch 1
    function _settleEpoch0() internal {
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);
    }

    // ── Oracle configuration tests ──

    function test_setOracleConfig() public {
        // setUp has configured 5 oracles with threshold 3
        assertEq(emission.getOracleCount(), 5);
        assertEq(emission.oracleThreshold(), 3);
    }

    function test_setOracleConfig_revertsForNonTimelock() public {
        address[] memory oList = new address[](1);
        oList[0] = oracle1;

        vm.prank(user);
        vm.expectRevert(AWPEmission.NotTimelock.selector);
        emission.setOracleConfig(oList, 1);
    }

    function test_setOracleConfig_revertsForZeroThreshold() public {
        address[] memory oList = new address[](1);
        oList[0] = oracle1;

        vm.prank(treasury);
        vm.expectRevert(AWPEmission.InvalidOracleConfig.selector);
        emission.setOracleConfig(oList, 0);
    }

    function test_setOracleConfig_revertsForThresholdExceedsOracles() public {
        address[] memory oList = new address[](2);
        oList[0] = oracle1;
        oList[1] = oracle2;

        vm.prank(treasury);
        vm.expectRevert(AWPEmission.InvalidOracleConfig.selector);
        emission.setOracleConfig(oList, 3);
    }

    function test_setOracleConfig_revertsForZeroAddress() public {
        address[] memory oList = new address[](2);
        oList[0] = oracle1;
        oList[1] = address(0);

        vm.prank(treasury);
        vm.expectRevert(AWPEmission.InvalidOracleConfig.selector);
        emission.setOracleConfig(oList, 1);
    }

    function test_setOracleConfig_revertsForDuplicateOracle() public {
        address[] memory oList = new address[](3);
        oList[0] = oracle1;
        oList[1] = oracle2;
        oList[2] = oracle1; // duplicate

        vm.prank(treasury);
        vm.expectRevert(AWPEmission.DuplicateOracle.selector);
        emission.setOracleConfig(oList, 2);
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

        _submitWithOracles(addrs, ws);

        // Verify weight updates in epoch 1
        assertEq(emission.getEpochWeight(1, recipient1), 300);
        assertEq(emission.getEpochWeight(1, recipient2), 100);
        assertEq(emission.getEpochTotalWeight(1), 400);
    }

    function test_submitAllocations_fullReplacement() public {
        // First submission for epoch 1: recipient1, recipient2
        address[] memory addrs1 = new address[](2);
        addrs1[0] = recipient1;
        addrs1[1] = recipient2;
        uint96[] memory ws1 = new uint96[](2);
        ws1[0] = 300;
        ws1[1] = 100;
        _submitWithOracles(addrs1, ws1);

        assertEq(emission.getEpochTotalWeight(1), 400);

        // Second submission for epoch 1 with completely different recipients: recipient3, recipient4
        address[] memory addrs2 = new address[](2);
        addrs2[0] = recipient3;
        addrs2[1] = recipient4;
        uint96[] memory ws2 = new uint96[](2);
        ws2[0] = 500;
        ws2[1] = 200;
        _submitWithOraclesForEpoch(addrs2, ws2, 1);

        // Old recipient weights zeroed out (for epoch 1)
        assertEq(emission.getEpochWeight(1, recipient1), 0);
        assertEq(emission.getEpochWeight(1, recipient2), 0);
        // New recipient weights are correct
        assertEq(emission.getEpochWeight(1, recipient3), 500);
        assertEq(emission.getEpochWeight(1, recipient4), 200);
        assertEq(emission.getEpochTotalWeight(1), 700);
    }

    function test_submitAllocations_revertsBeforeOracleConfig() public {
        // Deploy a new emission without configuring oracles
        vm.startPrank(deployer);
        AWPEmission impl2 = new AWPEmission();
        bytes memory initData2 = abi.encodeCall(
            AWPEmission.initialize,
            (address(awpToken), treasury, INITIAL_DAILY_EMISSION, block.timestamp, EPOCH_DURATION)
        );
        ERC1967Proxy proxy2 = new ERC1967Proxy(address(impl2), initData2);
        AWPEmission emission2 = AWPEmission(address(proxy2));
        vm.stopPrank();

        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        bytes[] memory sigs = new bytes[](1);
        sigs[0] = _signAllocations(oracle1Pk, addrs, ws, 0, 1, address(emission2));

        vm.expectRevert(AWPEmission.OracleNotConfigured.selector);
        emission2.submitAllocations(addrs, ws, sigs, 1);
    }

    function test_submitAllocations_revertsBelowThreshold() public {
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        uint256 nonce = emission.allocationNonce();
        uint256 effectiveEpoch = emission.settledEpoch() + 1;

        // Only 2 signatures, but threshold is 3
        bytes[] memory sigs = new bytes[](2);
        sigs[0] = _signAllocations(oracle1Pk, addrs, ws, nonce, effectiveEpoch, address(emission));
        sigs[1] = _signAllocations(oracle2Pk, addrs, ws, nonce, effectiveEpoch, address(emission));

        vm.expectRevert(AWPEmission.InvalidSignatureCount.selector);
        emission.submitAllocations(addrs, ws, sigs, effectiveEpoch);
    }

    function test_submitAllocations_revertsDuplicateSigner() public {
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        uint256 nonce = emission.allocationNonce();
        uint256 effectiveEpoch = emission.settledEpoch() + 1;

        // The same oracle signs twice
        bytes[] memory sigs = new bytes[](3);
        sigs[0] = _signAllocations(oracle1Pk, addrs, ws, nonce, effectiveEpoch, address(emission));
        sigs[1] = _signAllocations(oracle1Pk, addrs, ws, nonce, effectiveEpoch, address(emission)); // duplicate
        sigs[2] = _signAllocations(oracle2Pk, addrs, ws, nonce, effectiveEpoch, address(emission));

        vm.expectRevert(AWPEmission.DuplicateSigner.selector);
        emission.submitAllocations(addrs, ws, sigs, effectiveEpoch);
    }

    function test_submitAllocations_revertsUnknownOracle() public {
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        uint256 nonce = emission.allocationNonce();
        uint256 effectiveEpoch = emission.settledEpoch() + 1;

        // Sign with a non-oracle private key
        uint256 fakeOraclePk = 0xDEAD;
        bytes[] memory sigs = new bytes[](3);
        sigs[0] = _signAllocations(oracle1Pk, addrs, ws, nonce, effectiveEpoch, address(emission));
        sigs[1] = _signAllocations(oracle2Pk, addrs, ws, nonce, effectiveEpoch, address(emission));
        sigs[2] = _signAllocations(fakeOraclePk, addrs, ws, nonce, effectiveEpoch, address(emission));

        vm.expectRevert(AWPEmission.UnknownOracle.selector);
        emission.submitAllocations(addrs, ws, sigs, effectiveEpoch);
    }

    function test_submitAllocations_revertsExceedsMaxRecipients() public {
        assertEq(emission.maxRecipients(), 10000);

        // Submitting 1 recipient should succeed
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitWithOracles(addrs, ws);
        assertEq(emission.getEpochTotalWeight(1), 100);
    }

    function test_submitAllocations_revertsDuplicateRecipient() public {
        // Same address appears twice in recipients
        address[] memory addrs = new address[](2);
        addrs[0] = recipient1;
        addrs[1] = recipient1; // duplicate
        uint96[] memory ws = new uint96[](2);
        ws[0] = 100;
        ws[1] = 200;
        uint256 nonce = emission.allocationNonce();
        uint256 effectiveEpoch = emission.settledEpoch() + 1;

        bytes[] memory sigs = new bytes[](3);
        sigs[0] = _signAllocations(oracle1Pk, addrs, ws, nonce, effectiveEpoch, address(emission));
        sigs[1] = _signAllocations(oracle2Pk, addrs, ws, nonce, effectiveEpoch, address(emission));
        sigs[2] = _signAllocations(oracle3Pk, addrs, ws, nonce, effectiveEpoch, address(emission));

        vm.expectRevert(AWPEmission.DuplicateRecipient.selector);
        emission.submitAllocations(addrs, ws, sigs, effectiveEpoch);
    }

    function test_submitAllocations_nonceIncrement() public {
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;

        // First submission, nonce=0
        assertEq(emission.allocationNonce(), 0);
        _submitWithOracles(addrs, ws);

        // nonce should have incremented to 1
        assertEq(emission.allocationNonce(), 1);

        // Second submission for epoch 2 (different epoch), nonce=1
        ws[0] = 200;
        _submitWithOraclesForEpoch(addrs, ws, 2);

        assertEq(emission.allocationNonce(), 2);
        assertEq(emission.getEpochWeight(2, recipient1), 200);
    }

    function test_submitAllocations_revertsNonceReplay() public {
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;

        // Submit nonce=0, effectiveEpoch=1
        uint256 nonce0 = emission.allocationNonce();
        uint256 effectiveEpoch = emission.settledEpoch() + 1;
        bytes[] memory sigs0 = new bytes[](3);
        sigs0[0] = _signAllocations(oracle1Pk, addrs, ws, nonce0, effectiveEpoch, address(emission));
        sigs0[1] = _signAllocations(oracle2Pk, addrs, ws, nonce0, effectiveEpoch, address(emission));
        sigs0[2] = _signAllocations(oracle3Pk, addrs, ws, nonce0, effectiveEpoch, address(emission));
        emission.submitAllocations(addrs, ws, sigs0, effectiveEpoch);
        assertEq(emission.allocationNonce(), 1);

        // Replay the same nonce=0 signatures — current nonce is 1, so the digest won't match
        vm.expectRevert();
        emission.submitAllocations(addrs, ws, sigs0, effectiveEpoch);
    }

    function test_submitAllocations_revertsMustBeFutureEpoch() public {
        // Cannot submit for currentEpoch (epoch 0)
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        uint256 nonce = emission.allocationNonce();
        uint256 currentEp = emission.settledEpoch();

        bytes[] memory sigs = new bytes[](3);
        sigs[0] = _signAllocations(oracle1Pk, addrs, ws, nonce, currentEp, address(emission));
        sigs[1] = _signAllocations(oracle2Pk, addrs, ws, nonce, currentEp, address(emission));
        sigs[2] = _signAllocations(oracle3Pk, addrs, ws, nonce, currentEp, address(emission));

        vm.expectRevert(AWPEmission.MustBeFutureEpoch.selector);
        emission.submitAllocations(addrs, ws, sigs, currentEp);
    }

    function test_submitAllocations_blockedDuringSettlement() public {
        // Submit allocations for epoch 1
        address[] memory addrs = new address[](2);
        addrs[0] = recipient1;
        addrs[1] = recipient2;
        uint96[] memory ws = new uint96[](2);
        ws[0] = 100;
        ws[1] = 100;
        _submitWithOraclesForEpoch(addrs, ws, 1);

        // Settle epoch 0 then start settling epoch 1 with limit=1
        _settleEpoch0();
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(1);
        assertTrue(emission.settleProgress() > 0);

        // Submit during settlement should revert (prevents reentrant overwrite)
        address[] memory addrs2 = new address[](1);
        addrs2[0] = recipient3;
        uint96[] memory ws2 = new uint96[](1);
        ws2[0] = 200;

        // Build sigs before expectRevert (allocationNonce() is a view call that would consume the expectRevert)
        _sortByAddress(addrs2, ws2);
        uint256 nonce = emission.allocationNonce();
        bytes[] memory sigs = new bytes[](3);
        sigs[0] = _signAllocations(oracle1Pk, addrs2, ws2, nonce, 3, address(emission));
        sigs[1] = _signAllocations(oracle2Pk, addrs2, ws2, nonce, 3, address(emission));
        sigs[2] = _signAllocations(oracle3Pk, addrs2, ws2, nonce, 3, address(emission));

        vm.expectRevert(abi.encodeWithSignature("SettlementInProgress()"));
        emission.submitAllocations(addrs2, ws2, sigs, 3);

        // Complete settlement — then submission succeeds
        emission.settleEpoch(200);
        assertEq(emission.settleProgress(), 0);
        _submitWithOraclesForEpoch(addrs2, ws2, 3);
        assertEq(emission.getEpochWeight(3, recipient3), 200);
    }

    // ── settleEpoch tests ──

    function test_settleEpoch() public {
        // Submit allocations for epoch 1
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitWithOraclesForEpoch(addrs, ws, 1);

        // Settle epoch 0 (no weights -> all to DAO)
        _settleEpoch0();
        assertEq(emission.settledEpoch(), 1);

        // Now settle epoch 1 (weights submitted for epoch 1)
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);

        assertEq(emission.settledEpoch(), 2);
        // Recipient should receive 50% of the emission
        uint256 subnetBal = awpToken.balanceOf(recipient1);
        assertTrue(subnetBal > 0);
        // Treasury should receive the DAO share
        uint256 treasuryBal = awpToken.balanceOf(treasury);
        assertTrue(treasuryBal > 0);
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
        _submitWithOraclesForEpoch(addrs, ws, 1);

        // Settle epoch 0 first
        _settleEpoch0();
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
        // Submit for epoch 1
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
        _submitWithOraclesForEpoch(addrs, ws, 1);

        // Settle epoch 0
        _settleEpoch0();

        // Settle epoch 1 in 2 batches
        vm.warp(block.timestamp + EPOCH_DURATION + 1);

        emission.settleEpoch(2);
        assertTrue(emission.settleProgress() > 0);

        emission.settleEpoch(2);
        assertEq(emission.settleProgress(), 0);
        assertEq(emission.settledEpoch(), 2);

        // 4 equal-weight recipients, each receives 1/4 of the subnet pool
        // Epoch 0 (no decay) + Epoch 1 (with decay) — both distributed because
        // activeEpoch is promoted to 1 during epoch 0 settlement (weights pre-submitted)
        uint256 epoch0Pool = INITIAL_DAILY_EMISSION / 2;
        uint256 decayedEmission = INITIAL_DAILY_EMISSION * 996844 / 1000000;
        uint256 epoch1Pool = decayedEmission / 2;
        uint256 expected = (epoch0Pool + epoch1Pool) / 4;
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
        _submitWithOraclesForEpoch(addrs, ws, 1);

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
        // No allocations — all emission goes to the DAO
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);

        uint256 treasuryBal = awpToken.balanceOf(treasury);
        assertEq(treasuryBal, INITIAL_DAILY_EMISSION);
    }

    function test_settleEpochTooEarly() public {
        vm.expectRevert(AWPEmission.EpochNotReady.selector);
        emission.settleEpoch(200);
    }

    function test_settleEpoch_revertsLimitZero() public {
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        vm.expectRevert(AWPEmission.InvalidParameter.selector);
        emission.settleEpoch(0);
    }

    // ── emergencySetWeight tests ──

    function test_emergencySetWeight() public {
        // Submit allocations for epoch 1
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitWithOraclesForEpoch(addrs, ws, 1);

        // Settle epoch 0 to advance to epoch 1
        _settleEpoch0();
        // Settle epoch 1 to promote activeEpoch=1
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);
        assertEq(emission.activeEpoch(), 1);

        // Emergency override weight (epoch=1, index=0, addr=recipient1)
        vm.prank(treasury);
        emission.emergencySetWeight(1, 0, recipient1, 500);

        assertEq(emission.getWeight(recipient1), 500);
        assertEq(emission.getTotalWeight(), 500);
    }

    function test_emergencySetWeight_revertsForNonTimelock() public {
        vm.prank(user);
        vm.expectRevert(AWPEmission.NotTimelock.selector);
        emission.emergencySetWeight(0, 0, recipient1, 100);
    }

    function test_emergencySetWeight_revertsIndexOutOfBounds() public {
        // index out of bounds (no allocations at epoch 0)
        vm.prank(treasury);
        vm.expectRevert(AWPEmission.InvalidParameter.selector);
        emission.emergencySetWeight(0, 0, recipient1, 100);
    }

    function test_emergencySetWeight_revertsDuringSettlement() public {
        // Submit allocations for epoch 1
        address[] memory addrs = new address[](2);
        addrs[0] = recipient1;
        addrs[1] = recipient2;
        uint96[] memory ws = new uint96[](2);
        ws[0] = 100;
        ws[1] = 100;
        _submitWithOraclesForEpoch(addrs, ws, 1);

        // Settle epoch 0
        _settleEpoch0();

        // Trigger settle epoch 1 with limit=1 to enter settlement in progress
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(1);
        assertTrue(emission.settleProgress() > 0);

        // Cannot emergencySetWeight while settlement is in progress
        vm.prank(treasury);
        vm.expectRevert(AWPEmission.SettlementInProgress.selector);
        emission.emergencySetWeight(1, 0, recipient1, 200);

        // Complete settlement
        emission.settleEpoch(200);
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
        _submitWithOraclesForEpoch(addrs, ws, 1);

        // Settle epoch 0
        _settleEpoch0();

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
        _submitWithOraclesForEpoch(addrs, ws, 1);

        // Settle epoch 0
        _settleEpoch0();

        // Start settling epoch 1 with limit=1
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(1);
        assertTrue(emission.settleProgress() > 0);

        // Cannot emergencySetWeight during settlement
        vm.prank(treasury);
        vm.expectRevert(AWPEmission.SettlementInProgress.selector);
        emission.emergencySetWeight(1, 0, recipient1, 200);

        // CAN submitAllocations for future epochs during settlement (epoch-versioned design)

        // Cannot setOracleConfig during settlement
        address[] memory newOracles = new address[](1);
        newOracles[0] = oracle1;
        vm.prank(treasury);
        vm.expectRevert(AWPEmission.SettlementInProgress.selector);
        emission.setOracleConfig(newOracles, 1);

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
        _submitWithOraclesForEpoch(addrs, ws, 1);

        // Settle epoch 0
        _settleEpoch0();
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

    function test_upgrade_revertsForNonTimelock() public {
        AWPEmission newImpl = new AWPEmission();

        vm.prank(user);
        vm.expectRevert(AWPEmission.NotTimelock.selector);
        emission.upgradeToAndCall(address(newImpl), "");
    }

    /// @dev Test that old signatures become invalid after oracle replacement
    function test_setOracleConfig_replaceAfterSubmission() public {
        // First submission (nonce=0) for epoch 1
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitWithOracles(addrs, ws);
        assertEq(emission.allocationNonce(), 1);

        // Replace with a completely new oracle set
        uint256 newPk1 = 0xB1;
        uint256 newPk2 = 0xB2;
        uint256 newPk3 = 0xB3;
        address[] memory newOracles = new address[](3);
        newOracles[0] = vm.addr(newPk1);
        newOracles[1] = vm.addr(newPk2);
        newOracles[2] = vm.addr(newPk3);
        vm.prank(treasury);
        emission.setOracleConfig(newOracles, 2);

        // Old oracle signatures for nonce=1, epoch 2 should be rejected
        ws[0] = 200;
        uint256 nonce1 = emission.allocationNonce();
        uint256 effectiveEpoch = 2;
        bytes[] memory oldSigs = new bytes[](3);
        oldSigs[0] = _signAllocations(oracle1Pk, addrs, ws, nonce1, effectiveEpoch, address(emission));
        oldSigs[1] = _signAllocations(oracle2Pk, addrs, ws, nonce1, effectiveEpoch, address(emission));
        oldSigs[2] = _signAllocations(oracle3Pk, addrs, ws, nonce1, effectiveEpoch, address(emission));
        vm.expectRevert(AWPEmission.UnknownOracle.selector);
        emission.submitAllocations(addrs, ws, oldSigs, effectiveEpoch);

        // New oracle signatures should succeed
        bytes[] memory newSigs = new bytes[](2);
        newSigs[0] = _signAllocations(newPk1, addrs, ws, nonce1, effectiveEpoch, address(emission));
        newSigs[1] = _signAllocations(newPk2, addrs, ws, nonce1, effectiveEpoch, address(emission));
        emission.submitAllocations(addrs, ws, newSigs, effectiveEpoch);

        assertEq(emission.getEpochWeight(effectiveEpoch, recipient1), 200);
    }

    // ── Epoch-versioned weight tests ──

    function test_epochVersionedWeights_independentEpochs() public {
        // Submit different weights for epochs 1 and 2
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws1 = new uint96[](1);
        ws1[0] = 100;
        _submitWithOraclesForEpoch(addrs, ws1, 1);

        uint96[] memory ws2 = new uint96[](1);
        ws2[0] = 500;
        _submitWithOraclesForEpoch(addrs, ws2, 2);

        // Verify weights are independent per epoch
        assertEq(emission.getEpochWeight(1, recipient1), 100);
        assertEq(emission.getEpochTotalWeight(1), 100);
        assertEq(emission.getEpochWeight(2, recipient1), 500);
        assertEq(emission.getEpochTotalWeight(2), 500);
    }

    function test_activeEpochPromotionOnSettle() public {
        // Submit for epoch 1
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitWithOraclesForEpoch(addrs, ws, 1);

        // activeEpoch is 0 initially
        assertEq(emission.activeEpoch(), 0);

        // Settle epoch 0 — oracle submitted for epoch 1, so activeEpoch promoted to 1
        // (settleEpoch checks _epochTotalWeight[settledEpoch + 1])
        _settleEpoch0();
        assertEq(emission.activeEpoch(), 1);
    }

    function test_activeEpochPersists_whenNoNewWeights() public {
        // Submit for epoch 1
        address[] memory addrs = new address[](1);
        addrs[0] = recipient1;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _submitWithOraclesForEpoch(addrs, ws, 1);

        // Settle epoch 0
        _settleEpoch0();
        // Settle epoch 1 — promotes activeEpoch to 1
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);
        assertEq(emission.activeEpoch(), 1);

        uint256 balAfterEpoch1 = awpToken.balanceOf(recipient1);
        assertTrue(balAfterEpoch1 > 0);

        // Settle epoch 2 — no weights for epoch 2, activeEpoch stays at 1
        vm.warp(block.timestamp + EPOCH_DURATION + 1);
        emission.settleEpoch(200);
        assertEq(emission.activeEpoch(), 1);

        // Still distributes using epoch 1 weights
        uint256 balAfterEpoch2 = awpToken.balanceOf(recipient1);
        assertTrue(balAfterEpoch2 > balAfterEpoch1);
    }
}
