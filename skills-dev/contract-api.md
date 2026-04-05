# AWP Smart Contract API Reference

## Chain Info
- **Network**: Base mainnet (chain ID 8453) / BSC (chain ID 56)
- **RPC**: Chain RPC endpoint (from `RPC_URL` env var)
- **Block time**: ~2 seconds (Base) / ~3 seconds (BSC)
- **EVM version**: Cancun (supports transient storage)
- **DEX**: Uniswap V4 (Base) / PancakeSwap V4 (BSC)

---

## AWPRegistry — Unified Entry Point

> EIP-712 domain name: "AWPRegistry"

### Account System (V2)
```
register()                                                    // Optional; equivalent to setRecipient(msg.sender)
bind(address target)                                          // Tree-based binding with anti-cycle check
bindFor(address user, address target, uint256 deadline, uint8 v, bytes32 r, bytes32 s)  // Gasless EIP-712
setRecipient(address recipient)                               // Set reward recipient
setRecipientFor(address user, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) // Gasless EIP-712
grantDelegate(address delegate)                               // Grant delegation to an address
revokeDelegate(address delegate)                              // Revoke delegation from an address
resolveRecipient(address addr) view                           // Walks boundTo chain to root
isRegistered(address addr) view                               // boundTo[addr] != 0 || recipient[addr] != 0
```

### Staking
> **Note:** Allocation functions (allocate, deallocate, reallocate) have moved to **AWPAllocator**. Deposit/withdraw via **veAWP**.

### Worknet Lifecycle
```
registerWorknet(WorknetParams params) → uint256 worknetId         // Anyone (costs AWP). worknetManager=0 auto-deploys WorknetManager proxy.
registerWorknetFor(address user, WorknetParams params, uint256 deadline, uint8 v, bytes32 r, bytes32 s) // Gasless (requires prior AWP approve)
registerWorknetForWithPermit(user, params, deadline, permitV, permitR, permitS, registerV, registerR, registerS) // Fully gasless (ERC-2612 permit + EIP-712)
activateWorknet(uint256 worknetId)                               // NFT Owner: Pending → Active
pauseWorknet(uint256 worknetId)                                  // NFT Owner: Active → Paused
resumeWorknet(uint256 worknetId)                                 // NFT Owner: Paused → Active
```

### Governance (Guardian only)
```
banWorknet(uint256 worknetId)                                    // Active/Paused → Banned
unbanWorknet(uint256 worknetId)                                  // Banned → Active
deregisterWorknet(uint256 worknetId)                             // Delete
setInitialAlphaPrice(uint256 price)
setGuardian(address g)
setWorknetTokenFactory(address factory)                          // Replace factory for new worknets
setWorknetManagerImpl(address impl)                              // Set/update auto-deploy impl
```

### View Functions
```
getWorknet(uint256 worknetId) → WorknetInfo                       // Lifecycle state only (lpPool, status, timestamps)
getWorknetFull(uint256 worknetId) → WorknetFullInfo               // Combined: AWPRegistry state + AWPWorkNet identity
getActiveWorknetCount() → uint256
getActiveWorknetIdAt(uint256 index) → uint256
isWorknetActive(uint256 worknetId) → bool
nextWorknetId() → uint256
getAgentInfo(address agent, uint256 worknetId) → AgentInfo
getAgentsInfo(address[] agents, uint256 worknetId) → AgentInfo[]
getRegistry() → (awpToken, awpWorkNet, worknetTokenFactory, awpEmission, lpManager, awpAllocator, veAWP, treasury, guardian)
resolveRecipient(address addr) → address                       // Walks boundTo chain to root
isRegistered(address addr) → bool                              // boundTo[addr] != 0 || recipient[addr] != 0
boundTo(address addr) → address
recipient(address addr) → address
delegates(address user, address delegate) → bool
nonces(address) → uint256
extractChainId(uint256 worknetId) pure → uint256       // Extract chainId from global worknetId
extractLocalId(uint256 worknetId) pure → uint256       // Extract local counter from global worknetId
```

> **WorknetId encoding**: WorknetId is now globally unique across chains: `(block.chainid << 64) | localCounter`. Use `extractChainId()` and `extractLocalId()` to decode.

### Emergency
```
pause()     // Guardian only
unpause()   // Timelock only
```

---

## AWPEmission — Emission Engine (UUPS Proxy)

### Guardian Submission
```
submitAllocations(address[] recipients, uint96[] weights, uint256 effectiveEpoch)
// onlyGuardian (cross-chain multisig)
// effectiveEpoch must be >= settledEpoch
// Full replacement: clears old allocations for that epoch, writes new ones
// Allowed during settlement (epoch-versioned design)
// 100% emission to recipients; Guardian includes treasury in recipients for DAO share
```

### Epoch Settlement
```
settleEpoch(uint256 limit)
// Anyone can call. Processes up to `limit` recipients per call.
// Phase 1 (first call): initialize epoch, snapshot weights, check settledEpoch <= currentEpoch()
// Phase 2 (each call): mint AWP to recipients proportionally via mintAndCall
// Phase 3 (final call): advance settledEpoch
```

