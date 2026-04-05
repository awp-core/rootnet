#!/usr/bin/env bash
# AWP — One-click deploy: mine vanity salts → deploy contracts → verify → generate api/.env
#
# Flow:
#   1. Build contracts
#   2. Compute initCodeHashes (tiered, re-run per tier as addresses become known)
#   3. Mine vanity salts via cast create2 (scripts/vanity/mine.sh)
#   4. Deploy all contracts with mined salts
#   5. Verify on BscScan (optional)
#   6. Generate api/.env
#   7. Initialize database
#
# Usage:
#   ./scripts/deploy.sh                  # Full: mine + deploy + verify + api config
#   ./scripts/deploy.sh --dry-run        # Simulate deployment (no broadcast)
#   ./scripts/deploy.sh --skip-mine      # Skip salt mining (use existing SALT_* from .env)
#   ./scripts/deploy.sh --verify-only    # Verify existing contracts only
#
# Input:  contracts/.env + scripts/vanity/salt.json
# Output: api/.env (ready for API services)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
CONTRACTS_DIR="$ROOT_DIR/contracts"
API_DIR="$ROOT_DIR/api"
VANITY_DIR="$SCRIPT_DIR/vanity"
SALT_JSON="$VANITY_DIR/salt.json"
OUTPUT_ENV="$API_DIR/.env"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*" >&2; exit 1; }
step()  { echo -e "\n${CYAN}══════ $* ══════${NC}"; }

# ─── Prerequisites ───

command -v forge >/dev/null 2>&1 || error "forge not found"
command -v cast  >/dev/null 2>&1 || error "cast not found"
command -v jq    >/dev/null 2>&1 || error "jq not found"

# ─── Load config ───

[[ -f "$CONTRACTS_DIR/.env" ]] || error "Missing $CONTRACTS_DIR/.env"
# Source .env but do NOT override variables already exported by deploy-multichain.sh
# Use a temp file that wraps each line with ${VAR:-value} pattern
while IFS='=' read -r key value; do
    [[ -z "$key" || "$key" == \#* ]] && continue
    # Only set if not already in environment
    if [[ -z "${!key+x}" ]]; then
        export "$key=$value"
    fi
done < <(grep -v '^\s*#' "$CONTRACTS_DIR/.env" | grep -v '^\s*$' | grep '=')

# ─── Parse flags ───

DRY_RUN=""
SKIP_MINE=""
VERIFY_ONLY=""
for arg in "$@"; do
    case "$arg" in
        --dry-run)     DRY_RUN="1" ;;
        --skip-mine)   SKIP_MINE=1 ;;
        --verify-only) VERIFY_ONLY=1 ;;
    esac
done

# ─── Validate required vars ───

for var in ETH_RPC_URL DEPLOYER_PRIVATE_KEY GUARDIAN LIQUIDITY_POOL AIRDROP POOL_MANAGER POSITION_MANAGER PERMIT2 CL_SWAP_ROUTER GENESIS_TIME; do
    [[ -n "${!var:-}" ]] || error "Missing required env var: $var"
done

DEPLOYER=$(cast wallet address --private-key "$DEPLOYER_PRIVATE_KEY")
DATABASE_URL="${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/awp?sslmode=disable}"
REDIS_URL="${REDIS_URL:-redis://localhost:6379/0}"
HTTP_ADDR="${HTTP_ADDR:-:8080}"
KEEPER_PRIVATE_KEY="${KEEPER_PRIVATE_KEY:-}"
RELAYER_PRIVATE_KEY="${RELAYER_PRIVATE_KEY:-}"
VANITY_RULE="${VANITY_RULE:-0}"
ETHERSCAN_API_KEY="${ETHERSCAN_API_KEY:-}"

if [[ -n "$VERIFY_ONLY" ]]; then
    # Skip to verification
    [[ -f "$OUTPUT_ENV" ]] || error "No $OUTPUT_ENV found. Deploy first."
    set -a; source "$OUTPUT_ENV"; set +a
    FACTORY_ADDRESS="${ALPHA_FACTORY_ADDRESS:-}"
    AWP_EMISSION_IMPL="${AWP_EMISSION_IMPL:-}"
    CHAIN_ID=$(cast chain-id --rpc-url "$ETH_RPC_URL" 2>/dev/null)
    if [ -z "$CHAIN_ID" ]; then echo "ERROR: Cannot determine CHAIN_ID from RPC"; exit 1; fi
    cd "$CONTRACTS_DIR"
    # fall through to verification section below
