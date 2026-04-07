# Worknet Developer Guide — AWP

> This document is intended for worknet developers (Coordinator / Worknet Contract), providing all the information needed to interact with AWP protocol contracts and APIs.

---

## 1. Architecture Overview

```
User
  |
  +-- bind(target) --> Tree-based binding (optional, every address is implicitly a root)
  +-- veAWP.deposit(AWP, lockDuration) --> Deposit AWP, receive position NFT
  +-- awpAllocator.allocate(staker, agent, worknetId, amount) --> Allocate stake to a worknet
  +-- awpRegistry.registerWorknet(params) --> Register new worknet (deploy WorknetToken + create LP)
        |
        v
AWPRegistry (Unified Entry Point -- account system + worknet management)
  |
  +-- AWPAllocator (UUPS proxy) -- Allocation logic + gasless EIP-712 support
  +-- veAWP -- ERC721 position NFT (deposit/withdraw AWP)
  +-- AWPWorkNet -- Worknet identity NFT (on-chain metadata)
  +-- LPManager -- DEX V4 CL pool, permanently locked LP
        |
        v
AWPEmission (UUPS Proxy -- Independent Emission Engine)
  |
  +-- Guardian multi-sig --> submitAllocations(recipients[], weights[], effectiveEpoch)
  +-- settleEpoch(limit) --> Batch-mint AWP to recipients via mintAndCall
        |
        v
Worknet Contract (What you develop)
  |
  +-- Receive AWP emission (minted to your contract address each epoch)
  +-- Mint WorknetTokens (you are the sole minter)
  +-- Distribute rewards to miners (based on contribution)
```

**Key design:** AWPEmission is a generic address-to-weight distribution engine. It does not know about worknets -- it simply mints AWP to addresses proportional to their Guardian-assigned weights. Your worknet contract receives AWP because the Guardian includes your contract address in the allocation list.

---

## 2. Worknet Registration Flow

### 2.1 Prerequisites

1. **Option A**: Deploy your own worknet manager contract, OR use `address(0)` to auto-deploy the default `WorknetManager` proxy
2. Prepare AWP -- LP creation cost = `initialAlphaMint * initialAlphaPrice / 1e18` (default: 100M * 0.001 = 100,000 AWP)
3. Query current price: `awpRegistry.initialAlphaPrice()` -- calculate the actual AWP amount needed
4. Call `AWPToken.approve(awpRegistry, awpAmount)` to authorize AWPRegistry for the transfer
5. No mandatory registration needed -- every address can register worknets directly

### 2.2 Registration

Three methods are available:

**Method 1: Direct registration (user pays gas + AWP)**

```solidity
IAWPRegistry.WorknetParams memory params = IAWPRegistry.WorknetParams({
    name: "My Worknet Alpha",      // WorknetToken name (1-64 bytes)
    symbol: "MWALPHA",             // WorknetToken symbol (1-16 bytes)
    worknetManager: address(0),    // address(0) = auto-deploy WorknetManager proxy
    salt: bytes32(0),              // 0 = use worknetId as CREATE2 salt (default)
    minStake: 0,                   // Minimum stake for agents (0 = no minimum)
    skillsURI: ""                  // Skills file URI (updatable later via AWPWorkNet)
});
uint256 worknetId = awpRegistry.registerWorknet(params);
```

**Method 2: Gasless registration (`registerWorknetFor`) -- relayer pays gas, user signs EIP-712 and pays AWP**

```solidity
awpRegistry.registerWorknetFor(
    user,          // Signer address
    params,        // WorknetParams (same struct as above)
    deadline,      // Signature expiry timestamp
    v, r, s        // EIP-712 signature
);
```

**Method 3: Fully gasless (`registerWorknetForWithPermit`) -- relayer pays gas, user signs two messages (ERC-2612 permit + EIP-712), zero gas**

```solidity
awpRegistry.registerWorknetForWithPermit(
    user,          // Signer address
    params,        // WorknetParams
    deadline,      // Shared deadline for both signatures
    permitV, permitR, permitS,     // ERC-2612 AWP permit signature
    registerV, registerR, registerS // EIP-712 registerWorknet signature
);
```

**Frontend (wagmi/viem):**

```javascript
const { writeContract } = useWriteContract();
await writeContract({
  address: AWP_REGISTRY_ADDRESS,
  abi: awpRegistryABI,
  functionName: 'registerWorknet',
  args: [{
    name, symbol,
    worknetManager: '0x0000000000000000000000000000000000000000',
    salt: '0x00...00',
    minStake: 0n,
    skillsURI: ''
  }]
});
```

> **Gasless relay**: Users can also register via `POST /api/relay/register-worknet` -- the relayer pays gas, the user signs an EIP-712 message and pays AWP.

### 2.2b Vanity WorknetToken Address (optional)

If the factory was deployed with a non-zero `vanityRule`, you can pre-mine a CREATE2 salt off-chain to get a custom-looking WorknetToken address.

