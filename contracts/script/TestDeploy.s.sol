// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AlphaTokenFactory} from "../src/token/AlphaTokenFactory.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {WorknetNFT} from "../src/core/WorknetNFT.sol";
import {LPManager} from "../src/core/LPManager.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {Treasury} from "../src/governance/Treasury.sol";

/// @title TestDeploy — Simplified deployment (used for E2E tests)
contract TestDeploy is Script {
    function run() external {
        // NOTE: Minimal deployment for testing. AWPDAO, WorknetManager impl, and defaultWorknetManagerImpl are omitted.
        address deployer = msg.sender;
        vm.startBroadcast();

        AWPToken awp = new AWPToken("AWP", "AWP", deployer);
        awp.initialMint(200_000_000 * 1e18);
        AlphaTokenFactory factory = new AlphaTokenFactory(deployer, 0);

        address[] memory p = new address[](0);
        address[] memory e = new address[](1);
        e[0] = address(0);
        Treasury treasury = new Treasury(0, p, e, deployer);

        AWPRegistry awpRegistryImpl = new AWPRegistry();
        AWPRegistry awpRegistry = AWPRegistry(address(new ERC1967Proxy(
            address(awpRegistryImpl),
            abi.encodeCall(AWPRegistry.initialize, (deployer, address(treasury), address(0x99)))
        )));
        WorknetNFT nft = new WorknetNFT("AWPSUB", "AWPSUB", address(awpRegistry));
        LPManager lp = new LPManager(address(awpRegistry), address(0), address(0), address(0), address(awp));

        AWPEmission em;
        {
            AWPEmission emImpl = new AWPEmission();
            bytes memory initData = abi.encodeCall(AWPEmission.initialize, (address(awp), deployer, 31_600_000e18, block.timestamp, 1 days, deployer));
            em = AWPEmission(address(new ERC1967Proxy(address(emImpl), initData)));
        }

        // Deploy StakingVault, then StakeNFT
        StakingVault svImpl = new StakingVault();
        StakingVault sv = StakingVault(address(new ERC1967Proxy(
            address(svImpl), abi.encodeCall(StakingVault.initialize, (address(awpRegistry), deployer))
        )));
        StakeNFT stakeNft = new StakeNFT(address(awp), address(sv), address(awpRegistry));

        awp.addMinter(address(em));
        awp.renounceAdmin();
        factory.setAddresses(address(awpRegistry));
        awpRegistry.initializeRegistry(address(awp), address(nft), address(factory), address(em), address(lp), address(sv), address(stakeNft), address(0), "");

        vm.stopBroadcast();

        // Output addresses to console
        console.log("AWPToken", address(awp));
        console.log("AWPRegistry", address(awpRegistry));
        console.log("WorknetNFT", address(nft));
        console.log("StakingVault", address(sv));
        console.log("StakeNFT", address(stakeNft));
        console.log("AWPEmission", address(em));
        console.log("Treasury", address(treasury));
        console.log("LPManager", address(lp));
    }
}