else

# ═══════════════════════════════════════════════
#  STEP 1: BUILD
# ═══════════════════════════════════════════════

step "Step 1/7 — Build contracts"
cd "$CONTRACTS_DIR"
forge build --force --quiet
info "Build complete"

# ═══════════════════════════════════════════════
#  STEP 2-3: TIERED SALT MINING
# ═══════════════════════════════════════════════

if [[ -z "$SKIP_MINE" ]]; then
    step "Step 2/7 — Compute initCodeHashes & mine vanity salts"

    # Generate salt.json from .env vanity patterns
    info "Generating $SALT_JSON from vanity patterns..."
    _vp() { echo "${!1:-}"; }  # read env var by name, empty if unset
    python3 - "$SALT_JSON" <<'PYEOF'
import json, os, sys
contracts = [
    {"name": "AWPToken",           "env_prefix": "AWP_TOKEN"},
    {"name": "AlphaTokenFactory",  "env_prefix": "ALPHA_FACTORY"},
    {"name": "AWPEmission_impl",   "env_prefix": "EMISSION_IMPL"},
    {"name": "AWPRegistry_impl",   "env_prefix": "AWP_REGISTRY_IMPL"},
    {"name": "WorknetManager_impl", "env_prefix": "WORKNET_MANAGER_IMPL"},
    {"name": "Treasury",           "env_prefix": "TREASURY"},
    {"name": "AWPRegistry",        "env_prefix": "AWP_REGISTRY"},
    {"name": "WorknetNFT",          "env_prefix": "WORKNET_NFT"},
    {"name": "LPManager",          "env_prefix": "LP_MANAGER"},
    {"name": "StakingVault_impl",   "env_prefix": "STAKING_VAULT_IMPL"},
    {"name": "StakingVault",       "env_prefix": "STAKING_VAULT"},
    {"name": "AWPEmission_proxy",  "env_prefix": "EMISSION_PROXY"},
    {"name": "StakeNFT",           "env_prefix": "STAKE_NFT"},
    {"name": "AWPDAO",             "env_prefix": "DAO"},
]
# Merge mode: preserve existing salt/address from salt.json if present
existing = {}
path = sys.argv[1]
try:
    with open(path) as f:
        old = json.load(f)
    for c in old.get("contracts", []):
        if c.get("salt"):
            existing[c["name"]] = {"salt": c["salt"], "address": c.get("address", "")}
except (FileNotFoundError, json.JSONDecodeError):
    pass

out = {"deployer": "0x4e59b44847b379578588920cA78FbF26c0B4956C", "contracts": []}
for c in contracts:
    prefix = os.environ.get(f"VANITY_PREFIX_{c['env_prefix']}", "")
    suffix = os.environ.get(f"VANITY_SUFFIX_{c['env_prefix']}", "")
    prev = existing.get(c["name"], {})
    out["contracts"].append({
        "name": c["name"],
        "prefix": prefix,
        "suffix": suffix,
        "initCodeHash": "",
        "salt": prev.get("salt", ""),
        "address": prev.get("address", "")
    })
with open(path, "w") as f:
    json.dump(out, f, indent=2)
