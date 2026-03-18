# StakeNFT + vAWP Voting + StakingVault Refactor — Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace deposit/withdraw with NFT-based staking positions, add sqrt-weighted voting, simplify StakingVault by removing STP and pending mechanism.

**Architecture:** StakeNFT (ERC721Enumerable) manages locked positions. StakingVault becomes pure allocation logic. AWPDAO uses custom counting module accepting tokenId[] for voting. All allocate/deallocate/reallocate are immediate (no pending).

**Tech Stack:** Solidity 0.8.24, OpenZeppelin 5.x, Foundry, Go 1.26

**Spec:** `docs/superpowers/specs/2026-03-16-stake-nft-vawp-design.md`

---

## Chunk 1: New Contracts + Interface Updates

### Task 1: Create IStakeNFT Interface

**Files:**
- Create: `contracts/src/interfaces/IStakeNFT.sol`

- [ ] **Step 1: Write interface**

```solidity
interface IStakeNFT {
    struct Position {
        uint128 amount;
        uint64 lockEndEpoch;
        uint64 createdEpoch;
    }

    function deposit(uint256 amount, uint64 lockEpochs) external returns (uint256 tokenId);
    function depositFor(address user, uint256 amount, uint64 lockEpochs) external returns (uint256 tokenId);
    function addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndEpoch) external;
    function withdraw(uint256 tokenId) external;

    function positions(uint256 tokenId) external view returns (uint128 amount, uint64 lockEndEpoch, uint64 createdEpoch);
    function getUserTotalStaked(address user) external view returns (uint256);
    function getVotingPower(uint256 tokenId) external view returns (uint256);
    function getUserVotingPower(address user) external view returns (uint256);
    function totalVotingPower() external view returns (uint256);
    function remainingEpochs(uint256 tokenId) external view returns (uint64);

    event Deposited(address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndEpoch);
    event PositionIncreased(uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndEpoch);
    event Withdrawn(address indexed user, uint256 indexed tokenId, uint256 amount);
}
```

### Task 2: Create StakeNFT Contract

**Files:**
- Create: `contracts/src/core/StakeNFT.sol`

- [ ] **Step 1: Write full StakeNFT implementation**

Key elements:
- Inherit `ERC721Enumerable`
- Storage: `Position` struct (amount, lockEndEpoch, createdEpoch), `_userTotalStaked` mapping (O(1) accumulator), `_nextTokenId`
- Immutables: `awpToken`, `stakingVault`, `awpEmission`
- Constants: `MAX_WEIGHT_EPOCHS = 54`, `MIN_LOCK_EPOCHS = 1`
- `onlyRootNet` modifier for `depositFor`

Functions:
- `deposit(amount, lockEpochs)`: transfer AWP, mint NFT, `lockEndEpoch = currentEpoch + 1 + lockEpochs`, update `_userTotalStaked`
- `depositFor(user, amount, lockEpochs)`: onlyRootNet, transfers AWP from `user` (not msg.sender), mints NFT to `user`
- `addToPosition(tokenId, amount, newLockEndEpoch)`: ownerOf check, extend only (>= current lockEnd, > currentEpoch), update accumulator
- `withdraw(tokenId)`: ownerOf check, `currentEpoch >= lockEndEpoch`, check `_userTotalStaked - amount >= userTotalAllocated`, update accumulator BEFORE burn
- `getUserTotalStaked(user)`: O(1) read from `_userTotalStaked`
- `getVotingPower(tokenId)`: `amount * Math.sqrt(min(remaining, 54))`
- `getUserVotingPower(user)`: iterate ERC721Enumerable, sum voting powers
- `totalVotingPower()`: upper bound = `awpToken.balanceOf(address(this)) * Math.sqrt(54)`
- `remainingEpochs(tokenId)`: `lockEndEpoch > currentEpoch ? lockEndEpoch - currentEpoch : 0`
- `_update` override: on transfer (not mint/burn), update `_userTotalStaked` for both parties, check sender has enough for allocations

