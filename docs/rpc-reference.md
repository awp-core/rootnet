# AWP JSON-RPC 2.0 API Reference

> **Endpoint**: `POST https://api.awp.sh/v2`
> **Discovery**: `GET https://api.awp.sh/v2` (returns `rpc.discover`)
> **Protocol**: JSON-RPC 2.0 (single and batch requests supported, max batch size: 20)
> **Chain**: Multi-chain (Base 8453, Ethereum 1, Arbitrum 42161, BSC 56). Most methods accept optional `chainId` parameter; omit for default chain.

---

## Request Format

```json
{
  "jsonrpc": "2.0",
  "method": "namespace.method",
  "params": { ... },
  "id": 1
}
```

## Response Format

```json
{
  "jsonrpc": "2.0",
  "result": { ... },
  "id": 1
}
```

## Error Codes

| Code | Meaning |
|------|---------|
| -32700 | Parse error |
| -32600 | Invalid request |
| -32601 | Method not found |
| -32602 | Invalid params |
| -32603 | Internal error |
| -32001 | Resource not found |

## Common Parameter Types

| Type | Description | Example |
|------|-------------|---------|
| `address` | 0x-prefixed 40 hex chars (case-insensitive) | `"0xAbC...123"` |
| `worknetId` | Globally unique ID: `(chainId << 64) \| localId` | `"36316036596842561537"` |
| `chainId` | Chain ID integer; omit or 0 for default | `8453` |
| `page` | 1-indexed page number (default 1) | `1` |
| `limit` | Items per page (default 20, max 100) | `50` |

---

## Methods

### stats

#### `stats.global`
Get global protocol statistics across all chains.

**Params**: none

**Response**:
```json
{
  "totalUsers": 1234,
  "totalWorknets": 56,
  "totalStaked": "1000000000000000000000000",
  "totalEmitted": "31600000000000000000000000",
  "chains": 4
}
```

---

### registry

#### `registry.get`
Get all contract addresses and EIP-712 domain info for a chain.

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `chainId` | integer | no | Chain ID |

**Response**:
```json
{
  "chainId": 8453,
  "awpRegistry": "0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A",
  "awpToken": "0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1",
  "awpEmission": "0x3C9cB73f8B81083882c5308Cce4F31f93600EaA9",
  "awpAllocator": "0xE8A204fD9c94C7E28bE11Af02fc4A4AC294Df29b",
  "veAWP": "0x4E119560632698Bab67cFAB5d8EC0A373363ba2d",
  "awpWorkNet": "0xB9F03539BE496d09c4d7964921d674B8763f5233",
  "lpManager": "0x3034E029e61e8c2fc525A7bC5E267Ad3837D72e3",
  "worknetTokenFactory": "0xB2e4897eD77d0f5BFa3140B9989594de09a8037c",
  "dao": "0x6a074aC9823c47f86EE4Fc7F62e4217Bc9C76004",
  "treasury": "0x82562023a053025F3201785160CaE6051efD759e",
  "eip712Domain": {
    "name": "AWPRegistry",
    "version": "1",
    "chainId": 8453,
    "verifyingContract": "0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A"
  }
}
```

---

### health

#### `health.check`
Basic health check.

**Params**: none

**Response**: `{"status": "ok"}`

#### `health.detailed`
Detailed health status including per-chain indexer/keeper state.

**Params**: none

---

### chains

#### `chains.list`
List all supported chains.

**Params**: none

**Response**: Array of `{chainId, name, dex, explorer}`

---

### users

#### `users.list`
List users (paginated, per-chain).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `chainId` | integer | no | Chain ID |
| `page` | integer | no | Page number |
| `limit` | integer | no | Items per page |

#### `users.listGlobal`
List users across all chains (deduplicated).

**Params**: `page`, `limit`

#### `users.count`
Get total user count.

**Params**: `chainId` (optional)

#### `users.get`
Get user details (balance, bound agents, recipient).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `address` | string | yes | User address |
| `chainId` | integer | no | Chain ID |

#### `users.getPortfolio`
Get full user portfolio (identity, balance, positions, allocations, delegates).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `address` | string | yes | User address |
| `chainId` | integer | no | Chain ID |

#### `users.getDelegates`
Get agents bound to a user (delegate tree).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `address` | string | yes | User address |
| `chainId` | integer | no | Chain ID |

---

### address

#### `address.check`
Check address registration status, binding, and recipient.

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `address` | string | yes | Address |
| `chainId` | integer | no | Chain ID |

**Response**:
```json
{
  "address": "0x...",
  "isRegisteredUser": true,
  "isRegisteredAgent": false,
  "boundTo": "0x...",
  "recipient": "0x...",
  "hasDelegate": false
}
```

