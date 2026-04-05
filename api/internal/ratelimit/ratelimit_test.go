package ratelimit

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// getTestRedis creates a test Redis client; skips tests if unavailable
func getTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	url := os.Getenv("TEST_REDIS_URL")
	if url == "" {
		url = "redis://localhost:6379/1"
	}
	opt, err := redis.ParseURL(url)
	if err != nil {
		t.Skipf("invalid TEST_REDIS_URL: %v", err)
	}
	rdb := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis unavailable: %v", err)
	}
	t.Cleanup(func() {
		_ = rdb.FlushDB(context.Background())
		_ = rdb.Close()
	})
	return rdb
}

func newTestLimiter(t *testing.T) (*Limiter, *redis.Client) {
	t.Helper()
	rdb := getTestRedis(t)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	return NewLimiter(rdb, logger), rdb
}

// ════════════════════════════════════════════════════════════
// CheckAndIncrement tests
// ════════════════════════════════════════════════════════════

func TestCheckAndIncrement_AllowsUnderLimit(t *testing.T) {
	limiter, rdb := newTestLimiter(t)
	ctx := context.Background()

	// Set custom rate limit: limit=5, window=3600s
	_ = rdb.HSet(ctx, configKey, "test_allow", "5:3600").Err()
	// Clear cache
	limiter.cacheMu.Lock()
	limiter.cache = nil
	limiter.cacheMu.Unlock()

	// All 3 requests should pass
	for i := 0; i < 3; i++ {
		exceeded, err := limiter.CheckAndIncrement(ctx, "test_allow", "10.0.0.1")
		if err != nil {
			t.Fatalf("request %d: unexpected error: %v", i+1, err)
		}
		if exceeded {
			t.Fatalf("request %d: should not be exceeded with limit=5", i+1)
		}
	}
}

func TestCheckAndIncrement_BlocksOverLimit(t *testing.T) {
	limiter, rdb := newTestLimiter(t)
	ctx := context.Background()

	_ = rdb.HSet(ctx, configKey, "test_block", "5:3600").Err()
	limiter.cacheMu.Lock()
	limiter.cache = nil
	limiter.cacheMu.Unlock()

	// First 5 should pass
	for i := 0; i < 5; i++ {
		exceeded, err := limiter.CheckAndIncrement(ctx, "test_block", "10.0.0.2")
		if err != nil {
			t.Fatalf("request %d: unexpected error: %v", i+1, err)
		}
		if exceeded {
			t.Fatalf("request %d: should not be exceeded yet", i+1)
		}
	}

	// 6th should be blocked
	exceeded, err := limiter.CheckAndIncrement(ctx, "test_block", "10.0.0.2")
	if err != nil {
		t.Fatalf("request 6: unexpected error: %v", err)
	}
	if !exceeded {
		t.Fatal("request 6: should be exceeded with limit=5")
	}
}

func TestCheckAndIncrement_ResetAfterTTL(t *testing.T) {
	limiter, rdb := newTestLimiter(t)
	ctx := context.Background()

	// Use 2-second window for testing
	_ = rdb.HSet(ctx, configKey, "test_ttl", "2:2").Err()
	limiter.cacheMu.Lock()
	limiter.cache = nil
	limiter.cacheMu.Unlock()

	// Fill the limit
	for i := 0; i < 2; i++ {
		exceeded, err := limiter.CheckAndIncrement(ctx, "test_ttl", "10.0.0.3")
		if err != nil {
			t.Fatalf("request %d: unexpected error: %v", i+1, err)
		}
		if exceeded {
			t.Fatalf("request %d: should not be exceeded", i+1)
		}
	}

	// Should be blocked
	exceeded, _ := limiter.CheckAndIncrement(ctx, "test_ttl", "10.0.0.3")
	if !exceeded {
		t.Fatal("should be exceeded after filling limit")
	}

	// Wait for TTL to expire
	time.Sleep(3 * time.Second)

	// Should be allowed again
	exceeded, err := limiter.CheckAndIncrement(ctx, "test_ttl", "10.0.0.3")
	if err != nil {
		t.Fatalf("after TTL reset: unexpected error: %v", err)
	}
	if exceeded {
		t.Fatal("after TTL reset: should be allowed again")
	}
}