Import `Math` from OpenZeppelin for `sqrt`.

### Task 3: Update IStakingVault Interface

**Files:**
- Modify: `contracts/src/interfaces/IStakingVault.sol`

- [ ] **Step 1: Rewrite interface**

Remove: `deposit`, `requestWithdraw`, `cancelWithdraw`, `withdraw`, `setCooldownPeriod`, `WithdrawRequest`, `withdrawRequests`, `cooldownPeriod`, `userBalance`, `SubnetSTP`, `getSubnetSTP`

Update signatures — remove `epoch` parameter:
- `allocate(address user, address agent, uint256 subnetId, uint256 amount)`
- `deallocate(address user, address agent, uint256 subnetId, uint256 amount)`
- `reallocate(address user, address fromAgent, uint256 fromSubnetId, address toAgent, uint256 toSubnetId, uint256 amount)`
- `freezeAgentAllocations(address user, address agent)`
- `getAgentStake(address user, address agent, uint256 subnetId) → uint256`

Keep: `userTotalAllocated`, `subnetTotalStake`, `getSubnetTotalStake`, `getAgentSubnets`

### Task 4: Rewrite StakingVault

**Files:**
- Rewrite: `contracts/src/core/StakingVault.sol`

- [ ] **Step 1: Write simplified StakingVault**

Delete entirely: `Allocation` struct, `_resolve`, `_effectiveAmount`, `_updateSTP`, `SubnetSTP`, `WithdrawRequest`, deposit/withdraw functions, cooldown, STP.

New storage:
```solidity
mapping(address => mapping(address => mapping(uint256 => uint128))) private _allocations;
mapping(address => uint256) public userTotalAllocated;
mapping(uint256 => uint256) public subnetTotalStake;
// agentSubnets tracking for freezeAgentAllocations
mapping(address => mapping(address => EnumerableSet.UintSet)) private agentSubnets;

address public immutable rootNet;
address public stakeNFT;
```

Functions (all `onlyRootNet`):
- `allocate(user, agent, subnetId, amount)`: check `amount > 0 && amount <= type(uint128).max`, check `StakeNFT.getUserTotalStaked(user) - userTotalAllocated[user] >= amount`, add to allocation + totals, add to agentSubnets
- `deallocate(user, agent, subnetId, amount)`: check amount, subtract from allocation + totals
- `reallocate(user, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)`: immediate atomic: subtract from source, add to destination, adjust subnetTotalStake, no change to userTotalAllocated
- `freezeAgentAllocations(user, agent)`: iterate agentSubnets, zero each allocation, subtract from totals
- `getAgentStake(user, agent, subnetId)`: direct `_allocations` read (no epoch, no resolve)
- `getUnallocated(user)`: `StakeNFT.getUserTotalStaked(user) - userTotalAllocated[user]`

### Task 5: Update RootNet

**Files:**
- Modify: `contracts/src/RootNet.sol`

- [ ] **Step 1: Remove deposit/withdraw functions**

Delete: `deposit()`, `requestWithdraw()`, `cancelWithdraw()`, `withdraw()`, `setCooldownPeriod()`

- [ ] **Step 2: Add stakeNFT to registry**

Add `address public stakeNFT` to registry fields.
Update `initializeRegistry` to accept stakeNFT address.
Update `updateAddress` to handle `"stakeNFT"` key.
Update `getRegistry` to return stakeNFT.

- [ ] **Step 3: Refactor registerAndStake**

New signature: `registerAndStake(uint256 depositAmount, uint64 lockEpochs, address agent, uint256 subnetId, uint256 allocateAmount)`
Call `IStakeNFT(stakeNFT).depositFor(msg.sender, depositAmount, lockEpochs)` instead of deposit.

- [ ] **Step 4: Simplify allocate/deallocate/reallocate**

