// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {WorknetTokenFactory} from "../src/token/WorknetTokenFactory.sol";

contract DeployWorknetTokenFactory is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
    bytes32 constant SALT = 0x4436ef272f01550d498dea59b30e1ce1b93eab1ba724a779d7cf9741b1b474a2;
    address constant EXPECTED = 0x000058EF25751Bb3687eB314185B46b942bE00AF;

    function run() external {
        uint256 deployerPk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPk);
        uint64 vanityRule = uint64(vm.envUint("VANITY_RULE"));

        console.log("deployer:", deployer);
        console.log("vanityRule:", vanityRule);

        // Build initcode
        bytes memory initcode = abi.encodePacked(
            type(WorknetTokenFactory).creationCode,
            abi.encode(deployer, vanityRule)
        );

        vm.startBroadcast(deployerPk);

        // Deploy via deterministic CREATE2 deployer
        (bool success,) = CREATE2_DEPLOYER.call(abi.encodePacked(SALT, initcode));
        require(success, "CREATE2 deploy failed");

        vm.stopBroadcast();

        // Verify address
        address deployed = address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, SALT, keccak256(initcode))
        ))));
        require(deployed == EXPECTED, "address mismatch");

        console.log("\nDeployed:", deployed);
        console.log("vanityRule:", WorknetTokenFactory(deployed).vanityRule());
        console.log("owner:", WorknetTokenFactory(deployed).owner());
    }
}
