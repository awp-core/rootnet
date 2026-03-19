// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AlphaTokenFactory} from "../src/token/AlphaTokenFactory.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {AccessManager} from "../src/core/AccessManager.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {SubnetNFT} from "../src/core/SubnetNFT.sol";
import {LPManager} from "../src/core/LPManager.sol";
import {RootNet} from "../src/RootNet.sol";
import {Treasury} from "../src/governance/Treasury.sol";

/// @title TestDeploy — Simplified deployment (used for E2E tests)
contract TestDeploy is Script {
    function run() external {
        address deployer = msg.sender;
        vm.startBroadcast();

        AWPToken awp = new AWPToken("AWP", "AWP", deployer);
        AlphaTokenFactory factory = new AlphaTokenFactory(deployer, 0);

        address[] memory p = new address[](0);
        address[] memory e = new address[](1);
        e[0] = address(0);
        Treasury treasury = new Treasury(0, p, e, deployer);

        RootNet rootNet = new RootNet(deployer, address(treasury), address(0x99));
        SubnetNFT nft = new SubnetNFT("AWPSUB", "AWPSUB", address(rootNet));
        AccessManager am = new AccessManager(address(rootNet));
        LPManager lp = new LPManager(address(rootNet), address(0), address(0), address(0), address(awp));

        AWPEmission em;
        {
            AWPEmission emImpl = new AWPEmission();
            bytes memory initData = abi.encodeCall(AWPEmission.initialize, (address(awp), address(treasury), 15_800_000e18, block.timestamp, 1 days));
            em = AWPEmission(address(new ERC1967Proxy(address(emImpl), initData)));
        }

        // Deploy StakingVault, then StakeNFT (setStakeNFT called in initializeRegistry)
        StakingVault sv = new StakingVault(address(rootNet));
        StakeNFT stakeNft = new StakeNFT(address(awp), address(sv), address(rootNet));

        awp.addMinter(address(em));
        awp.renounceAdmin();
        factory.setAddresses(address(rootNet));
        rootNet.initializeRegistry(address(awp), address(nft), address(factory), address(em), address(lp), address(am), address(sv), address(stakeNft), address(0), "");

        vm.stopBroadcast();

        // Output addresses to console
        console.log("AWPToken", address(awp));
        console.log("RootNet", address(rootNet));
        console.log("SubnetNFT", address(nft));
        console.log("AccessManager", address(am));
        console.log("StakingVault", address(sv));
        console.log("StakeNFT", address(stakeNft));
        console.log("AWPEmission", address(em));
        console.log("Treasury", address(treasury));
        console.log("LPManager", address(lp));
    }
}
