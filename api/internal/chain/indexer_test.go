package chain_test

import (
	"context"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// testDB holds the test database environment
type testDB struct {
	pool *pgxpool.Pool
	q    *gen.Queries
	t    *testing.T
}

// setupTestDB connects to the test database and truncates all tables
func setupTestDB(t *testing.T) *testDB {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/awp_test?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skipf("skipping test: cannot connect to database: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Skipf("skipping test: database unavailable: %v", err)
	}

	db := &testDB{
		pool: pool,
		q:    gen.New(pool),
		t:    t,
	}
	db.truncate()

	t.Cleanup(func() {
		db.truncate()
		pool.Close()
	})

	return db
}

// truncate clears all tables
func (db *testDB) truncate() {
	db.t.Helper()
	_, err := db.pool.Exec(context.Background(), `TRUNCATE TABLE
		recipient_awp_distributions, stake_positions, stake_allocations,
		user_balances, epochs,
		subnets, proposals, users, sync_states`)
	if err != nil {
		db.t.Fatalf("TRUNCATE failed: %v", err)
	}
}

// numericFromInt64 converts an int64 to pgtype.Numeric
const testChainID int64 = 8453

func numericFromInt64(v int64) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   big.NewInt(v),
		Exp:   0,
		Valid: true,
	}
}

// assertNumericEqual asserts that a pgtype.Numeric equals the expected int64 value
func assertNumericEqual(t *testing.T, label string, got pgtype.Numeric, want int64) {
	t.Helper()
	if !got.Valid {
		t.Fatalf("%s: numeric value invalid (Valid=false)", label)
	}
	// Convert Numeric to big.Int for comparison (valid for integer cases with Exp=0)
	wantBig := big.NewInt(want)
	gotBig := got.Int
	if got.Exp != 0 {
		// Handle possible exponent offset
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(got.Exp)), nil)
		if got.Exp > 0 {
			gotBig = new(big.Int).Mul(gotBig, scale)
		}
		// Negative exponent means decimal fraction; tests only use integers so this should not occur
	}
	if gotBig.Cmp(wantBig) != 0 {
		t.Fatalf("%s: expected %d, got %s", label, want, gotBig.String())
	}
}

// --- Test scenarios: simulate database operations for indexer event processing ---

func TestIndexerScenario_Bound(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	agentAddr := "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	targetAddr := "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

	// Simulate Bound event: agent binds to target
	err := db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID,
		Address: agentAddr,
		BoundTo: targetAddr,
	})
	if err != nil {
		t.Fatalf("UpsertUserBinding failed: %v", err)
	}

	err = db.q.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: testChainID, UserAddress: agentAddr})
	if err != nil {
		t.Fatalf("InitUserBalance failed: %v", err)
	}

	// Verify user was created with correct binding
	user, err := db.q.GetUser(ctx, gen.GetUserParams{ChainID: testChainID, Address: agentAddr})
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if strings.TrimSpace(user.Address) != agentAddr {
		t.Fatalf("user address mismatch: expected %s, got %s", agentAddr, user.Address)
	}
	if strings.TrimSpace(user.BoundTo) != targetAddr {
		t.Fatalf("bound_to mismatch: expected %s, got %s", targetAddr, user.BoundTo)
	}

	// Verify balance initialized to 0
	bal, err := db.q.GetUserBalance(ctx, gen.GetUserBalanceParams{ChainID: testChainID, UserAddress: agentAddr})
	if err != nil {
		t.Fatalf("GetUserBalance failed: %v", err)
	}
	assertNumericEqual(t, "total_allocated", bal.TotalAllocated, 0)
}

func TestIndexerScenario_BindAndLookup(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	ownerAddr := "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	agentAddr := "0xcccccccccccccccccccccccccccccccccccccccc"

	// Simulate Bound event: agent binds to owner
	err := db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID,
		Address: agentAddr,
		BoundTo: ownerAddr,
	})
	if err != nil {
		t.Fatalf("UpsertUserBinding failed: %v", err)
	}

	// Verify agent binding
	user, err := db.q.GetUser(ctx, gen.GetUserParams{ChainID: testChainID, Address: agentAddr})
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if strings.TrimSpace(user.Address) != agentAddr {
		t.Fatalf("address mismatch: expected %s, got %s", agentAddr, user.Address)
	}
	if strings.TrimSpace(user.BoundTo) != ownerAddr {
		t.Fatalf("bound_to mismatch: expected %s, got %s", ownerAddr, user.BoundTo)
	}

	// Verify query by bound_to (owner)
	agents, err := db.q.GetUsersByBoundTo(ctx, gen.GetUsersByBoundToParams{ChainID: testChainID, BoundTo: ownerAddr})
	if err != nil {
		t.Fatalf("GetUsersByBoundTo failed: %v", err)
	}
	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}
	if agents[0].Address != agentAddr {
		t.Fatalf("agent address mismatch: expected %s, got %s", agentAddr, agents[0].Address)
	}
}

func TestIndexerScenario_VeAWPDeposited(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	userAddr := "0xdddddddddddddddddddddddddddddddddddddd"

	// Register user
	_ = db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: userAddr, BoundTo: ""})

	// Simulate VeAWP.Deposited event: create stake position
	err := db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID:      1,
		Owner:        userAddr,
		Amount:       numericFromInt64(1000),
		LockEndTime: 10,
		CreatedAt: 100,
	})
	if err != nil {
		t.Fatalf("InsertStakePosition failed: %v", err)
	}

	pos, err := db.q.GetStakePosition(ctx, gen.GetStakePositionParams{ChainID: testChainID, TokenID: 1})
	if err != nil {
		t.Fatalf("GetStakePosition failed: %v", err)
	}
	assertNumericEqual(t, "position amount", pos.Amount, 1000)
	if strings.TrimSpace(pos.Owner) != userAddr {
		t.Fatalf("owner mismatch: expected %s, got %s", userAddr, pos.Owner)
	}

	// Create another position
	err = db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID:      2,
		Owner:        userAddr,
		Amount:       numericFromInt64(500),
		LockEndTime: 20,
		CreatedAt: 200,
	})
	if err != nil {
		t.Fatalf("second InsertStakePosition failed: %v", err)
	}

	// Verify total staked
	total, err := db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: userAddr})
	if err != nil {
		t.Fatalf("GetUserTotalStaked failed: %v", err)
	}
	assertNumericEqual(t, "cumulative total staked", total, 1500)
}

func TestIndexerScenario_VeAWPPositionIncreased(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	userAddr := "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"

	// Setup: create a stake position
	_ = db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 10, Owner: userAddr, Amount: numericFromInt64(2000),
		LockEndTime: 10, CreatedAt: 100,
	})

	// Simulate PositionIncreased: update amount and lock
	err := db.q.UpdateStakePosition(ctx, gen.UpdateStakePositionParams{ChainID: testChainID,
		TokenID:      10,
		Amount:       numericFromInt64(3000),
		LockEndTime: 20,
	})
	if err != nil {
		t.Fatalf("UpdateStakePosition failed: %v", err)
	}

	pos, err := db.q.GetStakePosition(ctx, gen.GetStakePositionParams{ChainID: testChainID, TokenID: 10})
	if err != nil {
		t.Fatalf("GetStakePosition failed: %v", err)
	}
	assertNumericEqual(t, "updated amount", pos.Amount, 3000)
	if pos.LockEndTime != 20 {
		t.Fatalf("lockEndTime mismatch: expected 20, got %d", pos.LockEndTime)
	}
}

