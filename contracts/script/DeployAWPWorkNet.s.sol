// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPWorkNet} from "../src/core/AWPWorkNet.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

/// @title DeployAWPWorkNet — Deploy AWPWorkNet impl + proxy via CREATE2 (same address on all chains)
contract DeployAWPWorkNet is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
    address constant REGISTRY_PROXY = 0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A;
    address constant GUARDIAN = 0x000002bEfa6A1C99A710862Feb6dB50525dF00A3;

    // Mined vanity salt for proxy → 0x00000bfbdEf8533E5F3228c9C846522D906100A7
    bytes32 constant IMPL_SALT = bytes32(0);
    bytes32 constant PROXY_SALT = 0x734dd3a7e08031ac12031abe036ac311acba83b154bd585530e680065e1bc14d;

    function run() external {
        uint256 deployerPk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(deployerPk);

        // Step 1: Deploy AWPWorkNet impl via CREATE2 (salt=0)
        bytes memory implInitcode = abi.encodePacked(
            type(AWPWorkNet).creationCode,
            abi.encode(REGISTRY_PROXY)
        );

        address implAddr;
        assembly {
            implAddr := create2(0, add(implInitcode, 0x20), mload(implInitcode), 0) // salt=0
        }
        require(implAddr != address(0), "impl deploy failed");
        console.log("AWPWorkNet impl:", implAddr);

        // Step 2: Deploy ERC1967Proxy via CREATE2 (vanity salt)
        bytes memory initData = abi.encodeCall(
            AWPWorkNet.initialize, ("AWP WorkNet", "WORKN", GUARDIAN)
        );
        bytes memory proxyInitcode = abi.encodePacked(
            type(ERC1967Proxy).creationCode,
            abi.encode(implAddr, initData)
        );

        // Send to deterministic deployer
        (bool success, bytes memory ret) = CREATE2_DEPLOYER.call(
            abi.encodePacked(PROXY_SALT, proxyInitcode)
        );
        require(success, "proxy deploy failed");
        address proxyAddr = address(uint160(uint256(bytes32(ret)) >> 96));

        // If the above doesn't return the address properly, compute it
        if (proxyAddr == address(0)) {
            proxyAddr = address(uint160(uint256(keccak256(
                abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, PROXY_SALT, keccak256(proxyInitcode))
            ))));
        }
        console.log("AWPWorkNet proxy:", proxyAddr);

        vm.stopBroadcast();

        // Verify
        require(proxyAddr == 0x00000bfbdEf8533E5F3228c9C846522D906100A7, "proxy address mismatch!");
        console.log("\n=== Verification ===");
        console.log("name:", AWPWorkNet(proxyAddr).name());
        console.log("symbol:", AWPWorkNet(proxyAddr).symbol());
        console.log("guardian:", AWPWorkNet(proxyAddr).guardian());
        console.log("awpRegistry:", AWPWorkNet(proxyAddr).awpRegistry());
    }
}
