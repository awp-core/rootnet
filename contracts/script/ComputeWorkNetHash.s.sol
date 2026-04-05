// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPWorkNet} from "../src/core/AWPWorkNet.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract ComputeWorkNetHash is Script {
    // Standard CREATE2 deployer (available on all chains)
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

    address constant REGISTRY_PROXY = 0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A;
    address constant GUARDIAN = 0x000002bEfa6A1C99A710862Feb6dB50525dF00A3;

    function run() external view {
        // Step 1: AWPWorkNet impl — constructor(awpRegistry)
        bytes memory implInitcode = abi.encodePacked(
            type(AWPWorkNet).creationCode,
            abi.encode(REGISTRY_PROXY)
        );
        bytes32 implInitcodeHash = keccak256(implInitcode);
        console.log("=== AWPWorkNet impl ===");
        console.log("initcode hash:");
        console.logBytes32(implInitcodeHash);

        // Step 2: Predict impl address (we need it for the proxy initcode)
        // Use salt = 0 for impl (no vanity needed for impl)
        bytes32 implSalt = bytes32(0);
        address implAddr = address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, implSalt, implInitcodeHash)
        ))));
        console.log("predicted impl addr (salt=0):", implAddr);

        // Step 3: ERC1967Proxy — constructor(impl, initData)
        bytes memory initData = abi.encodeCall(
            AWPWorkNet.initialize, ("AWP WorkNet", "WORKN", GUARDIAN)
        );
        console.log("\n=== ERC1967Proxy (for vanity search) ===");
        console.log("impl addr:", implAddr);
        console.log("initData:");
        console.logBytes(initData);

        bytes memory proxyInitcode = abi.encodePacked(
            type(ERC1967Proxy).creationCode,
            abi.encode(implAddr, initData)
        );
        bytes32 proxyInitcodeHash = keccak256(proxyInitcode);
        console.log("\nproxy initcode hash:");
        console.logBytes32(proxyInitcodeHash);
        console.log("\nUse this hash with cast create2:");
        console.log("cast create2 --starts-with 0000 --ends-with 00a7 --deployer 0x4e59b44847b379578588920cA78FbF26c0B4956C --init-code-hash <hash>");
    }
}
