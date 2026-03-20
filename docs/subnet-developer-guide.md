# Subnet Developer Guide — AWP

> This document is intended for subnet developers (Coordinator / Subnet Contract), providing all the information needed to interact with AWPRegistry contracts and APIs.

---

## 1. Architecture Overview

```
User
  │
  ├── bind(target) → Tree-based binding (optional, every address is implicitly a root)
  ├── stakeNFT.deposit(AWP, lockDuration) → Deposit AWP, receive position NFT
  ├── allocate(staker, agent, subnetId, amount) → Allocate stake to a subnet (via AWPRegistry)
  └── registerSubnet(params) → Register new subnet (deploy Alpha + create LP)
        │
        ▼
AWPRegistry (Unified Entry Point — account system + allocation + subnet management)
  │
  ├── StakingVault — Pure allocation logic (no deposit/withdraw)
  ├── StakeNFT — ERC721 position NFT (deposit/withdraw AWP)
  ├── SubnetNFT — Subnet NFT (ownership)
  └── LPManager — PancakeSwap V4 CL pool, permanently locked LP
        │
        ▼
AWPEmission (UUPS Proxy — Independent Emission Engine)
  │
  ├── Oracle multi-sig → submitAllocations(recipients[], weights[], signatures[], effectiveEpoch)
  ├── settleEpoch(limit) → Batch-mint AWP to recipients and DAO
  └── emergencySetWeight(epoch, index, addr, weight) → DAO override via Timelock
        │
        ▼
Subnet Contract (What you develop)
  │
  ├── Receive AWP emission (minted to your contract address each epoch)
  ├── Mint Alpha tokens (you are the sole minter)
  └── Distribute rewards to miners (based on contribution)
```

**Key design:** AWPEmission is a generic address→weight distribution engine. It does not know about subnets — it simply mints AWP to addresses proportional to their oracle-assigned weights. Your subnet contract receives AWP because the oracle network includes your contract address in the allocation list.

---

## 2. Subnet Registration Flow

### 2.1 Prerequisites

1. **Option A**: Deploy your own subnet manager contract, OR use `address(0)` to auto-deploy the default `SubnetManager` proxy
2. Prepare AWP — LP creation cost = `INITIAL_ALPHA_MINT × initialAlphaPrice / 1e18` (default: 100M × 0.01 = 1M AWP)
3. Query current price: `awpRegistry.initialAlphaPrice()` → calculate the actual AWP amount needed
4. Call `AWPToken.approve(awpRegistry, awpAmount)` to authorize AWPRegistry for the transfer
5. No mandatory registration needed — every address can register subnets directly

### 2.2 Registration

```solidity
// Solidity — auto-deploy SubnetManager (subnetManager = address(0))
IAWPRegistry.SubnetParams memory params = IAWPRegistry.SubnetParams({
    name: "My Subnet Alpha",     // Alpha Token name (1-64 bytes)
    symbol: "MSALPHA",           // Alpha Token symbol (1-16 bytes)
    subnetManager: address(0),   // address(0) = auto-deploy SubnetManager proxy
    salt: bytes32(0),            // 0 = use subnetId as CREATE2 salt (default)
    minStake: 0                  // Minimum stake for agents (0 = no minimum)
});
uint256 subnetId = awpRegistry.registerSubnet(params);
```

```javascript
// Frontend (wagmi/viem)
const { writeContract } = useWriteContract();
await writeContract({
  address: AWP_REGISTRY_ADDRESS,
  abi: awpRegistryABI,
  functionName: 'registerSubnet',
  args: [{ name, symbol, subnetManager: '0x0000000000000000000000000000000000000000',
            salt: '0x00...00', minStake: 0n }]
});
```

> **Gasless registration**: Users can also register via `POST /api/relay/register-subnet` — the relayer pays gas, the user signs an EIP-712 message and pays AWP.

### 2.2b Vanity Alpha Token Address (optional)

If the factory was deployed with a non-zero `vanityRule`, you can pre-mine a CREATE2 salt off-chain to get a custom-looking Alpha token address.

**Step 1 — Understand the vanity rule**

Query `AlphaTokenFactory.vanityRule()` to see what pattern is enforced. A value of `0` means no validation is performed and you can skip this section.

