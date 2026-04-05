// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {MerkleProof} from "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import {WorknetManagerBase, IERC1363Receiver} from "../src/worknets/WorknetManagerBase.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @title MockAlphaForWM — Minimal mock for WorknetToken used by WorknetManager tests
contract MockAlphaForWM {
    mapping(address => uint256) public balances;
    mapping(address => mapping(address => uint256)) public allowances;

    function mint(address to, uint256 amount) external {
        balances[to] += amount;
    }

    function burn(uint256 amount) external {
        balances[msg.sender] -= amount;
    }

    function balanceOf(address account) external view returns (uint256) {
        return balances[account];
    }

    function transfer(address to, uint256 amount) external returns (bool) {
        balances[msg.sender] -= amount;
        balances[to] += amount;
        return true;
    }

    function approve(address spender, uint256 amount) external returns (bool) {
        allowances[msg.sender][spender] = amount;
        return true;
    }
}

/// @title TestableWorknetManager — Minimal concrete impl for unit testing (no DEX)
contract TestableWorknetManager is WorknetManagerBase {
    constructor() { _disableInitializers(); }

    function initialize(address worknetToken_, bytes32 poolId_, address admin_) external initializer {
        __WorknetManagerBase_init(worknetToken_, poolId_, admin_);
    }

    function _addSingleSidedLiquidity(uint256) internal pure override {}
    function _buybackAndBurn(uint256, uint256) internal pure override {}
}

