# Account & Staking System

## Overview

The account system is split across three contracts:

- **AWPRegistry** (UUPS proxy, EIP-712 domain "AWPRegistry") — binding tree, recipient resolution, delegation
- **veAWP** (ERC721, non-upgradeable) — AWP deposit/withdraw as position NFTs
- **AWPAllocator** (UUPS proxy, EIP-712 domain "AWPAllocator") — allocation management

Every address is implicitly a root (`boundTo == address(0)`). No mandatory registration. The first call to `bind()` or `setRecipient()` increments `registeredCount` and emits `UserRegistered`.

---

## AWPRegistry Storage

```
boundTo[addr]                                   // address — parent, address(0) = root
recipient[addr]                                 // address — reward recipient, address(0) = self
delegates[staker][delegate]                     // bool — delegation authorization
nonces[addr]                                    // uint256 — EIP-712 replay protection
registeredCount                                 // uint256 — total registered users
```

## veAWP Storage

```
positions[tokenId]                              // Position { uint128 amount, uint64 lockEndTime, uint64 createdAt }
_userTotalStaked[addr]                          // uint256 — O(1) total staked balance
```

## AWPAllocator Storage

```
_allocations[staker][agent][worknetId]          // uint128 — allocation amount
_agentWorknets[staker][agent]                   // EnumerableSet.UintSet — worknet IDs with non-zero allocations
userTotalAllocated[staker]                      // uint256 — sum of all allocations
worknetTotalStake[worknetId]                    // uint256 — total stake across all stakers
nonces[addr]                                    // uint256 — EIP-712 replay protection
```

---

## Binding (AWPRegistry)

Binding and staking are fully decoupled. Bind/unbind never affects staking. Staking never affects binding.

**bind(target)**: Set `boundTo[msg.sender] = target`. Anti-cycle check walks up from target via `boundTo` chain; if `msg.sender` is encountered, reverts with `CycleDetected`. Maximum depth 256, otherwise `ChainTooDeep`. Self-bind reverts with `SelfBind`. Works regardless of current state — a root can bind, an already-bound address can switch target in a single call. First bind (when `boundTo == 0 && recipient == 0`) increments `registeredCount`.

**unbind()**: Set `boundTo[msg.sender] = address(0)`. Caller becomes root.

All chain traversals (bind anti-cycle, resolveRecipient) terminate when `boundTo[cur] == address(0)` or `boundTo[cur] == cur`.

### Gasless Binding

- **bindFor(agent, target, deadline, v, r, s)** — EIP-712 signed by agent
- **unbindFor(user, deadline, v, r, s)** — EIP-712 signed by user

---

## Recipient (AWPRegistry)

**setRecipient(addr)**: Set `recipient[msg.sender] = addr`. Any address can call at any time. The stored value only takes effect when the address is a root (`boundTo == address(0)` or `boundTo == self`). Setting recipient on a non-root stores the value but does not affect reward resolution until the address becomes root. First setRecipient (when `boundTo == 0 && recipient == 0`) increments `registeredCount`.

**resolveRecipient(addr)**: Walk up `boundTo` chain until root found. Return `recipient[root]` if non-zero, otherwise return root address itself. Maximum depth 256.

```
resolveRecipient(addr):
  cur = addr
  while boundTo[cur] != 0 && boundTo[cur] != cur:
    cur = boundTo[cur]
    if ++depth >= 256: revert ChainTooDeep
  return recipient[cur] != 0 ? recipient[cur] : cur
```

**batchResolveRecipients(addrs)**: Batch view function, resolves multiple addresses in a single call.

**isRegistered(addr)**: Returns true if `boundTo[addr] != 0 || recipient[addr] != 0`.

### Gasless Recipient

- **setRecipientFor(user, recipient, deadline, v, r, s)** — EIP-712 signed by user

---

## Delegation (AWPRegistry)

Delegates can act on behalf of a staker for allocation operations on AWPAllocator. AWPAllocator reads `AWPRegistry.delegates(staker, caller)` cross-contract.

**grantDelegate(delegate)**: `delegates[msg.sender][delegate] = true`.

**revokeDelegate(delegate)**: Require `delegate != msg.sender` (`CannotRevokeSelf`). `delegates[msg.sender][delegate] = false`.

### Gasless Delegation

- **grantDelegateFor(user, delegate, deadline, v, r, s)** — EIP-712 signed by user
- **revokeDelegateFor(user, delegate, deadline, v, r, s)** — EIP-712 signed by user

---

## Staking (veAWP)

AWP staking uses position NFTs (ERC721, symbol "sAWP"). Each position has an amount, a lock end time, and a creation timestamp. Positions are transferable.

