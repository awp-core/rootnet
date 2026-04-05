// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {WorknetTokenFactory} from "../src/token/WorknetTokenFactory.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";

contract DeployFactoryAndRegistryV2 is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

    // Factory
    bytes32 constant FACTORY_SALT = 0x9fa06123192c9e00ea03dfc8dee19e8fc1ce41b7705e35ead4076ed7a9a3751c;
    address constant EXPECTED_FACTORY = 0x0000D4996BDBb99c772e3fA9f0e94AB52AAFFAC7;

    // Registry impl
    bytes32 constant REGISTRY_SALT = 0x8f135c59746bd716cb229e20a48c8ed0ed1aa12de17ec410031f598d5bfe0e01;
    address constant EXPECTED_REGISTRY = 0x00000d52427c825d4aE72195535Ca5e210c8001a;

    // Immutables for Registry
    address constant AWP_TOKEN    = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant AWP_WORKNET  = 0x00000bfbdEf8533E5F3228c9C846522D906100A7;
    address constant EMISSION     = 0x3C9cB73f8B81083882c5308Cce4F31f93600EaA9;
    address constant LP_MANAGER   = 0x00001961b9AcCD86b72DE19Be24FaD6f7c5b00A2;
    address constant ALLOCATOR    = 0x0000D6BB5e040E35081b3AaF59DD71b21C9800AA;
    address constant VEAWP        = 0x0000b534C63D78212f1BDCc315165852793A00A8;
    address constant TREASURY     = 0x82562023a053025F3201785160CaE6051efD759e;
    address constant REGISTRY_PROXY = 0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A;

    function run() external {
        uint256 deployerPk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPk);
        uint64 vanityRule = uint64(vm.envUint("VANITY_RULE"));

        vm.startBroadcast(deployerPk);

        // 1. Deploy WorknetTokenFactory via CREATE2
        bytes memory factoryInitcode = abi.encodePacked(
            type(WorknetTokenFactory).creationCode,
            abi.encode(deployer, vanityRule)
        );
        (bool ok1,) = CREATE2_DEPLOYER.call(abi.encodePacked(FACTORY_SALT, factoryInitcode));
        require(ok1, "factory deploy failed");
        address factoryAddr = _predict(FACTORY_SALT, keccak256(factoryInitcode));
        require(factoryAddr == EXPECTED_FACTORY, "factory addr mismatch");
        console.log("1. Factory:", factoryAddr);

        // 2. Configure Factory
        WorknetTokenFactory(factoryAddr).setAddresses(REGISTRY_PROXY);
        console.log("2. Factory configured + ownership renounced");

        // 3. Deploy AWPRegistry impl via CREATE2 (with new factory address)
        bytes memory registryInitcode = abi.encodePacked(
            type(AWPRegistry).creationCode,
            abi.encode(AWP_TOKEN, AWP_WORKNET, factoryAddr, EMISSION, LP_MANAGER, ALLOCATOR, VEAWP, TREASURY)
        );
        (bool ok2,) = CREATE2_DEPLOYER.call(abi.encodePacked(REGISTRY_SALT, registryInitcode));
        require(ok2, "registry impl deploy failed");
        address registryImpl = _predict(REGISTRY_SALT, keccak256(registryInitcode));
        require(registryImpl == EXPECTED_REGISTRY, "registry addr mismatch");
        console.log("3. Registry impl:", registryImpl);

        vm.stopBroadcast();

        // Verify
        console.log("\n=== Verification ===");
        console.log("Factory.vanityRule:", WorknetTokenFactory(factoryAddr).vanityRule());
        console.log("Factory.awpRegistry:", WorknetTokenFactory(factoryAddr).awpRegistry());
        console.log("Factory.configured:", WorknetTokenFactory(factoryAddr).configured());
        console.log("Registry.awpToken:", AWPRegistry(registryImpl).awpToken());
        console.log("Registry.worknetTokenFactory:", AWPRegistry(registryImpl).worknetTokenFactory());
        console.log("Registry.veAWP:", AWPRegistry(registryImpl).veAWP());
    }

    function _predict(bytes32 salt, bytes32 h) internal pure returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, salt, h)
        ))));
    }
}
