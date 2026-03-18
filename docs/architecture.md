# RootNet Implementation Guide — Claude Code Development Document

> **Version**: 9.0  
> **Project**: awp-rootnet  
> **Stack**: Solidity 0.8.24 (Foundry) + Go 1.26 (Chi + sqlc + pgx) + PostgreSQL  
> **Target**: BSC Testnet → BSC Mainnet  
> **Principle**: Maximize reuse of OpenZeppelin 5.x

---

## 1. Project Structure

```
awp-rootnet/
├── CLAUDE.md
├── docs/
│   └── architecture.md
│
├── contracts/
│   ├── src/
│   │   ├── RootNet.sol                    # Main contract: unified entry + subnet/staking management
│   │   ├── token/
│   │   │   ├── AWPToken.sol               # ERC20Votes + minter (10B)
│   │   │   ├── AlphaToken.sol             # Dual minter (10B/subnet)
│   │   │   ├── AlphaTokenFactory.sol      # CREATE2 full deployment
│   │   │   └── AWPEmission.sol            # AWP emission (mint on demand)
│   │   ├── core/
│   │   │   ├── SubnetNFT.sol              # Pure ERC721
│   │   │   ├── AccessManager.sol          # User registration + Agent permissions
│   │   │   ├── StakingVault.sol           # Pure allocation logic (onlyRootNet)
│   │   │   ├── StakeNFT.sol               # ERC721 position NFT (deposit/withdraw AWP)
│   │   │   └── LPManager.sol              # PancakeSwap V4
│   │   ├── governance/
│   │   │   ├── AWPDAO.sol            # OZ Governor + StakeNFT-based voting
│   │   │   └── Treasury.sol              # OZ TimelockController
│   │   └── interfaces/
│   │       ├── IRootNet.sol
│   │       ├── IAWPToken.sol
│   │       ├── IAlphaToken.sol
│   │       ├── IAWPEmission.sol
│   │       ├── ISubnetNFT.sol
│   │       ├── IAccessManager.sol
│   │       ├── IStakingVault.sol
│   │       ├── ILPManager.sol
│   │       └── IPancakeV4.sol
│   ├── test/
│   │   ├── AWPToken.t.sol
│   │   ├── AWPEmission.t.sol
│   │   ├── AlphaTokenFactory.t.sol
│   │   ├── RootNet.t.sol
│   │   ├── SubnetNFT.t.sol
│   │   ├── AccessManager.t.sol
│   │   ├── StakingVault.t.sol
│   │   ├── LPManager.t.sol
│   │   ├── AWPDAO.t.sol
│   │   └── Integration.t.sol
│   ├── script/
│   │   └── Deploy.s.sol
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
│   │   │   │   └── governance.go
│   │   │   └── ws/hub.go                  # WebSocket hub (Redis Pub/Sub)
│   │   ├── chain/
│   │   │   ├── client.go                  # go-ethereum RPC client
│   │   │   ├── bindings/                  # abigen-generated contract bindings
│   │   │   │   ├── rootnet.go
│   │   │   │   ├── awptoken.go
│   │   │   │   ├── alphatoken.go
│   │   │   │   ├── subnetnft.go
│   │   │   │   ├── stakingvault.go
│   │   │   │   └── awpemission.go
│   │   │   ├── indexer.go                 # Event listening + DB writes
│   │   │   └── keeper.go                  # settleEpoch + token prices
│   │   ├── db/
│   │   │   ├── query/                     # sqlc SQL files
│   │   │   │   ├── user.sql
│   │   │   │   ├── agent.sql
│   │   │   │   ├── subnet.sql
│   │   │   │   ├── staking.sql
│   │   │   │   ├── emission.sql
│   │   │   │   └── governance.sql
│   │   │   ├── migrations/                # Atlas migration files
│   │   │   │   └── schema.hcl
│   │   │   └── gen/                       # sqlc-generated code (never hand-written)
│   │   │       ├── db.go
│   │   │       ├── models.go
│   │   │       └── query.sql.go
│   │   └── service/                       # Business logic layer
│   │       ├── user.go
│   │       ├── agent.go
│   │       ├── staking.go
│   │       ├── subnet.go
│   │       └── emission.go
│   ├── go.mod
│   ├── go.sum
│   ├── sqlc.yaml
│   ├── atlas.hcl
│   ├── Makefile                           # generate, migrate, build, test
│   └── Dockerfile
│
├── dashboard/
├── docker-compose.yaml
├── .env.example
└── README.md
```

---

## 2. Architecture Overview

```
                          ┌──────────────┐
                          │  AWPDAO │
                          └──────┬───────┘
                          ┌──────▼───────┐
                          │   Treasury   │
                          └──────┬───────┘
                                 │ onlyTimelock
  User / Agent ─────────► ┌─────▼──────┐
                          │  RootNet   │
                          │ Unified    │
                          │  Entry     │
                          └─┬──┬──┬──┬──┬─┘
                            │  │  │  │  │
              ┌─────────────┘  │  │  │  └──────────┐
              ▼                ▼  ▼  ▼              ▼
        ┌──────────┐   ┌────────┐ ┌──────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
        │ Access   │   │Subnet │ │Stake │ │ LP       │ │ AWP      │ │ Stake    │
        │ Manager  │   │ NFT   │ │Vault │ │ Manager  │ │ Emission │ │ NFT      │
        └──────────┘   └───────┘ └──────┘ └──────────┘ └──────────┘ └──────────┘
                                               │              │           │
                                         PancakeSwap V4   mint AWP   deposit/withdraw AWP

Contracts: 11

Emission Model:
  → AWPEmission is a generic address→weight distribution engine (weight management, epoch settlement, batch minting, decay)
  → AWPEmission has no subnet concepts; RootNet manages activeSubnetIds locally and does NOT notify AWPEmission on lifecycle changes
  → Exponential decay: emission *= 0.996844/epoch, ~99% released in 4 years
  → Per-epoch AWP: 50% → subnets (by weight), 50% → DAO Treasury
  → settleEpoch() on AWPEmission, callable by anyone (Keeper calls it)
  → Subnet contracts decide how to use AWP (add LP, distribute to miners, buy back Alpha, etc.)
  → Alpha emission managed by subnet contracts (they have mint permission)
  → RootNet holds no AWP, holds no Alpha (pure control layer)
```

### OZ Reuse

```
AWPToken           → ERC20, ERC20Permit, ERC20Votes, ERC20Burnable
AlphaToken         → ERC20Upgradeable, ERC20BurnableUpgradeable, Initializable
AlphaTokenFactory  → Ownable (no Clones; uses CREATE2 full deployment)
AWPEmission        → Initializable, UUPSUpgradeable, ReentrancyGuardUpgradeable, EIP712Upgradeable
RootNet            → Pausable, ReentrancyGuard, EnumerableSet
SubnetNFT          → ERC721
AccessManager      → EnumerableSet
StakingVault       → EnumerableSet (pure allocation, no deposit/withdraw)
StakeNFT           → ERC721, ReentrancyGuard (position NFT, deposit/withdraw AWP)
AWPDAO        → Governor, GovernorSettings, GovernorTimelockControl (overrides _getVotes/_countVote for StakeNFT-based voting)
Treasury           → TimelockController
```

### Permission Model

```
Three Key Types:
  Owner:           Stake AWP via StakeNFT + allocate via RootNet + manage Agents + transfer NFT + vote with tokenId[]
  Agent (default): Mining authentication only
  Agent (Manager): Mining + remove Agents + set Manager + allocate stake

Inter-contract Permissions:
  onlyRootNet   → AccessManager / StakingVault / SubnetNFT / LPManager
  onlyTimelock  → RootNet admin functions / AWPEmission governance (setWeight, setOracleConfig, etc.)
  onlyGuardian  → Can only pause the contract
```

---

## 3. Tokenomics

```
AWP:
  Total Supply: 10,000,000,000 (10B), MAX_SUPPLY hard cap
    → 50% minted in constructor to deployer (distributed via transfer)
    → 50% minted on demand by AWPEmission (exponential decay, ~99% in 4 years)

  Allocation:
    50%  5B   → Mining emission (AWPEmission mint on demand, exponential decay, ~99% in 4 years)
    20%  2B   → DAO Treasury (initial funds)
    10%  1B   → Team (4-year vesting, 1-year cliff)
    7.5% 750M → Investors (2-year vesting, 6-month cliff)
    10%  1B   → Initial liquidity (AWP/BNB main pool)
    2.5% 250M → Airdrop + early contributors

  Mining Emission Curve:
    Each epoch, AWPEmission mints AWP with exponential decay:
      daily_emission = initial_daily × e^(-λ × day)
      
      λ = 0.00316
      initial_daily ≈ 15,800,000 AWP

      | Time       | Daily Emission | Cumulative         |
      |------------|---------------|--------------------|
      | D1         | ~15.8M        |                    |
      | D30        | ~14.4M        |                    |
      | D90        | ~11.8M        |                    |
      | D180       | ~8.9M         |                    |
      | D365       | ~5.0M         |                    |
      | Year 1     |               | ~1.72B (34.3%)     |
      | Year 4     |               | ~4.95B (99%)       |
    
    On-chain implementation (per-epoch decay):
      decayFactor = 996844 (corresponds to e^(-0.00316), precision 1e6)
      currentEmission = lastEmission * 996844 / 1000000
      → No on-chain e^x needed, just multiplication

  Per-Epoch AWP Distribution:
    50% → Subnets (by governanceWeight, AWPEmission mints to subnetContract)
    50% → DAO Treasury (equal match, AWPEmission mints to Treasury)
    
    Example D1:
      daily_emission = 15.8M AWP
      Subnets receive: 7.9M AWP (distributed by weight to subnet contracts)
      DAO receives: 7.9M AWP

Alpha (per subnet):
  Max Supply: 10,000,000,000 (10B)
  Minted at registration: 100,000,000 (100M) for LP creation
  Subsequent minting: managed by subnet contract (emission curve and rules are subnet-autonomous)

LP Creation Cost:
  Initial price: 0.01 AWP/Alpha (adjustable by DAO)
  Alpha mint: 100M
  AWP cost: 0.01 × 100M = 1,000,000 AWP (1M)
  LP: 1M AWP + 100M Alpha → PancakeSwap V4
```