func TestCheckAndIncrement_CounterNotIncrementedWhenBlocked(t *testing.T) {
	limiter, rdb := newTestLimiter(t)
	ctx := context.Background()

	_ = rdb.HSet(ctx, configKey, "test_decr", "3:3600").Err()
	limiter.cacheMu.Lock()
	limiter.cache = nil
	limiter.cacheMu.Unlock()

	ip := "10.0.0.4"
	key := "rl:test_decr:" + ip

	// Fill the limit
	for i := 0; i < 3; i++ {
		_, _ = limiter.CheckAndIncrement(ctx, "test_decr", ip)
	}

	// Blocked request
	exceeded, _ := limiter.CheckAndIncrement(ctx, "test_decr", ip)
	if !exceeded {
		t.Fatal("should be exceeded")
	}

	// Verify counter is still 3 (DECR rolled back the blocked INCR)
	val, err := rdb.Get(ctx, key).Int64()
	if err != nil {
		t.Fatalf("failed to read counter: %v", err)
	}
	if val != 3 {
		t.Fatalf("expected counter=3 after blocked request (DECR should rollback), got %d", val)
	}
}

// ════════════════════════════════════════════════════════════
// GetConfig tests
// ════════════════════════════════════════════════════════════

func TestGetConfig_DefaultValues(t *testing.T) {
	limiter, _ := newTestLimiter(t)
	ctx := context.Background()

	// Without Redis config, should return compile-time default values
	cfg := limiter.GetConfig(ctx, "relay")
	if cfg.Limit != 100 {
		t.Errorf("expected default relay limit=100, got %d", cfg.Limit)
	}
	if cfg.Window != 1*time.Hour {
		t.Errorf("expected default relay window=1h, got %v", cfg.Window)
	}

	cfg = limiter.GetConfig(ctx, "upload_salts")
	if cfg.Limit != 5 {
		t.Errorf("expected default upload_salts limit=5, got %d", cfg.Limit)
	}
}

func TestGetConfig_CustomValues(t *testing.T) {
	limiter, rdb := newTestLimiter(t)
	ctx := context.Background()

	// Set custom configuration in Redis
	_ = rdb.HSet(ctx, configKey, "relay", "200:7200").Err()
	// Clear cache to force re-read
	limiter.cacheMu.Lock()
	limiter.cache = nil
	limiter.cacheMu.Unlock()

	cfg := limiter.GetConfig(ctx, "relay")
	if cfg.Limit != 200 {
		t.Errorf("expected custom relay limit=200, got %d", cfg.Limit)
	}
	if cfg.Window != 7200*time.Second {
		t.Errorf("expected custom relay window=7200s, got %v", cfg.Window)
	}
}

func TestGetConfig_CachesTTL(t *testing.T) {
	limiter, rdb := newTestLimiter(t)
	ctx := context.Background()

	// First read (populates cache)
	_ = rdb.HSet(ctx, configKey, "relay", "300:1800").Err()
	limiter.cacheMu.Lock()
	limiter.cache = nil
	limiter.cacheMu.Unlock()

	cfg1 := limiter.GetConfig(ctx, "relay")
	if cfg1.Limit != 300 {
		t.Fatalf("expected limit=300, got %d", cfg1.Limit)
	}

	// Update Redis (but cache is still valid within 10s TTL)
	_ = rdb.HSet(ctx, configKey, "relay", "999:60").Err()

	// Within cache TTL, should still return old value
	cfg2 := limiter.GetConfig(ctx, "relay")
	if cfg2.Limit != 300 {
		t.Errorf("expected cached limit=300 within TTL, got %d", cfg2.Limit)
	}
	if cfg2.Window != 1800*time.Second {
		t.Errorf("expected cached window=1800s within TTL, got %v", cfg2.Window)
	}
}

// ════════════════════════════════════════════════════════════
// parseConfig tests (package-internal access)
// ════════════════════════════════════════════════════════════

func TestParseConfig_Valid(t *testing.T) {
	cfg, ok := parseConfig("100:3600")
	if !ok {
		t.Fatal("expected parseConfig to succeed")
	}
	if cfg.Limit != 100 {
		t.Errorf("expected limit=100, got %d", cfg.Limit)
	}
	if cfg.Window != 3600*time.Second {
		t.Errorf("expected window=3600s, got %v", cfg.Window)
	}
}

func TestParseConfig_Invalid(t *testing.T) {
	cases := []string{
		"",
		"abc",
		"100",
		"abc:3600",
		"100:abc",
		"-1:3600",
		"100:-1",
	}
	for _, c := range cases {
		_, ok := parseConfig(c)
		if ok {
			t.Errorf("parseConfig(%q) should fail", c)
		}
	}
}
