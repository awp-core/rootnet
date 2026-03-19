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
import {TimelockController} from "@openzeppelin/contracts/governance/TimelockController.sol";
import {SubnetManager} from "../src/subnets/SubnetManager.sol";

/// @title Deploy — Deterministic deployment via CREATE2 factory (0x4e59b44847b379578588920cA78FbF26c0B4956C)
/// @dev Each contract has its own salt read from .env (e.g. SALT_AWP_TOKEN, SALT_ROOTNET, etc.)
///      for vanity address mining. Use the companion predict script to find salts off-chain.
///      Same salt + same initcode = same address on any EVM chain.
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
        uint64 vanityRule = uint64(vm.envOr("VANITY_RULE", uint256(0)));

        // Read per-contract salts (pre-mined off-chain for vanity addresses)
        bytes32 saltAWPToken = _readSalt("SALT_AWP_TOKEN");
        bytes32 saltFactory = _readSalt("SALT_ALPHA_FACTORY");
        bytes32 saltTreasury = _readSalt("SALT_TREASURY");
        bytes32 saltRootNet = _readSalt("SALT_ROOTNET");
        bytes32 saltSubnetNFT = _readSalt("SALT_SUBNET_NFT");
        bytes32 saltAccessMgr = _readSalt("SALT_ACCESS_MANAGER");
        bytes32 saltLPManager = _readSalt("SALT_LP_MANAGER");
        bytes32 saltEmissionImpl = _readSalt("SALT_EMISSION_IMPL");
        bytes32 saltEmissionProxy = _readSalt("SALT_EMISSION_PROXY");
        bytes32 saltVault = _readSalt("SALT_STAKING_VAULT");
        bytes32 saltStakeNFT = _readSalt("SALT_STAKE_NFT");
        bytes32 saltDAO = _readSalt("SALT_DAO");
        bytes32 saltSubnetMgrImpl = _readSalt("SALT_SUBNET_MANAGER_IMPL");

        console.log("Deployer:", deployer);

        vm.startBroadcast(deployerPrivateKey);

        // ═══════════════════════════════════════════════
        //  Step 1: AWPToken
        // ═══════════════════════════════════════════════
        AWPToken awp = AWPToken(_create2(
            saltAWPToken,
            abi.encodePacked(type(AWPToken).creationCode, abi.encode("AWP Token", "AWP", deployer))
        ));
        console.log("AWPToken:", address(awp));

        // ═══════════════════════════════════════════════
        //  Step 2: AlphaTokenFactory
        // ═══════════════════════════════════════════════
        AlphaTokenFactory factory = AlphaTokenFactory(_create2(
            saltFactory,
            abi.encodePacked(type(AlphaTokenFactory).creationCode, abi.encode(deployer, vanityRule))
        ));
        console.log("AlphaTokenFactory:", address(factory));

        // ═══════════════════════════════════════════════
        //  Step 3: Treasury
        // ═══════════════════════════════════════════════
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

        // ═══════════════════════════════════════════════
        //  Step 4: RootNet
        // ═══════════════════════════════════════════════
        RootNet rootNet = RootNet(_create2(
            saltRootNet,
            abi.encodePacked(type(RootNet).creationCode, abi.encode(deployer, address(treasury), guardian))
        ));
        console.log("RootNet:", address(rootNet));

        // ═══════════════════════════════════════════════
        //  Step 5: SubnetNFT
        // ═══════════════════════════════════════════════
        SubnetNFT nft = SubnetNFT(_create2(
            saltSubnetNFT,
            abi.encodePacked(type(SubnetNFT).creationCode, abi.encode("AWP Subnet", "AWPSUB", address(rootNet)))
        ));
        console.log("SubnetNFT:", address(nft));

        // ═══════════════════════════════════════════════
        //  Step 6: AccessManager
        // ═══════════════════════════════════════════════
        AccessManager access = AccessManager(_create2(
            saltAccessMgr,
            abi.encodePacked(type(AccessManager).creationCode, abi.encode(address(rootNet)))
        ));
        console.log("AccessManager:", address(access));

        // ═══════════════════════════════════════════════
        //  Step 7: LPManager
        // ═══════════════════════════════════════════════
        LPManager lp = LPManager(_create2(
            saltLPManager,
            abi.encodePacked(type(LPManager).creationCode, abi.encode(address(rootNet), poolManager, positionManager, permit2Addr, address(awp)))
        ));
        console.log("LPManager:", address(lp));

        // ═══════════════════════════════════════════════
        //  Step 8: AWPEmission (UUPS proxy)
        // ═══════════════════════════════════════════════
        AWPEmission emission;
        {
            AWPEmission emissionImpl = AWPEmission(_create2(
                saltEmissionImpl,
                abi.encodePacked(type(AWPEmission).creationCode)
            ));
            console.log("AWPEmission impl:", address(emissionImpl));

            bytes memory initData = abi.encodeCall(AWPEmission.initialize, (address(awp), address(treasury), INITIAL_DAILY_EMISSION, block.timestamp, EPOCH_DURATION));
            emission = AWPEmission(_create2(
                saltEmissionProxy,
                abi.encodePacked(type(ERC1967Proxy).creationCode, abi.encode(address(emissionImpl), initData))
            ));
        }
        console.log("AWPEmission proxy:", address(emission));

        // ═══════════════════════════════════════════════
        //  Step 9: StakingVault + StakeNFT
        // ═══════════════════════════════════════════════
        //   StakingVault.stakeNFT is set via setStakeNFT() (called in initializeRegistry)
        //   This breaks the circular CREATE2 dependency — both can have vanity addresses.
        StakingVault vault = StakingVault(_create2(
            saltVault,
            abi.encodePacked(type(StakingVault).creationCode, abi.encode(address(rootNet)))
        ));
        console.log("StakingVault:", address(vault));

        StakeNFT stakeNft = StakeNFT(_create2(
            saltStakeNFT,
            abi.encodePacked(type(StakeNFT).creationCode, abi.encode(address(awp), address(vault), address(rootNet)))
        ));
        console.log("StakeNFT:", address(stakeNft));

        // ═══════════════════════════════════════════════
        //  Step 10: SubnetManager implementation (shared by all auto-deployed subnet proxies)
        // ═══════════════════════════════════════════════
        SubnetManager subnetMgrImpl = SubnetManager(_create2(
            saltSubnetMgrImpl,
            abi.encodePacked(type(SubnetManager).creationCode)
        ));
        console.log("SubnetManager impl:", address(subnetMgrImpl));

        // ═══════════════════════════════════════════════
        //  Step 11: AWPDAO
        // ═══════════════════════════════════════════════
        AWPDAO dao; // scoped to avoid stack-too-deep
        dao = AWPDAO(payable(_create2(
            saltDAO,
            abi.encodePacked(type(AWPDAO).creationCode, abi.encode(
                address(stakeNft), address(awp),
                address(treasury),
                uint48(7200), uint32(50400), uint256(4)
            ))
        )));
        console.log("AWPDAO:", address(dao));

        // ═══════════════════════════════════════════════
        //  Step 12: Roles + Admin
        // ═══════════════════════════════════════════════
        treasury.grantRole(treasury.PROPOSER_ROLE(), address(dao));
        treasury.grantRole(treasury.CANCELLER_ROLE(), address(dao));
        treasury.renounceRole(treasury.DEFAULT_ADMIN_ROLE(), deployer);
        console.log("Roles granted + Treasury admin renounced");

        awp.addMinter(address(emission));
        awp.renounceAdmin();
        console.log("AWP minter set + admin locked");

        factory.setAddresses(address(rootNet));
        console.log("Factory configured");

        // ═══════════════════════════════════════════════
        //  Step 13: Initialize registry
        // ═══════════════════════════════════════════════
        rootNet.initializeRegistry(
            address(awp), address(nft), address(factory), address(emission),
            address(lp), address(access), address(vault), address(stakeNft),
            address(subnetMgrImpl)
        );
        console.log("Registry initialized");

        // ═══════════════════════════════════════════════
        //  Step 14: Token distribution
        // ═══════════════════════════════════════════════
        awp.transfer(address(treasury), 90_000_000 * 1e18);
        awp.transfer(liquidityPool, 10_000_000 * 1e18);
        awp.transfer(airdropAddr, 100_000_000 * 1e18);
        console.log("Tokens distributed");

        vm.stopBroadcast();

        // ═══════════════════════════════════════════════
        //  Verification
        // ═══════════════════════════════════════════════
        require(awp.admin() == address(0), "Admin should be renounced");
        require(awp.minters(address(emission)), "Emission should be minter");
        require(rootNet.registryInitialized(), "Registry should be initialized");
        // Note: deployer balance check removed — distribution addresses may equal deployer in test/staging

        console.log("=== Deployment Complete ===");
    }
}