---

## 4. Contract Implementation

### 4.1 AWPToken.sol

```
Inherits: ERC20, ERC20Permit, ERC20Votes, ERC20Burnable
Implements: IERC1363 (transferAndCall / approveAndCall)

MAX_SUPPLY = 10_000_000_000 * 1e18
INITIAL_MINT = 5_000_000_000 * 1e18                    // 50% non-mining portion

Storage:
  mapping(address => bool) public minters;
  address public admin;                                // deployer → renounced after setup

constructor(name, symbol, deployer)
  → admin = deployer
  → _mint(deployer, INITIAL_MINT)                      // 50% minted directly to deployer for distribution via transfer

/// @notice Add a minter (admin only)
addMinter(address minter) external
  → require(msg.sender == admin)
  → minters[minter] = true

/// @notice Admin renounces control (called after setup, minter list permanently locked)
renounceAdmin() external
  → require(msg.sender == admin)
  → admin = address(0)

mint(address to, uint256 amount) external
  → require(minters[msg.sender], "Not minter")
  → require(totalSupply() + amount <= MAX_SUPPLY)
  → _mint(to, amount)

// ERC20Burnable: burn(amount), burnFrom(account, amount) — provided by OZ

// ERC1363 Callback:
transferAndCall(address to, uint256 amount, bytes calldata data) → bool
approveAndCall(address spender, uint256 amount, bytes calldata data) → bool

_update(...) override(ERC20, ERC20Votes)
nonces(...) override(ERC20Permit, Nonces)
```

### 4.2 AlphaToken.sol — Dual Minter + Callback + Burnable

```
Inherits: ERC20Upgradeable, ERC20BurnableUpgradeable, Initializable
Implements: IERC1363 (transferAndCall / approveAndCall)

Storage:
  subnetId: uint256
  admin: address                                  // RootNet (admin permission permanently retained)
  mapping(address => bool) public minters;        // Only { RootNet (temporary), subnetContract }
  mapping(address => bool) public minterPaused;   // Admin can pause/resume minter (for ban)
  bool public mintersLocked;                      // Once true, no new minters can be added
  MAX_SUPPLY = 10_000_000_000 * 1e18

Functions:
  initialize(name, symbol, subnetId, admin)
    → initializer
    → admin = _admin
    → minters[_admin] = true                      // RootNet initially has mint permission

  /// @notice Set subnet contract as sole minter, RootNet gives up mint, permanently locked
  setSubnetMinter(address subnetContract) external
    → require(msg.sender == admin)
    → require(!mintersLocked, "Minters locked")
    → if (subnetContract != address(0)):
        minters[subnetContract] = true
    → minters[admin] = false                      // RootNet gives up mint
    → mintersLocked = true                         // Permanently locked, no more minters

  /// @notice Admin can pause/resume minter (for ban/unban, does not change minters list)
  setMinterPaused(address minter, bool paused) external
    → require(msg.sender == admin)
    → minterPaused[minter] = paused

  mint(address to, uint256 amount) external
    → require(minters[msg.sender], "Not minter")
    → require(!minterPaused[msg.sender], "Minter paused")
    → require(totalSupply() + amount <= MAX_SUPPLY)
    → _mint(to, amount)

  // ERC20BurnableUpgradeable: burn(amount), burnFrom(account, amount) — provided by OZ

  // ERC1363 Callback:
  transferAndCall(to, amount, data) → transfer + onTransferReceived callback
  approveAndCall(spender, amount, data) → approve + onApprovalReceived callback
```

### 4.3 AlphaTokenFactory.sol

```
Inherits: Ownable
No Clones / EIP-1167 proxy. Each AlphaToken is a standalone CREATE2-deployed contract.

Storage:
  rootNet: address          // set once via setAddresses
  configured: bool
  vanityRule: uint64        // immutable, set in constructor

Vanity rule encoding (per position, 8 positions packed into uint64):
  [prefix0][prefix1][prefix2][prefix3][suffix0][suffix1][suffix2][suffix3]
  (high byte = prefix0, low byte = suffix3)
  Value meaning:
    0-9   → require digit '0'-'9'
    10-15 → require lowercase hex 'a'-'f' (EIP-55 must NOT uppercase)
    16-21 → require uppercase hex 'A'-'F' (EIP-55 must uppercase)
    >=22  → wildcard (no check)
  vanityRule=0 → skip all validation
  Example "A1????cafe": 0x1001FFFF0C0A0F0E

Functions:
  constructor(deployer, vanityRule)
    → No implementation parameter (no Clones)
  setAddresses(rootNet) onlyOwner → configured=true → renounceOwnership()
  deploy(subnetId, name, symbol, admin, salt) → require(msg.sender == rootNet)
    → effectiveSalt = (salt == bytes32(0)) ? bytes32(subnetId) : salt
    → new AlphaToken{salt: effectiveSalt}()
    → token.initialize(name, symbol, subnetId, admin)
    → if (vanityRule != 0): _validateVanityAddress(address(token))
    → return address(token)
```

### 4.2 AlphaToken.sol — time-cap fix

After `setSubnetMinter(subnetContract)`:
  - `supplyAtLock = totalSupply()` — snapshots pre-mint supply so LP mint is excluded from time cap
  - `createdAt = block.timestamp` — resets clock so subnet minter can mint immediately
  - `mintersLocked = true` — no further minter changes
  This prevents the 4-day lockout that would otherwise occur if `createdAt` was set at `initialize()`.
  Admin minting (initial LP liquidity) is always exempt from the time-based cap.

### 4.4 AWPEmission.sol — Generic Address→Weight Distribution Engine (UUPS Proxy)