**Step 2 — Get a salt via API (recommended)**

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

**Step 3 — Pass the salt in SubnetParams**

```solidity
IAWPRegistry.SubnetParams memory params = IAWPRegistry.SubnetParams({
    // ... other fields ...
    salt: 0xYourMinedSalt000000000000000000000000000000000000000000000000
});
```

The factory will deploy the Alpha token at the vanity address and then validate it matches the configured rule. If validation fails, the transaction reverts with `InvalidVanityAddress()`.

### 2.3 Automatic Steps After Registration

| Step | Action | Result |
|------|--------|--------|
| 1 | Calculate LP cost | `INITIAL_ALPHA_MINT × initialAlphaPrice / 1e18` |
| 2 | Transfer AWP from user to LPManager | `safeTransferFrom(user, lpManager, awpAmount)` |
| 3 | Deploy AlphaToken (CREATE2) | Standalone contract via factory; salt from params (or subnetId if 0); admin = AWPRegistry |
| 4 | Create PancakeSwap V4 CL pool | Full-range liquidity, 1% fee, permanently locked |
| 5 | Auto-deploy SubnetManager proxy (if subnetManager=0) | ERC1967Proxy → defaultSubnetManagerImpl, admin = user |
| 6 | Set subnet manager as sole Alpha minter | `setSubnetMinter(sc)` permanently locked |
| 7 | Mint SubnetNFT with identity data | name, subnetManager, alphaToken, minStake stored on-chain |
| 8 | Store lifecycle state | lpPool, status=Pending, createdAt |

### 2.4 Activate Subnet

After registration, the subnet status is `Pending`. The NFT Owner must manually activate it:

```solidity
awpRegistry.activateSubnet(subnetId);
// Status: Pending → Active
// Added to AWPRegistry's activeSubnetIds set
```

**Note:** Activation makes the subnet visible as "Active" in the protocol, but it does NOT automatically include it in AWP emission. To receive emissions, the oracle network must include your subnet contract address in their `submitAllocations` submission. Contact the oracle operators or DAO to request inclusion.

---

## 3. SubnetManager — Default Subnet Contract API

> When `subnetManager = address(0)` is passed to `registerSubnet`, AWPRegistry auto-deploys a `SubnetManager` proxy. The registrant (`msg.sender` or the EIP-712 signer for gasless) becomes `DEFAULT_ADMIN_ROLE`.

### 3.0 Roles

| Role | Bytes32 | Purpose |
|------|---------|---------|
| `DEFAULT_ADMIN_ROLE` | `0x00` | Grant/revoke all roles. Assigned to subnet registrant at deployment. |
| `MERKLE_ROLE` | `keccak256("MERKLE_ROLE")` | Submit Merkle roots for Alpha distribution |
| `STRATEGY_ROLE` | `keccak256("STRATEGY_ROLE")` | Choose AWP handling strategy + manually execute |
| `TRANSFER_ROLE` | `keccak256("TRANSFER_ROLE")` | Transfer any ERC20 held by the contract |

**Post-registration setup:**
```solidity
// Grant operator roles (admin calls these)
IAccessControl(subnetManager).grantRole(MERKLE_ROLE, merkleOperator);
IAccessControl(subnetManager).grantRole(STRATEGY_ROLE, strategyOperator);
IAccessControl(subnetManager).grantRole(TRANSFER_ROLE, treasuryMultisig);
```

### 3.1 Merkle Distribution (MERKLE_ROLE)

Distribute Alpha tokens to users via Merkle proofs. The contract mints directly — no pre-funding needed.

```
setMerkleRoot(uint32 epoch, bytes32 root)     // Submit root for an epoch (e.g. 20260318)
```

**Building the Merkle tree (off-chain):**
```javascript
// Leaf format: keccak256(keccak256(abi.encode(account, amount)))
import { StandardMerkleTree } from "@openzeppelin/merkle-tree";

const values = [
  ["0xAlice...", "1000000000000000000000"],  // 1000 Alpha
  ["0xBob...",   "500000000000000000000"],   // 500 Alpha
];
const tree = StandardMerkleTree.of(values, ["address", "uint256"]);
const root = tree.root;
// Submit: subnetManager.setMerkleRoot(20260318, root);

// Get proof for Alice:
const proof = tree.getProof(["0xAlice...", "1000000000000000000000"]);
```