**Step 1 -- Understand the vanity rule**

Query `WorknetTokenFactory.vanityRule()` to see what pattern is enforced. A value of `0` means no validation is performed and you can skip this section.

The vanity rule is a `uint64` encoding 8 hex positions (4 prefix + 4 suffix). Per position: 0-9 = digit, 10-15 = lowercase a-f, 16-21 = uppercase A-F (EIP-55), >= 22 = wildcard.

**Step 2 -- Get a salt via API (recommended)**

The API maintains a pre-mined salt pool. Request a salt:
```bash
curl -X POST https://tapi.awp.sh/api/vanity/compute-salt
# Response: {"salt": "0x530c...", "address": "0xa1...cafe", "source": "pool", "elapsed": "1ms"}
```

Or mine offline using parameters from the API:
```bash
# Get mining params
curl https://tapi.awp.sh/api/vanity/mining-params
# {"factoryAddress": "0xAe8E...", "initCodeHash": "0xec76...", "vanityRule": "0x0A01..."}

# Mine with cast
cast create2 --deployer <factoryAddress> --init-code-hash <initCodeHash> \
  --starts-with a1 --ends-with cafe --case-sensitive

# Upload to pool
curl -X POST https://tapi.awp.sh/api/vanity/upload-salts \
  -H "Content-Type: application/json" \
  -d '{"salts": [{"salt": "0x...", "address": "0x..."}]}'
```

**Step 3 -- Pass the salt in WorknetParams**

```solidity
IAWPRegistry.WorknetParams memory params = IAWPRegistry.WorknetParams({
    // ... other fields ...
    salt: 0xYourMinedSalt000000000000000000000000000000000000000000000000,
    // ...
});
```

The factory will deploy the WorknetToken at the vanity address and then validate it matches the configured rule. If validation fails, the transaction reverts with `InvalidVanityAddress()`.

### 2.3 Automatic Steps After Registration

| Step | Action | Result |
|------|--------|--------|
| 1 | Calculate LP cost | `initialAlphaMint * initialAlphaPrice / 1e18` |
| 2 | Transfer AWP from user to LPManager | `safeTransferFrom(user, lpManager, awpAmount)` |
| 3 | Deploy WorknetToken (CREATE2) | Standalone contract via factory; salt from params (or worknetId if 0); admin = AWPRegistry |
| 4 | Create DEX V4 CL pool | Full-range liquidity, 1% fee, permanently locked |
| 5 | Auto-deploy WorknetManager proxy (if worknetManager=0) | ERC1967Proxy with defaultWorknetManagerImpl, admin = user |
| 6 | Set worknet manager as sole WorknetToken minter | `setWorknetMinter(sc)` permanently locked |
| 7 | Mint AWPWorkNet with identity data | name, worknetManager, worknetToken, minStake, skillsURI stored on-chain |
| 8 | Store lifecycle state | lpPool, status=Pending, createdAt |

### 2.4 Activate Worknet

After registration, the worknet status is `Pending`. The NFT Owner must manually activate it:

```solidity
awpRegistry.activateWorknet(worknetId);
// Status: Pending -> Active
// Added to AWPRegistry's activeWorknetIds set
```

Gasless activation is also available via `activateWorknetFor(user, worknetId, deadline, v, r, s)` or `POST /api/relay/activate-worknet`.

**Note:** Activation makes the worknet visible as "Active" in the protocol, but it does NOT automatically include it in AWP emission. To receive emissions, the Guardian multi-sig must include your worknet contract address in their `submitAllocations` submission.

### 2.5 Worknet Lifecycle

```
Pending --> Active    (activateWorknet by NFT owner)
Active  --> Paused    (pauseWorknet by NFT owner)
Paused  --> Active    (resumeWorknet by NFT owner)
Active  --> Banned    (banWorknet by Guardian)
Banned  --> Active    (unbanWorknet by Guardian, checks MAX_ACTIVE_WORKNETS)
Active  --> Deregistered  (deregisterWorknet by Guardian)
Paused  --> Deregistered  (deregisterWorknet by Guardian)
Banned  --> Deregistered  (deregisterWorknet by Guardian)
```

### 2.6 Multi-Chain

WorknetId is globally unique: `(block.chainid << 64) | localCounter`. Users can allocate stake on their staking chain to any chain's worknet. Each chain has its own AWPEmission, and the Guardian coordinates emission quotas across chains.

---

## 3. WorknetManager -- Default Worknet Contract API

> When `worknetManager = address(0)` is passed to `registerWorknet`, AWPRegistry auto-deploys a `WorknetManager` proxy. The registrant (`msg.sender` or the EIP-712 signer for gasless) becomes `DEFAULT_ADMIN_ROLE`.

### 3.0 Roles

