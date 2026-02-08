// Package delivery provides the delivery service for the Pizza Vibe application.
// It handles delivering pizza orders by calling the delivery agent and forwarding updates.
package delivery

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

// TestDeliverEndpointReturnsAccepted tests that the /deliver endpoint accepts a valid delivery request.
func TestDeliverEndpointReturnsAccepted(t *testing.T) {
	d := NewDelivery()
	router := chi.NewRouter()
	router.Post("/deliver", d.HandleDeliver)

	orderID := uuid.New()
	req := DeliverRequest{
		OrderID: orderID,
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 1},
		},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/deliver", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rr.Code)
	}
}

// TestDeliverEndpointInvalidJSON tests that the /deliver endpoint returns bad request for invalid JSON.
func TestDeliverEndpointInvalidJSON(t *testing.T) {
	d := NewDelivery()
	router := chi.NewRouter()
	router.Post("/deliver", d.HandleDeliver)

	httpReq := httptest.NewRequest(http.MethodPost, "/deliver", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

// TestDeliverEndpointEmptyOrderItems tests that the /deliver endpoint returns bad request for empty order items.
func TestDeliverEndpointEmptyOrderItems(t *testing.T) {
	d := NewDelivery()
	router := chi.NewRouter()
	router.Post("/deliver", d.HandleDeliver)

	orderID := uuid.New()
	req := DeliverRequest{
		OrderID:    orderID,
		OrderItems: []OrderItem{},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/deliver", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

// TestDeliverEndpointReturnsResponse tests that the /deliver endpoint returns a proper response body.
func TestDeliverEndpointReturnsResponse(t *testing.T) {
	d := NewDelivery()
	router := chi.NewRouter()
	router.Post("/deliver", d.HandleDeliver)

	orderID := uuid.New()
	req := DeliverRequest{
		OrderID: orderID,
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 2},
			{PizzaType: "Pepperoni", Quantity: 1},
		},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/deliver", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rr.Code)
	}

	var resp DeliverResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.OrderID != orderID {
		t.Errorf("expected orderID %s, got %s", orderID, resp.OrderID)
	}

	if resp.Status != "delivering" {
		t.Errorf("expected status 'delivering', got %s", resp.Status)
	}
}

// TestDeliverCallsDeliveryAgent tests that delivery calls the delivery-agent streaming endpoint.
func TestDeliverCallsDeliveryAgent(t *testing.T) {
	// Track requests received by the mock delivery-agent
	agentRequests := make(chan AgentDeliverRequest, 10)
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/deliver/stream" && r.Method == http.MethodPost {
			var req AgentDeliverRequest
			json.NewDecoder(r.Body).Decode(&req)
			agentRequests <- req

			// Send SSE response with delivery updates
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.WriteHeader(http.StatusOK)

			flusher, ok := w.(http.Flusher)
			if !ok {
				t.Error("expected http.Flusher")
				return
			}

			updates := []DeliveryUpdate{
				{Type: "action", Action: "checking_bikes", Message: "Checking available bikes for delivery", ToolName: "getBikes", ToolInput: "{}"},
				{Type: "action", Action: "reserving_bike", Message: "Reserving bike for delivery: bike-1", ToolName: "reserveBike", ToolInput: `{"bikeId":"bike-1","user":"delivery-agent-dave"}`},
				{Type: "action", Action: "checking_bike_status", Message: "Checking bike status: bike-1", ToolName: "getBike", ToolInput: `{"bikeId":"bike-1"}`},
				{Type: "result", Action: "completed", Message: "Delivery completed successfully"},
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

	d := NewDeliveryWithConfig(DeliveryConfig{
		StoreURL:         storeServer.URL,
		DeliveryAgentURL: agentServer.URL,
	})

	router := chi.NewRouter()
	router.Post("/deliver", d.HandleDeliver)

	orderID := uuid.New()
	req := DeliverRequest{
		OrderID: orderID,
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 1},
		},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/deliver", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpReq)

	// Wait for agent request
	timeout := time.After(5 * time.Second)
	select {
	case agentReq := <-agentRequests:
		if agentReq.OrderID != orderID.String() {
			t.Errorf("expected agent request orderId %s, got %s", orderID.String(), agentReq.OrderID)
		}
	case <-timeout:
		t.Fatal("timed out waiting for agent request")
	}

	// Wait for events until DELIVERED
	var events []OrderEvent
	for {
		select {
		case event := <-eventsReceived:
			events = append(events, event)
			if event.Status == "DELIVERED" {
				// Verify we got streaming updates + DELIVERED
				if len(events) < 4 { // At least 3 action events + DELIVERED
					t.Errorf("expected at least 4 events, got %d", len(events))
				}
				// All events should have correct orderID and source
				for _, e := range events {
					if e.OrderID != orderID {
						t.Errorf("expected event orderID %s, got %s", orderID, e.OrderID)
					}
					if e.Source != "delivery" {
						t.Errorf("expected event source 'delivery', got '%s'", e.Source)
					}
				}
				return
			}
		case <-timeout:
			t.Fatalf("timed out waiting for events, received %d events so far", len(events))
		}
	}
}

// TestDeliverStreamsUpdatesToStore tests that the delivery service streams delivery updates to the store.
func TestDeliverStreamsUpdatesToStore(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/deliver/stream" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.WriteHeader(http.StatusOK)

			flusher, ok := w.(http.Flusher)
			if !ok {
				t.Error("expected http.Flusher")
				return
			}

			updates := []DeliveryUpdate{
				{Type: "action", Action: "checking_bikes", Message: "Checking available bikes for delivery"},
				{Type: "action", Action: "reserving_bike", Message: "Reserving bike for delivery"},
				{Type: "action", Action: "checking_bike_status", Message: "Checking bike status"},
				{Type: "result", Action: "completed", Message: "Delivery completed"},
			}

			for _, update := range updates {
				data, _ := json.Marshal(update)
				w.Write([]byte("data: " + string(data) + "\n\n"))
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

	d := NewDeliveryWithConfig(DeliveryConfig{
		StoreURL:         storeServer.URL,
		DeliveryAgentURL: agentServer.URL,
	})

	router := chi.NewRouter()
	router.Post("/deliver", d.HandleDeliver)

	orderID := uuid.New()
	req := DeliverRequest{
		OrderID: orderID,
		OrderItems: []OrderItem{
			{PizzaType: "Margherita", Quantity: 1},
		},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/deliver", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpReq)

	// Collect events until we see DELIVERED
	var events []OrderEvent
	timeout := time.After(5 * time.Second)

	for {
		select {
		case event := <-eventsReceived:
			events = append(events, event)
			if event.Status == "DELIVERED" {
				// Verify we got streaming updates
				hasCheckingBikes := false
				hasReservingBike := false
				hasCheckingBikeStatus := false
				for _, e := range events {
					if e.Status == "checking_bikes" {
						hasCheckingBikes = true
					}
					if e.Status == "reserving_bike" {
						hasReservingBike = true
					}
					if e.Status == "checking_bike_status" {
						hasCheckingBikeStatus = true
					}
				}
				if !hasCheckingBikes {
					t.Error("expected checking_bikes event")
				}
				if !hasReservingBike {
					t.Error("expected reserving_bike event")
				}
				if !hasCheckingBikeStatus {
					t.Error("expected checking_bike_status event")
				}
				return
			}
		case <-timeout:
			t.Fatalf("timed out waiting for events, got %d events: %+v", len(events), events)
		}
	}
}

// TestDeliverForwardsFullDeliveryUpdateData tests that the delivery service sends message, toolName, and toolInput
// to the store along with the action status.
func TestDeliverForwardsFullDeliveryUpdateData(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/deliver/stream" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			flusher, _ := w.(http.Flusher)

			updates := []DeliveryUpdate{
				{Type: "action", Action: "checking_bikes", Message: "Checking available bikes for delivery", ToolName: "getBikes", ToolInput: "{}"},
				{Type: "action", Action: "reserving_bike", Message: "Reserving bike for delivery: bike-1", ToolName: "reserveBike", ToolInput: `{"bikeId":"bike-1","user":"delivery-agent-dave"}`},
				{Type: "action", Action: "checking_bike_status", Message: "Checking bike status: bike-1", ToolName: "getBike", ToolInput: `{"bikeId":"bike-1"}`},
				{Type: "result", Action: "completed", Message: "Delivery completed successfully"},
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

	d := NewDeliveryWithConfig(DeliveryConfig{
		StoreURL:         storeServer.URL,
		DeliveryAgentURL: agentServer.URL,
	})

	router := chi.NewRouter()
	router.Post("/deliver", d.HandleDeliver)

	orderID := uuid.New()
	req := DeliverRequest{
		OrderID:    orderID,
		OrderItems: []OrderItem{{PizzaType: "Margherita", Quantity: 1}},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/deliver", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httpReq)

	// Collect events until DELIVERED
	var events []OrderEvent
	timeout := time.After(5 * time.Second)
	for {
		select {
		case event := <-eventsReceived:
			events = append(events, event)
			if event.Status == "DELIVERED" {
				// Verify that action events include message and toolName
				hasRichData := false
				for _, e := range events {
					if e.Status == "checking_bikes" && e.Message != "" && e.ToolName != "" {
						hasRichData = true
					}
				}
				if !hasRichData {
					t.Error("expected events to include message and toolName for action events")
				}

				// Verify reserving_bike has toolInput
				for _, e := range events {
					if e.Status == "reserving_bike" {
						if e.ToolInput == "" {
							t.Error("expected reserving_bike event to include toolInput")
						}
						if e.Message == "" {
							t.Error("expected reserving_bike event to include message")
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

// TestDeliverDoesNotForwardPartialEvents tests that partial (token-level) events from
// the delivery agent are not forwarded to the store.
func TestDeliverDoesNotForwardPartialEvents(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/deliver/stream" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			flusher, _ := w.(http.Flusher)

			updates := []DeliveryUpdate{
				{Type: "action", Action: "checking_bikes", Message: "Checking available bikes"},
				{Type: "partial", Message: "I am now"},
				{Type: "partial", Message: " checking the"},
				{Type: "partial", Message: " bike status"},
				{Type: "action", Action: "reserving_bike", Message: "Reserving bike"},
				{Type: "result", Action: "completed", Message: "Done"},
			}

			for _, update := range updates {
				data, _ := json.Marshal(update)
				w.Write([]byte("data: " + string(data) + "\n\n"))
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

	d := NewDeliveryWithConfig(DeliveryConfig{
		StoreURL:         storeServer.URL,
		DeliveryAgentURL: agentServer.URL,
	})

	router := chi.NewRouter()
	router.Post("/deliver", d.HandleDeliver)

	orderID := uuid.New()
	req := DeliverRequest{
		OrderID:    orderID,
		OrderItems: []OrderItem{{PizzaType: "Margherita", Quantity: 1}},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/deliver", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httpReq)

	// Collect all events until DELIVERED
	var events []OrderEvent
	timeout := time.After(5 * time.Second)
	for {
		select {
		case event := <-eventsReceived:
			events = append(events, event)
			if event.Status == "DELIVERED" {
				// Should only have action events + DELIVERED (no partial events)
				// Expected: checking_bikes, reserving_bike, DELIVERED = 3 events
				if len(events) != 3 {
					t.Errorf("expected 3 events (2 actions + DELIVERED), got %d: %+v", len(events), events)
				}
				return
			}
		case <-timeout:
			t.Fatalf("timed out waiting for events, got %d events: %+v", len(events), events)
		}
	}
}

// TestParseSSEStreamWithoutSpace tests parsing SSE events in "data:value" format (no space after colon),
// which is the format Quarkus RESTEasy Reactive produces.
func TestParseSSEStreamWithoutSpace(t *testing.T) {
	d := NewDelivery()
	updates := make(chan DeliveryUpdate, 10)

	// Simulate Quarkus SSE format: "data:{json}\n\n" (no space after data:)
	sseData := `data:{"type":"action","action":"checking_bikes","message":"Checking available bikes"}

data:{"type":"action","action":"reserving_bike","message":"Reserving bike"}

data:{"type":"result","action":"completed","message":"Delivery completed"}

`
	body := bytes.NewReader([]byte(sseData))

	go func() {
		defer close(updates)
		err := d.parseSSEStream(body, updates)
		if err != nil {
			t.Errorf("parseSSEStream returned error: %v", err)
		}
	}()

	var received []DeliveryUpdate
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
	if received[0].Action != "checking_bikes" {
		t.Errorf("expected first action 'checking_bikes', got %q", received[0].Action)
	}
	if received[1].Action != "reserving_bike" {
		t.Errorf("expected second action 'reserving_bike', got %q", received[1].Action)
	}
	if received[2].Type != "result" {
		t.Errorf("expected third type 'result', got %q", received[2].Type)
	}
}

// TestParseSSEStreamMultiLineEvent tests parsing SSE events that include id: and event: fields
// before the data: field, as some SSE implementations produce.
func TestParseSSEStreamMultiLineEvent(t *testing.T) {
	d := NewDelivery()
	updates := make(chan DeliveryUpdate, 10)

	// Multi-line SSE events with id: and event: fields
	sseData := "id:1\nevent:message\ndata:{\"type\":\"action\",\"action\":\"checking_bikes\",\"message\":\"Checking bikes\"}\n\nid:2\nevent:message\ndata:{\"type\":\"result\",\"action\":\"completed\",\"message\":\"Done\"}\n\n"
	body := bytes.NewReader([]byte(sseData))

	go func() {
		defer close(updates)
		err := d.parseSSEStream(body, updates)
		if err != nil {
			t.Errorf("parseSSEStream returned error: %v", err)
		}
	}()

	var received []DeliveryUpdate
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
	if received[0].Action != "checking_bikes" {
		t.Errorf("expected first action 'checking_bikes', got %q", received[0].Action)
	}
	if received[1].Type != "result" {
		t.Errorf("expected second type 'result', got %q", received[1].Type)
	}
}

// TestDeliverSendsDeliveredEventAtEnd tests that a DELIVERED event is always sent
// after the agent completes, regardless of streaming content.
func TestDeliverSendsDeliveredEventAtEnd(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/deliver/stream" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			flusher, _ := w.(http.Flusher)

			// Only send a result event (no action events)
			update := DeliveryUpdate{Type: "result", Action: "completed", Message: "Delivery completed"}
			data, _ := json.Marshal(update)
			w.Write([]byte("data: " + string(data) + "\n\n"))
			flusher.Flush()
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

	d := NewDeliveryWithConfig(DeliveryConfig{
		StoreURL:         storeServer.URL,
		DeliveryAgentURL: agentServer.URL,
	})

	router := chi.NewRouter()
	router.Post("/deliver", d.HandleDeliver)

	orderID := uuid.New()
	req := DeliverRequest{
		OrderID:    orderID,
		OrderItems: []OrderItem{{PizzaType: "Margherita", Quantity: 1}},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/deliver", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httpReq)

	// Wait for DELIVERED event
	timeout := time.After(5 * time.Second)
	select {
	case event := <-eventsReceived:
		if event.Status != "DELIVERED" {
			t.Errorf("expected DELIVERED event, got %s", event.Status)
		}
		if event.Source != "delivery" {
			t.Errorf("expected source 'delivery', got %s", event.Source)
		}
	case <-timeout:
		t.Fatal("timed out waiting for DELIVERED event")
	}
}
