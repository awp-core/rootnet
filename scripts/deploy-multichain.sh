#!/bin/bash
set -euo pipefail

# AWP Multi-Chain Deployment Wrapper
# Usage:
#   ./scripts/deploy-multichain.sh <chainName>     Deploy to a specific chain
#   ./scripts/deploy-multichain.sh --all            Deploy to all chains
#   ./scripts/deploy-multichain.sh --list           List available chains

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
CHAINS_FILE="$PROJECT_DIR/chains.yaml"

# GENESIS_TIME must be set for deterministic AWPEmission proxy address
if [ -z "${GENESIS_TIME:-}" ]; then
    echo "ERROR: GENESIS_TIME must be set (Unix timestamp for emission epoch start)"
    exit 1
fi
export GENESIS_TIME

if [[ ! -f "$CHAINS_FILE" ]]; then
    echo "ERROR: $CHAINS_FILE not found"
    exit 1
fi

# Require python3 + pyyaml
if ! python3 -c "import yaml" 2>/dev/null; then
    echo "ERROR: python3 + pyyaml required. Install: pip install pyyaml"
    exit 1
fi

list_chains() {
    python3 -c "
import yaml
chains = yaml.safe_load(open('$CHAINS_FILE'))['chains']
print('Available chains:')
for name, cfg in chains.items():
    print(f'  {name:12s}  chainId={cfg[\"chainId\"]:>6}  {cfg[\"name\"]}  initialMint={cfg[\"initialMint\"]}')
"
}

deploy_chain() {
    local chain_name="$1"

    echo ""
    echo "═══════════════════════════════════════════════"
    echo "  Deploying AWP to: $chain_name"
    echo "═══════════════════════════════════════════════"

    # Read chain config from YAML
    # IMPORTANT: Do NOT use os.path.expandvars for rpcUrl — ETH_RPC_URL gets overwritten per chain.
    # Instead, read the raw ${VAR_NAME} and resolve it from the original .env file.
    local cfg
    cfg=$(ENV_FILE="$PROJECT_DIR/contracts/.env" YAML_FILE="$CHAINS_FILE" CHAIN="$chain_name" python3 << 'PYEOF'
import yaml, json, re, os

# Load original .env values (not polluted by deploy loop)
env = {}
try:
    with open(os.environ["ENV_FILE"]) as f:
        for line in f:
            line = line.strip()
            if line and not line.startswith("#") and "=" in line:
                k, v = line.split("=", 1)
                env[k.strip()] = v.strip()
except FileNotFoundError:
    pass

chains = yaml.safe_load(open(os.environ["YAML_FILE"]))["chains"]
chain_name = os.environ["CHAIN"]
if chain_name not in chains:
    import sys; print("ERROR: chain not found", file=sys.stderr); sys.exit(1)
c = chains[chain_name]
# Resolve ${VAR} references from original .env, not current shell
raw = c["rpcUrl"]
m = re.match(r'^\$\{(\w+)\}$', raw)
if m:
    c["rpcUrl"] = env.get(m.group(1), raw)
print(json.dumps(c))
PYEOF
)

    if [[ -z "$cfg" ]]; then
        echo "ERROR: Chain '$chain_name' not found in $CHAINS_FILE"
        exit 1
    fi

    # Resolve RPC URL from chain config (expandvars already done in Python above)
    local resolved_rpc=$(echo "$cfg" | jq -r '.rpcUrl')
    # If the resolved URL is empty or still a ${VAR} reference (expansion failed), try reading from original .env
    if [[ -z "$resolved_rpc" || "$resolved_rpc" == \$* ]]; then
        error "Failed to resolve rpcUrl for $chain_name. Check chains.yaml and .env."
    fi
    export ETH_RPC_URL="$resolved_rpc"
    export POOL_MANAGER=$(echo "$cfg" | jq -r '.poolManager')
    export POSITION_MANAGER=$(echo "$cfg" | jq -r '.positionManager')
    export PERMIT2=$(echo "$cfg" | jq -r '.permit2')
    export CL_SWAP_ROUTER=$(echo "$cfg" | jq -r '.swapRouter')
    export STATE_VIEW=$(echo "$cfg" | jq -r '.stateView // empty')
    export INITIAL_MINT=$(echo "$cfg" | jq -r '.initialMint')

    local dex=$(echo "$cfg" | jq -r '.dex')
    echo "Chain: $(echo "$cfg" | jq -r '.name') (chainId=$(echo "$cfg" | jq -r '.chainId'))"
    echo "DEX: $dex"
    echo "Initial mint: ${INITIAL_MINT}M AWP"
    echo "RPC: $ETH_RPC_URL"
    echo ""

    # First chain: full mine+deploy; subsequent chains reuse salts(--skip-mine)
    cd "$PROJECT_DIR"
    if [[ "${_FIRST_CHAIN_MINED:-}" == "1" ]]; then
        ./scripts/deploy.sh --skip-mine
    else
        ./scripts/deploy.sh
        export _FIRST_CHAIN_MINED=1
    fi

    cp "$PROJECT_DIR/api/.env" "$PROJECT_DIR/api/.env.${chain_name}" 2>/dev/null || true

    if [ "$dex" = "pancakeswap_v4" ]; then
        echo "NOTE: LPManager and WorknetManager addresses differ on this chain (different DEX bytecode)"
    fi

    echo ""
    echo "✓ Deployment to $chain_name complete"
    echo ""
}

# Parse args
case "${1:-}" in
    --list)
        list_chains
        ;;
    --all)
        for chain_name in $(python3 -c "import yaml; [print(k) for k in yaml.safe_load(open('$CHAINS_FILE'))['chains']]"); do
            deploy_chain "$chain_name"
        done
        echo "═══════════════════════════════════════════════"
        echo "  All chains deployed successfully"
        echo "═══════════════════════════════════════════════"
        ;;
    "")
        echo "Usage: $0 <chainName|--all|--list>"
        echo ""
        list_chains
        ;;
    *)
        deploy_chain "$1"
        ;;
esac
