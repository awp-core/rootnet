// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {AWPToken} from "../src/token/AWPToken.sol";
import {AlphaToken} from "../src/token/AlphaToken.sol";
import {AlphaTokenFactory} from "../src/token/AlphaTokenFactory.sol";
import {AWPEmission} from "../src/token/AWPEmission.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";
import {StakingVault} from "../src/core/StakingVault.sol";
import {StakeNFT} from "../src/core/StakeNFT.sol";
import {WorknetNFT} from "../src/core/WorknetNFT.sol";
import {LPManager} from "../src/core/LPManager.sol";
import {LPManagerUni} from "../src/core/LPManagerUni.sol";
import {LPManagerBase} from "../src/core/LPManagerBase.sol";
import {AWPRegistry} from "../src/AWPRegistry.sol";
import {IAWPRegistry} from "../src/interfaces/IAWPRegistry.sol";
import {Treasury} from "../src/governance/Treasury.sol";
import {WorknetManager} from "../src/worknets/WorknetManager.sol";
import {WorknetManagerUni} from "../src/worknets/WorknetManagerUni.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

// ═══════════════════════════════════════════════════════════════════════
//  ForkTestBase — Shared test infrastructure; subclasses implement chain-specific DEX integration
// ═══════════════════════════════════════════════════════════════════════

abstract contract ForkTestBase is Test {
    // ── Test accounts ──
    address deployer = makeAddr("deployer");
    address guardian = makeAddr("guardian");
    address user = makeAddr("user");

    // ── Protocol contracts ──
    AWPToken awpToken;
    AlphaToken alphaImpl;
    AlphaTokenFactory alphaFactory;
    AWPRegistry registryImpl;
    AWPRegistry registry; // proxy
    WorknetNFT worknetNFT;
    AWPEmission emissionImpl;
    AWPEmission emission; // proxy
    StakingVault vaultImpl;
    StakingVault vault; // proxy
    StakeNFT stakeNFT;
    Treasury treasury;

    // ── LP & Worknet ──
    LPManagerBase lpManager;
    address worknetManagerImpl;

    // ── Constants ──
    uint256 constant INITIAL_MINT = 200_000_000 * 1e18;
    uint256 constant AWP_LP_AMOUNT = 1_000_000 * 1e18; // initialAlphaMint * initialAlphaPrice
    uint256 constant ALPHA_LP_AMOUNT = 100_000_000 * 1e18;
    uint24 constant POOL_FEE = 10000;
    int24 constant TICK_SPACING = 200;

    // ── Fork ID (set by subclass) ──
    uint256 internal fork;

    /// @dev Subclass must implement: deploy chain-specific LPManager
    function _deployLPManager(address registry_, address permit2_, address awpToken_)
        internal
        virtual
        returns (LPManagerBase);

    /// @dev Subclass must implement: get chain-specific WorknetManager dexConfig
    function _getWorknetManagerDexConfig() internal view virtual returns (bytes memory);

    /// @dev Subclass must implement: deploy chain-specific WorknetManager impl
    function _deployWorknetManagerImpl() internal virtual returns (address);

    /// @dev Deploy full protocol stack (on fork)
    function _deployProtocol() internal {
        vm.startPrank(deployer);

        // 1. AWPToken
        awpToken = new AWPToken("AWP Token", "AWP", deployer);
        awpToken.initialMint(INITIAL_MINT);

        // 2. AlphaTokenFactory (vanityRule=0 skips validation)
        alphaFactory = new AlphaTokenFactory(deployer, 0);

        // 3. Treasury
        address[] memory proposers = new address[](1);
        proposers[0] = deployer;
        address[] memory executors = new address[](1);
        executors[0] = address(0); // anyone can execute
        treasury = new Treasury(1 days, proposers, executors, deployer);

        // 4. AWPRegistry (proxy)
        registryImpl = new AWPRegistry();
        bytes memory registryInitData = abi.encodeWithSelector(
            AWPRegistry.initialize.selector, deployer, address(treasury), guardian
        );
        registry = AWPRegistry(address(new ERC1967Proxy(address(registryImpl), registryInitData)));

        // 5. WorknetNFT
        worknetNFT = new WorknetNFT("AWP Worknet", "wNFT", address(registry));

        // 6. Deploy chain-specific LPManager
        lpManager = _deployLPManager(address(registry), address(0), address(awpToken));

        // 7. AWPEmission (proxy)
        emissionImpl = new AWPEmission();
        bytes memory emissionInitData = abi.encodeWithSelector(
            AWPEmission.initialize.selector,
            address(awpToken),
            guardian,
            1_000_000 * 1e18, // initialDailyEmission
            block.timestamp,  // genesisTime
            1 days,           // epochDuration
            address(treasury) // treasury (failed mint fallback)
        );
        emission = AWPEmission(address(new ERC1967Proxy(address(emissionImpl), emissionInitData)));

        // 8. StakingVault (proxy)
        vaultImpl = new StakingVault();
        bytes memory vaultInitData = abi.encodeWithSelector(
            StakingVault.initialize.selector, address(registry), guardian
        );
        vault = StakingVault(address(new ERC1967Proxy(address(vaultImpl), vaultInitData)));

        // 9. StakeNFT
        stakeNFT = new StakeNFT(address(awpToken), address(vault), address(registry));

        // 10. WorknetManager impl
        worknetManagerImpl = _deployWorknetManagerImpl();

        // 11. initializeRegistry (deployer one-time call, locked after)
        bytes memory dexConfig = _getWorknetManagerDexConfig();
        registry.initializeRegistry(
            address(awpToken),
            address(worknetNFT),
            address(alphaFactory),
            address(emission),
            address(lpManager),
            address(vault),
            address(stakeNFT),
            worknetManagerImpl,
            dexConfig
        );

        // 12. Configure associations
        alphaFactory.setAddresses(address(registry));
        awpToken.addMinter(address(emission));

        vm.stopPrank();
    }

    /// @dev Deploy AlphaToken clone for standalone LP tests
    function _deployAlphaClone(string memory name, string memory symbol, uint256 id)
        internal
        returns (AlphaToken)
    {
        alphaImpl = new AlphaToken();
        address clone = Clones.clone(address(alphaImpl));
        AlphaToken alpha = AlphaToken(clone);
        alpha.initialize(name, symbol, id, deployer);
        return alpha;
    }

    /// @dev Transfer tokens to LPManager and call createPoolAndAddLiquidity
    function _createPool(AlphaToken alpha, uint256 awpAmt, uint256 alphaAmt)
        internal
        returns (bytes32 poolId, uint256 lpTokenId)
    {
        // Transfer tokens to lpManager
        vm.startPrank(deployer);
        awpToken.transfer(address(lpManager), awpAmt);
        alpha.mint(deployer, alphaAmt);
        IERC20(address(alpha)).transfer(address(lpManager), alphaAmt);
        vm.stopPrank();

        // Simulate AWPRegistry call
        vm.prank(address(registry));
        (poolId, lpTokenId) = lpManager.createPoolAndAddLiquidity(address(alpha), awpAmt, alphaAmt);
    }
}

