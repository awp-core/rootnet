// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title MockWorknetManager — Minimal UUPS impl for tests
contract MockWorknetManager is UUPSUpgradeable {
    address public alphaToken;
    bytes32 public poolId;
    address public admin;

    constructor() { _disableInitializers(); }

    function initialize(address alphaToken_, bytes32 poolId_, address admin_) external initializer {
        __UUPSUpgradeable_init();
        alphaToken = alphaToken_;
        poolId = poolId_;
        admin = admin_;
    }

    function _authorizeUpgrade(address) internal pure override {}
}
