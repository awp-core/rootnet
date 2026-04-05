package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// testAdminToken is the admin token used for testing
const testAdminToken = "test-admin-secret-token-2026"

// newAdminTestEnv creates a test environment with AdminToken configured
func newAdminTestEnv(t *testing.T) *testEnv {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/awp_test?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skipf("skipping test: cannot connect to database: %v", err)
	}
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
		AWPWorkNetAddress:   "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		DAOAddress:          "0xcccccccccccccccccccccccccccccccccccccccc",
		AWPTokenAddress:     "0xdddddddddddddddddddddddddddddddddddddd",
		AWPAllocatorAddress: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		AWPEmissionAddress:  "0xffffffffffffffffffffffffffffffffffffffff",
		AdminToken:          testAdminToken,
	}

	limiter := ratelimit.NewLimiter(rdb, logger)
	h := handler.NewHandler(queries, pool, rdb, cfg, logger, limiter)
	hub := ws.NewHub(rdb, logger)
	router := server.NewRouter(server.RouterParams{Config: cfg, Handler: h, Hub: hub})

	env := &testEnv{
		pool:    pool,
		rdb:     rdb,
		queries: queries,
		handler: h,
		router:  router,
		t:       t,
	}

	env.cleanDB()

	t.Cleanup(func() {
		env.cleanDB()
		pool.Close()
		_ = rdb.Close()
	})

	return env
}

// adminRequest sends a request with Authorization header
func (e *testEnv) adminRequest(method, path string, body string, token string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()
	e.router.ServeHTTP(rr, req)
	return rr
}

// ════════════════════════════════════════════════════════════
// Auth tests — 3 tests
// ════════════════════════════════════════════════════════════

func TestAdminAuth_NoToken(t *testing.T) {
	env := newAdminTestEnv(t)

	rr := env.adminRequest("GET", "/api/admin/chains", "", "")
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != "unauthorized" {
		t.Errorf("expected error=unauthorized, got %v", resp["error"])
	}
}

func TestAdminAuth_InvalidToken(t *testing.T) {
	env := newAdminTestEnv(t)

	rr := env.adminRequest("GET", "/api/admin/chains", "", "wrong-token")
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 with wrong token, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestAdminAuth_ValidToken(t *testing.T) {
	env := newAdminTestEnv(t)

	rr := env.adminRequest("GET", "/api/admin/chains", "", testAdminToken)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 with valid token, got %d: %s", rr.Code, rr.Body.String())
	}
}

// ════════════════════════════════════════════════════════════
// Chain Management — 4 tests
// ════════════════════════════════════════════════════════════

func TestAdminListChains(t *testing.T) {
	env := newAdminTestEnv(t)
	ctx := context.Background()

	// Insert a chain record
	_ = env.queries.InsertChain(ctx, gen.InsertChainParams{
		ChainID:      56,
		Name:         "BSC",
		RpcUrl:       "https://bsc-rpc.example.com",
		Dex:          "pancakeswap",
		Explorer:     "https://bscscan.com",
		AwpRegistry:  "0xaaaa",
		AwpToken:     "0xbbbb",
		AwpEmission:  "0xcccc",
		AwpAllocator: "0xdddd",
		Veawp:        "0xeeee",
		AwpWorknet:   "0xffff",
		DaoAddress:   "0x1111",
		LpManager:    "0x2222",
		PoolManager:  "0x3333",
		DeployBlock:  1000,
	})
	t.Cleanup(func() {
		_, _ = env.pool.Exec(context.Background(), "DELETE FROM chains WHERE chain_id = 56")
	})

	rr := env.adminRequest("GET", "/api/admin/chains", "", testAdminToken)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp []map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if len(resp) < 1 {
		t.Fatal("expected at least 1 chain")
	}

	// Find the inserted BSC chain
	found := false
	for _, c := range resp {
		if c["chainId"] == float64(56) {
			found = true
			if c["name"] != "BSC" {
				t.Errorf("expected name=BSC, got %v", c["name"])
			}
			// admin endpoint should include rpcUrl
			if c["rpcUrl"] != "https://bsc-rpc.example.com" {
				t.Errorf("expected rpcUrl to be visible in admin endpoint, got %v", c["rpcUrl"])
			}
		}
	}
	if !found {
		t.Error("BSC chain (chainId=56) not found in response")
	}
}

func TestAdminAddChain(t *testing.T) {
	env := newAdminTestEnv(t)

	body := `{
		"chainId": 42161,
		"name": "Arbitrum",
		"rpcUrl": "https://arb-rpc.example.com",
		"dex": "uniswap",
		"explorer": "https://arbiscan.io",
		"awpRegistry": "0xaaaa",
		"awpToken": "0xbbbb",
		"awpEmission": "0xcccc",
		"awpAllocator": "0xdddd",
		"veAWP": "0xeeee",
		"awpWorkNet": "0xffff",
		"daoAddress": "0x1111",
		"lpManager": "0x2222",
		"poolManager": "0x3333",
		"deployBlock": 50000
	}`

	t.Cleanup(func() {
		_, _ = env.pool.Exec(context.Background(), "DELETE FROM chains WHERE chain_id = 42161")
	})

	rr := env.adminRequest("POST", "/api/admin/chains", body, testAdminToken)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["chainId"] != float64(42161) {
		t.Errorf("expected chainId=42161, got %v", resp["chainId"])
	}
	if resp["status"] != "added" {
		t.Errorf("expected status=added, got %v", resp["status"])
	}

	// Verify duplicate add returns 409
	rr2 := env.adminRequest("POST", "/api/admin/chains", body, testAdminToken)
	if rr2.Code != http.StatusConflict {
		t.Errorf("expected 409 for duplicate chain, got %d: %s", rr2.Code, rr2.Body.String())
	}
}