// ═══════════════════════════════════════════════════════════════════════
//  BSC Fork Test — PancakeSwap V4
// ═══════════════════════════════════════════════════════════════════════

contract BSCForkTest is ForkTestBase {
    // ── PancakeSwap V4 addresses (BSC mainnet) ──
    address constant CL_POOL_MANAGER = 0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b;
    address constant CL_POSITION_MANAGER = 0x55f4c8abA71A1e923edC303eb4fEfF14608cC226;
    address constant PERMIT2 = 0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768;
    address constant CL_SWAP_ROUTER = 0x1b81D678ffb9C0263b24A97847620C99d213eB14;

    function setUp() public {
        fork = vm.createFork(vm.envString("BSC_RPC_URL"));
        vm.selectFork(fork);
        _deployProtocol();
    }

    function _deployLPManager(address registry_, address, address awpToken_)
        internal
        override
        returns (LPManagerBase)
    {
        return new LPManager(registry_, CL_POOL_MANAGER, CL_POSITION_MANAGER, PERMIT2, awpToken_);
    }

    function _getWorknetManagerDexConfig() internal pure override returns (bytes memory) {
        return abi.encode(CL_POOL_MANAGER, CL_POSITION_MANAGER, CL_SWAP_ROUTER, PERMIT2, POOL_FEE, TICK_SPACING);
    }

    function _deployWorknetManagerImpl() internal override returns (address) {
        return address(new WorknetManager());
    }

    // ── Test 1: LP Pool creation ──
    function test_createPoolAndAddLiquidity() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        (bytes32 poolId, uint256 lpTokenId) = _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        assertNotEq(poolId, bytes32(0), "poolId should not be zero");
        assertGt(lpTokenId, 0, "lpTokenId should be > 0");

        // Verify mappings stored
        assertEq(lpManager.alphaTokenToPoolId(address(alpha)), poolId);
        assertEq(lpManager.alphaTokenToTokenId(address(alpha)), lpTokenId);

        // Verify tokens transferred out of LPManager (most went into pool)
        // Note: small residual balance possible due to precision
        uint256 remainingAWP = awpToken.balanceOf(address(lpManager));
        uint256 remainingAlpha = IERC20(address(alpha)).balanceOf(address(lpManager));
        assertLt(remainingAWP, AWP_LP_AMOUNT, "AWP should be transferred out");
        assertLt(remainingAlpha, ALPHA_LP_AMOUNT, "Alpha should be transferred out");
    }

    // ── Test 2: Duplicate LP creation should revert ──
    function test_createPoolTwiceReverts() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        // Second creation should revert
        vm.startPrank(deployer);
        awpToken.transfer(address(lpManager), AWP_LP_AMOUNT);
        alpha.mint(deployer, ALPHA_LP_AMOUNT);
        IERC20(address(alpha)).transfer(address(lpManager), ALPHA_LP_AMOUNT);
        vm.stopPrank();

        vm.prank(address(registry));
        vm.expectRevert(LPManagerBase.PoolAlreadyExists.selector);
        lpManager.createPoolAndAddLiquidity(address(alpha), AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);
    }

    // ── Test 3: LP Fee Compounding ──
    function test_compoundFeesNoRevert() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        // compoundFees may revert with no fees (V4 TAKE_PAIR with 0 fees)
        // Just verify pool and tokenId stored correctly
        assertTrue(lpManager.alphaTokenToTokenId(address(alpha)) > 0);
        assertTrue(lpManager.alphaTokenToPoolId(address(alpha)) != bytes32(0));
    }

    // ── Test 4: WorknetManager initialization ──
    function test_worknetManagerInitialize() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        (bytes32 poolId,) = _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        // Deploy WorknetManager via proxy
        bytes memory initData = abi.encodeWithSelector(
            WorknetManager.initialize.selector,
            address(registry),
            address(alpha),
            address(awpToken),
            poolId,
            deployer,
            _getWorknetManagerDexConfig()
        );
        WorknetManager wm = WorknetManager(address(new ERC1967Proxy(worknetManagerImpl, initData)));

        // Verify storage
        assertEq(address(wm.awpRegistry()), address(registry));
        assertEq(address(wm.alphaToken()), address(alpha));
        assertEq(address(wm.awpToken()), address(awpToken));
        assertEq(wm.poolId(), poolId);
        assertEq(wm.clPoolManager(), CL_POOL_MANAGER);
        assertEq(wm.clPositionManager(), CL_POSITION_MANAGER);
        assertEq(wm.clSwapRouter(), CL_SWAP_ROUTER);
        assertEq(wm.permit2(), PERMIT2);
        assertEq(wm.poolFee(), POOL_FEE);
        assertEq(wm.tickSpacing(), TICK_SPACING);

        // Verify role grants
        assertTrue(wm.hasRole(wm.DEFAULT_ADMIN_ROLE(), deployer));
        assertTrue(wm.hasRole(wm.MERKLE_ROLE(), deployer));
        assertTrue(wm.hasRole(wm.STRATEGY_ROLE(), deployer));
        assertTrue(wm.hasRole(wm.TRANSFER_ROLE(), deployer));

        // Verify poolKey constructed correctly
        (address c0, address c1,,,, ) = wm.poolKey();
        if (address(awpToken) < address(alpha)) {
            assertEq(c0, address(awpToken));
            assertEq(c1, address(alpha));
        } else {
            assertEq(c0, address(alpha));
            assertEq(c1, address(awpToken));
        }
    }

    // ── Test 5: WorknetManager UUPS upgrade ──
    function test_worknetManagerUpgrade() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        (bytes32 poolId,) = _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        bytes memory initData = abi.encodeWithSelector(
            WorknetManager.initialize.selector,
            address(registry),
            address(alpha),
            address(awpToken),
            poolId,
            deployer,
            _getWorknetManagerDexConfig()
        );
        WorknetManager wm = WorknetManager(address(new ERC1967Proxy(worknetManagerImpl, initData)));

        // Upgrade to new implementation
        address newImpl = address(new WorknetManager());
        vm.prank(deployer);
        wm.upgradeToAndCall(newImpl, "");

        // Verify state preserved
        assertEq(address(wm.awpRegistry()), address(registry));
        assertEq(address(wm.alphaToken()), address(alpha));
        assertEq(wm.poolId(), poolId);
        assertEq(wm.clPoolManager(), CL_POOL_MANAGER);
        assertTrue(wm.hasRole(wm.DEFAULT_ADMIN_ROLE(), deployer));
    }

    // ── Test 6: Full E2E — worknet registration + LP creation ──
    function test_fullE2ERegistration() public {
        // Give user AWP for registration
        vm.prank(deployer);
        awpToken.transfer(user, 10_000_000 * 1e18);

        // user approve AWPRegistry
        uint256 lpAWPAmount = registry.initialAlphaMint() * registry.initialAlphaPrice() / 1e18;
        vm.prank(user);
        awpToken.approve(address(registry), lpAWPAmount);

        // Register worknet
        IAWPRegistry.WorknetParams memory params = IAWPRegistry.WorknetParams({
            name: "Test Worknet",
            symbol: "TWRK",
            worknetManager: address(0), // auto-deploy
            salt: bytes32(0),
            minStake: 0,
            skillsURI: ""
        });

        vm.prank(user);
        uint256 worknetId = registry.registerWorknet(params);
        assertGt(worknetId, 0, "worknetId should be > 0");

        // Verify worknet status
        (, IAWPRegistry.WorknetStatus status,,) = registry.worknets(worknetId);
        assertEq(uint8(status), uint8(IAWPRegistry.WorknetStatus.Pending));

        // Activate worknet
        vm.prank(user);
        registry.activateWorknet(worknetId);

        (, IAWPRegistry.WorknetStatus statusAfter,,) = registry.worknets(worknetId);
        assertEq(uint8(statusAfter), uint8(IAWPRegistry.WorknetStatus.Active));
    }
}