```
Purpose: UUPS 可升级代理。通用 address→weight 分发引擎。不感知子网概念；
         任何地址（subnet contract、DAO 等）均可作为 recipient。
         oracle 多签共识提交权重（submitAllocations），结算 epoch（settleEpoch(limit)），
         批量铸造 AWP，指数衰减。唯一 AWP minter。
         DAO Timelock 管理 oracle 配置和参数；Timelock 可紧急覆盖权重（emergencySetWeight）。
         RootNet 不通知 AWPEmission 任何子网生命周期事件。
Inherits: Initializable, UUPSUpgradeable, ReentrancyGuardUpgradeable, EIP712Upgradeable

Storage:
  awpToken: IAWPToken
  treasury: address

  // Epoch-versioned packed allocations (V3: no per-address mapping; arrays per epoch)
  mapping(uint256 => address[]) internal _epochRecipients;    // epoch → recipient list
  mapping(uint256 => uint96[]) internal _epochWeights;        // epoch → weight list (parallel array)
  mapping(uint256 => uint256) public _epochTotalWeight;       // epoch → sum of weights
  uint256 public activeEpoch;                                 // latest epoch with submitted allocations
  uint256 public maxRecipients = 10000;

  // Oracle 多签共识（EIP-712 签名）
  address[] public oracles;
  uint256 public threshold;                                // 生效所需最小签名数
  uint256 public currentNonce;                             // 防重放

  // Emission（指数衰减）
  uint256 public settledEpoch = 0;
  uint256 public currentDailyEmission;                     // 当前每日排放量（衰减中）
  uint256 public constant DECAY_FACTOR = 996844;           // e^(-0.00316) × 1e6
  uint256 public constant DECAY_PRECISION = 1000000;
  uint256 public constant EMISSION_SPLIT_BPS = 5000;       // 50% recipients, 50% DAO
  uint256 public settleIndex = 0;
  uint256 public epochEmissionLocked;
  uint256 private _epochRecipientPool;
  uint256 private _epochRecipientMinted;
  uint256 public settleProgress = 0;                       // 0=not settling, >0=in progress (uint256 not bool)

  // Epoch timing (owned by AWPEmission, immutables)
  uint256 public immutable genesisTime;                   // Set in initialize
  uint256 public immutable epochDuration;                  // Set in initialize (86400 = 1 day)
  // currentEpoch() = (block.timestamp - genesisTime) / epochDuration

  address public rootNet;

  modifier onlyTimelock() { require(msg.sender == treasury); _; }

Functions:

  initialize(awpToken, treasury, rootNet, initialDailyEmission, genesisTime_, epochDuration_) external initializer
    → genesisTime = genesisTime_
    → epochDuration = epochDuration_
    // Note: oracles and threshold configured post-deploy via setOracleConfig

  // ═══════════════
  //  Oracle 权重提交（收集足够签名后自动生效）
  // ═══════════════

  submitAllocations(address[] recipients, uint96[] weights, bytes[] signatures, uint256 effectiveEpoch) external
    → EIP-712 typed-data hash = keccak256(abi.encode(ALLOCATIONS_TYPEHASH, keccak256(abi.encodePacked(recipients)),
        keccak256(abi.encodePacked(weights)), effectiveEpoch, currentNonce))
    → verify each signature recovers to a distinct oracle address
    → require(valid signature count >= threshold)
    → require(recipients.length == weights.length)
    → require(effectiveEpoch > activeEpoch, "epoch must advance")
    → _applyAllocations(recipients, weights, effectiveEpoch)
    → currentNonce++
    → emit AllocationsSubmitted(currentNonce - 1, recipients, weights, effectiveEpoch)

  _applyAllocations(address[] recipients, uint96[] weights, uint256 effectiveEpoch) internal
    → store recipients and weights into _epochRecipients[effectiveEpoch] and _epochWeights[effectiveEpoch]
    → compute and store _epochTotalWeight[effectiveEpoch] from scratch
    → activeEpoch = effectiveEpoch
    → emit GovernanceWeightUpdated(address indexed recipient, uint96 weight) per entry

  // ═══════════════
  //  紧急权重覆盖（onlyTimelock）
  // ═══════════════

  emergencySetWeight(uint256 epoch, uint256 index, address addr, uint96 weight) external onlyTimelock
    → update weights[addr] for the given epoch and recompute totalWeight
    → emit GovernanceWeightUpdated(addr, weight)

  // ═══════════════
  //  Epoch Settlement（任何人可调用；Keeper 负责调用）
  // ═══════════════

  settleEpoch(uint256 limit) external nonReentrant
    → Phase 1: 初始化（O(1)，无循环）
      if (settleProgress == 0):
        uint256 current = currentEpoch()
        require(current > settledEpoch, "epoch not yet elapsed")
        if (settledEpoch > 0): currentDailyEmission = currentDailyEmission * DECAY_FACTOR / DECAY_PRECISION
        epochEmissionLocked = min(currentDailyEmission, awpRemaining)
        _epochRecipientPool = epochEmissionLocked * EMISSION_SPLIT_BPS / 10000
        _epochRecipientMinted = 0; settleProgress = 1; settleIndex = 0

    → Phase 2: 批量铸造（每次处理 limit 个 recipients，使用 epoch-versioned 数组）
      epoch = settledEpoch + 1
      recs = _epochRecipients[epoch]; wts = _epochWeights[epoch]; tw = _epochTotalWeight[epoch]
      end = min(settleIndex + limit, recs.length)
      for i in settleIndex..end:
        addr = recs[i]
        w = wts[i]
        if w == 0 || tw == 0 || addr == address(0): continue
        awpShare = _epochRecipientPool * w / tw
        if awpShare > 0:
          _mintTo(addr, awpShare)
          _epochRecipientMinted += awpShare
          emit RecipientAWPDistributed(settledEpoch, addr, awpShare)
      settleIndex = end

    → Phase 3: 完成 — 铸造 DAO 份额
      if settleIndex >= recs.length:
        daoShare = epochEmissionLocked - _epochRecipientMinted
        if (daoShare > 0): _mintTo(treasury, daoShare)
        emit DAOMatchDistributed(settledEpoch, daoShare)
        settleProgress = 0; settledEpoch++
        emit EpochSettled(settledEpoch - 1, epochEmissionLocked, recs.length)

  // ═══════════════
  //  治理（onlyTimelock）
  // ═══════════════

  setOracleConfig(address[] oracles, uint256 threshold) external onlyTimelock
    → 更新 oracle 列表和阈值
    → emit OracleConfigUpdated(oracles, threshold)

  // ═══════════════
  //  UUPS 升级授权
  // ═══════════════

  _authorizeUpgrade(address newImpl) internal override onlyTimelock

  // ═══════════════
  //  Internal
  // ═══════════════

  _mintTo(address to, uint256 amount) internal
    → require(to != address(0) && amount > 0)
    → IAWPToken(awpToken).mint(to, amount)

  // ═══════════════
  //  View
  // ═══════════════

  getRecipientCount(uint256 epoch) view → _epochRecipients[epoch].length
  getRecipientAt(uint256 epoch, uint256 index) view → _epochRecipients[epoch][index]

Events:
  RecipientAWPDistributed(uint256 indexed epoch, address indexed recipient, uint256 awpAmount)
  DAOMatchDistributed(uint256 indexed epoch, uint256 amount)
  EpochSettled(uint256 indexed epoch, uint256 totalEmission, uint256 recipientCount)
  GovernanceWeightUpdated(address indexed recipient, uint96 weight)
  AllocationsSubmitted(uint256 indexed nonce, address[] recipients, uint96[] weights, uint256 effectiveEpoch)
  OracleConfigUpdated(address[] oracles, uint256 threshold)
```

### 4.5 AccessManager.sol — Internal Contract

```
RootNet-only callable (onlyRootNet)

Storage:
  rootNet: address

  // User registration
  mapping(address => bool) public isRegistered;
  mapping(address => uint64) public registeredAt;
  uint256 public totalUsers;

  // Agent
  agentOwner: mapping(address => address)
  userAgents: mapping(address => EnumerableSet.AddressSet)
  isManager: mapping(address => bool)
  rewardRecipients: mapping(address => address)
  MAX_AGENTS_PER_USER = 32

Functions:

  // ── User Registration ──
  register(address user) external onlyRootNet
    → require(!isRegistered[user])
    → require(agentOwner[user] == address(0), "Address is agent")
    → isRegistered[user] = true
    → registeredAt[user] = uint64(block.timestamp)
    → totalUsers++

  // ── Agent Registration (Agent calls itself, proving private key ownership) ──
  registerAgent(address agent, address user) external onlyRootNet
    → require(!isRegistered[agent], "Agent is a registered user")
    → require(isRegistered[user], "User not registered")
    → require(agentOwner[agent] == address(0), "Already bound")
    → require(agent != user)
    → require(userAgents[user].length() < MAX_AGENTS_PER_USER)
    → agentOwner[agent] = user
    → userAgents[user].add(agent)

  // ── Agent Management ──
  removeAgent(address user, address agent, address operator) external onlyRootNet
    → require(agentOwner[agent] == user)
    → require(agent != operator, "Cannot remove self")
    → delete agentOwner[agent]
    → delete isManager[agent]
    → userAgents[user].remove(agent)

  setManager(address user, address agent, bool _isManager, address operator) external onlyRootNet
    → require(agentOwner[agent] == user)
    → if (!_isManager) require(agent != operator, "Cannot revoke self")
    → isManager[agent] = _isManager

  setRewardRecipient(address user, address recipient) external onlyRootNet
    → require(recipient != address(0))
    → rewardRecipients[user] = recipient

  // ── View ──
  getOwner(address addr) external view → address
    → if (agentOwner[addr] != address(0)) return agentOwner[addr]
    → if (isRegistered[addr]) return addr
    → return address(0)

  isRegisteredUser(address) view → bool
  isRegisteredAgent(address) view → bool
  isKnownAddress(address) view → bool
  getAgents(address user) view → address[]
  isAgent(address user, address agent) view → bool
    → return agentOwner[agent] == user || (agent == user && isRegistered[user])
  isManagerAgent(address) view → bool
  getRewardRecipient(address user) view → address
    → return rewardRecipients[user] != address(0) ? rewardRecipients[user] : user
  getTotalUsers() view → uint256
```

### 4.6 StakingVault.sol — Internal Contract (Pure Allocation)

```
Write functions RootNet-only (onlyRootNet)

Core Design:
  → Pure allocation logic only. No deposit/withdraw/cooldown/STP/pending/freeze.
  → Stake bound to (user, agent, subnetId) triple
  → Allocations are plain uint128 (no Allocation struct)
  → allocate/deallocate/reallocate all take effect immediately (no pending mechanism)
  → removeAgent deallocates specified subnets immediately (no freeze epoch tracking)
  → User's stakeable balance comes from StakeNFT._userTotalStaked (queried via totalStakedOf(user))

Storage:
  rootNet: address
  stakeNFT: IStakeNFT                           // For querying user's total staked balance

  // Allocations (plain uint128, no struct)
  mapping(address => mapping(address => mapping(uint256 => uint128))) public allocations;
  // allocations[user][agent][subnetId] = amount
  mapping(address => uint128) public userTotalAllocated;

  // Subnet totals
  mapping(uint256 => uint256) public subnetTotalStake;

Functions:

  // Allocate (immediate)
  allocate(address user, address agent, uint256 subnetId, uint128 amount) external onlyRootNet
    → require(amount > 0 && getUnallocated(user) >= amount)
    → allocations[user][agent][subnetId] += amount
    → userTotalAllocated[user] += amount
    → subnetTotalStake[subnetId] += amount

  deallocate(address user, address agent, uint256 subnetId, uint128 amount) external onlyRootNet
    → require(allocations[user][agent][subnetId] >= amount)
    → allocations[user][agent][subnetId] -= amount
    → userTotalAllocated[user] -= amount
    → subnetTotalStake[subnetId] -= amount

  // Reallocate (immediate, no pending mechanism)
  reallocate(user, fromAgent, fromSubnetId, toAgent, toSubnetId, amount) external onlyRootNet
    → require(allocations[user][fromAgent][fromSubnetId] >= amount)
    → allocations[user][fromAgent][fromSubnetId] -= amount
    → allocations[user][toAgent][toSubnetId] += amount
    → subnetTotalStake[fromSubnetId] -= amount
    → subnetTotalStake[toSubnetId] += amount

  // Deallocate agent's specified subnets immediately (on removeAgent)
  deallocateAgent(address user, address agent, uint256[] subnetIds) external onlyRootNet
    → for each subnetId in subnetIds:
        amt = allocations[user][agent][subnetId]
        if amt > 0:
          allocations[user][agent][subnetId] = 0
          userTotalAllocated[user] -= amt
          subnetTotalStake[subnetId] -= amt

  // View
  getUnallocated(address user) view → stakeNFT.totalStakedOf(user) - userTotalAllocated[user]
  getAgentStake(address user, address agent, uint256 subnetId) view → uint128
  getSubnetTotalStake(uint256 subnetId) view → uint256
```