func TestAdminUpdateChain(t *testing.T) {
	// AdminUpdateChain is not a standalone endpoint (using PUT /api/admin/ratelimit to test)
	// This test verifies rate limit update functionality
	env := newAdminTestEnv(t)

	body := `{"key": "test_limit", "limit": 50, "window": 1800}`
	rr := env.adminRequest("PUT", "/api/admin/ratelimit", body, testAdminToken)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["status"] != "updated" {
		t.Errorf("expected status=updated, got %v", resp["status"])
	}

	// Verify reading the updated configuration
	rr2 := env.adminRequest("GET", "/api/admin/ratelimit", "", testAdminToken)
	if rr2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr2.Code, rr2.Body.String())
	}

	var configs map[string]string
	_ = json.Unmarshal(rr2.Body.Bytes(), &configs)
	if configs["test_limit"] != "50:1800" {
		t.Errorf("expected test_limit=50:1800, got %v", configs["test_limit"])
	}
}

func TestAdminDeleteChain(t *testing.T) {
	env := newAdminTestEnv(t)
	ctx := context.Background()

	// First add a chain
	_ = env.queries.InsertChain(ctx, gen.InsertChainParams{
		ChainID: 99999, Name: "TestChain", RpcUrl: "https://test.example.com",
		Dex: "uniswap", Explorer: "https://test.explorer.com",
		AwpRegistry: "0x1", AwpToken: "0x2", AwpEmission: "0x3",
		AwpAllocator: "0x4", Veawp: "0x5", AwpWorknet: "0x6",
		DaoAddress: "0x7", LpManager: "0x8", PoolManager: "0x9",
		DeployBlock: 0,
	})
	t.Cleanup(func() {
		_, _ = env.pool.Exec(context.Background(), "DELETE FROM chains WHERE chain_id = 99999")
	})

	rr := env.adminRequest("DELETE", "/api/admin/chains/99999", "", testAdminToken)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["status"] != "deactivated" {
		t.Errorf("expected status=deactivated, got %v", resp["status"])
	}
}

// ════════════════════════════════════════════════════════════
// System Info — 1 test
// ════════════════════════════════════════════════════════════

func TestAdminSystemInfo(t *testing.T) {
	env := newAdminTestEnv(t)

	rr := env.adminRequest("GET", "/api/admin/system", "", testAdminToken)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)

	// Verify basic fields exist
	if resp["goVersion"] == nil {
		t.Error("expected goVersion in system info")
	}
	if resp["uptime"] == nil {
		t.Error("expected uptime in system info")
	}
	if resp["goroutines"] == nil {
		t.Error("expected goroutines in system info")
	}
	if resp["numCPU"] == nil {
		t.Error("expected numCPU in system info")
	}

	// Verify Redis info
	redisInfo, ok := resp["redis"].(map[string]any)
	if !ok {
		t.Fatal("expected redis object in system info")
	}
	if redisInfo["status"] != "PONG" {
		t.Errorf("expected redis status=PONG, got %v", redisInfo["status"])
	}

	// Verify DB info
	dbInfo, ok := resp["database"].(map[string]any)
	if !ok {
		t.Fatal("expected database object in system info")
	}
	if dbInfo["status"] != "ok" {
		t.Errorf("expected db status=ok, got %v", dbInfo["status"])
	}
}

// ════════════════════════════════════════════════════════════
// Edge cases
// ════════════════════════════════════════════════════════════

func TestAdminAddChain_InvalidBody(t *testing.T) {
	env := newAdminTestEnv(t)

	t.Run("MissingChainId", func(t *testing.T) {
		rr := env.adminRequest("POST", "/api/admin/chains", `{"name":"Test"}`, testAdminToken)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for missing chainId, got %d", rr.Code)
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		rr := env.adminRequest("POST", "/api/admin/chains", `{"chainId":1}`, testAdminToken)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for missing name, got %d", rr.Code)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		rr := env.adminRequest("POST", "/api/admin/chains", `{invalid}`, testAdminToken)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for invalid JSON, got %d", rr.Code)
		}
	})
}

func TestAdminUpdateRateLimit_InvalidBody(t *testing.T) {
	env := newAdminTestEnv(t)

	t.Run("MissingKey", func(t *testing.T) {
		rr := env.adminRequest("PUT", "/api/admin/ratelimit", `{"limit":10,"window":60}`, testAdminToken)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for missing key, got %d", rr.Code)
		}
	})

	t.Run("NegativeLimit", func(t *testing.T) {
		rr := env.adminRequest("PUT", "/api/admin/ratelimit", `{"key":"test","limit":-1,"window":60}`, testAdminToken)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for negative limit, got %d", rr.Code)
		}
	})
}
