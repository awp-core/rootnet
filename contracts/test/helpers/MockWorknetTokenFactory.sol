// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {WorknetToken} from "../../src/token/WorknetToken.sol";

/// @title MockWorknetTokenFactory — Deploys WorknetToken for unit tests
/// @dev Implements the pending* getters that WorknetToken reads during construction.
contract MockWorknetTokenFactory {
    string public pendingName;
    string public pendingSymbol;
    uint256 public pendingWorknetId;

    function deploy(string memory name_, string memory symbol_, uint256 worknetId_) external returns (WorknetToken) {
        pendingName = name_;
        pendingSymbol = symbol_;
        pendingWorknetId = worknetId_;
        WorknetToken token = new WorknetToken();
        delete pendingName;
        delete pendingSymbol;
        delete pendingWorknetId;
        return token;
    }
}
