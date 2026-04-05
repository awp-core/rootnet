// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title MockWorknetManager — Minimal UUPS impl for tests
contract MockWorknetManager is Initializable, UUPSUpgradeable {
    address public worknetToken;
    bytes32 public poolId;
    address public admin;

    constructor() { _disableInitializers(); }

    function initialize(address worknetToken_, bytes32 poolId_, address admin_) external initializer {
        worknetToken = worknetToken_;
        poolId = poolId_;
        admin = admin_;
    }

    function _authorizeUpgrade(address) internal pure override {}
}
