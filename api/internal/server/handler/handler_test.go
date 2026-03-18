package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/cortexia/rootnet/api/internal/config"
	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/cortexia/rootnet/api/internal/server"
	"github.com/cortexia/rootnet/api/internal/server/handler"
	"github.com/cortexia/rootnet/api/internal/server/ws"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// testEnv holds the test environment
type testEnv struct {
	pool    *pgxpool.Pool
	rdb     *redis.Client
	queries *gen.Queries
	handler *handler.Handler
	router  http.Handler
	t       *testing.T
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/cortexia_test?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skipf("skipping test: cannot connect to database: %v", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Skipf("skipping test: database unavailable: %v", err)
	}

	redisURL := os.Getenv("TEST_REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/1"
	}

	opt, _ := redis.ParseURL(redisURL)
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		pool.Close()
		t.Skipf("skipping test: Redis unavailable: %v", err)
	}

	queries := gen.New(pool)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		TreasuryAddress:     "0x1234567890abcdef1234567890abcdef12345678",
		RootNetAddress:      "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		SubnetNFTAddress:    "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		DAOAddress:          "0xcccccccccccccccccccccccccccccccccccccccc",
		AWPTokenAddress:     "0xdddddddddddddddddddddddddddddddddddddd",
		StakingVaultAddress: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		AWPEmissionAddress:  "0xffffffffffffffffffffffffffffffffffffffff",
	}

	h := handler.NewHandler(queries, rdb, cfg, logger)
	hub := ws.NewHub(rdb, logger)
	router := server.NewRouter(server.RouterParams{Handler: h, Hub: hub})

	env := &testEnv{
		pool:    pool,
		rdb:     rdb,
		queries: queries,
		handler: h,
		router:  router,
		t:       t,
	}

	// Clean up database
	env.cleanDB()

	t.Cleanup(func() {
		env.cleanDB()
		pool.Close()
		_ = rdb.Close()
	})

	return env
}

func (e *testEnv) cleanDB() {
	e.t.Helper()
	ctx := context.Background()
	// Truncate all tables in one statement
	_, err := e.pool.Exec(ctx, `TRUNCATE TABLE
		recipient_awp_distributions, stake_positions, stake_allocations,
		user_balances, user_reward_recipients, agents, epochs,
		subnets, proposals, users, sync_states`)
	if err != nil {
		e.t.Logf("cleanDB TRUNCATE failed: %v", err)
	}
	if err := e.rdb.FlushDB(ctx).Err(); err != nil {
		e.t.Logf("cleanDB FlushDB failed: %v", err)
	}
}

func (e *testEnv) request(method, path string, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rr := httptest.NewRecorder()
	e.router.ServeHTTP(rr, req)
	return rr
}

// numericFromInt64 converts an int64 to pgtype.Numeric
func numericFromInt64(v int64) pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(v), Exp: 0, Valid: true}
}

// ════════════════════════════════════════════════════════════
// System — 2 tests
// ════════════════════════════════════════════════════════════

func TestHealth(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/health", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["status"] != "ok" {
		t.Errorf("expected status=ok, got %s", resp["status"])
	}
}

func TestGetRegistry(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/registry", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)

	// Verify all contract addresses are returned from config
	if resp["rootNet"] != "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Errorf("unexpected rootNet: %s", resp["rootNet"])
	}
	if resp["subnetNFT"] != "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" {
		t.Errorf("unexpected subnetNFT: %s", resp["subnetNFT"])
	}
	if resp["dao"] != "0xcccccccccccccccccccccccccccccccccccccccc" {
		t.Errorf("unexpected dao: %s", resp["dao"])
	}
	if resp["awpToken"] != "0xdddddddddddddddddddddddddddddddddddddd" {
		t.Errorf("unexpected awpToken: %s", resp["awpToken"])
	}
	if resp["stakingVault"] != "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		t.Errorf("unexpected stakingVault: %s", resp["stakingVault"])
	}
	if resp["awpEmission"] != "0xffffffffffffffffffffffffffffffffffffffff" {
		t.Errorf("unexpected awpEmission: %s", resp["awpEmission"])
	}
	if resp["treasury"] != "0x1234567890abcdef1234567890abcdef12345678" {
		t.Errorf("unexpected treasury: %s", resp["treasury"])
	}
}

// ════════════════════════════════════════════════════════════
// Users — 6 tests
// ════════════════════════════════════════════════════════════