| Role | Bytes32 | Purpose |
|------|---------|---------|
| `DEFAULT_ADMIN_ROLE` | `0x00` | Grant/revoke all roles. Assigned to worknet registrant at deployment. |
| `MERKLE_ROLE` | `keccak256("MERKLE_ROLE")` | Submit Merkle roots for WorknetToken distribution |
| `STRATEGY_ROLE` | `keccak256("STRATEGY_ROLE")` | Choose AWP handling strategy + manually execute |
| `TRANSFER_ROLE` | `keccak256("TRANSFER_ROLE")` | Transfer any ERC20 held by the contract |

**Post-registration setup:**
```solidity
// Grant operator roles (admin calls these)
IAccessControl(worknetManager).grantRole(MERKLE_ROLE, merkleOperator);
IAccessControl(worknetManager).grantRole(STRATEGY_ROLE, strategyOperator);
IAccessControl(worknetManager).grantRole(TRANSFER_ROLE, treasuryMultisig);
```

### 3.1 Merkle Distribution (MERKLE_ROLE)

Distribute WorknetTokens to users via Merkle proofs. The contract mints directly -- no pre-funding needed.

```
setMerkleRoot(uint32 epoch, bytes32 root)     // Submit root for an epoch (e.g. 20260318)
```

**Building the Merkle tree (off-chain):**
```javascript
// Leaf format: keccak256(keccak256(abi.encode(account, amount)))
import { StandardMerkleTree } from "@openzeppelin/merkle-tree";

const values = [
  ["0xAlice...", "1000000000000000000000"],  // 1000 WorknetToken
  ["0xBob...",   "500000000000000000000"],   // 500 WorknetToken
];
const tree = StandardMerkleTree.of(values, ["address", "uint256"]);
const root = tree.root;
// Submit: worknetManager.setMerkleRoot(20260318, root);

// Get proof for Alice:
const proof = tree.getProof(["0xAlice...", "1000000000000000000000"]);
```

**User claiming:**
```
claim(uint32 epoch, uint256 amount, bytes32[] proof)  // Anyone; mints WorknetToken to resolved recipient
isClaimed(uint32 epoch, address account) -> bool       // Check claim status
```

The `claim` function resolves the recipient via `awpRegistry.resolveRecipient(msg.sender)`, walking the bind chain to the root address.

```javascript
// Frontend claim
await worknetManager.claim(20260318, parseEther("1000"), proof);
```

### 3.2 AWP Strategy (STRATEGY_ROLE)

Configure how incoming AWP emission is handled.

```
setStrategy(AWPStrategy strategy)              // Set the active strategy
executeStrategy(uint256 amount)                // Manually execute on AWP held by contract
```

| Strategy | Enum Value | Behavior |
|----------|-----------|----------|
| `Reserve` | 0 | AWP stays in contract; distribute manually or via Merkle |
| `AddLiquidity` | 1 | Add AWP as single-sided liquidity (0 to current price range) |
| `BuybackBurn` | 2 | Swap AWP to WorknetToken via DEX V4, then burn WorknetToken |

```javascript
// Set strategy to BuybackBurn
await worknetManager.setStrategy(2);

// AWP emission now auto-triggers BuybackBurn via ERC1363 onTransferReceived
// No manual executeStrategy needed -- AWPEmission.settleEpoch uses mintAndCall
```

**Auto-trigger flow (ERC1363):**
```
AWPEmission.settleEpoch(limit)
  -> awpToken.mintAndCall(worknetManager, amount, "")
    -> WorknetManager.onTransferReceived()
      -> if strategy == AddLiquidity -> DEX V4 single-sided LP
      -> if strategy == BuybackBurn  -> swap AWP to WorknetToken, burn WorknetToken
      -> if strategy == Reserve      -> no-op (AWP stays in contract)
```

**Additional configuration:**
- `setSlippageTolerance(uint256 bps)` -- set slippage for buyback swaps (default 500 = 5%, max 5000 = 50%)
- `setStrategyPaused(bool paused)` -- emergency pause strategy execution (DEFAULT_ADMIN_ROLE)
- `setMinStrategyAmount(uint256 amount)` -- minimum AWP for strategy execution; below this, onTransferReceived defaults to Reserve

### 3.3 Token Transfer (TRANSFER_ROLE)

```
transferToken(address token, address to, uint256 amount)  // Transfer any ERC20 out
batchTransferToken(address token, address[] recipients, uint256[] amounts)  // Batch transfer
```

```javascript
// Send reserved AWP to a multisig
await worknetManager.transferToken(AWP_TOKEN, multisig, parseEther("10000"));
```

### 3.4 View Functions

```
worknetToken() -> address                        // WorknetToken for this worknet
awpToken() -> address                          // AWP token
poolId() -> bytes32                            // DEX V4 pool ID
currentStrategy() -> AWPStrategy               // Current AWP handling strategy (0/1/2)
merkleRoots(uint32 epoch) -> bytes32           // Merkle root for an epoch
claimed(uint32 epoch, address account) -> bool // Alias for isClaimed
slippageBps() -> uint256                       // Current slippage tolerance in bps
strategyPaused() -> bool                       // Whether strategy execution is paused
minStrategyAmount() -> uint256                 // Minimum AWP for auto-strategy
```

