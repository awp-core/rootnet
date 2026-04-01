package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/cortexia/rootnet/api/internal/ratelimit"
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
		dbURL = "postgres://postgres:postgres@localhost:5432/awp_test?sslmode=disable"
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
		ChainID:             31337,
		TreasuryAddress:     "0x1234567890abcdef1234567890abcdef12345678",
		AWPRegistryAddress:  "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		SubnetNFTAddress:    "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		DAOAddress:          "0xcccccccccccccccccccccccccccccccccccccccc",
		AWPTokenAddress:     "0xdddddddddddddddddddddddddddddddddddddd",
		StakingVaultAddress: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		AWPEmissionAddress:  "0xffffffffffffffffffffffffffffffffffffffff",
	}

	limiter := ratelimit.NewLimiter(rdb, logger)
	h := handler.NewHandler(queries, rdb, cfg, logger, limiter)
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
		user_balances, epochs,
		subnets, proposals, users, sync_states, indexed_blocks, vanity_salts`)
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
	if resp["awpRegistry"] != "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Errorf("unexpected awpRegistry: %s", resp["awpRegistry"])
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

	// Insert 3 users
	for i := range []int{0, 1, 2} {
		addr := strings.Replace("0x0000000000000000000000000000000000000001", "1", string(rune('1'+i)), 1)
		_ = env.queries.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{
			ChainID: 31337,
			Address: addr,
			BoundTo: "",
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
	_ = env.queries.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{
		ChainID: 31337,
		Address: "0x0000000000000000000000000000000000000001",
		BoundTo: "",
	})
	_ = env.queries.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{
		ChainID: 31337,
		Address: "0x0000000000000000000000000000000000000002",
		BoundTo: "",
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
	_ = env.queries.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{
		ChainID: 31337,
		Address: addr,
		BoundTo: "",
	})

	// Initialize balance and create stake position
	_ = env.queries.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: 31337, UserAddress: addr})
	_ = env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		ChainID: 31337, TokenID: 1, Owner: addr, Amount: numericFromInt64(5000),
		LockEndTime: 50, CreatedAt: 1000,
	})

	// Insert agent
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		"0x00000000000000000000000000000000000000a1", addr,
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
	_ = env.queries.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{
		ChainID: 31337,
		Address: addr,
		BoundTo: "",
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
	if resp["isRegistered"] != false {
		t.Error("expected isRegistered=false")
	}
	if resp["boundTo"] != "" {
		t.Errorf("expected empty boundTo, got %v", resp["boundTo"])
	}
	if resp["recipient"] != "" {
		t.Errorf("expected empty recipient, got %v", resp["recipient"])
	}
}

func TestCheckAddressUnbound(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	addr := "0x0000000000000000000000000000000000000001"
	_ = env.queries.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{
		ChainID: 31337,
		Address: addr,
		BoundTo: "",
	})

	rr := env.request("GET", "/api/address/"+addr+"/check", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	// Address exists but has no binding and no recipient, so isRegistered=false
	if resp["isRegistered"] != false {
		t.Error("expected isRegistered=false for address with empty binding and recipient")
	}
}

func TestCheckAddressBound(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	agentAddr := "0x00000000000000000000000000000000000000a1"
	ownerAddr := "0x0000000000000000000000000000000000000001"
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		agentAddr, ownerAddr,
	)

	rr := env.request("GET", "/api/address/"+agentAddr+"/check", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["isRegistered"] != true {
		t.Error("expected isRegistered=true for bound address")
	}
	if !strings.Contains(resp["boundTo"].(string), "0000001") {
		t.Errorf("expected boundTo to contain owner, got %v", resp["boundTo"])
	}
}

func TestCheckAddressWithRecipient(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	userAddr := "0x00000000000000000000000000000000000000a2"
	recipientAddr := "0x0000000000000000000000000000000000000002"
	_ = env.queries.UpsertUserRecipient(ctx, gen.UpsertUserRecipientParams{
		ChainID: 31337,
		Address:   userAddr,
		Recipient: recipientAddr,
	})

	rr := env.request("GET", "/api/address/"+userAddr+"/check", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["isRegistered"] != true {
		t.Error("expected isRegistered=true for address with recipient")
	}
	if !strings.Contains(resp["recipient"].(string), "0000002") {
		t.Errorf("expected recipient to contain address, got %v", resp["recipient"])
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
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		agent1, ownerAddr,
	)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		agent2, ownerAddr,
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
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		agentAddr, ownerAddr,
	)

	// GetAgentDetail uses the /by-owner/{owner}/{agent} route — returns user record
	rr := env.request("GET", "/api/agents/by-owner/"+ownerAddr+"/"+agentAddr, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if !strings.Contains(resp["address"].(string), "a1") {
		t.Errorf("unexpected address: %v", resp["address"])
	}
	if !strings.Contains(resp["bound_to"].(string), "0000001") {
		t.Errorf("unexpected bound_to: %v", resp["bound_to"])
	}
}

func TestLookupAgent(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	ownerAddr := "0x0000000000000000000000000000000000000001"
	agentAddr := "0x00000000000000000000000000000000000000a1"
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		agentAddr, ownerAddr,
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
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		agent1, ownerAddr,
	)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		agent2, ownerAddr,
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
	_ = env.queries.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: 31337, UserAddress: addr})
	_ = env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		ChainID: 31337, TokenID: 1, Owner: addr, Amount: numericFromInt64(10001),
		LockEndTime: 50, CreatedAt: 100,
	})
	_ = env.queries.AddUserAllocated(ctx, gen.AddUserAllocatedParams{
		ChainID: 31337, UserAddress:    addr,
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
		ChainID: 31337, TokenID: 1, Owner: addr, Amount: numericFromInt64(5000),
		LockEndTime: 50, CreatedAt: 100,
	})
	_ = env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		ChainID: 31337, TokenID: 2, Owner: addr, Amount: numericFromInt64(3000),
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
			ChainID: 31337, UserAddress:  userAddr,
			AgentAddress: agentAddr,
			SubnetID:     numericFromInt64(i),
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
		ChainID:      31337,
		UserAddress:  userAddr,
		AgentAddress: agentAddr,
		SubnetID:     numericFromInt64(1),
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
		ChainID:      31337,
		UserAddress:  user1,
		AgentAddress: agent1,
		SubnetID:     numericFromInt64(1),
		Amount:       numericFromInt64(3001),
	})
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID:      31337,
		UserAddress:  user2,
		AgentAddress: agent2,
		SubnetID:     numericFromInt64(1),
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
		`INSERT INTO subnets (chain_id, subnet_id, owner, name, symbol, subnet_contract, alpha_token, status, created_at)
		 VALUES (31337, $1, $2, $3, $4, $5, $6, $7, $8)`,
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
		"INSERT INTO epochs (chain_id, epoch_id, start_time, daily_emission) VALUES (31337, $1, $2, $3)",
		1, 1000, 15800000,
	)

	// Insert distribution records (keyed by subnet_contract address matching insertSubnet)
	subnetContract := "0x00000000000000000000000000000000000000c1"
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO recipient_awp_distributions (chain_id, epoch_id, recipient, awp_amount) VALUES (31337, $1, $2, $3)",
		1, subnetContract, 7900000,
	)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO recipient_awp_distributions (chain_id, epoch_id, recipient, awp_amount) VALUES (31337, $1, $2, $3)",
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
		`INSERT INTO subnets (chain_id, subnet_id, owner, name, symbol, subnet_contract, alpha_token, status, created_at)
		 VALUES (31337, $1, $2, $3, $4, $5, $6, $7, $8)`,
		1, "0xowner", "Sub1", "S1", "0xsc1", "0xalpha1", "Active", 1000)

	// Insert allocation
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID:      31337,
		UserAddress:  "0xuser1",
		AgentAddress: "0x1111111111111111111111111111111111111111",
		SubnetID:     numericFromInt64(1),
		Amount:       numericFromInt64(5001),
	})

	rr := env.request("GET", "/api/subnets/1/agents/0x1111111111111111111111111111111111111111", "")
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
	rr := env.request("GET", "/api/subnets/999/agents/0x2222222222222222222222222222222222222222", "")
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
	// When keeper is not running, Redis has no emission_current cache → 503
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 (no keeper cache), got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] == nil {
		t.Errorf("expected error message in response")
	}
}

func TestGetCurrentEmissionWithCache(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	data := `{"epoch":"5","dailyEmission":"15000000","totalWeight":"1000"}`
	env.rdb.Set(ctx, "emission_current:31337", data, 0)

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
			"INSERT INTO epochs (chain_id, epoch_id, start_time, daily_emission) VALUES (31337, $1, $2, $3)",
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
	env.rdb.Set(ctx, "awp_info:31337", data, 0)

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
		"INSERT INTO proposals (chain_id, proposal_id, proposer, description, status, votes_for, votes_against) VALUES (31337, $1, $2, $3, $4, $5, $6)",
		"0x0000000000000000000000000000000000000000000000000000000000000001",
		"0x0000000000000000000000000000000000000001",
		"Proposal 1", "Active", 0, 0,
	)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO proposals (chain_id, proposal_id, proposer, description, status, votes_for, votes_against) VALUES (31337, $1, $2, $3, $4, $5, $6)",
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"0x0000000000000000000000000000000000000001",
		"Proposal 2", "Executed", 0, 0,
	)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO proposals (chain_id, proposal_id, proposer, description, status, votes_for, votes_against) VALUES (31337, $1, $2, $3, $4, $5, $6)",
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
		"INSERT INTO proposals (chain_id, proposal_id, proposer, description, status, votes_for, votes_against) VALUES (31337, $1, $2, $3, $4, $5, $6)",
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

	// Step 1: Register user with recipient (so address check shows isRegistered=true)
	err := env.queries.UpsertUserRecipient(ctx, gen.UpsertUserRecipientParams{
		ChainID: 31337,
		Address:   userAddr,
		Recipient: userAddr, // self-recipient
	})
	if err != nil {
		t.Fatalf("UpsertUserRecipient failed: %v", err)
	}

	// Step 2: Initialize balance
	err = env.queries.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: 31337, UserAddress: userAddr})
	if err != nil {
		t.Fatalf("InitUserBalance failed: %v", err)
	}

	// Step 3: Create stake position
	err = env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		ChainID: 31337, TokenID: 1, Owner: userAddr, Amount: numericFromInt64(100001),
		LockEndTime: 50, CreatedAt: 1000,
	})
	if err != nil {
		t.Fatalf("InsertStakePosition failed: %v", err)
	}

	// Step 4: Insert agent
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		agentAddr, userAddr,
	)

	// Step 5: Insert stake allocations
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID: 31337, UserAddress: userAddr,
		AgentAddress: agentAddr,
		SubnetID:     numericFromInt64(1),
		Amount:       numericFromInt64(30000),
	})
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID: 31337, UserAddress: userAddr,
		AgentAddress: agentAddr,
		SubnetID:     numericFromInt64(2),
		Amount:       numericFromInt64(20000),
	})

	// Update total_allocated
	err = env.queries.AddUserAllocated(ctx, gen.AddUserAllocatedParams{
		ChainID: 31337, UserAddress:    userAddr,
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
		if resp["isRegistered"] != true {
			t.Error("expected isRegistered=true in E2E flow")
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
		"INSERT INTO epochs (chain_id, epoch_id, start_time, daily_emission) VALUES (31337, $1, $2, $3)",
		1, 3000, 15800000,
	)

	// Step 4: Insert distribution record (by subnet_contract address, matching the address set in insertSubnet)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO recipient_awp_distributions (chain_id, epoch_id, recipient, awp_amount) VALUES (31337, $1, $2, $3)",
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
	recipientAddr := "0x0000000000000000000000000000000000001234"

	// User has a recipient set (makes isRegistered=true)
	_ = env.queries.UpsertUserRecipient(ctx, gen.UpsertUserRecipientParams{
		ChainID: 31337,
		Address:   userAddr,
		Recipient: recipientAddr,
	})
	_ = env.queries.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: 31337, UserAddress: userAddr})
	// Agent is bound to user
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		agentAddr, userAddr,
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
		if resp["isRegistered"] != true {
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

// ════════════════════════════════════════════════════════════
// JSON-RPC 2.0 — /v2
// ════════════════════════════════════════════════════════════

func TestJSONRPC_Discover(t *testing.T) {
	env := newTestEnv(t)

	body := `{"jsonrpc":"2.0","method":"rpc.discover","id":1}`
	rr := env.request("POST", "/v2", body)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["jsonrpc"] != "2.0" {
		t.Error("expected jsonrpc 2.0")
	}
	result, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatal("expected result object")
	}
	methods, ok := result["methods"].([]any)
	if !ok || len(methods) == 0 {
		t.Fatal("expected non-empty methods array")
	}

	// 验证至少包含关键方法
	methodNames := make(map[string]bool)
	for _, m := range methods {
		if mm, ok := m.(map[string]any); ok {
			methodNames[mm["name"].(string)] = true
		}
	}
	for _, expected := range []string{
		"registry.get", "staking.getBalance", "subnets.list",
		"emission.getCurrent", "tokens.getAWP", "governance.listProposals",
	} {
		if !methodNames[expected] {
			t.Errorf("missing method in discover: %s", expected)
		}
	}
}

func TestJSONRPC_HealthCheck(t *testing.T) {
	env := newTestEnv(t)

	body := `{"jsonrpc":"2.0","method":"health.check","id":2}`
	rr := env.request("POST", "/v2", body)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	result := resp["result"].(map[string]any)
	if result["status"] != "ok" {
		t.Errorf("expected status ok, got %v", result["status"])
	}
}

func TestJSONRPC_RegistryGet(t *testing.T) {
	env := newTestEnv(t)

	body := `{"jsonrpc":"2.0","method":"registry.get","id":3}`
	rr := env.request("POST", "/v2", body)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	result := resp["result"].(map[string]any)
	if result["awpRegistry"] != "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Errorf("unexpected awpRegistry: %v", result["awpRegistry"])
	}
	if result["chainId"].(float64) != 31337 {
		t.Errorf("unexpected chainId: %v", result["chainId"])
	}
}

func TestJSONRPC_MethodNotFound(t *testing.T) {
	env := newTestEnv(t)

	body := `{"jsonrpc":"2.0","method":"nonexistent.method","id":4}`
	rr := env.request("POST", "/v2", body)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	errObj := resp["error"].(map[string]any)
	if errObj["code"].(float64) != -32601 {
		t.Errorf("expected method not found error (-32601), got %v", errObj["code"])
	}
}

func TestJSONRPC_InvalidParams(t *testing.T) {
	env := newTestEnv(t)

	// 无效地址
	body := `{"jsonrpc":"2.0","method":"users.get","params":{"address":"invalid"},"id":5}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	errObj := resp["error"].(map[string]any)
	if errObj["code"].(float64) != -32602 {
		t.Errorf("expected invalid params error (-32602), got %v", errObj["code"])
	}
}

