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

// poll executes one complete scan cycle
func (idx *Indexer) poll(ctx context.Context) error {
	// 1. Read the last synced block
	q := gen.New(idx.pool)
	state, err := q.GetSyncState(ctx, syncKey)
	if err != nil {
		// On first run sync_states may have no record; start from deploy block
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

	// 2. Get the latest block number on-chain
	latestBlock, err := idx.chain.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}
	if fromBlock > latestBlock {
		return nil // already at the latest, nothing to do
	}

	// Confirmation depth: lag behind chain tip to avoid reorg-affected blocks
	const confirmationDepth = 15
	if latestBlock > confirmationDepth {
		latestBlock -= confirmationDepth
	}
	if fromBlock > latestBlock {
		return nil
	}

	// Limit single query range to avoid RPC restrictions
	const maxBlockRange = 500
	toBlock := latestBlock
	if toBlock-fromBlock > maxBlockRange {
		toBlock = fromBlock + maxBlockRange
	}

	slog.Info("scanning blocks", "from", fromBlock, "to", toBlock)

	// 3. Filter logs from RootNet, AWPEmission, StakeNFT, and SubnetNFT contracts
	logs, err := idx.chain.Eth.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
		Addresses: []common.Address{idx.chain.RootNetAddr, idx.chain.AWPEmissionAddr, idx.chain.StakeNFTAddr, idx.chain.SubnetNFTAddr},
	})
	if err != nil {
		return fmt.Errorf("failed to filter logs: %w", err)
	}

	// 4. Process all events within a PostgreSQL transaction
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
			slog.Warn("failed to process log", "txHash", lg.TxHash.Hex(), "topic", lg.Topics[0].Hex(), "error", err)
			continue
		}
		events = append(events, evts...)
	}

	// 5. Update sync progress
	if err := qtx.UpsertSyncState(ctx, gen.UpsertSyncStateParams{
		ContractName: syncKey,
		LastBlock:    int64(toBlock),
	}); err != nil {
		return fmt.Errorf("failed to update sync state: %w", err)
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

// processLog parses a single log entry and performs the corresponding database writes, returning Redis events to publish
func (idx *Indexer) processLog(ctx context.Context, q *gen.Queries, lg types.Log) ([]redisEvent, error) {
	rootNet := idx.chain.RootNet
	awpEmission := idx.chain.AWPEmission
	stakeNFT := idx.chain.StakeNFT

	// Attempt to match each event signature
	// UserRegistered
	if evt, err := rootNet.ParseUserRegistered(lg); err == nil {
		if err := q.InsertUser(ctx, gen.InsertUserParams{
			Address:      strings.ToLower(evt.User.Hex()),
			RegisteredAt: int64(lg.BlockNumber),
		}); err != nil {
			return nil, fmt.Errorf("InsertUser: %w", err)
		}
		if err := q.InitUserBalance(ctx, strings.ToLower(evt.User.Hex())); err != nil {
			return nil, fmt.Errorf("InitUserBalance: %w", err)
		}
		return []redisEvent{makeEvent("UserRegistered", lg, map[string]interface{}{
			"user": evt.User.Hex(),
		})}, nil
	}

	// AgentBound (covers both first bind and rebind)
	if evt, err := rootNet.ParseAgentBound(lg); err == nil {
		if err := q.UpsertAgent(ctx, gen.UpsertAgentParams{
			AgentAddress: strings.ToLower(evt.Agent.Hex()),
			OwnerAddress: strings.ToLower(evt.Principal.Hex()),
		}); err != nil {
			return nil, fmt.Errorf("UpsertAgent: %w", err)
		}
		return []redisEvent{makeEvent("AgentBound", lg, map[string]interface{}{
			"principal":    evt.Principal.Hex(),
			"agent":        evt.Agent.Hex(),
			"oldPrincipal": evt.OldPrincipal.Hex(),
		})}, nil
	}

	// AgentRemoved
	if evt, err := rootNet.ParseAgentRemoved(lg); err == nil {
		if err := q.UpdateAgentRemoved(ctx, gen.UpdateAgentRemovedParams{
			AgentAddress: strings.ToLower(evt.Agent.Hex()),
			RemovedAt:    pgtype.Int8{Int64: int64(lg.BlockNumber), Valid: true},
		}); err != nil {
			return nil, fmt.Errorf("UpdateAgentRemoved: %w", err)
		}
		if err := q.FreezeAgentAllocations(ctx, gen.FreezeAgentAllocationsParams{
			UserAddress:  strings.ToLower(evt.User.Hex()),
			AgentAddress: strings.ToLower(evt.Agent.Hex()),
		}); err != nil {
			return nil, fmt.Errorf("FreezeAgentAllocations: %w", err)
		}
		return []redisEvent{makeEvent("AgentRemoved", lg, map[string]interface{}{
			"user":     evt.User.Hex(),
			"agent":    evt.Agent.Hex(),
			"operator": evt.Operator.Hex(),
		})}, nil
	}

	// AgentUnbound (agent voluntarily unbinds from principal — freeze allocations same as AgentRemoved)
	if evt, err := rootNet.ParseAgentUnbound(lg); err == nil {
		if err := q.UpdateAgentRemoved(ctx, gen.UpdateAgentRemovedParams{
			AgentAddress: strings.ToLower(evt.Agent.Hex()),
			RemovedAt:    pgtype.Int8{Int64: int64(lg.BlockNumber), Valid: true},
		}); err != nil {
			return nil, fmt.Errorf("UpdateAgentRemoved (unbound): %w", err)
		}
		if err := q.FreezeAgentAllocations(ctx, gen.FreezeAgentAllocationsParams{
			UserAddress:  strings.ToLower(evt.Principal.Hex()),
			AgentAddress: strings.ToLower(evt.Agent.Hex()),
		}); err != nil {
			return nil, fmt.Errorf("FreezeAgentAllocations (unbound): %w", err)
		}
		return []redisEvent{makeEvent("AgentUnbound", lg, map[string]interface{}{
			"principal": evt.Principal.Hex(),
			"agent":     evt.Agent.Hex(),
		})}, nil
	}

	// DelegationUpdated
	if evt, err := rootNet.ParseDelegationUpdated(lg); err == nil {
		if err := q.UpdateAgentManager(ctx, gen.UpdateAgentManagerParams{
			AgentAddress: strings.ToLower(evt.Agent.Hex()),
			IsManager:    evt.IsManager,
		}); err != nil {
			return nil, fmt.Errorf("UpdateAgentManager: %w", err)
		}
		return []redisEvent{makeEvent("DelegationUpdated", lg, map[string]interface{}{
			"user":      evt.User.Hex(),
			"agent":     evt.Agent.Hex(),
			"isManager": evt.IsManager,
			"operator":  evt.Operator.Hex(),
		})}, nil
	}

	// RewardRecipientUpdated
	if evt, err := rootNet.ParseRewardRecipientUpdated(lg); err == nil {
		if err := q.UpsertRewardRecipient(ctx, gen.UpsertRewardRecipientParams{
			UserAddress:      strings.ToLower(evt.User.Hex()),
			RecipientAddress: strings.ToLower(evt.Recipient.Hex()),
		}); err != nil {
			return nil, fmt.Errorf("UpsertRewardRecipient: %w", err)
		}
		return []redisEvent{makeEvent("RewardRecipientUpdated", lg, map[string]interface{}{
			"user":      evt.User.Hex(),
			"recipient": evt.Recipient.Hex(),
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
			"tokenId":       evt.TokenId.String(),
			"addedAmount":   evt.AddedAmount.String(),
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

	// Allocated
	if evt, err := rootNet.ParseAllocated(lg); err == nil {
		if err := q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
			UserAddress:  strings.ToLower(evt.User.Hex()),
			AgentAddress: strings.ToLower(evt.Agent.Hex()),
			SubnetID:     evt.SubnetId.Int64(),
			Amount:       bigIntToNumeric(evt.Amount),
		}); err != nil {
			return nil, fmt.Errorf("UpsertStakeAllocation: %w", err)
		}
		if err := q.AddUserAllocated(ctx, gen.AddUserAllocatedParams{
			UserAddress:    strings.ToLower(evt.User.Hex()),
			TotalAllocated: bigIntToNumeric(evt.Amount),
		}); err != nil {
			return nil, fmt.Errorf("AddUserAllocated: %w", err)
		}
		return []redisEvent{makeEvent("Allocated", lg, map[string]interface{}{
			"user":     evt.User.Hex(),
			"agent":    evt.Agent.Hex(),
			"subnetId": evt.SubnetId.String(),
			"amount":   evt.Amount.String(),
			"operator": evt.Operator.Hex(),
		})}, nil
	}

	// Deallocated
	if evt, err := rootNet.ParseDeallocated(lg); err == nil {
		if err := q.SubtractStakeAllocation(ctx, gen.SubtractStakeAllocationParams{
			UserAddress:  strings.ToLower(evt.User.Hex()),
			AgentAddress: strings.ToLower(evt.Agent.Hex()),
			SubnetID:     evt.SubnetId.Int64(),
			Amount:       bigIntToNumeric(evt.Amount),
		}); err != nil {
			return nil, fmt.Errorf("SubtractStakeAllocation: %w", err)
		}
		if err := q.SubtractUserAllocated(ctx, gen.SubtractUserAllocatedParams{
			UserAddress:    strings.ToLower(evt.User.Hex()),
			TotalAllocated: bigIntToNumeric(evt.Amount),
		}); err != nil {
			return nil, fmt.Errorf("SubtractUserAllocated: %w", err)
		}
		return []redisEvent{makeEvent("Deallocated", lg, map[string]interface{}{
			"user":     evt.User.Hex(),
			"agent":    evt.Agent.Hex(),
			"subnetId": evt.SubnetId.String(),
			"amount":   evt.Amount.String(),
			"operator": evt.Operator.Hex(),
		})}, nil
	}

	// Reallocated — replaces the old ReallocationQueued; update DB allocations and publish event
	if evt, err := rootNet.ParseReallocated(lg); err == nil {
		// Subtract from source allocation
		if err := q.SubtractStakeAllocation(ctx, gen.SubtractStakeAllocationParams{
			UserAddress:  strings.ToLower(evt.User.Hex()),
			AgentAddress: strings.ToLower(evt.FromAgent.Hex()),
			SubnetID:     evt.FromSubnet.Int64(),
			Amount:       bigIntToNumeric(evt.Amount),
		}); err != nil {
			return nil, fmt.Errorf("SubtractStakeAllocation(Reallocated): %w", err)
		}
		// Add to destination allocation
		if err := q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
			UserAddress:  strings.ToLower(evt.User.Hex()),
			AgentAddress: strings.ToLower(evt.ToAgent.Hex()),
			SubnetID:     evt.ToSubnet.Int64(),
			Amount:       bigIntToNumeric(evt.Amount),
		}); err != nil {
			return nil, fmt.Errorf("UpsertStakeAllocation(Reallocated): %w", err)
		}
		return []redisEvent{makeEvent("Reallocated", lg, map[string]interface{}{
			"user":       evt.User.Hex(),
			"fromAgent":  evt.FromAgent.Hex(),
			"fromSubnet": evt.FromSubnet.String(),
			"toAgent":    evt.ToAgent.Hex(),
			"toSubnet":   evt.ToSubnet.String(),
			"amount":     evt.Amount.String(),
			"operator":   evt.Operator.Hex(),
		})}, nil
	}

	// SubnetRegistered
	if evt, err := rootNet.ParseSubnetRegistered(lg); err == nil {
		if err := q.InsertSubnet(ctx, gen.InsertSubnetParams{
			SubnetID:       evt.SubnetId.Int64(),
			Owner:          strings.ToLower(evt.Owner.Hex()),
			Name:           evt.Name,
			Symbol:         evt.Symbol,
			MetadataUri:    pgtype.Text{String: evt.MetadataURI, Valid: evt.MetadataURI != ""},
			SubnetContract: strings.ToLower(evt.SubnetManager.Hex()),
			CoordinatorUrl: pgtype.Text{String: evt.CoordinatorURL, Valid: evt.CoordinatorURL != ""},
			SkillsUri:      pgtype.Text{Valid: false},
			MinStake:       bigIntToNumeric(big.NewInt(0)), // default 0; updated via SubnetNFT.MinStakeUpdated event
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
			"subnetId":       evt.SubnetId.String(),
			"owner":          evt.Owner.Hex(),
			"name":           evt.Name,
			"symbol":         evt.Symbol,
			"metadataURI":    evt.MetadataURI,
			"subnetManager":  evt.SubnetManager.Hex(),
			"alphaToken":     evt.AlphaToken.Hex(),
			"coordinatorURL": evt.CoordinatorURL,
		})}, nil
	}

	// LPCreated
	if evt, err := rootNet.ParseLPCreated(lg); err == nil {
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

	// MetadataUpdated
	if evt, err := rootNet.ParseMetadataUpdated(lg); err == nil {
		if err := q.UpdateSubnetMetadata(ctx, gen.UpdateSubnetMetadataParams{
			SubnetID:       evt.SubnetId.Int64(),
			MetadataUri:    pgtype.Text{String: evt.MetadataURI, Valid: evt.MetadataURI != ""},
			CoordinatorUrl: pgtype.Text{String: evt.CoordinatorURL, Valid: evt.CoordinatorURL != ""},
		}); err != nil {
			return nil, fmt.Errorf("UpdateSubnetMetadata: %w", err)
		}
		return []redisEvent{makeEvent("MetadataUpdated", lg, map[string]interface{}{
			"subnetId":       evt.SubnetId.String(),
			"metadataURI":    evt.MetadataURI,
			"coordinatorURL": evt.CoordinatorURL,
		})}, nil
	}

	// SubnetActivated
	if evt, err := rootNet.ParseSubnetActivated(lg); err == nil {
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
	if evt, err := rootNet.ParseSubnetPaused(lg); err == nil {
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
	if evt, err := rootNet.ParseSubnetResumed(lg); err == nil {
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
	if evt, err := rootNet.ParseSubnetBanned(lg); err == nil {
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
	if evt, err := rootNet.ParseSubnetUnbanned(lg); err == nil {
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
	if evt, err := rootNet.ParseSubnetDeregistered(lg); err == nil {
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
			EpochID:        evt.Epoch.Int64(),
			StartTime:      epochStartTime,
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
