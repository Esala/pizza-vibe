package store

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TestPostOrder verifies that POST /order creates a new order and returns it with a UUID.
func TestPostOrder(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)

	// Create order request payload
	reqBody := CreateOrderRequest{
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 2},
			{PizzaType: "Pepperoni", Quantity: 1},
		},
		OrderData: "Ring the doorbell",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", rec.Code)
	}

	var response Order
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Verify OrderID is a valid UUID
	if response.OrderID == uuid.Nil {
		t.Error("expected valid OrderID, got nil UUID")
	}

	// Verify order items
	if len(response.OrderItems) != 2 {
		t.Errorf("expected 2 order items, got %d", len(response.OrderItems))
	}

	// Verify order status is set to "pending"
	if response.OrderStatus != "pending" {
		t.Errorf("expected OrderStatus 'pending', got '%s'", response.OrderStatus)
	}

	// Verify order data
	if response.OrderData != "Ring the doorbell" {
		t.Errorf("expected OrderData 'Ring the doorbell', got '%s'", response.OrderData)
	}
}

// TestPostOrderInvalidJSON verifies that POST /order returns 400 for invalid JSON.
func TestPostOrderInvalidJSON(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", rec.Code)
	}
}

// TestPostOrderEmptyItems verifies that POST /order returns 400 when no items provided.
func TestPostOrderEmptyItems(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)

	reqBody := CreateOrderRequest{
		OrderItems: []OrderItem{},
		OrderData:  "Empty order",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", rec.Code)
	}
}

// TestPostEvents verifies that POST /events updates order status from kitchen/delivery events.
func TestPostEvents(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)
	router.Post("/events", store.HandleEvent)

	// First, create an order
	orderReq := CreateOrderRequest{
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 1},
		},
		OrderData: "Test order",
	}
	orderBody, _ := json.Marshal(orderReq)
	createReq := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(orderBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	var createdOrder Order
	json.Unmarshal(createRec.Body.Bytes(), &createdOrder)

	// Now send an event to update the order status
	eventReq := OrderEvent{
		OrderID: createdOrder.OrderID,
		Status:  "cooking",
		Source:  "kitchen",
	}
	eventBody, _ := json.Marshal(eventReq)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(eventBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", rec.Code)
	}

	// Verify the order status was updated
	order, exists := store.GetOrder(createdOrder.OrderID)
	if !exists {
		t.Fatal("order not found")
	}
	if order.OrderStatus != "cooking" {
		t.Errorf("expected OrderStatus 'cooking', got '%s'", order.OrderStatus)
	}
}

// TestPostEventsInvalidOrderID verifies that POST /events returns 404 for non-existent order.
func TestPostEventsInvalidOrderID(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/events", store.HandleEvent)

	eventReq := OrderEvent{
		OrderID: uuid.New(), // Non-existent order
		Status:  "cooking",
		Source:  "kitchen",
	}
	eventBody, _ := json.Marshal(eventReq)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(eventBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404 Not Found, got %d", rec.Code)
	}
}

// TestPostEventsInvalidJSON verifies that POST /events returns 400 for invalid JSON.
func TestPostEventsInvalidJSON(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/events", store.HandleEvent)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", rec.Code)
	}
}

// TestPostOrderCallsKitchenService verifies that POST /order calls the kitchen service to cook the order.
func TestPostOrderCallsKitchenService(t *testing.T) {
	// Create a mock kitchen server
	var receivedRequest CookRequest
	kitchenCalled := make(chan bool, 1)
	kitchenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedRequest)
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"status": "cooking"})
		kitchenCalled <- true
	}))
	defer kitchenServer.Close()

	store := NewStore()
	store.SetKitchenURL(kitchenServer.URL)

	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)

	reqBody := CreateOrderRequest{
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 2},
		},
		OrderData: "Test order",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", rec.Code)
	}

	// Wait for the kitchen to be called (with timeout)
	select {
	case <-kitchenCalled:
		// Kitchen was called
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for kitchen service to be called")
	}

	// Verify the kitchen received the correct order data
	var response Order
	json.Unmarshal(rec.Body.Bytes(), &response)

	if receivedRequest.OrderID != response.OrderID {
		t.Errorf("expected kitchen to receive orderId %s, got %s", response.OrderID, receivedRequest.OrderID)
	}

	if len(receivedRequest.OrderItems) != 1 {
		t.Errorf("expected kitchen to receive 1 order item, got %d", len(receivedRequest.OrderItems))
	}

	if receivedRequest.OrderItems[0].PizzaType != "Margherita" {
		t.Errorf("expected kitchen to receive Margherita pizza, got %s", receivedRequest.OrderItems[0].PizzaType)
	}
}

