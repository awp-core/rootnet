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
}

const configKey = "ratelimit:config"

// NewLimiter creates a Limiter with default configurations.
// Defaults are used when the Redis config key is absent.
func NewLimiter(rdb *redis.Client, logger *slog.Logger) *Limiter {
	return &Limiter{
		rdb:    rdb,
		logger: logger,
		defaults: map[string]Config{
			"relay":         {Limit: 100, Window: 1 * time.Hour},
			"upload_salts":  {Limit: 5, Window: 1 * time.Hour},
			"compute_salt":  {Limit: 20, Window: 1 * time.Hour},
			"ws_connect":    {Limit: 10, Window: 0}, // 0 = concurrent count, not time-windowed
		},
	}
}

// GetConfig reads the current config for a limiter name.
// Returns the Redis-stored value if present, otherwise the compiled default.
func (l *Limiter) GetConfig(ctx context.Context, name string) Config {
	val, err := l.rdb.HGet(ctx, configKey, name).Result()
	if err == nil {
		if cfg, ok := parseConfig(val); ok {
			return cfg
		}
		l.logger.Warn("invalid rate limit config in Redis, using default", "name", name, "value", val)
	}
	if def, ok := l.defaults[name]; ok {
		return def
	}
	return Config{Limit: 100, Window: time.Hour} // ultimate fallback
}

// CheckAndIncrement atomically checks the rate limit and increments the counter in a single
// Lua script. Returns (exceeded bool, err). This prevents TOCTOU races where concurrent
// requests all see count=0 before any increment.
func (l *Limiter) CheckAndIncrement(ctx context.Context, name string, ip string) (bool, error) {
	cfg := l.GetConfig(ctx, name)
	key := fmt.Sprintf("rl:%s:%s", name, ip)

	luaScript := redis.NewScript(`
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
	result, err := luaScript.Run(ctx, l.rdb, []string{key}, int(cfg.Window.Seconds()), cfg.Limit).Int64()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// CheckIP checks if IP has exceeded the rate limit (read-only, for backward compat).
func (l *Limiter) CheckIP(ctx context.Context, name string, ip string) (bool, error) {
	cfg := l.GetConfig(ctx, name)
	key := fmt.Sprintf("rl:%s:%s", name, ip)
	count, err := l.rdb.Get(ctx, key).Int64()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return count >= int64(cfg.Limit), nil
}

// RecordSuccess increments the rate limit counter for a successful operation.
// Uses background context so the increment is not cancelled by client disconnect.
func (l *Limiter) RecordSuccess(name string, ip string) {
	cfg := l.GetConfig(context.Background(), name)
	key := fmt.Sprintf("rl:%s:%s", name, ip)

	luaScript := redis.NewScript(`
		local count = redis.call('INCR', KEYS[1])
		if count == 1 then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end
		return count
	`)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := luaScript.Run(ctx, l.rdb, []string{key}, int(cfg.Window.Seconds())).Err(); err != nil {
		l.logger.Error("failed to record rate limit", "name", name, "error", err)
	}
}

// Record increments the counter immediately (for non-success-only counting like upload-salts).
func (l *Limiter) Record(ctx context.Context, name string, ip string) {
	cfg := l.GetConfig(ctx, name)
	key := fmt.Sprintf("rl:%s:%s", name, ip)

	luaScript := redis.NewScript(`
		local count = redis.call('INCR', KEYS[1])
		if count == 1 then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end
		return count
	`)
	if err := luaScript.Run(ctx, l.rdb, []string{key}, int(cfg.Window.Seconds())).Err(); err != nil {
		l.logger.Error("failed to record rate limit", "name", name, "error", err)
	}
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
