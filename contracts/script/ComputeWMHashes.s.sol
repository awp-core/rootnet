// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {WorknetManager} from "../src/worknets/WorknetManager.sol";
import {WorknetManagerUni} from "../src/worknets/WorknetManagerUni.sol";

contract ComputeWMHashes is Script {
    function run() external pure {
        // Base: permit2, poolMgr, posMgr, swapRouter, stateView
        console.log("WMUni(Base):", vm.toString(keccak256(abi.encodePacked(
            type(WorknetManagerUni).creationCode,
            abi.encode(
                0x000000000022D473030F116dDEE9F6B43aC78BA3,
                0x498581fF718922c3f8e6A244956aF099B2652b2b,
                0x7C5f5A4bBd8fD63184577525326123B519429bDc,
                0x6fF5693b99212Da76ad316178A184AB56D299b43,
                0xA3c0c9b65baD0b08107Aa264b0f3dB444b867A71
            )
        ))));

        // ETH
        console.log("WMUni(ETH):", vm.toString(keccak256(abi.encodePacked(
            type(WorknetManagerUni).creationCode,
            abi.encode(
                0x000000000022D473030F116dDEE9F6B43aC78BA3,
                0x000000000004444c5dc75cB358380D2e3dE08A90,
                0xbD216513d74C8cf14cf4747E6AaA6420FF64ee9e,
                0x66a9893cC07D91D95644AEDD05D03f95e1dBA8Af,
                0x7fFE42C4a5DEeA5b0feC41C94C136Cf115597227
            )
        ))));

        // ARB
        console.log("WMUni(ARB):", vm.toString(keccak256(abi.encodePacked(
            type(WorknetManagerUni).creationCode,
            abi.encode(
                0x000000000022D473030F116dDEE9F6B43aC78BA3,
                0x360E68faCcca8cA495c1B759Fd9EEe466db9FB32,
                0xd88F38F930b7952f2DB2432Cb002E7abbF3dD869,
                0xa51afAF359d044F8e56fE74B9575f23142cD4B76,
                0x76fd297e2d437cd7F76A5F2B02a5ce11c663A86e
            )
        ))));

        // BSC: permit2, poolMgr, posMgr, swapRouter (no stateView)
        console.log("WM(BSC):", vm.toString(keccak256(abi.encodePacked(
            type(WorknetManager).creationCode,
            abi.encode(
                0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768,
                0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b,
                0x55f4c8abA71A1e923edC303eb4fEfF14608cC226,
                0x1b81D678ffb9C0263b24A97847620C99d213eB14
            )
        ))));
    }
}
