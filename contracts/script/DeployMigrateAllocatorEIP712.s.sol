// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {MigrateAllocatorEIP712} from "./MigrateAllocatorEIP712.sol";

contract DeployMigrateAllocatorEIP712 is Script {
    function run() external {
        uint256 deployerPk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(deployerPk);

        MigrateAllocatorEIP712 migration = new MigrateAllocatorEIP712();
        console.log("Migration impl deployed:", address(migration));

        vm.stopBroadcast();
    }
}
