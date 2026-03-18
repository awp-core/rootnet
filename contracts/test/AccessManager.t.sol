// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {AccessManager} from "../src/core/AccessManager.sol";

contract AccessManagerTest is Test {
    AccessManager public am;

    address public rootNet = makeAddr("rootNet");
    address public alice = makeAddr("alice");
    address public bob = makeAddr("bob");
    address public charlie = makeAddr("charlie");
    address public eve = makeAddr("eve");

    function setUp() public {
        am = new AccessManager(rootNet);
    }

    // ──────────────────────────────────────────────
    // register
    // ──────────────────────────────────────────────

    function test_register_success() public {
        vm.prank(rootNet);
        am.register(alice);

        assertTrue(am.isRegistered(alice));
        assertEq(am.registeredAt(alice), uint64(block.timestamp));
        assertEq(am.totalUsers(), 1);
    }

    function test_register_alreadyRegistered_reverts() public {
        vm.startPrank(rootNet);
        am.register(alice);

        vm.expectRevert(AccessManager.AlreadyRegistered.selector);
        am.register(alice);
        vm.stopPrank();
    }

    function test_register_addressIsAgent_reverts() public {
        // First bind charlie as bobs agent (bind auto-registers bob)
        vm.startPrank(rootNet);
        am.bind(charlie, bob);

        // charlie is already an agent, cannot register as user
        vm.expectRevert(AccessManager.AddressIsAgent.selector);
        am.register(charlie);
        vm.stopPrank();
    }

    function test_register_onlyRootNet() public {
        vm.prank(alice);
        vm.expectRevert(AccessManager.NotRootNet.selector);
        am.register(alice);
    }

    function test_register_incrementsTotalUsers() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.register(bob);
        vm.stopPrank();

        assertEq(am.totalUsers(), 2);
    }

    // ──────────────────────────────────────────────
    // bind
    // ──────────────────────────────────────────────

    function test_bind_success_newAgent() public {
        vm.startPrank(rootNet);
        am.register(alice);
        address oldPrincipal = am.bind(bob, alice);
        vm.stopPrank();

        assertEq(oldPrincipal, address(0));
        assertEq(am.agentOwner(bob), alice);
        address[] memory agents = am.getAgents(alice);
        assertEq(agents.length, 1);
        assertEq(agents[0], bob);
    }

    function test_bind_autoRegistersPrincipal() public {
        // alice not yet registered; bind should auto-register her
        assertFalse(am.isRegistered(alice));

        vm.prank(rootNet);
        am.bind(bob, alice);

        assertTrue(am.isRegistered(alice));
        assertEq(am.agentOwner(bob), alice);
        assertEq(am.totalUsers(), 1);
    }

    function test_bind_rebind_success() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.register(charlie);
        am.bind(bob, alice);

        // rebind bob from alice to charlie
        address oldPrincipal = am.bind(bob, charlie);
        vm.stopPrank();

        assertEq(oldPrincipal, alice);
        assertEq(am.agentOwner(bob), charlie);
        assertEq(am.getAgents(alice).length, 0);
        assertEq(am.getAgents(charlie).length, 1);
    }

    function test_bind_agentIsRegisteredPrincipal_reverts() public {
        vm.startPrank(rootNet);
        am.register(alice);

        // alice is already a principal, cannot become an agent
        vm.expectRevert(AccessManager.AddressIsPrincipal.selector);
        am.bind(alice, bob);
        vm.stopPrank();
    }

    function test_bind_principalIsAgent_reverts() public {
        vm.startPrank(rootNet);
        // bob is agent of alice
        am.bind(bob, alice);

        // can't make charlie an agent of bob (bob is an agent, not a principal)
        vm.expectRevert(AccessManager.AddressIsAgent.selector);
        am.bind(charlie, bob);
        vm.stopPrank();
    }

    function test_bind_agentIsSelf_reverts() public {
        vm.prank(rootNet);
        vm.expectRevert(AccessManager.AgentIsSelf.selector);
        am.bind(alice, alice);
    }

    function test_bind_selfRebind_noOp() public {
        // Rebinding to the same principal should be a no-op (no state change)
        vm.startPrank(rootNet);
        am.register(alice);
        am.bind(bob, alice);

        address[] memory agentsBefore = am.getAgents(alice);
        address oldPrincipal = am.bind(bob, alice); // same principal
        vm.stopPrank();

        // Should return alice as oldPrincipal (no change)
        assertEq(oldPrincipal, alice);
        // State unchanged
        assertEq(am.agentOwner(bob), alice);
        assertEq(am.getAgents(alice).length, agentsBefore.length);
    }

    function test_bind_manyAgents_noLimit() public {
        // No per-principal agent limit; 100+ agents should be fine
        vm.startPrank(rootNet);
        am.register(alice);
        for (uint256 i = 0; i < 100; i++) {
            address agent = address(uint160(0xA000 + i));
            am.bind(agent, alice);
        }
        vm.stopPrank();

        assertEq(am.getAgents(alice).length, 100);
    }

    function test_bind_onlyRootNet() public {
        vm.prank(alice);
        vm.expectRevert(AccessManager.NotRootNet.selector);
        am.bind(bob, alice);
    }

    // ──────────────────────────────────────────────
    // removeAgent
    // ──────────────────────────────────────────────

    function test_removeAgent_success() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.bind(bob, alice);
        am.setManager(alice, bob, true, alice);

        // alice (operator) removes bob
        am.removeAgent(alice, bob, alice);
        vm.stopPrank();

        assertEq(am.agentOwner(bob), address(0));
        assertFalse(am.isManager(bob));
        assertEq(am.getAgents(alice).length, 0);
    }

    function test_removeAgent_wrongOwner_reverts() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.register(charlie);
        am.bind(bob, alice);

        // charlie is not bobs owner
        vm.expectRevert(AccessManager.NotAgentOwner.selector);
        am.removeAgent(charlie, bob, charlie);
        vm.stopPrank();
    }

    function test_removeAgent_selfRemoval_reverts() public {
        vm.startPrank(rootNet);
        am.bind(bob, alice);

        // bob cannot remove itself
        vm.expectRevert(AccessManager.CannotRemoveSelf.selector);
        am.removeAgent(alice, bob, bob);
        vm.stopPrank();
    }

    function test_removeAgent_onlyRootNet() public {
        vm.prank(rootNet);
        am.bind(bob, alice);

        vm.prank(alice);
        vm.expectRevert(AccessManager.NotRootNet.selector);
        am.removeAgent(alice, bob, alice);
    }

    // ──────────────────────────────────────────────
    // setManager
    // ──────────────────────────────────────────────

    function test_setManager_grant_success() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.bind(bob, alice);

        am.setManager(alice, bob, true, alice);
        vm.stopPrank();

        assertTrue(am.isManager(bob));
    }

    function test_setManager_revoke_success() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.bind(bob, alice);
        am.setManager(alice, bob, true, alice);

        // alice revokes bobs manager
        am.setManager(alice, bob, false, alice);
        vm.stopPrank();

        assertFalse(am.isManager(bob));
    }

    function test_setManager_revokeSelf_reverts() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.bind(bob, alice);
        am.setManager(alice, bob, true, alice);

        // bob cannot revoke its own manager privileges
        vm.expectRevert(AccessManager.CannotRevokeSelf.selector);
        am.setManager(alice, bob, false, bob);
        vm.stopPrank();
    }

    function test_setManager_wrongOwner_reverts() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.register(charlie);
        am.bind(bob, alice);

        vm.expectRevert(AccessManager.NotAgentOwner.selector);
        am.setManager(charlie, bob, true, charlie);
        vm.stopPrank();
    }

    function test_setManager_selfGrant_allowed() public {
        // agent granting itself manager is allowed (only revoking self is blocked)
        vm.startPrank(rootNet);
        am.register(alice);
        am.bind(bob, alice);

        am.setManager(alice, bob, true, bob);
        vm.stopPrank();

        assertTrue(am.isManager(bob));
    }

    function test_setManager_onlyRootNet() public {
        vm.prank(rootNet);
        am.bind(bob, alice);

        vm.prank(alice);
        vm.expectRevert(AccessManager.NotRootNet.selector);
        am.setManager(alice, bob, true, alice);
    }

    // ──────────────────────────────────────────────
    // setRewardRecipient
    // ──────────────────────────────────────────────

    function test_setRewardRecipient_success() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.setRewardRecipient(alice, charlie);
        vm.stopPrank();

        assertEq(am.rewardRecipients(alice), charlie);
    }

    function test_setRewardRecipient_zeroAddress_reverts() public {
        vm.startPrank(rootNet);
        am.register(alice);

        vm.expectRevert(AccessManager.InvalidRecipient.selector);
        am.setRewardRecipient(alice, address(0));
        vm.stopPrank();
    }

    function test_setRewardRecipient_onlyRootNet() public {
        vm.prank(alice);
        vm.expectRevert(AccessManager.NotRootNet.selector);
        am.setRewardRecipient(alice, charlie);
    }

    // ──────────────────────────────────────────────
    // View functions
    // ──────────────────────────────────────────────

    function test_getOwner_returnsOwnerForAgent() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.bind(bob, alice);
        vm.stopPrank();

        assertEq(am.getOwner(bob), alice);
    }

    function test_getOwner_returnsSelfForUser() public {
        vm.prank(rootNet);
        am.register(alice);

        assertEq(am.getOwner(alice), alice);
    }

    function test_getOwner_returnsZeroForUnknown() public view {
        assertEq(am.getOwner(alice), address(0));
    }

    function test_isRegisteredUser() public {
        assertFalse(am.isRegisteredUser(alice));

        vm.prank(rootNet);
        am.register(alice);

        assertTrue(am.isRegisteredUser(alice));
    }

    function test_isRegisteredAgent() public {
        assertFalse(am.isRegisteredAgent(bob));

        vm.startPrank(rootNet);
        am.register(alice);
        am.bind(bob, alice);
        vm.stopPrank();

        assertTrue(am.isRegisteredAgent(bob));
    }

    function test_isKnownAddress() public {
        assertFalse(am.isKnownAddress(alice));
        assertFalse(am.isKnownAddress(bob));

        vm.startPrank(rootNet);
        am.register(alice);
        vm.stopPrank();

        assertTrue(am.isKnownAddress(alice));
        assertFalse(am.isKnownAddress(bob));

        vm.prank(rootNet);
        am.bind(bob, alice);

        assertTrue(am.isKnownAddress(bob));
    }

    function test_isAgent_agentBelongsToUser() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.bind(bob, alice);
        vm.stopPrank();

        assertTrue(am.isAgent(alice, bob));
    }

    function test_isAgent_userIsSelfAgent() public {
        vm.prank(rootNet);
        am.register(alice);

        // user itself is considered its own agent
        assertTrue(am.isAgent(alice, alice));
    }

    function test_isAgent_unrelatedReturnsFalse() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.register(charlie);
        am.bind(bob, alice);
        vm.stopPrank();

        assertFalse(am.isAgent(charlie, bob));
    }

    function test_isAgent_unregisteredUserReturnsFalse() public view {
        assertFalse(am.isAgent(alice, alice));
    }

    function test_isManagerAgent() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.bind(bob, alice);
        vm.stopPrank();

        assertFalse(am.isManagerAgent(bob));

        vm.prank(rootNet);
        am.setManager(alice, bob, true, alice);

        assertTrue(am.isManagerAgent(bob));
    }

    function test_getRewardRecipient_defaultReturnsSelf() public {
        vm.prank(rootNet);
        am.register(alice);

        assertEq(am.getRewardRecipient(alice), alice);
    }

    function test_getRewardRecipient_customRecipient() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.setRewardRecipient(alice, charlie);
        vm.stopPrank();

        assertEq(am.getRewardRecipient(alice), charlie);
    }

    function test_getTotalUsers() public {
        assertEq(am.getTotalUsers(), 0);

        vm.startPrank(rootNet);
        am.register(alice);
        am.register(bob);
        vm.stopPrank();

        assertEq(am.getTotalUsers(), 2);
    }

    // ──────────────────────────────────────────────
    // unbind
    // ──────────────────────────────────────────────

    function test_unbind_success() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.bind(bob, alice);
        am.setManager(alice, bob, true, alice);

        address oldPrincipal = am.unbind(bob);
        vm.stopPrank();

        assertEq(oldPrincipal, alice);
        // bob is now unregistered
        assertEq(am.agentOwner(bob), address(0));
        assertFalse(am.isManager(bob));
        assertFalse(am.isRegisteredAgent(bob));
        assertFalse(am.isKnownAddress(bob));
        // alice's agent list is empty
        assertEq(am.getAgents(alice).length, 0);
    }

    function test_unbind_notBound_reverts() public {
        vm.prank(rootNet);
        vm.expectRevert(AccessManager.NotBound.selector);
        am.unbind(bob);
    }

    function test_unbind_onlyRootNet() public {
        vm.prank(rootNet);
        am.bind(bob, alice);

        vm.prank(bob);
        vm.expectRevert(AccessManager.NotRootNet.selector);
        am.unbind(bob);
    }

    function test_unbind_agentCanRebindAfterUnbind() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.register(charlie);
        am.bind(bob, alice);
        am.unbind(bob);

        // bob is now unregistered, can bind to charlie
        address oldPrincipal = am.bind(bob, charlie);
        vm.stopPrank();

        assertEq(oldPrincipal, address(0)); // fresh bind
        assertEq(am.agentOwner(bob), charlie);
    }

    function test_bind_autoRegister_incrementsTotalUsers() public {
        // bind auto-registers principal, should increment totalUsers
        assertEq(am.getTotalUsers(), 0);

        vm.prank(rootNet);
        am.bind(bob, alice);

        assertEq(am.getTotalUsers(), 1);
        assertTrue(am.isRegistered(alice));
    }

    function test_getAgents_empty() public {
        vm.prank(rootNet);
        am.register(alice);

        assertEq(am.getAgents(alice).length, 0);
    }

    function test_getAgents_multipleAgents() public {
        vm.startPrank(rootNet);
        am.register(alice);
        am.bind(bob, alice);
        am.bind(charlie, alice);
        vm.stopPrank();

        address[] memory agents = am.getAgents(alice);
        assertEq(agents.length, 2);
    }

    // ──────────────────────────────────────────────
    // rootNet immutable
    // ──────────────────────────────────────────────

    function test_rootNet_isSet() public view {
        assertEq(am.rootNet(), rootNet);
    }
}
