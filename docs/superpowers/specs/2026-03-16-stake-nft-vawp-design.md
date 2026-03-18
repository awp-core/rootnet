# StakeNFT + vAWP Voting + StakingVault Refactor — Design Spec

## Overview

Replace the current StakingVault deposit/withdraw/cooldown system with an NFT-based staking model. Each stake position is an ERC721 token with locked amount and expiry. Voting power is derived from NFT positions using `amount * sqrt(min(remainingEpochs, 54))`. The DAO Governor uses a custom counting module that accepts tokenId arrays. StakingVault is simplified by removing STP, pending dual-slot mechanism, and the Allocation struct.

## Design Decisions

| Decision | Choice |
|----------|--------|
| Multiple positions per user | Yes — each deposit creates a new NFT |
| Add to existing position | Yes — addToPosition with optional lock extension (only extend, not shorten) |
| Withdraw | Only after lock expiry, burns NFT, full amount |
| Lock duration cap | None (but voting power caps at √54) |
| Epoch duration | 1 week |
| Lock start | Current epoch does not count; lockEndEpoch = currentEpoch + 1 + lockEpochs |
| Minimum lock | 1 epoch (= 1 week effective) |
| NFT transferable | Yes — transferring transfers the locked AWP and future voting power |
| Transfer restriction | Sender's remaining NFT total >= userTotalAllocated (else must deallocate first) |
| Allocation model | Method C — allocations bound to user, NFT only manages amount + lock |
| Allocation timing | Immediate — no pending/dual-slot. allocate/deallocate/reallocate all instant |
| Voting mechanism | No delegate, no checkpoint. Vote by submitting tokenId[] with params |
| Voting eligibility | ownerOf(tokenId) == voter AND remainingEpochs >= 1 AND minted before proposal |
| Voting power | amount * sqrt(min(remainingEpochs, 54)) |
| Double-vote prevention | mapping(proposalId => mapping(tokenId => bool)) |
| Governor | OZ Governor with custom counting module (GovernorVotes removed) |
| vAWP token | None — StakeNFT provides voting power calculation directly |
| STP | Deleted — coordinator calculates from events |

## New Contract: StakeNFT

**Inherits:** ERC721Enumerable

**Storage:**
```solidity
struct Position {
    uint128 amount;       // Staked AWP
    uint64  lockEndEpoch; // Unlock epoch (inclusive: withdrawable when currentEpoch >= lockEndEpoch)
    uint64  createdEpoch; // Epoch when position was minted (for proposal eligibility)
}

mapping(uint256 => Position) public positions;
uint256 private _nextTokenId = 1;

// O(1) balance tracking — avoids ERC721Enumerable iteration for balance checks
mapping(address => uint256) private _userTotalStaked;

IERC20 public immutable awpToken;
address public immutable stakingVault;
address public immutable awpEmission;  // for currentEpoch()

uint64 public constant MAX_WEIGHT_EPOCHS = 54;
uint64 public constant MIN_LOCK_EPOCHS = 1;
```

**Key design: `_userTotalStaked` accumulator**

Instead of iterating ERC721Enumerable to sum positions (O(n)), maintain a `_userTotalStaked` mapping updated on deposit/withdraw/transfer. This provides O(1) reads for allocation checks and transfer hooks.

- `deposit`: `_userTotalStaked[msg.sender] += amount`
- `withdraw`: `_userTotalStaked[msg.sender] -= position.amount`
- `transfer`: `_userTotalStaked[from] -= position.amount; _userTotalStaked[to] += position.amount`

**Functions:**