func TestListUsersEmpty(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/users/?page=1&limit=10", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var users []any
	_ = json.Unmarshal(rr.Body.Bytes(), &users)
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

func TestListUsersWithData(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// Insert 3 users (registered_at determines sort order, DESC)
	for i, ts := range []int64{1000, 2000, 3000} {
		addr := strings.Replace("0x0000000000000000000000000000000000000001", "1", string(rune('1'+i)), 1)
		_ = env.queries.InsertUser(ctx, gen.InsertUserParams{
			Address:      addr,
			RegisteredAt: ts,
		})
	}

	t.Run("AllUsers", func(t *testing.T) {
		rr := env.request("GET", "/api/users/?page=1&limit=10", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var users []map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &users)
		if len(users) != 3 {
			t.Errorf("expected 3 users, got %d", len(users))
		}
	})

	t.Run("Pagination", func(t *testing.T) {
		rr := env.request("GET", "/api/users/?page=2&limit=2", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var users []map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &users)
		if len(users) != 1 {
			t.Errorf("expected 1 user on page 2 (limit=2), got %d", len(users))
		}
	})
}

func TestGetUserCount(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// Insert 2 users
	_ = env.queries.InsertUser(ctx, gen.InsertUserParams{
		Address:      "0x0000000000000000000000000000000000000001",
		RegisteredAt: 1000,
	})
	_ = env.queries.InsertUser(ctx, gen.InsertUserParams{
		Address:      "0x0000000000000000000000000000000000000002",
		RegisteredAt: 2000,
	})

	rr := env.request("GET", "/api/users/count", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]int64
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["count"] != 2 {
		t.Errorf("expected count=2, got %d", resp["count"])
	}
}

func TestGetUser(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	addr := "0x0000000000000000000000000000000000000001"

	// Insert user
	_ = env.queries.InsertUser(ctx, gen.InsertUserParams{
		Address:      addr,
		RegisteredAt: 1000,
	})

	// Initialize balance and create stake position
	_ = env.queries.InitUserBalance(ctx, addr)
	_ = env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		TokenID: 1, Owner: addr, Amount: numericFromInt64(5000),
		LockEndTime: 50, CreatedAt: 1000,
	})

	// Insert agent
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO agents (agent_address, owner_address, is_manager) VALUES ($1, $2, $3)",
		"0x00000000000000000000000000000000000000a1", addr, false,
	)

	rr := env.request("GET", "/api/users/"+addr, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)

	// Verify user field
	user, ok := resp["user"].(map[string]any)
	if !ok {
		t.Fatal("expected user field in response")
	}
	if !strings.Contains(user["address"].(string), "0000001") {
		t.Errorf("unexpected address: %v", user["address"])
	}

	// Verify balance field exists
	if resp["balance"] == nil {
		t.Error("expected balance field in response")
	}

	// Verify agents field
	agents, ok := resp["agents"].([]any)
	if !ok {
		t.Fatal("expected agents array in response")
	}
	if len(agents) != 1 {
		t.Errorf("expected 1 agent, got %d", len(agents))
	}
}

func TestGetUserNotFound(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/users/0x0000000000000000000000000000000000000099", "")
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestGetUserCaseInsensitive(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// Insert lowercase address
	addr := "0x000000000000000000000000000000000000abcd"
	_ = env.queries.InsertUser(ctx, gen.InsertUserParams{
		Address:      addr,
		RegisteredAt: 1000,
	})

	// Query with uppercase (handler normalizes to lowercase via normalizeAddr)
	rr := env.request("GET", "/api/users/0x000000000000000000000000000000000000ABCD", "")
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 for case-insensitive lookup, got %d: %s", rr.Code, rr.Body.String())
	}
}

// ════════════════════════════════════════════════════════════
// Address Lookup — 4 tests
// ════════════════════════════════════════════════════════════

func TestCheckAddressUnknown(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/address/0x0000000000000000000000000000000000000099/check", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["isRegisteredUser"] != false {
		t.Error("expected isRegisteredUser=false")
	}
	if resp["isRegisteredAgent"] != false {
		t.Error("expected isRegisteredAgent=false")
	}
	if resp["isManager"] != false {
		t.Error("expected isManager=false")
	}
}

func TestCheckAddressUser(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	addr := "0x0000000000000000000000000000000000000001"
	_ = env.queries.InsertUser(ctx, gen.InsertUserParams{
		Address:      addr,
		RegisteredAt: 1000,
	})

	rr := env.request("GET", "/api/address/"+addr+"/check", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["isRegisteredUser"] != true {
		t.Error("expected isRegisteredUser=true")
	}
	if resp["isRegisteredAgent"] != false {
		t.Error("expected isRegisteredAgent=false for user-only address")
	}
}

func TestCheckAddressAgent(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	agentAddr := "0x00000000000000000000000000000000000000a1"
	ownerAddr := "0x0000000000000000000000000000000000000001"
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO agents (agent_address, owner_address, is_manager) VALUES ($1, $2, $3)",
		agentAddr, ownerAddr, false,
	)

	rr := env.request("GET", "/api/address/"+agentAddr+"/check", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["isRegisteredAgent"] != true {
		t.Error("expected isRegisteredAgent=true")
	}
	if resp["ownerAddress"] == nil || !strings.Contains(resp["ownerAddress"].(string), "0000001") {
		t.Errorf("expected ownerAddress to contain owner, got %v", resp["ownerAddress"])
	}
	if resp["isManager"] != false {
		t.Error("expected isManager=false")
	}
}

func TestCheckAddressManagerAgent(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	agentAddr := "0x00000000000000000000000000000000000000a2"
	ownerAddr := "0x0000000000000000000000000000000000000002"
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO agents (agent_address, owner_address, is_manager) VALUES ($1, $2, $3)",
		agentAddr, ownerAddr, true,
	)

	rr := env.request("GET", "/api/address/"+agentAddr+"/check", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["isRegisteredAgent"] != true {
		t.Error("expected isRegisteredAgent=true")
	}
	if resp["isManager"] != true {
		t.Error("expected isManager=true for manager agent")
	}
}

// ════════════════════════════════════════════════════════════
// Agents — 7 tests
// ════════════════════════════════════════════════════════════

func TestGetAgentsByOwnerEmpty(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/agents/by-owner/0x0000000000000000000000000000000000000001", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var agents []any
	_ = json.Unmarshal(rr.Body.Bytes(), &agents)
	if len(agents) != 0 {
		t.Errorf("expected 0 agents, got %d", len(agents))
	}
}

func TestGetAgentsByOwnerWithData(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	ownerAddr := "0x0000000000000000000000000000000000000001"
	agent1 := "0x00000000000000000000000000000000000000a1"
	agent2 := "0x00000000000000000000000000000000000000a2"

	_, _ = env.pool.Exec(ctx,
		"INSERT INTO agents (agent_address, owner_address, is_manager) VALUES ($1, $2, $3)",
		agent1, ownerAddr, false,
	)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO agents (agent_address, owner_address, is_manager) VALUES ($1, $2, $3)",
		agent2, ownerAddr, true,
	)

	rr := env.request("GET", "/api/agents/by-owner/"+ownerAddr, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var agents []map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &agents)
	if len(agents) != 2 {
		t.Errorf("expected 2 agents, got %d", len(agents))
	}
}

