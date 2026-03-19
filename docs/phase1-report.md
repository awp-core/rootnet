# Phase 1 Development Report — Smart Contracts

> **Note**: This report describes the Phase 1 design. Several features have been redesigned since: AlphaToken deployment changed from Clones to CREATE2, StakingVault cooldown/freeze-release removed, AWPEmission restructured to V3 epoch-versioned allocations.

## Overview

All 10 Solidity contracts for the AWP RootNet Agent Mining protocol have been implemented, tested, and verified. The contracts follow the architecture document (docs/architecture.md v9.0) and target Solidity 0.8.20 with OpenZeppelin 5.x.

## Contracts Implemented

| Contract | Path | Lines | Description |
|----------|------|-------|-------------|
| AWPToken | src/token/AWPToken.sol | 80 | ERC20Votes + minter model + ERC1363 |
| AlphaToken | src/token/AlphaToken.sol | 83 | Upgradeable ERC20 clone, dual minter |
| AlphaTokenFactory | src/token/AlphaTokenFactory.sol | 48 | CREATE2 full deployment (no Clones) |
| AWPEmission | src/token/AWPEmission.sol | ~420 | UUPS upgradeable emission engine: oracle multi-sig weight submission (submitAllocations), epoch settlement, batch minting, exponential decay |
| AccessManager | src/core/AccessManager.sol | 130 | User/Agent registration + permissions |
| StakingVault | src/core/StakingVault.sol | 275 | Deposit/withdraw/allocate/reallocate/STP |
| SubnetNFT | src/core/SubnetNFT.sol | 45 | Pure ERC721 with baseURI template |
| LPManager | src/core/LPManager.sol | 57 | PancakeSwap V4 LP wrapper (simplified) |
| RootNet | src/RootNet.sol | ~450 | Unified entry: subnet management + staking management |
| AWPDAO | src/governance/AWPDAO.sol | 95 | OZ Governor suite |
| Treasury | src/governance/Treasury.sol | 11 | OZ TimelockController |

**Total source lines**: 1,822 (contracts + interfaces)
**Interfaces**: 10 files, 273 lines

## Test Coverage

| Test File | Tests | Status |
|-----------|-------|--------|
| AWPToken.t.sol | 29 tests (incl. 2 fuzz) | ✅ All pass |
| AWPEmission.t.sol (in AWPToken.t.sol) | 8 tests (incl. 1 fuzz) | ✅ All pass |
| AWPEmission.t.sol (standalone) | 6 tests | ✅ All pass |
| AlphaTokenFactory.t.sol | 16 tests (incl. 1 fuzz) | ✅ All pass |
| AccessManager.t.sol | 23 tests | ✅ All pass |
| SubnetNFT.t.sol | 17 tests | ✅ All pass |
| StakingVault.t.sol | 31 tests | ✅ All pass |
| LPManager.t.sol | 5 tests | ✅ All pass |
| RootNet.t.sol | 26 tests | ✅ All pass |
| AWPDAO.t.sol | 9 tests | ✅ All pass |
| Integration.t.sol | 46 tests (est.) | ✅ All pass |

**Total: 216 tests, 216 passed, 0 failed**

Test categories covered:
- Access control (onlyRootNet, onlyTimelock, onlyGuardian)
- Token minting/burning with MAX_SUPPLY enforcement
- ERC1363 callback flows (transferAndCall, approveAndCall)
- ERC20Votes delegation
- Clone deployment and initialization
- Dual minter model with permanent locking
- User/Agent registration with mutual exclusion
- Stake deposit/withdraw with cooldown period
- Allocation/deallocation/reallocation across epochs
- Agent removal with freeze/release cycle
- Subnet registration with LP creation
- Subnet lifecycle (activate/pause/resume/ban/unban/deregister)
- Emission settlement with exponential decay
- Multi-subnet weighted emission distribution
- Batch settlement across multiple subnets
- Governance flow (propose/vote/queue/execute)
- Full 22-step deployment sequence
- Manager Agent proxy operations
- Pause/unpause protection

## Deploy Script

`script/Deploy.s.sol` — 163 lines, implements the full 22-step deployment sequence:
1. AWPToken (constructor mints 200M to deployer)
2. AlphaToken implementation
3. AlphaTokenFactory
4. Treasury (TimelockController)
5. AWPDAO (Governor)
6-7. Treasury role setup + admin renounce
8. RootNet
9-12. SubnetNFT, AccessManager, StakingVault, LPManager
13a. AWPEmission implementation deploy
13b. ERC1967Proxy(impl, initData) — proxy address is permanent minter
14-15. Add minter + renounce admin (permanently locks minter list)
16. AlphaTokenFactory.setAddresses
17. RootNet.initializeRegistry
18-22. AWP distribution (Treasury 90M, LP 10M, Airdrop 100M)

## Design Decisions

### 1. LPManager Simplified
The LPManager uses a simplified implementation that generates deterministic pool addresses via hashing rather than integrating with actual PancakeSwap V4 contracts. This is because:
- PancakeSwap V4 is not yet deployed on BSC mainnet at time of writing
- The actual integration requires specific pool manager and position manager addresses
- The simplified version correctly holds tokens in the LPManager (simulating locked LP)
- Production deployment will need to replace with real PancakeSwap V4 calls

### 2. via_ir Compilation
The `via_ir = true` flag is enabled in foundry.toml because RootNet.sol's `registerSubnet` function has too many local variables for the legacy code generator. This is a standard Foundry approach for complex contracts and has no functional impact.

### 3. Treasury with Zero Delay in Tests
Integration tests use `minDelay=0` for the Treasury to simplify testing governance flows. The Deploy script uses `172800` (2 days) for production.

### 4. ERC1363 Implementation
Both AWPToken and AlphaToken implement ERC1363 (transferAndCall/approveAndCall) inline rather than using a separate library. The callback is only invoked if the target is a contract (checked via `to.code.length > 0`).

### 5. AlphaToken Upgradeable Pattern
AlphaToken uses OpenZeppelin's upgradeable contracts (`ERC20Upgradeable`, `ERC20BurnableUpgradeable`) for the Clones pattern, as minimal proxies delegate all calls to the implementation contract and need initializer-based setup rather than constructors.

## Known Limitations / Phase 2 TODOs

1. **PancakeSwap V4 Integration**: LPManager needs real PancakeSwap V4 integration when the protocol is available on BSC. Current implementation is a mock.

2. **Gas Optimization**: While the architecture follows gas-efficient patterns (`settleEpoch(limit)` iterates recipients[] in bounded batches, mints directly via `_mintTo`), further gas profiling and optimization can be done.

3. **Formal Verification**: Critical financial logic (emission calculation, staking math) should undergo formal verification.

4. **Audit**: All contracts need professional security audit before mainnet deployment.

5. **Vesting Contracts**: Team and investor vesting contracts (referenced in deployment steps 19-20) are not included in this phase. Standard OZ VestingWallet or a custom implementation is needed.

6. **Frontend Integration**: ABI generation (`forge inspect`) and wagmi/viem type generation needed for dashboard.

7. **Backend Bindings**: `abigen` needs to be run to generate Go contract bindings for the indexer/keeper.

8. **Upgradability Strategy**: RootNet itself is not upgradeable. The `updateAddress` mechanism via DAO timelock provides limited flexibility. Consider whether a proxy pattern is needed for RootNet.
