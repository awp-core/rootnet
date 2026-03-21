package chain

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/cortexia/rootnet/api/internal/chain/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

// Keeper is the on-chain scheduled task executor responsible for settling epochs and updating caches
type Keeper struct {
	client      *ethclient.Client
	awpEmission *bindings.AWPEmission
	awpToken    *bindings.AWPToken
	key         *ecdsa.PrivateKey
	chainID     *big.Int
	cron        *cron.Cron
	redis       *redis.Client
	logger      *slog.Logger
	cancel      context.CancelFunc

	// Serializes sendSettleEpoch to prevent concurrent nonce collisions
	txMu sync.Mutex

	// Cached RPC results — shared between trySettleEpoch and updateTokenPrices within the same tick
	cacheMu            sync.Mutex
	cachedCurrentEpoch *big.Int
	cachedSettledEpoch *big.Int
	cacheTime          time.Time
}

// emissionInfo is used to cache the current emission info in Redis
type emissionInfo struct {
	Epoch         string `json:"epoch"`
	SettledEpoch  string `json:"settledEpoch"`
	DailyEmission string `json:"dailyEmission"`
	TotalWeight   string `json:"totalWeight"`
}

// awpInfo is used to cache AWP token info in Redis
type awpInfo struct {
	TotalSupply string `json:"totalSupply"`
	MaxSupply   string `json:"maxSupply"`
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	client *ethclient.Client,
	awpEmissionAddr common.Address,
	awpTokenAddr common.Address,
	key *ecdsa.PrivateKey,
	chainID *big.Int,
	rdb *redis.Client,
	logger *slog.Logger,
) (*Keeper, error) {
	awpEmission, err := bindings.NewAWPEmission(awpEmissionAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to bind AWPEmission contract: %w", err)
	}

	awpToken, err := bindings.NewAWPToken(awpTokenAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to bind AWPToken contract: %w", err)
	}

	return &Keeper{
		client:      client,
		awpEmission: awpEmission,
		awpToken:    awpToken,
		key:         key,
		chainID:     chainID,
		cron:        cron.New(),
		redis:       rdb,
		logger:      logger,
	}, nil
}

// Start launches the scheduled task scheduler.
func (k *Keeper) Start(_ context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	k.cancel = cancel

	// Attempt to settle epoch every 30 seconds
	if _, err := k.cron.AddFunc("@every 30s", func() {
		k.trySettleEpoch(ctx)
	}); err != nil {
		k.logger.Error("failed to add cron", "error", err)
	}

	// Update token price cache every 25s to keep ahead of the shortest Redis TTL (30s for emission_current)
	if _, err := k.cron.AddFunc("@every 25s", func() {
		k.updateTokenPrices(ctx)
	}); err != nil {
		k.logger.Error("failed to add cron", "error", err)
	}

	k.cron.Start()
	k.logger.Info("Keeper scheduled tasks started")
}

// Stop halts the scheduled task scheduler
func (k *Keeper) Stop() {
	if k.cancel != nil {
		k.cancel()
	}
	k.cron.Stop()
	k.logger.Info("Keeper scheduled tasks stopped")
}

// auth builds an on-chain transaction signer
func (k *Keeper) auth() (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(k.key, k.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction signer: %w", err)
	}
	return auth, nil
}

// fetchEpochs returns cached currentEpoch and settledEpoch, refreshing from RPC if stale (>10s).
// Thread-safe: protected by cacheMu since trySettleEpoch and updateTokenPrices run in concurrent cron goroutines.
func (k *Keeper) fetchEpochs() (currentEpoch, settledEpoch *big.Int, err error) {
	k.cacheMu.Lock()
	defer k.cacheMu.Unlock()

	if time.Since(k.cacheTime) < 10*time.Second && k.cachedCurrentEpoch != nil {
		return k.cachedCurrentEpoch, k.cachedSettledEpoch, nil
	}
	currentEpoch, err = k.awpEmission.CurrentEpoch(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("CurrentEpoch: %w", err)
	}
	settledEpoch, err = k.awpEmission.SettledEpoch(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("SettledEpoch: %w", err)
	}
	k.cachedCurrentEpoch = currentEpoch
	k.cachedSettledEpoch = settledEpoch
	k.cacheTime = time.Now()
	return currentEpoch, settledEpoch, nil
}