contract WorknetManagerTest is Test {
    TestableWorknetManager public wmImpl;
    WorknetManagerBase public wm;
    MockAlphaForWM public alpha;
    address public admin = address(this);
    address public merkleAdmin = makeAddr("merkleAdmin");
    address public strategyAdmin = makeAddr("strategyAdmin");
    address public transferAdmin = makeAddr("transferAdmin");
    address public alice = makeAddr("alice");
    address public bob = makeAddr("bob");

    bytes32 constant POOL_ID = bytes32(uint256(42));

    function setUp() public {
        alpha = new MockAlphaForWM();

        wmImpl = new TestableWorknetManager();
        wm = WorknetManagerBase(address(new ERC1967Proxy(
            address(wmImpl),
            abi.encodeCall(TestableWorknetManager.initialize, (address(alpha), POOL_ID, admin))
        )));

        // Grant additional roles
        wm.grantRole(wm.MERKLE_ROLE(), merkleAdmin);
        wm.grantRole(wm.STRATEGY_ROLE(), strategyAdmin);
        wm.grantRole(wm.TRANSFER_ROLE(), transferAdmin);
    }

    // ═══════════════════════════════════════════════
    //  Initialization
    // ═══════════════════════════════════════════════

    function test_initialize() public view {
        assertEq(address(wm.worknetToken()), address(alpha));
        assertEq(wm.poolId(), POOL_ID);
        assertEq(wm.slippageBps(), 500);
        assertFalse(wm.strategyPaused());
        assertTrue(wm.hasRole(wm.DEFAULT_ADMIN_ROLE(), admin));
        assertTrue(wm.hasRole(wm.MERKLE_ROLE(), admin));
        assertTrue(wm.hasRole(wm.STRATEGY_ROLE(), admin));
        assertTrue(wm.hasRole(wm.TRANSFER_ROLE(), admin));
    }

    // ═══════════════════════════════════════════════
    //  Merkle Distribution
    // ═══════════════════════════════════════════════

    function test_setMerkleRoot() public {
        bytes32 root = keccak256("root");
        vm.prank(merkleAdmin);
        wm.setMerkleRoot(0, root);

        assertEq(wm.merkleRoots(0), root);
    }

    function test_setMerkleRoot_zeroRoot_reverts() public {
        vm.prank(merkleAdmin);
        vm.expectRevert(WorknetManagerBase.ZeroRoot.selector);
        wm.setMerkleRoot(0, bytes32(0));
    }

    function test_setMerkleRoot_alreadySet_reverts() public {
        bytes32 root = keccak256("root");
        vm.prank(merkleAdmin);
        wm.setMerkleRoot(0, root);

        vm.prank(merkleAdmin);
        vm.expectRevert(WorknetManagerBase.RootAlreadySet.selector);
        wm.setMerkleRoot(0, keccak256("another"));
    }

    function test_setMerkleRoot_notRole_reverts() public {
        vm.prank(alice);
        vm.expectRevert();
        wm.setMerkleRoot(0, keccak256("root"));
    }

    function test_claim() public {
        uint256 amount = 100e18;
        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(alice, amount))));
        bytes32 root = leaf;

        vm.prank(merkleAdmin);
        wm.setMerkleRoot(0, root);

        // Mock AWPRegistry.resolveRecipient(alice) to return alice
        vm.mockCall(
            wm.awpRegistry(),
            abi.encodeWithSignature("resolveRecipient(address)", alice),
            abi.encode(alice)
        );

        bytes32[] memory proof = new bytes32[](0);
        vm.prank(alice);
        wm.claim(0, amount, proof);

        assertTrue(wm.claimed(0, alice));
        assertEq(alpha.balances(alice), amount);
    }

    function test_claim_alreadyClaimed_reverts() public {
        uint256 amount = 100e18;
        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(alice, amount))));
        bytes32 root = leaf;

        vm.prank(merkleAdmin);
        wm.setMerkleRoot(0, root);

        vm.mockCall(
            wm.awpRegistry(),
            abi.encodeWithSignature("resolveRecipient(address)", alice),
            abi.encode(alice)
        );

        bytes32[] memory proof = new bytes32[](0);
        vm.prank(alice);
        wm.claim(0, amount, proof);

        vm.prank(alice);
        vm.expectRevert(WorknetManagerBase.AlreadyClaimed.selector);
        wm.claim(0, amount, proof);
    }

    function test_claim_noRoot_reverts() public {
        bytes32[] memory proof = new bytes32[](0);
        vm.prank(alice);
        vm.expectRevert(WorknetManagerBase.NoRootForEpoch.selector);
        wm.claim(0, 100e18, proof);
    }

    function test_claim_invalidProof_reverts() public {
        bytes32 root = keccak256("real_root");
        vm.prank(merkleAdmin);
        wm.setMerkleRoot(0, root);

        vm.mockCall(
            wm.awpRegistry(),
            abi.encodeWithSignature("resolveRecipient(address)", alice),
            abi.encode(alice)
        );

        bytes32[] memory proof = new bytes32[](0);
        vm.prank(alice);
        vm.expectRevert(WorknetManagerBase.InvalidProof.selector);
        wm.claim(0, 100e18, proof);
    }

    function test_claim_twoLeafTree() public {
        uint256 amountA = 100e18;
        uint256 amountB = 200e18;

        bytes32 leafA = keccak256(bytes.concat(keccak256(abi.encode(alice, amountA))));
        bytes32 leafB = keccak256(bytes.concat(keccak256(abi.encode(bob, amountB))));

        // Standard sorted Merkle root
        bytes32 root;
        if (leafA <= leafB) {
            root = keccak256(abi.encodePacked(leafA, leafB));
        } else {
            root = keccak256(abi.encodePacked(leafB, leafA));
        }

        vm.prank(merkleAdmin);
        wm.setMerkleRoot(1, root);

        vm.mockCall(
            wm.awpRegistry(),
            abi.encodeWithSignature("resolveRecipient(address)", alice),
            abi.encode(alice)
        );
        vm.mockCall(
            wm.awpRegistry(),
            abi.encodeWithSignature("resolveRecipient(address)", bob),
            abi.encode(bob)
        );

        // Alice claims with proof = [leafB]
        bytes32[] memory proofA = new bytes32[](1);
        proofA[0] = leafB;
        vm.prank(alice);
        wm.claim(1, amountA, proofA);
        assertEq(alpha.balances(alice), amountA);

        // Bob claims with proof = [leafA]
        bytes32[] memory proofB = new bytes32[](1);
        proofB[0] = leafA;
        vm.prank(bob);
        wm.claim(1, amountB, proofB);
        assertEq(alpha.balances(bob), amountB);
    }

    function test_isClaimed() public view {
        assertFalse(wm.isClaimed(0, alice));
    }

    // ═══════════════════════════════════════════════
    //  Strategy
    // ═══════════════════════════════════════════════

    function test_setStrategy() public {
        vm.prank(strategyAdmin);
        wm.setStrategy(WorknetManagerBase.AWPStrategy.BuybackBurn);
        assertEq(uint8(wm.currentStrategy()), uint8(WorknetManagerBase.AWPStrategy.BuybackBurn));
    }

    function test_setStrategy_notRole_reverts() public {
        vm.prank(alice);
        vm.expectRevert();
        wm.setStrategy(WorknetManagerBase.AWPStrategy.AddLiquidity);
    }

    function test_executeStrategy_reserve_noOp() public {
        // Default strategy is Reserve
        vm.prank(strategyAdmin);
        wm.executeStrategy(100e18, 0);
        // No revert = success (Reserve does nothing)
    }

    function test_executeStrategy_paused_reverts() public {
        wm.setStrategyPaused(true);

        vm.prank(strategyAdmin);
        vm.expectRevert(WorknetManagerBase.StrategyIsPaused.selector);
        wm.executeStrategy(100e18, 0);
    }

    function test_executeStrategy_zeroAmount_reverts() public {
        vm.prank(strategyAdmin);
        vm.expectRevert(WorknetManagerBase.ZeroAmount.selector);
        wm.executeStrategy(0, 0);
    }

    // ═══════════════════════════════════════════════
    //  ERC1363 Receiver
    // ═══════════════════════════════════════════════

    function test_onTransferReceived_nonAWP_noop() public {
        // Call from non-AWP token — should just return selector
        vm.prank(address(alpha));
        bytes4 ret = wm.onTransferReceived(address(0), address(0), 100e18, "");
        assertEq(ret, IERC1363Receiver.onTransferReceived.selector);
    }

    function test_onTransferReceived_fromAWP() public {
        // Mock call as AWP token
        vm.prank(wm.awpToken());
        bytes4 ret = wm.onTransferReceived(address(0), address(0), 100e18, "");
        assertEq(ret, IERC1363Receiver.onTransferReceived.selector);
    }

    // ═══════════════════════════════════════════════
    //  Token Transfer (TRANSFER_ROLE)
    // ═══════════════════════════════════════════════

    function test_transferToken() public {
        // Give wm some mock tokens
        vm.mockCall(
            address(0x1234),
            abi.encodeWithSignature("transfer(address,uint256)", alice, 50e18),
            abi.encode(true)
        );

        vm.prank(transferAdmin);
        wm.transferToken(address(0x1234), alice, 50e18);
    }

    function test_transferToken_notRole_reverts() public {
        vm.prank(alice);
        vm.expectRevert();
        wm.transferToken(address(0x1234), alice, 50e18);
    }

    function test_batchTransferToken_arrayMismatch_reverts() public {
        address[] memory recipients = new address[](2);
        uint256[] memory amounts = new uint256[](1);

        vm.prank(transferAdmin);
        vm.expectRevert(WorknetManagerBase.ArrayLengthMismatch.selector);
        wm.batchTransferToken(address(0x1234), recipients, amounts);
    }

    // ═══════════════════════════════════════════════
    //  Configuration
    // ═══════════════════════════════════════════════

    function test_setSlippageTolerance() public {
        vm.prank(strategyAdmin);
        wm.setSlippageTolerance(300);
        assertEq(wm.slippageBps(), 300);
    }

    function test_setSlippageTolerance_zero_reverts() public {
        vm.prank(strategyAdmin);
        vm.expectRevert(WorknetManagerBase.InvalidSlippage.selector);
        wm.setSlippageTolerance(0);
    }

    function test_setSlippageTolerance_over5000_reverts() public {
        vm.prank(strategyAdmin);
        vm.expectRevert(WorknetManagerBase.InvalidSlippage.selector);
        wm.setSlippageTolerance(5001);
    }

    function test_setStrategyPaused() public {
        wm.setStrategyPaused(true);
        assertTrue(wm.strategyPaused());

        wm.setStrategyPaused(false);
        assertFalse(wm.strategyPaused());
    }

    function test_setStrategyPaused_notAdmin_reverts() public {
        vm.prank(alice);
        vm.expectRevert();
        wm.setStrategyPaused(true);
    }

    function test_setMinStrategyAmount() public {
        vm.prank(strategyAdmin);
        wm.setMinStrategyAmount(1e18);
        assertEq(wm.minStrategyAmount(), 1e18);
    }
}
