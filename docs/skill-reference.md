# AWP Protocol — Skill Development Reference

> This document provides everything needed to build a Claude Code skill that interacts with the AWP (Agent Working Protocol) smart contracts and API. It is user-facing only — no admin/guardian functions are exposed.

---

## 1. Protocol Overview

AWP is a multi-chain agent mining protocol deployed on **Base (8453)**, **Ethereum (1)**, **Arbitrum (42161)**, and **BSC (56)**. All contract addresses are identical across all 4 chains.

**Core concepts:**
- **AWP Token**: ERC20 governance/utility token (10B max supply)
- **Worknet**: An autonomous agent network. Each worknet has its own Alpha token and LP pool
- **StakeNFT**: Users deposit AWP with a lock period to mint position NFTs
- **Allocation**: Users allocate their staked AWP to agents in worknets via StakingVault
- **Emission**: AWP minted daily and distributed to worknets by weight
- **Binding**: Tree-based account linking — agents bind to owners, rewards flow to root

---

## 2. Contract Addresses (Same on all 4 chains)

```
AWPToken:           0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1
AWPRegistry:        0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A
AWPEmission:        0x3C9cB73f8B81083882c5308Cce4F31f93600EaA9
StakingVault:       0xE8A204fD9c94C7E28bE11Af02fc4A4AC294Df29b
StakeNFT:           0x4E119560632698Bab67cFAB5d8EC0A373363ba2d
WorknetNFT:         0xB9F03539BE496d09c4d7964921d674B8763f5233
AlphaTokenFactory:  0xB2e4897eD77d0f5BFa3140B9989594de09a8037c
Treasury:           0x82562023a053025F3201785160CaE6051efD759e
AWPDAO:             0x6a074aC9823c47f86EE4Fc7F62e4217Bc9C76004
```

---

## 3. API — JSON-RPC 2.0

**Endpoint**: `POST https://api.awp.sh/v2`
**Discovery**: `GET https://api.awp.sh/v2`
**WebSocket**: `wss://api.awp.sh/ws/live`
**Batch**: Up to 20 requests per batch, executed concurrently.

### Request/Response Format

```json
// Request
{"jsonrpc": "2.0", "method": "namespace.method", "params": {...}, "id": 1}

// Success
{"jsonrpc": "2.0", "result": {...}, "id": 1}

// Error
{"jsonrpc": "2.0", "error": {"code": -32601, "message": "method not found"}, "id": 1}
```

### Error Codes

| Code | Meaning |
|------|---------|
| -32700 | Parse error |
| -32600 | Invalid request |
| -32601 | Method not found |
| -32602 | Invalid params |
| -32603 | Internal error |
| -32001 | Resource not found |

### Common Parameters

- `address`: `string` — 0x-prefixed, 40 hex chars, case-insensitive
- `worknetId`: `string` — Globally unique: `(chainId << 64) | localId`, passed as decimal string
- `chainId`: `integer` — Optional on most methods; omit or 0 for default chain
- `page`: `integer` — 1-indexed (default 1)
- `limit`: `integer` — Items per page (default 20, max 100)

---

### 3.1 System

| Method | Params | Description |
|--------|--------|-------------|
| `stats.global` | — | Global protocol stats (users, worknets, staked, emitted, chains) |
| `registry.get` | `chainId?` | All contract addresses + EIP-712 domain info |
| `health.check` | — | `{"status": "ok"}` |
| `health.detailed` | — | Per-chain indexer/keeper status |
| `chains.list` | — | Supported chains array |

### 3.2 Users

| Method | Params | Description |
|--------|--------|-------------|
| `users.list` | `chainId?`, `page?`, `limit?` | Paginated user list |
| `users.listGlobal` | `page?`, `limit?` | Cross-chain deduplicated |
| `users.count` | `chainId?` | Total user count |
| `users.get` | `address` (required), `chainId?` | User details (balance, agents, recipient) |
| `users.getPortfolio` | `address` (required), `chainId?` | Full portfolio: identity + balance + positions + allocations + delegates |
| `users.getDelegates` | `address` (required), `chainId?` | Agents bound to user |

### 3.3 Address & Nonce

| Method | Params | Description |
|--------|--------|-------------|
| `address.check` | `address` (required), `chainId?` | Registration, binding, recipient status |
| `address.resolveRecipient` | `address` (required), `chainId?` | Walk bind chain to root, return recipient |
| `address.batchResolveRecipients` | `addresses[]` (required, max 500), `chainId?` | Batch resolve |
| `nonce.get` | `address` (required), `chainId?` | AWPRegistry EIP-712 nonce |
| `nonce.getStaking` | `address` (required), `chainId?` | StakingVault EIP-712 nonce |

### 3.4 Agents

