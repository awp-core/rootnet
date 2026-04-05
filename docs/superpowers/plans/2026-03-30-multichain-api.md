# Multi-Chain API/Indexer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the API and Indexer support multiple EVM chains simultaneously — one indexer goroutine per chain, unified DB with chain_id columns, chain-aware API endpoints.

**Architecture:** Config loads chains from `chains.yaml`. Indexer spawns one polling goroutine per chain. DB schema adds `chain_id` to all chain-scoped tables (composite primary keys). API endpoints accept optional `?chainId=` filter and return `chain_id`/`chain_name` in responses. New `GET /api/chains` endpoint lists supported chains.

**Tech Stack:** Go 1.26, Chi v5, sqlc, pgx/v5, PostgreSQL, Redis, gopkg.in/yaml.v3

---

### Task 1: Chain Config Loader

**Files:**
- Create: `api/internal/config/chains.go`
- Modify: `api/internal/config/config.go`

- [ ] **Step 1: Create chains.go**

Create `api/internal/config/chains.go`:

```go
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ChainConfig holds per-chain deployment configuration
type ChainConfig struct {
	ChainID         int64  `yaml:"chainId"`
	Name            string `yaml:"name"`
	RPCURL          string `yaml:"rpcUrl"`
	DEX             string `yaml:"dex"`
	InitialMint     int64  `yaml:"initialMint"`
	Explorer        string `yaml:"explorer"`
	PoolManager     string `yaml:"poolManager"`
	PositionManager string `yaml:"positionManager"`
	Permit2         string `yaml:"permit2"`
	SwapRouter      string `yaml:"swapRouter"`
	StateView       string `yaml:"stateView"`
}

// ChainsFile is the top-level structure of chains.yaml
type ChainsFile struct {
	Chains map[string]ChainConfig `yaml:"chains"`
}

// LoadChains reads chains.yaml and resolves env vars in rpcUrl fields
func LoadChains(path string) ([]ChainConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read chains.yaml: %w", err)
	}
	var file ChainsFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse chains.yaml: %w", err)
	}
	chains := make([]ChainConfig, 0, len(file.Chains))
	for _, cfg := range file.Chains {
		cfg.RPCURL = os.ExpandEnv(cfg.RPCURL)
		chains = append(chains, cfg)
	}
	return chains, nil
}
```

- [ ] **Step 2: Add `ChainsFile` field to Config**

In `api/internal/config/config.go`, add:

```go
	// Multi-chain config file path (optional; if set, overrides single-chain CHAIN_ID/RPC_URL)
	ChainsFile string `env:"CHAINS_FILE" envDefault:""`
```

- [ ] **Step 3: Add go dependency**

```bash
cd api && go get gopkg.in/yaml.v3
```

- [ ] **Step 4: Build to verify**

```bash
go build ./...
```

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: chain config loader — LoadChains from chains.yaml

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 2: DB Schema — Add chain_id to All Tables

**Files:**
- Modify: `api/internal/db/schema.sql`

- [ ] **Step 1: Rewrite schema with chain_id columns**

The key changes:
- `users`: add `chain_id BIGINT NOT NULL`, change PK to `(chain_id, address)`
- `subnets`: add `chain_id BIGINT NOT NULL` (subnet_id already globally unique, keep as PK)
- `stake_allocations`: add `chain_id BIGINT NOT NULL`, change PK to `(chain_id, user_address, agent_address, subnet_id)`
- `user_balances`: add `chain_id BIGINT NOT NULL`, change PK to `(chain_id, user_address)`
- `stake_positions`: add `chain_id BIGINT NOT NULL`, change PK to `(chain_id, token_id)`
- `epochs`: add `chain_id BIGINT NOT NULL`, change PK to `(chain_id, epoch_id)`
- `recipient_awp_distributions`: add `chain_id BIGINT NOT NULL`, update unique index
- `sync_states`: add `chain_id BIGINT NOT NULL`, change PK to `(chain_id, contract_name)`
- `indexed_blocks`: add `chain_id BIGINT NOT NULL`, change PK to `(chain_id, block_number)`
- `vanity_salts`: add `chain_id BIGINT NOT NULL`

- [ ] **Step 2: Commit**

```bash
git add api/internal/db/schema.sql
git commit -m "schema: add chain_id to all tables for multi-chain indexing

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 3: SQL Queries — Add chain_id Parameter

**Files:**
- Modify: `api/internal/db/query/user.sql`
- Modify: `api/internal/db/query/subnet.sql`
- Modify: `api/internal/db/query/staking.sql`
- Modify: `api/internal/db/query/stake_position.sql`
- Modify: `api/internal/db/query/emission.sql`
- Modify: `api/internal/db/query/governance.sql`
- Modify: `api/internal/db/query/vanity_salt.sql`

- [ ] **Step 1: Update ALL queries to include chain_id**

Every INSERT/UPDATE/SELECT needs the `chain_id` parameter where the table has it. Key pattern:

```sql
-- name: GetUser :one
SELECT address, bound_to, recipient, registered_at, chain_id FROM users
WHERE address = $1 AND chain_id = $2;

