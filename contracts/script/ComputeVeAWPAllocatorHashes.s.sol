// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {veAWP} from "../src/core/veAWP.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {StubAllocator} from "./StubAllocator.sol";

contract ComputeVeAWPAllocatorHashes is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
    address constant AWP_TOKEN = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant REGISTRY_PROXY = 0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A;
    address constant GUARDIAN = 0x000002bEfa6A1C99A710862Feb6dB50525dF00A3;

    function run() external {
        address deployer = vm.addr(vm.envUint("DEPLOYER_PRIVATE_KEY"));
        console.log("deployer:", deployer);

        // Step 1: StubAllocator impl
        bytes32 stubHash = keccak256(type(StubAllocator).creationCode);
        address stubImpl = _predict(bytes32(0), stubHash);
        console.log("\n=== StubAllocator impl (salt=0) ===");
        console.log("addr:", stubImpl);

        // Step 2: Proxy initcode hash — guardian=deployer so deployer can upgrade
        bytes memory initData = abi.encodeCall(StubAllocator.initialize, (REGISTRY_PROXY, deployer));
        bytes memory proxyInitcode = abi.encodePacked(
            type(ERC1967Proxy).creationCode,
            abi.encode(stubImpl, initData)
        );
        bytes32 proxyHash = keccak256(proxyInitcode);
        console.log("\n=== AWPAllocator proxy ===");
        console.log("initcode hash:");
        console.logBytes32(proxyHash);

        // Step 3: veAWP needs allocator proxy address (mined separately)
        console.log("\n=== veAWP ===");
        console.log("After mining allocator proxy, compute:");
        console.log("  keccak256(veAWP.creationCode + encode(AWP_TOKEN, PROXY_ADDR, GUARDIAN))");
    }

    function _predict(bytes32 salt, bytes32 initcodeHash) internal pure returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, salt, initcodeHash)
        ))));
    }
}
