#!/usr/bin/env bash
# AWP — Admin CLI
#
# Hot-update rate limits, view system status, manage salt pool, and more.
# All rate limits are stored in Redis and take effect instantly (no restart).
#
# Usage:
#   ./scripts/admin.sh <command> [args...]
#
# Commands:
#   status                          Show API health, indexer sync, salt pool, rate limits
#   limits                          Show all current rate limit configs
#   set-limit <name> <limit> <window_sec>   Set a rate limit (hot-update)
#   reset-limit <name>              Reset a rate limit to compiled default
#   reset-all-limits                Reset all rate limits to compiled defaults
#   salt-pool                       Show salt pool status
#   salt-mine <count>               Mine and upload vanity salts
#   redis-flush-ratelimits          Clear all IP rate limit counters (not configs)
#   ws-clients                      Show WebSocket connection stats
#   help                            Show this help
#
# Rate limit names:
#   relay          Gasless relay transactions (default: 100/3600s)
#   upload_salts   Vanity salt uploads (default: 5/3600s)
#   compute_salt   Vanity salt computation (default: 20/3600s)
#   ws_connect     WebSocket connections per IP (default: 10, window=0 for concurrent)
#
# Environment:
#   API_URL        API base URL (default: https://tapi.awp.sh)
#   REDIS_URL      Redis URL (default: redis://localhost:6379/0)

set -euo pipefail

API_URL="${API_URL:-https://tapi.awp.sh}"
REDIS_URL="${REDIS_URL:-redis://localhost:6379/0}"

# Parse Redis URL → redis-cli args
_redis_cli() {
    local host port db password
    # Strip redis:// prefix
    local url="${REDIS_URL#redis://}"
    # Extract password if present (format: :password@host:port/db)
    if [[ "$url" == *"@"* ]]; then
        local auth="${url%%@*}"
        password="${auth#:}"
        url="${url#*@}"
    fi
    # Extract host:port/db
    host="${url%%:*}"
    local rest="${url#*:}"
    port="${rest%%/*}"
    db="${rest#*/}"
    [[ -z "$host" ]] && host="localhost"
    [[ -z "$port" || "$port" == "$rest" ]] && port="6379"
    [[ -z "$db" || "$db" == "$rest" ]] && db="0"

    local args=(-h "$host" -p "$port" -n "$db")
    [[ -n "${password:-}" ]] && args+=(-a "$password")
    redis-cli "${args[@]}" "$@"
}

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
BOLD='\033[1m'
NC='\033[0m'

_header() { echo -e "\n${CYAN}═══ $* ═══${NC}"; }
_ok()     { echo -e "${GREEN}✓${NC} $*"; }
_warn()   { echo -e "${YELLOW}!${NC} $*"; }
_err()    { echo -e "${RED}✗${NC} $*"; }

# ── Compiled defaults (must match ratelimit.go) ──
declare -A DEFAULTS=(
    [relay]="100:3600"
    [upload_salts]="5:3600"
    [compute_salt]="20:3600"
    [ws_connect]="10:0"
)

_format_limit() {
    local val="$1"
    local limit="${val%%:*}"
    local window="${val##*:}"
    if [[ "$window" == "0" ]]; then
        echo "${limit} concurrent"
    elif [[ "$window" == "3600" ]]; then
        echo "${limit}/hour"
    elif [[ "$window" == "86400" ]]; then
        echo "${limit}/day"
    else
        echo "${limit}/${window}s"
    fi
}

# ═══════════════════════════════════════════
#  Commands
# ═══════════════════════════════════════════

cmd_status() {
    _header "API Health"
    local health
    health=$(curl -sf "${API_URL}/api/health" 2>/dev/null) && _ok "API: $health" || _err "API unreachable"

    _header "Indexer Sync"
    local sync_block
    sync_block=$(_redis_cli --no-auth-warning GET "indexer:last_block" 2>/dev/null) || true
    if [[ -n "$sync_block" ]]; then
        _ok "Last indexed block: $sync_block"
    else
        _warn "Sync state not in Redis (check DB sync_states table)"
    fi

    _header "Salt Pool"
    local count
    count=$(curl -sf "${API_URL}/api/vanity/salts/count" 2>/dev/null | python3 -c "import json,sys;print(json.load(sys.stdin).get('available','?'))" 2>/dev/null) || count="?"
    echo "  Available: $count"

    _header "Rate Limits"
    cmd_limits

    _header "Redis Info"
    local info
    info=$(_redis_cli --no-auth-warning INFO keyspace 2>/dev/null | grep "^db") || true
    echo "  $info"
}

cmd_limits() {
    echo -e "  ${BOLD}Name              Current          Default${NC}"
    echo "  ────────────────────────────────────────────"
    for name in relay upload_salts compute_salt ws_connect; do
        local current default
        current=$(_redis_cli --no-auth-warning HGET ratelimit:config "$name" 2>/dev/null) || current=""
        default="${DEFAULTS[$name]}"
        local display
        if [[ -n "$current" ]]; then
            display="$(_format_limit "$current")  ${YELLOW}(custom)${NC}"
        else
            display="$(_format_limit "$default")  (default)"
        fi
        printf "  %-17s %-20b %s\n" "$name" "$display" "$(_format_limit "$default")"
    done

    echo ""
    echo "  Active IP counters:"
    for name in relay upload_salts compute_salt; do
        local count
        count=$(_redis_cli --no-auth-warning KEYS "rl:${name}:*" 2>/dev/null | wc -l) || count=0
        echo "    $name: $count IPs tracked"
    done
}

