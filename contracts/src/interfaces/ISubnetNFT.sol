// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";

/// @title ISubnetNFT — Subnet NFT interface with on-chain metadata
interface ISubnetNFT is IERC721 {
    struct SubnetData {
        string name;
        address subnetManager;
        address alphaToken;
        string skillsURI;
        uint128 minStake;
        address owner;
    }

    function mint(address to, uint256 tokenId, string calldata name_, address subnetManager_, address alphaToken_, uint128 minStake_) external;
    function burn(uint256 tokenId) external;
    function setBaseURI(string memory uri) external;
    function getSubnetManager(uint256 tokenId) external view returns (address);
    function getAlphaToken(uint256 tokenId) external view returns (address);
    function getMinStake(uint256 tokenId) external view returns (uint128);
    function getSubnetData(uint256 tokenId) external view returns (SubnetData memory);
}
