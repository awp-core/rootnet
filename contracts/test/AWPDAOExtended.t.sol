// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {IGovernor} from "@openzeppelin/contracts/governance/IGovernor.sol";
import {TimelockController} from "@openzeppelin/contracts/governance/TimelockController.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {StakingVault} from "../src/core/StakingVault.sol";

/// @title AWPDAOExtended — Supplementary edge case and integration tests for AWPDAO
contract AWPDAOExtendedTest is Test {
    AWPDAO public dao;
    Treasury public treasury;
    AWPToken public awpToken;
    StakeNFT public stakeNFT;
    StakingVault public vault;

    address public deployer = makeAddr("deployer");
    address public guardian = makeAddr("guardian");
    address public voter = makeAddr("voter");
    address public voter2 = makeAddr("voter2");

    uint256 public constant MIN_DELAY = 1 days;
    uint48 public constant VOTING_DELAY = 1; // 1 block
    uint32 public constant VOTING_PERIOD = 100; // 100 blocks (fast for tests)
    uint256 public constant QUORUM_PERCENTAGE = 4; // 4%

    // ── Reusable proposal params ──
    address[] internal _targets;
    uint256[] internal _values;
    bytes[] internal _calldatas;

    function setUp() public {
        vm.startPrank(deployer);

        // Deploy AWPToken
        awpToken = new AWPToken("AWP", "AWP", deployer);
        awpToken.initialMint(200_000_000 * 1e18);

        // Deploy StakingVault + StakeNFT
        vault = StakingVault(address(new ERC1967Proxy(
            address(new StakingVault()), abi.encodeCall(StakingVault.initialize, (address(this), address(this)))
        )));
        stakeNFT = new StakeNFT(address(awpToken), address(vault), address(this));
        vm.stopPrank();
        vault.setStakeNFT(address(stakeNFT));
        vm.startPrank(deployer);

        // Deploy Treasury
        address[] memory proposers = new address[](0);
        address[] memory executors = new address[](1);
        executors[0] = address(0); // anyone can execute
        treasury = new Treasury(MIN_DELAY, proposers, executors, deployer);

        // Deploy AWPDAO
        dao = new AWPDAO(
            address(stakeNFT),
            address(awpToken),
            TimelockController(payable(address(treasury))),
            VOTING_DELAY,
            VOTING_PERIOD,
            QUORUM_PERCENTAGE
        );

        // Grant DAO proposer and canceller roles
        treasury.grantRole(treasury.PROPOSER_ROLE(), address(dao));
        treasury.grantRole(treasury.CANCELLER_ROLE(), address(dao));

        // Transfer Treasury admin to guardian, deployer renounces
        treasury.grantRole(treasury.DEFAULT_ADMIN_ROLE(), guardian);
        treasury.renounceRole(treasury.DEFAULT_ADMIN_ROLE(), deployer);

        // Give voter tokens and stake (large amount, ensures quorum)
        uint256 voterAmount = 100_000_000 * 1e18;
        awpToken.transfer(voter, voterAmount);

        // Give voter2 tokens and stake (small amount)
        uint256 voter2Amount = 1_000_000 * 1e18;
        awpToken.transfer(voter2, voter2Amount);

        vm.stopPrank();

        // voter stakes — tokenId = 1
        vm.startPrank(voter);
        awpToken.approve(address(stakeNFT), voterAmount);
        stakeNFT.deposit(voterAmount, 52 weeks);
        vm.stopPrank();

        // voter2 stakes — tokenId = 2
        vm.startPrank(voter2);
        awpToken.approve(address(stakeNFT), voter2Amount);
        stakeNFT.deposit(voter2Amount, 52 weeks);
        vm.stopPrank();

        // Initialize reusable no-op proposal params
        _targets = new address[](1);
        _targets[0] = address(treasury);
        _values = new uint256[](1);
        _calldatas = new bytes[](1);
    }

    // ═══════════════════════════════════════════════
    //  Helper functions
    // ═══════════════════════════════════════════════

    /// @dev Create proposal and advance to voting period
    function _createAndActivateProposal(string memory desc) internal returns (uint256 proposalId) {
        // Ensure position createdAt < proposalCreatedAt
        vm.warp(block.timestamp + 1 days);
        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        uint256[] memory proposerTokenIds = new uint256[](1);
        proposerTokenIds[0] = 1;
        vm.prank(voter);
        proposalId = dao.proposeWithTokens(_targets, _values, _calldatas, desc, proposerTokenIds);

        // Advance to voting period
        vm.roll(block.number + VOTING_DELAY + 1);
        vm.warp(block.timestamp + (VOTING_DELAY + 1) * 12);
    }

    /// @dev voter casts for-vote on proposal
    function _voteFor(uint256 proposalId) internal {
        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = 1;
        vm.prank(voter);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));
    }

    /// @dev Advance past voting period end
    function _advancePastVoting() internal {
        vm.roll(block.number + VOTING_PERIOD + 1);
        vm.warp(block.timestamp + (VOTING_PERIOD + 1) * 12);
    }

    // ═══════════════════════════════════════════════
    //  1. Treasury Guardian emergency backdoor
    // ═══════════════════════════════════════════════

    function test_guardianCanGrantProposerRole() public {
        address newDAO = makeAddr("newDAO");
        bytes32 proposerRole = treasury.PROPOSER_ROLE();
        vm.prank(guardian);
        treasury.grantRole(proposerRole, newDAO);
        assertTrue(treasury.hasRole(proposerRole, newDAO));
    }

    function test_guardianCanRevokeProposerRole() public {
        bytes32 proposerRole = treasury.PROPOSER_ROLE();
        vm.prank(guardian);
        treasury.revokeRole(proposerRole, address(dao));
        assertFalse(treasury.hasRole(proposerRole, address(dao)));
    }

    function test_nonGuardianCannotGrantRoles() public {
        address rando = makeAddr("rando");
        bytes32 proposerRole = treasury.PROPOSER_ROLE();
        vm.prank(rando);
        vm.expectRevert();
        treasury.grantRole(proposerRole, rando);
    }

    // ═══════════════════════════════════════════════
    //  2. Proposal edge cases
    // ═══════════════════════════════════════════════

    function test_proposeWithTokens_insufficientVotingPower() public {
        // voter2 only has 1M AWP, below proposalThreshold
        // Need a user with voting power < threshold
        address tinyVoter = makeAddr("tinyVoter");
        vm.prank(deployer);
        awpToken.transfer(tinyVoter, 100 * 1e18); // 100 AWP — far below threshold

        vm.startPrank(tinyVoter);
        awpToken.approve(address(stakeNFT), 100 * 1e18);
        stakeNFT.deposit(100 * 1e18, 52 weeks); // tokenId = 3
        vm.stopPrank();

        vm.warp(block.timestamp + 1 days);
        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = 3;

        vm.prank(tinyVoter);
        vm.expectRevert(); // GovernorInsufficientProposerVotes
        dao.proposeWithTokens(_targets, _values, _calldatas, "tiny proposal", tokenIds);
    }

    function test_proposeWithTokens_expiredLock() public {
        // Create a short-lock position and wait for expiry
        address shortVoter = makeAddr("shortVoter");
        vm.prank(deployer);
        awpToken.transfer(shortVoter, 50_000_000 * 1e18);

        vm.startPrank(shortVoter);
        awpToken.approve(address(stakeNFT), 50_000_000 * 1e18);
        stakeNFT.deposit(50_000_000 * 1e18, 1 days); // tokenId = 3, min lock
        vm.stopPrank();

        // Wait for lock to expire
        vm.warp(block.timestamp + 2 days);
        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = 3;

        // After lock expiry, voting power = 0 (remainingTime = 0), below threshold
        vm.prank(shortVoter);
        vm.expectRevert(); // GovernorInsufficientProposerVotes (voting power = 0)
        dao.proposeWithTokens(_targets, _values, _calldatas, "expired lock proposal", tokenIds);
    }

    function test_proposeWithTokens_mintedAfterProposal_voteReverts() public {
        // Create proposal first, then mint NFT, then try voting with new NFT -> MintedAfterProposal
        uint256 proposalId = _createAndActivateProposal("anti-manipulation test");

        // Now create a new stake position (createdAt >= proposalCreatedAt)
        address lateMinter = makeAddr("lateMinter");
        vm.prank(deployer);
        awpToken.transfer(lateMinter, 10_000_000 * 1e18);
        vm.startPrank(lateMinter);
        awpToken.approve(address(stakeNFT), 10_000_000 * 1e18);
        stakeNFT.deposit(10_000_000 * 1e18, 52 weeks); // tokenId = 3
        vm.stopPrank();

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = 3;

        vm.prank(lateMinter);
        vm.expectRevert(AWPDAO.MintedAfterProposal.selector);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));
    }

    function test_propose_standardBlocked() public {
        vm.expectRevert(AWPDAO.UseProposeWithTokens.selector);
        dao.propose(_targets, _values, _calldatas, "blocked");
    }

    // ═══════════════════════════════════════════════
    //  3. Voting edge cases
    // ═══════════════════════════════════════════════

    function test_doubleVoteWithSameTokenId() public {
        uint256 proposalId = _createAndActivateProposal("double vote test");

        // First vote
        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = 1;
        vm.prank(voter);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));

        // Second vote with same tokenId -> TokenAlreadyVoted
        vm.prank(voter);
        vm.expectRevert(AWPDAO.TokenAlreadyVoted.selector);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));
    }

    function test_voteWithTokenNotOwned() public {
        uint256 proposalId = _createAndActivateProposal("not owner test");

        // voter2 tries to use voter's tokenId=1 -> NotTokenOwner
        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = 1; // Owned by voter, not voter2
        vm.prank(voter2);
        vm.expectRevert(AWPDAO.NotTokenOwner.selector);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));
    }

    function test_voteWithExpiredLock() public {
        // Create short-lock position
        address shortVoter = makeAddr("shortVoter");
        vm.prank(deployer);
        awpToken.transfer(shortVoter, 50_000_000 * 1e18);
        vm.startPrank(shortVoter);
        awpToken.approve(address(stakeNFT), 50_000_000 * 1e18);
        stakeNFT.deposit(50_000_000 * 1e18, 1 days); // tokenId = 3
        vm.stopPrank();

        // Wait so createdAt < proposalCreatedAt but lock not yet expired
        vm.warp(block.timestamp + 12 hours);
        vm.roll(block.number + 1);

        // Create proposal (tokenId=3 createdAt < block.timestamp)
        uint256[] memory proposerTokenIds = new uint256[](1);
        proposerTokenIds[0] = 1;
        vm.prank(voter);
        uint256 proposalId = dao.proposeWithTokens(_targets, _values, _calldatas, "expired lock vote", proposerTokenIds);

        // Advance to voting period, also letting shortVoter lock expire
        vm.roll(block.number + VOTING_DELAY + 1);
        vm.warp(block.timestamp + 2 days); // Past 1-day lock

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = 3;
        vm.prank(shortVoter);
        vm.expectRevert(AWPDAO.LockExpired.selector);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));
    }

    function test_voteWithTokenMintedAfterProposal() public {
        uint256 proposalId = _createAndActivateProposal("minted after test");

        // Mint new position after proposal creation
        address lateMinter = makeAddr("lateMinter");
        vm.prank(deployer);
        awpToken.transfer(lateMinter, 10_000_000 * 1e18);
        vm.startPrank(lateMinter);
        awpToken.approve(address(stakeNFT), 10_000_000 * 1e18);
        stakeNFT.deposit(10_000_000 * 1e18, 52 weeks); // tokenId = 3
        vm.stopPrank();

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = 3;
        vm.prank(lateMinter);
        vm.expectRevert(AWPDAO.MintedAfterProposal.selector);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));
    }

    function test_batchVotingWithMultipleTokenIds() public {
        // Create second stake position for voter
        vm.prank(deployer);
        awpToken.transfer(voter, 10_000_000 * 1e18);
        vm.startPrank(voter);
        awpToken.approve(address(stakeNFT), 10_000_000 * 1e18);
        stakeNFT.deposit(10_000_000 * 1e18, 52 weeks); // tokenId = 3
        vm.stopPrank();

        uint256 proposalId = _createAndActivateProposal("batch vote test");

        // Vote with multiple tokenIds at once
        uint256[] memory tokenIds = new uint256[](2);
        tokenIds[0] = 1;
        tokenIds[1] = 3;
        vm.prank(voter);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));

        // Confirm voting power is sum of both positions
        (, uint256 forVotes,) = dao.proposalVotes(proposalId);
        assertGt(forVotes, 0, "batch voting power should be > 0");

        // Confirm single tokenId has less power
        uint256 singlePower = stakeNFT.getVotingPower(1);
        uint256 secondPower = stakeNFT.getVotingPower(3);
        assertEq(forVotes, singlePower + secondPower, "batch power should equal sum of individual powers");
    }

    // ═══════════════════════════════════════════════
    //  4. Quorum
    // ═══════════════════════════════════════════════

    function test_setQuorumPercentViaGovernance() public {
        // Construct setQuorumPercent call as proposal content
        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(AWPDAO.setQuorumPercent, (10));
        string memory desc = "Set quorum to 10%";

        // Ensure createdAt < proposalCreatedAt
        vm.warp(block.timestamp + 1 days);
        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        uint256[] memory proposerTokenIds = new uint256[](1);
        proposerTokenIds[0] = 1;
        vm.prank(voter);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, desc, proposerTokenIds);

        // Advance to voting period
        vm.roll(block.number + VOTING_DELAY + 1);
        vm.warp(block.timestamp + (VOTING_DELAY + 1) * 12);

        // Vote
        _voteFor(proposalId);

        // End voting
        _advancePastVoting();
        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Succeeded));

        // Queue
        bytes32 descHash = keccak256(bytes(desc));
        dao.queue(targets, values, calldatas, descHash);

        // Wait for timelock delay
        vm.warp(block.timestamp + MIN_DELAY + 1);

        // Execute
        dao.execute(targets, values, calldatas, descHash);
        assertEq(dao.quorumPercent(), 10, "quorum should be updated to 10%");
    }

    function test_proposalBarelyReachesQuorum() public {
        // 4% quorum. totalVotingPower = balanceOf(stakeNFT) * 7.
        // voter has 100M, voter2 has 1M. voter alone easily exceeds 4%.
        uint256 proposalId = _createAndActivateProposal("barely quorum");

        // Only voter votes (100M)
        _voteFor(proposalId);
        _advancePastVoting();

        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Succeeded), "should succeed with quorum");
    }

    function test_proposalMissesQuorum() public {
        // Only voter2 votes (1M / total ~101M = ~1%, below 4% quorum)
        vm.warp(block.timestamp + 1 days);
        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        uint256[] memory proposerTokenIds = new uint256[](1);
        proposerTokenIds[0] = 1;
        vm.prank(voter);
        uint256 proposalId = dao.proposeWithTokens(_targets, _values, _calldatas, "miss quorum", proposerTokenIds);

        vm.roll(block.number + VOTING_DELAY + 1);
        vm.warp(block.timestamp + (VOTING_DELAY + 1) * 12);

        // Only voter2 votes
        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = 2;
        vm.prank(voter2);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));

        _advancePastVoting();

        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Defeated), "should be defeated without quorum");
    }

    // ═══════════════════════════════════════════════
    //  5. Timelock integration
    // ═══════════════════════════════════════════════

    function test_queueAndExecuteWithDelay() public {
        // Send ETH from Treasury
        address payable target = payable(makeAddr("timelockTarget"));
        vm.deal(address(treasury), 1 ether);

        address[] memory targets = new address[](1);
        targets[0] = target;
        uint256[] memory values = new uint256[](1);
        values[0] = 1 ether;
        bytes[] memory calldatas = new bytes[](1);
        string memory desc = "Timelock delay test";

        vm.warp(block.timestamp + 1 days);
        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        uint256[] memory proposerTokenIds = new uint256[](1);
        proposerTokenIds[0] = 1;
        vm.prank(voter);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, desc, proposerTokenIds);

        vm.roll(block.number + VOTING_DELAY + 1);
        vm.warp(block.timestamp + (VOTING_DELAY + 1) * 12);
        _voteFor(proposalId);
        _advancePastVoting();

        bytes32 descHash = keccak256(bytes(desc));
        dao.queue(targets, values, calldatas, descHash);
        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Queued));

        vm.warp(block.timestamp + MIN_DELAY + 1);
        dao.execute(targets, values, calldatas, descHash);
        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Executed));
        assertEq(target.balance, 1 ether);
    }

    function test_executeBeforeDelayReverts() public {
        address payable target = payable(makeAddr("earlyExecTarget"));
        vm.deal(address(treasury), 1 ether);

        address[] memory targets = new address[](1);
        targets[0] = target;
        uint256[] memory values = new uint256[](1);
        values[0] = 1 ether;
        bytes[] memory calldatas = new bytes[](1);
        string memory desc = "Early execute test";

        vm.warp(block.timestamp + 1 days);
        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        uint256[] memory proposerTokenIds = new uint256[](1);
        proposerTokenIds[0] = 1;
        vm.prank(voter);
        dao.proposeWithTokens(targets, values, calldatas, desc, proposerTokenIds);

        vm.roll(block.number + VOTING_DELAY + 1);
        vm.warp(block.timestamp + (VOTING_DELAY + 1) * 12);
        _voteFor(dao.hashProposal(targets, values, calldatas, keccak256(bytes(desc))));
        _advancePastVoting();

        bytes32 descHash = keccak256(bytes(desc));
        dao.queue(targets, values, calldatas, descHash);

        // Execute without waiting for delay -> should revert
        vm.expectRevert();
        dao.execute(targets, values, calldatas, descHash);
    }

    function test_cancelQueuedProposal() public {
        address payable target = payable(makeAddr("cancelTarget"));
        vm.deal(address(treasury), 1 ether);

        address[] memory targets = new address[](1);
        targets[0] = target;
        uint256[] memory values = new uint256[](1);
        values[0] = 1 ether;
        bytes[] memory calldatas = new bytes[](1);
        string memory desc = "Cancel test";

        vm.warp(block.timestamp + 1 days);
        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        uint256[] memory proposerTokenIds = new uint256[](1);
        proposerTokenIds[0] = 1;
        vm.prank(voter);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, desc, proposerTokenIds);

        vm.roll(block.number + VOTING_DELAY + 1);
        vm.warp(block.timestamp + (VOTING_DELAY + 1) * 12);
        _voteFor(proposalId);
        _advancePastVoting();

        bytes32 descHash = keccak256(bytes(desc));
        dao.queue(targets, values, calldatas, descHash);
        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Queued));

        // Cancel directly on Treasury (DAO has CANCELLER_ROLE)
        // salt = bytes20(dao) ^ descHash (GovernorTimelockControl._timelockSalt)
        bytes32 salt = bytes20(address(dao)) ^ descHash;
        bytes32 opId = treasury.hashOperationBatch(targets, values, calldatas, 0, salt);
        vm.prank(address(dao));
        treasury.cancel(opId);

        // Proposal should become Canceled
        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Canceled));
    }

    // ═══════════════════════════════════════════════
    //  6. hasVoted behavior
    // ═══════════════════════════════════════════════

    function test_hasVotedAlwaysReturnsFalse() public {
        uint256 proposalId = _createAndActivateProposal("hasVoted test");
        _voteFor(proposalId);

        // hasVoted always returns false (by design: per-tokenId tracking, not per-account)
        assertFalse(dao.hasVoted(proposalId, voter), "hasVoted should always return false by design");
    }

    function test_hasVotedWithToken_tracksPerTokenId() public {
        uint256 proposalId = _createAndActivateProposal("per-token tracking test");
        _voteFor(proposalId);

        // hasVotedWithToken tracks per-tokenId voting state
        assertTrue(dao.hasVotedWithToken(proposalId, 1), "tokenId 1 should be marked as voted");
        assertFalse(dao.hasVotedWithToken(proposalId, 2), "tokenId 2 should not be marked as voted");
    }

    // ═══════════════════════════════════════════════
    //  Extra: castVote / castVoteWithReason blocked
    // ═══════════════════════════════════════════════

    function test_castVote_reverts() public {
        vm.expectRevert(AWPDAO.UseCastVoteWithParams.selector);
        dao.castVote(0, 1);
    }

    function test_castVoteWithReason_reverts() public {
        vm.expectRevert(AWPDAO.UseCastVoteWithParams.selector);
        dao.castVoteWithReason(0, 1, "reason");
    }

    // ═══════════════════════════════════════════════
    //  Extra: setQuorumPercent boundary values
    // ═══════════════════════════════════════════════

    function test_setQuorumPercent_directCallReverts() public {
        // setQuorumPercent can only be called via governance
        vm.prank(deployer);
        vm.expectRevert();
        dao.setQuorumPercent(10);
    }
}
