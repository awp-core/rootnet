// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {Base64} from "@openzeppelin/contracts/utils/Base64.sol";

/// @title SubnetNFT — Subnet NFT with on-chain metadata
/// @notice Each subnet = one NFT. Carries immutable identity (name, subnetManager, alphaToken)
///         plus owner-updatable fields (skillsURI, minStake). Status/lifecycle managed by AWPRegistry.
/// @dev tokenId = subnetId, assigned by AWPRegistry as (block.chainid << 64) | _nextLocalId++.
contract SubnetNFT is ERC721 {
    using Strings for uint256;

    /// @notice Immutable subnet identity fields (set once at mint, never changed)
    struct SubnetIdentity {
        string name;              // Subnet / Alpha token name
        address subnetManager;   // Subnet contract address (holds Alpha minting rights)
        address alphaToken;       // Alpha token address
    }

    /// @notice Owner-updatable subnet metadata
    struct SubnetMeta {
        string skillsURI;         // Skills file URI (OpenClaw skill discovery)
        uint128 minStake;         // Minimum stake requirement for agents (0 = no minimum)
    }

    /// @notice Full subnet data stored in NFT (returned by getSubnetData)
    struct SubnetData {
        string name;
        address subnetManager;
        address alphaToken;
        string skillsURI;
        uint128 minStake;
        address owner;
    }

    /// @notice AWPRegistry contract address (immutable)
    address public immutable awpRegistry;

    /// @dev NFT metadata base URI
    string private _baseTokenURI;

    /// @dev tokenId => immutable identity
    mapping(uint256 => SubnetIdentity) private _identity;

    /// @dev tokenId => owner-updatable metadata
    mapping(uint256 => SubnetMeta) private _meta;

    error NotAWPRegistry();
    error NotTokenOwner();
    error TokenNotExist();

    event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI);
    event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake);

    modifier onlyAWPRegistry() {
        if (msg.sender != awpRegistry) revert NotAWPRegistry();
        _;
    }

    constructor(string memory name_, string memory symbol_, address awpRegistry_) ERC721(name_, symbol_) {
        awpRegistry = awpRegistry_;
    }

    // ═══════════════════════════════════════════════
    // ── AWPRegistry-only writes ──
    // ═══════════════════════════════════════════════

    /// @notice Mint subnet NFT with identity data + initial metadata (called by AWPRegistry during registerSubnet)
    function mint(
        address to,
        uint256 tokenId,
        string calldata name_,
        address subnetManager_,
        address alphaToken_,
        uint128 minStake_,
        string calldata skillsURI_
    ) external onlyAWPRegistry {
        _mint(to, tokenId);
        _identity[tokenId] = SubnetIdentity({
            name: name_,
            subnetManager: subnetManager_,
            alphaToken: alphaToken_
        });
        if (minStake_ > 0) {
            _meta[tokenId].minStake = minStake_;
            emit MinStakeUpdated(tokenId, minStake_);
        }
        if (bytes(skillsURI_).length > 0) {
            _meta[tokenId].skillsURI = skillsURI_;
            emit SkillsURIUpdated(tokenId, skillsURI_);
        }
    }

    /// @notice Burn subnet NFT (called by AWPRegistry on deregister)
    function burn(uint256 tokenId) external onlyAWPRegistry {
        _burn(tokenId);
        delete _identity[tokenId];
        delete _meta[tokenId];
    }

    /// @notice Set the base URI for NFT metadata
    function setBaseURI(string memory uri) external onlyAWPRegistry {
        _baseTokenURI = uri;
    }

    // ═══════════════════════════════════════════════
    // ── Owner-updatable fields ──
    // ═══════════════════════════════════════════════

    /// @notice Update skills URI (only NFT owner)
    function setSkillsURI(uint256 tokenId, string calldata skillsURI_) external {
        if (ownerOf(tokenId) != msg.sender) revert NotTokenOwner();
        _meta[tokenId].skillsURI = skillsURI_;
        emit SkillsURIUpdated(tokenId, skillsURI_);
    }

    /// @notice Update minimum stake requirement (only NFT owner)
    function setMinStake(uint256 tokenId, uint128 minStake_) external {
        if (ownerOf(tokenId) != msg.sender) revert NotTokenOwner();
        _meta[tokenId].minStake = minStake_;
        emit MinStakeUpdated(tokenId, minStake_);
    }

    // ═══════════════════════════════════════════════
    // ── View functions ──
    // ═══════════════════════════════════════════════

    /// @notice Get full subnet data from NFT
    function getSubnetData(uint256 tokenId) external view returns (SubnetData memory) {
        address owner = _ownerOf(tokenId);
        if (owner == address(0)) revert TokenNotExist();
        SubnetIdentity storage id_ = _identity[tokenId];
        SubnetMeta storage meta_ = _meta[tokenId];
        return SubnetData({
            name: id_.name,
            subnetManager: id_.subnetManager,
            alphaToken: id_.alphaToken,
            skillsURI: meta_.skillsURI,
            minStake: meta_.minStake,
            owner: owner
        });
    }

    /// @notice Get subnet contract address
    function getSubnetManager(uint256 tokenId) external view returns (address) {
        return _identity[tokenId].subnetManager;
    }

    /// @notice Get alpha token address
    function getAlphaToken(uint256 tokenId) external view returns (address) {
        return _identity[tokenId].alphaToken;
    }

    /// @notice Get minimum stake requirement
    function getMinStake(uint256 tokenId) external view returns (uint128) {
        return _meta[tokenId].minStake;
    }

    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        _requireOwned(tokenId);

        // If baseURI is set, use it (external metadata server)
        if (bytes(_baseTokenURI).length > 0) {
            return string.concat(_baseTokenURI, tokenId.toString());
        }

        // Otherwise, generate on-chain JSON metadata
        SubnetIdentity storage id = _identity[tokenId];
        SubnetMeta storage meta = _meta[tokenId];

        string memory json = string.concat(
            '{"name":"', id.name,
            '","description":"AWP Subnet #', tokenId.toString(),
            '","external_url":"', meta.skillsURI,
            '","attributes":[',
                '{"trait_type":"Subnet Manager","value":"', Strings.toHexString(id.subnetManager),
                '"},{"trait_type":"Alpha Token","value":"', Strings.toHexString(id.alphaToken),
                '"},{"trait_type":"Min Stake","value":"', uint256(meta.minStake).toString(),
                '"},{"trait_type":"Chain ID","value":"', (tokenId >> 64).toString(),
                '"},{"trait_type":"Local ID","value":"', (tokenId & ((1 << 64) - 1)).toString(),
            '"}]}'
        );

        return string.concat("data:application/json;base64,", Base64.encode(bytes(json)));
    }

    function _baseURI() internal view override returns (string memory) {
        return _baseTokenURI;
    }
}
