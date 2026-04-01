// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IAWPEmission {
    function settleEpoch(uint256 limit) external;
    function submitAllocations(address[] calldata recipients, uint96[] calldata weights, uint256 effectiveEpoch) external;
    function emergencySetWeight(uint256 epoch_, uint256 index, address addr, uint96 weight) external;
    function setGuardian(address g) external;
    function setDecayFactor(uint256 newDecayFactor) external;
    function setEmissionSplitBps(uint256 newSplitBps) external;
    function guardian() external view returns (address);
    function decayFactor() external view returns (uint256);
    function emissionSplitBps() external view returns (uint256);
    function settledEpoch() external view returns (uint256);
    function activeEpoch() external view returns (uint256);
    function currentDailyEmission() external view returns (uint256);
    function settleProgress() external view returns (uint256);
    function getRecipientCount() external view returns (uint256);
    function getRecipient(uint256 index) external view returns (address);
    function getWeight(address addr) external view returns (uint96);
    function getTotalWeight() external view returns (uint256);
    function getEpochRecipientCount(uint256 epoch) external view returns (uint256);
    function getEpochWeight(uint256 epoch, address addr) external view returns (uint96);
    function getEpochTotalWeight(uint256 epoch) external view returns (uint256);

    event AllocationsSubmitted(uint256 indexed nonce, address[] recipients, uint96[] weights);
    event GovernanceWeightUpdated(address indexed addr, uint96 weight);
    event RecipientAWPDistributed(uint256 indexed epoch, address indexed recipient, uint256 awpAmount);
    event DAOMatchDistributed(uint256 indexed epoch, uint256 amount);
    event EpochSettled(uint256 indexed epoch, uint256 totalEmission, uint256 recipientCount);
}
