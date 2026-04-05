// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {LPManagerStub} from "../src/core/LPManagerStub.sol";
import {LPManagerUni} from "../src/core/LPManagerUni.sol";
import {LPManager} from "../src/core/LPManager.sol";
import {WorknetManager} from "../src/worknets/WorknetManager.sol";
import {WorknetManagerUni} from "../src/worknets/WorknetManagerUni.sol";

contract ComputeAllHashes is Script {
    address constant FACTORY = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

    function _predict(bytes32 salt, bytes memory code) internal pure returns (address) {
        return address(uint160(uint256(keccak256(abi.encodePacked(bytes1(0xff), FACTORY, salt, keccak256(code))))));
    }

    function run() external view {
        address deployer = vm.addr(vm.envUint("DEPLOYER_PRIVATE_KEY"));
        console.log("Deployer:", deployer);

        // Stub: constructor(deployer) — deployer is immutable
        bytes memory stubCode = abi.encodePacked(type(LPManagerStub).creationCode, abi.encode(deployer));
        console.log("STUB_HASH:", vm.toString(keccak256(stubCode)));

        // If SALT_LP_STUB is set, predict stub addr and compute proxy hash
        bytes32 stubSalt = vm.envOr("SALT_LP_STUB", bytes32(0));
        if (stubSalt != bytes32(0)) {
            address stub = _predict(stubSalt, stubCode);
            console.log("STUB_ADDR:", stub);

            bytes memory proxyCode = abi.encodePacked(
                type(ERC1967Proxy).creationCode,
                abi.encode(stub, abi.encodeCall(LPManagerStub.initialize, ()))
            );
            console.log("PROXY_HASH:", vm.toString(keccak256(proxyCode)));

            bytes32 proxySalt = vm.envOr("SALT_LP_PROXY", bytes32(0));
            if (proxySalt != bytes32(0)) {
                console.log("PROXY_ADDR:", _predict(proxySalt, proxyCode));
            }
        }

        // LPManager impls (chain-specific constructor args)
        // Base
        console.log("LP_UNI_BASE:", vm.toString(keccak256(abi.encodePacked(type(LPManagerUni).creationCode,
            abi.encode(0x000000000022D473030F116dDEE9F6B43aC78BA3, 0x498581fF718922c3f8e6A244956aF099B2652b2b, 0x7C5f5A4bBd8fD63184577525326123B519429bDc)))));
        // ETH
        console.log("LP_UNI_ETH:", vm.toString(keccak256(abi.encodePacked(type(LPManagerUni).creationCode,
            abi.encode(0x000000000022D473030F116dDEE9F6B43aC78BA3, 0x000000000004444c5dc75cB358380D2e3dE08A90, 0xbD216513d74C8cf14cf4747E6AaA6420FF64ee9e)))));
        // ARB
        console.log("LP_UNI_ARB:", vm.toString(keccak256(abi.encodePacked(type(LPManagerUni).creationCode,
            abi.encode(0x000000000022D473030F116dDEE9F6B43aC78BA3, 0x360E68faCcca8cA495c1B759Fd9EEe466db9FB32, 0xd88F38F930b7952f2DB2432Cb002E7abbF3dD869)))));
        // BSC
        console.log("LP_PCS_BSC:", vm.toString(keccak256(abi.encodePacked(type(LPManager).creationCode,
            abi.encode(0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768, 0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b, 0x55f4c8abA71A1e923edC303eb4fEfF14608cC226)))));

        // WorknetManager impls
        console.log("WM_UNI_BASE:", vm.toString(keccak256(abi.encodePacked(type(WorknetManagerUni).creationCode,
            abi.encode(0x000000000022D473030F116dDEE9F6B43aC78BA3, 0x498581fF718922c3f8e6A244956aF099B2652b2b, 0x7C5f5A4bBd8fD63184577525326123B519429bDc, 0x6fF5693b99212Da76ad316178A184AB56D299b43, 0xA3c0c9b65baD0b08107Aa264b0f3dB444b867A71)))));
        console.log("WM_UNI_ETH:", vm.toString(keccak256(abi.encodePacked(type(WorknetManagerUni).creationCode,
            abi.encode(0x000000000022D473030F116dDEE9F6B43aC78BA3, 0x000000000004444c5dc75cB358380D2e3dE08A90, 0xbD216513d74C8cf14cf4747E6AaA6420FF64ee9e, 0x66a9893cC07D91D95644AEDD05D03f95e1dBA8Af, 0x7fFE42C4a5DEeA5b0feC41C94C136Cf115597227)))));
        console.log("WM_UNI_ARB:", vm.toString(keccak256(abi.encodePacked(type(WorknetManagerUni).creationCode,
            abi.encode(0x000000000022D473030F116dDEE9F6B43aC78BA3, 0x360E68faCcca8cA495c1B759Fd9EEe466db9FB32, 0xd88F38F930b7952f2DB2432Cb002E7abbF3dD869, 0xa51afAF359d044F8e56fE74B9575f23142cD4B76, 0x76fd297e2d437cd7F76A5F2B02a5ce11c663A86e)))));
        console.log("WM_PCS_BSC:", vm.toString(keccak256(abi.encodePacked(type(WorknetManager).creationCode,
            abi.encode(0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768, 0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b, 0x55f4c8abA71A1e923edC303eb4fEfF14608cC226, 0x1b81D678ffb9C0263b24A97847620C99d213eB14)))));
    }
}
