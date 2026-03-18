// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {AWPEmission} from "../../src/token/AWPEmission.sol";

/// @dev EIP-712 signing base class for emission-related tests to inherit
abstract contract EmissionSigningHelper is Test {
    bytes32 private constant ALLOCATION_TYPEHASH =
        keccak256("SubmitAllocations(address[] recipients,uint96[] weights,uint256 nonce,uint256 effectiveEpoch)");

    function _signAllocations(
        uint256 pk,
        address[] memory recipients,
        uint96[] memory weights,
        uint256 nonce,
        uint256 effectiveEpoch,
        address emissionAddr
    ) internal view returns (bytes memory) {
        bytes32 structHash = keccak256(
            abi.encode(
                ALLOCATION_TYPEHASH,
                _hashAddressArray(recipients),
                _hashUint96Array(weights),
                nonce,
                effectiveEpoch
            )
        );
        bytes32 digest = _getEmissionDigest(structHash, emissionAddr);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(pk, digest);
        return abi.encodePacked(r, s, v);
    }

    function _hashAddressArray(address[] memory arr) internal pure returns (bytes32) {
        bytes32[] memory hashes = new bytes32[](arr.length);
        for (uint256 i = 0; i < arr.length; i++) {
            hashes[i] = bytes32(uint256(uint160(arr[i])));
        }
        return keccak256(abi.encodePacked(hashes));
    }

    function _hashUint96Array(uint96[] memory arr) internal pure returns (bytes32) {
        bytes32[] memory hashes = new bytes32[](arr.length);
        for (uint256 i = 0; i < arr.length; i++) {
            hashes[i] = bytes32(uint256(arr[i]));
        }
        return keccak256(abi.encodePacked(hashes));
    }

    function _getEmissionDigest(bytes32 structHash, address emissionAddr) internal view returns (bytes32) {
        bytes32 domainSeparator = keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
                keccak256("AWPEmission"),
                keccak256("2"),
                block.chainid,
                emissionAddr
            )
        );
        return keccak256(abi.encodePacked("\x19\x01", domainSeparator, structHash));
    }
}