### 4.6b StakeNFT.sol — ERC721 Position NFT

```
Inherits: ERC721, ReentrancyGuard

Core Design:
  → Users deposit AWP with a lock period (timestamp-based, lockDuration in seconds)
  → Each position is an NFT with (amount, lockEndTime, createdAt)
  → NFTs are transferable (ERC721)
  → O(1) balance tracking via _userTotalStaked accumulator
  → addToPosition increases amount on existing position
  → withdraw burns NFT after lock expires

Storage:
  awpToken: IERC20
  rootNet: IRootNet

  struct Position {
    uint128 amount;
    uint64 lockEndTime;
    uint64 createdAt;
  }
  mapping(uint256 => Position) public positions;   // tokenId → Position
  mapping(address => uint256) public _userTotalStaked;  // O(1) accumulator
  uint256 private _nextTokenId = 1;

Functions:

  // Deposit AWP with lock period → mint position NFT
  deposit(uint128 amount, uint64 lockDuration) external nonReentrant → uint256 tokenId
    → require(amount > 0 && lockDuration > 0)
    → IERC20(awpToken).transferFrom(msg.sender, address(this), amount)
    → tokenId = _nextTokenId++
    → uint64 lockEndTime = uint64(block.timestamp) + lockDuration
    → positions[tokenId] = Position(amount, lockEndTime, uint64(block.timestamp))
    → _userTotalStaked[msg.sender] += amount
    → _mint(msg.sender, tokenId)
    → emit Deposited(msg.sender, tokenId, amount, lockEndTime)

  // Deposit for another address
  depositFor(address to, uint128 amount, uint64 lockDuration) external nonReentrant → uint256 tokenId
    → same as deposit but mints to `to`

  // Add more AWP to existing position (optionally extend lock)
  addToPosition(uint256 tokenId, uint128 amount, uint64 newLockEndTime) external nonReentrant
    → require(ownerOf(tokenId) == msg.sender)
    → IERC20(awpToken).transferFrom(msg.sender, address(this), amount)
    → positions[tokenId].amount += amount
    → if (newLockEndTime > positions[tokenId].lockEndTime):
        positions[tokenId].lockEndTime = newLockEndTime
    → _userTotalStaked[msg.sender] += amount
    → emit PositionIncreased(tokenId, amount, positions[tokenId].lockEndTime)

  // Withdraw after lock expires (burns NFT)
  withdraw(uint256 tokenId) external nonReentrant
    → require(ownerOf(tokenId) == msg.sender)
    → require(positions[tokenId].lockEndTime <= block.timestamp)
    → uint128 amount = positions[tokenId].amount
    → _userTotalStaked[msg.sender] -= amount
    → delete positions[tokenId]
    → _burn(tokenId)
    → IERC20(awpToken).transfer(msg.sender, amount)
    → emit Withdrawn(msg.sender, tokenId, amount)

  // Voting power for AWPDAO
  getVotingPower(uint256 tokenId) view → uint256
    → Position p = positions[tokenId]
    → if (p.lockEndTime <= block.timestamp) return 0
    → uint256 remainingTime = p.lockEndTime - block.timestamp
    → uint256 capped = min(remainingTime, 54 weeks)    // MAX_WEIGHT_SECONDS = 54 weeks
    → return p.amount * sqrt(capped / 7 days)

  // O(1) total staked query
  totalStakedOf(address user) view → uint256
    → return _userTotalStaked[user]

  // Override _update to maintain _userTotalStaked on transfer
  _update(address to, uint256 tokenId, address auth) internal override
    → address from = _ownerOf(tokenId)
    → if (from != address(0)):
        _userTotalStaked[from] -= positions[tokenId].amount
    → if (to != address(0)):
        _userTotalStaked[to] += positions[tokenId].amount
    → super._update(to, tokenId, auth)

Events:
  Deposited(address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndTime)
  PositionIncreased(uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndTime)
  Withdrawn(address indexed user, uint256 indexed tokenId, uint256 amount)
```

### 4.7 SubnetNFT.sol — Pure NFT

```
Inherits: ERC721 (~50 lines)
modifier onlyRootNet()

Storage:
  string public baseURI;                     // e.g. "https://api.awp.network/subnets/"

mint(to, tokenId) external onlyRootNet
burn(tokenId) external onlyRootNet
setBaseURI(string uri) external onlyRootNet
tokenURI(tokenId) → string.concat(baseURI, Strings.toString(tokenId))
Overrides: _update, supportsInterface
```

### 4.8 LPManager.sol — PancakeSwap V4

```
modifier onlyRootNet()

Storage:
  rootNet, poolManager, positionManager, awpToken: address

Functions:
  /// Create LP at registration (two-sided, full range, one-time)
  createPoolAndAddLiquidity(alphaToken, awpAmount, alphaAmount)
    → Create PancakeSwap V4 pool + initialize price + full-range two-sided LP
    → LP NFT stays in LPManager (permanently locked)
    → return (pool, lpTokenId)

  // No removeLiquidity — LP permanently locked
  // No collectFees — fees auto-compound
```

### 4.9 RootNet.sol — Main Contract

