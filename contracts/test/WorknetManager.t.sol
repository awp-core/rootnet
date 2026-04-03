// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, Vm} from "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {WorknetManager} from "../src/worknets/WorknetManager.sol";
import {IERC1363Receiver} from "../src/interfaces/IERC1363Receiver.sol";

// ═══════════════════════════════════════════════════
//  Mock: AWPRegistry — resolveRecipient returns mapped address
// ═══════════════════════════════════════════════════
contract MockAWPRegistry {
    mapping(address => address) public recipientOverrides;

    function setRecipientOverride(address from, address to) external {
        recipientOverrides[from] = to;
    }

    function resolveRecipient(address addr) external view returns (address) {
        address r = recipientOverrides[addr];
        return r == address(0) ? addr : r;
    }

    function batchResolveRecipients(address[] calldata addrs) external view returns (address[] memory result) {
        result = new address[](addrs.length);
        for (uint256 i = 0; i < addrs.length; i++) {
            address r = recipientOverrides[addrs[i]];
            result[i] = r == address(0) ? addrs[i] : r;
        }
    }
}

// ═══════════════════════════════════════════════════
//  Mock: AlphaToken — records mint calls
// ═══════════════════════════════════════════════════
contract MockAlphaToken {
    mapping(address => uint256) public minted;
    uint256 public totalMinted;

    function mint(address to, uint256 amount) external {
        minted[to] += amount;
        totalMinted += amount;
    }

    function burn(uint256) external {}
    function balanceOf(address) external pure returns (uint256) { return 0; }
}

// ═══════════════════════════════════════════════════
//  Mock: AWPToken — ERC20 + transferAndCall support
// ═══════════════════════════════════════════════════
contract MockAWPToken {
    string public name = "AWP";
    string public symbol = "AWP";
    uint8 public decimals = 18;
    uint256 public totalSupply;
    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    function mint(address to, uint256 amount) external {
        balanceOf[to] += amount;
        totalSupply += amount;
    }

    function transfer(address to, uint256 amount) external returns (bool) {
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;
        return true;
    }

    function transferFrom(address from, address to, uint256 amount) external returns (bool) {
        if (allowance[from][msg.sender] != type(uint256).max) {
            allowance[from][msg.sender] -= amount;
        }
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        return true;
    }

    function approve(address spender, uint256 amount) external returns (bool) {
        allowance[msg.sender][spender] = amount;
        return true;
    }

    /// @dev Mock transferAndCall: transfer then call receiver's onTransferReceived
    function transferAndCall(address to, uint256 amount, bytes calldata data) external returns (bool) {
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;
        bytes4 ret = IERC1363Receiver(to).onTransferReceived(msg.sender, msg.sender, amount, data);
        require(ret == IERC1363Receiver.onTransferReceived.selector, "ERC1363: bad return");
        return true;
    }
}

// ═══════════════════════════════════════════════════
//  Mock: WorknetManagerV2 — for UUPS upgrade tests
// ═══════════════════════════════════════════════════
contract WorknetManagerV2 is WorknetManager {
    uint256 public v2Value;

    function setV2Value(uint256 val) external {
        v2Value = val;
    }

    function version() external pure returns (string memory) {
        return "2.0";
    }
}

