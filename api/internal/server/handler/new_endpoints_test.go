package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/cortexia/rootnet/api/internal/db/gen"
)

// ════════════════════════════════════════════════════════════
// Portfolio — 2 tests
// ════════════════════════════════════════════════════════════

func TestUsersGetPortfolio_Empty(t *testing.T) {
	env := newTestEnv(t)
	addr := "0x0000000000000000000000000000000000000099"

	rr := env.request("GET", "/api/users/"+addr+"/portfolio", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["address"] != addr {
		t.Errorf("expected address=%s, got %v", addr, resp["address"])
	}
	if resp["isRegistered"] != false {
		t.Errorf("expected isRegistered=false for new user, got %v", resp["isRegistered"])
	}
	if resp["boundTo"] != "" {
		t.Errorf("expected empty boundTo, got %v", resp["boundTo"])
	}
	if resp["recipient"] != "" {
		t.Errorf("expected empty recipient, got %v", resp["recipient"])
	}

	// Balance should be zero
	balance, ok := resp["balance"].(map[string]any)
	if !ok {
		t.Fatal("expected balance object in response")
	}
	if balance["totalStaked"] != "0" {
		t.Errorf("expected totalStaked=0, got %v", balance["totalStaked"])
	}
}

func TestUsersGetPortfolio_WithData(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	addr := "0x0000000000000000000000000000000000000001"
	agentAddr := "0x00000000000000000000000000000000000000a1"

	// Create user and set recipient
	_ = env.queries.UpsertUserRecipient(ctx, gen.UpsertUserRecipientParams{
		ChainID:   31337,
		Address:   addr,
		Recipient: "0x0000000000000000000000000000000000000002",
	})

	// Create balance and staking
	_ = env.queries.InitUserBalance(ctx, gen.InitUserBalanceParams{ChainID: 31337, UserAddress: addr})
	_ = env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		ChainID: 31337, TokenID: 1, Owner: addr, Amount: numericFromInt64(10000),
		LockEndTime: 50, CreatedAt: 100,
	})

	// Create allocation
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID:      31337,
		UserAddress:  addr,
		AgentAddress: agentAddr,
		SubnetID:     numericFromInt64(1),
		Amount:       numericFromInt64(5000),
	})

	// Bind agent (mock delegates)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		agentAddr, addr,
	)

	rr := env.request("GET", "/api/users/"+addr+"/portfolio", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["isRegistered"] != true {
		t.Errorf("expected isRegistered=true, got %v", resp["isRegistered"])
	}
	if !strings.Contains(resp["recipient"].(string), "0000002") {
		t.Errorf("expected recipient to contain 0000002, got %v", resp["recipient"])
	}

	// Verify positions exist
	positions, ok := resp["positions"].([]any)
	if !ok {
		t.Fatal("expected positions array in response")
	}
	if len(positions) != 1 {
		t.Errorf("expected 1 position, got %d", len(positions))
	}

	// Verify delegates exist
	delegates, ok := resp["delegates"].([]any)
	if !ok {
		t.Fatal("expected delegates array in response")
	}
	if len(delegates) != 1 {
		t.Errorf("expected 1 delegate, got %d", len(delegates))
	}
}

// ════════════════════════════════════════════════════════════
// Subnets Ranked — 2 tests
// ════════════════════════════════════════════════════════════

func TestSubnetsListRanked_Empty(t *testing.T) {
	env := newTestEnv(t)

	rr := env.request("GET", "/api/subnets/ranked?page=1&limit=10", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp []any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if len(resp) != 0 {
		t.Errorf("expected 0 ranked worknets, got %d", len(resp))
	}
}

func TestSubnetsListRanked_WithData(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// Insert two subnets
	insertSubnet(t, env.pool, 1, "0x0000000000000000000000000000000000000001", "NetA", "NA", "Active")
	insertSubnet(t, env.pool, 2, "0x0000000000000000000000000000000000000001", "NetB", "NB", "Active")

	// Give NetA more stake (use values without trailing zeros to avoid numeric format issues)
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID:      31337,
		UserAddress:  "0x0000000000000000000000000000000000000001",
		AgentAddress: "0x00000000000000000000000000000000000000a1",
		SubnetID:     numericFromInt64(1),
		Amount:       numericFromInt64(10001),
	})
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID:      31337,
		UserAddress:  "0x0000000000000000000000000000000000000001",
		AgentAddress: "0x00000000000000000000000000000000000000a1",
		SubnetID:     numericFromInt64(2),
		Amount:       numericFromInt64(3001),
	})

	rr := env.request("GET", "/api/subnets/ranked?page=1&limit=10", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp []map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if len(resp) != 2 {
		t.Fatalf("expected 2 ranked worknets, got %d", len(resp))
	}

	// First should be the one with the most stake (NetA = 10001)
	if resp[0]["name"] != "NetA" {
		t.Errorf("expected first ranked worknet to be NetA, got %v", resp[0]["name"])
	}
	if resp[0]["totalStake"] != "10001" {
		t.Errorf("expected totalStake=10001, got %v", resp[0]["totalStake"])
	}
}