#### `address.resolveRecipient`
Resolve the effective recipient by walking the bind chain to root (on-chain call).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `address` | string | yes | Address |
| `chainId` | integer | no | Chain ID for on-chain read |

**Response**: `{"address": "0x...", "resolvedRecipient": "0x..."}`

#### `address.batchResolveRecipients`
Batch resolve recipients (max 500 addresses, on-chain call).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `addresses` | array\<string\> | yes | Address list (max 500) |
| `chainId` | integer | no | Chain ID |

**Response**: Array of `{"address": "0x...", "resolvedRecipient": "0x..."}`

---

### nonce

#### `nonce.get`
Get AWPRegistry EIP-712 nonce (for bind, setRecipient, registerWorknet, grantDelegate, revokeDelegate, unbind).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `address` | string | yes | Address |
| `chainId` | integer | no | Chain ID |

**Response**: `{"nonce": 42}`

#### `nonce.getStaking`
Get AWPAllocator EIP-712 nonce (for allocateFor, deallocateFor).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `address` | string | yes | Address |
| `chainId` | integer | no | Chain ID |

**Response**: `{"nonce": 5}`

---

### agents

#### `agents.getByOwner`
Get all agents bound to an owner.

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `owner` | string | yes | Owner address |
| `chainId` | integer | no | Chain ID |

#### `agents.getDetail`
Get agent details (owner, bound worknets, allocations).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `agent` | string | yes | Agent address |
| `chainId` | integer | no | Chain ID |

#### `agents.lookup`
Look up agent owner address.

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `agent` | string | yes | Agent address |
| `chainId` | integer | no | Chain ID |

**Response**: `{"ownerAddress": "0x..."}`

#### `agents.batchInfo`
Batch query agent info and stake in a worknet (max 100 agents).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `agents` | array\<string\> | yes | Agent addresses (max 100) |
| `worknetId` | string | yes | Worknet ID |
| `chainId` | integer | no | Chain ID |

---

### staking

#### `staking.getBalance`
Get user AWP staking balance (staked/allocated/available).

**Params**: `address` (required), `chainId` (optional)

#### `staking.getUserBalanceGlobal`
Get user staking balance aggregated across all chains.

**Params**: `address` (required)

#### `staking.getPositions`
Get user veAWP positions (per-chain).

**Params**: `address` (required), `chainId` (optional)

#### `staking.getPositionsGlobal`
Get user veAWP positions across all chains.

**Params**: `address` (required)

#### `staking.getAllocations`
Get user stake allocations (paginated).

**Params**: `address` (required), `chainId` (optional), `page`, `limit`

#### `staking.getFrozen`
Get user frozen allocations (deprecated — always returns empty).

**Params**: `address` (required), `chainId` (optional)

#### `staking.getPending`
Get pending allocation changes (always returns empty array).

**Params**: none

#### `staking.getAgentSubnetStake`
Get agent's stake amount in a specific worknet (cross-chain, no chainId needed).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `agent` | string | yes | Agent address |
| `worknetId` | string | yes | Worknet ID |

**Response**: `{"amount": "1000000000000000000000"}`

#### `staking.getAgentSubnets`
Get all worknets an agent participates in (cross-chain).

**Params**: `agent` (required)

#### `staking.getSubnetTotalStake`
Get worknet total stake across all agents (cross-chain).

**Params**: `worknetId` (required)

**Response**: `{"total": "5000000000000000000000000"}`

---

### subnets (worknets)

#### `subnets.list`
List worknets (paginated, optional status filter). `chainId=0` returns all chains.

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `status` | string | no | `Pending`, `Active`, `Paused`, `Banned` |
| `chainId` | integer | no | 0 = all chains |
| `page` | integer | no | Page number |
| `limit` | integer | no | Items per page |

#### `subnets.listRanked`
List worknets ranked by total stake.

**Params**: `chainId` (optional), `page`, `limit`

#### `subnets.search`
Search worknets by name or symbol (case-insensitive ILIKE).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `query` | string | yes | Search string (1-100 chars) |
| `chainId` | integer | no | Chain ID |
| `page` | integer | no | Page number |
| `limit` | integer | no | Items per page |

#### `subnets.getByOwner`
Get worknets owned by an address.

**Params**: `owner` (required), `chainId` (optional), `page`, `limit`

#### `subnets.get`
Get worknet details.

**Params**: `worknetId` (required)

#### `subnets.getSkills`
Get worknet skills URI.

**Params**: `worknetId` (required)

#### `subnets.getEarnings`
Get worknet AWP earnings history (paginated).

