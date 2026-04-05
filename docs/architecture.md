# AWP Protocol Architecture Reference

> **Version**: 10.0
> **Project**: AWP -- Agent Mining protocol
> **Stack**: Solidity 0.8.33 (Foundry, optimizer_runs=800, via_ir=true, evm_version=cancun) + Go 1.26 (Chi + sqlc + pgx) + PostgreSQL
> **Target**: Base (Uniswap V4), Ethereum, Arbitrum, BSC (PancakeSwap V4)
> **Principle**: Maximize reuse of OpenZeppelin 5.x

---

## 1. Project Structure

```
awp-registry/
├── CLAUDE.md
├── docs/
│   └── architecture.md
│
├── contracts/
│   ├── src/
│   │   ├── AWPRegistry.sol                    # Main contract: unified entry + worknet management + account system (UUPS proxy)
│   │   ├── token/
│   │   │   ├── AWPToken.sol                   # ERC20Votes + minter (10B)
│   │   │   ├── AlphaToken.sol                 # Dual minter (10B/worknet)
│   │   │   ├── AlphaTokenFactory.sol          # CREATE2 full deployment (no Clones)
│   │   │   └── AWPEmission.sol                # UUPS proxy: Guardian-only emission engine
│   │   ├── core/
│   │   │   ├── WorknetNFT.sol                 # ERC721 with on-chain identity storage
│   │   │   ├── StakingVault.sol               # UUPS proxy + EIP-712 gasless allocation
│   │   │   ├── StakeNFT.sol                   # ERC721 position NFT (deposit/withdraw AWP)
│   │   │   └── LPManager.sol                  # Uniswap V4 / PancakeSwap V4 LP
│   │   ├── worknets/
│   │   │   └── WorknetManager.sol             # Default worknet contract (UUPS, behind ERC1967Proxy)
│   │   ├── governance/
│   │   │   ├── AWPDAO.sol                     # OZ Governor + StakeNFT-based voting
│   │   │   └── Treasury.sol                   # OZ TimelockController
│   │   └── interfaces/
│   │       ├── IAWPRegistry.sol
│   │       ├── IAWPToken.sol
│   │       ├── IAlphaToken.sol
│   │       ├── IAWPEmission.sol
│   │       ├── IWorknetNFT.sol
│   │       ├── IStakingVault.sol
│   │       ├── IStakeNFT.sol
│   │       ├── ILPManager.sol
│   │       └── IAlphaTokenFactory.sol
│   ├── test/
│   │   ├── AWPEmission.t.sol
│   │   ├── AlphaTokenFactory.t.sol
│   │   ├── AWPRegistry.t.sol
│   │   ├── AWPRegistryExtended.t.sol
│   │   ├── AWPDAOExtended.t.sol
│   │   ├── AWPDAO.t.sol
│   │   ├── WorknetNFT.t.sol
│   │   ├── WorknetManager.t.sol
│   │   ├── StakingVault.t.sol
│   │   ├── StakingVaultExtended.t.sol
│   │   ├── Integration.t.sol
│   │   ├── E2E.t.sol
│   │   ├── MultiChainE2E.t.sol
│   │   └── ForkTest.t.sol
│   ├── script/
│   │   ├── Deploy.s.sol
│   │   ├── TestDeploy.s.sol
│   │   ├── Predict.s.sol
│   │   └── InitCodeHashes.s.sol
│   ├── foundry.toml
│   └── remappings.txt
│
├── api/
│   ├── cmd/
│   │   ├── api/main.go                    # HTTP server entry
│   │   ├── indexer/main.go                # Chain indexer entry
│   │   └── keeper/main.go                 # Keeper bot entry
│   ├── internal/
│   │   ├── config/config.go               # Env config (caarlos0/env)
│   │   ├── server/
│   │   │   ├── server.go                  # Chi router + middleware
│   │   │   ├── handler/                   # HTTP handlers (grouped by domain)
│   │   │   │   ├── user.go
│   │   │   │   ├── agent.go
│   │   │   │   ├── staking.go
│   │   │   │   ├── subnet.go
│   │   │   │   ├── emission.go
│   │   │   │   ├── token.go
│   │   │   │   ├── governance.go
│   │   │   │   ├── relay.go
│   │   │   │   └── vanity.go
│   │   │   └── ws/hub.go                  # WebSocket hub (Redis Pub/Sub)
│   │   ├── chain/
│   │   │   ├── client.go                  # go-ethereum RPC client
│   │   │   ├── bindings/                  # abigen-generated contract bindings
│   │   │   ├── indexer.go                 # Event listening + DB writes
│   │   │   └── keeper.go                  # settleEpoch + token prices
│   │   ├── db/
│   │   │   ├── query/                     # sqlc SQL files
│   │   │   ├── migrations/                # Atlas migration files
│   │   │   └── gen/                       # sqlc-generated code
│   │   └── service/                       # Business logic layer
│   ├── go.mod
│   ├── go.sum
│   ├── sqlc.yaml
│   ├── atlas.hcl
│   ├── Makefile
│   └── Dockerfile
│
├── scripts/
│   ├── deploy-multichain.sh
│   ├── deploy_safe.json
│   └── admin.sh
├── dashboard/
├── docker-compose.yaml
├── .env.example
└── README.md
```

---

## 2. Architecture Overview