kept = sum(1 for c in out["contracts"] if c["salt"])
print(f"Generated {path} with {len(contracts)} contracts ({kept} salts preserved)")
PYEOF

    # Helper: run InitCodeHashes.s.sol and parse output
    # Resolves chain-specific variants: "LPManager (Uniswap)" → "LPManager" on Uni chains
    compute_hashes() {
        local chain_id
        chain_id=$(cast chain-id --rpc-url "$ETH_RPC_URL" 2>/dev/null || echo "0")
        local dex_suffix="Uniswap"
        [[ "$chain_id" == "56" || "$chain_id" == "97" ]] && dex_suffix="PancakeSwap"

        forge script script/InitCodeHashes.s.sol --rpc-url "$ETH_RPC_URL" 2>&1 | grep ': 0x' | while read -r line; do
            name="${line%%:*}"
            name=$(echo "$name" | xargs) # trim
            hash="${line##*: }"
            # Skip Deployer line
            [[ "$name" == "Deployer" ]] && continue
            # For chain-specific contracts, only emit the matching variant
            if [[ "$name" == *"(PancakeSwap)"* || "$name" == *"(Uniswap)"* ]]; then
                if [[ "$name" != *"($dex_suffix)"* ]]; then
                    continue # skip wrong DEX variant
                fi
                # Strip variant suffix: "LPManager (Uniswap)" → "LPManager"
                name="${name%% (*}"
            fi
            # Map proxy names to salt.json names
            [[ "$name" == "AWPRegistry_proxy" ]] && name="AWPRegistry"
            [[ "$name" == "StakingVault_proxy" ]] && name="StakingVault"
            echo "$name $hash"
        done
    }

    # Helper: update salt.json with initCodeHash for a contract
    set_hash() {
        local name="$1" hash="$2"
        local idx
        idx=$(jq --arg n "$name" '[.contracts[].name] | index($n)' "$SALT_JSON")
        if [[ "$idx" != "null" && -n "$idx" ]]; then
            local tmp=$(mktemp)
            jq ".contracts[$idx].initCodeHash = \"$hash\"" "$SALT_JSON" > "$tmp" && mv "$tmp" "$SALT_JSON"
        fi
    }

    # Helper: read mined address from salt.json
    get_addr() {
        jq -r --arg n "$1" '.contracts[] | select(.name==$n) | .address // empty' "$SALT_JSON"
    }

    # Helper: export mined address as ADDR_* env var for next tier
    export_addr() {
        local name="$1"
        local addr=$(get_addr "$name")
        [[ -n "$addr" ]] || return
        case "$name" in
            AWPToken)          export ADDR_AWP_TOKEN="$addr" ;;
            AlphaTokenFactory) export ADDR_FACTORY="$addr" ;;
            AWPEmission_impl)  export ADDR_EMISSION_IMPL="$addr" ;;
            Treasury)          export ADDR_TREASURY="$addr" ;;
            AWPRegistry)       export ADDR_AWP_REGISTRY="$addr" ;;
            WorknetNFT)         export ADDR_WORKNET_NFT="$addr" ;;
            StakingVault)      export ADDR_STAKING_VAULT="$addr" ;;
            StakeNFT)          export ADDR_STAKE_NFT="$addr" ;;
            AWPRegistry_impl)  export ADDR_AWP_REGISTRY_IMPL="$addr" ;;
            StakingVault_impl) export ADDR_STAKING_VAULT_IMPL="$addr" ;;
        esac
    }

    # Note: salts are preserved from previous runs (merge mode in Python block above).
    # To force fresh mining, delete salt.json before running.

    # Tier definitions: mine in order, export addresses between tiers
    TIERS=(
        "AWPToken AlphaTokenFactory AWPEmission_impl AWPRegistry_impl StakingVault_impl WorknetManager_impl"
        "Treasury"
        "AWPRegistry"
        "WorknetNFT LPManager StakingVault AWPEmission_proxy"
        "StakeNFT"
        "AWPDAO"
    )

    for tier_idx in "${!TIERS[@]}"; do
        tier="${TIERS[$tier_idx]}"
        info "Tier $((tier_idx + 1))/6: $tier"

        # Compute hashes for this tier
        while read -r name hash; do
            for t in $tier; do
                if [[ "$name" == "$t" ]]; then
                    set_hash "$name" "$hash"
                fi
            done
        done < <(compute_hashes)

        # Mine salts for contracts in this tier that have patterns
        cd "$VANITY_DIR"
        ./mine.sh "$SALT_JSON"
        cd "$CONTRACTS_DIR"

        # Export mined addresses for next tier
        for t in $tier; do
            export_addr "$t"
        done
    done

    info "All salts mined!"

    # Write SALT_* env vars to contracts/.env for Deploy.s.sol
    sed -i '/^SALT_/d' "$CONTRACTS_DIR/.env"
    {
        echo ""
        echo "# ── Mined CREATE2 salts (auto-generated by deploy.sh) ──"
        for i in $(seq 0 $(($(jq '.contracts | length' "$SALT_JSON") - 1))); do
            NAME=$(jq -r ".contracts[$i].name" "$SALT_JSON")
            SALT=$(jq -r ".contracts[$i].salt // empty" "$SALT_JSON")
            if [[ -n "$SALT" ]]; then
                case "$NAME" in
                    AWPToken)          echo "SALT_AWP_TOKEN=$SALT" ;;
                    AlphaTokenFactory) echo "SALT_ALPHA_FACTORY=$SALT" ;;
                    AWPEmission_impl)  echo "SALT_EMISSION_IMPL=$SALT" ;;
                    AWPEmission_proxy) echo "SALT_EMISSION_PROXY=$SALT" ;;
                    Treasury)          echo "SALT_TREASURY=$SALT" ;;
                    AWPRegistry)       echo "SALT_AWP_REGISTRY=$SALT" ;;
                    WorknetNFT)         echo "SALT_WORKNET_NFT=$SALT" ;;
                    LPManager)         echo "SALT_LP_MANAGER=$SALT" ;;
                    StakingVault)      echo "SALT_STAKING_VAULT=$SALT" ;;
                    StakeNFT)          echo "SALT_STAKE_NFT=$SALT" ;;
                    AWPDAO)            echo "SALT_DAO=$SALT" ;;
                    WorknetManager_impl) echo "SALT_WORKNET_MANAGER_IMPL=$SALT" ;;
                    AWPRegistry_impl)  echo "SALT_AWP_REGISTRY_IMPL=$SALT" ;;
                    StakingVault_impl) echo "SALT_STAKING_VAULT_IMPL=$SALT" ;;
                esac
            fi
        done
    } >> "$CONTRACTS_DIR/.env"

    # Reload only SALT_* vars from .env (don't overwrite DEX addresses from deploy-multichain.sh)
    while IFS='=' read -r key value; do
        [[ "$key" == SALT_* ]] && export "$key=$value"
    done < <(grep '^SALT_' "$CONTRACTS_DIR/.env")

