package e2e_test

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cortexia/rootnet/api/internal/chain"
	"github.com/cortexia/rootnet/api/internal/chain/bindings"
	"github.com/cortexia/rootnet/api/internal/config"
	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/cortexia/rootnet/api/internal/ratelimit"
	"github.com/cortexia/rootnet/api/internal/server"
	"github.com/cortexia/rootnet/api/internal/server/handler"
	"github.com/cortexia/rootnet/api/internal/server/ws"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// e2eEnv is the end-to-end test environment: Anvil local chain + Indexer + API
type e2eEnv struct {
	client   *ethclient.Client
	deployer *ecdsa.PrivateKey
	chainID  *big.Int

	// Contract bindings (created from addresses)
	awpRegistry *bindings.AWPRegistry
	awpToken *bindings.AWPToken
	stakeNFT *bindings.StakeNFT

	// Contract addresses
	awpRegistryAddr  common.Address
	awpAddr      common.Address
	svAddr       common.Address
	emAddr       common.Address
	nftAddr      common.Address
	trsAddr      common.Address
	stakeNFTAddr common.Address

	// Backend
	pool    *pgxpool.Pool
	rdb     *redis.Client
	queries *gen.Queries
	indexer *chain.Indexer
	router  http.Handler

	t *testing.T
}

func newE2EEnv(t *testing.T) *e2eEnv {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/awp_test?sslmode=disable"
	}
	redisURL := os.Getenv("TEST_REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/1"
	}

	// Connect to pre-started Anvil + pre-deployed contracts
	// Prerequisites:
	//   anvil --port 18545
	//   forge script script/TestDeploy.s.sol:TestDeploy --rpc-url http://127.0.0.1:18545 --broadcast --private-key ac09...
	client, err := ethclient.Dial("http://127.0.0.1:18545")
	if err != nil {
		t.Skipf("skipping E2E test: cannot connect to Anvil: %v", err)
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		t.Skipf("skipping E2E test: Anvil not ready: %v", err)
	}
	deployerKey, _ := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")

	// Read contract addresses from broadcast file
	addresses := readDeployedAddresses(t, "/home/ubuntu/code/Cortexia/contracts")

	awpRegistryAddr := common.HexToAddress(addresses["AWPRegistry"])
	awpAddr := common.HexToAddress(addresses["AWPToken"])
	svAddr := common.HexToAddress(addresses["StakingVault"])
	emAddr := common.HexToAddress(addresses["AWPEmission"])
	nftAddr := common.HexToAddress(addresses["SubnetNFT"])
	trsAddr := common.HexToAddress(addresses["Treasury"])

	stakeNFTAddr := common.HexToAddress(addresses["StakeNFT"])

	awpRegistry, _ := bindings.NewAWPRegistry(awpRegistryAddr, client)
	awpToken, _ := bindings.NewAWPToken(awpAddr, client)
	stakeNFT, _ := bindings.NewStakeNFT(stakeNFTAddr, client)

	// Connect to DB
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skipf("skipping E2E test: DB unavailable: %v", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Skipf("skipping E2E test: DB connection failed: %v", err)
	}

	// Redis
	opt, _ := redis.ParseURL(redisURL)
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		pool.Close()
		t.Skipf("skipping E2E test: Redis unavailable: %v", err)
	}

	env := &e2eEnv{
		client:      client,
		deployer:    deployerKey,
		chainID:     chainID,
		awpRegistry:  awpRegistry,
		awpToken:     awpToken,
		stakeNFT:     stakeNFT,
		awpRegistryAddr:  awpRegistryAddr,
		awpAddr:      awpAddr,
		svAddr:       svAddr,
		emAddr:       emAddr,
		nftAddr:      nftAddr,
		trsAddr:      trsAddr,
		stakeNFTAddr: stakeNFTAddr,
		pool:        pool,
		rdb:         rdb,
		queries:     gen.New(pool),
		t:           t,
	}

	env.cleanDB()
	env.setupIndexer()
	env.setupRouter()

	t.Cleanup(func() {
		env.cleanDB()
		pool.Close()
		_ = rdb.Close()
		client.Close()
	})

	return env
}