func TestIndexerScenario_VeAWPWithdrawn(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	userAddr := "0xffffffffffffffffffffffffffffffffffffffff"

	// Setup: create a stake position
	_ = db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 20, Owner: userAddr, Amount: numericFromInt64(5000),
		LockEndTime: 5, CreatedAt: 100,
	})

	// Simulate Withdrawn event (burn)
	err := db.q.BurnStakePosition(ctx, gen.BurnStakePositionParams{ChainID: testChainID, TokenID: 20})
	if err != nil {
		t.Fatalf("BurnStakePosition failed: %v", err)
	}

	pos, err := db.q.GetStakePosition(ctx, gen.GetStakePositionParams{ChainID: testChainID, TokenID: 20})
	if err != nil {
		t.Fatalf("GetStakePosition failed: %v", err)
	}
	if !pos.Burned {
		t.Fatal("position should be marked as burned")
	}

	// Burned positions should not appear in user positions list
	positions, err := db.q.GetUserStakePositions(ctx, gen.GetUserStakePositionsParams{ChainID: testChainID, Owner: userAddr})
	if err != nil {
		t.Fatalf("GetUserStakePositions failed: %v", err)
	}
	if len(positions) != 0 {
		t.Fatalf("expected 0 active positions, got %d", len(positions))
	}

	// Total staked should be 0 (burned excluded)
	total, err := db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: userAddr})
	if err != nil {
		t.Fatalf("GetUserTotalStaked failed: %v", err)
	}
	assertNumericEqual(t, "total staked after burn", total, 0)
}

func TestIndexerScenario_VeAWPTransfer(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	aliceAddr := "0x1111111111111111111111111111111111111111"
	bobAddr := "0x2222222222222222222222222222222222222222"

	// Setup: create a position owned by alice
	_ = db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 30, Owner: aliceAddr, Amount: numericFromInt64(1000),
		LockEndTime: 10, CreatedAt: 100,
	})

	// Simulate ERC721 Transfer: alice -> bob
	err := db.q.UpdateStakePositionOwner(ctx, gen.UpdateStakePositionOwnerParams{ChainID: testChainID,
		TokenID: 30,
		Owner:   bobAddr,
	})
	if err != nil {
		t.Fatalf("UpdateStakePositionOwner failed: %v", err)
	}

	pos, err := db.q.GetStakePosition(ctx, gen.GetStakePositionParams{ChainID: testChainID, TokenID: 30})
	if err != nil {
		t.Fatalf("GetStakePosition failed: %v", err)
	}
	if strings.TrimSpace(pos.Owner) != bobAddr {
		t.Fatalf("owner mismatch: expected %s, got %s", bobAddr, pos.Owner)
	}

	// Alice should have 0, bob should have 1000
	aliceTotal, _ := db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: aliceAddr})
	assertNumericEqual(t, "alice total after transfer", aliceTotal, 0)

	bobTotal, _ := db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: bobAddr})
	assertNumericEqual(t, "bob total after transfer", bobTotal, 1000)
}

func TestIndexerScenario_Allocated(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	userAddr := "0x2222222222222222222222222222222222222222"
	agentAddr := "0x3333333333333333333333333333333333333333"
	subnetID := numericFromInt64(1)

	// Setup: register user + balance
	_ = db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: userAddr, BoundTo: ""})
	_ = db.q.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: testChainID, UserAddress: userAddr})

	// Simulate Allocated event
	err := db.q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{ChainID: testChainID,
		UserAddress:  userAddr,
		AgentAddress: agentAddr,
		SubnetID:     subnetID,
		Amount:       numericFromInt64(1000),
	})
	if err != nil {
		t.Fatalf("UpsertStakeAllocation failed: %v", err)
	}

	err = db.q.AddUserAllocated(ctx, gen.AddUserAllocatedParams{ChainID: testChainID,
		UserAddress:    userAddr,
		TotalAllocated: numericFromInt64(1000),
	})
	if err != nil {
		t.Fatalf("AddUserAllocated failed: %v", err)
	}

	// Verify agent-subnet stake
	stake, err := db.q.GetAgentSubnetStake(ctx, gen.GetAgentSubnetStakeParams{ChainID: testChainID,
		AgentAddress: agentAddr,
		SubnetID:     subnetID,
	})
	if err != nil {
		t.Fatalf("GetAgentSubnetStake failed: %v", err)
	}
	assertNumericEqual(t, "agent_subnet_stake", stake, 1000)

	// Verify user allocated has increased
	bal, err := db.q.GetUserBalance(ctx, gen.GetUserBalanceParams{ChainID: testChainID, UserAddress: userAddr})
	if err != nil {
		t.Fatalf("GetUserBalance failed: %v", err)
	}
	assertNumericEqual(t, "total_allocated", bal.TotalAllocated, 1000)

	// Allocate another 500 (UpsertStakeAllocation should accumulate)
	err = db.q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{ChainID: testChainID,
		UserAddress:  userAddr,
		AgentAddress: agentAddr,
		SubnetID:     subnetID,
		Amount:       numericFromInt64(500),
	})
	if err != nil {
		t.Fatalf("second UpsertStakeAllocation failed: %v", err)
	}
	_ = db.q.AddUserAllocated(ctx, gen.AddUserAllocatedParams{ChainID: testChainID,
		UserAddress:    userAddr,
		TotalAllocated: numericFromInt64(500),
	})

	stake, err = db.q.GetAgentSubnetStake(ctx, gen.GetAgentSubnetStakeParams{ChainID: testChainID,
		AgentAddress: agentAddr,
		SubnetID:     subnetID,
	})
	if err != nil {
		t.Fatalf("GetAgentSubnetStake (after accumulation) failed: %v", err)
	}
	assertNumericEqual(t, "agent_subnet_stake after accumulation", stake, 1500)
}

func TestIndexerScenario_Deallocated(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	userAddr := "0x4444444444444444444444444444444444444444"
	agentAddr := "0x5555555555555555555555555555555555555555"
	subnetID := numericFromInt64(2)

	// Setup: register user + allocation 2000
	_ = db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: userAddr, BoundTo: ""})
	_ = db.q.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: testChainID, UserAddress: userAddr})
	_ = db.q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{ChainID: testChainID,
		UserAddress:  userAddr,
		AgentAddress: agentAddr,
		SubnetID:     subnetID,
		Amount:       numericFromInt64(2000),
	})
	_ = db.q.AddUserAllocated(ctx, gen.AddUserAllocatedParams{ChainID: testChainID,
		UserAddress:    userAddr,
		TotalAllocated: numericFromInt64(2000),
	})

	// Simulate Deallocated event: reduce by 700
	err := db.q.SubtractStakeAllocation(ctx, gen.SubtractStakeAllocationParams{ChainID: testChainID,
		UserAddress:  userAddr,
		AgentAddress: agentAddr,
		SubnetID:     subnetID,
		Amount:       numericFromInt64(700),
	})
	if err != nil {
		t.Fatalf("SubtractStakeAllocation failed: %v", err)
	}

	err = db.q.SubtractUserAllocated(ctx, gen.SubtractUserAllocatedParams{ChainID: testChainID,
		UserAddress:    userAddr,
		TotalAllocated: numericFromInt64(700),
	})
	if err != nil {
		t.Fatalf("SubtractUserAllocated failed: %v", err)
	}

	// Verify stake has decreased
	stake, err := db.q.GetAgentSubnetStake(ctx, gen.GetAgentSubnetStakeParams{ChainID: testChainID,
		AgentAddress: agentAddr,
		SubnetID:     subnetID,
	})
	if err != nil {
		t.Fatalf("GetAgentSubnetStake failed: %v", err)
	}
	assertNumericEqual(t, "stake after reduction", stake, 1300)

	// Verify allocated has decreased
	bal, err := db.q.GetUserBalance(ctx, gen.GetUserBalanceParams{ChainID: testChainID, UserAddress: userAddr})
	if err != nil {
		t.Fatalf("GetUserBalance failed: %v", err)
	}
	assertNumericEqual(t, "allocated after reduction", bal.TotalAllocated, 1300)
}

