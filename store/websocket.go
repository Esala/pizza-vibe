package store

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// OrderUpdate represents an order status update sent to WebSocket clients.
type OrderUpdate struct {
	OrderID uuid.UUID `json:"orderId"`
	Status  string    `json:"status"`
	Source  string    `json:"source"`
}

// WebSocketEvent represents the event format sent to frontend clients via WebSocket.
type WebSocketEvent struct {
	OrderID   uuid.UUID `json:"orderId"`
	Status    string    `json:"status"`
	Source    string    `json:"source"`
	Timestamp string    `json:"timestamp"`
}

// WebSocketHub manages WebSocket client connections and broadcasts messages.
type WebSocketHub struct {
	mu      sync.RWMutex
	clients map[string]*websocket.Conn
}

// NewWebSocketHub creates a new WebSocketHub instance.
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients: make(map[string]*websocket.Conn),
	}
}

// AddClient registers a new WebSocket client connection with a client ID.
func (h *WebSocketHub) AddClient(clientID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[clientID] = conn
}

// RemoveClient unregisters a WebSocket client connection by client ID.
func (h *WebSocketHub) RemoveClient(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, clientID)
}

// HasClient checks if a client with the given ID is registered.
func (h *WebSocketHub) HasClient(clientID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, exists := h.clients[clientID]
	return exists
}

// Broadcast sends a message to all connected WebSocket clients.
func (h *WebSocketHub) Broadcast(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for clientID, conn := range h.clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			slog.Error("websocket write error", "clientId", clientID, "error", err)
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
// to receive order updates. Requires a clientId query parameter.
func (s *Store) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("clientId")
	if clientID == "" {
		http.Error(w, "clientId query parameter is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade error", "error", err)
		return
	}

	s.hub.AddClient(clientID, conn)
	slog.Info("websocket client connected", "clientId", clientID)

	// Keep connection open and handle disconnection
	go func() {
		defer func() {
			s.hub.RemoveClient(clientID)
			conn.Close()
			slog.Info("websocket client disconnected", "clientId", clientID)
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

// BroadcastOrderUpdate sends an order update to all connected WebSocket clients
// using the WebSocketEvent format.
func (s *Store) BroadcastOrderUpdate(update OrderUpdate) {
	event := WebSocketEvent{
		OrderID:   update.OrderID,
		Status:    update.Status,
		Source:    update.Source,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	message, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal websocket event", "error", err)
		return
	}
	s.hub.Broadcast(message)
}
