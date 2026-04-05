#!/bin/bash
set -euo pipefail

# AWP — Verify all contracts on all configured chains
# Usage:
#   ./scripts/verify-all.sh                  # Verify all chains from chains.yaml
#   ./scripts/verify-all.sh base             # Verify a specific chain
#   ETHERSCAN_API_KEY=xxx ./scripts/verify-all.sh

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

if [[ -z "${ETHERSCAN_API_KEY:-}" ]]; then
    echo "ERROR: ETHERSCAN_API_KEY not set"
    echo "Usage: ETHERSCAN_API_KEY=xxx $0 [chainName]"
    exit 1
fi

# Source the api/.env for contract addresses (same addresses on all chains)
if [[ -f "$PROJECT_DIR/api/.env" ]]; then
    set -a; source "$PROJECT_DIR/api/.env"; set +a
else
    echo "ERROR: api/.env not found. Run deploy.sh first."
    exit 1
fi

# Also source contracts/.env for deployer/guardian/vanity config
if [[ -f "$PROJECT_DIR/contracts/.env" ]]; then
    set -a; source "$PROJECT_DIR/contracts/.env"; set +a
fi

cd "$PROJECT_DIR/contracts"

verify_chain() {
    local chain_name="$1" chain_id="$2" rpc_url="$3"

    echo ""
    echo "═══════════════════════════════════════════════"
    echo "  Verifying contracts on: $chain_name (chainId=$chain_id)"
    echo "═══════════════════════════════════════════════"

    # Determine verifier URL
    local verifier_url=""
    case "$chain_id" in
        56)     verifier_url="https://api.bscscan.com/api" ;;
        97)     verifier_url="https://api-testnet.bscscan.com/api" ;;
        8453)   verifier_url="https://api.basescan.org/api" ;;
        84532)  verifier_url="https://api-sepolia.basescan.org/api" ;;
        1)      verifier_url="https://api.etherscan.io/api" ;;
        42161)  verifier_url="https://api.arbiscan.io/api" ;;
        421614) verifier_url="https://api-sepolia.arbiscan.io/api" ;;
    esac

    local vf="--etherscan-api-key $ETHERSCAN_API_KEY"
    [[ -n "$verifier_url" ]] && vf="$vf --verifier-url $verifier_url"

    vc() {
        local addr="$1" contract="$2" args="${3:-}" label="${4:-$contract}"
        [[ -z "$addr" ]] && { echo "  SKIP $label — no address"; return; }
        local cmd="forge verify-contract $addr $contract --chain-id $chain_id $vf --watch"
        [[ -n "$args" ]] && cmd="$cmd --constructor-args $args"
        echo "  Verifying $label at $addr ..."
        eval "$cmd" 2>&1 | tail -1 || echo "  WARN: $label verification failed (may already be verified)"
    }

    DEPLOYER="${DEPLOYER:-$(cast wallet address --private-key $DEPLOYER_PRIVATE_KEY 2>/dev/null || echo '')}"
    local initial_mint_wei
    initial_mint_wei=$(echo "${INITIAL_MINT:-200000000} * 1000000000000000000" | bc)

    # Standalone contracts
    vc "$AWP_TOKEN_ADDRESS" "src/token/AWPToken.sol:AWPToken" \
        "$(cast abi-encode 'c(string,string,address,uint256)' 'AWP Token' 'AWP' "$DEPLOYER" "$initial_mint_wei")" "AWPToken"
    vc "$ALPHA_FACTORY_ADDRESS" "src/token/AlphaTokenFactory.sol:AlphaTokenFactory" \
        "$(cast abi-encode 'c(address,uint64)' "$DEPLOYER" "${VANITY_RULE:-0}")" "AlphaTokenFactory"
    vc "$TREASURY_ADDRESS" "src/governance/Treasury.sol:Treasury" \
        "$(cast abi-encode 'c(uint256,address[],address[],address)' 172800 '[]' '[0x0000000000000000000000000000000000000000]' "$DEPLOYER")" "Treasury"
    vc "$WORKNETNFT_ADDRESS" "src/core/WorknetNFT.sol:WorknetNFT" \
        "$(cast abi-encode 'c(string,string,address)' 'AWP Worknet' 'AWPSUB' "$AWP_REGISTRY_ADDRESS")" "WorknetNFT"
    vc "$STAKE_NFT_ADDRESS" "src/core/StakeNFT.sol:StakeNFT" \
        "$(cast abi-encode 'c(address,address,address)' "$AWP_TOKEN_ADDRESS" "$STAKING_VAULT_ADDRESS" "$AWP_REGISTRY_ADDRESS")" "StakeNFT"
    vc "$DAO_ADDRESS" "src/governance/AWPDAO.sol:AWPDAO" \
        "$(cast abi-encode 'c(address,address,address,uint48,uint32,uint256)' "$STAKE_NFT_ADDRESS" "$AWP_TOKEN_ADDRESS" "$TREASURY_ADDRESS" 7200 50400 4)" "AWPDAO"

    # UUPS implementation contracts (no constructor args)
    [[ -n "${AWP_REGISTRY_IMPL:-}" ]] && vc "$AWP_REGISTRY_IMPL" "src/AWPRegistry.sol:AWPRegistry" "" "AWPRegistry (impl)"
    [[ -n "${AWP_EMISSION_IMPL:-}" ]] && vc "$AWP_EMISSION_IMPL" "src/token/AWPEmission.sol:AWPEmission" "" "AWPEmission (impl)"
    [[ -n "${STAKING_VAULT_IMPL:-}" ]] && vc "$STAKING_VAULT_IMPL" "src/core/StakingVault.sol:StakingVault" "" "StakingVault (impl)"

    # DEX-specific LP + WorknetManager impls (no constructor args — UUPS pattern)
    if [[ "$chain_id" == "56" || "$chain_id" == "97" ]]; then
        [[ -n "${LP_MANAGER_IMPL:-}" ]] && vc "$LP_MANAGER_IMPL" "src/core/LPManager.sol:LPManager" "" "LPManager impl (PCS)"
        [[ -n "${WORKNET_MANAGER_IMPL:-}" ]] && vc "$WORKNET_MANAGER_IMPL" "src/worknets/WorknetManager.sol:WorknetManager" "" "WorknetManager impl (PCS)"
    else
        [[ -n "${LP_MANAGER_IMPL:-}" ]] && vc "$LP_MANAGER_IMPL" "src/core/LPManagerUni.sol:LPManagerUni" "" "LPManager impl (Uni)"
        [[ -n "${WORKNET_MANAGER_IMPL:-}" ]] && vc "$WORKNET_MANAGER_IMPL" "src/worknets/WorknetManagerUni.sol:WorknetManagerUni" "" "WorknetManager impl (Uni)"
    fi

    echo "  ✓ $chain_name verification complete"
}

