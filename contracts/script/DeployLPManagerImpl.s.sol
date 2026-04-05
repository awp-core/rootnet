// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {LPManagerUni} from "../src/core/LPManagerUni.sol";
import {LPManager} from "../src/core/LPManager.sol";

contract DeployLPManagerImpl is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

    function run() external {
        uint256 chainId = block.chainid;
        uint256 deployerPk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(deployerPk);

        address deployed;

        if (chainId == 8453) {
            bytes32 salt = 0x23858a9994a990b557ead370beadd71174c430c85ef7c4e58a7e05f62545615a;
            bytes memory initcode = abi.encodePacked(type(LPManagerUni).creationCode, abi.encode(
                0x000000000022D473030F116dDEE9F6B43aC78BA3,
                0x498581fF718922c3f8e6A244956aF099B2652b2b,
                0x7C5f5A4bBd8fD63184577525326123B519429bDc
            ));
            (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(salt, initcode));
            require(ok, "deploy failed");
            deployed = _predict(salt, keccak256(initcode));
            require(deployed == 0x0000fabe9E310Ae8a41449dcA937d1aD998000A5, "addr mismatch");

        } else if (chainId == 1) {
            bytes32 salt = 0xa56321abc115f255d95ab18b45f8143c4f53aef3f2c4e8cfc65157630c95a411;
            bytes memory initcode = abi.encodePacked(type(LPManagerUni).creationCode, abi.encode(
                0x000000000022D473030F116dDEE9F6B43aC78BA3,
                0x000000000004444c5dc75cB358380D2e3dE08A90,
                0xbD216513d74C8cf14cf4747E6AaA6420FF64ee9e
            ));
            (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(salt, initcode));
            require(ok, "deploy failed");
            deployed = _predict(salt, keccak256(initcode));
            require(deployed == 0x0000fa2a290DFa5f72Fa4f9fb5d502cd310600A5, "addr mismatch");

        } else if (chainId == 42161) {
            bytes32 salt = 0xc06787c312a237fecf2736787229fa6345acaa5b052a865096ca0628d9a5d420;
            bytes memory initcode = abi.encodePacked(type(LPManagerUni).creationCode, abi.encode(
                0x000000000022D473030F116dDEE9F6B43aC78BA3,
                0x360E68faCcca8cA495c1B759Fd9EEe466db9FB32,
                0xd88F38F930b7952f2DB2432Cb002E7abbF3dD869
            ));
            (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(salt, initcode));
            require(ok, "deploy failed");
            deployed = _predict(salt, keccak256(initcode));
            require(deployed == 0x00009016a913EE5535A4F6AD31cA20CA43a000A5, "addr mismatch");

        } else if (chainId == 56) {
            bytes32 salt = 0x8d3f81303127178fe0816b3c19419738f79ccafc076c0fd4cc9a70b068f8348c;
            bytes memory initcode = abi.encodePacked(type(LPManager).creationCode, abi.encode(
                0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768,
                0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b,
                0x55f4c8abA71A1e923edC303eb4fEfF14608cC226
            ));
            (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(salt, initcode));
            require(ok, "deploy failed");
            deployed = _predict(salt, keccak256(initcode));
            require(deployed == 0x00009d0D53408e5ecA86E96B6ba9988571Fd00A5, "addr mismatch");

        } else {
            revert("unsupported chain");
        }

        vm.stopBroadcast();
        console.log("LPManager impl deployed:", deployed);
    }

    function _predict(bytes32 salt, bytes32 h) internal pure returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, salt, h)
        ))));
    }
}