```
                          ┌──────────────┐
                          │   AWPDAO     │
                          └──────┬───────┘
                          ┌──────▼───────┐
                          │   Treasury   │
                          └──────┬───────┘
                                 │
  User ──────────────────► ┌─────▼──────┐       Guardian (Safe 3/5 multisig)
                          │ AWPRegistry │           │
                          │ (UUPS      │ ◄─────────┘ (pause, setGuardian, upgrades)
                          │  Proxy)    │
                          └─┬──┬──┬──┬──┘
                            │  │  │  │
              ┌─────────────┘  │  │  └──────────┐
              ▼                ▼  ▼              ▼
        ┌────────┐   ┌──────────┐ ┌──────────┐ ┌──────────┐
        │Worknet │   │ Staking  │ │ LP       │ │ AWP      │
        │ NFT    │   │ Vault    │ │ Manager  │ │ Emission │
        └────────┘   │(UUPS+712)│ └──────────┘ │(UUPS)    │
                     └──────────┘       │       └──────────┘
        ┌────────┐        │        V4 DEX          │
        │Worknet │   ┌────▼───┐                mint AWP
        │Manager │   │ Stake  │
        │(UUPS)  │   │ NFT    │
        └────────┘   └────────┘
                    deposit/withdraw AWP

Contracts: 11 (AWPRegistry, AWPEmission, StakingVault, WorknetNFT, WorknetManager,
               StakeNFT, LPManager, AlphaTokenFactory, AWPToken, AlphaToken, AWPDAO+Treasury)

Key design:
  - AWPEmission is a generic address->weight distribution engine (no worknet awareness)
  - Guardian (Safe multisig) submits weights directly; no Oracle multi-sig, no Timelock dependency for emission
  - 100% of epoch emission goes to Guardian-submitted recipients (Guardian includes treasury for DAO share)
  - settleEpoch() callable by anyone (Keeper calls it); 3-phase batched settlement
  - All UUPS upgrades controlled by Guardian (not Timelock)
  - Allocation functions live on StakingVault (not AWPRegistry)
  - WorknetId globally unique: (block.chainid << 64) | localCounter
```

### OZ Reuse

```
AWPToken            -> ERC20, ERC20Permit, ERC20Votes, ERC20Burnable, IERC1363
AlphaToken          -> ERC20, ERC20Permit, ERC20Burnable (standalone CREATE2 deployment)
AlphaTokenFactory   -> Ownable (CREATE2 full deployment, no Clones)
AWPEmission         -> Initializable, UUPSUpgradeable, ReentrancyGuardUpgradeable (no EIP712)
AWPRegistry         -> Initializable, UUPSUpgradeable, PausableUpgradeable, ReentrancyGuardUpgradeable, EIP712Upgradeable
WorknetNFT          -> ERC721
StakingVault        -> Initializable, UUPSUpgradeable, ReentrancyGuardUpgradeable, EIP712Upgradeable, EnumerableSet
StakeNFT            -> ERC721, ReentrancyGuard
WorknetManager      -> Initializable, UUPSUpgradeable, AccessControlUpgradeable, ReentrancyGuardUpgradeable, IERC1363Receiver
AWPDAO              -> Governor, GovernorSettings, GovernorTimelockControl
Treasury            -> TimelockController
```

### Permission Model

```
Account System V2 (no mandatory registration):
  Every address:   Implicitly a root; can bind(target) to form delegation trees
  register():      Removed from AWPRegistry
  Delegates:       grantDelegate(delegate) / revokeDelegate(delegate) for staking operations
  Staker:          allocate/deallocate/reallocate on StakingVault -- staker is explicit parameter

Inter-contract Permissions:
  onlyAWPRegistry   -> WorknetNFT / LPManager
  onlyGuardian      -> AWPRegistry (pause, setGuardian, UUPS upgrade)
                    -> AWPEmission (submitAllocations, config, setGuardian, UUPS upgrade)
                    -> StakingVault (UUPS upgrade)
  onlyTimelock      -> AWPRegistry admin functions (setInitialAlphaPrice, setImmunityPeriod, ban/unban/deregister)
  staker or delegate -> StakingVault (allocate, deallocate, reallocate)

Guardian = Safe multisig 0x000002bEfa6A1C99A710862Feb6dB50525dF00A3 (3/5 threshold)
```

---

## 3. Tokenomics

```
AWP:
  Total Supply: 10,000,000,000 (10B), MAX_SUPPLY hard cap
    -> INITIAL_MINT configurable per chain (constructor param, immutable)
    -> Remainder minted on demand by AWPEmission (exponential decay, ~99% in 4 years)

  Mining Emission Curve:
    Each epoch, AWPEmission mints AWP with exponential decay:
      On-chain: currentEmission = lastEmission * 996844 / 1000000
      decayFactor = 996844 (corresponds to ~e^(-0.00316), precision 1e6)
      decayFactor is Guardian-configurable (min 900000)

    | Time       | Approximate Cumulative |
    |------------|----------------------|
    | Year 1     | ~34%                 |
    | Year 4     | ~99%                 |

  Per-Epoch AWP Distribution:
    100% -> Recipients (by weight). Guardian includes treasury address in recipients for DAO share.
    No hardcoded 50/50 split. Guardian controls the ratio via submitted weights.

Alpha (per worknet):
  Max Supply: 10,000,000,000 (10B)
  Minted at registration: INITIAL_ALPHA_MINT (100M) for LP creation
  Subsequent minting: managed by worknet contract (WorknetManager has Merkle claim)
  Standalone CREATE2 deployment (not proxy)
  supplyAtLock snapshot + createdAt reset at setWorknetMinter

LP Creation Cost:
  Initial price: initialAlphaPrice (default 0.01 AWP/Alpha)
  Alpha mint: 100M
  AWP cost: 0.01 x 100M = 1M AWP
  LP: 1M AWP + 100M Alpha -> V4 DEX (Uniswap on Base/ETH/Arb, PancakeSwap on BSC)
```

