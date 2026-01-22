package store

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

// CreateOrderRequest represents the request body for creating a new order.
type CreateOrderRequest struct {
	OrderItems []OrderItem `json:"orderItems"`
	OrderData  string      `json:"orderData"`
}

// Store manages pizza orders and provides HTTP handlers for the store service.
type Store struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*Order
	hub    *WebSocketHub
}

// NewStore creates a new Store instance with initialized order storage and WebSocket hub.
func NewStore() *Store {
	return &Store{
		orders: make(map[uuid.UUID]*Order),
		hub:    NewWebSocketHub(),
	}
}

// HandleCreateOrder handles POST /order requests to create new pizza orders.
// It validates the request, generates a UUID for the order, and stores it.
func (s *Store) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate that at least one item is provided
	if len(req.OrderItems) == 0 {
		http.Error(w, "Order must contain at least one item", http.StatusBadRequest)
		return
	}

	// Create new order with generated UUID
	order := &Order{
		OrderID:     uuid.New(),
		OrderItems:  req.OrderItems,
		OrderData:   req.OrderData,
		OrderStatus: "pending",
	}

	// Store the order
	s.mu.Lock()
	s.orders[order.OrderID] = order
	s.mu.Unlock()

	// Return the created order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetOrder retrieves an order by its UUID.
func (s *Store) GetOrder(orderID uuid.UUID) (*Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, exists := s.orders[orderID]
	return order, exists
}

// UpdateOrderStatus updates the status of an existing order.
func (s *Store) UpdateOrderStatus(orderID uuid.UUID, status string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, exists := s.orders[orderID]
	if !exists {
		return false
	}
	order.OrderStatus = status
	return true
}

// OrderEvent represents an event received from kitchen or delivery services.
type OrderEvent struct {
	OrderID uuid.UUID `json:"orderId"`
	Status  string    `json:"status"`
	Source  string    `json:"source"` // "kitchen" or "delivery"
}

// HandleEvent handles POST /events requests to receive order updates
// from kitchen and delivery services. It updates the order status and
// broadcasts the update to all connected WebSocket clients.
func (s *Store) HandleEvent(w http.ResponseWriter, r *http.Request) {
	var event OrderEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update the order status
	if !s.UpdateOrderStatus(event.OrderID, event.Status) {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Broadcast the update to WebSocket clients
	s.BroadcastOrderUpdate(OrderUpdate{
		OrderID: event.OrderID,
		Status:  event.Status,
		Source:  event.Source,
	})

	w.WriteHeader(http.StatusOK)
}