// trySettleEpoch attempts to settle the current epoch (calls AWPEmission.SettleEpoch)
func (k *Keeper) trySettleEpoch(ctx context.Context) {
	settleProgress, err := k.awpEmission.SettleProgress(nil)
	if err != nil {
		k.logger.Error("failed to read SettleProgress state", "error", err)
		return
	}

	// If settlement is in progress (progress > 0), call settleEpoch to continue processing
	if settleProgress.Cmp(big.NewInt(0)) > 0 {
		k.logger.Info("epoch settlement in progress, continuing settleEpoch")
		k.sendSettleEpoch(ctx)
		return
	}

	// Compare settledEpoch vs currentEpoch (cached, shared with updateTokenPrices)
	currentEpoch, settledEpoch, err := k.fetchEpochs()
	if err != nil {
		k.logger.Error("failed to fetch epoch data", "error", err)
		return
	}

	if currentEpoch.Cmp(settledEpoch) > 0 {
		k.logger.Info("new epoch ready, executing settleEpoch",
			"settledEpoch", settledEpoch.String(),
			"currentEpoch", currentEpoch.String(),
		)
		k.sendSettleEpoch(ctx)
	}
}

// sendSettleEpoch sends the settle-epoch transaction (calls AWPEmission.SettleEpoch)
func (k *Keeper) sendSettleEpoch(ctx context.Context) {
	k.txMu.Lock()
	defer k.txMu.Unlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	txAuth, err := k.auth()
	if err != nil {
		k.logger.Error("failed to build transaction signer", "error", err)
		return
	}
	txAuth.Context = timeoutCtx

	tx, err := k.awpEmission.SettleEpoch(txAuth, big.NewInt(200))
	if err != nil {
		k.logger.Error("failed to call SettleEpoch", "error", err)
		return
	}

	k.logger.Info("SettleEpoch transaction sent", "txHash", tx.Hash().Hex())
}

// updateTokenPrices updates token-related caches in Redis
func (k *Keeper) updateTokenPrices(ctx context.Context) {
	// Read epoch data (cached, shared with trySettleEpoch — avoids redundant RPC calls)
	currentEpoch, settledEpoch, err := k.fetchEpochs()
	if err != nil {
		k.logger.Error("failed to fetch epoch data", "error", err)
		return
	}

	dailyEmission, err := k.awpEmission.CurrentDailyEmission(nil)
	if err != nil {
		k.logger.Error("failed to read CurrentDailyEmission", "error", err)
		return
	}

	totalWeight, err := k.awpEmission.GetTotalWeight(nil)
	if err != nil {
		k.logger.Error("failed to read GetTotalWeight", "error", err)
		return
	}

	// Write emission_current cache, TTL=30s (per Redis Key Spec)
	emData, err := json.Marshal(emissionInfo{
		Epoch:         currentEpoch.String(),
		SettledEpoch:  settledEpoch.String(),
		DailyEmission: dailyEmission.String(),
		TotalWeight:   totalWeight.String(),
	})
	if err != nil {
		k.logger.Error("failed to serialize emission info", "error", err)
		return
	}

	if err := k.redis.Set(ctx, "emission_current", emData, 30*time.Second).Err(); err != nil {
		k.logger.Error("failed to write emission_current cache", "error", err)
	}

	// Write awp_info cache, TTL=1m (per Redis Key Spec)
	totalSupply, err := k.awpToken.TotalSupply(nil)
	if err != nil {
		k.logger.Error("failed to read AWPToken TotalSupply", "error", err)
		return
	}
	maxSupply, err := k.awpToken.MAXSUPPLY(nil)
	if err != nil {
		k.logger.Error("failed to read AWPToken MaxSupply", "error", err)
		return
	}
	awpData, err := json.Marshal(awpInfo{
		TotalSupply: totalSupply.String(),
		MaxSupply:   maxSupply.String(),
	})
	if err != nil {
		k.logger.Error("failed to serialize AWP info", "error", err)
		return
	}

	if err := k.redis.Set(ctx, "awp_info", awpData, 1*time.Minute).Err(); err != nil {
		k.logger.Error("failed to write awp_info cache", "error", err)
	}

	// TODO: Implement alpha price updates per Redis Key Spec
	// alpha_price:{subnetId} → JSON, TTL=10m
	// Requires CLPoolManager bindings to read sqrtPriceX96 from on-chain pools

	k.logger.Debug("token caches updated",
		"epoch", currentEpoch.String(),
		"settledEpoch", settledEpoch.String(),
		"dailyEmission", dailyEmission.String(),
		"totalWeight", totalWeight.String(),
	)
}