**Params**: `worknetId` (required), `page`, `limit`

#### `subnets.getAgentInfo`
Get agent staking info in a worknet.

**Params**: `worknetId` (required), `agent` (required)

#### `subnets.listAgents`
List agents in a worknet ranked by stake.

**Params**: `worknetId` (required), `chainId` (optional), `page`, `limit`

---

### emission

#### `emission.getCurrent`
Get current emission data (epoch, daily emission, total weight).

**Params**: `chainId` (optional)

#### `emission.getSchedule`
Get emission projections (30/90/365 day forecasts with decay).

**Params**: `chainId` (optional)

#### `emission.getGlobalSchedule`
Get emission schedule aggregated across all chains.

**Params**: none

#### `emission.listEpochs`
List settled epochs (paginated).

**Params**: `chainId` (optional), `page`, `limit`

#### `emission.getEpochDetail`
Get epoch detail with per-recipient distributions.

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `epochId` | integer | yes | Epoch ID |
| `chainId` | integer | no | Chain ID |

---

### tokens

#### `tokens.getAWP`
Get AWP token info (totalSupply, maxSupply, minters).

**Params**: `chainId` (optional)

#### `tokens.getAWPGlobal`
Get AWP token info aggregated across all chains.

**Params**: none

#### `tokens.getAlphaInfo`
Get worknet WorknetToken info (name, symbol, totalSupply, worknetManager).

**Params**: `worknetId` (required)

#### `tokens.getAlphaPrice`
Get worknet WorknetToken price (from LP pool, cached in Redis).

**Params**: `worknetId` (required)

---

### governance

#### `governance.listProposals`
List governance proposals (per-chain, paginated, optional status filter).

**Params**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `status` | string | no | `Active`, `Canceled`, `Defeated`, `Succeeded`, `Queued`, `Expired`, `Executed` |
| `chainId` | integer | no | Chain ID |
| `page` | integer | no | Page number |
| `limit` | integer | no | Items per page |

#### `governance.listAllProposals`
List proposals across all chains.

**Params**: `status` (optional), `page`, `limit`

#### `governance.getProposal`
Get proposal details.

**Params**: `proposalId` (required), `chainId` (optional)

#### `governance.getTreasury`
Get Treasury contract address.

**Params**: none

**Response**: `{"treasuryAddress": "0x82562023a053025F3201785160CaE6051efD759e"}`

---

## Batch Request Example

```json
[
  {"jsonrpc": "2.0", "method": "users.get", "params": {"address": "0xAbC..."}, "id": 1},
  {"jsonrpc": "2.0", "method": "staking.getBalance", "params": {"address": "0xAbC..."}, "id": 2},
  {"jsonrpc": "2.0", "method": "emission.getCurrent", "params": {}, "id": 3}
]
```

Batch requests execute concurrently. Max 20 requests per batch.

---

## WebSocket

**Endpoint**: `wss://api.awp.sh/ws/live`

Real-time event stream. Supports optional address-based filtering via `watchAddresses` subscription message.

Events pushed: `UserRegistered`, `Bound`, `Unbound`, `RecipientSet`, `Deposited`, `Withdrawn`, `Allocated`, `Deallocated`, `Reallocated`, `WorknetRegistered`, `WorknetActivated`, `EpochSettled`, `RecipientAWPDistributed`, `AllocationsSubmitted`, etc.

---

## Rate Limits

- Nonce endpoints: configurable via Redis (`ratelimit:config` hash)
- Relay endpoints: 100 req/IP/hour (default)
- Batch agent info: rate limited per IP
- All limits hot-updatable via admin API

---

## Deployed Contract Addresses (Identical on all 4 chains)

| Contract | Address |
|----------|---------|
| AWPToken | `0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1` |
| AWPRegistry (proxy) | `0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A` |
| AWPEmission (proxy) | `0x3C9cB73f8B81083882c5308Cce4F31f93600EaA9` |
| AWPAllocator (proxy) | `0xE8A204fD9c94C7E28bE11Af02fc4A4AC294Df29b` |
| veAWP | `0x4E119560632698Bab67cFAB5d8EC0A373363ba2d` |
| AWPWorkNet | `0xB9F03539BE496d09c4d7964921d674B8763f5233` |
| WorknetTokenFactory | `0xB2e4897eD77d0f5BFa3140B9989594de09a8037c` |
| Treasury | `0x82562023a053025F3201785160CaE6051efD759e` |
| AWPDAO | `0x6a074aC9823c47f86EE4Fc7F62e4217Bc9C76004` |
| Guardian (Safe 3/5) | `0x000002bEfa6A1C99A710862Feb6dB50525dF00A3` |