**deposit(amount, lockDuration)**: Transfer AWP to veAWP contract. Mint a new NFT to `msg.sender`. Minimum lock duration: 1 day. Position struct: `{ uint128 amount, uint64 lockEndTime, uint64 createdAt }`. Increments `_userTotalStaked[msg.sender]`.

**depositWithPermit(amount, lockDuration, deadline, v, r, s)**: Gasless deposit using ERC-2612 permit (no prior approve tx).

**depositFor(user, amount, lockDuration)**: Only callable by AWPRegistry (used during worknet registration). AWP is transferred from the user.

**addToPosition(tokenId, amount, newLockEndTime)**: Add AWP to an existing position and/or extend the lock. Blocked if the position lock has expired (`PositionExpired`). Adding tokens resets `createdAt` to `block.timestamp` (prevents voting power manipulation). Lock can only be extended, never shortened (`LockCannotShorten`).

**withdraw(tokenId)**: Burn NFT and return AWP to owner. Requires lock expired (`LockNotExpired`). Requires remaining `_userTotalStaked` covers `AWPAllocator.userTotalAllocated` (`InsufficientUnallocated`).

### Transfer Safety

The `_update` hook maintains `_userTotalStaked` on every mint, burn, and transfer. On transfer, it verifies the sender's remaining staked amount still covers their allocations in AWPAllocator.

### Voting Power

```
votingPower(tokenId) = amount * sqrt(min(remainingTime, 54 weeks) / 7 days)
```

`totalVotingPower()` returns an upper-bound estimate: `awpToken.balanceOf(stakeNFT) * 7` (assumes all positions locked at maximum weight). Governance may adjust quorum to compensate.

---

## Allocation (AWPAllocator)

Allocation functions live on AWPAllocator (UUPS proxy, EIP-712 domain "AWPAllocator"). All operations are immediate (no cooldown).

**allocate(staker, agent, worknetId, amount)**: Caller must be staker or an authorized delegate (reads `AWPRegistry.delegates` cross-contract). Checks available balance via `veAWP.getUserTotalStaked(staker) - userTotalAllocated[staker] >= amount`. Allocations stored as `uint128`. `worknetId == 0` is rejected (`ZeroWorknetId`). Cross-chain: worknetId from any chain is accepted (globally unique, no on-chain status check). Updates:
- `_allocations[staker][agent][worknetId] += amount`
- `userTotalAllocated[staker] += amount`
- `worknetTotalStake[worknetId] += amount`
- `_agentWorknets[staker][agent].add(worknetId)`

**deallocate(staker, agent, worknetId, amount)**: Same permission check. Updates are the inverse of allocate. Removes worknetId from `_agentWorknets` when fully deallocated.

**reallocate(staker, fromAgent, fromWorknetId, toAgent, toWorknetId, amount)**: Atomic move from one (agent, worknet) pair to another. `userTotalAllocated` is unchanged (it is a move). Both source and destination `worknetId == 0` rejected.

Agent can be any address. No restriction based on binding relationship.

### Gasless Allocation

- **allocateFor(staker, agent, worknetId, amount, deadline, v, r, s)** — EIP-712 signed by staker
- **deallocateFor(staker, agent, worknetId, amount, deadline, v, r, s)** — EIP-712 signed by staker

### Query Functions

- `getAgentStake(user, agent, worknetId)` — current allocation amount
- `getAgentWorknets(user, agent)` — all worknet IDs with non-zero allocations
- `getWorknetTotalStake(worknetId)` — total stake across all stakers

---

## Events

### AWPRegistry

```
UserRegistered(address indexed user)
Bound(address indexed addr, address indexed target)
Unbound(address indexed addr)
RecipientSet(address indexed addr, address recipient)
DelegateGranted(address indexed staker, address indexed delegate)
DelegateRevoked(address indexed staker, address indexed delegate)
```

### veAWP

```
Deposited(address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndTime)
PositionIncreased(uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndTime)
Withdrawn(address indexed user, uint256 indexed tokenId, uint256 amount)
```

### AWPAllocator

```
Allocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator)
Deallocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator)
Reallocated(address indexed staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount, address operator)
```

---

## EIP-712 Domains

| Contract     | Domain Name    | Version | Nonce Endpoint              |
|-------------|----------------|---------|----------------------------|
| AWPRegistry  | "AWPRegistry"  | "1"     | `nonces(address)` on AWPRegistry  |
| AWPAllocator | "AWPAllocator" | "1"     | `nonces(address)` on AWPAllocator |

Each contract maintains its own independent nonce counter per signer.
