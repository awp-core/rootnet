package chain

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	// syncKey is the key name for sync state in the sync_states table
	syncKey = "indexer"
	// pollInterval is the interval for polling new on-chain blocks
	pollInterval = 3 * time.Second
	// redisChannel is the Redis Pub/Sub channel name
	redisChannel = "chain_events"
)

// Indexer is the on-chain event indexer responsible for scanning block logs and writing to PostgreSQL
type Indexer struct {
	chain         *Client
	pool          *pgxpool.Pool
	rds           *redis.Client
	deployBlock   int64
	genesisTime   int64 // AWPEmission genesisTime (unix seconds)
	epochDuration int64 // AWPEmission epochDuration (seconds)
}

// NewIndexer creates an event indexer instance
// deployBlock is used as the start block on first run when sync_states is empty.
func NewIndexer(chain *Client, pool *pgxpool.Pool, rds *redis.Client, deployBlock int64) (*Indexer, error) {
	// Cache genesisTime and epochDuration (immutable, read once)
	gt, err := chain.AWPEmission.GenesisTime(nil)
	if err != nil || gt == nil {
		return nil, fmt.Errorf("failed to read AWPEmission.genesisTime: %w", err)
	}
	ed, err := chain.AWPEmission.EpochDuration(nil)
	if err != nil || ed == nil {
		return nil, fmt.Errorf("failed to read AWPEmission.epochDuration: %w", err)
	}
	return &Indexer{
		chain:         chain,
		pool:          pool,
		rds:           rds,
		deployBlock:   deployBlock,
		genesisTime:   gt.Int64(),
		epochDuration: ed.Int64(),
	}, nil
}

// redisEvent is the event format published to Redis
type redisEvent struct {
	Type        string      `json:"type"`
	BlockNumber uint64      `json:"blockNumber"`
	TxHash      string      `json:"txHash"`
	Data        interface{} `json:"data"`
}

// Run starts the indexer main loop: read sync progress → filter logs → process events → update progress → publish Redis
func (idx *Indexer) Run(ctx context.Context) error {
	slog.Info("indexer started")
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("indexer received stop signal")
			return ctx.Err()
		case <-ticker.C:
			if err := idx.poll(ctx); err != nil {
				slog.Error("poll failed", "error", err)
			}
		}
	}
}

// poll executes one complete scan cycle using optimistic indexing with parent hash verification.
// No fixed confirmation depth — processes up to chain tip, detects reorgs via block hash chain.
func (idx *Indexer) poll(parentCtx context.Context) error {
	// Per-poll timeout to prevent indefinite hangs on RPC failures
	ctx, cancel := context.WithTimeout(parentCtx, 60*time.Second)
	defer cancel()

	// 1. Read the last synced block
	q := gen.New(idx.pool)
	state, err := q.GetSyncState(ctx, syncKey)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("failed to read sync state: %w", err)
		}
		if idx.deployBlock > 0 {
			state.LastBlock = idx.deployBlock - 1
		} else {
			state.LastBlock = 0
		}
	}
	fromBlock := uint64(state.LastBlock) + 1

	// 2. Get the latest block number on-chain (no confirmation depth — process up to tip)
	latestBlock, err := idx.chain.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}
	if fromBlock > latestBlock {
		return nil
	}

	// 3. Reorg detection: verify the parent hash chain before processing new blocks
	if fromBlock > 1 {
		reorgBlock, err := idx.detectReorg(ctx, q, fromBlock-1)
		if err != nil {
			return fmt.Errorf("reorg detection failed: %w", err)
		}
		if reorgBlock > 0 {
			slog.Warn("reorg detected, rolling back", "forkPoint", reorgBlock)
			if err := idx.rollback(ctx, reorgBlock); err != nil {
				return fmt.Errorf("reorg rollback failed: %w", err)
			}
			fromBlock = uint64(reorgBlock) + 1
		}
	}

	// Limit single query range to avoid RPC restrictions
	const maxBlockRange = 500
	toBlock := latestBlock
	if toBlock-fromBlock > maxBlockRange {
		toBlock = fromBlock + maxBlockRange
	}

	slog.Info("scanning blocks", "from", fromBlock, "to", toBlock)

	// 4. Filter logs from all monitored contracts
	logs, err := idx.chain.Eth.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
		Addresses: []common.Address{idx.chain.AWPRegistryAddr, idx.chain.AWPEmissionAddr, idx.chain.StakeNFTAddr, idx.chain.SubnetNFTAddr, idx.chain.AWPDAOAddr},
	})
	if err != nil {
		return fmt.Errorf("failed to filter logs: %w", err)
	}

	// 5. Process all events within a PostgreSQL transaction
	tx, err := idx.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin database transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := gen.New(tx)
	var events []redisEvent

	for i := range logs {
		lg := logs[i]
		if len(lg.Topics) == 0 {
			continue
		}
		evts, err := idx.processLog(ctx, qtx, lg)
		if err != nil {
			return fmt.Errorf("failed to process log at block %d tx %s: %w", lg.BlockNumber, lg.TxHash.Hex(), err)
		}
		events = append(events, evts...)
	}

	// 6. Record block hashes for reorg detection
	if err := idx.recordBlockHashes(ctx, qtx, fromBlock, toBlock); err != nil {
		return fmt.Errorf("failed to record block hashes: %w", err)
	}

	// 7. Update sync progress
	if err := qtx.UpsertSyncState(ctx, gen.UpsertSyncStateParams{
		ContractName: syncKey,
		LastBlock:    int64(toBlock),
	}); err != nil {
		return fmt.Errorf("failed to update sync state: %w", err)
	}

	// Prune old block hashes (keep last 256 blocks to limit DB growth)
	if int64(toBlock) > 256 {
		_ = qtx.PruneIndexedBlocks(ctx, int64(toBlock)-256)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 6. Publish events to Redis (after successful transaction commit)
	for _, evt := range events {
		payload, err := json.Marshal(evt)
		if err != nil {
			slog.Error("failed to serialize Redis event", "type", evt.Type, "error", err)
			continue
		}
		if err := idx.rds.Publish(ctx, redisChannel, payload).Err(); err != nil {
			slog.Error("failed to publish Redis event", "type", evt.Type, "error", err)
		}
	}

	slog.Info("scan complete", "blocks", toBlock-fromBlock+1, "logs", len(logs), "events", len(events))
	return nil
}