### 3.5 Complete Operator Workflow

```
1. Register worknet -> auto-deploy WorknetManager -> you are admin
2. Grant roles:
   - MERKLE_ROLE to your backend that computes reward distributions
   - STRATEGY_ROLE to your governance multisig
   - TRANSFER_ROLE to your treasury multisig
3. Set AWP strategy (default: Reserve):
   - worknetManager.setStrategy(2)  // BuybackBurn
4. Each epoch, AWP emission arrives automatically via mintAndCall
5. Periodically compute Merkle tree of WorknetToken rewards:
   - worknetManager.setMerkleRoot(epoch, root)
6. Users claim WorknetTokens:
   - worknetManager.claim(epoch, amount, proof)
7. Monitor via WebSocket:
   - Subscribe to: SkillsURIUpdated, MinStakeUpdated, Allocated, Deallocated
8. Update worknet settings (NFT owner only):
   - awpWorkNet.setSkillsURI(worknetId, "ipfs://...")
   - awpWorkNet.setMinStake(worknetId, 100e18)
   - awpWorkNet.setMetadataURI(worknetId, "https://...")
```

---

## 4. Worknet Contract Development

### 4.1 Default WorknetManager (auto-deployed)

If you pass `worknetManager = address(0)` during registration, AWPRegistry auto-deploys a **WorknetManager** proxy with built-in features:

- **Merkle Distribution**: Submit Merkle roots per epoch, users claim WorknetTokens via proof
- **AWP Strategy**: Choose how incoming AWP is handled:
  - `Reserve` -- keep in contract for manual distribution
  - `AddLiquidity` -- add single-sided liquidity to DEX V4
  - `BuybackBurn` -- buy WorknetToken from pool and burn it
- **ERC1363 Receiver**: AWP received via `mintAndCall` (from AWPEmission) auto-triggers the current strategy
- **Token Transfer**: Move any ERC20 held by the contract (single or batch)
- **AccessControl**: Three roles -- `MERKLE_ROLE`, `STRATEGY_ROLE`, `TRANSFER_ROLE` (you are the admin)
- **UUPS Upgradeable**: The worknet owner (DEFAULT_ADMIN_ROLE) can upgrade the implementation

After registration, grant roles to your operator addresses:
```solidity
IAccessControl(worknetManager).grantRole(MERKLE_ROLE, operatorAddress);
IAccessControl(worknetManager).grantRole(STRATEGY_ROLE, operatorAddress);
```

### 4.2 Custom Worknet Contract (advanced)

If you need custom logic, deploy your own contract and pass its address as `worknetManager`:

```solidity
contract MyWorknetContract is IERC1363Receiver {
    IWorknetToken public worknetToken;
    IERC20 public awpToken;

    // AWP emission arrives via mintAndCall -- implement IERC1363Receiver to auto-process
    function onTransferReceived(address, address, uint256 amount, bytes calldata)
        external returns (bytes4) {
        // Your custom AWP handling logic
        return IERC1363Receiver.onTransferReceived.selector;
    }

    // Mint WorknetToken to miners (you are the sole minter after setWorknetMinter)
    function mintWorknetTokenToMiner(address miner, uint256 amount) external {
        worknetToken.mint(miner, amount);
    }
}
```

**Important:** After registration, `setWorknetMinter` permanently locks the WorknetToken minter to your worknet contract. This cannot be changed. If your contract has a bug, use `minterPaused` (set by AWPRegistry admin) to halt minting.

### 4.3 AWPWorkNet -- On-Chain Identity

Each worknet is represented by an NFT (`tokenId = worknetId`) that stores both immutable and updatable data.

**Immutable fields (set at mint, never changed):**
- `name` -- worknet / WorknetToken name
- `worknetManager` -- worknet contract address
- `worknetToken` -- WorknetToken address

**Owner-updatable fields:**
- `skillsURI` -- skills file URI for agent discovery (`setSkillsURI`)
- `minStake` -- minimum stake requirement for agents (`setMinStake`). Stored on-chain but NOT enforced by AWPAllocator.allocate; used as off-chain/coordinator reference only.
- `metadataURI` -- custom metadata JSON URI (`setMetadataURI`). Overrides on-chain JSON generation.

**tokenURI resolution (3-tier):**
1. Per-token `metadataURI` (if set by worknet owner)
2. Global `baseURI` (if set by AWPRegistry governance)
3. On-chain generated Base64 JSON with attributes (Worknet Manager, WorknetToken, Min Stake, Chain ID, Local ID)

### 4.4 Key Permissions