func TestIndexerScenario_SubnetRegistered(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Simulate SubnetRegistered event
	err := db.q.InsertSubnet(ctx, gen.InsertSubnetParams{ChainID: testChainID,
		SubnetID: numericFromInt64(1),
		Owner:          "0xaaaa000000000000000000000000000000000001",
		Name:           "TestSubnet",
		Symbol:         "TST",
		SubnetContract: "0xbbbb000000000000000000000000000000000001",
		MinStake:       pgtype.Numeric{Int: big.NewInt(0), Exp: 0, Valid: true},
		WorknetToken:     "0xcccc000000000000000000000000000000000001",
		LpPool:         pgtype.Text{Valid: false},
		CreatedAt:      200,
		ImmunityEndsAt: pgtype.Int8{Valid: false},
	})
	if err != nil {
		t.Fatalf("InsertSubnet failed: %v", err)
	}

	// Verify
	subnet, err := db.q.GetSubnet(ctx, numericFromInt64(1))
	if err != nil {
		t.Fatalf("GetSubnet failed: %v", err)
	}
	if subnet.Owner != "0xaaaa000000000000000000000000000000000001" {
		t.Fatalf("Owner mismatch: %s", subnet.Owner)
	}
	if subnet.Name != "TestSubnet" {
		t.Fatalf("Name mismatch: %s", subnet.Name)
	}
	if subnet.Symbol != "TST" {
		t.Fatalf("Symbol mismatch: %s", subnet.Symbol)
	}
	if subnet.Status != "Pending" {
		t.Fatalf("newly registered subnet should have Pending status, got: %s", subnet.Status)
	}
	if subnet.LpPool.Valid {
		t.Fatal("newly registered subnet LpPool should be NULL")
	}
	if subnet.Burned {
		t.Fatal("newly registered subnet should not be marked as burned")
	}
}

func TestIndexerScenario_SubnetLifecycle(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Register subnet
	_ = db.q.InsertSubnet(ctx, gen.InsertSubnetParams{ChainID: testChainID,
		SubnetID: numericFromInt64(10),
		Owner:          "0xaaaa000000000000000000000000000000000010",
		Name:           "LifecycleNet",
		Symbol:         "LCN",
		SubnetContract: "0xbbbb000000000000000000000000000000000010",
		MinStake:       pgtype.Numeric{Int: big.NewInt(0), Exp: 0, Valid: true},
		WorknetToken:     "0xcccc000000000000000000000000000000000010",
		LpPool:         pgtype.Text{Valid: false},
		CreatedAt:      300,
		ImmunityEndsAt: pgtype.Int8{Valid: false},
	})

	// 1. Activate subnet (SubnetActivated)
	err := db.q.UpdateSubnetActivated(ctx, gen.UpdateSubnetActivatedParams{
		SubnetID: numericFromInt64(10),
		ActivatedAt: pgtype.Int8{Int64: 350, Valid: true},
	})
	if err != nil {
		t.Fatalf("UpdateSubnetActivated failed: %v", err)
	}

	subnet, err := db.q.GetSubnet(ctx, numericFromInt64(10))
	if err != nil {
		t.Fatalf("GetSubnet failed: %v", err)
	}
	if subnet.Status != "Active" {
		t.Fatalf("status after activation should be Active, got: %s", subnet.Status)
	}
	if !subnet.ActivatedAt.Valid || subnet.ActivatedAt.Int64 != 350 {
		t.Fatalf("ActivatedAt mismatch: %v", subnet.ActivatedAt)
	}

	// 2. Pause subnet (SubnetPaused)
	err = db.q.UpdateSubnetStatus(ctx, gen.UpdateSubnetStatusParams{
		SubnetID: numericFromInt64(10),
		Status:   "Paused",
	})
	if err != nil {
		t.Fatalf("UpdateSubnetStatus(Paused) failed: %v", err)
	}

	subnet, err = db.q.GetSubnet(ctx, numericFromInt64(10))
	if err != nil {
		t.Fatalf("GetSubnet failed: %v", err)
	}
	if subnet.Status != "Paused" {
		t.Fatalf("status after pause should be Paused, got: %s", subnet.Status)
	}

	// 3. Resume subnet (SubnetResumed)
	err = db.q.UpdateSubnetStatus(ctx, gen.UpdateSubnetStatusParams{
		SubnetID: numericFromInt64(10),
		Status:   "Active",
	})
	if err != nil {
		t.Fatalf("UpdateSubnetStatus(Active) failed: %v", err)
	}

	subnet, err = db.q.GetSubnet(ctx, numericFromInt64(10))
	if err != nil {
		t.Fatalf("GetSubnet failed: %v", err)
	}
	if subnet.Status != "Active" {
		t.Fatalf("status after resume should be Active, got: %s", subnet.Status)
	}

	// 4. Ban subnet (SubnetBanned)
	err = db.q.UpdateSubnetStatus(ctx, gen.UpdateSubnetStatusParams{
		SubnetID: numericFromInt64(10),
		Status:   "Banned",
	})
	if err != nil {
		t.Fatalf("UpdateSubnetStatus(Banned) failed: %v", err)
	}

	subnet, err = db.q.GetSubnet(ctx, numericFromInt64(10))
	if err != nil {
		t.Fatalf("GetSubnet failed: %v", err)
	}
	if subnet.Status != "Banned" {
		t.Fatalf("status after ban should be Banned, got: %s", subnet.Status)
	}

	// 5. Unban subnet (SubnetUnbanned)
	err = db.q.UpdateSubnetStatus(ctx, gen.UpdateSubnetStatusParams{
		SubnetID: numericFromInt64(10),
		Status:   "Active",
	})
	if err != nil {
		t.Fatalf("UpdateSubnetStatus(Active) failed: %v", err)
	}

	subnet, err = db.q.GetSubnet(ctx, numericFromInt64(10))
	if err != nil {
		t.Fatalf("GetSubnet failed: %v", err)
	}
	if subnet.Status != "Active" {
		t.Fatalf("status after unban should be Active, got: %s", subnet.Status)
	}

	// 6. Deregister subnet (SubnetDeregistered)
	err = db.q.UpdateSubnetBurned(ctx, numericFromInt64(10))
	if err != nil {
		t.Fatalf("UpdateSubnetBurned failed: %v", err)
	}

	subnet, err = db.q.GetSubnet(ctx, numericFromInt64(10))
	if err != nil {
		t.Fatalf("GetSubnet failed: %v", err)
	}
	if !subnet.Burned {
		t.Fatal("should be marked as burned after deregistration")
	}
}

func TestIndexerScenario_SubnetLP(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Register subnet
	_ = db.q.InsertSubnet(ctx, gen.InsertSubnetParams{ChainID: testChainID,
		SubnetID: numericFromInt64(20),
		Owner:          "0xaaaa000000000000000000000000000000000020",
		Name:           "LPTestNet",
		Symbol:         "LPT",
		SubnetContract: "0xbbbb000000000000000000000000000000000020",
		MinStake:       pgtype.Numeric{Int: big.NewInt(0), Exp: 0, Valid: true},
		WorknetToken:     "0xcccc000000000000000000000000000000000020",
		LpPool:         pgtype.Text{Valid: false},
		CreatedAt:      400,
		ImmunityEndsAt: pgtype.Int8{Valid: false},
	})

	// Simulate LPCreated event
	lpPoolAddr := "0xdddd000000000000000000000000000000000020"
	err := db.q.UpdateSubnetLP(ctx, gen.UpdateSubnetLPParams{
		SubnetID: numericFromInt64(20),
		LpPool:   pgtype.Text{String: lpPoolAddr, Valid: true},
	})
	if err != nil {
		t.Fatalf("UpdateSubnetLP failed: %v", err)
	}

	subnet, err := db.q.GetSubnet(ctx, numericFromInt64(20))
	if err != nil {
		t.Fatalf("GetSubnet failed: %v", err)
	}
	if !subnet.LpPool.Valid || strings.TrimSpace(subnet.LpPool.String) != lpPoolAddr {
		t.Fatalf("LpPool mismatch: %v", subnet.LpPool)
	}
}

