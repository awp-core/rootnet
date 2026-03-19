#!/usr/bin/env bash
# Mine vanity salts for Alpha token deployment and upload to API.
#
# Fetches mining parameters (factory, initCodeHash, vanityRule) from the API,
# mines EIP-55 compliant salts using cast create2, and uploads them via
# POST /api/vanity/upload-salts.
#
# Usage:
#   ./mine-alpha-salts.sh [count] [api_url]
#
# Examples:
#   ./mine-alpha-salts.sh 20                              # Mine 20, upload to default API
#   ./mine-alpha-salts.sh 50 https://tapi.awp.sh          # Mine 50, custom API
#   ./mine-alpha-salts.sh 10 http://localhost:8001         # Mine 10, local
#
# Prerequisites:
#   - cast (Foundry)
#   - curl, jq

set -euo pipefail

COUNT="${1:-20}"
API_BASE="${2:-https://tapi.awp.sh}"
API="${API_BASE}/api"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*" >&2; exit 1; }

command -v cast  >/dev/null 2>&1 || error "cast (Foundry) not found"
command -v curl  >/dev/null 2>&1 || error "curl not found"
command -v jq    >/dev/null 2>&1 || error "jq not found"

# ── Fetch mining parameters from API ──

info "Fetching mining params from ${API}/vanity/mining-params ..."
PARAMS=$(curl -sf "${API}/vanity/mining-params") || error "Failed to fetch mining params"

FACTORY=$(echo "$PARAMS" | jq -r '.factoryAddress')
INIT_HASH=$(echo "$PARAMS" | jq -r '.initCodeHash')
VANITY_RULE=$(echo "$PARAMS" | jq -r '.vanityRule')

[[ -n "$FACTORY" && "$FACTORY" != "null" ]] || error "Missing factoryAddress"
[[ -n "$INIT_HASH" && "$INIT_HASH" != "null" ]] || error "Missing initCodeHash"

info "Factory:      $FACTORY"
info "InitCodeHash: $INIT_HASH"
info "VanityRule:   $VANITY_RULE"

# ── Decode vanity rule into cast create2 flags ──

# Rule is uint64 packed: 8 bytes = 4 prefix + 4 suffix
# Per byte: 0-9=digit, 10-15=lowercase a-f, 16-21=uppercase A-F, >=22=wildcard
decode_position() {
    local val=$1
    if (( val <= 9 )); then
        printf '%d' "$val"
    elif (( val <= 15 )); then
        printf "\\x$(printf '%02x' $((val - 10 + 97)))"  # a-f
    elif (( val <= 21 )); then
        printf "\\x$(printf '%02x' $((val - 16 + 65)))"  # A-F
    else
        echo ""  # wildcard
    fi
}

RULE_HEX=$(printf '%016x' "$(( $VANITY_RULE ))" 2>/dev/null || echo "0000000000000000")

PREFIX=""
SUFFIX=""
CASE_SENSITIVE=false

for i in 0 1 2 3; do
    byte_hex="${RULE_HEX:$((i*2)):2}"
    val=$((16#$byte_hex))
    if (( val < 22 )); then
        ch=$(decode_position $val)
        PREFIX="${PREFIX}${ch}"
        if (( val >= 10 )); then CASE_SENSITIVE=true; fi
    fi
done

for i in 4 5 6 7; do
    byte_hex="${RULE_HEX:$((i*2)):2}"
    val=$((16#$byte_hex))
    if (( val < 22 )); then
        ch=$(decode_position $val)
        SUFFIX="${SUFFIX}${ch}"
        if (( val >= 10 )); then CASE_SENSITIVE=true; fi
    fi
done

if [[ -z "$PREFIX" && -z "$SUFFIX" ]]; then
    error "VanityRule has no constraints (all wildcards)"
fi

info "Pattern:      prefix='${PREFIX}' suffix='${SUFFIX}' case_sensitive=${CASE_SENSITIVE}"

# Build cast create2 args
CAST_ARGS=(create2 --deployer "$FACTORY" --init-code-hash "$INIT_HASH")
[[ -n "$PREFIX" ]] && CAST_ARGS+=(--starts-with "$PREFIX")
[[ -n "$SUFFIX" ]] && CAST_ARGS+=(--ends-with "$SUFFIX")
[[ "$CASE_SENSITIVE" == "true" ]] && CAST_ARGS+=(--case-sensitive)

echo -e "\n${CYAN}══════ Mining $COUNT salts ══════${NC}\n"

# ── Mine salts with EIP-55 post-validation ──

SALTS_JSON="["
MINED=0
ATTEMPTS=0

while (( MINED < COUNT )); do
    RESULT=$(cast "${CAST_ARGS[@]}" 2>&1) || { warn "cast create2 failed, retrying..."; continue; }

    SALT=$(echo "$RESULT" | grep "Salt:" | awk '{print $2}')
    ADDR=$(echo "$RESULT" | grep "Address:" | awk '{print $2}')

    [[ -n "$SALT" && -n "$ADDR" ]] || { warn "Failed to parse cast output"; continue; }

    ATTEMPTS=$((ATTEMPTS + 1))

    # EIP-55 post-validation: verify against contract's predictDeployAddress
    # (cast create2 --case-sensitive only checks literal match, not EIP-55 checksum)
    PREDICTED=$(cast call "$FACTORY" "predictDeployAddress(bytes32)(address)" "$SALT" --rpc-url "https://bsc-dataseed.binance.org" 2>/dev/null || echo "")

    if [[ -z "$PREDICTED" || "$(echo "$PREDICTED" | tr '[:upper:]' '[:lower:]')" != "$(echo "$ADDR" | tr '[:upper:]' '[:lower:]')" ]]; then
        warn "Address mismatch for salt $SALT, skipping"
        continue
    fi

    # Use on-chain predicted address (has correct EIP-55 checksum from Solidity)
    ADDR="$PREDICTED"

    MINED=$((MINED + 1))
    info "[$MINED/$COUNT] $ADDR (attempt $ATTEMPTS)"

    # Append to JSON array
    if (( MINED > 1 )); then SALTS_JSON+=","; fi
    SALTS_JSON+="{\"salt\":\"$SALT\",\"address\":\"$ADDR\"}"
done

SALTS_JSON+="]"

# ── Upload to API ──

echo -e "\n${CYAN}══════ Uploading $MINED salts ══════${NC}\n"

UPLOAD_RESULT=$(curl -sf -X POST "${API}/vanity/upload-salts" \
    -H "Content-Type: application/json" \
    -d "{\"salts\": $SALTS_JSON}") || error "Upload failed"

INSERTED=$(echo "$UPLOAD_RESULT" | jq -r '.inserted')
REJECTED=$(echo "$UPLOAD_RESULT" | jq -r '.rejected')

info "Inserted: $INSERTED"
[[ "$REJECTED" != "0" ]] && warn "Rejected: $REJECTED"

# ── Pool status ──

POOL_COUNT=$(curl -sf "${API}/vanity/salts/count" | jq -r '.available' 2>/dev/null || echo "?")
info "Salt pool: $POOL_COUNT available"

echo -e "\n${GREEN}Done!${NC}"