cmd_set_limit() {
    local name="${1:-}"
    local limit="${2:-}"
    local window="${3:-}"

    if [[ -z "$name" || -z "$limit" || -z "$window" ]]; then
        echo "Usage: $0 set-limit <name> <limit> <window_seconds>"
        echo "  e.g.: $0 set-limit relay 200 7200    # 200 requests per 2 hours"
        echo "  e.g.: $0 set-limit ws_connect 20 0   # 20 concurrent WS connections"
        echo ""
        echo "Names: relay, upload_salts, compute_salt, ws_connect"
        exit 1
    fi

    if [[ -z "${DEFAULTS[$name]+x}" ]]; then
        _err "Unknown limit name: $name"
        echo "Valid names: ${!DEFAULTS[*]}"
        exit 1
    fi

    _redis_cli --no-auth-warning HSET ratelimit:config "$name" "${limit}:${window}" >/dev/null
    _ok "Set $name = $(_format_limit "${limit}:${window}") (effective immediately)"
}

cmd_reset_limit() {
    local name="${1:-}"
    if [[ -z "$name" ]]; then
        echo "Usage: $0 reset-limit <name>"
        exit 1
    fi
    _redis_cli --no-auth-warning HDEL ratelimit:config "$name" >/dev/null
    _ok "Reset $name to default: $(_format_limit "${DEFAULTS[$name]}")"
}

cmd_reset_all() {
    _redis_cli --no-auth-warning DEL ratelimit:config >/dev/null
    _ok "All rate limits reset to compiled defaults"
    cmd_limits
}

cmd_salt_pool() {
    _header "Salt Pool"
    local count
    count=$(curl -sf "${API_URL}/api/vanity/salts/count" 2>/dev/null | python3 -c "import json,sys;print(json.load(sys.stdin).get('available','?'))" 2>/dev/null) || count="?"
    echo "  Available: $count"

    local params
    params=$(curl -sf "${API_URL}/api/vanity/mining-params" 2>/dev/null) || params="{}"
    echo "  Factory:      $(echo "$params" | python3 -c "import json,sys;print(json.load(sys.stdin).get('factoryAddress','?'))" 2>/dev/null)"
    echo "  InitCodeHash: $(echo "$params" | python3 -c "import json,sys;print(json.load(sys.stdin).get('initCodeHash','?'))" 2>/dev/null)"
    echo "  VanityRule:   $(echo "$params" | python3 -c "import json,sys;print(json.load(sys.stdin).get('vanityRule','?'))" 2>/dev/null)"
}

cmd_salt_mine() {
    local count="${1:-20}"
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    if [[ -x "$script_dir/vanity/mine-alpha-salts.sh" ]]; then
        "$script_dir/vanity/mine-alpha-salts.sh" "$count" "$API_URL"
    else
        _err "mine-alpha-salts.sh not found at $script_dir/vanity/"
        exit 1
    fi
}

cmd_flush_ratelimits() {
    local count=0
    for prefix in "rl:relay:" "rl:upload_salts:" "rl:compute_salt:"; do
        local keys
        keys=$(_redis_cli --no-auth-warning KEYS "${prefix}*" 2>/dev/null) || true
        if [[ -n "$keys" ]]; then
            echo "$keys" | while read -r key; do
                _redis_cli --no-auth-warning DEL "$key" >/dev/null
                count=$((count + 1))
            done
        fi
    done
    _ok "Flushed all IP rate limit counters (configs preserved)"
}

cmd_ws_clients() {
    _header "WebSocket Connections"
    local ws_limit
    ws_limit=$(_redis_cli --no-auth-warning HGET ratelimit:config ws_connect 2>/dev/null) || ws_limit=""
    if [[ -n "$ws_limit" ]]; then
        echo "  Limit: $(_format_limit "$ws_limit") per IP (custom)"
    else
        echo "  Limit: $(_format_limit "${DEFAULTS[ws_connect]}") per IP (default)"
    fi
    echo "  (Connection count is tracked in-memory, not queryable via CLI)"
}

cmd_help() {
    head -30 "$0" | grep -E "^#" | sed 's/^# \?//'
}

# ═══════════════════════════════════════════
#  Main
# ═══════════════════════════════════════════

case "${1:-help}" in
    status)              cmd_status ;;
    limits)              cmd_limits ;;
    set-limit)           cmd_set_limit "${2:-}" "${3:-}" "${4:-}" ;;
    reset-limit)         cmd_reset_limit "${2:-}" ;;
    reset-all-limits)    cmd_reset_all ;;
    salt-pool)           cmd_salt_pool ;;
    salt-mine)           cmd_salt_mine "${2:-20}" ;;
    redis-flush-ratelimits) cmd_flush_ratelimits ;;
    ws-clients)          cmd_ws_clients ;;
    help|--help|-h)      cmd_help ;;
    *)
        _err "Unknown command: $1"
        echo "Run '$0 help' for usage"
        exit 1
        ;;
esac