func TestIndexerScenario_EpochSettled(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Simulate EpochSettled event
	err := db.q.UpsertEpoch(ctx, gen.UpsertEpochParams{ChainID: testChainID,
		EpochID:       1,
		StartTime:     1000,
		DailyEmission: numericFromInt64(50000),
		DaoEmission:   pgtype.Numeric{Valid: false},
	})
	if err != nil {
		t.Fatalf("UpsertEpoch failed: %v", err)
	}

	epoch, err := db.q.GetEpoch(ctx, gen.GetEpochParams{ChainID: testChainID, EpochID: 1})
	if err != nil {
		t.Fatalf("GetEpoch failed: %v", err)
	}
	if epoch.EpochID != 1 {
		t.Fatalf("EpochID mismatch: %d", epoch.EpochID)
	}
	if epoch.StartTime != 1000 {
		t.Fatalf("StartTime mismatch: %d", epoch.StartTime)
	}
	assertNumericEqual(t, "daily_emission", epoch.DailyEmission, 50000)

	// Simulate DAOMatchDistributed event: update dao_emission
	err = db.q.UpdateEpochDAO(ctx, gen.UpdateEpochDAOParams{ChainID: testChainID,
		EpochID:     1,
		DaoEmission: numericFromInt64(5000),
	})
	if err != nil {
		t.Fatalf("UpdateEpochDAO failed: %v", err)
	}

	epoch, err = db.q.GetEpoch(ctx, gen.GetEpochParams{ChainID: testChainID, EpochID: 1})
	if err != nil {
		t.Fatalf("GetEpoch failed: %v", err)
	}
	assertNumericEqual(t, "dao_emission", epoch.DaoEmission, 5000)
}

func TestIndexerScenario_SubnetAWPDistributed(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create epoch first
	_ = db.q.UpsertEpoch(ctx, gen.UpsertEpochParams{ChainID: testChainID,
		EpochID:       5,
		StartTime:     5000,
		DailyEmission: numericFromInt64(100000),
		DaoEmission:   pgtype.Numeric{Valid: false},
	})

	// Simulate AWPDistributed event (by recipient address)
	recipient1 := "0x0000000000000000000000000000000000000001"
	recipient2 := "0x0000000000000000000000000000000000000002"
	err := db.q.InsertRecipientAWPDistribution(ctx, gen.InsertRecipientAWPDistributionParams{ChainID: testChainID,
		EpochID:   5,
		Recipient: recipient1,
		AwpAmount: numericFromInt64(25000),
	})
	if err != nil {
		t.Fatalf("InsertRecipientAWPDistribution failed: %v", err)
	}

	// Add another recipient distribution
	err = db.q.InsertRecipientAWPDistribution(ctx, gen.InsertRecipientAWPDistributionParams{ChainID: testChainID,
		EpochID:   5,
		Recipient: recipient2,
		AwpAmount: numericFromInt64(30000),
	})
	if err != nil {
		t.Fatalf("InsertRecipientAWPDistribution(recipient2) failed: %v", err)
	}

	// Verify via GetRecipientEarnings query
	earnings, err := db.q.GetRecipientEarnings(ctx, gen.GetRecipientEarningsParams{ChainID: testChainID,
		Recipient: recipient1,
		Limit:     10,
		Offset:    0,
	})
	if err != nil {
		t.Fatalf("GetRecipientEarnings failed: %v", err)
	}
	if len(earnings) != 1 {
		t.Fatalf("expected 1 earnings record, got %d", len(earnings))
	}
	assertNumericEqual(t, "recipient1_awp", earnings[0].AwpAmount, 25000)

	// Verify epoch distribution list
	dists, err := db.q.GetEpochDistributions(ctx, gen.GetEpochDistributionsParams{ChainID: testChainID, EpochID: 5})
	if err != nil {
		t.Fatalf("GetEpochDistributions failed: %v", err)
	}
	if len(dists) != 2 {
		t.Fatalf("expected 2 distribution records, got %d", len(dists))
	}
}

func TestIndexerScenario_SyncState(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Write sync state for the first time
	err := db.q.UpsertSyncState(ctx, gen.UpsertSyncStateParams{ChainID: testChainID,
		ContractName: "indexer",
		LastBlock:    1000,
	})
	if err != nil {
		t.Fatalf("UpsertSyncState failed: %v", err)
	}

	state, err := db.q.GetSyncState(ctx, gen.GetSyncStateParams{ChainID: testChainID, ContractName: "indexer"})
	if err != nil {
		t.Fatalf("GetSyncState failed: %v", err)
	}
	if state.LastBlock != 1000 {
		t.Fatalf("LastBlock mismatch: expected 1000, got %d", state.LastBlock)
	}

	// Update sync state
	err = db.q.UpsertSyncState(ctx, gen.UpsertSyncStateParams{ChainID: testChainID,
		ContractName: "indexer",
		LastBlock:    2000,
	})
	if err != nil {
		t.Fatalf("second UpsertSyncState failed: %v", err)
	}

	state, err = db.q.GetSyncState(ctx, gen.GetSyncStateParams{ChainID: testChainID, ContractName: "indexer"})
	if err != nil {
		t.Fatalf("GetSyncState failed: %v", err)
	}
	if state.LastBlock != 2000 {
		t.Fatalf("LastBlock mismatch: expected 2000, got %d", state.LastBlock)
	}

	// Query non-existent key
	_, err = db.q.GetSyncState(ctx, gen.GetSyncStateParams{ChainID: testChainID, ContractName: "nonexistent"})
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got: %v", err)
	}
}

func TestIndexerScenario_Unbound(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	ownerAddr := "0x6666666666666666666666666666666666666666"
	agentAddr := "0x7777777777777777777777777777777777777777"
	subnetID := numericFromInt64(5)

	// Setup: bind agent to owner, create stake allocation
	_ = db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: agentAddr, BoundTo: ownerAddr})
	_ = db.q.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: testChainID, UserAddress: ownerAddr})
	_ = db.q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{ChainID: testChainID,
		UserAddress:  ownerAddr,
		AgentAddress: agentAddr,
		SubnetID:     subnetID,
		Amount:       numericFromInt64(3000),
	})

	// Simulate Unbound event: clear binding
	err := db.q.ClearUserBinding(ctx, gen.ClearUserBindingParams{ChainID: testChainID, Address: agentAddr})
	if err != nil {
		t.Fatalf("ClearUserBinding failed: %v", err)
	}

	// Verify binding is cleared
	user, err := db.q.GetUser(ctx, gen.GetUserParams{ChainID: testChainID, Address: agentAddr})
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if strings.TrimSpace(user.BoundTo) != "" {
		t.Fatalf("BoundTo should be empty after unbind, got %s", user.BoundTo)
	}

	// Freeze agent allocations (simulating what indexer does on unbind)
	err = db.q.FreezeAgentAllocations(ctx, gen.FreezeAgentAllocationsParams{ChainID: testChainID,
		UserAddress:  ownerAddr,
		AgentAddress: agentAddr,
	})
	if err != nil {
		t.Fatalf("FreezeAgentAllocations failed: %v", err)
	}

	// Verify allocations are frozen (GetAgentSubnetStake excludes frozen)
	stake, err := db.q.GetAgentSubnetStake(ctx, gen.GetAgentSubnetStakeParams{ChainID: testChainID,
		AgentAddress: agentAddr,
		SubnetID:     subnetID,
	})
	if err != nil {
		t.Fatalf("GetAgentSubnetStake failed: %v", err)
	}
	assertNumericEqual(t, "stake after freeze (excluding frozen)", stake, 0)

	// Verify frozen allocation still exists
	frozen, err := db.q.GetFrozenByUser(ctx, gen.GetFrozenByUserParams{ChainID: testChainID, UserAddress: ownerAddr})
	if err != nil {
		t.Fatalf("GetFrozenByUser failed: %v", err)
	}
	if len(frozen) != 1 {
		t.Fatalf("expected 1 frozen record, got %d", len(frozen))
	}
	assertNumericEqual(t, "frozen amount", frozen[0].Amount, 3000)
}

