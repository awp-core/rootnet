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

// Relayer submits gasless transactions using a relayer private key (bindFor / setRecipientFor / registerWorknetFor)
// Uses mutex to serialize tx submissions, preventing nonce collisions
type Relayer struct {
	client       *ethclient.Client
	awpRegistry  *bindings.AWPRegistry
	awpAllocator *bindings.AWPAllocator
	veAWPHelper  *bindings.VeAWPHelper
	key          *ecdsa.PrivateKey
	chainID      *big.Int
	logger       *slog.Logger
	mu           sync.Mutex // serializes tx submissions to prevent nonce collisions
	rdb          *redis.Client // for tx status tracking
}

// RelayTxStatus is the transaction status stored in Redis
type RelayTxStatus struct {
	Status    string `json:"status"` // "pending", "confirmed", "failed"
	TxHash   string `json:"txHash"`
	BlockNum uint64 `json:"blockNumber,omitempty"`
}


// NewRelayer creates a Relayer instance
// VeAWPHelper address — same on all chains (CREATE2 deterministic deployment)
var VeAWPHelperAddr = common.HexToAddress("0x0000561EDE5C1Ba0b81cE585964050bEAE730001")

func NewRelayer(
	client *ethclient.Client,
	awpRegistryAddr common.Address,
	awpAllocatorAddr common.Address,
	key *ecdsa.PrivateKey,
	chainID *big.Int,
	rdb *redis.Client,
	logger *slog.Logger,
) (*Relayer, error) {
	awpRegistry, err := bindings.NewAWPRegistry(awpRegistryAddr, client)
	if err != nil {
		return nil, fmt.Errorf("bind AWPRegistry: %w", err)
	}
	awpAllocator, err := bindings.NewAWPAllocator(awpAllocatorAddr, client)
	if err != nil {
		return nil, fmt.Errorf("bind AWPAllocator: %w", err)
	}
	veAWPHelper, err := bindings.NewVeAWPHelper(VeAWPHelperAddr, client)
	if err != nil {
		return nil, fmt.Errorf("bind VeAWPHelper: %w", err)
	}
	return &Relayer{
		client:       client,
		awpRegistry:  awpRegistry,
		awpAllocator: awpAllocator,
		veAWPHelper:  veAWPHelper,
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

// trackTx asynchronously tracks transaction receipt and writes the result to Redis
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
				// Timeout, mark as failed
				status := RelayTxStatus{Status: "failed", TxHash: txHash}
				data, _ := json.Marshal(status)
				r.rdb.Set(context.Background(), key, data, 10*time.Minute)
				return
			case <-ticker.C:
				receipt, err := r.client.TransactionReceipt(ctx, hash)
				if err != nil {
					continue // still pending
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

// GetTxStatus queries the relay transaction status
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
func (r *Relayer) RelayAllocate(ctx context.Context, staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.awpAllocator.AllocateFor(auth, staker, agent, worknetId, amount, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("AllocateFor tx: %w", err)
	}

	r.logger.Info("relay allocateFor sent", "txHash", tx.Hash().Hex(), "staker", staker.Hex(), "agent", agent.Hex())
	r.trackTx(tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// RelayDeallocate relays a deallocateFor transaction (V2: staker instead of user)
func (r *Relayer) RelayDeallocate(ctx context.Context, staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.awpAllocator.DeallocateFor(auth, staker, agent, worknetId, amount, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("DeallocateFor tx: %w", err)
	}

	r.logger.Info("relay deallocateFor sent", "txHash", tx.Hash().Hex(), "staker", staker.Hex(), "agent", agent.Hex())
	r.trackTx(tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// NOTE: activateWorknet is Guardian-only (onlyGuardian modifier) — not relayable.
// There is no gasless activateWorknetFor variant in the contract.

// RelayRegisterSubnet relays a fully gasless registerSubnetForWithPermit transaction
// User signs two off-chain messages: (1) ERC-2612 permit for AWP, (2) EIP-712 registerSubnet
func (r *Relayer) RelayRegisterSubnet(
	ctx context.Context,
	user common.Address,
	params bindings.IAWPRegistryWorknetParams,
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

	tx, err := r.awpRegistry.RegisterWorknetForWithPermit(
		auth, user, params, deadline,
		permitV, permitR, permitS,
		registerV, registerR, registerS,
	)
	if err != nil {
		return "", fmt.Errorf("RegisterWorknetForWithPermit tx: %w", err)
	}

	r.logger.Info("relay registerSubnetFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex(), "name", params.Name)
	r.trackTx(tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// RelayGrantDelegate relays a grantDelegateFor transaction
func (r *Relayer) RelayGrantDelegate(ctx context.Context, user common.Address, delegate common.Address, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.awpRegistry.GrantDelegateFor(auth, user, delegate, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("GrantDelegateFor tx: %w", err)
	}

	r.logger.Info("relay grantDelegateFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex(), "delegate", delegate.Hex())
	r.trackTx(tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// RelayRevokeDelegate relays a revokeDelegateFor transaction
func (r *Relayer) RelayRevokeDelegate(ctx context.Context, user common.Address, delegate common.Address, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.awpRegistry.RevokeDelegateFor(auth, user, delegate, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("RevokeDelegateFor tx: %w", err)
	}

	r.logger.Info("relay revokeDelegateFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex(), "delegate", delegate.Hex())
	r.trackTx(tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// RelayUnbind relays an unbindFor transaction
func (r *Relayer) RelayUnbind(ctx context.Context, user common.Address, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.awpRegistry.UnbindFor(auth, user, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("UnbindFor tx: %w", err)
	}

	r.logger.Info("relay unbindFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex())
	r.trackTx(tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// RelayStake relays a gasless stake via VeAWPHelper.depositFor
// User signs one ERC-2612 permit (AWP → VeAWPHelper), relayer pays gas.
func (r *Relayer) RelayStake(
	ctx context.Context,
	user common.Address,
	amount *big.Int,
	lockDuration uint64,
	deadline *big.Int,
	v uint8, rs [32]byte, ss [32]byte,
) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 60*time.Second) // staking is heavier
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.veAWPHelper.DepositFor(auth, user, amount, lockDuration, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("VeAWPHelper.DepositFor tx: %w", err)
	}

	r.logger.Info("relay stake sent", "txHash", tx.Hash().Hex(), "user", user.Hex(), "amount", amount.String())
	r.trackTx(tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}
