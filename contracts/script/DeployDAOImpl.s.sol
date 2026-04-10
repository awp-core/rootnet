// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;
import {Script, console} from "forge-std/Script.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";

contract DeployDAOImpl is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
    bytes32 constant SALT = 0xb46750de51b09c59e06913013c28adb43b3ab58a7e6d96df353fea5f4ffdd99e;
    address constant EXPECTED = 0x0000B393c3CCEdD0391EeB24808c579725500DA0;
    address constant VEAWP = 0x0000b534C63D78212f1BDCc315165852793A00A8;

    function run() external {
        bytes memory initcode = abi.encodePacked(type(AWPDAO).creationCode, abi.encode(VEAWP));
        vm.startBroadcast();
        (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(SALT, initcode));
        require(ok, "CREATE2 failed");
        vm.stopBroadcast();
        address deployed = address(uint160(uint256(keccak256(abi.encodePacked(
            bytes1(0xff), CREATE2_DEPLOYER, SALT, keccak256(initcode))))));
        require(deployed == EXPECTED, "address mismatch");
        console.log("AWPDAO impl deployed at:", deployed);
    }
}