func TestIndexerScenario_RecipientSet(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	userAddr := "0xaaaa111111111111111111111111111111111111"
	recipientAddr := "0xbbbb111111111111111111111111111111111111"

	// Simulate RecipientSet event
	err := db.q.UpsertUserRecipient(ctx, gen.UpsertUserRecipientParams{ChainID: testChainID,
		Address:   userAddr,
		Recipient: recipientAddr,
	})
	if err != nil {
		t.Fatalf("UpsertUserRecipient failed: %v", err)
	}

	user, err := db.q.GetUser(ctx, gen.GetUserParams{ChainID: testChainID, Address: userAddr})
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if strings.TrimSpace(user.Recipient) != recipientAddr {
		t.Fatalf("Recipient mismatch: expected %s, got %s", recipientAddr, user.Recipient)
	}

	// Update recipient
	newRecipient := "0xcccc111111111111111111111111111111111111"
	err = db.q.UpsertUserRecipient(ctx, gen.UpsertUserRecipientParams{ChainID: testChainID,
		Address:   userAddr,
		Recipient: newRecipient,
	})
	if err != nil {
		t.Fatalf("UpsertUserRecipient (update) failed: %v", err)
	}

	user, err = db.q.GetUser(ctx, gen.GetUserParams{ChainID: testChainID, Address: userAddr})
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if strings.TrimSpace(user.Recipient) != newRecipient {
		t.Fatalf("Recipient after update mismatch: expected %s, got %s", newRecipient, user.Recipient)
	}
}

func TestIndexerScenario_TransactionAtomicity(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	userAddr := "0xdddd111111111111111111111111111111111111"

	// Setup
	_ = db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: userAddr, BoundTo: ""})
	_ = db.q.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: testChainID, UserAddress: userAddr})

	// Simulate indexer transaction mode: execute multiple operations in a transaction
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		t.Fatalf("Begin failed: %v", err)
	}

	qtx := gen.New(tx)

	// Insert a stake position within a transaction
	err = qtx.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 100, Owner: userAddr, Amount: numericFromInt64(5000),
		LockEndTime: 10, CreatedAt: 100,
	})
	if err != nil {
		_ = tx.Rollback(ctx)
		t.Fatalf("InsertStakePosition within transaction failed: %v", err)
	}

	// Update sync state
	err = qtx.UpsertSyncState(ctx, gen.UpsertSyncStateParams{ChainID: testChainID,
		ContractName: "indexer",
		LastBlock:    500,
	})
	if err != nil {
		_ = tx.Rollback(ctx)
		t.Fatalf("UpsertSyncState within transaction failed: %v", err)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	// Verify all operations took effect
	pos, err := db.q.GetStakePosition(ctx, gen.GetStakePositionParams{ChainID: testChainID, TokenID: 100})
	if err != nil {
		t.Fatalf("GetStakePosition failed: %v", err)
	}
	assertNumericEqual(t, "position amount after transaction", pos.Amount, 5000)

	state, err := db.q.GetSyncState(ctx, gen.GetSyncStateParams{ChainID: testChainID, ContractName: "indexer"})
	if err != nil {
		t.Fatalf("GetSyncState failed: %v", err)
	}
	if state.LastBlock != 500 {
		t.Fatalf("LastBlock mismatch: expected 500, got %d", state.LastBlock)
	}
}

func TestIndexerScenario_TransactionRollback(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	userAddr := "0xeeee111111111111111111111111111111111111"

	// Setup: create a stake position
	_ = db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 200, Owner: userAddr, Amount: numericFromInt64(5000),
		LockEndTime: 10, CreatedAt: 100,
	})

	// Perform some operations in a transaction then rollback
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		t.Fatalf("Begin failed: %v", err)
	}

	qtx := gen.New(tx)
	_ = qtx.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 201, Owner: userAddr, Amount: numericFromInt64(9999),
		LockEndTime: 20, CreatedAt: 200,
	})

	// Rollback
	err = tx.Rollback(ctx)
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	// Verify only original position exists
	total, err := db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: userAddr})
	if err != nil {
		t.Fatalf("GetUserTotalStaked failed: %v", err)
	}
	assertNumericEqual(t, "total staked after rollback", total, 5000)
}

// numericToBigInt converts pgtype.Numeric to *big.Int (handles Exp field)
func numericToBigInt(t *testing.T, n pgtype.Numeric) *big.Int {
	t.Helper()
	if !n.Valid {
		t.Fatal("numeric value is invalid (Valid=false)")
	}
	result := new(big.Int).Set(n.Int)
	if n.Exp > 0 {
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n.Exp)), nil)
		result.Mul(result, scale)
	}
	return result
}

// --- Integration tests: complex multi-step scenarios and invariant validation ---

func TestIntegration_FullUserLifecycle(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	alice := "0xa11ce00000000000000000000000000000000001"
	agentA := "0xage0a00000000000000000000000000000000001"
	subnet1 := numericFromInt64(1)

	// 1. Register user + initialize balance
	err := db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: alice, BoundTo: ""})
	if err != nil {
		t.Fatalf("UpsertUserBinding failed: %v", err)
	}
	err = db.q.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: testChainID, UserAddress: alice})
	if err != nil {
		t.Fatalf("InitUserBalance failed: %v", err)
	}

	// 2. Bind agent to alice
	_ = db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: agentA, BoundTo: alice})

	// 3. Simulate VeAWP.Deposited: create stake position of 10000
	err = db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 1, Owner: alice, Amount: numericFromInt64(10000),
		LockEndTime: 50, CreatedAt: 100,
	})
	if err != nil {
		t.Fatalf("InsertStakePosition failed: %v", err)
	}

	// 4. Simulate Allocated event: allocate 5000 to (agentA, subnet1)
	err = db.q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{ChainID: testChainID,
		UserAddress:  alice,
		AgentAddress: agentA,
		SubnetID:     subnet1,
		Amount:       numericFromInt64(5000),
	})
	if err != nil {
		t.Fatalf("UpsertStakeAllocation failed: %v", err)
	}
	err = db.q.AddUserAllocated(ctx, gen.AddUserAllocatedParams{ChainID: testChainID,
		UserAddress:    alice,
		TotalAllocated: numericFromInt64(5000),
	})
	if err != nil {
		t.Fatalf("AddUserAllocated failed: %v", err)
	}

	// 5. Simulate Deallocated event: reduce by 2000
	err = db.q.SubtractStakeAllocation(ctx, gen.SubtractStakeAllocationParams{ChainID: testChainID,
		UserAddress:  alice,
		AgentAddress: agentA,
		SubnetID:     subnet1,
		Amount:       numericFromInt64(2000),
	})
	if err != nil {
		t.Fatalf("SubtractStakeAllocation failed: %v", err)
	}
	err = db.q.SubtractUserAllocated(ctx, gen.SubtractUserAllocatedParams{ChainID: testChainID,
		UserAddress:    alice,
		TotalAllocated: numericFromInt64(2000),
	})
	if err != nil {
		t.Fatalf("SubtractUserAllocated failed: %v", err)
	}

	// 6. Verify final state: totalStaked=10000, allocated=3000
	totalStaked, err := db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: alice})
	if err != nil {
		t.Fatalf("GetUserTotalStaked failed: %v", err)
	}
	assertNumericEqual(t, "total staked", totalStaked, 10000)

	bal, err := db.q.GetUserBalance(ctx, gen.GetUserBalanceParams{ChainID: testChainID, UserAddress: alice})
	if err != nil {
		t.Fatalf("GetUserBalance failed: %v", err)
	}
	assertNumericEqual(t, "final allocated", bal.TotalAllocated, 3000)

	// unallocated = totalStaked - totalAllocated = 10000 - 3000 = 7000
	stakedBig := numericToBigInt(t, totalStaked)
	allocBig := numericToBigInt(t, bal.TotalAllocated)
	unallocated := new(big.Int).Sub(stakedBig, allocBig)
	if unallocated.Cmp(big.NewInt(7000)) != 0 {
		t.Fatalf("unallocated balance mismatch: expected 7000, got %s", unallocated.String())
	}

	// 7. Verify stake allocation amount=3000
	alloc, err := db.q.GetStakeAllocation(ctx, gen.GetStakeAllocationParams{ChainID: testChainID,
		UserAddress:  alice,
		AgentAddress: agentA,
		SubnetID:     subnet1,
	})
	if err != nil {
		t.Fatalf("GetStakeAllocation failed: %v", err)
	}
	assertNumericEqual(t, "stake allocation amount", alloc.Amount, 3000)
}

