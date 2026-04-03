// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";

/// @title IWorknetNFT — Worknet NFT interface with on-chain metadata
interface IWorknetNFT is IERC721 {
    struct WorknetData {
        string name;
        address worknetManager;
        address alphaToken;
        string skillsURI;
        uint128 minStake;
        address owner;
    }

    function mint(address to, uint256 tokenId, string calldata name_, address worknetManager_, address alphaToken_, uint128 minStake_, string calldata skillsURI_) external;
    function burn(uint256 tokenId) external;
    function setBaseURI(string memory uri) external;
    function getWorknetManager(uint256 tokenId) external view returns (address);
    function getAlphaToken(uint256 tokenId) external view returns (address);
    function getMinStake(uint256 tokenId) external view returns (uint128);
    function getWorknetData(uint256 tokenId) external view returns (WorknetData memory);
    function setSkillsURI(uint256 tokenId, string calldata skillsURI_) external;
    function setMinStake(uint256 tokenId, uint128 minStake_) external;
}