| Method | Params | Description |
|--------|--------|-------------|
| `agents.getByOwner` | `owner` (required), `chainId?` | All agents bound to owner |
| `agents.getDetail` | `agent` (required), `chainId?` | Agent details |
| `agents.lookup` | `agent` (required), `chainId?` | Returns `{"ownerAddress": "0x..."}` |
| `agents.batchInfo` | `agents[]` (required, max 100), `worknetId` (required), `chainId?` | Batch agent info + stake |

### 3.5 Staking

| Method | Params | Description |
|--------|--------|-------------|
| `staking.getBalance` | `address` (required), `chainId?` | Staked / allocated / available |
| `staking.getUserBalanceGlobal` | `address` (required) | Aggregated across all chains |
| `staking.getPositions` | `address` (required), `chainId?` | StakeNFT positions |
| `staking.getPositionsGlobal` | `address` (required) | Positions across all chains |
| `staking.getAllocations` | `address` (required), `chainId?`, `page?`, `limit?` | Allocation records |
| `staking.getAgentSubnetStake` | `agent` (required), `worknetId` (required) | Agent's stake in worknet (cross-chain) |
| `staking.getAgentSubnets` | `agent` (required) | All worknets agent participates in |
| `staking.getSubnetTotalStake` | `worknetId` (required) | Total stake in worknet |

### 3.6 Worknets

| Method | Params | Description |
|--------|--------|-------------|
| `subnets.list` | `status?`, `chainId?`, `page?`, `limit?` | Filter: `Pending`/`Active`/`Paused`/`Banned` |
| `subnets.listRanked` | `chainId?`, `page?`, `limit?` | Ranked by total stake |
| `subnets.search` | `query` (required, 1-100 chars), `chainId?`, `page?`, `limit?` | Name/symbol search |
| `subnets.getByOwner` | `owner` (required), `chainId?`, `page?`, `limit?` | Worknets owned by address |
| `subnets.get` | `worknetId` (required) | Worknet details |
| `subnets.getSkills` | `worknetId` (required) | Skills URI |
| `subnets.getEarnings` | `worknetId` (required), `page?`, `limit?` | AWP earnings history |
| `subnets.getAgentInfo` | `worknetId` (required), `agent` (required) | Agent info in worknet |
| `subnets.listAgents` | `worknetId` (required), `chainId?`, `page?`, `limit?` | Agents ranked by stake |

### 3.7 Emission

| Method | Params | Description |
|--------|--------|-------------|
| `emission.getCurrent` | `chainId?` | Current epoch, daily emission, total weight |
| `emission.getSchedule` | `chainId?` | 30/90/365 day projections with decay |
| `emission.getGlobalSchedule` | — | Aggregated across all chains |
| `emission.listEpochs` | `chainId?`, `page?`, `limit?` | Settled epochs |
| `emission.getEpochDetail` | `epochId` (required), `chainId?` | Per-recipient distributions |

### 3.8 Tokens

| Method | Params | Description |
|--------|--------|-------------|
| `tokens.getAWP` | `chainId?` | AWP supply, max supply |
| `tokens.getAWPGlobal` | — | Aggregated across chains |
| `tokens.getAlphaInfo` | `worknetId` (required) | Alpha token info |
| `tokens.getAlphaPrice` | `worknetId` (required) | Price from LP pool (cached 10min) |

### 3.9 Governance

| Method | Params | Description |
|--------|--------|-------------|
| `governance.listProposals` | `status?`, `chainId?`, `page?`, `limit?` | Filter: `Active`/`Canceled`/`Defeated`/`Succeeded`/`Queued`/`Expired`/`Executed` |
| `governance.listAllProposals` | `status?`, `page?`, `limit?` | Cross-chain |
| `governance.getProposal` | `proposalId` (required), `chainId?` | Proposal details |
| `governance.getTreasury` | — | Treasury address |

---

## 4. User-Facing Smart Contract Functions

### 4.1 AWPRegistry — Account System

```solidity
// Binding
function bind(address target) external;
function unbind() external;

// Recipient
function setRecipient(address addr) external;

// Delegation
function grantDelegate(address delegate) external;
function revokeDelegate(address delegate) external;

// View
function resolveRecipient(address addr) external view returns (address);
function batchResolveRecipients(address[] calldata addrs) external view returns (address[] memory);
function isRegistered(address addr) external view returns (bool);
function boundTo(address) external view returns (address);
function recipient(address) external view returns (address);
function delegates(address user, address delegate) external view returns (bool);
function nonces(address) external view returns (uint256);
```

### 4.2 AWPRegistry — Worknet Management