Remove `currentEpoch()` parameter from StakingVault calls.
`reallocate`: emit `Reallocated` event instead of `ReallocationQueued`.
Add new event: `Reallocated(address indexed user, address fromAgent, uint256 fromSubnet, address toAgent, uint256 toSubnet, uint256 amount, address operator)`
Remove old event: `ReallocationQueued`

- [ ] **Step 5: Simplify removeAgent and getAgentInfo**

`removeAgent`: call `freezeAgentAllocations(user, agent)` without epoch param.
`getAgentInfo` / `getAgentsInfo`: call `getAgentStake(user, agent, subnetId)` without epoch param.

- [ ] **Step 6: Remove stale events and imports**

Remove events: `Deposited`, `WithdrawRequested`, `Withdrawn`, `WithdrawCancelled`, `ReallocationQueued`
Update IRootNet.sol if these events are defined there.

### Task 6: Rewrite AWPDAO

**Files:**
- Rewrite: `contracts/src/governance/AWPDAO.sol`

- [ ] **Step 1: Write custom Governor with NFT voting**

Remove inheritance: `GovernorVotes`, `GovernorCountingSimple`, `GovernorVotesQuorumFraction`
Keep: `Governor`, `GovernorSettings`, `GovernorTimelockControl`

Storage:
- `IStakeNFT public stakeNFT`
- `IERC20 public awpToken`
- `uint256 public quorumPercent`
- `mapping(uint256 => mapping(uint256 => bool)) public hasVotedWithToken`
- Vote tracking: `mapping(uint256 => ProposalVote)` with `forVotes`, `againstVotes`, `abstainVotes`, `hasVoted[address]`

Override `_getVotes`: return sentinel `1`
Override `_countVote`: decode `tokenId[]` from params, verify ownership + lock + created epoch, sum voting power
Override `quorum`: `awpToken.balanceOf(address(stakeNFT)) * Math.sqrt(54) * quorumPercent / 100`
Override `proposalThreshold`: `1_000_000e18`
Override `_voteSucceeded`: standard majority check
Override `_quorumReached`: compare accumulated for+abstain votes against quorum

### Task 7: Update Deploy Scripts

**Files:**
- Modify: `contracts/script/Deploy.s.sol`
- Modify: `contracts/script/TestDeploy.s.sol`

- [ ] **Step 1: Add StakeNFT to deployment**

Insert StakeNFT deployment between LPManager and AWPEmission.
Change `EPOCH_DURATION` from `1 days` to `7 days`.
Update `initializeRegistry` to pass StakeNFT address.
Update AWPDAO constructor to pass StakeNFT and awpToken.

- [ ] **Step 2: Verify source compilation**

Run: `cd /home/ubuntu/code/Cortexia/contracts && /home/ubuntu/.foundry/bin/forge build --skip test 2>&1`

- [ ] **Step 3: Commit**

```bash
git add contracts/src/ contracts/script/
git commit -m "feat: StakeNFT + StakingVault refactor + AWPDAO custom voting — source contracts"
```

## Chunk 2: Test Updates

### Task 8: Rewrite StakingVault and StakeNFT Tests

**Files:**
- Create: `contracts/test/StakeNFT.t.sol`
- Rewrite: `contracts/test/StakingVault.t.sol`

- [ ] **Step 1: Write StakeNFT tests**

Tests:
- `test_deposit`: deposit AWP, verify NFT minted, position stored, _userTotalStaked updated
- `test_deposit_multiplePosistions`: create 3 NFTs, verify getUserTotalStaked sums correctly
- `test_addToPosition_amount`: add AWP to existing position
- `test_addToPosition_extendLock`: extend lock (only forward, not backward)
- `test_addToPosition_revertExpiredExtension`: newLockEndEpoch <= currentEpoch reverts
- `test_withdraw`: after lock expiry, withdraw burns NFT and returns AWP
- `test_withdraw_revertBeforeExpiry`: withdraw before lock end reverts
- `test_withdraw_revertInsufficientUnallocated`: cannot withdraw if remaining < allocated
- `test_transfer`: NFT transfer updates _userTotalStaked for both parties
- `test_transfer_revertInsufficientUnallocated`: transfer blocked if sender would be undercollateralized
- `test_getVotingPower`: verify sqrt calculation at various remaining epochs
- `test_getVotingPower_cappedAt54`: remaining > 54 caps at sqrt(54)
- `test_totalVotingPower`: upper bound = totalStaked * sqrt(54)

