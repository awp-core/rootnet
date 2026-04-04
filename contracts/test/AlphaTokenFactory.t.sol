// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {AlphaTokenFactory} from "../src/token/AlphaTokenFactory.sol";
import {AlphaToken} from "../src/token/AlphaToken.sol";

contract AlphaTokenFactoryTest is Test {
    AlphaTokenFactory public factory;
    address public deployer = address(this);
    address public registry = makeAddr("registry");

    function setUp() public {
        factory = new AlphaTokenFactory(deployer, 0); // vanityRule=0 → no validation
        factory.setAddresses(registry);
    }

    // ═══════════════════════════════════════════════
    //  Construction & Configuration
    // ═══════════════════════════════════════════════

    function test_configured() public view {
        assertTrue(factory.configured());
        assertEq(factory.awpRegistry(), registry);
        assertEq(factory.vanityRule(), 0);
    }

    function test_ownerRenounced() public view {
        assertEq(factory.owner(), address(0));
    }

    function test_bytecodeHash() public view {
        bytes32 expected = keccak256(type(AlphaToken).creationCode);
        assertEq(factory.ALPHA_BYTECODE_HASH(), expected);
    }

    // ═══════════════════════════════════════════════
    //  Deploy
    // ═══════════════════════════════════════════════

    function test_deploy() public {
        vm.prank(registry);
        address token = factory.deploy(1, "Alpha1", "A1", registry, bytes32(0));

        assertTrue(token != address(0));
        AlphaToken alpha = AlphaToken(token);
        assertEq(alpha.name(), "Alpha1");
        assertEq(alpha.symbol(), "A1");
        assertEq(alpha.worknetId(), 1);
        assertEq(alpha.admin(), registry);
        assertTrue(alpha.minters(registry));
    }

    function test_deploy_customSalt() public {
        bytes32 salt = bytes32(uint256(42));
        vm.prank(registry);
        address token = factory.deploy(1, "Alpha1", "A1", registry, salt);
        assertTrue(token != address(0));
    }

    function test_deploy_notRegistry_reverts() public {
        vm.expectRevert(AlphaTokenFactory.NotAWPRegistry.selector);
        factory.deploy(1, "Alpha1", "A1", registry, bytes32(0));
    }

    function test_deploy_deterministic() public {
        bytes32 salt = bytes32(uint256(100));
        address predicted = factory.predictDeployAddress(salt);

        vm.prank(registry);
        address deployed = factory.deploy(1, "Alpha", "A", registry, salt);

        assertEq(deployed, predicted);
    }

    function test_deploy_defaultSalt_usesWorknetId() public {
        uint256 wid = 845300000001;
        address predicted = factory.predictDeployAddress(bytes32(wid));

        vm.prank(registry);
        address deployed = factory.deploy(wid, "Alpha", "A", registry, bytes32(0));

        assertEq(deployed, predicted);
    }

    function test_deploy_sameSalt_reverts() public {
        vm.prank(registry);
        factory.deploy(1, "A1", "A1", registry, bytes32(0));

        vm.prank(registry);
        vm.expectRevert(); // CREATE2 collision
        factory.deploy(1, "A2", "A2", registry, bytes32(0));
    }

    // ═══════════════════════════════════════════════
    //  Vanity Validation
    // ═══════════════════════════════════════════════

    function test_vanityRule_deploy_noRule() public {
        // vanityRule=0 means no validation — already setUp with 0
        vm.prank(registry);
        address token = factory.deploy(999, "V", "V", registry, bytes32(0));
        assertTrue(token != address(0));
    }

    function test_vanityFactory_withRule() public {
        // Create a factory with all-wildcard rule (0xFFFFFFFFFFFFFFFF)
        // All positions >= 22 are wildcard → should pass any address
        uint64 allWild = 0xFFFFFFFFFFFFFFFF;
        AlphaTokenFactory vFactory = new AlphaTokenFactory(deployer, allWild);
        vFactory.setAddresses(registry);

        vm.prank(registry);
        address token = vFactory.deploy(1, "V", "V", registry, bytes32(0));
        assertTrue(token != address(0));
    }

    // ═══════════════════════════════════════════════
    //  setAddresses
    // ═══════════════════════════════════════════════

    function test_setAddresses_notOwner_reverts() public {
        AlphaTokenFactory f2 = new AlphaTokenFactory(deployer, 0);
        vm.prank(makeAddr("random"));
        vm.expectRevert();
        f2.setAddresses(registry);
    }

    function test_setAddresses_onlyOnce() public {
        // factory already configured in setUp, owner renounced
        vm.expectRevert(); // Ownable: caller is not the owner (owner is address(0))
        factory.setAddresses(makeAddr("other"));
    }
}
