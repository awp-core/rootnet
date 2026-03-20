// Package ratelimit provides a Redis-backed, hot-updatable rate limiter.
//
// Limits are stored in Redis as a JSON hash at key "ratelimit:config".
// Each field is a limiter name (e.g. "relay", "upload_salts") with value "limit:window_seconds".
// Updating a limit is instant: redis-cli HSET ratelimit:config relay "200:3600"
//
// If a config key is missing from Redis, the compiled default is used.
package ratelimit

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config defines a rate limit rule: max requests per window per IP.
type Config struct {
	Limit  int           // max requests per window
	Window time.Duration // window duration
}

// Limiter is a Redis-backed, hot-reloadable rate limiter.
type Limiter struct {
	rdb      *redis.Client
	logger   *slog.Logger
	defaults map[string]Config // compiled defaults (fallback)

	// In-process config cache (reduces 1 Redis RTT per rate-limited request)
	cacheMu   sync.RWMutex
	cache     map[string]Config
	cacheTime time.Time
}

const configKey = "ratelimit:config"

// Package-level Lua scripts (allocated once, SHA1 cached by go-redis)
var (
	luaCheckAndIncr = redis.NewScript(`
		local count = redis.call('INCR', KEYS[1])
		if count == 1 then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end
		if count > tonumber(ARGV[2]) then
			redis.call('DECR', KEYS[1])
			return 1
		end
		return 0
	`)
	luaIncr = redis.NewScript(`
		local count = redis.call('INCR', KEYS[1])
		if count == 1 then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end
		return count
	`)
)

// NewLimiter creates a Limiter with default configurations.
// Defaults are used when the Redis config key is absent.
func NewLimiter(rdb *redis.Client, logger *slog.Logger) *Limiter {
	return &Limiter{
		rdb:    rdb,
		logger: logger,
		defaults: map[string]Config{
			"relay":            {Limit: 100, Window: 1 * time.Hour},
			"upload_salts":     {Limit: 5, Window: 1 * time.Hour},
			"compute_salt":     {Limit: 20, Window: 1 * time.Hour},
			"batch_agent_info": {Limit: 30, Window: 1 * time.Hour},
			"ws_connect":       {Limit: 10, Window: 0}, // 0 = concurrent count, not time-windowed
		},
	}
}

// GetConfig reads the current config for a limiter name.
// Uses a 10-second in-process cache to avoid a Redis RTT on every rate-limited request.
func (l *Limiter) GetConfig(ctx context.Context, name string) Config {
	// Check in-process cache first (10s TTL)
	l.cacheMu.RLock()
	if l.cache != nil && time.Since(l.cacheTime) < 10*time.Second {
		if cfg, ok := l.cache[name]; ok {
			l.cacheMu.RUnlock()
			return cfg
		}
	}
	l.cacheMu.RUnlock()

	// Cache miss or expired — read all configs from Redis in one call
	l.cacheMu.Lock()
	defer l.cacheMu.Unlock()
	// Double-check after acquiring write lock
	if l.cache != nil && time.Since(l.cacheTime) < 10*time.Second {
		if cfg, ok := l.cache[name]; ok {
			return cfg
		}
	}
	newCache := make(map[string]Config)
	for k, v := range l.defaults {
		newCache[k] = v
	}
	all, err := l.rdb.HGetAll(ctx, configKey).Result()
	if err == nil {
		for k, v := range all {
			if cfg, ok := parseConfig(v); ok {
				newCache[k] = cfg
			}
		}
	}
	l.cache = newCache
	l.cacheTime = time.Now()

	if cfg, ok := newCache[name]; ok {
		return cfg
	}
	return Config{Limit: 100, Window: time.Hour} // ultimate fallback
}

// CheckAndIncrement atomically checks the rate limit and increments the counter in a single
// Lua script. Returns (exceeded bool, err). This prevents TOCTOU races where concurrent
// requests all see count=0 before any increment.
func (l *Limiter) CheckAndIncrement(ctx context.Context, name string, ip string) (bool, error) {
	cfg := l.GetConfig(ctx, name)
	key := fmt.Sprintf("rl:%s:%s", name, ip)

	result, err := luaCheckAndIncr.Run(ctx, l.rdb, []string{key}, int(cfg.Window.Seconds()), cfg.Limit).Int64()
	if err != nil {
		// Fail-closed: treat Redis errors as "exceeded" to prevent abuse during outages
		l.logger.Error("rate limit Redis error, blocking request (fail-closed)", "name", name, "error", err)
		return true, err
	}
	return result == 1, nil
}


// FormatError returns a user-friendly rate limit exceeded message.
func (l *Limiter) FormatError(ctx context.Context, name string) string {
	cfg := l.GetConfig(ctx, name)
	return fmt.Sprintf("rate limit exceeded: max %d requests per %.0fs", cfg.Limit, cfg.Window.Seconds())
}

// GetClientIP extracts the real client IP from the request.
func GetClientIP(r *http.Request) string {
	ip := r.RemoteAddr
	if host, _, err := net.SplitHostPort(ip); err == nil {
		ip = host
	}
	return ip
}

// parseConfig parses "limit:window_seconds" string.
func parseConfig(val string) (Config, bool) {
	parts := strings.SplitN(val, ":", 2)
	if len(parts) != 2 {
		return Config{}, false
	}
	limit, err := strconv.Atoi(parts[0])
	if err != nil || limit < 0 {
		return Config{}, false
	}
	windowSec, err := strconv.Atoi(parts[1])
	if err != nil || windowSec < 0 {
		return Config{}, false
	}
	return Config{Limit: limit, Window: time.Duration(windowSec) * time.Second}, true
}
