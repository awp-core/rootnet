package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/redis/go-redis/v9"
)

// allocationWatch subscribes to allocation changes for a specific (agent, worknetId)
type allocationWatch struct {
	Agent    string `json:"agent"`
	SubnetID string `json:"worknetId"`
}

// Client represents a WebSocket client connection
type Client struct {
	hub              *Hub
	conn             *websocket.Conn
	send             chan []byte
	filters          map[string]bool // event type filters
	allocWatches     map[string]bool // "agent:worknetId" exact match
	subnetWatches    map[string]bool // "worknetId" whole-subnet match (when agent is omitted)
	watchAddresses   map[string]bool // user-level address filter (lowercase address)
	ip               string          // client IP for connection limit tracking
}

// AllocationQuerier queries current allocation for (agent, worknetId) (injected by handler layer)
type AllocationQuerier interface {
	GetAgentSubnetStakeWS(ctx context.Context, chainID int64, agent string, subnetID string) (string, error)
}

// Hub maintains all active WebSocket clients and broadcasts Redis messages to them
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	rdb        *redis.Client
	logger     *slog.Logger
	connsByIP  map[string]int // per-IP WebSocket connection counter
	ipMu       sync.Mutex    // guards connsByIP
	allocQuery AllocationQuerier // optional: for enriching allocation push
	chainID    int64             // chain ID for DB queries
}

// filterMessage represents a subscription filter message sent by a client
type filterMessage struct {
	Subscribe        []string          `json:"subscribe,omitempty"`
	WatchAllocations []allocationWatch `json:"watchAllocations,omitempty"`
	WatchAddresses   []string          `json:"watchAddresses,omitempty"`
}