---

## 4. Contract Implementation

### 4.1 AWPToken.sol

```
Inherits: ERC20, ERC20Permit, ERC20Votes, ERC20Burnable
Implements: IERC1363 (transferAndCall / approveAndCall)

MAX_SUPPLY = 10_000_000_000 * 1e18
INITIAL_MINT: constructor param (immutable, configurable per chain)

Storage:
  mapping(address => bool) public minters;
  address public admin;                          // deployer -> renounced after setup

constructor(name, symbol, deployer, initialMint)
  -> admin = deployer
  -> _mint(deployer, initialMint)

addMinter(address minter) external            // admin only
renounceAdmin() external                      // admin only, permanent
mint(address to, uint256 amount) external     // minters only, capped at MAX_SUPPLY
mintAndCall(address to, uint256 amount, bytes data) external
  -> mint + ERC1363 onTransferReceived callback on recipient

Deployer is never a minter; only receives INITIAL_MINT in constructor for distribution via transfer.
```

### 4.2 AlphaToken.sol -- Dual Minter + Callback + Burnable

```
Inherits: ERC20, ERC20Permit, ERC20Burnable
Standalone CREATE2 deployment (not proxy, not upgradeable)

Storage:
  worknetId: uint256
  admin: address                                // AWPRegistry
  mapping(address => bool) public minters;
  mapping(address => bool) public minterPaused;
  bool public mintersLocked;
  uint256 public supplyAtLock;                   // Snapshot at setWorknetMinter
  uint64 public createdAt;                       // Reset at setWorknetMinter
  MAX_SUPPLY = 10_000_000_000 * 1e18

Key functions:
  setWorknetMinter(address worknetContract) external
    -> require(msg.sender == admin && !mintersLocked)
    -> minters[worknetContract] = true
    -> minters[admin] = false
    -> supplyAtLock = totalSupply()              // Snapshot pre-mint supply
    -> createdAt = uint64(block.timestamp)       // Reset clock for worknet minter
    -> mintersLocked = true

  setMinterPaused(address minter, bool paused) external
    -> require(msg.sender == admin)              // For ban/unban
```

### 4.3 AlphaTokenFactory.sol

```
Inherits: Ownable
CREATE2 full deployment (no Clones / EIP-1167 proxy). Each AlphaToken is standalone.

Storage:
  rootNet: address          // set once via setAddresses
  configured: bool
  vanityRule: uint64        // immutable, set in constructor

Vanity rule encoding (8 positions packed into uint64):
  [prefix0][prefix1][prefix2][prefix3][suffix0][suffix1][suffix2][suffix3]
  Value meaning:
    0-9   -> digit '0'-'9'
    10-15 -> lowercase hex 'a'-'f' (EIP-55 must NOT uppercase)
    16-21 -> uppercase hex 'A'-'F' (EIP-55 must uppercase)
    >=22  -> wildcard
  vanityRule=0 -> skip all validation
  Example "A1????cafe": 0x1001FFFF0C0A0F0E

deploy(worknetId, name, symbol, admin, salt) -> require(msg.sender == awpRegistry)
  -> effectiveSalt = (salt == bytes32(0)) ? bytes32(worknetId) : salt
  -> new AlphaToken{salt: effectiveSalt}()
  -> if (vanityRule != 0): _validateVanityAddress(address(token))
```

### 4.4 AWPEmission.sol -- Guardian-Only Emission Engine (UUPS Proxy, V3)

