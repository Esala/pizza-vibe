// Package kitchen provides the kitchen service for the Pizza Vibe application.
// It handles cooking pizza orders by processing order items.
package kitchen

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
