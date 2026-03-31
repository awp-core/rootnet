package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/redis/go-redis/v9"
)

// mockAllocQuerier 模拟分配查询接口
type mockAllocQuerier struct {
	stakeMap map[string]string // "agent:subnetId" -> amount
}

func (m *mockAllocQuerier) GetAgentSubnetStakeWS(ctx context.Context, chainID int64, agent string, subnetID string) (string, error) {
	key := agent + ":" + subnetID
	if v, ok := m.stakeMap[key]; ok {
		return v, nil
	}
	return "0", nil
}

// newTestRedis creates a Redis client connected to the local test instance; skips the test if unavailable
func newTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skipf("skipping test: cannot connect to Redis: %v", err)
	}
	t.Cleanup(func() { _ = rdb.Close() })
	return rdb
}

// newTestHub creates a Hub instance for testing
func newTestHub(t *testing.T) (*Hub, *redis.Client) {
	t.Helper()
	rdb := newTestRedis(t)
	logger := slog.Default()
	h := NewHub(rdb, logger)
	return h, rdb
}

// startHub starts Hub.Run in the background and returns a cancel function
func startHub(t *testing.T, h *Hub) context.CancelFunc {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	go h.Run(ctx)
	// Wait for the Run loop to start
	time.Sleep(100 * time.Millisecond)
	t.Cleanup(cancel)
	return cancel
}