// readDeployedAddresses reads contract addresses from forge broadcast output
func readDeployedAddresses(t *testing.T, contractsDir string) map[string]string {
	t.Helper()
	// forge broadcast output is at broadcast/TestDeploy.s.sol/<chainId>/run-latest.json
	// Anvil default chainId = 31337
	runFile := fmt.Sprintf("%s/broadcast/TestDeploy.s.sol/31337/run-latest.json", contractsDir)
	data, err := os.ReadFile(runFile)
	if err != nil {
		t.Fatalf("failed to read deployment result: %v", err)
	}

	var result struct {
		Transactions []struct {
			ContractName string `json:"contractName"`
			ContractAddress string `json:"contractAddress"`
			TransactionType string `json:"transactionType"`
		} `json:"transactions"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to parse deployment result: %v", err)
	}

	addrs := make(map[string]string)
	for _, tx := range result.Transactions {
		if tx.TransactionType == "CREATE" && tx.ContractName != "" {
			addrs[tx.ContractName] = tx.ContractAddress
		}
	}

	required := []string{"AWPRegistry", "AWPToken", "StakingVault", "AWPEmission", "SubnetNFT", "Treasury", "StakeNFT"}
	for _, name := range required {
		if addrs[name] == "" {
			t.Fatalf("deployment result missing contract: %s (found: %v)", name, addrs)
		}
	}

	return addrs
}

func (e *e2eEnv) cleanDB() {
	ctx := context.Background()
	_, _ = e.pool.Exec(ctx, `TRUNCATE TABLE
		recipient_awp_distributions, stake_positions, stake_allocations,
		user_balances, epochs,
		subnets, proposals, users, sync_states`)
	_ = e.rdb.FlushDB(ctx).Err()
}

func (e *e2eEnv) auth() *bind.TransactOpts {
	auth, _ := bind.NewKeyedTransactorWithChainID(e.deployer, e.chainID)
	return auth
}

func (e *e2eEnv) setupIndexer() {
	addrs := map[string]string{
		"AWPRegistry":  e.awpRegistryAddr.Hex(),
		"AWPToken":     e.awpAddr.Hex(),
		"AWPEmission":  e.emAddr.Hex(),
		"StakingVault": e.svAddr.Hex(),
		"SubnetNFT":    e.nftAddr.Hex(),
		"AWPDAO":       common.Address{}.Hex(),
		"StakeNFT":     e.stakeNFTAddr.Hex(),
	}
	chainClient, err := chain.NewClient(context.Background(), "http://127.0.0.1:18545", addrs)
	if err != nil {
		e.t.Fatalf("failed to create chain.Client: %v", err)
	}
	e.indexer, _ = chain.NewIndexer(chainClient, e.pool, e.rdb, 0)
}

func (e *e2eEnv) setupRouter() {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		AWPRegistryAddress:  e.awpRegistryAddr.Hex(),
		AWPTokenAddress:     e.awpAddr.Hex(),
		StakingVaultAddress: e.svAddr.Hex(),
		AWPEmissionAddress:  e.emAddr.Hex(),
		SubnetNFTAddress:    e.nftAddr.Hex(),
		TreasuryAddress:     e.trsAddr.Hex(),
		StakeNFTAddress:     e.stakeNFTAddr.Hex(),
	}
	limiter := ratelimit.NewLimiter(e.rdb, logger)
	h := handler.NewHandler(e.queries, e.rdb, cfg, logger, limiter)
	hub := ws.NewHub(e.rdb, logger)
	e.router = server.NewRouter(server.RouterParams{Handler: h, Hub: hub})
}

// runIndexer executes one indexer poll
func (e *e2eEnv) runIndexer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = e.indexer.Run(ctx)
}

func (e *e2eEnv) request(method, path, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rr := httptest.NewRecorder()
	e.router.ServeHTTP(rr, req)
	return rr
}

// ════════════════════════════════════════════
//  E2E: Contract → Indexer → DB → API
// ════════════════════════════════════════════

func TestE2E_Chain_RegisterUser(t *testing.T) {
	env := newE2EEnv(t)
	deployer := crypto.PubkeyToAddress(env.deployer.PublicKey)

	// Register on-chain
	_, err := env.awpRegistry.Register(env.auth())
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	time.Sleep(500 * time.Millisecond) // wait for auto-mine

	// Indexer processes
	env.runIndexer()

	// API verification
	addr := strings.ToLower(deployer.Hex())

	t.Run("GetUser", func(t *testing.T) {
		rr := env.request("GET", "/api/users/"+addr, "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("CheckAddress", func(t *testing.T) {
		rr := env.request("GET", "/api/address/"+addr+"/check", "")
		var resp map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["isRegisteredUser"] != true {
			t.Errorf("expected isRegisteredUser=true, got %v", resp)
		}
	})

	t.Run("UserCount", func(t *testing.T) {
		rr := env.request("GET", "/api/users/count", "")
		var resp map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["count"] != float64(1) {
			t.Errorf("expected count=1, got %v", resp["count"])
		}
	})
}

func TestE2E_Chain_DepositAWP(t *testing.T) {
	env := newE2EEnv(t)
	deployer := crypto.PubkeyToAddress(env.deployer.PublicKey)
	addr := strings.ToLower(deployer.Hex())

	// Register
	_, _ = env.awpRegistry.Register(env.auth())
	time.Sleep(300 * time.Millisecond)

	// Approve + Deposit via StakeNFT
	deposit := new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	_, _ = env.awpToken.Approve(env.auth(), env.stakeNFTAddr, deposit)
	time.Sleep(300 * time.Millisecond)
	_, err := env.stakeNFT.Deposit(env.auth(), deposit, 10) // lock for 10 epochs
	if err != nil {
		t.Fatalf("deposit failed: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// Indexer
	env.runIndexer()

	// API balance verification
	rr := env.request("GET", "/api/staking/user/"+addr+"/balance", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var bal map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &bal)
	if bal["totalStaked"] == "0" || bal["totalStaked"] == nil {
		t.Errorf("expected non-zero totalStaked, got %v", bal["totalStaked"])
	}
	t.Logf("API balance after on-chain deposit: %v", bal)
}

func TestE2E_Chain_RegisterAgent(t *testing.T) {
	env := newE2EEnv(t)
	deployer := crypto.PubkeyToAddress(env.deployer.PublicKey)

	// Register user
	_, _ = env.awpRegistry.Register(env.auth())
	time.Sleep(300 * time.Millisecond)

	// Use Anvil's 2nd account as Agent
	agentKey, _ := crypto.HexToECDSA("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	agentAddr := crypto.PubkeyToAddress(agentKey.PublicKey)
	agentAuth, _ := bind.NewKeyedTransactorWithChainID(agentKey, env.chainID)

	_, err := env.awpRegistry.Bind(agentAuth, deployer)
	if err != nil {
		t.Fatalf("bind (registerAgent) failed: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	env.runIndexer()

	t.Run("AgentsByOwner", func(t *testing.T) {
		rr := env.request("GET", "/api/agents/by-owner/"+strings.ToLower(deployer.Hex()), "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
		var agents []map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &agents)
		if len(agents) != 1 {
			t.Fatalf("expected 1 agent, got %d: %s", len(agents), rr.Body.String())
		}
	})

	t.Run("LookupAgent", func(t *testing.T) {
		rr := env.request("GET", "/api/agents/lookup/"+strings.ToLower(agentAddr.Hex()), "")
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("CheckAgent", func(t *testing.T) {
		rr := env.request("GET", "/api/address/"+strings.ToLower(agentAddr.Hex())+"/check", "")
		var resp map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["isRegisteredAgent"] != true {
			t.Errorf("expected isRegisteredAgent=true, got %v", resp)
		}
	})
}
