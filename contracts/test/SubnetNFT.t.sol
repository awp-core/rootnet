// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {SubnetNFT} from "../src/core/SubnetNFT.sol";
import {IERC721Errors} from "@openzeppelin/contracts/interfaces/draft-IERC6093.sol";

contract SubnetNFTTest is Test {
    SubnetNFT public nft;

    address public awpRegistry = makeAddr("awpRegistry");
    address public alice = makeAddr("alice");
    address public bob = makeAddr("bob");

    function setUp() public {
        nft = new SubnetNFT("AWP Subnet", "SUBNET", awpRegistry);
    }

    // ──────────────────────────────────────────────
    // Constructor
    // ──────────────────────────────────────────────

    function test_constructor() public view {
        assertEq(nft.name(), "AWP Subnet");
        assertEq(nft.symbol(), "SUBNET");
        assertEq(nft.awpRegistry(), awpRegistry);
    }

    // ──────────────────────────────────────────────
    // mint
    // ──────────────────────────────────────────────

    function test_mint_success() public {
        vm.prank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");

        assertEq(nft.ownerOf(1), alice);
        assertEq(nft.balanceOf(alice), 1);
    }

    function test_mint_multipleTokens() public {
        vm.startPrank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");
        nft.mint(alice, 2, "Test", address(0x1), address(0x2), 0, "");
        nft.mint(bob, 3, "Test", address(0x1), address(0x2), 0, "");
        vm.stopPrank();

        assertEq(nft.balanceOf(alice), 2);
        assertEq(nft.balanceOf(bob), 1);
        assertEq(nft.ownerOf(2), alice);
        assertEq(nft.ownerOf(3), bob);
    }

    function test_mint_onlyAWPRegistry() public {
        vm.prank(alice);
        vm.expectRevert(SubnetNFT.NotAWPRegistry.selector);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");
    }

    function test_mint_duplicateTokenId_reverts() public {
        vm.startPrank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");

        vm.expectRevert(abi.encodeWithSelector(IERC721Errors.ERC721InvalidSender.selector, address(0)));
        nft.mint(bob, 1, "Test", address(0x1), address(0x2), 0, "");
        vm.stopPrank();
    }

    // ──────────────────────────────────────────────
    // burn
    // ──────────────────────────────────────────────

    function test_burn_success() public {
        vm.startPrank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");
        nft.burn(1);
        vm.stopPrank();

        assertEq(nft.balanceOf(alice), 0);

        // ownerOf should revert since token no longer exists
        vm.expectRevert(abi.encodeWithSelector(IERC721Errors.ERC721NonexistentToken.selector, uint256(1)));
        nft.ownerOf(1);
    }

    function test_burn_nonexistentToken_reverts() public {
        vm.prank(awpRegistry);
        vm.expectRevert(abi.encodeWithSelector(IERC721Errors.ERC721NonexistentToken.selector, uint256(999)));
        nft.burn(999);
    }

    function test_burn_onlyAWPRegistry() public {
        vm.prank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");

        vm.prank(alice);
        vm.expectRevert(SubnetNFT.NotAWPRegistry.selector);
        nft.burn(1);
    }

    // ──────────────────────────────────────────────
    // setBaseURI
    // ──────────────────────────────────────────────

    function test_setBaseURI_success() public {
        vm.startPrank(awpRegistry);
        nft.mint(alice, 42, "Test", address(0x1), address(0x2), 0, "");
        nft.setBaseURI("https://api.cortexia.io/subnet/");
        vm.stopPrank();

        assertEq(nft.tokenURI(42), "https://api.cortexia.io/subnet/42");
    }

    function test_setBaseURI_updateURI() public {
        vm.startPrank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");
        nft.setBaseURI("https://old.example.com/");
        assertEq(nft.tokenURI(1), "https://old.example.com/1");

        nft.setBaseURI("https://new.example.com/meta/");
        vm.stopPrank();

        assertEq(nft.tokenURI(1), "https://new.example.com/meta/1");
    }

    function test_setBaseURI_onlyAWPRegistry() public {
        vm.prank(alice);
        vm.expectRevert(SubnetNFT.NotAWPRegistry.selector);
        nft.setBaseURI("https://evil.com/");
    }

    // ──────────────────────────────────────────────
    // tokenURI
    // ──────────────────────────────────────────────

    function test_tokenURI_emptyBaseURI() public {
        vm.prank(awpRegistry);
        nft.mint(alice, 7, "Test", address(0x1), address(0x2), 0, "");

        // when baseURI is empty, tokenURI returns empty + tokenId.toString()
        assertEq(nft.tokenURI(7), "7");
    }

    function test_tokenURI_withBaseURI() public {
        vm.startPrank(awpRegistry);
        nft.setBaseURI("ipfs://QmHash/");
        nft.mint(alice, 100, "Test", address(0x1), address(0x2), 0, "");
        vm.stopPrank();

        assertEq(nft.tokenURI(100), "ipfs://QmHash/100");
    }

    function test_tokenURI_nonexistentToken_reverts() public {
        vm.expectRevert(abi.encodeWithSelector(IERC721Errors.ERC721NonexistentToken.selector, uint256(1)));
        nft.tokenURI(1);
    }

    function test_tokenURI_afterBurn_reverts() public {
        vm.startPrank(awpRegistry);
        nft.mint(alice, 5, "Test", address(0x1), address(0x2), 0, "");
        nft.burn(5);
        vm.stopPrank();

        vm.expectRevert(abi.encodeWithSelector(IERC721Errors.ERC721NonexistentToken.selector, uint256(5)));
        nft.tokenURI(5);
    }

    // ──────────────────────────────────────────────
    // ERC721 standard behavior
    // ──────────────────────────────────────────────

    function test_transferFrom() public {
        vm.prank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");

        vm.prank(alice);
        nft.transferFrom(alice, bob, 1);

        assertEq(nft.ownerOf(1), bob);
    }

    function test_approve_and_transferFrom() public {
        vm.prank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");

        vm.prank(alice);
        nft.approve(bob, 1);

        vm.prank(bob);
        nft.transferFrom(alice, bob, 1);

        assertEq(nft.ownerOf(1), bob);
    }

    // ──────────────────────────────────────────────
    // setSkillsURI
    // ──────────────────────────────────────────────

    function test_setSkillsURI_success() public {
        vm.prank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");

        vm.prank(alice);
        nft.setSkillsURI(1, "https://example.com/skills");

        SubnetNFT.SubnetData memory data = nft.getSubnetData(1);
        assertEq(data.skillsURI, "https://example.com/skills");
    }

    function test_setSkillsURI_onlyOwner() public {
        vm.prank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");

        vm.prank(bob);
        vm.expectRevert(SubnetNFT.NotTokenOwner.selector);
        nft.setSkillsURI(1, "https://evil.com");
    }

    // ──────────────────────────────────────────────
    // setMinStake
    // ──────────────────────────────────────────────

    function test_setMinStake_success() public {
        vm.prank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");

        vm.prank(alice);
        nft.setMinStake(1, 1000);

        assertEq(nft.getMinStake(1), 1000);
    }

    function test_setMinStake_onlyOwner() public {
        vm.prank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 0, "");

        vm.prank(bob);
        vm.expectRevert(SubnetNFT.NotTokenOwner.selector);
        nft.setMinStake(1, 1000);
    }

    // ──────────────────────────────────────────────
    // mint with skillsURI and minStake
    // ──────────────────────────────────────────────

    function test_mint_withSkillsURI() public {
        vm.prank(awpRegistry);
        nft.mint(alice, 1, "Test", address(0x1), address(0x2), 500, "https://skills.io");

        SubnetNFT.SubnetData memory data = nft.getSubnetData(1);
        assertEq(data.skillsURI, "https://skills.io");
        assertEq(data.minStake, 500);
    }
}