func TestJSONRPC_BatchRequest(t *testing.T) {
	env := newTestEnv(t)

	body := `[
		{"jsonrpc":"2.0","method":"health.check","id":1},
		{"jsonrpc":"2.0","method":"registry.get","id":2},
		{"jsonrpc":"2.0","method":"users.count","id":3}
	]`
	rr := env.request("POST", "/v2", body)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var responses []map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &responses)
	if len(responses) != 3 {
		t.Fatalf("expected 3 responses in batch, got %d", len(responses))
	}

	// 每个响应都应有 result（无 error）
	for i, resp := range responses {
		if resp["error"] != nil {
			t.Errorf("batch response %d has error: %v", i, resp["error"])
		}
		if resp["result"] == nil {
			t.Errorf("batch response %d missing result", i)
		}
	}
}

func TestJSONRPC_UsersGet(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	addr := "0x0000000000000000000000000000000000001234"
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO users (chain_id, address, registered_at) VALUES (31337, $1, 100)", addr)

	body := `{"jsonrpc":"2.0","method":"users.get","params":{"address":"0x0000000000000000000000000000000000001234"},"id":6}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)
	user := result["user"].(map[string]any)
	if user["address"] != addr {
		t.Errorf("expected address %s, got %v", addr, user["address"])
	}
}

func TestJSONRPC_SubnetsGet(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	_, _ = env.pool.Exec(ctx,
		`INSERT INTO subnets (chain_id, subnet_id, owner, name, symbol, subnet_contract, alpha_token, status, created_at)
		 VALUES (31337, 1, '0xowner', 'TestSubnet', 'TS', '0xsc', '0xat', 'Active', 100)`)

	body := `{"jsonrpc":"2.0","method":"subnets.get","params":{"subnetId":"1"},"id":7}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)
	if result["name"] != "TestSubnet" {
		t.Errorf("expected name TestSubnet, got %v", result["name"])
	}
}

