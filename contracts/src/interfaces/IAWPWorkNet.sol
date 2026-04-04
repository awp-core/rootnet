// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";

/// @title IAWPWorkNet — UUPS upgradeable worknet NFT interface
interface IAWPWorkNet is IERC721 {
    struct WorknetIdentity {
        string name;
        string symbol;
        address worknetManager;
        address worknetToken;
        bytes32 lpPool;
    }

    struct WorknetMeta {
        string skillsURI;
        uint128 minStake;
        string imageURI;
        string metadataURI;
    }

    struct WorknetData {
        string name;
        string symbol;
        address worknetManager;
        address worknetToken;
        bytes32 lpPool;
        string skillsURI;
        uint128 minStake;
        string imageURI;
        address owner;
    }

    // ── AWPRegistry-only writes ──
    function mint(
        address to, uint256 tokenId,
        string calldata name_, string calldata symbol_,
        address worknetManager_, address worknetToken_, bytes32 lpPool_,
        uint128 minStake_, string calldata skillsURI_
    ) external;
    function setBaseURI(string memory uri) external;

    // ── View ──
    function getWorknetData(uint256 tokenId) external view returns (WorknetData memory);
    function getWorknetIdentity(uint256 tokenId) external view returns (WorknetIdentity memory);
    function getWorknetMeta(uint256 tokenId) external view returns (WorknetMeta memory);
    function getWorknetManager(uint256 tokenId) external view returns (address);
    function getWorknetToken(uint256 tokenId) external view returns (address);
    function getLPPool(uint256 tokenId) external view returns (bytes32);
    function getMinStake(uint256 tokenId) external view returns (uint128);
}