func TestGetAgentDetail(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	ownerAddr := "0x0000000000000000000000000000000000000001"
	agentAddr := "0x00000000000000000000000000000000000000a1"

	_, _ = env.pool.Exec(ctx,
		"INSERT INTO agents (agent_address, owner_address, is_manager) VALUES ($1, $2, $3)",
		agentAddr, ownerAddr, false,
	)

	// GetAgentDetail uses the /by-owner/{owner}/{agent} route but internally only queries agent
	rr := env.request("GET", "/api/agents/by-owner/"+ownerAddr+"/"+agentAddr, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if !strings.Contains(resp["agent_address"].(string), "a1") {
		t.Errorf("unexpected agent_address: %v", resp["agent_address"])
	}
	if !strings.Contains(resp["owner_address"].(string), "0000001") {
		t.Errorf("unexpected owner_address: %v", resp["owner_address"])
	}
}

func TestLookupAgent(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	ownerAddr := "0x0000000000000000000000000000000000000001"
	agentAddr := "0x00000000000000000000000000000000000000a1"
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO agents (agent_address, owner_address, is_manager) VALUES ($1, $2, $3)",
		agentAddr, ownerAddr, false,
	)

	rr := env.request("GET", "/api/agents/lookup/"+agentAddr, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if !strings.Contains(resp["ownerAddress"], "0000001") {
		t.Errorf("expected ownerAddress to contain owner, got %s", resp["ownerAddress"])
	}
}

func TestLookupAgentNotFound(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/agents/lookup/0x0000000000000000000000000000000000000099", "")
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestBatchAgentInfo(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	agent1 := "0x00000000000000000000000000000000000000a1"
	agent2 := "0x00000000000000000000000000000000000000a2"
	ownerAddr := "0x0000000000000000000000000000000000000001"

	_, _ = env.pool.Exec(ctx,
		"INSERT INTO agents (agent_address, owner_address, is_manager) VALUES ($1, $2, $3)",
		agent1, ownerAddr, false,
	)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO agents (agent_address, owner_address, is_manager) VALUES ($1, $2, $3)",
		agent2, ownerAddr, true,
	)

	body := `{"agents":["` + agent1 + `","` + agent2 + `"],"subnetId":1}`
	rr := env.request("POST", "/api/agents/batch-info", body)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var results []map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &results)
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestBatchAgentInfoTooMany(t *testing.T) {
	env := newTestEnv(t)

	// Build 101 addresses
	agents := make([]string, 101)
	for i := range agents {
		agents[i] = "\"0x0000000000000000000000000000000000000001\""
	}
	body := `{"agents":[` + strings.Join(agents, ",") + `],"subnetId":1}`

	rr := env.request("POST", "/api/agents/batch-info", body)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for >100 agents, got %d: %s", rr.Code, rr.Body.String())
	}
}

// ════════════════════════════════════════════════════════════
// Staking — 9 tests
// ════════════════════════════════════════════════════════════

