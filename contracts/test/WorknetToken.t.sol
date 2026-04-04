// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {WorknetToken} from "../src/token/WorknetToken.sol";
import {MockWorknetTokenFactory} from "./helpers/MockWorknetTokenFactory.sol";
import {ERC1363Utils} from "@openzeppelin/contracts/token/ERC20/utils/ERC1363Utils.sol";
import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";
import {IERC1363Receiver} from "@openzeppelin/contracts/interfaces/IERC1363Receiver.sol";
import {IERC1363Spender} from "@openzeppelin/contracts/interfaces/IERC1363Spender.sol";

contract WorknetTokenTest is Test {
    WorknetToken public alpha;
    MockWorknetTokenFactory public tokenFactory;
    address public worknetMgr = makeAddr("worknetMgr");
    address public alice = makeAddr("alice");

    uint256 constant WORKNET_ID = 845300000001;

    function setUp() public {
        tokenFactory = new MockWorknetTokenFactory();
        alpha = tokenFactory.deploy("Alpha", "ALPHA", WORKNET_ID);
    }

    // ═══════════════════════════════════════════════
    //  Construction
    // ═══════════════════════════════════════════════

    function test_constructor() public view {
        assertEq(alpha.name(), "Alpha");
        assertEq(alpha.symbol(), "ALPHA");
        assertEq(alpha.worknetId(), WORKNET_ID);
        assertEq(alpha.minter(), address(0));
        assertEq(alpha.createdAt(), uint64(block.timestamp));
        assertFalse(alpha.initialized());
    }

    function test_MAX_SUPPLY() public view {
        assertEq(alpha.MAX_SUPPLY(), 10_000_000_000 * 1e18);
    }

    // ═══════════════════════════════════════════════
    //  Open minting (before setMinter — atomic with deploy)
    // ═══════════════════════════════════════════════

    function test_mint_beforeInit_openAccess() public {
        // Before init, anyone can mint (deploy + init in same tx, no front-run window)
        alpha.mint(alice, 1000e18);
        assertEq(alpha.balanceOf(alice), 1000e18);

        vm.prank(alice);
        alpha.mint(alice, 500e18);
        assertEq(alpha.balanceOf(alice), 1500e18);
    }

    function test_mint_exceedsMaxSupply_reverts() public {
        alpha.mint(alice, alpha.MAX_SUPPLY());
        vm.expectRevert(WorknetToken.ExceedsMaxSupply.selector);
        alpha.mint(alice, 1);
    }

    // ═══════════════════════════════════════════════
    //  setMinter
    // ═══════════════════════════════════════════════

    function test_setMinter() public {
        alpha.mint(alice, 1000e18);
        alpha.setMinter(worknetMgr);

        assertTrue(alpha.initialized());
        assertEq(alpha.minter(), worknetMgr);
        assertEq(alpha.supplyAtLock(), 1000e18);
    }

    function test_setMinter_zeroAddress_reverts() public {
        vm.expectRevert(WorknetToken.ZeroAddress.selector);
        alpha.setMinter(address(0));
    }

    function test_setMinter_twice_reverts() public {
        alpha.setMinter(worknetMgr);
        vm.expectRevert(WorknetToken.AlreadyInitialized.selector);
        alpha.setMinter(address(1));
    }

    function test_afterInit_onlyMinterCanMint() public {
        alpha.setMinter(worknetMgr);
        vm.expectRevert(WorknetToken.NotMinter.selector);
        alpha.mint(alice, 100e18);

        vm.prank(worknetMgr);
        alpha.mint(alice, 100e18);
        assertEq(alpha.balanceOf(alice), 100e18);
    }

    // ═══════════════════════════════════════════════
    //  Time-based minting cap
    // ═══════════════════════════════════════════════

    function test_timeBasedCap_sameBlock() public {
        alpha.setMinter(worknetMgr);

        uint256 maxMintable = alpha.MAX_SUPPLY() * 1 / 365 days;
        vm.prank(worknetMgr);
        alpha.mint(alice, maxMintable);
        assertEq(alpha.balanceOf(alice), maxMintable);
    }

    function test_timeBasedCap_afterOneYear() public {
        alpha.setMinter(worknetMgr);
        vm.warp(block.timestamp + 365 days);

        uint256 maxSupply = alpha.MAX_SUPPLY();
        vm.prank(worknetMgr);
        alpha.mint(alice, maxSupply);
        assertEq(alpha.balanceOf(alice), maxSupply);
    }

    function test_timeBasedCap_exceeds_reverts() public {
        alpha.setMinter(worknetMgr);
        vm.warp(block.timestamp + 30 days);

        uint256 maxMintable = alpha.MAX_SUPPLY() * 30 days / 365 days;
        vm.prank(worknetMgr);
        alpha.mint(alice, maxMintable);

        vm.prank(worknetMgr);
        vm.expectRevert(WorknetToken.ExceedsMintableLimit.selector);
        alpha.mint(alice, 1e18);
    }

    function test_timeBasedCap_withPreMint() public {
        uint256 preMint = 1_000_000_000e18;
        alpha.mint(alice, preMint);
        alpha.setMinter(worknetMgr);

        assertEq(alpha.supplyAtLock(), preMint);

        vm.warp(block.timestamp + 365 days);
        uint256 remaining = alpha.MAX_SUPPLY() - preMint;
        vm.prank(worknetMgr);
        alpha.mint(alice, remaining);
        assertEq(alpha.totalSupply(), alpha.MAX_SUPPLY());
    }

    function test_timeBasedCap_burnFreesHeadroom() public {
        alpha.setMinter(worknetMgr);
        vm.warp(block.timestamp + 30 days);

        uint256 cap = alpha.MAX_SUPPLY() * 30 days / 365 days;
        vm.prank(worknetMgr);
        alpha.mint(alice, cap);

        // Burn half
        vm.prank(alice);
        alpha.burn(cap / 2);

        // Can mint again (burned tokens freed headroom)
        vm.prank(worknetMgr);
        alpha.mint(alice, cap / 2);
    }

    function test_currentMintableLimit_beforeInit() public view {
        assertEq(alpha.currentMintableLimit(), 0);
    }

    function test_currentMintableLimit_afterInit() public {
        alpha.setMinter(worknetMgr);
        vm.warp(block.timestamp + 30 days);
        uint256 expected = alpha.MAX_SUPPLY() * 30 days / 365 days;
        assertEq(alpha.currentMintableLimit(), expected);
    }

    // ═══════════════════════════════════════════════
    //  ERC1363 (OZ)
    // ═══════════════════════════════════════════════

    function test_transferAndCall_toEOA_reverts() public {
        alpha.mint(address(this), 100e18);
        vm.expectRevert(abi.encodeWithSelector(ERC1363Utils.ERC1363InvalidReceiver.selector, alice));
        alpha.transferAndCall(alice, 50e18, "");
    }

    function test_transferAndCall_toContract() public {
        MockReceiver receiver = new MockReceiver();
        alpha.mint(address(this), 100e18);
        alpha.transferAndCall(address(receiver), 50e18, "hello");
        assertEq(alpha.balanceOf(address(receiver)), 50e18);
        assertTrue(receiver.called());
    }

    function test_transferAndCall_badCallback_reverts() public {
        BadReceiver bad = new BadReceiver();
        alpha.mint(address(this), 100e18);
        vm.expectRevert(abi.encodeWithSelector(ERC1363Utils.ERC1363InvalidReceiver.selector, address(bad)));
        alpha.transferAndCall(address(bad), 50e18, "");
    }

    function test_transferFromAndCall() public {
        MockReceiver receiver = new MockReceiver();
        alpha.mint(address(this), 100e18);
        alpha.approve(alice, 50e18);
        vm.prank(alice);
        alpha.transferFromAndCall(address(this), address(receiver), 50e18, "");
        assertEq(alpha.balanceOf(address(receiver)), 50e18);
    }

    function test_approveAndCall_toContract() public {
        MockSpender spender = new MockSpender();
        alpha.approveAndCall(address(spender), 100e18, "data");
        assertEq(alpha.allowance(address(this), address(spender)), 100e18);
        assertTrue(spender.called());
    }

    function test_approveAndCall_badCallback_reverts() public {
        BadSpender bad = new BadSpender();
        vm.expectRevert(abi.encodeWithSelector(ERC1363Utils.ERC1363InvalidSpender.selector, address(bad)));
        alpha.approveAndCall(address(bad), 100e18, "");
    }

    // ═══════════════════════════════════════════════
    //  ERC165 + Burn
    // ═══════════════════════════════════════════════

    function test_supportsInterface_ERC165() public view {
        assertTrue(alpha.supportsInterface(type(IERC165).interfaceId));
    }

    function test_burn() public {
        alpha.mint(address(this), 100e18);
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