// broadcastEvent is used to parse the event type, chainId, and data from a broadcast message
type broadcastEvent struct {
	Type    string                 `json:"type"`
	ChainID int64                  `json:"chainId"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// broadcastEventType parses only the type field (for scenarios where data is not needed)
type broadcastEventType struct {
	Type    string `json:"type"`
	ChainID int64  `json:"chainId"`
}

// NewHub creates a new WebSocket Hub instance
func NewHub(rdb *redis.Client, logger *slog.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rdb:        rdb,
		logger:     logger,
		connsByIP:  make(map[string]int),
	}
}

// SetAllocationQuerier injects the allocation query interface (used to include current balance in push)
// Must be called before Hub.Run()
func (h *Hub) SetAllocationQuerier(q AllocationQuerier, chainID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.allocQuery = q
	h.chainID = chainID
}

// Run starts the Hub's main loop, subscribes to the Redis channel, and handles client register/unregister/broadcast
func (h *Hub) Run(ctx context.Context) {
	pubsub := h.rdb.Subscribe(ctx, "chain_events")
	redisCh := pubsub.Channel()
	h.logger.Info("WebSocket Hub started, listening on chain_events channel")

	for {
		select {
		case <-ctx.Done():
			h.logger.Info("WebSocket Hub shutting down")
			_ = pubsub.Close()
			h.mu.Lock()
			for client := range h.clients {
				close(client.send)
				client.conn.Close(websocket.StatusGoingAway, "server shutting down")
				h.ipMu.Lock()
				h.connsByIP[client.ip]--
				if h.connsByIP[client.ip] <= 0 {
					delete(h.connsByIP, client.ip)
				}
				h.ipMu.Unlock()
				delete(h.clients, client)
			}
			h.mu.Unlock()
			return

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.logger.Debug("client connected", "total", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			h.ipMu.Lock()
			h.connsByIP[client.ip]--
			if h.connsByIP[client.ip] <= 0 {
				delete(h.connsByIP, client.ip)
			}
			h.ipMu.Unlock()
			h.logger.Debug("client disconnected", "total", len(h.clients))

		case msg := <-h.broadcast:
			h.broadcastToClients(msg)

		case redisMsg, ok := <-redisCh:
			if !ok {
				h.logger.Warn("Redis subscription channel closed, attempting reconnect...")
				_ = pubsub.Close()
				// Reconnect loop with exponential backoff — keep client connections alive
				backoff := time.Second
				const maxBackoff = 30 * time.Second
				for {
					select {
					case <-ctx.Done():
						h.logger.Info("WebSocket Hub shutting down during reconnect")
						h.mu.Lock()
						for client := range h.clients {
							close(client.send)
							client.conn.Close(websocket.StatusGoingAway, "server shutting down")
							h.ipMu.Lock()
							h.connsByIP[client.ip]--
							if h.connsByIP[client.ip] <= 0 {
								delete(h.connsByIP, client.ip)
							}
							h.ipMu.Unlock()
							delete(h.clients, client)
						}
						h.mu.Unlock()
						return
					default:
					}
					time.Sleep(backoff)
					pubsub = h.rdb.Subscribe(ctx, "chain_events")
					// Verify the connection actually succeeded
					pingCtx, pingCancel := context.WithTimeout(ctx, 3*time.Second)
					err := pubsub.Ping(pingCtx)
					pingCancel()
					if err == nil {
						redisCh = pubsub.Channel()
						h.logger.Info("Redis Pub/Sub reconnected")
						break
					}
					h.logger.Warn("Redis reconnect failed, retrying", "backoff", backoff, "error", err)
					_ = pubsub.Close()
					if backoff < maxBackoff {
						backoff *= 2
						if backoff > maxBackoff {
							backoff = maxBackoff
						}
					}
				}
				continue // restart the main select loop with new redisCh
			}
			// Broadcast Redis message to all clients
			h.broadcastToClients([]byte(redisMsg.Payload))
		}
	}
}

// allocEventTypes are event types that trigger allocation push
var allocEventTypes = map[string]bool{
	"Allocated":   true,
	"Deallocated": true,
	"Reallocated": true,
}

// allocDelivery records allocation changes to push (collected outside lock, DB queried outside lock, sent outside lock)
type allocDelivery struct {
	client *Client
	key    string // "agent:worknetId"
}

// broadcastToClients sends a message to all connected clients (respecting filter settings)
func (h *Hub) broadcastToClients(msg []byte) {
	// Parse type first (ignoring data type to avoid unmarshal failures from type mismatch)
	var evtType broadcastEventType
	eventType := ""
	if err := json.Unmarshal(msg, &evtType); err == nil {
		eventType = evtType.Type
	}

	// Only parse full data for allocation events
	var evt broadcastEvent
	var allocKeys []string
	if allocEventTypes[eventType] {
		if err := json.Unmarshal(msg, &evt); err == nil && evt.Data != nil {
			allocKeys = extractAllocKeys(eventType, evt.Data)
		}
	}

	// Extract involved worknetId list (for subnet-level matching)
	var allocSubnets []string
	if len(allocKeys) > 0 {
		seen := make(map[string]bool)
		for _, key := range allocKeys {
			if parts := strings.SplitN(key, ":", 2); len(parts) == 2 {
				if !seen[parts[1]] {
					allocSubnets = append(allocSubnets, parts[1])
					seen[parts[1]] = true
				}
			}
		}
	}

	// Phase 1: Collect client list to push while holding lock (no DB queries)
	var deliveries []allocDelivery

	h.mu.Lock()
	for client := range h.clients {
		hasWatches := len(client.allocWatches) > 0 || len(client.subnetWatches) > 0
		matched := false

		// 1. Allocation subscription push (collect matching keys, do not send)
		if hasWatches && len(allocKeys) > 0 {
			matchedKeys := make(map[string]bool)

			// 1a. Exact match agent:worknetId (no break, collect all matches)
			for _, key := range allocKeys {
				if client.allocWatches[key] {
					matchedKeys[key] = true
				}
			}

			// 1b. Subnet-level match (collect all keys matching the subnet)
			for _, sid := range allocSubnets {
				if client.subnetWatches[sid] {
					for _, key := range allocKeys {
						if strings.HasSuffix(key, ":"+sid) && !matchedKeys[key] {
							matchedKeys[key] = true
						}
					}
				}
			}

			for key := range matchedKeys {
				deliveries = append(deliveries, allocDelivery{client: client, key: key})
				matched = true
			}
		}

		// 2. Regular type-filtered push (send original message directly, no DB query needed)
		if !matched {
			// 2a. Address-level filter: if client set watchAddresses, only push events involving those addresses
			if len(client.watchAddresses) > 0 {
				if evt.Data == nil && allocEventTypes[eventType] {
					// Already parsed above
				} else if evt.Data == nil {
					// Need to parse data to check address fields
					var addrEvt broadcastEvent
					if err := json.Unmarshal(msg, &addrEvt); err == nil {
						evt = addrEvt
					}
				}
				if !matchesWatchAddresses(evt.Data, client.watchAddresses) {
					continue
				}
			}

			if len(client.filters) > 0 && eventType != "" {
				if !client.filters[eventType] {
					continue
				}
			}
			select {
			case client.send <- msg:
			default:
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
	h.mu.Unlock()

	// Phase 2: Do DB queries and send enriched messages outside lock (deduplicated by key)
	if len(deliveries) > 0 {
		// 2a: Deduplicate by watchKey, query DB only once per key
		enrichCache := make(map[string][]byte)
		for _, d := range deliveries {
			if _, ok := enrichCache[d.key]; !ok {
				enriched := h.enrichAllocEvent(evt, d.key)
				if data, err := json.Marshal(enriched); err == nil {
					enrichCache[d.key] = data
				}
			}
		}

		// 2b: Send cached results to all matching clients
		for _, d := range deliveries {
			data, ok := enrichCache[d.key]
			if !ok {
				continue
			}
			select {
			case d.client.send <- data:
			default:
				h.mu.Lock()
				if _, ok := h.clients[d.client]; ok {
					delete(h.clients, d.client)
					close(d.client.send)
				}
				h.mu.Unlock()
			}
		}
	}
}

// addressFields are field names in event data that may contain addresses
var addressFields = []string{"staker", "user", "agent", "owner", "from", "to", "proposer", "recipient", "fromAgent", "toAgent"}

// matchesWatchAddresses checks if address fields in event data match the client watchAddresses
func matchesWatchAddresses(data map[string]interface{}, watchAddresses map[string]bool) bool {
	if data == nil || len(watchAddresses) == 0 {
		return false
	}
	for _, field := range addressFields {
		if v, ok := data[field]; ok {
			if addr, ok := v.(string); ok && len(addr) == 42 {
				if watchAddresses[strings.ToLower(addr)] {
					return true
				}
			}
		}
	}
	return false
}

// extractAllocKeys extracts the involved "agent:worknetId" key list from event data
func extractAllocKeys(eventType string, data map[string]interface{}) []string {
	var keys []string
	agent, _ := data["agent"].(string)
	subnetID, _ := data["worknetId"].(string)
	if agent != "" && subnetID != "" {
		keys = append(keys, strings.ToLower(agent)+":"+subnetID)
	}
	// Reallocated involves two (agent, subnet) pairs
	if eventType == "Reallocated" {
		fromAgent, _ := data["fromAgent"].(string)
		fromSubnet, _ := data["fromSubnet"].(string)
		toAgent, _ := data["toAgent"].(string)
		toSubnet, _ := data["toSubnet"].(string)
		if fromAgent != "" && fromSubnet != "" {
			keys = append(keys, strings.ToLower(fromAgent)+":"+fromSubnet)
		}
		if toAgent != "" && toSubnet != "" {
			keys = append(keys, strings.ToLower(toAgent)+":"+toSubnet)
		}
	}
	return keys
}

// enrichAllocEvent enriches an allocation event with current balance (called outside lock)
func (h *Hub) enrichAllocEvent(evt broadcastEvent, watchKey string) map[string]interface{} {
	result := map[string]interface{}{
		"type":        "AllocationChanged",
		"chainId":     evt.ChainID,
		"sourceEvent": evt.Type,
		"data":        evt.Data,
	}

	parts := strings.SplitN(watchKey, ":", 2)
	if len(parts) == 2 {
		agent, subnetID := parts[0], parts[1]
		result["agent"] = agent
		result["worknetId"] = subnetID

		// Use chainID from event; if 0, fall back to Hub default chainID
		cid := evt.ChainID
		h.mu.RLock()
		q := h.allocQuery
		if cid == 0 {
			cid = h.chainID
		}
		h.mu.RUnlock()

		if q != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if currentStake, err := q.GetAgentSubnetStakeWS(ctx, cid, agent, subnetID); err == nil {
				result["currentStake"] = currentStake
			}
		}
	}

	return result
}

// HandleConnect handles the WebSocket connection upgrade request
func (h *Hub) HandleConnect(w http.ResponseWriter, r *http.Request) {
	// Per-IP connection limit
	ip := r.RemoteAddr
	if host, _, err := net.SplitHostPort(ip); err == nil {
		ip = host
	}
	h.ipMu.Lock()
	maxConns := 10 // default
	if val, err := h.rdb.HGet(r.Context(), "ratelimit:config", "ws_connect").Result(); err == nil {
		if parts := strings.SplitN(val, ":", 2); len(parts) >= 1 {
			if n, err := strconv.Atoi(parts[0]); err == nil && n > 0 {
				maxConns = n
			}
		}
	}
	if h.connsByIP[ip] >= maxConns {
		h.ipMu.Unlock()
		http.Error(w, "too many WebSocket connections", http.StatusTooManyRequests)
		return
	}
	h.connsByIP[ip]++
	h.ipMu.Unlock()

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		h.logger.Error("WebSocket upgrade failed", "error", err)
		// Roll back the IP counter since the connection was not established
		h.ipMu.Lock()
		h.connsByIP[ip]--
		if h.connsByIP[ip] <= 0 {
			delete(h.connsByIP, ip)
		}
		h.ipMu.Unlock()
		return
	}

	// Limit incoming message size to 4KB
	conn.SetReadLimit(4096)

	client := &Client{
		hub:            h,
		conn:           conn,
		send:           make(chan []byte, 256),
		filters:        make(map[string]bool),
		allocWatches:   make(map[string]bool),
		subnetWatches:  make(map[string]bool),
		watchAddresses: make(map[string]bool),
		ip:             ip,
	}

	h.register <- client

	// Start read and write goroutines
	// Use context.Background() instead of r.Context() because the HTTP handler context
	// is cancelled when the handler returns, but the WebSocket connection outlives it.
	// The readPump/writePump goroutines manage their own lifecycle via connection close.
	go h.readPump(context.Background(), client)
	go h.writePump(context.Background(), client)
}

// readPump reads messages from the client and handles subscription filter updates
func (h *Hub) readPump(ctx context.Context, client *Client) {
	defer func() {
		h.unregister <- client
		_ = client.conn.Close(websocket.StatusNormalClosure, "connection closed")
	}()

	for {
		_, data, err := client.conn.Read(ctx)
		if err != nil {
			// Client disconnected or read error
			h.logger.Debug("failed to read client message", "error", err)
			return
		}

		// Parse the filter subscription message sent by the client
		var fm filterMessage
		if err := json.Unmarshal(data, &fm); err != nil {
			h.logger.Debug("failed to parse client message", "error", err)
			continue
		}

		if len(fm.Subscribe) > 0 && len(fm.Subscribe) <= 50 {
			// Update the client's event filters (max 50 event types)
			h.mu.Lock()
			client.filters = make(map[string]bool, len(fm.Subscribe))
			for _, eventType := range fm.Subscribe {
				if len(eventType) <= 64 { // max 64 chars per event type name
					client.filters[eventType] = true
				}
			}
			h.mu.Unlock()
			h.logger.Debug("client updated filters", "filters", fm.Subscribe)
		}

		// Handle address-level filter subscription
		if len(fm.WatchAddresses) > 0 && len(fm.WatchAddresses) <= 50 {
			h.mu.Lock()
			client.watchAddresses = make(map[string]bool, len(fm.WatchAddresses))
			for _, addr := range fm.WatchAddresses {
				if len(addr) == 42 && addr[:2] == "0x" {
					client.watchAddresses[strings.ToLower(addr)] = true
				}
			}
			h.mu.Unlock()
			h.logger.Debug("client updated address watches", "count", len(client.watchAddresses))
		}

		// Handle allocation watch subscription
		if len(fm.WatchAllocations) > 0 && len(fm.WatchAllocations) <= 100 {
			h.mu.Lock()
			client.allocWatches = make(map[string]bool)
			client.subnetWatches = make(map[string]bool)
			for _, w := range fm.WatchAllocations {
				if w.SubnetID == "" || len(w.SubnetID) > 30 {
					continue
				}
				if w.Agent == "" {
					// Omitting agent = subscribe to entire subnet
					client.subnetWatches[w.SubnetID] = true
				} else if len(w.Agent) == 42 {
					key := strings.ToLower(w.Agent) + ":" + w.SubnetID
					client.allocWatches[key] = true
				}
			}
			h.mu.Unlock()
			h.logger.Debug("client updated allocation watches",
				"exact", len(client.allocWatches), "subnet", len(client.subnetWatches))
		}
	}
}

// writePump writes messages from the send channel into the WebSocket connection
func (h *Hub) writePump(ctx context.Context, client *Client) {
	defer func() {
		_ = client.conn.Close(websocket.StatusNormalClosure, "connection closed")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-client.send:
			if !ok {
				// Send channel closed
				return
			}
			if err := client.conn.Write(ctx, websocket.MessageText, msg); err != nil {
				h.logger.Debug("failed to write client message", "error", err)
				return
			}
		}
	}
}
