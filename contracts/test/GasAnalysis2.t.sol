// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;
import {Test, console} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {veAWP} from "../src/core/veAWP.sol";
import {AWPAllocator} from "../src/core/AWPAllocator.sol";

contract GasAnalysis2 is Test {
    address constant AWP = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant VEAWP = 0x0000b534C63D78212f1BDCc315165852793A00A8;
    address constant ALLOCATOR = 0x0000D6BB5e040E35081b3AaF59DD71b21C9800AA;

    function setUp() public {
        vm.createSelectFork(vm.envString("BASE_RPC_URL"));
    }

    function test_nft_transfer_hook_cost() public {
        // 创建一个 position 然后测量 transfer hook 的 external call 开销
        address alice = makeAddr("alice");
        address bob = makeAddr("bob");
        deal(AWP, alice, 10_000 ether);
        
        vm.startPrank(alice);
        IERC20(AWP).approve(VEAWP, 10_000 ether);
        uint256 tokenId = veAWP(VEAWP).deposit(10_000 ether, 90 days);
        vm.stopPrank();

        // 测量 NFT transfer（包含 _update hook → AWPAllocator.userTotalAllocated 外部调用）
        uint256 g0 = gasleft();
        vm.prank(alice);
        veAWP(VEAWP).transferFrom(alice, bob, tokenId);
        uint256 transferGas = g0 - gasleft();
        console.log("NFT transfer (with hook):   ", transferGas);

        // 对比：测量纯 AWPAllocator.userTotalAllocated 调用开销
        g0 = gasleft();
        AWPAllocator(ALLOCATOR).userTotalAllocated(alice);
        uint256 allocCallGas = g0 - gasleft();
        console.log("AWPAllocator.userTotal:     ", allocCallGas);

        // 直接存款（无 helper）
        deal(AWP, bob, 10_000 ether);
        vm.startPrank(bob);
        IERC20(AWP).approve(VEAWP, 10_000 ether);
        g0 = gasleft();
        veAWP(VEAWP).deposit(10_000 ether, 90 days);
        uint256 directDeposit = g0 - gasleft();
        console.log("Direct deposit (no helper): ", directDeposit);
        vm.stopPrank();
    }
}