-- name: ListUsers :many
SELECT address, bound_to, recipient, registered_at, chain_id FROM users
ORDER BY registered_at DESC LIMIT $1 OFFSET $2;
-- (no chain_id filter — returns all chains, API can filter)

-- name: InsertSubnet :exec
INSERT INTO subnets (subnet_id, chain_id, owner, ...) VALUES ($1, $2, ...)
ON CONFLICT (subnet_id) DO NOTHING;
```

For queries that aggregate across chains (e.g., `GetSubnetTotalStake`), use no chain_id filter — sum across all chains.

- [ ] **Step 2: Regenerate sqlc**

```bash
cd api && $(go env GOPATH)/bin/sqlc generate
```

- [ ] **Step 3: Build to verify**

```bash
go build ./...
```

Fix any compilation errors from changed generated types.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat: add chain_id to all SQL queries for multi-chain

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 4: Multi-Chain Indexer

**Files:**
- Modify: `api/internal/chain/indexer.go`
- Modify: `api/cmd/indexer/main.go`

- [ ] **Step 1: Update Indexer to accept chainId**

The `Indexer` struct needs a `chainId int64` field. All DB writes pass this chainId. The `NewIndexer` constructor takes `chainId`.

- [ ] **Step 2: Update indexer main.go for multi-chain**

When `CHAINS_FILE` env var is set, load chains from YAML and spawn one indexer goroutine per chain. When not set, fall back to single-chain mode (backward compatible).

```go
func main() {
    // If CHAINS_FILE is set, run multi-chain indexer
    // Otherwise, run single-chain (backward compat)
}
```

Each goroutine creates its own `chain.Client` and `chain.Indexer` with the chain's RPC URL and chainId.

- [ ] **Step 3: Build and test**

```bash
go build ./cmd/indexer
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat: multi-chain indexer — one goroutine per chain from chains.yaml

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 5: API Handlers — Chain-Aware Endpoints

**Files:**
- Modify: `api/internal/server/handler/handler.go` — GetRegistry, new GetChains
- Modify: `api/internal/server/handler/subnet.go` — chain_id in responses
- Modify: `api/internal/server/handler/staking.go` — global balance aggregation
- Modify: `api/internal/server/server.go` — new route

- [ ] **Step 1: Add `GET /api/chains` endpoint**

Returns the list of supported chains from chains.yaml:

```json
[
  {"chainId": 8453, "name": "Base", "explorer": "https://basescan.org"},
  {"chainId": 1, "name": "Ethereum", "explorer": "https://etherscan.io"},
  ...
]
```

- [ ] **Step 2: Update subnet responses to include chain_id and chain_name**

`GET /api/subnets` returns `chain_id` and `chain_name` for each subnet. The chain_name is resolved from the loaded chains config.

Add optional `?chainId=8453` query parameter to filter by chain.

- [ ] **Step 3: Add global staking balance endpoint**

`GET /api/staking/user/{addr}/global` aggregates StakeNFT positions and allocations across all chains.

- [ ] **Step 4: Update /api/registry to support ?chainId**

`GET /api/registry?chainId=8453` returns contracts for that specific chain. Without filter, returns the primary chain.

- [ ] **Step 5: Register new routes in server.go**

```go
r.Get("/chains", h.GetChains)
```

- [ ] **Step 6: Build and test**

```bash
go build ./...
go test ./internal/server/handler/ -count=1 -timeout 60s
```

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "feat: chain-aware API — /api/chains, chainId filter, global staking balance

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 6: Multi-Chain Keeper

**Files:**
- Modify: `api/cmd/keeper/main.go`

- [ ] **Step 1: Update keeper for multi-chain**

Same pattern as indexer: when `CHAINS_FILE` is set, spawn one keeper goroutine per chain. Each keeper manages its own chain's AWPEmission settlement and cache.

Redis keys become chain-scoped: `emission_current:{chainId}`, `awp_info:{chainId}`.

- [ ] **Step 2: Build**

```bash
go build ./cmd/keeper
```

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "feat: multi-chain keeper — per-chain emission settlement and cache

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 7: Update Documentation and Tests

**Files:**
- Modify: `api/.env.example` or create `api/.env.multichain.example`
- Modify: `skills-dev/config.md`
- Modify: `skills-dev/rest-api.md`

- [ ] **Step 1: Add multi-chain API docs**

Document new endpoints:
- `GET /api/chains`
- `?chainId=` filter on `/api/subnets`, `/api/registry`
- `GET /api/staking/user/{addr}/global`
- Chain-scoped Redis keys

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "docs: multi-chain API documentation

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

## Task Summary

| Task | Description | Estimate |
|------|-------------|----------|
| 1 | Chain config loader (chains.go) | 3 min |
| 2 | DB schema — chain_id columns | 5 min |
| 3 | SQL queries — chain_id params | 10 min |
| 4 | Multi-chain indexer | 10 min |
| 5 | Chain-aware API endpoints | 10 min |
| 6 | Multi-chain keeper | 5 min |
| 7 | Documentation | 3 min |
