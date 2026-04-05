// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {veAWP} from "../src/core/veAWP.sol";
import {AWPAllocator} from "../src/core/AWPAllocator.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {StubAllocator} from "./StubAllocator.sol";

contract DeployVeAWPAllocator is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
    address constant AWP_TOKEN = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant REGISTRY_PROXY = 0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A;
    address constant GUARDIAN = 0x000002bEfa6A1C99A710862Feb6dB50525dF00A3;

    // Mined salts
    bytes32 constant STUB_SALT = bytes32(0);
    bytes32 constant PROXY_SALT = 0x49c15b64c8474d08aae832898f8b1d6e833c3f4641d1dadff7015ad14ba698d6;
    bytes32 constant VEAWP_SALT = 0x119be13ba2e5ca641f33997ae6cf99020ddbf4587ba332a4b600daddf7774a61;

    // Expected addresses
    address constant EXPECTED_PROXY = 0x0000D6BB5e040E35081b3AaF59DD71b21C9800AA;
    address constant EXPECTED_VEAWP = 0x0000b534C63D78212f1BDCc315165852793A00A8;

    function run() external {
        uint256 deployerPk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPk);

        vm.startBroadcast(deployerPk);

        // ── Step 1: Deploy StubAllocator impl via CREATE2 (salt=0) ──
        bytes memory stubInitcode = type(StubAllocator).creationCode;
        address stubImpl;
        (bool ok1,) = CREATE2_DEPLOYER.call(abi.encodePacked(STUB_SALT, stubInitcode));
        require(ok1, "stub deploy failed");
        stubImpl = _predict(STUB_SALT, keccak256(stubInitcode));
        console.log("1. StubAllocator impl:", stubImpl);

        // ── Step 2: Deploy ERC1967Proxy via CREATE2 (vanity salt) ──
        bytes memory initData = abi.encodeCall(StubAllocator.initialize, (REGISTRY_PROXY, deployer));
        bytes memory proxyInitcode = abi.encodePacked(
            type(ERC1967Proxy).creationCode,
            abi.encode(stubImpl, initData)
        );
        (bool ok2,) = CREATE2_DEPLOYER.call(abi.encodePacked(PROXY_SALT, proxyInitcode));
        require(ok2, "proxy deploy failed");
        address allocatorProxy = _predict(PROXY_SALT, keccak256(proxyInitcode));
        console.log("2. AWPAllocator proxy:", allocatorProxy);
        require(allocatorProxy == EXPECTED_PROXY, "proxy addr mismatch");

        // ── Step 3: Deploy veAWP via CREATE2 (vanity salt) ──
        bytes memory veawpInitcode = abi.encodePacked(
            type(veAWP).creationCode,
            abi.encode(AWP_TOKEN, allocatorProxy, GUARDIAN)
        );
        (bool ok3,) = CREATE2_DEPLOYER.call(abi.encodePacked(VEAWP_SALT, veawpInitcode));
        require(ok3, "veAWP deploy failed");
        address veawpAddr = _predict(VEAWP_SALT, keccak256(veawpInitcode));
        console.log("3. veAWP:", veawpAddr);
        require(veawpAddr == EXPECTED_VEAWP, "veAWP addr mismatch");

        // ── Step 4: Deploy real AWPAllocator impl (with correct veAWP immutable) ──
        AWPAllocator realImpl = new AWPAllocator(REGISTRY_PROXY, veawpAddr);
        console.log("4. AWPAllocator real impl:", address(realImpl));

        // ── Step 5: Upgrade proxy to real impl ──
        // deployer is guardian of stub, so can upgrade
        UUPSUpgradeable(allocatorProxy).upgradeToAndCall(address(realImpl), "");
        console.log("5. Proxy upgraded to real impl");

        // ── Step 6: Transfer guardian to real GUARDIAN ──
        // The stub's guardian is deployer; after upgrade, storage guardian is still deployer
        // Real AWPAllocator's guardian is storage-based, set during initialize (which already ran with deployer)
        // We need to set guardian to the real GUARDIAN
        // AWPAllocator has setGuardian(address) callable by current guardian
        AWPAllocator(allocatorProxy).setGuardian(GUARDIAN);
        console.log("6. Guardian transferred to:", GUARDIAN);

        vm.stopBroadcast();

        // ── Verify ──
        console.log("\n=== Verification ===");
        console.log("veAWP.awpToken:", address(veAWP(veawpAddr).awpToken()));
        console.log("veAWP.awpAllocator:", veAWP(veawpAddr).awpAllocator());
        console.log("veAWP.guardian:", veAWP(veawpAddr).guardian());
        console.log("Allocator.awpRegistry:", AWPAllocator(allocatorProxy).awpRegistry());
        console.log("Allocator.veAWP:", AWPAllocator(allocatorProxy).veAWP());
        console.log("Allocator.guardian:", AWPAllocator(allocatorProxy).guardian());
    }

    function _predict(bytes32 salt, bytes32 initcodeHash) internal pure returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, salt, initcodeHash)
        ))));
    }
}