**User claiming:**
```
claim(uint32 epoch, uint256 amount, bytes32[] proof)  // Anyone; mints Alpha to msg.sender
isClaimed(uint32 epoch, address account) → bool        // Check claim status
```

```javascript
// Frontend claim
await subnetManager.claim(20260318, parseEther("1000"), proof);
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
| `AddLiquidity` | 1 | Add AWP as single-sided liquidity (0 ~ current price range) |
| `BuybackBurn` | 2 | Swap AWP → Alpha via PancakeSwap V4, then burn Alpha |

```javascript
// Set strategy to BuybackBurn
await subnetManager.setStrategy(2);

// AWP emission now auto-triggers BuybackBurn via ERC1363 onTransferReceived
// No manual executeStrategy needed — AWPEmission.settleEpoch uses mintAndCall
```

**Auto-trigger flow:**
```
AWPEmission.settleEpoch(limit)
  → awpToken.mintAndCall(subnetManager, amount, "")
    → SubnetManager.onTransferReceived()
      → if strategy == AddLiquidity → PancakeSwap V4 single-sided LP
      → if strategy == BuybackBurn  → swap AWP→Alpha, burn Alpha
      → if strategy == Reserve      → no-op (AWP stays in contract)
```

### 3.3 Token Transfer (TRANSFER_ROLE)

```
transferToken(address token, address to, uint256 amount)  // Transfer any ERC20 out
```

```javascript
// Send reserved AWP to a multisig
await subnetManager.transferToken(AWP_TOKEN, multisig, parseEther("10000"));
```

### 3.4 View Functions

```
alphaToken() → address                        // Alpha token for this subnet
awpToken() → address                          // AWP token
poolId() → bytes32                            // PancakeSwap V4 pool ID
currentStrategy() → AWPStrategy               // Current AWP handling strategy (0/1/2)
merkleRoots(uint32 epoch) → bytes32           // Merkle root for an epoch
claimed(uint32 epoch, address account) → bool  // Alias for isClaimed
```

### 3.5 Complete Operator Workflow

```
1. Register subnet → auto-deploy SubnetManager → you are admin
2. Grant roles:
   - MERKLE_ROLE to your backend that computes reward distributions
   - STRATEGY_ROLE to your governance multisig
   - TRANSFER_ROLE to your treasury multisig
3. Set AWP strategy (default: Reserve):
   - subnetManager.setStrategy(2)  // BuybackBurn
4. Each epoch, AWP emission arrives automatically via mintAndCall
5. Periodically compute Merkle tree of Alpha rewards:
   - subnetManager.setMerkleRoot(epoch, root)
6. Users claim Alpha tokens:
   - subnetManager.claim(epoch, amount, proof)
7. Monitor via WebSocket:
   - Subscribe to: SkillsURIUpdated, MinStakeUpdated, Allocated, Deallocated
8. Update subnet settings:
   - subnetNFT.setSkillsURI(subnetId, "ipfs://...")
   - subnetNFT.setMinStake(subnetId, 100e18)
