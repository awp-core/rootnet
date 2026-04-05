// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title LPManagerStub — Minimal UUPS impl for proxy initialization only.
/// @dev Deploy once, use as initial impl for ERC1967Proxy on all chains.
///      Only deployer can upgrade to chain-specific LPManager/LPManagerUni impl.
///      deployer is immutable (in bytecode, not storage) — no slot collision after upgrade.
contract LPManagerStub is Initializable, UUPSUpgradeable {
    address public constant awpRegistry = 0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A;

    /// @dev Deployer baked into bytecode — zero storage footprint on the proxy
    address public immutable deployer;

    error NotDeployer();

    constructor(address deployer_) {
        deployer = deployer_;
        _disableInitializers();
    }

    function initialize() external initializer {}

    /// @dev Only deployer can upgrade — prevents front-running between deploy and upgrade.
    function _authorizeUpgrade(address) internal view override {
        if (msg.sender != deployer) revert NotDeployer();
    }
}