// ════════════════════════════════════════════════════════════
// Subnets Agents — 2 tests
// ════════════════════════════════════════════════════════════

func TestSubnetsListAgents_Empty(t *testing.T) {
	env := newTestEnv(t)
	insertSubnet(t, env.pool, 1, "0x0000000000000000000000000000000000000001", "Sub1", "S1", "Active")

	rr := env.request("GET", "/api/subnets/1/agents?page=1&limit=10", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp []any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if len(resp) != 0 {
		t.Errorf("expected 0 agents, got %d", len(resp))
	}
}

func TestSubnetsListAgents_WithData(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	insertSubnet(t, env.pool, 1, "0x0000000000000000000000000000000000000001", "Sub1", "S1", "Active")

	agent1 := "0x00000000000000000000000000000000000000a1"
	agent2 := "0x00000000000000000000000000000000000000a2"
	user := "0x0000000000000000000000000000000000000001"

	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID: 31337, UserAddress: user, AgentAddress: agent1,
		SubnetID: numericFromInt64(1), Amount: numericFromInt64(8000),
	})
	_ = env.queries.UpsertStakeAllocation(ctx, gen.UpsertStakeAllocationParams{
		ChainID: 31337, UserAddress: user, AgentAddress: agent2,
		SubnetID: numericFromInt64(1), Amount: numericFromInt64(2000),
	})

	rr := env.request("GET", "/api/subnets/1/agents?page=1&limit=10", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp []map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if len(resp) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(resp))
	}

	// Sorted by stake: agent1 should be first
	if resp[0]["agentAddress"] != agent1 {
		t.Errorf("expected first agent=%s, got %v", agent1, resp[0]["agentAddress"])
	}
}

// ════════════════════════════════════════════════════════════
// Subnets Search — 2 tests
// ════════════════════════════════════════════════════════════

