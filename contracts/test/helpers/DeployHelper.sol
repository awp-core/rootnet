// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

// ── Protocol contracts ──
import {AWPToken} from "../../src/token/AWPToken.sol";
import {AlphaTokenFactory} from "../../src/token/AlphaTokenFactory.sol";
import {AWPEmission} from "../../src/token/AWPEmission.sol";
import {AWPRegistry} from "../../src/AWPRegistry.sol";
import {AWPAllocator} from "../../src/core/AWPAllocator.sol";
import {AWPWorkNet} from "../../src/core/AWPWorkNet.sol";
import {veAWP} from "../../src/core/veAWP.sol";
import {Treasury} from "../../src/governance/Treasury.sol";
import {AWPDAO} from "../../src/governance/AWPDAO.sol";
import {IAWPRegistry} from "../../src/interfaces/IAWPRegistry.sol";
import {TimelockControllerUpgradeable} from "@openzeppelin/contracts-upgradeable/governance/TimelockControllerUpgradeable.sol";

// ── Mocks ──
import {MockLPManager} from "./MockLPManager.sol";
import {MockWorknetManager} from "./MockWorknetManager.sol";

/// @title DeployHelper — Shared setUp for all AWP protocol tests
/// @dev Solves circular dependencies via stub proxy pattern:
///      1. Deploy AWPToken + Treasury (no deps)
///      2. Predict proxy addresses for AWPWorkNet, AWPAllocator via CREATE2-like nonce prediction
///      3. Deploy AWPRegistry impl with all 8 immutable addresses
///      4. Deploy AWPRegistry proxy
///      5. Deploy AWPWorkNet proxy (needs awpRegistry)
///      6. Deploy veAWP (needs awpAllocator — predict or deploy allocator first)
///      7. Deploy AWPAllocator proxy (needs awpRegistry + veAWP)
///      8. Deploy AWPEmission proxy
///      9. Wire everything up
abstract contract DeployHelper is Test {
    // ── Accounts ──
    address public deployer;
    address public guardian;
    address public alice;
    address public bob;
    address public relayer;

    // ── Protocol contracts ──
    AWPToken public awp;
    Treasury public treasury;
    AWPRegistry public awpRegistry;
    AWPWorkNet public awpWorkNet;
    AWPAllocator public awpAllocator;
    veAWP public veAwp;
    AWPEmission public awpEmission;
    AlphaTokenFactory public factory;
    MockLPManager public lpManager;
    AWPDAO public dao;

    // ── Proxy addresses (predicted) ──
    address public registryProxy;
    address public workNetProxy;
    address public allocatorProxy;
    address public emissionProxy;

    // ── Constants ──
    uint256 constant INITIAL_MINT = 200_000_000 * 1e18;
    uint256 constant GENESIS_TIME = 1_700_000_000;
    uint256 constant INITIAL_DAILY_EMISSION = 31_600_000 * 1e18;

    function _deployAll() internal {
        deployer = address(this);
        guardian = makeAddr("guardian");
        alice = makeAddr("alice");
        bob = makeAddr("bob");
        relayer = makeAddr("relayer");

        vm.deal(deployer, 100 ether);
        vm.deal(guardian, 100 ether);
        vm.deal(alice, 100 ether);
        vm.deal(bob, 100 ether);

        // ── Step 1: Independent contracts ──
        awp = new AWPToken("AWP Token", "AWP", deployer);
        awp.initialMint(INITIAL_MINT);

        address[] memory proposers = new address[](0);
        address[] memory executors = new address[](1);
        executors[0] = address(0);
        treasury = new Treasury(0, proposers, executors, deployer);

        factory = new AlphaTokenFactory(deployer, 0);

        // ── Step 2: Predict proxy addresses ──
        // We deploy stubs for contracts that have circular deps.
        // The trick: deploy proxies with a minimal stub, then upgrade to real impl.

        // Predict nonce-based addresses for proxies we'll deploy
        uint256 startNonce = vm.getNonce(deployer);

        // We'll deploy in this order:
        // nonce+0: AWPWorkNet stub impl
        // nonce+1: AWPWorkNet proxy
        // nonce+2: MockLPManager
        // nonce+3: AWPAllocator stub (minimal Initializable+UUPS)
        // nonce+4: AWPAllocator proxy
        // nonce+5: veAWP
        // nonce+6: AWPEmission impl
        // nonce+7: AWPEmission proxy
        // nonce+8: AWPRegistry impl
        // nonce+9: AWPRegistry proxy

        workNetProxy = vm.computeCreateAddress(deployer, startNonce + 1);
        allocatorProxy = vm.computeCreateAddress(deployer, startNonce + 4);
        address veAwpAddr = vm.computeCreateAddress(deployer, startNonce + 5);
        emissionProxy = vm.computeCreateAddress(deployer, startNonce + 7);
        registryProxy = vm.computeCreateAddress(deployer, startNonce + 9);

        // ── Step 3: Deploy AWPWorkNet (UUPS proxy) ──
        // Stub impl: just needs awpRegistry immutable
        AWPWorkNet workNetImpl = new AWPWorkNet(registryProxy); // nonce+0
        awpWorkNet = AWPWorkNet(address(new ERC1967Proxy( // nonce+1
            address(workNetImpl),
            abi.encodeCall(AWPWorkNet.initialize, ("AWP WorkNet", "WORKN", guardian))
        )));

        // ── Step 4: MockLPManager ──
        lpManager = new MockLPManager(registryProxy, address(awp)); // nonce+2

        // ── Step 5: AWPAllocator (UUPS proxy) — needs awpRegistry + veAWP ──
        AWPAllocator allocatorImpl = new AWPAllocator(registryProxy, veAwpAddr); // nonce+3
        awpAllocator = AWPAllocator(address(new ERC1967Proxy( // nonce+4
            address(allocatorImpl),
            abi.encodeCall(AWPAllocator.initialize, (registryProxy, guardian))
        )));

        // ── Step 6: veAWP (non-upgradeable) — needs awpToken + awpAllocator + guardian ──
        veAwp = new veAWP(address(awp), allocatorProxy, guardian); // nonce+5

        // ── Step 7: AWPEmission (UUPS proxy) ──
        AWPEmission emissionImpl = new AWPEmission(address(awp)); // nonce+6
        awpEmission = AWPEmission(address(new ERC1967Proxy( // nonce+7
            address(emissionImpl),
            abi.encodeCall(AWPEmission.initialize, (
                address(awp), guardian, INITIAL_DAILY_EMISSION, GENESIS_TIME, 1 days, address(treasury)
            ))
        )));

        // ── Step 8: AWPRegistry (UUPS proxy) — needs all 8 immutables ──
        AWPRegistry registryImpl = new AWPRegistry( // nonce+8
            address(awp),
            address(awpWorkNet),
            address(factory),
            address(awpEmission),
            address(lpManager),
            address(awpAllocator),
            address(veAwp),
            address(treasury)
        );
        awpRegistry = AWPRegistry(address(new ERC1967Proxy( // nonce+9
            address(registryImpl),
            abi.encodeCall(AWPRegistry.initialize, (deployer, address(treasury), guardian))
        )));

        // ── Step 9: Verify predicted addresses ──
        require(address(awpWorkNet) == workNetProxy, "WorkNet proxy address mismatch");
        require(address(awpAllocator) == allocatorProxy, "Allocator proxy address mismatch");
        require(address(veAwp) == veAwpAddr, "veAWP address mismatch");
        require(address(awpEmission) == emissionProxy, "Emission proxy address mismatch");
        require(address(awpRegistry) == registryProxy, "Registry proxy address mismatch");

        // ── Step 10: Post-deploy setup ──
        // Set WorknetManager impl (real UUPS mock, not address(1))
        MockWorknetManager wmImpl = new MockWorknetManager();
        vm.prank(guardian);
        awpRegistry.setWorknetManagerImpl(address(wmImpl));

        // AWP minter
        awp.addMinter(address(awpEmission));
        awp.renounceAdmin();

        // AlphaTokenFactory
        factory.setAddresses(address(awpRegistry));

        // Fund test accounts
        awp.transfer(alice, 10_000_000 * 1e18);
        awp.transfer(bob, 10_000_000 * 1e18);

        // Set initial alpha price/mint (already set in initialize defaults)
        // initialAlphaPrice = 1e15, initialAlphaMint = 1e27
    }

    /// @dev Register a worknet as `user` and return worknetId
    function _registerWorknet(address user) internal returns (uint256) {
        return _registerWorknet(user, "TestWorknet", "TWN");
    }

    function _registerWorknet(address user, string memory name, string memory symbol) internal returns (uint256) {
        uint256 cost = awpRegistry.initialAlphaMint() * awpRegistry.initialAlphaPrice() / 1e18;
        vm.startPrank(user);
        awp.approve(address(awpRegistry), cost);
        uint256 wid = awpRegistry.registerWorknet(IAWPRegistry.WorknetParams({
            name: name,
            symbol: symbol,
            worknetManager: address(0),
            salt: bytes32(0),
            minStake: 0,
            skillsURI: ""
        }));
        vm.stopPrank();
        return wid;
    }

    /// @dev Activate a worknet as guardian
    function _activateWorknet(uint256 worknetId) internal {
        vm.prank(guardian);
        awpRegistry.activateWorknet(worknetId);
    }
}