### Governance (Guardian only)
```
upgradeToAndCall(address newImpl, bytes data)                   // UUPS upgrade
```

### View Functions
```
currentEpoch() → uint256                                       // Time-based: (block.timestamp - genesisTime) / epochDuration (1 day)
genesisTime() → uint256                                        // Immutable, set at initialization
epochDuration() → uint256                                      // Immutable, 86400 (1 day)
settledEpoch() → uint256                                       // Next epoch to settle (starts at 0)
currentDailyEmission() → uint256                               // Current epoch emission (wei)
settleProgress() → uint256                                     // 0=idle, >0=in progress
epochEmissionLocked() → uint256
allocationNonce() → uint256
maxRecipients() → uint256
awpRegistry() → address
getRecipientCount() → uint256
getRecipient(uint256 index) → address
getWeight(address addr) → uint96                               // Active epoch weight (O(n) scan)
getTotalWeight() → uint256                                     // Active epoch total weight
getEpochRecipientCount(uint256 epoch) → uint256
getEpochWeight(uint256 epoch, address addr) → uint96
getEpochTotalWeight(uint256 epoch) → uint256
```

---

## Data Structures

```solidity
enum WorknetStatus { Pending, Active, Paused, Banned }

// AWPRegistry lifecycle state only — identity data in AWPWorkNet
struct WorknetInfo {
    bytes32 lpPool;            // DEX V4 PoolId (Uniswap or PancakeSwap)
    WorknetStatus status;
    uint64 createdAt;
    uint64 activatedAt;
}

// Combined: AWPRegistry state + AWPWorkNet identity
struct WorknetFullInfo {
    address worknetManager;     // Worknet manager contract (WorknetToken minter)
    address worknetToken;      // WorknetToken address
    bytes32 lpPool;
    WorknetStatus status;
    uint64 createdAt;
    uint64 activatedAt;
    string name;
    string skillsURI;
    uint128 minStake;
    address owner;
}

struct WorknetParams {
    string name;               // WorknetToken name (1-64 bytes)
    string symbol;             // WorknetToken symbol (1-16 bytes)
    address worknetManager;     // address(0) = auto-deploy WorknetManager proxy
    bytes32 salt;              // CREATE2 salt; bytes32(0) = use worknetId as salt
    uint128 minStake;          // Minimum stake for agents (0 = no minimum)
    string skillsURI;          // Skills file URI (IPFS/HTTPS)
}

struct AgentInfo {
    address boundTo;
    bool isValid;
    uint256 stake;
    address recipient;
}
```

---

## veAWP — ERC721 Position NFT

```
deposit(uint256 amount, uint64 lockDuration) → uint256 tokenId   // Deposit AWP + mint position NFT; lockDuration in seconds
depositFor(address user, uint256 amount, uint64 lockDuration) → uint256 tokenId // onlyAWPRegistry
addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime) // Add AWP and/or extend lock; reverts PositionExpired if lock expired
withdraw(uint256 tokenId)                                        // Withdraw after lock expires (burns NFT)
```

### View Functions
```
positions(uint256 tokenId) → (uint128 amount, uint64 lockEndTime, uint64 createdAt)
getUserTotalStaked(address user) → uint256                       // O(1) total staked balance
getVotingPower(uint256 tokenId) → uint256                        // amount * sqrt(min(remainingTime, 54 weeks) / 7 days)
getUserVotingPower(address user, uint256[] tokenIds) → uint256   // Sum voting power for given tokenIds
totalVotingPower() → uint256                                     // Upper-bound estimate
getPositionForVoting(uint256 tokenId) → (address owner, uint128 amount, uint64 lockEndTime, uint64 createdAt, uint64 remaining, uint256 votingPower)
remainingTime(uint256 tokenId) → uint64                          // Remaining lock time in seconds
```

---

## AWPAllocator — Allocation Logic (UUPS Proxy)

> EIP-712 domain name: "StakingVault"

```
allocate(address staker, address agent, uint256 worknetId, uint256 amount)    // Staker or delegate (reads AWPRegistry.delegates)
deallocate(address staker, address agent, uint256 worknetId, uint256 amount)  // Staker or delegate
reallocate(address staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount) // Staker or delegate
allocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) // Gasless EIP-712
deallocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) // Gasless EIP-712
```

### View Functions
```
userTotalAllocated(address staker) → uint256
getAgentStake(address staker, address agent, uint256 worknetId) → uint256
worknetTotalStake(uint256 worknetId) → uint256
getWorknetTotalStake(uint256 worknetId) → uint256
getAgentWorknets(address staker, address agent) → uint256[]
```

---

## AWPWorkNet — ERC721 with On-Chain Identity

```
// AWPRegistry-only writes
mint(address to, uint256 tokenId, string name_, address worknetManager_, address worknetToken_, uint128 minStake_, string skillsURI_) // onlyAWPRegistry
burn(uint256 tokenId)                                          // onlyAWPRegistry
setBaseURI(string uri)                                         // onlyAWPRegistry

// Owner-updatable
setSkillsURI(uint256 tokenId, string skillsURI_)               // NFT owner only
setMinStake(uint256 tokenId, uint128 minStake_)                // NFT owner only
```

