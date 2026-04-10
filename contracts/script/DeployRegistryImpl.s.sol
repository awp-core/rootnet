// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script, console} from "forge-std/Script.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";

contract DeployRegistryImpl is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
    bytes32 constant SALT = 0x94b7e4c4569ad6ebf0a3eaa942def841b19c5a18a24052a1284298696e608558;
    address constant EXPECTED = 0x00002644AA845171B0DBDCEf0B3B0E2119e5001a;

    function run() external {
        bytes memory initcode = abi.encodePacked(
            type(AWPRegistry).creationCode,
            abi.encode(
                0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1, // awpToken
                0x00000bfbdEf8533E5F3228c9C846522D906100A7, // awpWorkNet
                0x00000a82b06Ea5b5BdD6003fbfb9602FA531CAFE, // worknetTokenFactory (NEW)
                0x3C9cB73f8B81083882c5308Cce4F31f93600EaA9, // awpEmission
                0x00001961b9AcCD86b72DE19Be24FaD6f7c5b00A2, // lpManager
                0x0000D6BB5e040E35081b3AaF59DD71b21C9800AA, // awpAllocator
                0x0000b534C63D78212f1BDCc315165852793A00A8, // veAWP
                0x82562023a053025F3201785160CaE6051efD759e  // treasury
            )
        );

        vm.startBroadcast();
        (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(SALT, initcode));
        require(ok, "CREATE2 failed");
        vm.stopBroadcast();

        address deployed = address(uint160(uint256(keccak256(abi.encodePacked(
            bytes1(0xff), CREATE2_DEPLOYER, SALT, keccak256(initcode)
        )))));
        require(deployed == EXPECTED, "address mismatch");
        console.log("AWPRegistry impl deployed at:", deployed);
    }
}
