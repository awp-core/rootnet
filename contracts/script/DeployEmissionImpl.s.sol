// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";

/// @title DeployEmissionImpl — Deploy new AWPEmission implementation via deterministic CREATE2
/// @dev Same salt on all chains → identical impl address
contract DeployEmissionImpl is Script {
    bytes32 constant SALT = 0xb1729b95e89705c1ccf369240e9d5a17feaf66177106289f2fdbcb79261649e6;

    function run() external {
        bytes memory initCode = type(AWPEmission).creationCode;

        uint256 pk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(pk);

        (bool ok, bytes memory ret) = CREATE2_FACTORY.call(abi.encodePacked(SALT, initCode));
        require(ok && ret.length == 20, "CREATE2 deploy failed");

        address deployed;
        assembly { deployed := mload(add(ret, 20)) }

        vm.stopBroadcast();

        console.log("AWPEmission new impl:", deployed);
        console.log("Chain ID:", block.chainid);
    }
}
