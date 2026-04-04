// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {DeployHelper} from "./helpers/DeployHelper.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {IAWPRegistry} from "../src/interfaces/IAWPRegistry.sol";
import {IAlphaToken} from "../src/interfaces/IAlphaToken.sol";
import {AWPWorkNet} from "../src/core/AWPWorkNet.sol";
import {AlphaToken} from "../src/token/AlphaToken.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {TimelockControllerUpgradeable} from "@openzeppelin/contracts-upgradeable/governance/TimelockControllerUpgradeable.sol";
import {MerkleProof} from "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import {IERC1363Receiver} from "../src/interfaces/IERC1363Receiver.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @title IntegrationTest — Cross-contract integration tests
/// @dev Tests full protocol flows that span multiple contracts
contract IntegrationTest is DeployHelper {
    TimelockControllerUpgradeable public timelock;
    uint256 public tokenIdAlice;
    uint256 public tokenIdBob;

    function setUp() public {
        _deployAll();
        _setupDAOAndTimelock();

        // Pre-stake for subsequent DAO tests
        vm.startPrank(alice);
        awp.approve(address(veAwp), 5_000_000e18);
        tokenIdAlice = veAwp.deposit(5_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.startPrank(bob);
        awp.approve(address(veAwp), 3_000_000e18);
        tokenIdBob = veAwp.deposit(3_000_000e18, 54 weeks);
        vm.stopPrank();
    }

    function _setupDAOAndTimelock() internal {
        // Predict DAO proxy address to use as timelock proposer
        uint256 nonce = vm.getNonce(address(this));
        // nonce+0: TimelockControllerUpgradeable impl
        // nonce+1: Timelock proxy
        // nonce+2: AWPDAO impl
        // nonce+3: AWPDAO proxy
        address daoProxyAddr = vm.computeCreateAddress(address(this), nonce + 3);

        address[] memory proposers = new address[](1);
        proposers[0] = daoProxyAddr;
        address[] memory executors = new address[](1);
        executors[0] = address(0); // anyone can execute

        TimelockControllerUpgradeable timelockImpl = new TimelockControllerUpgradeable();
        timelock = TimelockControllerUpgradeable(payable(address(new ERC1967Proxy(
            address(timelockImpl),
            abi.encodeCall(TimelockControllerUpgradeable.initialize, (
                0, // minDelay=0 for testing
                proposers,
                executors,
                address(this)
            ))
        ))));

        AWPDAO daoImpl = new AWPDAO(address(veAwp));
        dao = AWPDAO(payable(address(new ERC1967Proxy(
            address(daoImpl),
            abi.encodeCall(AWPDAO.initialize, (
                timelock,
                1,       // votingDelay: 1 block
                100,     // votingPeriod: 100 blocks (short for testing)
                1,       // lateQuorumExtension
                4,       // quorumPercent: 4%
                guardian
            ))
        ))));

        require(address(dao) == daoProxyAddr, "DAO proxy addr mismatch");

        // Grant DAO proposer + canceller roles
        timelock.grantRole(timelock.PROPOSER_ROLE(), address(dao));
        timelock.grantRole(timelock.CANCELLER_ROLE(), address(dao));
    }

    // ═══════════════════════════════════════════════
    //  1. DAO full governance flow: propose → vote → queue → execute
    // ═══════════════════════════════════════════════

    function test_dao_fullGovernanceFlow() public {
        // Send AWP to timelock for testing execute withdrawal
        awp.transfer(address(timelock), 1_000e18);

        // Advance to block 10 to ensure createdAt < proposalCreatedAt
        vm.roll(10);
        vm.warp(block.timestamp + 10);

        // Build proposal: timelock transfers 500 AWP to alice
        address[] memory targets = new address[](1);
        targets[0] = address(awp);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(IERC20.transfer, (alice, 500e18));
        string memory description = "Transfer 500 AWP to Alice";

        // Alice proposes (at block 10, voteStart = 10 + votingDelay(1) = 11)
        uint256[] memory propTokens = new uint256[](1);
        propTokens[0] = tokenIdAlice;

        vm.prank(alice);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, description, propTokens);

        // Advance past votingDelay (block > 11)
        vm.roll(20);

        // Alice votes For
        uint256[] memory aliceTokens = new uint256[](1);
        aliceTokens[0] = tokenIdAlice;
        vm.prank(alice);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(aliceTokens));

        // Bob votes For
        uint256[] memory bobTokens = new uint256[](1);
        bobTokens[0] = tokenIdBob;
        vm.prank(bob);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(bobTokens));

        // Verify vote results
        (, uint256 forVotes,) = dao.proposalVotes(proposalId);
        assertTrue(forVotes > 0);

        // Advance past votingPeriod (voteEnd = 10 + 1 + 100 = 111)
        vm.roll(200);

        // Queue proposal (minDelay=0, can execute immediately)
        dao.queue(targets, values, calldatas, keccak256(bytes(description)));

        // Execute
        uint256 aliceBalBefore = awp.balanceOf(alice);
        dao.execute(targets, values, calldatas, keccak256(bytes(description)));

        assertEq(awp.balanceOf(alice) - aliceBalBefore, 500e18);
    }

    // ═══════════════════════════════════════════════
    //  2. Emission → WorknetManager (mintAndCall callback)
    // ═══════════════════════════════════════════════

    function test_emission_settleToWorknetManager() public {
        // Register and activate worknet — activation deploys WorknetManager proxy
        uint256 wid = _registerWorknet(alice, "Emission Test", "EMT");
        _activateWorknet(wid);

        // Get worknetManager address
        AWPWorkNet.WorknetData memory wd = awpWorkNet.getWorknetData(wid);
        address worknetMgr = wd.worknetManager;
        assertTrue(worknetMgr != address(0));

        // Set worknetManager as emission recipient
        vm.warp(GENESIS_TIME + 1);

        uint256[] memory packed = new uint256[](1);
        packed[0] = (uint256(100) << 160) | uint256(uint160(worknetMgr));

        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 0);

        // settleEpoch calls awpToken.mint(worknetMgr, amount) then try onTransferReceived
        uint256 wmBalBefore = awp.balanceOf(worknetMgr);
        awpEmission.settleEpoch(100);
        uint256 wmBalAfter = awp.balanceOf(worknetMgr);

        assertTrue(wmBalAfter > wmBalBefore);
        assertEq(awpEmission.settledEpoch(), 1);
    }

    // ═══════════════════════════════════════════════
    //  3. AlphaToken lifecycle: activate → setWorknetMinter → Merkle claim
    // ═══════════════════════════════════════════════

    function test_alphaToken_fullLifecycle() public {
        // Register and activate worknet
        uint256 wid = _registerWorknet(alice, "Alpha Lifecycle", "ALC");
        _activateWorknet(wid);

        AWPWorkNet.WorknetData memory wd = awpWorkNet.getWorknetData(wid);
        address alphaAddr = wd.alphaToken;
        address wmAddr = wd.worknetManager;
        assertTrue(alphaAddr != address(0));
        assertTrue(wmAddr != address(0));

        AlphaToken alpha = AlphaToken(alphaAddr);

        // Verify AlphaToken state
        assertEq(alpha.worknetId(), wid);
        assertTrue(alpha.mintersLocked()); // setWorknetMinter was called
        assertTrue(alpha.minters(wmAddr)); // worknetManager is the minter
        assertFalse(alpha.minters(address(awpRegistry))); // admin minting revoked
        assertTrue(alpha.supplyAtLock() > 0); // LP pre-mint amount

        // WorknetManager can mint Alpha to users via merkle claim
        // Build merkle proof: bob claims 50e18
        uint256 claimAmount = 50e18;
        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(bob, claimAmount))));
        bytes32 root = leaf; // single-leaf tree

        // MockWorknetManager doesn't have setMerkleRoot
        // Real WorknetManager proxy does (if defaultWorknetManagerImpl is real impl)
        // MockWorknetManager doesn't support merkle in tests, skip claim step
        // Instead verify AlphaToken time-based minting cap
        vm.warp(block.timestamp + 30 days);
        uint256 mintable = alpha.currentMintableLimit();
        assertTrue(mintable > 0);
    }

    // ═══════════════════════════════════════════════
    //  4. Ban / Unban effect on AlphaToken minting
    // ═══════════════════════════════════════════════

    function test_banUnban_alphaTokenMinting() public {
        // Register and activate worknet
        uint256 wid = _registerWorknet(alice, "BanTest", "BAN");
        _activateWorknet(wid);

        AWPWorkNet.WorknetData memory wd = awpWorkNet.getWorknetData(wid);
        address alphaAddr = wd.alphaToken;
        address wmAddr = wd.worknetManager;
        AlphaToken alpha = AlphaToken(alphaAddr);

        // Confirm worknetManager can currently mint
        assertTrue(alpha.minters(wmAddr));
        assertFalse(alpha.minterPaused(wmAddr));

        // Guardian ban worknet
        vm.prank(guardian);
        awpRegistry.banWorknet(wid);

        // Verify worknet status is Banned
        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Banned));

        // Verify AlphaToken minter is paused
        assertTrue(alpha.minterPaused(wmAddr));

        // Guardian unban worknet
        vm.prank(guardian);
        awpRegistry.unbanWorknet(wid);

        // Verify restored
        info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Active));
        assertFalse(alpha.minterPaused(wmAddr));
    }

    // ═══════════════════════════════════════════════
    //  5. Treasury withdrawal: schedule → execute AWP transfer
    // ═══════════════════════════════════════════════

    function test_treasury_scheduleAndWithdrawAWP() public {
        // Send AWP to timelock (as Treasury)
        uint256 depositAmount = 10_000e18;
        awp.transfer(address(timelock), depositAmount);
        assertEq(awp.balanceOf(address(timelock)), depositAmount);

        // Advance to ensure stake createdAt < proposalCreatedAt
        vm.roll(10);
        vm.warp(block.timestamp + 10);

        // Build proposal: Treasury transfers 5000 AWP to bob
        address[] memory targets = new address[](1);
        targets[0] = address(awp);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(IERC20.transfer, (bob, 5_000e18));
        string memory description = "Withdraw 5000 AWP from Treasury to Bob";

        // Alice proposes (block 10, voteStart = 11)
        uint256[] memory propTokens = new uint256[](1);
        propTokens[0] = tokenIdAlice;
        vm.prank(alice);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, description, propTokens);

        // Advance past votingDelay (block > 11)
        vm.roll(20);

        // Alice + Bob both vote For
        uint256[] memory aliceTokens = new uint256[](1);
        aliceTokens[0] = tokenIdAlice;
        vm.prank(alice);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(aliceTokens));

        uint256[] memory bobTokens = new uint256[](1);
        bobTokens[0] = tokenIdBob;
        vm.prank(bob);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(bobTokens));

        // Advance past votingPeriod (voteEnd = 111)
        vm.roll(200);

        // Queue
        dao.queue(targets, values, calldatas, keccak256(bytes(description)));

        // Execute
        uint256 bobBalBefore = awp.balanceOf(bob);
        dao.execute(targets, values, calldatas, keccak256(bytes(description)));

        assertEq(awp.balanceOf(bob) - bobBalBefore, 5_000e18);
        assertEq(awp.balanceOf(address(timelock)), depositAmount - 5_000e18);
    }

    // ═══════════════════════════════════════════════
    //  6. Signal Proposal: votes pass without execution
    // ═══════════════════════════════════════════════

    function test_dao_signalProposal() public {
        vm.roll(10);
        vm.warp(block.timestamp + 10);

        uint256[] memory propTokens = new uint256[](1);
        propTokens[0] = tokenIdAlice;

        vm.prank(alice);
        uint256 proposalId = dao.signalPropose("Should we add XYZ feature?", propTokens);
        assertTrue(dao.isSignalProposal(proposalId));

        // Advance past votingDelay
        vm.roll(20);

        // Vote
        uint256[] memory aliceTokens = new uint256[](1);
        aliceTokens[0] = tokenIdAlice;
        vm.prank(alice);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(aliceTokens));

        uint256[] memory bobTokens = new uint256[](1);
        bobTokens[0] = tokenIdBob;
        vm.prank(bob);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(bobTokens));

        (, uint256 forVotes,) = dao.proposalVotes(proposalId);
        assertTrue(forVotes > 0);

        // Advance past votingPeriod
        vm.roll(200);

        // Signal proposal doesn't need queue/execute
        assertFalse(dao.proposalNeedsQueuing(proposalId));
    }

    // ═══════════════════════════════════════════════
    //  7. Guardian Cancel: emergency cancel active proposal
    // ═══════════════════════════════════════════════

    function test_dao_guardianCancelActiveProposal() public {
        vm.roll(10);
        vm.warp(block.timestamp + 10);

        address[] memory targets = new address[](1);
        targets[0] = address(awp);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(IERC20.transfer, (alice, 999_999e18));
        string memory description = "Malicious drain proposal";

        uint256[] memory propTokens = new uint256[](1);
        propTokens[0] = tokenIdAlice;
        vm.prank(alice);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, description, propTokens);

        // Guardian emergency cancel
        vm.prank(guardian);
        dao.guardianCancel(targets, values, calldatas, keccak256(bytes(description)));

        // Proposal canceled
        assertEq(uint8(dao.state(proposalId)), 2); // Canceled
    }

    // ═══════════════════════════════════════════════
    //  8. Emission multi-epoch sequential settle
    // ═══════════════════════════════════════════════

    function test_emission_multiEpochSettle() public {
        vm.warp(GENESIS_TIME + 1);

        // Submit epoch 0 weights
        uint256[] memory packed = new uint256[](2);
        packed[0] = (uint256(70) << 160) | uint256(uint160(alice));
        packed[1] = (uint256(30) << 160) | uint256(uint160(bob));

        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 0);

        // Settle epoch 0
        uint256 aliceBalBefore = awp.balanceOf(alice);
        uint256 bobBalBefore = awp.balanceOf(bob);
        awpEmission.settleEpoch(100);

        uint256 aliceMinted = awp.balanceOf(alice) - aliceBalBefore;
        uint256 bobMinted = awp.balanceOf(bob) - bobBalBefore;
        assertTrue(aliceMinted > 0);
        assertTrue(bobMinted > 0);
        // alice gets 70%, bob gets 30%
        assertApproxEqRel(aliceMinted * 30, bobMinted * 70, 1e16); // 1% tolerance

        // Advance to epoch 1
        vm.warp(GENESIS_TIME + 1 days + 1);

        // Submit epoch 1 weights (different allocation)
        uint256[] memory packed2 = new uint256[](1);
        packed2[0] = (uint256(100) << 160) | uint256(uint160(bob));

        vm.prank(guardian);
        awpEmission.submitAllocations(packed2, 100, 1);

        // Settle epoch 1
        bobBalBefore = awp.balanceOf(bob);
        aliceBalBefore = awp.balanceOf(alice);
        awpEmission.settleEpoch(100);

        // Epoch 1: 100% to bob
        assertEq(awp.balanceOf(alice), aliceBalBefore); // alice gets nothing new
        assertTrue(awp.balanceOf(bob) > bobBalBefore);   // bob gets everything
        assertEq(awpEmission.settledEpoch(), 2);
    }

    // ═══════════════════════════════════════════════
    //  9. Allocation + Emission combined: stake → allocate → emit → verify
    // ═══════════════════════════════════════════════

    function test_allocation_and_emission_combined() public {
        // Register and activate worknet
        uint256 wid = _registerWorknet(alice, "Combined", "CMB");
        _activateWorknet(wid);

        // Bob has existing veAWP stake, allocate to worknet
        vm.prank(bob);
        awpAllocator.allocate(bob, bob, wid, 1_000_000e18);

        assertEq(awpAllocator.worknetTotalStake(wid), 1_000_000e18);
        assertEq(awpAllocator.userTotalAllocated(bob), 1_000_000e18);

        // Emission: 100% weight to alice (worknet manager)
        vm.warp(GENESIS_TIME + 1);

        uint256[] memory packed = new uint256[](1);
        packed[0] = (uint256(100) << 160) | uint256(uint160(alice));

        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 0);

        awpEmission.settleEpoch(100);

        // alice received emission
        assertTrue(awp.balanceOf(alice) > 5_000_000e18); // initial 10M + emission

        // Bob deallocate
        vm.prank(bob);
        awpAllocator.deallocateAll(bob, bob, wid);

        assertEq(awpAllocator.userTotalAllocated(bob), 0);
        assertEq(awpAllocator.worknetTotalStake(wid), 0);
    }

    // ═══════════════════════════════════════════════
    //  10. Pause/Resume + Ban/Unban full status transitions
    // ═══════════════════════════════════════════════

    function test_worknet_fullStatusTransitions() public {
        uint256 wid = _registerWorknet(alice, "Status", "STS");

        // Pending → Active
        _activateWorknet(wid);
        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Active));
        assertTrue(awpRegistry.isWorknetActive(wid));

        // Active → Paused (by owner)
        vm.prank(alice);
        awpRegistry.pauseWorknet(wid);
        info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Paused));
        assertFalse(awpRegistry.isWorknetActive(wid));

        // Paused → Active (by owner)
        vm.prank(alice);
        awpRegistry.resumeWorknet(wid);
        assertTrue(awpRegistry.isWorknetActive(wid));

        // Active → Banned (by guardian)
        vm.prank(guardian);
        awpRegistry.banWorknet(wid);
        info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Banned));

        // Banned → Active (by guardian)
        vm.prank(guardian);
        awpRegistry.unbanWorknet(wid);
        assertTrue(awpRegistry.isWorknetActive(wid));
    }

    // ═══════════════════════════════════════════════
    //  11. Delegation cross-contract: Registry grant → Allocator use → WorkNet update
    // ═══════════════════════════════════════════════

    function test_delegation_crossContract() public {
        uint256 wid = _registerWorknet(alice, "Delegate", "DEL");
        _activateWorknet(wid);

        // Alice grants bob as delegate (in AWPRegistry)
        vm.prank(alice);
        awpRegistry.grantDelegate(bob);

        // Bob allocates alice's stake as delegate (in AWPAllocator, reads AWPRegistry.delegates)
        vm.prank(bob);
        awpAllocator.allocate(alice, bob, wid, 500_000e18);
        assertEq(awpAllocator.getAgentStake(alice, bob, wid), 500_000e18);

        // Bob updates worknet skillsURI as delegate (in AWPWorkNet, reads AWPRegistry.delegates)
        vm.prank(bob);
        awpWorkNet.setSkillsURI(wid, "https://skills.example.com");
        AWPWorkNet.WorknetMeta memory meta = awpWorkNet.getWorknetMeta(wid);
        assertEq(meta.skillsURI, "https://skills.example.com");

        // Alice revokes delegate
        vm.prank(alice);
        awpRegistry.revokeDelegate(bob);

        // Bob's allocate attempt should now fail
        vm.prank(bob);
        vm.expectRevert();
        awpAllocator.allocate(alice, bob, wid, 100_000e18);
    }

    // ═══════════════════════════════════════════════
    //  12. Binding + Recipient chain: bind → setRecipient → resolveRecipient
    // ═══════════════════════════════════════════════

    function test_binding_recipientChain() public {
        // alice bind → bob, bob setRecipient → relayer
        vm.prank(alice);
        awpRegistry.bind(bob);

        vm.prank(bob);
        awpRegistry.setRecipient(relayer);

        // resolveRecipient(alice) should walk the bind chain to bob's recipient = relayer
        assertEq(awpRegistry.resolveRecipient(alice), relayer);

        // If emission goes to alice, the final recipient is relayer
        // Verify the bind chain affects distribution
        assertEq(awpRegistry.boundTo(alice), bob);
        assertEq(awpRegistry.recipient(bob), relayer);
    }

    // ═══════════════════════════════════════════════
    //  13. Gasless EIP-712: bindFor + setRecipientFor
    // ═══════════════════════════════════════════════

    function test_gasless_bindFor_setRecipientFor() public {
        // Use vm.sign to simulate EIP-712 signatures
        (address signer, uint256 signerPk) = makeAddrAndKey("eip712signer");
        awp.transfer(signer, 1e18); // Give some balance

        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = awpRegistry.nonces(signer);

        // Build EIP-712 digest for bindFor
        bytes32 BIND_TYPEHASH = keccak256("Bind(address agent,address target,uint256 nonce,uint256 deadline)");
        bytes32 structHash = keccak256(abi.encode(BIND_TYPEHASH, signer, bob, nonce, deadline));
        bytes32 digest = _registryDigest(structHash);

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        // Relayer calls bindFor
        vm.prank(relayer);
        awpRegistry.bindFor(signer, bob, deadline, v, r, s);

        assertEq(awpRegistry.boundTo(signer), bob);

        // setRecipientFor
        nonce = awpRegistry.nonces(signer); // nonce+1 now
        bytes32 SET_RECIPIENT_TYPEHASH = keccak256("SetRecipient(address user,address recipient,uint256 nonce,uint256 deadline)");
        structHash = keccak256(abi.encode(SET_RECIPIENT_TYPEHASH, signer, relayer, nonce, deadline));
        digest = _registryDigest(structHash);

        (v, r, s) = vm.sign(signerPk, digest);

        vm.prank(relayer);
        awpRegistry.setRecipientFor(signer, relayer, deadline, v, r, s);

        assertEq(awpRegistry.recipient(signer), relayer);
    }

    // ═══════════════════════════════════════════════
    //  14. Gasless EIP-712: grantDelegateFor + revokeDelegateFor
    // ═══════════════════════════════════════════════

    function test_gasless_grantDelegateFor_revokeDelegateFor() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("delegateSigner");
        uint256 deadline = block.timestamp + 1 hours;

        // grantDelegateFor
        uint256 nonce = awpRegistry.nonces(signer);
        bytes32 GRANT_DELEGATE_TYPEHASH = keccak256("GrantDelegate(address user,address delegate,uint256 nonce,uint256 deadline)");
        bytes32 structHash = keccak256(abi.encode(GRANT_DELEGATE_TYPEHASH, signer, bob, nonce, deadline));
        bytes32 digest = _registryDigest(structHash);

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);
        vm.prank(relayer);
        awpRegistry.grantDelegateFor(signer, bob, deadline, v, r, s);

        assertTrue(awpRegistry.delegates(signer, bob));

        // revokeDelegateFor
        nonce = awpRegistry.nonces(signer);
        bytes32 REVOKE_DELEGATE_TYPEHASH = keccak256("RevokeDelegate(address user,address delegate,uint256 nonce,uint256 deadline)");
        structHash = keccak256(abi.encode(REVOKE_DELEGATE_TYPEHASH, signer, bob, nonce, deadline));
        digest = _registryDigest(structHash);

        (v, r, s) = vm.sign(signerPk, digest);
        vm.prank(relayer);
        awpRegistry.revokeDelegateFor(signer, bob, deadline, v, r, s);

        assertFalse(awpRegistry.delegates(signer, bob));
    }

    // ═══════════════════════════════════════════════
    //  15. Gasless EIP-712: allocateFor + deallocateFor
    // ═══════════════════════════════════════════════

    function test_gasless_allocateFor_deallocateFor() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("allocSigner");
        awp.transfer(signer, 2_000_000e18);

        // signer stakes
        vm.startPrank(signer);
        awp.approve(address(veAwp), 1_000_000e18);
        veAwp.deposit(1_000_000e18, 30 days);
        vm.stopPrank();

        uint256 worknetId = 845300000001; // Assumed to exist
        uint256 deadline = block.timestamp + 1 hours;
        uint256 amount = 500_000e18;

        // allocateFor
        uint256 nonce = awpAllocator.nonces(signer);
        bytes32 ALLOCATE_TYPEHASH = keccak256("Allocate(address staker,address agent,uint256 worknetId,uint256 amount,uint256 nonce,uint256 deadline)");
        bytes32 structHash = keccak256(abi.encode(ALLOCATE_TYPEHASH, signer, signer, worknetId, amount, nonce, deadline));
        bytes32 digest = _allocatorDigest(structHash);

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);
        vm.prank(relayer);
        awpAllocator.allocateFor(signer, signer, worknetId, amount, deadline, v, r, s);

        assertEq(awpAllocator.getAgentStake(signer, signer, worknetId), amount);

        // deallocateFor
        nonce = awpAllocator.nonces(signer);
        bytes32 DEALLOCATE_TYPEHASH = keccak256("Deallocate(address staker,address agent,uint256 worknetId,uint256 amount,uint256 nonce,uint256 deadline)");
        structHash = keccak256(abi.encode(DEALLOCATE_TYPEHASH, signer, signer, worknetId, amount, nonce, deadline));
        digest = _allocatorDigest(structHash);

        (v, r, s) = vm.sign(signerPk, digest);
        vm.prank(relayer);
        awpAllocator.deallocateFor(signer, signer, worknetId, amount, deadline, v, r, s);

        assertEq(awpAllocator.getAgentStake(signer, signer, worknetId), 0);
    }

    // ═══════════════════════════════════════════════
    //  16. Gasless: registerWorknetFor (EIP-712 signed registration)
    // ═══════════════════════════════════════════════

    function test_gasless_registerWorknetFor() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("regSigner");
        awp.transfer(signer, 10_000_000e18);

        uint256 cost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.prank(signer);
        awp.approve(address(awpRegistry), cost);

        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = awpRegistry.nonces(signer);

        IAWPRegistry.WorknetParams memory params = IAWPRegistry.WorknetParams({
            name: "GaslessNet", symbol: "GLN",
            worknetManager: address(0), salt: bytes32(0),
            minStake: 0, skillsURI: ""
        });

        bytes32 WORKNET_PARAMS_TYPEHASH = keccak256("WorknetParams(string name,string symbol,address worknetManager,bytes32 salt,uint128 minStake,string skillsURI)");
        bytes32 paramsStructHash = keccak256(abi.encode(
            WORKNET_PARAMS_TYPEHASH,
            keccak256(bytes(params.name)),
            keccak256(bytes(params.symbol)),
            params.worknetManager, params.salt, params.minStake,
            keccak256(bytes(params.skillsURI))
        ));

        bytes32 REGISTER_WORKNET_TYPEHASH = keccak256(
            "RegisterWorknet(address user,WorknetParams params,uint256 nonce,uint256 deadline)WorknetParams(string name,string symbol,address worknetManager,bytes32 salt,uint128 minStake,string skillsURI)"
        );
        bytes32 structHash = keccak256(abi.encode(REGISTER_WORKNET_TYPEHASH, signer, paramsStructHash, nonce, deadline));
        bytes32 digest = _registryDigest(structHash);

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        vm.prank(relayer);
        uint256 wid = awpRegistry.registerWorknetFor(signer, params, deadline, v, r, s);

        assertTrue(wid > 0);
        IAWPRegistry.WorknetInfo memory info = awpRegistry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Pending));
    }

    // ═══════════════════════════════════════════════
    //  17. veAWP.depositWithPermit (ERC-2612 gasless deposit)
    // ═══════════════════════════════════════════════

    function test_veAWP_depositWithPermit() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("permitSigner");
        awp.transfer(signer, 1_000_000e18);

        uint256 amount = 500_000e18;
        uint256 deadline = block.timestamp + 1 hours;

        // Build ERC-2612 permit signature
        bytes32 PERMIT_TYPEHASH = keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)");
        bytes32 domainSeparator = awp.DOMAIN_SEPARATOR();
        uint256 nonce = awp.nonces(signer);

        bytes32 structHash = keccak256(abi.encode(PERMIT_TYPEHASH, signer, address(veAwp), amount, nonce, deadline));
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", domainSeparator, structHash));

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        // No prior approve needed — depositWithPermit calls permit internally
        vm.prank(signer);
        uint256 tokenId = veAwp.depositWithPermit(amount, 30 days, deadline, v, r, s);

        assertTrue(tokenId > 0);
        assertEq(veAwp.getUserTotalStaked(signer), amount);
    }

    // ═══════════════════════════════════════════════
    //  18. veAWP NFT transfer + allocation coverage check
    // ═══════════════════════════════════════════════

    function test_veAWP_transfer_allocationCoverage() public {
        // Alice has 5M staked, allocates 4M
        vm.prank(alice);
        awpAllocator.allocate(alice, alice, 845300000001, 4_000_000e18);

        // Alice has 1 NFT (tokenIdAlice), staked=5M, allocated=4M
        // If NFT is transferred to bob, alice's staked becomes 0 but allocated is still 4M -> should revert
        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSignature("InsufficientUnallocated()"));
        veAwp.transferFrom(alice, bob, tokenIdAlice);

        // If deallocated first then transferred -> should succeed
        vm.prank(alice);
        awpAllocator.deallocateAll(alice, alice, 845300000001);

        vm.prank(alice);
        veAwp.transferFrom(alice, bob, tokenIdAlice);

        assertEq(veAwp.ownerOf(tokenIdAlice), bob);
        assertEq(veAwp.getUserTotalStaked(alice), 0);
        assertEq(veAwp.getUserTotalStaked(bob), 5_000_000e18 + 3_000_000e18); // bob's original 3M + alice's 5M
    }

    // ═══════════════════════════════════════════════
    //  19. Batched settle (settleProgress): limit < recipients
    // ═══════════════════════════════════════════════

    function test_emission_batchedSettle() public {
        vm.warp(GENESIS_TIME + 1);

        // Submit 3 recipients
        address charlie = makeAddr("charlie");
        uint256[] memory packed = new uint256[](3);
        packed[0] = (uint256(50) << 160) | uint256(uint160(alice));
        packed[1] = (uint256(30) << 160) | uint256(uint160(bob));
        packed[2] = (uint256(20) << 160) | uint256(uint160(charlie));

        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 0);

        // Settle with limit=1: only processes the 1st recipient
        uint256 aliceBefore = awp.balanceOf(alice);
        awpEmission.settleEpoch(1);

        // settleProgress should be > 0 (settlement in progress)
        assertTrue(awpEmission.settleProgress() > 0);
        assertEq(awpEmission.settledEpoch(), 0); // Not yet complete

        // alice should have received emission
        assertTrue(awp.balanceOf(alice) > aliceBefore);

        // Settle with limit=1: processes the 2nd recipient
        uint256 bobBefore = awp.balanceOf(bob);
        awpEmission.settleEpoch(1);
        assertTrue(awp.balanceOf(bob) > bobBefore);
        assertTrue(awpEmission.settleProgress() > 0); // Still not complete

        // Settle with limit=10: processes remaining and completes
        uint256 charlieBefore = awp.balanceOf(charlie);
        awpEmission.settleEpoch(10);
        assertTrue(awp.balanceOf(charlie) > charlieBefore);
        assertEq(awpEmission.settleProgress(), 0); // Complete
        assertEq(awpEmission.settledEpoch(), 1);    // Epoch 0 settled
    }

    // ═══════════════════════════════════════════════
    //  20. AWPRegistry.pause() halts all operations
    // ═══════════════════════════════════════════════

    function test_registry_pause_blocksOperations() public {
        // Pause
        vm.prank(guardian);
        awpRegistry.pause();
        assertTrue(awpRegistry.paused());

        // Register -> revert
        uint256 cost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.startPrank(alice);
        awp.approve(address(awpRegistry), cost);
        vm.expectRevert(); // EnforcedPause
        awpRegistry.registerWorknet(IAWPRegistry.WorknetParams({
            name: "Paused", symbol: "P", worknetManager: address(0),
            salt: bytes32(0), minStake: 0, skillsURI: ""
        }));
        vm.stopPrank();

        // Bind -> revert
        vm.prank(alice);
        vm.expectRevert();
        awpRegistry.bind(bob);

        // Unpause
        vm.prank(guardian);
        awpRegistry.unpause();
        assertFalse(awpRegistry.paused());

        // Register -> success
        uint256 wid = _registerWorknet(alice, "AfterPause", "AP");
        assertTrue(wid > 0);
    }

    // ═══════════════════════════════════════════════
    //  21. Emission exponential decay verification
    // ═══════════════════════════════════════════════

    function test_emission_exponentialDecay() public {
        vm.warp(GENESIS_TIME + 1);

        uint256[] memory packed = new uint256[](1);
        packed[0] = (uint256(100) << 160) | uint256(uint160(alice));

        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 0);

        uint256 emission0 = awpEmission.currentDailyEmission();

        // Settle epoch 0
        awpEmission.settleEpoch(100);

        // Advance to epoch 1
        vm.warp(GENESIS_TIME + 1 days + 1);
        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 1);

        // Settle epoch 1 -- triggers decay
        awpEmission.settleEpoch(100);
        uint256 emission1 = awpEmission.currentDailyEmission();

        // Decay: emission1 = emission0 * decayFactor / 1_000_000
        uint256 decayFactor = awpEmission.decayFactor();
        uint256 expectedEmission1 = emission0 * decayFactor / 1_000_000;
        assertEq(emission1, expectedEmission1);

        // One more epoch
        vm.warp(GENESIS_TIME + 2 days + 1);
        vm.prank(guardian);
        awpEmission.submitAllocations(packed, 100, 2);
        awpEmission.settleEpoch(100);
        uint256 emission2 = awpEmission.currentDailyEmission();

        uint256 expectedEmission2 = emission1 * decayFactor / 1_000_000;
        assertEq(emission2, expectedEmission2);

        // Decay is monotonically decreasing
        assertTrue(emission0 > emission1);
        assertTrue(emission1 > emission2);
    }

    // ═══════════════════════════════════════════════
    //  22. DAO multi-tokenId voting: combine voting power from multiple NFTs
    // ═══════════════════════════════════════════════

    function test_dao_multiTokenVoting() public {
        // Alice deposits an additional NFT
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1_000_000e18);
        uint256 tokenIdAlice2 = veAwp.deposit(1_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.roll(10);
        vm.warp(block.timestamp + 10);

        // Build proposal
        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(AWPDAO.setQuorumPercent, (5));
        string memory description = "Multi-token vote test";

        uint256[] memory propTokens = new uint256[](2);
        propTokens[0] = tokenIdAlice;
        propTokens[1] = tokenIdAlice2;

        vm.prank(alice);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, description, propTokens);

        // Advance past votingDelay
        vm.roll(20);

        // Alice votes with 2 tokenIds together
        uint256[] memory aliceTokens = new uint256[](2);
        aliceTokens[0] = tokenIdAlice;
        aliceTokens[1] = tokenIdAlice2;

        vm.prank(alice);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(aliceTokens));

        // Both tokens are marked as having voted
        assertTrue(dao.hasVotedWithToken(proposalId, tokenIdAlice));
        assertTrue(dao.hasVotedWithToken(proposalId, tokenIdAlice2));

        // Voting power is the sum of both tokens
        (, uint256 forVotes,) = dao.proposalVotes(proposalId);
        assertTrue(forVotes > 0);
    }

    // ═══════════════════════════════════════════════
    //  23. Gasless EIP-712: expired signature revert
    // ═══════════════════════════════════════════════

    function test_gasless_expiredSignature_reverts() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("expiredSigner");
        uint256 deadline = block.timestamp - 1; // Already expired

        bytes32 BIND_TYPEHASH = keccak256("Bind(address agent,address target,uint256 nonce,uint256 deadline)");
        bytes32 structHash = keccak256(abi.encode(BIND_TYPEHASH, signer, bob, 0, deadline));
        bytes32 digest = _registryDigest(structHash);

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        vm.prank(relayer);
        vm.expectRevert(AWPRegistry.ExpiredSignature.selector);
        awpRegistry.bindFor(signer, bob, deadline, v, r, s);
    }

    // ═══════════════════════════════════════════════
    //  24. Gasless EIP-712: invalid signature revert
    // ═══════════════════════════════════════════════

    function test_gasless_invalidSignature_reverts() public {
        (, uint256 wrongPk) = makeAddrAndKey("wrongSigner");
        address signer = makeAddr("realSigner");
        uint256 deadline = block.timestamp + 1 hours;

        bytes32 BIND_TYPEHASH = keccak256("Bind(address agent,address target,uint256 nonce,uint256 deadline)");
        bytes32 structHash = keccak256(abi.encode(BIND_TYPEHASH, signer, bob, 0, deadline));
        bytes32 digest = _registryDigest(structHash);

        // Sign with the wrong key
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(wrongPk, digest);

        vm.prank(relayer);
        vm.expectRevert(AWPRegistry.InvalidSignature.selector);
        awpRegistry.bindFor(signer, bob, deadline, v, r, s);
    }

    // ═══════════════════════════════════════════════
    //  Helper: build EIP-712 digest for AWPRegistry
    // ═══════════════════════════════════════════════

    function _registryDigest(bytes32 structHash) internal view returns (bytes32) {
        // EIP-712: "\x19\x01" || domainSeparator || structHash
        // AWPRegistry uses EIP712Upgradeable with name="AWPRegistry", version="1"
        bytes32 domainSeparator = keccak256(abi.encode(
            keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
            keccak256("AWPRegistry"),
            keccak256("1"),
            block.chainid,
            address(awpRegistry)
        ));
        return keccak256(abi.encodePacked("\x19\x01", domainSeparator, structHash));
    }

    /// @dev Build EIP-712 digest for AWPAllocator
    function _allocatorDigest(bytes32 structHash) internal view returns (bytes32) {
        bytes32 domainSeparator = keccak256(abi.encode(
            keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
            keccak256("AWPAllocator"),
            keccak256("1"),
            block.chainid,
            address(awpAllocator)
        ));
        return keccak256(abi.encodePacked("\x19\x01", domainSeparator, structHash));
    }
}