| Permission | Who Holds It | Description |
|------------|-------------|-------------|
| WorknetToken minting | **Your worknet contract** | `setWorknetMinter` permanently locked -- only you can mint WorknetToken |
| WorknetToken minting pause | AWPRegistry (admin) | `setMinterPaused(true)` pauses your minting when banned |
| AWP emission receipt | Guardian-assigned | Your address must be in the Guardian's `submitAllocations` list to receive AWP |
| WorknetToken MAX_SUPPLY | 10B | Independent cap per worknet |

### 4.5 Ban and Unban

- **Ban**: Guardian calls `banWorknet(worknetId)` -- your WorknetToken minting is paused, worknet removed from AWPRegistry's active list
- **Unban**: Treasury calls `unbanWorknet(worknetId)` -- status restores to **Active**, re-added to active list (checks MAX_ACTIVE_WORKNETS = 10,000)
- During a ban, your existing AWP and WorknetToken holdings are **unaffected** -- only new WorknetToken minting is blocked
- Guardian typically excludes banned worknets from their allocation submissions

### 4.6 Deregistration

- Guardian calls `deregisterWorknet(worknetId)`
- Users must manually deallocate from deregistered worknets (deallocate has no status check)
- Frontend should alert users on `WorknetDeregistered` events
- AWPWorkNet is burned on deregistration

---

## 5. API Reference

> Base URL: `https://tapi.awp.sh` (or your self-hosted address)

### 5.1 System

#### `GET /api/health`
```json
{"status": "ok"}
```

#### `GET /api/registry`
Returns all protocol contract addresses (excludes implementation contracts), plus EIP-712 domain info:
```json
{
  "awpRegistry": "0x...",
  "awpToken": "0x...",
  "awpEmission": "0x...",
  "awpAllocator": "0x...",
  "veAWP": "0x...",
  "awpWorkNet": "0x...",
  "lpManager": "0x...",
  "worknetTokenFactory": "0x...",
  "dao": "0x...",
  "treasury": "0x...",
  "chainId": 8453,
  "eip712Domain": {
    "name": "AWPRegistry",
    "version": "1",
    "chainId": 8453,
    "verifyingContract": "0x..."
  }
}
```

### 5.2 Users

#### `GET /api/users/{address}`
```json
{
  "user": { "address": "0x...", "registered_at": 1710000000 },
  "balance": { "user_address": "0x...", "total_staked": "5000", "total_allocated": "3000" },
  "rewardRecipient": { "user_address": "0x...", "recipient_address": "0x..." }
}
```

#### `GET /api/users/count`
```json
{"count": 1234}
```

#### `GET /api/users/?page=1&limit=20`
Paginated user list.

#### `GET /api/address/{address}/check`
```json
{
  "isRegistered": true,
  "boundTo": "0x...",
  "recipient": "0x..."
}
```

### 5.3 Nonces

#### `GET /api/nonce/{address}`
AWPRegistry nonce (for bind/setRecipient/registerWorknet EIP-712 signatures).

#### `GET /api/staking-nonce/{address}`
AWPAllocator nonce (for allocate/deallocate EIP-712 signatures).

### 5.4 Staking

#### `GET /api/staking/user/{address}/balance`
```json
{
  "totalStaked": "10000000000000000000000",
  "totalAllocated": "5000000000000000000000",
  "unallocated": "5000000000000000000000"
}
```
> `totalStaked` is computed from veAWP positions. No `withdrawRequest` field (lock period replaces cooldown).

#### `GET /api/staking/user/{address}/positions`
```json
[
  {"token_id": 1, "amount": "5000000000000000000000", "lock_end_time": 1710604800, "created_at": 1710000000},
  {"token_id": 7, "amount": "5000000000000000000000", "lock_end_time": 1713196800, "created_at": 1710345600}
]
```

#### `GET /api/staking/user/{address}/allocations?page=1&limit=20`
```json
[
  {"user_address": "0x...", "agent_address": "0x...", "worknet_id": 1, "amount": "5000", "frozen": false}
]
```

#### `GET /api/staking/agent/{agent}/worknet/{worknetId}`
**Commonly used by Coordinators** -- query total stake of an agent on a specific worknet:
```json
{"amount": "5000000000000000000000"}
```

#### `GET /api/staking/agent/{agent}/worknets`
```json
[{"worknet_id": 1, "amount": "5000"}, {"worknet_id": 3, "amount": "2000"}]
```

#### `GET /api/staking/worknet/{worknetId}/total`
```json
{"total": "50000000000000000000000"}
```

### 5.5 Worknets

#### `GET /api/worknets/?status=Active&page=1&limit=20`
```json
[
  {
    "worknet_id": 1,
    "owner": "0x...",
    "name": "My Worknet",
    "symbol": "MWRK",
    "worknet_contract": "0x...",
    "alpha_token": "0x...",
    "lp_pool": "0x...",
    "status": "Active",
    "created_at": 1710000000,
    "activated_at": 1710000100
  }
]
```

#### `GET /api/worknets/{worknetId}`
Single worknet detail.

