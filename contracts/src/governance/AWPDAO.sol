// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Governor} from "@openzeppelin/contracts/governance/Governor.sol";
import {GovernorSettings} from "@openzeppelin/contracts/governance/extensions/GovernorSettings.sol";
import {GovernorTimelockControl} from "@openzeppelin/contracts/governance/extensions/GovernorTimelockControl.sol";
import {TimelockController} from "@openzeppelin/contracts/governance/TimelockController.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IStakeNFT} from "../interfaces/IStakeNFT.sol";

/// @title AWPDAO — Governance contract with NFT-based voting
/// @notice Uses StakeNFT positions for voting power. Voters submit tokenId arrays.
///         Replaces GovernorVotes + GovernorCountingSimple with a custom counting module.
/// @dev Inheritance chain: Governor, GovernorSettings, GovernorTimelockControl.
///      _getVotes returns 1 (sentinel); actual power is computed in _countVote.
///      clock() uses block.number (default). No GovernorVotes dependency.
///      Anti-manipulation uses createdAt timestamp instead of epoch.
contract AWPDAO is
    Governor,
    GovernorSettings,
    GovernorTimelockControl
{
    // ── Vote types (matches Governor Bravo ordering) ──
    enum VoteType {
        Against,
        For,
        Abstain
    }

    // ── Vote counting storage ──
    struct ProposalVote {
        uint256 againstVotes;
        uint256 forVotes;
        uint256 abstainVotes;
    }

    /// @notice StakeNFT contract for position queries
    IStakeNFT public immutable stakeNFT;

    /// @notice AWP token contract (for quorum calculation)
    IERC20 public immutable awpToken;

    /// @notice Quorum percentage (e.g. 4 = 4% of total voting power)
    uint256 public quorumPercent;

    /// @notice Double-vote prevention per tokenId per proposal
    mapping(uint256 proposalId => mapping(uint256 tokenId => bool)) public hasVotedWithToken;

    /// @notice Timestamp at proposal creation time (for anti-manipulation check)
    mapping(uint256 proposalId => uint256) public proposalCreatedAt;

    /// @notice Whether a proposal is signal-only (no on-chain execution, skips Timelock)
    mapping(uint256 proposalId => bool) public isSignalProposal;

    /// @notice Snapshot of totalVotingPower at proposal creation (for quorum calculation)
    mapping(uint256 proposalId => uint256) public proposalTotalVotingPower;

    /// @dev Vote tallies per proposal
    mapping(uint256 proposalId => ProposalVote) private _proposalVotes;

    // ── Errors ──
    error NoTokens();
    error NotTokenOwner();
    error TokenAlreadyVoted();
    error LockExpired();
    error MintedAfterProposal();
    error InsufficientVotingPower();
    error UseProposeWithTokens();

    /// @notice Constructor
    /// @param stakeNFT_ StakeNFT contract address
    /// @param awpToken_ AWP token contract address
    /// @param timelock_ TimelockController (Treasury) contract
    /// @param votingDelay_ Delay from proposal creation to voting start (blocks)
    /// @param votingPeriod_ Voting duration (blocks)
    /// @param quorumPercent_ Quorum percentage (e.g. 4 for 4%)
    constructor(
        address stakeNFT_,
        address awpToken_,
        TimelockController timelock_,
        uint48 votingDelay_,
        uint32 votingPeriod_,
        uint256 quorumPercent_
    )
        Governor("AWPDAO")
        GovernorSettings(votingDelay_, votingPeriod_, 1_000_000 * 1e18)
        GovernorTimelockControl(timelock_)
    {
        stakeNFT = IStakeNFT(stakeNFT_);
        awpToken = IERC20(awpToken_);
        quorumPercent = quorumPercent_;
    }

    // ═══════════════════════════════════════════════
    //  Clock (replaces GovernorVotes clock)
    // ═══════════════════════════════════════════════

    /// @notice Returns the current block number as the clock value
    function clock() public view override returns (uint48) {
        return uint48(block.number);
    }

    // solhint-disable-next-line func-name-mixedcase
    function CLOCK_MODE() public pure override returns (string memory) {
        return "mode=blocknumber&from=default";
    }

    // ═══════════════════════════════════════════════
    //  Custom counting module
    // ═══════════════════════════════════════════════

    // solhint-disable-next-line func-name-mixedcase
    function COUNTING_MODE() public pure override returns (string memory) {
        return "support=bravo&quorum=for,abstain&params=tokenIds";
    }

    /// @notice Always returns false — per-tokenId tracking in hasVotedWithToken is the real guard.
    ///         Returning false bypasses OZ Governor's per-account check so users can submit
    ///         multiple castVoteWithParams calls with different tokenId batches.
    function hasVoted(uint256, address) public pure override returns (bool) {
        return false;
    }

    /// @dev Block castVote (no params) — users must use castVoteWithReasonAndParams with tokenId array
    error UsecastVoteWithParams();
    error InvalidQuorumPercent();
    error ZeroTotalVotingPower();

    function castVote(uint256, uint8) public pure override returns (uint256) {
        revert UsecastVoteWithParams();
    }

    function castVoteWithReason(uint256, uint8, string calldata) public pure override returns (uint256) {
        revert UsecastVoteWithParams();
    }

    function castVoteBySig(uint256, uint8, address, bytes memory) public pure override returns (uint256) {
        revert UsecastVoteWithParams();
    }

    function castVoteWithReasonAndParamsBySig(uint256, uint8, address, string calldata, bytes memory, bytes memory) public pure override returns (uint256) {
        revert UsecastVoteWithParams();
    }

    /// @notice Get vote tallies for a proposal
    function proposalVotes(uint256 proposalId)
        public
        view
        returns (uint256 againstVotes, uint256 forVotes, uint256 abstainVotes)
    {
        ProposalVote storage pv = _proposalVotes[proposalId];
        return (pv.againstVotes, pv.forVotes, pv.abstainVotes);
    }

    /// @dev Sentinel value — actual voting power is computed in _countVote from tokenId params
    function _getVotes(address, uint256, bytes memory) internal pure override returns (uint256) {
        return 1;
    }

    /// @dev Custom vote counting: decode tokenId array from params, verify ownership + eligibility,
    ///      accumulate voting power. Each tokenId can only vote once per proposal.
    function _countVote(
        uint256 proposalId,
        address account,
        uint8 support,
        uint256, // totalWeight (sentinel, ignored)
        bytes memory params
    ) internal override returns (uint256) {
        ProposalVote storage proposalVote = _proposalVotes[proposalId];

        // Allow multiple castVoteWithParams calls from the same account (for different tokenIds)
        // The per-account hasVoted is set on first call; subsequent calls with new tokenIds are still allowed
        // because we track per-tokenId voting in hasVotedWithToken
        uint256[] memory tokenIds = abi.decode(params, (uint256[]));
        if (tokenIds.length == 0) revert NoTokens();

        uint256 totalPower = 0;
        uint256 propCreatedAt = proposalCreatedAt[proposalId];

        for (uint256 i = 0; i < tokenIds.length;) {
            uint256 tid = tokenIds[i];

            // Single external call replaces 3 separate calls (positions + remainingTime + getVotingPower)
            (address tokenOwner,,, uint64 createdAt, uint64 remainingSeconds, uint256 power) =
                stakeNFT.getPositionForVoting(tid);

            // Verify caller owns the token
            if (tokenOwner != account) revert NotTokenOwner();

            // Prevent double-voting with the same token
            if (hasVotedWithToken[proposalId][tid]) revert TokenAlreadyVoted();

            // Token must still have remaining lock
            if (remainingSeconds < 1) revert LockExpired();

            // Anti-manipulation: only NFTs minted before the proposal can vote
            if (createdAt >= propCreatedAt) revert MintedAfterProposal();

            hasVotedWithToken[proposalId][tid] = true;
            totalPower += power;
            unchecked { ++i; }
        }

        // Accumulate votes by support type
        if (support == uint8(VoteType.Against)) {
            proposalVote.againstVotes += totalPower;
        } else if (support == uint8(VoteType.For)) {
            proposalVote.forVotes += totalPower;
        } else if (support == uint8(VoteType.Abstain)) {
            proposalVote.abstainVotes += totalPower;
        } else {
            revert GovernorInvalidVoteType();
        }

        return totalPower;
    }

    /// @dev Check if quorum is reached for a proposal (uses snapshot of totalVotingPower at proposal creation)
    function _quorumReached(uint256 proposalId) internal view override returns (bool) {
        ProposalVote storage pv = _proposalVotes[proposalId];
        uint256 q = proposalTotalVotingPower[proposalId] * quorumPercent / 100;
        return q <= pv.forVotes + pv.abstainVotes;
    }

    /// @dev Check if a proposal vote succeeded (forVotes > againstVotes)
    function _voteSucceeded(uint256 proposalId) internal view override returns (bool) {
        ProposalVote storage pv = _proposalVotes[proposalId];
        return pv.forVotes > pv.againstVotes;
    }

    // ═══════════════════════════════════════════════
    //  Quorum and threshold
    // ═══════════════════════════════════════════════

    /// @notice Update quorum percentage (only via governance through Timelock)
    function setQuorumPercent(uint256 newQuorumPercent) external onlyGovernance {
        if (newQuorumPercent == 0 || newQuorumPercent > 100) revert InvalidQuorumPercent();
        quorumPercent = newQuorumPercent;
    }

    /// @notice Quorum = totalVotingPower * quorumPercent / 100
    /// @dev Uses stakeNFT.totalVotingPower() which only counts properly deposited AWP,
    ///      not accidentally-sent tokens.
    function quorum(uint256) public view override returns (uint256) {
        return stakeNFT.totalVotingPower() * quorumPercent / 100;
    }

    /// @notice Minimum voting power required to create a proposal (configurable via GovernorSettings)
    function proposalThreshold() public view override(Governor, GovernorSettings) returns (uint256) {
        return super.proposalThreshold();
    }

    /// @dev Standard propose reverts — use proposeWithTokens instead
    function propose(
        address[] memory,
        uint256[] memory,
        bytes[] memory,
        string memory
    ) public pure override returns (uint256) {
        revert UseProposeWithTokens();
    }

    /// @notice Propose with token proof: caller must supply their tokenIds to prove threshold
    /// @param targets Target contract addresses
    /// @param values ETH values
    /// @param calldatas Encoded function calls
    /// @param description Proposal description
    /// @param tokenIds Proposer's StakeNFT tokenIds to prove voting power >= proposalThreshold
    function proposeWithTokens(
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        string memory description,
        uint256[] memory tokenIds
    ) public returns (uint256) {
        address proposer = _msgSender();

        // Check description restriction
        if (!_isValidDescriptionForProposer(proposer, description)) {
            revert GovernorRestrictedProposer(proposer);
        }

        // Check proposer has sufficient voting power via StakeNFT
        uint256 votingPower = stakeNFT.getUserVotingPower(proposer, tokenIds);
        uint256 threshold = proposalThreshold();
        if (votingPower < threshold) {
            revert GovernorInsufficientProposerVotes(proposer, votingPower, threshold);
        }

        uint256 proposalId = _propose(targets, values, calldatas, description, proposer);
        proposalCreatedAt[proposalId] = block.timestamp;
        uint256 totalVP = stakeNFT.totalVotingPower();
        if (totalVP == 0) revert ZeroTotalVotingPower();
        proposalTotalVotingPower[proposalId] = totalVP;
        return proposalId;
    }

    /// @notice Create a signal-only proposal (no on-chain execution, skips Timelock)
    /// @dev Uses a single no-op target (OZ Governor requires targets.length > 0).
    ///      After voting succeeds, call execute() with matching arrays to finalize.
    ///      Both _queueOperations and _executeOperations are overridden to skip Timelock for signal proposals.
    /// @param description Proposal description (the content being voted on)
    /// @param tokenIds Proposer's StakeNFT tokenIds to prove voting power >= proposalThreshold
    function signalPropose(
        string memory description,
        uint256[] memory tokenIds
    ) public returns (uint256) {
        address proposer = _msgSender();

        // Check description restriction (consistent with proposeWithTokens)
        if (!_isValidDescriptionForProposer(proposer, description)) {
            revert GovernorRestrictedProposer(proposer);
        }

        // Check proposer has sufficient voting power
        uint256 votingPower = stakeNFT.getUserVotingPower(proposer, tokenIds);
        uint256 threshold = proposalThreshold();
        if (votingPower < threshold) {
            revert GovernorInsufficientProposerVotes(proposer, votingPower, threshold);
        }

        // Create proposal with a single no-op target (OZ Governor requires targets.length > 0)
        // _executeOperations is overridden to skip execution for signal proposals
        address[] memory targets = new address[](1);
        targets[0] = address(this);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);

        uint256 proposalId = _propose(targets, values, calldatas, description, proposer);
        proposalCreatedAt[proposalId] = block.timestamp;
        uint256 totalVP = stakeNFT.totalVotingPower();
        if (totalVP == 0) revert ZeroTotalVotingPower();
        proposalTotalVotingPower[proposalId] = totalVP;
        isSignalProposal[proposalId] = true;

        return proposalId;
    }

    // ═══════════════════════════════════════════════
    //  Internal helpers
    // ═══════════════════════════════════════════════

    // ═══════════════════════════════════════════════
    //  Required overrides (resolve multiple-inheritance diamond)
    // ═══════════════════════════════════════════════

    function votingDelay() public view override(Governor, GovernorSettings) returns (uint256) {
        return super.votingDelay();
    }

    function votingPeriod() public view override(Governor, GovernorSettings) returns (uint256) {
        return super.votingPeriod();
    }

    function state(uint256 proposalId)
        public
        view
        override(Governor, GovernorTimelockControl)
        returns (ProposalState)
    {
        return super.state(proposalId);
    }

    /// @dev Signal proposals skip Timelock queuing; executable proposals go through Timelock
    function proposalNeedsQueuing(uint256 proposalId)
        public
        view
        override(Governor, GovernorTimelockControl)
        returns (bool)
    {
        if (isSignalProposal[proposalId]) return false;
        return super.proposalNeedsQueuing(proposalId);
    }

    function _queueOperations(
        uint256 proposalId,
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        bytes32 descriptionHash
    ) internal override(Governor, GovernorTimelockControl) returns (uint48) {
        if (isSignalProposal[proposalId]) {
            // Signal proposals: no Timelock scheduling, return 0 (no eta)
            return 0;
        }
        return super._queueOperations(proposalId, targets, values, calldatas, descriptionHash);
    }

    function _executeOperations(
        uint256 proposalId,
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        bytes32 descriptionHash
    ) internal override(Governor, GovernorTimelockControl) {
        if (isSignalProposal[proposalId]) {
            // Signal proposals: no-op execution (empty targets, skip Timelock)
            return;
        }
        super._executeOperations(proposalId, targets, values, calldatas, descriptionHash);
    }

    function _cancel(
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        bytes32 descriptionHash
    ) internal override(Governor, GovernorTimelockControl) returns (uint256) {
        return super._cancel(targets, values, calldatas, descriptionHash);
    }

    function _executor() internal view override(Governor, GovernorTimelockControl) returns (address) {
        return super._executor();
    }
}