// detectReorg walks back from lastBlock checking stored block hashes against the chain.
// Returns the fork point block number (>0 if reorg detected, 0 if no reorg).
// Since we store one hash per batch (toBlock only), ErrNoRows for intermediate blocks
// is expected and skipped — the walk-back continues until it finds a stored hash to verify.
func (idx *Indexer) detectReorg(ctx context.Context, q *gen.Queries, lastBlock uint64) (int64, error) {
	const maxReorgDepth = 64
	mismatchSeen := false

	for i := uint64(0); i < maxReorgDepth && lastBlock > i; i++ {
		checkBlock := lastBlock - i
		storedHash, err := q.GetIndexedBlockHash(ctx, int64(checkBlock))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue // no hash stored for this block (not a batch boundary), skip
			}
			return 0, fmt.Errorf("get indexed block hash %d: %w", checkBlock, err)
		}

		// Fetch actual block hash from chain
		header, err := idx.chain.Eth.HeaderByNumber(ctx, new(big.Int).SetUint64(checkBlock))
		if err != nil {
			return 0, fmt.Errorf("get header for block %d: %w", checkBlock, err)
		}

		chainHash := header.Hash().Hex()
		if strings.EqualFold(strings.TrimSpace(storedHash), chainHash) {
			if !mismatchSeen {
				return 0, nil // latest stored hash matches, no reorg
			}
			// Found the common ancestor — everything after this was reorged
			return int64(checkBlock), nil
		}
		mismatchSeen = true
		slog.Warn("block hash mismatch", "block", checkBlock, "stored", storedHash, "chain", chainHash)
	}

	if !mismatchSeen {
		return 0, nil // no stored hashes found within lookback range (first run)
	}
	// All stored hashes diverged — deep reorg, reset to earliest checked block
	return int64(lastBlock - maxReorgDepth + 1), nil
}