#### `GET /api/worknets/{worknetId}/earnings?page=1&limit=20`
Worknet AWP emission history (queried by worknet contract address from `recipient_awp_distributions`):
```json
[
  {"epoch_id": 5, "recipient": "0x1234...", "awp_amount": "7900000000000000000000000"}
]
```

#### `GET /api/worknets/{worknetId}/agents/{agent}`
**Used by Coordinators** -- query agent stake info on a worknet:
```json
{"agent": "0x...", "worknetId": 1, "stake": "5000000000000000000000"}
```

### 5.6 Emission

#### `GET /api/emission/current`
Read from Redis cache (updated by Keeper every 30 seconds):
```json
{"epoch": "42", "dailyEmission": "15000000000000000000000000", "totalWeight": "5000"}
```

#### `GET /api/emission/schedule`
Projected emissions for 30/90/365 days:
```json
{
  "currentDailyEmission": "15800000000000000000000000",
  "projections": [
    {"days": 30, "totalEmission": "452000000...", "finalDailyRate": "14300000..."},
    {"days": 90, "totalEmission": "1200000000...", "finalDailyRate": "12000000..."},
    {"days": 365, "totalEmission": "3500000000...", "finalDailyRate": "5000000..."}
  ]
}
```

#### `GET /api/emission/epochs?page=1&limit=20`
```json
[
  {"epoch_id": 42, "start_time": 1710000000, "daily_emission": "15000000..."}
]
```

### 5.7 Tokens

#### `GET /api/tokens/awp`
```json
{"totalSupply": "5015800000000000000000000000", "maxSupply": "10000000000000000000000000000"}
```

#### `GET /api/tokens/alpha/{worknetId}`
```json
{"worknetId": 1, "name": "My Worknet Alpha", "symbol": "MWALPHA", "worknetToken": "0x..."}
```

#### `GET /api/tokens/alpha/{worknetId}/price`
```json
{"priceInAWP": "0.015", "reserve0": "...", "reserve1": "...", "updatedAt": "..."}
```

### 5.8 Governance

#### `GET /api/governance/proposals?status=Active&page=1&limit=20`
#### `GET /api/governance/proposals/{proposalId}`
#### `GET /api/governance/treasury`
```json
{"treasuryAddress": "0x..."}
```

### 5.9 Relay (Gasless Transactions)

> Rate limit: 100 requests per IP per 1 hour (configurable via Redis). Requires `RELAYER_PRIVATE_KEY` configured on the server.

#### `POST /api/relay/bind`
Relayer submits `bindFor()` on-chain (tree-based binding, user signs EIP-712, relayer pays gas):
```json
// Request
{"user": "0x...", "target": "0x...", "deadline": 1742400000, "signature": "0x...130 hex chars"}
// Response
{"txHash": "0x..."}
```

#### `POST /api/relay/set-recipient`
Relayer submits `setRecipientFor()` on-chain for the user:
```json
// Request
{"user": "0x...", "recipient": "0x...", "deadline": 1742400000, "signature": "0x...130 hex chars"}
// Response
{"txHash": "0x..."}
```

#### `POST /api/relay/register-worknet`
Fully gasless worknet registration -- user signs two messages (ERC-2612 permit + EIP-712), relayer pays all gas:
```json
// Request
{"user": "0x...", "name": "EVO Alpha", "symbol": "EVO",
 "worknetManager": "0x0000...0000",
 "salt": "0x...", "minStake": "0", "skillsURI": "", "deadline": 1742400000,
 "permitSignature": "0x...(ERC-2612 AWP permit)",
 "registerSignature": "0x...(EIP-712 registerWorknet)"}
// Response
{"txHash": "0x..."}
```

#### `POST /api/relay/allocate`
Gasless allocation via AWPAllocator EIP-712:
```json
{"user": "0x...", "agent": "0x...", "worknetId": "1", "amount": "5000", "deadline": 1742400000, "signature": "0x..."}
```

#### `POST /api/relay/deallocate`
Gasless deallocation via AWPAllocator EIP-712:
```json
{"user": "0x...", "agent": "0x...", "worknetId": "1", "amount": "5000", "deadline": 1742400000, "signature": "0x..."}
```

#### `POST /api/relay/activate-worknet`
Gasless worknet activation:
```json
{"user": "0x...", "worknetId": "1", "deadline": 1742400000, "signature": "0x..."}
```

### 5.10 Vanity

#### `GET /api/vanity/mining-params`
```json
{"factoryAddress": "0x...", "initCodeHash": "0x...", "vanityRule": "0x..."}
```

#### `POST /api/vanity/compute-salt`
Request a pre-mined vanity salt from the pool (DB pool first, `cast` fallback):
```json
{"salt": "0x...", "address": "0x...", "source": "pool", "elapsed": "1ms"}
```

#### `POST /api/vanity/upload-salts`
Upload mined salts to the pool (rate limited: 5/hr/IP).

#### `GET /api/vanity/salts`
List available salts.

#### `GET /api/vanity/salts/count`
Count of available salts in pool.

