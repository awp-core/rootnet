# Gasless Staking API

## POST /api/relay/stake

Gasless AWP staking via [VeAWPHelper](../contracts/src/core/VeAWPHelper.sol). The user signs a single ERC-2612 permit off-chain; the relayer pays gas and executes the deposit on behalf of the user.

**Contract:** `VeAWPHelper` — `0x0000561EDE5C1Ba0b81cE585964050bEAE730001` (same address on all 4 chains)

---

## Request

```http
POST https://api.awp.sh/api/relay/stake
Content-Type: application/json
```

```json
{
  "chainId": 8453,
  "user": "0x1234...abcd",
  "amount": "10000000000000000000000",
  "lockDuration": 7776000,
  "deadline": 1775649600,
  "signature": "0x..."
}
```

| Field | Type | Description |
|-------|------|-------------|
| `chainId` | int64 | Chain ID (8453 Base, 1 Ethereum, 42161 Arbitrum, 56 BSC) |
| `user` | string | User address (permit signer, receives the veAWP NFT) |
| `amount` | string | AWP amount in wei (e.g. `"10000000000000000000000"` = 10,000 AWP) |
| `lockDuration` | uint64 | Lock period in seconds (minimum 86400 = 1 day) |
| `deadline` | uint64 | Permit expiry timestamp (must be > current block.timestamp) |
| `signature` | string | ERC-2612 permit signature (65 bytes, 0x-prefixed hex) |

## Response

**Success (200):**
```json
{
  "txHash": "0xfa18553e2ce7876d5e3c6069137193cfb781248f914c494cf3f28dcdbad28168"
}
```

**Error (400/429):**
```json
{
  "error": "AWP transfer to helper failed"
}
```

---

## How to Sign the Permit