func TestIntegration_BalanceInvariant(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	user := "0x1a0b0c0d0e0f0000000000000000000000000001"
	agent := "0x2a0b0c0d0e0f0000000000000000000000000001"
	subnetID := numericFromInt64(1)

	// Setup user
	_ = db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: user, BoundTo: ""})
	_ = db.q.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: testChainID, UserAddress: user})

	// 1. Stake 7000 via two NFT positions
	_ = db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 1, Owner: user, Amount: numericFromInt64(5000),
		LockEndTime: 50, CreatedAt: 100,
	})
	_ = db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 2, Owner: user, Amount: numericFromInt64(2000),
		LockEndTime: 50, CreatedAt: 200,
	})

	// 2. Allocate 3000
	_ = db.q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{ChainID: testChainID,
		UserAddress: user, AgentAddress: agent, SubnetID: subnetID,
		Amount: numericFromInt64(3000),
	})
	_ = db.q.AddUserAllocated(ctx, gen.AddUserAllocatedParams{ChainID: testChainID,
		UserAddress: user, TotalAllocated: numericFromInt64(3000),
	})

	// 3. Allocate 1000 more (allocated=4000)
	_ = db.q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{ChainID: testChainID,
		UserAddress: user, AgentAddress: agent, SubnetID: subnetID,
		Amount: numericFromInt64(1000),
	})
	_ = db.q.AddUserAllocated(ctx, gen.AddUserAllocatedParams{ChainID: testChainID,
		UserAddress: user, TotalAllocated: numericFromInt64(1000),
	})

	// 4. Deallocate 500 (allocated=3500)
	_ = db.q.SubtractStakeAllocation(ctx, gen.SubtractStakeAllocationParams{ChainID: testChainID,
		UserAddress: user, AgentAddress: agent, SubnetID: subnetID,
		Amount: numericFromInt64(500),
	})
	_ = db.q.SubtractUserAllocated(ctx, gen.SubtractUserAllocatedParams{ChainID: testChainID,
		UserAddress: user, TotalAllocated: numericFromInt64(500),
	})

	// 5. Verify invariants
	totalStaked, err := db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: user})
	if err != nil {
		t.Fatalf("GetUserTotalStaked failed: %v", err)
	}
	assertNumericEqual(t, "total_staked", totalStaked, 7000)

	bal, err := db.q.GetUserBalance(ctx, gen.GetUserBalanceParams{ChainID: testChainID, UserAddress: user})
	if err != nil {
		t.Fatalf("GetUserBalance failed: %v", err)
	}
	assertNumericEqual(t, "total_allocated", bal.TotalAllocated, 3500)

	// Invariant 1: totalStaked >= totalAllocated
	stakedBig := numericToBigInt(t, totalStaked)
	allocBig := numericToBigInt(t, bal.TotalAllocated)
	if stakedBig.Cmp(allocBig) < 0 {
		t.Fatalf("invariant violated: totalStaked(%s) < totalAllocated(%s)",
			stakedBig.String(), allocBig.String())
	}

	// Invariant 2: totalStaked - totalAllocated == unallocated (no drift)
	unallocated := new(big.Int).Sub(stakedBig, allocBig)
	expected := big.NewInt(3500)
	if unallocated.Cmp(expected) != 0 {
		t.Fatalf("invariant violated: unallocated expected %s, got %s", expected.String(), unallocated.String())
	}
}

func TestIntegration_MultiUserIsolation(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	alice := "0x3a0b0c0d0e0f0000000000000000000000000001"
	bob := "0x3a0b0c0d0e0f0000000000000000000000000002"
	agent := "0x3a0b0c0d0e0f0000000000000000000000000003"
	subnetID := numericFromInt64(1)

	// 1. Create alice and bob
	_ = db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: alice, BoundTo: ""})
	_ = db.q.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: testChainID, UserAddress: alice})
	_ = db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: bob, BoundTo: ""})
	_ = db.q.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: testChainID, UserAddress: bob})

	// 2. Stake 5000 for alice, 3000 for bob (via veAWP positions)
	_ = db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 1, Owner: alice, Amount: numericFromInt64(5000),
		LockEndTime: 50, CreatedAt: 100,
	})
	_ = db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 2, Owner: bob, Amount: numericFromInt64(3000),
		LockEndTime: 50, CreatedAt: 101,
	})

	// 3. Alice allocates 2000
	_ = db.q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{ChainID: testChainID,
		UserAddress: alice, AgentAddress: agent, SubnetID: subnetID,
		Amount: numericFromInt64(2000),
	})
	_ = db.q.AddUserAllocated(ctx, gen.AddUserAllocatedParams{ChainID: testChainID,
		UserAddress: alice, TotalAllocated: numericFromInt64(2000),
	})

	// 4. Verify bob's staked amount is unaffected
	bobStaked, err := db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: bob})
	if err != nil {
		t.Fatalf("GetUserTotalStaked(bob) failed: %v", err)
	}
	assertNumericEqual(t, "bob total_staked", bobStaked, 3000)

	bobBal, err := db.q.GetUserBalance(ctx, gen.GetUserBalanceParams{ChainID: testChainID, UserAddress: bob})
	if err != nil {
		t.Fatalf("GetUserBalance(bob) failed: %v", err)
	}
	assertNumericEqual(t, "bob total_allocated", bobBal.TotalAllocated, 0)

	// Verify alice's state is correct
	aliceStaked, err := db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: alice})
	if err != nil {
		t.Fatalf("GetUserTotalStaked(alice) failed: %v", err)
	}
	assertNumericEqual(t, "alice total_staked", aliceStaked, 5000)

	aliceBal, err := db.q.GetUserBalance(ctx, gen.GetUserBalanceParams{ChainID: testChainID, UserAddress: alice})
	if err != nil {
		t.Fatalf("GetUserBalance(alice) failed: %v", err)
	}
	assertNumericEqual(t, "alice total_allocated", aliceBal.TotalAllocated, 2000)
}

