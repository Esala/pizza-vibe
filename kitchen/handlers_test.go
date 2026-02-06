// Package kitchen provides the kitchen service for the Pizza Vibe application.
// It handles cooking pizza orders by processing order items.
package kitchen

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

// TestCookEndpointReturnsAccepted tests that the /cook endpoint accepts a valid cook request.
func TestCookEndpointReturnsAccepted(t *testing.T) {
	kitchen := NewKitchen()
	router := chi.NewRouter()
	router.Post("/cook", kitchen.HandleCook)

	orderID := uuid.New()
	req := CookRequest{
		OrderID: orderID,
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 1},
		},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/cook", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rr.Code)
	}
}

// TestCookEndpointInvalidJSON tests that the /cook endpoint returns bad request for invalid JSON.
func TestCookEndpointInvalidJSON(t *testing.T) {
	kitchen := NewKitchen()
	router := chi.NewRouter()
	router.Post("/cook", kitchen.HandleCook)

	httpReq := httptest.NewRequest(http.MethodPost, "/cook", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

// TestCookEndpointEmptyOrderItems tests that the /cook endpoint returns bad request for empty order items.
func TestCookEndpointEmptyOrderItems(t *testing.T) {
	kitchen := NewKitchen()
	router := chi.NewRouter()
	router.Post("/cook", kitchen.HandleCook)

	orderID := uuid.New()
	req := CookRequest{
		OrderID:    orderID,
		OrderItems: []OrderItem{},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/cook", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

// TestCookEndpointReturnsResponse tests that the /cook endpoint returns a proper response body.
func TestCookEndpointReturnsResponse(t *testing.T) {
	kitchen := NewKitchen()
	router := chi.NewRouter()
	router.Post("/cook", kitchen.HandleCook)

	orderID := uuid.New()
	req := CookRequest{
		OrderID: orderID,
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 2},
			{PizzaType: "Pepperoni", Quantity: 1},
		},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/cook", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rr.Code)
	}

	var resp CookResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.OrderID != orderID {
		t.Errorf("expected orderID %s, got %s", orderID, resp.OrderID)
	}

	if resp.Status != "cooking" {
		t.Errorf("expected status 'cooking', got %s", resp.Status)
	}
}

// TestCookCallsCookingAgentForEachPizza tests that cooking calls the cooking-agent for each pizza.
func TestCookCallsCookingAgentForEachPizza(t *testing.T) {
	// Track requests received by the mock cooking-agent
	agentRequests := make(chan AgentCookRequest, 10)
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cook/stream" && r.Method == http.MethodPost {
			var req AgentCookRequest
			json.NewDecoder(r.Body).Decode(&req)
			agentRequests <- req

			// Send SSE response with cooking updates
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.WriteHeader(http.StatusOK)

			flusher, ok := w.(http.Flusher)
			if !ok {
				t.Error("expected http.Flusher")
				return
			}

			// Send a few updates
			updates := []CookingUpdate{
				{Type: "action", Action: "checking_inventory", Message: "Checking available ingredients"},
				{Type: "action", Action: "acquiring_ingredients", Message: "Acquiring mozzarella"},
				{Type: "action", Action: "reserving_oven", Message: "Reserving oven for cooking"},
				{Type: "result", Action: "completed", Message: "Pizza cooked successfully"},
			}

			for _, update := range updates {
				data, _ := json.Marshal(update)
				w.Write([]byte("data: " + string(data) + "\n\n"))
				flusher.Flush()
			}
		}
	}))
	defer agentServer.Close()

	// Track events received by the mock store
	eventsReceived := make(chan OrderEvent, 50)
	storeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/events" {
			var event OrderEvent
			json.NewDecoder(r.Body).Decode(&event)
			eventsReceived <- event
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer storeServer.Close()

	kitchen := NewKitchenWithConfig(KitchenConfig{
		StoreURL:        storeServer.URL,
		CookingAgentURL: agentServer.URL,
	})

	router := chi.NewRouter()
	router.Post("/cook", kitchen.HandleCook)

	orderID := uuid.New()
	req := CookRequest{
		OrderID: orderID,
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 2},
			{PizzaType: "Pepperoni", Quantity: 1},
		},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/cook", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpReq)

	// Wait for agent requests (should be 3: 2 Margherita + 1 Pepperoni)
	var requests []AgentCookRequest
	timeout := time.After(5 * time.Second)
	for i := 0; i < 3; i++ {
		select {
		case req := <-agentRequests:
			requests = append(requests, req)
		case <-timeout:
			t.Fatalf("timed out waiting for agent requests, got %d", len(requests))
		}
	}

	if len(requests) != 3 {
		t.Errorf("expected 3 agent requests, got %d", len(requests))
	}

	// Wait for events - we expect multiple streaming updates per pizza + DONE
	// Each pizza has 4 updates, 3 pizzas = 12 updates + 3 "cooked" events + 1 DONE = minimum 4 per pizza
	var events []OrderEvent
	for {
		select {
		case event := <-eventsReceived:
			events = append(events, event)
			if event.Status == "DONE" {
				// Verify we got events for the cooking process
				if len(events) < 4 { // At least 3 cooked events + DONE
					t.Errorf("expected at least 4 events, got %d", len(events))
				}
				return
			}
		case <-timeout:
			t.Fatalf("timed out waiting for events, got %d events", len(events))
		}
	}
}

