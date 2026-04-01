// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {AWPToken} from "../src/token/AWPToken.sol";

import {IERC1363Receiver, IERC1363Spender} from "../src/interfaces/IERC1363Receiver.sol";

// ── Helper contracts: ERC1363 callback receivers ──

contract MockReceiver is IERC1363Receiver {
    address public lastOperator;
    address public lastFrom;
    uint256 public lastAmount;
    bytes public lastData;

    function onTransferReceived(address operator, address from, uint256 amount, bytes calldata data)
        external
        returns (bytes4)
    {
        lastOperator = operator;
        lastFrom = from;
        lastAmount = amount;
        lastData = data;
        return IERC1363Receiver.onTransferReceived.selector;
    }
}

contract MockSpender is IERC1363Spender {
    address public lastOwner;
    uint256 public lastAmount;
    bytes public lastData;

    function onApprovalReceived(address owner, uint256 amount, bytes calldata data) external returns (bytes4) {
        lastOwner = owner;
        lastAmount = amount;
        lastData = data;
        return IERC1363Spender.onApprovalReceived.selector;
    }
}

contract RejectingReceiver is IERC1363Receiver {
    function onTransferReceived(address, address, uint256, bytes calldata) external pure returns (bytes4) {
        return bytes4(0xdeadbeef);
    }
}

contract RejectingSpender is IERC1363Spender {
    function onApprovalReceived(address, uint256, bytes calldata) external pure returns (bytes4) {
        return bytes4(0xdeadbeef);
    }
}

// ── AWPToken tests ──

