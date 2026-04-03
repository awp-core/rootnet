// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {IAWPToken} from "../interfaces/IAWPToken.sol";
import {IAWPEmission} from "../interfaces/IAWPEmission.sol";
import {IERC1363Receiver} from "../interfaces/IERC1363Receiver.sol";

/// @title AWPEmission V3 — UUPS upgradeable emission contract (epoch-versioned weight distribution engine)
/// @notice Guardian-only weight submission, no Oracle/Timelock dependency.
/// @dev Epoch-versioned design: submitAllocations writes to a future epoch slot without clearing old data.
///      settleEpoch promotes the latest submitted weights as activeEpoch when available.
///      Guardian (cross-chain multisig) submits weights directly — no Oracle signatures or Timelock.
///      100% of epoch emission goes to recipients; Guardian includes treasury in recipients for DAO share.
///      Anyone can call settleEpoch to trigger settlement.
///      AWPEmission now owns its own epoch timing (baseEpoch + baseTime + epochDuration).
contract AWPEmission is Initializable, UUPSUpgradeable, ReentrancyGuardUpgradeable, IAWPEmission {

    // ══════════════════════════════════════════════
    //  Storage layout — V3 (fresh proxy deployment, epoch-versioned weights)
    // ══════════════════════════════════════════════

    /// @dev Reserved slot 0: was awpRegistry, kept for UUPS proxy upgrade safety
    uint256 private __reserved_slot0;               // slot 0

    /// @notice AWP token contract reference
    IAWPToken public awpToken;                      // slot 1

    /// @notice Treasury address — queryable on-chain reference; Guardian includes it in recipient list for DAO share
    address public treasury;                        // slot 2

    /// @notice Epoch duration in seconds (default 1 day = 86400)
    uint256 public epochDuration;                   // slot 3 (reused from freed slot)

    /// @notice Next epoch to settle (starts at 0; incremented after each settlement)
    uint256 public settledEpoch;                    // slot 4

    /// @notice Base timestamp for epoch calculation (initially = genesisTime, adjusted on epochDuration change)
    uint256 public baseTime;                        // slot 5

    /// @notice Current epoch daily emission amount (AWP wei)
    uint256 public currentDailyEmission;            // slot 6

    /// @notice Epoch of the currently active (most recently promoted) weights
    uint256 public activeEpoch;                     // slot 7

    /// @notice Epoch => packed allocations array: [32-bit reserved | 64-bit weight | 160-bit address]
    mapping(uint256 => uint256[]) internal _epochAllocations;  // slot 8

    /// @dev Freed slot (was _epochWeights mapping of mapping)
    uint256 private __freed_slot9;                             // slot 9

    /// @notice Epoch => sum of all recipient weights
    mapping(uint256 => uint256) internal _epochTotalWeight;    // slot 10

    /// @notice Settlement progress: 0 = idle, >0 = in progress (settleProgress - 1 = last processed index)
    uint256 public settleProgress;                  // slot 11

    /// @notice Total epoch emission locked in Phase 1
    uint256 public epochEmissionLocked;             // slot 12

    /// @dev Number of recipients snapshotted in Phase 1
    uint256 private _snapshotLen;                   // slot 13

    /// @dev Total weight snapshotted in Phase 1
    uint256 private _snapshotWeight;                // slot 14

    /// @dev Recipient pool total snapshotted in Phase 1
    uint256 private _snapshotPool;                  // slot 15

    /// @dev Cumulative AWP minted in Phase 2
    uint256 private _epochMinted;                   // slot 16

    /// @dev Freed slot 17: was oracles[] array (Guardian replaced Oracle multi-sig)
    address[] private __freed_oracles;               // slot 17

    /// @notice Base epoch offset for epoch calculation (adjusted when epochDuration changes)
    uint256 public baseEpoch;                        // slot 18

    /// @notice Epoch pause: 0 = not paused; >0 = resume timestamp. Packed with frozenEpoch in slot 19.
    uint64 public pausedUntil;
    /// @notice Epoch value frozen at pause time (= settledEpoch when paused)
    uint64 public frozenEpoch;                       // slot 19 (packed: pausedUntil + frozenEpoch = 16 bytes)

    /// @notice Maximum number of recipients allowed
    uint256 public maxRecipients;                   // slot 20

    /// @dev Freed slot 21: was awpRegistry reference
    uint256 private __freed_slot21;                 // slot 21

    /// @dev Freed slot 22: was isOracleMap (Guardian replaced Oracle)
    mapping(address => bool) private __freed_isOracleMap; // slot 22

    /// @notice Cached AWP MAX_SUPPLY (set in initialize, avoids repeated cross-contract call)
    uint256 private _cachedMaxSupply;               // slot 23

    /// @notice Decay factor per epoch (default 996844 / 1000000 ≈ 0.3156% decay)
    uint256 public decayFactor;                     // slot 24

    /// @dev Freed slot 25: was emissionSplitBps (split now fully dynamic via Guardian recipients)
    uint256 private __freed_slot25;                 // slot 25

    /// @notice Guardian address — manages emission config across chains (cross-chain multisig)
    address public guardian;                         // slot 26

    /// @dev Reserved storage gap for upgrades
    uint256[33] private __gap;                      // slots 27-59

    // ══════════════════════════════════════════════
    //  Constants
    // ══════════════════════════════════════════════

    /// @notice Exponential decay factor denominator
    uint256 public constant DECAY_PRECISION = 1000000;

    // ══════════════════════════════════════════════
    //  Error definitions
    // ══════════════════════════════════════════════

    /// @dev Recipient address is the zero address
    error InvalidRecipient();
    /// @dev Weight is zero
    error ZeroWeight();
    /// @dev All epochs up to the current time-based epoch have already been settled
    error EpochNotReady();
    /// @dev All AWP has been minted
    error MiningComplete();
    /// @dev Settlement is in progress; weights/recipients cannot be modified
    error SettlementInProgress();
    /// @dev Too many recipients
    error TooManyRecipients(uint256 count, uint256 max);
    /// @dev Decay factor out of valid range
    error InvalidDecayFactor();
    /// @dev Zero address passed
    error ZeroAddress();
    /// @dev Limit parameter is zero
    error ZeroLimit();
    /// @dev Epoch duration is zero
    error ZeroEpochDuration();
    /// @dev effectiveEpoch must be >= settledEpoch (not already settled)
    error MustBeFutureEpoch();
    /// @dev Caller is not the Guardian
    error NotGuardian();
    /// @dev Genesis time has not been reached yet
    error GenesisNotReached();
    /// @dev Patch index exceeds allocation array length
    error IndexOutOfBounds(uint256 index, uint256 length);


    // ══════════════════════════════════════════════
    //  Modifiers
    // ══════════════════════════════════════════════

    /// @dev Only the Guardian (cross-chain multisig) may call
    modifier onlyGuardian() {
        if (msg.sender != guardian) revert NotGuardian();
        _;
    }

    // ══════════════════════════════════════════════
    //  Constructor + initialization
    // ══════════════════════════════════════════════

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Initialize the emission contract (called on proxy deployment)
    /// @param awpToken_ AWP token contract address
    /// @param guardian_ Guardian address (cross-chain Safe multisig) — controls weights, decay, upgrades
    /// @param initialDailyEmission_ Daily emission for the first epoch (wei)
    /// @param genesisTime_ Genesis timestamp for epoch calculation
    /// @param epochDuration_ Epoch duration in seconds (default 86400 = 1 day)
    /// @param treasury_ Treasury address (fallback for failed recipient mints)
    function initialize(
        address awpToken_,
        address guardian_,
        uint256 initialDailyEmission_,
        uint256 genesisTime_,
        uint256 epochDuration_,
        address treasury_
    ) external initializer {
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();

        if (epochDuration_ == 0) revert ZeroEpochDuration();
        if (awpToken_ == address(0) || guardian_ == address(0)) revert ZeroAddress();
        awpToken = IAWPToken(awpToken_);
        guardian = guardian_;
        treasury = treasury_;
        currentDailyEmission = initialDailyEmission_;
        baseTime = genesisTime_;
        epochDuration = epochDuration_;
        maxRecipients = 10000;
        decayFactor = 996844;
        _cachedMaxSupply = IAWPToken(awpToken_).MAX_SUPPLY();
    }

    // ══════════════════════════════════════════════
    //  Epoch calculation (self-contained)
    // ══════════════════════════════════════════════

    /// @notice Current epoch number, derived from base time and epoch duration.
    ///         When paused: returns frozenEpoch (= settledEpoch at pause time).
    ///         When pause expired: resumes from frozenEpoch at pausedUntil timestamp.
    function currentEpoch() public view returns (uint256) {
        uint64 pu = pausedUntil;
        if (pu == 0) {
            // Normal: not paused
            if (block.timestamp < baseTime) revert GenesisNotReached();
            return baseEpoch + (block.timestamp - baseTime) / epochDuration;
        }
        if (block.timestamp < pu) {
            // Paused: epoch frozen
            return uint256(frozenEpoch);
        }
        // Pause expired: resume from frozenEpoch at pausedUntil
        return uint256(frozenEpoch) + (block.timestamp - uint256(pu)) / epochDuration;
    }

    // ══════════════════════════════════════════════
    //  UUPS upgrade authorization
    // ══════════════════════════════════════════════

    /// @dev Only the Guardian may authorize an upgrade
    function _authorizeUpgrade(address) internal view override {
        if (msg.sender != guardian) revert NotGuardian();
    }

    // ══════════════════════════════════════════════
    //  Guardian weight submission (replaces Oracle multi-sig)
    // ══════════════════════════════════════════════

    /// @notice Submit recipient weight allocations for a future epoch (Guardian only)
    /// @dev Resets existing allocations for the epoch, then writes packed entries.
    ///      Each element in packed_: [32-bit zero | 64-bit weight | 160-bit address]
    ///      Empty array is allowed — means this chain has no emission for this epoch.
    ///      totalWeight is the cross-chain global total (sum of all chains' weights).
    ///      Guardian ensures no duplicates and correct encoding off-chain.
    /// @param packed_ Array of packed (weight, address) values — same encoding as storage
    /// @param totalWeight_ Cross-chain global total weight
    /// @param effectiveEpoch The future epoch these weights take effect in
    function submitAllocations(
        uint256[] calldata packed_,
        uint256 totalWeight_,
        uint256 effectiveEpoch
    ) external onlyGuardian {
        if (settleProgress != 0) revert SettlementInProgress();
        if (effectiveEpoch < settledEpoch) revert MustBeFutureEpoch();
        if (packed_.length > maxRecipients) revert TooManyRecipients(packed_.length, maxRecipients);

        delete _epochAllocations[effectiveEpoch];
        _pushPacked(effectiveEpoch, packed_);
        _epochTotalWeight[effectiveEpoch] = totalWeight_;

        emit AllocationsSubmitted(effectiveEpoch, packed_, totalWeight_);
    }

    /// @notice Append additional allocations to an existing epoch (Guardian only)
    /// @dev totalWeight is NOT updated — it was set in submitAllocations.
    ///      Same packed encoding as submitAllocations.
    /// @param packed_ Array of packed (weight, address) values
    /// @param effectiveEpoch The epoch to append to
    function appendAllocations(
        uint256[] calldata packed_,
        uint256 effectiveEpoch
    ) external onlyGuardian {
        if (settleProgress != 0) revert SettlementInProgress();
        if (effectiveEpoch < settledEpoch) revert MustBeFutureEpoch();
        if (packed_.length == 0) revert ZeroLimit();

        uint256 existingLen = _epochAllocations[effectiveEpoch].length;
        if (existingLen + packed_.length > maxRecipients) {
            revert TooManyRecipients(existingLen + packed_.length, maxRecipients);
        }

        _pushPacked(effectiveEpoch, packed_);

        emit AllocationsAppended(effectiveEpoch, packed_);
    }

    /// @notice Modify existing allocations in-place (Guardian only)
    /// @dev Each element in patches_: [32-bit index | 64-bit weight | 160-bit address]
    ///      Same layout as storage but top 32 bits = array index instead of reserved zero.
    ///      weight=0 is allowed (soft-delete, skipped at settle).
    ///      address(0) keeps existing address unchanged (weight-only update).
    ///      totalWeight is replaced unconditionally.
    /// @param patches_ Array of packed (index, weight, address) values
    /// @param newTotalWeight_ New cross-chain global total weight
    /// @param effectiveEpoch The epoch to modify
    function modifyAllocations(
        uint256[] calldata patches_,
        uint256 newTotalWeight_,
        uint256 effectiveEpoch
    ) external onlyGuardian {
        if (settleProgress != 0) revert SettlementInProgress();
        if (effectiveEpoch < settledEpoch) revert MustBeFutureEpoch();

        uint256[] storage arr = _epochAllocations[effectiveEpoch];
        uint256 arrLen = arr.length;
        uint256 len = patches_.length;

        for (uint256 i = 0; i < len;) {
            uint256 p = patches_[i];
            uint256 idx = p >> 224;                              // top 32 bits
            uint64 newWeight = uint64(p >> 160);                 // middle 64 bits
            address newAddr = address(uint160(p));               // bottom 160 bits

            if (idx >= arrLen) revert IndexOutOfBounds(idx, arrLen);

            if (newAddr == address(0)) {
                address existingAddr = address(uint160(arr[idx]));
                arr[idx] = (uint256(newWeight) << 160) | uint256(uint160(existingAddr));
            } else {
                arr[idx] = (uint256(newWeight) << 160) | uint256(uint160(newAddr));
            }
            unchecked { ++i; }
        }

        _epochTotalWeight[effectiveEpoch] = newTotalWeight_;

        emit AllocationsModified(effectiveEpoch, patches_, newTotalWeight_);
    }

    /// @dev Push packed entries to storage. Guardian (trusted) ensures valid encoding off-chain.
    function _pushPacked(uint256 effectiveEpoch, uint256[] calldata packed_) internal {
        uint256[] storage arr = _epochAllocations[effectiveEpoch];
        uint256 len = packed_.length;
        for (uint256 i = 0; i < len;) {
            arr.push(packed_[i]);
            unchecked { ++i; }
        }
    }

    // ══════════════════════════════════════════════
    //  Guardian parameter configuration
    // ══════════════════════════════════════════════

    /// @notice Minimum decay factor (90% of DECAY_PRECISION = max 10% decay per epoch)
    uint256 public constant MIN_DECAY_FACTOR = 900000;


    /// @notice Update the maximum number of recipients per epoch (Guardian only)
    function setMaxRecipients(uint256 newMax) external onlyGuardian {
        if (newMax == 0) revert ZeroLimit();
        maxRecipients = newMax;
        emit MaxRecipientsUpdated(newMax);
    }

    /// @notice Update the per-epoch decay factor (Guardian only)
    /// @param newDecayFactor Must be >= MIN_DECAY_FACTOR and < DECAY_PRECISION
    function setDecayFactor(uint256 newDecayFactor) external onlyGuardian {
        if (newDecayFactor < MIN_DECAY_FACTOR || newDecayFactor >= DECAY_PRECISION) revert InvalidDecayFactor();
        decayFactor = newDecayFactor;
        emit DecayFactorUpdated(newDecayFactor);
    }

    /// @notice Update epoch duration (Guardian only)
    /// @dev Per-epoch emission is NOT adjusted — total cumulative emission stays the same,
    ///      but the time to exhaust it scales proportionally with epoch length.
    /// @param newDuration New epoch duration in seconds (must be > 0)
    function setEpochDuration(uint256 newDuration) external onlyGuardian {
        _checkResume();
        if (newDuration == 0) revert ZeroEpochDuration();
        if (settleProgress != 0) revert SettlementInProgress();

        uint256 oldDuration = epochDuration;
        // Anchor to settledEpoch so unsettled epoch timing is recalculated under the new duration.
        // settledEpoch's start time = baseTime + (settledEpoch - baseEpoch) * oldDuration
        uint256 newBaseTime = baseTime + (settledEpoch - baseEpoch) * oldDuration;
        baseEpoch = settledEpoch;
        baseTime = newBaseTime;
        epochDuration = newDuration;

        emit EpochDurationUpdated(oldDuration, newDuration);
    }

    // ══════════════════════════════════════════════
    //  Epoch pause
    // ══════════════════════════════════════════════

    /// @notice Pause epoch counting until a specified timestamp (Guardian only).
    ///         During pause, currentEpoch() returns settledEpoch (frozen).
    ///         Can be called again to update resumeTime. Pass 0 to resume immediately.
    error InvalidResumeTime();

    function pauseEpochUntil(uint64 resumeTime) external onlyGuardian {
        _checkResume();
        // resumeTime must be in the future (or 0 for immediate resume)
        if (resumeTime != 0 && resumeTime <= uint64(block.timestamp)) revert InvalidResumeTime();
        if (pausedUntil == 0 && resumeTime != 0) {
            frozenEpoch = uint64(settledEpoch);
        } else if (pausedUntil != 0 && resumeTime == 0) {
            // Immediate resume: rebase to now
            baseEpoch = uint256(frozenEpoch);
            baseTime = block.timestamp;
            frozenEpoch = 0;
        }
        pausedUntil = resumeTime;
        emit EpochPausedUntil(resumeTime, frozenEpoch);
    }

    /// @dev Clean up expired pause: absorb pause duration into baseTime/baseEpoch
    function _checkResume() internal {
        uint64 pu = pausedUntil;
        if (pu != 0 && block.timestamp >= uint256(pu)) {
            baseEpoch = uint256(frozenEpoch);
            baseTime = uint256(pu);
            pausedUntil = 0;
            frozenEpoch = 0;
        }
    }

    // ══════════════════════════════════════════════
    //  Guardian self-management
    // ══════════════════════════════════════════════

    /// @notice Update the treasury address for failed mint fallback (Guardian only)
    function setTreasury(address t) external onlyGuardian {
        if (t == address(0)) revert ZeroAddress();
        treasury = t;
        emit TreasuryUpdated(t);
    }

    /// @notice Update the guardian address (only Guardian may call — self-sovereign)
    /// @dev Guardian manages itself. If Guardian keys are lost, there is no on-chain recovery path.
    function setGuardian(address g) external onlyGuardian {
        if (g == address(0)) revert ZeroAddress();
        guardian = g;
        emit GuardianUpdated(g);
    }

    // ══════════════════════════════════════════════
    //  DEPRECATED — Oracle configuration removed (Guardian replaces Oracle)
    // ══════════════════════════════════════════════
    // Oracle infrastructure (oracles[], oracleThreshold, isOracleMap, submitAllocations with signatures)
    // has been replaced by Guardian-only submitAllocations. The storage slots are freed but preserved
    // for UUPS proxy upgrade compatibility.

    // ══════════════════════════════════════════════
    //  Emission settlement (3-phase design, callable by anyone)
    // ══════════════════════════════════════════════

    /// @notice Execute epoch settlement: batch-mint recipient AWP emission
    /// @param limit Maximum number of recipients to process in this call
    function settleEpoch(uint256 limit) external nonReentrant {
        _checkResume();
        if (limit == 0) revert ZeroLimit();

        // ── Phase 1: Initialization (O(1), no loops) ──
        if (settleProgress == 0) {
            // settledEpoch = next epoch to settle; can settle if <= currentEpoch
            if (settledEpoch > currentEpoch()) revert EpochNotReady();

            // Exponential decay: skip for the very first epoch (settledEpoch == 0)
            if (settledEpoch > 0) {
                currentDailyEmission = currentDailyEmission * decayFactor / DECAY_PRECISION;
            }

            // Calculate actual epoch emission (capped at remaining AWP mintable supply)
            uint256 awpRemaining = _cachedMaxSupply - awpToken.totalSupply();
            epochEmissionLocked = currentDailyEmission > awpRemaining ? awpRemaining : currentDailyEmission;
            if (epochEmissionLocked == 0) revert MiningComplete();

            // Promote activeEpoch if new weights exist for the epoch being settled.
            // O(1): activeEpoch persists in storage; if no new weights, the previous
            // allocation carries forward automatically (supports weekly submissions).
            if (_epochTotalWeight[settledEpoch] > 0) {
                activeEpoch = settledEpoch;
            }

            // 100% emission to recipients (Guardian can include treasury in recipients for DAO share)
            _snapshotPool = epochEmissionLocked;
            _epochMinted = 0;

            // Snapshot weight and recipient count from the active epoch
            uint256 ae = activeEpoch;
            _snapshotWeight = _epochTotalWeight[ae];
            _snapshotLen = _epochAllocations[ae].length;

            settleProgress = 1;
        }

        // ── Phase 2: Batch direct minting ──
        uint256 start = settleProgress - 1;
        uint256 snapshotLen = _snapshotLen;
        uint256 end;
        unchecked { end = start + limit; }
        if (end > snapshotLen || end < start) end = snapshotLen; // overflow-safe clamp

        uint256 pool = _snapshotPool;
        uint256 tw = _snapshotWeight;
        uint256 minted = _epochMinted;
        uint256 epoch = settledEpoch; // epoch being settled (before increment in Phase 3)
        uint256 ae = activeEpoch;

        // Use epoch-locked budget minus already-minted as the cap (not live totalSupply)
        uint256 awpRemaining = epochEmissionLocked - minted;

        if (tw > 0) {
            uint256[] storage allocations = _epochAllocations[ae];
            IAWPToken token = awpToken;
            for (uint256 i = start; i < end;) {
                uint256 packed = allocations[i];
                uint256 weight = uint64(packed >> 160);
                if (weight > 0) {
                    uint256 share = pool * weight / tw;
                    if (share > 0 && awpRemaining > 0) {
                        uint256 toMint = share > awpRemaining ? awpRemaining : share;
                        address recipient = address(uint160(packed));
                        // mint always succeeds; best-effort ERC1363 callback for contract recipients.
                        token.mint(recipient, toMint);
                        if (recipient.code.length > 0) {
                            try IERC1363Receiver(recipient).onTransferReceived(address(this), address(this), toMint, "") {} catch {}
                        }
                        minted += toMint;
                        awpRemaining -= toMint;
                        emit RecipientAWPDistributed(epoch, recipient, toMint);
                    }
                }
                unchecked { ++i; }
            }
        }
        _epochMinted = minted;

        // ── Phase 3: Finalize ──
        if (end >= snapshotLen) {
            _snapshotPool = 0;
            _epochMinted = 0;
            settleProgress = 0;
            unchecked { settledEpoch++; }

            emit EpochSettled(epoch, minted, snapshotLen);
        } else {
            settleProgress = end + 1;
        }
    }

    // ══════════════════════════════════════════════
    //  View
    // ══════════════════════════════════════════════

    /// @notice Get the number of recipients in the active epoch
    function getRecipientCount() external view returns (uint256) {
        return _epochAllocations[activeEpoch].length;
    }

    /// @notice Get a recipient address by index from the active epoch (unpacked)
    function getRecipient(uint256 index) external view returns (address) {
        return address(uint160(_epochAllocations[activeEpoch][index]));
    }

    /// @notice Get the weight of a recipient in the active epoch (O(n) scan)
    function getWeight(address addr) external view returns (uint96) {
        return _scanWeight(_epochAllocations[activeEpoch], addr);
    }

    /// @notice Get total weight in the active epoch
    function getTotalWeight() external view returns (uint256) {
        return _epochTotalWeight[activeEpoch];
    }

    /// @notice Get the number of recipients for a specific epoch
    function getEpochRecipientCount(uint256 epoch) external view returns (uint256) {
        return _epochAllocations[epoch].length;
    }

    /// @notice Get the weight of a recipient for a specific epoch (O(n) scan, view only)
    function getEpochWeight(uint256 epoch, address addr) external view returns (uint96) {
        return _scanWeight(_epochAllocations[epoch], addr);
    }

    /// @notice Get total weight for a specific epoch
    function getEpochTotalWeight(uint256 epoch) external view returns (uint256) {
        return _epochTotalWeight[epoch];
    }

    // ══════════════════════════════════════════════
    //  Internal
    // ══════════════════════════════════════════════

    /// @dev Scan a packed allocations array to find the weight for a given address
    function _scanWeight(uint256[] storage allocations, address addr) internal view returns (uint96) {
        uint256 len = allocations.length;
        for (uint256 i = 0; i < len;) {
            uint256 packed = allocations[i];
            if (address(uint160(packed)) == addr) {
                return uint96(uint64(packed >> 160));
            }
            unchecked { ++i; }
        }
        return 0;
    }
}
