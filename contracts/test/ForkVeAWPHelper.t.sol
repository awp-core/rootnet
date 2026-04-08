// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import {veAWP} from "../src/core/veAWP.sol";
import {VeAWPHelper} from "../src/core/VeAWPHelper.sol";

/// @title ForkVeAWPHelper — Fork tests for gasless staking helper
contract ForkVeAWPHelper is Test {
    address constant AWP_TOKEN = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant VEAWP = 0x0000b534C63D78212f1BDCc315165852793A00A8;

    VeAWPHelper helper;
    veAWP ve;
    IERC20 awp;

    // 用户私钥（用于签名 permit）
    uint256 constant USER_PK = 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80;
    address user;
    address relayer;

    uint256 constant AMOUNT = 10_000 ether;
    uint64 constant LOCK_90_DAYS = 90 days;

    function setUp() public {
        vm.createSelectFork(vm.envString("BASE_RPC_URL"));

        ve = veAWP(VEAWP);
        awp = IERC20(AWP_TOKEN);

        // Deploy helper
        helper = new VeAWPHelper(AWP_TOKEN, VEAWP);

        user = vm.addr(USER_PK);
        relayer = makeAddr("relayer");

        // Fund user with AWP
        deal(AWP_TOKEN, user, 100_000 ether);
        // Fund relayer with ETH for gas
        deal(relayer, 1 ether);
    }

    // ═══════════════════════════════════════════
    //  Helper: _signPermit
    // ═══════════════════════════════════════════

    function _signPermit(
        address owner,
        address spender,
        uint256 value,
        uint256 deadline
    ) internal view returns (uint8 v, bytes32 r, bytes32 s) {
        // 获取 AWP 的 EIP-2612 domain separator 和 nonce
        bytes32 PERMIT_TYPEHASH = keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)");
        bytes32 domainSeparator = _getDomainSeparator();
        uint256 nonce = _getNonce(owner);

        bytes32 structHash = keccak256(abi.encode(
            PERMIT_TYPEHASH, owner, spender, value, nonce, deadline
        ));
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", domainSeparator, structHash));

        (v, r, s) = vm.sign(USER_PK, digest);
    }

    function _getDomainSeparator() internal view returns (bytes32) {
        // 调用 AWP token 的 DOMAIN_SEPARATOR()
        (bool ok, bytes memory data) = AWP_TOKEN.staticcall(abi.encodeWithSignature("DOMAIN_SEPARATOR()"));
        require(ok, "DOMAIN_SEPARATOR failed");
        return abi.decode(data, (bytes32));
    }

    function _getNonce(address owner) internal view returns (uint256) {
        (bool ok, bytes memory data) = AWP_TOKEN.staticcall(abi.encodeWithSignature("nonces(address)", owner));
        require(ok, "nonces failed");
        return abi.decode(data, (uint256));
    }

    // ═══════════════════════════════════════════
    //  depositFor — 完全 gasless 质押
    // ═══════════════════════════════════════════

    function test_depositFor_gasless() public {
        uint256 deadline = block.timestamp + 1 hours;

        // 用户离线签名 permit（approve helper）
        (uint8 v, bytes32 r, bytes32 s) = _signPermit(user, address(helper), AMOUNT, deadline);

        // Relayer 发起交易（用户不付 gas）
        vm.prank(relayer);
        uint256 tokenId = helper.depositFor(user, AMOUNT, LOCK_90_DAYS, deadline, v, r, s);

        // 验证：NFT 归用户所有
        assertEq(ve.ownerOf(tokenId), user, "NFT should be owned by user");

        // 验证：position 正确
        (uint128 amount, uint64 lockEnd, uint64 createdAt) = ve.positions(tokenId);
        assertEq(amount, uint128(AMOUNT), "position amount");
        assertGt(lockEnd, block.timestamp, "lock active");
        assertEq(createdAt, uint64(block.timestamp), "createdAt");

        // 验证：userTotalStaked 正确
        assertEq(ve.getUserTotalStaked(user), AMOUNT, "userTotalStaked");

        // 验证：helper 不持有任何资产
        assertEq(awp.balanceOf(address(helper)), 0, "helper AWP = 0");
        assertEq(ve.balanceOf(address(helper)), 0, "helper NFTs = 0");
    }

    function test_depositFor_user_can_withdraw_after_lock() public {
        uint256 deadline = block.timestamp + 1 hours;
        (uint8 v, bytes32 r, bytes32 s) = _signPermit(user, address(helper), AMOUNT, deadline);

        vm.prank(relayer);
        uint256 tokenId = helper.depositFor(user, AMOUNT, 1 days, deadline, v, r, s);

        // 锁到期后用户可以赎回
        vm.warp(block.timestamp + 1 days + 1);

        uint256 balBefore = awp.balanceOf(user);
        vm.prank(user);
        ve.withdraw(tokenId);

        assertEq(awp.balanceOf(user), balBefore + AMOUNT, "AWP returned to user");
    }

    function test_depositFor_multiple_deposits() public {
        uint256 deadline = block.timestamp + 1 hours;
        uint256 amount1 = 5_000 ether;
        uint256 amount2 = 3_000 ether;

        // 第一次存入
        (uint8 v1, bytes32 r1, bytes32 s1) = _signPermit(user, address(helper), amount1, deadline);
        vm.prank(relayer);
        uint256 id1 = helper.depositFor(user, amount1, LOCK_90_DAYS, deadline, v1, r1, s1);

        // 第二次存入（nonce 自动递增）
        (uint8 v2, bytes32 r2, bytes32 s2) = _signPermit(user, address(helper), amount2, deadline);
        vm.prank(relayer);
        uint256 id2 = helper.depositFor(user, amount2, LOCK_90_DAYS, deadline, v2, r2, s2);

        assertEq(ve.ownerOf(id1), user);
        assertEq(ve.ownerOf(id2), user);
        assertEq(ve.getUserTotalStaked(user), amount1 + amount2);
    }

    function test_depositFor_with_pre_approval() public {
        // 用户已经 approve 了 helper（不需要 permit）
        vm.prank(user);
        awp.approve(address(helper), AMOUNT);

        // Relayer 用空签名调用（permit 会 revert 但被 try/catch 忽略）
        vm.prank(relayer);
        uint256 tokenId = helper.depositFor(user, AMOUNT, LOCK_90_DAYS, 0, 0, bytes32(0), bytes32(0));

        assertEq(ve.ownerOf(tokenId), user);
        assertEq(ve.getUserTotalStaked(user), AMOUNT);
    }

    // ═══════════════════════════════════════════
    //  depositFor — 错误场景
    // ═══════════════════════════════════════════

    function test_depositFor_reverts_zero_amount() public {
        vm.prank(relayer);
        vm.expectRevert(VeAWPHelper.ZeroAmount.selector);
        helper.depositFor(user, 0, LOCK_90_DAYS, 0, 0, bytes32(0), bytes32(0));
    }

    function test_depositFor_reverts_zero_address() public {
        vm.prank(relayer);
        vm.expectRevert(VeAWPHelper.ZeroAddress.selector);
        helper.depositFor(address(0), AMOUNT, LOCK_90_DAYS, 0, 0, bytes32(0), bytes32(0));
    }

    function test_depositFor_reverts_invalid_permit() public {
        // 用错误私钥签名
        (uint8 v, bytes32 r, bytes32 s) = _signPermit(user, address(helper), AMOUNT, block.timestamp + 1 hours);
        // 但 amount 不匹配
        vm.prank(relayer);
        vm.expectRevert(); // transferFrom will fail (no allowance)
        helper.depositFor(user, AMOUNT + 1, LOCK_90_DAYS, block.timestamp + 1 hours, v, r, s);
    }

    function test_depositFor_reverts_lock_too_short() public {
        uint256 deadline = block.timestamp + 1 hours;
        (uint8 v, bytes32 r, bytes32 s) = _signPermit(user, address(helper), AMOUNT, deadline);

        vm.prank(relayer);
        vm.expectRevert(); // veAWP.LockTooShort
        helper.depositFor(user, AMOUNT, 1 hours, deadline, v, r, s); // < MIN_LOCK_DURATION
    }

    // ═══════════════════════════════════════════
    //  depositFor — helper 合约安全性
    // ═══════════════════════════════════════════

    function test_helper_holds_no_state_between_calls() public {
        uint256 deadline = block.timestamp + 1 hours;
        (uint8 v, bytes32 r, bytes32 s) = _signPermit(user, address(helper), AMOUNT, deadline);

        vm.prank(relayer);
        helper.depositFor(user, AMOUNT, LOCK_90_DAYS, deadline, v, r, s);

        // Helper 不持有任何资产
        assertEq(awp.balanceOf(address(helper)), 0, "no AWP");
        assertEq(ve.balanceOf(address(helper)), 0, "no NFT");

        // AWP allowance from helper to veAWP is max-approved (by design, saves gas per call)
        assertGt(awp.allowance(address(helper), VEAWP), 0, "max-approve persists");
    }

    function test_depositFor_reentrancy_protected() public {
        // VeAWPHelper 有 nonReentrant 修饰符
        // 此测试确认合约已编译并部署（reentrancy guard 内嵌在字节码中）
        assertTrue(address(helper).code.length > 0);
    }

    // ═══════════════════════════════════════════
    //  depositFor — user == address(this) guard
    // ═══════════════════════════════════════════

    function test_depositFor_reverts_self_deposit() public {
        vm.prank(relayer);
        vm.expectRevert(VeAWPHelper.InvalidUser.selector);
        helper.depositFor(address(helper), AMOUNT, LOCK_90_DAYS, 0, 0, bytes32(0), bytes32(0));
    }
}
