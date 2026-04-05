// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {WorknetToken} from "./WorknetToken.sol";

/// @title WorknetTokenFactory — deploys WorknetToken via CREATE2 with configurable vanity address validation
/// @notice Uses CREATE2 to deploy an independent WorknetToken instance for each worknet.
///         WorknetToken has no constructor args — reads params from factory via callback during construction.
///         This makes creationCode constant: salt determines address regardless of token name/symbol/worknetId.
///         Vanity salts are universal — mine once, reuse for any worknet.
/// @dev Vanity rule encoding per position (uint8):
///      0-9:   match digit '0'-'9'
///      10-15: match lowercase hex 'a'-'f' (EIP-55 must NOT uppercase)
///      16-21: match uppercase hex 'A'-'F' (EIP-55 must uppercase)
///      >=22:  wildcard — no check for this position
///      8 positions packed into a single uint64:
///      [prefix0, prefix1, prefix2, prefix3, suffix0, suffix1, suffix2, suffix3]
contract WorknetTokenFactory is Ownable {
    /// @notice AWPRegistry contract address (immutable after setAddresses)
    address public awpRegistry;

    /// @notice Whether configuration is complete (true after setAddresses)
    bool public configured;

    /// @notice Vanity address rule: 8 positions packed into uint64
    uint64 public immutable vanityRule;

    /// @notice keccak256(type(WorknetToken).creationCode) — constant, no constructor args
    bytes32 public immutable WORKNET_TOKEN_BYTECODE_HASH;

    // ── Pending deploy params (set before CREATE2, read by WorknetToken constructor, deleted after) ──
    string public pendingName;
    string public pendingSymbol;
    uint256 public pendingWorknetId;

    error NotAWPRegistry();
    error InvalidVanityAddress();
    error SaltAlreadyUsed();

    constructor(address deployer_, uint64 vanityRule_) Ownable(deployer_) {
        vanityRule = vanityRule_;
        WORKNET_TOKEN_BYTECODE_HASH = keccak256(type(WorknetToken).creationCode);
    }

    /// @notice Set the AWPRegistry address and renounce ownership (can only be called once)
    function setAddresses(address awpRegistry_) external onlyOwner {
        awpRegistry = awpRegistry_;
        configured = true;
        renounceOwnership();
    }

    /// @notice Deploy a new WorknetToken instance via CREATE2 (only AWPRegistry may call)
    /// @param worknetId_ Worknet ID
    /// @param name_     Token name
    /// @param symbol_   Token symbol
    /// @param salt_     CREATE2 salt; if bytes32(0), uses worknetId as salt
    /// @return Address of the newly created WorknetToken contract
    function deploy(uint256 worknetId_, string memory name_, string memory symbol_, bytes32 salt_)
        external
        returns (address)
    {
        if (msg.sender != awpRegistry) revert NotAWPRegistry();

        // Store params for WorknetToken constructor callback
        pendingName = name_;
        pendingSymbol = symbol_;
        pendingWorknetId = worknetId_;

        bytes32 effectiveSalt = salt_ == bytes32(0) ? bytes32(worknetId_) : salt_;
        WorknetToken token = new WorknetToken{salt: effectiveSalt}();

        // Clear pending params (gas refund)
        delete pendingName;
        delete pendingSymbol;
        delete pendingWorknetId;

        // Validate vanity address if rule is configured
        if (vanityRule != 0) {
            _validateVanityAddress(address(token));
        }

        return address(token);
    }

    /// @notice Predict the CREATE2 deployment address for a given salt
    /// @dev Salt is universal — address is independent of token name/symbol/worknetId
    function predictDeployAddress(bytes32 salt_) external view returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), address(this), salt_, WORKNET_TOKEN_BYTECODE_HASH)
        ))));
    }

    /// @notice Check if a salt is valid: no address collision + vanity rule compliance
    /// @param salt_ The effective salt (after bytes32(0) → worknetId substitution)
    /// @return predicted The address that would be deployed
    function validateSalt(bytes32 salt_) external view returns (address predicted) {
        predicted = address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), address(this), salt_, WORKNET_TOKEN_BYTECODE_HASH)
        ))));
        if (predicted.code.length > 0) revert SaltAlreadyUsed();
        if (vanityRule != 0) {
            _validateVanityAddress(predicted);
        }
    }

    /// @notice Validate that an address satisfies the configured vanity rule with EIP-55 case checking
    function _validateVanityAddress(address addr) internal view {
        uint64 rule = vanityRule;
        uint160 val = uint160(addr);

        bytes32 eip55Hash;
        bool needsHash;
        for (uint256 i = 0; i < 8;) {
            uint8 e = uint8(rule >> (56 - i * 8));
            if (e >= 10 && e < 22) { needsHash = true; break; }
            unchecked { ++i; }
        }
        if (needsHash) {
            eip55Hash = _computeEip55Hash(val);
        }

        for (uint256 i = 0; i < 8;) {
            uint8 expected = uint8(rule >> (56 - i * 8));
            if (expected < 22) {
                uint256 pos = i < 4 ? i : 32 + i;
                uint8 nibble = uint8(val >> ((39 - pos) * 4)) & 0x0f;
                uint8 want = expected < 16 ? expected : expected - 6;
                if (nibble != want) revert InvalidVanityAddress();
                if (expected >= 10) {
                    uint8 hashByte = uint8(eip55Hash[pos >> 1]);
                    uint8 hashNibble = (pos & 1 == 0) ? (hashByte >> 4) : (hashByte & 0x0f);
                    if ((expected >= 16) != (hashNibble >= 8)) revert InvalidVanityAddress();
                }
            }
            unchecked { ++i; }
        }
    }

    function _computeEip55Hash(uint160 val) internal pure returns (bytes32) {
        bytes memory hex40 = new bytes(40);
        bytes16 hexChars = "0123456789abcdef";
        for (uint256 i = 0; i < 20;) {
            uint8 b = uint8(val >> (152 - i * 8));
            hex40[i * 2] = hexChars[b >> 4];
            hex40[i * 2 + 1] = hexChars[b & 0x0f];
            unchecked { ++i; }
        }
        return keccak256(hex40);
    }
}