// ════════════════════════════════════════════════════════════
// JSON-RPC 2.0 — Protocol-level tests
// ════════════════════════════════════════════════════════════

func TestJSONRPC_GETReturnsDiscover(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("GET", "/v2", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	result := resp["result"].(map[string]any)
	methods := result["methods"].([]any)
	if len(methods) == 0 {
		t.Fatal("GET /v2 should return rpc.discover with methods")
	}
}

func TestJSONRPC_InvalidJSON(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("POST", "/v2", `{invalid json!!!`)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	errObj := resp["error"].(map[string]any)
	if errObj["code"].(float64) != -32700 {
		t.Errorf("expected parse error (-32700), got %v", errObj["code"])
	}
}

func TestJSONRPC_MissingVersion(t *testing.T) {
	env := newTestEnv(t)

	// 没有 jsonrpc 字段
	body := `{"method":"health.check","id":1}`
	rr := env.request("POST", "/v2", body)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	errObj := resp["error"].(map[string]any)
	if errObj["code"].(float64) != -32600 {
		t.Errorf("expected invalid request error (-32600), got %v", errObj["code"])
	}

	// jsonrpc 字段值错误
	body2 := `{"jsonrpc":"1.0","method":"health.check","id":2}`
	rr2 := env.request("POST", "/v2", body2)
	var resp2 map[string]any
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp2)
	errObj2 := resp2["error"].(map[string]any)
	if errObj2["code"].(float64) != -32600 {
		t.Errorf("expected invalid request error (-32600) for wrong version, got %v", errObj2["code"])
	}
}

