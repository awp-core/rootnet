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
)

// Relayer submits gasless transactions using a relayer private key (bindFor / setRecipientFor / registerSubnetFor)
// Uses mutex to serialize tx submissions, preventing nonce collisions
type Relayer struct {
	client       *ethclient.Client
	awpRegistry  *bindings.AWPRegistry
	stakingVault *bindings.StakingVault
	key          *ecdsa.PrivateKey
	chainID      *big.Int
	logger       *slog.Logger
	mu           sync.Mutex // serializes tx submissions to prevent nonce collisions
	rdb          *redis.Client // for tx status tracking
}

// RelayTxStatus Redis 中存储的交易状态
type RelayTxStatus struct {
	Status    string `json:"status"` // "pending", "confirmed", "failed"
	TxHash   string `json:"txHash"`
	BlockNum uint64 `json:"blockNumber,omitempty"`
}


// NewRelayer creates a Relayer instance
func NewRelayer(
	client *ethclient.Client,
	awpRegistryAddr common.Address,
	stakingVaultAddr common.Address,
	key *ecdsa.PrivateKey,
	chainID *big.Int,
	rdb *redis.Client,
	logger *slog.Logger,
) (*Relayer, error) {
	awpRegistry, err := bindings.NewAWPRegistry(awpRegistryAddr, client)
	if err != nil {
		return nil, fmt.Errorf("bind AWPRegistry: %w", err)
	}
	stakingVault, err := bindings.NewStakingVault(stakingVaultAddr, client)
	if err != nil {
		return nil, fmt.Errorf("bind StakingVault: %w", err)
	}
	return &Relayer{
		client:       client,
		awpRegistry:  awpRegistry,
		stakingVault: stakingVault,
		key:          key,
		chainID:      chainID,
		rdb:          rdb,
		logger:       logger,
	}, nil
}

// CheckNonce reads the on-chain nonce for an address and returns it.
// Used by relay handlers to pre-check that a signature's nonce is still valid before submitting.
func (r *Relayer) CheckNonce(user common.Address) (*big.Int, error) {
	return r.awpRegistry.Nonces(nil, user)
}

// trackTx 异步追踪交易 receipt，结果写入 Redis
func (r *Relayer) trackTx(txHash string) {
	if r.rdb == nil {
		return
	}
	key := "relay:status:" + txHash
	status := RelayTxStatus{Status: "pending", TxHash: txHash}
	data, _ := json.Marshal(status)
	r.rdb.Set(context.Background(), key, data, 10*time.Minute)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		hash := common.HexToHash(txHash)
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				// 超时，标记失败
				status := RelayTxStatus{Status: "failed", TxHash: txHash}
				data, _ := json.Marshal(status)
				r.rdb.Set(context.Background(), key, data, 10*time.Minute)
				return
			case <-ticker.C:
				receipt, err := r.client.TransactionReceipt(ctx, hash)
				if err != nil {
					continue // 还在 pending
				}
				st := "confirmed"
				if receipt.Status == 0 {
					st = "failed"
				}
				result := RelayTxStatus{Status: st, TxHash: txHash, BlockNum: receipt.BlockNumber.Uint64()}
				data, _ := json.Marshal(result)
				r.rdb.Set(context.Background(), key, data, 10*time.Minute)
				return
			}
		}
	}()
}

// GetTxStatus 查询 relay 交易状态
func (r *Relayer) GetTxStatus(ctx context.Context, txHash string) (*RelayTxStatus, error) {
	if r.rdb == nil {
		return nil, fmt.Errorf("redis not configured")
	}
	key := "relay:status:" + txHash
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var status RelayTxStatus
	if err := json.Unmarshal([]byte(val), &status); err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *Relayer) auth(ctx context.Context) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(r.key, r.chainID)
	if err != nil {
		return nil, fmt.Errorf("create tx signer: %w", err)
	}
	auth.Context = ctx
	return auth, nil
}

// RelayBind relays a bindFor transaction (V2: agent binds to target)
func (r *Relayer) RelayBind(ctx context.Context, agent common.Address, target common.Address, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.awpRegistry.BindFor(auth, agent, target, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("BindFor tx: %w", err)
	}

	txHashHex := tx.Hash().Hex()
	r.logger.Info("relay bindFor sent", "txHash", txHashHex, "agent", agent.Hex(), "target", target.Hex())
	r.trackTx(txHashHex)
	return txHashHex, nil
}

// RelaySetRecipient relays a setRecipientFor transaction (V2: renamed from setRewardRecipientFor)
func (r *Relayer) RelaySetRecipient(ctx context.Context, user common.Address, recipient common.Address, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.awpRegistry.SetRecipientFor(auth, user, recipient, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("SetRecipientFor tx: %w", err)
	}

	r.logger.Info("relay setRecipientFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex(), "recipient", recipient.Hex())
	r.trackTx(tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// RelayAllocate relays an allocateFor transaction (V2: staker instead of user)
func (r *Relayer) RelayAllocate(ctx context.Context, staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.stakingVault.AllocateFor(auth, staker, agent, subnetId, amount, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("AllocateFor tx: %w", err)
	}

	r.logger.Info("relay allocateFor sent", "txHash", tx.Hash().Hex(), "staker", staker.Hex(), "agent", agent.Hex())
	r.trackTx(tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// RelayDeallocate relays a deallocateFor transaction (V2: staker instead of user)
func (r *Relayer) RelayDeallocate(ctx context.Context, staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.stakingVault.DeallocateFor(auth, staker, agent, subnetId, amount, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("DeallocateFor tx: %w", err)
	}

	r.logger.Info("relay deallocateFor sent", "txHash", tx.Hash().Hex(), "staker", staker.Hex(), "agent", agent.Hex())
	r.trackTx(tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// RelayActivateSubnet relays an activateSubnetFor transaction
func (r *Relayer) RelayActivateSubnet(ctx context.Context, user common.Address, subnetId *big.Int, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.awpRegistry.ActivateSubnetFor(auth, user, subnetId, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("ActivateSubnetFor tx: %w", err)
	}

	r.logger.Info("relay activateSubnetFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex(), "subnetId", subnetId.String())
	r.trackTx(tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// RelayRegisterSubnet relays a fully gasless registerSubnetForWithPermit transaction
// User signs two off-chain messages: (1) ERC-2612 permit for AWP, (2) EIP-712 registerSubnet
func (r *Relayer) RelayRegisterSubnet(
	ctx context.Context,
	user common.Address,
	params bindings.IAWPRegistrySubnetParams,
	deadline *big.Int,
	permitV uint8, permitR [32]byte, permitS [32]byte,
	registerV uint8, registerR [32]byte, registerS [32]byte,
) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 60*time.Second) // subnet registration is slower
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.awpRegistry.RegisterSubnetForWithPermit(
		auth, user, params, deadline,
		permitV, permitR, permitS,
		registerV, registerR, registerS,
	)
	if err != nil {
		return "", fmt.Errorf("RegisterSubnetForWithPermit tx: %w", err)
	}

	r.logger.Info("relay registerSubnetFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex(), "name", params.Name)
	r.trackTx(tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}