### 5.11 WebSocket

#### `WS /ws/live`

Connect to receive real-time on-chain event push:

```javascript
const ws = new WebSocket('wss://tapi.awp.sh/ws/live');

// Optional: set event filters (only receive events you're interested in)
ws.send(JSON.stringify({
  subscribe: ["RecipientAWPDistributed", "EpochSettled", "Allocated", "Deallocated", "AllocationsSubmitted"]
}));

// Receive events
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  // data = { type: "RecipientAWPDistributed", blockNumber: 12345, txHash: "0x...", data: { ... } }
};
```

**Event Types:**

| Event | Data | Coordinator Relevant |
|-------|------|:---:|
| `Bound` | `{user, target, oldTarget}` | Yes |
| `RecipientUpdated` | `{user, recipient}` | |
| `DelegateGranted` | `{user, delegate}` | Yes |
| `DelegateRevoked` | `{user, delegate}` | Yes |
| `Deposited` (veAWP) | `{user, tokenId, amount, lockEndTime}` | |
| `PositionIncreased` (veAWP) | `{tokenId, addedAmount, newLockEndTime}` | |
| `Withdrawn` (veAWP) | `{user, tokenId, amount}` | |
| `Allocated` | `{user, agent, worknetId, amount}` | Yes |
| `Deallocated` | `{user, agent, worknetId, amount}` | Yes |
| `Reallocated` | `{user, fromAgent, fromWorknet, toAgent, toWorknet, amount}` | Yes |
| `WorknetRegistered` | `{worknetId, owner, name, symbol, worknetManager, worknetToken}` | |
| `LPCreated` | `{worknetId, poolId, awpAmount, alphaAmount}` | |
| `SkillsURIUpdated` | `{worknetId, skillsURI}` | Yes |
| `MinStakeUpdated` | `{worknetId, minStake}` | Yes |
| `WorknetActivated` | `{worknetId}` | Yes |
| `WorknetPaused` | `{worknetId}` | Yes |
| `WorknetResumed` | `{worknetId}` | Yes |
| `WorknetBanned` | `{worknetId}` | Yes |
| `WorknetUnbanned` | `{worknetId}` | Yes |
| `WorknetDeregistered` | `{worknetId}` | Yes |
| `RecipientAWPDistributed` | `{epoch, recipient, awpAmount}` | Yes |
| `EpochSettled` | `{epoch, totalEmission, recipientCount}` | Yes |
| `AllocationsSubmitted` | `{nonce, recipients, weights}` | Yes |

---

## 6. Contract Interaction Reference

### 6.1 On-Chain Reads (view functions, no gas)

```solidity
// Query full worknet info (combines AWPRegistry state + AWPWorkNet identity)
IAWPRegistry.WorknetFullInfo memory worknet = awpRegistry.getWorknetFull(worknetId);
// worknet.worknetManager, worknet.worknetToken, worknet.status, worknet.owner, worknet.minStake

// Query current epoch
uint256 epoch = awpEmission.currentEpoch();

// Check if worknet is active
bool active = awpRegistry.isWorknetActive(worknetId);

// Query emission weight
uint96 weight = awpEmission.getWeight(worknetManagerAddress);
uint256 totalWeight = awpEmission.getTotalWeight();

// AWPWorkNet: get full worknet data
AWPWorkNet.WorknetData memory data = awpWorkNet.getWorknetData(worknetId);

// AWPWorkNet: update skills URI or minStake (NFT owner only)
awpWorkNet.setSkillsURI(worknetId, "ipfs://QmNewSkills...");
awpWorkNet.setMinStake(worknetId, 100e18); // 100 AWP minimum
awpWorkNet.setMetadataURI(worknetId, "https://metadata.example.com/1");
```

### 6.2 Recommended Coordinator Data Sync Pattern

```
Cold Start:
  1. GET /api/worknets/{worknetId} --> Worknet basic info
  2. GET /api/staking/worknet/{worknetId}/total --> Total stake on worknet
  3. GET /api/staking/agent/{agent}/worknet/{worknetId} --> Per-agent stake

Runtime Incremental Updates:
  1. WebSocket subscribe: ["Bound", "Allocated", "Deallocated", "Reallocated",
                           "DelegateGranted", "DelegateRevoked"]
  2. Update local agent cache on each event
  3. Periodically GET /api/staking/worknet/{worknetId}/total to verify total stake

Reward Distribution:
  1. Calculate reward ratio based on agent stake and contribution score
  2. Resolve recipients: awpRegistry.resolveRecipient(agent) or
     awpRegistry.batchResolveRecipients(agents)
  3. Transfer AWP/WorknetToken to the resolved recipient addresses
```

### 6.3 Emission Timeline

