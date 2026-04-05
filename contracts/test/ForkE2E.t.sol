// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {IAWPRegistry} from "../src/interfaces/IAWPRegistry.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {AWPAllocator} from "../src/core/AWPAllocator.sol";
import {AWPWorkNet} from "../src/core/AWPWorkNet.sol";
import {veAWP} from "../src/core/veAWP.sol";
import {IveAWP} from "../src/interfaces/IveAWP.sol";
import {AWPDAO} from "../src/governance/AWPDAO.sol";
import {WorknetTokenFactory} from "../src/token/WorknetTokenFactory.sol";
import {WorknetToken} from "../src/token/WorknetToken.sol";
import {IWorknetToken} from "../src/interfaces/IWorknetToken.sol";
import {TimelockControllerUpgradeable} from "@openzeppelin/contracts-upgradeable/governance/TimelockControllerUpgradeable.sol";

interface IWorknetManager {
    function setMerkleRoot(uint32 epoch, bytes32 root) external;
    function claim(uint32 epoch, uint256 amount, bytes32[] calldata proof) external;
    function isClaimed(uint32 epoch, address account) external view returns (bool);
}

interface IWorknetManagerConfig {
    function setStrategy(uint8 strategy) external;
    function setSlippageTolerance(uint16 bps) external;
    function setMinStrategyAmount(uint256 amount) external;
    function setStrategyPaused(bool paused) external;
    function transferToken(address token, address to, uint256 amount) external;
}