func TestGetBalanceDefault(t *testing.T) {
	env := newTestEnv(t)
	// Should return zero values when no balance record exists
	rr := env.request("GET", "/api/staking/user/0x0000000000000000000000000000000000000099/balance", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["totalStaked"] != "0" {
		t.Errorf("expected totalStaked=0, got %v", resp["totalStaked"])
	}
	if resp["totalAllocated"] != "0" {
		t.Errorf("expected totalAllocated=0, got %v", resp["totalAllocated"])
	}
	if resp["unallocated"] != "0" {
		t.Errorf("expected unallocated=0, got %v", resp["unallocated"])
	}
}

func TestGetBalanceWithData(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	addr := "0x0000000000000000000000000000000000000001"
	_ = env.queries.InitUserBalance(ctx, addr)
	_ = env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		TokenID: 1, Owner: addr, Amount: numericFromInt64(10001),
		LockEndTime: 50, CreatedAt: 100,
	})
	_ = env.queries.AddUserAllocated(ctx, gen.AddUserAllocatedParams{
		UserAddress:    addr,
		TotalAllocated: numericFromInt64(3001),
	})

	rr := env.request("GET", "/api/staking/user/"+addr+"/balance", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["totalStaked"] != "10001" {
		t.Errorf("expected totalStaked=10001, got %v", resp["totalStaked"])
	}
	if resp["totalAllocated"] != "3001" {
		t.Errorf("expected totalAllocated=3001, got %v", resp["totalAllocated"])
	}
	if resp["unallocated"] != "7000" {
		t.Errorf("expected unallocated=7000, got %v", resp["unallocated"])
	}
}

func TestGetStakePositions(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	addr := "0x0000000000000000000000000000000000000001"
	_ = env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		TokenID: 1, Owner: addr, Amount: numericFromInt64(5000),
		LockEndTime: 50, CreatedAt: 100,
	})
	_ = env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		TokenID: 2, Owner: addr, Amount: numericFromInt64(3000),
		LockEndTime: 100, CreatedAt: 200,
	})

	rr := env.request("GET", "/api/staking/user/"+addr+"/positions", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var positions []map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &positions)
	if len(positions) != 2 {
		t.Fatalf("expected 2 positions, got %d", len(positions))
	}
}

func TestGetAllocationsEmpty(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/staking/user/0x0000000000000000000000000000000000000001/allocations?page=1&limit=10", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var allocs []any
	_ = json.Unmarshal(rr.Body.Bytes(), &allocs)
	if len(allocs) != 0 {
		t.Errorf("expected 0 allocations, got %d", len(allocs))
	}
}

func TestGetAllocationsWithData(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	userAddr := "0x0000000000000000000000000000000000000001"
	agentAddr := "0x00000000000000000000000000000000000000a1"

	// Insert 3 allocation records
	for i := int64(1); i <= 3; i++ {
		_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
			UserAddress:  userAddr,
			AgentAddress: agentAddr,
			SubnetID:     i,
			Amount:       numericFromInt64(1000 * i),
		})
	}

	t.Run("AllAllocations", func(t *testing.T) {
		rr := env.request("GET", "/api/staking/user/"+userAddr+"/allocations?page=1&limit=10", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var allocs []any
		_ = json.Unmarshal(rr.Body.Bytes(), &allocs)
		if len(allocs) != 3 {
			t.Errorf("expected 3 allocations, got %d", len(allocs))
		}
	})

	t.Run("Pagination", func(t *testing.T) {
		rr := env.request("GET", "/api/staking/user/"+userAddr+"/allocations?page=2&limit=2", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var allocs []any
		_ = json.Unmarshal(rr.Body.Bytes(), &allocs)
		if len(allocs) != 1 {
			t.Errorf("expected 1 allocation on page 2, got %d", len(allocs))
		}
	})
}

func TestGetPending(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/staking/user/0x0000000000000000000000000000000000000001/pending", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var pending []any
	_ = json.Unmarshal(rr.Body.Bytes(), &pending)
	if len(pending) != 0 {
		t.Errorf("expected empty pending array, got %d items", len(pending))
	}
}

func TestGetFrozenEmpty(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/staking/user/0x0000000000000000000000000000000000000001/frozen", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var frozen []any
	_ = json.Unmarshal(rr.Body.Bytes(), &frozen)
	if len(frozen) != 0 {
		t.Errorf("expected empty frozen array, got %d items", len(frozen))
	}
}

func TestGetAgentSubnetStake(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	agentAddr := "0x00000000000000000000000000000000000000a1"
	userAddr := "0x0000000000000000000000000000000000000001"

	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		UserAddress:  userAddr,
		AgentAddress: agentAddr,
		SubnetID:     1,
		Amount:       numericFromInt64(5001),
	})

	rr := env.request("GET", "/api/staking/agent/"+agentAddr+"/subnet/1", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["amount"] != "5001" {
		t.Errorf("expected amount=5001, got %s", resp["amount"])
	}
}

func TestGetSubnetTotalStake(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	agent1 := "0x00000000000000000000000000000000000000a1"
	agent2 := "0x00000000000000000000000000000000000000a2"
	user1 := "0x0000000000000000000000000000000000000001"
	user2 := "0x0000000000000000000000000000000000000002"

	// Two allocation records for the same subnet
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		UserAddress:  user1,
		AgentAddress: agent1,
		SubnetID:     1,
		Amount:       numericFromInt64(3001),
	})
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		UserAddress:  user2,
		AgentAddress: agent2,
		SubnetID:     1,
		Amount:       numericFromInt64(7002),
	})

	rr := env.request("GET", "/api/staking/subnet/1/total", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["total"] != "10003" {
		t.Errorf("expected total=10003, got %s", resp["total"])
	}
}

// ════════════════════════════════════════════════════════════
// Subnets — 5 tests
// ════════════════════════════════════════════════════════════