```

---

## 3. Subnet Contract Development

### 3.1 Default SubnetManager (auto-deployed)

If you pass `subnetManager = address(0)` during registration, AWPRegistry auto-deploys a **SubnetManager** proxy with built-in features:

- **Merkle Distribution**: Submit Merkle roots per epoch, users claim Alpha tokens via proof
- **AWP Strategy**: Choose how incoming AWP is handled:
  - `Reserve` — keep in contract for manual distribution
  - `AddLiquidity` — add single-sided liquidity to PancakeSwap V4
  - `BuybackBurn` — buy Alpha from pool and burn it
- **ERC1363 Receiver**: AWP received via `mintAndCall` (from AWPEmission) auto-triggers the current strategy
- **Token Transfer**: Move any ERC20 held by the contract
- **AccessControl**: Three roles — `MERKLE_ROLE`, `STRATEGY_ROLE`, `TRANSFER_ROLE` (you are the admin)

After registration, grant roles to your operator addresses:
```solidity
IAccessControl(subnetManager).grantRole(MERKLE_ROLE, operatorAddress);
IAccessControl(subnetManager).grantRole(STRATEGY_ROLE, operatorAddress);
```

### 3.2 Custom Subnet Contract (advanced)

If you need custom logic, deploy your own contract and pass its address as `subnetManager`:

```solidity
contract MySubnetContract {
    IAlphaToken public alphaToken;
    IERC20 public awpToken;

    // AWP emission arrives via mintAndCall — implement IERC1363Receiver to auto-process
    function onTransferReceived(address, address, uint256 amount, bytes calldata)
        external returns (bytes4) {
        // Your custom AWP handling logic
        return IERC1363Receiver.onTransferReceived.selector;
    }

    // Mint Alpha to miners (you are the sole minter)
    function mintAlphaToMiner(address miner, uint256 amount) external {
        alphaToken.mint(miner, amount);
    }
}
```

### 3.2 Key Permissions

| Permission | Who Holds It | Description |
|------------|-------------|-------------|
| Alpha minting | **Your subnet contract** | `setSubnetMinter` permanently locked — only you can mint Alpha |
| Alpha minting pause | AWPRegistry (admin) | `setMinterPaused(true)` pauses your minting when banned |
| AWP emission receipt | Oracle-assigned | Your address must be in the oracle's `submitAllocations` list to receive AWP |
| Alpha MAX_SUPPLY | 10B | Independent cap per subnet |

### 3.3 Ban and Unban

- **Ban**: After a DAO proposal passes, Treasury (Timelock) calls `banSubnet(subnetId)` → your Alpha minting is paused, subnet removed from AWPRegistry's active list
- **Unban**: Treasury calls `unbanSubnet(subnetId)` → status restores to **Active**, re-added to active list
- During a ban, your existing AWP and Alpha holdings are **unaffected** — only new Alpha minting is blocked
- Note: The direct caller is the Treasury contract (Timelock), not the DAO contract
- Oracle operators typically exclude banned subnets from their allocation submissions

---

## 4. API Reference

> Base URL: `https://tapi.awp.sh` (or your self-hosted address)

### 4.1 System

#### `GET /api/health`
```json
{"status": "ok"}
```

#### `GET /api/registry`
Returns all 9 protocol contract addresses (excludes implementation contracts):
```json
{
  "awpRegistry": "0x...",
  "awpToken": "0x...",
  "awpEmission": "0x...",
  "stakingVault": "0x...",
  "stakeNFT": "0x...",
  "subnetNFT": "0x...",
  "lpManager": "0x...",
  "alphaTokenFactory": "0x...",
  "dao": "0x...",
  "treasury": "0x..."
}
```

### 4.2 Users

#### `GET /api/users/{address}`
```json
{
  "user": { "address": "0x...", "registered_at": 1710000000 },
  "balance": { "user_address": "0x...", "total_staked": "5000", "total_allocated": "3000" },
  "rewardRecipient": { "user_address": "0x...", "recipient_address": "0x..." },
  "agents": [{ "agent_address": "0x...", "owner_address": "0x...", "is_manager": false, "removed": false }]
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

### 4.3 Staking

#### `GET /api/staking/user/{address}/balance`
```json
{
  "totalStaked": "10000000000000000000000",
  "totalAllocated": "5000000000000000000000",
  "unallocated": "5000000000000000000000"
}
```
> `totalStaked` is computed from StakeNFT positions. No `withdrawRequest` field (lock period replaces cooldown).

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
  {"user_address": "0x...", "agent_address": "0x...", "subnet_id": 1, "amount": "5000", "frozen": false}
]
```

#### `GET /api/staking/agent/{agent}/subnet/{subnetId}`
**Commonly used by Coordinators** — query total stake of an agent on a specific subnet:
```json
{"amount": "5000000000000000000000"}
```

#### `GET /api/staking/agent/{agent}/subnets`
```json
[{"subnet_id": 1, "amount": "5000"}, {"subnet_id": 3, "amount": "2000"}]
```

#### `GET /api/staking/subnet/{subnetId}/total`
```json
{"total": "50000000000000000000000"}
```

### 4.5 Subnets

