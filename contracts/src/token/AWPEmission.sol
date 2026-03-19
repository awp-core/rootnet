// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {EIP712Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {IAWPToken} from "../interfaces/IAWPToken.sol";
import {IAWPEmission} from "../interfaces/IAWPEmission.sol";

/// @title AWPEmission V3 — UUPS upgradeable emission contract (epoch-versioned weight distribution engine)
/// @notice Manages epoch-versioned recipient weights, computes exponential decay emission, and batch-mints AWP to recipients and the DAO.
/// @dev Epoch-versioned design: submitAllocations writes to a future epoch slot without clearing old data.
///      settleEpoch promotes the latest submitted weights as activeEpoch when available.
///      Adds multi-oracle batch weight submission (EIP-712 signature verification).
///      DAO configures parameters via Timelock calls to emergencySetWeight.
///      Anyone can call settleEpoch to trigger settlement.
///      AWPEmission now owns its own epoch timing (genesisTime + epochDuration).
contract AWPEmission is Initializable, UUPSUpgradeable, ReentrancyGuardUpgradeable, EIP712Upgradeable, IAWPEmission {

    // ══════════════════════════════════════════════
    //  Storage layout — V3 (fresh proxy deployment, epoch-versioned weights)
    // ══════════════════════════════════════════════

    /// @dev Reserved slot 0: was rootNet, kept for UUPS proxy upgrade safety
    uint256 private __reserved_slot0;               // slot 0

    /// @notice AWP token contract reference
    IAWPToken public awpToken;                      // slot 1

    /// @notice Treasury (Timelock) address — holds governance operation rights + receives DAO share
    address public treasury;                        // slot 2

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

    /// @notice List of oracle addresses
    address[] public oracles;                       // slot 17

    /// @notice Multisig threshold (requires >= threshold valid signatures)
    uint256 public oracleThreshold;                 // slot 18

    /// @notice Allocation submission nonce, prevents replay attacks
    uint256 public allocationNonce;                 // slot 19

    /// @notice Maximum number of recipients allowed
    uint256 public maxRecipients;                   // slot 20

    /// @dev Reserved slot 21: was rootNet reference (no longer needed)
    uint256 private __freed_slot21;                 // slot 21

    /// @dev Reserved storage gap for upgrades
    uint256[38] private __gap;                      // slots 22-59

    // ══════════════════════════════════════════════
    //  Constants
    // ══════════════════════════════════════════════

    /// @notice Exponential decay factor numerator
    uint256 public constant DECAY_FACTOR = 996844;

    /// @notice Exponential decay factor denominator
    uint256 public constant DECAY_PRECISION = 1000000;

    /// @notice Emission split ratio (basis points): 5000 = 50% to recipients
    uint256 public constant EMISSION_SPLIT_BPS = 5000;

    // ══════════════════════════════════════════════
    //  EIP-712 type hashes
    // ══════════════════════════════════════════════

    /// @dev EIP-712 type hash: SubmitAllocations(address[] recipients, uint96[] weights, uint256 nonce, uint256 effectiveEpoch)
    bytes32 private constant ALLOCATION_TYPEHASH =
        keccak256("SubmitAllocations(address[] recipients,uint96[] weights,uint256 nonce,uint256 effectiveEpoch)");

    // ══════════════════════════════════════════════
    //  Error definitions
    // ══════════════════════════════════════════════

    /// @dev Caller is not the Timelock
    error NotTimelock();
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
    /// @dev Oracle is not configured
    error OracleNotConfigured();
    /// @dev Oracle configuration is invalid (invalid threshold or address list)
    error InvalidOracleConfig();
    /// @dev Number of valid signatures is below the threshold
    error InvalidSignatureCount();
    /// @dev Duplicate address in oracle array
    error DuplicateOracle();
    /// @dev Signer is not a registered oracle
    error UnknownOracle();
    /// @dev Same signer appears more than once
    error DuplicateSigner();
    /// @dev Array lengths do not match
    error ArrayLengthMismatch();
    /// @dev Parameter value is invalid
    error InvalidParameter();
    /// @dev Duplicate address in recipient list
    error DuplicateRecipient();
    /// @dev effectiveEpoch must be a future epoch
    error MustBeFutureEpoch();

    // ══════════════════════════════════════════════
    //  Modifiers
    // ══════════════════════════════════════════════

    /// @dev Only the Timelock may call
    modifier onlyTimelock() {
        if (msg.sender != treasury) revert NotTimelock();
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
    /// @param treasury_ Treasury address (Timelock) — receives DAO share + holds governance rights
    /// @param initialDailyEmission_ Daily emission for the first epoch (wei)
    /// @param genesisTime_ Genesis timestamp for epoch calculation
    /// @param epochDuration_ Epoch duration in seconds (default 86400 = 1 day)
    function initialize(
        address awpToken_,
        address treasury_,
        uint256 initialDailyEmission_,
        uint256 genesisTime_,
        uint256 epochDuration_
    ) external initializer {
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
        __EIP712_init("AWPEmission", "2");

        awpToken = IAWPToken(awpToken_);
        treasury = treasury_;
        currentDailyEmission = initialDailyEmission_;
        genesisTime = genesisTime_;
        epochDuration = epochDuration_;
        maxRecipients = 10000;
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

    /// @dev Only the Timelock may authorize an upgrade
    function _authorizeUpgrade(address) internal override onlyTimelock {}

    // ══════════════════════════════════════════════
    //  Multi-oracle batch weight submission (EIP-712 signature verification)
    // ══════════════════════════════════════════════

    /// @notice Submit oracle-signed recipient weight allocations for a future epoch
    /// @dev Requires >= oracleThreshold valid oracle signatures using EIP-712 structured signing.
    ///      Packs address (160 bits) and weight (96 bits) into a single uint256 per recipient.
    ///      Recipients MUST be sorted in ascending address order (duplicates rejected).
    /// @param recipients_ Array of recipient addresses (must be sorted ascending)
    /// @param weights_ Corresponding weight array (uint96)
    /// @param signatures Array of oracle signatures (65 bytes each)
    /// @param effectiveEpoch The future epoch these weights take effect in
    function submitAllocations(
        address[] calldata recipients_,
        uint96[] calldata weights_,
        bytes[] calldata signatures,
        uint256 effectiveEpoch
    ) external {
        // Cannot submit during active settlement (prevents reentrant overwrite via mintAndCall callback)
        if (settleProgress != 0) revert SettlementInProgress();
        // Allow submission for any epoch > settledEpoch (including already-elapsed epochs for catch-up)
        if (effectiveEpoch <= settledEpoch) revert MustBeFutureEpoch();
        // Check that oracle is configured
        if (oracleThreshold == 0) revert OracleNotConfigured();
        // Check that signature count does not exceed total oracle count
        if (signatures.length > oracles.length) revert InvalidSignatureCount();
        // Check array length match
        if (recipients_.length != weights_.length) revert ArrayLengthMismatch();
        // Check recipient count does not exceed the maximum
        if (recipients_.length > maxRecipients) revert InvalidParameter();

        // Construct EIP-712 struct hash (arrays encoded per EIP-712 spec: keccak256 of concatenated padded elements)
        bytes32 structHash = keccak256(
            abi.encode(
                ALLOCATION_TYPEHASH,
                _hashAddressArray(recipients_),
                _hashUint96Array(weights_),
                allocationNonce,
                effectiveEpoch
            )
        );
        bytes32 digest = _hashTypedDataV4(structHash);

        // Verify signatures: check each signature comes from a registered oracle with no duplicates
        uint256 validCount = 0;
        address[] memory seen = new address[](signatures.length);
        for (uint256 i = 0; i < signatures.length;) {
            address signer = ECDSA.recover(digest, signatures[i]);
            if (!_isOracle(signer)) revert UnknownOracle();
            // Check for duplicate signatures
            for (uint256 j = 0; j < validCount;) {
                if (seen[j] == signer) revert DuplicateSigner();
                unchecked { ++j; }
            }
            seen[validCount] = signer;
            unchecked { ++validCount; ++i; }
        }
        if (validCount < oracleThreshold) revert InvalidSignatureCount();

        // Clear previous submission for same effectiveEpoch (O(1) length reset)
        delete _epochAllocations[effectiveEpoch];
        _epochTotalWeight[effectiveEpoch] = 0;

        // Build packed array in memory, then assign to storage in one shot
        uint256 len = recipients_.length;
        uint256[] memory packed = new uint256[](len);
        uint256 tw = 0;
        for (uint256 i = 0; i < len;) {
            address addr = recipients_[i];
            uint96 w = weights_[i];
            if (addr == address(0)) revert InvalidRecipient();
            if (w == 0) revert InvalidAmount();
            // Duplicate check: require sorted ascending order
            if (i > 0 && uint160(addr) <= uint160(recipients_[i - 1])) revert DuplicateRecipient();
            packed[i] = (uint256(uint160(addr)) << 96) | uint256(w);
            tw += w;
            unchecked { ++i; }
        }
        _epochAllocations[effectiveEpoch] = packed;
        _epochTotalWeight[effectiveEpoch] = tw;

        // Increment nonce to prevent replay
        allocationNonce++;

        emit AllocationsSubmitted(allocationNonce - 1, recipients_, weights_);
    }

    // ══════════════════════════════════════════════
    //  Governance parameter configuration (onlyTimelock)
    // ══════════════════════════════════════════════

    /// @notice Emergency overwrite a recipient entry at a given index (Timelock only, used when oracle fails)
    /// @dev Caller provides epoch and index for O(1) access. Replaces both address and weight at that slot.
    /// @param epoch_ Target epoch
    /// @param index Index in the packed allocations array
    /// @param addr New recipient address
    /// @param weight New weight value (uint96)
    function emergencySetWeight(uint256 epoch_, uint256 index, address addr, uint96 weight) external onlyTimelock {
        if (settleProgress != 0) revert SettlementInProgress();
        if (addr == address(0)) revert InvalidRecipient();

        uint256[] storage allocs = _epochAllocations[epoch_];
        if (index >= allocs.length) revert InvalidParameter();

        uint96 oldW = uint96(allocs[index]);
        allocs[index] = (uint256(uint160(addr)) << 96) | uint256(weight);
        _epochTotalWeight[epoch_] = _epochTotalWeight[epoch_] - oldW + weight;

        emit GovernanceWeightUpdated(addr, weight);
    }

    // ══════════════════════════════════════════════
    //  Oracle configuration (onlyTimelock)
    // ══════════════════════════════════════════════

    /// @notice Set the oracle address list and multisig threshold
    /// @param oracles_ New oracle address array
    /// @param threshold_ New multisig threshold
    function setOracleConfig(address[] calldata oracles_, uint256 threshold_) external onlyTimelock {
        if (settleProgress != 0) revert SettlementInProgress();
        if (threshold_ == 0 || threshold_ > oracles_.length) revert InvalidOracleConfig();

        // Check for zero addresses and duplicates
        for (uint256 i = 0; i < oracles_.length;) {
            if (oracles_[i] == address(0)) revert InvalidOracleConfig();
            for (uint256 j = 0; j < i;) {
                if (oracles_[j] == oracles_[i]) revert DuplicateOracle();
                unchecked { ++j; }
            }
            unchecked { ++i; }
        }

        // Clear old array
        delete oracles;

        // Copy new array
        for (uint256 i = 0; i < oracles_.length;) {
            oracles.push(oracles_[i]);
            unchecked { ++i; }
        }
        oracleThreshold = threshold_;

        emit OracleConfigUpdated(oracles_, threshold_);
    }

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
                currentDailyEmission = currentDailyEmission * DECAY_FACTOR / DECAY_PRECISION;
            }

            // Calculate actual epoch emission (capped at remaining AWP mintable supply)
            uint256 awpRemaining = awpToken.MAX_SUPPLY() - awpToken.totalSupply();
            epochEmissionLocked = currentDailyEmission > awpRemaining ? awpRemaining : currentDailyEmission;
            if (epochEmissionLocked == 0) revert MiningComplete();

            // Promote activeEpoch if new weights were submitted for settledEpoch
            if (_epochTotalWeight[settledEpoch] > 0) {
                activeEpoch = settledEpoch;
            }

            // Recipient pool = total emission × 50%
            _snapshotPool = epochEmissionLocked * EMISSION_SPLIT_BPS / 10000;
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
        uint256 maxSupply = awpToken.MAX_SUPPLY();
        uint256 awpRemaining = maxSupply - awpToken.totalSupply();

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

        // ── Phase 3: Finalize — mint DAO share ──
        if (end >= snapshotLen) {
            uint256 daoShare = minted >= epochEmissionLocked
                ? 0
                : epochEmissionLocked - minted;
            uint256 actualDaoMinted = 0;
            if (daoShare > 0 && awpRemaining > 0) {
                actualDaoMinted = daoShare > awpRemaining ? awpRemaining : daoShare;
                awpToken.mint(treasury, actualDaoMinted);
            }
            emit DAOMatchDistributed(epoch, actualDaoMinted);

            _snapshotPool = 0;
            _epochMinted = 0;
            settleProgress = 0;
            unchecked { settledEpoch++; }

            emit EpochSettled(epoch, epochEmissionLocked, snapshotLen);
        } else {
            settleProgress = end + 1;
        }
    }

    // ══════════════════════════════════════════════
    //  View
    // ══════════════════════════════════════════════

    /// @notice Get the number of oracles
    function getOracleCount() external view returns (uint256) {
        return oracles.length;
    }

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

    /// @dev Check whether an address is a registered oracle
    function _isOracle(address addr) internal view returns (bool) {
        for (uint256 i = 0; i < oracles.length;) {
            if (oracles[i] == addr) return true;
            unchecked { ++i; }
        }
        return false;
    }

    /// @dev EIP-712 compliant hash of address[]: keccak256 of concatenated 32-byte padded addresses.
    function _hashAddressArray(address[] calldata arr) internal pure returns (bytes32) {
        uint256 len = arr.length;
        uint256 ptr;
        assembly {
            ptr := mload(0x40)
            mstore(0x40, add(ptr, mul(len, 32))) // advance free memory pointer
        }
        for (uint256 i = 0; i < len;) {
            bytes32 v = bytes32(uint256(uint160(arr[i])));
            assembly { mstore(add(ptr, mul(i, 32)), v) }
            unchecked { ++i; }
        }
        bytes32 result;
        assembly { result := keccak256(ptr, mul(len, 32)) }
        return result;
    }

    /// @dev EIP-712 compliant hash of uint96[]: keccak256 of concatenated 32-byte padded values.
    function _hashUint96Array(uint96[] calldata arr) internal pure returns (bytes32) {
        uint256 len = arr.length;
        uint256 ptr;
        assembly {
            ptr := mload(0x40)
            mstore(0x40, add(ptr, mul(len, 32))) // advance free memory pointer
        }
        for (uint256 i = 0; i < len;) {
            bytes32 v = bytes32(uint256(arr[i]));
            assembly { mstore(add(ptr, mul(i, 32)), v) }
            unchecked { ++i; }
        }
        bytes32 result;
        assembly { result := keccak256(ptr, mul(len, 32)) }
        return result;
    }

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