else
    step "Step 2/7 — Skipping salt mining (--skip-mine)"
fi

# ═══════════════════════════════════════════════
#  STEP 2.5: VALIDATE VANITY ADDRESSES
# ═══════════════════════════════════════════════
# Validate mined salts match .env VANITY_PREFIX/SUFFIX rules
# Runs after both mining and --skip-mine paths

validate_vanity() {
    local name="$1" env_prefix="$2"
    local want_prefix="${!env_prefix:-}"
    local suffix_var="${env_prefix/PREFIX/SUFFIX}"
    local want_suffix="${!suffix_var:-}"

    # No vanity rule, skip
    [[ -z "$want_prefix" && -z "$want_suffix" ]] && return 0

    # Read address from salt.json
    local addr
    addr=$(jq -r --arg n "$name" '.contracts[] | select(.name==$n) | .address // empty' "$SALT_JSON")
    if [[ -z "$addr" ]]; then
        warn "[$name] No address in salt.json — cannot validate vanity"
        return 1
    fi

    # Strip 0x prefix, keep original EIP-55 case for case-sensitive matching
    local hex="${addr#0x}"

    local ok=1
    if [[ -n "$want_prefix" ]]; then
        local actual_prefix="${hex:0:${#want_prefix}}"
        if [[ "$actual_prefix" != "$want_prefix" ]]; then
            warn "[$name] Vanity prefix mismatch: want=$want_prefix got=$actual_prefix (addr=$addr)"
            ok=0
        fi
    fi
    if [[ -n "$want_suffix" ]]; then
        local actual_suffix="${hex: -${#want_suffix}}"
        if [[ "$actual_suffix" != "$want_suffix" ]]; then
            warn "[$name] Vanity suffix mismatch: want=$want_suffix got=$actual_suffix (addr=$addr)"
            ok=0
        fi
    fi

    if [[ "$ok" == "1" ]]; then
        info "[$name] Vanity OK: $addr (prefix=$want_prefix suffix=$want_suffix)"
    else
        error "[$name] Vanity validation failed! Re-run without --skip-mine to re-mine salts."
    fi
}