contract AWPTokenTest is Test {
    AWPToken public token;
    address public deployer;
    address public alice;
    address public bob;
    address public minter;

    uint256 constant MAX_SUPPLY = 10_000_000_000 * 1e18;
    uint256 constant INITIAL_MINT = 200_000_000 * 1e18;

    function setUp() public {
        deployer = makeAddr("deployer");
        alice = makeAddr("alice");
        bob = makeAddr("bob");
        minter = makeAddr("minter");

        vm.startPrank(deployer);
        token = new AWPToken("AWP Token", "AWP", deployer);
        token.initialMint(200_000_000 * 1e18);
        vm.stopPrank();
    }

    // ── Constructor tests ──

    function test_constructor_mintsInitialSupply() public view {
        assertEq(token.balanceOf(deployer), INITIAL_MINT);
        assertEq(token.totalSupply(), INITIAL_MINT);
    }

    function test_constructor_setsAdmin() public view {
        assertEq(token.admin(), deployer);
    }

    function test_constructor_setsNameAndSymbol() public view {
        assertEq(token.name(), "AWP Token");
        assertEq(token.symbol(), "AWP");
    }

    // ── addMinter tests ──

    function test_addMinter_success() public {
        vm.prank(deployer);
        token.addMinter(minter);
        assertTrue(token.minters(minter));
    }

    function test_addMinter_revertIfNotAdmin() public {
        vm.prank(alice);
        vm.expectRevert(AWPToken.NotAdmin.selector);
        token.addMinter(minter);
    }

    // ── renounceAdmin tests ──

    function test_renounceAdmin_setsAdminToZero() public {
        vm.prank(deployer);
        token.renounceAdmin();
        assertEq(token.admin(), address(0));
    }

    function test_renounceAdmin_revertIfNotAdmin() public {
        vm.prank(alice);
        vm.expectRevert(AWPToken.NotAdmin.selector);
        token.renounceAdmin();
    }

    function test_renounceAdmin_preventsAddMinter() public {
        vm.prank(deployer);
        token.renounceAdmin();

        // After renouncing admin, nobody can add minters
        vm.prank(deployer);
        vm.expectRevert(AWPToken.NotAdmin.selector);
        token.addMinter(minter);
    }

    // ── mint tests ──

    function test_mint_success() public {
        vm.prank(deployer);
        token.addMinter(minter);

        uint256 amount = 1000 * 1e18;
        vm.prank(minter);
        token.mint(alice, amount);

        assertEq(token.balanceOf(alice), amount);
        assertEq(token.totalSupply(), INITIAL_MINT + amount);
    }

    function test_mint_revertIfNotMinter() public {
        vm.prank(alice);
        vm.expectRevert(AWPToken.NotMinter.selector);
        token.mint(alice, 100);
    }

    function test_mint_revertIfExceedsMaxSupply() public {
        vm.prank(deployer);
        token.addMinter(minter);

        // Remaining mintable amount = MAX_SUPPLY - INITIAL_MINT = 5B
        uint256 remaining = MAX_SUPPLY - INITIAL_MINT;

        vm.prank(minter);
        vm.expectRevert(AWPToken.ExceedsMaxSupply.selector);
        token.mint(alice, remaining + 1);
    }

    function test_mint_exactlyReachesMaxSupply() public {
        vm.prank(deployer);
        token.addMinter(minter);

        uint256 remaining = MAX_SUPPLY - INITIAL_MINT;

        vm.prank(minter);
        token.mint(alice, remaining);

        assertEq(token.totalSupply(), MAX_SUPPLY);
    }

    function test_mint_revertAfterMaxSupplyReached() public {
        vm.prank(deployer);
        token.addMinter(minter);

        uint256 remaining = MAX_SUPPLY - INITIAL_MINT;
        vm.prank(minter);
        token.mint(alice, remaining);

        // Minting 1 more should fail
        vm.prank(minter);
        vm.expectRevert(AWPToken.ExceedsMaxSupply.selector);
        token.mint(alice, 1);
    }

    function test_mint_deployerIsNotMinterByDefault() public {
        vm.prank(deployer);
        vm.expectRevert(AWPToken.NotMinter.selector);
        token.mint(alice, 100);
    }

    // ── burn / burnFrom tests ──

    function test_burn_reducesSupply() public {
        uint256 burnAmount = 1000 * 1e18;

        vm.prank(deployer);
        token.burn(burnAmount);

        assertEq(token.balanceOf(deployer), INITIAL_MINT - burnAmount);
        assertEq(token.totalSupply(), INITIAL_MINT - burnAmount);
    }

    function test_burnFrom_withApproval() public {
        uint256 transferAmount = 5000 * 1e18;
        uint256 burnAmount = 2000 * 1e18;

        vm.prank(deployer);
        token.transfer(alice, transferAmount);

        vm.prank(alice);
        token.approve(bob, burnAmount);

        vm.prank(bob);
        token.burnFrom(alice, burnAmount);

        assertEq(token.balanceOf(alice), transferAmount - burnAmount);
    }

    function test_burnFrom_revertWithoutApproval() public {
        vm.prank(deployer);
        token.transfer(alice, 1000 * 1e18);

        vm.prank(bob);
        vm.expectRevert();
        token.burnFrom(alice, 500 * 1e18);
    }

    // ── ERC20Permit nonces test ──

    function test_noncesWork() public view {
        assertEq(token.nonces(deployer), 0);
    }

    // ── ERC1363 transferAndCall tests ──

    function test_transferAndCall_toContract() public {
        MockReceiver receiver = new MockReceiver();
        uint256 amount = 100 * 1e18;
        bytes memory data = abi.encodePacked("hello");

        vm.prank(deployer);
        bool success = token.transferAndCall(address(receiver), amount, data);

        assertTrue(success);
        assertEq(token.balanceOf(address(receiver)), amount);
        assertEq(receiver.lastOperator(), deployer);
        assertEq(receiver.lastFrom(), deployer);
        assertEq(receiver.lastAmount(), amount);
        assertEq(receiver.lastData(), data);
    }

    function test_transferAndCall_toEOA() public {
        uint256 amount = 100 * 1e18;
        vm.etch(alice, ""); // Ensure alice is EOA (no code) — needed for fork tests

        vm.prank(deployer);
        bool success = token.transferAndCall(alice, amount, "");

        assertTrue(success);
        assertEq(token.balanceOf(alice), amount);
    }

    function test_transferAndCall_revertIfRejected() public {
        RejectingReceiver receiver = new RejectingReceiver();
        uint256 amount = 100 * 1e18;

        vm.prank(deployer);
        vm.expectRevert(AWPToken.InvalidCallback.selector);
        token.transferAndCall(address(receiver), amount, "");
    }

    // ── ERC1363 approveAndCall tests ──

    function test_approveAndCall_toContract() public {
        MockSpender spender = new MockSpender();
        uint256 amount = 200 * 1e18;
        bytes memory data = abi.encodePacked("approve");

        vm.prank(deployer);
        bool success = token.approveAndCall(address(spender), amount, data);

        assertTrue(success);
        assertEq(token.allowance(deployer, address(spender)), amount);
        assertEq(spender.lastOwner(), deployer);
        assertEq(spender.lastAmount(), amount);
        assertEq(spender.lastData(), data);
    }

    function test_approveAndCall_toEOA() public {
        uint256 amount = 200 * 1e18;
        vm.etch(alice, ""); // Ensure alice is EOA (no code) — needed for fork tests

        vm.prank(deployer);
        bool success = token.approveAndCall(alice, amount, "");

        assertTrue(success);
        assertEq(token.allowance(deployer, alice), amount);
    }

    function test_approveAndCall_revertIfRejected() public {
        RejectingSpender spender = new RejectingSpender();
        uint256 amount = 200 * 1e18;

        vm.prank(deployer);
        vm.expectRevert(AWPToken.InvalidCallback.selector);
        token.approveAndCall(address(spender), amount, "");
    }

    // ── Fuzz tests ──

    // ── initialMint tests ──

    function test_initialMint_success() public {
        // Deploy a fresh token (no initialMint in setUp)
        address freshDeployer = makeAddr("freshDeployer");
        vm.prank(freshDeployer);
        AWPToken freshToken = new AWPToken("AWP", "AWP", freshDeployer);

        assertEq(freshToken.totalSupply(), 0);
        assertEq(freshToken.balanceOf(freshDeployer), 0);

        vm.prank(freshDeployer);
        freshToken.initialMint(200_000_000 * 1e18);

        assertEq(freshToken.balanceOf(freshDeployer), 200_000_000 * 1e18);
        assertEq(freshToken.totalSupply(), 200_000_000 * 1e18);
        assertTrue(freshToken.initialMinted());
    }

    function test_initialMint_onlyAdmin() public {
        address freshDeployer = makeAddr("freshDeployer2");
        vm.prank(freshDeployer);
        AWPToken freshToken = new AWPToken("AWP", "AWP", freshDeployer);

        vm.prank(alice);
        vm.expectRevert(AWPToken.NotAdmin.selector);
        freshToken.initialMint(100 * 1e18);
    }

    function test_initialMint_calledOnce() public {
        // setUp already called initialMint on `token`
        assertTrue(token.initialMinted());

        vm.prank(deployer);
        vm.expectRevert(AWPToken.AlreadyMinted.selector);
        token.initialMint(1 * 1e18);
    }

    function test_initialMint_zero() public {
        address freshDeployer = makeAddr("freshDeployer3");
        vm.prank(freshDeployer);
        AWPToken freshToken = new AWPToken("AWP", "AWP", freshDeployer);

        vm.prank(freshDeployer);
        freshToken.initialMint(0);

        assertTrue(freshToken.initialMinted());
        assertEq(freshToken.totalSupply(), 0);
        assertEq(freshToken.balanceOf(freshDeployer), 0);
    }

    function test_initialMint_exceedsMaxSupply() public {
        address freshDeployer = makeAddr("freshDeployer4");
        vm.prank(freshDeployer);
        AWPToken freshToken = new AWPToken("AWP", "AWP", freshDeployer);

        vm.prank(freshDeployer);
        vm.expectRevert(AWPToken.ExceedsMaxSupply.selector);
        freshToken.initialMint(MAX_SUPPLY + 1);
    }

    function test_noConstructorMint() public {
        // After constructor (before initialMint), totalSupply == 0
        address freshDeployer = makeAddr("freshDeployer5");
        vm.prank(freshDeployer);
        AWPToken freshToken = new AWPToken("AWP", "AWP", freshDeployer);

        assertEq(freshToken.totalSupply(), 0);
        assertEq(freshToken.balanceOf(freshDeployer), 0);
        assertFalse(freshToken.initialMinted());
    }

    // ── Fuzz tests ──

    function testFuzz_mint_respectsMaxSupply(uint256 amount) public {
        vm.prank(deployer);
        token.addMinter(minter);

        uint256 remaining = MAX_SUPPLY - INITIAL_MINT;
        amount = bound(amount, 1, remaining);

        vm.prank(minter);
        token.mint(alice, amount);

        assertLe(token.totalSupply(), MAX_SUPPLY);
    }

    function testFuzz_burn_reducesSupply(uint256 burnAmount) public {
        burnAmount = bound(burnAmount, 1, INITIAL_MINT);

        vm.prank(deployer);
        token.burn(burnAmount);

        assertEq(token.totalSupply(), INITIAL_MINT - burnAmount);
    }
}
