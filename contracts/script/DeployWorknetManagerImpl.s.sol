// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {WorknetManagerUni} from "../src/worknets/WorknetManagerUni.sol";
import {WorknetManager} from "../src/worknets/WorknetManager.sol";

contract DeployWorknetManagerImpl is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

    function run() external {
        uint256 chainId = block.chainid;
        uint256 deployerPk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(deployerPk);

        address deployed;

        if (chainId == 8453) {
            // Base — WorknetManagerUni
            bytes32 salt = 0x321f343552ceec41f138327c85f826a66c6662f726ebb353bb56aa90ab79abfa;
            bytes memory initcode = abi.encodePacked(
                type(WorknetManagerUni).creationCode,
                abi.encode(
                    0x000000000022D473030F116dDEE9F6B43aC78BA3, // permit2
                    0x498581fF718922c3f8e6A244956aF099B2652b2b, // poolManager
                    0x7C5f5A4bBd8fD63184577525326123B519429bDc, // positionManager
                    0xA3c0c9b65baD0b08107Aa264b0f3dB444b867A71  // stateView
                )
            );
            (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(salt, initcode));
            require(ok, "deploy failed");
            deployed = _predict(salt, keccak256(initcode));
            require(deployed == 0x000011EE4117c52dC0Eb146cBC844cb155B200A9, "addr mismatch");

        } else if (chainId == 1) {
            // ETH — WorknetManagerUni
            bytes32 salt = 0x2f95e3341dd397c2818b5bbc29801d2fae26db138003fd7e233408356a191a94;
            bytes memory initcode = abi.encodePacked(
                type(WorknetManagerUni).creationCode,
                abi.encode(
                    0x000000000022D473030F116dDEE9F6B43aC78BA3, // permit2
                    0x000000000004444c5dc75cB358380D2e3dE08A90, // poolManager
                    0xbD216513d74C8cf14cf4747E6AaA6420FF64ee9e, // positionManager
                    0x7fFE42C4a5DEeA5b0feC41C94C136Cf115597227  // stateView
                )
            );
            (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(salt, initcode));
            require(ok, "deploy failed");
            deployed = _predict(salt, keccak256(initcode));
            require(deployed == 0x0000DD4841bB4e66AF61A5E35204C1606b4a00A9, "addr mismatch");

        } else if (chainId == 42161) {
            // ARB — WorknetManagerUni
            bytes32 salt = 0x64a852ee0b61961173f3a1e24d456034e3b8a719622d773ab3d5c0b75f907e0b;
            bytes memory initcode = abi.encodePacked(
                type(WorknetManagerUni).creationCode,
                abi.encode(
                    0x000000000022D473030F116dDEE9F6B43aC78BA3, // permit2
                    0x360E68faCcca8cA495c1B759Fd9EEe466db9FB32, // poolManager
                    0xd88F38F930b7952f2DB2432Cb002E7abbF3dD869, // positionManager
                    0x76fd297e2d437cd7F76A5F2B02a5ce11c663A86e  // stateView
                )
            );
            (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(salt, initcode));
            require(ok, "deploy failed");
            deployed = _predict(salt, keccak256(initcode));
            require(deployed == 0x000055Ca7d29e8dC7eDEF3892849347214a300A9, "addr mismatch");

        } else if (chainId == 56) {
            // BSC — WorknetManager (PancakeSwap)
            bytes32 salt = 0xde7719868f54bdcdbac875f5c0efef8b9ec111bda20d4bed1fadc89cd1bb1b05;
            bytes memory initcode = abi.encodePacked(
                type(WorknetManager).creationCode,
                abi.encode(
                    0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768, // permit2
                    0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b, // clPoolManager
                    0x55f4c8abA71A1e923edC303eb4fEfF14608cC226, // clPositionManager
                    0x1b81D678ffb9C0263b24A97847620C99d213eB14  // clSwapRouter
                )
            );
            (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(salt, initcode));
            require(ok, "deploy failed");
            deployed = _predict(salt, keccak256(initcode));
            require(deployed == 0x0000269C10feF9B603A228b075F8C99BAE5b00A9, "addr mismatch");

        } else {
            revert("unsupported chain");
        }

        vm.stopBroadcast();
        console.log("WorknetManager impl deployed:", deployed);
    }

    function _predict(bytes32 salt, bytes32 h) internal pure returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, salt, h)
        ))));
    }
}
