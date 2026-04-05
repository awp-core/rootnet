// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {TimelockControllerUpgradeable} from "@openzeppelin/contracts-upgradeable/governance/TimelockControllerUpgradeable.sol";

contract DeployAWPDAO is Script {
    address constant CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
    address constant VEAWP = 0x0000b534C63D78212f1BDCc315165852793A00A8;
    address constant TREASURY = 0x82562023a053025F3201785160CaE6051efD759e;
    address constant GUARDIAN = 0x000002bEfa6A1C99A710862Feb6dB50525dF00A3;

    bytes32 constant IMPL_SALT = bytes32(0);
    bytes32 constant PROXY_SALT = 0xe74b5c39fb0938e7a923456a551af80a26b1f1cbac3c444a47ac5d4f1f0cf2a9;
    address constant EXPECTED_PROXY = 0x00006879f79f3Da189b5D0fF6e58ad0127Cc0DA0;

    function run() external {
        uint256 deployerPk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(deployerPk);

        // Step 1: Deploy AWPDAO impl via CREATE2 (salt=0)
        bytes memory implInitcode = abi.encodePacked(
            type(AWPDAO).creationCode,
            abi.encode(VEAWP)
        );
        (bool ok1,) = CREATE2_DEPLOYER.call(abi.encodePacked(IMPL_SALT, implInitcode));
        require(ok1, "impl deploy failed");
        address implAddr = _predict(IMPL_SALT, keccak256(implInitcode));
        console.log("1. AWPDAO impl:", implAddr);

        // Step 2: Deploy ERC1967Proxy via CREATE2 (vanity salt)
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
        (bool ok2,) = CREATE2_DEPLOYER.call(abi.encodePacked(PROXY_SALT, proxyInitcode));
        require(ok2, "proxy deploy failed");
        address proxyAddr = _predict(PROXY_SALT, keccak256(proxyInitcode));
        console.log("2. AWPDAO proxy:", proxyAddr);
        require(proxyAddr == EXPECTED_PROXY, "proxy addr mismatch");

        vm.stopBroadcast();

        // Verify
        AWPDAO dao = AWPDAO(payable(proxyAddr));
        console.log("\n=== Verification ===");
        console.log("name:", dao.name());
        console.log("veAWP:", address(dao.veAWP()));
        console.log("timelock:", dao.timelock());
        console.log("votingDelay:", dao.votingDelay());
        console.log("votingPeriod:", dao.votingPeriod());
        console.log("quorumPercent:", dao.quorumPercent());
        console.log("guardian:", dao.guardian());
    }

    function _predict(bytes32 salt, bytes32 h) internal pure returns (address) {
        return address(uint160(uint256(keccak256(
            abi.encodePacked(bytes1(0xff), CREATE2_DEPLOYER, salt, h)
        ))));
    }
}
