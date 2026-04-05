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
| AWPToken | `0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1` |
| AWPRegistry | `0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A` |
| Treasury | `0x82562023a053025F3201785160CaE6051efD759e` |
| WorknetTokenFactory | `0x0000D4996BDBb99c772e3fA9f0e94AB52AAFFAC7` |
| AWPWorkNet | `0x00000bfbdEf8533E5F3228c9C846522D906100A7` |
| LPManagerUni | `0x00001961b9AcCD86b72DE19Be24FaD6f7c5b00A2` |
| AWPEmission (proxy) | `0x3C9cB73f8B81083882c5308Cce4F31f93600EaA9` |
| AWPAllocator | `0x0000D6BB5e040E35081b3AaF59DD71b21C9800AA` |
| veAWP | `0x0000b534C63D78212f1BDCc315165852793A00A8` |
| WorknetManagerUni (impl) | `0x000011EE4117c52dC0Eb146cBC844cb155B200A9` |
| AWPDAO | `0x00006879f79f3Da189b5D0fF6e58ad0127Cc0DA0` |

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
| `AWP_ALLOCATOR_ADDRESS` | Yes | AWPAllocator contract address |
| `VEAWP_ADDRESS` | Yes | veAWP contract address |
| `AWP_WORKNET_ADDRESS` | Yes | AWPWorkNet contract address |
| `LP_MANAGER_ADDRESS` | Yes | LPManager/LPManagerUni contract address |
| `WORKNET_TOKEN_FACTORY_ADDRESS` | Yes | WorknetTokenFactory address (also enables `/api/vanity/*`) |
| `DAO_ADDRESS` | Yes | AWPDAO contract address |
| `TREASURY_ADDRESS` | Yes | Treasury (TimelockController) address |
| `GUARDIAN_ADDRESS` | Yes | Guardian (cross-chain multisig) address |
| `DEPLOY_BLOCK` | Yes | Block number at which contracts were deployed (indexer start) |
| `KEEPER_PRIVATE_KEY` | No | Private key for keeper service (epoch settlement) |
| `RELAYER_PRIVATE_KEY` | No | Private key for gasless relay (enables `/api/relay/*`) |
| `WORKNET_TOKEN_BYTECODE_HASH` | No | `keccak256(WorknetToken.creationCode)` hex (enables vanity mining) |
| `VANITY_RULE` | No | `WorknetTokenFactory.vanityRule()` uint64 hex (e.g. `0x1001FFFF12101514`) |

## DEX Integration

> DEX addresses (PoolManager, PositionManager, Permit2, SwapRouter/UniversalRouter, StateView)
> are chain-specific and configured at deployment. The deploy script auto-selects:
> - **Base/Ethereum**: Uniswap V4 → LPManagerUni + WorknetManagerUni
> - **BSC**: PancakeSwap V4 → LPManager + WorknetManager

## Protocol Constants

| Constant | Value | Description |
|----------|-------|-------------|
| AWP MAX_SUPPLY | 10B (10^28 wei) | Total AWP supply cap |
| WorknetToken MAX_SUPPLY | 10B per worknet | Independent per-worknet cap |
| INITIAL_DAILY_EMISSION | 15.8M AWP | First epoch daily emission |
| EPOCH_DURATION | 1 day (86400s) | Time between settlements (daily epochs, AWPEmission only) |
| DECAY_FACTOR | 0.996844 per epoch | ~0.3156% decay each epoch |
| EMISSION_SPLIT | 100% to recipients | Guardian includes treasury in recipients for DAO share |
| MAX_ACTIVE_WORKNETS | 10,000 | AWPRegistry active worknet limit |
| maxRecipients | 10,000 | AWPEmission recipient limit |
| MAX_WEIGHT_SECONDS | 54 weeks (32,659,200s) | Max time for voting power sqrt |
| Timelock Delay | 2 days | Governance execution delay |
| Pool Fee | 1% (10000) | DEX V4 CL fee |
| Tick Spacing | 200 | DEX V4 CL tick spacing |
| WorknetTokenFactory.vanityRule | uint64, immutable | 0 = no validation; non-zero = 8-position EIP-55 pattern |

## Token Decimals

All tokens use **18 decimals**. Amounts in the API and contracts are in **wei** (smallest unit).

```
1 AWP = 1,000,000,000,000,000,000 wei = 10^18 wei
1 WorknetToken = 10^18 wei
```

## Redis Keys (internal)

| Key | TTL | Updated By | Content |
|-----|-----|------------|---------|
| `emission_current:{chainId}` | 30s | Keeper (25s interval) | `{epoch, settledEpoch, dailyEmission, totalWeight}` |
| `awp_info:{chainId}` | 1m | Keeper (25s interval) | `{totalSupply, maxSupply}` |
| `worknet_token_price:{worknetId}` | 10m | External | `{priceInAWP, reserve0, reserve1, updatedAt}` |
| `ratelimit:config` | — | admin.sh / Redis CLI | Rate limit configs (hash: name → "limit:window_seconds") |
| `rl:relay:{ip}` | 1h | Relay handler | Relay IP request counter |
| `rl:upload_salts:{ip}` | 1h | Salt handler | Upload salts IP counter |
| `rl:compute_salt:{ip}` | 1h | Vanity handler | Compute salt IP counter |
| `chain_events` (Pub/Sub) | — | Indexer | Real-time event stream |

> Rate limits are configured via Redis hash `ratelimit:config` (hot-updatable, no restart). Defaults compiled into the ratelimit package. Use `scripts/admin.sh` to manage.
