// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {AlphaToken} from "../src/token/AlphaToken.sol";
import {AlphaTokenFactory} from "../src/token/AlphaTokenFactory.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {IERC1363Receiver, IERC1363Spender} from "../src/interfaces/IERC1363Receiver.sol";

// ── Helper contracts ──

contract MockAlphaReceiver is IERC1363Receiver {
    address public lastOperator;
    uint256 public lastAmount;

    function onTransferReceived(address operator, address from, uint256 amount, bytes calldata)
        external
        returns (bytes4)
    {
        lastOperator = operator;
        lastAmount = amount;
        return IERC1363Receiver.onTransferReceived.selector;
    }
}

contract MockAlphaSpender is IERC1363Spender {
    address public lastOwner;
    uint256 public lastAmount;

    function onApprovalReceived(address owner, uint256 amount, bytes calldata) external returns (bytes4) {
        lastOwner = owner;
        lastAmount = amount;
        return IERC1363Spender.onApprovalReceived.selector;
    }
}

contract RejectingAlphaReceiver is IERC1363Receiver {
    function onTransferReceived(address, address, uint256, bytes calldata) external pure returns (bytes4) {
        return bytes4(0xdeadbeef);
    }
}

// ── AlphaToken tests ──

contract AlphaTokenTest is Test {
    AlphaToken public token;
    address public admin;
    address public alice;
    address public bob;
    address public subnetManager;

    uint256 constant MAX_SUPPLY = 10_000_000_000 * 1e18;
    uint256 constant SUBNET_ID = 42;

    function setUp() public {
        admin = makeAddr("admin");
        alice = makeAddr("alice");
        bob = makeAddr("bob");
        subnetManager = makeAddr("subnetManager");

        // Create and initialize directly (constructor no longer disables initializers)
        token = new AlphaToken();
        token.initialize("Alpha 42", "ALPHA42", SUBNET_ID, admin);
    }

    // ── initialize tests ──

    function test_initialize_setsState() public view {
        assertEq(token.name(), "Alpha 42");
        assertEq(token.symbol(), "ALPHA42");
        assertEq(token.subnetId(), SUBNET_ID);
        assertEq(token.admin(), admin);
    }

    function test_initialize_adminIsMinter() public view {
        assertTrue(token.minters(admin));
    }

    function test_initialize_cannotBeCalledTwice() public {
        vm.expectRevert();
        token.initialize("X", "X", 1, admin);
    }

    // ── mint tests ──

    function test_mint_byAdmin() public {
        uint256 amount = 1000 * 1e18;

        vm.prank(admin);
        token.mint(alice, amount);

        assertEq(token.balanceOf(alice), amount);
        assertEq(token.totalSupply(), amount);
    }

    function test_mint_revertIfNotMinter() public {
        vm.prank(alice);
        vm.expectRevert(AlphaToken.NotMinter.selector);
        token.mint(alice, 100);
    }

    function test_mint_revertIfExceedsMaxSupply() public {
        vm.prank(admin);
        vm.expectRevert(AlphaToken.ExceedsMaxSupply.selector);
        token.mint(alice, MAX_SUPPLY + 1);
    }

    function test_mint_exactlyReachesMaxSupply() public {
        vm.prank(admin);
        token.mint(alice, MAX_SUPPLY);
        assertEq(token.totalSupply(), MAX_SUPPLY);
    }

    function test_mint_revertIfMinterPaused() public {
        vm.prank(admin);
        token.setMinterPaused(admin, true);

        vm.prank(admin);
        vm.expectRevert(AlphaToken.MinterPaused.selector);
        token.mint(alice, 100);
    }

    // ── setSubnetMinter tests ──

    function test_setSubnetMinter_locksForever() public {
        vm.prank(admin);
        token.setSubnetMinter(subnetManager);

        assertTrue(token.mintersLocked());
        assertTrue(token.minters(subnetManager));
        assertFalse(token.minters(admin)); // admin loses minting rights
    }

    function test_setSubnetMinter_revertIfNotAdmin() public {
        vm.prank(alice);
        vm.expectRevert(AlphaToken.NotAdmin.selector);
        token.setSubnetMinter(subnetManager);
    }

    function test_setSubnetMinter_revertIfAlreadyLocked() public {
        vm.prank(admin);
        token.setSubnetMinter(subnetManager);

        vm.prank(admin);
        vm.expectRevert(AlphaToken.MintersLocked.selector);
        token.setSubnetMinter(alice);
    }

    function test_setSubnetMinter_adminCannotMintAfterLock() public {
        vm.prank(admin);
        token.setSubnetMinter(subnetManager);

        vm.prank(admin);
        vm.expectRevert(AlphaToken.NotMinter.selector);
        token.mint(alice, 100);
    }

    function test_setSubnetMinter_newMinterCanMint() public {
        vm.prank(admin);
        token.setSubnetMinter(subnetManager);

        vm.warp(block.timestamp + 1 days);
        vm.prank(subnetManager);
        token.mint(alice, 500 * 1e18);

        assertEq(token.balanceOf(alice), 500 * 1e18);
    }

    function test_setSubnetMinter_withZeroAddress_reverts() public {
        // address(0) is not allowed as minter, preventing token lock-up
        vm.prank(admin);
        vm.expectRevert(AlphaToken.NotMinter.selector);
        token.setSubnetMinter(address(0));
    }

    // ── setMinterPaused tests ──

    function test_setMinterPaused_pauseAndUnpause() public {
        vm.prank(admin);
        token.setSubnetMinter(subnetManager);

        // Pause
        vm.prank(admin);
        token.setMinterPaused(subnetManager, true);

        vm.prank(subnetManager);
        vm.expectRevert(AlphaToken.MinterPaused.selector);
        token.mint(alice, 100);

        // Resume
        vm.prank(admin);
        token.setMinterPaused(subnetManager, false);

        vm.warp(block.timestamp + 1 days);
        vm.prank(subnetManager);
        token.mint(alice, 100 * 1e18);
        assertEq(token.balanceOf(alice), 100 * 1e18);
    }

    function test_setMinterPaused_revertIfNotAdmin() public {
        vm.prank(alice);
        vm.expectRevert(AlphaToken.NotAdmin.selector);
        token.setMinterPaused(admin, true);
    }

    // ── burn tests ──

    function test_burn_reducesSupply() public {
        uint256 mintAmount = 1000 * 1e18;
        uint256 burnAmount = 400 * 1e18;

        vm.prank(admin);
        token.mint(alice, mintAmount);

        vm.prank(alice);
        token.burn(burnAmount);

        assertEq(token.balanceOf(alice), mintAmount - burnAmount);
        assertEq(token.totalSupply(), mintAmount - burnAmount);
    }

    function test_burnFrom_withApproval() public {
        uint256 amount = 1000 * 1e18;

        vm.prank(admin);
        token.mint(alice, amount);

        vm.prank(alice);
        token.approve(bob, 500 * 1e18);

        vm.prank(bob);
        token.burnFrom(alice, 500 * 1e18);

        assertEq(token.balanceOf(alice), 500 * 1e18);
    }

    // ── ERC1363 callback tests ──

    function test_transferAndCall_success() public {
        MockAlphaReceiver receiver = new MockAlphaReceiver();
        uint256 amount = 100 * 1e18;

        vm.prank(admin);
        token.mint(alice, amount);

        vm.prank(alice);
        bool success = token.transferAndCall(address(receiver), amount, "");

        assertTrue(success);
        assertEq(token.balanceOf(address(receiver)), amount);
        assertEq(receiver.lastOperator(), alice);
        assertEq(receiver.lastAmount(), amount);
    }

    function test_transferAndCall_toEOA() public {
        uint256 amount = 100 * 1e18;

        vm.prank(admin);
        token.mint(alice, amount);

        vm.prank(alice);
        bool success = token.transferAndCall(bob, amount, "");

        assertTrue(success);
        assertEq(token.balanceOf(bob), amount);
    }

    function test_transferAndCall_revertIfRejected() public {
        RejectingAlphaReceiver receiver = new RejectingAlphaReceiver();
        uint256 amount = 100 * 1e18;

        vm.prank(admin);
        token.mint(alice, amount);

        vm.prank(alice);
        vm.expectRevert(AlphaToken.InvalidCallback.selector);
        token.transferAndCall(address(receiver), amount, "");
    }

    function test_approveAndCall_success() public {
        MockAlphaSpender spender = new MockAlphaSpender();
        uint256 amount = 200 * 1e18;

        vm.prank(admin);
        token.mint(alice, amount);

        vm.prank(alice);
        bool success = token.approveAndCall(address(spender), amount, "");

        assertTrue(success);
        assertEq(token.allowance(alice, address(spender)), amount);
        assertEq(spender.lastOwner(), alice);
        assertEq(spender.lastAmount(), amount);
    }

    // ── Fuzz tests ──

    function testFuzz_mint_respectsMaxSupply(uint256 amount) public {
        amount = bound(amount, 1, MAX_SUPPLY);

        vm.prank(admin);
        token.mint(alice, amount);

        assertLe(token.totalSupply(), MAX_SUPPLY);
    }
}

