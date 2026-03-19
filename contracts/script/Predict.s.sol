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
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {SubnetManager} from "../src/subnets/SubnetManager.sol";

/// @title Predict — Print the deterministic addresses for all contracts given the current .env salts
/// @dev Run: forge script script/Predict.s.sol
///      Use this to verify vanity addresses before deployment.
///      No broadcast — pure read-only address computation.
contract Predict is Script {
    address constant DETERMINISTIC_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

    uint256 constant INITIAL_DAILY_EMISSION = 15_800_000 * 1e18;
    uint256 constant EPOCH_DURATION = 1 days;
    uint256 constant TIMELOCK_DELAY = 172800;

    function _predict(bytes32 salt, bytes memory creationCode) internal pure returns (address) {
        bytes32 hash = keccak256(abi.encodePacked(bytes1(0xff), DETERMINISTIC_DEPLOYER, salt, keccak256(creationCode)));
        return address(uint160(uint256(hash)));
    }

    function _readSalt(string memory key) internal view returns (bytes32) {
        return vm.envOr(key, bytes32(0));
    }

    function run() external view {
        address deployer = vm.addr(vm.envUint("DEPLOYER_PRIVATE_KEY"));
        address guardian = vm.envAddress("GUARDIAN");
        address poolManager = vm.envAddress("POOL_MANAGER");
        address positionManager = vm.envAddress("POSITION_MANAGER");
        address permit2Addr = vm.envAddress("PERMIT2");

        // ── Predict all addresses in deployment order ──

        // AWPToken
        address awp = _predict(
            _readSalt("SALT_AWP_TOKEN"),
            abi.encodePacked(type(AWPToken).creationCode, abi.encode("AWP Token", "AWP", deployer))
        );
        console.log("AWPToken:          ", awp);

        // AlphaTokenFactory
        address factory = _predict(
            _readSalt("SALT_ALPHA_FACTORY"),
            abi.encodePacked(type(AlphaTokenFactory).creationCode, abi.encode(deployer, uint64(0x1001FFFF12101514)))
        );
        console.log("AlphaTokenFactory: ", factory);

        // Treasury
        address treasury;
        {
            address[] memory proposers = new address[](0);
            address[] memory executors = new address[](1);
            executors[0] = address(0);
            treasury = _predict(
                _readSalt("SALT_TREASURY"),
                abi.encodePacked(type(Treasury).creationCode, abi.encode(TIMELOCK_DELAY, proposers, executors, deployer))
            );
        }
        console.log("Treasury:          ", treasury);

        // RootNet
        address rootNet = _predict(
            _readSalt("SALT_ROOTNET"),
            abi.encodePacked(type(RootNet).creationCode, abi.encode(deployer, treasury, guardian))
        );
        console.log("RootNet:           ", rootNet);

        // SubnetNFT
        address subnetNFT = _predict(
            _readSalt("SALT_SUBNET_NFT"),
            abi.encodePacked(type(SubnetNFT).creationCode, abi.encode("AWP Subnet", "AWPSUB", rootNet))
        );
        console.log("SubnetNFT:         ", subnetNFT);

        // AccessManager
        address accessMgr = _predict(
            _readSalt("SALT_ACCESS_MANAGER"),
            abi.encodePacked(type(AccessManager).creationCode, abi.encode(rootNet))
        );
        console.log("AccessManager:     ", accessMgr);

        // LPManager
        address lpMgr = _predict(
            _readSalt("SALT_LP_MANAGER"),
            abi.encodePacked(type(LPManager).creationCode, abi.encode(rootNet, poolManager, positionManager, permit2Addr, awp))
        );
        console.log("LPManager:         ", lpMgr);

        // AWPEmission
        address emissionImpl = _predict(
            _readSalt("SALT_EMISSION_IMPL"),
            abi.encodePacked(type(AWPEmission).creationCode)
        );
        console.log("AWPEmission impl:  ", emissionImpl);

        // Note: genesisTime is set at deploy time; using 0 as placeholder for prediction
        bytes memory initData = abi.encodeCall(AWPEmission.initialize, (awp, treasury, INITIAL_DAILY_EMISSION, uint256(0), EPOCH_DURATION));
        address emissionProxy = _predict(
            _readSalt("SALT_EMISSION_PROXY"),
            abi.encodePacked(type(ERC1967Proxy).creationCode, abi.encode(emissionImpl, initData))
        );
        console.log("AWPEmission proxy: ", emissionProxy);

        // StakingVault (stakeNFT set via setStakeNFT after deploy — no circular dependency)
        address vault = _predict(
            _readSalt("SALT_STAKING_VAULT"),
            abi.encodePacked(type(StakingVault).creationCode, abi.encode(rootNet))
        );
        console.log("StakingVault:      ", vault);

        // StakeNFT
        address stakeNft = _predict(
            _readSalt("SALT_STAKE_NFT"),
            abi.encodePacked(type(StakeNFT).creationCode, abi.encode(awp, vault, rootNet))
        );
        console.log("StakeNFT:          ", stakeNft);

        // SubnetManager impl
        address subnetMgrImpl = _predict(
            _readSalt("SALT_SUBNET_MANAGER_IMPL"),
            abi.encodePacked(type(SubnetManager).creationCode)
        );
        console.log("SubnetManager impl:", subnetMgrImpl);

        // AWPDAO
        address dao = _predict(
            _readSalt("SALT_DAO"),
            abi.encodePacked(type(AWPDAO).creationCode, abi.encode(
                stakeNft, awp, treasury,
                uint48(7200), uint32(50400), uint256(4)
            ))
        );
        console.log("AWPDAO:            ", dao);
    }
}