```
Inherits: Pausable, ReentrancyGuard, EnumerableSet

Storage:

  // Address registry
  awpToken, subnetNFT, alphaTokenFactory, awpEmission, lpManager,
  accessManager, stakingVault, stakeNFT, treasury, guardian: address
  _deployer: address                                     // Temporary; zeroed after initializeRegistry
  registryInitialized: bool

  // Note: Epoch logic has been moved to AWPEmission. RootNet no longer has genesisTime/epochDuration/currentEpoch().

  // Subnet data (on-chain stores only essential data; name/symbol/metadataURI/coordinatorURL recorded via events, stored in DB)
  struct SubnetInfo {
    address subnetContract;                            // Subnet contract (has Alpha mint permission + receives AWP emission)
    address alphaToken;
    bytes32 lpPool;                                    // Pool identifier (bytes32, not address)
    SubnetStatus status;
    uint64 createdAt; uint64 activatedAt; uint64 immunityEndsAt;
  }
  enum SubnetStatus { Pending, Active, Paused, Banned }

  struct SubnetParams {
    string name;                                       // Alpha Token name (stored in event, not on-chain)
    string symbol;                                     // Alpha Token symbol (stored in event)
    string metadataURI;                                // IPFS URI (stored in event)
    address subnetContract;                            // Subnet contract address
    string coordinatorURL;                             // Coordinator service URL (stored in event)
    string skillsURI;                                  // Skills description URI (stored in event)
    bytes32 salt;                                      // CREATE2 salt for vanity address (0 = use subnetId)
  }

  mapping(uint256 => SubnetInfo) public subnets;
  uint256 private _nextSubnetId = 1;
  uint256 public initialAlphaPrice = 1e16;            // 0.01 AWP
  uint256 public constant INITIAL_ALPHA_MINT = 100_000_000 * 1e18;  // 100M
  uint256 public immunityPeriod = 30 days;

  // Active subnet tracking (local to RootNet; AWPEmission is not notified)
  EnumerableSet.UintSet private activeSubnetIds;
  uint256 public constant MAX_ACTIVE_SUBNETS = 10000;

  // Note: Epoch logic lives entirely in AWPEmission (genesisTime, epochDuration, currentEpoch()).
  // Emission state (settledEpoch, totalWeight, etc.) lives in AWPEmission.
  // RootNet does NOT call any lifecycle functions on AWPEmission.

  // Permissions
  modifier onlyTimelock() { require(msg.sender == treasury); _; }
  modifier onlyGuardian() { require(msg.sender == guardian); _; }

Functions:

  // ═══════════════
  //  Identity Resolution
  // ═══════════════

  _resolveOwnerOrManager() internal view → address
    → if (accessManager.isRegisteredUser(msg.sender)): return msg.sender
    → owner = accessManager.getOwner(msg.sender)
    → require(owner != address(0), "Unknown address")
    → require(accessManager.isManagerAgent(msg.sender), "Not manager")
    → return owner

  _resolveOwnerOnly() internal view → address
    → require(accessManager.isRegisteredUser(msg.sender), "Not registered")
    → return msg.sender

  _validateAgent(address user, address agent) internal view
    → require(accessManager.isAgent(user, agent))

  // ═══════════════
  //  Registry
  // ═══════════════

  constructor(deployer, treasury, guardian)
    → _deployer = deployer  // Temporary; zeroed after initializeRegistry

  initializeRegistry(...) external
    → require(msg.sender == _deployer && !registryInitialized)
    → // Set all internal contract addresses
    → registryInitialized = true
    → _deployer = address(0)                                 // Permanently renounced; only DAO can update thereafter
  updateAddress(bytes32 key, address addr) external onlyTimelock
  getRegistry() external view

  // ═══════════════
  //  User Registration
  // ═══════════════

  register() external nonReentrant whenNotPaused
    → accessManager.register(msg.sender)
    → emit UserRegistered(msg.sender)

  /// @notice One-click register + deposit (via StakeNFT) + allocate
  registerAndStake(uint128 depositAmount, uint64 lockDuration, address agent, uint256 subnetId,
                   uint128 allocateAmount) external nonReentrant whenNotPaused
    → if (!accessManager.isRegisteredUser(msg.sender)):
        accessManager.register(msg.sender); emit UserRegistered(msg.sender)
    → if depositAmount > 0:
        IERC20(awpToken).transferFrom(msg.sender, stakeNFT, depositAmount)
        stakeNFT.depositFor(msg.sender, depositAmount, lockDuration)
    → if allocateAmount > 0 && agent != address(0) && subnetId > 0:
        _validateAgent(msg.sender, agent)
        stakingVault.allocate(msg.sender, agent, subnetId, allocateAmount)
        emit Allocated(msg.sender, agent, subnetId, allocateAmount, msg.sender)

  // ═══════════════
  //  Agent Registration (Agent calls itself)
  // ═══════════════

  registerAgent(address user) external nonReentrant whenNotPaused
    → accessManager.registerAgent(msg.sender, user)
    → emit AgentRegistered(user, msg.sender)

  // ═══════════════
  //  Agent Management (Owner or Manager)
  // ═══════════════

  removeAgent(address agent, uint256[] subnetIds) external nonReentrant
    → user = _resolveOwnerOrManager()
    → stakingVault.deallocateAgent(user, agent, subnetIds)
    → accessManager.removeAgent(user, agent, msg.sender)
    → emit AgentRemoved(user, agent, msg.sender)

  setManager(address agent, bool _isManager) external
    → user = _resolveOwnerOrManager()
    → accessManager.setManager(user, agent, _isManager, msg.sender)
    → emit ManagerUpdated(user, agent, _isManager, msg.sender)

  setRewardRecipient(address recipient) external
    → user = _resolveOwnerOnly()
    → accessManager.setRewardRecipient(user, recipient)
    → emit RewardRecipientUpdated(user, recipient)

  // ═══════════════
  //  Staking: Allocation (Owner or Manager)
  // ═══════════════
  // Note: Deposit/withdraw is handled by StakeNFT directly. RootNet only manages allocations.

  allocate(address agent, uint256 subnetId, uint128 amount) external
    → user = _resolveOwnerOrManager()
    → _validateAgent(user, agent)
    → stakingVault.allocate(user, agent, subnetId, amount)
    → emit Allocated(user, agent, subnetId, amount, msg.sender)

  deallocate(address agent, uint256 subnetId, uint128 amount) external
    → user = _resolveOwnerOrManager()
    → stakingVault.deallocate(user, agent, subnetId, amount)
    → emit Deallocated(user, agent, subnetId, amount, msg.sender)

  reallocate(fromAgent, fromSubnetId, toAgent, toSubnetId, amount) external
    → user = _resolveOwnerOrManager()
    → _validateAgent(user, toAgent)
    → stakingVault.reallocate(user, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)
    → emit Reallocated(user, fromAgent, fromSubnetId, toAgent, toSubnetId, amount, msg.sender)

  // ═══════════════
  //  Subnet Registration (AWP payment + auto LP)
  // ═══════════════
  // ⚠️ Prerequisite: User must first call AWPToken.approve(rootNet, lpAWPAmount)

  registerSubnet(SubnetParams calldata params) external nonReentrant whenNotPaused → uint256
    → require(bytes(params.name).length > 0 && bytes(params.name).length <= 64)
    → require(bytes(params.symbol).length > 0 && bytes(params.symbol).length <= 16)
    → require(params.subnetContract != address(0), "Subnet contract required")
    → // Does not require msg.sender to be registered
    →
    → // 1. Calculate LP creation cost
    → uint256 lpAWPAmount = INITIAL_ALPHA_MINT * initialAlphaPrice / 1e18  // 100M × 0.01 = 1M AWP
    →
    → // 2. AWP transferred directly from user to LPManager (does not pass through RootNet)
    → IERC20(awpToken).transferFrom(msg.sender, lpManager, lpAWPAmount)
    →
    → // 3. Mint NFT
    → uint256 subnetId = _nextSubnetId++
    → ISubnetNFT(subnetNFT).mint(msg.sender, subnetId)
    →
    → // 4. Deploy Alpha Token (admin = RootNet)
    → address alphaToken = IAlphaTokenFactory(alphaTokenFactory)
        .deploy(subnetId, params.name, params.symbol, address(this), params.salt)
    →
    → // 5. RootNet mints Alpha directly to LPManager (RootNet is initial minter)
    → IAlphaToken(alphaToken).mint(lpManager, INITIAL_ALPHA_MINT)
    →
    → // 6. LPManager creates LP (AWP + Alpha already in place)
    → (address pool, ) = ILPManager(lpManager)
        .createPoolAndAddLiquidity(alphaToken, lpAWPAmount, INITIAL_ALPHA_MINT)
    →
    → // 7. Set subnet contract as sole minter + RootNet gives up mint (permanently locked)
    → IAlphaToken(alphaToken).setSubnetMinter(params.subnetContract)
    →   // Internal: minters[subnetContract] = true, minters[RootNet] = false, mintersLocked = true
    →
    → // 8. Store (no strings stored on-chain; recorded via events)
    → subnets[subnetId] = SubnetInfo(
        subnetContract=params.subnetContract,
        alphaToken=alphaToken, lpPool=poolId,
        status=Pending, createdAt=now, immunityEndsAt=now+immunityPeriod)
    → // AWPEmission is NOT notified; subnet starts Pending and is not yet in activeSubnetIds
    → emit SubnetRegistered(subnetId, msg.sender, params.name, params.symbol,
          params.metadataURI, params.subnetContract, alphaToken, params.coordinatorURL)
    → emit LPCreated(subnetId, pool, lpAWPAmount, INITIAL_ALPHA_MINT)
    → return subnetId

  // ═══════════════
  //  Subnet Lifecycle
  // ═══════════════
  // ⚠️ subnetContract is immutable after registration
  // ⚠️ Alpha minter locked at registration (only subnetContract); cannot add/remove

  activateSubnet(uint256 subnetId) external
    → require(ownerOf == msg.sender && status == Pending)
    → status = Active; activatedAt = uint64(now)
    → activeSubnetIds.add(subnetId)                        // Local tracking only; no AWPEmission call
    → emit SubnetActivated(subnetId)

  pauseSubnet(uint256 subnetId) external
    → require(ownerOf == msg.sender && status == Active)
    → status = Paused
    → activeSubnetIds.remove(subnetId)                     // Local tracking only; no AWPEmission call
    → emit SubnetPaused(subnetId)

  resumeSubnet(uint256 subnetId) external
    → require(ownerOf == msg.sender && status == Paused)
    → require(activeSubnetIds.length() < MAX_ACTIVE_SUBNETS)
    → status = Active
    → activeSubnetIds.add(subnetId)                        // Local tracking only; no AWPEmission call
    → emit SubnetResumed(subnetId)

  banSubnet(uint256 subnetId) external onlyTimelock
    → require(status == Active || status == Paused)
    → address sc = subnets[subnetId].subnetContract
    → if (sc != address(0)):
        IAlphaToken(subnets[subnetId].alphaToken).setMinterPaused(sc, true)
    → activeSubnetIds.remove(subnetId)                     // Local tracking only; no AWPEmission call
    → status = Banned
    → emit SubnetBanned(subnetId)

  unbanSubnet(uint256 subnetId) external onlyTimelock
    → require(status == Banned)
    → require(activeSubnetIds.length() < MAX_ACTIVE_SUBNETS)
    → address sc = subnets[subnetId].subnetContract
    → if (sc != address(0)):
        IAlphaToken(subnets[subnetId].alphaToken).setMinterPaused(sc, false)
    → status = Active
    → activeSubnetIds.add(subnetId)                        // Local tracking only; no AWPEmission call
    → emit SubnetUnbanned(subnetId)

  deregisterSubnet(uint256 subnetId) external onlyTimelock
    → require(block.timestamp > immunityEndsAt)
    → address sc = subnets[subnetId].subnetContract
    → address alphaToken = subnets[subnetId].alphaToken
    → if (sc != address(0)):
        IAlphaToken(alphaToken).setMinterPaused(sc, true)
    → activeSubnetIds.remove(subnetId)                     // Local tracking only; no AWPEmission call
    → delete subnets[subnetId]
    → ISubnetNFT(subnetNFT).burn(subnetId)
    → emit SubnetDeregistered(subnetId)

  // ═══════════════
  //  Subnet Parameters
  // ═══════════════

  // Owner (no metadata stored on-chain; emits event for Indexer to write to DB)
  updateMetadata(uint256 subnetId, string metadataURI, string coordinatorURL) external
    → require(ownerOf == msg.sender)
    → emit MetadataUpdated(subnetId, metadataURI, coordinatorURL)
  // ⚠️ subnetContract is immutable (permanently locked at registration)

  // DAO
  setInitialAlphaPrice(uint256 price) external onlyTimelock
  setGuardian(address g) external onlyTimelock
  setImmunityPeriod(uint256 p) external onlyTimelock
  // ⚠️ setCooldownPeriod removed: no cooldown in StakeNFT model (lock period replaces cooldown)
  // ⚠️ setGovernanceWeight 已移除：权重由 AWPEmission oracle 多签（submitAllocations）或 Timelock（emergencySetWeight）管理

  // ═══════════════
  //  Emission（委托给 AWPEmission）
  // ═══════════════
  // settleEpoch(limit) 在 AWPEmission 上，任何人可调用（Keeper 负责调用 AWPEmission.settleEpoch(200)）
  // Emission 事件（RecipientAWPDistributed, DAOMatchDistributed, EpochSettled）由 AWPEmission 发出
  // setOracleConfig 在 AWPEmission（onlyTimelock）；epochDuration 现在是 AWPEmission 的 immutable
  // No setBatchSize/setMaxActiveSubnets on AWPEmission; limit is passed per-call to settleEpoch(limit)

  // ═══════════════
  //  Subnet Queries (for Coordinator)
  // ═══════════════

  getAgentInfo(address agent, uint256 subnetId) external view
    → (owner, isValid, stake, rewardRecipient)

  getAgentsInfo(address[] agents, uint256 subnetId) external view
    → Batch query

  // ═══════════════
  //  View
  // ═══════════════

  getSubnet(uint256) view
  getActiveSubnetCount() view → activeSubnetIds.length()
  getActiveSubnetIdAt(uint256 index) view → activeSubnetIds.at(index)
  isSubnetActive(uint256) view → status == Active
  nextSubnetId() view

  pause() external onlyGuardian; unpause() external onlyTimelock

Events:
  UserRegistered(address indexed user)
  AgentRegistered(address indexed user, address indexed agent)
  AgentRemoved(address indexed user, address indexed agent, address operator)
  ManagerUpdated(address indexed user, address indexed agent, bool isManager, address operator)
  RewardRecipientUpdated(address indexed user, address recipient)
  Allocated(address indexed user, address indexed agent, uint256 indexed subnetId, uint128 amount, address operator)
  Deallocated(address indexed user, address indexed agent, uint256 indexed subnetId, uint128 amount, address operator)
  Reallocated(address indexed user, address fromAgent, uint256 fromSubnet,
              address toAgent, uint256 toSubnet, uint128 amount, address operator)
  // Note: Deposited/Withdrawn/PositionIncreased events now emit from StakeNFT
  // Note: WithdrawRequested/WithdrawCancelled removed (no cooldown in StakeNFT model)
  SubnetRegistered(uint256 indexed subnetId, address indexed owner, string name,
                   string symbol, string metadataURI, address subnetContract,
                   address alphaToken, string coordinatorURL)
  LPCreated(uint256 indexed subnetId, address pool, uint256 awpAmount, uint256 alphaAmount)
  MetadataUpdated(uint256 indexed subnetId, string metadataURI, string coordinatorURL)
  SubnetActivated(uint256 indexed subnetId)
  SubnetPaused(uint256 indexed subnetId)
  SubnetResumed(uint256 indexed subnetId)
  SubnetBanned(uint256 indexed subnetId)
  SubnetUnbanned(uint256 indexed subnetId)
  SubnetDeregistered(uint256 indexed subnetId)
  // Note: GovernanceWeightUpdated, SubnetAWPDistributed, DAOMatchDistributed, EpochSettled
  // now emit from AWPEmission (see section 4.4)
```

