// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";

contract DeployRegistryImpl is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
    bytes32 constant SALT = 0xb1085e1295d5ceaf9614c81dbdc805579fe30343179fc5b88ba037008bbba884;
    address constant EXPECTED = 0x00001a5aa14608d8D9971656f94022c88D4a001a;

    // 8 immutable addresses
    address constant AWP_TOKEN    = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant AWP_WORKNET  = 0x00000bfbdEf8533E5F3228c9C846522D906100A7;
    address constant FACTORY      = 0x000058EF25751Bb3687eB314185B46b942bE00AF;
    address constant EMISSION     = 0x3C9cB73f8B81083882c5308Cce4F31f93600EaA9;
    address constant LP_MANAGER   = 0x00001961b9AcCD86b72DE19Be24FaD6f7c5b00A2;
    address constant ALLOCATOR    = 0x0000D6BB5e040E35081b3AaF59DD71b21C9800AA;
    address constant VEAWP        = 0x0000b534C63D78212f1BDCc315165852793A00A8;
    address constant TREASURY     = 0x82562023a053025F3201785160CaE6051efD759e;

    function run() external {
        uint256 deployerPk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(deployerPk);

        bytes memory initcode = abi.encodePacked(
            type(AWPRegistry).creationCode,
            abi.encode(AWP_TOKEN, AWP_WORKNET, FACTORY, EMISSION, LP_MANAGER, ALLOCATOR, VEAWP, TREASURY)
        );

        (bool ok,) = CREATE2_DEPLOYER.call(abi.encodePacked(SALT, initcode));
        require(ok, "deploy failed");

        address deployed = _predict(SALT, keccak256(initcode));
        require(deployed == EXPECTED, "address mismatch");

        vm.stopBroadcast();

        console.log("AWPRegistry impl:", deployed);
        AWPRegistry reg = AWPRegistry(deployed);
        console.log("  awpToken:", reg.awpToken());
        console.log("  awpWorkNet:", reg.awpWorkNet());
        console.log("  worknetTokenFactory:", reg.worknetTokenFactory());
        console.log("  awpEmission:", reg.awpEmission());
        console.log("  lpManager:", reg.lpManager());
        console.log("  awpAllocator:", reg.awpAllocator());
        console.log("  veAWP:", reg.veAWP());
        console.log("  treasury:", reg.treasury());
    }

    function _predict(bytes32 salt, bytes32 h) internal pure returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, salt, h)
        ))));
    }
}