# Parse args
CHAIN_NAME="${1:-}"

if [[ -n "$CHAIN_NAME" ]]; then
    # Verify single chain
    if [[ ! -f "$PROJECT_DIR/chains.yaml" ]]; then
        echo "ERROR: chains.yaml not found"
        exit 1
    fi
    cfg=$(python3 -c "
import yaml, json, os
chains = yaml.safe_load(open('$PROJECT_DIR/chains.yaml'))['chains']
c = chains['$CHAIN_NAME']
c['rpcUrl'] = os.path.expandvars(c['rpcUrl'])
print(json.dumps(c))
")
    chain_id=$(echo "$cfg" | jq -r '.chainId')
    rpc_url=$(echo "$cfg" | jq -r '.rpcUrl')
    verify_chain "$CHAIN_NAME" "$chain_id" "$rpc_url"
else
    # Verify all chains
    if [[ -f "$PROJECT_DIR/chains.yaml" ]]; then
        for name in $(python3 -c "import yaml; [print(k) for k in yaml.safe_load(open('$PROJECT_DIR/chains.yaml'))['chains']]"); do
            cfg=$(python3 -c "
import yaml, json, os
chains = yaml.safe_load(open('$PROJECT_DIR/chains.yaml'))['chains']
c = chains['$name']
c['rpcUrl'] = os.path.expandvars(c['rpcUrl'])
print(json.dumps(c))
")
            chain_id=$(echo "$cfg" | jq -r '.chainId')
            rpc_url=$(echo "$cfg" | jq -r '.rpcUrl')
            verify_chain "$name" "$chain_id" "$rpc_url"
        done
    else
        # Single chain from env
        chain_id="${CHAIN_ID:-$(cast chain-id --rpc-url "$ETH_RPC_URL" 2>/dev/null || echo "8453")}"
        verify_chain "default" "$chain_id" "$ETH_RPC_URL"
    fi
fi

echo ""
echo "═══════════════════════════════════════════════"
echo "  All verifications complete"
echo "═══════════════════════════════════════════════"
