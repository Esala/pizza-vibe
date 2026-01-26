package store

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// OrderUpdate represents an order status update sent to WebSocket clients.
type OrderUpdate struct {
	OrderID uuid.UUID `json:"orderId"`
	Status  string    `json:"status"`
	Source  string    `json:"source"`
}

// WebSocketHub manages WebSocket client connections and broadcasts messages.
type WebSocketHub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]bool
}

// NewWebSocketHub creates a new WebSocketHub instance.
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients: make(map[*websocket.Conn]bool),
	}
}

// AddClient registers a new WebSocket client connection.
func (h *WebSocketHub) AddClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = true
}

// RemoveClient unregisters a WebSocket client connection.
func (h *WebSocketHub) RemoveClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
}

// Broadcast sends a message to all connected WebSocket clients.
func (h *WebSocketHub) Broadcast(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for conn := range h.clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			slog.Error("websocket write error", "error", err)
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// HandleWebSocket handles WebSocket connection requests from frontend clients.
// It upgrades the HTTP connection to WebSocket and registers the client
// to receive order updates.
func (s *Store) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade error", "error", err)
		return
	}

	s.hub.AddClient(conn)

	// Keep connection open and handle disconnection
	go func() {
		defer func() {
			s.hub.RemoveClient(conn)
			conn.Close()
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

// BroadcastOrderUpdate sends an order update to all connected WebSocket clients.
func (s *Store) BroadcastOrderUpdate(update OrderUpdate) {
	message, err := json.Marshal(update)
	if err != nil {
		slog.Error("failed to marshal order update", "error", err)
		return
	}
	s.hub.Broadcast(message)
}