func TestJSONRPC_EmptyBatch(t *testing.T) {
	env := newTestEnv(t)
	rr := env.request("POST", "/v2", `[]`)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	errObj := resp["error"].(map[string]any)
	if errObj["code"].(float64) != -32600 {
		t.Errorf("expected invalid request error (-32600), got %v", errObj["code"])
	}
	if msg, ok := errObj["message"].(string); !ok || msg != "empty batch" {
		t.Errorf("expected 'empty batch' message, got %v", errObj["message"])
	}
}

func TestJSONRPC_BatchSizeExceeded(t *testing.T) {
	env := newTestEnv(t)

	// 构建 21 个请求的批量数组
	batch := "["
	for i := 1; i <= 21; i++ {
		if i > 1 {
			batch += ","
		}
		batch += fmt.Sprintf(`{"jsonrpc":"2.0","method":"health.check","id":%d}`, i)
	}
	batch += "]"

	rr := env.request("POST", "/v2", batch)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	errObj := resp["error"].(map[string]any)
	if errObj["code"].(float64) != -32600 {
		t.Errorf("expected invalid request error (-32600), got %v", errObj["code"])
	}
}

func TestJSONRPC_NullParams(t *testing.T) {
	env := newTestEnv(t)

	// null params — 对无参方法应正常工作
	body := `{"jsonrpc":"2.0","method":"health.check","params":null,"id":1}`
	rr := env.request("POST", "/v2", body)
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error with null params: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)
	if result["status"] != "ok" {
		t.Errorf("expected status ok, got %v", result["status"])
	}

	// 缺少 params 字段 — 同样应正常
	body2 := `{"jsonrpc":"2.0","method":"health.check","id":2}`
	rr2 := env.request("POST", "/v2", body2)
	var resp2 map[string]any
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp2)
	if resp2["error"] != nil {
		t.Fatalf("unexpected error with missing params: %v", resp2["error"])
	}
}

func TestJSONRPC_IDPreservation(t *testing.T) {
	env := newTestEnv(t)

	// 字符串 ID
	body1 := `{"jsonrpc":"2.0","method":"health.check","id":"my-string-id"}`
	rr1 := env.request("POST", "/v2", body1)
	var resp1 map[string]any
	_ = json.Unmarshal(rr1.Body.Bytes(), &resp1)
	if resp1["id"] != "my-string-id" {
		t.Errorf("expected string id 'my-string-id', got %v", resp1["id"])
	}

	// 数字 ID
	body2 := `{"jsonrpc":"2.0","method":"health.check","id":42}`
	rr2 := env.request("POST", "/v2", body2)
	var resp2 map[string]any
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp2)
	if resp2["id"].(float64) != 42 {
		t.Errorf("expected numeric id 42, got %v", resp2["id"])
	}

	// null ID
	body3 := `{"jsonrpc":"2.0","method":"health.check","id":null}`
	rr3 := env.request("POST", "/v2", body3)
	var resp3 map[string]any
	_ = json.Unmarshal(rr3.Body.Bytes(), &resp3)
	if resp3["id"] != nil {
		t.Errorf("expected null id, got %v", resp3["id"])
	}
}

// ════════════════════════════════════════════════════════════
// JSON-RPC 2.0 — Staking method tests
// ════════════════════════════════════════════════════════════

func TestJSONRPC_StakingGetBalance(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	addr := "0x0000000000000000000000000000000000aabb01"

	// 插入 stake position（总质押 = 10001）— 使用 numericFromInt64 保证 Exp=0
	if err := env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		ChainID: 31337, TokenID: 1, Owner: addr, Amount: numericFromInt64(7001),
		LockEndTime: 99999, CreatedAt: 100,
	}); err != nil {
		t.Fatalf("InsertStakePosition 1 failed: %v", err)
	}
	if err := env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		ChainID: 31337, TokenID: 2, Owner: addr, Amount: numericFromInt64(3000),
		LockEndTime: 99999, CreatedAt: 200,
	}); err != nil {
		t.Fatalf("InsertStakePosition 2 failed: %v", err)
	}

	// 插入 user_balances（已分配 = 4001）
	if err := env.queries.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: 31337, UserAddress: addr}); err != nil {
		t.Fatalf("InitUserBalance failed: %v", err)
	}
	if err := env.queries.AddUserAllocated(ctx, gen.AddUserAllocatedParams{
		ChainID: 31337, UserAddress: addr, TotalAllocated: numericFromInt64(4001),
	}); err != nil {
		t.Fatalf("AddUserAllocated failed: %v", err)
	}

	body := fmt.Sprintf(`{"jsonrpc":"2.0","method":"staking.getBalance","params":{"address":"%s"},"id":1}`, addr)
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)
	if result["totalStaked"] != "10001" {
		t.Errorf("expected totalStaked=10001, got %v", result["totalStaked"])
	}
	if result["totalAllocated"] != "4001" {
		t.Errorf("expected totalAllocated=4001, got %v", result["totalAllocated"])
	}
	if result["unallocated"] != "6000" {
		t.Errorf("expected unallocated=6000, got %v", result["unallocated"])
	}
}

