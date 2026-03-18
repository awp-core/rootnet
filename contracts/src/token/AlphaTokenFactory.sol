// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {AlphaToken} from "./AlphaToken.sol";

/// @title AlphaTokenFactory — deploys AlphaToken via CREATE2 with configurable vanity address validation
/// @notice Uses CREATE2 to deploy an independent AlphaToken instance for each subnet.
///         Vanity rules are configured at deployment: 4 prefix positions + 4 suffix positions.
///         Each position encodes a required hex character with EIP-55 case sensitivity.
/// @dev Vanity rule encoding per position (uint8):
///      0-9:   match digit '0'-'9'
///      10-15: match lowercase hex 'a'-'f' (EIP-55 must NOT uppercase)
///      16-21: match uppercase hex 'A'-'F' (EIP-55 must uppercase)
///      >=22:  wildcard — no check for this position
///      8 positions packed into a single uint64:
///      [prefix0, prefix1, prefix2, prefix3, suffix0, suffix1, suffix2, suffix3]
///      high byte = prefix0, low byte = suffix3
///      Default rule "0xA1????CAFE": 0x1001FFFF12101514
contract AlphaTokenFactory is Ownable {
    /// @notice RootNet contract address (immutable after setAddresses)
    address public rootNet;

    /// @notice Whether configuration is complete (true after setAddresses)
    bool public configured;

    /// @notice Vanity address rule: 8 positions packed into uint64
    /// @dev Byte layout: [prefix0][prefix1][prefix2][prefix3][suffix0][suffix1][suffix2][suffix3]
    uint64 public immutable vanityRule;

    /// @notice Precomputed keccak256(type(AlphaToken).creationCode) for predictDeployAddress
    bytes32 public immutable ALPHA_BYTECODE_HASH;

    /// @dev Caller is not RootNet
    error NotRootNet();

    /// @dev Deployed address does not satisfy vanity constraints
    error InvalidVanityAddress();

    /// @notice Deploy the factory contract
    /// @param deployer_   Initial owner (deployer)
    /// @param vanityRule_ Packed vanity rule (8 × uint8, see encoding above). 0 = no validation.
    constructor(address deployer_, uint64 vanityRule_) Ownable(deployer_) {
        vanityRule = vanityRule_;
        ALPHA_BYTECODE_HASH = keccak256(type(AlphaToken).creationCode);
    }

    /// @notice Set the RootNet address and renounce ownership (can only be called once)
    /// @param rootNet_ RootNet contract address
    function setAddresses(address rootNet_) external onlyOwner {
        rootNet = rootNet_;
        configured = true;
        renounceOwnership();
    }

    /// @notice Deploy a new AlphaToken instance via CREATE2 (only RootNet may call)
    /// @param subnetId_ Subnet ID
    /// @param name_     Token name
    /// @param symbol_   Token symbol
    /// @param admin_    Admin address for the instance (typically RootNet itself)
    /// @param salt_     CREATE2 salt; if bytes32(0), uses subnetId as salt
    /// @return Address of the newly created AlphaToken contract
    function deploy(uint256 subnetId_, string memory name_, string memory symbol_, address admin_, bytes32 salt_)
        external
        returns (address)
    {
        if (msg.sender != rootNet) revert NotRootNet();

        bytes32 effectiveSalt = salt_ == bytes32(0) ? bytes32(subnetId_) : salt_;

        AlphaToken token = new AlphaToken{salt: effectiveSalt}();
        token.initialize(name_, symbol_, subnetId_, admin_);

        // Validate vanity address if rule is configured
        if (vanityRule != 0) {
            _validateVanityAddress(address(token));
        }

        return address(token);
    }

    /// @notice Predict the CREATE2 deployment address for a given salt (for off-chain vanity salt search)
    /// @param salt_ The CREATE2 salt to test
    /// @return The address that would be deployed with this salt
    function predictDeployAddress(bytes32 salt_) external view returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), address(this), salt_, ALPHA_BYTECODE_HASH)
        ))));
    }

    /// @notice Validate that an address satisfies the configured vanity rule with EIP-55 case checking
    /// @param addr The address to validate
    function _validateVanityAddress(address addr) internal view {
        // Convert address to lowercase hex string (40 chars)
        bytes memory hex40 = _toHexString(addr);

        // Compute EIP-55 checksum hash
        bytes32 hash = keccak256(hex40);

        uint64 rule = vanityRule;

        // Check 4 prefix positions (hex40[0..3])
        for (uint256 i = 0; i < 4;) {
            uint8 expected = uint8(rule >> (56 - i * 8));
            if (expected < 22) {
                _checkPosition(hex40, hash, i, expected);
            }
            unchecked { ++i; }
        }

        // Check 4 suffix positions (hex40[36..39])
        for (uint256 i = 0; i < 4;) {
            uint8 expected = uint8(rule >> (24 - i * 8));
            if (expected < 22) {
                _checkPosition(hex40, hash, 36 + i, expected);
            }
            unchecked { ++i; }
        }
    }

    /// @dev Check a single hex position against the expected value with EIP-55 case rules
    /// @param hex40    Lowercase hex string (40 chars)
    /// @param hash     keccak256(hex40) for EIP-55 checksum
    /// @param pos      Position in hex40 (0-39)
    /// @param expected Expected value: 0-9=digit, 10-15=lowercase a-f, 16-21=uppercase A-F
    function _checkPosition(bytes memory hex40, bytes32 hash, uint256 pos, uint8 expected) internal pure {
        bytes1 c = hex40[pos];

        if (expected <= 9) {
            // Expect digit '0'-'9'
            if (c != bytes1(uint8(0x30 + expected))) revert InvalidVanityAddress();
        } else if (expected <= 15) {
            // Expect lowercase letter 'a'-'f' (EIP-55 must NOT uppercase)
            uint8 letter = expected - 10; // 0=a, 1=b, ..., 5=f
            if (c != bytes1(uint8(0x61 + letter))) revert InvalidVanityAddress();
            // EIP-55 check: hash nibble at this position must be < 8 (stays lowercase)
            uint8 hashNibble = _getHashNibble(hash, pos);
            if (hashNibble >= 8) revert InvalidVanityAddress();
        } else {
            // expected 16-21: Expect uppercase letter 'A'-'F' (EIP-55 must uppercase)
            uint8 letter = expected - 16; // 0=A, 1=B, ..., 5=F
            // Raw hex must be the corresponding lowercase letter
            if (c != bytes1(uint8(0x61 + letter))) revert InvalidVanityAddress();
            // EIP-55 check: hash nibble at this position must be >= 8 (uppercased)
            uint8 hashNibble = _getHashNibble(hash, pos);
            if (hashNibble < 8) revert InvalidVanityAddress();
        }
    }

    /// @dev Extract the hash nibble at a given hex position for EIP-55
    function _getHashNibble(bytes32 hash, uint256 pos) internal pure returns (uint8) {
        uint8 hashByte = uint8(hash[pos / 2]);
        return (pos % 2 == 0) ? (hashByte >> 4) : (hashByte & 0x0f);
    }

    /// @notice Convert an address to its lowercase hex string representation (40 characters, no "0x" prefix)
    function _toHexString(address addr) internal pure returns (bytes memory) {
        bytes memory hex40 = new bytes(40);
        uint160 value = uint160(addr);
        bytes16 hexChars = "0123456789abcdef";
        for (uint256 i = 0; i < 20;) {
            uint8 b = uint8(value >> (152 - i * 8));
            hex40[i * 2] = hexChars[b >> 4];
            hex40[i * 2 + 1] = hexChars[b & 0x0f];
            unchecked { ++i; }
        }
        return hex40;
    }
}
