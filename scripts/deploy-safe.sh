#!/usr/bin/env bash
# Deploy a Safe multisig wallet with vanity address support.
# Deploys to multiple chains with identical address (same salt + same params = same CREATE2 address).
#
# Usage:
#   ./scripts/deploy-safe.sh --mine                    # Mine vanity salt only
#   ./scripts/deploy-safe.sh --deploy base,bsc         # Deploy to specific chains
#   ./scripts/deploy-safe.sh --deploy all              # Deploy to all chains in chains.yaml
#   ./scripts/deploy-safe.sh --predict                 # Predict address without deploying
#
# Configuration (env vars or contracts/.env):
#   SAFE_OWNERS="0xAlice,0xBob,0xCharlie"   # Comma-separated owner addresses
#   SAFE_THRESHOLD=2                         # Required signatures (default: 2)
#   SAFE_VANITY_PREFIX="0000"               # Vanity prefix (optional)
#   SAFE_VANITY_SUFFIX=""                   # Vanity suffix (optional)
#   SAFE_SALT_NONCE=0                       # Salt nonce (0 = mine for vanity, or use pre-mined value)
#   DEPLOYER_PRIVATE_KEY=0x...              # Key to pay gas for deployment tx

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# ─── Colors ───
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*" >&2; exit 1; }

# ─── Safe official deployed addresses (same on all EVM chains) ───
SAFE_PROXY_FACTORY="0x4e1DCf7AD4e460CfD30791CCC4F9c8a4f820ec67"
SAFE_SINGLETON="0x41675C099F32341bf84BFc5382aF534df5C7461a"

# ─── Load env ───
[[ -f "$ROOT_DIR/contracts/.env" ]] && { set -a; source "$ROOT_DIR/contracts/.env"; set +a; }

OWNERS="${SAFE_OWNERS:-}"
THRESHOLD="${SAFE_THRESHOLD:-2}"
VANITY_PREFIX="${SAFE_VANITY_PREFIX:-}"
VANITY_SUFFIX="${SAFE_VANITY_SUFFIX:-}"
SALT_NONCE="${SAFE_SALT_NONCE:-0}"

# ─── Validate ───
command -v cast >/dev/null 2>&1 || error "cast (Foundry) not found"
command -v python3 >/dev/null 2>&1 || error "python3 not found"
[[ -n "$OWNERS" ]] || error "SAFE_OWNERS not set (comma-separated addresses)"
[[ -n "$DEPLOYER_PRIVATE_KEY" ]] || error "DEPLOYER_PRIVATE_KEY not set"

