# Staking Allocation Migration: AWPRegistry → StakingVault

**Date:** 2026-03-30
**Status:** Approved

---

## Goal

Move allocate/deallocate/reallocate (including gasless variants) from AWPRegistry to StakingVault. Reduces AWPRegistry bytecode by ~3-4KB and shortens the allocation call chain.

## Changes

### StakingVault (expand)

- Inherit `UUPSUpgradeable` + `EIP712Upgradeable` (becomes UUPS proxy)
- `allocate(staker, agent, subnetId, amount)` becomes public — checks `msg.sender == staker || isDelegate(staker, msg.sender)`
- `deallocate(staker, agent, subnetId, amount)` becomes public — same auth
- `reallocate(staker, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)` becomes public — same auth
- New `allocateFor(staker, agent, subnetId, amount, deadline, v, r, s)` — EIP-712 gasless
- New `deallocateFor(staker, agent, subnetId, amount, deadline, v, r, s)` — EIP-712 gasless
- Delegate check: `IAWPRegistry(awpRegistry).delegates(staker, msg.sender)` — cross-contract read
- Events `Allocated`, `Deallocated`, `Reallocated` emitted from StakingVault
- `_authorizeUpgrade` restricted to onlyTimelock (read treasury from AWPRegistry)
- EIP-712 domain: name="StakingVault", version="1"
- `nonces` mapping for replay protection
- Storage gap for future upgrades
- `initialize(awpRegistry_, stakeNFT_)` with initializer

### AWPRegistry (shrink)

- Remove: `allocate`, `deallocate`, `reallocate`, `allocateFor`, `deallocateFor`
- Remove: `ALLOCATE_TYPEHASH`, `DEALLOCATE_TYPEHASH` constants
- Keep: `delegates` mapping + `grantDelegate` + `revokeDelegate` (delegation is account system, stays in AWPRegistry)
- Keep: all subnet, account, and governance functions

### Deploy.s.sol

- StakingVault deployed as impl + ERC1967Proxy (like AWPRegistry, AWPEmission)
- `initialize(awpRegistry, stakeNFT)` called via proxy init data
- Add SALT_STAKING_VAULT_IMPL env var

### API relay.go

- Relay endpoints `/api/relay/allocate` and `/api/relay/deallocate` now call StakingVault (not AWPRegistry)
- EIP-712 domain changes: verifyingContract = StakingVault proxy address, name = "StakingVault"
- Relayer needs StakingVault binding

### Indexer

- Listen for `Allocated`/`Deallocated`/`Reallocated` events from StakingVault address (not AWPRegistry)

## Not Changed

- StakeNFT (deposit/withdraw)
- AWPRegistry account system (bind/setRecipient/delegate)
- AWPRegistry subnet lifecycle
- SubnetManager/LPManager