// rollback resets the indexer state to the given fork point atomically within a transaction.
// Allocation tracking uses additive upserts, so we must truncate stake_allocations and
// user_balances on rollback to prevent double-counting when events are replayed.
func (idx *Indexer) rollback(ctx context.Context, forkPoint int64) error {
	tx, err := idx.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin rollback tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := gen.New(tx)
	if err := qtx.UpsertSyncState(ctx, gen.UpsertSyncStateParams{
		ContractName: syncKey,
		LastBlock:    forkPoint,
	}); err != nil {
		return fmt.Errorf("reset sync state: %w", err)
	}
	if err := qtx.DeleteIndexedBlocksAfter(ctx, forkPoint); err != nil {
		return fmt.Errorf("delete indexed blocks after %d: %w", forkPoint, err)
	}
	// Delete only rows written after the fork point (scoped rollback, not global truncate)
	if err := qtx.DeleteStakeAllocationsAfterBlock(ctx, forkPoint); err != nil {
		return fmt.Errorf("delete stake_allocations after %d: %w", forkPoint, err)
	}
	if err := qtx.DeleteUserBalancesAfterBlock(ctx, forkPoint); err != nil {
		return fmt.Errorf("delete user_balances after %d: %w", forkPoint, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit rollback tx: %w", err)
	}
	slog.Info("rollback complete", "forkPoint", forkPoint)
	return nil
}

// recordBlockHashes stores the hash of the last block in the batch for reorg detection.
// Only 1 RPC call per poll cycle. This is sufficient because blockchain is a hash chain —
// a reorg at any depth propagates to change all subsequent block hashes.
func (idx *Indexer) recordBlockHashes(ctx context.Context, q *gen.Queries, fromBlock, toBlock uint64) error {
	header, err := idx.chain.Eth.HeaderByNumber(ctx, new(big.Int).SetUint64(toBlock))
	if err != nil {
		return fmt.Errorf("get header for block %d: %w", toBlock, err)
	}
	return q.UpsertIndexedBlock(ctx, gen.UpsertIndexedBlockParams{
		BlockNumber: int64(toBlock),
		BlockHash:   header.Hash().Hex(),
	})
}

// processLog parses a single log entry and performs the corresponding database writes, returning Redis events to publish
func (idx *Indexer) processLog(ctx context.Context, q *gen.Queries, lg types.Log) ([]redisEvent, error) {
	awpRegistry := idx.chain.AWPRegistry
	awpEmission := idx.chain.AWPEmission
	stakeNFT := idx.chain.StakeNFT

	// ── AWPRegistry Account System V2 events ──

	// UserRegistered(address indexed user)
	if evt, err := awpRegistry.ParseUserRegistered(lg); err == nil {
		if err := q.SetUserRegisteredAt(ctx, gen.SetUserRegisteredAtParams{
			Address:      strings.ToLower(evt.User.Hex()),
			RegisteredAt: int64(lg.BlockNumber),
		}); err != nil {
			return nil, fmt.Errorf("SetUserRegisteredAt: %w", err)
		}
		return []redisEvent{makeEvent("UserRegistered", lg, map[string]interface{}{
			"user": evt.User.Hex(),
		})}, nil
	}

	// Bound(address indexed addr, address indexed target)
	if evt, err := awpRegistry.ParseBound(lg); err == nil {
		if err := q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{
			Address: strings.ToLower(evt.Addr.Hex()),
			BoundTo: strings.ToLower(evt.Target.Hex()),
		}); err != nil {
			return nil, fmt.Errorf("UpsertUserBinding: %w", err)
		}
		if err := q.InitUserBalance(ctx, strings.ToLower(evt.Addr.Hex())); err != nil {
			return nil, fmt.Errorf("InitUserBalance (Bound): %w", err)
		}
		return []redisEvent{makeEvent("Bound", lg, map[string]interface{}{
			"addr":   evt.Addr.Hex(),
			"target": evt.Target.Hex(),
		})}, nil
	}

	// Unbound(address indexed addr)
	if evt, err := awpRegistry.ParseUnbound(lg); err == nil {
		if err := q.ClearUserBinding(ctx, strings.ToLower(evt.Addr.Hex())); err != nil {
			return nil, fmt.Errorf("ClearUserBinding: %w", err)
		}
		return []redisEvent{makeEvent("Unbound", lg, map[string]interface{}{
			"addr": evt.Addr.Hex(),
		})}, nil
	}

	// RecipientSet(address indexed addr, address recipient)
	if evt, err := awpRegistry.ParseRecipientSet(lg); err == nil {
		if err := q.UpsertUserRecipient(ctx, gen.UpsertUserRecipientParams{
			Address:   strings.ToLower(evt.Addr.Hex()),
			Recipient: strings.ToLower(evt.Recipient.Hex()),
		}); err != nil {
			return nil, fmt.Errorf("UpsertUserRecipient: %w", err)
		}
		return []redisEvent{makeEvent("RecipientSet", lg, map[string]interface{}{
			"addr":      evt.Addr.Hex(),
			"recipient": evt.Recipient.Hex(),
		})}, nil
	}

	// DelegateGranted(address indexed staker, address indexed delegate)
	if evt, err := awpRegistry.ParseDelegateGranted(lg); err == nil {
		return []redisEvent{makeEvent("DelegateGranted", lg, map[string]interface{}{
			"staker":   evt.Staker.Hex(),
			"delegate": evt.Delegate.Hex(),
		})}, nil
	}

	// DelegateRevoked(address indexed staker, address indexed delegate)
	if evt, err := awpRegistry.ParseDelegateRevoked(lg); err == nil {
		return []redisEvent{makeEvent("DelegateRevoked", lg, map[string]interface{}{
			"staker":   evt.Staker.Hex(),
			"delegate": evt.Delegate.Hex(),
		})}, nil
	}

	// StakeNFT.Deposited — new stake position created
	if evt, err := stakeNFT.ParseDeposited(lg); err == nil {
		// Read position from chain to get the actual createdAt timestamp
		pos, err := idx.chain.StakeNFT.Positions(nil, evt.TokenId)
		if err != nil {
			return nil, fmt.Errorf("failed to read position for createdAt: %w", err)
		}
		if err := q.InsertStakePosition(ctx, gen.InsertStakePositionParams{
			TokenID:     evt.TokenId.Int64(),
			Owner:       strings.ToLower(evt.User.Hex()),
			Amount:      bigIntToNumeric(evt.Amount),
			LockEndTime: int64(evt.LockEndTime),
			CreatedAt:   int64(pos.CreatedAt),
		}); err != nil {
			return nil, fmt.Errorf("InsertStakePosition: %w", err)
		}
		return []redisEvent{makeEvent("Deposited", lg, map[string]interface{}{
			"user":        evt.User.Hex(),
			"tokenId":     evt.TokenId.String(),
			"amount":      evt.Amount.String(),
			"lockEndTime": evt.LockEndTime,
		})}, nil
	}

	// StakeNFT.PositionIncreased — position amount/lock updated
	if evt, err := stakeNFT.ParsePositionIncreased(lg); err == nil {
		// Read updated position from chain to get new total amount
		pos, err := idx.chain.StakeNFT.Positions(nil, evt.TokenId)
		if err != nil {
			return nil, fmt.Errorf("failed to read position: %w", err)
		}
		if err := q.UpdateStakePosition(ctx, gen.UpdateStakePositionParams{
			TokenID:     evt.TokenId.Int64(),
			Amount:      bigIntToNumeric(pos.Amount),
			LockEndTime: int64(pos.LockEndTime),
		}); err != nil {
			return nil, fmt.Errorf("UpdateStakePosition: %w", err)
		}
		return []redisEvent{makeEvent("PositionIncreased", lg, map[string]interface{}{
			"tokenId":        evt.TokenId.String(),
			"addedAmount":    evt.AddedAmount.String(),
			"newLockEndTime": evt.NewLockEndTime,
		})}, nil
	}

	// StakeNFT.Withdrawn — position burned
	if evt, err := stakeNFT.ParseWithdrawn(lg); err == nil {
		if err := q.BurnStakePosition(ctx, evt.TokenId.Int64()); err != nil {
			return nil, fmt.Errorf("BurnStakePosition: %w", err)
		}
		return []redisEvent{makeEvent("Withdrawn", lg, map[string]interface{}{
			"user":    evt.User.Hex(),
			"tokenId": evt.TokenId.String(),
			"amount":  evt.Amount.String(),
		})}, nil
	}

	// StakeNFT.Transfer — NFT ownership transfer (ERC721 Transfer event)
	// Guard on address to avoid matching SubnetNFT Transfer (same event signature)
	if lg.Address == idx.chain.StakeNFTAddr {
		if evt, err := stakeNFT.ParseTransfer(lg); err == nil {
			// Skip mint (from=0) and burn (to=0) — handled by Deposited/Withdrawn
			zeroAddr := common.Address{}
			if evt.From != zeroAddr && evt.To != zeroAddr {
				if err := q.UpdateStakePositionOwner(ctx, gen.UpdateStakePositionOwnerParams{
					TokenID: evt.TokenId.Int64(),
					Owner:   strings.ToLower(evt.To.Hex()),
				}); err != nil {
					return nil, fmt.Errorf("UpdateStakePositionOwner: %w", err)
				}
				return []redisEvent{makeEvent("StakeNFTTransfer", lg, map[string]interface{}{
					"from":    evt.From.Hex(),
					"to":      evt.To.Hex(),
					"tokenId": evt.TokenId.String(),
				})}, nil
			}
			return nil, nil
		}
	}

	// Allocated (V2: staker instead of user)
	if evt, err := awpRegistry.ParseAllocated(lg); err == nil {
		if err := q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
			UserAddress:  strings.ToLower(evt.Staker.Hex()),
			AgentAddress: strings.ToLower(evt.Agent.Hex()),
			SubnetID:     evt.SubnetId.Int64(),
			Amount:       bigIntToNumeric(evt.Amount),
			UpdatedBlock: int64(lg.BlockNumber),
		}); err != nil {
			return nil, fmt.Errorf("UpsertStakeAllocation: %w", err)
		}
		if err := q.AddUserAllocated(ctx, gen.AddUserAllocatedParams{
			UserAddress:    strings.ToLower(evt.Staker.Hex()),
			TotalAllocated: bigIntToNumeric(evt.Amount),
			UpdatedBlock:   int64(lg.BlockNumber),
		}); err != nil {
			return nil, fmt.Errorf("AddUserAllocated: %w", err)
		}
		return []redisEvent{makeEvent("Allocated", lg, map[string]interface{}{
			"staker":   evt.Staker.Hex(),
			"agent":    evt.Agent.Hex(),
			"subnetId": evt.SubnetId.String(),
			"amount":   evt.Amount.String(),
			"operator": evt.Operator.Hex(),
		})}, nil
	}

	// Deallocated (V2: staker instead of user)
	if evt, err := awpRegistry.ParseDeallocated(lg); err == nil {
		if err := q.SubtractStakeAllocation(ctx, gen.SubtractStakeAllocationParams{
			UserAddress:  strings.ToLower(evt.Staker.Hex()),
			AgentAddress: strings.ToLower(evt.Agent.Hex()),
			SubnetID:     evt.SubnetId.Int64(),
			Amount:       bigIntToNumeric(evt.Amount),
			UpdatedBlock: int64(lg.BlockNumber),
		}); err != nil {
			return nil, fmt.Errorf("SubtractStakeAllocation: %w", err)
		}
		if err := q.SubtractUserAllocated(ctx, gen.SubtractUserAllocatedParams{
			UserAddress:    strings.ToLower(evt.Staker.Hex()),
			TotalAllocated: bigIntToNumeric(evt.Amount),
			UpdatedBlock:   int64(lg.BlockNumber),
		}); err != nil {
			return nil, fmt.Errorf("SubtractUserAllocated: %w", err)
		}
		return []redisEvent{makeEvent("Deallocated", lg, map[string]interface{}{
			"staker":   evt.Staker.Hex(),
			"agent":    evt.Agent.Hex(),
			"subnetId": evt.SubnetId.String(),
			"amount":   evt.Amount.String(),
			"operator": evt.Operator.Hex(),
		})}, nil
	}

	// Reallocated (V2: staker instead of user)
	if evt, err := awpRegistry.ParseReallocated(lg); err == nil {
		// Subtract from source allocation
		if err := q.SubtractStakeAllocation(ctx, gen.SubtractStakeAllocationParams{
			UserAddress:  strings.ToLower(evt.Staker.Hex()),
			AgentAddress: strings.ToLower(evt.FromAgent.Hex()),
			SubnetID:     evt.FromSubnet.Int64(),
			Amount:       bigIntToNumeric(evt.Amount),
			UpdatedBlock: int64(lg.BlockNumber),
		}); err != nil {
			return nil, fmt.Errorf("SubtractStakeAllocation(Reallocated): %w", err)
		}
		// Add to destination allocation
		if err := q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
			UserAddress:  strings.ToLower(evt.Staker.Hex()),
			AgentAddress: strings.ToLower(evt.ToAgent.Hex()),
			SubnetID:     evt.ToSubnet.Int64(),
			Amount:       bigIntToNumeric(evt.Amount),
			UpdatedBlock: int64(lg.BlockNumber),
		}); err != nil {
			return nil, fmt.Errorf("UpsertStakeAllocation(Reallocated): %w", err)
		}
		return []redisEvent{makeEvent("Reallocated", lg, map[string]interface{}{
			"staker":     evt.Staker.Hex(),
			"fromAgent":  evt.FromAgent.Hex(),
			"fromSubnet": evt.FromSubnet.String(),
			"toAgent":    evt.ToAgent.Hex(),
			"toSubnet":   evt.ToSubnet.String(),
			"amount":     evt.Amount.String(),
			"operator":   evt.Operator.Hex(),
		})}, nil
	}

	// SubnetRegistered
	if evt, err := awpRegistry.ParseSubnetRegistered(lg); err == nil {
		// Read skillsURI and minStake from SubnetNFT on-chain (not included in event)
		skillsURI := ""
		minStake := big.NewInt(0)
		if nftData, nftErr := idx.chain.SubnetNFT.GetSubnetData(nil, evt.SubnetId); nftErr == nil {
			skillsURI = nftData.SkillsURI
			if nftData.MinStake != nil {
				minStake = nftData.MinStake
			}
		}
		if err := q.InsertSubnet(ctx, gen.InsertSubnetParams{
			SubnetID:       evt.SubnetId.Int64(),
			Owner:          strings.ToLower(evt.Owner.Hex()),
			Name:           evt.Name,
			Symbol:         evt.Symbol,
			SubnetContract: strings.ToLower(evt.SubnetManager.Hex()),
			SkillsUri:      pgtype.Text{String: skillsURI, Valid: skillsURI != ""},
			MinStake:       bigIntToNumeric(minStake),
			AlphaToken:     strings.ToLower(evt.AlphaToken.Hex()),
			LpPool:         pgtype.Text{Valid: false},
			CreatedAt:      int64(lg.BlockNumber),
			ImmunityEndsAt: pgtype.Int8{Valid: false},
		}); err != nil {
			return nil, fmt.Errorf("InsertSubnet: %w", err)
		}
		// Mark matching salt as used in pool (best-effort; pool may not contain this address)
		if markErr := q.MarkSaltUsedByAddress(ctx, gen.MarkSaltUsedByAddressParams{
			SubnetID: pgtype.Int8{Int64: evt.SubnetId.Int64(), Valid: true},
			Lower:    strings.ToLower(evt.AlphaToken.Hex()),
		}); markErr != nil {
			slog.Warn("mark salt used failed (non-critical)", "error", markErr, "alphaToken", evt.AlphaToken.Hex())
		}
		return []redisEvent{makeEvent("SubnetRegistered", lg, map[string]interface{}{
			"subnetId":      evt.SubnetId.String(),
			"owner":         evt.Owner.Hex(),
			"name":          evt.Name,
			"symbol":        evt.Symbol,
			"subnetManager": evt.SubnetManager.Hex(),
			"alphaToken":    evt.AlphaToken.Hex(),
		})}, nil
	}

	// LPCreated
	if evt, err := awpRegistry.ParseLPCreated(lg); err == nil {
		poolIdHex := "0x" + hex.EncodeToString(evt.PoolId[:])
		if err := q.UpdateSubnetLP(ctx, gen.UpdateSubnetLPParams{
			SubnetID: evt.SubnetId.Int64(),
			LpPool:   pgtype.Text{String: poolIdHex, Valid: true},
		}); err != nil {
			return nil, fmt.Errorf("UpdateSubnetLP: %w", err)
		}
		return []redisEvent{makeEvent("LPCreated", lg, map[string]interface{}{
			"subnetId":    evt.SubnetId.String(),
			"poolId":      poolIdHex,
			"awpAmount":   evt.AwpAmount.String(),
			"alphaAmount": evt.AlphaAmount.String(),
		})}, nil
	}

	// SubnetActivated
	if evt, err := awpRegistry.ParseSubnetActivated(lg); err == nil {
		if err := q.UpdateSubnetActivated(ctx, gen.UpdateSubnetActivatedParams{
			SubnetID:    evt.SubnetId.Int64(),
			ActivatedAt: pgtype.Int8{Int64: int64(lg.BlockNumber), Valid: true},
		}); err != nil {
			return nil, fmt.Errorf("UpdateSubnetActivated: %w", err)
		}
		return []redisEvent{makeEvent("SubnetActivated", lg, map[string]interface{}{
			"subnetId": evt.SubnetId.String(),
		})}, nil
	}

	// SubnetPaused
	if evt, err := awpRegistry.ParseSubnetPaused(lg); err == nil {
		if err := q.UpdateSubnetStatus(ctx, gen.UpdateSubnetStatusParams{
			SubnetID: evt.SubnetId.Int64(),
			Status:   "Paused",
		}); err != nil {
			return nil, fmt.Errorf("UpdateSubnetStatus(Paused): %w", err)
		}
		return []redisEvent{makeEvent("SubnetPaused", lg, map[string]interface{}{
			"subnetId": evt.SubnetId.String(),
		})}, nil
	}

	// SubnetResumed
	if evt, err := awpRegistry.ParseSubnetResumed(lg); err == nil {
		if err := q.UpdateSubnetStatus(ctx, gen.UpdateSubnetStatusParams{
			SubnetID: evt.SubnetId.Int64(),
			Status:   "Active",
		}); err != nil {
			return nil, fmt.Errorf("UpdateSubnetStatus(Active): %w", err)
		}
		return []redisEvent{makeEvent("SubnetResumed", lg, map[string]interface{}{
			"subnetId": evt.SubnetId.String(),
		})}, nil
	}

	// SubnetBanned
	if evt, err := awpRegistry.ParseSubnetBanned(lg); err == nil {
		if err := q.UpdateSubnetStatus(ctx, gen.UpdateSubnetStatusParams{
			SubnetID: evt.SubnetId.Int64(),
			Status:   "Banned",
		}); err != nil {
			return nil, fmt.Errorf("UpdateSubnetStatus(Banned): %w", err)
		}
		return []redisEvent{makeEvent("SubnetBanned", lg, map[string]interface{}{
			"subnetId": evt.SubnetId.String(),
		})}, nil
	}

	// SubnetUnbanned
	if evt, err := awpRegistry.ParseSubnetUnbanned(lg); err == nil {
		if err := q.UpdateSubnetStatus(ctx, gen.UpdateSubnetStatusParams{
			SubnetID: evt.SubnetId.Int64(),
			Status:   "Active",
		}); err != nil {
			return nil, fmt.Errorf("UpdateSubnetStatus(Active): %w", err)
		}
		return []redisEvent{makeEvent("SubnetUnbanned", lg, map[string]interface{}{
			"subnetId": evt.SubnetId.String(),
		})}, nil
	}

	// SubnetDeregistered
	if evt, err := awpRegistry.ParseSubnetDeregistered(lg); err == nil {
		if err := q.UpdateSubnetBurned(ctx, evt.SubnetId.Int64()); err != nil {
			return nil, fmt.Errorf("UpdateSubnetBurned: %w", err)
		}
		return []redisEvent{makeEvent("SubnetDeregistered", lg, map[string]interface{}{
			"subnetId": evt.SubnetId.String(),
		})}, nil
	}

	// ── SubnetNFT events ──

	subnetNFT := idx.chain.SubnetNFT

	// SkillsURIUpdated (emitted from SubnetNFT)
	if evt, err := subnetNFT.ParseSkillsURIUpdated(lg); err == nil {
		if err := q.UpdateSubnetSkillsURI(ctx, gen.UpdateSubnetSkillsURIParams{
			SubnetID:  evt.TokenId.Int64(),
			SkillsUri: pgtype.Text{String: evt.SkillsURI, Valid: evt.SkillsURI != ""},
		}); err != nil {
			return nil, fmt.Errorf("UpdateSubnetSkillsURI: %w", err)
		}
		return []redisEvent{makeEvent("SkillsURIUpdated", lg, map[string]interface{}{
			"subnetId":  evt.TokenId.String(),
			"skillsURI": evt.SkillsURI,
		})}, nil
	}

	// MinStakeUpdated (emitted from SubnetNFT)
	if evt, err := subnetNFT.ParseMinStakeUpdated(lg); err == nil {
		if err := q.UpdateSubnetMinStake(ctx, gen.UpdateSubnetMinStakeParams{
			SubnetID: evt.TokenId.Int64(),
			MinStake: bigIntToNumeric(evt.MinStake),
		}); err != nil {
			return nil, fmt.Errorf("UpdateSubnetMinStake: %w", err)
		}
		return []redisEvent{makeEvent("MinStakeUpdated", lg, map[string]interface{}{
			"subnetId": evt.TokenId.String(),
			"minStake": evt.MinStake.String(),
		})}, nil
	}

	// SubnetNFT.Transfer — subnet ownership transfer (ERC721 Transfer event)
	// Guard on address to avoid matching StakeNFT Transfer (same event signature)
	if lg.Address == idx.chain.SubnetNFTAddr {
		if evt, err := subnetNFT.ParseTransfer(lg); err == nil {
			// Skip mint (from=0) and burn (to=0) — handled by SubnetRegistered/SubnetDeregistered
			zeroAddr := common.Address{}
			if evt.From != zeroAddr && evt.To != zeroAddr {
				if err := q.UpdateSubnetOwner(ctx, gen.UpdateSubnetOwnerParams{
					SubnetID: evt.TokenId.Int64(),
					Owner:    strings.ToLower(evt.To.Hex()),
				}); err != nil {
					return nil, fmt.Errorf("UpdateSubnetOwner: %w", err)
				}
				return []redisEvent{makeEvent("SubnetNFTTransfer", lg, map[string]interface{}{
					"from":     evt.From.Hex(),
					"to":       evt.To.Hex(),
					"subnetId": evt.TokenId.String(),
				})}, nil
			}
			return nil, nil
		}
	}

	// ── AWPEmission events ──

	// GovernanceWeightUpdated (emitted from AWPEmission) — weight data lives on-chain; only publish Redis event
	if evt, err := awpEmission.ParseGovernanceWeightUpdated(lg); err == nil {
		return []redisEvent{makeEvent("GovernanceWeightUpdated", lg, map[string]interface{}{
			"addr":   evt.Addr.Hex(),
			"weight": evt.Weight.String(),
		})}, nil
	}

	// RecipientAWPDistributed (emitted from AWPEmission)
	if evt, err := awpEmission.ParseRecipientAWPDistributed(lg); err == nil {
		if err := q.InsertRecipientAWPDistribution(ctx, gen.InsertRecipientAWPDistributionParams{
			EpochID:   evt.Epoch.Int64(),
			Recipient: strings.ToLower(evt.Recipient.Hex()),
			AwpAmount: bigIntToNumeric(evt.AwpAmount),
		}); err != nil {
			return nil, fmt.Errorf("InsertRecipientAWPDistribution: %w", err)
		}
		return []redisEvent{makeEvent("RecipientAWPDistributed", lg, map[string]interface{}{
			"epoch":     evt.Epoch.String(),
			"recipient": evt.Recipient.Hex(),
			"awpAmount": evt.AwpAmount.String(),
		})}, nil
	}

	// DAOMatchDistributed (emitted from AWPEmission)
	if evt, err := awpEmission.ParseDAOMatchDistributed(lg); err == nil {
		if err := q.UpdateEpochDAO(ctx, gen.UpdateEpochDAOParams{
			EpochID:     evt.Epoch.Int64(),
			DaoEmission: bigIntToNumeric(evt.Amount),
		}); err != nil {
			return nil, fmt.Errorf("UpdateEpochDAO: %w", err)
		}
		return []redisEvent{makeEvent("DAOMatchDistributed", lg, map[string]interface{}{
			"epoch":  evt.Epoch.String(),
			"amount": evt.Amount.String(),
		})}, nil
	}

	// EpochSettled (emitted from AWPEmission)
	if evt, err := awpEmission.ParseEpochSettled(lg); err == nil {
		// start_time = genesisTime + epochId * epochDuration (time-based epoch, not block number)
		epochStartTime := idx.genesisTime + evt.Epoch.Int64()*idx.epochDuration
		if err := q.UpsertEpoch(ctx, gen.UpsertEpochParams{
			EpochID:       evt.Epoch.Int64(),
			StartTime:     epochStartTime,
			DailyEmission: bigIntToNumeric(evt.TotalEmission),
			DaoEmission:   pgtype.Numeric{Valid: false},
		}); err != nil {
			return nil, fmt.Errorf("UpsertEpoch: %w", err)
		}
		return []redisEvent{makeEvent("EpochSettled", lg, map[string]interface{}{
			"epoch":          evt.Epoch.String(),
			"totalEmission":  evt.TotalEmission.String(),
			"recipientCount": evt.RecipientCount.String(),
		})}, nil
	}

	// AllocationsSubmitted (emitted from AWPEmission)
	if evt, err := awpEmission.ParseAllocationsSubmitted(lg); err == nil {
		addrs := make([]string, len(evt.Recipients))
		for i, a := range evt.Recipients {
			addrs[i] = a.Hex()
		}
		ws := make([]string, len(evt.Weights))
		for i, w := range evt.Weights {
			ws[i] = w.String()
		}
		return []redisEvent{makeEvent("AllocationsSubmitted", lg, map[string]interface{}{
			"nonce":      evt.Nonce.String(),
			"recipients": addrs,
			"weights":    ws,
		})}, nil
	}

	// OracleConfigUpdated (emitted from AWPEmission)
	if evt, err := awpEmission.ParseOracleConfigUpdated(lg); err == nil {
		addrs := make([]string, len(evt.Oracles))
		for i, o := range evt.Oracles {
			addrs[i] = o.Hex()
		}
		return []redisEvent{makeEvent("OracleConfigUpdated", lg, map[string]interface{}{
			"oracles":   addrs,
			"threshold": evt.Threshold.String(),
		})}, nil
	}

	// ── AWPRegistry governance events (notification-only, no DB writes) ──
	if lg.Address == idx.chain.AWPRegistryAddr {
		topic := lg.Topics[0]
		switch topic {
		// GuardianUpdated(address indexed newGuardian)
		case common.HexToHash("0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf"):
			newGuardian := common.BytesToAddress(lg.Topics[1].Bytes())
			return []redisEvent{makeEvent("GuardianUpdated", lg, map[string]interface{}{
				"newGuardian": newGuardian.Hex(),
			})}, nil
		// InitialAlphaPriceUpdated(uint256 newPrice)
		case common.HexToHash("0xab7ee876750d22d253d0b38988caea5f6285a832697e4889d9beb36515dde34e"):
			newPrice := new(big.Int).SetBytes(lg.Data)
			return []redisEvent{makeEvent("InitialAlphaPriceUpdated", lg, map[string]interface{}{
				"newPrice": newPrice.String(),
			})}, nil
		// ImmunityPeriodUpdated(uint256 newPeriod)
		case common.HexToHash("0x49b186851943e5bbcefec9411c3238262c6e102e4000142f8f060143d1b8724c"):
			newPeriod := new(big.Int).SetBytes(lg.Data)
			return []redisEvent{makeEvent("ImmunityPeriodUpdated", lg, map[string]interface{}{
				"newPeriod": newPeriod.String(),
			})}, nil
		// AlphaTokenFactoryUpdated(address indexed newFactory)
		case common.HexToHash("0xca3b5054bdfbf81973dd36029b7ef8c5479d0739433700df6b2e6d690ead4a3e"):
			newFactory := common.BytesToAddress(lg.Topics[1].Bytes())
			return []redisEvent{makeEvent("AlphaTokenFactoryUpdated", lg, map[string]interface{}{
				"newFactory": newFactory.Hex(),
			})}, nil
		// DefaultSubnetManagerImplUpdated(address indexed newImpl)
		case common.HexToHash("0xa37cb79f631c6bb2a11d965d06cce40e3c936eba1649879b8ffa233c0219f949"):
			newImpl := common.BytesToAddress(lg.Topics[1].Bytes())
			return []redisEvent{makeEvent("DefaultSubnetManagerImplUpdated", lg, map[string]interface{}{
				"newImpl": newImpl.Hex(),
			})}, nil
		}
	}

	// ── AWPDAO events (notification-only, no DB writes) ──

	awpDAO := idx.chain.AWPDAO

	// ProposalCreated
	if evt, err := awpDAO.ParseProposalCreated(lg); err == nil {
		return []redisEvent{makeEvent("ProposalCreated", lg, map[string]interface{}{
			"proposalId":  evt.ProposalId.String(),
			"proposer":    evt.Proposer.Hex(),
			"voteStart":   evt.VoteStart.String(),
			"voteEnd":     evt.VoteEnd.String(),
			"description": evt.Description,
		})}, nil
	}

	// VoteCast
	if evt, err := awpDAO.ParseVoteCast(lg); err == nil {
		return []redisEvent{makeEvent("VoteCast", lg, map[string]interface{}{
			"voter":      evt.Voter.Hex(),
			"proposalId": evt.ProposalId.String(),
			"support":    evt.Support,
			"weight":     evt.Weight.String(),
			"reason":     evt.Reason,
		})}, nil
	}

	// ProposalExecuted
	if evt, err := awpDAO.ParseProposalExecuted(lg); err == nil {
		return []redisEvent{makeEvent("ProposalExecuted", lg, map[string]interface{}{
			"proposalId": evt.ProposalId.String(),
		})}, nil
	}

	// ProposalCanceled
	if evt, err := awpDAO.ParseProposalCanceled(lg); err == nil {
		return []redisEvent{makeEvent("ProposalCanceled", lg, map[string]interface{}{
			"proposalId": evt.ProposalId.String(),
		})}, nil
	}

	// ProposalQueued
	if evt, err := awpDAO.ParseProposalQueued(lg); err == nil {
		return []redisEvent{makeEvent("ProposalQueued", lg, map[string]interface{}{
			"proposalId": evt.ProposalId.String(),
			"etaSeconds": evt.EtaSeconds.String(),
		})}, nil
	}

	// Unrecognized event (may be Paused/Unpaused or other events that don't need to be stored)
	return nil, nil
}

// bigIntToNumeric converts a *big.Int to pgtype.Numeric
func bigIntToNumeric(v *big.Int) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   new(big.Int).Set(v),
		Exp:   0,
		Valid: true,
	}
}

// makeEvent constructs an event structure to be published to Redis
func makeEvent(eventType string, lg types.Log, data map[string]interface{}) redisEvent {
	return redisEvent{
		Type:        eventType,
		BlockNumber: lg.BlockNumber,
		TxHash:      lg.TxHash.Hex(),
		Data:        data,
	}
}
