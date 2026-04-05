// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";

interface IEmission {
    function settledEpoch() external view returns (uint256);
    function currentEpoch() external view returns (uint256);
    function pausedUntil() external view returns (uint64);
    function frozenEpoch() external view returns (uint64);
    function pauseEpochUntil(uint64 resumeTime) external;
    function settleEpoch(uint256 limit) external;
    function guardian() external view returns (address);
}

contract SimPause is Script {
    function run() external {
        IEmission e = IEmission(0x3C9cB73f8B81083882c5308Cce4F31f93600EaA9);
        address guardian = e.guardian();
        address newImpl = 0x0000F62411D27f708a89b6b616993ff75B5E00ae;

        console.log("BEFORE upgrade:");
        console.log("  settledEpoch:", e.settledEpoch());
        console.log("  currentEpoch:", e.currentEpoch());

        // Upgrade to new impl
        vm.prank(guardian);
        (bool ok,) = address(e).call(abi.encodeWithSignature("upgradeToAndCall(address,bytes)", newImpl, bytes("")));
        require(ok, "upgrade failed");
        console.log("  upgraded OK");

        // Pause for 1 day
        vm.prank(guardian);
        e.pauseEpochUntil(uint64(block.timestamp + 86400));

        console.log("AFTER pause:");
        console.log("  settledEpoch:", e.settledEpoch());
        console.log("  currentEpoch:", e.currentEpoch());
        console.log("  frozenEpoch:", uint256(e.frozenEpoch()));

        // Try settle
        try e.settleEpoch(100) {
            console.log("  settle: SUCCEEDED (BAD)");
        } catch {
            console.log("  settle: REVERTED (GOOD)");
        }

        // Resume
        vm.prank(guardian);
        e.pauseEpochUntil(0);

        console.log("AFTER resume:");
        console.log("  currentEpoch:", e.currentEpoch());
        console.log("  settledEpoch:", e.settledEpoch());

        // Try settle after resume
        try e.settleEpoch(100) {
            console.log("  settle: SUCCEEDED (GOOD)");
        } catch {
            console.log("  settle: REVERTED (check why)");
        }
    }
}
