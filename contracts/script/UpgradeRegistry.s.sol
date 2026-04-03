// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";

contract UpgradeRegistry is Script {
    function run() external {
        uint256 deployerPk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        bytes32 salt = vm.envOr("UPGRADE_SALT", bytes32(0xddf10a97313079280542e037d42ad055215880acde85c21266a0ee71ff639605));

        vm.startBroadcast(deployerPk);

        // Deploy new implementation via CREATE2 for deterministic address
        AWPRegistry newImpl = new AWPRegistry{salt: salt}();
        console.log("New AWPRegistry impl:", address(newImpl));

        vm.stopBroadcast();
    }
}