#### `GET /api/subnets/?status=Active&page=1&limit=20`
```json
[
  {
    "subnet_id": 1,
    "owner": "0x...",
    "name": "My Subnet",
    "symbol": "MSUB",
    "subnet_contract": "0x...",
    "alpha_token": "0x...",
    "lp_pool": "0x...",
    "status": "Active",
    "created_at": 1710000000,
    "activated_at": 1710000100
  }
]
```

#### `GET /api/subnets/{subnetId}`
Single subnet detail.

#### `GET /api/subnets/{subnetId}/earnings?page=1&limit=20`
Subnet AWP emission history (queried by subnet contract address from `recipient_awp_distributions`):
```json
[
  {"epoch_id": 5, "recipient": "0x1234...", "awp_amount": "7900000000000000000000000"}
]
```

#### `GET /api/subnets/{subnetId}/agents/{agent}`
**Used by Coordinators** — query agent stake info on a subnet:
```json
{"agent": "0x...", "subnetId": 1, "stake": "5000000000000000000000"}
```

### 4.6 Emission

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
  {"epoch_id": 42, "start_time": 1710000000, "daily_emission": "15000000...", "dao_emission": "7500000..."}
]
```

### 4.7 Tokens

#### `GET /api/tokens/awp`
```json
{"totalSupply": "5015800000000000000000000000", "maxSupply": "10000000000000000000000000000"}
```

#### `GET /api/tokens/alpha/{subnetId}`
```json
{"subnetId": 1, "name": "My Subnet Alpha", "symbol": "MSALPHA", "alphaToken": "0x..."}
```

#### `GET /api/tokens/alpha/{subnetId}/price`
```json
{"priceInAWP": "0.015", "reserve0": "...", "reserve1": "...", "updatedAt": "..."}
```

### 4.8 Governance

#### `GET /api/governance/proposals?status=Active&page=1&limit=20`
#### `GET /api/governance/proposals/{proposalId}`
#### `GET /api/governance/treasury`
```json
{"treasuryAddress": "0x..."}
```

### 4.9 Relay (Gasless Transactions)

> Rate limit: 100 requests per IP per 1 hour. Requires `RELAYER_PRIVATE_KEY` configured on the server.

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

#### `POST /api/relay/register-subnet`
Fully gasless subnet registration — user signs two messages (ERC-2612 permit + EIP-712), relayer pays all gas:
```json
// Request
{"user": "0x...", "name": "EVO Alpha", "symbol": "EVO",
 "subnetManager": "0x0000...0000",
 "salt": "0x...", "minStake": "0", "deadline": 1742400000,
 "permitSignature": "0x...(ERC-2612 AWP permit)",
 "registerSignature": "0x...(EIP-712 registerSubnet)"}
// Response
{"txHash": "0x..."}
```

### 4.10 WebSocket

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
| `Bound` | `{user, target, oldTarget}` | ✅ |
| `RecipientUpdated` | `{user, recipient}` | |
| `DelegateGranted` | `{user, delegate}` | ✅ |
| `DelegateRevoked` | `{user, delegate}` | ✅ |
| `Deposited` (StakeNFT) | `{user, tokenId, amount, lockEndTime}` | |
| `PositionIncreased` (StakeNFT) | `{tokenId, addedAmount, newLockEndTime}` | |
| `Withdrawn` (StakeNFT) | `{user, tokenId, amount}` | |
| `Allocated` | `{user, agent, subnetId, amount}` | ✅ |
| `Deallocated` | `{user, agent, subnetId, amount}` | ✅ |
| `Reallocated` | `{user, fromAgent, fromSubnet, toAgent, toSubnet, amount}` | ✅ |
| `SubnetRegistered` | `{subnetId, owner, name, symbol, subnetManager, alphaToken}` | |
| `LPCreated` | `{subnetId, poolId, awpAmount, alphaAmount}` | |
| `SkillsURIUpdated` | `{subnetId, skillsURI}` | ✅ |
| `MinStakeUpdated` | `{subnetId, minStake}` | ✅ |
| `SubnetActivated` | `{subnetId}` | ✅ |
| `SubnetPaused` | `{subnetId}` | ✅ |
| `SubnetResumed` | `{subnetId}` | ✅ |
| `SubnetBanned` | `{subnetId}` | ✅ |
| `SubnetUnbanned` | `{subnetId}` | ✅ |
| `SubnetDeregistered` | `{subnetId}` | ✅ |
| `GovernanceWeightUpdated` | `{addr, weight}` | ✅ |
| `RecipientAWPDistributed` | `{epoch, recipient, awpAmount}` | ✅ |
| `DAOMatchDistributed` | `{epoch, amount}` | |
| `EpochSettled` | `{epoch, totalEmission, recipientCount}` | ✅ |
| `AllocationsSubmitted` | `{nonce, recipients, weights}` | ✅ |
| `OracleConfigUpdated` | `{oracles, threshold}` | |

---

## 5. Contract Interaction Reference

### 5.1 On-Chain Reads (view functions, no gas)

```solidity
// Query full subnet info (combines AWPRegistry state + SubnetNFT identity)
IAWPRegistry.SubnetFullInfo memory subnet = awpRegistry.getSubnetFull(subnetId);
// subnet.subnetManager, subnet.alphaToken, subnet.status, subnet.owner, subnet.minStake

