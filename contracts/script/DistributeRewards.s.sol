// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

interface IWorknetManager {
    function transferToken(address token, address to, uint256 amount) external;
}

/// @title DistributeRewards — Read reward2.csv, transfer 100k AWP to deployer, distribute rest by weight
/// @dev Run: forge script script/DistributeRewards.s.sol --tc DistributeRewards --rpc-url $BASE_RPC_URL --broadcast --slow
contract DistributeRewards is Script {
    address constant WORKNET_MANAGER = 0xede5BDBD8F9C68d7Bf4d6Fae980865aA97F7D3fE;
    address constant AWP_TOKEN = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant DEPLOYER = 0x000000000000b67A7D3E5D8D88B5397b8c197861;
    uint256 constant DEPLOYER_RESERVE = 100_000 ether; // 100,000 AWP

    function run() external {
        // 1. Read CSV
        string memory csv = vm.readFile("../reward2.csv");
        string[] memory lines = vm.split(csv, "\n");

        uint256 count;
        address[] memory recipients = new address[](lines.length);
        uint256[] memory weights = new uint256[](lines.length);
        uint256 totalWeight;

        for (uint256 i = 0; i < lines.length; i++) {
            if (bytes(lines[i]).length == 0) continue;
            string[] memory parts = vm.split(lines[i], ",");
            address addr = vm.parseAddress(parts[0]);
            uint256 w = vm.parseUint(parts[1]);
            recipients[count] = addr;
            weights[count] = w;
            totalWeight += w;
            count++;
        }

        console.log("Recipients:", count);
        console.log("Total weight:", totalWeight);

        // 2. Check AWP balance in WorknetManager
        uint256 awpBalance = IERC20(AWP_TOKEN).balanceOf(WORKNET_MANAGER);
        console.log("AWP in WorknetManager:", awpBalance / 1e18);
        require(awpBalance > DEPLOYER_RESERVE, "Insufficient AWP balance");

        uint256 distributable = awpBalance - DEPLOYER_RESERVE;
        console.log("Deployer reserve:", DEPLOYER_RESERVE / 1e18);
        console.log("Distributable AWP:", distributable / 1e18);

        // 3. Dry run: show first 5 distributions
        console.log("--- Preview (first 5) ---");
        for (uint256 i = 0; i < 5 && i < count; i++) {
            uint256 share = distributable * weights[i] / totalWeight;
            console.log(recipients[i], share / 1e18, "AWP");
        }

        // 4. Execute
        uint256 pk = vm.envUint("ADMIN_PRIVATE_KEY");
        vm.startBroadcast(pk);

        IWorknetManager wm = IWorknetManager(WORKNET_MANAGER);

        // Transfer 100k AWP to deployer
        wm.transferToken(AWP_TOKEN, DEPLOYER, DEPLOYER_RESERVE);
        console.log("Transferred 100k AWP to deployer");

        // Distribute rest by weight
        uint256 totalSent;
        for (uint256 i = 0; i < count; i++) {
            uint256 share = distributable * weights[i] / totalWeight;
            if (share == 0) continue;
            wm.transferToken(AWP_TOKEN, recipients[i], share);
            totalSent += share;
        }

        vm.stopBroadcast();

        console.log("Total distributed:", totalSent / 1e18, "AWP");
        console.log("Dust remaining:", (distributable - totalSent) / 1e18, "AWP");
    }
}
