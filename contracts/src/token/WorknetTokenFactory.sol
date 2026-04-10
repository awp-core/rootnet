// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {WorknetToken} from "./WorknetToken.sol";

/// @title WorknetTokenFactory — deploys WorknetToken via CREATE2 with vanity address validation
/// @notice Uses transient storage for pending params (TSTORE/TLOAD) to save ~40k gas vs regular SSTORE.
///         WorknetToken has no constructor args — reads params from factory via callback during construction.
/// @dev Vanity rule: 8 positions packed into uint64. See _validateVanityAddress for encoding.
contract WorknetTokenFactory is Ownable {
    address public awpRegistry;
    bool public configured;
    uint64 public immutable vanityRule;
    bytes32 public immutable WORKNET_TOKEN_BYTECODE_HASH;

    // Transient storage slots for pending deploy params (EIP-1153, auto-cleared after tx)
    // Using fixed slots to avoid keccak overhead
    bytes32 private constant _SLOT_NAME1    = bytes32(uint256(keccak256("WorknetTokenFactory.pendingName1")) - 1);
    bytes32 private constant _SLOT_NAME2    = bytes32(uint256(keccak256("WorknetTokenFactory.pendingName2")) - 1);
    bytes32 private constant _SLOT_NAME_LEN = bytes32(uint256(keccak256("WorknetTokenFactory.pendingNameLen")) - 1);
    bytes32 private constant _SLOT_SYM      = bytes32(uint256(keccak256("WorknetTokenFactory.pendingSymbol")) - 1);
    bytes32 private constant _SLOT_SYM_LEN  = bytes32(uint256(keccak256("WorknetTokenFactory.pendingSymbolLen")) - 1);
    bytes32 private constant _SLOT_WID      = bytes32(uint256(keccak256("WorknetTokenFactory.pendingWorknetId")) - 1);

    error NotAWPRegistry();
    error InvalidVanityAddress();
    error SaltAlreadyUsed();

    constructor(address deployer_, uint64 vanityRule_) Ownable(deployer_) {
        vanityRule = vanityRule_;
        WORKNET_TOKEN_BYTECODE_HASH = keccak256(type(WorknetToken).creationCode);
    }

    function setAddresses(address awpRegistry_) external onlyOwner {
        awpRegistry = awpRegistry_;
        configured = true;
        renounceOwnership();
    }

    /// @notice Deploy a new WorknetToken via CREATE2 (only AWPRegistry may call)
    function deploy(uint256 worknetId_, string memory name_, string memory symbol_, bytes32 salt_)
        external
        returns (address)
    {
        if (msg.sender != awpRegistry) revert NotAWPRegistry();

        // Store params in transient storage (WorknetToken reads via pendingName/Symbol/WorknetId)
        _tstore_string(name_, _SLOT_NAME1, _SLOT_NAME2, _SLOT_NAME_LEN);
        _tstore_string(symbol_, _SLOT_SYM, bytes32(0), _SLOT_SYM_LEN);
        _tstore_uint(_SLOT_WID, worknetId_);

        bytes32 effectiveSalt = salt_ == bytes32(0) ? bytes32(worknetId_) : salt_;
        WorknetToken token = new WorknetToken{salt: effectiveSalt}();

        // No delete needed — transient storage auto-clears after tx

        if (vanityRule != 0) {
            _validateVanityAddress(address(token));
        }

        return address(token);
    }

    // ── Public getters for WorknetToken constructor callback (reads from transient storage) ──

    function pendingName() external view returns (string memory) {
        return _tload_string(_SLOT_NAME1, _SLOT_NAME2, _SLOT_NAME_LEN);
    }

    function pendingSymbol() external view returns (string memory) {
        return _tload_string(_SLOT_SYM, bytes32(0), _SLOT_SYM_LEN);
    }

    function pendingWorknetId() external view returns (uint256 v) {
        bytes32 slot = _SLOT_WID;
        assembly { v := tload(slot) }
    }

    // ── Transient storage helpers ──

    function _tstore_uint(bytes32 slot, uint256 val) internal {
        assembly { tstore(slot, val) }
    }

    function _tstore_string(string memory s, bytes32 slot1, bytes32 slot2, bytes32 lenSlot) internal {
        uint256 len = bytes(s).length;
        assembly {
            tstore(lenSlot, len)
            tstore(slot1, mload(add(s, 32)))
            if and(gt(len, 32), iszero(iszero(slot2))) {
                tstore(slot2, mload(add(s, 64)))
            }
        }
    }

    function _tload_string(bytes32 slot1, bytes32 slot2, bytes32 lenSlot) internal view returns (string memory s) {
        uint256 len;
        bytes32 d1;
        bytes32 d2;
        assembly {
            len := tload(lenSlot)
            d1 := tload(slot1)
            if and(gt(len, 32), iszero(iszero(slot2))) {
                d2 := tload(slot2)
            }
        }
        s = new string(len);
        assembly {
            mstore(add(s, 32), d1)
            if gt(len, 32) {
                mstore(add(s, 64), d2)
            }
        }
    }

    // ── Predict / Validate ──

    function predictDeployAddress(bytes32 salt_) external view returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), address(this), salt_, WORKNET_TOKEN_BYTECODE_HASH)
        ))));
    }

    function validateSalt(bytes32 salt_) external view returns (address predicted) {
        predicted = address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), address(this), salt_, WORKNET_TOKEN_BYTECODE_HASH)
        ))));
        if (predicted.code.length > 0) revert SaltAlreadyUsed();
        if (vanityRule != 0) {
            _validateVanityAddress(predicted);
        }
    }

    // ── Vanity validation (unchanged) ──

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
