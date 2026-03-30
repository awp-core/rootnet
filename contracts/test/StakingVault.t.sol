// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract StakingVaultTest is Test {
    StakingVault public vault;
    StakeNFT public stakeNFT;
    AWPToken public awp;

    address public deployer = address(this);
    address public user1 = makeAddr("user1");
    address public user2 = makeAddr("user2");
    address public agent1 = makeAddr("agent1");
    address public agent2 = makeAddr("agent2");
    address public nonAWPRegistry = makeAddr("nonAWPRegistry");

    uint256 public constant SUBNET_1 = 1;
    uint256 public constant SUBNET_2 = 2;
    uint256 public constant DEPOSIT_AMOUNT = 1000 ether;
    uint256 public constant EPOCH_DURATION = 7 days;
    uint256 public genesisTime;

    /// @dev This test contract acts as awpRegistry; StakeNFT calls awpRegistry.currentEpoch()
    function currentEpoch() external view returns (uint256) {
        return (block.timestamp - genesisTime) / EPOCH_DURATION;
    }

    /// @dev Required by StakingVault._isAuthorized → IAWPRegistryDelegates.delegates()
    function delegates(address, address) external pure returns (bool) {
        return false;
    }

    /// @dev Required by StakingVault._authorizeUpgrade → IAWPRegistryDelegates.treasury()
    function treasury() external pure returns (address) {
        return address(0);
    }

    function setUp() public {
        genesisTime = block.timestamp;

        // Deploy AWPToken (deployer gets INITIAL_MINT)
        awp = new AWPToken("AWP", "AWP", deployer, 200_000_000 * 1e18);

        // Deploy StakingVault + StakeNFT (circular dependency)
        // This test contract (address(this)) acts as awpRegistry
        uint64 nonce = vm.getNonce(deployer);
        address predictedVault = vm.computeCreateAddress(deployer, nonce);
        address predictedStakeNFT = vm.computeCreateAddress(deployer, nonce + 1);

        vault = StakingVault(address(new ERC1967Proxy(
            address(new StakingVault()), abi.encodeCall(StakingVault.initialize, (address(this)))
        )));
        stakeNFT = new StakeNFT(address(awp), address(vault), address(this));
        vault.setStakeNFT(address(stakeNFT));

        // Give user1 AWP and have them deposit into StakeNFT
        awp.transfer(user1, 10_000 ether);
        vm.startPrank(user1);
        awp.approve(address(stakeNFT), 10_000 ether);
        stakeNFT.deposit(DEPOSIT_AMOUNT, 52 weeks);
        vm.stopPrank();
    }

    // ══════════════════════════════════════════════
    // Allocate tests
    // ══════════════════════════════════════════════

    function test_allocate_basic() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 300 ether);

        assertEq(vault.getAgentStake(user1, agent1, SUBNET_1), 300 ether);
        assertEq(vault.userTotalAllocated(user1), 300 ether);
        assertEq(vault.subnetTotalStake(SUBNET_1), 300 ether);
    }

    function test_allocate_moreThanUnallocated_reverts() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 800 ether);

        // Only 200 unallocated, allocating 300 should fail
        vm.prank(user1);
        vm.expectRevert(StakingVault.InsufficientUnallocated.selector);
        vault.allocate(user1, agent2, SUBNET_2, 300 ether);
    }

    function test_allocate_zeroAmount_reverts() public {
        vm.prank(user1);
        vm.expectRevert(StakingVault.InvalidAmount.selector);
        vault.allocate(user1, agent1, SUBNET_1, 0);
    }

    // ══════════════════════════════════════════════
    // Deallocate tests
    // ══════════════════════════════════════════════

    function test_deallocate_basic() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 500 ether);
        vm.prank(user1);
        vault.deallocate(user1, agent1, SUBNET_1, 200 ether);

        assertEq(vault.getAgentStake(user1, agent1, SUBNET_1), 300 ether);
        assertEq(vault.userTotalAllocated(user1), 300 ether);
        assertEq(vault.subnetTotalStake(SUBNET_1), 300 ether);
    }

    function test_deallocate_full_zerosStake() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 500 ether);
        vm.prank(user1);
        vault.deallocate(user1, agent1, SUBNET_1, 500 ether);

        assertEq(vault.getAgentStake(user1, agent1, SUBNET_1), 0);
    }

    function test_deallocate_moreThanAllocated_reverts() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 200 ether);

        vm.prank(user1);
        vm.expectRevert(StakingVault.InsufficientAllocation.selector);
        vault.deallocate(user1, agent1, SUBNET_1, 300 ether);
    }

    // ══════════════════════════════════════════════
    // Reallocate tests (immediate, no dual-slot)
    // ══════════════════════════════════════════════

    function test_reallocate_immediate() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 500 ether);

        vm.prank(user1);
        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 200 ether);

        // Immediate effect
        assertEq(vault.getAgentStake(user1, agent1, SUBNET_1), 300 ether);
        assertEq(vault.getAgentStake(user1, agent2, SUBNET_2), 200 ether);

        // Subnet totals
        assertEq(vault.subnetTotalStake(SUBNET_1), 300 ether);
        assertEq(vault.subnetTotalStake(SUBNET_2), 200 ether);

        // userTotalAllocated unchanged
        assertEq(vault.userTotalAllocated(user1), 500 ether);
    }

    function test_reallocate_multipleAccumulate() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 500 ether);

        vm.prank(user1);
        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 100 ether);
        vm.prank(user1);
        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 150 ether);

        assertEq(vault.getAgentStake(user1, agent1, SUBNET_1), 250 ether);
        assertEq(vault.getAgentStake(user1, agent2, SUBNET_2), 250 ether);
    }

    function test_reallocate_insufficientAllocation_reverts() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 100 ether);

        vm.prank(user1);
        vm.expectRevert(StakingVault.InsufficientAllocation.selector);
        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 200 ether);
    }

    function test_reallocate_zeroAmount_reverts() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 100 ether);

        vm.prank(user1);
        vm.expectRevert(StakingVault.InvalidAmount.selector);
        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 0);
    }

    // ══════════════════════════════════════════════
    // Freeze Agent allocations tests
    // ══════════════════════════════════════════════

    function test_freezeAgentAllocations_immediateRelease() public {
        vm.startPrank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 300 ether);
        vault.allocate(user1, agent1, SUBNET_2, 200 ether);
        vm.stopPrank();

        vault.freezeAgentAllocations(user1, agent1);

        // Allocations zeroed
        assertEq(vault.getAgentStake(user1, agent1, SUBNET_1), 0);
        assertEq(vault.getAgentStake(user1, agent1, SUBNET_2), 0);

        // Subnet totals reduced
        assertEq(vault.subnetTotalStake(SUBNET_1), 0);
        assertEq(vault.subnetTotalStake(SUBNET_2), 0);

        // userTotalAllocated released
        assertEq(vault.userTotalAllocated(user1), 0);
    }

    function test_freezeAgentAllocations_agentSubnetsCleared() public {
        vm.startPrank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 300 ether);
        vault.allocate(user1, agent1, SUBNET_2, 200 ether);
        vm.stopPrank();

        assertEq(vault.getAgentSubnets(user1, agent1).length, 2);

        vault.freezeAgentAllocations(user1, agent1);

        // Set must be fully cleared after freeze
        assertEq(vault.getAgentSubnets(user1, agent1).length, 0);
    }

    function test_deallocate_full_clearsAgentSubnets() public {
        vm.prank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 500 ether);
        assertEq(vault.getAgentSubnets(user1, agent1).length, 1);

        vm.prank(user1);
        vault.deallocate(user1, agent1, SUBNET_1, 500 ether);

        // Subnet removed from set after full deallocation
        assertEq(vault.getAgentSubnets(user1, agent1).length, 0);
    }

    function test_freezeAfterReallocate_setsConsistent() public {
        vm.startPrank(user1);
        vault.allocate(user1, agent1, SUBNET_1, 500 ether);

        // Reallocate everything from agent1/SUBNET_1 to agent2/SUBNET_2
        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 500 ether);
        vm.stopPrank();

        // agent1 should have no subnets left
        assertEq(vault.getAgentSubnets(user1, agent1).length, 0);
        // agent2 should have SUBNET_2
        assertEq(vault.getAgentSubnets(user1, agent2).length, 1);

        // Freeze agent1 — should be a no-op (no allocations)
        vault.freezeAgentAllocations(user1, agent1);
        assertEq(vault.userTotalAllocated(user1), 500 ether);

        // Freeze agent2 — should clear everything
        vault.freezeAgentAllocations(user1, agent2);
        assertEq(vault.getAgentStake(user1, agent2, SUBNET_2), 0);
        assertEq(vault.getAgentSubnets(user1, agent2).length, 0);
        assertEq(vault.userTotalAllocated(user1), 0);
    }

    // ══════════════════════════════════════════════
    // onlyAWPRegistry access control tests
    // ══════════════════════════════════════════════

    function test_notAuthorized_allocate() public {
        vm.prank(nonAWPRegistry);
        vm.expectRevert(StakingVault.NotAuthorized.selector);
        vault.allocate(user1, agent1, SUBNET_1, 100 ether);
    }

    function test_notAuthorized_deallocate() public {
        vm.prank(nonAWPRegistry);
        vm.expectRevert(StakingVault.NotAuthorized.selector);
        vault.deallocate(user1, agent1, SUBNET_1, 100 ether);
    }

    function test_notAuthorized_reallocate() public {
        vm.prank(nonAWPRegistry);
        vm.expectRevert(StakingVault.NotAuthorized.selector);
        vault.reallocate(user1, agent1, SUBNET_1, agent2, SUBNET_2, 100 ether);
    }

    function test_onlyAWPRegistry_freezeAgentAllocations() public {
        vm.prank(nonAWPRegistry);
        vm.expectRevert(StakingVault.NotAWPRegistry.selector);
        vault.freezeAgentAllocations(user1, agent1);
    }

    // ══════════════════════════════════════════════
    // Gasless EIP-712 allocateFor / deallocateFor tests
    // ══════════════════════════════════════════════

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

    function test_allocateFor_gasless() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("gaslessSigner");

        // 给 signer 代币并质押
        awp.transfer(signer, 10_000 ether);
        vm.startPrank(signer);
        awp.approve(address(stakeNFT), 10_000 ether);
        stakeNFT.deposit(1000 ether, 52 weeks);
        vm.stopPrank();

        uint256 amount = 500 ether;
        uint256 nonce = vault.nonces(signer);
        uint256 deadline = block.timestamp + 1 hours;

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Allocate(address staker,address agent,uint256 subnetId,uint256 amount,uint256 nonce,uint256 deadline)"),
            signer, agent1, SUBNET_1, amount, nonce, deadline
        ));
        bytes32 digest = _getVaultDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        // relayer 提交 gasless 交易
        vm.prank(user2);
        vault.allocateFor(signer, agent1, SUBNET_1, amount, deadline, v, r, s);

        assertEq(vault.getAgentStake(signer, agent1, SUBNET_1), amount);
        assertEq(vault.nonces(signer), nonce + 1);
    }

    function test_allocateFor_expiredSignature_reverts() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("expiredSigner");

        awp.transfer(signer, 10_000 ether);
        vm.startPrank(signer);
        awp.approve(address(stakeNFT), 10_000 ether);
        stakeNFT.deposit(1000 ether, 52 weeks);
        vm.stopPrank();

        uint256 deadline = block.timestamp - 1; // 过期

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Allocate(address staker,address agent,uint256 subnetId,uint256 amount,uint256 nonce,uint256 deadline)"),
            signer, agent1, SUBNET_1, uint256(500 ether), uint256(0), deadline
        ));
        bytes32 digest = _getVaultDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        vm.prank(user2);
        vm.expectRevert(StakingVault.ExpiredSignature.selector);
        vault.allocateFor(signer, agent1, SUBNET_1, 500 ether, deadline, v, r, s);
    }

    function test_allocateFor_wrongSigner_reverts() public {
        (address signer, ) = makeAddrAndKey("wrongSignerUser");
        (, uint256 wrongPk) = makeAddrAndKey("wrongKey");

        awp.transfer(signer, 10_000 ether);
        vm.startPrank(signer);
        awp.approve(address(stakeNFT), 10_000 ether);
        stakeNFT.deposit(1000 ether, 52 weeks);
        vm.stopPrank();

        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = vault.nonces(signer);

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Allocate(address staker,address agent,uint256 subnetId,uint256 amount,uint256 nonce,uint256 deadline)"),
            signer, agent1, SUBNET_1, uint256(500 ether), nonce, deadline
        ));
        bytes32 digest = _getVaultDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(wrongPk, digest); // 使用错误的私钥签名

        vm.prank(user2);
        vm.expectRevert(StakingVault.InvalidSignature.selector);
        vault.allocateFor(signer, agent1, SUBNET_1, 500 ether, deadline, v, r, s);
    }

    function test_allocateFor_replayProtection() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("replaySigner");

        awp.transfer(signer, 10_000 ether);
        vm.startPrank(signer);
        awp.approve(address(stakeNFT), 10_000 ether);
        stakeNFT.deposit(2000 ether, 52 weeks);
        vm.stopPrank();

        uint256 amount = 500 ether;
        uint256 nonce = vault.nonces(signer);
        uint256 deadline = block.timestamp + 1 hours;

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Allocate(address staker,address agent,uint256 subnetId,uint256 amount,uint256 nonce,uint256 deadline)"),
            signer, agent1, SUBNET_1, amount, nonce, deadline
        ));
        bytes32 digest = _getVaultDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        // 第一次调用成功
        vm.prank(user2);
        vault.allocateFor(signer, agent1, SUBNET_1, amount, deadline, v, r, s);

        // 第二次使用相同签名应该失败（nonce 已递增）
        vm.prank(user2);
        vm.expectRevert(StakingVault.InvalidSignature.selector);
        vault.allocateFor(signer, agent1, SUBNET_1, amount, deadline, v, r, s);
    }

    function test_deallocateFor_gasless() public {
        (address signer, uint256 signerPk) = makeAddrAndKey("deallocSigner");

        awp.transfer(signer, 10_000 ether);
        vm.startPrank(signer);
        awp.approve(address(stakeNFT), 10_000 ether);
        stakeNFT.deposit(1000 ether, 52 weeks);
        vault.allocate(signer, agent1, SUBNET_1, 500 ether);
        vm.stopPrank();

        uint256 amount = 200 ether;
        uint256 nonce = vault.nonces(signer);
        uint256 deadline = block.timestamp + 1 hours;

        bytes32 structHash = keccak256(abi.encode(
            keccak256("Deallocate(address staker,address agent,uint256 subnetId,uint256 amount,uint256 nonce,uint256 deadline)"),
            signer, agent1, SUBNET_1, amount, nonce, deadline
        ));
        bytes32 digest = _getVaultDigest(structHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, digest);

        vm.prank(user2);
        vault.deallocateFor(signer, agent1, SUBNET_1, amount, deadline, v, r, s);

        assertEq(vault.getAgentStake(signer, agent1, SUBNET_1), 300 ether);
        assertEq(vault.nonces(signer), nonce + 1);
    }

    // ══════════════════════════════════════════════
    // UUPS upgrade authorization tests
    // ══════════════════════════════════════════════

    function test_vaultUpgradeByNonTreasury_reverts() public {
        StakingVault newImpl = new StakingVault();
        vm.prank(user1);
        vm.expectRevert(StakingVault.NotAWPRegistry.selector);
        vault.upgradeToAndCall(address(newImpl), "");
    }
}
