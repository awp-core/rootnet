// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {AlphaToken} from "../src/token/AlphaToken.sol";
import {IERC1363Receiver, IERC1363Spender} from "../src/interfaces/IERC1363Receiver.sol";

/// @title AlphaTokenTest — Unit tests for worknet AlphaToken
contract AlphaTokenTest is Test {
    AlphaToken public alpha;
    address public admin = address(this);
    address public worknetMgr = makeAddr("worknetMgr");
    address public alice = makeAddr("alice");

    uint256 constant WORKNET_ID = 845300000001;

    function setUp() public {
        alpha = new AlphaToken();
        alpha.initialize("Alpha", "ALPHA", WORKNET_ID, admin);
    }

    // ═══════════════════════════════════════════════
    //  Initialization
    // ═══════════════════════════════════════════════

    function test_initialize() public view {
        assertEq(alpha.name(), "Alpha");
        assertEq(alpha.symbol(), "ALPHA");
        assertEq(alpha.worknetId(), WORKNET_ID);
        assertEq(alpha.admin(), admin);
        assertEq(alpha.createdAt(), uint64(block.timestamp));
        assertTrue(alpha.minters(admin));
        assertFalse(alpha.mintersLocked());
    }

    function test_initialize_twice_reverts() public {
        vm.expectRevert();
        alpha.initialize("X", "X", 1, admin);
    }

    function test_MAX_SUPPLY() public view {
        assertEq(alpha.MAX_SUPPLY(), 10_000_000_000 * 1e18);
    }

    // ═══════════════════════════════════════════════
    //  Admin minting (before lock)
    // ═══════════════════════════════════════════════

    function test_mint_asAdmin() public {
        alpha.mint(alice, 1000e18);
        assertEq(alpha.balanceOf(alice), 1000e18);
        assertEq(alpha.totalSupply(), 1000e18);
    }

    function test_mint_notMinter_reverts() public {
        vm.prank(alice);
        vm.expectRevert(AlphaToken.NotMinter.selector);
        alpha.mint(alice, 100e18);
    }

    function test_mint_exceedsMaxSupply_reverts() public {
        alpha.mint(alice, alpha.MAX_SUPPLY());
        vm.expectRevert(AlphaToken.ExceedsMaxSupply.selector);
        alpha.mint(alice, 1);
    }

    // ═══════════════════════════════════════════════
    //  setWorknetMinter
    // ═══════════════════════════════════════════════

    function test_setWorknetMinter() public {
        alpha.mint(alice, 1000e18); // pre-mint for LP

        alpha.setWorknetMinter(worknetMgr);

        assertTrue(alpha.mintersLocked());
        assertTrue(alpha.minters(worknetMgr));
        assertFalse(alpha.minters(admin));
        assertEq(alpha.supplyAtLock(), 1000e18);
        assertEq(alpha.grossMintedSinceLock(), 0);
    }

    function test_setWorknetMinter_notAdmin_reverts() public {
        vm.prank(alice);
        vm.expectRevert(AlphaToken.NotAdmin.selector);
        alpha.setWorknetMinter(worknetMgr);
    }

    function test_setWorknetMinter_zeroAddress_reverts() public {
        vm.expectRevert(AlphaToken.ZeroAddress.selector);
        alpha.setWorknetMinter(address(0));
    }

    function test_setWorknetMinter_twice_reverts() public {
        alpha.setWorknetMinter(worknetMgr);

        vm.expectRevert(AlphaToken.MintersLocked.selector);
        alpha.setWorknetMinter(address(1));
    }

    // ═══════════════════════════════════════════════
    //  Time-based minting cap (after lock)
    // ═══════════════════════════════════════════════

    function test_timeBasedCap_sameBlock() public {
        alpha.setWorknetMinter(worknetMgr);

        // Same block: elapsed=0 → forced to 1 → maxMintable = MAX_SUPPLY * 1 / 365 days
        uint256 maxMintable = alpha.MAX_SUPPLY() * 1 / 365 days;
        vm.prank(worknetMgr);
        alpha.mint(alice, maxMintable);
        assertEq(alpha.balanceOf(alice), maxMintable);
    }

    function test_timeBasedCap_afterOneYear() public {
        alpha.setWorknetMinter(worknetMgr);

        vm.warp(block.timestamp + 365 days);
        // After 1 year: maxMintable = MAX_SUPPLY (no pre-minted supply at lock)
        vm.startPrank(worknetMgr);
        // Mint in batches to stay within per-call limits
        alpha.mint(alice, alpha.MAX_SUPPLY());
        vm.stopPrank();
        assertEq(alpha.balanceOf(alice), alpha.MAX_SUPPLY());
    }

    function test_timeBasedCap_exceeds_reverts() public {
        alpha.setWorknetMinter(worknetMgr);

        vm.warp(block.timestamp + 30 days);
        uint256 maxMintable = alpha.MAX_SUPPLY() * 30 days / 365 days;

        vm.prank(worknetMgr);
        alpha.mint(alice, maxMintable);

        // One more should fail
        vm.prank(worknetMgr);
        vm.expectRevert(AlphaToken.ExceedsMintableLimit.selector);
        alpha.mint(alice, 1e18);
    }

    function test_timeBasedCap_withPreMint() public {
        // Pre-mint 1B tokens
        uint256 preMint = 1_000_000_000e18;
        alpha.mint(alice, preMint);
        alpha.setWorknetMinter(worknetMgr);

        assertEq(alpha.supplyAtLock(), preMint);

        vm.warp(block.timestamp + 365 days);
        // Max mintable after lock = (MAX_SUPPLY - supplyAtLock) * elapsed / 365 days
        uint256 remaining = alpha.MAX_SUPPLY() - preMint;
        vm.prank(worknetMgr);
        alpha.mint(alice, remaining);
        assertEq(alpha.totalSupply(), alpha.MAX_SUPPLY());
    }

    function test_currentMintableLimit_beforeLock() public view {
        assertEq(alpha.currentMintableLimit(), 0);
    }

    function test_currentMintableLimit_afterLock() public {
        alpha.setWorknetMinter(worknetMgr);

        vm.warp(block.timestamp + 30 days);
        uint256 limit = alpha.currentMintableLimit();
        uint256 expected = alpha.MAX_SUPPLY() * 30 days / 365 days;
        assertEq(limit, expected);
    }

    // ═══════════════════════════════════════════════
    //  setMinterPaused
    // ═══════════════════════════════════════════════

    function test_setMinterPaused() public {
        alpha.setWorknetMinter(worknetMgr);

        alpha.setMinterPaused(worknetMgr, true);

        vm.prank(worknetMgr);
        vm.expectRevert(AlphaToken.MinterPaused.selector);
        alpha.mint(alice, 100e18);
    }

    function test_setMinterPaused_resume() public {
        alpha.setWorknetMinter(worknetMgr);

        alpha.setMinterPaused(worknetMgr, true);
        alpha.setMinterPaused(worknetMgr, false);

        vm.prank(worknetMgr);
        alpha.mint(alice, 100e18);
        assertEq(alpha.balanceOf(alice), 100e18);
    }

    function test_setMinterPaused_notAdmin_reverts() public {
        vm.prank(alice);
        vm.expectRevert(AlphaToken.NotAdmin.selector);
        alpha.setMinterPaused(worknetMgr, true);
    }

    // ═══════════════════════════════════════════════
    //  ERC1363 callbacks
    // ═══════════════════════════════════════════════

    function test_transferAndCall_toEOA() public {
        alpha.mint(admin, 100e18);
        alpha.transferAndCall(alice, 50e18, "");
        assertEq(alpha.balanceOf(alice), 50e18);
    }

    function test_transferAndCall_toContract() public {
        MockReceiver receiver = new MockReceiver();
        alpha.mint(admin, 100e18);
        alpha.transferAndCall(address(receiver), 50e18, "hello");
        assertEq(alpha.balanceOf(address(receiver)), 50e18);
        assertTrue(receiver.called());
    }

    function test_transferAndCall_badCallback_reverts() public {
        BadReceiver bad = new BadReceiver();
        alpha.mint(admin, 100e18);
        vm.expectRevert(AlphaToken.InvalidCallback.selector);
        alpha.transferAndCall(address(bad), 50e18, "");
    }

    function test_approveAndCall_toContract() public {
        MockSpender spender = new MockSpender();
        alpha.approveAndCall(address(spender), 100e18, "data");
        assertEq(alpha.allowance(admin, address(spender)), 100e18);
        assertTrue(spender.called());
    }

    function test_approveAndCall_badCallback_reverts() public {
        BadSpender bad = new BadSpender();
        vm.expectRevert(AlphaToken.InvalidCallback.selector);
        alpha.approveAndCall(address(bad), 100e18, "");
    }

    // ═══════════════════════════════════════════════
    //  Burn
    // ═══════════════════════════════════════════════

    function test_burn() public {
        alpha.mint(admin, 100e18);
        alpha.burn(50e18);
        assertEq(alpha.totalSupply(), 50e18);
    }
}

// ── Test helpers ──

contract MockReceiver is IERC1363Receiver {
    bool public called;
    function onTransferReceived(address, address, uint256, bytes calldata) external returns (bytes4) {
        called = true;
        return IERC1363Receiver.onTransferReceived.selector;
    }
}

contract BadReceiver is IERC1363Receiver {
    function onTransferReceived(address, address, uint256, bytes calldata) external pure returns (bytes4) {
        return bytes4(0xdeadbeef);
    }
}

contract MockSpender is IERC1363Spender {
    bool public called;
    function onApprovalReceived(address, uint256, bytes calldata) external returns (bytes4) {
        called = true;
        return IERC1363Spender.onApprovalReceived.selector;
    }
}

contract BadSpender is IERC1363Spender {
    function onApprovalReceived(address, uint256, bytes calldata) external pure returns (bytes4) {
        return bytes4(0xdeadbeef);
    }
}