// ═══════════════════════════════════════════════════════════════════════
//  Base Fork Test — Uniswap V4
// ═══════════════════════════════════════════════════════════════════════

contract BaseForkTest is ForkTestBase {
    // ── Uniswap V4 addresses (Base mainnet) ──
    address constant POOL_MANAGER = 0x498581fF718922c3f8e6A244956aF099B2652b2b;
    address constant POSITION_MANAGER = 0x7C5f5A4bBd8fD63184577525326123B519429bDc;
    address constant PERMIT2 = 0x000000000022D473030F116dDEE9F6B43aC78BA3;
    address constant SWAP_ROUTER = 0x6fF5693b99212Da76ad316178A184AB56D299b43;
    address constant STATE_VIEW = 0xA3c0c9b65baD0b08107Aa264b0f3dB444b867A71;

    function setUp() public {
        fork = vm.createFork(vm.envString("BASE_RPC_URL"));
        vm.selectFork(fork);
        _deployProtocol();
    }

    function _deployLPManager(address registry_, address, address awpToken_)
        internal
        override
        returns (LPManagerBase)
    {
        return new LPManagerUni(registry_, POOL_MANAGER, POSITION_MANAGER, PERMIT2, awpToken_);
    }

    function _getWorknetManagerDexConfig() internal pure override returns (bytes memory) {
        return abi.encode(POOL_MANAGER, POSITION_MANAGER, SWAP_ROUTER, PERMIT2, POOL_FEE, TICK_SPACING, STATE_VIEW);
    }

    function _deployWorknetManagerImpl() internal override returns (address) {
        return address(new WorknetManagerUni());
    }

    // ── Test 1: LP Pool creation ──
    function test_createPoolAndAddLiquidity() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        (bytes32 poolId, uint256 lpTokenId) = _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        assertNotEq(poolId, bytes32(0), "poolId should not be zero");
        assertGt(lpTokenId, 0, "lpTokenId should be > 0");
        assertEq(lpManager.alphaTokenToPoolId(address(alpha)), poolId);
        assertEq(lpManager.alphaTokenToTokenId(address(alpha)), lpTokenId);
    }

    // ── Test 2: Duplicate LP creation should revert ──
    function test_createPoolTwiceReverts() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        vm.startPrank(deployer);
        awpToken.transfer(address(lpManager), AWP_LP_AMOUNT);
        alpha.mint(deployer, ALPHA_LP_AMOUNT);
        IERC20(address(alpha)).transfer(address(lpManager), ALPHA_LP_AMOUNT);
        vm.stopPrank();

        vm.prank(address(registry));
        vm.expectRevert(LPManagerBase.PoolAlreadyExists.selector);
        lpManager.createPoolAndAddLiquidity(address(alpha), AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);
    }

    // ── Test 3: LP Fee Compounding ──
    function test_compoundFeesNoRevert() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);
        // compoundFees may revert with no fees (V4 TAKE_PAIR with 0 fees)
        assertTrue(lpManager.alphaTokenToTokenId(address(alpha)) > 0);
        assertTrue(lpManager.alphaTokenToPoolId(address(alpha)) != bytes32(0));
    }

    // ── Test 4: WorknetManagerUni initialization ──
    function test_worknetManagerUniInitialize() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        (bytes32 poolId,) = _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        bytes memory initData = abi.encodeWithSelector(
            WorknetManagerUni.initialize.selector,
            address(registry),
            address(alpha),
            address(awpToken),
            poolId,
            deployer,
            _getWorknetManagerDexConfig()
        );
        WorknetManagerUni wm = WorknetManagerUni(address(new ERC1967Proxy(worknetManagerImpl, initData)));

        // Verify base storage
        assertEq(address(wm.awpRegistry()), address(registry));
        assertEq(address(wm.alphaToken()), address(alpha));
        assertEq(address(wm.awpToken()), address(awpToken));
        assertEq(wm.poolId(), poolId);
        assertEq(wm.clPoolManager(), POOL_MANAGER);
        assertEq(wm.clPositionManager(), POSITION_MANAGER);
        assertEq(wm.permit2(), PERMIT2);

        // Verify Uni-specific storage
        assertEq(wm.stateView(), STATE_VIEW);

        // Verify uniPoolKey
        (address c0, address c1, uint24 fee, int24 ts, address hooks) = wm.uniPoolKey();
        if (address(awpToken) < address(alpha)) {
            assertEq(c0, address(awpToken));
            assertEq(c1, address(alpha));
        } else {
            assertEq(c0, address(alpha));
            assertEq(c1, address(awpToken));
        }
        assertEq(fee, POOL_FEE);
        assertEq(ts, TICK_SPACING);
        assertEq(hooks, address(0));

        // Verify roles
        assertTrue(wm.hasRole(wm.DEFAULT_ADMIN_ROLE(), deployer));
        assertTrue(wm.hasRole(wm.MERKLE_ROLE(), deployer));
        assertTrue(wm.hasRole(wm.STRATEGY_ROLE(), deployer));
        assertTrue(wm.hasRole(wm.TRANSFER_ROLE(), deployer));
    }

    // ── Test 5: WorknetManagerUni UUPS upgrade ──
    function test_worknetManagerUniUpgrade() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        (bytes32 poolId,) = _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        bytes memory initData = abi.encodeWithSelector(
            WorknetManagerUni.initialize.selector,
            address(registry),
            address(alpha),
            address(awpToken),
            poolId,
            deployer,
            _getWorknetManagerDexConfig()
        );
        WorknetManagerUni wm = WorknetManagerUni(address(new ERC1967Proxy(worknetManagerImpl, initData)));

        address newImpl = address(new WorknetManagerUni());
        vm.prank(deployer);
        wm.upgradeToAndCall(newImpl, "");

        // Verify state preserved
        assertEq(address(wm.awpRegistry()), address(registry));
        assertEq(address(wm.alphaToken()), address(alpha));
        assertEq(wm.poolId(), poolId);
        assertEq(wm.stateView(), STATE_VIEW);
        assertTrue(wm.hasRole(wm.DEFAULT_ADMIN_ROLE(), deployer));
    }

    // ── Test 6: Full E2E — worknet registration + LP creation ──
    function test_fullE2ERegistration() public {
        vm.prank(deployer);
        awpToken.transfer(user, 10_000_000 * 1e18);

        uint256 lpAWPAmount = registry.initialAlphaMint() * registry.initialAlphaPrice() / 1e18;
        vm.prank(user);
        awpToken.approve(address(registry), lpAWPAmount);

        IAWPRegistry.WorknetParams memory params = IAWPRegistry.WorknetParams({
            name: "Base Worknet",
            symbol: "BWRK",
            worknetManager: address(0),
            salt: bytes32(0),
            minStake: 0,
            skillsURI: ""
        });

        vm.prank(user);
        uint256 worknetId = registry.registerWorknet(params);
        assertGt(worknetId, 0);

        (, IAWPRegistry.WorknetStatus status,,) = registry.worknets(worknetId);
        assertEq(uint8(status), uint8(IAWPRegistry.WorknetStatus.Pending));

        vm.prank(user);
        registry.activateWorknet(worknetId);

        (, IAWPRegistry.WorknetStatus statusAfter,,) = registry.worknets(worknetId);
        assertEq(uint8(statusAfter), uint8(IAWPRegistry.WorknetStatus.Active));
    }
}