### 4.10 AWPDAO.sol + Treasury.sol

```
AWPDAO: Inherits OZ Governor, GovernorSettings, GovernorTimelockControl.
  Overrides _getVotes and _countVote for StakeNFT-based voting (no delegate/checkpoint).
  No rootNet dependency (removed).
  → Voters submit tokenId[] arrays (StakeNFT position NFTs)
  → Voting power = amount * sqrt(min(remainingTime, 54 weeks) / 7 days)
  → Anti-manipulation: only NFTs with createdAt < proposalCreatedAt can vote (timestamp-based)
  → proposalCreatedAt (timestamp, not epoch)
  → Per-tokenId double-vote prevention (mapping(proposalId => mapping(tokenId => bool)))
  → No delegate mechanism; voting power is non-transferable except via NFT transfer
  → proposeWithTokens(targets, values, calldatas, description, tokenIds): executable proposal via Timelock
  → propose() is blocked (reverts); must use proposeWithTokens
  → signalPropose(description, tokenIds): signal-only proposal (vote-only, no Timelock, no on-chain execution)
  → Proposal lifecycle: Pending → Active → Succeeded/Defeated → Queued → Executed (proposeWithTokens)
  → Signal proposals: Pending → Active → Succeeded/Defeated (no queue/execute)

Treasury: OZ TimelockController with zero custom code (no changes)
RootNet holds no AWP and no Alpha.
AWP emission minted on demand by AWPEmission.
DAO-matched AWP minted by AWPEmission to Treasury.
```

---

## 5. Coordinator

```
Coordinator is the off-chain operations service for a subnet, deployed by the subnet Owner.

Responsibilities:
  1. Identity Verification — Listen to RootNet events to maintain Agent cache, verify Agent heartbeat signatures
  2. Task Management — Assign tasks, collect results, evaluate quality, compute contribution scores
  3. Reward Distribution — Subnet contract mints Alpha + receives AWP → distributes to miners' rewardRecipient
     → Subnet contract receives AWP each epoch (minted by AWPEmission on demand)
     → Subnet contract mints Alpha independently (has minter permission)
     → Distribution logic is entirely subnet-autonomous

How subnets query RootNet:
  Cold start: getAgentsInfo() full pull
  Running: Listen to events for incremental cache updates (AgentRegistered, Allocated, Deallocated, ...)

Different subnet types have different Coordinator logic:
  Benchmark: Generate problems / solve / verify
  DATA Mining: Collect data / verify / deduplicate
  AI Arena: Match opponents / ELO scoring
```

---

## 6. Contract Deployment

```
Step 1:  AWPToken("AWP Token", "AWP", deployer)
         → constructor mints 5B to deployer
Step 2:  (skipped — no AlphaToken impl deployment needed; CREATE2 deploys inline)
Step 3:  AlphaTokenFactory(deployer, vanityRule)                // vanityRule=0 to disable
Step 4:  Treasury(172800, [], [address(0)], deployer)
Step 5:  AWPDAO(awpToken, treasury, stakeNFT, ...)  // 6 params, no rootNet
Step 6:  Treasury.grantRole(PROPOSER+CANCELLER, awpDAO)
Step 7:  Treasury.renounceRole(ADMIN, deployer)
Step 8:  RootNet(deployer, treasury, guardian)  // No epochDuration (epoch moved to AWPEmission)
Step 9:  SubnetNFT("AWP Subnet", "CXSUB", rootNet)
Step 10: AccessManager(rootNet)
Step 11: StakingVault(rootNet, stakeNFT)
Step 11b: StakeNFT(awpToken, stakingVault, rootNet)
Step 12: LPManager(rootNet, poolManager, positionManager, awpToken)
Step 13a: AWPEmission impl = new AWPEmission()                  // 部署实现合约（不初始化）
Step 13b: initData = AWPEmission.initialize.selector(awpToken, treasury, rootNet,
                       initialDailyEmission=15_800_000e18, genesisTime_, epochDuration_=86400)
          awpEmission = new ERC1967Proxy(impl, initData)         // 部署 UUPS 代理并初始化
          // Note: oracles and threshold configured post-deploy via AWPEmission.setOracleConfig()
Step 14: AWPToken.addMinter(awpEmission)                         // 代理地址获得铸币权                        // Emission contract gets mint permission
Step 15: AWPToken.renounceAdmin()                               // Permanently lock minter list
Step 16: AlphaTokenFactory.setAddresses(rootNet)
Step 17: RootNet.initializeRegistry(all addresses, including awpEmission and stakeNFT — 10 addresses total)

// After all contracts deployed, deployer distributes non-mining portion via transfer:
Step 18: AWPToken.transfer(treasury, 2_000_000_000e18)          // 20% DAO Treasury
Step 19: AWPToken.transfer(teamVesting, 1_000_000_000e18)       // 10% Team
Step 20: AWPToken.transfer(investorVesting, 750_000_000e18)     // 7.5% Investors
Step 21: AWPToken.transfer(liquidityPool, 1_000_000_000e18)     // 10% Initial liquidity
Step 22: AWPToken.transfer(airdrop, 250_000_000e18)             // 2.5% Airdrop

// At this point:
// AWPToken minters = { awpEmission } (sole minter)
// AWPToken admin = address(0) (permanently locked)
// AWPToken totalSupply = 5B (deployer has transferred everything out)
// Remaining 5B minted by AWPEmission via exponential decay (~99% in 4 years)
// Deployer was never a minter; only received 5B in constructor for distribution via transfer
```

---

## 7. API (Go)

### Tech Stack

```
Go 1.26

Core:
  chi/v5           → HTTP router + middleware (CORS, rate-limit, recover, logger)
  pgx/v5           → PostgreSQL driver (native, not through database/sql adapter)
  sqlc             → SQL → Go code generation (compile-time type safety, zero ORM)
  go-ethereum      → On-chain interaction + abigen contract bindings

Data:
  Atlas            → Declarative DB migration (HCL schema → auto diff + apply)
  go-redis/v9      → Redis cache + Pub/Sub (inter-process communication)

Realtime:
  github.com/coder/websocket → WebSocket (context-aware, production-grade)

Infrastructure:
  log/slog         → Structured logging (Go stdlib, zero deps)
  caarlos0/env     → Env vars → struct (tag parsing)
  uber-go/fx       → Dependency injection + lifecycle management (graceful shutdown)
  robfig/cron/v3   → Keeper scheduled tasks

Observability:
  OpenTelemetry    → Traces + Metrics (Prometheus/Jaeger integration)

Code Generation:
  abigen           → Solidity ABI → Go contract bindings
  sqlc             → SQL → Go query functions + model structs

Build:
  Makefile         → generate / migrate / build / test / lint
  golangci-lint    → Static analysis
  Docker multi-stage build
```

### Architecture

```
API is a read-only service + on-chain data indexer.

Write operations handled by frontend direct-to-chain:
  Frontend wagmi/viem → useWriteContract({ address: rootNet, abi, functionName, args })
  → User signs in wallet → Sent directly to BSC
  → Does not pass through backend (no /tx/ routes)

Backend: three independent processes:
  cmd/api/main.go      → HTTP read API + WebSocket (stateless, horizontally scalable)
  cmd/indexer/main.go  → Chain Indexer (single instance, processes block events sequentially → writes to PostgreSQL)
  cmd/keeper/main.go   → Keeper Bot (single instance, executes scheduled on-chain operations)

Inter-process communication:
  Indexer → Redis Pub/Sub → API (WebSocket broadcast)
  Indexer PUBLISHES events to Redis channel "chain_events" after writing to DB
  API process SUBSCRIBES to "chain_events" → broadcasts to WebSocket clients on receipt
  → Supports horizontal scaling with multiple API instances

Shared: internal/ packages for db/, chain/, config/
Independently deployed, sharing data via PostgreSQL + Redis.
```

