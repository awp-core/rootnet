// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {GovernorUpgradeable} from "@openzeppelin/contracts-upgradeable/governance/GovernorUpgradeable.sol";
import {GovernorSettingsUpgradeable} from "@openzeppelin/contracts-upgradeable/governance/extensions/GovernorSettingsUpgradeable.sol";
import {GovernorTimelockControlUpgradeable} from "@openzeppelin/contracts-upgradeable/governance/extensions/GovernorTimelockControlUpgradeable.sol";
import {GovernorPreventLateQuorumUpgradeable} from "@openzeppelin/contracts-upgradeable/governance/extensions/GovernorPreventLateQuorumUpgradeable.sol";
import {TimelockControllerUpgradeable} from "@openzeppelin/contracts-upgradeable/governance/TimelockControllerUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IveAWP} from "../interfaces/IveAWP.sol";

/// @title AWPDAO — UUPS upgradeable governance with veAWP NFT-based voting
/// @notice veAWP holders vote with their position NFTs. Voters submit tokenId arrays.
///         Supports gasless voting via castVoteWithReasonAndParamsBySig (OZ built-in EIP-712).
/// @dev Inheritance: GovernorUpgradeable, GovernorSettingsUpgradeable, GovernorTimelockControlUpgradeable,
///      GovernorPreventLateQuorumUpgradeable, UUPSUpgradeable.
///      veAWP is immutable (proxy address). Guardian (storage, changeable) or Executor can upgrade.
contract AWPDAO is
    GovernorUpgradeable,
    GovernorSettingsUpgradeable,
    GovernorTimelockControlUpgradeable,
    GovernorPreventLateQuorumUpgradeable,
    UUPSUpgradeable
{
    // ── Vote types ──
    enum VoteType { Against, For, Abstain }

    // ── Vote counting storage ──
    struct ProposalVote {
        uint128 againstVotes;
        uint128 forVotes;
        uint128 abstainVotes;
    }

    // ── Immutables ──

    /// @notice veAWP contract for position queries
    IveAWP public immutable veAWP;

    // ── Storage ──

    /// @notice Guardian address — can upgrade DAO and emergency-cancel proposals
    address public guardian;

    /// @notice Quorum percentage (e.g. 4 = 4% of total voting power)
    uint256 public quorumPercent;

    /// @notice Double-vote prevention per tokenId per proposal
    mapping(uint256 proposalId => mapping(uint256 tokenId => bool)) public hasVotedWithToken;

    /// @notice Timestamp at proposal creation (for anti-manipulation check)
    mapping(uint256 proposalId => uint256) public proposalCreatedAt;

    /// @notice Whether a proposal is signal-only (no on-chain execution, skips Timelock)
    mapping(uint256 proposalId => bool) public isSignalProposal;

    /// @notice Snapshot of totalVotingPower at proposal creation (for quorum)
    mapping(uint256 proposalId => uint256) public proposalTotalVotingPower;

    /// @dev Vote tallies per proposal (time-weighted power — determines outcome)
    mapping(uint256 proposalId => ProposalVote) private _proposalVotes;

    /// @dev Raw staked amount voted For+Abstain per proposal (determines quorum)
    mapping(uint256 proposalId => uint256) private _proposalQuorumVotes;

    /// @dev Reserved storage gap
    uint256[41] private __gap;

    // ── Errors ──

    error NoTokens();
    error NotTokenOwner();
    error TokenAlreadyVoted();
    error LockExpired();
    error MintedAfterProposal();
    error UseProposeWithTokens();
    error UseCastVoteWithParams();
    error InvalidQuorumPercent();
    error ZeroTotalVotingPower();
    error NotGuardianOrExecutor();
    error NotGuardian();

    // ── Events ──

    event QuorumPercentUpdated(uint256 newQuorumPercent);
    event GuardianUpdated(address indexed newGuardian);

    // ═══════════════════════════════════════════════
    //  Constructor / Initialize
    // ═══════════════════════════════════════════════

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor(address veAWP_) {
        veAWP = IveAWP(veAWP_);
        _disableInitializers();
    }

    /// @notice Initialize the DAO (called once via proxy)
    function initialize(
        TimelockControllerUpgradeable timelock_,
        uint48 votingDelay_,
        uint32 votingPeriod_,
        uint48 lateQuorumExtension_,
        uint256 quorumPercent_,
        address guardian_
    ) external initializer {
        __Governor_init("AWPDAO");
        __GovernorSettings_init(votingDelay_, votingPeriod_, 200_000 * 1e18);
        __GovernorTimelockControl_init(timelock_);
        __GovernorPreventLateQuorum_init(lateQuorumExtension_);
        if (quorumPercent_ == 0 || quorumPercent_ > 100) revert InvalidQuorumPercent();
        quorumPercent = quorumPercent_;
        guardian = guardian_;
    }

    /// @notice Migration: update voting parameters (called via upgradeToAndCall)
    function migrateVotingParams(uint48 newVotingDelay, uint32 newVotingPeriod) external reinitializer(2) {
        _setVotingDelay(newVotingDelay);
        _setVotingPeriod(newVotingPeriod);
    }

    /// @dev UUPS upgrade authorization — Guardian or Executor
    function _authorizeUpgrade(address) internal view override {
        if (msg.sender != guardian && msg.sender != _executor()) revert NotGuardianOrExecutor();
    }

    /// @notice Update guardian address (Guardian or Executor). Pass address(0) to permanently remove guardian.
    function setGuardian(address g) external {
        if (msg.sender != guardian && msg.sender != _executor()) revert NotGuardianOrExecutor();
        guardian = g;
        emit GuardianUpdated(g);
    }

    // ═══════════════════════════════════════════════
    //  Clock
    // ═══════════════════════════════════════════════

    function clock() public view override returns (uint48) {
        return uint48(block.timestamp);
    }

    // solhint-disable-next-line func-name-mixedcase
    function CLOCK_MODE() public pure override returns (string memory) {
        return "mode=timestamp";
    }

    // ═══════════════════════════════════════════════
    //  Custom counting module
    // ═══════════════════════════════════════════════

    // solhint-disable-next-line func-name-mixedcase
    function COUNTING_MODE() public pure override returns (string memory) {
        return "support=bravo&quorum=for,abstain&params=tokenIds";
    }

    /// @notice Always returns false — per-tokenId tracking is the real guard
    function hasVoted(uint256, address) public pure override returns (bool) {
        return false;
    }

    /// @dev Block simple vote functions — must use castVoteWithReasonAndParams (with tokenId array in params)
    function castVote(uint256, uint8) public pure override returns (uint256) {
        revert UseCastVoteWithParams();
    }

    function castVoteWithReason(uint256, uint8, string calldata) public pure override returns (uint256) {
        revert UseCastVoteWithParams();
    }

    function castVoteBySig(uint256, uint8, address, bytes memory) public pure override returns (uint256) {
        revert UseCastVoteWithParams();
    }

    // NOTE: castVoteWithReasonAndParamsBySig is NOT blocked — it supports gasless voting.
    // The voter signs EIP-712 with params = abi.encode(uint256[] tokenIds).
    // OZ Governor handles signature verification; _countVote decodes tokenIds from params.

    /// @notice Get vote tallies
    function proposalVotes(uint256 proposalId)
        public view returns (uint256 againstVotes, uint256 forVotes, uint256 abstainVotes)
    {
        ProposalVote storage pv = _proposalVotes[proposalId];
        return (pv.againstVotes, pv.forVotes, pv.abstainVotes);
    }

    /// @dev Sentinel — actual voting power computed in _countVote
    function _getVotes(address, uint256, bytes memory) internal pure override returns (uint256) {
        return 1;
    }

    /// @dev Custom vote counting: decode tokenId array, verify ownership + eligibility.
    ///      Accumulates time-weighted power (for outcome) and raw amount (for quorum).
    function _countVote(
        uint256 proposalId,
        address account,
        uint8 support,
        uint256,
        bytes memory params
    ) internal override returns (uint256) {
        ProposalVote storage proposalVote = _proposalVotes[proposalId];

        uint256[] memory tokenIds = abi.decode(params, (uint256[]));
        if (tokenIds.length == 0) revert NoTokens();

        uint256 totalPower;
        uint256 totalAmount;
        uint256 propCreatedAt = proposalCreatedAt[proposalId];

        for (uint256 i = 0; i < tokenIds.length;) {
            uint256 tid = tokenIds[i];

            (address tokenOwner, uint128 amount,, uint64 createdAt, uint64 remainingSeconds, uint256 power) =
                veAWP.getPositionForVoting(tid);

            if (tokenOwner != account) revert NotTokenOwner();
            if (hasVotedWithToken[proposalId][tid]) revert TokenAlreadyVoted();
            if (remainingSeconds < 1) revert LockExpired();
            if (createdAt >= propCreatedAt) revert MintedAfterProposal();

            hasVotedWithToken[proposalId][tid] = true;
            totalPower += power;
            totalAmount += amount;
            unchecked { ++i; }
        }

        // Time-weighted power → determines vote outcome (For vs Against)
        if (support == uint8(VoteType.Against)) {
            proposalVote.againstVotes += uint128(totalPower);
        } else if (support == uint8(VoteType.For)) {
            proposalVote.forVotes += uint128(totalPower);
        } else if (support == uint8(VoteType.Abstain)) {
            proposalVote.abstainVotes += uint128(totalPower);
        } else {
            revert GovernorInvalidVoteType();
        }

        // Raw staked amount → determines quorum (For + Abstain only)
        if (support != uint8(VoteType.Against)) {
            _proposalQuorumVotes[proposalId] += totalAmount;
        }

        return totalPower;
    }

    /// @dev Quorum based on raw staked AWP amount (not time-weighted power).
    ///      4% means 4% of total staked AWP must vote For or Abstain.
    function _quorumReached(uint256 proposalId) internal view override returns (bool) {
        uint256 q = proposalTotalVotingPower[proposalId] * quorumPercent / 100;
        return q <= _proposalQuorumVotes[proposalId];
    }

    function _voteSucceeded(uint256 proposalId) internal view override returns (bool) {
        ProposalVote storage pv = _proposalVotes[proposalId];
        return pv.forVotes > pv.againstVotes;
    }

    // ═══════════════════════════════════════════════
    //  Guardian emergency cancel
    // ═══════════════════════════════════════════════

    /// @notice Guardian can cancel any active/queued proposal (emergency safety mechanism)
    /// @param targets Proposal targets (must match original proposal)
    /// @param values Proposal values
    /// @param calldatas Proposal calldatas
    /// @param descriptionHash keccak256 of proposal description
    function guardianCancel(
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        bytes32 descriptionHash
    ) external returns (uint256) {
        if (msg.sender != guardian) revert NotGuardian();
        return _cancel(targets, values, calldatas, descriptionHash);
    }

    // ═══════════════════════════════════════════════
    //  Quorum and proposals
    // ═══════════════════════════════════════════════

    /// @notice Update quorum percentage (Executor or Guardian)
    function setQuorumPercent(uint256 newQuorumPercent) external {
        if (msg.sender != _executor() && msg.sender != guardian) revert NotGuardianOrExecutor();
        if (newQuorumPercent == 0 || newQuorumPercent > 100) revert InvalidQuorumPercent();
        quorumPercent = newQuorumPercent;
        emit QuorumPercentUpdated(newQuorumPercent);
    }

    /// @notice Quorum for a specific proposal (snapshot at creation; 0 for unknown = live estimate)
    function quorum(uint256 proposalId) public view override returns (uint256) {
        uint256 snapshot = proposalTotalVotingPower[proposalId];
        uint256 totalVP = snapshot > 0 ? snapshot : veAWP.totalVotingPower();
        return totalVP * quorumPercent / 100;
    }

    function proposalThreshold() public view override(GovernorUpgradeable, GovernorSettingsUpgradeable) returns (uint256) {
        return super.proposalThreshold();
    }

    /// @dev Standard propose reverts — use proposeWithTokens
    function propose(
        address[] memory, uint256[] memory, bytes[] memory, string memory
    ) public pure override returns (uint256) {
        revert UseProposeWithTokens();
    }

    /// @notice Propose with token proof
    function proposeWithTokens(
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        string memory description,
        uint256[] memory tokenIds
    ) public returns (uint256) {
        address proposer = _validateProposer(description, tokenIds);
        uint256 proposalId = _propose(targets, values, calldatas, description, proposer);
        _snapshotProposal(proposalId);
        return proposalId;
    }

    /// @notice Signal-only proposal (no on-chain execution)
    function signalPropose(
        string memory description,
        uint256[] memory tokenIds
    ) public returns (uint256) {
        address proposer = _validateProposer(description, tokenIds);

        address[] memory targets = new address[](1);
        targets[0] = address(this);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);

        uint256 proposalId = _propose(targets, values, calldatas, description, proposer);
        _snapshotProposal(proposalId);
        isSignalProposal[proposalId] = true;
        return proposalId;
    }

    // ═══════════════════════════════════════════════
    //  Internal
    // ═══════════════════════════════════════════════

    function _validateProposer(string memory description, uint256[] memory) internal view returns (address) {
        address proposer = _msgSender();
        if (!_isValidDescriptionForProposer(proposer, description)) {
            revert GovernorRestrictedProposer(proposer);
        }
        // Use raw staked amount (not time-weighted power) for threshold check,
        // consistent with quorum which also uses raw staked amounts.
        uint256 stakedAmount = veAWP.getUserTotalStaked(proposer);
        uint256 threshold = proposalThreshold();
        if (stakedAmount < threshold) {
            revert GovernorInsufficientProposerVotes(proposer, stakedAmount, threshold);
        }
        return proposer;
    }

    function _snapshotProposal(uint256 proposalId) internal {
        proposalCreatedAt[proposalId] = block.timestamp;
        uint256 totalVP = veAWP.totalVotingPower();
        if (totalVP == 0) revert ZeroTotalVotingPower();
        proposalTotalVotingPower[proposalId] = totalVP;
    }

    // ═══════════════════════════════════════════════
    //  Required overrides (diamond resolution)
    // ═══════════════════════════════════════════════

    function votingDelay() public view override(GovernorUpgradeable, GovernorSettingsUpgradeable) returns (uint256) {
        return super.votingDelay();
    }

    function votingPeriod() public view override(GovernorUpgradeable, GovernorSettingsUpgradeable) returns (uint256) {
        return super.votingPeriod();
    }

    function state(uint256 proposalId)
        public view override(GovernorUpgradeable, GovernorTimelockControlUpgradeable) returns (ProposalState)
    {
        return super.state(proposalId);
    }

    function proposalNeedsQueuing(uint256 proposalId)
        public view override(GovernorUpgradeable, GovernorTimelockControlUpgradeable) returns (bool)
    {
        if (isSignalProposal[proposalId]) return false;
        return super.proposalNeedsQueuing(proposalId);
    }

    function proposalDeadline(uint256 proposalId)
        public view override(GovernorUpgradeable, GovernorPreventLateQuorumUpgradeable) returns (uint256)
    {
        return super.proposalDeadline(proposalId);
    }

    function _tallyUpdated(uint256 proposalId)
        internal override(GovernorUpgradeable, GovernorPreventLateQuorumUpgradeable)
    {
        super._tallyUpdated(proposalId);
    }

    function _queueOperations(
        uint256 proposalId,
        address[] memory targets, uint256[] memory values,
        bytes[] memory calldatas, bytes32 descriptionHash
    ) internal override(GovernorUpgradeable, GovernorTimelockControlUpgradeable) returns (uint48) {
        if (isSignalProposal[proposalId]) return 0;
        return super._queueOperations(proposalId, targets, values, calldatas, descriptionHash);
    }

    function _executeOperations(
        uint256 proposalId,
        address[] memory targets, uint256[] memory values,
        bytes[] memory calldatas, bytes32 descriptionHash
    ) internal override(GovernorUpgradeable, GovernorTimelockControlUpgradeable) {
        if (isSignalProposal[proposalId]) return;
        super._executeOperations(proposalId, targets, values, calldatas, descriptionHash);
    }

    function _cancel(
        address[] memory targets, uint256[] memory values,
        bytes[] memory calldatas, bytes32 descriptionHash
    ) internal override(GovernorUpgradeable, GovernorTimelockControlUpgradeable) returns (uint256) {
        return super._cancel(targets, values, calldatas, descriptionHash);
    }

    function _executor() internal view override(GovernorUpgradeable, GovernorTimelockControlUpgradeable) returns (address) {
        return super._executor();
    }
}