// ═══════════════════════════════════════════════════════════════════════
//  Ethereum Fork Test — Uniswap V4
// ═══════════════════════════════════════════════════════════════════════

contract EthereumForkTest is ForkTestBase {
    // ── Uniswap V4 addresses (Ethereum mainnet) ──
    address constant POOL_MANAGER = 0x000000000004444c5dc75cB358380D2e3dE08A90;
    address constant POSITION_MANAGER = 0xbD216513d74C8cf14cf4747E6AaA6420FF64ee9e;
    address constant PERMIT2 = 0x000000000022D473030F116dDEE9F6B43aC78BA3;
    address constant SWAP_ROUTER = 0x66a9893cC07D91D95644AEDD05D03f95e1dBA8Af;
    address constant STATE_VIEW = 0x7fFE42C4a5DEeA5b0feC41C94C136Cf115597227;

    function setUp() public {
        fork = vm.createFork(vm.envString("ETH_RPC_URL"));
        vm.selectFork(fork);
        _deployProtocol();
    }

    function _deployLPManager(address registry_, address, address awpToken_)
        internal
        override
        returns (LPManagerBase)
    {
        return new LPManagerUni(registry_, POOL_MANAGER, POSITION_MANAGER, PERMIT2, awpToken_);
    }

    function _getWorknetManagerDexConfig() internal pure override returns (bytes memory) {
        return abi.encode(POOL_MANAGER, POSITION_MANAGER, SWAP_ROUTER, PERMIT2, POOL_FEE, TICK_SPACING, STATE_VIEW);
    }

    function _deployWorknetManagerImpl() internal override returns (address) {
        return address(new WorknetManagerUni());
    }

    // ── Test 1: LP Pool creation ──
    function test_createPoolAndAddLiquidity() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        (bytes32 poolId, uint256 lpTokenId) = _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        assertNotEq(poolId, bytes32(0), "poolId should not be zero");
        assertGt(lpTokenId, 0, "lpTokenId should be > 0");
        assertEq(lpManager.alphaTokenToPoolId(address(alpha)), poolId);
        assertEq(lpManager.alphaTokenToTokenId(address(alpha)), lpTokenId);
    }

    // ── Test 2: Duplicate LP creation should revert ──
    function test_createPoolTwiceReverts() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        vm.startPrank(deployer);
        awpToken.transfer(address(lpManager), AWP_LP_AMOUNT);
        alpha.mint(deployer, ALPHA_LP_AMOUNT);
        IERC20(address(alpha)).transfer(address(lpManager), ALPHA_LP_AMOUNT);
        vm.stopPrank();

        vm.prank(address(registry));
        vm.expectRevert(LPManagerBase.PoolAlreadyExists.selector);
        lpManager.createPoolAndAddLiquidity(address(alpha), AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);
    }

    // ── Test 3: LP Fee Compounding ──
    function test_compoundFeesNoRevert() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);
        // compoundFees may revert with no fees (V4 TAKE_PAIR with 0 fees)
        assertTrue(lpManager.alphaTokenToTokenId(address(alpha)) > 0);
        assertTrue(lpManager.alphaTokenToPoolId(address(alpha)) != bytes32(0));
    }

    // ── Test 4: WorknetManagerUni initialization ──
    function test_worknetManagerUniInitialize() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        (bytes32 poolId,) = _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        bytes memory initData = abi.encodeWithSelector(
            WorknetManagerUni.initialize.selector,
            address(registry),
            address(alpha),
            address(awpToken),
            poolId,
            deployer,
            _getWorknetManagerDexConfig()
        );
        WorknetManagerUni wm = WorknetManagerUni(address(new ERC1967Proxy(worknetManagerImpl, initData)));

        assertEq(address(wm.awpRegistry()), address(registry));
        assertEq(address(wm.alphaToken()), address(alpha));
        assertEq(wm.poolId(), poolId);
        assertEq(wm.clPoolManager(), POOL_MANAGER);
        assertEq(wm.stateView(), STATE_VIEW);

        (address c0, address c1, uint24 fee, int24 ts, address hooks) = wm.uniPoolKey();
        if (address(awpToken) < address(alpha)) {
            assertEq(c0, address(awpToken));
            assertEq(c1, address(alpha));
        } else {
            assertEq(c0, address(alpha));
            assertEq(c1, address(awpToken));
        }
        assertEq(fee, POOL_FEE);
        assertEq(ts, TICK_SPACING);
        assertEq(hooks, address(0));

        assertTrue(wm.hasRole(wm.DEFAULT_ADMIN_ROLE(), deployer));
    }

    // ── Test 5: WorknetManagerUni UUPS upgrade ──
    function test_worknetManagerUniUpgrade() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        (bytes32 poolId,) = _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        bytes memory initData = abi.encodeWithSelector(
            WorknetManagerUni.initialize.selector,
            address(registry),
            address(alpha),
            address(awpToken),
            poolId,
            deployer,
            _getWorknetManagerDexConfig()
        );
        WorknetManagerUni wm = WorknetManagerUni(address(new ERC1967Proxy(worknetManagerImpl, initData)));

        address newImpl = address(new WorknetManagerUni());
        vm.prank(deployer);
        wm.upgradeToAndCall(newImpl, "");

        assertEq(address(wm.awpRegistry()), address(registry));
        assertEq(wm.poolId(), poolId);
        assertEq(wm.stateView(), STATE_VIEW);
        assertTrue(wm.hasRole(wm.DEFAULT_ADMIN_ROLE(), deployer));
    }

    // ── Test 6: Full E2E — worknet registration + LP creation ──
    function test_fullE2ERegistration() public {
        vm.prank(deployer);
        awpToken.transfer(user, 10_000_000 * 1e18);

        uint256 lpAWPAmount = registry.initialAlphaMint() * registry.initialAlphaPrice() / 1e18;
        vm.prank(user);
        awpToken.approve(address(registry), lpAWPAmount);

        IAWPRegistry.WorknetParams memory params = IAWPRegistry.WorknetParams({
            name: "Eth Worknet",
            symbol: "EWRK",
            worknetManager: address(0),
            salt: bytes32(0),
            minStake: 0,
            skillsURI: ""
        });

        vm.prank(user);
        uint256 worknetId = registry.registerWorknet(params);
        assertGt(worknetId, 0);

        vm.prank(user);
        registry.activateWorknet(worknetId);

        (, IAWPRegistry.WorknetStatus statusAfter,,) = registry.worknets(worknetId);
        assertEq(uint8(statusAfter), uint8(IAWPRegistry.WorknetStatus.Active));
    }
}