```
Purpose: UUPS upgradeable proxy. Generic address->weight distribution engine.
         Guardian (cross-chain Safe multisig) submits epoch-versioned packed allocations.
         settleEpoch(limit) batch-mints AWP with exponential decay.
         No Oracle signatures. No EIP-712 inheritance. No Timelock dependency.
         100% emission to recipients; Guardian includes treasury in recipients for DAO share.
Inherits: Initializable, UUPSUpgradeable, ReentrancyGuardUpgradeable

Storage:
  awpToken: IAWPToken                             // slot 1
  treasury: address                               // slot 2 (queryable reference)
  epochDuration: uint256                          // slot 3 (default 86400 = 1 day, configurable)
  settledEpoch: uint256                           // slot 4 (next epoch to settle, starts at 0)
  baseTime: uint256                               // slot 5 (initially = genesisTime)
  currentDailyEmission: uint256                   // slot 6
  activeEpoch: uint256                            // slot 7 (most recently promoted weights)
  _epochAllocations: mapping(uint256 => uint256[]) // slot 8 (packed entries per epoch)
  _epochTotalWeight: mapping(uint256 => uint256)   // slot 10 (cross-chain global total per epoch)
  settleProgress: uint256                         // slot 11 (0 = idle)
  epochEmissionLocked: uint256                    // slot 12
  baseEpoch: uint256                              // slot 18
  pausedUntil: uint64                             // slot 19 (packed)
  frozenEpoch: uint64                             // slot 19 (packed)
  maxRecipients: uint256                          // slot 20 (default 10000)
  decayFactor: uint256                            // slot 24 (default 996844)
  guardian: address                               // slot 26

Constants:
  DECAY_PRECISION = 1000000
  MIN_DECAY_FACTOR = 900000

Packed allocation format:
  uint256: [32-bit reserved | 64-bit weight | 160-bit address]
  modifyAllocations uses: [32-bit index | 64-bit weight | 160-bit address]

Epoch calculation:
  currentEpoch() = baseEpoch + (block.timestamp - baseTime) / epochDuration
  When paused: returns frozenEpoch
  When pause expired: resumes from frozenEpoch at pausedUntil timestamp

Functions:

  initialize(awpToken_, guardian_, initialDailyEmission_, genesisTime_, epochDuration_, treasury_)

  // Guardian weight submission
  submitAllocations(uint256[] packed_, uint256 totalWeight_, uint256 effectiveEpoch) onlyGuardian
    -> Replaces existing allocations for the epoch
    -> totalWeight_ is cross-chain global total (sum of all chains)
    -> Empty array allowed (no emission for this chain's epoch)

  appendAllocations(uint256[] packed_, uint256 effectiveEpoch) onlyGuardian
    -> Append additional recipients to an existing epoch
    -> totalWeight is NOT updated (was set in submitAllocations)

  modifyAllocations(uint256[] patches_, uint256 newTotalWeight_, uint256 effectiveEpoch) onlyGuardian
    -> In-place patch: top 32 bits = array index
    -> weight=0 allowed (soft-delete, skipped at settle)
    -> address(0) keeps existing address (weight-only update)

  // Settlement (3-phase, callable by anyone)
  settleEpoch(uint256 limit) external nonReentrant
    Phase 1 (init):
      -> Applies decay: currentDailyEmission *= decayFactor / DECAY_PRECISION
      -> Caps at remaining AWP mintable supply
      -> Promotes activeEpoch if new weights exist; otherwise carries forward
    Phase 2 (batch mint):
      -> Iterates up to `limit` recipients per call
      -> awpToken.mint(recipient, amount) + best-effort ERC1363 callback
      -> Failed callbacks silently ignored (try/catch)
    Phase 3 (finalize):
      -> Resets progress, increments settledEpoch

  // Guardian config
  setMaxRecipients(uint256) onlyGuardian
  setDecayFactor(uint256) onlyGuardian           // >= MIN_DECAY_FACTOR, < DECAY_PRECISION
  setEpochDuration(uint256) onlyGuardian         // Rebases epoch calculation
  setTreasury(address) onlyGuardian
  setGuardian(address) onlyGuardian              // Self-sovereign
  pauseEpochUntil(uint64 resumeTime) onlyGuardian

  _authorizeUpgrade(address) -> onlyGuardian

Events:
  AllocationsSubmitted(uint256 indexed effectiveEpoch, uint256[] packed, uint256 totalWeight)
  AllocationsAppended(uint256 indexed effectiveEpoch, uint256[] packed)
  AllocationsModified(uint256 indexed effectiveEpoch, uint256[] patches, uint256 newTotalWeight)
  RecipientAWPDistributed(uint256 indexed epoch, address indexed recipient, uint256 awpAmount)
  EpochSettled(uint256 indexed epoch, uint256 totalEmission, uint256 recipientCount)
  EpochDurationUpdated(uint256 oldDuration, uint256 newDuration)
  EpochPausedUntil(uint64 resumeTime, uint64 frozenEpoch)
  GuardianUpdated(address indexed newGuardian)
  TreasuryUpdated(address indexed newTreasury)
  MaxRecipientsUpdated(uint256 newMax)
  DecayFactorUpdated(uint256 newDecayFactor)

Removed from V2/V1:
  Oracle multi-sig, EIP-712 signatures, emissionSplitBps, emergencySetWeight, daoShare,
  setOracleConfig, DAOMatchDistributed, GovernanceWeightUpdated
```

### 4.5 WorknetNFT.sol -- ERC721 with On-Chain Identity

```
Inherits: ERC721
modifier onlyAWPRegistry()

Storage:
  struct WorknetIdentity {
    string name;              // Worknet / Alpha token name
    address worknetManager;   // Worknet contract address
    address alphaToken;       // Alpha token address
  }
  struct WorknetMeta {
    string skillsURI;         // Skills file URI
    uint128 minStake;         // Minimum stake (off-chain reference only, not enforced on-chain)
    string metadataURI;       // Custom metadata JSON URI
  }

  tokenId = worknetId (globally unique: (block.chainid << 64) | localCounter)

Functions:
  mint(to, tokenId, name, worknetManager, alphaToken, skillsURI, minStake) onlyAWPRegistry
  burn(tokenId) onlyAWPRegistry
  setSkillsURI(tokenId, uri) -> onlyOwner; emit SkillsURIUpdated
  setMinStake(tokenId, minStake) -> onlyOwner; emit MinStakeUpdated
  setMetadataURI(tokenId, uri) -> onlyOwner; emit MetadataURIUpdated
  setBaseURI(string uri) onlyAWPRegistry

  tokenURI 3-tier: per-token metadataURI -> global baseURI -> on-chain Base64 JSON
  Lifecycle status managed by AWPRegistry, not WorknetNFT.
```

### 4.6 StakingVault.sol -- UUPS Proxy with EIP-712 Gasless Support

```
Inherits: Initializable, UUPSUpgradeable, ReentrancyGuardUpgradeable, EIP712Upgradeable
EIP-712 domain name: "StakingVault"

Core Design:
  -> UUPS proxy with gasless support (allocateFor, deallocateFor)
  -> Allocation functions live here (moved from AWPRegistry)
  -> Auth: caller must be staker or staker's delegate (reads AWPRegistry.delegates cross-contract)
  -> (staker, agent, worknetId) triple
  -> All operations immediate (no pending mechanism)
  -> worknetId=0 rejected
  -> Cross-chain allocate: worknetId from any chain, no on-chain status check
  -> Auto-enumerates agent worknets via EnumerableSet

Storage:
  awpRegistry: address
  guardian: address                                     // upgrade auth
  stakeNFT: address                                     // for balance checks
  _allocations: mapping(address => mapping(address => mapping(uint256 => uint128)))
  _agentWorknets: mapping(address => mapping(address => EnumerableSet.UintSet))
  userTotalAllocated: mapping(address => uint256)
  worknetTotalStake: mapping(uint256 => uint256)
  nonces: mapping(address => uint256)                   // EIP-712 replay prevention

Functions:
  allocate(staker, agent, worknetId, amount)             // staker or delegate
  deallocate(staker, agent, worknetId, amount)           // staker or delegate
  reallocate(staker, fromAgent, fromWorknetId, toAgent, toWorknetId, amount)
  allocateFor(staker, agent, worknetId, amount, deadline, v, r, s)    // gasless EIP-712
  deallocateFor(staker, agent, worknetId, amount, deadline, v, r, s)  // gasless EIP-712
  getUnallocated(user) view -> stakeNFT.totalStakedOf(user) - userTotalAllocated[user]

  _authorizeUpgrade(address) -> onlyGuardian

Note: freezeAgentAllocations has been removed.
```

