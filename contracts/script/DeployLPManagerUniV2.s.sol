// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {LPManagerUni} from "../src/core/LPManagerUni.sol";

contract DeployLPManagerUniV2 is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

    function run() external {
        uint256 chainId = block.chainid;
        uint256 deployerPk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(deployerPk);

        address deployed;

        if (chainId == 8453) {
            bytes32 salt = 0x27a6a21d26992c8b1e7ca7491b048f2b52e0b67d439b293c8ab27759dba830ee;
            bytes memory initcode = abi.encodePacked(type(LPManagerUni).creationCode, abi.encode(
                0x000000000022D473030F116dDEE9F6B43aC78BA3,
                0x498581fF718922c3f8e6A244956aF099B2652b2b,
                0x7C5f5A4bBd8fD63184577525326123B519429bDc,
                0xA3c0c9b65baD0b08107Aa264b0f3dB444b867A71
            ));
            (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(salt, initcode));
            require(ok, "deploy failed");
            deployed = _predict(salt, keccak256(initcode));
            require(deployed == 0x0000Da504981F02e1bDa04A48E6dded5BAC200a5, "addr mismatch");

        } else if (chainId == 1) {
            bytes32 salt = 0xb86545d7dbc5757fd15295e10dd9e2e72b185c4061d8bdaf3bf86407ad649ed4;
            bytes memory initcode = abi.encodePacked(type(LPManagerUni).creationCode, abi.encode(
                0x000000000022D473030F116dDEE9F6B43aC78BA3,
                0x000000000004444c5dc75cB358380D2e3dE08A90,
                0xbD216513d74C8cf14cf4747E6AaA6420FF64ee9e,
                0x7fFE42C4a5DEeA5b0feC41C94C136Cf115597227
            ));
            (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(salt, initcode));
            require(ok, "deploy failed");
            deployed = _predict(salt, keccak256(initcode));
            require(deployed == 0x0000086f63BF08612aB305C0e97E95aD797600a5, "addr mismatch");

        } else if (chainId == 42161) {
            bytes32 salt = 0x2c2f72450708f4565a46c1d988b3426c813c2ac40fa4072b1dee7121e54997db;
            bytes memory initcode = abi.encodePacked(type(LPManagerUni).creationCode, abi.encode(
                0x000000000022D473030F116dDEE9F6B43aC78BA3,
                0x360E68faCcca8cA495c1B759Fd9EEe466db9FB32,
                0xd88F38F930b7952f2DB2432Cb002E7abbF3dD869,
                0x76fd297e2d437cd7F76A5F2B02a5ce11c663A86e
            ));
            (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(salt, initcode));
            require(ok, "deploy failed");
            deployed = _predict(salt, keccak256(initcode));
            require(deployed == 0x0000d2F325a08b2d8656E880d4c325E13bb600a5, "addr mismatch");

        } else {
            revert("BSC uses LPManager(PCS), not LPManagerUni");
        }

        vm.stopBroadcast();
        console.log("LPManagerUni v2 deployed:", deployed);
    }

    function _predict(bytes32 salt, bytes32 h) internal pure returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, salt, h)
        ))));
    }
}
