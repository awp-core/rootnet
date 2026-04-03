// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {Base64} from "@openzeppelin/contracts/utils/Base64.sol";

/// @title WorknetNFT — Worknet NFT with on-chain metadata
/// @notice Each worknet = one NFT. Carries immutable identity (name, worknetManager, alphaToken)
///         plus owner-updatable fields (skillsURI, minStake). Status/lifecycle managed by AWPRegistry.
/// @dev tokenId = worknetId, assigned by AWPRegistry as (block.chainid << 64) | _nextLocalId++.
contract WorknetNFT is ERC721 {
    using Strings for uint256;

    /// @notice Immutable worknet identity fields (set once at mint, never changed)
    struct WorknetIdentity {
        string name;              // Worknet / Alpha token name
        address worknetManager;   // Worknet contract address (holds Alpha minting rights)
        address alphaToken;       // Alpha token address
    }

    /// @notice Owner-updatable worknet metadata
    struct WorknetMeta {
        string skillsURI;         // Skills file URI (OpenClaw skill discovery)
        uint128 minStake;         // Minimum stake requirement for agents (0 = no minimum)
        string metadataURI;       // Custom metadata JSON URI (overrides on-chain generation)
    }

    /// @notice Full worknet data stored in NFT (returned by getWorknetData)
    struct WorknetData {
        string name;
        address worknetManager;
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
    mapping(uint256 => WorknetIdentity) private _identity;

    /// @dev tokenId => owner-updatable metadata
    mapping(uint256 => WorknetMeta) private _meta;

    error NotAWPRegistry();
    error NotTokenOwner();
    error JsonUnsafeCharacter();

    /// @dev Reject strings containing JSON-unsafe characters (", \, control chars 0x00-0x1F)
    function _rejectJsonUnsafe(bytes memory b) internal pure {
        for (uint256 i = 0; i < b.length;) {
            bytes1 c = b[i];
            if (c == 0x22 || c == 0x5c || c < 0x20) revert JsonUnsafeCharacter();
            unchecked { ++i; }
        }
    }
    error TokenNotExist();

    event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI);
    event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake);
    event MetadataURIUpdated(uint256 indexed tokenId, string metadataURI);

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

    /// @notice Mint worknet NFT with identity data + initial metadata (called by AWPRegistry during registerWorknet)
    function mint(
        address to,
        uint256 tokenId,
        string calldata name_,
        address worknetManager_,
        address alphaToken_,
        uint128 minStake_,
        string calldata skillsURI_
    ) external onlyAWPRegistry {
        _mint(to, tokenId);
        _identity[tokenId] = WorknetIdentity({
            name: name_,
            worknetManager: worknetManager_,
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

    /// @notice Burn worknet NFT (called by AWPRegistry on deregister)
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
        _rejectJsonUnsafe(bytes(skillsURI_));
        _meta[tokenId].skillsURI = skillsURI_;
        emit SkillsURIUpdated(tokenId, skillsURI_);
    }

    /// @notice Update minimum stake requirement (only NFT owner)
    function setMinStake(uint256 tokenId, uint128 minStake_) external {
        if (ownerOf(tokenId) != msg.sender) revert NotTokenOwner();
        _meta[tokenId].minStake = minStake_;
        emit MinStakeUpdated(tokenId, minStake_);
    }

    /// @notice Set custom metadata URI for this worknet NFT (only NFT owner)
    /// @dev When set, tokenURI returns this URI instead of on-chain generated JSON.
    ///      Set to empty string to revert to on-chain metadata.
    function setMetadataURI(uint256 tokenId, string calldata metadataURI_) external {
        if (ownerOf(tokenId) != msg.sender) revert NotTokenOwner();
        _meta[tokenId].metadataURI = metadataURI_;
        emit MetadataURIUpdated(tokenId, metadataURI_);
    }

    // ═══════════════════════════════════════════════
    // ── View functions ──
    // ═══════════════════════════════════════════════

    /// @notice Get full worknet data from NFT
    function getWorknetData(uint256 tokenId) external view returns (WorknetData memory) {
        address owner = _ownerOf(tokenId);
        if (owner == address(0)) revert TokenNotExist();
        WorknetIdentity storage id_ = _identity[tokenId];
        WorknetMeta storage meta_ = _meta[tokenId];
        return WorknetData({
            name: id_.name,
            worknetManager: id_.worknetManager,
            alphaToken: id_.alphaToken,
            skillsURI: meta_.skillsURI,
            minStake: meta_.minStake,
            owner: owner
        });
    }

    /// @notice Get worknet contract address
    function getWorknetManager(uint256 tokenId) external view returns (address) {
        return _identity[tokenId].worknetManager;
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

        // Priority 1: per-token custom metadata URI (set by worknet owner)
        if (bytes(_meta[tokenId].metadataURI).length > 0) {
            return _meta[tokenId].metadataURI;
        }

        // Priority 2: global baseURI (set by AWPRegistry governance)
        if (bytes(_baseTokenURI).length > 0) {
            return string.concat(_baseTokenURI, tokenId.toString());
        }

        // Priority 3: on-chain generated JSON metadata
        WorknetIdentity storage id = _identity[tokenId];
        WorknetMeta storage meta = _meta[tokenId];

        string memory json = string.concat(
            '{"name":"', id.name,
            '","description":"AWP Worknet #', tokenId.toString(),
            '","external_url":"', meta.skillsURI,
            '","attributes":[',
                '{"trait_type":"Worknet Manager","value":"', Strings.toHexString(id.worknetManager),
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
