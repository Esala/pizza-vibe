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

// TestCookSendsUpdateEventsToStore tests that cooking sends update events to the store service.
func TestCookSendsUpdateEventsToStore(t *testing.T) {
	// Track events received by the mock store
	eventsReceived := make(chan OrderEvent, 10)
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
		CookingTimeFunc: func() int { return 1 }, // 1 second per pizza (fast for test)
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

	// Wait for events
	var events []OrderEvent
	timeout := time.After(5 * time.Second)
	for {
		select {
		case event := <-eventsReceived:
			events = append(events, event)
			if event.Status == "DONE" {
				// Got the DONE event, check results
				if len(events) < 2 {
					t.Errorf("expected at least 2 events (update + DONE), got %d", len(events))
				}
				// Last event should be DONE
				if events[len(events)-1].Status != "DONE" {
					t.Errorf("expected last event to be 'DONE', got '%s'", events[len(events)-1].Status)
				}
				// All events should have correct orderID and source
				for _, e := range events {
					if e.OrderID != orderID {
						t.Errorf("expected event orderID %s, got %s", orderID, e.OrderID)
					}
					if e.Source != "kitchen" {
						t.Errorf("expected event source 'kitchen', got '%s'", e.Source)
					}
				}
				return
			}
		case <-timeout:
			t.Fatal("timed out waiting for events")
		}
	}
}