if [[ -f "$SALT_JSON" ]]; then
    step "Step 2.5/7 — Validate vanity addresses"
    validate_vanity "AWPToken" "VANITY_PREFIX_AWP_TOKEN"
    validate_vanity "AWPRegistry" "VANITY_PREFIX_AWP_REGISTRY"
    # Add more contract vanity rule validations here
fi

# ═══════════════════════════════════════════════
#  STEP 4: DEPLOY
# ═══════════════════════════════════════════════

step "Step 3/7 — Record deploy block"
DEPLOY_BLOCK=$(cast block-number --rpc-url "$ETH_RPC_URL")
info "Starting block: $DEPLOY_BLOCK"

step "Step 4/7 — Deploy contracts"
[[ -n "$DRY_RUN" ]] && info "Dry-run mode"

DEPLOY_OUTPUT=$(forge script script/Deploy.s.sol \
    --rpc-url "$ETH_RPC_URL" \
    $([ -z "$DRY_RUN" ] && echo "--broadcast") 2>&1) || { echo "$DEPLOY_OUTPUT"; error "Deployment failed"; }

# Parse addresses
parse_address() {
    echo "$DEPLOY_OUTPUT" | grep -i "$1" | grep -oE '0x[0-9a-fA-F]{40}' | tail -1
}

AWP_TOKEN_ADDRESS=$(parse_address "AWPToken:")
FACTORY_ADDRESS=$(parse_address "AlphaTokenFactory:")
TREASURY_ADDRESS=$(parse_address "Treasury:")
AWP_REGISTRY_IMPL=$(parse_address "AWPRegistry impl:")
AWP_REGISTRY_ADDRESS=$(parse_address "AWPRegistry proxy:")
SUBNETNFT_ADDRESS=$(parse_address "WorknetNFT:")
LP_MANAGER_ADDRESS=$(parse_address "LPManager")
AWP_EMISSION_IMPL=$(parse_address "AWPEmission impl:")
AWP_EMISSION_ADDRESS=$(parse_address "AWPEmission proxy:")
STAKING_VAULT_IMPL=$(parse_address "StakingVault impl:")
STAKING_VAULT_ADDRESS=$(parse_address "StakingVault proxy:")
STAKENFT_ADDRESS=$(parse_address "StakeNFT:")
DAO_ADDRESS=$(parse_address "AWPDAO:")
WORKNET_MANAGER_IMPL=$(parse_address "WorknetManager impl")

for var in AWP_TOKEN_ADDRESS FACTORY_ADDRESS TREASURY_ADDRESS AWP_REGISTRY_IMPL AWP_REGISTRY_ADDRESS SUBNETNFT_ADDRESS \
           LP_MANAGER_ADDRESS AWP_EMISSION_IMPL AWP_EMISSION_ADDRESS STAKING_VAULT_IMPL STAKING_VAULT_ADDRESS \
           STAKENFT_ADDRESS DAO_ADDRESS WORKNET_MANAGER_IMPL; do
    [[ -n "${!var:-}" ]] || error "Failed to parse $var"
done

info "All contracts deployed!"

# Get actual deploy block
if [[ -z "$DRY_RUN" ]]; then
    CHAIN_ID=$(cast chain-id --rpc-url "$ETH_RPC_URL" 2>/dev/null)
    if [ -z "$CHAIN_ID" ]; then echo "ERROR: Cannot determine CHAIN_ID from RPC"; exit 1; fi
    RUN_FILE="$CONTRACTS_DIR/broadcast/Deploy.s.sol/$CHAIN_ID/run-latest.json"
    if [[ -f "$RUN_FILE" ]]; then
        ACTUAL_BLOCK=$(jq -r '.receipts[0].blockNumber // empty' "$RUN_FILE" 2>/dev/null | xargs printf "%d" 2>/dev/null || true)
        [[ "${ACTUAL_BLOCK:-0}" -gt 0 ]] && DEPLOY_BLOCK="$ACTUAL_BLOCK"
    fi
fi

