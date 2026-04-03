# E2E Tests

## How to Run

E2E tests require Anvil local chain + PostgreSQL + Redis. Since contracts compile slowly with `via_ir=true` (~60s),
it is recommended to run on CI environments or machines with sufficient computing power.

### Prerequisites

1. Start Anvil: `anvil --port 18545`
2. Deploy contracts: `cd contracts && forge script script/TestDeploy.s.sol:TestDeploy --rpc-url http://127.0.0.1:18545 --broadcast`
3. Start PostgreSQL + Redis
4. Create test database: `createdb awp_test`
5. Apply schema: `psql awp_test < api/internal/db/schema.sql`

### Run

```bash
TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/awp_test?sslmode=disable" \
TEST_REDIS_URL="redis://localhost:6379/1" \
go test ./test/e2e/ -v -timeout 180s
```
