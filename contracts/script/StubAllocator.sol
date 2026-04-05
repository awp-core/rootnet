// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @dev Minimal stub to occupy the proxy slot. Same initialize(address,address) selector as AWPAllocator.
contract StubAllocator is Initializable, UUPSUpgradeable {
    address public guardian;
    constructor() { _disableInitializers(); }
    function initialize(address, address guardian_) external initializer {
        guardian = guardian_;
    }
    function _authorizeUpgrade(address) internal view override {
        require(msg.sender == guardian);
    }
}