# Compute AlphaToken initCodeHash
# Read initcode hash from deployed AlphaTokenFactory (most reliable — includes Solidity metadata)
ALPHA_INITCODE_HASH=$(cast call "$FACTORY_ADDRESS" "ALPHA_BYTECODE_HASH()(bytes32)" --rpc-url "$ETH_RPC_URL" 2>/dev/null || echo "")

# ═══════════════════════════════════════════════
#  STEP 5: GENERATE API .env
# ═══════════════════════════════════════════════

step "Step 5/7 — Generate $OUTPUT_ENV"

cat > "$OUTPUT_ENV" <<ENVFILE
# AWP — API Configuration
# Generated by scripts/deploy.sh at $(date -u +"%Y-%m-%dT%H:%M:%SZ")

# ── Database ──
DATABASE_URL=$DATABASE_URL

# ── Redis ──
REDIS_URL=$REDIS_URL

# ── HTTP server ──
HTTP_ADDR=$HTTP_ADDR

# ── Chain ──
CHAIN_ID=$CHAIN_ID
RPC_URL=$ETH_RPC_URL

# ── Contract addresses (all 10 protocol contracts) ──
AWP_REGISTRY_ADDRESS=$AWP_REGISTRY_ADDRESS
AWP_TOKEN_ADDRESS=$AWP_TOKEN_ADDRESS
AWP_EMISSION_ADDRESS=$AWP_EMISSION_ADDRESS
STAKING_VAULT_ADDRESS=$STAKING_VAULT_ADDRESS
STAKE_NFT_ADDRESS=$STAKENFT_ADDRESS
SUBNETNFT_ADDRESS=$SUBNETNFT_ADDRESS
LP_MANAGER_ADDRESS=$LP_MANAGER_ADDRESS
ALPHA_FACTORY_ADDRESS=$FACTORY_ADDRESS
DAO_ADDRESS=$DAO_ADDRESS
TREASURY_ADDRESS=$TREASURY_ADDRESS

# ── Indexer start block ──
DEPLOY_BLOCK=$DEPLOY_BLOCK

# ── Vanity address mining ──
ALPHA_INITCODE_HASH=${ALPHA_INITCODE_HASH:-}
VANITY_RULE=$VANITY_RULE

# ── Gasless relay (fill private key to enable) ──
RELAYER_PRIVATE_KEY=$RELAYER_PRIVATE_KEY

# ── Keeper (fill private key to enable) ──
KEEPER_PRIVATE_KEY=$KEEPER_PRIVATE_KEY

# ── For verification (internal) ──
AWP_REGISTRY_IMPL=$AWP_REGISTRY_IMPL
AWP_EMISSION_IMPL=$AWP_EMISSION_IMPL
STAKING_VAULT_IMPL=$STAKING_VAULT_IMPL
WORKNET_MANAGER_IMPL=$WORKNET_MANAGER_IMPL
ENVFILE

chmod 600 "$OUTPUT_ENV"
info "Config written to $OUTPUT_ENV (chmod 600)"

# ═══════════════════════════════════════════════
#  STEP 6: DATABASE
# ═══════════════════════════════════════════════

step "Step 6/7 — Database setup"
if command -v psql >/dev/null 2>&1 && [[ -z "$DRY_RUN" ]]; then
    psql "$DATABASE_URL" -f "$API_DIR/internal/db/schema.sql" 2>/dev/null || info "Schema already exists"
    psql "$DATABASE_URL" -c "
        INSERT INTO sync_states (contract_name, last_block)
        VALUES ('indexer', $((DEPLOY_BLOCK - 1)))
        ON CONFLICT (contract_name) DO UPDATE SET last_block = EXCLUDED.last_block;
    " 2>/dev/null && info "sync_states set to block $((DEPLOY_BLOCK - 1))" \
                  || warn "Failed to set sync_states"
else
    warn "Skipping DB setup (psql not found or dry-run)"
fi

fi # end of non-verify-only block

# ═══════════════════════════════════════════════
#  STEP 7: VERIFY
# ═══════════════════════════════════════════════

step "Step 7/7 — Verify contracts"

if [[ -z "${ETHERSCAN_API_KEY:-}" ]]; then
    warn "ETHERSCAN_API_KEY not set — skipping verification"
    warn "Set it in contracts/.env and re-run with --verify-only"
