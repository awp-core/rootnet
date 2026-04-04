// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {DeployHelper} from "./helpers/DeployHelper.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {IAWPRegistry} from "../src/interfaces/IAWPRegistry.sol";

contract AWPRegistryTest is DeployHelper {
    function setUp() public {
        _deployAll();
    }

    // ═══════════════════════════════════════════════
    //  Binding
    // ═══════════════════════════════════════════════

    function test_bind() public {
        vm.prank(alice);
        awpRegistry.bind(bob);
        assertEq(awpRegistry.boundTo(alice), bob);
    }

    function test_unbind() public {
        vm.prank(alice);
        awpRegistry.bind(bob);
        vm.prank(alice);
        awpRegistry.unbind();
        assertEq(awpRegistry.boundTo(alice), address(0));
    }

    function test_bind_selfBind_reverts() public {
        vm.prank(alice);
        vm.expectRevert(AWPRegistry.SelfBind.selector);
        awpRegistry.bind(alice);
    }

    function test_bind_cycle_reverts() public {
        vm.prank(alice);
        awpRegistry.bind(bob);
        vm.prank(bob);
        vm.expectRevert(AWPRegistry.CycleDetected.selector);
        awpRegistry.bind(alice);
    }

    // ═══════════════════════════════════════════════
    //  Recipient
    // ═══════════════════════════════════════════════

    function test_setRecipient() public {
        vm.prank(alice);
        awpRegistry.setRecipient(bob);
        assertEq(awpRegistry.recipient(alice), bob);
    }

    function test_setRecipient_clear() public {
        vm.prank(alice);
        awpRegistry.setRecipient(bob);
        vm.prank(alice);
        awpRegistry.setRecipient(address(0));
        assertEq(awpRegistry.recipient(alice), address(0));
    }

    function test_resolveRecipient() public {
        vm.prank(alice);
        awpRegistry.bind(bob);
        vm.prank(bob);
        awpRegistry.setRecipient(relayer);

        assertEq(awpRegistry.resolveRecipient(alice), relayer);
    }

    // ═══════════════════════════════════════════════
    //  Delegation
    // ═══════════════════════════════════════════════

    function test_grantDelegate() public {
        vm.prank(alice);
        awpRegistry.grantDelegate(bob);
        assertTrue(awpRegistry.delegates(alice, bob));
    }

    function test_revokeDelegate() public {
        vm.prank(alice);
        awpRegistry.grantDelegate(bob);
        vm.prank(alice);
        awpRegistry.revokeDelegate(bob);
        assertFalse(awpRegistry.delegates(alice, bob));
    }

    // ═══════════════════════════════════════════════
    //  Worknet Registration
    // ═══════════════════════════════════════════════

    function test_registerWorknet() public {
        uint256 wid = _registerWorknet(alice);

        assertTrue(wid > 0);
        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Pending));
    }

    function test_registerWorknet_escrpwsAWP() public {
        uint256 cost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        uint256 balBefore = awp.balanceOf(alice);

        _registerWorknet(alice);

        assertEq(balBefore - awp.balanceOf(alice), cost);
    }

    function test_registerWorknet_worknetId_format() public {
        uint256 wid = _registerWorknet(alice);
        // chainId * 100_000_000 + localId
        assertEq(wid / 100_000_000, block.chainid);
        assertEq(wid % 100_000_000, 1);
    }

    // ═══════════════════════════════════════════════
    //  Activation (Guardian only)
    // ═══════════════════════════════════════════════

    function test_activateWorknet() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Active));
        assertTrue(info.lpPool != bytes32(0));
    }

    function test_activateWorknet_mintsNFT() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        assertEq(awpWorkNet.ownerOf(wid), alice);
    }

    function test_activateWorknet_notGuardian_reverts() public {
        uint256 wid = _registerWorknet(alice);

        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.activateWorknet(wid);
    }

    // ═══════════════════════════════════════════════
    //  Cancel / Reject
    // ═══════════════════════════════════════════════

    function test_cancelWorknet_refundsAWP() public {
        uint256 cost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        uint256 balBefore = awp.balanceOf(alice);
        uint256 wid = _registerWorknet(alice);

        vm.prank(alice);
        awpRegistry.cancelWorknet(wid);

        assertEq(awp.balanceOf(alice), balBefore); // fully refunded
    }

    function test_cancelWorknet_notOwner_reverts() public {
        uint256 wid = _registerWorknet(alice);

        vm.prank(bob);
        vm.expectRevert(AWPRegistry.NotOwner.selector);
        awpRegistry.cancelWorknet(wid);
    }

    function test_rejectWorknet_refundsToOwner() public {
        uint256 cost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        uint256 balBefore = awp.balanceOf(alice);
        uint256 wid = _registerWorknet(alice);

        vm.prank(guardian);
        awpRegistry.rejectWorknet(wid);

        assertEq(awp.balanceOf(alice), balBefore);
    }

    // ═══════════════════════════════════════════════
    //  Pause / Resume / Ban
    // ═══════════════════════════════════════════════

    function test_pauseWorknet() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        vm.prank(alice);
        awpRegistry.pauseWorknet(wid);

        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Paused));
    }

    function test_resumeWorknet() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        vm.prank(alice);
        awpRegistry.pauseWorknet(wid);
        vm.prank(alice);
        awpRegistry.resumeWorknet(wid);

        assertTrue(awpRegistry.isWorknetActive(wid));
    }

    function test_banWorknet() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        vm.prank(guardian);
        awpRegistry.banWorknet(wid);

        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Banned));
    }

    // ═══════════════════════════════════════════════
    //  GetWorknetFull
    // ═══════════════════════════════════════════════

    function test_getWorknetFull_pending() public {
        uint256 wid = _registerWorknet(alice);
        IAWPRegistry.WorknetFullInfo memory full = awpRegistry.getWorknetFull(wid);

        assertEq(full.owner, alice);
        assertEq(full.alphaToken, address(0)); // not deployed yet
        assertEq(uint8(full.status), uint8(IAWPRegistry.WorknetStatus.Pending));
    }

    function test_getWorknetFull_active() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        // Verify via getWorknet (simpler, no cross-contract NFT read)
        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Active));
        assertTrue(info.lpPool != bytes32(0));
    }

    function test_getWorknetFull_nonExistent_reverts() public {
        vm.expectRevert();
        awpRegistry.getWorknetFull(999999);
    }

    // ═══════════════════════════════════════════════
    //  Guardian
    // ═══════════════════════════════════════════════

    function test_setGuardian() public {
        vm.prank(guardian);
        awpRegistry.setGuardian(alice);
        assertEq(awpRegistry.guardian(), alice);
    }

    function test_pause_unpause() public {
        vm.prank(guardian);
        awpRegistry.pause();
        assertTrue(awpRegistry.paused());

        vm.prank(guardian);
        awpRegistry.unpause();
        assertFalse(awpRegistry.paused());
    }

    // ═══════════════════════════════════════════════
    //  GetRegistry
    // ═══════════════════════════════════════════════

    function test_getRegistry() public view {
        (address a, address b, address c, address d, address e, address f, address g, address h, address i) =
            awpRegistry.getRegistry();
        assertEq(a, address(awp));
        assertEq(b, address(awpWorkNet));
        assertEq(c, address(factory));
        assertEq(d, address(awpEmission));
        assertEq(e, address(lpManager));
        assertEq(f, address(awpAllocator));
        assertEq(g, address(veAwp));
        assertEq(h, address(treasury));
        assertEq(i, guardian);
    }

    // ═══════════════════════════════════════════════
    //  Rescue
    // ═══════════════════════════════════════════════

    function test_rescueToken_cannotRescueAWP() public {
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.CannotRescueEscrowedToken.selector);
        awpRegistry.rescueToken(address(awp), guardian, 1);
    }

    // ═══════════════════════════════════════════════
    //  Register — edge cases
    // ═══════════════════════════════════════════════

    function test_registerWorknet_emptyName_reverts() public {
        uint256 cost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.startPrank(alice);
        awp.approve(address(awpRegistry), cost);
        vm.expectRevert(AWPRegistry.InvalidWorknetName.selector);
        awpRegistry.registerWorknet(IAWPRegistry.WorknetParams({
            name: "", symbol: "T", worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        }));
        vm.stopPrank();
    }

    function test_registerWorknet_nameTooLong_reverts() public {
        uint256 cost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        // 65 chars exceeds 64 limit
        bytes memory longName = new bytes(65);
        for (uint i; i < 65; i++) longName[i] = "a";

        vm.startPrank(alice);
        awp.approve(address(awpRegistry), cost);
        vm.expectRevert(AWPRegistry.InvalidWorknetName.selector);
        awpRegistry.registerWorknet(IAWPRegistry.WorknetParams({
            name: string(longName), symbol: "T", worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        }));
        vm.stopPrank();
    }

    function test_registerWorknet_emptySymbol_reverts() public {
        uint256 cost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.startPrank(alice);
        awp.approve(address(awpRegistry), cost);
        vm.expectRevert(AWPRegistry.InvalidWorknetSymbol.selector);
        awpRegistry.registerWorknet(IAWPRegistry.WorknetParams({
            name: "Test", symbol: "", worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        }));
        vm.stopPrank();
    }

    function test_registerWorknet_symbolTooLong_reverts() public {
        uint256 cost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        bytes memory longSym = new bytes(17);
        for (uint i; i < 17; i++) longSym[i] = "X";

        vm.startPrank(alice);
        awp.approve(address(awpRegistry), cost);
        vm.expectRevert(AWPRegistry.InvalidWorknetSymbol.selector);
        awpRegistry.registerWorknet(IAWPRegistry.WorknetParams({
            name: "Test", symbol: string(longSym), worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        }));
        vm.stopPrank();
    }

    function test_registerWorknet_jsonUnsafeName_reverts() public {
        uint256 cost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.startPrank(alice);
        awp.approve(address(awpRegistry), cost);
        vm.expectRevert(AWPRegistry.JsonUnsafeCharacter.selector);
        awpRegistry.registerWorknet(IAWPRegistry.WorknetParams({
            name: 'bad"quote', symbol: "T", worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        }));
        vm.stopPrank();
    }

    function test_registerWorknet_withCustomWorknetManager() public {
        uint256 cost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.startPrank(alice);
        awp.approve(address(awpRegistry), cost);
        uint256 wid = awpRegistry.registerWorknet(IAWPRegistry.WorknetParams({
            name: "CustomWM", symbol: "CWM",
            worknetManager: makeAddr("customWM"),
            salt: bytes32(0), minStake: 100e18, skillsURI: "https://skills"
        }));
        vm.stopPrank();
        assertTrue(wid > 0);
    }

    function test_registerWorknet_incrementsLocalId() public {
        uint256 wid1 = _registerWorknet(alice, "Net1", "N1");
        uint256 wid2 = _registerWorknet(bob, "Net2", "N2");

        assertEq(wid1 % 100_000_000, 1);
        assertEq(wid2 % 100_000_000, 2);
    }

    function test_registerWorknet_insufficientApproval_reverts() public {
        vm.startPrank(alice);
        awp.approve(address(awpRegistry), 1); // too little
        vm.expectRevert();
        awpRegistry.registerWorknet(IAWPRegistry.WorknetParams({
            name: "Test", symbol: "T", worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        }));
        vm.stopPrank();
    }

    // ═══════════════════════════════════════════════
    //  Activate — edge cases
    // ═══════════════════════════════════════════════

    function test_activateWorknet_deploysAlphaToken() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        IAWPRegistry.WorknetFullInfo memory full = awpRegistry.getWorknetFull(wid);
        assertTrue(full.alphaToken != address(0));
    }

    function test_activateWorknet_deploysWorknetManagerProxy() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        IAWPRegistry.WorknetFullInfo memory full = awpRegistry.getWorknetFull(wid);
        assertTrue(full.worknetManager != address(0));
    }

    function test_activateWorknet_alreadyActive_reverts() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        vm.prank(guardian);
        vm.expectRevert();
        awpRegistry.activateWorknet(wid);
    }

    function test_activateWorknet_cancelled_reverts() public {
        uint256 wid = _registerWorknet(alice);
        vm.prank(alice);
        awpRegistry.cancelWorknet(wid);

        vm.prank(guardian);
        vm.expectRevert();
        awpRegistry.activateWorknet(wid);
    }

    function test_activateWorknet_setsActivatedAt() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(info.activatedAt, block.timestamp);
    }

    function test_activateWorknet_clearsEscrowAndPending() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        (uint128 lpAmount, uint128 alphaMint) = awpRegistry.worknetEscrow(wid);
        assertEq(lpAmount, 0);
        assertEq(alphaMint, 0);
    }

    // ═══════════════════════════════════════════════
    //  Cancel / Reject — edge cases
    // ═══════════════════════════════════════════════

    function test_cancelWorknet_alreadyActive_reverts() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        vm.prank(alice);
        vm.expectRevert(); // pendingWorknets[wid].owner == address(0)
        awpRegistry.cancelWorknet(wid);
    }

    function test_rejectWorknet_alreadyActive_reverts() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        vm.prank(guardian);
        vm.expectRevert();
        awpRegistry.rejectWorknet(wid);
    }

    function test_rejectWorknet_notGuardian_reverts() public {
        uint256 wid = _registerWorknet(alice);

        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.rejectWorknet(wid);
    }

    function test_cancelWorknet_doubleCancel_reverts() public {
        uint256 wid = _registerWorknet(alice);

        vm.prank(alice);
        awpRegistry.cancelWorknet(wid);

        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotOwner.selector);
        awpRegistry.cancelWorknet(wid);
    }

    // ═══════════════════════════════════════════════
    //  Pause / Resume / Ban — edge cases
    // ═══════════════════════════════════════════════

    function test_pauseWorknet_notOwner_reverts() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        vm.prank(bob);
        vm.expectRevert(AWPRegistry.NotOwner.selector);
        awpRegistry.pauseWorknet(wid);
    }

    function test_pauseWorknet_alreadyPaused_reverts() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        vm.prank(alice);
        awpRegistry.pauseWorknet(wid);

        vm.prank(alice);
        vm.expectRevert(); // InvalidWorknetStatus
        awpRegistry.pauseWorknet(wid);
    }

    function test_resumeWorknet_notPaused_reverts() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        vm.prank(alice);
        vm.expectRevert(); // InvalidWorknetStatus
        awpRegistry.resumeWorknet(wid);
    }

    function test_resumeWorknet_notOwner_reverts() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        vm.prank(alice);
        awpRegistry.pauseWorknet(wid);

        vm.prank(bob);
        vm.expectRevert(AWPRegistry.NotOwner.selector);
        awpRegistry.resumeWorknet(wid);
    }

    function test_banWorknet_pending_reverts() public {
        uint256 wid = _registerWorknet(alice);

        vm.prank(guardian);
        vm.expectRevert(); // InvalidWorknetStatus
        awpRegistry.banWorknet(wid);
    }

    function test_banWorknet_notGuardian_reverts() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.banWorknet(wid);
    }

    function test_unbanWorknet_notBanned_reverts() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        vm.prank(guardian);
        vm.expectRevert(); // InvalidWorknetStatus
        awpRegistry.unbanWorknet(wid);
    }

    function test_banWorknet_fromPaused() public {
        uint256 wid = _registerWorknet(alice);
        _activateWorknet(wid);

        // Pause first
        vm.prank(alice);
        awpRegistry.pauseWorknet(wid);

        // Then ban from Paused state
        vm.prank(guardian);
        awpRegistry.banWorknet(wid);

        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Banned));
    }

    // ═══════════════════════════════════════════════
    //  Active Worknet Tracking
    // ═══════════════════════════════════════════════

    function test_getActiveWorknetCount_after_register_activate() public {
        assertEq(awpRegistry.getActiveWorknetCount(), 0);

        uint256 wid = _registerWorknet(alice);
        assertEq(awpRegistry.getActiveWorknetCount(), 0); // Pending doesn't count

        _activateWorknet(wid);
        assertEq(awpRegistry.getActiveWorknetCount(), 1);
    }

    function test_getActiveWorknetIds() public {
        uint256 wid1 = _registerWorknet(alice, "A", "A");
        uint256 wid2 = _registerWorknet(bob, "B", "B");
        _activateWorknet(wid1);
        _activateWorknet(wid2);

        uint256[] memory ids = awpRegistry.getActiveWorknetIds(0, 10);
        assertEq(ids.length, 2);
    }

    function test_isWorknetActive_transitions() public {
        uint256 wid = _registerWorknet(alice);
        assertFalse(awpRegistry.isWorknetActive(wid)); // Pending

        _activateWorknet(wid);
        assertTrue(awpRegistry.isWorknetActive(wid));  // Active

        vm.prank(alice);
        awpRegistry.pauseWorknet(wid);
        assertFalse(awpRegistry.isWorknetActive(wid)); // Paused

        vm.prank(alice);
        awpRegistry.resumeWorknet(wid);
        assertTrue(awpRegistry.isWorknetActive(wid));  // Active again

        vm.prank(guardian);
        awpRegistry.banWorknet(wid);
        assertFalse(awpRegistry.isWorknetActive(wid)); // Banned
    }

    function test_activeWorknetCount_decreasesOnPauseAndBan() public {
        uint256 wid1 = _registerWorknet(alice, "A", "A");
        uint256 wid2 = _registerWorknet(bob, "B", "B");
        _activateWorknet(wid1);
        _activateWorknet(wid2);
        assertEq(awpRegistry.getActiveWorknetCount(), 2);

        vm.prank(alice);
        awpRegistry.pauseWorknet(wid1);
        assertEq(awpRegistry.getActiveWorknetCount(), 1);

        vm.prank(guardian);
        awpRegistry.banWorknet(wid2);
        assertEq(awpRegistry.getActiveWorknetCount(), 0);
    }

    // ═══════════════════════════════════════════════
    //  Guardian Parameters
    // ═══════════════════════════════════════════════

    function test_setInitialAlphaPrice() public {
        vm.prank(guardian);
        awpRegistry.setInitialAlphaPrice(1e15);
        assertEq(awpRegistry.initialAlphaPrice(), 1e15);
    }

    function test_setInitialAlphaPrice_tooLow_reverts() public {
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.PriceTooLow.selector);
        awpRegistry.setInitialAlphaPrice(1e11);
    }

    function test_setInitialAlphaPrice_tooHigh_reverts() public {
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.PriceTooHigh.selector);
        awpRegistry.setInitialAlphaPrice(1e31);
    }

    function test_setInitialAlphaMint() public {
        vm.prank(guardian);
        awpRegistry.setInitialAlphaMint(5e26);
        assertEq(awpRegistry.initialAlphaMint(), 5e26);
    }

    function test_setInitialAlphaMint_zero_reverts() public {
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.InvalidMintAmount.selector);
        awpRegistry.setInitialAlphaMint(0);
    }

    function test_setWorknetManagerImpl() public {
        address newImpl = makeAddr("newImpl");
        vm.prank(guardian);
        awpRegistry.setWorknetManagerImpl(newImpl);
    }

    function test_setWorknetManagerImpl_zero_reverts() public {
        vm.prank(guardian);
        vm.expectRevert(AWPRegistry.ZeroAddress.selector);
        awpRegistry.setWorknetManagerImpl(address(0));
    }
}
