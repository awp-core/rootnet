# Account & Staking System

## Storage

```
boundTo[addr]                                   // address — parent, 0 = root
recipient[addr]                                 // address — reward recipient, 0 = self
delegates[staker][delegate]                     // bool
staked[staker]                                  // uint256
committed[staker]                               // uint256
allocations[staker][worker][subnetId]           // uint256
workerStake[worker][subnetId]                   // uint256
```

Every address is implicitly a root (boundTo == 0). No registration required.

## Binding

Binding and staking are fully decoupled. Bind/unbind never affects staking. Staking never affects binding.

**bind(target)**: Set boundTo[msg.sender] = target. Anti-cycle check: walk up from target via boundTo chain until boundTo == 0; if msg.sender is encountered, revert. Works regardless of current state — root can bind, already-bound can switch to a different target in a single call.

**unbind()**: Set boundTo[msg.sender] = 0. Caller becomes root.

All chain traversals (bind anti-cycle, resolveRecipient) terminate when boundTo[cur] == 0 or boundTo[cur] == cur.

## Recipient

**setRecipient(addr)**: Set recipient[msg.sender] = addr. Any address can call at any time. The stored value only takes effect when the address is a root (boundTo == 0 or boundTo == self). Setting recipient on a non-root address stores the value but does not affect reward resolution until the address becomes root.

**resolveRecipient(addr)**: Walk up boundTo chain until root found. Return recipient[root] if non-zero, otherwise return root address itself.

```
resolveRecipient(addr):
  cur = addr
  while boundTo[cur] != 0 && boundTo[cur] != cur:
    cur = boundTo[cur]
  return recipient[cur] != 0 ? recipient[cur] : cur
```

## Staking

Any address can stake. Staking is independent of binding.

**deposit(amount)**: Transfer AWP to vault. staked[msg.sender] += amount.

**requestWithdraw(amount)**: Require staked[msg.sender] - committed[msg.sender] >= amount. Start cooldown.

**withdraw()**: After cooldown, transfer AWP back.

## Allocation

**allocate(staker, worker, subnetId, amount)**: Caller must be staker or delegates[staker][caller]. Require staked[staker] - committed[staker] >= amount. Update:
- allocations[staker][worker][subnetId] += amount
- committed[staker] += amount
- workerStake[worker][subnetId] += amount

**deallocate(staker, worker, subnetId, amount)**: Same permission check. Require allocations[staker][worker][subnetId] >= amount. Update:
- allocations[staker][worker][subnetId] -= amount
- committed[staker] -= amount
- workerStake[worker][subnetId] -= amount

Worker can be any address. No restriction based on binding relationship.

## Delegation

**grantDelegate(delegate)**: delegates[msg.sender][delegate] = true.

**revokeDelegate(delegate)**: Require delegate != msg.sender. delegates[msg.sender][delegate] = false.

## Coordinator Query

**workerStake[worker][subnetId]** provides O(1) lookup for a worker's total stake on a subnet. This is the denormalized sum of allocations[*][worker][subnetId] across all stakers, maintained incrementally by allocate/deallocate.

**resolveRecipient(worker)** provides the reward destination for any worker address.

## Events

```
Bound(address indexed addr, address indexed target)
Unbound(address indexed addr)
RecipientSet(address indexed addr, address recipient)
Deposited(address indexed staker, uint256 amount)
WithdrawRequested(address indexed staker, uint256 amount)
Withdrawn(address indexed staker, uint256 amount)
Allocated(address indexed staker, address indexed worker, uint256 indexed subnetId, uint256 amount, address operator)
Deallocated(address indexed staker, address indexed worker, uint256 indexed subnetId, uint256 amount, address operator)
DelegateGranted(address indexed staker, address indexed delegate)
DelegateRevoked(address indexed staker, address indexed delegate)
```