// ── AlphaTokenFactory tests ──

contract AlphaTokenFactoryTest is Test {
    AlphaTokenFactory public factory;
    address public deployer;
    address public rootNet;
    address public alice;

    function setUp() public {
        deployer = makeAddr("deployer");
        rootNet = makeAddr("rootNet");
        alice = makeAddr("alice");

        vm.prank(deployer);
        factory = new AlphaTokenFactory(deployer, 0);
    }

    // ── Constructor tests ──

    function test_constructor_setsOwner() public view {
        assertEq(factory.owner(), deployer);
    }

    function test_constructor_notConfigured() public view {
        assertFalse(factory.configured());
        assertEq(factory.rootNet(), address(0));
    }

    // ── setAddresses tests ──

    function test_setAddresses_configuresAndRenounces() public {
        vm.prank(deployer);
        factory.setAddresses(rootNet);

        assertTrue(factory.configured());
        assertEq(factory.rootNet(), rootNet);
        assertEq(factory.owner(), address(0)); // ownership has been renounced
    }

    function test_setAddresses_revertIfNotOwner() public {
        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        factory.setAddresses(rootNet);
    }

    function test_setAddresses_cannotBeCalledTwice() public {
        vm.prank(deployer);
        factory.setAddresses(rootNet);

        // Ownership has been renounced; calling again should fail
        vm.prank(deployer);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, deployer));
        factory.setAddresses(alice);
    }

    // ── deploy tests ──

    function test_deploy_createsToken() public {
        vm.prank(deployer);
        factory.setAddresses(rootNet);

        vm.prank(rootNet);
        address token = factory.deploy(1, "Subnet Alpha 1", "SA1", rootNet, bytes32(0));

        assertTrue(token != address(0));

        AlphaToken alphaToken = AlphaToken(token);
        assertEq(alphaToken.name(), "Subnet Alpha 1");
        assertEq(alphaToken.symbol(), "SA1");
        assertEq(alphaToken.subnetId(), 1);
        assertEq(alphaToken.admin(), rootNet);
    }

    function test_deploy_revertIfNotRootNet() public {
        vm.prank(deployer);
        factory.setAddresses(rootNet);

        vm.prank(alice);
        vm.expectRevert(AlphaTokenFactory.NotRootNet.selector);
        factory.deploy(1, "X", "X", alice, bytes32(0));
    }

    function test_deploy_multipleSubnets() public {
        vm.prank(deployer);
        factory.setAddresses(rootNet);

        vm.startPrank(rootNet);
        address token1 = factory.deploy(1, "Alpha 1", "A1", rootNet, bytes32(0));
        address token2 = factory.deploy(2, "Alpha 2", "A2", rootNet, bytes32(0));
        address token3 = factory.deploy(3, "Alpha 3", "A3", rootNet, bytes32(0));
        vm.stopPrank();

        assertTrue(token1 != token2);
        assertTrue(token2 != token3);

        assertEq(AlphaToken(token1).subnetId(), 1);
        assertEq(AlphaToken(token2).subnetId(), 2);
        assertEq(AlphaToken(token3).subnetId(), 3);
    }

    function test_deploy_adminIsMinter() public {
        vm.prank(deployer);
        factory.setAddresses(rootNet);

        vm.prank(rootNet);
        address token = factory.deploy(1, "Alpha", "A", rootNet, bytes32(0));

        AlphaToken alphaToken = AlphaToken(token);
        assertTrue(alphaToken.minters(rootNet));
    }

    function test_deploy_canMint() public {
        vm.prank(deployer);
        factory.setAddresses(rootNet);

        vm.prank(rootNet);
        address token = factory.deploy(1, "Alpha", "A", rootNet, bytes32(0));

        AlphaToken alphaToken = AlphaToken(token);
        vm.prank(rootNet);
        alphaToken.mint(alice, 1000e18);
        assertEq(alphaToken.balanceOf(alice), 1000e18);
    }

    function test_deploy_cannotReinitialize() public {
        vm.prank(deployer);
        factory.setAddresses(rootNet);

        vm.prank(rootNet);
        address token = factory.deploy(1, "Alpha", "A", rootNet, bytes32(0));

        vm.expectRevert();
        AlphaToken(token).initialize("Hacked", "H", 999, alice);
    }

    // ── Full lifecycle integration test ──

    function test_fullLifecycle_deployAndLockMinter() public {
        // 1. Configure factory
        vm.prank(deployer);
        factory.setAddresses(rootNet);

        // 2. RootNet deploys AlphaToken
        vm.prank(rootNet);
        address token = factory.deploy(1, "Alpha 1", "A1", rootNet, bytes32(0));
        AlphaToken alphaToken = AlphaToken(token);

        // 3. RootNet (admin) mints some tokens first
        vm.prank(rootNet);
        alphaToken.mint(alice, 1000 * 1e18);
        assertEq(alphaToken.balanceOf(alice), 1000 * 1e18);

        // 4. Set subnet contract as minter, RootNet relinquishes minting rights
        address subnetManager = makeAddr("subnetManager");
        vm.prank(rootNet);
        alphaToken.setSubnetMinter(subnetManager);

        // 5. RootNet can no longer mint
        vm.prank(rootNet);
        vm.expectRevert(AlphaToken.NotMinter.selector);
        alphaToken.mint(alice, 100);

        // 6. Subnet contract can mint
        vm.warp(block.timestamp + 1 days);
        vm.prank(subnetManager);
        alphaToken.mint(alice, 500 * 1e18);
        assertEq(alphaToken.balanceOf(alice), 1500 * 1e18);

        // 7. Ban subnet (pause minter)
        vm.prank(rootNet);
        alphaToken.setMinterPaused(subnetManager, true);

        vm.prank(subnetManager);
        vm.expectRevert(AlphaToken.MinterPaused.selector);
        alphaToken.mint(alice, 100);

        // 8. Unban
        vm.prank(rootNet);
        alphaToken.setMinterPaused(subnetManager, false);

        vm.warp(block.timestamp + 1 days);
        vm.prank(subnetManager);
        alphaToken.mint(alice, 200 * 1e18);
        assertEq(alphaToken.balanceOf(alice), 1700 * 1e18);

        // 9. Permanently locked; cannot set minter again
        vm.prank(rootNet);
        vm.expectRevert(AlphaToken.MintersLocked.selector);
        alphaToken.setSubnetMinter(alice);
    }

    // ── Fuzz tests ──

    function testFuzz_deploy_differentSubnetIds(uint256 subnetId) public {
        subnetId = bound(subnetId, 0, 10000);

        vm.prank(deployer);
        factory.setAddresses(rootNet);

        vm.prank(rootNet);
        address token = factory.deploy(subnetId, "Alpha", "A", rootNet, bytes32(0));

        assertEq(AlphaToken(token).subnetId(), subnetId);
    }
}

