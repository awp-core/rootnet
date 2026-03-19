// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @dev Mock LPManager for non-fork tests, does not depend on PancakeSwap
contract MockLPManager {
    address public immutable awpRegistry;
    IERC20 public immutable awpToken;
    mapping(address => bytes32) public alphaTokenToPool;

    error NotAWPRegistry();
    error PoolAlreadyExists();

    modifier onlyAWPRegistry() {
        if (msg.sender != awpRegistry) revert NotAWPRegistry();
        _;
    }

    constructor(address awpRegistry_, address awpToken_) {
        awpRegistry = awpRegistry_;
        awpToken = IERC20(awpToken_);
    }

    function createPoolAndAddLiquidity(address alphaToken, uint256, uint256)
        external
        onlyAWPRegistry
        returns (bytes32 poolId, uint256 lpTokenId)
    {
        if (alphaTokenToPool[alphaToken] != bytes32(0)) revert PoolAlreadyExists();
        poolId = keccak256(abi.encodePacked(alphaToken, address(awpToken)));
        alphaTokenToPool[alphaToken] = poolId;
        lpTokenId = 0;
    }
}