func TestIntegration_EpochDistributionAccounting(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// 1. Insert epoch, daily_emission=1000
	err := db.q.UpsertEpoch(ctx, gen.UpsertEpochParams{ChainID: testChainID,
		EpochID:       10,
		StartTime:     10000,
		DailyEmission: numericFromInt64(1000),
		DaoEmission:   pgtype.Numeric{Valid: false},
	})
	if err != nil {
		t.Fatalf("UpsertEpoch failed: %v", err)
	}

	// 2. Insert distribution records for 3 recipients (r1=300, r2=200, r3=500)
	distributions := []struct {
		recipient string
		amount    int64
	}{
		{"0x0000000000000000000000000000000000000001", 300},
		{"0x0000000000000000000000000000000000000002", 200},
		{"0x0000000000000000000000000000000000000003", 500},
	}
	for _, d := range distributions {
		err := db.q.InsertRecipientAWPDistribution(ctx, gen.InsertRecipientAWPDistributionParams{ChainID: testChainID,
			EpochID:   10,
			Recipient: d.recipient,
			AwpAmount: numericFromInt64(d.amount),
		})
		if err != nil {
			t.Fatalf("InsertRecipientAWPDistribution(%s) failed: %v", d.recipient, err)
		}
	}

	// 3. Query each recipient's earnings
	for _, d := range distributions {
		earnings, err := db.q.GetRecipientEarnings(ctx, gen.GetRecipientEarningsParams{ChainID: testChainID,
			Recipient: d.recipient, Limit: 10, Offset: 0,
		})
		if err != nil {
			t.Fatalf("GetRecipientEarnings(%s) failed: %v", d.recipient, err)
		}
		if len(earnings) != 1 {
			t.Fatalf("recipient %s: expected 1 earnings record, got %d", d.recipient, len(earnings))
		}
		assertNumericEqual(t, "recipient earnings", earnings[0].AwpAmount, d.amount)
	}

	// 4. Query epoch distribution list, verify total = 1000
	dists, err := db.q.GetEpochDistributions(ctx, gen.GetEpochDistributionsParams{ChainID: testChainID, EpochID: 10})
	if err != nil {
		t.Fatalf("GetEpochDistributions failed: %v", err)
	}
	if len(dists) != 3 {
		t.Fatalf("expected 3 distribution records, got %d", len(dists))
	}
	total := big.NewInt(0)
	for _, d := range dists {
		total.Add(total, numericToBigInt(t, d.AwpAmount))
	}
	if total.Cmp(big.NewInt(1000)) != 0 {
		t.Fatalf("distribution total mismatch: expected 1000, got %s", total.String())
	}

	// 5. Update dao_emission=0 (all rewards distributed to subnets)
	err = db.q.UpdateEpochDAO(ctx, gen.UpdateEpochDAOParams{ChainID: testChainID,
		EpochID:     10,
		DaoEmission: numericFromInt64(0),
	})
	if err != nil {
		t.Fatalf("UpdateEpochDAO failed: %v", err)
	}

	// 6. Verify epoch data consistency: dao_emission has been recorded
	epoch, err := db.q.GetEpoch(ctx, gen.GetEpochParams{ChainID: testChainID, EpochID: 10})
	if err != nil {
		t.Fatalf("GetEpoch failed: %v", err)
	}
	assertNumericEqual(t, "dao_emission", epoch.DaoEmission, 0)
	assertNumericEqual(t, "daily_emission", epoch.DailyEmission, 1000)
}

func TestIntegration_VeAWPLifecycle(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	user := "0x4a0b0c0d0e0f0000000000000000000000000001"

	// 1. Deposit (create position) 5000
	err := db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 1, Owner: user, Amount: numericFromInt64(5000),
		LockEndTime: 50, CreatedAt: 100,
	})
	if err != nil {
		t.Fatalf("InsertStakePosition failed: %v", err)
	}

	total, err := db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: user})
	if err != nil {
		t.Fatalf("GetUserTotalStaked failed: %v", err)
	}
	assertNumericEqual(t, "staked after first deposit", total, 5000)

	// 2. Add to position (PositionIncreased: 5000 -> 7000)
	err = db.q.UpdateStakePosition(ctx, gen.UpdateStakePositionParams{ChainID: testChainID,
		TokenID: 1, Amount: numericFromInt64(7000), LockEndTime: 60,
	})
	if err != nil {
		t.Fatalf("UpdateStakePosition failed: %v", err)
	}

	total, err = db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: user})
	if err != nil {
		t.Fatalf("GetUserTotalStaked failed: %v", err)
	}
	assertNumericEqual(t, "staked after increase", total, 7000)

	// 3. Withdraw (burn) position
	err = db.q.BurnStakePosition(ctx, gen.BurnStakePositionParams{ChainID: testChainID, TokenID: 1})
	if err != nil {
		t.Fatalf("BurnStakePosition failed: %v", err)
	}

	total, err = db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: user})
	if err != nil {
		t.Fatalf("GetUserTotalStaked failed: %v", err)
	}
	assertNumericEqual(t, "staked after burn", total, 0)

	// 4. Verify position is burned
	pos, err := db.q.GetStakePosition(ctx, gen.GetStakePositionParams{ChainID: testChainID, TokenID: 1})
	if err != nil {
		t.Fatalf("GetStakePosition failed: %v", err)
	}
	if !pos.Burned {
		t.Fatal("position should be burned")
	}
}

func TestIntegration_AgentRemovalAndReallocation(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	user := "0x5a0b0c0d0e0f0000000000000000000000000001"
	agent := "0x5a0b0c0d0e0f0000000000000000000000000002"
	subnet1 := numericFromInt64(1)

	// 1. Create user, agent, and stake position
	_ = db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: user, BoundTo: ""})
	_ = db.q.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: testChainID, UserAddress: user})
	_ = db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 1, Owner: user, Amount: numericFromInt64(10000),
		LockEndTime: 50, CreatedAt: 100,
	})
	_ = db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: agent, BoundTo: user})

	// 2. Allocate 5000 to (agent, subnet1)
	err := db.q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{ChainID: testChainID,
		UserAddress: user, AgentAddress: agent, SubnetID: subnet1,
		Amount: numericFromInt64(5000),
	})
	if err != nil {
		t.Fatalf("UpsertStakeAllocation failed: %v", err)
	}
	_ = db.q.AddUserAllocated(ctx, gen.AddUserAllocatedParams{ChainID: testChainID,
		UserAddress: user, TotalAllocated: numericFromInt64(5000),
	})

	// 3. FreezeAgentAllocations — set frozen=true
	err = db.q.FreezeAgentAllocations(ctx, gen.FreezeAgentAllocationsParams{ChainID: testChainID,
		UserAddress:  user,
		AgentAddress: agent,
	})
	if err != nil {
		t.Fatalf("FreezeAgentAllocations failed: %v", err)
	}

	// 4. Verify frozen allocations exist
	frozen, err := db.q.GetFrozenByUser(ctx, gen.GetFrozenByUserParams{ChainID: testChainID, UserAddress: user})
	if err != nil {
		t.Fatalf("GetFrozenByUser failed: %v", err)
	}
	if len(frozen) != 1 {
		t.Fatalf("expected 1 frozen record, got %d", len(frozen))
	}
	assertNumericEqual(t, "frozen amount", frozen[0].Amount, 5000)

	// Verify GetAgentSubnetStake excludes frozen (returns 0)
	stake, err := db.q.GetAgentSubnetStake(ctx, gen.GetAgentSubnetStakeParams{ChainID: testChainID,
		AgentAddress: agent, SubnetID: subnet1,
	})
	if err != nil {
		t.Fatalf("GetAgentSubnetStake failed: %v", err)
	}
	assertNumericEqual(t, "non-frozen stake after freeze", stake, 0)

	// 5. DeleteFrozenAllocations — simulate executePendingOperations release
	err = db.q.DeleteFrozenAllocations(ctx, gen.DeleteFrozenAllocationsParams{ChainID: testChainID, UserAddress: user})
	if err != nil {
		t.Fatalf("DeleteFrozenAllocations failed: %v", err)
	}

	// 6. Verify allocations have been deleted
	frozen, err = db.q.GetFrozenByUser(ctx, gen.GetFrozenByUserParams{ChainID: testChainID, UserAddress: user})
	if err != nil {
		t.Fatalf("GetFrozenByUser(after delete) failed: %v", err)
	}
	if len(frozen) != 0 {
		t.Fatalf("expected 0 frozen records, got %d", len(frozen))
	}

	// GetStakeAllocation should also return ErrNoRows
	_, err = db.q.GetStakeAllocation(ctx, gen.GetStakeAllocationParams{ChainID: testChainID,
		UserAddress: user, AgentAddress: agent, SubnetID: subnet1,
	})
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows (allocation deleted), got: %v", err)
	}

	// 7. SubtractUserAllocated to release allocated quota
	err = db.q.SubtractUserAllocated(ctx, gen.SubtractUserAllocatedParams{ChainID: testChainID,
		UserAddress: user, TotalAllocated: numericFromInt64(5000),
	})
	if err != nil {
		t.Fatalf("SubtractUserAllocated failed: %v", err)
	}

	// Verify final state: allocated=0, totalStaked=10000
	totalStaked, err := db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: user})
	if err != nil {
		t.Fatalf("GetUserTotalStaked failed: %v", err)
	}
	assertNumericEqual(t, "final total_staked", totalStaked, 10000)

	bal, err := db.q.GetUserBalance(ctx, gen.GetUserBalanceParams{ChainID: testChainID, UserAddress: user})
	if err != nil {
		t.Fatalf("GetUserBalance failed: %v", err)
	}
	assertNumericEqual(t, "final total_allocated", bal.TotalAllocated, 0)
}

