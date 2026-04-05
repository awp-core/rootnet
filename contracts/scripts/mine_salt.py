#!/usr/bin/env python3
"""Mine CREATE2 vanity salts using all CPU cores."""

import sys
import os
import hashlib
import multiprocessing
import time
import struct

FACTORY = bytes.fromhex("4e59b44847b379578588920cA78FbF26c0B4956C")

def keccak256(data: bytes) -> bytes:
    from Crypto.Hash import keccak
    k = keccak.new(digest_bits=256)
    k.update(data)
    return k.digest()

def mine_worker(args):
    init_code_hash_hex, prefix_hex, suffix_hex, worker_id, total_workers, batch_size = args

    init_code_hash = bytes.fromhex(init_code_hash_hex)
    prefix = bytes.fromhex(prefix_hex) if prefix_hex else b""
    suffix = bytes.fromhex(suffix_hex) if suffix_hex else b""

    prefix_len = len(prefix)
    suffix_len = len(suffix)

    base = b'\xff' + FACTORY

    counter = worker_id
    checked = 0

    while True:
        salt = counter.to_bytes(32, 'big')
        data = base + salt + init_code_hash
        h = keccak256(data)
        addr = h[12:]  # last 20 bytes

        match = True
        if prefix_len > 0 and addr[:prefix_len] != prefix:
            match = False
        if match and suffix_len > 0 and addr[-suffix_len:] != suffix:
            match = False

        if match:
            return salt.hex(), addr.hex()

        counter += total_workers
        checked += 1
        if checked % batch_size == 0:
            # Check if another worker found it
            if os.path.exists("/tmp/salt_found"):
                return None, None

def main():
    if len(sys.argv) < 2:
        print("Usage: mine_salt.py <initCodeHash> [prefix_hex] [suffix_hex]")
        print("Example: mine_salt.py abc123...def 0000 00a2")
        sys.exit(1)

    init_code_hash = sys.argv[1]
    prefix = sys.argv[2] if len(sys.argv) > 2 else ""
    suffix = sys.argv[3] if len(sys.argv) > 3 else ""

    # Clean up
    if os.path.exists("/tmp/salt_found"):
        os.remove("/tmp/salt_found")

    ncpu = multiprocessing.cpu_count()
    batch_size = 100_000

    print(f"Mining with {ncpu} workers...")
    print(f"initCodeHash: {init_code_hash}")
    print(f"prefix: {prefix or '(none)'}, suffix: {suffix or '(none)'}")

    # Estimate difficulty
    bits = len(prefix) * 4 + len(suffix) * 4
    est = 2 ** bits
    print(f"Difficulty: ~{est:,} attempts (~{bits} bits)")

    start = time.time()

    args = [(init_code_hash, prefix, suffix, i, ncpu, batch_size) for i in range(ncpu)]

    with multiprocessing.Pool(ncpu) as pool:
        for result in pool.imap_unordered(mine_worker, args):
            salt, addr = result
            if salt:
                # Signal other workers
                open("/tmp/salt_found", "w").close()
                elapsed = time.time() - start
                print(f"\nFound in {elapsed:.1f}s!")
                print(f"Salt: 0x{salt}")
                print(f"Address: 0x{addr}")
                pool.terminate()
                return

    print("Not found")

if __name__ == "__main__":
    main()