```solidity
struct WorknetParams {
    string name;              // Alpha Token name (1-64 chars, no " or \)
    string symbol;            // Alpha Token symbol (1-16 chars, no " or \)
    address worknetManager;   // 0x0 = auto-deploy WorknetManager proxy
    bytes32 salt;             // CREATE2 salt (0x0 = use worknetId)
    uint128 minStake;         // Min stake for agents (reference only, not enforced)
    string skillsURI;         // Skills description URI
}

// Registration (costs initialAlphaMint * initialAlphaPrice AWP, currently 100,000 AWP)
// User must approve AWPRegistry for AWP spending before calling
function registerWorknet(WorknetParams calldata params) external returns (uint256 worknetId);

// Lifecycle (caller must be WorknetNFT owner)
function activateWorknet(uint256 worknetId) external;   // Pending → Active
function pauseWorknet(uint256 worknetId) external;      // Active → Paused
function resumeWorknet(uint256 worknetId) external;     // Paused → Active

// View
function getWorknet(uint256 worknetId) external view returns (WorknetInfo memory);
function getWorknetFull(uint256 worknetId) external view returns (WorknetFullInfo memory);
function getActiveWorknetCount() external view returns (uint256);
function isWorknetActive(uint256 worknetId) external view returns (bool);
function initialAlphaPrice() external view returns (uint256);
function initialAlphaMint() external view returns (uint256);
```

### 4.3 StakeNFT — AWP Staking

```solidity
// Deposit AWP, mint position NFT (user must approve StakeNFT for AWP spending)
function deposit(uint256 amount, uint64 lockDuration) external returns (uint256 tokenId);

// Deposit with ERC-2612 permit (no prior approve needed)
function depositWithPermit(uint256 amount, uint64 lockDuration, uint256 deadline, uint8 v, bytes32 r, bytes32 s) external returns (uint256 tokenId);

// Add AWP to existing position (blocked if lock expired)
function addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime) external;

// Withdraw after lock expires (burns NFT)
function withdraw(uint256 tokenId) external;

// View
function getUserTotalStaked(address user) external view returns (uint256);
function getVotingPower(uint256 tokenId) external view returns (uint256);
function remainingTime(uint256 tokenId) external view returns (uint64);
```

### 4.4 StakingVault — Allocation

```solidity
// Caller must be staker or staker's delegate
function allocate(address staker, address agent, uint256 worknetId, uint256 amount) external;
function deallocate(address staker, address agent, uint256 worknetId, uint256 amount) external;
function reallocate(address staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount) external;

// View
function userTotalAllocated(address user) external view returns (uint256);
function getAgentStake(address user, address agent, uint256 worknetId) external view returns (uint256);
function getAgentWorknets(address user, address agent) external view returns (uint256[] memory);
function nonces(address) external view returns (uint256);
```

---

## 5. Gasless Relay Endpoints

Users sign EIP-712 messages off-chain; relayer submits the transaction and pays gas.

**RPC not available for relay** — use REST relay endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `POST /api/relay/bind` | RelayBind | Gasless bind |
| `POST /api/relay/unbind` | RelayUnbind | Gasless unbind |
| `POST /api/relay/set-recipient` | RelaySetRecipient | Gasless setRecipient |
| `POST /api/relay/register` | RelayRegister | Gasless register (= setRecipientFor to self) |
| `POST /api/relay/allocate` | RelayAllocate | Gasless allocate |
| `POST /api/relay/deallocate` | RelayDeallocate | Gasless deallocate |
| `POST /api/relay/activate-subnet` | RelayActivateWorknet | Gasless activate worknet |
| `POST /api/relay/register-subnet` | RelayRegisterWorknet | Gasless register worknet (with permit) |
| `POST /api/relay/grant-delegate` | RelayGrantDelegate | Gasless grant delegate |
| `POST /api/relay/revoke-delegate` | RelayRevokeDelegate | Gasless revoke delegate |
| `GET /api/relay/status/{txHash}` | GetRelayStatus | Check relay tx status |

### Relay Request Format (example: bind)

```json
POST /api/relay/bind
{
  "agent": "0x...",      // signer address
  "target": "0x...",     // bind target
  "deadline": 1712345678,
  "v": 27,
  "r": "0x...",
  "s": "0x..."
}
```

### EIP-712 Signing

**AWPRegistry domain** (for bind, unbind, setRecipient, grantDelegate, revokeDelegate, registerWorknet, activateWorknet):
```json
{
  "name": "AWPRegistry",
  "version": "1",
  "chainId": <chainId>,
  "verifyingContract": "0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A"
}
```

**StakingVault domain** (for allocate, deallocate):
```json
{
  "name": "StakingVault",
  "version": "1",
  "chainId": <chainId>,
  "verifyingContract": "0xE8A204fD9c94C7E28bE11Af02fc4A4AC294Df29b"
}
```

### EIP-712 Type Definitions