### 4.7 StakeNFT.sol -- ERC721 Position NFT

```
Inherits: ERC721, ReentrancyGuard

Core Design:
  -> Deposit AWP with lock period (lockDuration in seconds) -> mint position NFT
  -> Each position = NFT with (amount, lockEndTime, createdAt)
  -> NFTs are transferable (ERC721)
  -> O(1) balance tracking via _userTotalStaked accumulator
  -> addToPosition blocked on expired locks (PositionExpired)
  -> withdraw burns NFT after lock expires

Functions:
  deposit(amount, lockDuration) -> tokenId
  depositFor(to, amount, lockDuration) -> tokenId
  addToPosition(tokenId, amount, newLockEndTime) -> blocked if lock expired
  withdraw(tokenId)
  getVotingPower(tokenId) view
    -> amount * sqrt(min(remainingTime, 54 weeks) / 7 days)
  totalStakedOf(user) view -> _userTotalStaked[user]
```

### 4.8 LPManager.sol -- DEX V4 LP

```
modifier onlyAWPRegistry()

Functions:
  createPoolAndAddLiquidity(alphaToken, awpAmount, alphaAmount)
    -> Create V4 pool + initialize price + full-range two-sided LP
    -> LP NFT stays in LPManager (permanently locked)
    -> return (pool, lpTokenId)
  compoundFees(alphaToken)
    -> Reinvests accumulated LP fees (called by Keeper cron)

No removeLiquidity. LP permanently locked.
LPManager/WorknetManager bytecode differs by DEX (Uniswap vs PancakeSwap).
BSC uses PancakeSwap V4; Base/ETH/Arb use Uniswap V4.
```

### 4.9 AWPRegistry.sol -- Main Contract (UUPS Proxy)

```
Inherits: Initializable, UUPSUpgradeable, PausableUpgradeable,
          ReentrancyGuardUpgradeable, EIP712Upgradeable
EIP-712 domain name: "AWPRegistry"

Storage:
  // Address registry
  awpToken, worknetNFT, alphaTokenFactory, awpEmission, lpManager,
  stakingVault, stakeNFT, treasury, guardian: address
  defaultWorknetManagerImpl: address
  dexConfig: bytes

  // Account System V2
  registeredCount: uint256
  mapping(address => address) public boundTo;
  mapping(address => address) public recipient;
  mapping(address => mapping(address => bool)) public delegates;
  _deployer: address                               // Temporary; zeroed after initializeRegistry
  registryInitialized: bool

  // Worknet data (on-chain stores only lifecycle state)
  struct WorknetInfo {
    bytes32 lpPool;
    WorknetStatus status;
    uint64 createdAt; uint64 activatedAt; uint64 immunityEndsAt;
  }
  enum WorknetStatus { Pending, Active, Paused, Banned }
  mapping(uint256 => WorknetInfo) public worknets;
  uint256 private _nextLocalId = 1;                // local counter; worknetId = (chainid << 64) | _nextLocalId++
  uint256 public initialAlphaPrice = 1e16;
  uint256 public constant INITIAL_ALPHA_MINT = 100_000_000 * 1e18;
  EnumerableSet.UintSet private activeWorknetIds;
  uint256 public constant MAX_ACTIVE_WORKNETS = 10000;

  // EIP-712 nonces for gasless operations
  mapping(address => uint256) public nonces;

Functions:

  // Initialization
  initialize(deployer, treasury, guardian) external initializer
  initializeRegistry(...) external
    -> require(msg.sender == _deployer && !registryInitialized)
    -> 9 params: 8 addresses + dexConfig bytes
    -> registryInitialized = true; _deployer = address(0)

  // Account System V2 (register() removed)
  bind(address target) external whenNotPaused
    -> _checkNoCycle(msg.sender, target)
  setRecipient(address _recipient) external
  grantDelegate(address delegate) external
  revokeDelegate(address delegate) external
  resolveRecipient(address addr) view
    -> Walks boundTo chain to root, returns recipient[root] or root

  // Gasless account operations (EIP-712)
  bindFor(address user, address target, uint256 deadline, bytes sig) external
  setRecipientFor(address user, address rec, uint256 deadline, bytes sig) external
  registerWorknetFor(address user, WorknetParams params, uint256 deadline, bytes sig) external
  registerWorknetForWithPermit(address user, WorknetParams params, ...) external
    -> Fully gasless: ERC-2612 permit + EIP-712 in one tx

  // Worknet Registration (AWP payment + auto LP)
  registerWorknet(WorknetParams calldata params) external nonReentrant whenNotPaused -> uint256
    -> If worknetManager == address(0) and defaultWorknetManagerImpl set: auto-deploy proxy
    -> AWP transferFrom(user -> LPManager)
    -> Alpha mint(LPManager)
    -> LPManager.createPoolAndAddLiquidity()
    -> AlphaToken.setWorknetMinter(worknetManager)
    -> WorknetNFT.mint(owner, worknetId, name, worknetManager, alphaToken, skillsURI, minStake)
    -> Store WorknetInfo (lifecycle state only)

  // Worknet Lifecycle
  activateWorknet(worknetId) external          // owner, Pending -> Active
  pauseWorknet(worknetId) external             // owner, Active -> Paused
  resumeWorknet(worknetId) external            // owner, Paused -> Active
  banWorknet(worknetId) external onlyTimelock  // Active|Paused -> Banned (pauses minter)
  unbanWorknet(worknetId) external onlyTimelock // Banned -> Active (checks MAX_ACTIVE_WORKNETS)
  deregisterWorknet(worknetId) external onlyTimelock // After immunity period

  // Admin (onlyTimelock)
  setInitialAlphaPrice(uint256) external onlyTimelock
  setImmunityPeriod(uint256) external onlyTimelock

  // Guardian
  setGuardian(address g) external onlyGuardian
  pause() external onlyGuardian
  unpause() external onlyTimelock

  _authorizeUpgrade(address) -> onlyGuardian

Events:
  Bound, RecipientUpdated, DelegateGranted, DelegateRevoked,
  WorknetRegistered, LPCreated, WorknetActivated, WorknetPaused,
  WorknetResumed, WorknetBanned, WorknetUnbanned, WorknetDeregistered,
  GuardianUpdated
```

