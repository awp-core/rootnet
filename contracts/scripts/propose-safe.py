#!/usr/bin/env python3
"""Submit Safe multisig proposals via the Safe Transaction Service API."""

import json
import os
import sys
import subprocess
import requests

SAFE = "0x000002bEfa6A1C99A710862Feb6dB50525dF00A3"
REGISTRY = "0x0000F34Ed3594F54faABbCb2Ec45738DDD1c001A"
LP_PROXY = "0x00001961b9AcCD86b72DE19Be24FaD6f7c5b00A2"

# WorknetManager impl per chain
WM_IMPLS = {
    8453: "0x00945e7fd4110b9c56ab4a3c2f53b6fabe6485e5",  # Base
    1:    "0x0029aABD49BF9ec7a34CDbcf75486B19CFAC3Ea8",  # ETH
    42161:"0x00c428DCa1678e41Ed17Cc5AE3cF14430e2085A0",  # ARB
    56:   "0x00D87f2f81E20cB1583F46d94BC7a7ad8f2DAC78",  # BSC
}

# Safe TX nonces per chain (replacing existing proposals)
NONCES = {
    8453: {"setLPManager": 8, "setWorknetManagerImpl": 9},
    1:    {"setLPManager": 8, "setWorknetManagerImpl": 9},
    42161:{"setLPManager": 8, "setWorknetManagerImpl": 9},
    56:   {"setLPManager": 11, "setWorknetManagerImpl": 12},
}

API_URLS = {
    8453: "https://safe-transaction-base.safe.global",
    1:    "https://safe-transaction-mainnet.safe.global",
    42161:"https://safe-transaction-arbitrum.safe.global",
    56:   "https://safe-transaction-bsc.safe.global",
}

CHAIN_NAMES = {8453: "Base", 1: "ETH", 42161: "ARB", 56: "BSC"}


def encode_calldata(func_sig: str, *args) -> str:
    """Use cast to encode calldata."""
    cmd = ["cast", "calldata", func_sig] + list(args)
    return subprocess.check_output(cmd).decode().strip()


def get_safe_tx_hash(chain_id: int, rpc: str, to: str, data: str, nonce: int) -> str:
    """Get the Safe transaction hash using the Safe contract."""
    cmd = [
        "cast", "call", SAFE,
        "getTransactionHash(address,uint256,bytes,uint8,uint256,uint256,uint256,address,address,uint256)(bytes32)",
        to,           # to
        "0",          # value
        data,         # data
        "0",          # operation (CALL)
        "0",          # safeTxGas
        "0",          # baseGas
        "0",          # gasPrice
        "0x0000000000000000000000000000000000000000",  # gasToken
        "0x0000000000000000000000000000000000000000",  # refundReceiver
        str(nonce),   # nonce
        "--rpc-url", rpc
    ]
    return subprocess.check_output(cmd).decode().strip()


def sign_hash(tx_hash: str, private_key: str) -> str:
    """Sign a hash with cast wallet sign."""
    cmd = ["cast", "wallet", "sign", "--private-key", private_key, "--no-hash", tx_hash]
    return subprocess.check_output(cmd).decode().strip()


def propose_tx(api_url: str, sender: str, to: str, data: str, nonce: int, signature: str, safe_tx_hash: str):
    """Submit proposal to Safe Transaction Service."""
    url = f"{api_url}/api/v1/safes/{SAFE}/multisig-transactions/"
    payload = {
        "to": to,
        "value": "0",
        "data": data,
        "operation": 0,
        "safeTxGas": "0",
        "baseGas": "0",
        "gasPrice": "0",
        "gasToken": "0x0000000000000000000000000000000000000000",
        "refundReceiver": "0x0000000000000000000000000000000000000000",
        "nonce": nonce,
        "contractTransactionHash": safe_tx_hash,
        "sender": sender,
        "signature": signature,
        "origin": "AWP Deploy Script"
    }
    resp = requests.post(url, json=payload)
    return resp.status_code, resp.text


def main():
    private_key = os.environ.get("DEPLOYER_PRIVATE_KEY")
    if not private_key:
        print("Error: DEPLOYER_PRIVATE_KEY not set")
        sys.exit(1)

    sender = subprocess.check_output(
        ["cast", "wallet", "address", private_key]
    ).decode().strip()
    print(f"Sender (delegate): {sender}")

    rpcs = {
        8453: os.environ["BASE_RPC_URL"],
        1:    os.environ["ETH_RPC_URL"],
        42161:os.environ["ARB_RPC_URL"],
        56:   os.environ["BSC_RPC_URL"],
    }

    # Encode calldatas
    setLPManager_data = encode_calldata("setLPManager(address)", LP_PROXY)
    print(f"setLPManager calldata: {setLPManager_data[:20]}...")

    for chain_id, chain_name in CHAIN_NAMES.items():
        print(f"\n{'='*50}")
        print(f"  {chain_name} (chainId={chain_id})")
        print(f"{'='*50}")

        rpc = rpcs[chain_id]
        api_url = API_URLS[chain_id]
        nonces = NONCES[chain_id]
        wm_impl = WM_IMPLS[chain_id]

        setWMImpl_data = encode_calldata("setWorknetManagerImpl(address)", wm_impl)

        proposals = [
            ("setLPManager", nonces["setLPManager"], setLPManager_data),
            ("setWorknetManagerImpl", nonces["setWorknetManagerImpl"], setWMImpl_data),
        ]

        for name, nonce, data in proposals:
            print(f"\n  --- {name} (nonce={nonce}) ---")

            # Get Safe tx hash
            safe_tx_hash = get_safe_tx_hash(chain_id, rpc, REGISTRY, data, nonce)
            print(f"  safeTxHash: {safe_tx_hash}")

            # Sign
            sig = sign_hash(safe_tx_hash, private_key)
            print(f"  signature: {sig[:20]}...")

            # Submit
            status, resp_text = propose_tx(api_url, sender, REGISTRY, data, nonce, sig, safe_tx_hash)
            if status in (200, 201):
                print(f"  OK (status={status})")
            else:
                print(f"  FAILED (status={status}): {resp_text[:200]}")


if __name__ == "__main__":
    main()