else
    CHAIN_ID="${CHAIN_ID:-$(cast chain-id --rpc-url "$ETH_RPC_URL" 2>/dev/null || echo "56")}"
    VERIFIER_URL=""
    case "$CHAIN_ID" in
        56)    VERIFIER_URL="https://api.bscscan.com/api" ;;
        97)    VERIFIER_URL="https://api-testnet.bscscan.com/api" ;;
        8453)  VERIFIER_URL="https://api.basescan.org/api" ;;
        84532) VERIFIER_URL="https://api-sepolia.basescan.org/api" ;;
        1)     VERIFIER_URL="https://api.etherscan.io/api" ;;
        42161) VERIFIER_URL="https://api.arbiscan.io/api" ;;
        421614) VERIFIER_URL="https://api-sepolia.arbiscan.io/api" ;;
    esac
    # Select per-chain explorer API key (fallback to ETHERSCAN_API_KEY)
    case "$CHAIN_ID" in
        1) EXPLORER_KEY="${ETHERSCAN_API_KEY:-$ETHERSCAN_API_KEY}" ;;
        56) EXPLORER_KEY="${BSCSCAN_API_KEY:-$ETHERSCAN_API_KEY}" ;;
        8453) EXPLORER_KEY="${BASESCAN_API_KEY:-$ETHERSCAN_API_KEY}" ;;
        42161) EXPLORER_KEY="${ARBISCAN_API_KEY:-$ETHERSCAN_API_KEY}" ;;
        *) EXPLORER_KEY="$ETHERSCAN_API_KEY" ;;
    esac
    VF="--etherscan-api-key $EXPLORER_KEY"
    [[ -n "$VERIFIER_URL" ]] && VF="$VF --verifier-url $VERIFIER_URL"

    vc() {
        local addr="$1" contract="$2" args="${3:-}" label="${4:-$contract}"
        [[ -z "$addr" ]] && { warn "Skip $label — no address"; return; }
        local cmd="forge verify-contract $addr $contract --chain-id $CHAIN_ID $VF --watch"
        [[ -n "$args" ]] && cmd="$cmd --constructor-args $args"
        info "Verifying $label at $addr ..."
        eval "$cmd" 2>&1 | tail -3 || warn "Verification failed for $label (may already be verified)"
    }

    # ── Standalone contracts (direct constructor args) ──
    vc "$AWP_TOKEN_ADDRESS" "src/token/AWPToken.sol:AWPToken" \
        "$(cast abi-encode 'c(string,string,address)' 'AWP Token' 'AWP' "$DEPLOYER")" "AWPToken"
    vc "$FACTORY_ADDRESS" "src/token/AlphaTokenFactory.sol:AlphaTokenFactory" \
        "$(cast abi-encode 'c(address,uint64)' "$DEPLOYER" "$VANITY_RULE")" "AlphaTokenFactory"
    vc "$TREASURY_ADDRESS" "src/governance/Treasury.sol:Treasury" \
        "$(cast abi-encode 'c(uint256,address[],address[],address)' 172800 '[]' '[0x0000000000000000000000000000000000000000]' "$DEPLOYER")" "Treasury"
    vc "$SUBNETNFT_ADDRESS" "src/core/WorknetNFT.sol:WorknetNFT" \
        "$(cast abi-encode 'c(string,string,address)' 'AWP Worknet' 'AWPSUB' "$AWP_REGISTRY_ADDRESS")" "WorknetNFT"
    vc "$STAKENFT_ADDRESS" "src/core/StakeNFT.sol:StakeNFT" \
        "$(cast abi-encode 'c(address,address,address)' "$AWP_TOKEN_ADDRESS" "$STAKING_VAULT_ADDRESS" "$AWP_REGISTRY_ADDRESS")" "StakeNFT"
    vc "$DAO_ADDRESS" "src/governance/AWPDAO.sol:AWPDAO" \
        "$(cast abi-encode 'c(address,address,address,uint48,uint32,uint256)' "$STAKENFT_ADDRESS" "$AWP_TOKEN_ADDRESS" "$TREASURY_ADDRESS" 7200 50400 4)" "AWPDAO"

    # ── DEX-specific LP Manager (auto-detect chain) ──
    if [[ "$CHAIN_ID" == "56" || "$CHAIN_ID" == "97" ]]; then
        [[ -n "$POOL_MANAGER" ]] && vc "$LP_MANAGER_ADDRESS" "src/core/LPManager.sol:LPManager" \
            "$(""  # UUPS proxy, no constructor args for impl)" "LPManager (PancakeSwap)"
    else
        [[ -n "$POOL_MANAGER" ]] && vc "$LP_MANAGER_ADDRESS" "src/core/LPManagerUni.sol:LPManagerUni" \
            "$(""  # UUPS proxy, no constructor args for impl)" "LPManager (Uniswap)"
    fi

    # ── UUPS implementation contracts (no constructor args — _disableInitializers only) ──
    [[ -n "${AWP_REGISTRY_IMPL:-}" ]] && vc "$AWP_REGISTRY_IMPL" "src/AWPRegistry.sol:AWPRegistry" "" "AWPRegistry (impl)"
    [[ -n "${AWP_EMISSION_IMPL:-}" ]] && vc "$AWP_EMISSION_IMPL" "src/token/AWPEmission.sol:AWPEmission" "" "AWPEmission (impl)"
    [[ -n "${STAKING_VAULT_IMPL:-}" ]] && vc "$STAKING_VAULT_IMPL" "src/core/StakingVault.sol:StakingVault" "" "StakingVault (impl)"
    if [[ "$CHAIN_ID" == "56" || "$CHAIN_ID" == "97" ]]; then
        [[ -n "${WORKNET_MANAGER_IMPL:-}" ]] && vc "$WORKNET_MANAGER_IMPL" "src/worknets/WorknetManager.sol:WorknetManager" "" "WorknetManager (impl)"
    else
        [[ -n "${WORKNET_MANAGER_IMPL:-}" ]] && vc "$WORKNET_MANAGER_IMPL" "src/worknets/WorknetManagerUni.sol:WorknetManagerUni" "" "WorknetManagerUni (impl)"
    fi

    info "Verification complete!"
