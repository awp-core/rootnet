// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {LPManagerStub} from "../src/core/LPManagerStub.sol";
import {LPManagerUni} from "../src/core/LPManagerUni.sol";
import {LPManager} from "../src/core/LPManager.sol";
import {WorknetManager} from "../src/worknets/WorknetManager.sol";
import {WorknetManagerUni} from "../src/worknets/WorknetManagerUni.sol";

/// @title DeployLP - Deploy LPManager (stub+proxy+impl) and WorknetManager impl
/// @dev All deployments use CREATE2 via deterministic factory.
///      Phase 1: deploy stub + proxy (same on all chains)
///      Phase 2: deploy chain-specific impl + upgradeToAndCall
///      Phase 3: deploy WorknetManager impl
contract DeployLP is Script {
    address constant FACTORY = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

    function _create2(bytes32 salt, bytes memory code) internal returns (address deployed) {
        (bool ok, bytes memory res) = FACTORY.call(abi.encodePacked(salt, code));
        require(ok && res.length == 20, "CREATE2 failed");
        deployed = address(bytes20(res));
    }

    function _salt(string memory key) internal view returns (bytes32) {
        return vm.envBytes32(key);
    }

    function run() external {
        uint256 pk = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(pk);

        address permit2Addr = vm.envAddress("PERMIT2");
        address poolManager = vm.envAddress("POOL_MANAGER");
        address positionManager = vm.envAddress("POSITION_MANAGER");
        address clSwapRouter = vm.envAddress("CL_SWAP_ROUTER");
        address stateView;
        {
            string memory sv = vm.envOr("STATE_VIEW", string(""));
            if (bytes(sv).length > 0) stateView = vm.parseAddress(sv);
        }

        vm.startBroadcast(pk);

        // 1. Deploy stub (same on all chains — deployer is immutable in bytecode)
        address stub = _create2(
            _salt("SALT_LP_STUB"),
            abi.encodePacked(type(LPManagerStub).creationCode, abi.encode(deployer))
        );
        console.log("LPManagerStub:", stub);

        // 2. Deploy proxy (same on all chains)
        bytes memory initData = abi.encodeCall(LPManagerStub.initialize, ());
        address proxy = _create2(
            _salt("SALT_LP_PROXY"),
            abi.encodePacked(type(ERC1967Proxy).creationCode, abi.encode(stub, initData))
        );
        console.log("LPManager proxy:", proxy);

        // 3. Deploy chain-specific impl + upgrade
        if (block.chainid == 56 || block.chainid == 97) {
            address impl = _create2(
                _salt("SALT_LP_IMPL"),
                abi.encodePacked(type(LPManager).creationCode, abi.encode(permit2Addr, poolManager, positionManager))
            );
            console.log("LPManager impl (PCS):", impl);
            LPManagerStub(proxy).upgradeToAndCall(impl, "");
            console.log("Upgraded to PCS impl");
        } else {
            address impl = _create2(
                _salt("SALT_LP_IMPL"),
                abi.encodePacked(type(LPManagerUni).creationCode, abi.encode(permit2Addr, poolManager, positionManager))
            );
            console.log("LPManager impl (Uni):", impl);
            LPManagerStub(proxy).upgradeToAndCall(impl, "");
            console.log("Upgraded to Uni impl");
        }

        // 4. Deploy WorknetManager impl
        if (block.chainid == 56 || block.chainid == 97) {
            address wm = _create2(
                _salt("SALT_WM_IMPL"),
                abi.encodePacked(
                    type(WorknetManager).creationCode,
                    abi.encode(permit2Addr, poolManager, positionManager, clSwapRouter)
                )
            );
            console.log("WorknetManager impl (PCS):", wm);
        } else {
            address wm = _create2(
                _salt("SALT_WM_IMPL"),
                abi.encodePacked(
                    type(WorknetManagerUni).creationCode,
                    abi.encode(permit2Addr, poolManager, positionManager, clSwapRouter, stateView)
                )
            );
            console.log("WorknetManager impl (Uni):", wm);
        }

        vm.stopBroadcast();
        console.log("=== Done ===");
    }
}
