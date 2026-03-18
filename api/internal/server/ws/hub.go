package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/coder/websocket"

	"net/http"
)

// Client represents a WebSocket client connection
type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	send    chan []byte
	filters map[string]bool // event type filters
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
	}
}

// Run starts the Hub's main loop, subscribes to the Redis channel, and handles client register/unregister/broadcast
func (h *Hub) Run(ctx context.Context) {
	// Subscribe to the Redis chain_events channel
	pubsub := h.rdb.Subscribe(ctx, "chain_events")
	defer func() { _ = pubsub.Close() }()

	redisCh := pubsub.Channel()
	h.logger.Info("WebSocket Hub started, listening on chain_events channel")

	for {
		select {
		case <-ctx.Done():
			h.logger.Info("WebSocket Hub shutting down")
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
			h.logger.Debug("client disconnected", "total", len(h.clients))

		case msg := <-h.broadcast:
			h.broadcastToClients(msg)

		case redisMsg, ok := <-redisCh:
			if !ok {
				h.logger.Warn("Redis subscription channel closed")
				return
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
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		h.logger.Error("WebSocket upgrade failed", "error", err)
		return
	}

	client := &Client{
		hub:     h,
		conn:    conn,
		send:    make(chan []byte, 256),
		filters: make(map[string]bool),
	}

	h.register <- client

	// Start read and write goroutines
	go h.readPump(r.Context(), client)
	go h.writePump(r.Context(), client)
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
