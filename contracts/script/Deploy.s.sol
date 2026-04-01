// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AlphaTokenFactory} from "../src/token/AlphaTokenFactory.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {SubnetNFT} from "../src/core/SubnetNFT.sol";
import {LPManager} from "../src/core/LPManager.sol";
import {LPManagerUni} from "../src/core/LPManagerUni.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {TimelockController} from "@openzeppelin/contracts/governance/TimelockController.sol";
import {SubnetManager} from "../src/subnets/SubnetManager.sol";
import {SubnetManagerUni} from "../src/subnets/SubnetManagerUni.sol";

/// @title Deploy — Deterministic deployment via CREATE2 factory (0x4e59b44847b379578588920cA78FbF26c0B4956C)
/// @dev Account System V2: AccessManager removed, binding/delegation/recipient managed in AWPRegistry directly.
contract Deploy is Script {
    // ── Deterministic CREATE2 factory (deployed on all EVM chains) ──
    address constant DETERMINISTIC_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

    // ── Protocol constants ──
    uint256 constant INITIAL_DAILY_EMISSION = 15_800_000 * 1e18;
    uint256 constant EPOCH_DURATION = 1 days;
    uint256 constant TIMELOCK_DELAY = 172800; // 2 days

    /// @dev Deploy via the deterministic CREATE2 factory
    function _create2(bytes32 salt, bytes memory creationCode) internal returns (address deployed) {
        (bool success, bytes memory result) = DETERMINISTIC_DEPLOYER.call(abi.encodePacked(salt, creationCode));
        require(success && result.length == 20, "CREATE2 deploy failed");
        deployed = address(bytes20(result));
    }

    /// @dev Read a salt from env, default to bytes32(0) if not set
    function _readSalt(string memory key) internal view returns (bytes32) {
        return vm.envOr(key, bytes32(0));
    }

    function run() external {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);

        // Read addresses from environment
        address guardian = vm.envAddress("GUARDIAN");
        address liquidityPool = vm.envAddress("LIQUIDITY_POOL");
        address airdropAddr = vm.envAddress("AIRDROP");
        address poolManager = vm.envAddress("POOL_MANAGER");
        address positionManager = vm.envAddress("POSITION_MANAGER");
        address permit2Addr = vm.envAddress("PERMIT2");
        address clSwapRouter = vm.envAddress("CL_SWAP_ROUTER");
        address stateView = vm.envOr("STATE_VIEW", address(0));
        uint64 vanityRule = uint64(vm.envOr("VANITY_RULE", uint256(0)));

        // Read per-contract salts
        bytes32 saltAWPToken = _readSalt("SALT_AWP_TOKEN");
        bytes32 saltFactory = _readSalt("SALT_ALPHA_FACTORY");
        bytes32 saltTreasury = _readSalt("SALT_TREASURY");
        bytes32 saltAWPRegistryImpl = _readSalt("SALT_AWP_REGISTRY_IMPL");
        bytes32 saltAWPRegistry = _readSalt("SALT_AWP_REGISTRY");
        bytes32 saltSubnetNFT = _readSalt("SALT_SUBNET_NFT");
        bytes32 saltLPManager = _readSalt("SALT_LP_MANAGER");
        bytes32 saltEmissionImpl = _readSalt("SALT_EMISSION_IMPL");
        bytes32 saltEmissionProxy = _readSalt("SALT_EMISSION_PROXY");
        bytes32 saltVault = _readSalt("SALT_STAKING_VAULT");
        bytes32 saltStakeNFT = _readSalt("SALT_STAKE_NFT");
        bytes32 saltDAO = _readSalt("SALT_DAO");
        bytes32 saltSubnetMgrImpl = _readSalt("SALT_SUBNET_MANAGER_IMPL");

        console.log("Deployer:", deployer);

        vm.startBroadcast(deployerPrivateKey);

        // Step 1: AWPToken
        uint256 initialMint = vm.envOr("INITIAL_MINT", uint256(200_000_000)) * 1e18;
        AWPToken awp = AWPToken(_create2(
            saltAWPToken,
            abi.encodePacked(type(AWPToken).creationCode, abi.encode("AWP Token", "AWP", deployer, initialMint))
        ));
        console.log("AWPToken:", address(awp));

        // Step 2: AlphaTokenFactory
        AlphaTokenFactory factory = AlphaTokenFactory(_create2(
            saltFactory,
            abi.encodePacked(type(AlphaTokenFactory).creationCode, abi.encode(deployer, vanityRule))
        ));
        console.log("AlphaTokenFactory:", address(factory));

        // Step 3: Treasury
        Treasury treasury;
        {
            address[] memory proposers = new address[](0);
            address[] memory executors = new address[](1);
            executors[0] = address(0);
            treasury = Treasury(payable(_create2(
                saltTreasury,
                abi.encodePacked(type(Treasury).creationCode, abi.encode(TIMELOCK_DELAY, proposers, executors, deployer))
            )));
        }
        console.log("Treasury:", address(treasury));

        // Step 4: AWPRegistry (UUPS proxy)
        AWPRegistry awpRegistry;
        {
            AWPRegistry awpRegistryImpl = AWPRegistry(_create2(
                saltAWPRegistryImpl,
                abi.encodePacked(type(AWPRegistry).creationCode)
            ));
            console.log("AWPRegistry impl:", address(awpRegistryImpl));

            bytes memory registryInitData = abi.encodeCall(AWPRegistry.initialize, (deployer, address(treasury), guardian));
            awpRegistry = AWPRegistry(_create2(
                saltAWPRegistry,
                abi.encodePacked(type(ERC1967Proxy).creationCode, abi.encode(address(awpRegistryImpl), registryInitData))
            ));
        }
        console.log("AWPRegistry proxy:", address(awpRegistry));

        // Step 5: SubnetNFT
        SubnetNFT nft = SubnetNFT(_create2(
            saltSubnetNFT,
            abi.encodePacked(type(SubnetNFT).creationCode, abi.encode("AWP Subnet", "AWPSUB", address(awpRegistry)))
        ));
        console.log("SubnetNFT:", address(nft));

        // Step 6: LPManager (auto-select based on chain: BSC → PancakeSwap, other → Uniswap V4)
        address lpAddr;
        if (block.chainid == 56 || block.chainid == 97) {
            lpAddr = _create2(
                saltLPManager,
                abi.encodePacked(type(LPManager).creationCode, abi.encode(address(awpRegistry), poolManager, positionManager, permit2Addr, address(awp)))
            );
            console.log("LPManager (PancakeSwap):", lpAddr);
        } else {
            lpAddr = _create2(
                saltLPManager,
                abi.encodePacked(type(LPManagerUni).creationCode, abi.encode(address(awpRegistry), poolManager, positionManager, permit2Addr, address(awp)))
            );
            console.log("LPManager (Uniswap):", lpAddr);
        }

        // Step 7: AWPEmission (UUPS proxy)
        AWPEmission emission;
        {
            AWPEmission emissionImpl = AWPEmission(_create2(
                saltEmissionImpl,
                abi.encodePacked(type(AWPEmission).creationCode)
            ));
            console.log("AWPEmission impl:", address(emissionImpl));

            uint256 genesisTime = vm.envOr("GENESIS_TIME", block.timestamp);
            bytes memory initData = abi.encodeCall(AWPEmission.initialize, (address(awp), deployer, INITIAL_DAILY_EMISSION, genesisTime, EPOCH_DURATION));
            emission = AWPEmission(_create2(
                saltEmissionProxy,
                abi.encodePacked(type(ERC1967Proxy).creationCode, abi.encode(address(emissionImpl), initData))
            ));
        }
        console.log("AWPEmission proxy:", address(emission));

        // Step 8: StakingVault (UUPS proxy) + StakeNFT
        StakingVault vault;
        {
            bytes32 saltVaultImpl = _readSalt("SALT_STAKING_VAULT_IMPL");
            StakingVault vaultImpl = StakingVault(_create2(
                saltVaultImpl,
                abi.encodePacked(type(StakingVault).creationCode)
            ));
            console.log("StakingVault impl:", address(vaultImpl));

            bytes memory vaultInitData = abi.encodeCall(StakingVault.initialize, (address(awpRegistry), address(treasury)));
            vault = StakingVault(_create2(
                saltVault,
                abi.encodePacked(type(ERC1967Proxy).creationCode, abi.encode(address(vaultImpl), vaultInitData))
            ));
        }
        console.log("StakingVault proxy:", address(vault));

        StakeNFT stakeNft = StakeNFT(_create2(
            saltStakeNFT,
            abi.encodePacked(type(StakeNFT).creationCode, abi.encode(address(awp), address(vault), address(awpRegistry)))
        ));
        console.log("StakeNFT:", address(stakeNft));

        // Step 9: SubnetManager implementation (auto-select based on chain)
        // BSC (56/97) uses PancakeSwap V4 SubnetManager; other chains use Uniswap V4 SubnetManagerUni
        SubnetManager subnetMgrImpl;
        if (block.chainid == 56 || block.chainid == 97) {
            subnetMgrImpl = SubnetManager(_create2(
                saltSubnetMgrImpl,
                abi.encodePacked(type(SubnetManager).creationCode)
            ));
            console.log("SubnetManager impl (PancakeSwap):", address(subnetMgrImpl));
        } else {
            subnetMgrImpl = SubnetManager(address(SubnetManagerUni(_create2(
                saltSubnetMgrImpl,
                abi.encodePacked(type(SubnetManagerUni).creationCode)
            ))));
            console.log("SubnetManager impl (Uniswap):", address(subnetMgrImpl));
        }

        // Step 10: AWPDAO
        AWPDAO dao;
        dao = AWPDAO(payable(_create2(
            saltDAO,
            abi.encodePacked(type(AWPDAO).creationCode, abi.encode(
                address(stakeNft), address(awp),
                address(treasury),
                uint48(7200), uint32(50400), uint256(4)
            ))
        )));
        console.log("AWPDAO:", address(dao));

        // Step 11: Roles + Admin
        treasury.grantRole(treasury.PROPOSER_ROLE(), address(dao));
        treasury.grantRole(treasury.CANCELLER_ROLE(), address(dao));
        treasury.renounceRole(treasury.DEFAULT_ADMIN_ROLE(), deployer);
        // NOTE: EXECUTOR_ROLE is granted to address(0) (open execution) by design.
        // Anyone can execute queued proposals after the timelock delay.
        // Deployer was never granted EXECUTOR_ROLE, so no renounce needed.
        console.log("Roles granted + Treasury admin renounced (open executor by design)");

        awp.addMinter(address(emission));
        awp.renounceAdmin();
        console.log("AWP minter set + admin locked");

        factory.setAddresses(address(awpRegistry));
        console.log("Factory configured");

        // Step 12: Initialize registry (no accessManager)
        uint24 POOL_FEE = 10000;
        int24 TICK_SPACING = 200;
        bytes memory dexCfg;
        if (block.chainid == 56 || block.chainid == 97) {
            // PancakeSwap V4: 6-field dexConfig
            dexCfg = abi.encode(poolManager, positionManager, clSwapRouter, permit2Addr, uint24(POOL_FEE), int24(TICK_SPACING));
        } else {
            // Uniswap V4: 7-field dexConfig (extra: stateView)
            dexCfg = abi.encode(poolManager, positionManager, clSwapRouter, permit2Addr, uint24(POOL_FEE), int24(TICK_SPACING), stateView);
        }
        awpRegistry.initializeRegistry(
            address(awp), address(nft), address(factory), address(emission),
            lpAddr, address(vault), address(stakeNft),
            address(subnetMgrImpl), dexCfg
        );
        console.log("Registry initialized");

        // Step 13: Token distribution
        uint256 treasuryAmount = vm.envOr("DIST_TREASURY", uint256(90_000_000)) * 1e18;
        uint256 liquidityAmount = vm.envOr("DIST_LIQUIDITY", uint256(10_000_000)) * 1e18;
        uint256 airdropAmount = vm.envOr("DIST_AIRDROP", uint256(100_000_000)) * 1e18;
        if (treasuryAmount > 0) awp.transfer(address(treasury), treasuryAmount);
        if (liquidityAmount > 0) awp.transfer(liquidityPool, liquidityAmount);
        if (airdropAmount > 0) awp.transfer(airdropAddr, airdropAmount);
        console.log("Tokens distributed");

        vm.stopBroadcast();

        // Verification
        require(awp.admin() == address(0), "Admin should be renounced");
        require(awp.minters(address(emission)), "Emission should be minter");
        require(awpRegistry.registryInitialized(), "Registry should be initialized");
        require(awp.balanceOf(deployer) == initialMint - treasuryAmount - liquidityAmount - airdropAmount, "Deployer balance mismatch after distribution");

        console.log("=== Deployment Complete ===");
    }
}