### View Functions
```
getWorknetData(uint256 tokenId) → WorknetData                    // Full identity + metadata
getWorknetManager(uint256 tokenId) → address
getWorknetToken(uint256 tokenId) → address
getMinStake(uint256 tokenId) → uint128
ownerOf(uint256 tokenId) → address                             // Standard ERC721
tokenURI(uint256 tokenId) → string                             // baseURI + tokenId
```

### Events
```
SkillsURIUpdated(uint256 indexed tokenId, string skillsURI)
MinStakeUpdated(uint256 indexed tokenId, uint128 minStake)
```

---

## AWPToken — ERC20 + ERC1363 + Permit

```
mint(address to, uint256 amount)                               // Authorized minters only
mintAndCall(address to, uint256 amount, bytes data)             // Mint + ERC1363 onTransferReceived callback
transferAndCall(address to, uint256 amount, bytes data)         // Transfer + ERC1363 callback
approveAndCall(address spender, uint256 amount, bytes data)     // Approve + ERC1363 callback
burn(uint256 amount)
```

---

## WorknetManager — Default Worknet Contract (Proxy)

> Auto-deployed by AWPRegistry when `worknetManager = address(0)`. Uses AccessControl (OZ).

### Roles
```
DEFAULT_ADMIN_ROLE = 0x00                       // Grant/revoke roles; assigned to worknet registrant
MERKLE_ROLE = keccak256("MERKLE_ROLE")          // Submit Merkle roots
STRATEGY_ROLE = keccak256("STRATEGY_ROLE")      // Set/execute AWP strategy
TRANSFER_ROLE = keccak256("TRANSFER_ROLE")      // Transfer tokens out
```

### Merkle Distribution (MERKLE_ROLE)
```
setMerkleRoot(uint32 epoch, bytes32 root)                      // Submit Merkle root for an epoch
claim(uint32 epoch, uint256 amount, bytes32[] proof)           // Public; mints WorknetToken to msg.sender
isClaimed(uint32 epoch, address account) → bool
```
> Leaf = `keccak256(keccak256(abi.encode(account, amount)))`

### AWP Strategy (STRATEGY_ROLE)
```
setStrategy(AWPStrategy strategy)                              // 0=Reserve, 1=AddLiquidity, 2=BuybackBurn
executeStrategy(uint256 amount)                                // Manually execute on held AWP
```

### ERC1363 Receiver (auto-trigger)
```
onTransferReceived(operator, from, amount, data) → bytes4      // Called by AWPToken.mintAndCall; auto-executes strategy
```

### Token Transfer (TRANSFER_ROLE)
```
transferToken(address token, address to, uint256 amount)
```

### View Functions
```
worknetToken() → address
awpToken() → address
poolId() → bytes32
currentStrategy() → AWPStrategy
merkleRoots(uint32 epoch) → bytes32
```

---

## AWPDAO — Custom NFT-Based Voting

```
proposeWithTokens(address[] targets, uint256[] values, bytes[] calldatas, string description, uint256[] tokenIds) → uint256 proposalId
signalPropose(string description, uint256[] tokenIds) → uint256 proposalId
// Signal-only: no on-chain execution, skips Timelock. Vote result recorded on-chain.
// After vote succeeds, call execute(targets=[address(dao)], values=[0], calldatas=[""], descriptionHash) to finalize.
castVoteWithReasonAndParams(uint256 proposalId, uint8 support, string reason, bytes params)
// params = abi.encode(uint256[] tokenIds)
// Voting power = amount * sqrt(min(remainingTime, 54 weeks) / 7 days)
// Anti-manipulation: only NFTs with createdAt < proposalCreatedAt can vote (timestamp-based)
// Per-tokenId double-vote prevention
// castVote() and castVoteWithReason() are blocked — must use params variant
```

### View Functions
```
proposalVotes(uint256 proposalId) → (uint256 againstVotes, uint256 forVotes, uint256 abstainVotes)
hasVotedWithToken(uint256 proposalId, uint256 tokenId) → bool
isSignalProposal(uint256 proposalId) → bool
proposalCreatedAt(uint256 proposalId) → uint64                  // Timestamp when proposal was created
quorum(uint256) → uint256                                       // totalVotingPower * quorumPercent / 100
proposalThreshold() → uint256                                   // 1,000,000 AWP
```

---

## Constants

| Constant | Value |
|----------|-------|
| AWP MAX_SUPPLY | 10,000,000,000 × 10^18 |
| WorknetToken MAX_SUPPLY | 10,000,000,000 × 10^18 (per worknet) |
| INITIAL_DAILY_EMISSION | 15,800,000 × 10^18 |
| EPOCH_DURATION | 86,400 (1 day), immutable on AWPEmission |
| DECAY_FACTOR | 996844 / 1000000 per epoch |
| EMISSION_SPLIT | 100% to recipients (Guardian includes treasury) |
| MAX_ACTIVE_WORKNETS | 10,000 |
| maxRecipients | 10,000 |
| MAX_WEIGHT_SECONDS | 54 weeks (32,659,200s) |