func TestJSONRPC_StakingGetAllocations(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	addr := "0x0000000000000000000000000000000000aabb02"

	// 插入 3 条 allocation（用 raw SQL 避免 sqlc 参数问题）
	for i := int64(1); i <= 3; i++ {
		agent := fmt.Sprintf("0x00000000000000000000000000000000000000a%d", i)
		if _, err := env.pool.Exec(ctx,
			`INSERT INTO stake_allocations (chain_id, user_address, agent_address, subnet_id, amount) VALUES ($1, $2, $3, $4, $5)`,
			int64(31337), addr, agent, i*100, i*1000,
		); err != nil {
			t.Fatalf("insert allocation %d: %v", i, err)
		}
	}

	body := fmt.Sprintf(`{"jsonrpc":"2.0","method":"staking.getAllocations","params":{"address":"%s","page":1,"limit":2},"id":1}`, addr)
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result, ok := resp["result"].([]any)
	if !ok {
		t.Fatal("expected result to be array")
	}
	if len(result) != 2 {
		t.Errorf("expected 2 allocations (limit=2), got %d; body=%s", len(result), rr.Body.String())
	}

	// 第二页应有 1 条
	body2 := fmt.Sprintf(`{"jsonrpc":"2.0","method":"staking.getAllocations","params":{"address":"%s","page":2,"limit":2},"id":2}`, addr)
	rr2 := env.request("POST", "/v2", body2)
	var resp2 map[string]any
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp2)
	result2 := resp2["result"].([]any)
	if len(result2) != 1 {
		t.Errorf("expected 1 allocation on page 2, got %d; body=%s", len(result2), rr2.Body.String())
	}
}

func TestJSONRPC_StakingGetAgentSubnetStake(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	agent := "0x00000000000000000000000000000000000000a1"
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID: 31337, UserAddress: "0x0000000000000000000000000000000000user01",
		AgentAddress: agent, SubnetID: numericFromInt64(500), Amount: numericFromInt64(7777),
	})

	body := `{"jsonrpc":"2.0","method":"staking.getAgentSubnetStake","params":{"agent":"0x00000000000000000000000000000000000000a1","subnetId":"500"},"id":1}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)
	if result["amount"] != "7777" {
		t.Errorf("expected amount=7777, got %v", result["amount"])
	}
}

func TestJSONRPC_StakingGetAgentSubnets(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	agent := "0x00000000000000000000000000000000000000b1"
	// 插入同一 agent 在不同子网的 allocation
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID: 31337, UserAddress: "0x0000000000000000000000000000000000user01",
		AgentAddress: agent, SubnetID: numericFromInt64(100), Amount: numericFromInt64(5000),
	})
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID: 31337, UserAddress: "0x0000000000000000000000000000000000user01",
		AgentAddress: agent, SubnetID: numericFromInt64(200), Amount: numericFromInt64(3000),
	})

	body := `{"jsonrpc":"2.0","method":"staking.getAgentSubnets","params":{"agent":"0x00000000000000000000000000000000000000b1"},"id":1}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result, ok := resp["result"].([]any)
	if !ok {
		t.Fatal("expected result to be array")
	}
	if len(result) != 2 {
		t.Errorf("expected 2 subnets, got %d", len(result))
	}
}

func TestJSONRPC_StakingGetSubnetTotalStake(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// 插入多个 agent 在同一子网的 allocation — 使用不能被10整除的数避免 pgx NUMERIC Exp>0
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID: 31337, UserAddress: "0x0000000000000000000000000000000000user01",
		AgentAddress: "0x00000000000000000000000000000000000000c1",
		SubnetID: numericFromInt64(999), Amount: numericFromInt64(5001),
	})
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID: 31337, UserAddress: "0x0000000000000000000000000000000000user02",
		AgentAddress: "0x00000000000000000000000000000000000000c2",
		SubnetID: numericFromInt64(999), Amount: numericFromInt64(3002),
	})

	body := `{"jsonrpc":"2.0","method":"staking.getSubnetTotalStake","params":{"subnetId":"999"},"id":1}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)
	if result["total"] != "8003" {
		t.Errorf("expected total=8003, got %v", result["total"])
	}
}

// ════════════════════════════════════════════════════════════
// JSON-RPC 2.0 — Subnet method tests
// ════════════════════════════════════════════════════════════

func TestJSONRPC_SubnetsList(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// 插入不同状态的子网
	_, _ = env.pool.Exec(ctx,
		`INSERT INTO subnets (chain_id, subnet_id, owner, name, symbol, subnet_contract, alpha_token, status, created_at)
		 VALUES (31337, 1001, '0xowner1', 'Sub1', 'S1', '0xsc1', '0xat1', 'Active', 100)`)
	_, _ = env.pool.Exec(ctx,
		`INSERT INTO subnets (chain_id, subnet_id, owner, name, symbol, subnet_contract, alpha_token, status, created_at)
		 VALUES (31337, 1002, '0xowner2', 'Sub2', 'S2', '0xsc2', '0xat2', 'Pending', 200)`)
	_, _ = env.pool.Exec(ctx,
		`INSERT INTO subnets (chain_id, subnet_id, owner, name, symbol, subnet_contract, alpha_token, status, created_at)
		 VALUES (31337, 1003, '0xowner3', 'Sub3', 'S3', '0xsc3', '0xat3', 'Active', 300)`)

	// 测试不带过滤
	body := `{"jsonrpc":"2.0","method":"subnets.list","params":{},"id":1}`
	rr := env.request("POST", "/v2", body)
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].([]any)
	if len(result) != 3 {
		t.Errorf("expected 3 subnets, got %d", len(result))
	}

	// 测试按状态过滤
	body2 := `{"jsonrpc":"2.0","method":"subnets.list","params":{"status":"Active"},"id":2}`
	rr2 := env.request("POST", "/v2", body2)
	var resp2 map[string]any
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp2)
	if resp2["error"] != nil {
		t.Fatalf("unexpected error: %v", resp2["error"])
	}
	result2 := resp2["result"].([]any)
	if len(result2) != 2 {
		t.Errorf("expected 2 Active subnets, got %d", len(result2))
	}
}

func TestJSONRPC_SubnetsGetSkills(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	_, _ = env.pool.Exec(ctx,
		`INSERT INTO subnets (chain_id, subnet_id, owner, name, symbol, subnet_contract, alpha_token, status, created_at, skills_uri)
		 VALUES (31337, 2001, '0xowner', 'SkillSub', 'SK', '0xsc', '0xat', 'Active', 100, 'ipfs://QmSkillsHash')`)

	body := `{"jsonrpc":"2.0","method":"subnets.getSkills","params":{"subnetId":"2001"},"id":1}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)
	if result["skillsURI"] != "ipfs://QmSkillsHash" {
		t.Errorf("expected skillsURI=ipfs://QmSkillsHash, got %v", result["skillsURI"])
	}
	// subnetId 返回为 pgtype.Numeric，JSON 序列化为 number
	if fmt.Sprintf("%v", result["subnetId"]) != "2001" {
		t.Errorf("expected subnetId=2001, got %v", result["subnetId"])
	}
}