func TestIndexerScenario_VeAWPPositionDecreased(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	userAddr := "0xdede000000000000000000000000000000000001"

	// Setup: create a stake position with 5000
	err := db.q.InsertStakePosition(ctx, gen.InsertStakePositionParams{ChainID: testChainID,
		TokenID: 50, Owner: userAddr, Amount: numericFromInt64(5000),
		LockEndTime: 100, CreatedAt: 10,
	})
	if err != nil {
		t.Fatalf("InsertStakePosition failed: %v", err)
	}

	// Verify initial state
	total, err := db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: userAddr})
	if err != nil {
		t.Fatalf("GetUserTotalStaked failed: %v", err)
	}
	assertNumericEqual(t, "staked before decrease", total, 5000)

	// Simulate PositionDecreased: partial withdrawal reduces amount to 3000
	// (indexer calls UpdateStakePosition with remainingAmount and unchanged lockEndTime)
	err = db.q.UpdateStakePosition(ctx, gen.UpdateStakePositionParams{ChainID: testChainID,
		TokenID:     50,
		Amount:      numericFromInt64(3000),
		LockEndTime: 100, // unchanged by partialWithdraw
	})
	if err != nil {
		t.Fatalf("UpdateStakePosition (PositionDecreased) failed: %v", err)
	}

	// Verify amount decreased
	pos, err := db.q.GetStakePosition(ctx, gen.GetStakePositionParams{ChainID: testChainID, TokenID: 50})
	if err != nil {
		t.Fatalf("GetStakePosition failed: %v", err)
	}
	assertNumericEqual(t, "amount after decrease", pos.Amount, 3000)
	if pos.LockEndTime != 100 {
		t.Fatalf("lockEndTime should remain unchanged, got %d", pos.LockEndTime)
	}
	if pos.Burned {
		t.Fatal("position should NOT be burned after partial withdrawal")
	}

	// Verify total staked reflects the decrease
	total, err = db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: userAddr})
	if err != nil {
		t.Fatalf("GetUserTotalStaked failed: %v", err)
	}
	assertNumericEqual(t, "staked after decrease", total, 3000)

	// Decrease again to 1000
	err = db.q.UpdateStakePosition(ctx, gen.UpdateStakePositionParams{ChainID: testChainID,
		TokenID:     50,
		Amount:      numericFromInt64(1000),
		LockEndTime: 100,
	})
	if err != nil {
		t.Fatalf("UpdateStakePosition (second decrease) failed: %v", err)
	}

	total, err = db.q.GetUserTotalStaked(ctx, gen.GetUserTotalStakedParams{ChainID: testChainID, Owner: userAddr})
	if err != nil {
		t.Fatalf("GetUserTotalStaked (second decrease) failed: %v", err)
	}
	assertNumericEqual(t, "staked after second decrease", total, 1000)
}

func TestIndexerScenario_Reallocated(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	userAddr := "0xrea100c000000000000000000000000000000001"
	fromAgent := "0xrea100c000000000000000000000000000000002"
	toAgent := "0xrea100c000000000000000000000000000000003"
	fromSubnet := numericFromInt64(10)
	toSubnet := numericFromInt64(20)

	// Setup: register user + balance
	_ = db.q.UpsertUserBinding(ctx, gen.UpsertUserBindingParams{ChainID: testChainID, Address: userAddr, BoundTo: ""})
	_ = db.q.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: testChainID, UserAddress: userAddr})

	// Setup: allocate 5000 to (fromAgent, fromSubnet)
	err := db.q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{ChainID: testChainID,
		UserAddress:  userAddr,
		AgentAddress: fromAgent,
		SubnetID:     fromSubnet,
		Amount:       numericFromInt64(5000),
	})
	if err != nil {
		t.Fatalf("initial UpsertStakeAllocation failed: %v", err)
	}
	_ = db.q.AddUserAllocated(ctx, gen.AddUserAllocatedParams{ChainID: testChainID,
		UserAddress:    userAddr,
		TotalAllocated: numericFromInt64(5000),
	})

	// Simulate Reallocated event: move 2000 from (fromAgent, fromSubnet) to (toAgent, toSubnet)
	// Step 1: Subtract from source
	err = db.q.SubtractStakeAllocation(ctx, gen.SubtractStakeAllocationParams{ChainID: testChainID,
		UserAddress:  userAddr,
		AgentAddress: fromAgent,
		SubnetID:     fromSubnet,
		Amount:       numericFromInt64(2000),
	})
	if err != nil {
		t.Fatalf("SubtractStakeAllocation (Reallocated source) failed: %v", err)
	}

	// Step 2: Add to destination
	err = db.q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{ChainID: testChainID,
		UserAddress:  userAddr,
		AgentAddress: toAgent,
		SubnetID:     toSubnet,
		Amount:       numericFromInt64(2000),
	})
	if err != nil {
		t.Fatalf("UpsertStakeAllocation (Reallocated dest) failed: %v", err)
	}

	// Verify source decreased: 5000 - 2000 = 3000
	sourceStake, err := db.q.GetAgentSubnetStake(ctx, gen.GetAgentSubnetStakeParams{ChainID: testChainID,
		AgentAddress: fromAgent,
		SubnetID:     fromSubnet,
	})
	if err != nil {
		t.Fatalf("GetAgentSubnetStake (source) failed: %v", err)
	}
	assertNumericEqual(t, "source stake after reallocation", sourceStake, 3000)

	// Verify destination received: 2000
	destStake, err := db.q.GetAgentSubnetStake(ctx, gen.GetAgentSubnetStakeParams{ChainID: testChainID,
		AgentAddress: toAgent,
		SubnetID:     toSubnet,
	})
	if err != nil {
		t.Fatalf("GetAgentSubnetStake (dest) failed: %v", err)
	}
	assertNumericEqual(t, "dest stake after reallocation", destStake, 2000)

	// Verify total_allocated unchanged (reallocation does not change user's total allocation)
	bal, err := db.q.GetUserBalance(ctx, gen.GetUserBalanceParams{ChainID: testChainID, UserAddress: userAddr})
	if err != nil {
		t.Fatalf("GetUserBalance failed: %v", err)
	}
	assertNumericEqual(t, "total_allocated after reallocation", bal.TotalAllocated, 5000)

	// Reallocate remaining 3000 from source to same destination
	err = db.q.SubtractStakeAllocation(ctx, gen.SubtractStakeAllocationParams{ChainID: testChainID,
		UserAddress:  userAddr,
		AgentAddress: fromAgent,
		SubnetID:     fromSubnet,
		Amount:       numericFromInt64(3000),
	})
	if err != nil {
		t.Fatalf("SubtractStakeAllocation (full move) failed: %v", err)
	}
	err = db.q.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{ChainID: testChainID,
		UserAddress:  userAddr,
		AgentAddress: toAgent,
		SubnetID:     toSubnet,
		Amount:       numericFromInt64(3000),
	})
	if err != nil {
		t.Fatalf("UpsertStakeAllocation (full move) failed: %v", err)
	}

	// Source should be 0
	sourceStake, err = db.q.GetAgentSubnetStake(ctx, gen.GetAgentSubnetStakeParams{ChainID: testChainID,
		AgentAddress: fromAgent,
		SubnetID:     fromSubnet,
	})
	if err != nil {
		t.Fatalf("GetAgentSubnetStake (source depleted) failed: %v", err)
	}
	assertNumericEqual(t, "source stake fully moved", sourceStake, 0)

	// Destination should be 2000 + 3000 = 5000
	destStake, err = db.q.GetAgentSubnetStake(ctx, gen.GetAgentSubnetStakeParams{ChainID: testChainID,
		AgentAddress: toAgent,
		SubnetID:     toSubnet,
	})
	if err != nil {
		t.Fatalf("GetAgentSubnetStake (dest final) failed: %v", err)
	}
	assertNumericEqual(t, "dest stake final", destStake, 5000)
}