// TestParseSSEStreamWithoutSpace tests parsing SSE events in "data:value" format (no space after colon),
// which is the format Quarkus RESTEasy Reactive produces.
func TestParseSSEStreamWithoutSpace(t *testing.T) {
	kitchen := NewKitchen()
	updates := make(chan CookingUpdate, 10)

	// Simulate Quarkus SSE format: "data:{json}\n\n" (no space after data:)
	sseData := `data:{"type":"action","action":"checking_inventory","message":"Checking ingredients"}

data:{"type":"action","action":"reserving_oven","message":"Reserving oven"}

data:{"type":"result","action":"completed","message":"Done cooking"}

`
	body := bytes.NewReader([]byte(sseData))

	go func() {
		defer close(updates)
		err := kitchen.parseSSEStream(body, updates)
		if err != nil {
			t.Errorf("parseSSEStream returned error: %v", err)
		}
	}()

	var received []CookingUpdate
	timeout := time.After(2 * time.Second)
	for {
		select {
		case update, ok := <-updates:
			if !ok {
				goto done
			}
			received = append(received, update)
		case <-timeout:
			t.Fatal("timed out waiting for updates")
		}
	}
done:
	if len(received) != 3 {
		t.Fatalf("expected 3 updates, got %d: %+v", len(received), received)
	}
	if received[0].Action != "checking_inventory" {
		t.Errorf("expected first action 'checking_inventory', got %q", received[0].Action)
	}
	if received[1].Action != "reserving_oven" {
		t.Errorf("expected second action 'reserving_oven', got %q", received[1].Action)
	}
	if received[2].Type != "result" {
		t.Errorf("expected third type 'result', got %q", received[2].Type)
	}
}

// TestParseSSEStreamMultiLineEvent tests parsing SSE events that include id: and event: fields
// before the data: field, as some SSE implementations produce.
func TestParseSSEStreamMultiLineEvent(t *testing.T) {
	kitchen := NewKitchen()
	updates := make(chan CookingUpdate, 10)

	// Multi-line SSE events with id: and event: fields
	sseData := "id:1\nevent:message\ndata:{\"type\":\"action\",\"action\":\"checking_inventory\",\"message\":\"Checking\"}\n\nid:2\nevent:message\ndata:{\"type\":\"result\",\"action\":\"completed\",\"message\":\"Done\"}\n\n"
	body := bytes.NewReader([]byte(sseData))

	go func() {
		defer close(updates)
		err := kitchen.parseSSEStream(body, updates)
		if err != nil {
			t.Errorf("parseSSEStream returned error: %v", err)
		}
	}()

	var received []CookingUpdate
	timeout := time.After(2 * time.Second)
	for {
		select {
		case update, ok := <-updates:
			if !ok {
				goto done
			}
			received = append(received, update)
		case <-timeout:
			t.Fatal("timed out waiting for updates")
		}
	}
done:
	if len(received) != 2 {
		t.Fatalf("expected 2 updates, got %d: %+v", len(received), received)
	}
	if received[0].Action != "checking_inventory" {
		t.Errorf("expected first action 'checking_inventory', got %q", received[0].Action)
	}
	if received[1].Type != "result" {
		t.Errorf("expected second type 'result', got %q", received[1].Type)
	}
}

