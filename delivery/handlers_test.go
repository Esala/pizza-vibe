// Package delivery provides the delivery service for the Pizza Vibe application.
// It handles delivering pizza orders by simulating delivery with progress updates.
package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// TestDeliverSendsProgressEventsToStore tests that delivery sends percentage-based progress
// events and a final DELIVERED event to the store service.
func TestDeliverSendsProgressEventsToStore(t *testing.T) {
	// Track events received by the mock store
	eventsReceived := make(chan OrderEvent, 100)
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
		DeliveryTimeFunc: func() int { return 3 }, // 3 seconds delivery (fast for test)
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

	// Wait for events
	var events []OrderEvent
	timeout := time.After(10 * time.Second)
	for {
		select {
		case event := <-eventsReceived:
			events = append(events, event)
			if event.Status == "DELIVERED" {
				// Got the DELIVERED event, check results
				if len(events) < 2 {
					t.Errorf("expected at least 2 events (progress + DELIVERED), got %d", len(events))
				}
				// Last event should be DELIVERED
				lastEvent := events[len(events)-1]
				if lastEvent.Status != "DELIVERED" {
					t.Errorf("expected last event to be 'DELIVERED', got '%s'", lastEvent.Status)
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

// TestDeliverSendsPercentageUpdates tests that progress events include percentage completion.
func TestDeliverSendsPercentageUpdates(t *testing.T) {
	eventsReceived := make(chan OrderEvent, 100)
	storeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/events" {
			var event OrderEvent
			json.NewDecoder(r.Body).Decode(&event)
			eventsReceived <- event
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer storeServer.Close()

	deliveryTime := 5
	d := NewDeliveryWithConfig(DeliveryConfig{
		StoreURL:         storeServer.URL,
		DeliveryTimeFunc: func() int { return deliveryTime },
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

	// Collect all events
	var events []OrderEvent
	timeout := time.After(15 * time.Second)
	for {
		select {
		case event := <-eventsReceived:
			events = append(events, event)
			if event.Status == "DELIVERED" {
				// Verify percentage updates were sent
				// With 5 second delivery time, we expect events at 1s, 2s, 3s, 4s (progress) + DELIVERED
				progressEvents := events[:len(events)-1]
				if len(progressEvents) == 0 {
					t.Error("expected progress events with percentages before DELIVERED")
					return
				}

				// Each progress event should contain a percentage
				for i, pe := range progressEvents {
					expectedPercent := ((i + 1) * 100) / deliveryTime
					expectedStatus := fmt.Sprintf("delivering %d%%", expectedPercent)
					if pe.Status != expectedStatus {
						t.Errorf("event %d: expected status '%s', got '%s'", i, expectedStatus, pe.Status)
					}
				}
				return
			}
		case <-timeout:
			t.Fatalf("timed out waiting for events, received %d events so far", len(events))
		}
	}
}