// Query single agent info
AWPRegistry.AgentInfo memory info = awpRegistry.getAgentInfo(agentAddress, subnetId);

// Batch query (Coordinator cold start)
AWPRegistry.AgentInfo[] memory infos = awpRegistry.getAgentsInfo(agents, subnetId);

// Query current epoch
uint256 epoch = awpEmission.currentEpoch();

// Check if subnet is active
bool active = awpRegistry.isSubnetActive(subnetId);

// Query emission weight
uint96 weight = awpEmission.getWeight(subnetManagerAddress);
uint256 totalWeight = awpEmission.getTotalWeight();

// SubnetNFT: update skills URI or minStake (NFT owner only)
subnetNFT.setSkillsURI(subnetId, "ipfs://QmNewSkills...");
subnetNFT.setMinStake(subnetId, 100e18); // 100 AWP minimum
```

### 5.2 Recommended Coordinator Data Sync Pattern

```
Cold Start:
  1. GET /api/subnets/{subnetId} → Subnet basic info
  2. POST /api/agents/batch-info → Batch query all known agents' stakes
  3. Or on-chain: awpRegistry.getAgentsInfo(allAgents, subnetId)

Runtime Incremental Updates:
  1. WebSocket subscribe: ["Bound", "Allocated", "Deallocated", "Reallocated", "DelegateGranted", "DelegateRevoked"]
  2. Update local agent cache on each event
  3. Periodically GET /api/staking/subnet/{subnetId}/total to verify total stake

Reward Distribution:
  1. Calculate reward ratio based on agent stake and contribution score
  2. Query rewardRecipient: POST /api/agents/batch-info or on-chain awpRegistry.getAgentInfo(agent, subnetId)
  3. Transfer AWP/Alpha to the rewardRecipient address
```

### 5.3 Emission Timeline [DRAFT]

> **The emission mechanism has not been finalized. The timeline below is preliminary.**

```
Each Epoch (1 day / daily):
  1. Keeper (or anyone) calls AWPEmission.settleEpoch(limit) — initializes epoch on first call
  2. Processes up to `limit` recipients per call, using mintAndCall to send AWP
  3. If your contract implements IERC1363Receiver, onTransferReceived is called automatically
  4. SubnetManager auto-executes the configured AWP strategy (Reserve/AddLiquidity/BuybackBurn)
  5. Final call mints DAO share and advances epoch counter

Your subnet receives: epochEmission × 50% × (yourWeight / totalWeight)

