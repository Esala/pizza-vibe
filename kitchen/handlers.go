// Package kitchen provides the kitchen service for the Pizza Vibe application.
// It handles cooking pizza orders by processing order items with simulated cooking times.
package kitchen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// KitchenConfig contains configuration options for the Kitchen service.
type KitchenConfig struct {
	StoreURL        string
	CookingTimeFunc func() int // Returns cooking time in seconds for each item
}

// OrderEvent represents an event sent to the store service.
type OrderEvent struct {
	OrderID uuid.UUID `json:"orderId"`
	Status  string    `json:"status"`
	Source  string    `json:"source"`
}

// Kitchen manages pizza cooking operations and provides HTTP handlers for the kitchen service.
type Kitchen struct {
	rng             *rand.Rand
	storeURL        string
	httpClient      *http.Client
	cookingTimeFunc func() int
}

// NewKitchen creates a new Kitchen instance with a seeded random number generator.
func NewKitchen() *Kitchen {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Kitchen{
		rng:      rng,
		storeURL: "http://store:8080",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cookingTimeFunc: func() int { return rng.Intn(10) + 1 },
	}
}

// NewKitchenWithConfig creates a new Kitchen instance with the given configuration.
func NewKitchenWithConfig(config KitchenConfig) *Kitchen {
	k := NewKitchen()
	if config.StoreURL != "" {
		k.storeURL = config.StoreURL
	}
	if config.CookingTimeFunc != nil {
		k.cookingTimeFunc = config.CookingTimeFunc
	}
	return k
}

// HandleCook handles POST /cook requests to cook pizza order items.
// It validates the request and starts cooking the items asynchronously.
// Each item takes a random time from 1 to 10 seconds to cook.
func (k *Kitchen) HandleCook(w http.ResponseWriter, r *http.Request) {
	var req CookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate that at least one item is provided
	if len(req.OrderItems) == 0 {
		http.Error(w, "Order must contain at least one item", http.StatusBadRequest)
		return
	}

	slog.Info("cook request received", "orderId", req.OrderID, "items", len(req.OrderItems))

	// Start cooking in a goroutine (background; detach from request context)
	go k.cookItems(context.Background(), req.OrderID, req.OrderItems)

	// Return accepted response immediately
	resp := CookResponse{
		OrderID: req.OrderID,
		Status:  "cooking",
		Message: fmt.Sprintf("Started cooking %d item(s)", len(req.OrderItems)),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

// cookItems simulates cooking each order item with a random cooking time between 1-10 seconds.
// It logs the cooking progress and sends update events to the store service.
func (k *Kitchen) cookItems(ctx context.Context, orderID uuid.UUID, items []OrderItem) {
	for _, item := range items {
		for i := 0; i < item.Quantity; i++ {
			// Get cooking time
			cookingTime := k.cookingTimeFunc()
			startTime := time.Now()

			// Send update events every second while cooking
			for elapsed := 0; elapsed < cookingTime; elapsed++ {
				select {
				case <-ctx.Done():
					slog.Warn("cooking cancelled", "orderId", orderID, "error", ctx.Err())
					return
				default:
				}
				k.sendEvent(ctx, orderID, fmt.Sprintf("cooking %s (%d/%d)", item.PizzaType, i+1, item.Quantity))
				time.Sleep(1 * time.Second)
			}

			duration := time.Since(startTime)
			slog.Info("item cooked", "orderId", orderID, "pizzaType", item.PizzaType, "duration", duration.Round(time.Second))
		}
	}
	slog.Info("all items cooked", "orderId", orderID)

	// Send DONE event
	k.sendEvent(ctx, orderID, "DONE")
}

// sendEvent sends an event to the store service.
func (k *Kitchen) sendEvent(ctx context.Context, orderID uuid.UUID, status string) {
	event := OrderEvent{
		OrderID: orderID,
		Status:  status,
		Source:  "kitchen",
	}

	body, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal event", "orderId", orderID, "error", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, k.storeURL+"/events", bytes.NewReader(body))
	if err != nil {
		slog.Error("failed to create event request", "orderId", orderID, "error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := k.httpClient.Do(req)
	if err != nil {
		slog.Error("failed to send event to store", "orderId", orderID, "status", status, "error", err)
		return
	}
	defer resp.Body.Close()
}