// TestEventsAreTrackedPerOrderID verifies that events are tracked per order ID.
func TestEventsAreTrackedPerOrderID(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)
	router.Post("/events", store.HandleEvent)

	// Create an order
	orderReq := CreateOrderRequest{
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 1},
		},
		OrderData: "Test order",
	}
	orderBody, _ := json.Marshal(orderReq)
	createReq := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(orderBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	var createdOrder Order
	json.Unmarshal(createRec.Body.Bytes(), &createdOrder)

	// Send multiple events
	events := []OrderEvent{
		{OrderID: createdOrder.OrderID, Status: "cooking", Source: "kitchen"},
		{OrderID: createdOrder.OrderID, Status: "preparing pizza", Source: "kitchen"},
		{OrderID: createdOrder.OrderID, Status: "in oven", Source: "kitchen"},
	}

	for _, event := range events {
		eventBody, _ := json.Marshal(event)
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(eventBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}

	// Verify all events are tracked
	trackedEvents := store.GetOrderEvents(createdOrder.OrderID)
	if len(trackedEvents) != 3 {
		t.Errorf("expected 3 tracked events, got %d", len(trackedEvents))
	}

	// Verify event order
	if trackedEvents[0].Status != "cooking" {
		t.Errorf("expected first event status 'cooking', got '%s'", trackedEvents[0].Status)
	}
	if trackedEvents[1].Status != "preparing pizza" {
		t.Errorf("expected second event status 'preparing pizza', got '%s'", trackedEvents[1].Status)
	}
	if trackedEvents[2].Status != "in oven" {
		t.Errorf("expected third event status 'in oven', got '%s'", trackedEvents[2].Status)
	}
}

// TestDoneEventUpdatesStatusToCooked verifies that a DONE event updates order status to COOKED.
func TestDoneEventUpdatesStatusToCooked(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)
	router.Post("/events", store.HandleEvent)

	// Create an order
	orderReq := CreateOrderRequest{
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 1},
		},
		OrderData: "Test order",
	}
	orderBody, _ := json.Marshal(orderReq)
	createReq := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(orderBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	var createdOrder Order
	json.Unmarshal(createRec.Body.Bytes(), &createdOrder)

	// Send DONE event
	doneEvent := OrderEvent{
		OrderID: createdOrder.OrderID,
		Status:  "DONE",
		Source:  "kitchen",
	}
	eventBody, _ := json.Marshal(doneEvent)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(eventBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", rec.Code)
	}

	// Verify the order status is now COOKED
	order, exists := store.GetOrder(createdOrder.OrderID)
	if !exists {
		t.Fatal("order not found")
	}
	if order.OrderStatus != "COOKED" {
		t.Errorf("expected OrderStatus 'COOKED', got '%s'", order.OrderStatus)
	}
}

// TestGetOrders verifies that GET /orders returns all orders.
func TestGetOrders(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)
	router.Get("/orders", store.HandleGetOrders)

	// Create two orders
	orderReq1 := CreateOrderRequest{
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 2},
		},
		OrderData: "Order 1",
	}
	orderReq2 := CreateOrderRequest{
		OrderItems: []OrderItem{
			{PizzaType: "Pepperoni", Quantity: 1},
		},
		OrderData: "Order 2",
	}

	body1, _ := json.Marshal(orderReq1)
	req1 := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)

	body2, _ := json.Marshal(orderReq2)
	req2 := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	// GET /orders
	getReq := httptest.NewRequest(http.MethodGet, "/orders", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", getRec.Code)
	}

	var orders []Order
	err := json.Unmarshal(getRec.Body.Bytes(), &orders)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(orders))
	}
}

