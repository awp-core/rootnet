#!/usr/bin/env bash
# Vanity address miner for CREATE2 deterministic deployment.
#
# Uses cast create2 (Foundry) for all mining operations.
#
# Pattern rules (per contract in salt.json):
#   prefix with a-f/A-F  → cast create2 --case-sensitive (EIP-55)
#   prefix digits only    → cast create2 nibble matching
#   suffix                → cast create2 --ends-with
#   no pattern            → zero salt
#
# Usage:
#   ./mine.sh [salt.json] [threads]
#
# Prerequisites:
#   - cast (Foundry)
#   - jq

set -euo pipefail

CONFIG="${1:-salt.json}"
THREADS="${2:-$(nproc)}"

# ─── Colors ───
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }

# ─── Check prerequisites ───
command -v cast >/dev/null 2>&1 || { echo "Error: cast not found. Install Foundry." >&2; exit 1; }
command -v jq  >/dev/null 2>&1 || { echo "Error: jq not found." >&2; exit 1; }
[[ -f "$CONFIG" ]] || { echo "Error: Config file $CONFIG not found" >&2; exit 1; }

DEPLOYER=$(jq -r '.deployer // "0x4e59b44847b379578588920cA78FbF26c0B4956C"' "$CONFIG")
COUNT=$(jq '.contracts | length' "$CONFIG")

echo "════════════════════════════════════════════"
echo "  AWP Vanity Address Miner"
echo "════════════════════════════════════════════"
echo ""
echo "  Config:    $CONFIG"
echo "  Deployer:  $DEPLOYER"
echo "  Contracts: $COUNT"
echo "  Threads:   $THREADS (cast create2)"
echo ""

# ─── Helper: pattern with letters needs EIP-55 ---
pattern_needs_eip55() {
    echo "$1" | grep -qE '[a-fA-F]'
}

# ─── Mine ───
mine() {
    local name="$1" deployer="$2" init_code_hash="$3" prefix="$4" suffix="$5" case_sensitive="$6" idx="$7"

    local args=("--deployer" "$deployer" "--init-code-hash" "$init_code_hash" "--threads" "$THREADS")
    [[ "$case_sensitive" == "true" ]] && args+=("--case-sensitive")
    [[ -n "$prefix" ]] && args+=("--starts-with" "$prefix")
    [[ -n "$suffix" ]] && args+=("--ends-with" "$suffix")

    local mode="nibble"
    [[ "$case_sensitive" == "true" ]] && mode="EIP-55"
    echo -e "[$name] ${CYAN}cast create2${NC} prefix=\"$prefix\" suffix=\"$suffix\" ($mode, $THREADS threads)..."

    local start end elapsed_ms
    start=$(date +%s%N)
    local output
    output=$(cast create2 "${args[@]}" 2>&1)
    end=$(date +%s%N)
    elapsed_ms=$(( (end - start) / 1000000 ))

    local addr salt
    addr=$(echo "$output" | grep "^Address:" | grep -oE '0x[0-9a-fA-F]{40}')
    salt=$(echo "$output" | grep "^Salt:" | grep -oE '0x[0-9a-fA-F]{64}')

    if [[ -z "$addr" || -z "$salt" ]]; then
        warn "[$name] Mining failed:"
        echo "$output"
        return 1
    fi

    echo "[$name] Found in ${elapsed_ms}ms"
    echo "  Salt:    $salt"
    echo "  Address: $addr"

    local tmp
    tmp=$(mktemp)
    jq ".contracts[$idx].salt = \"$salt\" | .contracts[$idx].address = \"$addr\"" "$CONFIG" > "$tmp" && mv "$tmp" "$CONFIG"
}

# ─── Main loop ───
for i in $(seq 0 $((COUNT - 1))); do
    NAME=$(jq -r ".contracts[$i].name" "$CONFIG")
    INIT_CODE_HASH=$(jq -r ".contracts[$i].initCodeHash" "$CONFIG")
    PREFIX=$(jq -r ".contracts[$i].prefix // empty" "$CONFIG")
    SUFFIX=$(jq -r ".contracts[$i].suffix // empty" "$CONFIG")
    EXISTING_SALT=$(jq -r ".contracts[$i].salt // empty" "$CONFIG")

    if [[ -n "$EXISTING_SALT" ]]; then
        ADDR=$(jq -r ".contracts[$i].address" "$CONFIG")
        echo "[$NAME] Already mined: $ADDR"
        continue
    fi

    if [[ -z "$INIT_CODE_HASH" ]]; then
        echo "[$NAME] Missing initCodeHash, skipping"
        continue
    fi

    if [[ -z "$PREFIX" && -z "$SUFFIX" ]]; then
        # Use a random salt to avoid CreateCollision on redeployment
        RAND_SALT="0x$(openssl rand -hex 32)"
        ADDR=$(cast create2 --deployer "$DEPLOYER" --init-code-hash "$INIT_CODE_HASH" --salt "$RAND_SALT" 2>/dev/null | grep -oE '0x[0-9a-fA-F]{40}' || echo "unknown")
        TMP=$(mktemp)
        jq ".contracts[$i].salt = \"$RAND_SALT\" | .contracts[$i].address = \"$ADDR\"" "$CONFIG" > "$TMP" && mv "$TMP" "$CONFIG"
        echo "[$NAME] No pattern (random salt) → $ADDR"
        continue
    fi

    NEEDS_EIP55=false
    pattern_needs_eip55 "${PREFIX}${SUFFIX}" && NEEDS_EIP55=true

    if $NEEDS_EIP55; then
        mine "$NAME" "$DEPLOYER" "$INIT_CODE_HASH" "$PREFIX" "$SUFFIX" "true" "$i"
    else
        mine "$NAME" "$DEPLOYER" "$INIT_CODE_HASH" "$PREFIX" "$SUFFIX" "false" "$i"
    fi
done

echo ""
echo "Results written to $CONFIG"
echo ""
name_to_env_key() {
    case "$1" in
        AWPToken)           echo "AWP_TOKEN" ;;
        AlphaTokenFactory)  echo "ALPHA_FACTORY" ;;
        AWPEmission_impl)   echo "EMISSION_IMPL" ;;
        AWPEmission_proxy)  echo "EMISSION_PROXY" ;;
        Treasury)           echo "TREASURY" ;;
        AWPRegistry)        echo "AWP_REGISTRY" ;;
        WorknetNFT)          echo "WORKNET_NFT" ;;
        LPManager)          echo "LP_MANAGER" ;;
        StakingVault)       echo "STAKING_VAULT" ;;
        StakeNFT)           echo "STAKE_NFT" ;;
        AWPDAO)             echo "DAO" ;;
        WorknetManager_impl) echo "WORKNET_MANAGER_IMPL" ;;
        *)                  echo "$1" | tr '[:lower:]' '[:upper:]' ;;
    esac
}

echo "# Copy to contracts/.env:"
for i in $(seq 0 $((COUNT - 1))); do
    NAME=$(jq -r ".contracts[$i].name" "$CONFIG")
    SALT=$(jq -r ".contracts[$i].salt // empty" "$CONFIG")
    if [[ -n "$SALT" ]]; then
        ENV_KEY=$(name_to_env_key "$NAME")
        echo "SALT_${ENV_KEY}=$SALT"
    fi
done
