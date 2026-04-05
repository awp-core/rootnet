// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {LPManagerStub} from "../src/core/LPManagerStub.sol";
import {LPManagerUni} from "../src/core/LPManagerUni.sol";
import {LPManager} from "../src/core/LPManager.sol";
import {WorknetManager} from "../src/worknets/WorknetManager.sol";
import {WorknetManagerUni} from "../src/worknets/WorknetManagerUni.sol";

/// @title ComputeHashes - Print initCodeHashes for salt mining
contract ComputeHashes is Script {
    address constant FACTORY = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

    function _hash(string memory label, bytes memory code) internal pure {
        console.log(label, vm.toString(keccak256(code)));
    }

    function _predict(bytes32 salt, bytes memory code) internal pure returns (address) {
        return address(uint160(uint256(keccak256(abi.encodePacked(
            bytes1(0xff), FACTORY, salt, keccak256(code)
        )))));
    }

    function run() external view {
        address deployer = vm.addr(vm.envUint("DEPLOYER_PRIVATE_KEY"));
        console.log("Deployer:", deployer);

        // Step 1: Stub initCodeHash (no constructor args)
        bytes memory stubCode = type(LPManagerStub).creationCode;
        _hash("LPManagerStub:", stubCode);

        // Step 2: Given a stub address, compute proxy initCodeHash
        // Need stub address first. If SALT_LP_STUB is set, predict it.
        bytes32 stubSalt = vm.envOr("SALT_LP_STUB", bytes32(0));
        if (stubSalt != bytes32(0)) {
            address stub = _predict(stubSalt, stubCode);
            console.log("Stub predicted:", stub);

            bytes memory proxyCode = abi.encodePacked(
                type(ERC1967Proxy).creationCode,
                abi.encode(stub, abi.encodeCall(LPManagerStub.initialize, ()))
            );
            _hash("LPManager proxy:", proxyCode);

            // If proxy salt is also set, predict proxy address
            bytes32 proxySalt = vm.envOr("SALT_LP_PROXY", bytes32(0));
            if (proxySalt != bytes32(0)) {
                console.log("Proxy predicted:", _predict(proxySalt, proxyCode));
            }
        }

        // Step 3: Impl initCodeHashes (chain-specific)
        // Base
        _hash("LPManagerUni(Base):", abi.encodePacked(type(LPManagerUni).creationCode,
            abi.encode(0x000000000022D473030F116dDEE9F6B43aC78BA3, 0x498581fF718922c3f8e6A244956aF099B2652b2b, 0x7C5f5A4bBd8fD63184577525326123B519429bDc)));
        // ETH
        _hash("LPManagerUni(ETH):", abi.encodePacked(type(LPManagerUni).creationCode,
            abi.encode(0x000000000022D473030F116dDEE9F6B43aC78BA3, 0x000000000004444c5dc75cB358380D2e3dE08A90, 0xbD216513d74C8cf14cf4747E6AaA6420FF64ee9e)));
        // ARB
        _hash("LPManagerUni(ARB):", abi.encodePacked(type(LPManagerUni).creationCode,
            abi.encode(0x000000000022D473030F116dDEE9F6B43aC78BA3, 0x360E68faCcca8cA495c1B759Fd9EEe466db9FB32, 0xd88F38F930b7952f2DB2432Cb002E7abbF3dD869)));
        // BSC
        _hash("LPManager(BSC):", abi.encodePacked(type(LPManager).creationCode,
            abi.encode(0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768, 0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b, 0x55f4c8abA71A1e923edC303eb4fEfF14608cC226)));

        // Step 4: WorknetManager impls - need swap router and stateView
        // These will be computed per-chain at deploy time
    }
}