// insertSubnet is a helper that inserts a subnet record
func insertSubnet(t *testing.T, pool *pgxpool.Pool, id int64, owner, name, symbol, status string) {
	t.Helper()
	ctx := context.Background()
	_, err := pool.Exec(ctx,
		`INSERT INTO subnets (subnet_id, owner, name, symbol, subnet_contract, alpha_token, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		id, owner, name, symbol,
		"0x00000000000000000000000000000000000000c1",
		"0x00000000000000000000000000000000000000d1",
		status, 1000,
	)
	if err != nil {
		t.Fatalf("insertSubnet failed: %v", err)
	}
}

func TestListSubnetsEmpty(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/subnets/?page=1&limit=10", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var subnets []any
	_ = json.Unmarshal(rr.Body.Bytes(), &subnets)
	if len(subnets) != 0 {
		t.Errorf("expected 0 subnets, got %d", len(subnets))
	}
}

func TestListSubnetsWithStatusFilter(t *testing.T) {
	env := newTestEnv(t)
	ownerAddr := "0x0000000000000000000000000000000000000001"

	insertSubnet(t, env.pool, 1, ownerAddr, "Subnet1", "S1", "Active")
	insertSubnet(t, env.pool, 2, ownerAddr, "Subnet2", "S2", "Pending")
	insertSubnet(t, env.pool, 3, ownerAddr, "Subnet3", "S3", "Active")

	t.Run("FilterActive", func(t *testing.T) {
		rr := env.request("GET", "/api/subnets/?status=Active&page=1&limit=10", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var subnets []map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &subnets)
		if len(subnets) != 2 {
			t.Errorf("expected 2 Active subnets, got %d", len(subnets))
		}
		for _, s := range subnets {
			if s["status"] != "Active" {
				t.Errorf("expected status=Active, got %v", s["status"])
			}
		}
	})

	t.Run("FilterPending", func(t *testing.T) {
		rr := env.request("GET", "/api/subnets/?status=Pending&page=1&limit=10", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var subnets []map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &subnets)
		if len(subnets) != 1 {
			t.Errorf("expected 1 Pending subnet, got %d", len(subnets))
		}
	})

	t.Run("NoFilter", func(t *testing.T) {
		rr := env.request("GET", "/api/subnets/?page=1&limit=10", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var subnets []map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &subnets)
		if len(subnets) != 3 {
			t.Errorf("expected 3 subnets without filter, got %d", len(subnets))
		}
	})
}

func TestGetSubnet(t *testing.T) {
	env := newTestEnv(t)
	ownerAddr := "0x0000000000000000000000000000000000000001"
	insertSubnet(t, env.pool, 42, ownerAddr, "TestNet", "TN", "Active")

	rr := env.request("GET", "/api/subnets/42", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["subnet_id"] != float64(42) {
		t.Errorf("expected subnet_id=42, got %v", resp["subnet_id"])
	}
	if resp["name"] != "TestNet" {
		t.Errorf("expected name=TestNet, got %v", resp["name"])
	}
	if resp["symbol"] != "TN" {
		t.Errorf("expected symbol=TN, got %v", resp["symbol"])
	}
	if resp["status"] != "Active" {
		t.Errorf("expected status=Active, got %v", resp["status"])
	}
}

func TestGetSubnetNotFound(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/subnets/9999", "")
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestGetSubnetEarnings(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	ownerAddr := "0x0000000000000000000000000000000000000001"
	insertSubnet(t, env.pool, 1, ownerAddr, "Sub1", "S1", "Active")

	// Insert epoch
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO epochs (epoch_id, start_time, daily_emission) VALUES ($1, $2, $3)",
		1, 1000, 15800000,
	)

	// Insert distribution records (keyed by subnet_contract address matching insertSubnet)
	subnetContract := "0x00000000000000000000000000000000000000c1"
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO recipient_awp_distributions (epoch_id, recipient, awp_amount) VALUES ($1, $2, $3)",
		1, subnetContract, 7900000,
	)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO recipient_awp_distributions (epoch_id, recipient, awp_amount) VALUES ($1, $2, $3)",
		2, subnetContract, 3000000,
	)

	rr := env.request("GET", "/api/subnets/1/earnings?page=1&limit=10", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var earnings []map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &earnings)
	if len(earnings) != 2 {
		t.Errorf("expected 2 earnings records, got %d", len(earnings))
	}
}

func TestGetSubnetAgentInfo(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// Insert subnet
	_, _ = env.pool.Exec(ctx,
		`INSERT INTO subnets (subnet_id, owner, name, symbol, subnet_contract, alpha_token, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		1, "0xowner", "Sub1", "S1", "0xsc1", "0xalpha1", "Active", 1000)

	// Insert allocation
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		UserAddress:  "0xuser1",
		AgentAddress: "0xagent1",
		SubnetID:     1,
		Amount:       numericFromInt64(5001),
	})

	rr := env.request("GET", "/api/subnets/1/agents/0xagent1", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["stake"] != "5001" {
		t.Errorf("expected stake=5001, got %v", resp["stake"])
	}
}

