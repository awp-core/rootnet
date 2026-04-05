// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {Base64} from "@openzeppelin/contracts/utils/Base64.sol";

/// @title AWPWorkNet — UUPS upgradeable worknet ownership NFT with on-chain metadata
/// @notice Each worknet = one NFT. Minted at activation (not registration).
///         Identity fields are immutable after mint. Owner or delegate can update metadata fields.
///         For full customization (description, external_url, etc.) set metadataURI to off-chain JSON.
///         Implements ERC-4906 for metadata update notifications to marketplaces/indexers.
/// @dev tokenId = worknetId = (block.chainid << 64) | localCounter, assigned by AWPRegistry.
contract AWPWorkNet is UUPSUpgradeable, ERC721Upgradeable {
    using Strings for uint256;

    // ══════════════════════════════════════════════
    //  Data structures
    // ══════════════════════════════════════════════

    /// @notice Worknet identity (set once at mint, never changed)
    struct WorknetIdentity {
        string name;              // Worknet / Alpha token name
        string symbol;            // Alpha token symbol
        address worknetManager;   // WorknetManager proxy address
        address worknetToken;       // Worknet token address
        bytes32 lpPool;           // LP pool ID (AWP/Alpha pair)
    }

    /// @notice Owner-updatable worknet metadata
    struct WorknetMeta {
        string skillsURI;         // Agent skill discovery URI (read by agents, protocol core)
        uint128 minStake;         // Minimum stake hint for agents (not enforced on-chain)
        string imageURI;          // NFT image URI (displayed in wallets/marketplaces)
        string metadataURI;       // Full NFT metadata override (bypasses on-chain JSON when set)
    }

    /// @notice Full worknet data (returned by getWorknetData)
    struct WorknetData {
        string name;
        string symbol;
        address worknetManager;
        address worknetToken;
        bytes32 lpPool;
        string skillsURI;
        uint128 minStake;
        string imageURI;
        string metadataURI;
        address owner;
    }

    // ══════════════════════════════════════════════
    //  Constants
    // ══════════════════════════════════════════════

    uint256 public constant MAX_URI_LENGTH = 512;
    uint256 public constant MAX_SKILLS_URI_LENGTH = 2048;

    /// @dev Pre-computed selector for AWPRegistry.delegates(address,address) — avoids runtime keccak
    bytes4 private constant DELEGATES_SELECTOR = bytes4(keccak256("delegates(address,address)"));

    // ══════════════════════════════════════════════
    //  Immutables
    // ══════════════════════════════════════════════

    /// @notice AWPRegistry contract address
    address public immutable awpRegistry;

    // ══════════════════════════════════════════════
    //  Storage
    // ══════════════════════════════════════════════

    /// @notice Guardian address — can upgrade contract and set collection-level config
    address public guardian;

    /// @dev NFT metadata base URI
    string private _baseTokenURI;

    /// @dev Collection-level metadata URI (ERC-7572 contractURI)
    string private _contractURI;

    /// @dev tokenId => identity (immutable after mint)
    mapping(uint256 => WorknetIdentity) private _identity;

    /// @dev tokenId => owner-updatable metadata
    mapping(uint256 => WorknetMeta) private _meta;

    /// @dev Reserved storage gap for upgrades
    uint256[43] private __gap;

    // ══════════════════════════════════════════════
    //  Errors & Events
    // ══════════════════════════════════════════════

    error NotAWPRegistry();
    error NotGuardian();
    error NotAuthorized();
    error TokenNotExist();
    error JsonUnsafeCharacter();
    error StringTooLong();
    error ZeroAddress();

    /// @dev ERC-4906: emitted when token metadata changes (marketplaces/indexers refresh on this)
    event MetadataUpdate(uint256 _tokenId);

    event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI);
    event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake);
    event ImageURIUpdated(uint256 indexed tokenId, string imageURI);
    event MetadataURIUpdated(uint256 indexed tokenId, string metadataURI);
    event GuardianUpdated(address indexed newGuardian);
    event ContractURIUpdated(string uri);

    // ══════════════════════════════════════════════
    //  Modifiers
    // ══════════════════════════════════════════════

    modifier onlyAWPRegistry() {
        if (msg.sender != awpRegistry) revert NotAWPRegistry();
        _;
    }

    modifier onlyGuardian() {
        if (msg.sender != guardian) revert NotGuardian();
        _;
    }

    modifier tokenExists(uint256 tokenId) {
        if (_ownerOf(tokenId) == address(0)) revert TokenNotExist();
        _;
    }

    /// @dev Only the NFT owner (for destructive actions like burn)
    modifier onlyTokenOwner(uint256 tokenId) {
        if (ownerOf(tokenId) != msg.sender) revert NotAuthorized();
        _;
    }

    /// @dev Owner or delegate authorized via AWPRegistry.delegates(owner, msg.sender)
    modifier onlyTokenOwnerOrDelegate(uint256 tokenId) {
        address owner_ = ownerOf(tokenId);
        if (msg.sender != owner_) {
            (bool ok, bytes memory data) = awpRegistry.staticcall(
                abi.encodeWithSelector(DELEGATES_SELECTOR, owner_, msg.sender)
            );
            if (!ok || !abi.decode(data, (bool))) revert NotAuthorized();
        }
        _;
    }

    // ══════════════════════════════════════════════
    //  Constructor + Initialization
    // ══════════════════════════════════════════════

    /// @param awpRegistry_ AWPRegistry contract address (baked into bytecode)
    constructor(address awpRegistry_) {
        awpRegistry = awpRegistry_;
        _disableInitializers();
    }

    /// @notice Initialize the NFT collection (called once via proxy)
    function initialize(string memory name_, string memory symbol_, address guardian_) external initializer {
        __ERC721_init(name_, symbol_);
        if (guardian_ == address(0)) revert ZeroAddress();
        guardian = guardian_;
    }

    /// @dev UUPS upgrade authorization — uses local guardian storage (no external dependency)
    function _authorizeUpgrade(address) internal view override {
        if (msg.sender != guardian) revert NotGuardian();
    }

    // ══════════════════════════════════════════════
    //  AWPRegistry-only writes
    // ══════════════════════════════════════════════

    /// @notice Mint worknet NFT with complete identity (called by AWPRegistry during activateWorknet)
    function mint(
        address to,
        uint256 tokenId,
        string calldata name_,
        string calldata symbol_,
        address worknetManager_,
        address worknetToken_,
        bytes32 lpPool_,
        uint128 minStake_,
        string calldata skillsURI_
    ) external onlyAWPRegistry {
        _mint(to, tokenId);
        _identity[tokenId] = WorknetIdentity({
            name: name_,
            symbol: symbol_,
            worknetManager: worknetManager_,
            worknetToken: worknetToken_,
            lpPool: lpPool_
        });
        if (minStake_ > 0) {
            _meta[tokenId].minStake = minStake_;
        }
        if (bytes(skillsURI_).length > 0) {
            _meta[tokenId].skillsURI = skillsURI_;
        }
    }

    /// @notice Burn worknet NFT (only NFT owner — voluntary renouncement)
    function burn(uint256 tokenId) external onlyTokenOwner(tokenId) {
        _burn(tokenId);
        delete _identity[tokenId];
        delete _meta[tokenId];
    }

    // ══════════════════════════════════════════════
    //  Guardian management
    // ══════════════════════════════════════════════

    /// @notice Update guardian address (self-sovereign)
    function setGuardian(address g) external onlyGuardian {
        if (g == address(0)) revert ZeroAddress();
        guardian = g;
        emit GuardianUpdated(g);
    }

    /// @notice Set NFT metadata base URI (Guardian only)
    function setBaseURI(string calldata uri) external onlyGuardian {
        _baseTokenURI = uri;
    }

    /// @notice Set collection-level metadata URI (ERC-7572)
    function setContractURI(string calldata uri) external onlyGuardian {
        _contractURI = uri;
        emit ContractURIUpdated(uri);
    }

    /// @notice Collection-level metadata (ERC-7572 — used by marketplaces for collection info)
    function contractURI() external view returns (string memory) {
        return _contractURI;
    }

    // ══════════════════════════════════════════════
    //  Owner / Delegate updatable fields
    // ══════════════════════════════════════════════

    /// @notice Update skills URI — agent skill discovery endpoint
    function setSkillsURI(uint256 tokenId, string calldata v) external onlyTokenOwnerOrDelegate(tokenId) {
        if (bytes(v).length > MAX_SKILLS_URI_LENGTH) revert StringTooLong();
        _rejectJsonUnsafe(bytes(v));
        _meta[tokenId].skillsURI = v;
        emit SkillsURIUpdated(tokenId, v);
        emit MetadataUpdate(tokenId);
    }

    /// @notice Update minimum stake hint (not enforced on-chain)
    function setMinStake(uint256 tokenId, uint128 v) external onlyTokenOwnerOrDelegate(tokenId) {
        _meta[tokenId].minStake = v;
        emit MinStakeUpdated(tokenId, v);
        emit MetadataUpdate(tokenId);
    }

    /// @notice Update NFT image URI (displayed in wallets/marketplaces)
    function setImageURI(uint256 tokenId, string calldata v) external onlyTokenOwnerOrDelegate(tokenId) {
        if (bytes(v).length > MAX_URI_LENGTH) revert StringTooLong();
        _rejectJsonUnsafe(bytes(v));
        _meta[tokenId].imageURI = v;
        emit ImageURIUpdated(tokenId, v);
        emit MetadataUpdate(tokenId);
    }

    /// @notice Set custom metadata URI — overrides on-chain JSON entirely. Empty string reverts to on-chain.
    ///         Use this for full customization: description, external_url, animation_url, etc.
    function setMetadataURI(uint256 tokenId, string calldata v) external onlyTokenOwnerOrDelegate(tokenId) {
        if (bytes(v).length > MAX_URI_LENGTH) revert StringTooLong();
        _meta[tokenId].metadataURI = v;
        emit MetadataURIUpdated(tokenId, v);
        emit MetadataUpdate(tokenId);
    }

    // ══════════════════════════════════════════════
    //  View functions
    // ══════════════════════════════════════════════

    /// @notice Get full worknet data (identity + meta + owner)
    function getWorknetData(uint256 tokenId) external view tokenExists(tokenId) returns (WorknetData memory) {
        WorknetIdentity storage id_ = _identity[tokenId];
        WorknetMeta storage meta_ = _meta[tokenId];
        return WorknetData({
            name: id_.name,
            symbol: id_.symbol,
            worknetManager: id_.worknetManager,
            worknetToken: id_.worknetToken,
            lpPool: id_.lpPool,
            skillsURI: meta_.skillsURI,
            minStake: meta_.minStake,
            imageURI: meta_.imageURI,
            metadataURI: meta_.metadataURI,
            owner: _ownerOf(tokenId)
        });
    }

    /// @notice Get immutable identity only
    function getWorknetIdentity(uint256 tokenId) external view tokenExists(tokenId) returns (WorknetIdentity memory) {
        return _identity[tokenId];
    }

    /// @notice Get owner-updatable metadata only
    function getWorknetMeta(uint256 tokenId) external view tokenExists(tokenId) returns (WorknetMeta memory) {
        return _meta[tokenId];
    }

    /// @notice Get worknet manager proxy address
    function getWorknetManager(uint256 tokenId) external view tokenExists(tokenId) returns (address) {
        return _identity[tokenId].worknetManager;
    }

    /// @notice Get alpha token address
    function getWorknetToken(uint256 tokenId) external view tokenExists(tokenId) returns (address) {
        return _identity[tokenId].worknetToken;
    }

    /// @notice Get LP pool ID
    function getLPPool(uint256 tokenId) external view tokenExists(tokenId) returns (bytes32) {
        return _identity[tokenId].lpPool;
    }

    /// @notice Get minimum stake hint
    function getMinStake(uint256 tokenId) external view tokenExists(tokenId) returns (uint128) {
        return _meta[tokenId].minStake;
    }

    /// @notice Token URI — 3-tier priority: metadataURI → baseURI/{tokenId} → on-chain JSON
    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        _requireOwned(tokenId);

        WorknetMeta storage meta = _meta[tokenId];

        // Priority 1: per-token custom metadata URI (full override)
        if (bytes(meta.metadataURI).length > 0) {
            return meta.metadataURI;
        }

        // Priority 2: global baseURI + tokenId
        if (bytes(_baseTokenURI).length > 0) {
            return string.concat(_baseTokenURI, tokenId.toString());
        }

        // Priority 3: on-chain generated JSON
        return _buildOnChainURI(tokenId, meta);
    }

    // ══════════════════════════════════════════════
    //  Internal
    // ══════════════════════════════════════════════

    /// @dev ERC-165: declare support for ERC-4906 (Metadata Update Extension)
    function supportsInterface(bytes4 interfaceId) public view override returns (bool) {
        return interfaceId == bytes4(0x49064906) || super.supportsInterface(interfaceId);
    }

    /// @dev Build on-chain Base64 JSON metadata (ERC721 / OpenSea standard)
    function _buildOnChainURI(uint256 tokenId, WorknetMeta storage meta) internal view returns (string memory) {
        WorknetIdentity storage id = _identity[tokenId];

        string memory json = string.concat(
            '{"name":"', id.name,
            '","description":"AWP Worknet #', tokenId.toString(), '"'
        );

        if (bytes(meta.imageURI).length > 0) {
            json = string.concat(json, ',"image":"', meta.imageURI, '"');
        }

        json = string.concat(json,
            ',"attributes":[',
                '{"trait_type":"Symbol","value":"', id.symbol,
                '"},{"trait_type":"Worknet Manager","value":"', Strings.toHexString(id.worknetManager),
                '"},{"trait_type":"Worknet Token","value":"', Strings.toHexString(id.worknetToken),
                '"},{"trait_type":"LP Pool","value":"', uint256(id.lpPool).toHexString(32),
                '"},{"trait_type":"Min Stake","value":"', uint256(meta.minStake).toString(),
                '"},{"trait_type":"Skills URI","value":"', meta.skillsURI,
                '"},{"trait_type":"Chain ID","value":"', (tokenId >> 64).toString(),
                '"},{"trait_type":"Local ID","value":"', (tokenId & ((1 << 64) - 1)).toString(),
            '"}]}'
        );

        return string.concat("data:application/json;base64,", Base64.encode(bytes(json)));
    }

    /// @dev Reject strings containing JSON-unsafe characters (", \, control chars 0x00-0x1F per RFC 8259)
    function _rejectJsonUnsafe(bytes memory b) internal pure {
        for (uint256 i = 0; i < b.length;) {
            bytes1 c = b[i];
            if (c == 0x22 || c == 0x5c || c < 0x20) revert JsonUnsafeCharacter();
            unchecked { ++i; }
        }
    }
}
