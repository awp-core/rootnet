// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {DeployHelper} from "./helpers/DeployHelper.sol";
import {AWPWorkNet} from "../src/core/AWPWorkNet.sol";
import {IAWPWorkNet} from "../src/interfaces/IAWPWorkNet.sol";

contract AWPWorkNetTest is DeployHelper {
    uint256 worknetId;

    function setUp() public {
        _deployAll();
        worknetId = _registerWorknet(alice);
        _activateWorknet(worknetId);
    }

    function test_ownerIsAlice() public view {
        assertEq(awpWorkNet.ownerOf(worknetId), alice);
    }

    function test_getWorknetData() public view {
        AWPWorkNet.WorknetData memory data = awpWorkNet.getWorknetData(worknetId);
        assertEq(data.owner, alice);
        assertTrue(data.worknetManager != address(0));
        assertTrue(data.alphaToken != address(0));
        assertTrue(data.lpPool != bytes32(0));
    }

    function test_getWorknetIdentity() public view {
        AWPWorkNet.WorknetIdentity memory id = awpWorkNet.getWorknetIdentity(worknetId);
        assertEq(id.name, "TestWorknet");
        assertEq(id.symbol, "TWN");
    }

    function test_setSkillsURI() public {
        vm.prank(alice);
        awpWorkNet.setSkillsURI(worknetId, "https://example.com/skills");

        AWPWorkNet.WorknetMeta memory meta = awpWorkNet.getWorknetMeta(worknetId);
        assertEq(meta.skillsURI, "https://example.com/skills");
    }

    function test_setSkillsURI_delegate() public {
        vm.prank(alice);
        awpRegistry.grantDelegate(bob);

        vm.prank(bob);
        awpWorkNet.setSkillsURI(worknetId, "https://example.com/skills-by-delegate");

        AWPWorkNet.WorknetMeta memory meta = awpWorkNet.getWorknetMeta(worknetId);
        assertEq(meta.skillsURI, "https://example.com/skills-by-delegate");
    }

    function test_setSkillsURI_notAuthorized_reverts() public {
        vm.prank(bob);
        vm.expectRevert(AWPWorkNet.NotAuthorized.selector);
        awpWorkNet.setSkillsURI(worknetId, "hack");
    }

    function test_setMinStake() public {
        vm.prank(alice);
        awpWorkNet.setMinStake(worknetId, 100e18);
        assertEq(awpWorkNet.getMinStake(worknetId), 100e18);
    }

    function test_setImageURI() public {
        vm.prank(alice);
        awpWorkNet.setImageURI(worknetId, "ipfs://Qm...");

        AWPWorkNet.WorknetMeta memory meta = awpWorkNet.getWorknetMeta(worknetId);
        assertEq(meta.imageURI, "ipfs://Qm...");
    }

    function test_setMetadataURI() public {
        vm.prank(alice);
        awpWorkNet.setMetadataURI(worknetId, "https://metadata.example.com/1");

        AWPWorkNet.WorknetMeta memory meta = awpWorkNet.getWorknetMeta(worknetId);
        assertEq(meta.metadataURI, "https://metadata.example.com/1");
    }

    function test_tokenURI_onChainJSON() public view {
        string memory uri = awpWorkNet.tokenURI(worknetId);
        assertTrue(bytes(uri).length > 40);
    }

    function test_tokenURI_metadataURI_overrides() public {
        vm.prank(alice);
        awpWorkNet.setMetadataURI(worknetId, "https://custom.json");

        assertEq(awpWorkNet.tokenURI(worknetId), "https://custom.json");
    }

    function test_burn_byOwner() public {
        vm.prank(alice);
        awpWorkNet.burn(worknetId);

        vm.expectRevert();
        awpWorkNet.ownerOf(worknetId);
    }

    function test_burn_notOwner_reverts() public {
        vm.prank(bob);
        vm.expectRevert(AWPWorkNet.NotAuthorized.selector);
        awpWorkNet.burn(worknetId);
    }

    function test_stringTooLong_reverts() public {
        bytes memory longStr = new bytes(513);
        for (uint i = 0; i < 513; i++) longStr[i] = "a";

        vm.prank(alice);
        vm.expectRevert(AWPWorkNet.StringTooLong.selector);
        awpWorkNet.setImageURI(worknetId, string(longStr));
    }

    function test_jsonUnsafe_reverts() public {
        vm.prank(alice);
        vm.expectRevert(AWPWorkNet.JsonUnsafeCharacter.selector);
        awpWorkNet.setSkillsURI(worknetId, 'bad"quote');
    }

    function test_nonExistentToken_reverts() public {
        vm.expectRevert(AWPWorkNet.TokenNotExist.selector);
        awpWorkNet.getWorknetData(999999);
    }

    // ── Guardian ──

    function test_setGuardian() public {
        vm.prank(guardian);
        awpWorkNet.setGuardian(alice);
        assertEq(awpWorkNet.guardian(), alice);
    }

    function test_setContractURI() public {
        vm.prank(guardian);
        awpWorkNet.setContractURI("https://collection.json");
        assertEq(awpWorkNet.contractURI(), "https://collection.json");
    }

    function test_supportsInterface_ERC4906() public view {
        assertTrue(awpWorkNet.supportsInterface(bytes4(0x49064906)));
    }

    function test_supportsInterface_ERC721() public view {
        assertTrue(awpWorkNet.supportsInterface(bytes4(0x80ac58cd)));
    }
}
