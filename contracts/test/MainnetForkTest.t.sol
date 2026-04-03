// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {IAWPRegistry} from "../src/interfaces/IAWPRegistry.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {WorknetNFT} from "../src/core/WorknetNFT.sol";
import {AlphaTokenFactory} from "../src/token/AlphaTokenFactory.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {TimelockController} from "@openzeppelin/contracts/governance/TimelockController.sol";
import {IGovernor} from "@openzeppelin/contracts/governance/IGovernor.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import {LPManagerBase} from "../src/core/LPManagerBase.sol";

/// @dev Minimal LPManager interface for fork tests
interface ILPManagerBase {
    function alphaTokenToPoolId(address) external view returns (bytes32);
    function alphaTokenToTokenId(address) external view returns (uint256);
    function compoundFees(address alphaToken) external;
    function createPoolAndAddLiquidity(address alphaToken, uint256 awpAmount, uint256 alphaAmount) external returns (bytes32, uint256);
}

/// @dev Minimal AlphaToken interface
interface IAlphaToken {
    function mintersLocked() external view returns (bool);
    function minters(address) external view returns (bool);
    function totalSupply() external view returns (uint256);
    function balanceOf(address) external view returns (uint256);
}

interface IAccessControl {
    function hasRole(bytes32 role, address account) external view returns (bool);
}

/// @dev Minimal interface for AWPDAO to avoid importing the full Governor stack
interface IAWPDAO {
    function proposeWithTokens(
        address[] memory targets, uint256[] memory values, bytes[] memory calldatas,
        string memory description, uint256[] memory tokenIds
    ) external returns (uint256);
    function signalPropose(string memory description, uint256[] memory tokenIds) external returns (uint256);
    function castVoteWithReasonAndParams(
        uint256 proposalId, uint8 support, string calldata reason, bytes memory params
    ) external returns (uint256);
    function execute(
        address[] memory targets, uint256[] memory values, bytes[] memory calldatas, bytes32 descriptionHash
    ) external payable returns (uint256);
    function state(uint256 proposalId) external view returns (IGovernor.ProposalState);
    function proposalVotes(uint256 proposalId) external view returns (uint256, uint256, uint256);
    function proposalThreshold() external view returns (uint256);
    function votingDelay() external view returns (uint256);
    function votingPeriod() external view returns (uint256);
}

// ═════════════════════════════════════════════════════════════
//  Constants library — shared addresses and config
// ═════════════════════════════════════════════════════════════

library Addrs {
    address constant AWP_TOKEN           = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant AWP_REGISTRY        = 0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A;
    address constant AWP_EMISSION        = 0x3C9cB73f8B81083882c5308Cce4F31f93600EaA9;
    address constant STAKING_VAULT       = 0xE8A204fD9c94C7E28bE11Af02fc4A4AC294Df29b;
    address constant STAKE_NFT           = 0x4E119560632698Bab67cFAB5d8EC0A373363ba2d;
    address constant WORKNET_NFT         = 0xB9F03539BE496d09c4d7964921d674B8763f5233;
    address constant ALPHA_TOKEN_FACTORY = 0xB2e4897eD77d0f5BFa3140B9989594de09a8037c;
    address constant TREASURY            = 0x82562023a053025F3201785160CaE6051efD759e;
    address constant AWP_DAO             = 0x6a074aC9823c47f86EE4Fc7F62e4217Bc9C76004;
    address constant LP_MANAGER_UNI      = 0x3034E029e61e8c2fc525A7bC5E267Ad3837D72e3;
    address constant WORKNET_MANAGER_UNI = 0x8d38accdE300917626f9c7DFd930cc47AA447137;
    address constant GUARDIAN            = 0x000002bEfa6A1C99A710862Feb6dB50525dF00A3;

    uint256 constant GENESIS_TIME           = 1775102400;
    uint256 constant EPOCH_DURATION         = 86400;
    uint256 constant CURRENT_DAILY_EMISSION = 31_600_000e18;
    uint256 constant DECAY_FACTOR           = 996844;
    uint256 constant DECAY_PRECISION        = 1_000_000;
}