/// @title ForkE2E — End-to-end fork tests against live Base chain state
/// @dev Run: forge test --match-path test/ForkE2E.t.sol --fork-url $BASE_RPC_URL -v
contract ForkE2E is Test {
    // ── On-chain addresses (Base) ──
    AWPRegistry constant registry = AWPRegistry(0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A);
    AWPToken constant awp = AWPToken(0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1);
    veAWP constant veAwp = veAWP(0x0000b534C63D78212f1BDCc315165852793A00A8);
    AWPAllocator constant allocator = AWPAllocator(0x0000D6BB5e040E35081b3AaF59DD71b21C9800AA);
    AWPEmission constant emission = AWPEmission(0x3C9cB73f8B81083882c5308Cce4F31f93600EaA9);
    AWPWorkNet constant workNet = AWPWorkNet(0x00000bfbdEf8533E5F3228c9C846522D906100A7);
    WorknetTokenFactory constant factory = WorknetTokenFactory(0x000058EF25751Bb3687eB314185B46b942bE00AF);
    AWPDAO constant dao = AWPDAO(payable(0x00006879f79f3Da189b5D0fF6e58ad0127Cc0DA0));
    TimelockControllerUpgradeable constant treasury = TimelockControllerUpgradeable(payable(0x82562023a053025F3201785160CaE6051efD759e));
    address constant GUARDIAN = 0x000002bEfa6A1C99A710862Feb6dB50525dF00A3;

    address alice;
    address bob;
    address charlie;

    function setUp() public {
        alice = makeAddr("alice");
        bob = makeAddr("bob");
        charlie = makeAddr("charlie");

        // Fund test accounts with AWP
        deal(address(awp), alice, 10_000_000e18);
        deal(address(awp), bob, 5_000_000e18);
        deal(address(awp), charlie, 1_000_000e18);
        vm.deal(alice, 1 ether);
        vm.deal(bob, 1 ether);
        vm.deal(charlie, 1 ether);

        // Replace factory with vanityRule=0 version (immutable can't be vm.store'd)
        // Deploy a new factory with vanityRule=0, then etch its runtime code onto the existing address
        WorknetTokenFactory noVanity = new WorknetTokenFactory(address(this), 0);
        noVanity.setAddresses(address(registry));
        vm.etch(address(factory), address(noVanity).code);
    }

    // ═══════════════════════════════════════════════
    //  1. Registry — getRegistry returns correct addresses
    // ═══════════════════════════════════════════════

    function test_registry_getRegistry() public view {
        assertEq(registry.awpToken(), address(awp));
        assertEq(registry.awpWorkNet(), address(workNet));
        assertEq(registry.worknetTokenFactory(), address(factory));
        assertEq(registry.awpEmission(), address(emission));
        assertEq(registry.awpAllocator(), address(allocator));
        assertEq(registry.veAWP(), address(veAwp));
        assertEq(registry.treasury(), address(treasury));
    }

    // ═══════════════════════════════════════════════
    //  2. Binding + Recipient chain
    // ═══════════════════════════════════════════════

    function test_bind_setRecipient_resolveRecipient() public {
        vm.prank(alice);
        registry.bind(bob);

        vm.prank(bob);
        registry.setRecipient(charlie);

        assertEq(registry.resolveRecipient(alice), charlie);
        assertEq(registry.boundTo(alice), bob);
        assertEq(registry.recipient(bob), charlie);

        // Unbind
        vm.prank(alice);
        registry.unbind();
        assertEq(registry.resolveRecipient(alice), alice);
    }

    function test_bind_cycle_reverts() public {
        vm.prank(alice);
        registry.bind(bob);

        vm.prank(bob);
        vm.expectRevert(AWPRegistry.CycleDetected.selector);
        registry.bind(alice);
    }

    function test_bind_self_reverts() public {
        vm.prank(alice);
        vm.expectRevert(AWPRegistry.SelfBind.selector);
        registry.bind(alice);
    }

    // ═══════════════════════════════════════════════
    //  3. Delegation
    // ═══════════════════════════════════════════════

    function test_delegation() public {
        vm.prank(alice);
        registry.grantDelegate(bob);
        assertTrue(registry.delegates(alice, bob));

        vm.prank(alice);
        registry.revokeDelegate(bob);
        assertFalse(registry.delegates(alice, bob));
    }

    // ═══════════════════════════════════════════════
    //  4. veAWP — deposit, withdraw, transfer
    // ═══════════════════════════════════════════════

    function test_veAWP_deposit_withdraw() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1_000_000e18);
        uint256 tokenId = veAwp.deposit(1_000_000e18, 30 days);
        vm.stopPrank();

        (uint128 amount, uint64 lockEnd,) = veAwp.positions(tokenId);
        assertEq(amount, 1_000_000e18);
        assertTrue(lockEnd > block.timestamp);
        assertEq(veAwp.getUserTotalStaked(alice), 1_000_000e18);
        assertEq(veAwp.totalStaked(), 1_000_000e18);

        // Cannot withdraw before lock expires
        vm.prank(alice);
        vm.expectRevert(veAWP.LockNotExpired.selector);
        veAwp.withdraw(tokenId);

        // Warp past lock
        vm.warp(block.timestamp + 31 days);

        vm.prank(alice);
        veAwp.withdraw(tokenId);
        assertEq(veAwp.getUserTotalStaked(alice), 0);
        assertEq(veAwp.totalStaked(), 0);
    }

    function test_veAWP_depositWithPermit() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("permitUser");
        deal(address(awp), signer, 1_000_000e18);

        uint256 amount = 500_000e18;
        uint256 deadline = block.timestamp + 1 hours;

        bytes32 PERMIT_TYPEHASH = keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)");
        bytes32 domainSeparator = awp.DOMAIN_SEPARATOR();
        bytes32 structHash = keccak256(abi.encode(PERMIT_TYPEHASH, signer, address(veAwp), amount, awp.nonces(signer), deadline));
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", domainSeparator, structHash));

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        vm.prank(signer);
        uint256 tokenId = veAwp.depositWithPermit(amount, 30 days, deadline, v, r, s);
        assertGt(tokenId, 0);
        assertEq(veAwp.getUserTotalStaked(signer), amount);
    }

    function test_veAWP_partialWithdraw() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1_000_000e18);
        uint256 tokenId = veAwp.deposit(1_000_000e18, 1 days);
        vm.stopPrank();

        vm.warp(block.timestamp + 2 days);

        vm.prank(alice);
        veAwp.partialWithdraw(tokenId, 400_000e18);

        (uint128 remaining,,) = veAwp.positions(tokenId);
        assertEq(remaining, 600_000e18);
        assertEq(veAwp.getUserTotalStaked(alice), 600_000e18);
    }

    function test_veAWP_transfer_blocked_by_allocation() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1_000_000e18);
        uint256 tokenId = veAwp.deposit(1_000_000e18, 30 days);
        vm.stopPrank();

        // Allocate most of stake
        vm.prank(alice);
        allocator.allocate(alice, alice, 845300000001, 900_000e18);

        // Transfer should fail (insufficient unallocated)
        vm.prank(alice);
        vm.expectRevert(veAWP.InsufficientUnallocated.selector);
        veAwp.transferFrom(alice, bob, tokenId);

        // Deallocate then transfer succeeds
        vm.prank(alice);
        allocator.deallocateAll(alice, alice, 845300000001);

        vm.prank(alice);
        veAwp.transferFrom(alice, bob, tokenId);
        assertEq(veAwp.ownerOf(tokenId), bob);
    }

    // ═══════════════════════════════════════════════
    //  5. AWPAllocator — allocate, deallocate, reallocate
    // ═══════════════════════════════════════════════

    function test_allocator_full_flow() public {
        // Stake first
        vm.startPrank(alice);
        awp.approve(address(veAwp), 2_000_000e18);
        veAwp.deposit(2_000_000e18, 30 days);
        vm.stopPrank();

        uint256 wid = 845300000001;

        // Allocate
        vm.prank(alice);
        allocator.allocate(alice, alice, wid, 1_000_000e18);
        assertEq(allocator.getAgentStake(alice, alice, wid), 1_000_000e18);
        assertEq(allocator.userTotalAllocated(alice), 1_000_000e18);
        assertEq(allocator.worknetTotalStake(wid), 1_000_000e18);

        // Deallocate partial
        vm.prank(alice);
        allocator.deallocate(alice, alice, wid, 500_000e18);
        assertEq(allocator.getAgentStake(alice, alice, wid), 500_000e18);

        // Reallocate to another worknet
        uint256 wid2 = 845300000002;
        vm.prank(alice);
        allocator.reallocate(alice, alice, wid, alice, wid2, 500_000e18);
        assertEq(allocator.getAgentStake(alice, alice, wid), 0);
        assertEq(allocator.getAgentStake(alice, alice, wid2), 500_000e18);
    }

    function test_allocator_delegate() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 2_000_000e18);
        veAwp.deposit(2_000_000e18, 30 days);
        registry.grantDelegate(bob);
        vm.stopPrank();

        // Bob allocates on behalf of alice
        vm.prank(bob);
        allocator.allocate(alice, bob, 845300000001, 1_000_000e18);
        assertEq(allocator.getAgentStake(alice, bob, 845300000001), 1_000_000e18);
    }

    function test_allocator_insufficient_reverts() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 100e18);
        veAwp.deposit(100e18, 30 days);
        vm.stopPrank();

        vm.prank(alice);
        vm.expectRevert(AWPAllocator.InsufficientUnallocated.selector);
        allocator.allocate(alice, alice, 845300000001, 200e18);
    }

    // ═══════════════════════════════════════════════
    //  6. Worknet registration + activation (full lifecycle)
    // ═══════════════════════════════════════════════

    function test_worknet_register_activate() public {
        uint256 cost = registry.initialAlphaMint() * registry.initialAlphaPrice() / 1e18;

        vm.startPrank(alice);
        awp.approve(address(registry), cost);
        uint256 wid = registry.registerWorknet(IAWPRegistry.WorknetParams({
            name: "ForkTest Worknet",
            symbol: "FTW",
            worknetManager: address(0),
            salt: bytes32(0),
            minStake: 0,
            skillsURI: ""
        }));
        vm.stopPrank();

        // Verify pending
        IAWPRegistry.WorknetInfo memory info = registry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Pending));

        // Guardian activates
        vm.prank(GUARDIAN);
        registry.activateWorknet(wid);

        info = registry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Active));
        assertTrue(registry.isWorknetActive(wid));

        // Verify WorknetToken deployed
        IAWPRegistry.WorknetFullInfo memory full = registry.getWorknetFull(wid);
        assertTrue(full.worknetToken != address(0));
        assertTrue(full.worknetManager != address(0));

        WorknetToken wt = WorknetToken(full.worknetToken);
        assertEq(wt.name(), "ForkTest Worknet");
        assertEq(wt.symbol(), "FTW");
        assertTrue(wt.initialized());
        assertEq(wt.minter(), full.worknetManager);

        // Verify NFT minted
        assertEq(workNet.ownerOf(wid), alice);
    }

    function test_worknet_register_cancel_refund() public {
        uint256 cost = registry.initialAlphaMint() * registry.initialAlphaPrice() / 1e18;

        vm.startPrank(alice);
        awp.approve(address(registry), cost);
        uint256 balBefore = awp.balanceOf(alice);
        uint256 wid = registry.registerWorknet(IAWPRegistry.WorknetParams({
            name: "CancelMe", symbol: "CXL",
            worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        }));

        registry.cancelWorknet(wid);
        vm.stopPrank();

        assertEq(awp.balanceOf(alice), balBefore); // full refund
    }

    function test_worknet_reject_by_guardian() public {
        uint256 cost = registry.initialAlphaMint() * registry.initialAlphaPrice() / 1e18;

        vm.startPrank(alice);
        awp.approve(address(registry), cost);
        uint256 wid = registry.registerWorknet(IAWPRegistry.WorknetParams({
            name: "RejectMe", symbol: "REJ",
            worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        }));
        vm.stopPrank();

        uint256 balBefore = awp.balanceOf(alice);
        vm.prank(GUARDIAN);
        registry.rejectWorknet(wid);
        assertEq(awp.balanceOf(alice), balBefore + cost); // refund to owner
    }

    // ═══════════════════════════════════════════════
    //  7. Worknet status transitions
    // ═══════════════════════════════════════════════

    function test_worknet_pause_resume_ban_unban() public {
        uint256 wid = _activateWorknet(alice, "StatusTest", "STS");

        // Pause by owner
        vm.prank(alice);
        registry.pauseWorknet(wid);
        assertFalse(registry.isWorknetActive(wid));

        // Resume by owner
        vm.prank(alice);
        registry.resumeWorknet(wid);
        assertTrue(registry.isWorknetActive(wid));

        // Ban by guardian
        vm.prank(GUARDIAN);
        registry.banWorknet(wid);
        IAWPRegistry.WorknetInfo memory info = registry.getWorknet(wid);
        assertEq(uint8(info.status), uint8(IAWPRegistry.WorknetStatus.Banned));

        // Unban by guardian
        vm.prank(GUARDIAN);
        registry.unbanWorknet(wid);
        assertTrue(registry.isWorknetActive(wid));
    }

    // ═══════════════════════════════════════════════
    //  8. AWPEmission — submit + settle
    // ═══════════════════════════════════════════════

    function test_emission_submit_settle() public {
        // Use plain EOA addresses (no code) to avoid onTransferReceived issues
        address r1 = address(0xE001);
        address r2 = address(0xE002);

        uint256 epoch = emission.settledEpoch();

        uint256 baseTime = emission.baseTime();
        uint256 epochDuration = emission.epochDuration();
        uint256 baseEpoch = emission.baseEpoch();
        uint256 targetTime = baseTime + (epoch + 1 - baseEpoch) * epochDuration + 1;
        if (block.timestamp < targetTime) vm.warp(targetTime);

        uint256[] memory packed = new uint256[](2);
        packed[0] = (uint256(70) << 160) | uint256(uint160(r1));
        packed[1] = (uint256(30) << 160) | uint256(uint160(r2));

        vm.prank(GUARDIAN);
        emission.submitAllocations(packed, 100, epoch);

        emission.settleEpoch(100);

        assertTrue(awp.balanceOf(r1) > 0);
        assertTrue(awp.balanceOf(r2) > 0);
        assertEq(emission.settledEpoch(), epoch + 1);
    }

    function test_emission_batched_settle() public {
        address r1 = address(0xE011);
        address r2 = address(0xE012);
        address r3 = address(0xE013);

        uint256 epoch = emission.settledEpoch();

        uint256 baseTime = emission.baseTime();
        uint256 epochDuration = emission.epochDuration();
        uint256 baseEpoch = emission.baseEpoch();
        uint256 targetTime = baseTime + (epoch + 1 - baseEpoch) * epochDuration + 1;
        if (block.timestamp < targetTime) vm.warp(targetTime);

        uint256[] memory packed = new uint256[](3);
        packed[0] = (uint256(50) << 160) | uint256(uint160(r1));
        packed[1] = (uint256(30) << 160) | uint256(uint160(r2));
        packed[2] = (uint256(20) << 160) | uint256(uint160(r3));

        vm.prank(GUARDIAN);
        emission.submitAllocations(packed, 100, epoch);

        emission.settleEpoch(1);
        assertTrue(emission.settleProgress() > 0);

        emission.settleEpoch(1);
        assertTrue(emission.settleProgress() > 0);

        emission.settleEpoch(10);
        assertEq(emission.settleProgress(), 0);
        assertEq(emission.settledEpoch(), epoch + 1);
    }

    // ═══════════════════════════════════════════════
    //  9. AWPDAO — propose, vote, queue, execute
    // ═══════════════════════════════════════════════

    function test_dao_full_governance() public {
        // Stake for voting power
        vm.startPrank(alice);
        awp.approve(address(veAwp), 5_000_000e18);
        uint256 tokenIdA = veAwp.deposit(5_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.startPrank(bob);
        awp.approve(address(veAwp), 3_000_000e18);
        uint256 tokenIdB = veAwp.deposit(3_000_000e18, 54 weeks);
        vm.stopPrank();

        // Advance so createdAt < proposalCreatedAt
        vm.warp(block.timestamp + 10);

        // Send AWP to treasury (timelock) for the proposal to transfer
        deal(address(awp), address(treasury), 1_000e18);

        // Build proposal: treasury transfers 500 AWP to charlie
        address[] memory targets = new address[](1);
        targets[0] = address(awp);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(IERC20.transfer, (charlie, 500e18));
        string memory description = "Fork test: transfer 500 AWP";

        uint256[] memory propTokens = new uint256[](1);
        propTokens[0] = tokenIdA;

        vm.prank(alice);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, description, propTokens);

        // Advance past votingDelay (1 day)
        vm.warp(block.timestamp + 1 days + 1);

        // Vote
        uint256[] memory aliceTokens = new uint256[](1);
        aliceTokens[0] = tokenIdA;
        vm.prank(alice);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(aliceTokens));

        uint256[] memory bobTokens = new uint256[](1);
        bobTokens[0] = tokenIdB;
        vm.prank(bob);
        dao.castVoteWithReasonAndParams(proposalId, 1, "", abi.encode(bobTokens));

        // Advance past votingPeriod (7 days)
        vm.warp(block.timestamp + 7 days + 1);

        // Queue
        dao.queue(targets, values, calldatas, keccak256(bytes(description)));

        // Advance past treasury minDelay (2 days)
        vm.warp(block.timestamp + 2 days + 1);

        // Execute
        uint256 charlieBal = awp.balanceOf(charlie);
        dao.execute(targets, values, calldatas, keccak256(bytes(description)));
        assertEq(awp.balanceOf(charlie) - charlieBal, 500e18);
    }

    function test_dao_signal_proposal() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 5_000_000e18);
        uint256 tokenId = veAwp.deposit(5_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.warp(block.timestamp + 10);

        uint256[] memory tokens = new uint256[](1);
        tokens[0] = tokenId;

        vm.prank(alice);
        uint256 proposalId = dao.signalPropose("Should we do X?", tokens);
        assertTrue(dao.isSignalProposal(proposalId));
        assertFalse(dao.proposalNeedsQueuing(proposalId));
    }

    function test_dao_guardian_cancel() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 5_000_000e18);
        uint256 tokenId = veAwp.deposit(5_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.warp(block.timestamp + 10);

        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(AWPDAO.setQuorumPercent, (5));

        uint256[] memory tokens = new uint256[](1);
        tokens[0] = tokenId;

        vm.prank(alice);
        uint256 proposalId = dao.proposeWithTokens(targets, values, calldatas, "malicious", tokens);

        vm.prank(GUARDIAN);
        dao.guardianCancel(targets, values, calldatas, keccak256("malicious"));
        assertEq(uint8(dao.state(proposalId)), 2); // Canceled
    }

    // ═══════════════════════════════════════════════
    //  10. Gasless EIP-712 — bindFor, allocateFor
    // ═══════════════════════════════════════════════

    function test_gasless_bindFor() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("eip712user");

        uint256 deadline = block.timestamp + 1 hours;
        bytes32 BIND_TYPEHASH = keccak256("Bind(address agent,address target,uint256 nonce,uint256 deadline)");
        bytes32 structHash = keccak256(abi.encode(BIND_TYPEHASH, signer, bob, registry.nonces(signer), deadline));

        bytes32 domainSeparator = keccak256(abi.encode(
            keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
            keccak256("AWPRegistry"), keccak256("1"), block.chainid, address(registry)
        ));
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", domainSeparator, structHash));

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        vm.prank(charlie); // relayer
        registry.bindFor(signer, bob, deadline, v, r, s);
        assertEq(registry.boundTo(signer), bob);
    }

    // ═══════════════════════════════════════════════
    //  11. AWPWorkNet — metadata updates
    // ═══════════════════════════════════════════════

    function test_worknet_nft_metadata() public {
        uint256 wid = _activateWorknet(alice, "NFT Test", "NFT");

        vm.prank(alice);
        workNet.setSkillsURI(wid, "https://skills.test.com");

        vm.prank(alice);
        workNet.setMinStake(wid, 100e18);

        AWPWorkNet.WorknetMeta memory meta = workNet.getWorknetMeta(wid);
        assertEq(meta.skillsURI, "https://skills.test.com");
        assertEq(meta.minStake, 100e18);

        // tokenURI should return valid JSON
        string memory uri = workNet.tokenURI(wid);
        assertTrue(bytes(uri).length > 0);
    }

    function test_worknet_nft_delegate_update() public {
        uint256 wid = _activateWorknet(alice, "DelegateNFT", "DLG");

        vm.prank(alice);
        registry.grantDelegate(bob);

        // Bob as delegate can update metadata
        vm.prank(bob);
        workNet.setSkillsURI(wid, "https://delegate-update.com");

        AWPWorkNet.WorknetMeta memory meta = workNet.getWorknetMeta(wid);
        assertEq(meta.skillsURI, "https://delegate-update.com");
    }

    // ═══════════════════════════════════════════════
    //  12. Registry pause/unpause
    // ═══════════════════════════════════════════════

    function test_registry_pause() public {
        vm.prank(GUARDIAN);
        registry.pause();
        assertTrue(registry.paused());

        // Register should fail
        vm.startPrank(alice);
        awp.approve(address(registry), type(uint256).max);
        vm.expectRevert();
        registry.registerWorknet(IAWPRegistry.WorknetParams({
            name: "Paused", symbol: "P",
            worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        }));
        vm.stopPrank();

        // Unpause
        vm.prank(GUARDIAN);
        registry.unpause();
        assertFalse(registry.paused());
    }

    // ═══════════════════════════════════════════════
    //  13. Emission decay verification
    // ═══════════════════════════════════════════════

    function test_emission_decay() public {
        uint256 epoch = emission.settledEpoch();
        uint256 emissionBefore = emission.currentDailyEmission();
        uint256 decayFactor = emission.decayFactor();

        uint256[] memory packed = new uint256[](1);
        packed[0] = (uint256(100) << 160) | uint256(uint160(alice));

        // Settle current epoch
        vm.prank(GUARDIAN);
        emission.submitAllocations(packed, 100, epoch);

        if (epoch <= emission.currentEpoch()) {
            emission.settleEpoch(100);

            uint256 emissionAfter = emission.currentDailyEmission();
            if (epoch > 0) {
                // Decay applied
                assertEq(emissionAfter, emissionBefore * decayFactor / 1_000_000);
            }
        }
    }

    // ═══════════════════════════════════════════════
    //  14. WorknetTokenFactory — predictDeployAddress
    // ═══════════════════════════════════════════════

    function test_factory_predictAddress() public view {
        bytes32 salt = bytes32(uint256(42));
        address predicted = factory.predictDeployAddress(salt);
        assertTrue(predicted != address(0));

        // Same salt always gives same address (universal salt)
        address predicted2 = factory.predictDeployAddress(salt);
        assertEq(predicted, predicted2);
    }

    // ═══════════════════════════════════════════════
    //  15. Edge cases
    // ═══════════════════════════════════════════════

    function test_register_empty_name_reverts() public {
        uint256 cost = registry.initialAlphaMint() * registry.initialAlphaPrice() / 1e18;
        vm.startPrank(alice);
        awp.approve(address(registry), cost);
        vm.expectRevert(AWPRegistry.InvalidWorknetName.selector);
        registry.registerWorknet(IAWPRegistry.WorknetParams({
            name: "", symbol: "X",
            worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        }));
        vm.stopPrank();
    }

    function test_register_json_unsafe_name_reverts() public {
        uint256 cost = registry.initialAlphaMint() * registry.initialAlphaPrice() / 1e18;
        vm.startPrank(alice);
        awp.approve(address(registry), cost);
        vm.expectRevert(AWPRegistry.JsonUnsafeCharacter.selector);
        registry.registerWorknet(IAWPRegistry.WorknetParams({
            name: 'bad"quote', symbol: "X",
            worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        }));
        vm.stopPrank();
    }

    function test_veAWP_voting_power() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1_000_000e18);
        uint256 tokenId = veAwp.deposit(1_000_000e18, 54 weeks);
        vm.stopPrank();

        uint256 power = veAwp.getVotingPower(tokenId);
        assertTrue(power > 0);
        // power = amount * sqrt(min(remaining, 54weeks) / 7days)
        // remaining ≈ 54 weeks (just deposited), sqrt(54) ≈ 7.348
        // But remaining could be slightly less due to block time, so use wider bounds
        assertTrue(power > 6_000_000e18);
        assertTrue(power < 8_000_000e18);
    }

    function test_allocator_zero_worknetId_reverts() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 1_000_000e18);
        veAwp.deposit(1_000_000e18, 30 days);
        vm.stopPrank();

        vm.prank(alice);
        vm.expectRevert(AWPAllocator.ZeroWorknetId.selector);
        allocator.allocate(alice, alice, 0, 100e18);
    }

    // ═══════════════════════════════════════════════
    //  16. WorknetToken distribution via Merkle claim
    // ═══════════════════════════════════════════════

    function test_worknetToken_merkleClaim() public {
        uint256 wid = _activateWorknet(alice, "ClaimTest", "CLM");
        IAWPRegistry.WorknetFullInfo memory full = registry.getWorknetFull(wid);
        address wmAddr = full.worknetManager;

        // Build merkle tree: bob claims 100e18
        uint256 claimAmount = 100e18;
        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(bob, claimAmount))));
        bytes32 root = leaf; // single-leaf tree

        // Alice (worknet admin) sets merkle root — impersonate the WorknetManager admin
        vm.prank(alice);
        IWorknetManager(wmAddr).setMerkleRoot(0, root);

        // Bob claims
        bytes32[] memory proof = new bytes32[](0);
        vm.prank(bob);
        IWorknetManager(wmAddr).claim(0, claimAmount, proof);

        // Verify WorknetToken minted to bob (via resolveRecipient)
        IWorknetToken wt = IWorknetToken(full.worknetToken);
        assertEq(wt.balanceOf(bob), claimAmount);

        // Double claim reverts
        vm.prank(bob);
        vm.expectRevert();
        IWorknetManager(wmAddr).claim(0, claimAmount, proof);
    }

    function test_worknetToken_claim_with_binding() public {
        uint256 wid = _activateWorknet(alice, "BindClaim", "BCL");
        IAWPRegistry.WorknetFullInfo memory full = registry.getWorknetFull(wid);
        address wmAddr = full.worknetManager;

        // Bob binds to charlie → claim should mint to charlie
        vm.prank(bob);
        registry.bind(charlie);

        vm.prank(charlie);
        registry.setRecipient(address(0xBEEF));

        uint256 claimAmount = 50e18;
        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(bob, claimAmount))));

        vm.prank(alice);
        IWorknetManager(wmAddr).setMerkleRoot(1, leaf);

        vm.prank(bob);
        IWorknetManager(wmAddr).claim(1, claimAmount, new bytes32[](0));

        // resolveRecipient(bob) → follows bind chain → charlie → recipient 0xBEEF
        IWorknetToken wt = IWorknetToken(full.worknetToken);
        assertEq(wt.balanceOf(address(0xBEEF)), claimAmount);
    }

    function test_worknetToken_invalidProof_reverts() public {
        uint256 wid = _activateWorknet(alice, "BadProof", "BPF");
        IAWPRegistry.WorknetFullInfo memory full = registry.getWorknetFull(wid);
        address wmAddr = full.worknetManager;

        vm.prank(alice);
        IWorknetManager(wmAddr).setMerkleRoot(0, keccak256("real_root"));

        vm.prank(bob);
        vm.expectRevert();
        IWorknetManager(wmAddr).claim(0, 100e18, new bytes32[](0));
    }

    function test_worknetToken_timeCap() public {
        uint256 wid = _activateWorknet(alice, "TimeCap", "TCP");
        IAWPRegistry.WorknetFullInfo memory full = registry.getWorknetFull(wid);

        WorknetToken wt = WorknetToken(full.worknetToken);
        assertTrue(wt.initialized());
        assertTrue(wt.supplyAtLock() > 0);

        // Time-based cap: check currentMintableLimit grows with time
        uint256 limit0 = wt.currentMintableLimit();

        vm.warp(block.timestamp + 30 days);
        uint256 limit30 = wt.currentMintableLimit();

        assertTrue(limit30 > limit0);
    }

    // ═══════════════════════════════════════════════
    //  17. DAO voting edge cases
    // ═══════════════════════════════════════════════

    function test_dao_multi_tokenId_voting() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 6_000_000e18);
        uint256 t1 = veAwp.deposit(3_000_000e18, 54 weeks);
        uint256 t2 = veAwp.deposit(3_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.warp(block.timestamp + 10);

        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(AWPDAO.setQuorumPercent, (5));

        uint256[] memory propTokens = new uint256[](2);
        propTokens[0] = t1;
        propTokens[1] = t2;

        vm.prank(alice);
        uint256 pid = dao.proposeWithTokens(targets, values, calldatas, "multi-token", propTokens);

        vm.warp(block.timestamp + 1 days + 1);

        // Vote with both tokens
        uint256[] memory voteTokens = new uint256[](2);
        voteTokens[0] = t1;
        voteTokens[1] = t2;
        vm.prank(alice);
        dao.castVoteWithReasonAndParams(pid, 1, "", abi.encode(voteTokens));

        assertTrue(dao.hasVotedWithToken(pid, t1));
        assertTrue(dao.hasVotedWithToken(pid, t2));
        (, uint256 forVotes,) = dao.proposalVotes(pid);
        assertTrue(forVotes > 0);
    }

    function test_dao_double_vote_reverts() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 5_000_000e18);
        uint256 tokenId = veAwp.deposit(5_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.warp(block.timestamp + 10);

        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(AWPDAO.setQuorumPercent, (5));

        uint256[] memory tokens = new uint256[](1);
        tokens[0] = tokenId;

        vm.prank(alice);
        uint256 pid = dao.proposeWithTokens(targets, values, calldatas, "double-vote", tokens);

        vm.warp(block.timestamp + 1 days + 1);

        vm.prank(alice);
        dao.castVoteWithReasonAndParams(pid, 1, "", abi.encode(tokens));

        // Double vote reverts
        vm.prank(alice);
        vm.expectRevert(AWPDAO.TokenAlreadyVoted.selector);
        dao.castVoteWithReasonAndParams(pid, 1, "", abi.encode(tokens));
    }

    function test_dao_expired_lock_vote_reverts() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 5_000_000e18);
        uint256 tokenId = veAwp.deposit(5_000_000e18, 2 days);
        vm.stopPrank();

        vm.warp(block.timestamp + 10);

        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(AWPDAO.setQuorumPercent, (5));

        uint256[] memory tokens = new uint256[](1);
        tokens[0] = tokenId;

        vm.prank(alice);
        uint256 pid = dao.proposeWithTokens(targets, values, calldatas, "expired-lock", tokens);

        // Warp past lock but within votingPeriod
        vm.warp(block.timestamp + 3 days);

        vm.prank(alice);
        vm.expectRevert(AWPDAO.LockExpired.selector);
        dao.castVoteWithReasonAndParams(pid, 1, "", abi.encode(tokens));
    }

    function test_dao_proposal_defeated() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 5_000_000e18);
        uint256 tAlice = veAwp.deposit(5_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.startPrank(bob);
        awp.approve(address(veAwp), 5_000_000e18);
        uint256 tBob = veAwp.deposit(5_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.warp(block.timestamp + 10);

        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(AWPDAO.setQuorumPercent, (99));

        uint256[] memory tokens = new uint256[](1);
        tokens[0] = tAlice;

        vm.prank(alice);
        uint256 pid = dao.proposeWithTokens(targets, values, calldatas, "defeated", tokens);

        vm.warp(block.timestamp + 1 days + 1);

        // Alice votes For
        vm.prank(alice);
        dao.castVoteWithReasonAndParams(pid, 1, "", abi.encode(tokens));

        // Bob votes Against (more weight or equal)
        tokens[0] = tBob;
        vm.prank(bob);
        dao.castVoteWithReasonAndParams(pid, 0, "", abi.encode(tokens));

        vm.warp(block.timestamp + 7 days + 1);

        // Proposal defeated (forVotes <= againstVotes)
        assertEq(uint8(dao.state(pid)), 3); // Defeated
    }

    function test_dao_insufficient_stake_to_propose_reverts() public {
        // Charlie has 1M AWP, threshold is 200K — should work
        vm.startPrank(charlie);
        awp.approve(address(veAwp), 1_000_000e18);
        uint256 tokenId = veAwp.deposit(1_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.warp(block.timestamp + 10);

        // Now test with tiny stake (below 200K threshold)
        address tiny = makeAddr("tiny");
        deal(address(awp), tiny, 100e18);
        vm.startPrank(tiny);
        awp.approve(address(veAwp), 100e18);
        uint256 tinyId = veAwp.deposit(100e18, 54 weeks);
        vm.stopPrank();

        vm.warp(block.timestamp + 10);

        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(AWPDAO.setQuorumPercent, (5));

        uint256[] memory tokens = new uint256[](1);
        tokens[0] = tinyId;

        vm.prank(tiny);
        vm.expectRevert(); // GovernorInsufficientProposerVotes
        dao.proposeWithTokens(targets, values, calldatas, "too-small", tokens);
    }

    // ═══════════════════════════════════════════════
    //  18. veAWP — addToPosition, batchWithdraw
    // ═══════════════════════════════════════════════

    function test_veAWP_addToPosition() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 2_000_000e18);
        uint256 tokenId = veAwp.deposit(1_000_000e18, 30 days);

        // Add more funds + extend lock
        uint64 newLock = uint64(block.timestamp + 60 days);
        veAwp.addToPosition(tokenId, 500_000e18, newLock);
        vm.stopPrank();

        (uint128 amount, uint64 lockEnd,) = veAwp.positions(tokenId);
        assertEq(amount, 1_500_000e18);
        assertEq(lockEnd, newLock);
        assertEq(veAwp.getUserTotalStaked(alice), 1_500_000e18);
    }

    function test_veAWP_batchWithdraw() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 3_000_000e18);
        uint256 t1 = veAwp.deposit(1_000_000e18, 1 days);
        uint256 t2 = veAwp.deposit(1_000_000e18, 1 days);
        uint256 t3 = veAwp.deposit(1_000_000e18, 1 days);
        vm.stopPrank();

        vm.warp(block.timestamp + 2 days);

        uint256[] memory tokenIds = new uint256[](3);
        tokenIds[0] = t1;
        tokenIds[1] = t2;
        tokenIds[2] = t3;

        uint256 balBefore = awp.balanceOf(alice);
        vm.prank(alice);
        veAwp.batchWithdraw(tokenIds);

        assertEq(awp.balanceOf(alice) - balBefore, 3_000_000e18);
        assertEq(veAwp.getUserTotalStaked(alice), 0);
    }

    // ═══════════════════════════════════════════════
    //  19. Multi-worknet allocation tracking
    // ═══════════════════════════════════════════════

    function test_multi_worknet_allocation() public {
        uint256 wid1 = _activateWorknet(alice, "Net1", "N1");
        uint256 wid2 = _activateWorknet(bob, "Net2", "N2");

        vm.startPrank(charlie);
        awp.approve(address(veAwp), 1_000_000e18);
        veAwp.deposit(1_000_000e18, 30 days);

        // Allocate to two worknets
        allocator.allocate(charlie, charlie, wid1, 400_000e18);
        allocator.allocate(charlie, charlie, wid2, 300_000e18);
        vm.stopPrank();

        assertEq(allocator.userTotalAllocated(charlie), 700_000e18);
        assertEq(allocator.worknetTotalStake(wid1), 400_000e18);
        assertEq(allocator.worknetTotalStake(wid2), 300_000e18);

        // Cannot over-allocate
        vm.prank(charlie);
        vm.expectRevert(AWPAllocator.InsufficientUnallocated.selector);
        allocator.allocate(charlie, charlie, wid1, 400_000e18);
    }

    // ═══════════════════════════════════════════════
    //  20. Gasless registerWorknetFor
    // ═══════════════════════════════════════════════

    function test_gasless_registerWorknetFor() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("regSigner");
        deal(address(awp), signer, 10_000_000e18);

        uint256 cost = registry.initialAlphaMint() * registry.initialAlphaPrice() / 1e18;
        vm.prank(signer);
        awp.approve(address(registry), cost);

        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = registry.nonces(signer);

        IAWPRegistry.WorknetParams memory params = IAWPRegistry.WorknetParams({
            name: "Gasless", symbol: "GLS",
            worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        });

        bytes32 WORKNET_PARAMS_TYPEHASH = keccak256("WorknetParams(string name,string symbol,address worknetManager,bytes32 salt,uint128 minStake,string skillsURI)");
        bytes32 paramsHash = keccak256(abi.encode(
            WORKNET_PARAMS_TYPEHASH,
            keccak256(bytes(params.name)), keccak256(bytes(params.symbol)),
            params.worknetManager, params.salt, params.minStake, keccak256(bytes(params.skillsURI))
        ));
        bytes32 REGISTER_TYPEHASH = keccak256("RegisterWorknet(address user,WorknetParams params,uint256 nonce,uint256 deadline)WorknetParams(string name,string symbol,address worknetManager,bytes32 salt,uint128 minStake,string skillsURI)");
        bytes32 structHash = keccak256(abi.encode(REGISTER_TYPEHASH, signer, paramsHash, nonce, deadline));

        bytes32 domainSeparator = keccak256(abi.encode(
            keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
            keccak256("AWPRegistry"), keccak256("1"), block.chainid, address(registry)
        ));
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", domainSeparator, structHash));

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        vm.prank(charlie); // relayer
        uint256 wid = registry.registerWorknetFor(signer, params, deadline, v, r, s);
        assertTrue(wid > 0);
    }

    // ═══════════════════════════════════════════════
    //  21. WorknetManager — strategy config & role management
    // ═══════════════════════════════════════════════

    function test_worknetManager_config() public {
        uint256 wid = _activateWorknet(alice, "ConfigTest", "CFG");
        IAWPRegistry.WorknetFullInfo memory full = registry.getWorknetFull(wid);
        address wmAddr = full.worknetManager;

        // Alice is admin of the WorknetManager
        // Set strategy
        vm.prank(alice);
        IWorknetManagerConfig(wmAddr).setStrategy(1); // AddLiquidity

        // Set slippage
        vm.prank(alice);
        IWorknetManagerConfig(wmAddr).setSlippageTolerance(300);

        // Set min strategy amount
        vm.prank(alice);
        IWorknetManagerConfig(wmAddr).setMinStrategyAmount(1e18);

        // Pause strategy
        vm.prank(alice);
        IWorknetManagerConfig(wmAddr).setStrategyPaused(true);

        // Non-admin cannot configure
        vm.prank(bob);
        vm.expectRevert();
        IWorknetManagerConfig(wmAddr).setStrategy(0);
    }

    function test_worknetManager_transferToken() public {
        uint256 wid = _activateWorknet(alice, "TransferTest", "TFR");
        IAWPRegistry.WorknetFullInfo memory full = registry.getWorknetFull(wid);
        address wmAddr = full.worknetManager;

        // Give WM some AWP
        deal(address(awp), wmAddr, 1000e18);

        // Alice (admin with TRANSFER_ROLE) transfers
        vm.prank(alice);
        IWorknetManagerConfig(wmAddr).transferToken(address(awp), bob, 500e18);
        assertEq(awp.balanceOf(bob), 5_000_000e18 + 500e18);
    }

    // ═══════════════════════════════════════════════
    //  22. Gasless allocateFor / deallocateFor
    // ═══════════════════════════════════════════════

    function test_gasless_allocateFor_deallocateFor() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("allocSigner");
        deal(address(awp), signer, 2_000_000e18);

        vm.startPrank(signer);
        awp.approve(address(veAwp), 1_000_000e18);
        veAwp.deposit(1_000_000e18, 30 days);
        vm.stopPrank();

        uint256 worknetId = 845300000001;
        uint256 deadline = block.timestamp + 1 hours;
        uint256 amount = 500_000e18;

        // allocateFor
        bytes32 ALLOCATE_TYPEHASH = keccak256("Allocate(address staker,address agent,uint256 worknetId,uint256 amount,uint256 nonce,uint256 deadline)");
        // EIP-712 domain fixed via MigrateAllocatorEIP712 (reinitializer(2))
        bytes32 allocDomain = keccak256(abi.encode(
            keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
            keccak256("AWPAllocator"), keccak256("1"), block.chainid, address(allocator)
        ));

        uint256 nonce = allocator.nonces(signer);
        bytes32 structHash = keccak256(abi.encode(ALLOCATE_TYPEHASH, signer, signer, worknetId, amount, nonce, deadline));
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", allocDomain, structHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        vm.prank(charlie); // relayer
        allocator.allocateFor(signer, signer, worknetId, amount, deadline, v, r, s);
        assertEq(allocator.getAgentStake(signer, signer, worknetId), amount);

        // deallocateFor
        bytes32 DEALLOCATE_TYPEHASH = keccak256("Deallocate(address staker,address agent,uint256 worknetId,uint256 amount,uint256 nonce,uint256 deadline)");
        nonce = allocator.nonces(signer);
        structHash = keccak256(abi.encode(DEALLOCATE_TYPEHASH, signer, signer, worknetId, amount, nonce, deadline));
        digest = keccak256(abi.encodePacked("\x19\x01", allocDomain, structHash));
        (v, r, s) = vm.sign(signerPk, digest);

        vm.prank(charlie);
        allocator.deallocateFor(signer, signer, worknetId, amount, deadline, v, r, s);
        assertEq(allocator.getAgentStake(signer, signer, worknetId), 0);
    }

    // ═══════════════════════════════════════════════
    //  23. Gasless setRecipientFor / grantDelegateFor
    // ═══════════════════════════════════════════════

    function test_gasless_setRecipientFor() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("recipientSigner");
        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = registry.nonces(signer);

        bytes32 TYPEHASH = keccak256("SetRecipient(address user,address recipient,uint256 nonce,uint256 deadline)");
        bytes32 structHash = keccak256(abi.encode(TYPEHASH, signer, charlie, nonce, deadline));
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", _registryDomain(), structHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        vm.prank(bob); // relayer
        registry.setRecipientFor(signer, charlie, deadline, v, r, s);
        assertEq(registry.recipient(signer), charlie);
    }

    function test_gasless_grantDelegateFor() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("delegateSigner");
        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = registry.nonces(signer);

        bytes32 TYPEHASH = keccak256("GrantDelegate(address user,address delegate,uint256 nonce,uint256 deadline)");
        bytes32 structHash = keccak256(abi.encode(TYPEHASH, signer, bob, nonce, deadline));
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", _registryDomain(), structHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        vm.prank(charlie);
        registry.grantDelegateFor(signer, bob, deadline, v, r, s);
        assertTrue(registry.delegates(signer, bob));
    }

    // ═══════════════════════════════════════════════
    //  24. DAO — vote with post-proposal mint reverts
    // ═══════════════════════════════════════════════

    function test_dao_minted_after_proposal_reverts() public {
        vm.startPrank(alice);
        awp.approve(address(veAwp), 5_000_000e18);
        uint256 tAlice = veAwp.deposit(5_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.warp(block.timestamp + 10);

        address[] memory targets = new address[](1);
        targets[0] = address(dao);
        uint256[] memory values = new uint256[](1);
        bytes[] memory calldatas = new bytes[](1);
        calldatas[0] = abi.encodeCall(AWPDAO.setQuorumPercent, (5));

        uint256[] memory tokens = new uint256[](1);
        tokens[0] = tAlice;
        vm.prank(alice);
        uint256 pid = dao.proposeWithTokens(targets, values, calldatas, "post-mint", tokens);

        vm.warp(block.timestamp + 1 days + 1);

        // Bob deposits AFTER proposal creation
        vm.startPrank(bob);
        awp.approve(address(veAwp), 1_000_000e18);
        uint256 tBob = veAwp.deposit(1_000_000e18, 54 weeks);
        vm.stopPrank();

        tokens[0] = tBob;
        vm.prank(bob);
        vm.expectRevert(AWPDAO.MintedAfterProposal.selector);
        dao.castVoteWithReasonAndParams(pid, 1, "", abi.encode(tokens));
    }

    // ═══════════════════════════════════════════════
    //  25. DAO — quorum verification (4% of staked AWP)
    // ═══════════════════════════════════════════════

    function test_dao_quorum_check() public {
        // Stake 10M total
        vm.startPrank(alice);
        awp.approve(address(veAwp), 5_000_000e18);
        uint256 tAlice = veAwp.deposit(5_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.startPrank(bob);
        awp.approve(address(veAwp), 5_000_000e18);
        uint256 tBob = veAwp.deposit(5_000_000e18, 54 weeks);
        vm.stopPrank();

        vm.warp(block.timestamp + 10);

        // Verify quorum is 4% of 10M = 400K AWP
        // quorum(0) uses live totalVotingPower as fallback
        uint256 q = dao.quorum(0);
        assertEq(q, veAwp.totalVotingPower() * 4 / 100);
        assertEq(q, 10_000_000e18 * 4 / 100); // 400K AWP
    }

    // ═══════════════════════════════════════════════
    //  26. Multi-epoch emission with decay tracking
    // ═══════════════════════════════════════════════

    function test_emission_multi_epoch_decay() public {
        address r1 = address(0xF001);
        uint256 decayFactor = emission.decayFactor();
        uint256 epochDuration = emission.epochDuration();

        // First: unpause by warping past pausedUntil and triggering _checkResume
        uint64 pu = emission.pausedUntil();
        if (pu > 0) vm.warp(uint256(pu) + 1);

        // Trigger _checkResume via a settle or Guardian call
        // Submit for current settledEpoch to trigger resume
        uint256 epoch = emission.settledEpoch();
        uint256[] memory packed = new uint256[](1);
        packed[0] = (uint256(100) << 160) | uint256(uint160(r1));

        vm.prank(GUARDIAN);
        emission.submitAllocations(packed, 100, epoch);
        emission.settleEpoch(100);

        // Now settled 1 epoch. Settle 2 more with decay tracking
        uint256 prevEmission = emission.currentDailyEmission();
        for (uint256 i = 0; i < 2; i++) {
            vm.warp(block.timestamp + epochDuration + 1);

            uint256 ep = emission.settledEpoch();
            vm.prank(GUARDIAN);
            emission.submitAllocations(packed, 100, ep);
            emission.settleEpoch(100);

            uint256 curEmission = emission.currentDailyEmission();
            assertEq(curEmission, prevEmission * decayFactor / 1_000_000);
            prevEmission = curEmission;
        }

        assertTrue(awp.balanceOf(r1) > 0);
    }

    // ═══════════════════════════════════════════════
    //  27. Registry — Guardian parameter management
    // ═══════════════════════════════════════════════

    function test_registry_guardian_params() public {
        uint256 oldPrice = registry.initialAlphaPrice();
        uint256 oldMint = registry.initialAlphaMint();

        vm.prank(GUARDIAN);
        registry.setInitialAlphaPrice(2e15);
        assertEq(registry.initialAlphaPrice(), 2e15);

        vm.prank(GUARDIAN);
        registry.setInitialAlphaMint(500_000_000e18);
        assertEq(registry.initialAlphaMint(), 500_000_000e18);

        // Non-guardian reverts
        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        registry.setInitialAlphaPrice(1e15);

        // Restore
        vm.startPrank(GUARDIAN);
        registry.setInitialAlphaPrice(oldPrice);
        registry.setInitialAlphaMint(oldMint);
        vm.stopPrank();
    }

    // ═══════════════════════════════════════════════
    //  28. AWPWorkNet — burn by owner
    // ═══════════════════════════════════════════════

    function test_worknet_nft_burn() public {
        uint256 wid = _activateWorknet(alice, "BurnTest", "BRN");
        assertEq(workNet.ownerOf(wid), alice);

        vm.prank(alice);
        workNet.burn(wid);

        vm.expectRevert();
        workNet.ownerOf(wid);
    }

    // ═══════════════════════════════════════════════
    //  29. Expired signature reverts
    // ═══════════════════════════════════════════════

    function test_gasless_expired_signature_reverts() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("expiredSigner");
        uint256 deadline = block.timestamp - 1; // expired

        bytes32 TYPEHASH = keccak256("Bind(address agent,address target,uint256 nonce,uint256 deadline)");
        bytes32 structHash = keccak256(abi.encode(TYPEHASH, signer, bob, 0, deadline));
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", _registryDomain(), structHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        vm.prank(charlie);
        vm.expectRevert(AWPRegistry.ExpiredSignature.selector);
        registry.bindFor(signer, bob, deadline, v, r, s);
    }

    function test_gasless_invalid_signature_reverts() public {
        (, uint256 wrongPk) = makeAddrAndKey("wrongKey");
        address signer = makeAddr("realSigner");
        uint256 deadline = block.timestamp + 1 hours;

        bytes32 TYPEHASH = keccak256("Bind(address agent,address target,uint256 nonce,uint256 deadline)");
        bytes32 structHash = keccak256(abi.encode(TYPEHASH, signer, bob, 0, deadline));
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", _registryDomain(), structHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(wrongPk, digest);

        vm.prank(charlie);
        vm.expectRevert(AWPRegistry.InvalidSignature.selector);
        registry.bindFor(signer, bob, deadline, v, r, s);
    }

    // ═══════════════════════════════════════════════
    //  30. Allocator batch operations
    // ═══════════════════════════════════════════════

    function test_allocator_batchAllocate() public {
        uint256 wid1 = _activateWorknet(alice, "Batch1", "B1");
        uint256 wid2 = _activateWorknet(bob, "Batch2", "B2");

        vm.startPrank(charlie);
        awp.approve(address(veAwp), 1_000_000e18);
        veAwp.deposit(1_000_000e18, 30 days);

        address[] memory agents = new address[](2);
        agents[0] = charlie;
        agents[1] = charlie;
        uint256[] memory worknetIds = new uint256[](2);
        worknetIds[0] = wid1;
        worknetIds[1] = wid2;
        uint256[] memory amounts = new uint256[](2);
        amounts[0] = 300_000e18;
        amounts[1] = 200_000e18;

        allocator.batchAllocate(charlie, agents, worknetIds, amounts);
        vm.stopPrank();

        assertEq(allocator.getAgentStake(charlie, charlie, wid1), 300_000e18);
        assertEq(allocator.getAgentStake(charlie, charlie, wid2), 200_000e18);
        assertEq(allocator.userTotalAllocated(charlie), 500_000e18);
    }

    // ═══════════════════════════════════════════════
    //  31. Emission → WorknetManager callback
    // ═══════════════════════════════════════════════

    function test_emission_to_worknetManager() public {
        uint256 wid = _activateWorknet(alice, "EmitWM", "EWM");
        IAWPRegistry.WorknetFullInfo memory full = registry.getWorknetFull(wid);
        address wmAddr = full.worknetManager;

        uint256 epoch = emission.settledEpoch();
        uint256 baseTime = emission.baseTime();
        uint256 epochDuration = emission.epochDuration();
        uint256 baseEpoch = emission.baseEpoch();
        uint256 targetTime = baseTime + (epoch + 1 - baseEpoch) * epochDuration + 1;
        if (block.timestamp < targetTime) vm.warp(targetTime);

        // Set WM as emission recipient
        uint256[] memory packed = new uint256[](1);
        packed[0] = (uint256(100) << 160) | uint256(uint160(wmAddr));

        vm.prank(GUARDIAN);
        emission.submitAllocations(packed, 100, epoch);

        uint256 wmBalBefore = awp.balanceOf(wmAddr);
        emission.settleEpoch(100);
        uint256 wmBalAfter = awp.balanceOf(wmAddr);

        // WM received AWP emission
        assertTrue(wmBalAfter > wmBalBefore);
    }

    // ═══════════════════════════════════════════════
    //  Helper
    // ═══════════════════════════════════════════════

    function _registryDomain() internal view returns (bytes32) {
        return keccak256(abi.encode(
            keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
            keccak256("AWPRegistry"), keccak256("1"), block.chainid, address(registry)
        ));
    }

    function _activateWorknet(address user, string memory name, string memory symbol) internal returns (uint256) {
        uint256 cost = registry.initialAlphaMint() * registry.initialAlphaPrice() / 1e18;
        vm.startPrank(user);
        awp.approve(address(registry), cost);
        uint256 wid = registry.registerWorknet(IAWPRegistry.WorknetParams({
            name: name, symbol: symbol,
            worknetManager: address(0), salt: bytes32(0), minStake: 0, skillsURI: ""
        }));
        vm.stopPrank();

        vm.prank(GUARDIAN);
        registry.activateWorknet(wid);
        return wid;
    }
}