- [ ] **Step 2: Rewrite StakingVault tests**

Simplified tests (no pending, no STP, no cooldown):
- `test_allocate`: allocate and verify balances
- `test_deallocate`: deallocate and verify
- `test_reallocate`: immediate move between (agent, subnet) pairs
- `test_freezeAgentAllocations`: zero all allocations for an agent
- `test_allocate_revertExceedsAvailable`: cannot allocate more than staked
- `test_getAgentStake`: direct read, no epoch

### Task 9: Rewrite RootNet Tests

**Files:**
- Rewrite: `contracts/test/RootNet.t.sol`

- [ ] **Step 1: Update setUp with StakeNFT**

Deploy StakeNFT in test setUp. Update initializeRegistry.
Replace `deposit()` calls with `stakeNFT.deposit()`.
Remove withdraw/cooldown tests.

- [ ] **Step 2: Update allocation tests**

Remove epoch parameters. Update `test_reallocate` for immediate semantics.
Remove `test_requestWithdraw`, `test_cancelWithdraw`, `test_withdraw`.

### Task 10: Rewrite E2E and Integration Tests

**Files:**
- Rewrite: `contracts/test/E2E.t.sol`
- Rewrite: `contracts/test/Integration.t.sol`

- [ ] **Step 1: Update E2E _deploy() and helpers**

Add StakeNFT to deployment. Change epoch from 1 day to 7 days.
Replace `deposit(amount)` → `stakeNFT.deposit(amount, lockEpochs)` in helpers.
Replace `requestWithdraw/withdraw` flow with `stakeNFT.withdraw(tokenId)`.
Update `_settleEpoch()` helper for weekly epochs.
Update `test_e2e_daoGovernanceWeight` for NFT-based voting (castVoteWithReasonAndParams with tokenId[]).

- [ ] **Step 2: Update Integration tests**

Same pattern. Remove withdrawal flow tests. Add StakeNFT deployment.

### Task 11: Rewrite AWPDAO Tests

**Files:**
- Rewrite: `contracts/test/AWPDAO.t.sol`

- [ ] **Step 1: Write voting tests**

- `test_castVoteWithNFT`: stake, propose, vote with tokenId[], verify power = amount * sqrt(remaining)
- `test_voteRevertsExpiredLock`: NFT with remainingEpochs=0 cannot vote
- `test_voteRevertsNotOwner`: non-owner cannot vote with someone else's NFT
- `test_voteRevertsAlreadyVoted`: same tokenId cannot vote twice on same proposal
- `test_voteRevertsMintedAfterProposal`: NFT created after proposal cannot vote
- `test_quorum`: verify quorum calculation

- [ ] **Step 2: Run full test suite**

Run: `cd /home/ubuntu/code/Cortexia/contracts && /home/ubuntu/.foundry/bin/forge test --no-match-contract Fork 2>&1`
All tests must pass.

- [ ] **Step 3: Commit**

```bash
git add contracts/test/
git commit -m "test: StakeNFT + StakingVault + AWPDAO voting tests"
```

## Chunk 3: Go Backend + DB

### Task 12: Update DB Schema and Queries

**Files:**
- Modify: `api/internal/db/schema.sql`
- Modify: `api/internal/db/query/staking.sql`
- Create: `api/internal/db/query/stake_position.sql`

- [ ] **Step 1: Update schema**