# Parse owners into array
IFS=',' read -ra OWNER_ARRAY <<< "$OWNERS"
OWNER_COUNT=${#OWNER_ARRAY[@]}
[[ "$OWNER_COUNT" -ge 1 ]] || error "At least 1 owner required"
[[ "$THRESHOLD" -ge 1 && "$THRESHOLD" -le "$OWNER_COUNT" ]] || error "Invalid threshold: $THRESHOLD (owners: $OWNER_COUNT)"

# ─── Build initializer calldata ───
# Safe.setup(owners[], threshold, to, data, fallbackHandler, paymentToken, payment, paymentReceiver)
ZERO="0x0000000000000000000000000000000000000000"

# Format owners as Solidity array
OWNERS_SOL="[$(printf '%s,' "${OWNER_ARRAY[@]}" | sed 's/,$//' )]"

INITIALIZER=$(cast calldata \
    "setup(address[],uint256,address,bytes,address,address,uint256,address)" \
    "$OWNERS_SOL" \
    "$THRESHOLD" \
    "$ZERO" "0x" "$ZERO" "$ZERO" 0 "$ZERO")

info "Owners: ${OWNER_ARRAY[*]}"
info "Threshold: $THRESHOLD-of-$OWNER_COUNT"

# ─── Compute init code hash for CREATE2 prediction ───
# SafeProxy creation code = SafeProxy.creationCode + abi.encode(singleton)
# The factory uses: keccak256(abi.encodePacked(keccak256(initializer), saltNonce)) as salt
compute_address() {
    local nonce="$1"
    # The factory's CREATE2 salt = keccak256(keccak256(initializer) ++ uint256(nonce))
    local init_hash
    init_hash=$(cast keccak256 "$INITIALIZER")
    local packed_salt
    packed_salt=$(cast keccak256 "$(cast abi-encode 'f(bytes32,uint256)' "$init_hash" "$nonce")")

    # SafeProxy initcode = type(SafeProxy).creationCode ++ abi.encode(singleton)
    # We need the proxy bytecode — get it from the factory
    # Simpler: use cast create2 with the known init code hash
    # For Safe v1.4.1, the proxy init code hash with singleton is deterministic
    # We can compute the address directly:
    local proxy_init_code_hash
    proxy_init_code_hash=$(cast keccak256 "$(cast code $SAFE_PROXY_FACTORY --rpc-url https://eth.llamarpc.com 2>/dev/null | head -c 10)") 2>/dev/null || true

    # Simpler approach: call the factory's view function to predict
    # Unfortunately SafeProxyFactory doesn't have a pure predict function
    # Use the CREATE2 formula directly
    python3 -c "
from eth_abi import encode
from web3 import Web3

factory = '$SAFE_PROXY_FACTORY'
singleton = '$SAFE_SINGLETON'
initializer = bytes.fromhex('${INITIALIZER#0x}')
nonce = $nonce

# Compute salt the way SafeProxyFactory does
init_hash = Web3.keccak(initializer)
salt = Web3.keccak(encode(['bytes32', 'uint256'], [init_hash, nonce]))

# SafeProxy creation code (constant for v1.4.1)
# The proxy bytecode is: creation_code + abi.encode(singleton)
# SafeProxy creation code hash for singleton 0x41675C099F32341bf84BFc5382aF534df5C7461a
# This is a known constant for Safe v1.4.1
proxy_creation_code = bytes.fromhex('608060405234801561001057600080fd5b5060405161017138038061017183398101604081905261002f91610055565b600080546001600160a01b0319166001600160a01b0392909216919091179055610085565b60006020828403121561006757600080fd5b81516001600160a01b038116811461007e57600080fd5b9392505050565b60de806100936000396000f3fe6080604052600073ffffffffffffffffffffffffffffffffffffffff8154167fa619486e00000000000000000000000000000000000000000000000000000000823503606857005b3660008037600080366000845af43d6000803e80801560875760003df35b3d6000fd5b') + encode(['address'], [Web3.to_checksum_address(singleton)])
init_code_hash = Web3.keccak(proxy_creation_code)

# CREATE2 address
addr_bytes = Web3.keccak(
    b'\xff' + bytes.fromhex(factory[2:]) + salt + init_code_hash
)
address = Web3.to_checksum_address('0x' + addr_bytes.hex()[-40:])
print(address)
" 2>/dev/null
}

# ─── Simpler approach: use cast to simulate ───
predict_address() {
    local nonce="$1"
    local rpc="${2:-https://eth.llamarpc.com}"

    # Call factory's createProxyWithNonce as a static call to get the return address
    cast call "$SAFE_PROXY_FACTORY" \
        "createProxyWithNonce(address,bytes,uint256)(address)" \
        "$SAFE_SINGLETON" "$INITIALIZER" "$nonce" \
        --rpc-url "$rpc" 2>/dev/null || echo "FAILED"
}

# ─── Mine vanity salt ───
mine_vanity() {
    if [[ -z "$VANITY_PREFIX" && -z "$VANITY_SUFFIX" ]]; then
        info "No vanity rules set, using saltNonce=$SALT_NONCE"
        return
    fi

    info "Mining vanity address (prefix=$VANITY_PREFIX, suffix=$VANITY_SUFFIX)..."
    info "Using $(nproc) threads"

    local rpc="https://eth.llamarpc.com"
    local best_nonce=0
    local found=0

    # Parallel mining using background jobs
    local num_threads=$(nproc)
    local batch_size=1000

    for batch_start in $(seq 0 $batch_size 1000000); do
        if [[ "$found" -eq 1 ]]; then break; fi

        for thread in $(seq 0 $((num_threads - 1))); do
            if [[ "$found" -eq 1 ]]; then break; fi
            (
                for nonce in $(seq $((batch_start + thread)) $num_threads $((batch_start + batch_size - 1))); do
                    addr=$(predict_address "$nonce" "$rpc" 2>/dev/null)
                    [[ "$addr" == "FAILED" || -z "$addr" ]] && continue

                    addr_lower=$(echo "${addr#0x}" | tr '[:upper:]' '[:lower:]')

                    prefix_ok=1
                    suffix_ok=1
                    [[ -n "$VANITY_PREFIX" ]] && [[ "${addr_lower:0:${#VANITY_PREFIX}}" != "$VANITY_PREFIX" ]] && prefix_ok=0
                    [[ -n "$VANITY_SUFFIX" ]] && [[ "${addr_lower: -${#VANITY_SUFFIX}}" != "$VANITY_SUFFIX" ]] && suffix_ok=0

                    if [[ "$prefix_ok" -eq 1 && "$suffix_ok" -eq 1 ]]; then
                        echo "FOUND:$nonce:$addr" > /tmp/safe_vanity_result
                        exit 0
                    fi
                done
            ) &
        done
        wait

        if [[ -f /tmp/safe_vanity_result ]]; then
            local result
            result=$(cat /tmp/safe_vanity_result)
            rm -f /tmp/safe_vanity_result
            SALT_NONCE=$(echo "$result" | cut -d: -f2)
            local addr=$(echo "$result" | cut -d: -f3)
            info "Found vanity address!"
            info "  Salt nonce: $SALT_NONCE"
            info "  Address:    $addr"
            found=1
            break
        fi

        echo -ne "\r  Searched $((batch_start + batch_size)) nonces..."
    done

    if [[ "$found" -eq 0 ]]; then
        warn "Vanity address not found in 1M attempts. Using saltNonce=0"
        SALT_NONCE=0
    fi
}

# ─── Deploy to a single chain ───
deploy_chain() {
    local chain_name="$1"
    local rpc_url="$2"

    local chain_id
    chain_id=$(cast chain-id --rpc-url "$rpc_url" 2>/dev/null) || error "Cannot connect to $chain_name RPC"

    info "Deploying Safe to $chain_name (chainId=$chain_id)..."

    # Check if already deployed
    local predicted
    predicted=$(predict_address "$SALT_NONCE" "$rpc_url")
    local existing_code
    existing_code=$(cast code "$predicted" --rpc-url "$rpc_url" 2>/dev/null)
    if [[ ${#existing_code} -gt 4 ]]; then
        info "  Already deployed at $predicted — skipping"
        return
    fi

    # Deploy
    local tx_hash
    tx_hash=$(cast send "$SAFE_PROXY_FACTORY" \
        "createProxyWithNonce(address,bytes,uint256)(address)" \
        "$SAFE_SINGLETON" "$INITIALIZER" "$SALT_NONCE" \
        --private-key "$DEPLOYER_PRIVATE_KEY" \
        --rpc-url "$rpc_url" \
        --json 2>/dev/null | python3 -c "import sys,json; print(json.load(sys.stdin).get('transactionHash',''))" 2>/dev/null) || true

    if [[ -n "$tx_hash" ]]; then
        info "  TX: $tx_hash"
        # Wait for confirmation
        cast receipt "$tx_hash" --rpc-url "$rpc_url" --confirmations 1 >/dev/null 2>&1 || true
    fi

    # Verify
    local code
    code=$(cast code "$predicted" --rpc-url "$rpc_url" 2>/dev/null)
    if [[ ${#code} -gt 4 ]]; then
        echo -e "  ${GREEN}✓ Safe deployed: $predicted${NC}"
    else
        echo -e "  ${RED}✗ Deployment failed${NC}"
    fi
}

# ─── Resolve chain RPC URLs from chains.yaml ───
get_chain_rpcs() {
    local chain_filter="$1"
    python3 -c "
import yaml, os, sys, json
with open('$ROOT_DIR/chains.yaml') as f:
    chains = yaml.safe_load(f)['chains']
result = []
for name, cfg in chains.items():
    if '$chain_filter' != 'all' and name not in '$chain_filter'.split(','):
        continue
    rpc = os.path.expandvars(cfg['rpcUrl'])
    if rpc:
        result.append({'name': name, 'rpc': rpc})
json.dump(result, sys.stdout)
" 2>/dev/null
}

# ─── Main ───
echo ""
echo "════════════════════════════════════════════"
echo "  Safe Multisig Deployer"
echo "════════════════════════════════════════════"
echo ""

case "${1:-}" in
    --mine)
        mine_vanity
        echo ""
        echo "Add to contracts/.env:"
        echo "  SAFE_SALT_NONCE=$SALT_NONCE"
        ;;

    --predict)
        rpc="${2:-https://eth.llamarpc.com}"
        addr=$(predict_address "$SALT_NONCE" "$rpc")
        echo "Predicted Safe address: $addr"
        echo "  Owners:    ${OWNER_ARRAY[*]}"
        echo "  Threshold: $THRESHOLD-of-$OWNER_COUNT"
        echo "  Nonce:     $SALT_NONCE"
        ;;

    --deploy)
        chains="${2:-all}"

        # Predict address first
        first_rpc=$(get_chain_rpcs "$chains" | python3 -c "import sys,json; print(json.load(sys.stdin)[0]['rpc'])" 2>/dev/null)
        PREDICTED=$(predict_address "$SALT_NONCE" "$first_rpc")
        info "Safe address will be: $PREDICTED"
        info "  Owners: ${OWNER_ARRAY[*]}"
        info "  Threshold: $THRESHOLD-of-$OWNER_COUNT"
        info "  Salt nonce: $SALT_NONCE"
        echo ""

        # Deploy to each chain
        chains_json=$(get_chain_rpcs "$chains")
        count=$(echo "$chains_json" | python3 -c "import sys,json; print(len(json.load(sys.stdin)))")

        for i in $(seq 0 $((count - 1))); do
            name=$(echo "$chains_json" | python3 -c "import sys,json; print(json.load(sys.stdin)[$i]['name'])")
            rpc=$(echo "$chains_json" | python3 -c "import sys,json; print(json.load(sys.stdin)[$i]['rpc'])")
            deploy_chain "$name" "$rpc"
        done

        echo ""
        echo "════════════════════════════════════════════"
        echo "  Safe deployed: $PREDICTED"
        echo "════════════════════════════════════════════"
        echo ""
        echo "Next steps:"
        echo "  1. Open https://app.safe.global → 'Add existing Safe' → paste $PREDICTED"
        echo "  2. Transfer Guardian:"
        echo "     cast send \$AWP_EMISSION 'setGuardian(address)' $PREDICTED --private-key \$KEY --rpc-url \$RPC"
        echo "     cast send \$AWP_REGISTRY 'setGuardian(address)' $PREDICTED --private-key \$KEY --rpc-url \$RPC"
        echo ""
        ;;

    *)
        echo "Usage:"
        echo "  $0 --mine                  Mine vanity salt nonce"
        echo "  $0 --predict               Show predicted address"
        echo "  $0 --deploy base,bsc       Deploy to specific chains"
        echo "  $0 --deploy all            Deploy to all chains"
        echo ""
        echo "Config (env vars or contracts/.env):"
        echo "  SAFE_OWNERS=0xA,0xB,0xC    Owner addresses (required)"
        echo "  SAFE_THRESHOLD=2           Signatures required (default: 2)"
        echo "  SAFE_VANITY_PREFIX=0000    Vanity prefix (optional)"
        echo "  SAFE_VANITY_SUFFIX=        Vanity suffix (optional)"
        echo "  SAFE_SALT_NONCE=0          Pre-mined nonce (0 = auto)"
        ;;
esac