### 4.10 WorknetManager.sol -- Default Worknet Contract

```
Inherits: Initializable, UUPSUpgradeable, AccessControlUpgradeable,
          ReentrancyGuardUpgradeable, IERC1363Receiver
Deployed behind ERC1967Proxy by AWPRegistry when worknetManager=address(0).

Roles:
  MERKLE_ROLE    -> Submit Merkle roots for Alpha distribution
  STRATEGY_ROLE  -> AWP handling strategy
  TRANSFER_ROLE  -> Token transfers

AWP Strategy: Reserve / AddLiquidity / BuybackBurn
  onTransferReceived auto-executes strategy on AWP receipt via mintAndCall
  DEX addresses injected at init time via dexConfig (not hardcoded)

Merkle claim: users prove inclusion -> Alpha minted to user
```

### 4.11 AWPDAO.sol + Treasury.sol

```
AWPDAO: OZ Governor + GovernorSettings + GovernorTimelockControl
  Overrides _getVotes and _countVote for StakeNFT-based voting.
  No awpRegistry dependency. No delegate/checkpoint mechanism.
  -> Voters submit tokenId[] arrays (StakeNFT position NFTs)
  -> Voting power = amount * sqrt(min(remainingTime, 54 weeks) / 7 days)
  -> Anti-manipulation: only NFTs with createdAt < proposalCreatedAt (strict: >= blocks same-block)
  -> Per-tokenId double-vote prevention
  -> totalVotingPower > 0 required for proposal creation
  -> Two proposal types:
     proposeWithTokens: executable via Timelock
     signalPropose: vote-only, no on-chain execution
  -> propose() is blocked (reverts)

Treasury: OZ TimelockController (zero custom code)
```

---

## 5. Multi-Chain Deployment

```
Chains: Base (8453), Ethereum (1), Arbitrum (42161), BSC (56)
Deploy config: chains.yaml
Deploy script: scripts/deploy-multichain.sh <chainName|--all|--list>

Same CREATE2 salts -> identical addresses on all chains.
Exception: LPManager/WorknetManager bytecode differs by DEX, so BSC addresses
  differ for these two contracts. Use per-chain overrides in chains.yaml.

GENESIS_TIME must be set explicitly (no block.timestamp fallback).
AWPToken.initialMint configurable per chain.

WorknetId: (block.chainid << 64) | localCounter -- globally unique across all chains.
Allocate is local: user allocates on their staking chain to any chain's worknet.
Emission: per-chain AWPEmission, Guardian coordinates quotas via totalWeight.
DAO: per-chain AWPDAO + Treasury, off-chain aggregated voting.

Vanity addresses: contracts/.env defines VANITY_PREFIX_*/VANITY_SUFFIX_* for key contracts.
  deploy.sh mines salts via scripts/vanity/mine.sh before deployment.
  Changing constructors/init params invalidates initCodeHash -> must re-mine salts.
```

### Deployment Sequence

```
1.  AWPToken(name, symbol, deployer, initialMint)
2.  AlphaTokenFactory(deployer, vanityRule)
3.  Treasury(delay, [], [address(0)], deployer)
4.  AWPRegistry impl -> ERC1967Proxy(impl, initialize(deployer, treasury, guardian))
5.  WorknetNFT("AWP Worknet", "AWPWN", awpRegistry)
6.  LPManager(awpRegistry, poolManager, positionManager, awpToken)
7.  AWPEmission impl -> ERC1967Proxy(impl, initialize(awpToken, guardian, initialDailyEmission, genesisTime, epochDuration, treasury))
8.  StakingVault impl -> ERC1967Proxy(impl, initialize(awpRegistry, treasury))
9.  StakeNFT(awpToken, stakingVault, awpRegistry)
10. WorknetManager impl
11. AWPDAO(stakeNFT, treasury, ...) -- 6 params, no awpRegistry
12. Treasury.grantRole(PROPOSER+CANCELLER, awpDAO)
13. Treasury.renounceRole(ADMIN, deployer)
14. AWPToken.addMinter(awpEmissionProxy)
15. AWPToken.renounceAdmin()                    // Permanently lock minter list
16. AlphaTokenFactory.setAddresses(awpRegistry)
17. AWPRegistry.initializeRegistry(9 params: 8 addresses + dexConfig bytes)
18. AWP transfer distribution (treasury, liquidity, airdrop)

Post-deploy:
  AWPToken minters = { awpEmission } (sole minter)
  AWPToken admin = address(0) (permanently locked)
  Deployer was never a minter; only received INITIAL_MINT in constructor for distribution
```

