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

	// Convert http URL to ws URL (with clientId)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?clientId=test-client"

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

	// Connect to WebSocket (with clientId)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?clientId=update-client"
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

	// Parse the received WebSocket event
	var event WebSocketEvent
	if err := json.Unmarshal(message, &event); err != nil {
		t.Fatalf("failed to unmarshal WebSocket message: %v", err)
	}

	if event.OrderID != createdOrder.OrderID {
		t.Errorf("expected OrderID %v, got %v", createdOrder.OrderID, event.OrderID)
	}
	if event.Status != "cooking" {
		t.Errorf("expected status 'cooking', got '%s'", event.Status)
	}
	if event.Source != "kitchen" {
		t.Errorf("expected source 'kitchen', got '%s'", event.Source)
	}
	if event.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

// TestWebSocketConnectionWithClientID verifies that clients can connect with a client ID.
func TestWebSocketConnectionWithClientID(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Get("/ws", store.HandleWebSocket)

	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?clientId=client-123"

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect to WebSocket with clientId: %v", err)
	}
	defer conn.Close()

	// Verify the client is registered with its ID
	if !store.hub.HasClient("client-123") {
		t.Error("expected client 'client-123' to be registered in the hub")
	}
}

// TestWebSocketConnectionWithoutClientIDIsRejected verifies connections without clientId are rejected.
func TestWebSocketConnectionWithoutClientIDIsRejected(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Get("/ws", store.HandleWebSocket)

	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	_, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err == nil {
		t.Fatal("expected connection to fail without clientId")
	}
	if resp != nil && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

// TestWebSocketBroadcastToSpecificClient verifies updates are sent to the correct client.
func TestWebSocketBroadcastToSpecificClient(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)
	router.Post("/events", store.HandleEvent)
	router.Get("/ws", store.HandleWebSocket)

	server := httptest.NewServer(router)
	defer server.Close()

	// Connect client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?clientId=client-abc"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Create an order
	orderReq := CreateOrderRequest{
		OrderItems: []OrderItem{{PizzaType: "Margherita", Quantity: 1}},
		OrderData:  "Client ID test",
	}
	orderBody, _ := json.Marshal(orderReq)
	resp, err := http.Post(server.URL+"/order", "application/json", bytes.NewReader(orderBody))
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}
	defer resp.Body.Close()

	var createdOrder Order
	json.NewDecoder(resp.Body).Decode(&createdOrder)

	// Send an event
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

	// Read the WebSocket event
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read WebSocket message: %v", err)
	}

	var event WebSocketEvent
	if err := json.Unmarshal(message, &event); err != nil {
		t.Fatalf("failed to unmarshal WebSocket event: %v", err)
	}

	if event.OrderID != createdOrder.OrderID {
		t.Errorf("expected OrderID %v, got %v", createdOrder.OrderID, event.OrderID)
	}
	if event.Status != "cooking" {
		t.Errorf("expected status 'cooking', got '%s'", event.Status)
	}
	if event.Source != "kitchen" {
		t.Errorf("expected source 'kitchen', got '%s'", event.Source)
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

	// Connect two clients (with clientIds)
	wsURL1 := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?clientId=multi-client-1"
	wsURL2 := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?clientId=multi-client-2"

	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	if err != nil {
		t.Fatalf("failed to connect client 1: %v", err)
	}
	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
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

	var event1, event2 WebSocketEvent
	json.Unmarshal(msg1, &event1)
	json.Unmarshal(msg2, &event2)

	if event1.Status != "ready" || event2.Status != "ready" {
		t.Error("both clients should receive 'ready' status")
	}
}
