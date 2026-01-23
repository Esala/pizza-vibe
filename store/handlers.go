package store

import (
	"bytes"
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

// CookRequest represents the request sent to the kitchen service.
type CookRequest struct {
	OrderID    uuid.UUID   `json:"orderId"`
	OrderItems []OrderItem `json:"orderItems"`
}

// Store manages pizza orders and provides HTTP handlers for the store service.
type Store struct {
	mu         sync.RWMutex
	orders     map[uuid.UUID]*Order
	events     map[uuid.UUID][]OrderEvent
	hub        *WebSocketHub
	kitchenURL string
	httpClient *http.Client
}

// NewStore creates a new Store instance with initialized order storage and WebSocket hub.
func NewStore() *Store {
	return &Store{
		orders:     make(map[uuid.UUID]*Order),
		events:     make(map[uuid.UUID][]OrderEvent),
		hub:        NewWebSocketHub(),
		kitchenURL: "http://kitchen:8081",
		httpClient: &http.Client{},
	}
}

// SetKitchenURL sets the URL of the kitchen service.
func (s *Store) SetKitchenURL(url string) {
	s.kitchenURL = url
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

	// Call kitchen service to cook the order
	go s.callKitchenService(order)

	// Return the created order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// callKitchenService sends a cook request to the kitchen service.
func (s *Store) callKitchenService(order *Order) {
	cookReq := CookRequest{
		OrderID:    order.OrderID,
		OrderItems: order.OrderItems,
	}

	body, err := json.Marshal(cookReq)
	if err != nil {
		return
	}

	resp, err := s.httpClient.Post(s.kitchenURL+"/cook", "application/json", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer resp.Body.Close()
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

	// Map DONE status from kitchen to COOKED
	status := event.Status
	if event.Source == "kitchen" && event.Status == "DONE" {
		status = "COOKED"
	}

	// Update the order status
	if !s.UpdateOrderStatus(event.OrderID, status) {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Track the event
	s.trackEvent(event)

	// Broadcast the update to WebSocket clients
	s.BroadcastOrderUpdate(OrderUpdate{
		OrderID: event.OrderID,
		Status:  status,
		Source:  event.Source,
	})

	w.WriteHeader(http.StatusOK)
}

// trackEvent stores an event in the order's event history.
func (s *Store) trackEvent(event OrderEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[event.OrderID] = append(s.events[event.OrderID], event)
}

// GetOrderEvents retrieves all events for a given order ID.
func (s *Store) GetOrderEvents(orderID uuid.UUID) []OrderEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.events[orderID]
}

// HandleGetOrders handles GET /orders requests to retrieve all orders.
func (s *Store) HandleGetOrders(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]*Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// HandleGetEvents handles GET /events requests to retrieve events for a specific order.
func (s *Store) HandleGetEvents(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("orderId")
	if orderIDStr == "" {
		http.Error(w, "orderId query parameter is required", http.StatusBadRequest)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid orderId format", http.StatusBadRequest)
		return
	}

	events := s.GetOrderEvents(orderID)
	if events == nil {
		events = []OrderEvent{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
