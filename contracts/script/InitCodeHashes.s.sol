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
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {SubnetManager} from "../src/subnets/SubnetManager.sol";

/// @title InitCodeHashes — Compute initcode hashes for vanity salt mining (tiered)
/// @dev Run: forge script script/InitCodeHashes.s.sol
///      Outputs initCodeHash per contract. Constructor args MUST match Deploy.s.sol exactly.
///      Tiered: higher tiers depend on addresses from lower tiers.
///      Set ADDR_* env vars as you mine each tier, then re-run for the next tier.
contract InitCodeHashes is Script {
    uint256 constant INITIAL_DAILY_EMISSION = 15_800_000 * 1e18;
    uint256 constant EPOCH_DURATION = 1 days;
    uint256 constant TIMELOCK_DELAY = 172800;

    function _addr(string memory key) internal view returns (address) {
        return vm.envOr(key, address(0xdead));
    }

    function run() external view {
        address deployer = vm.addr(vm.envUint("DEPLOYER_PRIVATE_KEY"));
        address guardian = vm.envAddress("GUARDIAN");
        address poolManager = vm.envAddress("POOL_MANAGER");
        address positionManager = vm.envAddress("POSITION_MANAGER");
        address permit2Addr = vm.envAddress("PERMIT2");
        uint64 vanityRule = uint64(vm.envOr("VANITY_RULE", uint256(0)));

        // Previously mined addresses (set these as you mine each tier)
        address awp = _addr("ADDR_AWP_TOKEN");
        address treasury = _addr("ADDR_TREASURY");
        address rootNet = _addr("ADDR_ROOTNET");
        address emissionImpl = _addr("ADDR_EMISSION_IMPL");
        address vault = _addr("ADDR_STAKING_VAULT");
        address stakeNft = _addr("ADDR_STAKE_NFT");

        console.log("=== InitCode Hashes (for salt.json) ===");
        console.log("Deployer:", deployer);
        console.log("");

        // Tier 1: No cross-contract dependencies
        console.log("--- Tier 1 (no dependencies) ---");
        _logHash("AWPToken", abi.encodePacked(type(AWPToken).creationCode, abi.encode("AWP Token", "AWP", deployer)));
        _logHash("AlphaTokenFactory", abi.encodePacked(type(AlphaTokenFactory).creationCode, abi.encode(deployer, vanityRule)));
        _logHash("AWPEmission_impl", abi.encodePacked(type(AWPEmission).creationCode));
        _logHash("SubnetManager_impl", abi.encodePacked(type(SubnetManager).creationCode));

        // Tier 2: Treasury (no dependency on other deployed contracts)
        console.log("");
        console.log("--- Tier 2 ---");
        {
            address[] memory proposers = new address[](0);
            address[] memory executors = new address[](1);
            executors[0] = address(0);
            _logHash("Treasury", abi.encodePacked(type(Treasury).creationCode, abi.encode(TIMELOCK_DELAY, proposers, executors, deployer)));
        }

        // Tier 3: AWPRegistry (depends on Treasury)
        console.log("");
        console.log("--- Tier 3 (depends on Treasury) ---");
        _logHash("AWPRegistry", abi.encodePacked(type(AWPRegistry).creationCode, abi.encode(deployer, treasury, guardian)));

        // Tier 4: Depends on AWPRegistry + AWP
        console.log("");
        console.log("--- Tier 4 (depends on AWPRegistry + AWP) ---");
        _logHash("SubnetNFT", abi.encodePacked(type(SubnetNFT).creationCode, abi.encode("AWP Subnet", "AWPSUB", rootNet)));
        _logHash("AccessManager", abi.encodePacked(type(AccessManager).creationCode, abi.encode(rootNet)));
        _logHash("LPManager", abi.encodePacked(type(LPManager).creationCode, abi.encode(rootNet, poolManager, positionManager, permit2Addr, awp)));
        _logHash("StakingVault", abi.encodePacked(type(StakingVault).creationCode, abi.encode(rootNet)));

        // AWPEmission proxy (depends on impl + AWP + Treasury)
        // Use uint256(0) as genesisTime placeholder — actual value set at deploy time.
        // This matches Predict.s.sol convention: hash is stable for salt mining.
        // NOTE: genesisTime is uint256(0) placeholder. Deploy.s.sol uses block.timestamp,
        // so the actual proxy address will differ. AWPEmission_proxy cannot be vanity-mined
        // deterministically — re-mine with the actual timestamp immediately before deployment.
        bytes memory initData = abi.encodeCall(AWPEmission.initialize, (awp, treasury, INITIAL_DAILY_EMISSION, uint256(0), EPOCH_DURATION));
        _logHash("AWPEmission_proxy", abi.encodePacked(type(ERC1967Proxy).creationCode, abi.encode(emissionImpl, initData)));

        // Tier 5: StakeNFT (depends on StakingVault + AWPRegistry + AWP)
        console.log("");
        console.log("--- Tier 5 (depends on StakingVault) ---");
        _logHash("StakeNFT", abi.encodePacked(type(StakeNFT).creationCode, abi.encode(awp, vault, rootNet)));

        // Tier 6: DAO (depends on StakeNFT + AWP + Treasury)
        console.log("");
        console.log("--- Tier 6 (depends on StakeNFT) ---");
        _logHash("AWPDAO", abi.encodePacked(type(AWPDAO).creationCode, abi.encode(stakeNft, awp, treasury, uint48(7200), uint32(50400), uint256(4))));
    }

    function _logHash(string memory name, bytes memory initCode) internal pure {
        console.log(string.concat(name, ": ", vm.toString(keccak256(initCode))));
    }
}
