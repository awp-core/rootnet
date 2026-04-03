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
import {LPManagerUni} from "../src/core/LPManagerUni.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {WorknetManager} from "../src/worknets/WorknetManager.sol";
import {WorknetManagerUni} from "../src/worknets/WorknetManagerUni.sol";

/// @title Predict — Print the deterministic addresses for all contracts given the current .env salts
/// @dev Run: forge script script/Predict.s.sol
///      Use this to verify vanity addresses before deployment.
///      No broadcast — pure read-only address computation.
///      Constructor args MUST match Deploy.s.sol exactly.
contract Predict is Script {
    address constant DETERMINISTIC_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

    uint256 constant INITIAL_DAILY_EMISSION = 31_600_000 * 1e18;
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
        uint64 vanityRule = uint64(vm.envOr("VANITY_RULE", uint256(0)));

        // AWPToken (3-param constructor, no initialMint — matches Deploy.s.sol)
        address awp = _predict(
            _readSalt("SALT_AWP_TOKEN"),
            abi.encodePacked(type(AWPToken).creationCode, abi.encode("AWP Token", "AWP", deployer))
        );
        console.log("AWPToken:          ", awp);

        // AlphaTokenFactory (read VANITY_RULE from env)
        address factory = _predict(
            _readSalt("SALT_ALPHA_FACTORY"),
            abi.encodePacked(type(AlphaTokenFactory).creationCode, abi.encode(deployer, vanityRule))
        );
        console.log("AlphaTokenFactory: ", factory);

        // Treasury (matches Deploy.s.sol params)
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

        // AWPRegistry impl (no args — matches Deploy.s.sol)
        address awpRegistryImpl = _predict(
            _readSalt("SALT_AWP_REGISTRY_IMPL"),
            abi.encodePacked(type(AWPRegistry).creationCode)
        );
        console.log("AWPRegistry impl:  ", awpRegistryImpl);

        // AWPRegistry proxy (ERC1967Proxy with initialize data)
        address awpRegistry;
        {
            bytes memory registryInitData = abi.encodeCall(AWPRegistry.initialize, (deployer, treasury, guardian));
            awpRegistry = _predict(
                _readSalt("SALT_AWP_REGISTRY"),
                abi.encodePacked(type(ERC1967Proxy).creationCode, abi.encode(awpRegistryImpl, registryInitData))
            );
        }
        console.log("AWPRegistry proxy: ", awpRegistry);

        // WorknetNFT
        address worknetNFT = _predict(
            _readSalt("SALT_WORKNET_NFT"),
            abi.encodePacked(type(WorknetNFT).creationCode, abi.encode("AWP Worknet", "AWPSUB", awpRegistry))
        );
        console.log("WorknetNFT:         ", worknetNFT);

        // LPManager (chain-conditional: PancakeSwap on BSC, Uniswap on others)
        address lpMgr;
        if (block.chainid == 56 || block.chainid == 97) {
            lpMgr = _predict(
                _readSalt("SALT_LP_MANAGER"),
                abi.encodePacked(type(LPManager).creationCode, abi.encode(awpRegistry, poolManager, positionManager, permit2Addr, awp))
            );
            console.log("LPManager (Pancake):", lpMgr);
        } else {
            lpMgr = _predict(
                _readSalt("SALT_LP_MANAGER"),
                abi.encodePacked(type(LPManagerUni).creationCode, abi.encode(awpRegistry, poolManager, positionManager, permit2Addr, awp))
            );
            console.log("LPManager (Uni):   ", lpMgr);
        }

        // AWPEmission impl (no args)
        address emissionImpl = _predict(
            _readSalt("SALT_EMISSION_IMPL"),
            abi.encodePacked(type(AWPEmission).creationCode)
        );
        console.log("AWPEmission impl:  ", emissionImpl);

        // AWPEmission proxy (with guardian and GENESIS_TIME — matches Deploy.s.sol)
        uint256 genesisTime = vm.envUint("GENESIS_TIME");
        bytes memory initData = abi.encodeCall(AWPEmission.initialize, (awp, guardian, INITIAL_DAILY_EMISSION, genesisTime, EPOCH_DURATION, treasury));
        address emissionProxy = _predict(
            _readSalt("SALT_EMISSION_PROXY"),
            abi.encodePacked(type(ERC1967Proxy).creationCode, abi.encode(emissionImpl, initData))
        );
        console.log("AWPEmission proxy: ", emissionProxy);

        // StakingVault impl (no args)
        address vaultImpl = _predict(
            _readSalt("SALT_STAKING_VAULT_IMPL"),
            abi.encodePacked(type(StakingVault).creationCode)
        );
        console.log("StakingVault impl: ", vaultImpl);

        // StakingVault proxy (with awpRegistry and treasury — matches Deploy.s.sol)
        address vault;
        {
            bytes memory vaultInitData = abi.encodeCall(StakingVault.initialize, (awpRegistry, guardian));
            vault = _predict(
                _readSalt("SALT_STAKING_VAULT"),
                abi.encodePacked(type(ERC1967Proxy).creationCode, abi.encode(vaultImpl, vaultInitData))
            );
        }
        console.log("StakingVault proxy:", vault);

        // StakeNFT
        address stakeNft = _predict(
            _readSalt("SALT_STAKE_NFT"),
            abi.encodePacked(type(StakeNFT).creationCode, abi.encode(awp, vault, awpRegistry))
        );
        console.log("StakeNFT:          ", stakeNft);

        // WorknetManager impl (chain-conditional)
        address worknetMgrImpl;
        if (block.chainid == 56 || block.chainid == 97) {
            worknetMgrImpl = _predict(
                _readSalt("SALT_WORKNET_MANAGER_IMPL"),
                abi.encodePacked(type(WorknetManager).creationCode)
            );
            console.log("WorknetMgr (Pancake):", worknetMgrImpl);
        } else {
            worknetMgrImpl = _predict(
                _readSalt("SALT_WORKNET_MANAGER_IMPL"),
                abi.encodePacked(type(WorknetManagerUni).creationCode)
            );
            console.log("WorknetMgr (Uni):   ", worknetMgrImpl);
        }

        // AWPDAO (matches Deploy.s.sol params)
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