### DB Schema

```sql
CREATE TABLE users (
    address       CHAR(42) PRIMARY KEY,
    registered_at BIGINT NOT NULL
);

CREATE TABLE agents (
    agent_address CHAR(42) PRIMARY KEY,
    owner_address CHAR(42) NOT NULL,
    is_manager    BOOLEAN NOT NULL DEFAULT FALSE,
    removed       BOOLEAN NOT NULL DEFAULT FALSE,
    removed_at    BIGINT
);
CREATE INDEX idx_agents_owner ON agents(owner_address);
-- Note: No FK on owner_address (Agent may appear in events before User registers)
-- Note: removeAgent uses soft delete (removed=true), preserving history

CREATE TABLE user_reward_recipients (
    user_address      CHAR(42) PRIMARY KEY,
    recipient_address CHAR(42) NOT NULL
);

CREATE TABLE subnets (
    subnet_id        INTEGER PRIMARY KEY,
    owner            CHAR(42) NOT NULL,
    name             VARCHAR(64) NOT NULL,
    symbol           VARCHAR(16) NOT NULL,
    metadata_uri     TEXT,
    governance_weight INTEGER NOT NULL DEFAULT 0,
    subnet_contract  CHAR(42) NOT NULL,
    coordinator_url  TEXT,
    alpha_token      CHAR(42) NOT NULL,
    lp_pool          CHAR(42),
    status           VARCHAR(16) NOT NULL DEFAULT 'Pending',
    created_at       BIGINT NOT NULL,
    activated_at     BIGINT,
    immunity_ends_at BIGINT,
    burned           BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_subnets_owner ON subnets(owner);
CREATE INDEX idx_subnets_status ON subnets(status);

CREATE TABLE stake_allocations (
    user_address  CHAR(42) NOT NULL,
    agent_address CHAR(42) NOT NULL,
    subnet_id     INTEGER NOT NULL,
    amount        NUMERIC(78,0) NOT NULL DEFAULT 0,
    PRIMARY KEY (user_address, agent_address, subnet_id)
);
CREATE INDEX idx_sa_subnet ON stake_allocations(subnet_id);
-- Note: No FKs (Agent can be removed, subnet can be deregistered)

CREATE TABLE user_balances (
    user_address    CHAR(42) PRIMARY KEY,
    total_allocated NUMERIC(78,0) NOT NULL DEFAULT 0
);
-- Note: total_balance removed; totalStaked is computed from stake_positions

CREATE TABLE stake_positions (
    token_id      INTEGER PRIMARY KEY,
    owner         CHAR(42) NOT NULL,
    amount        NUMERIC(78,0) NOT NULL,
    lock_end_time  BIGINT NOT NULL,
    created_at     BIGINT NOT NULL
);
CREATE INDEX idx_sp_owner ON stake_positions(owner);
-- Note: withdraw_requests table removed (no cooldown in StakeNFT model)

CREATE TABLE epochs (
    epoch_id        INTEGER PRIMARY KEY,
    start_time      BIGINT NOT NULL,
    daily_emission  NUMERIC(78,0) NOT NULL,
    subnet_emission NUMERIC(78,0),
    dao_emission    NUMERIC(78,0)
);

CREATE TABLE subnet_awp_distributions (
    id         SERIAL PRIMARY KEY,
    epoch_id   INTEGER NOT NULL,
    subnet_id  INTEGER NOT NULL,
    awp_amount NUMERIC(78,0) NOT NULL
);
CREATE INDEX idx_sad_epoch ON subnet_awp_distributions(epoch_id);
CREATE INDEX idx_sad_subnet ON subnet_awp_distributions(subnet_id);

CREATE TABLE proposals (
    proposal_id   VARCHAR(66) PRIMARY KEY,
    proposer      CHAR(42) NOT NULL,
    description   TEXT,
    status        VARCHAR(16) NOT NULL,
    votes_for     NUMERIC(78,0) NOT NULL DEFAULT 0,
    votes_against NUMERIC(78,0) NOT NULL DEFAULT 0
);
CREATE INDEX idx_proposals_proposer ON proposals(proposer);

CREATE TABLE sync_states (
    contract_name VARCHAR(64) PRIMARY KEY,
    last_block    BIGINT NOT NULL DEFAULT 0
);
```

### Routes

```go
func NewRouter(h *handler.Handler, ws *ws.Hub) chi.Router {
    r := chi.NewRouter()
    r.Use(middleware.RealIP)
    r.Use(middleware.RequestID)
    r.Use(slogchi.Recovery)
    r.Use(middleware.Compress(5))
    r.Use(cors.Handler(cors.Options{...}))
    r.Use(otelchi.Middleware("api"))

    r.Route("/api", func(r chi.Router) {

        // ── System ──
        r.Get("/registry", h.GetRegistry)
        r.Get("/health", h.Health)

        // ── Users ──
        r.Route("/users", func(r chi.Router) {
            r.Get("/", h.ListUsers)                    // ?page=1&limit=20
            r.Get("/count", h.GetUserCount)
            r.Get("/{address}", h.GetUser)
        })

        // ── Address Lookup ──
        r.Get("/address/{address}/check", h.CheckAddress)
        // Returns: { isRegisteredUser, isRegisteredAgent, ownerAddress, isManager }

        // ── Agents ──
        r.Route("/agents", func(r chi.Router) {
            r.Get("/by-owner/{owner}", h.GetAgentsByOwner)      // All agents for a user
            r.Get("/by-owner/{owner}/{agent}", h.GetAgentDetail) // Single agent detail
            r.Get("/lookup/{agent}", h.LookupAgent)              // Look up owner by agent address
            r.Post("/batch-info", h.BatchAgentInfo)              // body: { agents: [], subnetId }
            // batch-info uses POST: body carries agents array; GET query string unfriendly for many addresses
        })

        // ── Staking ──
        r.Route("/staking", func(r chi.Router) {
            r.Get("/user/{address}/balance", h.GetBalance)
            // Returns: { totalStaked, totalAllocated, unallocated }
            r.Get("/user/{address}/positions", h.GetPositions)  // StakeNFT positions
            r.Get("/user/{address}/allocations", h.GetAllocations)  // ?page=1&limit=20
            r.Get("/agent/{agent}/subnet/{subnetId}", h.GetAgentSubnetStake)
            r.Get("/agent/{agent}/subnets", h.GetAgentSubnets)
            r.Get("/subnet/{subnetId}/total", h.GetSubnetTotalStake)
            // Returns: { subnetId, totalStake, stp }
        })

        // ── Subnets ──
        r.Route("/subnets", func(r chi.Router) {
            r.Get("/", h.ListSubnets)                    // ?status=Active&page=1&limit=20
            r.Get("/{subnetId}", h.GetSubnet)
            r.Get("/{subnetId}/earnings", h.GetSubnetEarnings)  // ?page=1&limit=20
            r.Get("/{subnetId}/agents/{agent}", h.GetSubnetAgentInfo)
            // GetSubnetAgentInfo: Used by Coordinator, returns owner+stake+rewardRecipient
        })

        // ── Emission ──
        r.Route("/emission", func(r chi.Router) {
            r.Get("/current", h.GetCurrentEmission)
            // Returns: { epoch, dailyEmission, decayFactor, awpTotalSupply, awpMaxSupply }
            r.Get("/schedule", h.GetEmissionSchedule)
            // Returns: Projected emissions for next 30/90/365 days
            r.Get("/epochs", h.ListEpochs)               // ?page=1&limit=20
        })

        // ── Tokens ──
        r.Route("/tokens", func(r chi.Router) {
            r.Get("/awp", h.GetAWPInfo)
            // Returns: { totalSupply, maxSupply, circulatingSupply, holders }
            r.Get("/alpha/{subnetId}", h.GetAlphaInfo)
            // Returns: { totalSupply, maxSupply, subnetContract, minterPaused }
            r.Get("/alpha/{subnetId}/price", h.GetAlphaPrice)
            // Returns: { priceInAWP, lpPool, reserve0, reserve1 } (read from PancakeSwap)
        })

        // ── Governance ──
        r.Route("/governance", func(r chi.Router) {
            r.Get("/proposals", h.ListProposals)          // ?status=Active&page=1&limit=20
            r.Get("/proposals/{proposalId}", h.GetProposal)
            r.Get("/treasury", h.GetTreasury)
            // Returns: { awpBalance, timelockDelay }
        })
    })

    // WebSocket
    r.Get("/ws/live", ws.HandleConnect)

    return r
}
```

### Chain Indexer

