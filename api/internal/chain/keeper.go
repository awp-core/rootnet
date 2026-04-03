package chain

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/cortexia/rootnet/api/internal/chain/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

// Keeper is the on-chain scheduled task executor responsible for settling epochs, updating caches, and compounding LP fees
type Keeper struct {
	client      *ethclient.Client
	awpEmission *bindings.AWPEmission
	awpToken    *bindings.AWPToken
	lpManager   *bindings.LPManagerBase
	poolManager *bindings.PoolManagerReader
	key         *ecdsa.PrivateKey
	chainID     *big.Int
	chainIDInt  int64
	cron        *cron.Cron
	redis       *redis.Client
	logger      *slog.Logger
	cancel      context.CancelFunc

	// Distributed lock value (unique per instance, for safe release)
	lockValue string

	// When true, skip settleEpoch (read-only mode: cache updates only)
	skipSettle bool

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
	lpManagerAddr common.Address,
	poolManagerAddr common.Address,
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

	var lpMgr *bindings.LPManagerBase
	if lpManagerAddr != (common.Address{}) {
		lpMgr, err = bindings.NewLPManagerBase(lpManagerAddr, client)
		if err != nil {
			return nil, fmt.Errorf("failed to bind LPManager contract: %w", err)
		}
	}

	var poolMgr *bindings.PoolManagerReader
	if poolManagerAddr != (common.Address{}) {
		poolMgr, err = bindings.NewPoolManagerReader(poolManagerAddr, client)
		if err != nil {
			return nil, fmt.Errorf("failed to bind PoolManager contract: %w", err)
		}
	}

	return &Keeper{
		client:      client,
		awpEmission: awpEmission,
		awpToken:    awpToken,
		lpManager:   lpMgr,
		poolManager: poolMgr,
		key:         key,
		chainID:     chainID,
		chainIDInt:  chainID.Int64(),
		cron:        cron.New(),
		redis:       rdb,
		logger:      logger,
	}, nil
}

// SetSkipSettle enables/disables settle-epoch execution (read-only mode when true)
func (k *Keeper) SetSkipSettle(skip bool) { k.skipSettle = skip }

// acquireLock acquires a Redis distributed lock to prevent concurrent keeper instances.
// Uses a unique lock value so releaseLock only deletes our own lock.
func (k *Keeper) acquireLock(ctx context.Context) bool {
	lockKey := fmt.Sprintf("keeper:lock:%d", k.chainIDInt)
	k.lockValue = fmt.Sprintf("%d-%d", os.Getpid(), time.Now().UnixNano())
	ok, err := k.redis.SetNX(ctx, lockKey, k.lockValue, 90*time.Second).Result()
	if err != nil {
		k.logger.Warn("failed to acquire keeper lock", "error", err)
		return false
	}
	return ok
}

// renewLock extends the distributed lock TTL only if we still own it.
func (k *Keeper) renewLock(ctx context.Context) {
	lockKey := fmt.Sprintf("keeper:lock:%d", k.chainIDInt)
	script := `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("EXPIRE", KEYS[1], ARGV[2]) else return 0 end`
	result, _ := k.redis.Eval(ctx, script, []string{lockKey}, k.lockValue, 90).Int64()
	if result == 0 {
		k.logger.Warn("lost keeper lock, another instance may have taken over", "chainId", k.chainIDInt)
	}
}

// releaseLock releases the distributed lock only if we own it (safe against stale lock deletion).
func (k *Keeper) releaseLock() {
	lockKey := fmt.Sprintf("keeper:lock:%d", k.chainIDInt)
	// Lua: only delete if the lock value matches ours
	script := `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) else return 0 end`
	k.redis.Eval(context.Background(), script, []string{lockKey}, k.lockValue)
}

// Start launches the scheduled task scheduler.
func (k *Keeper) Start(_ context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	k.cancel = cancel

	// Distributed lock: prevent multiple keeper instances from running simultaneously
	if !k.acquireLock(ctx) {
		k.logger.Error("another keeper instance is already running for this chain, exiting", "chainId", k.chainIDInt)
		return
	}

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

	// Update relayer native balance every 25s (aligned with token price updates)
	if _, err := k.cron.AddFunc("@every 25s", func() {
		k.updateRelayerBalance(ctx)
	}); err != nil {
		k.logger.Error("failed to add relayer balance cron", "error", err)
	}

	// Compound LP fees every 24 hours (if LPManager is configured)
	if k.lpManager != nil {
		if _, err := k.cron.AddFunc("@every 24h", func() {
			k.compoundAllFees(ctx)
		}); err != nil {
			k.logger.Error("failed to add compound fees cron", "error", err)
		}
	}

	// Renew distributed lock every 30s (TTL 60s, ensures lock does not expire)
	if _, err := k.cron.AddFunc("@every 30s", func() {
		k.renewLock(ctx)
	}); err != nil {
		k.logger.Error("failed to add lock renewal cron", "error", err)
	}

	k.cron.Start()
	k.logger.Info("Keeper scheduled tasks started")
}