fi

# ═══════════════════════════════════════════════
#  SUMMARY
# ═══════════════════════════════════════════════

echo ""
echo "════════════════════════════════════════════════"
echo "  AWP Deployment Complete"
echo "════════════════════════════════════════════════"
echo ""
echo "  AWPRegistry:      ${AWP_REGISTRY_ADDRESS:-N/A}"
echo "  AWPToken:         ${AWP_TOKEN_ADDRESS:-N/A}"
echo "  AWPEmission:      ${AWP_EMISSION_ADDRESS:-N/A}"
echo "  StakingVault:     ${STAKING_VAULT_ADDRESS:-N/A}"
echo "  StakeNFT:         ${STAKENFT_ADDRESS:-N/A}"
echo "  WorknetNFT:        ${SUBNETNFT_ADDRESS:-N/A}"
echo "  LPManager:        ${LP_MANAGER_ADDRESS:-N/A}"
echo "  Factory:          ${FACTORY_ADDRESS:-N/A}"
echo "  AWPDAO:           ${DAO_ADDRESS:-N/A}"
echo "  Treasury:         ${TREASURY_ADDRESS:-N/A}"
echo "  SubnetMgr impl:  ${WORKNET_MANAGER_IMPL:-N/A}"
echo ""
echo "  Deploy Block:     ${DEPLOY_BLOCK:-N/A}"
echo "  Vanity Rule:      $VANITY_RULE"
echo "  API Config:       $OUTPUT_ENV"
echo ""
echo "  Next: ./scripts/deploy-api.sh"
echo "════════════════════════════════════════════════"

echo ""
echo "=== POST-DEPLOY CHECKLIST ==="
echo "1. Transfer Guardian to Safe multisig: AWPEmission.setGuardian(safeAddress)"
echo "2. Transfer Guardian on AWPRegistry: AWPRegistry.setGuardian(safeAddress)"
echo "3. Verify all contracts on block explorer: ./scripts/verify-all.sh"
echo "4. Fund relayer address with native gas tokens"
echo ""