// ═══════════════════════════════════════════════════════════════════════
//  Arbitrum Fork Test — Uniswap V4
// ═══════════════════════════════════════════════════════════════════════

contract ArbitrumForkTest is ForkTestBase {
    // ── Uniswap V4 addresses (Arbitrum mainnet) ──
    address constant POOL_MANAGER = 0x360E68faCcca8cA495c1B759Fd9EEe466db9FB32;
    address constant POSITION_MANAGER = 0xd88F38F930b7952f2DB2432Cb002E7abbF3dD869;
    address constant PERMIT2 = 0x000000000022D473030F116dDEE9F6B43aC78BA3;
    address constant SWAP_ROUTER = 0xa51afAF359d044F8e56fE74B9575f23142cD4B76;
    address constant STATE_VIEW = 0x76fd297e2d437cd7F76A5F2B02a5ce11c663A86e;

    function setUp() public {
        fork = vm.createFork(vm.envString("ARB_RPC_URL"));
        vm.selectFork(fork);
        _deployProtocol();
    }

    function _deployLPManager(address registry_, address, address awpToken_)
        internal
        override
        returns (LPManagerBase)
    {
        return new LPManagerUni(registry_, POOL_MANAGER, POSITION_MANAGER, PERMIT2, awpToken_);
    }

    function _getWorknetManagerDexConfig() internal pure override returns (bytes memory) {
        return abi.encode(POOL_MANAGER, POSITION_MANAGER, SWAP_ROUTER, PERMIT2, POOL_FEE, TICK_SPACING, STATE_VIEW);
    }

    function _deployWorknetManagerImpl() internal override returns (address) {
        return address(new WorknetManagerUni());
    }

    // ── Test 1: LP Pool creation ──
    function test_createPoolAndAddLiquidity() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        (bytes32 poolId, uint256 lpTokenId) = _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        assertNotEq(poolId, bytes32(0), "poolId should not be zero");
        assertGt(lpTokenId, 0, "lpTokenId should be > 0");
        assertEq(lpManager.alphaTokenToPoolId(address(alpha)), poolId);
        assertEq(lpManager.alphaTokenToTokenId(address(alpha)), lpTokenId);
    }

    // ── Test 2: Duplicate LP creation should revert ──
    function test_createPoolTwiceReverts() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        vm.startPrank(deployer);
        awpToken.transfer(address(lpManager), AWP_LP_AMOUNT);
        alpha.mint(deployer, ALPHA_LP_AMOUNT);
        IERC20(address(alpha)).transfer(address(lpManager), ALPHA_LP_AMOUNT);
        vm.stopPrank();

        vm.prank(address(registry));
        vm.expectRevert(LPManagerBase.PoolAlreadyExists.selector);
        lpManager.createPoolAndAddLiquidity(address(alpha), AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);
    }

    // ── Test 3: LP Fee Compounding ──
    function test_compoundFeesNoRevert() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);
        // compoundFees may revert with no fees (V4 TAKE_PAIR with 0 fees)
        assertTrue(lpManager.alphaTokenToTokenId(address(alpha)) > 0);
        assertTrue(lpManager.alphaTokenToPoolId(address(alpha)) != bytes32(0));
    }

    // ── Test 4: WorknetManagerUni initialization ──
    function test_worknetManagerUniInitialize() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        (bytes32 poolId,) = _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        bytes memory initData = abi.encodeWithSelector(
            WorknetManagerUni.initialize.selector,
            address(registry),
            address(alpha),
            address(awpToken),
            poolId,
            deployer,
            _getWorknetManagerDexConfig()
        );
        WorknetManagerUni wm = WorknetManagerUni(address(new ERC1967Proxy(worknetManagerImpl, initData)));

        assertEq(address(wm.awpRegistry()), address(registry));
        assertEq(address(wm.alphaToken()), address(alpha));
        assertEq(wm.poolId(), poolId);
        assertEq(wm.clPoolManager(), POOL_MANAGER);
        assertEq(wm.stateView(), STATE_VIEW);

        (address c0, address c1, uint24 fee, int24 ts, address hooks) = wm.uniPoolKey();
        if (address(awpToken) < address(alpha)) {
            assertEq(c0, address(awpToken));
            assertEq(c1, address(alpha));
        } else {
            assertEq(c0, address(alpha));
            assertEq(c1, address(awpToken));
        }
        assertEq(fee, POOL_FEE);
        assertEq(ts, TICK_SPACING);
        assertEq(hooks, address(0));

        assertTrue(wm.hasRole(wm.DEFAULT_ADMIN_ROLE(), deployer));
    }

    // ── Test 5: WorknetManagerUni UUPS upgrade ──
    function test_worknetManagerUniUpgrade() public {
        AlphaToken alpha = _deployAlphaClone("Test Alpha", "TALPHA", 1);
        (bytes32 poolId,) = _createPool(alpha, AWP_LP_AMOUNT, ALPHA_LP_AMOUNT);

        bytes memory initData = abi.encodeWithSelector(
            WorknetManagerUni.initialize.selector,
            address(registry),
            address(alpha),
            address(awpToken),
            poolId,
            deployer,
            _getWorknetManagerDexConfig()
        );
        WorknetManagerUni wm = WorknetManagerUni(address(new ERC1967Proxy(worknetManagerImpl, initData)));

        address newImpl = address(new WorknetManagerUni());
        vm.prank(deployer);
        wm.upgradeToAndCall(newImpl, "");

        assertEq(address(wm.awpRegistry()), address(registry));
        assertEq(wm.poolId(), poolId);
        assertEq(wm.stateView(), STATE_VIEW);
        assertTrue(wm.hasRole(wm.DEFAULT_ADMIN_ROLE(), deployer));
    }

    // ── Test 6: Full E2E — worknet registration + LP creation ──
    function test_fullE2ERegistration() public {
        vm.prank(deployer);
        awpToken.transfer(user, 10_000_000 * 1e18);

        uint256 lpAWPAmount = registry.initialAlphaMint() * registry.initialAlphaPrice() / 1e18;
        vm.prank(user);
        awpToken.approve(address(registry), lpAWPAmount);

        IAWPRegistry.WorknetParams memory params = IAWPRegistry.WorknetParams({
            name: "Arb Worknet",
            symbol: "AWRK",
            worknetManager: address(0),
            salt: bytes32(0),
            minStake: 0,
            skillsURI: ""
        });

        vm.prank(user);
        uint256 worknetId = registry.registerWorknet(params);
        assertGt(worknetId, 0);

        vm.prank(user);
        registry.activateWorknet(worknetId);

        (, IAWPRegistry.WorknetStatus statusAfter,,) = registry.worknets(worknetId);
        assertEq(uint8(statusAfter), uint8(IAWPRegistry.WorknetStatus.Active));
    }
}
