# Account System V2 — Design Spec

> Date: 2026-03-19
> Status: Approved
> Scope: AWPRegistry + AccessManager removal + StakingVault adaptation + tests + Go backend

## Summary

Replace the Principal/Agent mutual-exclusion account model with a tree-based binding system. Remove mandatory registration. Decouple binding from staking. Simplify delegation.

## What Changes

### AWPRegistry — New Storage & Functions

**New storage (replaces AccessManager):**
```solidity
mapping(address => address) public boundTo;           // 0 = root
mapping(address => address) public recipient;          // 0 = self (unregistered)
mapping(address => mapping(address => bool)) public delegates;
```

**New functions:**
- `register()` — convenience: `setRecipient(msg.sender)` (marks as registered)
- `bind(address target)` — set `boundTo[msg.sender] = target`, anti-cycle walk
- `unbind()` — set `boundTo[msg.sender] = address(0)`
- `setRecipient(address addr)` — set `recipient[msg.sender] = addr`
- `resolveRecipient(address addr) view` — walk boundTo chain to root, return recipient
- `isRegistered(address addr) view` — `recipient[addr] != address(0)`
- `grantDelegate(address delegate)` — `delegates[msg.sender][delegate] = true`
- `revokeDelegate(address delegate)` — `delegates[msg.sender][delegate] = false`
- `allocate(address staker, address agent, uint256 subnetId, uint256 amount)` — staker explicit
- `deallocate(address staker, address agent, uint256 subnetId, uint256 amount)` — same
- `reallocate(address staker, ...)` — same pattern

**Gasless versions:**
- `bindFor(address agent, address target, ...)` — EIP-712
- `setRecipientFor(address user, address recipient, ...)` — EIP-712
- `allocateFor(address staker, address agent, uint256 subnetId, uint256 amount, ...)` — EIP-712
- `deallocateFor(...)` — same

**Permission model for allocate/deallocate:**
```
caller == staker || delegates[staker][caller]
```

**Removed functions:**
- `register()` as mandatory step (replaced by optional `setRecipient`)
- `registerFor()` (replaced by `setRecipientFor`)
- `registerAndStake()` (user calls `deposit` + `allocate` separately)
- `removeAgent()` (agents can unbind themselves; staker deallocates)
- `setDelegation()` (replaced by `grantDelegate`/`revokeDelegate`)

**Anti-cycle check in bind():**
```solidity
function bind(address target) external {
    require(target != address(0), "use unbind()");
    // Walk up from target; if we find msg.sender, it would create a cycle
    address cur = target;
    while (boundTo[cur] != address(0) && boundTo[cur] != cur) {
        require(boundTo[cur] != msg.sender, "cycle detected");
        cur = boundTo[cur];
    }
    boundTo[msg.sender] = target;
    emit Bound(msg.sender, target);
}
```

### AccessManager — Deleted

All identity/binding/delegation logic moves into AWPRegistry. The separate contract is no longer needed.

### StakingVault — Adapted

`allocate/deallocate/reallocate` now receive `staker` as explicit parameter (previously resolved from msg.sender via AccessManager). The vault itself doesn't do permission checks — AWPRegistry handles that before calling vault.

### StakeNFT — Unchanged

Position NFTs, lock periods, deposit/withdraw, voting power — all stay as-is.

### AWPDAO — Unchanged

StakeNFT-based voting stays as-is.

## What Doesn't Change

- StakeNFT (ERC721 positions, lockDuration, depositWithPermit)
- AWPDAO (StakeNFT position-based voting)
- AWPEmission (epoch settlement, oracle allocations)
- SubnetManager / SubnetNFT / LPManager / AlphaToken / AWPToken
- Subnet registration / activation / lifecycle

## Events

```
Bound(address indexed addr, address indexed target)
Unbound(address indexed addr)
RecipientSet(address indexed addr, address recipient)
Allocated(address indexed staker, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
Deallocated(address indexed staker, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
DelegateGranted(address indexed staker, address indexed delegate)
DelegateRevoked(address indexed staker, address indexed delegate)
```

## Migration Impact

| Component | Action |
|-----------|--------|
| AWPRegistry.sol | Major rewrite — new storage, new functions, remove AccessManager calls |
| AccessManager.sol | Delete |
| IAccessManager.sol | Delete |
| StakingVault.sol | Minor — allocate/deallocate receive staker explicitly |
| Tests | Rewrite AWPRegistry tests, remove AccessManager tests |
| Go indexer | New events (Bound, Unbound, DelegateGranted, etc.), remove old events |
| Go API handlers | Adapt to new data model |
| Deploy.s.sol | Remove AccessManager deployment step |
| docs | Update architecture, API reference |
