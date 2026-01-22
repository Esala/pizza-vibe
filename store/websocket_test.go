package store

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

// TestWebSocketConnection verifies that clients can connect via WebSocket.
func TestWebSocketConnection(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Get("/ws", store.HandleWebSocket)

	server := httptest.NewServer(router)
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Connection successful
	t.Log("WebSocket connection established successfully")
}

// TestWebSocketReceivesOrderUpdates verifies that WebSocket clients receive order updates.
func TestWebSocketReceivesOrderUpdates(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)
	router.Post("/events", store.HandleEvent)
	router.Get("/ws", store.HandleWebSocket)

	server := httptest.NewServer(router)
	defer server.Close()

	// Connect to WebSocket
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Create an order via HTTP
	orderReq := CreateOrderRequest{
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 1},
		},
		OrderData: "WebSocket test",
	}
	orderBody, _ := json.Marshal(orderReq)
	resp, err := http.Post(server.URL+"/order", "application/json", bytes.NewReader(orderBody))
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}
	defer resp.Body.Close()

	var createdOrder Order
	json.NewDecoder(resp.Body).Decode(&createdOrder)

	// Send an event to update the order
	eventReq := OrderEvent{
		OrderID: createdOrder.OrderID,
		Status:  "cooking",
		Source:  "kitchen",
	}
	eventBody, _ := json.Marshal(eventReq)
	resp2, err := http.Post(server.URL+"/events", "application/json", bytes.NewReader(eventBody))
	if err != nil {
		t.Fatalf("failed to send event: %v", err)
	}
	defer resp2.Body.Close()

	// Read the update from WebSocket with a timeout
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read WebSocket message: %v", err)
	}

	// Parse the received order update
	var update OrderUpdate
	if err := json.Unmarshal(message, &update); err != nil {
		t.Fatalf("failed to unmarshal WebSocket message: %v", err)
	}

	if update.OrderID != createdOrder.OrderID {
		t.Errorf("expected OrderID %v, got %v", createdOrder.OrderID, update.OrderID)
	}
	if update.Status != "cooking" {
		t.Errorf("expected status 'cooking', got '%s'", update.Status)
	}
}

// TestMultipleWebSocketClients verifies that multiple clients can receive updates.
func TestMultipleWebSocketClients(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)
	router.Post("/events", store.HandleEvent)
	router.Get("/ws", store.HandleWebSocket)

	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect two clients
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect client 1: %v", err)
	}
	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect client 2: %v", err)
	}
	defer conn2.Close()

	// Create an order
	orderReq := CreateOrderRequest{
		OrderItems: []OrderItem{{PizzaType: "Pepperoni", Quantity: 1}},
		OrderData:  "Multi-client test",
	}
	orderBody, _ := json.Marshal(orderReq)
	resp, _ := http.Post(server.URL+"/order", "application/json", bytes.NewReader(orderBody))
	var createdOrder Order
	json.NewDecoder(resp.Body).Decode(&createdOrder)
	resp.Body.Close()

	// Send an event
	eventReq := OrderEvent{
		OrderID: createdOrder.OrderID,
		Status:  "ready",
		Source:  "kitchen",
	}
	eventBody, _ := json.Marshal(eventReq)
	resp2, _ := http.Post(server.URL+"/events", "application/json", bytes.NewReader(eventBody))
	resp2.Body.Close()

	// Both clients should receive the update
	conn1.SetReadDeadline(time.Now().Add(2 * time.Second))
	conn2.SetReadDeadline(time.Now().Add(2 * time.Second))

	_, msg1, err := conn1.ReadMessage()
	if err != nil {
		t.Fatalf("client 1 failed to read: %v", err)
	}

	_, msg2, err := conn2.ReadMessage()
	if err != nil {
		t.Fatalf("client 2 failed to read: %v", err)
	}

	var update1, update2 OrderUpdate
	json.Unmarshal(msg1, &update1)
	json.Unmarshal(msg2, &update2)

	if update1.Status != "ready" || update2.Status != "ready" {
		t.Error("both clients should receive 'ready' status")
	}
}