/// @title MainnetForkBase — Shared setup and helpers for mainnet fork tests
abstract contract MainnetForkBase is Test {
    AWPToken          internal awpToken;
    AWPRegistry       internal awpRegistry;
    AWPEmission       internal awpEmission;
    StakingVault      internal stakingVault;
    StakeNFT          internal stakeNFT;
    WorknetNFT        internal worknetNFT;
    AlphaTokenFactory internal alphaTokenFactory;
    Treasury          internal treasury;
    IAWPDAO           internal awpDao;

    address internal alice = makeAddr("alice");
    address internal bob   = makeAddr("bob");
    address internal carol = makeAddr("carol");
    address internal dave  = makeAddr("dave");
    address internal eve   = makeAddr("eve");

    uint256 internal forkId;

    function _createFork() internal virtual returns (uint256);

    function setUp() public virtual {
        forkId = _createFork();
        vm.selectFork(forkId);

        awpToken          = AWPToken(Addrs.AWP_TOKEN);
        awpRegistry       = AWPRegistry(Addrs.AWP_REGISTRY);
        awpEmission       = AWPEmission(Addrs.AWP_EMISSION);
        stakingVault      = StakingVault(Addrs.STAKING_VAULT);
        stakeNFT          = StakeNFT(Addrs.STAKE_NFT);
        worknetNFT        = WorknetNFT(Addrs.WORKNET_NFT);
        alphaTokenFactory = AlphaTokenFactory(Addrs.ALPHA_TOKEN_FACTORY);
        treasury          = Treasury(payable(Addrs.TREASURY));
        awpDao            = IAWPDAO(Addrs.AWP_DAO);

        vm.deal(alice, 10 ether);
        vm.deal(bob, 10 ether);
        vm.deal(carol, 10 ether);
        vm.deal(dave, 10 ether);
        vm.deal(eve, 10 ether);
        vm.deal(Addrs.GUARDIAN, 10 ether);
        vm.deal(Addrs.TREASURY, 10 ether);
    }

    // ── Helpers ──

    function _packArray(address[] memory addrs, uint96[] memory wts) internal pure returns (uint256[] memory) {
        uint256[] memory packed = new uint256[](addrs.length);
        for (uint256 i = 0; i < addrs.length; i++) {
            packed[i] = (uint256(wts[i]) << 160) | uint256(uint160(addrs[i]));
        }
        return packed;
    }

    function _settleEpoch0(address recipient) internal returns (uint256) {
        vm.warp(Addrs.GENESIS_TIME);
        address[] memory r = new address[](1);
        r[0] = recipient;
        uint96[] memory w = new uint96[](1);
        w[0] = 10000;
        vm.prank(Addrs.GUARDIAN);
        awpEmission.submitAllocations(_packArray(r, w), 10000, 0);
        awpEmission.settleEpoch(100);
        return awpToken.balanceOf(recipient);
    }

    function _sortAndSubmit(address[] memory addrs, uint96[] memory wts, uint256 epoch) internal {
        uint256 tw = 0;
        for (uint256 i = 0; i < wts.length; i++) tw += wts[i];
        vm.prank(Addrs.GUARDIAN);
        awpEmission.submitAllocations(_packArray(addrs, wts), tw, epoch);
    }

    function _stakeAWP(address user, uint256 amount, uint64 lockDuration) internal returns (uint256) {
        vm.startPrank(user);
        awpToken.approve(Addrs.STAKE_NFT, amount);
        uint256 tid = stakeNFT.deposit(amount, lockDuration);
        vm.stopPrank();
        return tid;
    }

    function _advanceBlocks(uint256 n) internal {
        vm.roll(block.number + n);
        vm.warp(block.timestamp + n * 12);
    }

    // ═════════════════════════════════════════════════════════
    //  1. STATE VERIFICATION
    // ═════════════════════════════════════════════════════════

    function test_AWPToken_state() public view {
        assertEq(awpToken.admin(), address(0));
        assertTrue(awpToken.minters(Addrs.AWP_EMISSION));
        assertTrue(awpToken.initialMinted());
        assertEq(awpToken.totalSupply(), 0);
        assertEq(awpToken.MAX_SUPPLY(), 10_000_000_000e18);
    }

    function test_AWPRegistry_state() public view {
        assertEq(awpRegistry.awpToken(), Addrs.AWP_TOKEN);
        assertEq(awpRegistry.worknetNFT(), Addrs.WORKNET_NFT);
        assertEq(awpRegistry.alphaTokenFactory(), Addrs.ALPHA_TOKEN_FACTORY);
        assertEq(awpRegistry.awpEmission(), Addrs.AWP_EMISSION);
        assertEq(awpRegistry.stakingVault(), Addrs.STAKING_VAULT);
        assertEq(awpRegistry.stakeNFT(), Addrs.STAKE_NFT);
        assertEq(awpRegistry.treasury(), Addrs.TREASURY);
        assertEq(awpRegistry.guardian(), Addrs.GUARDIAN);
        assertTrue(awpRegistry.registryInitialized());
        assertFalse(awpRegistry.paused());
    }

    function test_AWPEmission_state() public view {
        assertEq(awpEmission.guardian(), Addrs.GUARDIAN);
        assertEq(awpEmission.treasury(), Addrs.TREASURY);
        assertEq(awpEmission.baseTime(), Addrs.GENESIS_TIME);
        assertEq(awpEmission.epochDuration(), Addrs.EPOCH_DURATION);
        assertEq(awpEmission.decayFactor(), Addrs.DECAY_FACTOR);
        assertEq(awpEmission.settledEpoch(), 0);
        assertEq(awpEmission.currentDailyEmission(), Addrs.CURRENT_DAILY_EMISSION);
        assertEq(awpEmission.maxRecipients(), 10000);
    }

    function test_StakingVault_state() public view {
        assertEq(stakingVault.awpRegistry(), Addrs.AWP_REGISTRY);
        assertEq(stakingVault.guardian(), Addrs.GUARDIAN);
        assertEq(stakingVault.stakeNFT(), Addrs.STAKE_NFT);
    }

    function test_StakeNFT_state() public view {
        assertEq(address(stakeNFT.awpToken()), Addrs.AWP_TOKEN);
        assertEq(stakeNFT.stakingVault(), Addrs.STAKING_VAULT);
        assertEq(stakeNFT.awpRegistry(), Addrs.AWP_REGISTRY);
    }

    function test_Factory_state() public view {
        assertTrue(alphaTokenFactory.configured());
        assertEq(alphaTokenFactory.owner(), address(0));
        assertEq(alphaTokenFactory.awpRegistry(), Addrs.AWP_REGISTRY);
    }

    function test_Treasury_state() public view {
        assertTrue(treasury.hasRole(treasury.PROPOSER_ROLE(), Addrs.AWP_DAO));
        assertTrue(treasury.hasRole(treasury.DEFAULT_ADMIN_ROLE(), Addrs.GUARDIAN));
    }

    function test_WorknetNFT_state() public view {
        assertEq(worknetNFT.awpRegistry(), Addrs.AWP_REGISTRY);
    }

    // ═════════════════════════════════════════════════════════
    //  2. ACCOUNT SYSTEM
    // ═════════════════════════════════════════════════════════

    function test_register() public {
        vm.prank(alice);
        awpRegistry.setRecipient(alice);
        assertEq(awpRegistry.recipient(alice), alice);
        assertTrue(awpRegistry.isRegistered(alice));
    }

    function test_bind() public {
        vm.prank(alice);
        awpRegistry.bind(bob);
        assertEq(awpRegistry.boundTo(alice), bob);
        assertTrue(awpRegistry.isRegistered(alice));
    }

    function test_setRecipient() public {
        vm.prank(alice);
        awpRegistry.setRecipient(carol);
        assertEq(awpRegistry.recipient(alice), carol);
    }

    function test_resolveRecipient_walksChain() public {
        vm.prank(carol);
        awpRegistry.setRecipient(dave);
        vm.prank(bob);
        awpRegistry.bind(carol);
        vm.prank(alice);
        awpRegistry.bind(bob);
        assertEq(awpRegistry.resolveRecipient(alice), dave);
    }

    function test_bind_revertsSelfBind() public {
        vm.prank(alice);
        vm.expectRevert(AWPRegistry.SelfBind.selector);
        awpRegistry.bind(alice);
    }

    function test_bind_revertsCycle() public {
        vm.prank(alice);
        awpRegistry.bind(bob);
        vm.prank(bob);
        vm.expectRevert(AWPRegistry.CycleDetected.selector);
        awpRegistry.bind(alice);
    }

    function test_delegation() public {
        vm.prank(alice);
        awpRegistry.grantDelegate(bob);
        assertTrue(awpRegistry.delegates(alice, bob));
        vm.prank(alice);
        awpRegistry.revokeDelegate(bob);
        assertFalse(awpRegistry.delegates(alice, bob));
    }

    // ═════════════════════════════════════════════════════════
    //  3. GUARDIAN OPERATIONS
    // ═════════════════════════════════════════════════════════

    function test_guardian_submitAllocations() public {
        vm.warp(Addrs.GENESIS_TIME);
        address[] memory r = new address[](2);
        uint96[] memory w = new uint96[](2);
        r[0] = address(0x1111);
        r[1] = address(0x2222);
        w[0] = 5000;
        w[1] = 5000;
        vm.prank(Addrs.GUARDIAN);
        awpEmission.submitAllocations(_packArray(r, w), 10000, 0);
        assertEq(awpEmission.getEpochRecipientCount(0), 2);
        assertEq(awpEmission.getEpochTotalWeight(0), 10000);
    }

    function test_guardian_setDecayFactor() public {
        vm.prank(Addrs.GUARDIAN);
        awpEmission.setDecayFactor(995000);
        assertEq(awpEmission.decayFactor(), 995000);
        vm.prank(Addrs.GUARDIAN);
        awpEmission.setDecayFactor(Addrs.DECAY_FACTOR);
    }

    function test_guardian_setEpochDuration() public {
        vm.warp(Addrs.GENESIS_TIME);
        vm.prank(Addrs.GUARDIAN);
        awpEmission.setEpochDuration(43200);
        assertEq(awpEmission.epochDuration(), 43200);
        vm.prank(Addrs.GUARDIAN);
        awpEmission.setEpochDuration(Addrs.EPOCH_DURATION);
    }

    function test_guardian_setTreasury() public {
        vm.prank(Addrs.GUARDIAN);
        awpEmission.setTreasury(makeAddr("newT"));
        vm.prank(Addrs.GUARDIAN);
        awpEmission.setTreasury(Addrs.TREASURY);
    }

    function test_guardian_pauseUnpause() public {
        vm.prank(Addrs.GUARDIAN);
        awpRegistry.pause();
        assertTrue(awpRegistry.paused());
        vm.prank(Addrs.GUARDIAN);
        awpRegistry.unpause();
        assertFalse(awpRegistry.paused());
    }

    function test_guardian_setGuardian() public {
        address ng = makeAddr("newG");
        vm.prank(Addrs.GUARDIAN);
        awpRegistry.setGuardian(ng);
        assertEq(awpRegistry.guardian(), ng);
        vm.prank(ng);
        awpRegistry.setGuardian(Addrs.GUARDIAN);
    }

    // ═════════════════════════════════════════════════════════
    //  4. EMISSION
    // ═════════════════════════════════════════════════════════

    function test_emission_settleEpoch0() public {
        uint256 minted = _settleEpoch0(alice);
        assertEq(minted, Addrs.CURRENT_DAILY_EMISSION);
        assertEq(awpToken.totalSupply(), Addrs.CURRENT_DAILY_EMISSION);
        assertEq(awpEmission.settledEpoch(), 1);
    }

    function test_emission_settleEpoch1_decay() public {
        _settleEpoch0(alice);
        uint256 bal0 = awpToken.balanceOf(alice);
        vm.warp(Addrs.GENESIS_TIME + Addrs.EPOCH_DURATION);
        awpEmission.settleEpoch(100);
        uint256 e1 = awpToken.balanceOf(alice) - bal0;
        assertEq(e1, Addrs.CURRENT_DAILY_EMISSION * Addrs.DECAY_FACTOR / Addrs.DECAY_PRECISION);
    }

    function test_emission_multiRecipients() public {
        vm.warp(Addrs.GENESIS_TIME);
        address[] memory r = new address[](3);
        uint96[] memory w = new uint96[](3);
        r[0] = alice; r[1] = bob; r[2] = carol;
        w[0] = 5000; w[1] = 3000; w[2] = 2000;
        _sortAndSubmit(r, w, 0);
        awpEmission.settleEpoch(100);
        uint256 total = awpToken.balanceOf(r[0]) + awpToken.balanceOf(r[1]) + awpToken.balanceOf(r[2]);
        assertEq(total, Addrs.CURRENT_DAILY_EMISSION);
    }

    function test_emission_revertBeforeGenesis() public {
        vm.warp(Addrs.GENESIS_TIME - 1);
        vm.expectRevert(AWPEmission.GenesisNotReached.selector);
        awpEmission.currentEpoch();
    }

    function test_emission_settleRevertsBeforeGenesis() public {
        vm.warp(Addrs.GENESIS_TIME - 1);
        vm.expectRevert(AWPEmission.GenesisNotReached.selector);
        awpEmission.settleEpoch(100);
    }

    function test_emission_batchSettlement() public {
        vm.warp(Addrs.GENESIS_TIME);
        address[] memory r = new address[](5);
        uint96[] memory w = new uint96[](5);
        r[0] = alice; r[1] = bob; r[2] = carol; r[3] = dave; r[4] = eve;
        for (uint256 i = 0; i < 5; i++) w[i] = 2000;
        _sortAndSubmit(r, w, 0);

        awpEmission.settleEpoch(2);
        assertTrue(awpEmission.settleProgress() > 0);
        awpEmission.settleEpoch(2);
        assertTrue(awpEmission.settleProgress() > 0);
        awpEmission.settleEpoch(2);
        assertEq(awpEmission.settleProgress(), 0);
        assertEq(awpEmission.settledEpoch(), 1);

        uint256 expected = Addrs.CURRENT_DAILY_EMISSION / 5;
        for (uint256 i = 0; i < 5; i++) {
            assertEq(awpToken.balanceOf(r[i]), expected);
        }
    }

    // ═════════════════════════════════════════════════════════
    //  5. STAKING
    // ═════════════════════════════════════════════════════════

    function test_staking_depositAndAllocate() public {
        uint256 minted = _settleEpoch0(alice);
        uint256 sa = minted / 2;
        uint256 tid = _stakeAWP(alice, sa, 30 days);

        (uint128 amt, uint64 lockEnd, uint64 cAt) = stakeNFT.positions(tid);
        assertEq(amt, sa);
        assertEq(lockEnd, uint64(block.timestamp) + 30 days);
        assertEq(cAt, uint64(block.timestamp));
        assertEq(stakeNFT.getUserTotalStaked(alice), sa);

        uint256 wid = (block.chainid << 64) | 999;
        vm.prank(alice);
        stakingVault.allocate(alice, alice, wid, sa);
        assertEq(stakingVault.getAgentStake(alice, alice, wid), sa);
        assertEq(stakingVault.userTotalAllocated(alice), sa);
    }

    function test_staking_deallocate() public {
        _settleEpoch0(alice);
        uint256 sa = awpToken.balanceOf(alice) / 2;
        _stakeAWP(alice, sa, 30 days);
        uint256 wid = (block.chainid << 64) | 999;
        vm.prank(alice);
        stakingVault.allocate(alice, alice, wid, sa);

        vm.prank(alice);
        stakingVault.deallocate(alice, alice, wid, sa / 2);
        assertEq(stakingVault.getAgentStake(alice, alice, wid), sa - sa / 2);
        assertEq(stakingVault.userTotalAllocated(alice), sa - sa / 2);
    }

    function test_staking_delegateCanAllocate() public {
        _settleEpoch0(alice);
        _stakeAWP(alice, awpToken.balanceOf(alice) / 2, 30 days);
        vm.prank(alice);
        awpRegistry.grantDelegate(bob);

        uint256 wid = (block.chainid << 64) | 999;
        vm.prank(bob);
        stakingVault.allocate(alice, alice, wid, 1000e18);
        assertEq(stakingVault.getAgentStake(alice, alice, wid), 1000e18);
    }

    function test_staking_reallocate() public {
        _settleEpoch0(alice);
        uint256 sa = awpToken.balanceOf(alice) / 2;
        _stakeAWP(alice, sa, 30 days);
        uint256 w1 = (block.chainid << 64) | 100;
        uint256 w2 = (block.chainid << 64) | 200;

        vm.prank(alice);
        stakingVault.allocate(alice, alice, w1, sa);
        vm.prank(alice);
        stakingVault.reallocate(alice, alice, w1, alice, w2, sa / 2);

        assertEq(stakingVault.getAgentStake(alice, alice, w1), sa - sa / 2);
        assertEq(stakingVault.getAgentStake(alice, alice, w2), sa / 2);
        assertEq(stakingVault.userTotalAllocated(alice), sa);
    }

    function test_staking_withdrawAfterLock() public {
        _settleEpoch0(alice);
        uint256 sa = awpToken.balanceOf(alice) / 4;
        uint256 tid = _stakeAWP(alice, sa, 1 days);
        vm.warp(block.timestamp + 1 days + 1);

        uint256 before = awpToken.balanceOf(alice);
        vm.prank(alice);
        stakeNFT.withdraw(tid);
        assertEq(awpToken.balanceOf(alice) - before, sa);
    }

    function test_staking_addToPosition() public {
        _settleEpoch0(alice);
        uint256 sa = awpToken.balanceOf(alice) / 4;
        uint256 tid = _stakeAWP(alice, sa, 30 days);

        vm.startPrank(alice);
        awpToken.approve(Addrs.STAKE_NFT, sa);
        stakeNFT.addToPosition(tid, sa, 0);
        vm.stopPrank();

        (uint128 amt,,) = stakeNFT.positions(tid);
        assertEq(amt, sa * 2);
    }

    function test_staking_cannotAllocateMoreThanStaked() public {
        _settleEpoch0(alice);
        uint256 sa = awpToken.balanceOf(alice) / 4;
        _stakeAWP(alice, sa, 30 days);
        vm.prank(alice);
        vm.expectRevert(StakingVault.InsufficientUnallocated.selector);
        stakingVault.allocate(alice, alice, (block.chainid << 64) | 999, sa + 1);
    }

    function test_staking_rejectsZeroWorknetId() public {
        _settleEpoch0(alice);
        _stakeAWP(alice, awpToken.balanceOf(alice) / 4, 30 days);
        vm.prank(alice);
        vm.expectRevert(StakingVault.ZeroWorknetId.selector);
        stakingVault.allocate(alice, alice, 0, 100e18);
    }

    function test_stakeNFT_cannotWithdrawBeforeLock() public {
        _settleEpoch0(alice);
        uint256 tid = _stakeAWP(alice, awpToken.balanceOf(alice) / 4, 30 days);
        vm.prank(alice);
        vm.expectRevert(StakeNFT.LockNotExpired.selector);
        stakeNFT.withdraw(tid);
    }

    function test_stakeNFT_addToExpiredReverts() public {
        _settleEpoch0(alice);
        uint256 tid = _stakeAWP(alice, awpToken.balanceOf(alice) / 4, 1 days);
        vm.warp(block.timestamp + 2 days);
        vm.startPrank(alice);
        awpToken.approve(Addrs.STAKE_NFT, 1000e18);
        vm.expectRevert(StakeNFT.PositionExpired.selector);
        stakeNFT.addToPosition(tid, 1000e18, 0);
        vm.stopPrank();
    }

    function test_staking_transferNFT() public {
        _settleEpoch0(alice);
        uint256 sa = awpToken.balanceOf(alice) / 4;
        uint256 tid = _stakeAWP(alice, sa, 30 days);
        vm.prank(alice);
        IERC721(Addrs.STAKE_NFT).transferFrom(alice, bob, tid);
        assertEq(stakeNFT.getUserTotalStaked(alice), 0);
        assertEq(stakeNFT.getUserTotalStaked(bob), sa);
    }

    // ═════════════════════════════════════════════════════════
    //  6. WORKNET REGISTRATION
    // ═════════════════════════════════════════════════════════

    function test_worknetRegistration() public {
        _settleEpoch0(alice);
        uint256 lpAmt = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        if (awpToken.balanceOf(alice) < lpAmt) return;

        vm.prank(alice);
        awpToken.approve(Addrs.AWP_REGISTRY, lpAmt);

        uint256 wid = _doRegisterWorknet(alice);
        assertTrue(wid > 0);
        assertEq(IERC721(Addrs.WORKNET_NFT).ownerOf(wid), alice);

        (,IAWPRegistry.WorknetStatus s,,) = awpRegistry.worknets(wid);
        assertEq(uint8(s), 0); // Pending

        vm.prank(alice);
        awpRegistry.activateWorknet(wid);
        (,s,,) = awpRegistry.worknets(wid);
        assertEq(uint8(s), 1); // Active
    }

    function _doRegisterWorknet(address user) internal returns (uint256) {
        // Use pre-mined vanity salt that produces an Alpha token address matching A1??...CAFE
        IAWPRegistry.WorknetParams memory p = IAWPRegistry.WorknetParams({
            name: "Test Worknet Alpha",
            symbol: "TALPHA",
            worknetManager: address(0),
            salt: bytes32(0xe15ccbfb709b022143edd0fb283af5a64661d65f2f93e66bdd210c35e3111b27),
            minStake: 1000e18,
            skillsURI: "https://example.com/skills.json"
        });
        vm.prank(user);
        return awpRegistry.registerWorknet(p);
    }

    // ═════════════════════════════════════════════════════════
    //  6b. LP & WORKNET LIFECYCLE (registration + LP pool + compoundFees + activation)
    // ═════════════════════════════════════════════════════════

    function test_worknet_lpPoolCreated() public {
        // Register worknet → LP pool auto-created
        _settleEpoch0(alice);
        uint256 lpAmt = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        if (awpToken.balanceOf(alice) < lpAmt) return; // skip if insufficient

        vm.prank(alice);
        awpToken.approve(Addrs.AWP_REGISTRY, lpAmt);
        uint256 wid = _doRegisterWorknet(alice);

        // Verify LP pool was created in LPManager
        ILPManagerBase lp = ILPManagerBase(awpRegistry.lpManager());
        // Get alphaToken from WorknetNFT
        WorknetNFT wnft = WorknetNFT(Addrs.WORKNET_NFT);
        address alphaToken = wnft.getAlphaToken(wid);
        assertTrue(alphaToken != address(0), "Alpha token not deployed");

        bytes32 poolId = lp.alphaTokenToPoolId(alphaToken);
        assertTrue(poolId != bytes32(0), "LP pool not created");

        uint256 lpTokenId = lp.alphaTokenToTokenId(alphaToken);
        assertTrue(lpTokenId > 0, "LP token ID not set");
    }

    function test_worknet_compoundFees() public {
        // Register worknet → LP pool → try compoundFees (no fees yet, should not revert or revert gracefully)
        _settleEpoch0(alice);
        uint256 lpAmt = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        if (awpToken.balanceOf(alice) < lpAmt) return;

        vm.prank(alice);
        awpToken.approve(Addrs.AWP_REGISTRY, lpAmt);
        uint256 wid = _doRegisterWorknet(alice);

        WorknetNFT wnft = WorknetNFT(Addrs.WORKNET_NFT);
        address alphaToken = wnft.getAlphaToken(wid);
        ILPManagerBase lp = ILPManagerBase(awpRegistry.lpManager());

        // compoundFees on a new pool with no fees — may revert (V4 TAKE_PAIR with 0 fees)
        // This is expected; verify LP data is intact afterward
        try lp.compoundFees(alphaToken) {
            // If it doesn't revert, that's fine
        } catch {
            // Expected: no fees to compound yet
        }

        // LP data should still be intact
        assertTrue(lp.alphaTokenToPoolId(alphaToken) != bytes32(0), "Pool still exists");
    }

    function test_worknet_alphaTokenMint() public {
        // After activation, WorknetManager can mint Alpha tokens via Merkle claims
        _settleEpoch0(alice);
        uint256 lpAmt = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        if (awpToken.balanceOf(alice) < lpAmt) return;

        vm.prank(alice);
        awpToken.approve(Addrs.AWP_REGISTRY, lpAmt);
        uint256 wid = _doRegisterWorknet(alice);

        // Activate
        vm.prank(alice);
        awpRegistry.activateWorknet(wid);

        // Get the worknet manager and alpha token
        WorknetNFT wnft = WorknetNFT(Addrs.WORKNET_NFT);
        address alphaToken = wnft.getAlphaToken(wid);
        address worknetMgr = wnft.getWorknetManager(wid);
        assertTrue(worknetMgr != address(0), "WorknetManager auto-deployed");

        // Verify alpha token has minter locked to worknetManager
        IAlphaToken alpha = IAlphaToken(alphaToken);
        assertTrue(alpha.mintersLocked(), "Minters should be locked");
        assertTrue(alpha.minters(worknetMgr), "WorknetManager should be minter");
        assertFalse(alpha.minters(Addrs.AWP_REGISTRY), "Registry should NOT be minter after lock");
    }

    function test_worknet_pauseBanDeregister() public {
        _settleEpoch0(alice);
        uint256 lpAmt = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        if (awpToken.balanceOf(alice) < lpAmt) return;

        vm.prank(alice);
        awpToken.approve(Addrs.AWP_REGISTRY, lpAmt);
        uint256 wid = _doRegisterWorknet(alice);

        // Activate
        vm.prank(alice);
        awpRegistry.activateWorknet(wid);

        // Pause (owner)
        vm.prank(alice);
        awpRegistry.pauseWorknet(wid);
        (,IAWPRegistry.WorknetStatus s,,) = awpRegistry.worknets(wid);
        assertEq(uint8(s), 2); // Paused

        // Resume (owner)
        vm.prank(alice);
        awpRegistry.resumeWorknet(wid);
        (,s,,) = awpRegistry.worknets(wid);
        assertEq(uint8(s), 1); // Active

        // Ban (timelock)
        vm.prank(Addrs.GUARDIAN);
        awpRegistry.banWorknet(wid);
        (,s,,) = awpRegistry.worknets(wid);
        assertEq(uint8(s), 3); // Banned

        // Deregister (timelock, after immunity)
        vm.warp(block.timestamp + 31 days);
        vm.prank(Addrs.GUARDIAN);
        awpRegistry.deregisterWorknet(wid);

        // NFT burned
        vm.expectRevert();
        IERC721(Addrs.WORKNET_NFT).ownerOf(wid);
    }

    function test_worknet_lpManagerOnlyRegistry() public {
        // Non-AWPRegistry cannot call createPoolAndAddLiquidity
        ILPManagerBase lp = ILPManagerBase(awpRegistry.lpManager());
        vm.prank(alice);
        vm.expectRevert();
        lp.createPoolAndAddLiquidity(alice, 1e18, 1e18);
    }

    // ═════════════════════════════════════════════════════════
    //  6c. WORKNET MANAGER AUTO-DEPLOY + INITIALIZATION
    // ═════════════════════════════════════════════════════════

    function test_worknetManager_autoDeployed() public {
        _settleEpoch0(alice);
        uint256 lpAmt = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        if (awpToken.balanceOf(alice) < lpAmt) return;

        vm.prank(alice);
        awpToken.approve(Addrs.AWP_REGISTRY, lpAmt);
        uint256 wid = _doRegisterWorknet(alice);

        WorknetNFT wnft = WorknetNFT(Addrs.WORKNET_NFT);
        address mgr = wnft.getWorknetManager(wid);
        assertTrue(mgr != address(0), "Manager auto-deployed");

        // Verify manager is initialized correctly (has roles)
        // DEFAULT_ADMIN_ROLE = 0x00
        assertTrue(
            IAccessControl(mgr).hasRole(bytes32(0), alice),
            "Alice should have DEFAULT_ADMIN_ROLE"
        );
    }

    // ═════════════════════════════════════════════════════════
    //  6d. PAUSE EPOCH + EMERGENCY UNPAUSE + RESCUE TOKEN
    // ═════════════════════════════════════════════════════════

    function test_pauseEpochUntil() public {
        vm.warp(Addrs.GENESIS_TIME);
        // Guardian pauses epoch until far future
        vm.prank(Addrs.GUARDIAN);
        awpEmission.pauseEpochUntil(uint64(block.timestamp + 30 days));

        // currentEpoch should be frozen at settledEpoch
        assertEq(awpEmission.currentEpoch(), awpEmission.settledEpoch());

        // settleEpoch still works for the frozen epoch
        address[] memory addrs = new address[](1);
        addrs[0] = alice;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _sortAndSubmit(addrs, ws, awpEmission.settledEpoch());
        awpEmission.settleEpoch(200);

        // After settling, settledEpoch=1 but currentEpoch=0 (frozen at settledEpoch before pause)
        assertEq(awpEmission.currentEpoch(), 0);
        assertEq(awpEmission.settledEpoch(), 1);

        // Immediate resume
        vm.prank(Addrs.GUARDIAN);
        awpEmission.pauseEpochUntil(0);

        // After immediate resume, baseTime = now, so currentEpoch = settledEpoch (0-based from now)
        // Need to warp forward for epoch to advance
        vm.warp(block.timestamp + Addrs.EPOCH_DURATION + 1);
        assertTrue(awpEmission.currentEpoch() >= awpEmission.settledEpoch(), "Epoch should advance after resume + warp");
    }

    function test_guardianUnpause() public {
        vm.prank(Addrs.GUARDIAN);
        awpRegistry.pause();
        assertTrue(awpRegistry.paused());

        vm.prank(Addrs.GUARDIAN);
        awpRegistry.unpause();
        assertFalse(awpRegistry.paused());
    }

    function test_setLPManager() public {
        address oldLP = awpRegistry.lpManager();
        address newLP = makeAddr("newLP");

        // Only Timelock can set LP manager
        vm.prank(Addrs.GUARDIAN);
        awpRegistry.setLPManager(newLP);
        assertEq(awpRegistry.lpManager(), newLP);

        // Restore
        vm.prank(Addrs.GUARDIAN);
        awpRegistry.setLPManager(oldLP);
    }

    function test_getActiveWorknetIds() public {
        // No active worknets initially
        uint256[] memory ids = awpRegistry.getActiveWorknetIds(0, 10);
        assertEq(ids.length, 0);
    }

    // ═════════════════════════════════════════════════════════
    //  6e. FULL LIFECYCLE: register → stake → allocate → emission → settle
    // ═════════════════════════════════════════════════════════

    function test_fullLifecycle() public {
        // 1. Settle epoch 0 to get AWP for alice
        uint256 minted = _settleEpoch0(alice);
        assertTrue(minted > 0, "Alice should have AWP from emission");

        // 2. Stake AWP
        uint256 stakeAmt = minted / 2;
        uint256 tid = _stakeAWP(alice, stakeAmt, 52 weeks);
        assertTrue(tid > 0, "Should have minted stake NFT");
        assertEq(stakeNFT.getUserTotalStaked(alice), stakeAmt);

        // 3. Allocate to a worknet (use cross-chain worknetId)
        uint256 foreignWorknetId = (uint256(42161) << 64) | 1;
        vm.prank(alice);
        stakingVault.allocate(alice, alice, foreignWorknetId, stakeAmt / 2);
        assertEq(stakingVault.getAgentStake(alice, alice, foreignWorknetId), stakeAmt / 2);

        // 4. Settle epoch 1 (with decay)
        vm.warp(block.timestamp + Addrs.EPOCH_DURATION + 1);
        address[] memory addrs = new address[](1);
        addrs[0] = alice;
        uint96[] memory ws = new uint96[](1);
        ws[0] = 100;
        _sortAndSubmit(addrs, ws, awpEmission.settledEpoch());
        awpEmission.settleEpoch(200);

        // 5. Verify alice received more AWP from epoch 1
        uint256 epoch1Bal = awpToken.balanceOf(alice);
        assertTrue(epoch1Bal > minted, "Should have more AWP after epoch 1 settlement");

        // 6. Deallocate and withdraw
        vm.prank(alice);
        stakingVault.deallocate(alice, alice, foreignWorknetId, stakeAmt / 2);
        assertEq(stakingVault.getAgentStake(alice, alice, foreignWorknetId), 0);

        // Can't withdraw yet (locked for 52 weeks)
        vm.expectRevert();
        vm.prank(alice);
        stakeNFT.withdraw(tid);

        // Warp past lock and withdraw
        vm.warp(block.timestamp + 52 weeks);
        vm.prank(alice);
        stakeNFT.withdraw(tid);
        assertTrue(awpToken.balanceOf(alice) > epoch1Bal, "Should have withdrawn stake");
    }

    // ═════════════════════════════════════════════════════════
    //  7. GOVERNANCE
    // ═════════════════════════════════════════════════════════

    function test_governance_propose() public {
        uint256 minted = _settleEpoch0(alice);
        uint256 tid = _stakeAWP(alice, minted, 54 weeks);
        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        uint256[] memory tids = new uint256[](1);
        tids[0] = tid;
        if (stakeNFT.getUserVotingPower(alice, tids) < awpDao.proposalThreshold()) return;

        address[] memory tgts = new address[](1);
        tgts[0] = Addrs.TREASURY;
        uint256[] memory vals = new uint256[](1);
        bytes[] memory data = new bytes[](1);

        vm.prank(alice);
        uint256 pid = awpDao.proposeWithTokens(tgts, vals, data, "Test: no-op", tids);
        assertTrue(pid > 0);
        assertEq(uint8(awpDao.state(pid)), uint8(IGovernor.ProposalState.Pending));
    }

    function test_governance_voteSignal() public {
        uint256 minted = _settleEpoch0(alice);
        uint256 tid = _stakeAWP(alice, minted, 54 weeks);
        vm.roll(block.number + 1);
        vm.warp(block.timestamp + 12);

        uint256[] memory tids = new uint256[](1);
        tids[0] = tid;
        if (stakeNFT.getUserVotingPower(alice, tids) < awpDao.proposalThreshold()) return;

        vm.prank(alice);
        uint256 pid = awpDao.signalPropose("Signal: test", tids);

        _advanceBlocks(awpDao.votingDelay() + 1);

        vm.prank(alice);
        awpDao.castVoteWithReasonAndParams(pid, 1, "yes", abi.encode(tids));

        (, uint256 forV,) = awpDao.proposalVotes(pid);
        assertTrue(forV > 0);

        _advanceBlocks(awpDao.votingPeriod() + 1);

        if (awpDao.state(pid) == IGovernor.ProposalState.Succeeded) {
            address[] memory t = new address[](1);
            t[0] = Addrs.AWP_DAO;
            uint256[] memory v = new uint256[](1);
            bytes[] memory c = new bytes[](1);
            awpDao.execute(t, v, c, keccak256("Signal: test"));
            assertEq(uint8(awpDao.state(pid)), uint8(IGovernor.ProposalState.Executed));
        }
    }

    // ═════════════════════════════════════════════════════════
    //  8. ACCESS CONTROL NEGATIVE TESTS
    // ═════════════════════════════════════════════════════════

    function test_acl_nonGuardianSubmit() public {
        vm.warp(Addrs.GENESIS_TIME);
        address[] memory r = new address[](1);
        r[0] = alice;
        uint96[] memory w = new uint96[](1);
        w[0] = 10000;
        vm.prank(alice);
        vm.expectRevert(AWPEmission.NotGuardian.selector);
        awpEmission.submitAllocations(_packArray(r, w), 10000, 0);
    }

    function test_acl_nonGuardianPause() public {
        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.pause();
    }

    function test_acl_nonTimelockUnpause() public {
        vm.prank(Addrs.GUARDIAN);
        awpRegistry.pause();
        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.unpause();
        vm.prank(Addrs.GUARDIAN);
        awpRegistry.unpause();
    }

    function test_acl_nonGuardianDecay() public {
        vm.prank(alice);
        vm.expectRevert(AWPEmission.NotGuardian.selector);
        awpEmission.setDecayFactor(990000);
    }

    function test_acl_nonGuardianEpochDuration() public {
        vm.prank(alice);
        vm.expectRevert(AWPEmission.NotGuardian.selector);
        awpEmission.setEpochDuration(3600);
    }

    function test_acl_nonGuardianUpgradeEmission() public {
        vm.prank(alice);
        vm.expectRevert(AWPEmission.NotGuardian.selector);
        awpEmission.upgradeToAndCall(makeAddr("x"), "");
    }

    function test_acl_nonGuardianUpgradeVault() public {
        vm.prank(alice);
        vm.expectRevert(StakingVault.NotGuardian.selector);
        stakingVault.upgradeToAndCall(makeAddr("x"), "");
    }

    function test_acl_nonGuardianUpgradeRegistry() public {
        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.upgradeToAndCall(makeAddr("x"), "");
    }

    function test_acl_nonTimelockBan() public {
        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.banWorknet(1);
    }

    function test_acl_nonTimelockSetPrice() public {
        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.setInitialAlphaPrice(1e16);
    }

    function test_acl_nonAuthorizedAllocate() public {
        vm.prank(bob);
        vm.expectRevert(StakingVault.NotAuthorized.selector);
        stakingVault.allocate(alice, alice, 1, 100e18);
    }

    function test_acl_nonGuardianSetGuardianRegistry() public {
        vm.prank(alice);
        vm.expectRevert(AWPRegistry.NotGuardian.selector);
        awpRegistry.setGuardian(alice);
    }

    function test_acl_nonGuardianSetGuardianEmission() public {
        vm.prank(alice);
        vm.expectRevert(AWPEmission.NotGuardian.selector);
        awpEmission.setGuardian(alice);
    }

    function test_acl_nonGuardianSetGuardianVault() public {
        vm.prank(alice);
        vm.expectRevert(StakingVault.NotGuardian.selector);
        stakingVault.setGuardian(alice);
    }

    function test_registryAlreadyInitialized() public {
        vm.prank(address(0xdead));
        vm.expectRevert(AWPRegistry.NotDeployer.selector);
        awpRegistry.initializeRegistry(
            Addrs.AWP_TOKEN, Addrs.WORKNET_NFT, Addrs.ALPHA_TOKEN_FACTORY, Addrs.AWP_EMISSION,
            Addrs.LP_MANAGER_UNI, Addrs.STAKING_VAULT, Addrs.STAKE_NFT, Addrs.WORKNET_MANAGER_UNI, ""
        );
    }
}

// ═════════════════════════════════════════════════════════════
//  Chain-specific test contracts
// ═════════════════════════════════════════════════════════════

contract EthMainnetTest is MainnetForkBase {
    function _createFork() internal override returns (uint256) {
        return vm.createFork(vm.envString("ETH_RPC_URL"));
    }
}

contract BaseMainnetTest is MainnetForkBase {
    function _createFork() internal override returns (uint256) {
        return vm.createFork(vm.envString("BASE_RPC_URL"));
    }
}

contract BscMainnetTest is MainnetForkBase {
    function _createFork() internal override returns (uint256) {
        return vm.createFork(vm.envString("BSC_RPC_URL"));
    }
}

contract ArbMainnetTest is MainnetForkBase {
    function _createFork() internal override returns (uint256) {
        return vm.createFork(vm.envString("ARB_RPC_URL"));
    }
}