// TestGetOrdersEmpty verifies that GET /orders returns empty array when no orders exist.
func TestGetOrdersEmpty(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Get("/orders", store.HandleGetOrders)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", rec.Code)
	}

	var orders []Order
	err := json.Unmarshal(rec.Body.Bytes(), &orders)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(orders) != 0 {
		t.Errorf("expected 0 orders, got %d", len(orders))
	}
}

// TestGetEventsForOrder verifies that GET /events returns events for a specific order.
func TestGetEventsForOrder(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)
	router.Post("/events", store.HandleEvent)
	router.Get("/events", store.HandleGetEvents)

	// Create an order
	orderReq := CreateOrderRequest{
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 1},
		},
		OrderData: "Test order",
	}
	orderBody, _ := json.Marshal(orderReq)
	createReq := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(orderBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	var createdOrder Order
	json.Unmarshal(createRec.Body.Bytes(), &createdOrder)

	// Send multiple events
	events := []OrderEvent{
		{OrderID: createdOrder.OrderID, Status: "cooking", Source: "kitchen"},
		{OrderID: createdOrder.OrderID, Status: "in oven", Source: "kitchen"},
		{OrderID: createdOrder.OrderID, Status: "DONE", Source: "kitchen"},
	}

	for _, event := range events {
		eventBody, _ := json.Marshal(event)
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(eventBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}

	// GET /events?orderId={orderId}
	getReq := httptest.NewRequest(http.MethodGet, "/events?orderId="+createdOrder.OrderID.String(), nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", getRec.Code)
	}

	var returnedEvents []OrderEvent
	err := json.Unmarshal(getRec.Body.Bytes(), &returnedEvents)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(returnedEvents) != 3 {
		t.Errorf("expected 3 events, got %d", len(returnedEvents))
	}

	if returnedEvents[0].Status != "cooking" {
		t.Errorf("expected first event status 'cooking', got '%s'", returnedEvents[0].Status)
	}
}

// TestGetEventsWithoutOrderId verifies that GET /events returns 400 when orderId is missing.
func TestGetEventsWithoutOrderId(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Get("/events", store.HandleGetEvents)

	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", rec.Code)
	}
}

// TestGetEventsInvalidOrderId verifies that GET /events returns 400 for invalid UUID.
func TestGetEventsInvalidOrderId(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Get("/events", store.HandleGetEvents)

	req := httptest.NewRequest(http.MethodGet, "/events?orderId=invalid-uuid", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", rec.Code)
	}
}

// TestGetEventsEmptyForOrder verifies that GET /events returns empty array for order without events.
func TestGetEventsEmptyForOrder(t *testing.T) {
	store := NewStore()
	router := chi.NewRouter()
	router.Post("/order", store.HandleCreateOrder)
	router.Get("/events", store.HandleGetEvents)

	// Create an order but don't send any events
	orderReq := CreateOrderRequest{
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 1},
		},
		OrderData: "Test order",
	}
	orderBody, _ := json.Marshal(orderReq)
	createReq := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(orderBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	var createdOrder Order
	json.Unmarshal(createRec.Body.Bytes(), &createdOrder)

	// GET /events?orderId={orderId}
	getReq := httptest.NewRequest(http.MethodGet, "/events?orderId="+createdOrder.OrderID.String(), nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", getRec.Code)
	}

	var returnedEvents []OrderEvent
	err := json.Unmarshal(getRec.Body.Bytes(), &returnedEvents)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(returnedEvents) != 0 {
		t.Errorf("expected 0 events, got %d", len(returnedEvents))
	}
}