---

## 6. Coordinator (Off-Chain)

```
Coordinator is the off-chain operations service for a worknet, deployed by the worknet owner.

Responsibilities:
  1. Identity Verification -- Listen to AWPRegistry events, maintain agent cache, verify heartbeat signatures
  2. Task Management -- Assign tasks, collect results, evaluate quality, compute contribution scores
  3. Reward Distribution -- WorknetManager mints Alpha + receives AWP -> distributes to miners' rewardRecipient

How worknets query AWPRegistry:
  Cold start: getAgentsInfo() full pull
  Running: Listen to events for incremental cache updates (Bound, Allocated, Deallocated, ...)

Different worknet types have different coordinator logic:
  Benchmark: Generate problems / solve / verify
  DATA Mining: Collect data / verify / deduplicate
  AI Arena: Match opponents / ELO scoring
```

---

## 7. API (Go)

### Tech Stack

```
Go 1.26

Core:
  chi/v5           -> HTTP router + middleware (CORS, rate-limit, recover, logger)
  pgx/v5           -> PostgreSQL driver (native)
  sqlc             -> SQL -> Go code generation
  go-ethereum      -> On-chain interaction + abigen contract bindings

Data:
  Atlas            -> Declarative DB migration
  go-redis/v9      -> Redis cache + Pub/Sub

Realtime:
  github.com/coder/websocket -> WebSocket

Infrastructure:
  log/slog         -> Structured logging
  caarlos0/env     -> Env vars -> struct
  uber-go/fx       -> Dependency injection + lifecycle management
  robfig/cron/v3   -> Keeper scheduled tasks

Code Generation:
  abigen           -> Solidity ABI -> Go contract bindings
  sqlc             -> SQL -> Go query functions
```

### Architecture

```
Backend is read-only + on-chain data indexer + gasless relay.

Write operations:
  Frontend wagmi/viem -> direct to chain (no backend proxy)
  Gasless relay: 10 relay endpoints for EIP-712 signed operations

Three independent Go processes:
  cmd/api/main.go      -> HTTP read API + WebSocket + relay
  cmd/indexer/main.go  -> Chain Indexer (event sync -> PostgreSQL)
  cmd/keeper/main.go   -> Keeper Bot (settleEpoch + compoundFees + token prices)

Inter-process communication:
  Indexer -> Redis Pub/Sub (channel: chain_events) -> API WebSocket broadcast
  Indexer uses optimistic indexing with block hash chain verification for reorg detection (max 64-block rollback)
```

### Routes

```
/api
  GET  /registry                        # Contract addresses + chainId + eip712Domain
  GET  /health                          # Health check
  GET  /health/detailed                 # Detailed health (DB, Redis, RPC)
  GET  /chains                          # Multi-chain info
  GET  /stats                           # Global protocol stats

  /users
    GET  /                              # List users
    GET  /count                         # User count
    GET  /global                        # Cross-chain user list
    GET  /{address}                     # User detail
    GET  /{address}/portfolio           # Portfolio (staking + allocations)
    GET  /{address}/delegates           # Delegates

  GET  /address/{address}/check         # Address type check
  GET  /address/{address}/resolve-recipient
  POST /address/batch-resolve-recipients
  GET  /nonce/{address}                 # AWPRegistry EIP-712 nonce
  GET  /staking-nonce/{address}         # StakingVault EIP-712 nonce

  /agents
    GET  /by-owner/{owner}
    GET  /by-owner/{owner}/{agent}
    GET  /lookup/{agent}
    POST /batch-info

  /staking
    GET  /user/{address}/balance
    GET  /user/{address}/balance/global
    GET  /user/{address}/positions
    GET  /user/{address}/positions/global
    GET  /user/{address}/allocations
    GET  /user/{address}/pending
    GET  /user/{address}/frozen
    GET  /agent/{agent}/subnet/{worknetId}
    GET  /agent/{agent}/subnets
    GET  /subnet/{worknetId}/total

  /subnets
    GET  /                              # List worknets
    GET  /ranked                        # Worknets by total stake
    GET  /search                        # Search worknets
    GET  /by-owner/{owner}
    GET  /{worknetId}
    GET  /{worknetId}/skills
    GET  /{worknetId}/earnings
    GET  /{worknetId}/agents
    GET  /{worknetId}/agents/{agent}

  /emission
    GET  /current
    GET  /schedule
    GET  /global-schedule
    GET  /epochs
    GET  /epochs/{epochId}

  /tokens
    GET  /awp
    GET  /awp/global
    GET  /alpha/{worknetId}
    GET  /alpha/{worknetId}/price

  /governance
    GET  /proposals
    GET  /proposals/global
    GET  /proposals/{proposalId}
    GET  /treasury

/api/admin (Bearer token auth)
  POST   /chains
  DELETE /chains/{chainId}
  GET    /chains
  PUT    /ratelimit
  GET    /ratelimit
  GET    /system

/api/relay (gasless EIP-712, rate-limited)
  POST /register                        # (Note: registerWorknet relay)
  POST /bind
  POST /set-recipient
  POST /allocate
  POST /deallocate
  POST /activate-subnet
  POST /register-subnet
  POST /grant-delegate
  POST /revoke-delegate
  POST /unbind
  GET  /status/{txHash}

/api/vanity
  GET  /mining-params
  POST /upload-salts
  GET  /salts
  GET  /salts/count
  POST /compute-salt                    # DB pool first, cast fallback

/api/announcements (public + admin CRUD)

/v2                                     # JSON-RPC 2.0 entry point

/ws/live                                # WebSocket real-time events
```