### deposit(uint256 amount, uint64 lockEpochs)
- `require(lockEpochs >= MIN_LOCK_EPOCHS)`
- `require(amount > 0)`
- Transfer AWP from msg.sender to this contract
- `lockEndEpoch = currentEpoch() + 1 + lockEpochs` (current epoch doesn't count)
- `createdEpoch = currentEpoch()`
- Mint NFT to msg.sender
- Store Position
- `_userTotalStaked[msg.sender] += amount`

### addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndEpoch)
- `require(ownerOf(tokenId) == msg.sender)`
- If amount > 0: transfer AWP, increase position.amount, `_userTotalStaked[msg.sender] += amount`
- If newLockEndEpoch > 0:
  - `require(newLockEndEpoch >= position.lockEndEpoch)` (only extend, not shorten)
  - `require(newLockEndEpoch > currentEpoch())` (must extend beyond current epoch)
  - `position.lockEndEpoch = newLockEndEpoch`
- At least one of amount or newLockEndEpoch must change

### withdraw(uint256 tokenId)
- `require(ownerOf(tokenId) == msg.sender)`
- `require(currentEpoch() >= position.lockEndEpoch)`
- **Check BEFORE burn**: `_userTotalStaked[msg.sender] - position.amount >= stakingVault.userTotalAllocated(msg.sender)`
- `_userTotalStaked[msg.sender] -= position.amount`
- Transfer AWP to msg.sender
- Burn NFT
- Delete position

### getUserTotalStaked(address user) → uint256
- O(1) read: `return _userTotalStaked[user]`

### getVotingPower(uint256 tokenId) → uint256
- `remaining = lockEndEpoch > currentEpoch ? lockEndEpoch - currentEpoch : 0`
- `effective = min(remaining, MAX_WEIGHT_EPOCHS)`
- `return amount * sqrt(effective)`
- Uses OpenZeppelin Math.sqrt

### getUserVotingPower(address user) → uint256
- Sum getVotingPower for all user's tokens (uses ERC721Enumerable iteration — OK for view functions)

### totalVotingPower() → uint256
- **Approximation for quorum**: `totalStakedAWP * sqrt(MAX_WEIGHT_EPOCHS)` (upper bound)
- Where `totalStakedAWP = awpToken.balanceOf(address(this))`
- This is an upper bound (assumes all positions locked at max), acceptable for quorum denominator
- Actual quorum requirement = `totalVotingPower() * quorumPercent / 100`

### remainingEpochs(uint256 tokenId) → uint64
- `lockEndEpoch > currentEpoch ? lockEndEpoch - currentEpoch : 0`

### Transfer hook (_update override)
```solidity
function _update(address to, uint256 tokenId, address auth) internal override returns (address from) {
    from = super._update(to, tokenId, auth);

    if (from != address(0) && to != address(0)) {
        // Transfer (not mint/burn): update accumulators and check allocation coverage
        uint128 amt = positions[tokenId].amount;
        _userTotalStaked[from] -= amt;
        _userTotalStaked[to] += amt;
        require(_userTotalStaked[from] >= stakingVault.userTotalAllocated(from), "InsufficientUnallocated");
    }
    // Note: mint and burn are handled in deposit() and withdraw() directly
}
```

**Important**: In OZ 5.x ERC721, `_update` is called during `_mint`, `_burn`, and `_transfer`. The `from != address(0) && to != address(0)` guard ensures the hook only runs for transfers. Mint/burn accumulator updates happen in `deposit()`/`withdraw()` instead, avoiding double-counting.

## StakingVault Changes

### Deleted:
- `deposit()`, `withdraw()`, `requestWithdraw()`, `cancelWithdraw()` — moved to StakeNFT
- `userBalance` mapping — replaced by StakeNFT.getUserTotalStaked() (O(1))
- `WithdrawRequest` struct, `withdrawRequests` mapping
- `cooldownPeriod`, `setCooldownPeriod()`
- `SubnetSTP` struct, `subnetSTP` mapping, `_updateSTP()`, `getSubnetSTP()`
- `Allocation` struct — replaced by plain `uint128`
- `pendingSub`, `pendingAdd`, `pendingMark` fields — removed (no more pending)
- `_resolve()`, `_effectiveAmount()` — removed
- `freezeAgentAllocations` force-resolve logic — simplified to direct zero

### Storage:
```solidity
mapping(address => mapping(address => mapping(uint256 => uint128))) private _allocations;

mapping(address => uint256) public userTotalAllocated;  // kept
mapping(uint256 => uint256) public subnetTotalStake;    // kept

address public immutable rootNet;
address public stakeNFT;  // reference to StakeNFT for balance checks
```

### allocate(user, agent, subnetId, amount)
- `require(amount > 0 && amount <= type(uint128).max)`
- `available = StakeNFT.getUserTotalStaked(user) - userTotalAllocated[user]` (O(1) call)
- `require(available >= amount)`
- `_allocations[user][agent][subnetId] += uint128(amount)`
- `userTotalAllocated[user] += amount`
- `subnetTotalStake[subnetId] += amount`

### deallocate(user, agent, subnetId, amount)
- `require(amount > 0 && amount <= type(uint128).max)`
- `require(_allocations[user][agent][subnetId] >= uint128(amount))`
- `_allocations[user][agent][subnetId] -= uint128(amount)`
- `userTotalAllocated[user] -= amount`
- `subnetTotalStake[subnetId] -= amount`

### reallocate(user, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)
- Immediate atomic operation (no pending)
- `require(_allocations[user][fromAgent][fromSubnetId] >= uint128(amount))`
- `_allocations[user][fromAgent][fromSubnetId] -= uint128(amount)`
- `_allocations[user][toAgent][toSubnetId] += uint128(amount)`
- `subnetTotalStake[fromSubnetId] -= amount`
- `subnetTotalStake[toSubnetId] += amount`
- `userTotalAllocated` unchanged (it's a move)

### freezeAgentAllocations(user, agent)
- Iterate agentSubnets[user][agent]
- For each: `amt = _allocations[user][agent][subnetId]`, zero it, subtract from totals
- `userTotalAllocated[user] -= totalFrozen`
- No force-resolve needed (no pending mechanism)

### getAgentStake(user, agent, subnetId) → uint256
- Direct read: `_allocations[user][agent][subnetId]`
- No epoch parameter needed (no pending to resolve)

## AWPDAO — Custom Counting Module

### Approach
Fork OZ Governor. Remove `GovernorVotes` and `GovernorVotesQuorumFraction` from inheritance. Override `_getVotes`, `_countVote`, `quorum`, and `proposalThreshold`.

### Governor inheritance chain:
```solidity
contract AWPDAO is
    Governor,
    GovernorSettings,
    GovernorTimelockControl,
    AWPVoteCounter          // custom module, replaces GovernorVotes + GovernorCountingSimple
```

### _getVotes override:
```solidity
function _getVotes(address, uint256, bytes memory) internal pure override returns (uint256) {
    return 1; // Nonzero sentinel to pass OZ internal guards; actual power computed in _countVote
}
```

### Storage:
```solidity
StakeNFT public stakeNFT;
mapping(uint256 proposalId => mapping(uint256 tokenId => bool)) public hasVotedWithToken;
```

### _countVote override:
```solidity
function _countVote(
    uint256 proposalId,
    address account,
    uint8 support,
    uint256 weight,     // ignored — sentinel from _getVotes
    bytes memory params
) internal override {
    uint256[] memory tokenIds = abi.decode(params, (uint256[]));
    require(tokenIds.length > 0, "no tokens");
    uint256 totalPower = 0;
    uint256 proposalBlock = proposalSnapshot(proposalId);

    for (uint i = 0; i < tokenIds.length; i++) {
        uint256 tid = tokenIds[i];
        require(stakeNFT.ownerOf(tid) == account, "not owner");
        require(!hasVotedWithToken[proposalId][tid], "already voted");
        require(stakeNFT.remainingEpochs(tid) >= 1, "lock expired");
        // Anti-manipulation: only NFTs minted before proposal can vote
        require(stakeNFT.positions(tid).createdEpoch < _epochAtBlock(proposalBlock), "minted after proposal");

        hasVotedWithToken[proposalId][tid] = true;
        totalPower += stakeNFT.getVotingPower(tid);
    }

    _accumulateVotes(proposalId, support, totalPower);
}
```

### Quorum:
```solidity
function quorum(uint256) public view override returns (uint256) {
    // Upper bound: all staked AWP at max lock
    uint256 maxPower = awpToken.balanceOf(address(stakeNFT)) * Math.sqrt(MAX_WEIGHT_EPOCHS);
    return maxPower * quorumPercent / 100;
}
```

### proposalThreshold:
```solidity
function proposalThreshold() public view override returns (uint256) {
    return 1_000_000e18; // 1M AWP staked minimum to propose (checked via getUserVotingPower)
}

function _checkProposal(address proposer) internal view {
    require(stakeNFT.getUserVotingPower(proposer) >= proposalThreshold(), "below threshold");
}
```

## RootNet Changes

### Deleted:
- `deposit()`, `requestWithdraw()`, `cancelWithdraw()`, `withdraw()` — users call StakeNFT directly
- `setCooldownPeriod()` — no cooldown, replaced by NFT lock

### registerAndStake — Refactored:
```solidity
function registerAndStake(
    uint256 depositAmount,
    uint64 lockEpochs,
    address agent,
    uint256 subnetId,
    uint256 allocateAmount
) external nonReentrant whenNotPaused {
    if (!IAccessManager(accessManager).isRegisteredUser(msg.sender)) {
        IAccessManager(accessManager).register(msg.sender);
        emit UserRegistered(msg.sender);
    }
    if (depositAmount > 0 && lockEpochs > 0) {
        // User must have pre-approved AWP to StakeNFT
        IStakeNFT(stakeNFT).depositFor(msg.sender, depositAmount, lockEpochs);
    }
    if (allocateAmount > 0 && agent != address(0) && subnetId > 0) {
        if (subnets[subnetId].status != SubnetStatus.Active) revert InvalidSubnetStatus();
        _validateAgent(msg.sender, agent);
        IStakingVault(stakingVault).allocate(msg.sender, agent, subnetId, allocateAmount);
        emit Allocated(msg.sender, agent, subnetId, allocateAmount, msg.sender);
    }
}
```

Note: StakeNFT needs a `depositFor(address user, uint256 amount, uint64 lockEpochs)` function callable by RootNet to support `registerAndStake`. AWP is transferred from `user` (not RootNet).

### Modified:
- `allocate()` — StakingVault reads StakeNFT for balance (O(1) via accumulator)
- `deallocate()` — simplified, no epoch param
- `reallocate()` — immediate (emit `Reallocated` instead of `ReallocationQueued`)
- `removeAgent()` — calls `StakingVault.freezeAgentAllocations` (simplified, no epoch param)
- `currentEpoch()` — still delegates to AWPEmission (now returns weekly epochs)
- `getAgentInfo()` / `getAgentsInfo()` — remove epoch parameter from `getAgentStake` call

### Events:
- Remove: `ReallocationQueued` — replace with `Reallocated(user, fromAgent, fromSubnet, toAgent, toSubnet, amount, operator)`
- Remove: `Deposited`, `WithdrawRequested`, `Withdrawn`, `WithdrawCancelled` — moved to StakeNFT
- Keep: `Allocated`, `Deallocated`

## AWPEmission Changes

### epochDuration
Deployment parameter changed from `1 days` to `7 days`. No contract code change needed — just the deploy script parameter.

## Events

### StakeNFT events:
```
event Deposited(address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndEpoch)
event PositionIncreased(uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndEpoch)
event Withdrawn(address indexed user, uint256 indexed tokenId, uint256 amount)
event Transfer(address indexed from, address indexed to, uint256 indexed tokenId) // ERC721 standard
```

### RootNet events (updated):
```
event Reallocated(address indexed user, address fromAgent, uint256 fromSubnet, address toAgent, uint256 toSubnet, uint256 amount, address operator)
```
Replaces `ReallocationQueued`.

## DB Schema Changes

### New table: `stake_positions`
```sql
CREATE TABLE stake_positions (
    token_id       BIGINT PRIMARY KEY,
    owner          CHAR(42) NOT NULL,
    amount         NUMERIC(78,0) NOT NULL,
    lock_end_epoch BIGINT NOT NULL,
    created_epoch  BIGINT NOT NULL,
    created_at     BIGINT NOT NULL,
    burned         BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_sp_owner ON stake_positions(owner);
```

### Remove from schema:
- `withdraw_requests` table — no more cooldown
- `user_balances.total_balance` column — replaced by sum of stake_positions

### Keep:
- `stake_allocations` table — unchanged
- `user_balances.total_allocated` column — kept

## Files Changed

| File | Action |
|------|--------|
| `contracts/src/core/StakeNFT.sol` | **Create** — ERC721Enumerable position NFT |
| `contracts/src/core/StakingVault.sol` | **Major refactor** — remove deposit/withdraw/STP/pending, simplify to pure allocation |
| `contracts/src/interfaces/IStakeNFT.sol` | **Create** |
| `contracts/src/interfaces/IStakingVault.sol` | **Update** — remove deleted functions |
| `contracts/src/RootNet.sol` | **Refactor** — remove deposit/withdraw, simplify allocate/reallocate |
| `contracts/src/governance/AWPDAO.sol` | **Rewrite** — custom counting module with tokenId[] voting |
| `contracts/script/Deploy.s.sol` | **Update** — deploy StakeNFT, epoch=7 days |
| `contracts/script/TestDeploy.s.sol` | **Update** |
| `contracts/test/` | **Rewrite** — all staking + governance tests |
| `api/internal/db/schema.sql` | **Update** — new table, remove withdraw_requests |
| `api/internal/chain/indexer.go` | **Update** — new events |
| `api/internal/chain/bindings/` | **Regenerate** |

## Security Considerations

1. **Transfer + allocation check**: O(1) `_userTotalStaked` accumulator ensures transfer reverts if sender's remaining amount < allocated. No enumeration ordering hazard.
2. **Withdraw ordering**: Balance check runs BEFORE burn. `_userTotalStaked` decremented before burn. No double-deduction risk.
3. **No flash loan voting**: Voting requires `remainingEpochs >= 1` (minimum 1-week lock). Flash loans cannot create locked positions.
4. **Anti-manipulation on voting**: `createdEpoch < epochAtProposal` check prevents stake-and-vote after seeing a proposal. NFTs minted after proposal creation cannot vote on that proposal.
5. **NFT pass-around within voting window**: Each tokenId can only vote once per proposal (`hasVotedWithToken`). Transferring an NFT to another address and voting again with the same tokenId is blocked. However, two colluding addresses with separate NFTs can each vote independently — this is inherent to transferable governance and is a known trade-off.
6. **Quorum based on upper bound**: `totalStaked * sqrt(54)` is an overestimate, meaning actual quorum requirement is stricter than the percentage suggests. This is conservative (harder to reach quorum = safer).
7. **Immediate reallocate**: No pending mechanism means no freeze/resolve complexity. Subnets monitor events for real-time allocation changes.
8. **Lock integrity**: lockEndEpoch can only increase (addToPosition enforces >= current and > currentEpoch). No way to shorten lock.
9. **Governor compatibility**: `GovernorVotes` removed from inheritance. `_getVotes` returns sentinel `1` to pass OZ guards. Actual power computed in `_countVote`.
