package chain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/cortexia/rootnet/api/internal/chain/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Relayer submits gasless transactions using a relayer private key (registerFor / bindFor / registerSubnetFor)
// Uses mutex to serialize tx submissions, preventing nonce collisions
type Relayer struct {
	client  *ethclient.Client
	rootNet *bindings.RootNet
	key     *ecdsa.PrivateKey
	chainID *big.Int
	logger  *slog.Logger
	mu      sync.Mutex // serializes tx submissions to prevent nonce collisions
}

// NewRelayer creates a Relayer instance
func NewRelayer(
	client *ethclient.Client,
	rootNetAddr common.Address,
	key *ecdsa.PrivateKey,
	chainID *big.Int,
	logger *slog.Logger,
) (*Relayer, error) {
	rootNet, err := bindings.NewRootNet(rootNetAddr, client)
	if err != nil {
		return nil, fmt.Errorf("bind RootNet: %w", err)
	}
	return &Relayer{
		client:  client,
		rootNet: rootNet,
		key:     key,
		chainID: chainID,
		logger:  logger,
	}, nil
}

func (r *Relayer) auth(ctx context.Context) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(r.key, r.chainID)
	if err != nil {
		return nil, fmt.Errorf("create tx signer: %w", err)
	}
	auth.Context = ctx
	return auth, nil
}

// RelayRegister relays a registerFor transaction
func (r *Relayer) RelayRegister(ctx context.Context, user common.Address, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.rootNet.RegisterFor(auth, user, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("RegisterFor tx: %w", err)
	}

	r.logger.Info("relay registerFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex())
	return tx.Hash().Hex(), nil
}

// RelayBind relays a bindFor transaction
func (r *Relayer) RelayBind(ctx context.Context, agent common.Address, principal common.Address, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.rootNet.BindFor(auth, agent, principal, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("BindFor tx: %w", err)
	}

	r.logger.Info("relay bindFor sent", "txHash", tx.Hash().Hex(), "agent", agent.Hex(), "principal", principal.Hex())
	return tx.Hash().Hex(), nil
}

// RelaySetRewardRecipient relays a setRewardRecipientFor transaction
func (r *Relayer) RelaySetRewardRecipient(ctx context.Context, user common.Address, recipient common.Address, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.rootNet.SetRewardRecipientFor(auth, user, recipient, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("SetRewardRecipientFor tx: %w", err)
	}

	r.logger.Info("relay setRewardRecipientFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex(), "recipient", recipient.Hex())
	return tx.Hash().Hex(), nil
}

// RelayAllocate relays an allocateFor transaction
func (r *Relayer) RelayAllocate(ctx context.Context, user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.rootNet.AllocateFor(auth, user, agent, subnetId, amount, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("AllocateFor tx: %w", err)
	}

	r.logger.Info("relay allocateFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex(), "agent", agent.Hex())
	return tx.Hash().Hex(), nil
}

// RelayDeallocate relays a deallocateFor transaction
func (r *Relayer) RelayDeallocate(ctx context.Context, user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, rs [32]byte, ss [32]byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	auth, err := r.auth(txCtx)
	if err != nil {
		return "", err
	}

	tx, err := r.rootNet.DeallocateFor(auth, user, agent, subnetId, amount, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("DeallocateFor tx: %w", err)
	}

	r.logger.Info("relay deallocateFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex(), "agent", agent.Hex())
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

	tx, err := r.rootNet.ActivateSubnetFor(auth, user, subnetId, deadline, v, rs, ss)
	if err != nil {
		return "", fmt.Errorf("ActivateSubnetFor tx: %w", err)
	}

	r.logger.Info("relay activateSubnetFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex(), "subnetId", subnetId.String())
	return tx.Hash().Hex(), nil
}

// RelayRegisterSubnet relays a fully gasless registerSubnetForWithPermit transaction
// User signs two off-chain messages: (1) ERC-2612 permit for AWP, (2) EIP-712 registerSubnet
func (r *Relayer) RelayRegisterSubnet(
	ctx context.Context,
	user common.Address,
	params bindings.IRootNetSubnetParams,
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

	tx, err := r.rootNet.RegisterSubnetForWithPermit(
		auth, user, params, deadline,
		permitV, permitR, permitS,
		registerV, registerR, registerS,
	)
	if err != nil {
		return "", fmt.Errorf("RegisterSubnetForWithPermit tx: %w", err)
	}

	r.logger.Info("relay registerSubnetFor sent", "txHash", tx.Hash().Hex(), "user", user.Hex(), "name", params.Name)
	return tx.Hash().Hex(), nil
}
