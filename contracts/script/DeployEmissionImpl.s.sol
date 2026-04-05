// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";

contract DeployEmissionImpl is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
    address constant AWP_TOKEN = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    bytes32 constant SALT = 0xd37af1bd9051f030d1411478b9e9ee40ad353b83f05bc1e99db913921c32ff34;
    address constant EXPECTED = 0x0000dd24c04AaC097149fb4Fc75D09a8378900ae;

    function run() external {
        uint256 deployerPk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(deployerPk);

        bytes memory initcode = abi.encodePacked(
            type(AWPEmission).creationCode,
            abi.encode(AWP_TOKEN)
        );
        (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(SALT, initcode));
        require(ok, "deploy failed");

        address deployed = _predict(SALT, keccak256(initcode));
        require(deployed == EXPECTED, "addr mismatch");

        vm.stopBroadcast();

        console.log("AWPEmission impl:", deployed);
        console.log("  awpToken:", address(AWPEmission(deployed).awpToken()));
        console.log("  cachedMaxSupply:", AWPEmission(deployed).cachedMaxSupply());
    }

    function _predict(bytes32 salt, bytes32 h) internal pure returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, salt, h)
        ))));
    }
}
