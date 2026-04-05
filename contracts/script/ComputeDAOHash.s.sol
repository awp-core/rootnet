// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {TimelockControllerUpgradeable} from "@openzeppelin/contracts-upgradeable/governance/TimelockControllerUpgradeable.sol";

contract ComputeDAOHash is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
    address constant VEAWP = 0x0000b534C63D78212f1BDCc315165852793A00A8;
    address constant TREASURY = 0x82562023a053025F3201785160CaE6051efD759e;
    address constant GUARDIAN = 0x000002bEfa6A1C99A710862Feb6dB50525dF00A3;

    function run() external view {
        // Step 1: AWPDAO impl (constructor takes veAWP)
        bytes memory implInitcode = abi.encodePacked(
            type(AWPDAO).creationCode,
            abi.encode(VEAWP)
        );
        bytes32 implHash = keccak256(implInitcode);
        address implAddr = _predict(bytes32(0), implHash);
        console.log("=== AWPDAO impl (salt=0) ===");
        console.log("addr:", implAddr);

        // Step 2: ERC1967Proxy initcode hash
        bytes memory initData = abi.encodeCall(
            AWPDAO.initialize,
            (
                TimelockControllerUpgradeable(payable(TREASURY)),
                1 days,   // votingDelay: 1 day
                7 days,   // votingPeriod: 7 days
                1 days,   // lateQuorumExtension: 1 day
                4,       // quorumPercent: 4%
                GUARDIAN
            )
        );
        bytes memory proxyInitcode = abi.encodePacked(
            type(ERC1967Proxy).creationCode,
            abi.encode(implAddr, initData)
        );
        bytes32 proxyHash = keccak256(proxyInitcode);
        console.log("\n=== AWPDAO proxy ===");
        console.log("initcode hash:");
        console.logBytes32(proxyHash);
    }

    function _predict(bytes32 salt, bytes32 h) internal pure returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, salt, h)
        ))));
    }
}
