#!/bin/bash
set -euo pipefail

# ============================================================
# Deploy LPManager (stub+proxy+impl) + WorknetManager impl
# All CREATE2 — deterministic addresses
# ============================================================

source .env

# ── Common salts (same on all chains) ──
export SALT_LP_STUB=0xceeb73d115fe3d54caec9bb0a2ae836eabd9e19175c961e2c7451d04ffd58b4f
export SALT_LP_PROXY=0xa12e1d132385103a6e743bec98d2d9f04573aa13b02cf512b11f6ec7590f2715
# Predicted addresses:
# Stub:  0x00000b3C797Bd2DE8a59644Ea8974423Ff374D91
# Proxy: 0x00001961b9AcCD86b72DE19Be24FaD6f7c5b00A2

# ── Per-chain DEX addresses ──
declare -A CHAIN_RPC CHAIN_ID CHAIN_PERMIT2 CHAIN_POOL_MGR CHAIN_POS_MGR CHAIN_SWAP_ROUTER CHAIN_STATE_VIEW
declare -A CHAIN_LP_IMPL_SALT CHAIN_WM_IMPL_SALT

# Base
CHAIN_RPC[base]="$BASE_RPC_URL"
CHAIN_ID[base]=8453
CHAIN_PERMIT2[base]=0x000000000022D473030F116dDEE9F6B43aC78BA3
CHAIN_POOL_MGR[base]=0x498581fF718922c3f8e6A244956aF099B2652b2b
CHAIN_POS_MGR[base]=0x7C5f5A4bBd8fD63184577525326123B519429bDc
CHAIN_SWAP_ROUTER[base]=0x6fF5693b99212Da76ad316178A184AB56D299b43
CHAIN_STATE_VIEW[base]=0xA3c0c9b65baD0b08107Aa264b0f3dB444b867A71
CHAIN_LP_IMPL_SALT[base]=0xe3a72ef6462e11898d5b2284fe96f7eb554dda6d2359bbcbda532a74f6a19f13
CHAIN_WM_IMPL_SALT[base]=0x4bb975d0fcf6f67590c11f7cf6150a507bdb805b82a12a418786a8b8226142be

# ETH
CHAIN_RPC[eth]="$ETH_RPC_URL"
CHAIN_ID[eth]=1
CHAIN_PERMIT2[eth]=0x000000000022D473030F116dDEE9F6B43aC78BA3
CHAIN_POOL_MGR[eth]=0x000000000004444c5dc75cB358380D2e3dE08A90
CHAIN_POS_MGR[eth]=0xbD216513d74C8cf14cf4747E6AaA6420FF64ee9e
CHAIN_SWAP_ROUTER[eth]=0x66a9893cC07D91D95644AEDD05D03f95e1dBA8Af
CHAIN_STATE_VIEW[eth]=0x7fFE42C4a5DEeA5b0feC41C94C136Cf115597227
CHAIN_LP_IMPL_SALT[eth]=0x4f72194917f27e34b28d902eed0ef62d669b117465d1e30f980c8c21bcfb331a
CHAIN_WM_IMPL_SALT[eth]=0xe100d66865fda0e8f8febdcd7e2308af81048b762012be387e502567fc0591a2

# ARB
CHAIN_RPC[arb]="$ARB_RPC_URL"
CHAIN_ID[arb]=42161
CHAIN_PERMIT2[arb]=0x000000000022D473030F116dDEE9F6B43aC78BA3
CHAIN_POOL_MGR[arb]=0x360E68faCcca8cA495c1B759Fd9EEe466db9FB32
CHAIN_POS_MGR[arb]=0xd88F38F930b7952f2DB2432Cb002E7abbF3dD869
CHAIN_SWAP_ROUTER[arb]=0xa51afAF359d044F8e56fE74B9575f23142cD4B76
CHAIN_STATE_VIEW[arb]=0x76fd297e2d437cd7F76A5F2B02a5ce11c663A86e
CHAIN_LP_IMPL_SALT[arb]=0xe641edfd0f3cf21699021c69498e0450f713b66009a364ac0136990b69780a89
CHAIN_WM_IMPL_SALT[arb]=0x68bd2475d1be7b3d9a4c1d84507f1520c6968d18b72cd648d3c72e0476391257

# BSC
CHAIN_RPC[bsc]="$BSC_RPC_URL"
CHAIN_ID[bsc]=56
CHAIN_PERMIT2[bsc]=0x31c2F6fcFf4F8759b3Bd5Bf0e1084A055615c768
CHAIN_POOL_MGR[bsc]=0xa0FfB9c1CE1Fe56963B0321B32E7A0302114058b
CHAIN_POS_MGR[bsc]=0x55f4c8abA71A1e923edC303eb4fEfF14608cC226
CHAIN_SWAP_ROUTER[bsc]=0x1b81D678ffb9C0263b24A97847620C99d213eB14
CHAIN_STATE_VIEW[bsc]=""
CHAIN_LP_IMPL_SALT[bsc]=0x95c567080c0932e533b8549f35baa2b0c4eff92b1f36632c59fda6e8e610d8e3
CHAIN_WM_IMPL_SALT[bsc]=0x471fcd983c4e998fed10310ae1dde85dc843bc2caef57d98eab1d0a78d48f978

# ── Deploy function ──
deploy_chain() {
    local chain=$1
    local rpc="${CHAIN_RPC[$chain]}"
    local chainid="${CHAIN_ID[$chain]}"

    echo ""
    echo "=========================================="
    echo "  Deploying on $chain (chainId=$chainid)"
    echo "=========================================="

    export PERMIT2="${CHAIN_PERMIT2[$chain]}"
    export POOL_MANAGER="${CHAIN_POOL_MGR[$chain]}"
    export POSITION_MANAGER="${CHAIN_POS_MGR[$chain]}"
    export CL_SWAP_ROUTER="${CHAIN_SWAP_ROUTER[$chain]}"
    export STATE_VIEW="${CHAIN_STATE_VIEW[$chain]}"
    export SALT_LP_IMPL="${CHAIN_LP_IMPL_SALT[$chain]}"
    export SALT_WM_IMPL="${CHAIN_WM_IMPL_SALT[$chain]}"

    forge script script/DeployLP.s.sol \
        --rpc-url "$rpc" \
        --chain-id "$chainid" \
        --broadcast \
        --slow \
        -vvv

    echo "=== $chain deployment complete ==="
}

# ── Parse args ──
if [ $# -eq 0 ]; then
    echo "Usage: $0 <chain|--all>"
    echo "Chains: base, eth, arb, bsc"
    exit 1
fi

if [ "$1" == "--all" ]; then
    for chain in base eth arb bsc; do
        deploy_chain "$chain"
    done
else
    deploy_chain "$1"
fi