// ═══════════════════════════════════════════════════
//  Test Suite
// ═══════════════════════════════════════════════════
contract WorknetManagerTest is Test {
    WorknetManager public manager;
    WorknetManager public impl;
    MockAWPRegistry public registry;
    MockAlphaToken public alpha;
    MockAWPToken public awp;

    address public admin = makeAddr("admin");
    address public merkleUser = makeAddr("merkleUser");
    address public strategyUser = makeAddr("strategyUser");
    address public transferUser = makeAddr("transferUser");
    address public nobody = makeAddr("nobody");
    address public recipient = makeAddr("recipient");

    // Mock DEX addresses (Reserve strategy won't actually call them)
    address public mockPoolManager = makeAddr("poolManager");
    address public mockPositionManager = makeAddr("positionManager");
    address public mockSwapRouter = makeAddr("swapRouter");
    address public mockPermit2 = makeAddr("permit2");
    uint24 public mockPoolFee = 500;
    int24 public mockTickSpacing = 10;

    bytes32 public constant MERKLE_ROLE = keccak256("MERKLE_ROLE");
    bytes32 public constant STRATEGY_ROLE = keccak256("STRATEGY_ROLE");
    bytes32 public constant TRANSFER_ROLE = keccak256("TRANSFER_ROLE");
    bytes32 public constant POOL_ID = keccak256("test-pool");

    function setUp() public {
        // Deploy mock dependencies
        registry = new MockAWPRegistry();
        alpha = new MockAlphaToken();
        awp = new MockAWPToken();

        // Set resolveRecipient mapping
        registry.setRecipientOverride(merkleUser, recipient);

        // Deploy WorknetManager implementation
        impl = new WorknetManager();

        // Encode dexConfig
        bytes memory dexConfig = abi.encode(
            mockPoolManager, mockPositionManager, mockSwapRouter, mockPermit2, mockPoolFee, mockTickSpacing
        );

        // Deploy ERC1967Proxy
        bytes memory initData = abi.encodeCall(
            WorknetManager.initialize,
            (address(registry), address(alpha), address(awp), POOL_ID, admin, dexConfig)
        );
        ERC1967Proxy proxy = new ERC1967Proxy(address(impl), initData);
        manager = WorknetManager(address(proxy));
    }

    // ═══════════════════════════════════════════════
    //  1. Initialize
    // ═══════════════════════════════════════════════

    function test_initialize_setsStorage() public view {
        assertEq(address(manager.awpRegistry()), address(registry));
        assertEq(address(manager.alphaToken()), address(alpha));
        assertEq(address(manager.awpToken()), address(awp));
        assertEq(manager.poolId(), POOL_ID);
        assertEq(manager.clPoolManager(), mockPoolManager);
        assertEq(manager.clPositionManager(), mockPositionManager);
        assertEq(manager.clSwapRouter(), mockSwapRouter);
        assertEq(manager.permit2(), mockPermit2);
        assertEq(manager.poolFee(), mockPoolFee);
        assertEq(manager.tickSpacing(), mockTickSpacing);
    }

    function test_initialize_setsRoles() public view {
        assertTrue(manager.hasRole(manager.DEFAULT_ADMIN_ROLE(), admin));
        assertTrue(manager.hasRole(MERKLE_ROLE, admin));
        assertTrue(manager.hasRole(STRATEGY_ROLE, admin));
        assertTrue(manager.hasRole(TRANSFER_ROLE, admin));
    }

    function test_initialize_setsPoolKey() public view {
        (address c0, address c1,,,, ) = manager.poolKey();
        // poolKey.currency0 should be smaller address
        if (address(awp) < address(alpha)) {
            assertEq(c0, address(awp));
            assertEq(c1, address(alpha));
        } else {
            assertEq(c0, address(alpha));
            assertEq(c1, address(awp));
        }
    }

    function test_initialize_defaultStrategyIsReserve() public view {
        assertEq(uint8(manager.currentStrategy()), uint8(WorknetManager.AWPStrategy.Reserve));
    }

    function test_initialize_cannotReinitialize() public {
        bytes memory dexConfig = abi.encode(
            mockPoolManager, mockPositionManager, mockSwapRouter, mockPermit2, mockPoolFee, mockTickSpacing
        );
        vm.expectRevert();
        manager.initialize(address(registry), address(alpha), address(awp), POOL_ID, admin, dexConfig);
    }

    function test_implementation_cannotInitialize() public {
        bytes memory dexConfig = abi.encode(
            mockPoolManager, mockPositionManager, mockSwapRouter, mockPermit2, mockPoolFee, mockTickSpacing
        );
        vm.expectRevert();
        impl.initialize(address(registry), address(alpha), address(awp), POOL_ID, admin, dexConfig);
    }

    // ═══════════════════════════════════════════════
    //  2. setMerkleRoot
    // ═══════════════════════════════════════════════

    function test_setMerkleRoot_success() public {
        bytes32 root = keccak256("root1");
        vm.prank(admin);
        vm.expectEmit(true, false, false, true, address(manager));
        emit WorknetManager.MerkleRootSet(1, root);
        manager.setMerkleRoot(1, root);
        assertEq(manager.merkleRoots(1), root);
    }

    function test_setMerkleRoot_revertOnZeroRoot() public {
        vm.prank(admin);
        vm.expectRevert(WorknetManager.ZeroRoot.selector);
        manager.setMerkleRoot(1, bytes32(0));
    }

    function test_setMerkleRoot_revertOnOverwrite() public {
        bytes32 root = keccak256("root1");
        vm.prank(admin);
        manager.setMerkleRoot(1, root);

        vm.prank(admin);
        vm.expectRevert(WorknetManager.RootAlreadySet.selector);
        manager.setMerkleRoot(1, keccak256("root2"));
    }

    function test_setMerkleRoot_revertUnauthorized() public {
        vm.prank(nobody);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, nobody, MERKLE_ROLE
            )
        );
        manager.setMerkleRoot(1, keccak256("root1"));
    }

    // ═══════════════════════════════════════════════
    //  3. claim — Merkle proof
    // ═══════════════════════════════════════════════

    /// @dev Build single-leaf Merkle tree and verify claim
    function test_claim_validProof_mintsToResolvedRecipient() public {
        uint256 amount = 1000 ether;
        // Double-hash leaf node
        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(merkleUser, amount))));
        // Single-leaf tree: root == leaf
        bytes32 root = leaf;

        vm.prank(admin);
        manager.setMerkleRoot(1, root);

        vm.prank(merkleUser);
        vm.expectEmit(true, true, false, true, address(manager));
        emit WorknetManager.Claimed(1, merkleUser, amount);
        manager.claim(1, amount, new bytes32[](0));

        // Alpha should mint to address from resolveRecipient
        assertEq(alpha.minted(recipient), amount);
        assertTrue(manager.isClaimed(1, merkleUser));
    }

    /// @dev Two-leaf Merkle tree test
    function test_claim_twoLeafTree() public {
        address user1 = merkleUser;
        address user2 = makeAddr("user2");
        uint256 amount1 = 500 ether;
        uint256 amount2 = 700 ether;

        bytes32 leaf1 = keccak256(bytes.concat(keccak256(abi.encode(user1, amount1))));
        bytes32 leaf2 = keccak256(bytes.concat(keccak256(abi.encode(user2, amount2))));

        // Merkle root = hash(sorted pair)
        bytes32 root;
        if (leaf1 <= leaf2) {
            root = keccak256(abi.encodePacked(leaf1, leaf2));
        } else {
            root = keccak256(abi.encodePacked(leaf2, leaf1));
        }

        vm.prank(admin);
        manager.setMerkleRoot(2, root);

        // user1 uses leaf2 as proof
        bytes32[] memory proof1 = new bytes32[](1);
        proof1[0] = leaf2;
        vm.prank(user1);
        manager.claim(2, amount1, proof1);
        assertEq(alpha.minted(recipient), amount1); // user1 -> recipient

        // user2 uses leaf1 as proof (no override, resolves to self)
        bytes32[] memory proof2 = new bytes32[](1);
        proof2[0] = leaf1;
        vm.prank(user2);
        manager.claim(2, amount2, proof2);
        assertEq(alpha.minted(user2), amount2);
    }

    function test_claim_revertInvalidProof() public {
        uint256 amount = 1000 ether;
        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(merkleUser, amount))));

        vm.prank(admin);
        manager.setMerkleRoot(1, leaf);

        // Construct proof with wrong amount
        vm.prank(merkleUser);
        vm.expectRevert(WorknetManager.InvalidProof.selector);
        manager.claim(1, amount + 1, new bytes32[](0));
    }

    function test_claim_revertDoubleClaim() public {
        uint256 amount = 1000 ether;
        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(merkleUser, amount))));

        vm.prank(admin);
        manager.setMerkleRoot(1, leaf);

        vm.prank(merkleUser);
        manager.claim(1, amount, new bytes32[](0));

        vm.prank(merkleUser);
        vm.expectRevert(WorknetManager.AlreadyClaimed.selector);
        manager.claim(1, amount, new bytes32[](0));
    }

    function test_claim_revertNoRoot() public {
        vm.prank(merkleUser);
        vm.expectRevert(WorknetManager.NoRootForEpoch.selector);
        manager.claim(99, 100, new bytes32[](0));
    }

    function test_claim_revertWrongSender() public {
        // Even though leaf is merkleUser, nobody's claim should fail due to invalid proof
        uint256 amount = 1000 ether;
        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(merkleUser, amount))));

        vm.prank(admin);
        manager.setMerkleRoot(1, leaf);

        vm.prank(nobody);
        vm.expectRevert(WorknetManager.InvalidProof.selector);
        manager.claim(1, amount, new bytes32[](0));
    }

    // ═══════════════════════════════════════════════
    //  4. setStrategy
    // ═══════════════════════════════════════════════

    function test_setStrategy_success() public {
        vm.prank(admin);
        vm.expectEmit(true, false, false, false, address(manager));
        emit WorknetManager.StrategyUpdated(WorknetManager.AWPStrategy.AddLiquidity);
        manager.setStrategy(WorknetManager.AWPStrategy.AddLiquidity);
        assertEq(uint8(manager.currentStrategy()), uint8(WorknetManager.AWPStrategy.AddLiquidity));
    }

    function test_setStrategy_canSwitchMultipleTimes() public {
        vm.startPrank(admin);
        manager.setStrategy(WorknetManager.AWPStrategy.BuybackBurn);
        assertEq(uint8(manager.currentStrategy()), uint8(WorknetManager.AWPStrategy.BuybackBurn));

        manager.setStrategy(WorknetManager.AWPStrategy.Reserve);
        assertEq(uint8(manager.currentStrategy()), uint8(WorknetManager.AWPStrategy.Reserve));
        vm.stopPrank();
    }

    function test_setStrategy_revertUnauthorized() public {
        vm.prank(nobody);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, nobody, STRATEGY_ROLE
            )
        );
        manager.setStrategy(WorknetManager.AWPStrategy.Reserve);
    }

    // ═══════════════════════════════════════════════
    //  5. executeStrategy — Reserve (no-op)
    // ═══════════════════════════════════════════════

    function test_executeStrategy_reserve_noOp() public {
        // Reserve strategy doesn't transfer tokens, doesn't emit AWPProcessed
        awp.mint(address(manager), 1000 ether);
        uint256 balBefore = awp.balanceOf(address(manager));

        vm.prank(admin);
        // Should not emit AWPProcessed
        vm.recordLogs();
        manager.executeStrategy(100 ether);

        // Balance unchanged
        assertEq(awp.balanceOf(address(manager)), balBefore);

        // Verify no AWPProcessed event
        Vm.Log[] memory logs = vm.getRecordedLogs();
        bytes32 awpProcessedTopic = keccak256("AWPProcessed(uint8,uint256)");
        for (uint256 i = 0; i < logs.length; i++) {
            assertTrue(logs[i].topics[0] != awpProcessedTopic, "Should not emit AWPProcessed for Reserve");
        }
    }

    function test_executeStrategy_revertZeroAmount() public {
        vm.prank(admin);
        vm.expectRevert(WorknetManager.ZeroAmount.selector);
        manager.executeStrategy(0);
    }

    function test_executeStrategy_revertUnauthorized() public {
        vm.prank(nobody);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, nobody, STRATEGY_ROLE
            )
        );
        manager.executeStrategy(100);
    }

    // ═══════════════════════════════════════════════
    //  6. onTransferReceived
    // ═══════════════════════════════════════════════

    function test_onTransferReceived_awpToken_reserve_returnsSelector() public {
        awp.mint(address(this), 500 ether);

        // Trigger onTransferReceived via transferAndCall
        awp.transferAndCall(address(manager), 500 ether, "");

        // Reserve strategy: tokens stay in contract
        assertEq(awp.balanceOf(address(manager)), 500 ether);
    }

    function test_onTransferReceived_nonAwpToken_returnsSelector() public {
        // Non-AWP token calling onTransferReceived should also return correct selector
        vm.prank(address(0xdead));
        bytes4 ret = manager.onTransferReceived(address(0), address(0), 100, "");
        assertEq(ret, IERC1363Receiver.onTransferReceived.selector);
    }

    function test_onTransferReceived_zeroAmount_noAction() public {
        // amount=0 should not trigger strategy
        vm.prank(address(awp));
        vm.recordLogs();
        bytes4 ret = manager.onTransferReceived(address(0), address(0), 0, "");
        assertEq(ret, IERC1363Receiver.onTransferReceived.selector);

        Vm.Log[] memory logs = vm.getRecordedLogs();
        bytes32 awpProcessedTopic = keccak256("AWPProcessed(uint8,uint256)");
        for (uint256 i = 0; i < logs.length; i++) {
            assertTrue(logs[i].topics[0] != awpProcessedTopic, "Should not emit AWPProcessed for zero amount");
        }
    }

    // ═══════════════════════════════════════════════
    //  7. transferToken
    // ═══════════════════════════════════════════════

    function test_transferToken_success() public {
        awp.mint(address(manager), 1000 ether);

        vm.prank(admin);
        vm.expectEmit(true, true, false, true, address(manager));
        emit WorknetManager.TokenTransferred(address(awp), recipient, 300 ether);
        manager.transferToken(address(awp), recipient, 300 ether);

        assertEq(awp.balanceOf(recipient), 300 ether);
        assertEq(awp.balanceOf(address(manager)), 700 ether);
    }

    function test_transferToken_revertUnauthorized() public {
        vm.prank(nobody);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, nobody, TRANSFER_ROLE
            )
        );
        manager.transferToken(address(awp), nobody, 1);
    }

    // ═══════════════════════════════════════════════
    //  8. UUPS upgrade
    // ═══════════════════════════════════════════════

    function test_upgrade_byAdmin() public {
        WorknetManagerV2 newImpl = new WorknetManagerV2();

        vm.prank(admin);
        manager.upgradeToAndCall(address(newImpl), "");

        // Verify upgrade succeeded
        WorknetManagerV2 upgraded = WorknetManagerV2(address(manager));
        assertEq(upgraded.version(), "2.0");

        // Original storage preserved
        assertEq(address(upgraded.awpRegistry()), address(registry));
        assertEq(upgraded.poolId(), POOL_ID);
    }

    function test_upgrade_revertUnauthorized() public {
        WorknetManagerV2 newImpl = new WorknetManagerV2();

        bytes32 defaultAdminRole = manager.DEFAULT_ADMIN_ROLE();
        vm.prank(nobody);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, nobody, defaultAdminRole
            )
        );
        manager.upgradeToAndCall(address(newImpl), "");
    }

    // ═══════════════════════════════════════════════
    //  9. Access control: granular role tests
    // ═══════════════════════════════════════════════

    function test_grantRole_onlyAdmin() public {
        address newMerkle = makeAddr("newMerkle");

        // admin can grant roles
        vm.prank(admin);
        manager.grantRole(MERKLE_ROLE, newMerkle);
        assertTrue(manager.hasRole(MERKLE_ROLE, newMerkle));

        // New role holder can call setMerkleRoot
        vm.prank(newMerkle);
        manager.setMerkleRoot(10, keccak256("root10"));
    }

    function test_revokeRole_removesAccess() public {
        vm.prank(admin);
        manager.revokeRole(MERKLE_ROLE, admin);

        vm.prank(admin);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, admin, MERKLE_ROLE
            )
        );
        manager.setMerkleRoot(20, keccak256("root20"));
    }

    function test_separateRoles_isolation() public {
        // Give merkleUser only MERKLE_ROLE
        vm.prank(admin);
        manager.grantRole(MERKLE_ROLE, merkleUser);

        // merkleUser can setMerkleRoot
        vm.prank(merkleUser);
        manager.setMerkleRoot(30, keccak256("root30"));

        // merkleUser cannot setStrategy (requires STRATEGY_ROLE)
        vm.prank(merkleUser);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, merkleUser, STRATEGY_ROLE
            )
        );
        manager.setStrategy(WorknetManager.AWPStrategy.Reserve);

        // merkleUser cannot transferToken (requires TRANSFER_ROLE)
        vm.prank(merkleUser);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, merkleUser, TRANSFER_ROLE
            )
        );
        manager.transferToken(address(awp), merkleUser, 1);
    }

    // ═══════════════════════════════════════════════
    //  10. isClaimed view
    // ═══════════════════════════════════════════════

    function test_isClaimed_falseBeforeClaim() public view {
        assertFalse(manager.isClaimed(1, merkleUser));
    }

    function test_isClaimed_trueAfterClaim() public {
        uint256 amount = 100 ether;
        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(merkleUser, amount))));

        vm.prank(admin);
        manager.setMerkleRoot(1, leaf);

        vm.prank(merkleUser);
        manager.claim(1, amount, new bytes32[](0));

        assertTrue(manager.isClaimed(1, merkleUser));
        // Other epochs unaffected
        assertFalse(manager.isClaimed(2, merkleUser));
    }

    // ═══════════════════════════════════════════════
    //  11. Fuzz: claim with random amounts
    // ═══════════════════════════════════════════════

    function testFuzz_claim(uint256 amount) public {
        vm.assume(amount > 0 && amount < type(uint128).max);

        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(merkleUser, amount))));

        vm.prank(admin);
        manager.setMerkleRoot(1, leaf);

        vm.prank(merkleUser);
        manager.claim(1, amount, new bytes32[](0));

        assertEq(alpha.minted(recipient), amount);
        assertTrue(manager.isClaimed(1, merkleUser));
    }

    // ═══════════════════════════════════════════════
    //  12. Multiple epochs
    // ═══════════════════════════════════════════════

    function test_claim_multipleEpochs() public {
        uint256 amount1 = 100 ether;
        uint256 amount2 = 200 ether;

        bytes32 leaf1 = keccak256(bytes.concat(keccak256(abi.encode(merkleUser, amount1))));
        bytes32 leaf2 = keccak256(bytes.concat(keccak256(abi.encode(merkleUser, amount2))));

        vm.startPrank(admin);
        manager.setMerkleRoot(1, leaf1);
        manager.setMerkleRoot(2, leaf2);
        vm.stopPrank();

        // Same user can claim across different epochs
        vm.prank(merkleUser);
        manager.claim(1, amount1, new bytes32[](0));

        vm.prank(merkleUser);
        manager.claim(2, amount2, new bytes32[](0));

        assertEq(alpha.minted(recipient), amount1 + amount2);
        assertTrue(manager.isClaimed(1, merkleUser));
        assertTrue(manager.isClaimed(2, merkleUser));
    }
}
