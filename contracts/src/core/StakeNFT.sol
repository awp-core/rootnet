// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {IStakeNFT} from "../interfaces/IStakeNFT.sol";
import {IStakingVault} from "../interfaces/IStakingVault.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/// @title StakeNFT — ERC721 NFT-based staking positions
/// @notice Each stake position is an ERC721 token with locked AWP amount and expiry.
///         Voting power is derived from positions using amount * sqrt(min(remainingTime, MAX_WEIGHT_DURATION) / VOTE_WEIGHT_DIVISOR).
/// @dev Key design: _userTotalStaked accumulator provides O(1) balance checks for allocation validation.
///      Transfer hook (_update override) maintains accumulator consistency and checks allocation coverage.
///      getUserVotingPower requires caller to pass tokenIds (no on-chain enumeration).
contract StakeNFT is ERC721, ReentrancyGuard, IStakeNFT {
    using SafeERC20 for IERC20;

    // ── Immutables ──

    /// @notice AWP token contract
    IERC20 public immutable awpToken;

    /// @notice StakingVault contract (for allocation balance checks)
    address public immutable stakingVault;

    /// @notice RootNet contract (for depositFor access control)
    address public immutable rootNet;

    // ── Constants ──

    /// @notice Maximum duration for voting power weight calculation (54 weeks)
    uint64 public constant MAX_WEIGHT_DURATION = 54 * 7 days;

    /// @notice Precomputed sqrt(MAX_WEIGHT_DURATION / VOTE_WEIGHT_DIVISOR) = sqrt(54) = 7
    uint256 public constant SQRT_MAX_WEIGHT_FACTOR = 7;

    /// @notice Minimum lock duration in seconds (1 day)
    uint64 public constant MIN_LOCK_DURATION = 1 days;

    /// @notice Divisor for voting power formula: amount * sqrt(remainingTime / VOTE_WEIGHT_DIVISOR)
    uint256 public constant VOTE_WEIGHT_DIVISOR = 7 days;

    // ── Storage ──

    /// @notice Position data per tokenId
    mapping(uint256 => Position) public positions;

    /// @dev Auto-incrementing token ID counter, starts at 1
    uint256 private _nextTokenId = 1;

    /// @notice O(1) balance tracking for allocation validation
    mapping(address => uint256) private _userTotalStaked;

    // ── Errors ──

    error InvalidAmount();
    error LockTooShort();
    error NotTokenOwner();
    error LockNotExpired();
    error PositionExpired();
    error InsufficientUnallocated();
    error NothingToUpdate();
    error LockCannotShorten();
    error LockMustExceedCurrentTime();
    error NotRootNet();

    // ── Modifiers ──

    modifier onlyRootNet() {
        if (msg.sender != rootNet) revert NotRootNet();
        _;
    }

    // ── Constructor ──

    /// @param awpToken_ AWP token address
    /// @param stakingVault_ StakingVault contract address
    /// @param rootNet_ RootNet contract address (for depositFor access control)
    constructor(
        address awpToken_,
        address stakingVault_,
        address rootNet_
    ) ERC721("AWP Stake Position", "sAWP") {
        awpToken = IERC20(awpToken_);
        stakingVault = stakingVault_;
        rootNet = rootNet_;
    }

    // ══════════════════════════════════════════════
    //  Write functions
    // ══════════════════════════════════════════════

    /// @inheritdoc IStakeNFT
    function deposit(uint256 amount, uint64 lockDuration) external returns (uint256 tokenId) {
        return _deposit(msg.sender, msg.sender, amount, lockDuration);
    }

    /// @notice Gasless deposit: user signs ERC-2612 permit off-chain, no prior approve tx needed
    /// @param amount AWP amount to deposit
    /// @param lockDuration Lock duration in seconds
    /// @param deadline Permit expiry timestamp
    /// @param v Permit signature v
    /// @param r Permit signature r
    /// @param s Permit signature s
    function depositWithPermit(
        uint256 amount, uint64 lockDuration,
        uint256 deadline, uint8 v, bytes32 r, bytes32 s
    ) external returns (uint256 tokenId) {
        IERC20Permit(address(awpToken)).permit(msg.sender, address(this), amount, deadline, v, r, s);
        return _deposit(msg.sender, msg.sender, amount, lockDuration);
    }

    /// @inheritdoc IStakeNFT
    function depositFor(address user, uint256 amount, uint64 lockDuration)
        external
        onlyRootNet
        returns (uint256 tokenId)
    {
        // AWP is transferred from the user (not RootNet)
        return _deposit(user, user, amount, lockDuration);
    }

    /// @dev Internal deposit logic: transfer AWP from `from`, mint NFT to `to`
    /// @param from Address to transfer AWP from
    /// @param to Address to mint NFT to
    /// @param amount AWP amount
    /// @param lockDuration Lock duration in seconds
    function _deposit(address from, address to, uint256 amount, uint64 lockDuration) internal returns (uint256 tokenId) {
        if (amount == 0) revert InvalidAmount();
        if (amount > type(uint128).max) revert InvalidAmount();
        if (lockDuration < MIN_LOCK_DURATION) revert LockTooShort();

        uint64 lockEndTime = uint64(block.timestamp) + lockDuration;

        // Transfer AWP from user to this contract
        awpToken.safeTransferFrom(from, address(this), amount);

        // Store position BEFORE mint so _update callback sees correct amount
        tokenId = _nextTokenId++;
        positions[tokenId] = Position({
            amount: uint128(amount),
            lockEndTime: lockEndTime,
            createdAt: uint64(block.timestamp)
        });

        // Mint NFT — _update handles _userTotalStaked accumulator
        _mint(to, tokenId);

        emit Deposited(to, tokenId, amount, lockEndTime);
    }

    /// @inheritdoc IStakeNFT
    function addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime) external {
        if (ownerOf(tokenId) != msg.sender) revert NotTokenOwner();

        Position storage pos = positions[tokenId];
        bool updated = false;

        if (amount > 0) {
            // Position lock must still be active to add tokens (prevents bypassing MIN_LOCK_DURATION)
            if (pos.lockEndTime <= uint64(block.timestamp)) revert PositionExpired();
            if (uint256(pos.amount) + amount > type(uint128).max) revert InvalidAmount();
            awpToken.safeTransferFrom(msg.sender, address(this), amount);
            pos.amount += uint128(amount);
            _userTotalStaked[msg.sender] += amount;
            updated = true;
        }

        if (newLockEndTime > 0) {
            if (newLockEndTime < pos.lockEndTime) revert LockCannotShorten();
            if (newLockEndTime <= uint64(block.timestamp)) revert LockMustExceedCurrentTime();
            pos.lockEndTime = newLockEndTime;
            updated = true;
        }

        if (!updated) revert NothingToUpdate();

        emit PositionIncreased(tokenId, amount, pos.lockEndTime);
    }

    /// @inheritdoc IStakeNFT
    function withdraw(uint256 tokenId) external nonReentrant {
        if (ownerOf(tokenId) != msg.sender) revert NotTokenOwner();

        Position memory pos = positions[tokenId];
        if (block.timestamp < pos.lockEndTime) revert LockNotExpired();

        uint256 amount = pos.amount;

        // Check: remaining staked amount must cover allocations
        if (_userTotalStaked[msg.sender] - amount < IStakingVault(stakingVault).userTotalAllocated(msg.sender)) {
            revert InsufficientUnallocated();
        }

        // CEI: burn NFT + delete storage BEFORE external transfer
        _burn(tokenId);
        delete positions[tokenId];
        emit Withdrawn(msg.sender, tokenId, amount);

        // Transfer AWP back to user (after all state changes)
        awpToken.safeTransfer(msg.sender, amount);
    }

    // ══════════════════════════════════════════════
    //  Transfer hook — maintain accumulators and check allocation coverage
    // ══════════════════════════════════════════════

    /// @dev Override _update to maintain _userTotalStaked accumulators on mint, burn, and transfer.
    function _update(address to, uint256 tokenId, address auth)
        internal
        override
        returns (address)
    {
        address from = super._update(to, tokenId, auth);
        uint128 amt = positions[tokenId].amount;

        if (from != address(0)) {
            // Sender loses stake (transfer or burn)
            _userTotalStaked[from] -= amt;
        }
        if (to != address(0)) {
            // Receiver gains stake (transfer or mint)
            _userTotalStaked[to] += amt;
        }
        // Transfer check: sender must still cover allocations
        if (from != address(0) && to != address(0)) {
            if (_userTotalStaked[from] < IStakingVault(stakingVault).userTotalAllocated(from)) {
                revert InsufficientUnallocated();
            }
        }

        return from;
    }

    // ══════════════════════════════════════════════
    //  View functions
    // ══════════════════════════════════════════════

    /// @inheritdoc IStakeNFT
    function getUserTotalStaked(address user) external view returns (uint256) {
        return _userTotalStaked[user];
    }

    /// @inheritdoc IStakeNFT
    function getVotingPower(uint256 tokenId) external view returns (uint256) {
        return _getVotingPower(tokenId);
    }

    /// @inheritdoc IStakeNFT
    function getUserVotingPower(address user, uint256[] calldata tokenIds) external view returns (uint256 total) {
        for (uint256 i = 0; i < tokenIds.length;) {
            if (ownerOf(tokenIds[i]) != user) revert NotTokenOwner();
            total += _getVotingPower(tokenIds[i]);
            unchecked { ++i; }
        }
    }

    /// @inheritdoc IStakeNFT
    /// @dev Intentional upper-bound estimate: assumes every position is locked at MAX_WEIGHT_DURATION.
    ///      This makes quorum harder to reach than the "true" value would, but is acceptable because
    ///      computing exact totalVotingPower on-chain would require iterating all positions (O(n)).
    ///      Governance may need to adjust quorumPercent over time to compensate.
    function totalVotingPower() external view returns (uint256) {
        // Upper bound: assumes all positions locked at maximum weight
        return awpToken.balanceOf(address(this)) * SQRT_MAX_WEIGHT_FACTOR;
    }

    /// @inheritdoc IStakeNFT
    function remainingTime(uint256 tokenId) external view returns (uint64) {
        return _remainingTime(tokenId);
    }

    /// @inheritdoc IStakeNFT
    function getPositionForVoting(uint256 tokenId) external view returns (
        address owner,
        uint128 amount,
        uint64 lockEndTime,
        uint64 createdAt,
        uint64 remainingSeconds,
        uint256 votingPower
    ) {
        owner = ownerOf(tokenId);
        Position memory pos = positions[tokenId];
        amount = pos.amount;
        lockEndTime = pos.lockEndTime;
        createdAt = pos.createdAt;
        uint256 remaining = pos.lockEndTime > block.timestamp ? pos.lockEndTime - block.timestamp : 0;
        uint256 effective = remaining < MAX_WEIGHT_DURATION ? remaining : MAX_WEIGHT_DURATION;
        remainingSeconds = uint64(remaining);
        votingPower = uint256(pos.amount) * Math.sqrt(effective / VOTE_WEIGHT_DIVISOR);
    }

    // ══════════════════════════════════════════════
    //  Internal helpers
    // ══════════════════════════════════════════════

    /// @dev Calculate voting power for a position: amount * sqrt(min(remainingTime, MAX_WEIGHT_DURATION) / VOTE_WEIGHT_DIVISOR)
    function _getVotingPower(uint256 tokenId) internal view returns (uint256) {
        Position memory pos = positions[tokenId];
        uint256 remaining = pos.lockEndTime > block.timestamp ? pos.lockEndTime - block.timestamp : 0;
        uint256 effective = remaining < MAX_WEIGHT_DURATION ? remaining : MAX_WEIGHT_DURATION;
        return uint256(pos.amount) * Math.sqrt(effective / VOTE_WEIGHT_DIVISOR);
    }

    /// @dev Get remaining lock time for a position in seconds
    function _remainingTime(uint256 tokenId) internal view returns (uint64) {
        uint64 lockEnd = positions[tokenId].lockEndTime;
        return lockEnd > uint64(block.timestamp) ? lockEnd - uint64(block.timestamp) : 0;
    }
}
