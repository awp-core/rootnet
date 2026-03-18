// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IStakeNFT — Interface for NFT-based staking positions
/// @notice Each stake position is an ERC721 token with locked amount and expiry.
///         Voting power is derived from NFT positions using amount * sqrt(min(remainingTime, MAX_WEIGHT_DURATION) / VOTE_WEIGHT_DIVISOR).
interface IStakeNFT {
    /// @notice Stake position data stored per tokenId
    struct Position {
        uint128 amount;       // Staked AWP amount
        uint64 lockEndTime;   // Unlock timestamp (withdrawable when block.timestamp >= lockEndTime)
        uint64 createdAt;     // Timestamp when position was minted (for proposal eligibility)
    }

    // ── Events ──

    event Deposited(address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndTime);
    event PositionIncreased(uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndTime);
    event Withdrawn(address indexed user, uint256 indexed tokenId, uint256 amount);

    // ── Write functions ──

    /// @notice Create a new stake position by depositing AWP with a lock duration
    /// @param amount AWP amount to deposit
    /// @param lockDuration Lock duration in seconds (minimum MIN_LOCK_DURATION)
    /// @return tokenId The minted NFT token ID
    function deposit(uint256 amount, uint64 lockDuration) external returns (uint256 tokenId);

    /// @notice Create a new stake position on behalf of a user (called by RootNet for registerAndStake)
    /// @dev AWP is transferred from the user (not msg.sender). Only callable by RootNet.
    /// @param user User address to mint the position to
    /// @param amount AWP amount to deposit
    /// @param lockDuration Lock duration in seconds
    /// @return tokenId The minted NFT token ID
    function depositFor(address user, uint256 amount, uint64 lockDuration) external returns (uint256 tokenId);

    /// @notice Add AWP and/or extend lock on an existing position
    /// @param tokenId Position token ID (must be owned by msg.sender)
    /// @param amount Additional AWP to add (0 to skip)
    /// @param newLockEndTime New lock end timestamp (0 to skip; must be >= current lockEndTime and > block.timestamp)
    function addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime) external;

    /// @notice Withdraw a position after lock expiry (burns the NFT)
    /// @param tokenId Position token ID (must be owned by msg.sender, lock must be expired)
    function withdraw(uint256 tokenId) external;

    // ── View functions ──

    /// @notice Get position data for a given token ID
    /// @param tokenId Position token ID
    /// @return amount Staked AWP amount
    /// @return lockEndTime Unlock timestamp
    /// @return createdAt Creation timestamp
    function positions(uint256 tokenId) external view returns (uint128 amount, uint64 lockEndTime, uint64 createdAt);

    /// @notice Get the total staked AWP for a user (O(1) via accumulator)
    /// @param user User address
    /// @return Total staked AWP amount
    function getUserTotalStaked(address user) external view returns (uint256);

    /// @notice Get the voting power of a specific position
    /// @param tokenId Position token ID
    /// @return Voting power = amount * sqrt(min(remainingTime, MAX_WEIGHT_DURATION) / VOTE_WEIGHT_DIVISOR)
    function getVotingPower(uint256 tokenId) external view returns (uint256);

    /// @notice Get the total voting power for a user given their tokenIds
    /// @param user User address (ownership verified per token)
    /// @param tokenIds Array of tokenIds owned by user
    /// @return Total voting power across the specified positions
    function getUserVotingPower(address user, uint256[] calldata tokenIds) external view returns (uint256);

    /// @notice Get the approximate total voting power (upper bound for quorum)
    /// @return totalStakedAWP * SQRT_MAX_WEIGHT_FACTOR
    function totalVotingPower() external view returns (uint256);

    /// @notice Get remaining time until unlock for a position
    /// @param tokenId Position token ID
    /// @return Remaining seconds (0 if expired)
    function remainingTime(uint256 tokenId) external view returns (uint64);

    /// @notice Get full position info for voting in a single call (reduces 3 external calls to 1)
    /// @param tokenId Position token ID
    /// @return owner Token owner address
    /// @return amount Staked AWP amount
    /// @return lockEndTime Unlock timestamp
    /// @return createdAt Creation timestamp
    /// @return remainingSeconds Remaining lock seconds
    /// @return votingPower Voting power for this position
    function getPositionForVoting(uint256 tokenId) external view returns (
        address owner,
        uint128 amount,
        uint64 lockEndTime,
        uint64 createdAt,
        uint64 remainingSeconds,
        uint256 votingPower
    );
}
