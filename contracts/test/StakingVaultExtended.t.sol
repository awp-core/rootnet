// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, Vm} from "forge-std/Test.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

/// @title StakingVaultExtended — Extended test coverage for StakingVault edge cases
contract StakingVaultExtendedTest is Test {
    StakingVault public vault;
    StakeNFT public stakeNFT;
    AWPToken public awp;

    address public deployer = address(this);
    address public user1 = makeAddr("user1");
    address public user2 = makeAddr("user2");
    address public user3 = makeAddr("user3");
    address public agent1 = makeAddr("agent1");
    address public agent2 = makeAddr("agent2");

    uint256 public constant WORKNET_1 = 1;
    uint256 public constant WORKNET_2 = 2;
    uint256 public constant WORKNET_3 = 3;
    uint256 public constant DEPOSIT_AMOUNT = 1000 ether;
    uint256 public constant EPOCH_DURATION = 7 days;
    uint256 public genesisTime;

    // ── Delegate state (test contract acts as awpRegistry) ──
    mapping(address => mapping(address => bool)) private _delegates;

    /// @dev StakeNFT calls awpRegistry.currentEpoch()
    function currentEpoch() external view returns (uint256) {
        return (block.timestamp - genesisTime) / EPOCH_DURATION;
    }

    /// @dev StakingVault._isAuthorized → IAWPRegistryDelegates.delegates()
    function delegates(address staker, address delegate) external view returns (bool) {
        return _delegates[staker][delegate];
    }

    /// @dev Test helper: set delegate relationship
    function _setDelegate(address staker, address delegate, bool status) internal {
        _delegates[staker][delegate] = status;
    }

    function setUp() public {
        genesisTime = block.timestamp;

        // Deploy AWPToken
        awp = new AWPToken("AWP", "AWP", deployer);
        awp.initialMint(200_000_000 * 1e18);

        // Deploy StakingVault proxy + StakeNFT
        vault = StakingVault(address(new ERC1967Proxy(
            address(new StakingVault()), abi.encodeCall(StakingVault.initialize, (address(this), address(this)))
        )));
        stakeNFT = new StakeNFT(address(awp), address(vault), address(this));
        vault.setStakeNFT(address(stakeNFT));

        // Give user1 tokens and stake
        awp.transfer(user1, 50_000 ether);
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 50_000 ether);
        stakeNFT.deposit(DEPOSIT_AMOUNT, 52 weeks);
        vm.stopPrank();

        // Give user2 tokens and stake
        awp.transfer(user2, 50_000 ether);
        vm.startPrank(user2);
        awp.approve(address(stakeNFT), 50_000 ether);
        stakeNFT.deposit(DEPOSIT_AMOUNT, 52 weeks);
        vm.stopPrank();

        // Give user3 tokens and stake
        awp.transfer(user3, 50_000 ether);
        vm.startPrank(user3);
        awp.approve(address(stakeNFT), 50_000 ether);
        stakeNFT.deposit(DEPOSIT_AMOUNT, 52 weeks);
        vm.stopPrank();
    }

    // ── EIP-712 helpers ──

    function _getVaultDigest(bytes32 structHash) internal view returns (bytes32) {
        bytes32 domainSeparator = keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
                keccak256("StakingVault"),
                keccak256("1"),
                block.chainid,
                address(vault)
            )
        );
        return keccak256(abi.encodePacked("\x19\x01", domainSeparator, structHash));
    }

    // ══════════════════════════════════════════════
    // 1. Reallocate edge cases
    // ══════════════════════════════════════════════

    function test_reallocate_fromWorknetIdZero_reverts() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 500 ether);

        vm.prank(user1);
        vm.expectRevert(StakingVault.ZeroWorknetId.selector);
        vault.reallocate(user1, agent1, 0, agent2, WORKNET_2, 200 ether);
    }

    function test_reallocate_toWorknetIdZero_reverts() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 500 ether);

        vm.prank(user1);
        vm.expectRevert(StakingVault.ZeroWorknetId.selector);
        vault.reallocate(user1, agent1, WORKNET_1, agent2, 0, 200 ether);
    }

    function test_reallocate_fullAmount_removesFromAgentWorknets() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 500 ether);

        assertEq(vault.getAgentWorknets(user1, agent1).length, 1);

        // Move full amount from agent1/WORKNET_1 to agent2/WORKNET_2
        vm.prank(user1);
        vault.reallocate(user1, agent1, WORKNET_1, agent2, WORKNET_2, 500 ether);

        // Source agent's worknet set should be cleared
        assertEq(vault.getAgentWorknets(user1, agent1).length, 0);
        assertEq(vault.getAgentStake(user1, agent1, WORKNET_1), 0);

        // Destination agent has new allocation
        assertEq(vault.getAgentWorknets(user1, agent2).length, 1);
        assertEq(vault.getAgentStake(user1, agent2, WORKNET_2), 500 ether);

        // Worknet totals are correct
        assertEq(vault.worknetTotalStake(WORKNET_1), 0);
        assertEq(vault.worknetTotalStake(WORKNET_2), 500 ether);

        // User total allocated unchanged
        assertEq(vault.userTotalAllocated(user1), 500 ether);
    }

    function test_reallocate_selfReallocate_sameAgentWorknet() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 500 ether);

        // Self-reallocation: same (agent, worknetId)
        vm.prank(user1);
        vault.reallocate(user1, agent1, WORKNET_1, agent1, WORKNET_1, 200 ether);

        // Amount unchanged (subtract 200, add 200)
        assertEq(vault.getAgentStake(user1, agent1, WORKNET_1), 500 ether);
        assertEq(vault.worknetTotalStake(WORKNET_1), 500 ether);
        assertEq(vault.userTotalAllocated(user1), 500 ether);
        assertEq(vault.getAgentWorknets(user1, agent1).length, 1);
    }

    function test_reallocate_betweenDifferentAgents_sameWorknet() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 500 ether);

        // Same worknet, move from agent1 to agent2
        vm.prank(user1);
        vault.reallocate(user1, agent1, WORKNET_1, agent2, WORKNET_1, 200 ether);

        assertEq(vault.getAgentStake(user1, agent1, WORKNET_1), 300 ether);
        assertEq(vault.getAgentStake(user1, agent2, WORKNET_1), 200 ether);

        // Worknet total unchanged (intra-worknet move)
        assertEq(vault.worknetTotalStake(WORKNET_1), 500 ether);
        assertEq(vault.userTotalAllocated(user1), 500 ether);
    }

    // ══════════════════════════════════════════════
    // ══════════════════════════════════════════════
    // 2. Multi-user stress tests
    // ══════════════════════════════════════════════

    function test_multipleUsers_sameAgentWorknet() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 300 ether);

        vm.prank(user2);
        vault.allocate(user2, agent1, WORKNET_1, 400 ether);

        vm.prank(user3);
        vault.allocate(user3, agent1, WORKNET_1, 200 ether);

        // worknetTotalStake is sum of all users
        assertEq(vault.worknetTotalStake(WORKNET_1), 900 ether);

        // Each user's allocation is correct
        assertEq(vault.getAgentStake(user1, agent1, WORKNET_1), 300 ether);
        assertEq(vault.getAgentStake(user2, agent1, WORKNET_1), 400 ether);
        assertEq(vault.getAgentStake(user3, agent1, WORKNET_1), 200 ether);
    }

    function test_worknetTotalStake_aggregation() public {
        // Multiple users allocating across multiple worknets
        vm.prank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 100 ether);

        vm.prank(user2);
        vault.allocate(user2, agent1, WORKNET_1, 200 ether);

        vm.prank(user1);
        vault.allocate(user1, agent2, WORKNET_2, 150 ether);

        vm.prank(user3);
        vault.allocate(user3, agent2, WORKNET_2, 250 ether);

        assertEq(vault.worknetTotalStake(WORKNET_1), 300 ether);
        assertEq(vault.worknetTotalStake(WORKNET_2), 400 ether);
    }

    function test_oneUserDeallocates_othersUnaffected() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 300 ether);

        vm.prank(user2);
        vault.allocate(user2, agent1, WORKNET_1, 400 ether);

        // user1 deallocates
        vm.prank(user1);
        vault.deallocate(user1, agent1, WORKNET_1, 300 ether);

        // user2 unaffected
        assertEq(vault.getAgentStake(user2, agent1, WORKNET_1), 400 ether);
        assertEq(vault.userTotalAllocated(user2), 400 ether);

        // Worknet total correctly reduced
        assertEq(vault.worknetTotalStake(WORKNET_1), 400 ether);

        // user1 zeroed
        assertEq(vault.getAgentStake(user1, agent1, WORKNET_1), 0);
        assertEq(vault.userTotalAllocated(user1), 0);
    }

    // ══════════════════════════════════════════════
    // 4. Amount boundary tests
    // ══════════════════════════════════════════════

    function test_allocate_exceedsUint128Max_reverts() public {
        uint256 tooLarge = uint256(type(uint128).max) + 1;

        vm.prank(user1);
        vm.expectRevert(StakingVault.AmountExceedsUint128.selector);
        vault.allocate(user1, agent1, WORKNET_1, tooLarge);
    }

    function test_allocate_exactUint128Max_insufficientStake() public {
        // uint128 max exceeds user stake, so should revert InsufficientUnallocated
        // Verifies uint128 max passes amount check but fails balance check
        vm.prank(user1);
        vm.expectRevert(StakingVault.InsufficientUnallocated.selector);
        vault.allocate(user1, agent1, WORKNET_1, type(uint128).max);
    }

    function test_deallocate_moreThanAllocated_reverts() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 200 ether);

        vm.prank(user1);
        vm.expectRevert(StakingVault.InsufficientAllocation.selector);
        vault.deallocate(user1, agent1, WORKNET_1, 201 ether);
    }

    function test_reallocate_exceedsUint128Max_reverts() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 500 ether);

        uint256 tooLarge = uint256(type(uint128).max) + 1;
        vm.prank(user1);
        vm.expectRevert(StakingVault.AmountExceedsUint128.selector);
        vault.reallocate(user1, agent1, WORKNET_1, agent2, WORKNET_2, tooLarge);
    }

    function test_deallocate_exceedsUint128Max_reverts() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 500 ether);

        uint256 tooLarge = uint256(type(uint128).max) + 1;
        vm.prank(user1);
        vm.expectRevert(StakingVault.AmountExceedsUint128.selector);
        vault.deallocate(user1, agent1, WORKNET_1, tooLarge);
    }

    // ══════════════════════════════════════════════
    // 5. Gasless replay protection
    // ══════════════════════════════════════════════

    function test_allocateFor_expiredDeadline_reverts() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("expiredDeadline");

        awp.transfer(signer, 10_000 ether);
        vm.startPrank(signer);
        awp.approve(address(stakeNFT), 10_000 ether);
        stakeNFT.deposit(1000 ether, 52 weeks);
        vm.stopPrank();

        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = vault.nonces(signer);

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Allocate(address staker,address agent,uint256 worknetId,uint256 amount,uint256 nonce,uint256 deadline)"),
            signer, agent1, WORKNET_1, uint256(500 ether), nonce, deadline
        ));
        bytes32 digest = _getVaultDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        // Advance past deadline
        vm.warp(deadline + 1);

        vm.prank(user2);
        vm.expectRevert(StakingVault.ExpiredSignature.selector);
        vault.allocateFor(signer, agent1, WORKNET_1, 500 ether, deadline, v, r, s);
    }

    function test_allocateFor_wrongSigner_reverts() public {
        (address signer,) = makeAddrAndKey("correctSigner");
        (, uint256 wrongPk) = makeAddrAndKey("wrongSigner");

        awp.transfer(signer, 10_000 ether);
        vm.startPrank(signer);
        awp.approve(address(stakeNFT), 10_000 ether);
        stakeNFT.deposit(1000 ether, 52 weeks);
        vm.stopPrank();

        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = vault.nonces(signer);

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Allocate(address staker,address agent,uint256 worknetId,uint256 amount,uint256 nonce,uint256 deadline)"),
            signer, agent1, WORKNET_1, uint256(500 ether), nonce, deadline
        ));
        bytes32 digest = _getVaultDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(wrongPk, digest);

        vm.prank(user2);
        vm.expectRevert(StakingVault.InvalidSignature.selector);
        vault.allocateFor(signer, agent1, WORKNET_1, 500 ether, deadline, v, r, s);
    }

    function test_allocateFor_reuseSameNonce_reverts() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("replayNonce");

        awp.transfer(signer, 10_000 ether);
        vm.startPrank(signer);
        awp.approve(address(stakeNFT), 10_000 ether);
        stakeNFT.deposit(2000 ether, 52 weeks);
        vm.stopPrank();

        uint256 amount = 300 ether;
        uint256 nonce = vault.nonces(signer);
        uint256 deadline = block.timestamp + 1 hours;

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Allocate(address staker,address agent,uint256 worknetId,uint256 amount,uint256 nonce,uint256 deadline)"),
            signer, agent1, WORKNET_1, amount, nonce, deadline
        ));
        bytes32 digest = _getVaultDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        // First call succeeds
        vm.prank(user2);
        vault.allocateFor(signer, agent1, WORKNET_1, amount, deadline, v, r, s);
        assertEq(vault.nonces(signer), nonce + 1);

        // Replay attack: reuse same signature
        vm.prank(user2);
        vm.expectRevert(StakingVault.InvalidSignature.selector);
        vault.allocateFor(signer, agent1, WORKNET_1, amount, deadline, v, r, s);
    }

    function test_deallocateFor_wrongNonce_reverts() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("wrongNonce");

        awp.transfer(signer, 10_000 ether);
        vm.startPrank(signer);
        awp.approve(address(stakeNFT), 10_000 ether);
        stakeNFT.deposit(1000 ether, 52 weeks);
        vault.allocate(signer, agent1, WORKNET_1, 500 ether);
        vm.stopPrank();

        uint256 amount = 200 ether;
        uint256 wrongNonce = vault.nonces(signer) + 1; // Intentionally wrong nonce
        uint256 deadline = block.timestamp + 1 hours;

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Deallocate(address staker,address agent,uint256 worknetId,uint256 amount,uint256 nonce,uint256 deadline)"),
            signer, agent1, WORKNET_1, amount, wrongNonce, deadline
        ));
        bytes32 digest = _getVaultDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        vm.prank(user2);
        vm.expectRevert(StakingVault.InvalidSignature.selector);
        vault.deallocateFor(signer, agent1, WORKNET_1, amount, deadline, v, r, s);
    }

    // ══════════════════════════════════════════════
    // 6. getAgentWorknets view function
    // ══════════════════════════════════════════════

    function test_getAgentWorknets_afterAllocate() public {
        vm.startPrank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 100 ether);
        vault.allocate(user1, agent1, WORKNET_2, 200 ether);
        vm.stopPrank();

        uint256[] memory worknets = vault.getAgentWorknets(user1, agent1);
        assertEq(worknets.length, 2);

        // Verify both worknets in set (order not guaranteed)
        bool hasW1 = false;
        bool hasW2 = false;
        for (uint256 i = 0; i < worknets.length; i++) {
            if (worknets[i] == WORKNET_1) hasW1 = true;
            if (worknets[i] == WORKNET_2) hasW2 = true;
        }
        assertTrue(hasW1, "Should contain WORKNET_1");
        assertTrue(hasW2, "Should contain WORKNET_2");
    }

    function test_getAgentWorknets_removedAfterFullDeallocate() public {
        vm.startPrank(user1);
        vault.allocate(user1, agent1, WORKNET_1, 100 ether);
        vault.allocate(user1, agent1, WORKNET_2, 200 ether);
        vm.stopPrank();

        // Fully deallocate WORKNET_1
        vm.prank(user1);
        vault.deallocate(user1, agent1, WORKNET_1, 100 ether);

        uint256[] memory worknets = vault.getAgentWorknets(user1, agent1);
        assertEq(worknets.length, 1);
        assertEq(worknets[0], WORKNET_2);
    }

    function test_getAgentWorknets_emptyForUnknownUser() public {
        address unknown = makeAddr("unknownUser");
        uint256[] memory worknets = vault.getAgentWorknets(unknown, agent1);
        assertEq(worknets.length, 0);
    }

    function test_getAgentWorknets_emptyForUnknownAgent() public {
        address unknownAgent = makeAddr("unknownAgent");
        uint256[] memory worknets = vault.getAgentWorknets(user1, unknownAgent);
        assertEq(worknets.length, 0);
    }

    // ══════════════════════════════════════════════
    // 7. Guardian rotation
    // ══════════════════════════════════════════════

    function test_setGuardian_byCurrentGuardian_succeeds() public {
        address newGuardian = makeAddr("newGuardian");

        // address(this) is current guardian
        vault.setGuardian(newGuardian);
        assertEq(vault.guardian(), newGuardian);
    }

    function test_setGuardian_byNonGuardian_reverts() public {
        address newGuardian = makeAddr("newGuardian");

        vm.prank(user1);
        vm.expectRevert(StakingVault.NotGuardian.selector);
        vault.setGuardian(newGuardian);
    }

    function test_setGuardian_toZeroAddress_reverts() public {
        vm.expectRevert(StakingVault.ZeroAddress.selector);
        vault.setGuardian(address(0));
    }

    function test_setGuardian_newGuardianCanUpgrade() public {
        address newGuardian = makeAddr("newGuardian2");

        // Rotate guardian
        vault.setGuardian(newGuardian);

        // Old guardian cannot upgrade
        StakingVault newImpl = new StakingVault();
        vm.expectRevert(StakingVault.NotGuardian.selector);
        vault.upgradeToAndCall(address(newImpl), "");

        // New guardian can upgrade
        StakingVault newImpl2 = new StakingVault();
        vm.prank(newGuardian);
        vault.upgradeToAndCall(address(newImpl2), "");

        // State preserved after upgrade
        assertEq(vault.guardian(), newGuardian);
        assertEq(vault.awpRegistry(), address(this));
    }
}