func TestGetSubnetAgentInfoNotFound(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/subnets/999/agents/0xnonexistent", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["stake"] != "0" {
		t.Errorf("expected stake=0, got %v", resp["stake"])
	}
}

// ════════════════════════════════════════════════════════════
// Emission — 4 tests
// ════════════════════════════════════════════════════════════

func TestGetCurrentEmissionCacheMiss(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/emission/current", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	// Cache miss should return empty object
	if len(resp) != 0 {
		t.Errorf("expected empty object on cache miss, got %d keys", len(resp))
	}
}

func TestGetCurrentEmissionWithCache(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	data := `{"epoch":"5","dailyEmission":"15000000","totalWeight":"1000"}`
	env.rdb.Set(ctx, "emission_current", data, 0)

	rr := env.request("GET", "/api/emission/current", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["epoch"] != "5" {
		t.Errorf("unexpected epoch: %v", resp["epoch"])
	}
	if resp["dailyEmission"] != "15000000" {
		t.Errorf("unexpected dailyEmission: %v", resp["dailyEmission"])
	}
	if resp["totalWeight"] != "1000" {
		t.Errorf("unexpected totalWeight: %v", resp["totalWeight"])
	}
}

func TestGetEmissionScheduleNoEpochs(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/emission/schedule", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)

	// Should use initial emission
	if resp["currentDailyEmission"] == nil {
		t.Fatal("expected currentDailyEmission in response")
	}

	// Verify projections structure
	projections, ok := resp["projections"].([]any)
	if !ok {
		t.Fatal("expected projections array in response")
	}
	if len(projections) != 3 {
		t.Errorf("expected 3 projections (30/90/365 days), got %d", len(projections))
	}

	// Verify each projection's fields
	for _, p := range projections {
		proj, ok := p.(map[string]any)
		if !ok {
			t.Fatal("expected projection to be object")
		}
		if proj["days"] == nil {
			t.Error("expected days in projection")
		}
		if proj["totalEmission"] == nil {
			t.Error("expected totalEmission in projection")
		}
		if proj["finalDailyRate"] == nil {
			t.Error("expected finalDailyRate in projection")
		}
	}
}

func TestListEpochsWithData(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// Insert 3 epochs
	for i := int64(1); i <= 3; i++ {
		_, _ = env.pool.Exec(ctx,
			"INSERT INTO epochs (epoch_id, start_time, daily_emission) VALUES ($1, $2, $3)",
			i, i*86400, 15800000-i*100000,
		)
	}

	rr := env.request("GET", "/api/emission/epochs?page=1&limit=10", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var epochs []map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &epochs)
	if len(epochs) != 3 {
		t.Errorf("expected 3 epochs, got %d", len(epochs))
	}
	// Should be ordered by epoch_id DESC
	if epochs[0]["epoch_id"] != float64(3) {
		t.Errorf("expected first epoch_id=3 (DESC order), got %v", epochs[0]["epoch_id"])
	}
}

// ════════════════════════════════════════════════════════════
// Tokens — 4 tests
// ════════════════════════════════════════════════════════════

func TestGetAWPInfoCacheMiss(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/tokens/awp", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if len(resp) != 0 {
		t.Errorf("expected empty object on cache miss, got %d keys", len(resp))
	}
}

func TestGetAWPInfoWithCache(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	data := `{"totalSupply":"5000000000","maxSupply":"10000000000"}`
	env.rdb.Set(ctx, "awp_info", data, 0)

	rr := env.request("GET", "/api/tokens/awp", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["totalSupply"] != "5000000000" {
		t.Errorf("unexpected totalSupply: %v", resp["totalSupply"])
	}
	if resp["maxSupply"] != "10000000000" {
		t.Errorf("unexpected maxSupply: %v", resp["maxSupply"])
	}
}

func TestGetAlphaInfo(t *testing.T) {
	env := newTestEnv(t)
	ownerAddr := "0x0000000000000000000000000000000000000001"
	insertSubnet(t, env.pool, 5, ownerAddr, "AlphaNet", "AN", "Active")

	rr := env.request("GET", "/api/tokens/alpha/5", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["subnetId"] != float64(5) {
		t.Errorf("expected subnetId=5, got %v", resp["subnetId"])
	}
	if resp["name"] != "AlphaNet" {
		t.Errorf("expected name=AlphaNet, got %v", resp["name"])
	}
	if resp["symbol"] != "AN" {
		t.Errorf("expected symbol=AN, got %v", resp["symbol"])
	}
	if resp["alphaToken"] == nil {
		t.Error("expected alphaToken in response")
	}
}

func TestGetAlphaPriceCacheMiss(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/tokens/alpha/1/price", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if len(resp) != 0 {
		t.Errorf("expected empty object on cache miss, got %d keys", len(resp))
	}
}

// ════════════════════════════════════════════════════════════
// Governance — 4 tests
// ════════════════════════════════════════════════════════════

func TestListProposalsEmpty(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/governance/proposals?page=1&limit=10", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var proposals []any
	_ = json.Unmarshal(rr.Body.Bytes(), &proposals)
	if len(proposals) != 0 {
		t.Errorf("expected 0 proposals, got %d", len(proposals))
	}
}

func TestListProposalsWithStatusFilter(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// Insert proposals with different statuses
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO proposals (proposal_id, proposer, description, status, votes_for, votes_against) VALUES ($1, $2, $3, $4, $5, $6)",
		"0x0000000000000000000000000000000000000000000000000000000000000001",
		"0x0000000000000000000000000000000000000001",
		"Proposal 1", "Active", 0, 0,
	)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO proposals (proposal_id, proposer, description, status, votes_for, votes_against) VALUES ($1, $2, $3, $4, $5, $6)",
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"0x0000000000000000000000000000000000000001",
		"Proposal 2", "Executed", 0, 0,
	)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO proposals (proposal_id, proposer, description, status, votes_for, votes_against) VALUES ($1, $2, $3, $4, $5, $6)",
		"0x0000000000000000000000000000000000000000000000000000000000000003",
		"0x0000000000000000000000000000000000000002",
		"Proposal 3", "Active", 0, 0,
	)

	t.Run("FilterActive", func(t *testing.T) {
		rr := env.request("GET", "/api/governance/proposals?status=Active&page=1&limit=10", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var proposals []map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &proposals)
		if len(proposals) != 2 {
			t.Errorf("expected 2 Active proposals, got %d", len(proposals))
		}
		for _, p := range proposals {
			if p["status"] != "Active" {
				t.Errorf("expected status=Active, got %v", p["status"])
			}
		}
	})

	t.Run("FilterExecuted", func(t *testing.T) {
		rr := env.request("GET", "/api/governance/proposals?status=Executed&page=1&limit=10", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var proposals []map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &proposals)
		if len(proposals) != 1 {
			t.Errorf("expected 1 Executed proposal, got %d", len(proposals))
		}
	})
}

