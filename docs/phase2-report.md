# Phase 2 Development Report — Go Backend

> **Note**: This is a historical report from Phase 2. Contract names have since been renamed: AlphaToken -> WorknetToken, SubnetNFT -> AWPWorkNet, StakingVault -> AWPAllocator, StakeNFT -> veAWP, SubnetManager -> WorknetManager. API event names changed from Subnet* to Worknet* prefix.

## Overview

The complete Go backend for AWP has been implemented, consisting of three independent processes (API server, Chain Indexer, Keeper Bot) sharing PostgreSQL and Redis for data and inter-process communication.

## Go Packages Implemented

| Package | Path | Lines | Description |
|---------|------|-------|-------------|
| config | internal/config/ | 41 | Environment variable configuration |
| chain/client | internal/chain/client.go | 138 | Ethereum RPC client with contract bindings |
| chain/indexer | internal/chain/indexer.go | 555 | Block event indexer → PostgreSQL + Redis Pub/Sub |
| chain/keeper | internal/chain/keeper.go | 268 | Scheduled on-chain operations + cache updates |
| handler | internal/server/handler/ | 891 | HTTP API handlers (8 files) |
| server | internal/server/server.go | 97 | Chi router setup with middleware |
| ws | internal/server/ws/hub.go | 210 | WebSocket hub with Redis Pub/Sub subscription |
| cmd/api | cmd/api/main.go | 81 | API server entry point (uber-go/fx) |
| cmd/indexer | cmd/indexer/main.go | 74 | Indexer entry point (uber-go/fx) |
| cmd/keeper | cmd/keeper/main.go | 94 | Keeper entry point (uber-go/fx) |

**Total hand-written Go code**: ~2,449 lines (excluding generated code)

### Generated Code

| Generated | Path | Tool |
|-----------|------|------|
| Contract bindings | internal/chain/bindings/ (7 files) | abigen |
| SQL query functions | internal/db/gen/ (8 files) | sqlc |

## API Routes Implemented

**31 routes** across 8 domain groups:

| Group | Routes | Endpoints |
|-------|--------|-----------|
| System | 2 | GET /health, GET /registry |
| Users | 3 | GET /users/, /users/count, /users/{address} |
| Address | 1 | GET /address/{address}/check |
| Agents | 4 | GET by-owner/{owner}, by-owner/{owner}/{agent}, lookup/{agent}, POST batch-info |
| Staking | 7 | GET balance, allocations, pending, frozen, agent stake, agent subnets, subnet total |
| Subnets | 4 | GET list, detail, earnings, agent info |
| Emission | 3 | GET current, schedule, epochs |
| Tokens | 3 | GET awp, alpha/{subnetId}, alpha/{subnetId}/price |
| Governance | 3 | GET proposals, proposal/{id}, treasury |
| WebSocket | 1 | GET /ws/live |

## sqlc Queries Generated

**66 queries** across 5 SQL files:

| File | Queries | Description |
|------|---------|-------------|
| user.sql | 19 | User registration, balances, reward recipients, withdrawals |
| agent.sql | 7 | Agent registration, lookup, management |
| subnet.sql | 13 | Subnet CRUD, lifecycle, metadata |
| staking.sql | 16 | Stake allocations, agent subnets, frozen handling |
| emission.sql | 9 | Epochs, distributions |
| governance.sql | 9 | Proposals, sync state |

## Test Results

**19 tests, 19 passed, 0 failed**

Tests cover:
- Health check endpoint
- User CRUD (list empty, count with data, not found)
- Address check (unknown, registered user)
- Subnet listing (empty)
- Epoch listing (empty)
- Emission schedule calculation (with no epoch data)
- Governance proposals (empty list, treasury address)
- Token endpoints (AWP info, Alpha price - both cache miss and hit)
- Redis cache integration (AWP info, emission current with cached data)
- Agent lookup (empty results, not found)
- Staking balance (not found returns default)

## Build Artifacts

All 3 binaries compile successfully:
- `bin/api` (19MB) — HTTP + WebSocket server
- `bin/indexer` (21MB) — Chain event indexer
- `bin/keeper` (21MB) — Keeper bot

## Lint Results

**0 issues** — golangci-lint v2 with errcheck, govet, staticcheck, unused enabled.

## Design Decisions

### 1. LPManager Simplified for Indexer
The indexer handles LPCreated events by simply updating the subnet's lp_pool field. Actual LP pool interaction (reserves, prices) is deferred to the Keeper's updateTokenPrices function, which will read from PancakeSwap V4 when available.

### 2. Emission Schedule Fallback
When no epochs exist in the DB (fresh deployment), GetEmissionSchedule uses the initial daily emission (15.8M AWP) as the starting point for projections, rather than returning an error.

### 3. WebSocket Library Migration
Migrated from deprecated `nhooyr.io/websocket` to `github.com/coder/websocket` (maintained fork, identical API).

### 4. Balance Not Found Handling
GET /staking/user/{address}/balance returns `{"totalBalance":"0","totalAllocated":"0","unallocated":"0"}` for unknown addresses instead of 404, since a user may exist on-chain but not yet indexed.

### 5. Batch Agent Info
The POST /agents/batch-info endpoint accepts a JSON body `{"agents":["0x..."],"subnetId":1}` and queries allocations for each agent, suitable for Coordinator cold-start data fetching.

### 6. Abigen Bindings
Generated Go bindings for 7 contracts (AWPRegistry, AWPToken, AlphaToken, SubnetNFT, StakingVault, AWPEmission, AWPDAO) directly from Foundry compiled ABI output.

### 7. Indexer Transaction Model
Each batch of events from a range of blocks is processed within a single PostgreSQL transaction. If any event handler fails, the entire batch is rolled back, ensuring atomicity. After successful commit, events are published to Redis Pub/Sub for WebSocket distribution.

### 8. Test Infrastructure
Tests use a separate `cortexia_test` database and Redis DB 1 to avoid interfering with development data. Tests are skipped if database/Redis is unavailable.

## TODOs for Phase 3

1. **Frontend Dashboard**: Next.js 14 + Tailwind + wagmi/viem dashboard consuming the API
2. **OpenTelemetry**: Add tracing and metrics (Prometheus/Jaeger integration)
3. **Docker Compose**: Multi-container setup with PostgreSQL, Redis, API, Indexer, Keeper
4. **PancakeSwap V4 Price Feed**: Implement actual pool reserve reading in Keeper's updateTokenPrices
5. **Rate Limiting**: Add per-IP rate limiting middleware to the API
6. **Pagination Metadata**: Add total count and page metadata to paginated responses
7. **WebSocket Filters**: Client-side event filtering is implemented but needs end-to-end testing
8. **Health Check Enhancement**: Add DB and Redis connectivity checks to /health
9. **Keeper Gas Management**: Implement gas price estimation and nonce management for on-chain transactions
10. **Event Replay**: Add ability to replay events from a specific block for recovery
11. **Integration Tests**: Add end-to-end tests with Anvil for full chain → indexer → API flow
