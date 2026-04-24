// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {Base64} from "@openzeppelin/contracts/utils/Base64.sol";
import {IveAWP} from "../interfaces/IveAWP.sol";
import {IAWPAllocator} from "../interfaces/IAWPAllocator.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/// @title veAWP — ERC721 NFT-based staking positions with on-chain metadata
/// @notice Each stake position is an ERC721 token with locked AWP amount and expiry.
///         Voting power = amount * sqrt(min(remainingTime, MAX_WEIGHT_DURATION) / VOTE_WEIGHT_DIVISOR).
/// @dev Key design: _userTotalStaked accumulator provides O(1) balance checks for allocation validation.
///      Transfer hook (_update override) maintains accumulator consistency and checks allocation coverage.
///      Not upgradeable — holds user funds, code-is-law guarantee.
contract veAWP is ERC721, ReentrancyGuard, IveAWP {
    using SafeERC20 for IERC20;
    using Strings for uint256;

    // ── Immutables ──

    /// @notice AWP token contract
    IERC20 public immutable awpToken;

    /// @notice AWPAllocator contract (for allocation balance checks)
    address public immutable awpAllocator;

    // ── Constants ──

    /// @notice Maximum duration for voting power weight calculation (54 weeks)
    uint64 public constant MAX_WEIGHT_DURATION = 54 * 7 days;

    /// @notice Minimum lock duration in seconds (1 day)
    uint64 public constant MIN_LOCK_DURATION = 1 days;

    /// @notice Divisor for voting power formula: amount * sqrt(remainingTime / VOTE_WEIGHT_DIVISOR)
    uint256 public constant VOTE_WEIGHT_DIVISOR = 7 days;

    // ── Storage ──

    /// @notice Guardian address — can rescue tokens and transfer guardian role
    address public guardian;

    /// @notice Position data per tokenId
    mapping(uint256 => Position) public positions;

    /// @dev Auto-incrementing token ID counter, starts at 1
    uint256 private _nextTokenId = 1;

    /// @notice O(1) balance tracking for allocation validation
    mapping(address => uint256) private _userTotalStaked;

    /// @notice Global total staked AWP (not affected by direct transfers to this contract)
    uint256 public totalStaked;

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
    error NotGuardian();
    error CannotRescueStakedToken();
    error PartialWithdrawExceedsBalance();
    error ZeroAddress();

    // ── Events ──

    event GuardianUpdated(address indexed newGuardian);
    event PositionDecreased(uint256 indexed tokenId, uint128 withdrawnAmount, uint128 remainingAmount);
    /// @dev ERC-4906: emitted when token metadata changes
    event MetadataUpdate(uint256 _tokenId);

    // ── Modifiers ──

    modifier onlyGuardian() {
        if (msg.sender != guardian) revert NotGuardian();
        _;
    }

    // ── Constructor ──

    /// @param awpToken_ AWP token address
    /// @param awpAllocator_ AWPAllocator contract address
    /// @param guardian_ Guardian address (can rescue tokens and transfer guardian role)
    constructor(
        address awpToken_,
        address awpAllocator_,
        address guardian_
    ) ERC721("veAWP Position", "veAWP") {
        awpToken = IERC20(awpToken_);
        awpAllocator = awpAllocator_;
        if (guardian_ == address(0)) revert ZeroAddress();
        guardian = guardian_;
    }

    // ══════════════════════════════════════════════
    //  Write functions
    // ══════════════════════════════════════════════

    /// @inheritdoc IveAWP
    function deposit(uint256 amount, uint64 lockDuration) external nonReentrant returns (uint256 tokenId) {
        return _deposit(msg.sender, amount, lockDuration);
    }

    /// @notice Gasless deposit: user signs ERC-2612 permit off-chain, no prior approve tx needed
    function depositWithPermit(
        uint256 amount, uint64 lockDuration,
        uint256 deadline, uint8 v, bytes32 r, bytes32 s
    ) external nonReentrant returns (uint256 tokenId) {
        try IERC20Permit(address(awpToken)).permit(msg.sender, address(this), amount, deadline, v, r, s) {} catch {}
        return _deposit(msg.sender, amount, lockDuration);
    }

    /// @inheritdoc IveAWP
    function addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime) external nonReentrant {
        if (ownerOf(tokenId) != msg.sender) revert NotTokenOwner();

        Position storage pos = positions[tokenId];
        bool updated = false;

        // Extend lock first (enables adding funds to expired positions in the same tx)
        if (newLockEndTime > 0) {
            if (newLockEndTime < pos.lockEndTime) revert LockCannotShorten();
            if (newLockEndTime <= uint64(block.timestamp)) revert LockMustExceedCurrentTime();
            pos.lockEndTime = newLockEndTime;
            updated = true;
        }

        if (amount > 0) {
            // Lock must be active (either originally or after extension above)
            if (pos.lockEndTime <= uint64(block.timestamp)) revert PositionExpired();
            if (uint256(pos.amount) + amount > type(uint128).max) revert InvalidAmount();
            awpToken.safeTransferFrom(msg.sender, address(this), amount);
            pos.amount += uint128(amount);
            _userTotalStaked[msg.sender] += amount;
            totalStaked += amount;
            pos.createdAt = uint64(block.timestamp); // Reset to prevent voting power manipulation
            updated = true;
        }

        if (!updated) revert NothingToUpdate();

        emit PositionIncreased(tokenId, amount, pos.lockEndTime);
        emit MetadataUpdate(tokenId);
    }

    /// @inheritdoc IveAWP
    function withdraw(uint256 tokenId) external nonReentrant {
        if (ownerOf(tokenId) != msg.sender) revert NotTokenOwner();

        Position memory pos = positions[tokenId];
        if (block.timestamp < pos.lockEndTime) revert LockNotExpired();

        uint256 amount = pos.amount;

        // Check: remaining staked amount must cover allocations
        if (_userTotalStaked[msg.sender] - amount < IAWPAllocator(awpAllocator).userTotalAllocated(msg.sender)) {
            revert InsufficientUnallocated();
        }

        // CEI: burn NFT + delete storage BEFORE external transfer
        _burn(tokenId);
        delete positions[tokenId];
        emit Withdrawn(msg.sender, tokenId, amount);

        awpToken.safeTransfer(msg.sender, amount);
    }

    /// @notice Withdraw part of a position's amount (lock must be expired). Keeps the NFT alive.
    /// @param tokenId Position token ID (must be owned by msg.sender)
    /// @param amount Amount to withdraw (must be < position amount; use withdraw() to take all)
    function partialWithdraw(uint256 tokenId, uint128 amount) external nonReentrant {
        if (ownerOf(tokenId) != msg.sender) revert NotTokenOwner();
        if (amount == 0) revert InvalidAmount();

        Position storage pos = positions[tokenId];
        if (block.timestamp < pos.lockEndTime) revert LockNotExpired();
        if (amount >= pos.amount) revert PartialWithdrawExceedsBalance();

        // Check: remaining staked amount must cover allocations
        if (_userTotalStaked[msg.sender] - amount < IAWPAllocator(awpAllocator).userTotalAllocated(msg.sender)) {
            revert InsufficientUnallocated();
        }

        pos.amount -= amount;
        _userTotalStaked[msg.sender] -= amount;
        totalStaked -= amount;

        emit PositionDecreased(tokenId, amount, pos.amount);
        emit MetadataUpdate(tokenId);

        awpToken.safeTransfer(msg.sender, amount);
    }

    /// @notice Withdraw multiple expired positions in a single transaction
    /// @param tokenIds Array of position token IDs to withdraw (all must be owned by msg.sender, all must be expired)
    function batchWithdraw(uint256[] calldata tokenIds) external nonReentrant {
        uint256 totalAmount;

        for (uint256 i = 0; i < tokenIds.length;) {
            uint256 tokenId = tokenIds[i];
            if (ownerOf(tokenId) != msg.sender) revert NotTokenOwner();

            Position memory pos = positions[tokenId];
            if (block.timestamp < pos.lockEndTime) revert LockNotExpired();

            totalAmount += pos.amount;

            _burn(tokenId);
            delete positions[tokenId];
            emit Withdrawn(msg.sender, tokenId, pos.amount);

            unchecked { ++i; }
        }

        // Check: remaining staked amount must cover allocations (single check after all burns)
        if (_userTotalStaked[msg.sender] < IAWPAllocator(awpAllocator).userTotalAllocated(msg.sender)) {
            revert InsufficientUnallocated();
        }

        // Single transfer for all withdrawn AWP
        if (totalAmount > 0) {
            awpToken.safeTransfer(msg.sender, totalAmount);
        }
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
            _userTotalStaked[from] -= amt;
            if (to == address(0)) totalStaked -= amt; // burn
        }
        if (to != address(0)) {
            _userTotalStaked[to] += amt;
            if (from == address(0)) totalStaked += amt; // mint
        }
        // Transfer check: sender must still cover allocations
        if (from != address(0) && to != address(0)) {
            if (_userTotalStaked[from] < IAWPAllocator(awpAllocator).userTotalAllocated(from)) {
                revert InsufficientUnallocated();
            }
        }

        return from;
    }

    // ══════════════════════════════════════════════
    //  View functions
    // ══════════════════════════════════════════════

    /// @inheritdoc IveAWP
    function getUserTotalStaked(address user) external view returns (uint256) {
        return _userTotalStaked[user];
    }

    /// @inheritdoc IveAWP
    function getVotingPower(uint256 tokenId) external view returns (uint256) {
        (uint256 vp,) = _votingPowerAndRemaining(positions[tokenId]);
        return vp;
    }

    /// @inheritdoc IveAWP
    function getUserVotingPower(address user, uint256[] calldata tokenIds) external view returns (uint256 total) {
        for (uint256 i = 0; i < tokenIds.length;) {
            if (ownerOf(tokenIds[i]) != user) revert NotTokenOwner();
            (uint256 vp,) = _votingPowerAndRemaining(positions[tokenIds[i]]);
            total += vp;
            unchecked { ++i; }
        }
    }

    /// @inheritdoc IveAWP
    /// @dev Returns totalStaked (not balanceOf) — immune to direct-transfer inflation.
    ///      Quorum is based on staked capital, not time-weighted voting power.
    function totalVotingPower() external view returns (uint256) {
        return totalStaked;
    }

    /// @inheritdoc IveAWP
    function remainingTime(uint256 tokenId) external view returns (uint64) {
        uint64 lockEnd = positions[tokenId].lockEndTime;
        return lockEnd > uint64(block.timestamp) ? lockEnd - uint64(block.timestamp) : 0;
    }

    /// @inheritdoc IveAWP
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
        (votingPower, remainingSeconds) = _votingPowerAndRemaining(pos);
    }

    /// @dev ERC-165: declare support for ERC-4906 (Metadata Update Extension)
    function supportsInterface(bytes4 interfaceId) public view override returns (bool) {
        return interfaceId == bytes4(0x49064906) || super.supportsInterface(interfaceId);
    }

    /// @notice Token URI — on-chain JSON metadata showing position details
    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        _requireOwned(tokenId);

        Position memory pos = positions[tokenId];
        (uint256 vp, uint64 remaining) = _votingPowerAndRemaining(pos);

        string memory json = string.concat(
            '{"name":"veAWP #', tokenId.toString(),
            '","description":"AWP staking position - ', _formatAWP(pos.amount), ' AWP locked',
            remaining > 0 ? string.concat(' for ', _formatDays(remaining)) : ', unlocked',
            '","attributes":[',
                '{"trait_type":"Amount","display_type":"number","value":', uint256(pos.amount).toString(),
                '},{"trait_type":"Lock End","display_type":"date","value":', uint256(pos.lockEndTime).toString(),
                '},{"trait_type":"Created At","display_type":"date","value":', uint256(pos.createdAt).toString(),
                '},{"trait_type":"Voting Power","display_type":"number","value":', vp.toString(),
                '},{"trait_type":"Remaining Days","display_type":"number","value":', (uint256(remaining) / 1 days).toString(),
            '}]}'
        );

        return string.concat("data:application/json;base64,", Base64.encode(bytes(json)));
    }

    // ══════════════════════════════════════════════
    //  Token rescue
    // ══════════════════════════════════════════════

    /// @notice Update guardian address (self-sovereign)
    function setGuardian(address g) external onlyGuardian {
        if (g == address(0)) revert ZeroAddress();
        guardian = g;
        emit GuardianUpdated(g);
    }

    /// @notice Rescue accidentally sent ERC20 tokens (Guardian only). Cannot rescue staked AWP.
    function rescueToken(address token, address to, uint256 amount) external onlyGuardian {
        if (token == address(awpToken)) revert CannotRescueStakedToken();
        IERC20(token).safeTransfer(to, amount);
    }

    // ══════════════════════════════════════════════
    //  Internal
    // ══════════════════════════════════════════════

    /// @dev Internal deposit logic: transfer AWP from user, mint NFT to user
    function _deposit(address user, uint256 amount, uint64 lockDuration) internal returns (uint256 tokenId) {
        if (amount == 0) revert InvalidAmount();
        if (amount > type(uint128).max) revert InvalidAmount();
        if (lockDuration < MIN_LOCK_DURATION) revert LockTooShort();

        uint64 lockEndTime = uint64(block.timestamp) + lockDuration;
        if (lockEndTime <= uint64(block.timestamp)) revert InvalidAmount(); // overflow guard

        awpToken.safeTransferFrom(user, address(this), amount);

        // Store position BEFORE mint so _update callback sees correct amount
        tokenId = _nextTokenId++;
        positions[tokenId] = Position({
            amount: uint128(amount),
            lockEndTime: lockEndTime,
            createdAt: uint64(block.timestamp)
        });

        _mint(user, tokenId);

        emit Deposited(user, tokenId, amount, lockEndTime);
    }

    /// @dev Calculate voting power and remaining seconds for a position
    function _votingPowerAndRemaining(Position memory pos) internal view returns (uint256 vp, uint64 remaining) {
        uint256 rem = pos.lockEndTime > block.timestamp ? pos.lockEndTime - block.timestamp : 0;
        uint256 effective = rem < MAX_WEIGHT_DURATION ? rem : MAX_WEIGHT_DURATION;
        vp = uint256(pos.amount) * Math.sqrt(effective / VOTE_WEIGHT_DIVISOR);
        remaining = uint64(rem);
    }

    /// @dev Format AWP amount with 18 decimals to human-readable string (e.g., "1000.00")
    function _formatAWP(uint128 amount) internal pure returns (string memory) {
        uint256 whole = uint256(amount) / 1e18;
        uint256 frac = (uint256(amount) % 1e18) / 1e16; // 2 decimal places
        if (frac == 0) return whole.toString();
        return string.concat(whole.toString(), ".", frac < 10 ? "0" : "", frac.toString());
    }

    /// @dev Format seconds to human-readable days string
    function _formatDays(uint64 seconds_) internal pure returns (string memory) {
        uint256 d = uint256(seconds_) / 1 days;
        if (d == 0) return "< 1 day";
        if (d == 1) return "1 day";
        return string.concat(d.toString(), " days");
    }
}
