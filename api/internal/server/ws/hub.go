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

// Client represents a WebSocket client connection
type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	send    chan []byte
	filters map[string]bool // event type filters
	ip      string          // client IP for connection limit tracking
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
	connsByIP map[string]int // per-IP WebSocket connection counter
	ipMu      sync.Mutex    // guards connsByIP
}

// filterMessage represents a subscription filter message sent by a client
type filterMessage struct {
	Subscribe []string `json:"subscribe"`
}

// broadcastEvent is used to parse the event type from a broadcast message
type broadcastEvent struct {
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

// broadcastToClients sends a message to all connected clients (respecting filter settings)
func (h *Hub) broadcastToClients(msg []byte) {
	// Parse event type for filtering
	var evt broadcastEvent
	eventType := ""
	if err := json.Unmarshal(msg, &evt); err == nil {
		eventType = evt.Type
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		// If the client has filters set, check whether the event type matches
		if len(client.filters) > 0 && eventType != "" {
			if !client.filters[eventType] {
				continue
			}
		}

		select {
		case client.send <- msg:
		default:
			// Send buffer is full; close this client
			delete(h.clients, client)
			close(client.send)
		}
	}
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
		hub:     h,
		conn:    conn,
		send:    make(chan []byte, 256),
		filters: make(map[string]bool),
		ip:      ip,
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

		if len(fm.Subscribe) > 0 {
			// Update the client's event filters
			h.mu.Lock()
			client.filters = make(map[string]bool, len(fm.Subscribe))
			for _, eventType := range fm.Subscribe {
				client.filters[eventType] = true
			}
			h.mu.Unlock()
			h.logger.Debug("client updated filters", "filters", fm.Subscribe)
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