Add `stake_positions` table. Remove `withdraw_requests` table. Remove `total_balance` from `user_balances`.

- [ ] **Step 2: Add stake position queries**

```sql
-- name: InsertStakePosition :exec
INSERT INTO stake_positions (token_id, owner, amount, lock_end_epoch, created_epoch, created_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateStakePositionOwner :exec
UPDATE stake_positions SET owner = $2 WHERE token_id = $1;

-- name: UpdateStakePosition :exec
UPDATE stake_positions SET amount = $2, lock_end_epoch = $3 WHERE token_id = $1;

-- name: BurnStakePosition :exec
UPDATE stake_positions SET burned = TRUE WHERE token_id = $1;

-- name: GetStakePosition :one
SELECT * FROM stake_positions WHERE token_id = $1;

-- name: GetUserStakePositions :many
SELECT * FROM stake_positions WHERE owner = $1 AND burned = FALSE ORDER BY token_id;

-- name: GetUserTotalStaked :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC(78,0) AS total FROM stake_positions WHERE owner = $1 AND burned = FALSE;
```

- [ ] **Step 3: Regenerate sqlc**

```bash
cd /home/ubuntu/code/Cortexia/api && /home/ubuntu/go/bin/sqlc generate
```

### Task 13: Regenerate Go Bindings and Update Indexer

**Files:**
- Regenerate: `api/internal/chain/bindings/`
- Modify: `api/internal/chain/indexer.go`
- Create: `api/internal/chain/bindings/stake_n_f_t.go`

- [ ] **Step 1: Regenerate all bindings**

```bash
cd /home/ubuntu/code/Cortexia/contracts
/home/ubuntu/.foundry/bin/forge inspect StakeNFT abi --json > /tmp/stake_nft_abi.json
/home/ubuntu/go/bin/abigen --abi /tmp/stake_nft_abi.json --pkg bindings --type StakeNFT --out /home/ubuntu/code/Cortexia/api/internal/chain/bindings/stake_n_f_t.go
/home/ubuntu/.foundry/bin/forge inspect RootNet abi --json > /tmp/rootnet_abi.json
/home/ubuntu/go/bin/abigen --abi /tmp/rootnet_abi.json --pkg bindings --type RootNet --out /home/ubuntu/code/Cortexia/api/internal/chain/bindings/root_net.go
/home/ubuntu/.foundry/bin/forge inspect AWPDAO abi --json > /tmp/awpdao_abi.json
/home/ubuntu/go/bin/abigen --abi /tmp/awpdao_abi.json --pkg bindings --type AWPDAO --out /home/ubuntu/code/Cortexia/api/internal/chain/bindings/a_w_p_d_a_o.go
```

- [ ] **Step 2: Update indexer for new events**

Add StakeNFT event parsing: `Deposited` (new format with tokenId), `PositionIncreased`, `Withdrawn` (new format with tokenId).
Update `Reallocated` event parsing (replaces `ReallocationQueued`).
Remove handlers for old events: `WithdrawRequested`, `WithdrawCancelled`, `Withdrawn` (old format), `Deposited` (old format from RootNet).
Add StakeNFT contract address to filter list.

- [ ] **Step 3: Update client.go**

Add `StakeNFTAddr` and `StakeNFT` binding to `Client` struct.

- [ ] **Step 4: Verify Go build**

Run: `cd /home/ubuntu/code/Cortexia/api && go build ./...`

- [ ] **Step 5: Commit**

```bash
git add api/
git commit -m "feat: StakeNFT Go bindings, indexer events, DB schema for NFT staking"
```

### Task 14: Final Verification

- [ ] **Step 1: Full Solidity test suite**

Run: `cd /home/ubuntu/code/Cortexia/contracts && /home/ubuntu/.foundry/bin/forge test --no-match-contract Fork`
All tests must pass.

- [ ] **Step 2: Go build**

Run: `cd /home/ubuntu/code/Cortexia/api && go build ./...`
Must compile clean.
