// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {IGovernor} from "@openzeppelin/contracts/governance/IGovernor.sol";
import {TimelockController} from "@openzeppelin/contracts/governance/TimelockController.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {StakingVault} from "../src/core/StakingVault.sol";

contract AWPDAOTest is Test {
    AWPDAO public dao;
    Treasury public treasury;
    AWPToken public awpToken;
    StakeNFT public stakeNFT;
    StakingVault public vault;

    address public deployer = makeAddr("deployer");
    address public voter = makeAddr("voter");

    uint256 public constant MIN_DELAY = 1 days;
    uint48 public constant VOTING_DELAY = 1; // 1 block
    uint32 public constant VOTING_PERIOD = 50400; // ~1 week in blocks
    uint256 public constant QUORUM_PERCENTAGE = 4; // 4%

    function setUp() public {
        vm.startPrank(deployer);

        // Deploy AWPToken
        awpToken = new AWPToken("AWP", "AWP", deployer);

        // Deploy StakingVault + StakeNFT (circular dependency)
        // This test contract (address(this)) acts as awpRegistry for access control
        uint64 nonce = vm.getNonce(deployer);
        address predictedVault = vm.computeCreateAddress(deployer, nonce);
        address predictedStakeNFT = vm.computeCreateAddress(deployer, nonce + 1);

        vault = new StakingVault(address(this));
        stakeNFT = new StakeNFT(address(awpToken), address(vault), address(this));
        vm.stopPrank();
        vault.setStakeNFT(address(stakeNFT));
        vm.startPrank(deployer);

        // Deploy Treasury
        address[] memory proposers = new address[](0);
        address[] memory executors = new address[](1);
        executors[0] = address(0); // anyone can execute
        treasury = new Treasury(MIN_DELAY, proposers, executors, deployer);

        // Deploy AWPDAO (no awpRegistry param)
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

        // Deployer renounces Treasury admin
        treasury.renounceRole(treasury.DEFAULT_ADMIN_ROLE(), deployer);

        // Transfer tokens to voter and have them stake via StakeNFT
        uint256 voterAmount = 100_000_000 * 1e18; // 100M AWP
        awpToken.transfer(voter, voterAmount);

        vm.stopPrank();

        // Voter deposits into StakeNFT to get voting power
        vm.startPrank(voter);
        awpToken.approve(address(stakeNFT), voterAmount);
        stakeNFT.deposit(voterAmount, 52 weeks); // 52 weeks lock
        vm.stopPrank();
    }

    /// @notice Verify Treasury roles are set correctly
    function test_treasuryRolesSetup() public view {
        assertTrue(treasury.hasRole(treasury.PROPOSER_ROLE(), address(dao)), "dao should be proposer");
        assertTrue(treasury.hasRole(treasury.CANCELLER_ROLE(), address(dao)), "dao should be canceller");
        assertFalse(treasury.hasRole(treasury.DEFAULT_ADMIN_ROLE(), deployer), "deployer should not be admin");
    }

    /// @notice Verify DAO basic configuration
    function test_daoConfiguration() public view {
        assertEq(dao.votingDelay(), VOTING_DELAY, "voting delay mismatch");
        assertEq(dao.votingPeriod(), VOTING_PERIOD, "voting period mismatch");
    }

    /// @notice Full governance flow: propose -> vote (with tokenIds) -> queue -> execute
    function test_fullGovernanceFlow() public {
        // Advance some time so position's createdAt < proposalCreatedAt
        vm.warp(block.timestamp + 1 days);

        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        // Prepare proposal: send ETH from Treasury
        address payable target = payable(makeAddr("target"));
        vm.etch(target, ""); // Ensure target is EOA — needed for fork tests
        uint256 sendAmount = 1 ether;
        vm.deal(address(treasury), sendAmount);

        address[] memory targets = new address[](1);
        targets[0] = target;
        uint256[] memory values = new uint256[](1);
        values[0] = sendAmount;
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = "";
        string memory description = "Transfer 1 ETH to target";

        // Create proposal (voter has sufficient voting power via StakeNFT)
        uint256[] memory proposerTokenIds = new uint256[](1);
        proposerTokenIds[0] = 1; // First token minted to voter
        vm.prank(voter);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, description, proposerTokenIds);

        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Pending), "should be Pending");

        // Advance to voting period
        vm.roll(block.number + VOTING_DELAY + 1);
        vm.warp(block.timestamp + (VOTING_DELAY + 1) * 12);

        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Active), "should be Active");

        // Vote with tokenIds (NFT-based voting)
        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = 1; // First token minted to voter
        bytes memory params = abi.encode(tokenIds);

        vm.prank(voter);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", params); // 1 = For

        // Advance past voting period
        vm.roll(block.number + VOTING_PERIOD + 1);
        vm.warp(block.timestamp + (VOTING_PERIOD + 1) * 12);

        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Succeeded), "should be Succeeded");

        // Queue
        bytes32 descriptionHash = keccak256(bytes(description));
        dao.queue(targets, values, calldatas, descriptionHash);

        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Queued), "should be Queued");

        // Advance past timelock delay
        vm.warp(block.timestamp + MIN_DELAY + 1);

        // Execute
        uint256 targetBalBefore = target.balance;
        dao.execute(targets, values, calldatas, descriptionHash);

        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Executed), "should be Executed");
        assertEq(target.balance, targetBalBefore + sendAmount, "target should receive ETH");
    }

    /// @notice Proposal defeated without quorum (insufficient voting power)
    function test_proposalDefeatedWithoutQuorum() public {
        // Create a small voter with insufficient voting power for quorum
        address smallVoter = makeAddr("smallVoter");
        vm.prank(deployer);
        awpToken.transfer(smallVoter, 1_000_000 * 1e18);

        vm.startPrank(smallVoter);
        awpToken.approve(address(stakeNFT), 1_000_000 * 1e18);
        stakeNFT.deposit(1_000_000 * 1e18, 52 weeks);
        vm.stopPrank();

        // Advance some time so position's createdAt < proposalCreatedAt
        vm.warp(block.timestamp + 1 days);

        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        // Propose
        address[] memory targets = new address[](1);
        targets[0] = address(treasury);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = "";

        uint256[] memory smallTokenIds = new uint256[](1);
        smallTokenIds[0] = 2; // Second token minted (after voter's token 1)
        vm.prank(smallVoter);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, "small proposal", smallTokenIds);

        // Advance to voting period
        vm.roll(block.number + VOTING_DELAY + 1);
        vm.warp(block.timestamp + (VOTING_DELAY + 1) * 12);

        // Small voter votes with their token
        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = 2; // Second token minted (after voter's token 1)
        bytes memory params = abi.encode(tokenIds);
        vm.prank(smallVoter);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", params);

        // End voting period
        vm.roll(block.number + VOTING_PERIOD + 1);
        vm.warp(block.timestamp + (VOTING_PERIOD + 1) * 12);

        // Should be defeated due to insufficient quorum
        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Defeated), "should be Defeated");
    }

    /// @notice Signal proposal: vote-only, no on-chain execution, skips Timelock
    function test_signalProposal() public {
        // Advance some time so position's createdAt < proposalCreatedAt
        vm.warp(block.timestamp + 1 days);
        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        // Create signal proposal
        uint256[] memory proposerTokenIds = new uint256[](1);
        proposerTokenIds[0] = 1;
        vm.prank(voter);
        uint256 proposalId = dao.signalPropose("Should we adopt proposal XYZ?", proposerTokenIds);

        assertTrue(dao.isSignalProposal(proposalId), "should be signal proposal");
        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Pending), "should be Pending");

        // Advance to voting period
        vm.roll(block.number + VOTING_DELAY + 1);
        vm.warp(block.timestamp + (VOTING_DELAY + 1) * 12);
        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Active), "should be Active");

        // Vote
        uint256[] memory tokenIds = new uint256[](1);
        tokenIds[0] = 1;
        vm.prank(voter);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(tokenIds));

        // Advance past voting period
        vm.roll(block.number + VOTING_PERIOD + 1);
        vm.warp(block.timestamp + (VOTING_PERIOD + 1) * 12);
        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Succeeded), "should be Succeeded");

        // Execute immediately — no queue/timelock needed
        // Must match the targets/values/calldatas used in signalPropose (single no-op entry)
        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        bytes32 descriptionHash = keccak256(bytes("Should we adopt proposal XYZ?"));

        dao.execute(targets, values, calldatas, descriptionHash);
        assertEq(uint256(dao.state(proposalId)), uint256(IGovernor.ProposalState.Executed), "should be Executed");

        // Verify vote tallies
        (uint256 against, uint256 forVotes, uint256 abstain) = dao.proposalVotes(proposalId);
        assertGt(forVotes, 0, "forVotes should be > 0");
        assertEq(against, 0, "againstVotes should be 0");
        assertEq(abstain, 0, "abstainVotes should be 0");
    }
}