The user signs an [ERC-2612](https://eips.ethereum.org/EIPS/eip-2612) permit that authorizes VeAWPHelper to transfer their AWP tokens.

### EIP-712 Domain

```json
{
  "name": "AWP Token",
  "version": "1",
  "chainId": 8453,
  "verifyingContract": "0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1"
}
```

> The domain is the **AWP Token** contract (not AWPRegistry or VeAWPHelper). This is standard ERC-2612 — the token contract verifies the permit.

### EIP-712 Types

```
Permit(address owner, address spender, uint256 value, uint256 nonce, uint256 deadline)
```

### Message Fields

| Field | Value |
|-------|-------|
| `owner` | User's address |
| `spender` | `0x0000561EDE5C1Ba0b81cE585964050bEAE730001` (VeAWPHelper) |
| `value` | AWP amount in wei |
| `nonce` | Current permit nonce from `AWPToken.nonces(user)` |
| `deadline` | Expiry timestamp (same value sent in the relay request) |

### Getting the Spender Address and Nonce

```bash
# 1. Get VeAWPHelper address from the registry API
curl -s -X POST https://api.awp.sh/v2 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"registry.get","params":{"chainId":8453},"id":1}'
# → result.veAWPHelper = "0x0000561EDE5C1Ba0b81cE585964050bEAE730001"

# 2. Get the user's permit nonce (on-chain query)
cast call 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1 \
  "nonces(address)(uint256)" 0xYOUR_ADDRESS \
  --rpc-url https://mainnet.base.org
```

---

## Signing Example (viem/TypeScript)

```typescript
import { parseEther } from 'viem'

const AWP_TOKEN = '0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1'
const VE_AWP_HELPER = '0x0000561EDE5C1Ba0b81cE585964050bEAE730001'

// 1. Get permit nonce
const nonce = await publicClient.readContract({
  address: AWP_TOKEN,
  abi: [{ name: 'nonces', type: 'function', inputs: [{ type: 'address' }], outputs: [{ type: 'uint256' }] }],
  functionName: 'nonces',
  args: [userAddress],
})

// 2. Sign ERC-2612 permit
const deadline = BigInt(Math.floor(Date.now() / 1000) + 3600) // 1 hour
const amount = parseEther('10000') // 10,000 AWP

const signature = await walletClient.signTypedData({
  domain: {
    name: 'AWP Token',
    version: '1',
    chainId: 8453,
    verifyingContract: AWP_TOKEN,
  },
  types: {
    Permit: [
      { name: 'owner', type: 'address' },
      { name: 'spender', type: 'address' },
      { name: 'value', type: 'uint256' },
      { name: 'nonce', type: 'uint256' },
      { name: 'deadline', type: 'uint256' },
    ],
  },
  primaryType: 'Permit',
  message: {
    owner: userAddress,
    spender: VE_AWP_HELPER,
    value: amount,
    nonce: nonce,
    deadline: deadline,
  },
})

// 3. Submit to relay
const response = await fetch('https://api.awp.sh/api/relay/stake', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    chainId: 8453,
    user: userAddress,
    amount: amount.toString(),
    lockDuration: 7776000, // 90 days
    deadline: Number(deadline),
    signature: signature,
  }),
})
const { txHash } = await response.json()
```

## Signing Example (Python/eth_account)

```python
from eth_account.messages import encode_typed_data
from eth_account import Account

AWP_TOKEN = "0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1"
VE_AWP_HELPER = "0x0000561EDE5C1Ba0b81cE585964050bEAE730001"

# Build EIP-712 typed data
typed_data = {
    "types": {
        "EIP712Domain": [
            {"name": "name", "type": "string"},
            {"name": "version", "type": "string"},
            {"name": "chainId", "type": "uint256"},
            {"name": "verifyingContract", "type": "address"},
        ],
        "Permit": [
            {"name": "owner", "type": "address"},
            {"name": "spender", "type": "address"},
            {"name": "value", "type": "uint256"},
            {"name": "nonce", "type": "uint256"},
            {"name": "deadline", "type": "uint256"},
        ],
    },
    "primaryType": "Permit",
    "domain": {
        "name": "AWP Token",
        "version": "1",
        "chainId": 8453,
        "verifyingContract": AWP_TOKEN,
    },
    "message": {
        "owner": user_address,
        "spender": VE_AWP_HELPER,
        "value": amount_wei,
        "nonce": nonce,         # from AWPToken.nonces(user)
        "deadline": deadline,
    },
}

# Sign
signed = Account.sign_typed_data(private_key, typed_data)
signature = signed.signature.hex()
```

---

## On-Chain Flow

```
User signs permit(owner=user, spender=VeAWPHelper, value=amount)
  │
  ▼
POST /api/relay/stake { user, amount, lockDuration, deadline, signature }
  │
  ▼
Relayer calls VeAWPHelper.depositFor(user, amount, lockDuration, deadline, v, r, s)
  │
  ├── 1. permit(user → VeAWPHelper)     — sets AWP allowance
  ├── 2. transferFrom(user → helper)     — pulls AWP from user
  ├── 3. veAWP.deposit(amount, lock)     — creates position, mints NFT to helper
  └── 4. transferFrom(helper → user)     — sends veAWP NFT to user
  │
  ▼
User receives veAWP position NFT (staked AWP, locked for lockDuration)
```

---

## Error Reference

| Selector | Error | Description |
|----------|-------|-------------|
| `0x1f2a2005` | `ZeroAmount()` | Amount must be > 0 |
| `0xd92e233d` | `ZeroAddress()` | User address must not be zero |
| `0xfd684c3b` | `InvalidUser()` | User cannot be the helper contract itself |
| `0x90b8ec18` | `TransferFailed()` | AWP transferFrom failed (insufficient balance or permit invalid) |
| `0xd4005715` | `LockTooShort()` | lockDuration < 1 day (86400 seconds) |
| `0x2c5211c6` | `InvalidAmount()` | Amount exceeds uint128 max |

---

## Contract Addresses

| Contract | Address | Chains |
|----------|---------|--------|
| VeAWPHelper | `0x0000561EDE5C1Ba0b81cE585964050bEAE730001` | All 4 |
| AWP Token | `0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1` | All 4 |
| veAWP | `0x0000b534C63D78212f1BDCc315165852793A00A8` | All 4 |

All addresses are identical across Base (8453), Ethereum (1), Arbitrum (42161), and BSC (56).

---

## Rate Limiting

The `/api/relay/stake` endpoint shares the relay rate limit: **100 requests per hour per IP** (configurable via Redis).

## Query Transaction Status

```http
GET https://api.awp.sh/api/relay/status/{txHash}
```

Returns `{ "status": "pending" | "confirmed" | "failed", "txHash": "0x...", "blockNumber": 12345 }`