Decay: each epoch emission × 0.996844 (~99% released in 4 years)
Initial daily emission: 15,800,000 AWP per epoch (1 epoch = 1 day)
```

**How your subnet gets weight:**
- Oracle network assigns weights to recipient addresses via `submitAllocations(recipients[], weights[], signatures[], effectiveEpoch)`
- Multiple oracles must sign the same allocation (EIP-712 multi-sig threshold)
- DAO can override a single recipient's entry via `emergencySetWeight(epoch, index, addr, weight)` through Timelock governance
- Your subnet contract address must be included in the oracle's recipient list to receive AWP

---

## 6. Contract Addresses (fill after deployment)

| Contract | Address | Description |
|----------|---------|-------------|
| AWPToken | `0x0000d0e38e9c6ba147b0098bb42007b942ef00a1` | Main token (ERC20Votes) |
| AWPRegistry | `0x00003a7fa04c3af3adba2dc3c6622277501400b1` | Unified entry point (allocation + subnet management) |
| SubnetNFT | `0x0f86ec2f2fbf234b00b18e66e7c5e00518091cda` | Subnet NFT (ERC721) |
| AlphaTokenFactory | `0x3ebe3168c898f4b05ebf0c0d17f4739e111e5164` | Alpha token deployer (CREATE2) |
| StakeNFT | `0x4f7e8d4487c0c514b72ed0e35ed707cb8acdce39` | Position NFT (ERC721, deposit/withdraw AWP) |
| AWPDAO | `0x7171211da849a2c569643fb1e8f5399ddd71939a` | Governor |
| Treasury | `0x9ee82684e4214edb405d930001e9058d1913d994` | Treasury (TimelockController) |
| SubnetManager (impl) | `0xE5771dC2a5a577CDFa6b939Af4F32Ad13CFc6D92` | Default subnet contract implementation |

---

## 7. Error Code Reference

### Contract Custom Errors

| Error | Contract | Trigger Condition |
|-------|----------|-------------------|
| `InvalidSubnetParams()` | AWPRegistry | name/symbol length invalid |
| `SubnetManagerRequired()` | AWPRegistry | subnetManager is zero and no defaultSubnetManagerImpl set |
| `NotTimelock()` | AWPRegistry | Caller is not the Treasury (Timelock) |
| `NotGuardian()` | AWPRegistry | Caller is not the Guardian |
| `PositionExpired()` | StakeNFT | Cannot add tokens to an expired lock position |
| `NotOwner()` | AWPRegistry | Non-NFT holder calling lifecycle function |
| `InvalidSubnetStatus()` | AWPRegistry | Status does not meet precondition |
| `MaxActiveSubnetsReached()` | AWPRegistry | Active subnet count exceeds 10000 |
| `ImmunityNotExpired()` | AWPRegistry | Attempt to deregister during immunity period |
| `InvalidAgent()` | AWPRegistry | Agent does not belong to the user |
| `NotRegistered()` | AWPRegistry | Caller is not a registered user |
| `ExpiredSignature()` | AWPRegistry | Gasless signature has expired |
| `InvalidSignature()` | AWPRegistry | Gasless signature verification failed |
| `InsufficientUnallocated()` | StakingVault | Insufficient unallocated balance |
| `InsufficientAllocation()` | StakingVault | Allocation insufficient for the reduction |
| `MinterPaused()` | AlphaToken | Alpha minting paused (subnet is banned) |
| `ExceedsMaxSupply()` | AlphaToken | Alpha minting exceeds 10B cap |
| `NotMinter()` | AlphaToken | Caller is not an authorized minter |
| `SettlementInProgress()` | AWPEmission | Cannot modify allocations during epoch settlement |
| `OracleNotConfigured()` | AWPEmission | Oracle is not yet configured |
| `InvalidSignatureCount()` | AWPEmission | Not enough valid oracle signatures |
| `EpochNotReady()` | AWPEmission | All epochs up to the current time-based epoch have been settled |

### API HTTP Status Codes

| Status Code | Meaning |
|-------------|---------|
| 200 | Success |
| 400 | Bad request (missing/invalid parameters) |
| 404 | Resource not found |
| 500 | Internal server error |

---

## 8. Development Checklist

### Subnet Contract Development

- [ ] Deploy subnet contract
- [ ] Implement AWP reward distribution logic
- [ ] Implement Alpha minting logic (proof of contribution)
- [ ] Call `awpRegistry.registerSubnet(params)` to register
- [ ] Call `awpRegistry.activateSubnet(subnetId)` to activate
- [ ] Contact oracle operators to include your subnet contract address in `submitAllocations` weight submissions
- [ ] Verify your contract receives AWP after the first epoch settlement

### Coordinator Development

- [ ] Implement agent identity verification (via API or contract query)
- [ ] Implement WebSocket event listener (subscribe to relevant events)
- [ ] Implement task assignment logic
- [ ] Implement contribution evaluation and scoring
- [ ] Implement reward calculation and distribution
- [ ] Implement cold-start data synchronization (batch-info API)
