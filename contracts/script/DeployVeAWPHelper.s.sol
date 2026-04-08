// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script, console} from "forge-std/Script.sol";
import {VeAWPHelper} from "../src/core/VeAWPHelper.sol";

contract DeployVeAWPHelper is Script {
    address constant FACTORY = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
    address constant AWP_TOKEN = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant VEAWP = 0x0000b534C63D78212f1BDCc315165852793A00A8;
    bytes32 constant SALT = 0xde024023eb8cbda3f6c5944efcd06a4909ebf456864ad537b7e7cd96ae6d9cd1;

    function run() external {
        bytes memory initcode = abi.encodePacked(
            type(VeAWPHelper).creationCode,
            abi.encode(AWP_TOKEN, VEAWP)
        );

        vm.startBroadcast();
        (bool ok,) = FACTORY.call(abi.encodePacked(SALT, initcode));
        require(ok, "CREATE2 failed");
        vm.stopBroadcast();

        address deployed = address(uint160(uint256(keccak256(abi.encodePacked(
            bytes1(0xff), FACTORY, SALT, keccak256(initcode)
        )))));
        console.log("Deployed at:", deployed);
    }
}
