// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {DeployHelper} from "./helpers/DeployHelper.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {TimelockControllerUpgradeable} from "@openzeppelin/contracts-upgradeable/governance/TimelockControllerUpgradeable.sol";
import {IGovernor} from "@openzeppelin/contracts/governance/IGovernor.sol";

contract AWPDAOTest is DeployHelper {
    uint256 public tokenIdAlice;
    uint256 public tokenIdBob;

    function setUp() public {
        _deployAll();
        _setupDAO();

        // Alice stakes 5M AWP for 30 days
        vm.startPrank(alice);
        awp.approve(address(veAwp), 5_000_000e18);
        tokenIdAlice = veAwp.deposit(5_000_000e18, 30 days);
        vm.stopPrank();

        // Bob stakes 3M AWP for 30 days
        vm.startPrank(bob);
        awp.approve(address(veAwp), 3_000_000e18);
        tokenIdBob = veAwp.deposit(3_000_000e18, 30 days);
        vm.stopPrank();
    }

    function _setupDAO() internal {
        // Deploy Treasury with DAO as proposer
        address[] memory proposers = new address[](1);
        address[] memory executors = new address[](1);
        executors[0] = address(0); // anyone can execute

        // We need to predict DAO proxy address for Treasury proposer
        uint256 nonce = vm.getNonce(address(this));
        // nonce+0: TimelockControllerUpgradeable impl
        // nonce+1: TimelockController proxy
        // nonce+2: AWPDAO impl
        // nonce+3: AWPDAO proxy
        address daoProxyAddr = vm.computeCreateAddress(address(this), nonce + 3);
        proposers[0] = daoProxyAddr;

        // Deploy timelock (upgradeable version via proxy)
        TimelockControllerUpgradeable timelockImpl = new TimelockControllerUpgradeable(); // nonce+0
        TimelockControllerUpgradeable timelock = TimelockControllerUpgradeable(payable(address(new ERC1967Proxy( // nonce+1
            address(timelockImpl),
            abi.encodeCall(TimelockControllerUpgradeable.initialize, (
                0, // minDelay=0 for tests
                proposers,
                executors,
                address(this) // admin
            ))
        ))));

        // Deploy DAO
        AWPDAO daoImpl = new AWPDAO(address(veAwp)); // nonce+2
        dao = AWPDAO(payable(address(new ERC1967Proxy( // nonce+3
            address(daoImpl),
            abi.encodeCall(AWPDAO.initialize, (
                timelock,
                1,       // votingDelay: 1 block
                50400,   // votingPeriod: ~7 days in blocks
                1,       // lateQuorumExtension: 1 block
                4,       // quorumPercent: 4%
                guardian
            ))
        ))));

        require(address(dao) == daoProxyAddr, "DAO proxy address mismatch");

        // Grant DAO the PROPOSER_ROLE on timelock
        timelock.grantRole(timelock.PROPOSER_ROLE(), address(dao));
        // Grant DAO the CANCELLER_ROLE on timelock
        timelock.grantRole(timelock.CANCELLER_ROLE(), address(dao));
    }

    // ═══════════════════════════════════════════════
    //  Initialization
    // ═══════════════════════════════════════════════

    function test_initialization() public view {
        assertEq(dao.name(), "AWPDAO");
        assertEq(dao.quorumPercent(), 4);
        assertEq(dao.guardian(), guardian);
        assertEq(address(dao.veAWP()), address(veAwp));
        assertEq(dao.votingDelay(), 1);
        assertEq(dao.votingPeriod(), 50400);
    }

    function test_COUNTING_MODE() public view {
        assertEq(dao.COUNTING_MODE(), "support=bravo&quorum=for,abstain&params=tokenIds");
    }

    function test_CLOCK_MODE() public view {
        assertEq(dao.CLOCK_MODE(), "mode=blocknumber&from=default");
    }

    // ═══════════════════════════════════════════════
    //  propose() blocked
    // ═══════════════════════════════════════════════

    function test_propose_blocked() public {
        address[] memory targets = new address[](1);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);

        vm.prank(alice);
        vm.expectRevert(AWPDAO.UseProposeWithTokens.selector);
        dao.propose(targets, values, calldatas, "test");
    }

    // ═══════════════════════════════════════════════
    //  castVote / castVoteWithReason blocked
    // ═══════════════════════════════════════════════

    function test_castVote_blocked() public {
        vm.expectRevert(AWPDAO.UseCastVoteWithParams.selector);
        dao.castVote(0, 1);
    }

    function test_castVoteWithReason_blocked() public {
        vm.expectRevert(AWPDAO.UseCastVoteWithParams.selector);
        dao.castVoteWithReason(0, 1, "reason");
    }

    // ═══════════════════════════════════════════════
    //  proposeWithTokens
    // ═══════════════════════════════════════════════

    function test_proposeWithTokens() public {
        // Need to advance one block so createdAt < proposalCreatedAt
        vm.roll(block.number + 2);
        vm.warp(block.timestamp + 2);

        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(AWPDAO.setQuorumPercent, (5));

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = tokenIdAlice;

        vm.prank(alice);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, "set quorum to 5%", tokenIds);

        assertTrue(proposalId != 0);
        assertTrue(dao.proposalCreatedAt(proposalId) > 0);
        assertTrue(dao.proposalTotalVotingPower(proposalId) > 0);
    }

    function test_proposeWithTokens_insufficientPower_reverts() public {
        vm.roll(block.number + 2);
        vm.warp(block.timestamp + 2);

        // Create a tiny stake that won't meet proposalThreshold
        address charlie = makeAddr("charlie");
        awp.transfer(charlie, 100e18);
        vm.startPrank(charlie);
        awp.approve(address(veAwp), 100e18);
        uint256 tinyTokenId = veAwp.deposit(100e18, 30 days);
        vm.stopPrank();

        vm.roll(block.number + 2);
        vm.warp(block.timestamp + 2);

        address[] memory targets = new address[](1);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = tinyTokenId;

        vm.prank(charlie);
        vm.expectRevert();
        dao.proposeWithTokens(targets, values, calldatas, "test", tokenIds);
    }

    // ═══════════════════════════════════════════════
    //  Voting
    // ═══════════════════════════════════════════════

    function test_vote_for() public {
        uint256 proposalId = _createProposal();

        // Advance past voting delay
        vm.roll(block.number + 2);

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = tokenIdBob;

        vm.prank(bob);
        dao.castVoteWithReasonAndParams(proposalId, 1, "I support", abi.encode(tokenIds));

        assertTrue(dao.hasVotedWithToken(proposalId, tokenIdBob));
        (, uint256 forVotes,) = dao.proposalVotes(proposalId);
        assertTrue(forVotes > 0);
    }

    function test_vote_against() public {
        uint256 proposalId = _createProposal();
        vm.roll(block.number + 2);

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = tokenIdBob;

        vm.prank(bob);
        dao.castVoteWithReasonAndParams(proposalId, 0, "", abi.encode(tokenIds));

        (uint256 againstVotes,,) = dao.proposalVotes(proposalId);
        assertTrue(againstVotes > 0);
    }

    function test_vote_abstain() public {
        uint256 proposalId = _createProposal();
        vm.roll(block.number + 2);

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = tokenIdBob;

        vm.prank(bob);
        dao.castVoteWithReasonAndParams(proposalId, 2, "", abi.encode(tokenIds));

        (,, uint256 abstainVotes) = dao.proposalVotes(proposalId);
        assertTrue(abstainVotes > 0);
    }

    function test_vote_notOwner_reverts() public {
        uint256 proposalId = _createProposal();
        vm.roll(block.number + 2);

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = tokenIdAlice; // alice's token

        vm.prank(bob);
        vm.expectRevert(AWPDAO.NotTokenOwner.selector);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));
    }

    function test_vote_doubleVote_reverts() public {
        uint256 proposalId = _createProposal();
        vm.roll(block.number + 2);

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = tokenIdBob;

        vm.prank(bob);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));

        vm.prank(bob);
        vm.expectRevert(AWPDAO.TokenAlreadyVoted.selector);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));
    }

    function test_vote_noTokens_reverts() public {
        uint256 proposalId = _createProposal();
        vm.roll(block.number + 2);

        uint256[] memory tokenIds = new uint256[](0);

        vm.prank(bob);
        vm.expectRevert(AWPDAO.NoTokens.selector);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));
    }

    function test_vote_mintedAfterProposal_reverts() public {
        uint256 proposalId = _createProposal();
        vm.roll(block.number + 2);

        // Charlie deposits after proposal
        address charlie = makeAddr("charlie");
        awp.transfer(charlie, 1_000_000e18);
        vm.startPrank(charlie);
        awp.approve(address(veAwp), 1_000_000e18);
        uint256 charlieToken = veAwp.deposit(1_000_000e18, 30 days);
        vm.stopPrank();

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = charlieToken;

        vm.prank(charlie);
        vm.expectRevert(AWPDAO.MintedAfterProposal.selector);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));
    }

    function test_vote_expiredLock_reverts() public {
        uint256 proposalId = _createProposal();
        vm.roll(block.number + 2);

        // Warp past bob's lock
        vm.warp(block.timestamp + 31 days);

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = tokenIdBob;

        vm.prank(bob);
        vm.expectRevert(AWPDAO.LockExpired.selector);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));
    }

    // ═══════════════════════════════════════════════
    //  hasVoted always false (per-token tracking)
    // ═══════════════════════════════════════════════

    function test_hasVoted_alwaysFalse() public view {
        assertFalse(dao.hasVoted(0, alice));
    }

    // ═══════════════════════════════════════════════
    //  Signal Proposal
    // ═══════════════════════════════════════════════

    function test_signalPropose() public {
        vm.roll(block.number + 2);
        vm.warp(block.timestamp + 2);

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = tokenIdAlice;

        vm.prank(alice);
        uint256 proposalId = dao.signalPropose("This is a signal proposal", tokenIds);

        assertTrue(proposalId != 0);
        assertTrue(dao.isSignalProposal(proposalId));
        assertFalse(dao.proposalNeedsQueuing(proposalId));
    }

    // ═══════════════════════════════════════════════
    //  Guardian
    // ═══════════════════════════════════════════════

    function test_setGuardian() public {
        vm.prank(guardian);
        dao.setGuardian(alice);
        assertEq(dao.guardian(), alice);
    }

    function test_setGuardian_notGuardianOrExecutor_reverts() public {
        vm.prank(alice);
        vm.expectRevert(AWPDAO.NotGuardianOrExecutor.selector);
        dao.setGuardian(alice);
    }

    function test_guardianCancel() public {
        uint256 proposalId = _createProposal();

        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(AWPDAO.setQuorumPercent, (5));

        vm.prank(guardian);
        dao.guardianCancel(targets, values, calldatas, keccak256("set quorum to 5%"));
    }

    function test_guardianCancel_notGuardian_reverts() public {
        address[] memory targets = new address[](1);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);

        vm.prank(alice);
        vm.expectRevert(AWPDAO.NotGuardian.selector);
        dao.guardianCancel(targets, values, calldatas, bytes32(0));
    }

    // ═══════════════════════════════════════════════
    //  Quorum
    // ═══════════════════════════════════════════════

    function test_setQuorumPercent() public {
        vm.prank(guardian);
        dao.setQuorumPercent(10);
        assertEq(dao.quorumPercent(), 10);
    }

    function test_setQuorumPercent_zero_reverts() public {
        vm.prank(guardian);
        vm.expectRevert(AWPDAO.InvalidQuorumPercent.selector);
        dao.setQuorumPercent(0);
    }

    function test_setQuorumPercent_over100_reverts() public {
        vm.prank(guardian);
        vm.expectRevert(AWPDAO.InvalidQuorumPercent.selector);
        dao.setQuorumPercent(101);
    }

    function test_quorum() public view {
        uint256 totalVP = veAwp.totalVotingPower();
        uint256 expectedQuorum = totalVP * 4 / 100;
        // quorum(0) uses live totalVotingPower since no snapshot for proposalId=0
        assertEq(dao.quorum(0), expectedQuorum);
    }

    // ═══════════════════════════════════════════════
    //  Helper
    // ═══════════════════════════════════════════════

    function _createProposal() internal returns (uint256) {
        vm.roll(block.number + 2);
        vm.warp(block.timestamp + 2);

        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(AWPDAO.setQuorumPercent, (5));

        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = tokenIdAlice;

        vm.prank(alice);
        return dao.proposeWithTokens(targets, values, calldatas, "set quorum to 5%", tokenIds);
    }
}
