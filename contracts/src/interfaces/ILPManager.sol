// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title ILPManager — Liquidity management interface
/// @notice Wraps PancakeSwap V4 LP creation; LP is permanently locked
interface ILPManager {
    /// @notice Create an LP pool and add full-range liquidity (called once during worknet registration)
    /// @param worknetToken Alpha token address
    /// @param awpAmount AWP amount
    /// @param alphaAmount Alpha amount
    /// @return poolId LP pool ID (bytes32)
    /// @return lpTokenId LP NFT ID (LP permanently locked in LPManager)
    function createPoolAndAddLiquidity(address worknetToken, uint256 awpAmount, uint256 alphaAmount)
        external
        returns (bytes32 poolId, uint256 lpTokenId);
}
