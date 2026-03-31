package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/redis/go-redis/v9"
)

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
