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
    address public worknetManager;

    uint256 constant MAX_SUPPLY = 10_000_000_000 * 1e18;
    uint256 constant WORKNET_ID = 42;

    function setUp() public {
        admin = makeAddr("admin");
        alice = makeAddr("alice");
        bob = makeAddr("bob");
        worknetManager = makeAddr("worknetManager");

        // Create and initialize directly (constructor no longer disables initializers)
        token = new AlphaToken();
        token.initialize("Alpha 42", "ALPHA42", WORKNET_ID, admin);
    }

    // ── initialize tests ──

    function test_initialize_setsState() public view {
        assertEq(token.name(), "Alpha 42");
        assertEq(token.symbol(), "ALPHA42");
        assertEq(token.worknetId(), WORKNET_ID);
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

    // ── setWorknetMinter tests ──

    function test_setWorknetMinter_locksForever() public {
        vm.prank(admin);
        token.setWorknetMinter(worknetManager);

        assertTrue(token.mintersLocked());
        assertTrue(token.minters(worknetManager));
        assertFalse(token.minters(admin)); // admin loses minting rights
    }

    function test_setWorknetMinter_revertIfNotAdmin() public {
        vm.prank(alice);
        vm.expectRevert(AlphaToken.NotAdmin.selector);
        token.setWorknetMinter(worknetManager);
    }

    function test_setWorknetMinter_revertIfAlreadyLocked() public {
        vm.prank(admin);
        token.setWorknetMinter(worknetManager);

        vm.prank(admin);
        vm.expectRevert(AlphaToken.MintersLocked.selector);
        token.setWorknetMinter(alice);
    }

    function test_setWorknetMinter_adminCannotMintAfterLock() public {
        vm.prank(admin);
        token.setWorknetMinter(worknetManager);

        vm.prank(admin);
        vm.expectRevert(AlphaToken.NotMinter.selector);
        token.mint(alice, 100);
    }

    function test_setWorknetMinter_newMinterCanMint() public {
        vm.prank(admin);
        token.setWorknetMinter(worknetManager);

        vm.warp(block.timestamp + 1 days);
        vm.prank(worknetManager);
        token.mint(alice, 500 * 1e18);

        assertEq(token.balanceOf(alice), 500 * 1e18);
    }

    function test_setWorknetMinter_withZeroAddress_reverts() public {
        // address(0) is not allowed as minter, preventing token lock-up
        vm.prank(admin);
        vm.expectRevert(AlphaToken.ZeroAddress.selector);
        token.setWorknetMinter(address(0));
    }

    // ── setMinterPaused tests ──

    function test_setMinterPaused_pauseAndUnpause() public {
        vm.prank(admin);
        token.setWorknetMinter(worknetManager);

        // Pause
        vm.prank(admin);
        token.setMinterPaused(worknetManager, true);

        vm.prank(worknetManager);
        vm.expectRevert(AlphaToken.MinterPaused.selector);
        token.mint(alice, 100);

        // Resume
        vm.prank(admin);
        token.setMinterPaused(worknetManager, false);

        vm.warp(block.timestamp + 1 days);
        vm.prank(worknetManager);
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
        vm.etch(bob, ""); // Ensure bob is EOA (no code) — needed for fork tests

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
    address public awpRegistry;
    address public alice;

    function setUp() public {
        deployer = makeAddr("deployer");
        awpRegistry = makeAddr("awpRegistry");
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
        assertEq(factory.awpRegistry(), address(0));
    }

    // ── setAddresses tests ──

    function test_setAddresses_configuresAndRenounces() public {
        vm.prank(deployer);
        factory.setAddresses(awpRegistry);

        assertTrue(factory.configured());
        assertEq(factory.awpRegistry(), awpRegistry);
        assertEq(factory.owner(), address(0)); // ownership has been renounced
    }

    function test_setAddresses_revertIfNotOwner() public {
        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        factory.setAddresses(awpRegistry);
    }

    function test_setAddresses_cannotBeCalledTwice() public {
        vm.prank(deployer);
        factory.setAddresses(awpRegistry);

        // Ownership has been renounced; calling again should fail
        vm.prank(deployer);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, deployer));
        factory.setAddresses(alice);
    }

    // ── deploy tests ──

    function test_deploy_createsToken() public {
        vm.prank(deployer);
        factory.setAddresses(awpRegistry);

        vm.prank(awpRegistry);
        address token = factory.deploy(1, "Worknet Alpha 1", "SA1", awpRegistry, bytes32(0));

        assertTrue(token != address(0));

        AlphaToken alphaToken = AlphaToken(token);
        assertEq(alphaToken.name(), "Worknet Alpha 1");
        assertEq(alphaToken.symbol(), "SA1");
        assertEq(alphaToken.worknetId(), 1);
        assertEq(alphaToken.admin(), awpRegistry);
    }

    function test_deploy_revertIfNotAWPRegistry() public {
        vm.prank(deployer);
        factory.setAddresses(awpRegistry);

        vm.prank(alice);
        vm.expectRevert(AlphaTokenFactory.NotAWPRegistry.selector);
        factory.deploy(1, "X", "X", alice, bytes32(0));
    }

    function test_deploy_multipleWorknets() public {
        vm.prank(deployer);
        factory.setAddresses(awpRegistry);

        vm.startPrank(awpRegistry);
        address token1 = factory.deploy(1, "Alpha 1", "A1", awpRegistry, bytes32(0));
        address token2 = factory.deploy(2, "Alpha 2", "A2", awpRegistry, bytes32(0));
        address token3 = factory.deploy(3, "Alpha 3", "A3", awpRegistry, bytes32(0));
        vm.stopPrank();

        assertTrue(token1 != token2);
        assertTrue(token2 != token3);

        assertEq(AlphaToken(token1).worknetId(), 1);
        assertEq(AlphaToken(token2).worknetId(), 2);
        assertEq(AlphaToken(token3).worknetId(), 3);
    }

    function test_deploy_adminIsMinter() public {
        vm.prank(deployer);
        factory.setAddresses(awpRegistry);

        vm.prank(awpRegistry);
        address token = factory.deploy(1, "Alpha", "A", awpRegistry, bytes32(0));

        AlphaToken alphaToken = AlphaToken(token);
        assertTrue(alphaToken.minters(awpRegistry));
    }

    function test_deploy_canMint() public {
        vm.prank(deployer);
        factory.setAddresses(awpRegistry);

        vm.prank(awpRegistry);
        address token = factory.deploy(1, "Alpha", "A", awpRegistry, bytes32(0));

        AlphaToken alphaToken = AlphaToken(token);
        vm.prank(awpRegistry);
        alphaToken.mint(alice, 1000e18);
        assertEq(alphaToken.balanceOf(alice), 1000e18);
    }

    function test_deploy_cannotReinitialize() public {
        vm.prank(deployer);
        factory.setAddresses(awpRegistry);

        vm.prank(awpRegistry);
        address token = factory.deploy(1, "Alpha", "A", awpRegistry, bytes32(0));

        vm.expectRevert();
        AlphaToken(token).initialize("Hacked", "H", 999, alice);
    }

    // ── Full lifecycle integration test ──

    function test_fullLifecycle_deployAndLockMinter() public {
        // 1. Configure factory
        vm.prank(deployer);
        factory.setAddresses(awpRegistry);

        // 2. AWPRegistry deploys AlphaToken
        vm.prank(awpRegistry);
        address token = factory.deploy(1, "Alpha 1", "A1", awpRegistry, bytes32(0));
        AlphaToken alphaToken = AlphaToken(token);

        // 3. AWPRegistry (admin) mints some tokens first
        vm.prank(awpRegistry);
        alphaToken.mint(alice, 1000 * 1e18);
        assertEq(alphaToken.balanceOf(alice), 1000 * 1e18);

        // 4. Set worknet contract as minter, AWPRegistry relinquishes minting rights
        address worknetManager = makeAddr("worknetManager");
        vm.prank(awpRegistry);
        alphaToken.setWorknetMinter(worknetManager);

        // 5. AWPRegistry can no longer mint
        vm.prank(awpRegistry);
        vm.expectRevert(AlphaToken.NotMinter.selector);
        alphaToken.mint(alice, 100);

        // 6. Worknet contract can mint
        vm.warp(block.timestamp + 1 days);
        vm.prank(worknetManager);
        alphaToken.mint(alice, 500 * 1e18);
        assertEq(alphaToken.balanceOf(alice), 1500 * 1e18);

        // 7. Ban worknet (pause minter)
        vm.prank(awpRegistry);
        alphaToken.setMinterPaused(worknetManager, true);

        vm.prank(worknetManager);
        vm.expectRevert(AlphaToken.MinterPaused.selector);
        alphaToken.mint(alice, 100);

        // 8. Unban
        vm.prank(awpRegistry);
        alphaToken.setMinterPaused(worknetManager, false);

        vm.warp(block.timestamp + 1 days);
        vm.prank(worknetManager);
        alphaToken.mint(alice, 200 * 1e18);
        assertEq(alphaToken.balanceOf(alice), 1700 * 1e18);

        // 9. Permanently locked; cannot set minter again
        vm.prank(awpRegistry);
        vm.expectRevert(AlphaToken.MintersLocked.selector);
        alphaToken.setWorknetMinter(alice);
    }

    // ── Fuzz tests ──

    function testFuzz_deploy_differentWorknetIds(uint256 worknetId) public {
        worknetId = bound(worknetId, 0, 10000);

        vm.prank(deployer);
        factory.setAddresses(awpRegistry);

        vm.prank(awpRegistry);
        address token = factory.deploy(worknetId, "Alpha", "A", awpRegistry, bytes32(0));

        assertEq(AlphaToken(token).worknetId(), worknetId);
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
        address rn = makeAddr("awpRegistry");
        h.setAddresses(rn);

        // deploy() skips _validateVanityAddress when vanityRule==0, so any address is accepted
        vm.prank(rn);
        address token = h.deploy(1, "Alpha", "A", rn, bytes32(0));
        assertTrue(token != address(0));
    }
}