func TestJSONRPC_SubnetsGetEarnings(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	_, _ = env.pool.Exec(ctx,
		`INSERT INTO subnets (chain_id, subnet_id, owner, name, symbol, subnet_contract, alpha_token, status, created_at)
		 VALUES (31337, 3001, '0xowner', 'EarnSub', 'ES', '0xsc', '0xat', 'Active', 100)`)

	// 验证空结果（无 epoch 数据）
	body := `{"jsonrpc":"2.0","method":"subnets.getEarnings","params":{"subnetId":"3001"},"id":1}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].([]any)
	if len(result) != 0 {
		t.Errorf("expected 0 earnings, got %d", len(result))
	}
}

// ════════════════════════════════════════════════════════════
// JSON-RPC 2.0 — Other method tests
// ════════════════════════════════════════════════════════════

func TestJSONRPC_AddressCheck(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	addr := "0x0000000000000000000000000000000000aabb03"
	_, _ = env.pool.Exec(ctx,
		`INSERT INTO users (chain_id, address, bound_to, recipient, registered_at)
		 VALUES (31337, $1, '0x0000000000000000000000000000000000owner1', '0x0000000000000000000000000000000000recip1', 100)`, addr)

	body := fmt.Sprintf(`{"jsonrpc":"2.0","method":"address.check","params":{"address":"%s"},"id":1}`, addr)
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)
	if result["isRegistered"] != true {
		t.Errorf("expected isRegistered=true, got %v", result["isRegistered"])
	}
	if result["boundTo"] != "0x0000000000000000000000000000000000owner1" {
		t.Errorf("expected boundTo, got %v", result["boundTo"])
	}
	if result["recipient"] != "0x0000000000000000000000000000000000recip1" {
		t.Errorf("expected recipient, got %v", result["recipient"])
	}

	// 未注册地址应返回 isRegistered=false
	body2 := `{"jsonrpc":"2.0","method":"address.check","params":{"address":"0x0000000000000000000000000000000000ff0000"},"id":2}`
	rr2 := env.request("POST", "/v2", body2)
	var resp2 map[string]any
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp2)
	result2 := resp2["result"].(map[string]any)
	if result2["isRegistered"] != false {
		t.Errorf("expected isRegistered=false for unknown address, got %v", result2["isRegistered"])
	}
}

func TestJSONRPC_AgentsGetByOwner(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	owner := "0x0000000000000000000000000000000000000b01"
	agent1 := "0x0000000000000000000000000000000000000b02"
	agent2 := "0x0000000000000000000000000000000000000b03"

	// 插入 owner
	_, _ = env.pool.Exec(ctx,
		`INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, '')`, owner)
	// 插入绑定到 owner 的 agent
	_, _ = env.pool.Exec(ctx,
		`INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2)`, agent1, owner)
	_, _ = env.pool.Exec(ctx,
		`INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2)`, agent2, owner)

	body := fmt.Sprintf(`{"jsonrpc":"2.0","method":"agents.getByOwner","params":{"owner":"%s"},"id":1}`, owner)
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].([]any)
	if len(result) != 2 {
		t.Errorf("expected 2 agents, got %d", len(result))
	}
}

func TestJSONRPC_AgentsLookup(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	agent := "0x0000000000000000000000000000000000000c01"
	owner := "0x0000000000000000000000000000000000000c02"

	_, _ = env.pool.Exec(ctx,
		`INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2)`, agent, owner)

	body := fmt.Sprintf(`{"jsonrpc":"2.0","method":"agents.lookup","params":{"agent":"%s"},"id":1}`, agent)
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)
	if result["ownerAddress"] != owner {
		t.Errorf("expected ownerAddress=%s, got %v", owner, result["ownerAddress"])
	}
}

func TestJSONRPC_ChainsList(t *testing.T) {
	env := newTestEnv(t)

	body := `{"jsonrpc":"2.0","method":"chains.list","id":1}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result, ok := resp["result"].([]any)
	if !ok || len(result) == 0 {
		t.Fatal("expected non-empty chains list")
	}
	chain := result[0].(map[string]any)
	if chain["chainId"].(float64) != 31337 {
		t.Errorf("expected chainId=31337, got %v", chain["chainId"])
	}
}

