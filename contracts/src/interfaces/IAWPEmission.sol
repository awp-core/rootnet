// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IAWPEmission {
    // ── Allocation management (Guardian only) ──
    function submitAllocations(uint256[] calldata packed, uint256 totalWeight, uint256 effectiveEpoch) external;
    function appendAllocations(uint256[] calldata packed, uint256 effectiveEpoch) external;
    function modifyAllocations(uint256[] calldata patches, uint256 newTotalWeight, uint256 effectiveEpoch) external;

    // ── Settlement (callable by anyone) ──
    function settleEpoch(uint256 limit) external;

    // ── Guardian configuration ──
    function setDecayFactor(uint256 newDecayFactor) external;
    function setMaxRecipients(uint256 newMax) external;
    function setEpochDuration(uint256 newDuration) external;
    function setTreasury(address t) external;
    function setGuardian(address g) external;
    function pauseEpochUntil(uint64 resumeTime) external;

    // ── View: state variables ──
    function guardian() external view returns (address);
    function treasury() external view returns (address);
    function decayFactor() external view returns (uint256);
    function currentDailyEmission() external view returns (uint256);
    function epochDuration() external view returns (uint256);
    function baseTime() external view returns (uint256);
    function baseEpoch() external view returns (uint256);
    function settledEpoch() external view returns (uint256);
    function activeEpoch() external view returns (uint256);
    function settleProgress() external view returns (uint256);
    function epochEmissionLocked() external view returns (uint256);
    function maxRecipients() external view returns (uint256);
    function pausedUntil() external view returns (uint64);
    function frozenEpoch() external view returns (uint64);

    // ── View: computed ──
    function currentEpoch() external view returns (uint256);
    function getRecipientCount() external view returns (uint256);
    function getRecipient(uint256 index) external view returns (address);
    function getWeight(address addr) external view returns (uint64);
    function getTotalWeight() external view returns (uint256);
    function getEpochRecipientCount(uint256 epoch) external view returns (uint256);
    function getEpochWeight(uint256 epoch, address addr) external view returns (uint64);
    function getEpochTotalWeight(uint256 epoch) external view returns (uint256);

    // ── Events ──
    event AllocationsSubmitted(uint256 indexed effectiveEpoch, uint256[] packed, uint256 totalWeight);
    event AllocationsAppended(uint256 indexed effectiveEpoch, uint256[] packed);
    event AllocationsModified(uint256 indexed effectiveEpoch, uint256[] patches, uint256 newTotalWeight);
    event RecipientAWPDistributed(uint256 indexed epoch, address indexed recipient, uint256 awpAmount);
    event EpochSettled(uint256 indexed epoch, uint256 totalEmission, uint256 recipientCount);
    event EpochDurationUpdated(uint256 oldDuration, uint256 newDuration);
    event EpochPausedUntil(uint64 resumeTime, uint64 frozenEpoch);
    event GuardianUpdated(address indexed newGuardian);
    event TreasuryUpdated(address indexed newTreasury);
    event MaxRecipientsUpdated(uint256 newMax);
    event DecayFactorUpdated(uint256 newDecayFactor);
}
