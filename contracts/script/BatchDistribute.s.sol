// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

interface IWorknetManager {
    function transferToken(address token, address to, uint256 amount) external;
    function grantRole(bytes32 role, address account) external;
    function revokeRole(bytes32 role, address account) external;
}

/// @title BatchTransfer — Batch transfer from WorknetManager
contract BatchTransfer {
    address public immutable worknetManager;
    address public immutable token;
    address public immutable admin;

    constructor(address wm_, address token_, address admin_) {
        worknetManager = wm_;
        token = token_;
        admin = admin_;
    }

    function execute(address[] calldata recipients, uint256[] calldata amounts) external {
        require(msg.sender == admin, "not admin");
        require(recipients.length == amounts.length, "length mismatch");
        IWorknetManager wm = IWorknetManager(worknetManager);
        for (uint256 i = 0; i < recipients.length; i++) {
            wm.transferToken(token, recipients[i], amounts[i]);
        }
    }
}

/// @title BatchDistributeScript — Deploy, grant, batch execute (80 per tx), revoke
contract BatchDistributeScript is Script {
    address constant WORKNET_MANAGER = 0xede5BDBD8F9C68d7Bf4d6Fae980865aA97F7D3fE;
    address constant AWP_TOKEN = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant DEPLOYER = 0x000000000000b67A7D3E5D8D88B5397b8c197861;
    uint256 constant DEPLOYER_RESERVE = 100_000 ether;
    uint256 constant BATCH_SIZE = 80;
    bytes32 constant TRANSFER_ROLE = keccak256("TRANSFER_ROLE");

    function run() external {
        // 1. Load CSV
        string memory csv = vm.readFile("../reward2.csv");
        string[] memory lines = vm.split(csv, "\n");

        uint256 count;
        address[] memory addrs = new address[](lines.length);
        uint256[] memory weights = new uint256[](lines.length);
        uint256 totalWeight;

        for (uint256 i = 0; i < lines.length; i++) {
            if (bytes(lines[i]).length == 0) continue;
            string[] memory parts = vm.split(lines[i], ",");
            addrs[count] = vm.parseAddress(parts[0]);
            weights[count] = vm.parseUint(parts[1]);
            totalWeight += weights[count];
            count++;
        }

        uint256 awpBalance = IERC20(AWP_TOKEN).balanceOf(WORKNET_MANAGER);
        uint256 distributable = awpBalance - DEPLOYER_RESERVE;

        console.log("Recipients:", count);
        console.log("Total weight:", totalWeight);
        console.log("AWP balance:", awpBalance / 1e18);
        console.log("Distributable:", distributable / 1e18);
        console.log("Batch size:", BATCH_SIZE);

        // 2. Build full recipient + amount arrays (deployer first)
        uint256 nonZero = 1;
        for (uint256 i = 0; i < count; i++) {
            if (distributable * weights[i] / totalWeight > 0) nonZero++;
        }

        address[] memory recipients = new address[](nonZero);
        uint256[] memory amounts = new uint256[](nonZero);
        uint256 idx;

        recipients[idx] = DEPLOYER;
        amounts[idx] = DEPLOYER_RESERVE;
        idx++;

        uint256 totalSent = DEPLOYER_RESERVE;
        for (uint256 i = 0; i < count; i++) {
            uint256 share = distributable * weights[i] / totalWeight;
            if (share == 0) continue;
            recipients[idx] = addrs[i];
            amounts[idx] = share;
            totalSent += share;
            idx++;
        }

        uint256 batches = (nonZero + BATCH_SIZE - 1) / BATCH_SIZE;
        console.log("Total transfers:", nonZero);
        console.log("Batches:", batches);
        console.log("Total AWP:", totalSent / 1e18);

        // 3. Execute
        uint256 adminPk = vm.envUint("ADMIN_PRIVATE_KEY");
        address admin = vm.addr(adminPk);
        vm.startBroadcast(adminPk);

        BatchTransfer batch = new BatchTransfer(WORKNET_MANAGER, AWP_TOKEN, admin);
        console.log("BatchTransfer:", address(batch));

        IWorknetManager(WORKNET_MANAGER).grantRole(TRANSFER_ROLE, address(batch));

        // Execute in batches
        for (uint256 b = 0; b < batches; b++) {
            uint256 start = b * BATCH_SIZE;
            uint256 end = start + BATCH_SIZE;
            if (end > nonZero) end = nonZero;
            uint256 size = end - start;

            address[] memory batchRecipients = new address[](size);
            uint256[] memory batchAmounts = new uint256[](size);
            for (uint256 i = 0; i < size; i++) {
                batchRecipients[i] = recipients[start + i];
                batchAmounts[i] = amounts[start + i];
            }

            batch.execute(batchRecipients, batchAmounts);
            console.log("Batch done:", b + 1, "size:", size);
        }

        IWorknetManager(WORKNET_MANAGER).revokeRole(TRANSFER_ROLE, address(batch));

        vm.stopBroadcast();

        console.log("Remaining AWP in WM:", IERC20(AWP_TOKEN).balanceOf(WORKNET_MANAGER) / 1e18);
    }
}