```
Each Epoch (1 day / daily):
  1. Keeper (or anyone) calls AWPEmission.settleEpoch(limit)
  2. Processes up to `limit` recipients per call, using mintAndCall to send AWP
  3. If your contract implements IERC1363Receiver, onTransferReceived is called automatically
  4. WorknetManager auto-executes the configured AWP strategy (Reserve/AddLiquidity/BuybackBurn)
  5. Multiple calls may be needed if there are more recipients than the limit

Your worknet receives: epochEmission * (yourWeight / totalWeight)

Decay: each epoch emission *= 0.996844 (~99% released in 4 years)
```

**How your worknet gets weight:**
- Guardian multi-sig submits epoch-versioned packed allocations via `submitAllocations(recipients[], weights[], effectiveEpoch)`
- 100% of emission goes to recipients by weight; Guardian includes treasury address for DAO share
- Your worknet contract address must be included in the Guardian's recipient list to receive AWP

---

## 7. Contract Addresses (fill after deployment)

| Contract | Address | Description |
|----------|---------|-------------|
| AWPToken | `0x...` | Main token (ERC20) |
| AWPRegistry | `0x...` | Unified entry point (account system + worknet management) |
| AWPWorkNet | `0x...` | Worknet identity NFT (ERC721) |
| WorknetTokenFactory | `0x...` | WorknetToken deployer (CREATE2) |
| veAWP | `0x...` | Position NFT (ERC721, deposit/withdraw AWP) |
| AWPAllocator | `0x...` | Allocation management (UUPS proxy) |
| AWPEmission | `0x...` | Emission engine (UUPS proxy) |
| AWPDAO | `0x...` | Governor |
| Treasury | `0x...` | Treasury (TimelockController) |
| WorknetManager (impl) | `0x...` | Default worknet contract implementation |

---

## 8. Error Code Reference

### Contract Custom Errors

| Error | Contract | Trigger Condition |
|-------|----------|-------------------|
| `InvalidWorknetParams()` | AWPRegistry | name/symbol length invalid |
| `WorknetManagerRequired()` | AWPRegistry | worknetManager is zero and no defaultWorknetManagerImpl set |
| `NotGuardian()` | AWPRegistry | Caller is not the Guardian |
| `NotOwner()` | AWPRegistry | Non-NFT holder calling lifecycle function |
| `InvalidWorknetStatus()` | AWPRegistry | Status does not meet precondition |
| `MaxActiveWorknetsReached()` | AWPRegistry | Active worknet count exceeds 10,000 |

| `PositionExpired()` | veAWP | Cannot add tokens to an expired lock position |
| `InsufficientUnallocated()` | AWPAllocator | Insufficient unallocated balance |
| `InsufficientAllocation()` | AWPAllocator | Allocation insufficient for the reduction |
| `MinterPaused()` | WorknetToken | WorknetToken minting paused (worknet is banned) |
| `ExceedsMaxSupply()` | WorknetToken | WorknetToken minting exceeds 10B cap |
| `NotMinter()` | WorknetToken | Caller is not an authorized minter |
| `ExpiredSignature()` | AWPRegistry | Gasless signature has expired |
| `InvalidSignature()` | AWPRegistry | Gasless signature verification failed |
| `EpochNotReady()` | AWPEmission | All epochs up to the current time-based epoch have been settled |
| `StrategyIsPaused()` | WorknetManager | Strategy execution is paused |
| `AlreadyClaimed()` | WorknetManager | User already claimed for this epoch |
| `InvalidProof()` | WorknetManager | Merkle proof verification failed |
| `RootAlreadySet()` | WorknetManager | Merkle root already submitted for this epoch |

### API HTTP Status Codes

| Status Code | Meaning |
|-------------|---------|
| 200 | Success |
| 400 | Bad request (missing/invalid parameters) |
| 404 | Resource not found |
| 429 | Rate limit exceeded |
| 500 | Internal server error |

---

## 9. Development Checklist

### Worknet Contract Development

- [ ] Deploy worknet contract (or use auto-deploy with `worknetManager = address(0)`)
- [ ] Implement AWP reward distribution logic (or use default WorknetManager strategies)
- [ ] Implement WorknetToken minting logic (proof of contribution via Merkle, or custom)
- [ ] Implement IERC1363Receiver.onTransferReceived for automatic AWP handling
- [ ] Call `awpRegistry.registerWorknet(params)` to register
- [ ] Call `awpRegistry.activateWorknet(worknetId)` to activate
- [ ] Configure AWPWorkNet metadata: setSkillsURI, setMinStake
- [ ] Request Guardian inclusion of your worknet contract address in `submitAllocations` weight submissions
- [ ] Verify your contract receives AWP after the first epoch settlement

### Coordinator Development

- [ ] Implement agent identity verification (via API or contract query)
- [ ] Implement WebSocket event listener (subscribe to relevant events)
- [ ] Implement task assignment logic
- [ ] Implement contribution evaluation and scoring
- [ ] Implement reward calculation and distribution
- [ ] Implement cold-start data synchronization (API + on-chain reads)
- [ ] Handle worknet lifecycle events (Paused, Banned, Deregistered)