// Stop halts the scheduled task scheduler and releases the distributed lock
func (k *Keeper) Stop() {
	if k.cancel != nil {
		k.cancel()
	}
	k.cron.Stop()
	k.releaseLock()
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

	if currentEpoch.Cmp(settledEpoch) >= 0 && !k.skipSettle {
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

	if err := k.redis.Set(ctx, fmt.Sprintf("emission_current:%d", k.chainIDInt), emData, 30*time.Second).Err(); err != nil {
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

	if err := k.redis.Set(ctx, fmt.Sprintf("awp_info:%d", k.chainIDInt), awpData, 1*time.Minute).Err(); err != nil {
		k.logger.Error("failed to write awp_info cache", "error", err)
	}

	// Update alpha token prices from on-chain pool state
	if k.lpManager != nil && k.poolManager != nil {
		k.updateAlphaPrices(ctx)
	}

	k.logger.Debug("token caches updated",
		"epoch", currentEpoch.String(),
		"settledEpoch", settledEpoch.String(),
		"dailyEmission", dailyEmission.String(),
		"totalWeight", totalWeight.String(),
	)
}

// alphaPrice is the cached price data for an Alpha token
type alphaPrice struct {
	PriceInAWP string `json:"priceInAWP"` // Alpha price denominated in AWP
	SqrtPriceX96 string `json:"sqrtPriceX96"` // raw sqrtPriceX96 for frontend use
	UpdatedAt  int64  `json:"updatedAt"`  // unix timestamp
}

// alphaPriceResult holds the computed price data for a single alpha token (used by worker goroutines)
type alphaPriceResult struct {
	key  string
	data []byte
}

// updateAlphaPrices reads sqrtPriceX96 from each active subnet's LP pool and caches the price.
// Called every 25s as part of updateTokenPrices.
// Uses a bounded worker pool (max 10 concurrent) to parallelize RPC calls across tokens.
// Redis keys are keyed by worknetId (alpha_price:{worknetId}) per spec.
func (k *Keeper) updateAlphaPrices(ctx context.Context) {
	// Read worknetId→alphaToken map from Redis (populated by indexer)
	subnetMapKey := fmt.Sprintf("active_alpha_subnet_map:%d", k.chainIDInt)
	data, err := k.redis.Get(ctx, subnetMapKey).Result()
	if err != nil {
		return // no active tokens, skip
	}

	var subnetMap map[string]string // worknetId → alphaToken
	if err := json.Unmarshal([]byte(data), &subnetMap); err != nil {
		return
	}

	if len(subnetMap) == 0 {
		return
	}

	type subnetToken struct {
		subnetID string
		token    string
	}
	entries := make([]subnetToken, 0, len(subnetMap))
	for sid, tok := range subnetMap {
		entries = append(entries, subnetToken{subnetID: sid, token: tok})
	}

	// Bounded worker pool: max 10 concurrent RPC goroutines
	const maxWorkers = 10
	sem := make(chan struct{}, maxWorkers)
	results := make(chan alphaPriceResult, len(entries))

	var wg sync.WaitGroup
	for _, entry := range entries {
		wg.Add(1)
		go func(e subnetToken) {
			defer wg.Done()
			sem <- struct{}{}        // acquire slot
			defer func() { <-sem }() // release slot

			token := common.HexToAddress(e.token)

			// Read poolId from LPManager (sequential: GetSlot0 depends on poolId)
			poolId, err := k.lpManager.AlphaTokenToPoolId(nil, token)
			if err != nil || poolId == [32]byte{} {
				return // no pool for this token
			}

			// Read sqrtPriceX96 from PoolManager
			slot0, err := k.poolManager.GetSlot0(nil, poolId)
			if err != nil {
				k.logger.Debug("failed to read pool slot0", "worknetId", e.subnetID, "alphaToken", e.token, "error", err)
				return
			}

			sqrtPrice := new(big.Int).SetBytes(slot0.SqrtPriceX96.Bytes())
			if sqrtPrice.Sign() == 0 {
				return // pool not initialized
			}

			// Compute price: (sqrtPriceX96 / 2^96)^2 = sqrtPriceX96^2 / 2^192
			sqrtPriceBig := slot0.SqrtPriceX96
			q96 := new(big.Int).Lsh(big.NewInt(1), 96)

			// price = sqrtPriceX96^2 / 2^192, scaled to 18 decimals
			numerator := new(big.Int).Mul(sqrtPriceBig, sqrtPriceBig)
			numerator.Mul(numerator, big.NewInt(1e18))
			denominator := new(big.Int).Mul(q96, q96)
			priceRaw := new(big.Int).Div(numerator, denominator)

			priceData, _ := json.Marshal(alphaPrice{
				PriceInAWP:   priceRaw.String(),
				SqrtPriceX96: sqrtPriceBig.String(),
				UpdatedAt:    time.Now().Unix(),
			})

			results <- alphaPriceResult{
				key:  fmt.Sprintf("alpha_price:%s", e.subnetID),
				data: priceData,
			}
		}(entry)
	}

	// Close results channel after all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results and write all to Redis via pipeline
	var collected []alphaPriceResult
	for r := range results {
		collected = append(collected, r)
	}

	if len(collected) > 0 {
		pipe := k.redis.Pipeline()
		for _, r := range collected {
			pipe.Set(ctx, r.key, r.data, 10*time.Minute)
		}
		if _, err := pipe.Exec(ctx); err != nil {
			k.logger.Error("failed to pipeline alpha prices to Redis", "error", err)
		}
	}
}

// compoundAllFees iterates active subnet alpha tokens and compounds LP fees for each.
// Called every 24h by cron. Each compoundFees call collects accrued trading fees from
// the LP position and reinvests them as additional liquidity.
func (k *Keeper) compoundAllFees(ctx context.Context) {
	// No longer holding lock for entire duration — lock per-tx to avoid blocking settleEpoch

	alphaTokensKey := fmt.Sprintf("active_alpha_tokens:%d", k.chainIDInt)
	data, err := k.redis.Get(ctx, alphaTokensKey).Result()
	if err != nil {
		k.logger.Debug("no active alpha tokens cached for compounding (will be populated by indexer)", "chainId", k.chainIDInt)
		return
	}

	var alphaTokens []string
	if err := json.Unmarshal([]byte(data), &alphaTokens); err != nil {
		k.logger.Error("failed to parse active alpha tokens", "error", err)
		return
	}

	if len(alphaTokens) == 0 {
		return
	}

	compounded := 0
	for _, tokenHex := range alphaTokens {
		token := common.HexToAddress(tokenHex)

		// Per-tx: lock -> auth (auto nonce) -> send -> unlock -> wait for confirmation
		k.txMu.Lock()
		auth, err := k.auth()
		if err != nil {
			k.txMu.Unlock()
			k.logger.Error("failed to create tx signer for compoundFees", "error", err)
			return
		}
		timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		auth.Context = timeoutCtx
		tx, txErr := k.lpManager.CompoundFees(auth, token)
		k.txMu.Unlock()

		if txErr != nil {
			cancel()
			k.logger.Debug("compoundFees skipped", "alphaToken", tokenHex, "error", txErr)
			continue
		}
		k.logger.Info("compoundFees tx sent", "alphaToken", tokenHex, "txHash", tx.Hash().Hex())

		// Wait for tx to be mined before sending next one to avoid nonce conflicts
		_, _ = bind.WaitMined(timeoutCtx, k.client, tx)
		cancel()
		compounded++
	}

	if compounded > 0 {
		k.logger.Info("LP fee compounding complete", "compounded", compounded, "total", len(alphaTokens))
	}
}

// updateRelayerBalance reads the relayer native token balance, writes to Redis, and alerts when balance is low
func (k *Keeper) updateRelayerBalance(ctx context.Context) {
	addr := crypto.PubkeyToAddress(k.key.PublicKey)
	balance, err := k.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		k.logger.Debug("failed to read relayer balance", "error", err)
		return
	}
	key := fmt.Sprintf("relayer_balance:%d", k.chainIDInt)
	k.redis.Set(ctx, key, balance.String(), 1*time.Minute)

	// Alert when balance is below 0.01 ETH/BNB
	threshold := big.NewInt(1e16)
	if balance.Cmp(threshold) < 0 {
		k.logger.Error("ALERT: relayer balance critically low",
			"chainId", k.chainIDInt,
			"balance", balance.String(),
			"address", addr.Hex(),
		)
	}
}