```
Bind(address agent, address target, uint256 nonce, uint256 deadline)
Unbind(address user, uint256 nonce, uint256 deadline)
SetRecipient(address user, address recipient, uint256 nonce, uint256 deadline)
GrantDelegate(address user, address delegate, uint256 nonce, uint256 deadline)
RevokeDelegate(address user, address delegate, uint256 nonce, uint256 deadline)
ActivateWorknet(address user, uint256 worknetId, uint256 nonce, uint256 deadline)
RegisterWorknet(address user, WorknetParams params, uint256 nonce, uint256 deadline)
  with WorknetParams(string name, string symbol, address worknetManager, bytes32 salt, uint128 minStake, string skillsURI)
Allocate(address staker, address agent, uint256 worknetId, uint256 amount, uint256 nonce, uint256 deadline)
Deallocate(address staker, address agent, uint256 worknetId, uint256 amount, uint256 nonce, uint256 deadline)
```

Nonces: fetch via `nonce.get` (AWPRegistry) or `nonce.getStaking` (StakingVault) before signing.

---

## 6. WebSocket Real-Time Events

**Endpoint**: `wss://api.awp.sh/ws/live`

Events pushed as JSON messages:

| Event | Fields | Description |
|-------|--------|-------------|
| `UserRegistered` | `user` | New user registration |
| `Bound` | `user`, `target` | User bound to target |
| `Unbound` | `user` | User unbound |
| `RecipientSet` | `user`, `recipient` | Recipient changed |
| `Deposited` | `user`, `tokenId`, `amount`, `lockEndTime` | AWP staked |
| `Withdrawn` | `user`, `tokenId`, `amount` | AWP withdrawn |
| `Allocated` | `staker`, `agent`, `worknetId`, `amount` | Stake allocated |
| `Deallocated` | `staker`, `agent`, `worknetId`, `amount` | Stake deallocated |
| `Reallocated` | `staker`, `fromAgent`, `fromWorknetId`, `toAgent`, `toWorknetId`, `amount` | Reallocation |
| `WorknetRegistered` | `worknetId`, `owner`, `name`, `symbol` | New worknet |
| `WorknetActivated` | `worknetId` | Worknet activated |
| `EpochSettled` | `epoch`, `totalEmission`, `recipientCount` | Epoch settled |
| `RecipientAWPDistributed` | `epoch`, `recipient`, `amount` | AWP distributed |
| `AllocationsSubmitted` | `epoch`, `totalWeight`, `recipients`, `weights` | Allocation weights submitted |

---

## 7. Key Protocol Parameters

| Parameter | Value | Description |
|-----------|-------|-------------|
| AWP Max Supply | 10,000,000,000 (10B) | Hard cap |
| Initial Daily Emission | 31,600,000 AWP | Per chain, decays daily |
| Decay Factor | 996844 / 1000000 | ~0.3156% decay per epoch |
| Epoch Duration | 1 day (86400s) | |
| Worknet Registration Cost | 100,000 AWP | `100M × 0.001` (initialAlphaMint × initialAlphaPrice) |
| Alpha Mint per Worknet | 100,000,000 | Minted to LP pool on registration |
| Min Lock Duration (StakeNFT) | 1 day | |
| Max Voting Weight Duration | 54 weeks | |
| Immunity Period | 30 days | Before worknet can be deregistered |
| Timelock Delay (Treasury) | 2 days | For DAO governance proposals |
| WorknetId Format | `(chainId << 64) \| localCounter` | Globally unique |

---

## 8. Common User Workflows

### 8.1 Stake AWP
1. `AWPToken.approve(StakeNFT, amount)` or use `depositWithPermit`
2. `StakeNFT.deposit(amount, lockDuration)` → receive position NFT (tokenId)
3. Check: `staking.getBalance(address)` or `staking.getPositions(address)`

### 8.2 Allocate Stake to Agent
1. `StakingVault.allocate(staker, agent, worknetId, amount)` — caller must be staker or delegate
2. Check: `staking.getAllocations(address)`

### 8.3 Register a Worknet
1. `AWPToken.approve(AWPRegistry, 100000e18)` (registration cost)
2. `AWPRegistry.registerWorknet({name, symbol, worknetManager: 0x0, salt: 0x0, minStake: 0, skillsURI: ""})` → worknetId
3. `AWPRegistry.activateWorknet(worknetId)` — makes it eligible for emission
4. Check: `subnets.get(worknetId)`

### 8.4 Bind Agent to Owner
1. `AWPRegistry.bind(ownerAddress)` — called by agent
2. Check: `address.check(agentAddress)` → shows `boundTo`

### 8.5 Set Reward Recipient
1. `AWPRegistry.setRecipient(recipientAddress)`
2. Check: `address.resolveRecipient(address)`

### 8.6 Vote on DAO Proposal
1. Must hold StakeNFT positions with `createdAt < proposalCreatedAt`
2. `AWPDAO.castVoteWithReasonAndParams(proposalId, support, reason, abi.encode(tokenIds))`
3. Check: `governance.getProposal(proposalId)`
