// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {EIP712Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import {ERC1967Utils} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Utils.sol";

/// @title MigrateAllocatorEIP712 — One-shot migration to fix empty EIP-712 domain
/// @dev Deploy → Guardian calls proxy.upgradeToAndCall(this, migrate(originalImpl))
///      migrate() writes EIP-712 storage then swaps impl back. Throwaway contract.
contract MigrateAllocatorEIP712 is Initializable, UUPSUpgradeable, EIP712Upgradeable {
    address public guardian;

    constructor() { _disableInitializers(); }

    function migrate(address originalImpl) external reinitializer(2) {
        __EIP712_init("AWPAllocator", "1");
        ERC1967Utils.upgradeToAndCall(originalImpl, "");
    }

    function _authorizeUpgrade(address) internal view override {
        require(msg.sender == guardian);
    }
}