// ── VanityHarness: exposes _validateVanityAddress for testing ──

contract VanityHarness is AlphaTokenFactory {
    constructor(uint64 rule) AlphaTokenFactory(msg.sender, rule) {}

    function checkVanity(address addr) external view {
        _validateVanityAddress(addr);
    }
}

// ── Vanity rule tests ──

contract VanityRuleTest is Test {
    function test_vanityRule_digit() public {
        // Rule: position 0 must be digit '0' (value=0), all others wildcard (255)
        // Packed: [0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF] = 0x00FFFFFFFFFFFFFF
        VanityHarness h = new VanityHarness(0x00FFFFFFFFFFFFFF);

        // Address starting with 0x01... (first nibble = 0) should pass
        // Use uint256 literal to avoid Solidity address-literal checksum parsing
        address addrStartsWith0 = address(uint160(uint256(0x01) << 152));
        h.checkVanity(addrStartsWith0);

        // Address starting with 0xA0... (first nibble = 0xA != 0) should fail
        address addrStartsWithA = address(uint160(uint256(0xA0) << 152));
        vm.expectRevert(AlphaTokenFactory.InvalidVanityAddress.selector);
        h.checkVanity(addrStartsWithA);
    }

    function test_vanityRule_wildcardAll() public {
        // Rule: all wildcards (0xFF per position = wildcard >=22)
        // vanityRule != 0 so validation is entered, but every position is wildcard → always passes
        VanityHarness h = new VanityHarness(0xFFFFFFFFFFFFFFFF);

        // Any address should pass regardless of content
        h.checkVanity(address(uint160(uint256(0x12) << 152)));
        h.checkVanity(address(uint160(uint256(0xAB) << 152)));
    }

    function test_vanityRule_zero_skipsValidation() public {
        // vanityRule == 0 means no validation in deploy(), but _validateVanityAddress
        // can still be called directly. With rule=0, all positions have expected=0 (digit '0')
        // which would fail for most addresses. This test confirms rule=0 behaviour in deploy().
        VanityHarness h = new VanityHarness(0);
        address rn = makeAddr("rootNet");
        h.setAddresses(rn);

        // deploy() skips _validateVanityAddress when vanityRule==0, so any address is accepted
        vm.prank(rn);
        address token = h.deploy(1, "Alpha", "A", rn, bytes32(0));
        assertTrue(token != address(0));
    }
}