func TestJSONRPC_GovernanceListProposals(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	_ = env.queries.InsertProposal(ctx, gen.InsertProposalParams{
		ChainID: 31337, ProposalID: "1001",
		Proposer: "0x0000000000000000000000000000000000prop01",
		Description: pgtype.Text{String: "Test proposal 1", Valid: true},
		Status: "Active",
	})
	_ = env.queries.InsertProposal(ctx, gen.InsertProposalParams{
		ChainID: 31337, ProposalID: "1002",
		Proposer: "0x0000000000000000000000000000000000prop02",
		Description: pgtype.Text{String: "Test proposal 2", Valid: true},
		Status: "Executed",
	})

	// 不过滤
	body := `{"jsonrpc":"2.0","method":"governance.listProposals","params":{},"id":1}`
	rr := env.request("POST", "/v2", body)
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].([]any)
	if len(result) != 2 {
		t.Errorf("expected 2 proposals, got %d", len(result))
	}

	// 按状态过滤
	body2 := `{"jsonrpc":"2.0","method":"governance.listProposals","params":{"status":"Active"},"id":2}`
	rr2 := env.request("POST", "/v2", body2)
	var resp2 map[string]any
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp2)
	result2 := resp2["result"].([]any)
	if len(result2) != 1 {
		t.Errorf("expected 1 Active proposal, got %d", len(result2))
	}
}

func TestJSONRPC_GovernanceGetProposal(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	_ = env.queries.InsertProposal(ctx, gen.InsertProposalParams{
		ChainID: 31337, ProposalID: "2001",
		Proposer: "0x0000000000000000000000000000000000prop01",
		Description: pgtype.Text{String: "Single proposal", Valid: true},
		Status: "Succeeded",
	})

	body := `{"jsonrpc":"2.0","method":"governance.getProposal","params":{"proposalId":"2001"},"id":1}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)
	if result["description"] != "Single proposal" {
		t.Errorf("expected description='Single proposal', got %v", result["description"])
	}
}

func TestJSONRPC_GovernanceTreasury(t *testing.T) {
	env := newTestEnv(t)

	body := `{"jsonrpc":"2.0","method":"governance.getTreasury","id":1}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)
	if result["treasuryAddress"] != "0x1234567890abcdef1234567890abcdef12345678" {
		t.Errorf("unexpected treasury address: %v", result["treasuryAddress"])
	}
}

func TestJSONRPC_EmissionSchedule(t *testing.T) {
	env := newTestEnv(t)

	body := `{"jsonrpc":"2.0","method":"emission.getSchedule","id":1}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)

	// 应有 currentDailyEmission 字段
	if result["currentDailyEmission"] == nil {
		t.Error("expected currentDailyEmission field")
	}

	// 应有 projections 数组，含 3 个条目（30/90/365天）
	projections, ok := result["projections"].([]any)
	if !ok || len(projections) != 3 {
		t.Fatalf("expected 3 projections, got %v", result["projections"])
	}

	// 验证每个 projection 结构
	for _, p := range projections {
		proj := p.(map[string]any)
		if proj["days"] == nil || proj["totalEmission"] == nil || proj["finalDailyRate"] == nil {
			t.Errorf("projection missing fields: %v", proj)
		}
	}
}

func TestJSONRPC_EmissionListEpochs(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// 空结果
	body := `{"jsonrpc":"2.0","method":"emission.listEpochs","params":{},"id":1}`
	rr := env.request("POST", "/v2", body)
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].([]any)
	if len(result) != 0 {
		t.Errorf("expected 0 epochs, got %d", len(result))
	}

	// 插入 epoch 数据
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO epochs (chain_id, epoch_id, start_time, daily_emission) VALUES (31337, $1, $2, $3)",
		0, 86400, 15800000)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO epochs (chain_id, epoch_id, start_time, daily_emission) VALUES (31337, $1, $2, $3)",
		1, 172800, 15750000)

	body2 := `{"jsonrpc":"2.0","method":"emission.listEpochs","params":{},"id":2}`
	rr2 := env.request("POST", "/v2", body2)
	var resp2 map[string]any
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp2)
	result2 := resp2["result"].([]any)
	if len(result2) != 2 {
		t.Errorf("expected 2 epochs, got %d", len(result2))
	}
}

// ════════════════════════════════════════════════════════════
// JSON-RPC 2.0 — Error case tests
// ════════════════════════════════════════════════════════════

func TestJSONRPC_UsersGetNotFound(t *testing.T) {
	env := newTestEnv(t)

	body := `{"jsonrpc":"2.0","method":"users.get","params":{"address":"0x0000000000000000000000000000000000099999"},"id":1}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	errObj := resp["error"].(map[string]any)
	if errObj["code"].(float64) != -32001 {
		t.Errorf("expected -32001 (not found), got %v", errObj["code"])
	}
	if errObj["message"] != "user not found" {
		t.Errorf("expected 'user not found', got %v", errObj["message"])
	}
}

func TestJSONRPC_SubnetsGetNotFound(t *testing.T) {
	env := newTestEnv(t)

	body := `{"jsonrpc":"2.0","method":"subnets.get","params":{"subnetId":"999999"},"id":1}`
	rr := env.request("POST", "/v2", body)

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	errObj := resp["error"].(map[string]any)
	if errObj["code"].(float64) != -32001 {
		t.Errorf("expected -32001 (not found), got %v", errObj["code"])
	}
	if errObj["message"] != "subnet not found" {
		t.Errorf("expected 'subnet not found', got %v", errObj["message"])
	}
}

