// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {LPManager} from "../src/core/LPManager.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AlphaToken} from "../src/token/AlphaToken.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @title LPManagerForkTest — BSC fork test verifying PancakeSwap V4 CL integration
contract LPManagerForkTest is Test {
    LPManager public lpManager;
    AWPToken public awpToken;
    AlphaToken public alphaToken;

    address public deployer = makeAddr("deployer");
    address public rootNet = makeAddr("rootNet");
    address public user = makeAddr("user");

    // BSC PancakeSwap V4 addresses
    address constant CL_POOL_MANAGER = 0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b;
    address constant CL_POSITION_MANAGER = 0x55f4c8abA71A1e923edC303eb4fEfF14608cC226;
    address constant PERMIT2 = 0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768;

    uint256 constant AWP_LP_AMOUNT = 1_000_000 * 1e18;
    uint256 constant ALPHA_LP_AMOUNT = 100_000_000 * 1e18;

    function setUp() public {
        vm.startPrank(deployer);

        // Deploy AWPToken (constructor mints 200M to deployer)
        awpToken = new AWPToken("AWP Token", "AWP", deployer);

        // Deploy AlphaToken and initialize
        // Create AlphaToken clone via Clones (implementation contract forbids direct initialize)
        AlphaToken alphaImpl = new AlphaToken();
        address alphaClone = Clones.clone(address(alphaImpl));
        alphaToken = AlphaToken(alphaClone);
        alphaToken.initialize("Test Alpha", "TALPHA", 1, deployer);
        // initialize automatically sets deployer as minter
        alphaToken.mint(deployer, ALPHA_LP_AMOUNT);

        // Deploy LPManager (5-parameter constructor)
        lpManager = new LPManager(rootNet, CL_POOL_MANAGER, CL_POSITION_MANAGER, PERMIT2, address(awpToken));

        // Transfer tokens to LPManager (simulating RootNet behavior)
        awpToken.transfer(address(lpManager), AWP_LP_AMOUNT);
        alphaToken.transfer(address(lpManager), ALPHA_LP_AMOUNT);

        vm.stopPrank();
    }

    /// @notice Create LP pool and add liquidity, verify poolId is non-zero and tokens have been transferred out
    function test_createPoolAndAddLiquidity() public {
        uint256 awpBefore = awpToken.balanceOf(address(lpManager));
        uint256 alphaBefore = alphaToken.balanceOf(address(lpManager));

        vm.prank(rootNet);
        (bytes32 poolId, uint256 lpTokenId) = lpManager.createPoolAndAddLiquidity(
            address(alphaToken), AWP_LP_AMOUNT, ALPHA_LP_AMOUNT
        );

        // Pool ID should not be zero
        assertTrue(poolId != bytes32(0), "poolId should not be zero");

        // LP NFT tokenId should have a value
        assertTrue(lpTokenId > 0, "lpTokenId should be > 0");

        // Tokens should have left LPManager (transferred to PancakeSwap)
        uint256 awpAfter = awpToken.balanceOf(address(lpManager));
        uint256 alphaAfter = alphaToken.balanceOf(address(lpManager));
        assertTrue(awpAfter < awpBefore, "AWP should have left LPManager");
        assertTrue(alphaAfter < alphaBefore, "Alpha should have left LPManager");

        // Pool ID has been recorded
        assertEq(lpManager.alphaTokenToPoolId(address(alphaToken)), poolId, "poolId should be stored");
    }

    /// @notice Creating a pool for the same Alpha token twice should revert
    function test_revertsDoubleCreate() public {
        vm.prank(rootNet);
        lpManager.createPoolAndAddLiquidity(address(alphaToken), AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        vm.prank(rootNet);
        vm.expectRevert(LPManager.PoolAlreadyExists.selector);
        lpManager.createPoolAndAddLiquidity(address(alphaToken), AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);
    }

    /// @notice Non-RootNet calls should revert
    function test_revertsNonRootNet() public {
        vm.prank(user);
        vm.expectRevert(LPManager.NotRootNet.selector);
        lpManager.createPoolAndAddLiquidity(address(alphaToken), AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);
    }
}