func TestSubnetsSearch_ByName(t *testing.T) {
	env := newTestEnv(t)

	insertSubnet(t, env.pool, 1, "0x0000000000000000000000000000000000000001", "AlphaNet", "AN", "Active")
	insertSubnet(t, env.pool, 2, "0x0000000000000000000000000000000000000001", "BetaNet", "BN", "Active")
	insertSubnet(t, env.pool, 3, "0x0000000000000000000000000000000000000001", "AlphaPlus", "AP", "Active")

	rr := env.request("GET", "/api/subnets/search?q=Alpha", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp []any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if len(resp) != 2 {
		t.Errorf("expected 2 results for 'Alpha', got %d", len(resp))
	}
}

func TestSubnetsSearch_NoMatch(t *testing.T) {
	env := newTestEnv(t)

	insertSubnet(t, env.pool, 1, "0x0000000000000000000000000000000000000001", "AlphaNet", "AN", "Active")

	rr := env.request("GET", "/api/subnets/search?q=NonExistent", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp []any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if len(resp) != 0 {
		t.Errorf("expected 0 results, got %d", len(resp))
	}
}

// ════════════════════════════════════════════════════════════
// Subnets By Owner
// ════════════════════════════════════════════════════════════

func TestSubnetsGetByOwner(t *testing.T) {
	env := newTestEnv(t)

	owner1 := "0x0000000000000000000000000000000000000001"
	owner2 := "0x0000000000000000000000000000000000000002"

	insertSubnet(t, env.pool, 1, owner1, "Net1", "N1", "Active")
	insertSubnet(t, env.pool, 2, owner1, "Net2", "N2", "Active")
	insertSubnet(t, env.pool, 3, owner2, "Net3", "N3", "Active")

	rr := env.request("GET", "/api/subnets/by-owner/"+owner1, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp []any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if len(resp) != 2 {
		t.Errorf("expected 2 subnets for owner1, got %d", len(resp))
	}
}

// ════════════════════════════════════════════════════════════
// Emission Epoch Detail — 2 tests
// ════════════════════════════════════════════════════════════

func TestEmissionGetEpochDetail_Empty(t *testing.T) {
	env := newTestEnv(t)

	rr := env.request("GET", "/api/emission/epochs/999", "")
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for non-existent epoch, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestEmissionGetEpochDetail_WithDistributions(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// Insert epoch (use values without trailing zeros to avoid numeric format issues)
	_, err := env.pool.Exec(ctx,
		"INSERT INTO epochs (chain_id, epoch_id, start_time, daily_emission) VALUES (31337, $1, $2, $3)",
		5, 86400, 15800001,
	)
	if err != nil {
		t.Fatalf("failed to insert epoch: %v", err)
	}

	// Insert distribution records
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO recipient_awp_distributions (chain_id, epoch_id, recipient, awp_amount) VALUES (31337, $1, $2, $3)",
		5, "0x00000000000000000000000000000000000000c1", 10000001,
	)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO recipient_awp_distributions (chain_id, epoch_id, recipient, awp_amount) VALUES (31337, $1, $2, $3)",
		5, "0x00000000000000000000000000000000000000c2", 5000001,
	)

	rr := env.request("GET", "/api/emission/epochs/5", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["epochId"] != float64(5) {
		t.Errorf("expected epochId=5, got %v", resp["epochId"])
	}
	if resp["dailyEmission"] != "15800001" {
		t.Errorf("expected dailyEmission=15800001, got %v", resp["dailyEmission"])
	}

	distributions, ok := resp["distributions"].([]any)
	if !ok {
		t.Fatal("expected distributions array")
	}
	if len(distributions) != 2 {
		t.Errorf("expected 2 distributions, got %d", len(distributions))
	}
}

// ════════════════════════════════════════════════════════════
// Staking Positions Global
// ════════════════════════════════════════════════════════════

func TestStakingGetPositionsGlobal(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	addr := "0x0000000000000000000000000000000000000001"

	// Insert positions on different chains (use values without trailing zeros to avoid numeric representation differences)
	_ = env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		ChainID: 31337, TokenID: 1, Owner: addr, Amount: numericFromInt64(5001),
		LockEndTime: 50, CreatedAt: 100,
	})
	_ = env.queries.InsertStakePosition(ctx, gen.InsertStakePositionParams{
		ChainID: 31337, TokenID: 2, Owner: addr, Amount: numericFromInt64(3001),
		LockEndTime: 100, CreatedAt: 200,
	})

	rr := env.request("GET", "/api/staking/user/"+addr+"/positions/global", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["totalStaked"] != "8002" {
		t.Errorf("expected totalStaked=8002, got %v", resp["totalStaked"])
	}
	if resp["count"] != float64(2) {
		t.Errorf("expected count=2, got %v", resp["count"])
	}

	positions, ok := resp["positions"].([]any)
	if !ok {
		t.Fatal("expected positions array")
	}
	if len(positions) != 2 {
		t.Fatalf("expected 2 positions, got %d", len(positions))
	}
}

// ════════════════════════════════════════════════════════════
// Delegates
// ════════════════════════════════════════════════════════════

func TestUsersGetDelegates(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	ownerAddr := "0x0000000000000000000000000000000000000001"
	agent1 := "0x00000000000000000000000000000000000000a1"
	agent2 := "0x00000000000000000000000000000000000000a2"

	// Bind two agents to owner
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		agent1, ownerAddr,
	)
	_, _ = env.pool.Exec(ctx,
		"INSERT INTO users (chain_id, address, bound_to) VALUES (31337, $1, $2) ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to",
		agent2, ownerAddr,
	)

	rr := env.request("GET", "/api/users/"+ownerAddr+"/delegates", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp []any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if len(resp) != 2 {
		t.Errorf("expected 2 delegates, got %d", len(resp))
	}
}

func TestUsersGetDelegates_Empty(t *testing.T) {
	env := newTestEnv(t)

	rr := env.request("GET", "/api/users/0x0000000000000000000000000000000000000099/delegates", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp []any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if len(resp) != 0 {
		t.Errorf("expected 0 delegates, got %d", len(resp))
	}
}
