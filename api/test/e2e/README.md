# E2E Tests

## 运行方式

E2E 测试需要 Anvil 本地链 + PostgreSQL + Redis。由于合约使用 `via_ir=true` 编译较慢（~60s），
建议在 CI 环境或有足够算力的机器上运行。

### 前置条件

1. 启动 Anvil: `anvil --port 18545`
2. 部署合约: `cd contracts && forge script script/TestDeploy.s.sol:TestDeploy --rpc-url http://127.0.0.1:18545 --broadcast`
3. 启动 PostgreSQL + Redis
4. 创建测试数据库: `createdb cortexia_test`
5. 应用 schema: `psql cortexia_test < api/internal/db/schema.sql`

### 运行

```bash
TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/cortexia_test?sslmode=disable" \
TEST_REDIS_URL="redis://localhost:6379/1" \
go test ./test/e2e/ -v -timeout 180s
```