func TestJSONRPC_InvalidSubnetID(t *testing.T) {
	env := newTestEnv(t)

	tests := []struct {
		name string
		body string
	}{
		{"zero", `{"jsonrpc":"2.0","method":"subnets.get","params":{"subnetId":"0"},"id":1}`},
		{"negative", `{"jsonrpc":"2.0","method":"subnets.get","params":{"subnetId":"-1"},"id":2}`},
		{"non-numeric", `{"jsonrpc":"2.0","method":"subnets.get","params":{"subnetId":"abc"},"id":3}`},
		{"empty", `{"jsonrpc":"2.0","method":"subnets.get","params":{"subnetId":""},"id":4}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := env.request("POST", "/v2", tc.body)
			var resp map[string]any
			_ = json.Unmarshal(rr.Body.Bytes(), &resp)
			errObj, ok := resp["error"].(map[string]any)
			if !ok {
				t.Fatal("expected error response")
			}
			if errObj["code"].(float64) != -32602 {
				t.Errorf("expected -32602, got %v", errObj["code"])
			}
		})
	}
}

func TestJSONRPC_MissingRequiredAddress(t *testing.T) {
	env := newTestEnv(t)

	tests := []struct {
		name string
		body string
	}{
		{"empty-params", `{"jsonrpc":"2.0","method":"users.get","params":{},"id":1}`},
		{"missing-address", `{"jsonrpc":"2.0","method":"staking.getBalance","params":{},"id":2}`},
		{"invalid-address", `{"jsonrpc":"2.0","method":"users.get","params":{"address":"not-an-address"},"id":3}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := env.request("POST", "/v2", tc.body)
			var resp map[string]any
			_ = json.Unmarshal(rr.Body.Bytes(), &resp)
			errObj, ok := resp["error"].(map[string]any)
			if !ok {
				t.Fatal("expected error response")
			}
			if errObj["code"].(float64) != -32602 {
				t.Errorf("expected -32602, got %v", errObj["code"])
			}
		})
	}
}

// ════════════════════════════════════════════════════════════
// JSON-RPC 2.0 — E2E flow test
// ════════════════════════════════════════════════════════════

func TestJSONRPC_E2E_StakingFlow(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	userAddr := "0x0000000000000000000000000000000000e2e001"
	agentAddr := "0x0000000000000000000000000000000000e2e002"

	// Step 1: 创建用户
	_, _ = env.pool.Exec(ctx,
		`INSERT INTO users (chain_id, address, bound_to, registered_at) VALUES (31337, $1, '', 100)`, userAddr)

	// 验证用户存在
	body := fmt.Sprintf(`{"jsonrpc":"2.0","method":"users.get","params":{"address":"%s"},"id":1}`, userAddr)
	rr := env.request("POST", "/v2", body)
	var resp1 map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp1)
	if resp1["error"] != nil {
		t.Fatalf("Step 1 failed: %v", resp1["error"])
	}

	// Step 2: 插入 stake positions（总质押 = 15001）
	if err := env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		ChainID: 31337, TokenID: 101, Owner: userAddr, Amount: numericFromInt64(10001),
		LockEndTime: 99999, CreatedAt: 100,
	}); err != nil {
		t.Fatalf("InsertStakePosition 1: %v", err)
	}
	if err := env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		ChainID: 31337, TokenID: 102, Owner: userAddr, Amount: numericFromInt64(5000),
		LockEndTime: 99999, CreatedAt: 200,
	}); err != nil {
		t.Fatalf("InsertStakePosition 2: %v", err)
	}

	// Step 3: 插入 allocations
	if err := env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID: 31337, UserAddress: userAddr, AgentAddress: agentAddr,
		SubnetID: numericFromInt64(500), Amount: numericFromInt64(8001), UpdatedBlock: 100,
	}); err != nil {
		t.Fatalf("UpsertStakeAllocation: %v", err)
	}
	if err := env.queries.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: 31337, UserAddress: userAddr}); err != nil {
		t.Fatalf("InitUserBalance: %v", err)
	}
	if err := env.queries.AddUserAllocated(ctx, gen.AddUserAllocatedParams{
		ChainID: 31337, UserAddress: userAddr, TotalAllocated: numericFromInt64(8001), UpdatedBlock: 100,
	}); err != nil {
		t.Fatalf("AddUserAllocated: %v", err)
	}

	// Step 4: 查询 balance
	body4 := fmt.Sprintf(`{"jsonrpc":"2.0","method":"staking.getBalance","params":{"address":"%s"},"id":4}`, userAddr)
	rr4 := env.request("POST", "/v2", body4)
	var resp4 map[string]any
	_ = json.Unmarshal(rr4.Body.Bytes(), &resp4)
	if resp4["error"] != nil {
		t.Fatalf("Step 4 failed: %v", resp4["error"])
	}
	balance := resp4["result"].(map[string]any)
	if balance["totalStaked"] != "15001" {
		t.Errorf("expected totalStaked=15001, got %v", balance["totalStaked"])
	}
	if balance["totalAllocated"] != "8001" {
		t.Errorf("expected totalAllocated=8001, got %v", balance["totalAllocated"])
	}
	if balance["unallocated"] != "7000" {
		t.Errorf("expected unallocated=7000, got %v", balance["unallocated"])
	}

	// Step 5: 查询 allocations
	body5 := fmt.Sprintf(`{"jsonrpc":"2.0","method":"staking.getAllocations","params":{"address":"%s"},"id":5}`, userAddr)
	rr5 := env.request("POST", "/v2", body5)
	var resp5 map[string]any
	_ = json.Unmarshal(rr5.Body.Bytes(), &resp5)
	allocs := resp5["result"].([]any)
	if len(allocs) != 1 {
		t.Fatalf("expected 1 allocation, got %d", len(allocs))
	}

	// Step 6: 查询 agent subnet stake
	body6 := fmt.Sprintf(`{"jsonrpc":"2.0","method":"staking.getAgentSubnetStake","params":{"agent":"%s","subnetId":"500"},"id":6}`, agentAddr)
	rr6 := env.request("POST", "/v2", body6)
	var resp6 map[string]any
	_ = json.Unmarshal(rr6.Body.Bytes(), &resp6)
	if resp6["error"] != nil {
		t.Fatalf("Step 6 failed: %v", resp6["error"])
	}
	stakeResult := resp6["result"].(map[string]any)
	if stakeResult["amount"] != "8001" {
		t.Errorf("expected agent subnet stake=8001, got %v", stakeResult["amount"])
	}
}
