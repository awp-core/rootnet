# Configuration Reference

## Chain Configuration

### Base Mainnet (Primary Deployment)

| Parameter | Value |
|-----------|-------|
| Network | Base Mainnet |
| Chain ID | 8453 |
| RPC URL | `https://mainnet.base.org` |
| Block Time | ~2 seconds |
| Explorer | https://basescan.org |
| DEX | Uniswap V4 |
| Permit2 | `0x000000000022D473030F116dDEE9F6B43aC78BA3` |

### BSC Mainnet (Alternate)

| Parameter | Value |
|-----------|-------|
| Network | BNB Smart Chain |
| Chain ID | 56 |
| Block Time | ~3 seconds |
| Explorer | https://bscscan.com |
| DEX | PancakeSwap V4 |
| Permit2 | `0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768` |

## Contract Addresses (Base Mainnet)

| Contract | Address |
|----------|---------|
| AWPToken | `0x0000d0e38e9c6ba147b0098bb42007b942ef00a1` |
| AWPRegistry | `0x00003a7fa04c3af3adba2dc3c6622277501400b1` |
| Treasury | `0x9ee82684e4214edb405d930001e9058d1913d994` |
| AlphaTokenFactory | `0x3ebe3168c898f4b05ebf0c0d17f4739e111e5164` |
| SubnetNFT | `0x0f86ec2f2fbf234b00b18e66e7c5e00518091cda` |
| LPManagerUni | `0x2703d681ff3f7c4dc9eeed6f3ebaba3e82f8ebae` |
| AWPEmission (proxy) | `0xd31b6fedf7e568091b7fcf3cb5aac86c3a0ef1cf` |
| StakingVault | `0x0367e9c2f79ab35dc65e6876405a747882296fca` |
| StakeNFT | `0x4f7e8d4487c0c514b72ed0e35ed707cb8acdce39` |
| SubnetManagerUni (impl) | `0x567882378dcc11ec0d763fc5ca6c862487bbe574` |
| AWPDAO | `0x7171211da849a2c569643fb1e8f5399ddd71939a` |

> Contract addresses are also available via `GET /api/registry` (returns all addresses + chainId).

## API Configuration

| Parameter | Value |
|-----------|-------|
| REST Base URL | `https://<api-host>/api` |
| WebSocket URL | `wss://<api-host>/ws/live` |
| HTTP Port | 8080 (default) |

## Environment Variables (API Server)

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `REDIS_URL` | Yes | Redis connection string |
| `RPC_URL` | Yes | Chain RPC endpoint |
| `CHAIN_ID` | Yes | Target chain ID (8453 for Base, 56 for BSC) |
| `AWP_REGISTRY_ADDRESS` | Yes | AWPRegistry contract address |
| `AWP_TOKEN_ADDRESS` | Yes | AWPToken contract address |
| `AWP_EMISSION_ADDRESS` | Yes | AWPEmission proxy address |
| `STAKING_VAULT_ADDRESS` | Yes | StakingVault contract address |
| `STAKE_NFT_ADDRESS` | Yes | StakeNFT contract address |
| `SUBNETNFT_ADDRESS` | Yes | SubnetNFT contract address |
| `LP_MANAGER_ADDRESS` | Yes | LPManager/LPManagerUni contract address |
| `ALPHA_FACTORY_ADDRESS` | Yes | AlphaTokenFactory address (also enables `/api/vanity/*`) |
| `DAO_ADDRESS` | Yes | AWPDAO contract address |
| `TREASURY_ADDRESS` | Yes | Treasury (TimelockController) address |
| `DEPLOY_BLOCK` | Yes | Block number at which contracts were deployed (indexer start) |
| `KEEPER_PRIVATE_KEY` | No | Private key for keeper service (epoch settlement) |
| `RELAYER_PRIVATE_KEY` | No | Private key for gasless relay (enables `/api/relay/*`) |
| `ALPHA_INITCODE_HASH` | No | `keccak256(AlphaToken.creationCode)` hex (enables vanity mining) |
| `VANITY_RULE` | No | `AlphaTokenFactory.vanityRule()` uint64 hex (e.g. `0x0A01FFFF0C0A0F0E`) |

## DEX Integration

> DEX addresses (PoolManager, PositionManager, Permit2, SwapRouter/UniversalRouter, StateView)
> are chain-specific and configured at deployment. The deploy script auto-selects:
> - **Base/Ethereum**: Uniswap V4 → LPManagerUni + SubnetManagerUni
> - **BSC**: PancakeSwap V4 → LPManager + SubnetManager

## Protocol Constants

| Constant | Value | Description |
|----------|-------|-------------|
| AWP MAX_SUPPLY | 10B (10^28 wei) | Total AWP supply cap |
| Alpha MAX_SUPPLY | 10B per subnet | Independent per-subnet cap |
| INITIAL_DAILY_EMISSION | 15.8M AWP | **[DRAFT]** First epoch daily emission |
| EPOCH_DURATION | 1 day (86400s) | **[DRAFT]** Time between settlements (daily epochs, AWPEmission only) |
| DECAY_FACTOR | 0.996844 per epoch | **[DRAFT]** ~0.3156% decay each epoch |
| EMISSION_SPLIT | 50/50 | **[DRAFT]** Recipients vs DAO |
| MAX_ACTIVE_SUBNETS | 10,000 | AWPRegistry active subnet limit |
| maxRecipients | 10,000 | AWPEmission recipient limit |
| MAX_WEIGHT_SECONDS | 54 weeks (32,659,200s) | Max time for voting power sqrt |
| Immunity Period | 30 days | Deregister protection |
| Timelock Delay | 2 days | Governance execution delay |
| Pool Fee | 1% (10000) | DEX V4 CL fee |
| Tick Spacing | 200 | DEX V4 CL tick spacing |
| Oracle Threshold | DAO-configured | Initially 3/5 recommended |
| AlphaTokenFactory.vanityRule | uint64, immutable | 0 = no validation; non-zero = 8-position EIP-55 pattern |

## Token Decimals

All tokens use **18 decimals**. Amounts in the API and contracts are in **wei** (smallest unit).

```
1 AWP = 1,000,000,000,000,000,000 wei = 10^18 wei
1 Alpha = 10^18 wei
```

## Redis Keys (internal)

| Key | TTL | Updated By | Content |
|-----|-----|------------|---------|
| `emission_current` | 30s | Keeper (25s interval) | `{epoch, settledEpoch, dailyEmission, totalWeight}` |
| `awp_info` | 1m | Keeper (25s interval) | `{totalSupply, maxSupply}` |
| `alpha_price:{subnetId}` | 10m | External | `{priceInAWP, reserve0, reserve1, updatedAt}` |
| `ratelimit:config` | — | admin.sh / Redis CLI | Rate limit configs (hash: name → "limit:window_seconds") |
| `rl:relay:{ip}` | 1h | Relay handler | Relay IP request counter |
| `rl:upload_salts:{ip}` | 1h | Salt handler | Upload salts IP counter |
| `rl:compute_salt:{ip}` | 1h | Vanity handler | Compute salt IP counter |
| `chain_events` (Pub/Sub) | — | Indexer | Real-time event stream |

> Rate limits are configured via Redis hash `ratelimit:config` (hot-updatable, no restart). Defaults compiled into the ratelimit package. Use `scripts/admin.sh` to manage.
