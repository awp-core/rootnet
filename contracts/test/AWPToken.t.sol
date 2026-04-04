// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {AWPToken} from "../src/token/AWPToken.sol";

contract AWPTokenTest is Test {
    AWPToken public awp;
    address public deployer = address(this);

    function setUp() public {
        awp = new AWPToken("AWP Token", "AWP", deployer);
    }

    function test_name() public view {
        assertEq(awp.name(), "AWP Token");
    }

    function test_initialMint() public {
        awp.initialMint(100e18);
        assertEq(awp.totalSupply(), 100e18);
        assertEq(awp.balanceOf(deployer), 100e18);
    }

    function test_initialMint_twice_reverts() public {
        awp.initialMint(100e18);
        vm.expectRevert();
        awp.initialMint(100e18);
    }

    function test_addMinter() public {
        awp.addMinter(address(1));
        assertTrue(awp.minters(address(1)));
    }

    function test_renounceAdmin() public {
        awp.renounceAdmin();
        assertEq(awp.admin(), address(0));
    }

    function test_mint_byMinter() public {
        awp.addMinter(address(1));
        vm.prank(address(1));
        awp.mint(deployer, 100e18);
        assertEq(awp.totalSupply(), 100e18);
    }

    function test_maxSupply() public view {
        assertEq(awp.MAX_SUPPLY(), 10_000_000_000 * 1e18);
    }
}
