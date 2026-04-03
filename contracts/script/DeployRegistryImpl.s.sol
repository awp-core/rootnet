// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";

/// @title DeployRegistryImpl — Deploy new AWPRegistry implementation via deterministic CREATE2
contract DeployRegistryImpl is Script {
    bytes32 constant SALT = 0x28e651489106ed4a60a8ff32a8d83c8c78e9e805965902822f8c7a48ade4df9e;

    function run() external {
        bytes memory initCode = type(AWPRegistry).creationCode;

        uint256 pk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(pk);

        (bool ok, bytes memory ret) = CREATE2_FACTORY.call(abi.encodePacked(SALT, initCode));
        require(ok && ret.length == 20, "CREATE2 deploy failed");

        address deployed;
        assembly { deployed := mload(add(ret, 20)) }

        vm.stopBroadcast();

        console.log("AWPRegistry new impl:", deployed);
        console.log("Chain ID:", block.chainid);
    }
}
