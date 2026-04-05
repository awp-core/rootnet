// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {AWPToken} from "../src/token/AWPToken.sol";

contract TreasuryTest is Test {
    Treasury public treasury;
    AWPToken public awp;
    address public admin = address(this);
    address public proposer = makeAddr("proposer");
    address public alice = makeAddr("alice");

    uint256 constant MIN_DELAY = 2 days;

    function setUp() public {
        awp = new AWPToken("AWP", "AWP", admin);
        awp.initialMint(1_000_000e18);

        address[] memory proposers = new address[](1);
        proposers[0] = proposer;
        address[] memory executors = new address[](1);
        executors[0] = address(0); // anyone can execute
        treasury = new Treasury(MIN_DELAY, proposers, executors, admin);

        // Transfer AWP to treasury
        awp.transfer(address(treasury), 500_000e18);
    }

    // ═══════════════════════════════════════════════
    //  Basic Properties
    // ═══════════════════════════════════════════════

    function test_minDelay() public view {
        assertEq(treasury.getMinDelay(), MIN_DELAY);
    }

    function test_hasProposerRole() public view {
        assertTrue(treasury.hasRole(treasury.PROPOSER_ROLE(), proposer));
    }

    function test_hasAdminRole() public view {
        assertTrue(treasury.hasRole(treasury.DEFAULT_ADMIN_ROLE(), admin));
    }

    function test_anyoneCanExecute() public view {
        assertTrue(treasury.hasRole(treasury.EXECUTOR_ROLE(), address(0)));
    }

    // ═══════════════════════════════════════════════
    //  Schedule + Execute
    // ═══════════════════════════════════════════════

    function test_scheduleAndExecute() public {
        bytes memory data = abi.encodeCall(awp.transfer, (alice, 100e18));
        bytes32 salt = bytes32(uint256(1));

        // Schedule
        vm.prank(proposer);
        treasury.schedule(address(awp), 0, data, bytes32(0), salt, MIN_DELAY);

        // Warp past delay
        vm.warp(block.timestamp + MIN_DELAY + 1);

        // Execute
        treasury.execute(address(awp), 0, data, bytes32(0), salt);

        assertEq(awp.balanceOf(alice), 100e18);
    }

    function test_execute_beforeDelay_reverts() public {
        bytes memory data = abi.encodeCall(awp.transfer, (alice, 100e18));
        bytes32 salt = bytes32(uint256(2));

        vm.prank(proposer);
        treasury.schedule(address(awp), 0, data, bytes32(0), salt, MIN_DELAY);

        vm.expectRevert();
        treasury.execute(address(awp), 0, data, bytes32(0), salt);
    }

    function test_schedule_notProposer_reverts() public {
        bytes memory data = abi.encodeCall(awp.transfer, (alice, 100e18));
        bytes32 salt = bytes32(uint256(3));

        vm.prank(alice);
        vm.expectRevert();
        treasury.schedule(address(awp), 0, data, bytes32(0), salt, MIN_DELAY);
    }

    // ═══════════════════════════════════════════════
    //  Cancel
    // ═══════════════════════════════════════════════

    function test_cancel() public {
        bytes memory data = abi.encodeCall(awp.transfer, (alice, 100e18));
        bytes32 salt = bytes32(uint256(4));

        vm.prank(proposer);
        treasury.schedule(address(awp), 0, data, bytes32(0), salt, MIN_DELAY);

        bytes32 opId = treasury.hashOperation(address(awp), 0, data, bytes32(0), salt);
        assertTrue(treasury.isOperationPending(opId));

        // Cancel (proposer has CANCELLER_ROLE by default in OZ TimelockController)
        vm.prank(proposer);
        treasury.cancel(opId);

        assertFalse(treasury.isOperationPending(opId));
    }

    // ═══════════════════════════════════════════════
    //  Admin
    // ═══════════════════════════════════════════════

    function test_grantProposerRole() public {
        treasury.grantRole(treasury.PROPOSER_ROLE(), alice);
        assertTrue(treasury.hasRole(treasury.PROPOSER_ROLE(), alice));
    }

    function test_updateDelay() public {
        // Self-call via schedule+execute
        bytes memory data = abi.encodeCall(treasury.updateDelay, (1 days));
        bytes32 salt = bytes32(uint256(5));

        vm.prank(proposer);
        treasury.schedule(address(treasury), 0, data, bytes32(0), salt, MIN_DELAY);

        vm.warp(block.timestamp + MIN_DELAY + 1);
        treasury.execute(address(treasury), 0, data, bytes32(0), salt);

        assertEq(treasury.getMinDelay(), 1 days);
    }

    // ═══════════════════════════════════════════════
    //  Receive ETH
    // ═══════════════════════════════════════════════

    function test_receiveETH() public {
        vm.deal(alice, 1 ether);
        vm.prank(alice);
        (bool ok,) = address(treasury).call{value: 1 ether}("");
        assertTrue(ok);
        assertEq(address(treasury).balance, 1 ether);
    }
}
