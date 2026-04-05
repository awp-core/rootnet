// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";

/// @title AWPEmissionFix — Temporary impl: reset settledEpoch + upgrade back to normal impl in one tx
contract AWPEmissionFix is AWPEmission(0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1) {
    function fixAndUpgrade(uint256 epoch, address newImpl) external {
        if (msg.sender != guardian) revert NotGuardian();
        settledEpoch = epoch;
        upgradeToAndCall(newImpl, "");
    }
}

contract FixBscEpoch is Script {
    function run() external {
        bytes32 salt = keccak256("awp-emission-fix-bsc-v1");
        bytes memory initCode = type(AWPEmissionFix).creationCode;

        uint256 pk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(pk);

        (bool ok, bytes memory ret) = CREATE2_FACTORY.call(abi.encodePacked(salt, initCode));
        require(ok && ret.length == 20, "CREATE2 deploy failed");

        address deployed;
        assembly { deployed := mload(add(ret, 20)) }

        vm.stopBroadcast();

        console.log("AWPEmissionFix deployed at:", deployed);
        console.log("Chain ID:", block.chainid);
    }
}
