// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @title MockLPManager — Minimal LP manager for tests (no real DEX)
contract MockLPManager {
    address public immutable awpRegistry;
    address public immutable awpToken;

    uint256 private _nextPoolCounter = 1;

    error NotAWPRegistry();

    constructor(address awpRegistry_, address awpToken_) {
        awpRegistry = awpRegistry_;
        awpToken = awpToken_;
    }

    function createPoolAndAddLiquidity(address, uint256, uint256)
        external
        returns (bytes32 poolId, uint256 lpTokenId)
    {
        if (msg.sender != awpRegistry) revert NotAWPRegistry();
        poolId = bytes32(_nextPoolCounter++);
        lpTokenId = uint256(poolId);
    }

    function worknetTokenToPoolId(address) external pure returns (bytes32) { return bytes32(0); }
    function worknetTokenToTokenId(address) external pure returns (uint256) { return 0; }
    function needsCompounding(address) external pure returns (bool, uint256) { return (false, 0); }
}