```go
// internal/chain/indexer.go

type Indexer struct {
    client   *ethclient.Client
    rootNet  *bindings.RootNet
    nft      *bindings.SubnetNFT
    dao      *bindings.AWPDAO
    queries  *gen.Queries
    pool     *pgxpool.Pool
    redis    *redis.Client                   // Event broadcast
    logger   *slog.Logger
}

func (idx *Indexer) Run(ctx context.Context) error {
    // 1. Read last_block from sync_states
    // 2. FilterLogs(fromBlock, toBlock) to fetch events
    // 3. Sort by block + logIndex, process sequentially
    // 4. Begin PostgreSQL transaction:
    //    a. Each event → corresponding sqlc query (Upsert/Insert)
    //    b. Update sync_states.last_block
    //    c. Commit transaction (atomicity)
    // 5. PUBLISH events to Redis channel "chain_events"
    // 6. Loop (3s interval / SubscribeNewHead)
}

// Event handling logic:
//   UserRegistered       → INSERT users
//   AgentRegistered      → INSERT agents
//   AgentRemoved         → UPDATE agents SET removed=true, removed_at=blockTime
//                          + DELETE stake_allocations WHERE user=? AND agent=? AND subnet_id IN (?)
//   ManagerUpdated       → UPDATE agents SET is_manager=?
//   RewardRecipientUpdated → UPSERT user_reward_recipients
//   Deposited (StakeNFT)  → INSERT stake_positions (tokenId, owner, amount, lockEndTime, createdAt)
//   PositionIncreased (StakeNFT) → UPDATE stake_positions SET amount += delta
//   Withdrawn (StakeNFT)  → DELETE stake_positions WHERE token_id = ?
//   StakeNFT Transfer     → UPDATE stake_positions SET owner = to
//   Allocated            → UPSERT stake_allocations + UPDATE user_balances.total_allocated
//   Deallocated          → UPDATE stake_allocations SET amount = amount - ? + UPDATE user_balances
//   Reallocated          → UPDATE stake_allocations for both from/to triples + UPDATE user_balances
//   (PendingOperationsExecuted removed — no freeze/pending mechanism)
//   SubnetRegistered     → INSERT subnets (+ name, symbol, metadata_uri, coordinator_url, alpha_token from event)
//                          immunity_ends_at = block.timestamp + RootNet.immunityPeriod()
//   LPCreated            → UPDATE subnets SET lp_pool
//   MetadataUpdated      → UPDATE subnets SET metadata_uri, coordinator_url
//   SubnetActivated      → UPDATE subnets SET status='Active', activated_at
//   SubnetPaused/Resumed → UPDATE subnets SET status
//   SubnetBanned/Unbanned → UPDATE subnets SET status
//   SubnetDeregistered   → UPDATE subnets SET burned=true
//   GovernanceWeightUpdated → UPDATE subnets SET governance_weight
//   SubnetAWPDistributed → INSERT subnet_awp_distributions
//   DAOMatchDistributed  → UPDATE epochs SET dao_emission
//   EpochSettled         → INSERT/UPDATE epochs
//   SubnetNFT Transfer   → UPDATE subnets SET owner
//   ProposalCreated      → INSERT proposals
//   VoteCast             → UPDATE proposals SET votes_for/against
//   ProposalQueued/Executed → UPDATE proposals SET status
```

### Keeper Bot

```go
// internal/chain/keeper.go

type Keeper struct {
    client      *ethclient.Client
    rootNet     *bindings.RootNet
    awpEmission *bindings.AWPEmission
    key         *ecdsa.PrivateKey
    cron     *cron.Cron
    queries  *gen.Queries                    // Query DB for pending users etc.
    redis    *redis.Client                   // Write cache (Alpha prices etc.)
    logger   *slog.Logger
}

func (k *Keeper) Start(ctx context.Context) {
    k.cron.AddFunc("@every 30s", k.trySettleEpoch)
    k.cron.AddFunc("@every 5m",  k.updateTokenPrices)
    k.cron.Start()
}

// trySettleEpoch:
//   1. Read AWPEmission.settleProgress()
//   2. If settleProgress>0 → continue calling AWPEmission.settleEpoch(200) (Phase 2/3)
//   3. If settleProgress==0 → check AWPEmission.currentEpoch() > AWPEmission.settledEpoch()
//      → if met, call AWPEmission.settleEpoch(200) (Phase 1)

// executeAllPending:
//   (Removed — no freeze/pending mechanism. removeAgent deallocates immediately.)

// updateTokenPrices:
//   1. Iterate active subnets
//   2. Read Alpha/AWP price from PancakeSwap V4 pools
//   3. Write to Redis cache (GET /tokens/alpha/{id}/price reads from Redis)
```

### WebSocket

```go
// internal/server/ws/hub.go

type Hub struct {
    clients   map[*Client]bool
    broadcast chan []byte
    mu        sync.RWMutex
    redis     *redis.Client
}

// On startup, SUBSCRIBE to Redis channel "chain_events"
// On message received, broadcast to all WebSocket clients

// Client connection: ws://host/ws/live
// Message format: { "type": "SubnetAWPDistributed", "data": { "epoch": 1, "subnetId": 3, "amount": "..." } }
// Clients can send filter messages: { "subscribe": ["SubnetAWPDistributed", "EpochSettled"] }
```

### Redis Key Spec

```
Cache (Keeper writes, API reads):
  alpha_price:{subnetId}       → JSON { priceInAWP, reserve0, reserve1, updatedAt }  TTL=10m
  awp_info                     → JSON { totalSupply, maxSupply, circulatingSupply }   TTL=1m
  emission_current             → JSON { epoch, dailyEmission, totalWeight }           TTL=30s
  subnet_total_stake:{subnetId} → string (NUMERIC)                                   TTL=1m

Pub/Sub (Indexer publishes, API subscribes):
  channel: chain_events        → JSON { type, blockNumber, txHash, data }
```

### Makefile

```makefile
.PHONY: generate migrate build test lint

generate:
	sqlc generate
	cd ../contracts && forge build
	abigen --abi out/RootNet.sol/RootNet.json --pkg bindings --out internal/chain/bindings/rootnet.go
	abigen --abi out/AWPToken.sol/AWPToken.json --pkg bindings --out internal/chain/bindings/awptoken.go
	# ... other contracts

migrate:
	atlas schema apply --env local

build:
	go build -o bin/api ./cmd/api
	go build -o bin/indexer ./cmd/indexer
	go build -o bin/keeper ./cmd/keeper

test:
	go test ./... -race -count=1

lint:
	golangci-lint run
```

---

## 8. Development Schedule

```
Phase 1 — Contracts (Week 1-2):
  Day 1:  Foundry + OZ + all interfaces + IPancakeV4
  Day 2:  AWPToken(minter) + AlphaToken(dual minter) + AlphaTokenFactory + AWPEmission
  Day 3:  SubnetNFT + AccessManager(registration+Agent+Manager)
  Day 4:  StakingVault(allocate/deallocate) + StakeNFT(deposit/withdraw/positions)
  Day 5:  LPManager(V4 LP creation)
  Day 6:  RootNet(registration+Agent+staking forwarding+subnet registration)
  Day 7:  RootNet(lifecycle+ban+settleEpoch+emission)
  Day 8:  Treasury + AWPDAO
  Day 9:  Integration.t.sol
  Day 10: Deploy.s.sol + BSC Testnet

Phase 2 — API (Week 3-4):
  Day 11: Go project init + sqlc + Atlas + abigen code generation
  Day 12: Chain client + Indexer (event listening + DB writes)
  Day 13: Read API handlers (user/agent/staking/subnet)
  Day 14: Read API handlers (emission/token/governance) + WebSocket
  Day 15: Keeper Bot (settleEpoch + token prices)
  Day 16: Redis cache + Pub/Sub (Indexer → API WebSocket)
  Day 17: OpenTelemetry + integration tests
  Day 18: Docker compose + end-to-end testing

Phase 3 — Frontend (Week 5)
```

---

## 9. Feature Checklist

| Feature | Contract | Status |
|---------|----------|--------|
| AWP 10B + voting power | AWPToken (ERC20Votes) | ✅ |
| Alpha 10B/subnet, dual minter | AlphaToken | ✅ |
| Alpha Token factory | AlphaTokenFactory (CREATE2) | ✅ |
| Subnet NFT | SubnetNFT (ERC721) | ✅ |
| User registration | AccessManager → RootNet | ✅ |
| Agent self-registration | AccessManager → RootNet | ✅ |
| User/Agent address mutual exclusion | AccessManager | ✅ |
| Manager manages Agents + allocation | RootNet + AccessManager | ✅ |
| Agent removal deallocates stake | StakingVault deallocateAgent | ✅ |
| Staking deposit/withdraw (NFT positions) | StakeNFT (ERC721) | ✅ |
| Staking allocate/deallocate (immediate) | StakingVault → RootNet | ✅ |
| Staking reallocate (immediate) | StakingVault → RootNet | ✅ |
| Staking (user,agent,subnet) triple | StakingVault | ✅ |
| NFT-based voting power | StakeNFT + AWPDAO | ✅ |
| Subnet registration (AWP + auto LP) | RootNet + LPManager | ✅ |
| LP permanently locked | LPManager | ✅ |
| Subnet Alpha minter permanently locked | AlphaToken (setSubnetMinter + mintersLocked) | ✅ |
| Subnet lifecycle + ban + deregister | RootNet | ✅ |
| governanceWeight | AWPEmission (oracle submitAllocations / Timelock emergencySetWeight) | ✅ |
| AWP emission (direct to subnet contracts) | AWPEmission | ✅ |
| AWP 50/50 split (subnet+DAO match) | AWPEmission | ✅ |
| Emission weight (oracle-managed) | AWPEmission (submitAllocations) | ✅ |
| Exponential decay emission curve | AWPEmission | ✅ |
| AWP mint on demand (no pre-mint) | AWPEmission (internal _mintTo) | ✅ |
| Batch settlement | AWPEmission | ✅ |
| ERC1363 Callback (AWP+Alpha) | AWPToken + AlphaToken | ✅ |
| ERC20Burnable (AWP+Alpha) | AWPToken + AlphaToken | ✅ |
| DAO governance | AWPDAO (Governor) | ✅ |
| DAO treasury + Timelock | Treasury (TimelockController) | ✅ |
| Guardian emergency pause | RootNet | ✅ |
| Subnet query (getAgentInfo) | RootNet | ✅ |
| Reward recipient address | AccessManager | ✅ |
| RootNet unified entry | RootNet | ✅ |