// TestCookForwardsFullCookingUpdateData tests that the kitchen sends message, toolName, and toolInput
// to the store along with the action status.
func TestCookForwardsFullCookingUpdateData(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cook/stream" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			flusher, _ := w.(http.Flusher)

			updates := []CookingUpdate{
				{Type: "action", Action: "checking_inventory", Message: "Checking available ingredients in inventory", ToolName: "getInventory", ToolInput: "{}"},
				{Type: "action", Action: "acquiring_ingredients", Message: "Acquiring ingredients: mozzarella", ToolName: "acquireItem", ToolInput: "{\"itemName\":\"mozzarella\"}"},
				{Type: "action", Action: "reserving_oven", Message: "Reserving oven for cooking: oven-1", ToolName: "reserveOven", ToolInput: "{\"ovenId\":\"oven-1\"}"},
				{Type: "result", Action: "completed", Message: "Pizza cooked successfully"},
			}

			for _, update := range updates {
				data, _ := json.Marshal(update)
				w.Write([]byte("data:" + string(data) + "\n\n"))
				flusher.Flush()
			}
		}
	}))
	defer agentServer.Close()

	eventsReceived := make(chan OrderEvent, 50)
	storeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/events" {
			var event OrderEvent
			json.NewDecoder(r.Body).Decode(&event)
			eventsReceived <- event
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer storeServer.Close()

	kitchen := NewKitchenWithConfig(KitchenConfig{
		StoreURL:        storeServer.URL,
		CookingAgentURL: agentServer.URL,
	})

	router := chi.NewRouter()
	router.Post("/cook", kitchen.HandleCook)

	orderID := uuid.New()
	req := CookRequest{
		OrderID:    orderID,
		OrderItems: []OrderItem{{PizzaType: "Margherita", Quantity: 1}},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/cook", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httpReq)

	// Collect events until DONE
	var events []OrderEvent
	timeout := time.After(5 * time.Second)
	for {
		select {
		case event := <-eventsReceived:
			events = append(events, event)
			if event.Status == "DONE" {
				// Verify that action events include message and toolName
				hasRichData := false
				for _, e := range events {
					if e.Status == "checking_inventory" && e.Message != "" && e.ToolName != "" {
						hasRichData = true
					}
				}
				if !hasRichData {
					t.Error("expected events to include message and toolName for action events")
				}

				// Verify reserving_oven has toolInput
				for _, e := range events {
					if e.Status == "reserving_oven" {
						if e.ToolInput == "" {
							t.Error("expected reserving_oven event to include toolInput")
						}
						if e.Message == "" {
							t.Error("expected reserving_oven event to include message")
						}
					}
				}
				return
			}
		case <-timeout:
			t.Fatalf("timed out waiting for events, got %d events: %+v", len(events), events)
		}
	}
}

// TestCookStreamsUpdatesToStore tests that the kitchen streams cooking updates to the store.
func TestCookStreamsUpdatesToStore(t *testing.T) {
	// Track requests received by the mock cooking-agent
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cook/stream" && r.Method == http.MethodPost {
			// Send SSE response with cooking updates
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.WriteHeader(http.StatusOK)

			flusher, ok := w.(http.Flusher)
			if !ok {
				t.Error("expected http.Flusher")
				return
			}

			// Send updates that should be forwarded to store
			updates := []CookingUpdate{
				{Type: "action", Action: "checking_inventory", Message: "Checking available ingredients"},
				{Type: "action", Action: "reserving_oven", Message: "Reserving oven for cooking"},
				{Type: "result", Action: "completed", Message: "Pizza cooked successfully"},
			}

			for _, update := range updates {
				data, _ := json.Marshal(update)
				w.Write([]byte("data: " + string(data) + "\n\n"))
				flusher.Flush()
			}
		}
	}))
	defer agentServer.Close()

	// Track events received by the mock store
	eventsReceived := make(chan OrderEvent, 50)
	storeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/events" {
			var event OrderEvent
			json.NewDecoder(r.Body).Decode(&event)
			eventsReceived <- event
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer storeServer.Close()

	kitchen := NewKitchenWithConfig(KitchenConfig{
		StoreURL:        storeServer.URL,
		CookingAgentURL: agentServer.URL,
	})

	router := chi.NewRouter()
	router.Post("/cook", kitchen.HandleCook)

	orderID := uuid.New()
	req := CookRequest{
		OrderID: orderID,
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 1},
		},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/cook", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpReq)

	// Collect events until we see DONE
	var events []OrderEvent
	timeout := time.After(5 * time.Second)

	for {
		select {
		case event := <-eventsReceived:
			events = append(events, event)
			if event.Status == "DONE" {
				// Verify we got streaming updates
				hasCheckingInventory := false
				hasReservingOven := false
				for _, e := range events {
					if e.Status == "checking_inventory" {
						hasCheckingInventory = true
					}
					if e.Status == "reserving_oven" {
						hasReservingOven = true
					}
				}
				if !hasCheckingInventory {
					t.Error("expected checking_inventory event")
				}
				if !hasReservingOven {
					t.Error("expected reserving_oven event")
				}
				return
			}
		case <-timeout:
			t.Fatalf("timed out waiting for events, got %d events: %+v", len(events), events)
		}
	}
}
