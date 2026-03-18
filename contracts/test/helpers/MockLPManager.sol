// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @dev Mock LPManager for non-fork tests, does not depend on PancakeSwap
contract MockLPManager {
    address public immutable rootNet;
    IERC20 public immutable awpToken;
    mapping(address => bytes32) public alphaTokenToPool;

    error NotRootNet();
    error PoolAlreadyExists();

    modifier onlyRootNet() {
        if (msg.sender != rootNet) revert NotRootNet();
        _;
    }

    constructor(address rootNet_, address awpToken_) {
        rootNet = rootNet_;
        awpToken = IERC20(awpToken_);
    }

    function createPoolAndAddLiquidity(address alphaToken, uint256, uint256)
        external
        onlyRootNet
        returns (bytes32 poolId, uint256 lpTokenId)
    {
        if (alphaTokenToPool[alphaToken] != bytes32(0)) revert PoolAlreadyExists();
        poolId = keccak256(abi.encodePacked(alphaToken, address(awpToken)));
        alphaTokenToPool[alphaToken] = poolId;
        lpTokenId = 0;
    }
}