// waitForClients polls until the number of registered clients in the Hub reaches the expected value
func waitForClients(h *Hub, expected int, timeout time.Duration) int {
	deadline := time.After(timeout)
	for {
		select {
		case <-deadline:
			h.mu.RLock()
			n := len(h.clients)
			h.mu.RUnlock()
			return n
		default:
			h.mu.RLock()
			n := len(h.clients)
			h.mu.RUnlock()
			if n >= expected {
				return n
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func TestNewHub(t *testing.T) {
	rdb := newTestRedis(t)
	logger := slog.Default()

	h := NewHub(rdb, logger)

	if h == nil {
		t.Fatal("NewHub returned nil")
	}
	if h.clients == nil {
		t.Fatal("clients map not initialized")
	}
	if h.broadcast == nil {
		t.Fatal("broadcast channel not initialized")
	}
	if h.register == nil {
		t.Fatal("register channel not initialized")
	}
	if h.unregister == nil {
		t.Fatal("unregister channel not initialized")
	}
	if h.rdb != rdb {
		t.Fatal("rdb not set correctly")
	}
	if h.logger != logger {
		t.Fatal("logger not set correctly")
	}
}

func TestHub_RegisterUnregister(t *testing.T) {
	h, _ := newTestHub(t)
	startHub(t, h)

	client := &Client{
		hub:     h,
		send:    make(chan []byte, 256),
		filters: make(map[string]bool),
	}

	// Register client
	h.register <- client
	time.Sleep(50 * time.Millisecond)

	h.mu.RLock()
	_, tracked := h.clients[client]
	count := len(h.clients)
	h.mu.RUnlock()

	if !tracked {
		t.Fatal("client should be tracked after registration")
	}
	if count != 1 {
		t.Fatalf("expected 1 client, got %d", count)
	}

	// Unregister client
	h.unregister <- client
	time.Sleep(50 * time.Millisecond)

	h.mu.RLock()
	_, tracked = h.clients[client]
	count = len(h.clients)
	h.mu.RUnlock()

	if tracked {
		t.Fatal("client should no longer be tracked after unregistration")
	}
	if count != 0 {
		t.Fatalf("expected 0 clients, got %d", count)
	}
}

func TestHub_BroadcastToClients(t *testing.T) {
	h, _ := newTestHub(t)
	// No need to start Run; test broadcastToClients directly

	client := &Client{
		hub:     h,
		send:    make(chan []byte, 256),
		filters: make(map[string]bool),
	}

	// Manually register client
	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	msg := []byte(`{"type":"TestEvent","data":"hello"}`)
	h.broadcastToClients(msg)

	select {
	case received := <-client.send:
		if string(received) != string(msg) {
			t.Fatalf("expected message %q, got %q", string(msg), string(received))
		}
	case <-time.After(time.Second):
		t.Fatal("timeout: broadcast message not received")
	}
}

func TestHub_FilteredBroadcast(t *testing.T) {
	h, _ := newTestHub(t)

	client := &Client{
		hub:  h,
		send: make(chan []byte, 256),
		filters: map[string]bool{
			"EpochSettled": true,
		},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// Send event that does not match the filter
	nonMatchMsg := []byte(`{"type":"StakeUpdated","data":"foo"}`)
	h.broadcastToClients(nonMatchMsg)

	select {
	case msg := <-client.send:
		t.Fatalf("should not receive non-matching event, but got: %s", string(msg))
	case <-time.After(100 * time.Millisecond):
		// Correct: no message received
	}

	// Send event that matches the filter
	matchMsg := []byte(`{"type":"EpochSettled","data":"bar"}`)
	h.broadcastToClients(matchMsg)

	select {
	case received := <-client.send:
		if string(received) != string(matchMsg) {
			t.Fatalf("expected message %q, got %q", string(matchMsg), string(received))
		}
	case <-time.After(time.Second):
		t.Fatal("timeout: matching broadcast message not received")
	}
}

func TestHub_HandleConnect(t *testing.T) {
	h, _ := newTestHub(t)
	startHub(t, h)

	// HandleConnect internally starts readPump/writePump using r.Context().
	// In httptest.NewServer, net/http cancels r.Context() once ServeHTTP returns,
	// causing readPump to error and trigger unregister immediately.
	// Block handler return with testDone channel to keep r.Context() alive.
	testDone := make(chan struct{})
	t.Cleanup(func() { close(testDone) })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.HandleConnect(w, r)
		<-testDone
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + srv.URL[len("http"):]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify WebSocket connection can be established successfully
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("WebSocket connection failed: %v", err)
	}
	defer func() { _ = conn.Close(websocket.StatusNormalClosure, "") }()

	// Wait for client registration
	count := waitForClients(h, 1, 3*time.Second)
	if count != 1 {
		t.Fatalf("expected 1 client registered via HandleConnect, got %d", count)
	}

	// Send subscription filter over WebSocket
	filterMsg, _ := json.Marshal(filterMessage{Subscribe: []string{"EpochSettled"}})
	if err := conn.Write(ctx, websocket.MessageText, filterMsg); err != nil {
		t.Fatalf("failed to send filter message: %v", err)
	}

	// Wait for readPump to process the filter
	time.Sleep(200 * time.Millisecond)

	// Broadcast a matching event
	matchMsg := []byte(`{"type":"EpochSettled","data":"test"}`)
	h.broadcast <- matchMsg

	// Read message via WebSocket
	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("failed to read WebSocket message: %v", err)
	}
	if string(data) != string(matchMsg) {
		t.Fatalf("expected %q, got %q", string(matchMsg), string(data))
	}
}

func TestHub_SlowClientEviction(t *testing.T) {
	h, _ := newTestHub(t)

	// Create a slow client with a buffer size of 1
	client := &Client{
		hub:     h,
		send:    make(chan []byte, 1),
		filters: make(map[string]bool),
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// Fill the send buffer
	msg1 := []byte(`{"type":"Msg1"}`)
	h.broadcastToClients(msg1)

	// Confirm the client is still tracked after the first message
	h.mu.RLock()
	_, tracked := h.clients[client]
	h.mu.RUnlock()
	if !tracked {
		t.Fatal("client should not be evicted after the first message")
	}

	// Broadcast again; buffer is full so the client should be evicted
	msg2 := []byte(`{"type":"Msg2"}`)
	h.broadcastToClients(msg2)

	h.mu.RLock()
	_, tracked = h.clients[client]
	count := len(h.clients)
	h.mu.RUnlock()

	if tracked {
		t.Fatal("slow client should be evicted")
	}
	if count != 0 {
		t.Fatalf("expected 0 clients, got %d", count)
	}
}

func TestHub_AllocationWatch(t *testing.T) {
	h, _ := newTestHub(t)

	watchClient := &Client{
		hub:          h,
		send:         make(chan []byte, 256),
		filters:      make(map[string]bool),
		allocWatches: map[string]bool{"0xagent1:12345": true},
	}
	normalClient := &Client{
		hub:          h,
		send:         make(chan []byte, 256),
		filters:      make(map[string]bool),
		allocWatches: make(map[string]bool),
	}

	h.mu.Lock()
	h.clients[watchClient] = true
	h.clients[normalClient] = true
	h.mu.Unlock()

	allocMsg := []byte(`{"type":"Allocated","blockNumber":100,"txHash":"0xabc","data":{"staker":"0xStaker1","agent":"0xAgent1","subnetId":"12345","amount":"1000","operator":"0xOp1"}}`)
	h.broadcastToClients(allocMsg)

	// watchClient 应收到 AllocationChanged
	select {
	case msg := <-watchClient.send:
		var evt map[string]interface{}
		if err := json.Unmarshal(msg, &evt); err != nil {
			t.Fatalf("failed to parse: %v", err)
		}
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
		if evt["sourceEvent"] != "Allocated" {
			t.Errorf("expected sourceEvent Allocated, got %v", evt["sourceEvent"])
		}
		if evt["agent"] != "0xagent1" {
			t.Errorf("expected agent 0xagent1, got %v", evt["agent"])
		}
	case <-time.After(time.Second):
		t.Fatal("watchClient should have received AllocationChanged")
	}

	// normalClient 应收到原始事件
	select {
	case msg := <-normalClient.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(msg, &evt)
		if evt["type"] != "Allocated" {
			t.Errorf("normalClient expected Allocated, got %v", evt["type"])
		}
	case <-time.After(time.Second):
		t.Fatal("normalClient should have received Allocated")
	}
}

func TestHub_AllocationWatch_Reallocated(t *testing.T) {
	h, _ := newTestHub(t)

	client := &Client{
		hub:          h,
		send:         make(chan []byte, 256),
		filters:      make(map[string]bool),
		allocWatches: map[string]bool{"0xtoagent:200": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	msg := []byte(`{"type":"Reallocated","blockNumber":60,"txHash":"0xr","data":{"staker":"0xS","fromAgent":"0xFromAgent","fromSubnet":"100","toAgent":"0xToAgent","toSubnet":"200","amount":"500","operator":"0xO"}}`)
	h.broadcastToClients(msg)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
	case <-time.After(time.Second):
		t.Fatal("should receive AllocationChanged for Reallocated to-side")
	}
}

func TestHub_SubnetWatch(t *testing.T) {
	h, _ := newTestHub(t)

	// 仅订阅 subnetId，不指定 agent
	client := &Client{
		hub:           h,
		send:          make(chan []byte, 256),
		filters:       make(map[string]bool),
		allocWatches:  make(map[string]bool),
		subnetWatches: map[string]bool{"12345": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// 任意 agent 在该 subnet 上分配都应推送
	msg := []byte(`{"type":"Allocated","blockNumber":200,"txHash":"0xaaa","data":{"staker":"0xS","agent":"0xAnyAgent","subnetId":"12345","amount":"999","operator":"0xO"}}`)
	h.broadcastToClients(msg)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
		if evt["subnetId"] != "12345" {
			t.Errorf("expected subnetId 12345, got %v", evt["subnetId"])
		}
	case <-time.After(time.Second):
		t.Fatal("subnet watch should receive AllocationChanged for any agent")
	}

	// 不同 subnet 不应触发推送（会走普通广播，但客户端无 filter 所以仍收到原始事件）
	otherMsg := []byte(`{"type":"Allocated","blockNumber":201,"txHash":"0xbbb","data":{"staker":"0xS","agent":"0xOther","subnetId":"99999","amount":"1","operator":"0xO"}}`)
	h.broadcastToClients(otherMsg)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		// 不匹配的 subnet，应以原始 Allocated 发送
		if evt["type"] != "Allocated" {
			t.Errorf("expected original Allocated, got %v", evt["type"])
		}
	case <-time.After(time.Second):
		t.Fatal("should still receive non-matching events as normal broadcast")
	}
}

// ============================================================================
// 1. extractAllocKeys 单元测试
// ============================================================================

func TestExtractAllocKeys_Allocated(t *testing.T) {
	data := map[string]interface{}{
		"staker":   "0xS",
		"agent":    "0xAgent1",
		"subnetId": "12345",
		"amount":   "1000",
	}
	keys := extractAllocKeys("Allocated", data)
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d: %v", len(keys), keys)
	}
	if keys[0] != "0xagent1:12345" {
		t.Errorf("expected 0xagent1:12345, got %s", keys[0])
	}
}

func TestExtractAllocKeys_Deallocated(t *testing.T) {
	data := map[string]interface{}{
		"staker":   "0xS",
		"agent":    "0xAgent2",
		"subnetId": "67890",
		"amount":   "500",
	}
	keys := extractAllocKeys("Deallocated", data)
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d: %v", len(keys), keys)
	}
	if keys[0] != "0xagent2:67890" {
		t.Errorf("expected 0xagent2:67890, got %s", keys[0])
	}
}

func TestExtractAllocKeys_Reallocated(t *testing.T) {
	data := map[string]interface{}{
		"staker":    "0xS",
		"agent":     "0xDirectAgent",
		"subnetId":  "100",
		"fromAgent": "0xFromAgent",
		"fromSubnet": "200",
		"toAgent":   "0xToAgent",
		"toSubnet":  "300",
		"amount":    "1000",
	}
	keys := extractAllocKeys("Reallocated", data)
	// 应返回 3 个 key: agent:subnetId, fromAgent:fromSubnet, toAgent:toSubnet
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d: %v", len(keys), keys)
	}
	expected := []string{
		"0xdirectagent:100",
		"0xfromagent:200",
		"0xtoagent:300",
	}
	for i, exp := range expected {
		if keys[i] != exp {
			t.Errorf("key[%d]: expected %s, got %s", i, exp, keys[i])
		}
	}
}

func TestExtractAllocKeys_MissingFields(t *testing.T) {
	// 缺少 agent
	data1 := map[string]interface{}{
		"subnetId": "12345",
	}
	keys1 := extractAllocKeys("Allocated", data1)
	if len(keys1) != 0 {
		t.Errorf("expected 0 keys when agent missing, got %d: %v", len(keys1), keys1)
	}

	// 缺少 subnetId
	data2 := map[string]interface{}{
		"agent": "0xAgent1",
	}
	keys2 := extractAllocKeys("Allocated", data2)
	if len(keys2) != 0 {
		t.Errorf("expected 0 keys when subnetId missing, got %d: %v", len(keys2), keys2)
	}

	// 两者都缺少
	data3 := map[string]interface{}{
		"staker": "0xS",
	}
	keys3 := extractAllocKeys("Allocated", data3)
	if len(keys3) != 0 {
		t.Errorf("expected 0 keys when both missing, got %d: %v", len(keys3), keys3)
	}
}

func TestExtractAllocKeys_CaseInsensitive(t *testing.T) {
	data := map[string]interface{}{
		"agent":    "0xABCDEF1234567890abcdef1234567890ABCDEF12",
		"subnetId": "999",
	}
	keys := extractAllocKeys("Allocated", data)
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
	// agent 应被转为小写
	if keys[0] != strings.ToLower("0xABCDEF1234567890abcdef1234567890ABCDEF12")+":999" {
		t.Errorf("expected lowercased agent, got %s", keys[0])
	}
}

// ============================================================================
// 2. enrichAllocEvent 单元测试
// ============================================================================

func TestEnrichAllocEvent_NoQuerier(t *testing.T) {
	h, _ := newTestHub(t)
	// 不设置 allocQuery

	evt := broadcastEvent{
		Type: "Allocated",
		Data: map[string]interface{}{
			"agent":    "0xagent1",
			"subnetId": "12345",
			"amount":   "1000",
		},
	}

	result := h.enrichAllocEvent(evt, "0xagent1:12345")

	if result["type"] != "AllocationChanged" {
		t.Errorf("expected type AllocationChanged, got %v", result["type"])
	}
	if _, exists := result["currentStake"]; exists {
		t.Error("currentStake should not be present without querier")
	}
}

func TestEnrichAllocEvent_WithQuerier(t *testing.T) {
	h, _ := newTestHub(t)
	h.SetAllocationQuerier(&mockAllocQuerier{
		stakeMap: map[string]string{
			"0xagent1:12345": "5000",
		},
	}, 56)

	evt := broadcastEvent{
		Type: "Allocated",
		Data: map[string]interface{}{
			"agent":    "0xagent1",
			"subnetId": "12345",
			"amount":   "1000",
		},
	}

	result := h.enrichAllocEvent(evt, "0xagent1:12345")

	if result["currentStake"] != "5000" {
		t.Errorf("expected currentStake 5000, got %v", result["currentStake"])
	}
}

func TestEnrichAllocEvent_Fields(t *testing.T) {
	h, _ := newTestHub(t)

	data := map[string]interface{}{
		"staker":   "0xS",
		"agent":    "0xagent1",
		"subnetId": "12345",
		"amount":   "1000",
	}
	evt := broadcastEvent{
		Type: "Allocated",
		Data: data,
	}

	result := h.enrichAllocEvent(evt, "0xagent1:12345")

	// 验证所有字段都存在
	requiredFields := []string{"type", "sourceEvent", "agent", "subnetId", "data"}
	for _, field := range requiredFields {
		if _, ok := result[field]; !ok {
			t.Errorf("missing required field: %s", field)
		}
	}

	if result["type"] != "AllocationChanged" {
		t.Errorf("expected AllocationChanged, got %v", result["type"])
	}
	if result["sourceEvent"] != "Allocated" {
		t.Errorf("expected sourceEvent Allocated, got %v", result["sourceEvent"])
	}
	if result["agent"] != "0xagent1" {
		t.Errorf("expected agent 0xagent1, got %v", result["agent"])
	}
	if result["subnetId"] != "12345" {
		t.Errorf("expected subnetId 12345, got %v", result["subnetId"])
	}
	// data 应为原始事件数据
	resultData, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data field should be map[string]interface{}")
	}
	if resultData["amount"] != "1000" {
		t.Errorf("expected data.amount 1000, got %v", resultData["amount"])
	}
}

// ============================================================================
// 3. 分配监听广播测试
// ============================================================================

func TestHub_AllocationWatch_Deallocated(t *testing.T) {
	h, _ := newTestHub(t)

	client := &Client{
		hub:          h,
		send:         make(chan []byte, 256),
		filters:      make(map[string]bool),
		allocWatches: map[string]bool{"0xagent1:12345": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	msg := []byte(`{"type":"Deallocated","blockNumber":100,"txHash":"0xabc","data":{"staker":"0xS","agent":"0xAgent1","subnetId":"12345","amount":"500","operator":"0xO"}}`)
	h.broadcastToClients(msg)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
		if evt["sourceEvent"] != "Deallocated" {
			t.Errorf("expected sourceEvent Deallocated, got %v", evt["sourceEvent"])
		}
	case <-time.After(time.Second):
		t.Fatal("should receive AllocationChanged for Deallocated event")
	}
}

func TestHub_AllocationWatch_FromSideReallocated(t *testing.T) {
	h, _ := newTestHub(t)

	// 监听 from-side 的 agent:subnet
	client := &Client{
		hub:          h,
		send:         make(chan []byte, 256),
		filters:      make(map[string]bool),
		allocWatches: map[string]bool{"0xfromagent:100": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	msg := []byte(`{"type":"Reallocated","blockNumber":60,"txHash":"0xr","data":{"staker":"0xS","fromAgent":"0xFromAgent","fromSubnet":"100","toAgent":"0xToAgent","toSubnet":"200","amount":"500","operator":"0xO"}}`)
	h.broadcastToClients(msg)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
		if evt["agent"] != "0xfromagent" {
			t.Errorf("expected agent 0xfromagent, got %v", evt["agent"])
		}
	case <-time.After(time.Second):
		t.Fatal("should receive AllocationChanged for Reallocated from-side")
	}
}

func TestHub_AllocationWatch_MultipleWatches(t *testing.T) {
	h, _ := newTestHub(t)

	client := &Client{
		hub:     h,
		send:    make(chan []byte, 256),
		filters: make(map[string]bool),
		allocWatches: map[string]bool{
			"0xagent1:111": true,
			"0xagent2:222": true,
			"0xagent3:333": true,
		},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// 第一个事件匹配第二个 watch
	msg1 := []byte(`{"type":"Allocated","blockNumber":1,"txHash":"0x1","data":{"staker":"0xS","agent":"0xAgent2","subnetId":"222","amount":"100","operator":"0xO"}}`)
	h.broadcastToClients(msg1)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
		if evt["agent"] != "0xagent2" {
			t.Errorf("expected agent 0xagent2, got %v", evt["agent"])
		}
	case <-time.After(time.Second):
		t.Fatal("should receive AllocationChanged for second watch")
	}

	// 第二个事件匹配第三个 watch
	msg2 := []byte(`{"type":"Allocated","blockNumber":2,"txHash":"0x2","data":{"staker":"0xS","agent":"0xAgent3","subnetId":"333","amount":"200","operator":"0xO"}}`)
	h.broadcastToClients(msg2)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["agent"] != "0xagent3" {
			t.Errorf("expected agent 0xagent3, got %v", evt["agent"])
		}
	case <-time.After(time.Second):
		t.Fatal("should receive AllocationChanged for third watch")
	}
}

func TestHub_AllocationWatch_ExactTakesPriority(t *testing.T) {
	h, _ := newTestHub(t)

	// 客户端同时设置精确匹配和子网级匹配
	client := &Client{
		hub:           h,
		send:          make(chan []byte, 256),
		filters:       make(map[string]bool),
		allocWatches:  map[string]bool{"0xagent1:12345": true},
		subnetWatches: map[string]bool{"12345": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	msg := []byte(`{"type":"Allocated","blockNumber":100,"txHash":"0xabc","data":{"staker":"0xS","agent":"0xAgent1","subnetId":"12345","amount":"1000","operator":"0xO"}}`)
	h.broadcastToClients(msg)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
		// 精确匹配优先，agent 应为精确 watch 的 agent
		if evt["agent"] != "0xagent1" {
			t.Errorf("expected agent 0xagent1, got %v", evt["agent"])
		}
	case <-time.After(time.Second):
		t.Fatal("should receive AllocationChanged")
	}

	// 应只收到一条消息（精确匹配优先，不重复发送子网级匹配）
	select {
	case extra := <-client.send:
		t.Fatalf("should not receive duplicate message, but got: %s", string(extra))
	case <-time.After(200 * time.Millisecond):
		// 正确：无多余消息
	}
}

func TestHub_AllocationWatch_WatchPlusTypeFilter(t *testing.T) {
	h, _ := newTestHub(t)

	// 客户端设置 type filter ["Allocated"] 和 allocWatch
	client := &Client{
		hub:          h,
		send:         make(chan []byte, 256),
		filters:      map[string]bool{"Allocated": true},
		allocWatches: map[string]bool{"0xagent1:12345": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// 匹配 allocWatch 的事件 → 应收到 AllocationChanged
	allocMsg := []byte(`{"type":"Allocated","blockNumber":100,"txHash":"0xabc","data":{"staker":"0xS","agent":"0xAgent1","subnetId":"12345","amount":"1000","operator":"0xO"}}`)
	h.broadcastToClients(allocMsg)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
	case <-time.After(time.Second):
		t.Fatal("should receive AllocationChanged for matching alloc watch")
	}

	// 不匹配 type filter 的非分配事件 → 不应收到
	nonAllocMsg := []byte(`{"type":"EpochSettled","data":"test"}`)
	h.broadcastToClients(nonAllocMsg)

	select {
	case extra := <-client.send:
		t.Fatalf("should not receive non-matching event, but got: %s", string(extra))
	case <-time.After(200 * time.Millisecond):
		// 正确：被 type filter 过滤
	}
}

func TestHub_AllocationWatch_NonAllocEvent(t *testing.T) {
	h, _ := newTestHub(t)

	client := &Client{
		hub:          h,
		send:         make(chan []byte, 256),
		filters:      make(map[string]bool),
		allocWatches: map[string]bool{"0xagent1:12345": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// 非分配事件（EpochSettled）应走普通广播路径
	msg := []byte(`{"type":"EpochSettled","data":{"epoch":"5"}}`)
	h.broadcastToClients(msg)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "EpochSettled" {
			t.Errorf("expected EpochSettled, got %v", evt["type"])
		}
	case <-time.After(time.Second):
		t.Fatal("non-alloc event should be received via normal broadcast")
	}
}

func TestHub_AllocationWatch_NoMatchFallsThrough(t *testing.T) {
	h, _ := newTestHub(t)

	// 监听某个 agent:subnet，但事件是另一个 agent:subnet
	client := &Client{
		hub:          h,
		send:         make(chan []byte, 256),
		filters:      make(map[string]bool),
		allocWatches: map[string]bool{"0xagent1:12345": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// 不匹配的分配事件 → 应作为原始事件发送（因为无 type filter）
	msg := []byte(`{"type":"Allocated","blockNumber":100,"txHash":"0xabc","data":{"staker":"0xS","agent":"0xOtherAgent","subnetId":"99999","amount":"1000","operator":"0xO"}}`)
	h.broadcastToClients(msg)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "Allocated" {
			t.Errorf("expected original Allocated (no match), got %v", evt["type"])
		}
	case <-time.After(time.Second):
		t.Fatal("non-matching alloc event should fall through to normal broadcast")
	}
}

// ============================================================================
// 4. 子网级监听测试
// ============================================================================

func TestHub_SubnetWatch_Deallocated(t *testing.T) {
	h, _ := newTestHub(t)

	client := &Client{
		hub:           h,
		send:          make(chan []byte, 256),
		filters:       make(map[string]bool),
		allocWatches:  make(map[string]bool),
		subnetWatches: map[string]bool{"12345": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	msg := []byte(`{"type":"Deallocated","blockNumber":100,"txHash":"0xabc","data":{"staker":"0xS","agent":"0xAnyAgent","subnetId":"12345","amount":"500","operator":"0xO"}}`)
	h.broadcastToClients(msg)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
		if evt["sourceEvent"] != "Deallocated" {
			t.Errorf("expected sourceEvent Deallocated, got %v", evt["sourceEvent"])
		}
	case <-time.After(time.Second):
		t.Fatal("subnet watch should trigger for Deallocated")
	}
}

func TestHub_SubnetWatch_Reallocated(t *testing.T) {
	h, _ := newTestHub(t)

	// 监听 toSubnet
	client := &Client{
		hub:           h,
		send:          make(chan []byte, 256),
		filters:       make(map[string]bool),
		allocWatches:  make(map[string]bool),
		subnetWatches: map[string]bool{"200": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	msg := []byte(`{"type":"Reallocated","blockNumber":60,"txHash":"0xr","data":{"staker":"0xS","fromAgent":"0xFromAgent","fromSubnet":"100","toAgent":"0xToAgent","toSubnet":"200","amount":"500","operator":"0xO"}}`)
	h.broadcastToClients(msg)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
		if evt["subnetId"] != "200" {
			t.Errorf("expected subnetId 200, got %v", evt["subnetId"])
		}
	case <-time.After(time.Second):
		t.Fatal("subnet watch should trigger for Reallocated to-subnet")
	}
}

func TestHub_SubnetWatch_MultipleSubnets(t *testing.T) {
	h, _ := newTestHub(t)

	client := &Client{
		hub:           h,
		send:          make(chan []byte, 256),
		filters:       make(map[string]bool),
		allocWatches:  make(map[string]bool),
		subnetWatches: map[string]bool{"111": true, "222": true, "333": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// 匹配第二个子网
	msg1 := []byte(`{"type":"Allocated","blockNumber":1,"txHash":"0x1","data":{"staker":"0xS","agent":"0xA","subnetId":"222","amount":"100","operator":"0xO"}}`)
	h.broadcastToClients(msg1)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
		if evt["subnetId"] != "222" {
			t.Errorf("expected subnetId 222, got %v", evt["subnetId"])
		}
	case <-time.After(time.Second):
		t.Fatal("should receive AllocationChanged for subnet 222")
	}

	// 匹配第三个子网
	msg2 := []byte(`{"type":"Allocated","blockNumber":2,"txHash":"0x2","data":{"staker":"0xS","agent":"0xB","subnetId":"333","amount":"200","operator":"0xO"}}`)
	h.broadcastToClients(msg2)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["subnetId"] != "333" {
			t.Errorf("expected subnetId 333, got %v", evt["subnetId"])
		}
	case <-time.After(time.Second):
		t.Fatal("should receive AllocationChanged for subnet 333")
	}
}

func TestHub_SubnetWatch_CombinedWithExact(t *testing.T) {
	h, _ := newTestHub(t)

	// 精确监听 subnet 111 上的 agent1，子网级监听 subnet 222
	client := &Client{
		hub:           h,
		send:          make(chan []byte, 256),
		filters:       make(map[string]bool),
		allocWatches:  map[string]bool{"0xagent1:111": true},
		subnetWatches: map[string]bool{"222": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// 精确匹配事件
	msg1 := []byte(`{"type":"Allocated","blockNumber":1,"txHash":"0x1","data":{"staker":"0xS","agent":"0xAgent1","subnetId":"111","amount":"100","operator":"0xO"}}`)
	h.broadcastToClients(msg1)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
		if evt["agent"] != "0xagent1" {
			t.Errorf("expected agent 0xagent1, got %v", evt["agent"])
		}
	case <-time.After(time.Second):
		t.Fatal("exact watch should trigger")
	}

	// 子网级匹配事件
	msg2 := []byte(`{"type":"Allocated","blockNumber":2,"txHash":"0x2","data":{"staker":"0xS","agent":"0xOtherAgent","subnetId":"222","amount":"200","operator":"0xO"}}`)
	h.broadcastToClients(msg2)

	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "AllocationChanged" {
			t.Errorf("expected AllocationChanged, got %v", evt["type"])
		}
		if evt["subnetId"] != "222" {
			t.Errorf("expected subnetId 222, got %v", evt["subnetId"])
		}
	case <-time.After(time.Second):
		t.Fatal("subnet watch should trigger")
	}
}

// ============================================================================
// 5. WebSocket 集成测试
// ============================================================================

func TestHub_WS_WatchAllocations_Integration(t *testing.T) {
	h, rdb := newTestHub(t)
	startHub(t, h)

	testDone := make(chan struct{})
	t.Cleanup(func() { close(testDone) })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.HandleConnect(w, r)
		<-testDone
	}))
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws"+srv.URL[len("http"):], nil)
	if err != nil {
		t.Fatalf("WebSocket connection failed: %v", err)
	}
	defer func() { _ = conn.Close(websocket.StatusNormalClosure, "") }()

	// 等待客户端注册
	count := waitForClients(h, 1, 3*time.Second)
	if count != 1 {
		t.Fatalf("expected 1 client, got %d", count)
	}

	// 发送 watchAllocations 订阅（地址必须 42 字符: 0x + 40 hex）
	watchMsg, _ := json.Marshal(filterMessage{
		WatchAllocations: []allocationWatch{
			{Agent: "0xAa11223344556677889900aabbccddeeff112233", SubnetID: "55555"},
		},
	})
	if err := conn.Write(ctx, websocket.MessageText, watchMsg); err != nil {
		t.Fatalf("failed to send watch message: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// 通过 Redis 发布分配事件（agent 大小写不同，extractAllocKeys 会 toLower）
	allocEvt := `{"type":"Allocated","blockNumber":500,"txHash":"0xtest","data":{"staker":"0xS","agent":"0xAa11223344556677889900aabbccddeeff112233","subnetId":"55555","amount":"2000","operator":"0xO"}}`
	if err := rdb.Publish(ctx, "chain_events", allocEvt).Err(); err != nil {
		t.Fatalf("failed to publish Redis event: %v", err)
	}

	// 读取 WebSocket 消息
	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("failed to read WebSocket message: %v", err)
	}

	var evt map[string]interface{}
	if err := json.Unmarshal(data, &evt); err != nil {
		t.Fatalf("failed to parse message: %v", err)
	}
	if evt["type"] != "AllocationChanged" {
		t.Errorf("expected AllocationChanged, got %v", evt["type"])
	}
	if evt["sourceEvent"] != "Allocated" {
		t.Errorf("expected sourceEvent Allocated, got %v", evt["sourceEvent"])
	}
}

func TestHub_WS_SubnetWatch_Integration(t *testing.T) {
	h, rdb := newTestHub(t)
	startHub(t, h)

	testDone := make(chan struct{})
	t.Cleanup(func() { close(testDone) })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.HandleConnect(w, r)
		<-testDone
	}))
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws"+srv.URL[len("http"):], nil)
	if err != nil {
		t.Fatalf("WebSocket connection failed: %v", err)
	}
	defer func() { _ = conn.Close(websocket.StatusNormalClosure, "") }()

	count := waitForClients(h, 1, 3*time.Second)
	if count != 1 {
		t.Fatalf("expected 1 client, got %d", count)
	}

	// 发送子网级 watchAllocations（agent 省略）
	watchMsg, _ := json.Marshal(filterMessage{
		WatchAllocations: []allocationWatch{
			{SubnetID: "77777"},
		},
	})
	if err := conn.Write(ctx, websocket.MessageText, watchMsg); err != nil {
		t.Fatalf("failed to send watch message: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// 通过 Redis 发布事件
	allocEvt := `{"type":"Allocated","blockNumber":600,"txHash":"0xsub","data":{"staker":"0xS","agent":"0xRandomAgent1234567890123456789012345678","subnetId":"77777","amount":"3000","operator":"0xO"}}`
	if err := rdb.Publish(ctx, "chain_events", allocEvt).Err(); err != nil {
		t.Fatalf("failed to publish Redis event: %v", err)
	}

	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("failed to read WebSocket message: %v", err)
	}

	var evt map[string]interface{}
	if err := json.Unmarshal(data, &evt); err != nil {
		t.Fatalf("failed to parse message: %v", err)
	}
	if evt["type"] != "AllocationChanged" {
		t.Errorf("expected AllocationChanged, got %v", evt["type"])
	}
	if evt["subnetId"] != "77777" {
		t.Errorf("expected subnetId 77777, got %v", evt["subnetId"])
	}
}

func TestHub_WS_UpdateSubscription(t *testing.T) {
	h, rdb := newTestHub(t)
	startHub(t, h)

	testDone := make(chan struct{})
	t.Cleanup(func() { close(testDone) })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.HandleConnect(w, r)
		<-testDone
	}))
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws"+srv.URL[len("http"):], nil)
	if err != nil {
		t.Fatalf("WebSocket connection failed: %v", err)
	}
	defer func() { _ = conn.Close(websocket.StatusNormalClosure, "") }()

	count := waitForClients(h, 1, 3*time.Second)
	if count != 1 {
		t.Fatalf("expected 1 client, got %d", count)
	}

	// 第一次订阅 agent1:11111（地址必须 42 字符: 0x + 40 hex）
	watch1, _ := json.Marshal(filterMessage{
		WatchAllocations: []allocationWatch{
			{Agent: "0x1111111111111111111111111111111111111111", SubnetID: "11111"},
		},
	})
	if err := conn.Write(ctx, websocket.MessageText, watch1); err != nil {
		t.Fatalf("failed to send first watch: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// 第二次订阅 agent2:22222（应替换第一次）
	watch2, _ := json.Marshal(filterMessage{
		WatchAllocations: []allocationWatch{
			{Agent: "0x2222222222222222222222222222222222222222", SubnetID: "22222"},
		},
	})
	if err := conn.Write(ctx, websocket.MessageText, watch2); err != nil {
		t.Fatalf("failed to send second watch: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// 发布匹配第一个订阅的事件 → 应作为原始事件发送（不再匹配）
	evt1 := `{"type":"Allocated","blockNumber":1,"txHash":"0x1","data":{"staker":"0xS","agent":"0x1111111111111111111111111111111111111111","subnetId":"11111","amount":"100","operator":"0xO"}}`
	if err := rdb.Publish(ctx, "chain_events", evt1).Err(); err != nil {
		t.Fatalf("failed to publish: %v", err)
	}

	_, data1, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}
	var parsed1 map[string]interface{}
	_ = json.Unmarshal(data1, &parsed1)
	// 第一个订阅已被替换，应收到原始 Allocated（不是 AllocationChanged）
	if parsed1["type"] != "Allocated" {
		t.Errorf("expected original Allocated (old watch replaced), got %v", parsed1["type"])
	}

	// 发布匹配第二个订阅的事件 → 应收到 AllocationChanged
	evt2 := `{"type":"Allocated","blockNumber":2,"txHash":"0x2","data":{"staker":"0xS","agent":"0x2222222222222222222222222222222222222222","subnetId":"22222","amount":"200","operator":"0xO"}}`
	if err := rdb.Publish(ctx, "chain_events", evt2).Err(); err != nil {
		t.Fatalf("failed to publish: %v", err)
	}

	_, data2, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}
	var parsed2 map[string]interface{}
	_ = json.Unmarshal(data2, &parsed2)
	if parsed2["type"] != "AllocationChanged" {
		t.Errorf("expected AllocationChanged for new watch, got %v", parsed2["type"])
	}
}

// ============================================================================
// 6. 边缘情况测试
// ============================================================================

func TestHub_AllocationWatch_SlowClientEviction(t *testing.T) {
	h, _ := newTestHub(t)

	// buffer 大小为 1，模拟慢客户端（使用非分配事件填满 buffer，然后用分配事件触发驱逐）
	client := &Client{
		hub:          h,
		send:         make(chan []byte, 1),
		filters:      make(map[string]bool),
		allocWatches: map[string]bool{"0xagent1:12345": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// 用非分配事件填满 buffer（走普通广播路径）
	msg1 := []byte(`{"type":"EpochSettled","data":{"epoch":"1"}}`)
	h.broadcastToClients(msg1)

	h.mu.RLock()
	_, tracked := h.clients[client]
	h.mu.RUnlock()
	if !tracked {
		t.Fatal("client should still be tracked after first message")
	}

	// 再发一条非分配事件：buffer 已满，应被驱逐
	msg2 := []byte(`{"type":"EpochSettled","data":{"epoch":"2"}}`)
	h.broadcastToClients(msg2)

	h.mu.RLock()
	_, tracked = h.clients[client]
	clientCount := len(h.clients)
	h.mu.RUnlock()

	if tracked {
		t.Fatal("slow client with alloc watches should be evicted when buffer full")
	}
	if clientCount != 0 {
		t.Fatalf("expected 0 clients, got %d", clientCount)
	}
}

func TestHub_AllocationWatch_EmptyData(t *testing.T) {
	h, _ := newTestHub(t)

	client := &Client{
		hub:          h,
		send:         make(chan []byte, 256),
		filters:      make(map[string]bool),
		allocWatches: map[string]bool{"0xagent1:12345": true},
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// 分配事件缺少 data 字段
	msg := []byte(`{"type":"Allocated","blockNumber":100,"txHash":"0xabc"}`)
	h.broadcastToClients(msg)

	// 因为无法提取 allocKeys，应走普通广播
	select {
	case received := <-client.send:
		var evt map[string]interface{}
		_ = json.Unmarshal(received, &evt)
		if evt["type"] != "Allocated" {
			t.Errorf("expected original Allocated, got %v", evt["type"])
		}
	case <-time.After(time.Second):
		t.Fatal("should receive message via normal broadcast")
	}
}

func TestHub_AllocationWatch_InvalidWatchParams(t *testing.T) {
	h, _ := newTestHub(t)
	startHub(t, h)

	testDone := make(chan struct{})
	t.Cleanup(func() { close(testDone) })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.HandleConnect(w, r)
		<-testDone
	}))
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws"+srv.URL[len("http"):], nil)
	if err != nil {
		t.Fatalf("WebSocket connection failed: %v", err)
	}
	defer func() { _ = conn.Close(websocket.StatusNormalClosure, "") }()

	count := waitForClients(h, 1, 3*time.Second)
	if count != 1 {
		t.Fatalf("expected 1 client, got %d", count)
	}

	// 发送无效的 watchAllocations（地址太短和太长）
	watchMsg, _ := json.Marshal(filterMessage{
		WatchAllocations: []allocationWatch{
			{Agent: "0xShort", SubnetID: "12345"},                                          // 地址太短（!= 42 字符）
			{Agent: "0xTooLongAddress1234567890123456789012345678901234567890", SubnetID: "12345"}, // 地址太长
			{Agent: "0xValid1234567890123456789012345678901234", SubnetID: ""},              // subnetId 为空
		},
	})
	if err := conn.Write(ctx, websocket.MessageText, watchMsg); err != nil {
		t.Fatalf("failed to send watch message: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// 验证客户端的 watches 为空（所有条目都应被过滤）
	h.mu.RLock()
	var allocCount, subnetCount int
	for c := range h.clients {
		allocCount = len(c.allocWatches)
		subnetCount = len(c.subnetWatches)
	}
	h.mu.RUnlock()

	if allocCount != 0 {
		t.Errorf("expected 0 allocWatches after invalid params, got %d", allocCount)
	}
	if subnetCount != 0 {
		t.Errorf("expected 0 subnetWatches after invalid params, got %d", subnetCount)
	}
}