func TestGetProposal(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	proposalID := "0x0000000000000000000000000000000000000000000000000000000000000042"
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO proposals (proposal_id, proposer, description, status, votes_for, votes_against) VALUES ($1, $2, $3, $4, $5, $6)",
		proposalID,
		"0x0000000000000000000000000000000000000001",
		"Test Proposal", "Active", 100, 50,
	)

	rr := env.request("GET", "/api/governance/proposals/"+proposalID, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["proposal_id"] == nil {
		t.Fatal("expected proposal_id in response")
	}
	if resp["status"] != "Active" {
		t.Errorf("expected status=Active, got %v", resp["status"])
	}
}

func TestGetTreasury(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/api/governance/treasury", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["treasuryAddress"] != "0x1234567890abcdef1234567890abcdef12345678" {
		t.Errorf("unexpected treasury address: %s", resp["treasuryAddress"])
	}
}

// ════════════════════════════════════════════════════════════
// E2E flow tests — 3 tests
// ════════════════════════════════════════════════════════════

func TestE2E_UserRegistrationToStaking(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	userAddr := "0x0000000000000000000000000000000000000001"
	agentAddr := "0x00000000000000000000000000000000000000a1"

	// Step 1: Register user
	err := env.queries.InsertUser(ctx, gen.InsertUserParams{
		Address:      userAddr,
		RegisteredAt: 1000,
	})
	if err != nil {
		t.Fatalf("InsertUser failed: %v", err)
	}

	// Step 2: Initialize balance
	err = env.queries.InitUserBalance(ctx, userAddr)
	if err != nil {
		t.Fatalf("InitUserBalance failed: %v", err)
	}

	// Step 3: Create stake position
	err = env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		TokenID: 1, Owner: userAddr, Amount: numericFromInt64(100001),
		LockEndTime: 50, CreatedAt: 1000,
	})
	if err != nil {
		t.Fatalf("InsertStakePosition failed: %v", err)
	}

	// Step 4: Insert agent
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO agents (agent_address, owner_address, is_manager) VALUES ($1, $2, $3)",
		agentAddr, userAddr, false,
	)

	// Step 5: Insert stake allocations
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		UserAddress:  userAddr,
		AgentAddress: agentAddr,
		SubnetID:     1,
		Amount:       numericFromInt64(30000),
	})
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		UserAddress:  userAddr,
		AgentAddress: agentAddr,
		SubnetID:     2,
		Amount:       numericFromInt64(20000),
	})

	// Update total_allocated
	err = env.queries.AddUserAllocated(ctx, gen.AddUserAllocatedParams{
		UserAddress:    userAddr,
		TotalAllocated: numericFromInt64(50001),
	})
	if err != nil {
		t.Fatalf("AddUserAllocated failed: %v", err)
	}

	// Verify: query balance
	t.Run("VerifyBalance", func(t *testing.T) {
		rr := env.request("GET", "/api/staking/user/"+userAddr+"/balance", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var resp map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["totalStaked"] != "100001" {
			t.Errorf("expected totalStaked=100001, got %v", resp["totalStaked"])
		}
		if resp["totalAllocated"] != "50001" {
			t.Errorf("expected totalAllocated=50001, got %v", resp["totalAllocated"])
		}
		if resp["unallocated"] != "50000" {
			t.Errorf("expected unallocated=50000, got %v", resp["unallocated"])
		}
	})

	// Verify: query allocation list
	t.Run("VerifyAllocations", func(t *testing.T) {
		rr := env.request("GET", "/api/staking/user/"+userAddr+"/allocations?page=1&limit=10", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var allocs []map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &allocs)
		if len(allocs) != 2 {
			t.Errorf("expected 2 allocations, got %d", len(allocs))
		}
	})

	// Verify: user detail includes agent
	t.Run("VerifyUserDetail", func(t *testing.T) {
		rr := env.request("GET", "/api/users/"+userAddr, "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var resp map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		agents, ok := resp["agents"].([]any)
		if !ok || len(agents) != 1 {
			t.Errorf("expected 1 agent in user detail, got %v", resp["agents"])
		}
	})

	// Verify: address check also detects user
	t.Run("VerifyAddressCheck", func(t *testing.T) {
		rr := env.request("GET", "/api/address/"+userAddr+"/check", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}
		var resp map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["isRegisteredUser"] != true {
			t.Error("expected isRegisteredUser=true in E2E flow")
		}
	})
}

func TestE2E_SubnetRegistrationToEmission(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	ownerAddr := "0x0000000000000000000000000000000000000001"

	// Step 1: Register subnet (Pending status)
	insertSubnet(t, env.pool, 1, ownerAddr, "TestSubnet", "TS", "Pending")

	// Step 2: Activate subnet
	_, err := env.pool.Exec(ctx,
		"UPDATE subnets SET status = 'Active', activated_at = $1 WHERE subnet_id = $2",
		2000, 1,
	)
	if err != nil {
		t.Fatalf("activate subnet failed: %v", err)
	}

	// Step 3: Insert epoch
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO epochs (epoch_id, start_time, daily_emission) VALUES ($1, $2, $3)",
		1, 3000, 15800000,
	)

	// Step 4: Insert distribution record (by subnet_contract address, matching the address set in insertSubnet)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO recipient_awp_distributions (epoch_id, recipient, awp_amount) VALUES ($1, $2, $3)",
		1, "0x00000000000000000000000000000000000000c1", 7900000,
	)

	// Verify: subnet detail
	t.Run("VerifySubnet", func(t *testing.T) {
		rr := env.request("GET", "/api/subnets/1", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var resp map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["status"] != "Active" {
			t.Errorf("expected status=Active after activation, got %v", resp["status"])
		}
	})

	// Verify: subnet earnings
	t.Run("VerifyEarnings", func(t *testing.T) {
		rr := env.request("GET", "/api/subnets/1/earnings?page=1&limit=10", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var earnings []map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &earnings)
		if len(earnings) != 1 {
			t.Errorf("expected 1 earnings record, got %d", len(earnings))
		}
	})

	// Verify: Alpha token info
	t.Run("VerifyAlphaInfo", func(t *testing.T) {
		rr := env.request("GET", "/api/tokens/alpha/1", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var resp map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["name"] != "TestSubnet" {
			t.Errorf("expected name=TestSubnet, got %v", resp["name"])
		}
	})

	// Verify: epoch list
	t.Run("VerifyEpochs", func(t *testing.T) {
		rr := env.request("GET", "/api/emission/epochs?page=1&limit=10", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var epochs []map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &epochs)
		if len(epochs) != 1 {
			t.Errorf("expected 1 epoch, got %d", len(epochs))
		}
	})
}

func TestE2E_AddressCaseNormalization(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// Insert all with lowercase
	userAddr := "0x000000000000000000000000000000000000abcd"
	agentAddr := "0x000000000000000000000000000000000000ef01"

	_ = env.queries.InsertUser(ctx, gen.InsertUserParams{
		Address:      userAddr,
		RegisteredAt: 1000,
	})
	_ = env.queries.InitUserBalance(ctx, userAddr)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO agents (agent_address, owner_address, is_manager) VALUES ($1, $2, $3)",
		agentAddr, userAddr, false,
	)

	// Query all routes using uppercase addresses
	upperUser := "0x000000000000000000000000000000000000ABCD"
	upperAgent := "0x000000000000000000000000000000000000EF01"

	t.Run("UserLookup", func(t *testing.T) {
		rr := env.request("GET", "/api/users/"+upperUser, "")
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 for case-insensitive user lookup, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("AddressCheck", func(t *testing.T) {
		rr := env.request("GET", "/api/address/"+upperUser+"/check", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}
		var resp map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["isRegisteredUser"] != true {
			t.Error("expected case-insensitive address check to find user")
		}
	})

	t.Run("AgentLookup", func(t *testing.T) {
		rr := env.request("GET", "/api/agents/lookup/"+upperAgent, "")
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 for case-insensitive agent lookup, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("BalanceLookup", func(t *testing.T) {
		rr := env.request("GET", "/api/staking/user/"+upperUser+"/balance", "")
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 for case-insensitive balance lookup, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("AgentsByOwner", func(t *testing.T) {
		rr := env.request("GET", "/api/agents/by-owner/"+upperUser, "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var agents []any
		_ = json.Unmarshal(rr.Body.Bytes(), &agents)
		if len(agents) != 1 {
			t.Errorf("expected 1 agent via case-insensitive owner lookup, got %d", len(agents))
		}
	})
}
