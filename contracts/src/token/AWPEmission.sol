// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {EIP712Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import {IAWPToken} from "../interfaces/IAWPToken.sol";
import {IAWPEmission} from "../interfaces/IAWPEmission.sol";

/// @title AWPEmission V3 — UUPS upgradeable emission contract (epoch-versioned weight distribution engine)
/// @notice Guardian-only weight submission, no Oracle/Timelock dependency.
/// @dev Epoch-versioned design: submitAllocations writes to a future epoch slot without clearing old data.
///      settleEpoch promotes the latest submitted weights as activeEpoch when available.
///      Guardian (cross-chain multisig) submits weights directly — no Oracle signatures or Timelock.
///      100% of epoch emission goes to recipients; Guardian includes treasury in recipients for DAO share.
///      Anyone can call settleEpoch to trigger settlement.
///      AWPEmission now owns its own epoch timing (genesisTime + epochDuration).
contract AWPEmission is Initializable, UUPSUpgradeable, ReentrancyGuardUpgradeable, EIP712Upgradeable, IAWPEmission {

    // ══════════════════════════════════════════════
    //  Storage layout — V3 (fresh proxy deployment, epoch-versioned weights)
    // ══════════════════════════════════════════════

    /// @dev Reserved slot 0: was awpRegistry, kept for UUPS proxy upgrade safety
    uint256 private __reserved_slot0;               // slot 0

    /// @notice AWP token contract reference
    IAWPToken public awpToken;                      // slot 1

    /// @dev Freed slot 2: was treasury address (DAO share now handled by including treasury in recipients)
    address private __freed_treasury;               // slot 2

    /// @notice Epoch duration in seconds (default 1 day = 86400)
    uint256 public epochDuration;                   // slot 3 (reused from freed slot)

    /// @notice Number of epochs that have been settled
    uint256 public settledEpoch;                    // slot 4

    /// @notice Genesis timestamp for epoch calculation
    uint256 public genesisTime;                     // slot 5 (reused from freed slot)

    /// @notice Current epoch daily emission amount (AWP wei)
    uint256 public currentDailyEmission;            // slot 6

    /// @notice Epoch of the currently active (most recently promoted) weights
    uint256 public activeEpoch;                     // slot 7

    /// @notice Epoch => packed allocations array: high 160 bits = address, low 96 bits = weight
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

    /// @dev Freed slot 18: was oracleThreshold
    uint256 private __freed_slot18;                  // slot 18

    /// @dev Freed slot 19: was allocationNonce
    uint256 private __freed_slot19;                  // slot 19

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

    /// @notice Guardian address — manages oracle config across chains (cross-chain multisig)
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
    /// @dev Mint amount is zero
    error InvalidAmount();
    /// @dev All epochs up to the current time-based epoch have already been settled
    error EpochNotReady();
    /// @dev All AWP has been minted
    error MiningComplete();
    /// @dev Settlement is in progress; weights/recipients cannot be modified
    error SettlementInProgress();
    /// @dev Array lengths do not match
    error ArrayLengthMismatch();
    /// @dev Parameter value is invalid
    error InvalidParameter();
    /// @dev Duplicate address in recipient list
    error DuplicateRecipient();
    /// @dev effectiveEpoch must be a future epoch
    error MustBeFutureEpoch();
    /// @dev Caller is not the Guardian
    error NotGuardian();

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
    function initialize(
        address awpToken_,
        address guardian_,
        uint256 initialDailyEmission_,
        uint256 genesisTime_,
        uint256 epochDuration_
    ) external initializer {
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
        __EIP712_init("AWPEmission", "2");

        awpToken = IAWPToken(awpToken_);
        guardian = guardian_;
        currentDailyEmission = initialDailyEmission_;
        genesisTime = genesisTime_;
        epochDuration = epochDuration_;
        maxRecipients = 10000;
        decayFactor = 996844;
        _cachedMaxSupply = IAWPToken(awpToken_).MAX_SUPPLY();
    }

    // ══════════════════════════════════════════════
    //  Epoch calculation (self-contained)
    // ══════════════════════════════════════════════

    /// @notice Current epoch number, derived from genesis time and epoch duration
    function currentEpoch() public view returns (uint256) {
        return (block.timestamp - genesisTime) / epochDuration;
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
    /// @dev Guardian is a cross-chain multisig that reads all chains' stake data off-chain,
    ///      computes weights, and submits identical allocations to every chain's AWPEmission.
    ///      Recipients MUST be sorted in ascending address order (duplicates rejected).
    /// @param recipients_ Array of recipient addresses (must be sorted ascending)
    /// @param weights_ Corresponding weight array (uint96)
    /// @param effectiveEpoch The future epoch these weights take effect in
    function submitAllocations(
        address[] calldata recipients_,
        uint96[] calldata weights_,
        uint256 effectiveEpoch
    ) external onlyGuardian {
        if (settleProgress != 0) revert SettlementInProgress();
        if (effectiveEpoch <= settledEpoch) revert MustBeFutureEpoch();
        if (recipients_.length != weights_.length) revert ArrayLengthMismatch();
        if (recipients_.length > maxRecipients) revert InvalidParameter();

        delete _epochAllocations[effectiveEpoch];
        _epochTotalWeight[effectiveEpoch] = 0;

        uint256 len = recipients_.length;
        uint256[] memory packed = new uint256[](len);
        uint256 tw = 0;
        for (uint256 i = 0; i < len;) {
            address addr = recipients_[i];
            uint96 w = weights_[i];
            if (addr == address(0)) revert InvalidRecipient();
            if (w == 0) revert InvalidAmount();
            if (i > 0 && uint160(addr) <= uint160(recipients_[i - 1])) revert DuplicateRecipient();
            packed[i] = (uint256(uint160(addr)) << 96) | uint256(w);
            tw += w;
            unchecked { ++i; }
        }
        _epochAllocations[effectiveEpoch] = packed;
        _epochTotalWeight[effectiveEpoch] = tw;

        emit AllocationsSubmitted(0, recipients_, weights_);
    }

    // ══════════════════════════════════════════════
    //  Guardian parameter configuration
    // ══════════════════════════════════════════════

    /// @notice Update the per-epoch decay factor (Guardian only)
    /// @param newDecayFactor Must be <= DECAY_PRECISION (no growth allowed)
    function setDecayFactor(uint256 newDecayFactor) external onlyGuardian {
        if (newDecayFactor == 0 || newDecayFactor >= DECAY_PRECISION) revert InvalidParameter();
        decayFactor = newDecayFactor;
    }

    // ══════════════════════════════════════════════
    //  Guardian self-management
    // ══════════════════════════════════════════════

    /// @notice Update the guardian address (only Guardian may call — self-sovereign)
    /// @dev Guardian manages itself. If Guardian keys are lost, Timelock can recover via UUPS upgrade.
    function setGuardian(address g) external onlyGuardian {
        guardian = g;
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

    /// @notice Execute epoch settlement: batch-mint recipient AWP emission + DAO share
    /// @param limit Maximum number of recipients to process in this call
    function settleEpoch(uint256 limit) external nonReentrant {
        if (limit == 0) revert InvalidParameter();

        // ── Phase 1: Initialization (O(1), no loops) ──
        if (settleProgress == 0) {
            if (settledEpoch >= currentEpoch()) revert EpochNotReady();

            // Exponential decay: multiply by decay factor every epoch starting from epoch 2
            if (settledEpoch > 0) {
                currentDailyEmission = currentDailyEmission * decayFactor / DECAY_PRECISION;
            }

            // Calculate actual epoch emission (capped at remaining AWP mintable supply)
            uint256 awpRemaining = _cachedMaxSupply - awpToken.totalSupply();
            epochEmissionLocked = currentDailyEmission > awpRemaining ? awpRemaining : currentDailyEmission;
            if (epochEmissionLocked == 0) revert MiningComplete();

            // Promote activeEpoch if new weights were submitted for the epoch being settled
            // Guardian writes to effectiveEpoch > settledEpoch, so check settledEpoch + 1
            if (_epochTotalWeight[settledEpoch + 1] > 0) {
                activeEpoch = settledEpoch + 1;
            }

            // 100% emission → recipients（Guardian 可将 treasury 加入 recipients 实现 DAO 分成）
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
        uint256 epoch = settledEpoch;
        uint256 ae = activeEpoch;

        // Cache remaining AWP supply once — reused in Phase 2 and Phase 3
        uint256 awpRemaining = _cachedMaxSupply - awpToken.totalSupply();

        if (tw > 0) {
            for (uint256 i = start; i < end;) {
                uint256 packed = _epochAllocations[ae][i];
                address recipient = address(uint160(packed >> 96));
                uint256 weight = uint96(packed);
                if (weight > 0) {
                    uint256 share = pool * weight / tw;
                    if (share > 0 && awpRemaining > 0) {
                        uint256 toMint = share > awpRemaining ? awpRemaining : share;
                        // mintAndCall triggers ERC1363 callback; fallback to plain mint if callback reverts
                        // Double try/catch: if both fail (e.g. minter revoked), skip recipient — share goes to DAO
                        bool mintOk;
                        try awpToken.mintAndCall(recipient, toMint, "") {
                            mintOk = true;
                        } catch {
                            try awpToken.mint(recipient, toMint) {
                                mintOk = true;
                            } catch {}
                        }
                        if (mintOk) {
                            minted += toMint;
                            awpRemaining -= toMint;
                            emit RecipientAWPDistributed(epoch, recipient, toMint);
                        }
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
        return address(uint160(_epochAllocations[activeEpoch][index] >> 96));
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
            if (address(uint160(packed >> 96)) == addr) {
                return uint96(packed);
            }
            unchecked { ++i; }
        }
        return 0;
    }
}
