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

// allocationWatch 订阅特定 (agent, subnetId) 分配变动
type allocationWatch struct {
	Agent    string `json:"agent"`
	SubnetID string `json:"subnetId"`
}

// Client represents a WebSocket client connection
type Client struct {
	hub              *Hub
	conn             *websocket.Conn
	send             chan []byte
	filters          map[string]bool // event type filters
	allocWatches     map[string]bool // "agent:subnetId" 精确匹配
	subnetWatches    map[string]bool // "subnetId" 整子网匹配（agent 省略时）
	ip               string          // client IP for connection limit tracking
}

// AllocationQuerier 查询 (agent, subnetId) 当前分配量（由 handler 层注入）
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
}

// broadcastEvent is used to parse the event type and data from a broadcast message
type broadcastEvent struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data,omitempty"`
}

// broadcastEventType 仅解析 type 字段（用于不关心 data 的场景）
type broadcastEventType struct {
	Type string `json:"type"`
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

// SetAllocationQuerier 注入分配查询接口（用于推送时附带当前余额）
func (h *Hub) SetAllocationQuerier(q AllocationQuerier, chainID int64) {
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
				// Reconnect loop with backoff — keep existing client connections alive
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
					time.Sleep(3 * time.Second)
					pubsub = h.rdb.Subscribe(ctx, "chain_events")
					redisCh = pubsub.Channel()
					h.logger.Info("Redis Pub/Sub reconnected")
					break
				}
				continue // restart the main select loop with new redisCh
			}
			// Broadcast Redis message to all clients
			h.broadcastToClients([]byte(redisMsg.Payload))
		}
	}
}

// allocEventTypes 是触发分配推送的事件类型
var allocEventTypes = map[string]bool{
	"Allocated":   true,
	"Deallocated": true,
	"Reallocated": true,
}

// broadcastToClients sends a message to all connected clients (respecting filter settings)
func (h *Hub) broadcastToClients(msg []byte) {
	// 先解析 type（不关心 data 类型，避免类型不匹配导致 unmarshal 失败）
	var evtType broadcastEventType
	eventType := ""
	if err := json.Unmarshal(msg, &evtType); err == nil {
		eventType = evtType.Type
	}

	// 仅对分配事件才解析完整 data
	var evt broadcastEvent
	var allocKeys []string
	if allocEventTypes[eventType] {
		if err := json.Unmarshal(msg, &evt); err == nil && evt.Data != nil {
			allocKeys = extractAllocKeys(eventType, evt.Data)
		}
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// 提取涉及的 subnetId 列表（用于子网级匹配）
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

	for client := range h.clients {
		sent := false
		evicted := false
		hasWatches := len(client.allocWatches) > 0 || len(client.subnetWatches) > 0

		// 1. 检查分配订阅推送（精确匹配 + 子网级匹配）
		if hasWatches && len(allocKeys) > 0 {
			// 1a. 精确匹配 agent:subnetId
			for _, key := range allocKeys {
				if client.allocWatches[key] {
					enriched := h.enrichAllocEvent(evt, key)
					if data, err := json.Marshal(enriched); err == nil {
						select {
						case client.send <- data:
							sent = true
						default:
							delete(h.clients, client)
							close(client.send)
							evicted = true
						}
					}
					break
				}
			}
			// 1b. 子网级匹配（agent 省略）
			if !sent && !evicted {
				for _, sid := range allocSubnets {
					if client.subnetWatches[sid] {
						for _, key := range allocKeys {
							if strings.HasSuffix(key, ":"+sid) {
								enriched := h.enrichAllocEvent(evt, key)
								if data, err := json.Marshal(enriched); err == nil {
									select {
									case client.send <- data:
										sent = true
									default:
										delete(h.clients, client)
										close(client.send)
										evicted = true
									}
								}
								break
							}
						}
						break
					}
				}
			}
		}

		// 2. 常规 type 过滤推送（跳过已发送或已驱逐的客户端）
		if !sent && !evicted {
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
}

// extractAllocKeys 从事件 data 中提取涉及的 "agent:subnetId" key 列表
func extractAllocKeys(eventType string, data map[string]interface{}) []string {
	var keys []string
	agent, _ := data["agent"].(string)
	subnetID, _ := data["subnetId"].(string)
	if agent != "" && subnetID != "" {
		keys = append(keys, strings.ToLower(agent)+":"+subnetID)
	}
	// Reallocated 涉及两个 (agent, subnet) 对
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

// enrichAllocEvent 为分配事件附加当前余额
func (h *Hub) enrichAllocEvent(evt broadcastEvent, watchKey string) map[string]interface{} {
	result := map[string]interface{}{
		"type":        "AllocationChanged",
		"sourceEvent": evt.Type,
		"data":        evt.Data,
	}

	// 解析 agent 和 subnetId
	parts := strings.SplitN(watchKey, ":", 2)
	if len(parts) == 2 {
		agent, subnetID := parts[0], parts[1]
		result["agent"] = agent
		result["subnetId"] = subnetID
		// 查询当前余额（如果 querier 已注入）
		if h.allocQuery != nil {
			if currentStake, err := h.allocQuery.GetAgentSubnetStakeWS(
				context.Background(), h.chainID, agent, subnetID,
			); err == nil {
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
		hub:           h,
		conn:          conn,
		send:          make(chan []byte, 256),
		filters:       make(map[string]bool),
		allocWatches:  make(map[string]bool),
		subnetWatches: make(map[string]bool),
		ip:            ip,
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

		// 处理分配监听订阅
		if len(fm.WatchAllocations) > 0 && len(fm.WatchAllocations) <= 100 {
			h.mu.Lock()
			client.allocWatches = make(map[string]bool)
			client.subnetWatches = make(map[string]bool)
			for _, w := range fm.WatchAllocations {
				if w.SubnetID == "" || len(w.SubnetID) > 30 {
					continue
				}
				if w.Agent == "" {
					// 省略 agent = 订阅整个子网
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