### DB Schema (Key Tables)

```sql
CREATE TABLE users (
    address       CHAR(42) PRIMARY KEY,
    bound_to      CHAR(42),
    recipient     CHAR(42),
    registered_at BIGINT NOT NULL
);

CREATE TABLE worknets (
    worknet_id       BIGINT PRIMARY KEY,   -- (chainid << 64) | localCounter
    owner            CHAR(42) NOT NULL,
    name             VARCHAR(64) NOT NULL,
    symbol           VARCHAR(16) NOT NULL,
    worknet_manager  CHAR(42) NOT NULL,
    alpha_token      CHAR(42) NOT NULL,
    lp_pool          CHAR(42),             -- nullable (WorknetRegistered precedes LPCreated)
    status           VARCHAR(16) NOT NULL DEFAULT 'Pending',
    created_at       BIGINT NOT NULL,
    activated_at     BIGINT,
    immunity_ends_at BIGINT,
    burned           BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE stake_allocations (
    user_address  CHAR(42) NOT NULL,
    agent_address CHAR(42) NOT NULL,
    worknet_id    BIGINT NOT NULL,
    amount        NUMERIC(78,0) NOT NULL DEFAULT 0,
    PRIMARY KEY (user_address, agent_address, worknet_id)
);

CREATE TABLE stake_positions (
    token_id      BIGINT PRIMARY KEY,
    owner         CHAR(42) NOT NULL,
    amount        NUMERIC(78,0) NOT NULL,
    lock_end_time BIGINT NOT NULL,
    created_at    BIGINT NOT NULL
);

CREATE TABLE epochs (
    epoch_id        INTEGER PRIMARY KEY,
    start_time      BIGINT NOT NULL,
    daily_emission  NUMERIC(78,0) NOT NULL,
    total_emission  NUMERIC(78,0)
);

CREATE TABLE vanity_salts (
    id          SERIAL PRIMARY KEY,
    salt        CHAR(66) NOT NULL,
    claimed     BOOLEAN NOT NULL DEFAULT FALSE
    -- FOR UPDATE SKIP LOCKED on claim
);

CREATE TABLE sync_states (
    contract_name VARCHAR(64) PRIMARY KEY,
    last_block    BIGINT NOT NULL DEFAULT 0
);
```

### Redis Key Spec

```
Cache (Keeper writes, API reads):
  alpha_price:{worknetId}       -> JSON { priceInAWP, ... }     TTL=10m
  awp_info:{chainId}            -> JSON { totalSupply, ... }    TTL=1m
  emission_current:{chainId}    -> JSON { epoch, ... }          TTL=30s

Rate limiting:
  ratelimit:config              -> hash, persistent, hot-updatable (admin.sh)
  rl:relay:{ip}                 -> counter, TTL=1h (100/IP/1h, atomic Lua INCR+EXPIRE)
  rl:upload_salts:{ip}          -> counter, TTL=1h (5/hr/IP)
  rl:compute_salt:{ip}          -> counter, TTL=1h (20/hr/IP)

Pub/Sub:
  channel: chain_events         -> JSON { type, blockNumber, txHash, data }
```

### Keeper Bot

```
Scheduled tasks:
  @every 30s  -> trySettleEpoch (calls AWPEmission.settleEpoch(200))
  @every 5m   -> updateTokenPrices (read Alpha/AWP from V4 pools -> Redis)
  (cron)      -> compoundFees (call LPManager.compoundFees per active worknet)
```

---

## 8. Bug Prevention Checklist

```
Contracts:
  - worknetId uses (block.chainid << 64) | _nextLocalId++
  - Tree-based binding with anti-cycle check (no address mutual exclusion)
  - bind(target) walks chain to detect cycles before binding
  - resolveRecipient(addr) walks boundTo chain to root for reward distribution
  - StakeNFT: only NFTs with createdAt < proposalCreatedAt can vote (strict)
  - StakeNFT: addToPosition blocked on expired locks (PositionExpired)
  - StakingVault: allocate/deallocate/reallocate all reject worknetId=0
  - WorknetNFT.minStake stored on-chain but NOT enforced by allocate (off-chain reference)
  - AWPDAO: totalVotingPower > 0 required for proposals
  - deregisterWorknet: users must manually deallocate (frontend should alert)
  - worknetManager==address(0) auto-deploys WorknetManager proxy if defaultWorknetManagerImpl set
  - setWorknetMinter permanently locked; ban uses minterPaused
  - AWPEmission weights submitted by Guardian; epoch-versioned packed allocations
  - AWPRegistry.unbanWorknet checks MAX_ACTIVE_WORKNETS before re-adding
  - AWP: deployer is never a minter; renounceAdmin permanently locks
  - settleEpoch has nonReentrant
  - AWPEmission: mint + best-effort ERC1363 callback (try/catch, failure does not revert)
  - All UUPS upgrades gated by onlyGuardian

API / DB:
  - DB lp_pool is nullable (WorknetRegistered precedes LPCreated)
  - DB vanity_salts: CREATE2 + vanityRule verification on upload; FOR UPDATE SKIP LOCKED on claim
  - DB: agents table removed; users table has bound_to and recipient columns
  - Permit2 BSC mainnet: 0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768
  - admin.sh: Hot-update rate limits, manage salt pool, view system status
```
