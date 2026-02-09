// Package delivery provides the delivery service for the Pizza Vibe application.
// It handles delivering pizza orders by simulating delivery with progress updates.
package delivery

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

// DeliveryConfig contains configuration options for the Delivery service.
type DeliveryConfig struct {
	StoreURL         string
	DeliveryTimeFunc func() int // Returns delivery time in seconds
}

// OrderEvent represents an event sent to the store service.
type OrderEvent struct {
	OrderID uuid.UUID `json:"orderId"`
	Status  string    `json:"status"`
	Source  string    `json:"source"`
}

// Delivery manages pizza delivery operations and provides HTTP handlers for the delivery service.
type Delivery struct {
	rng              *rand.Rand
	storeURL         string
	httpClient       *http.Client
	deliveryTimeFunc func() int
}

// NewDelivery creates a new Delivery instance with a seeded random number generator.
// The default delivery time is a random interval between 5 and 20 seconds.
func NewDelivery() *Delivery {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Delivery{
		rng:      rng,
		storeURL: "http://store:8080",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		deliveryTimeFunc: func() int { return rng.Intn(16) + 5 },
	}
}

// NewDeliveryWithConfig creates a new Delivery instance with the given configuration.
func NewDeliveryWithConfig(config DeliveryConfig) *Delivery {
	d := NewDelivery()
	if config.StoreURL != "" {
		d.storeURL = config.StoreURL
	}
	if config.DeliveryTimeFunc != nil {
		d.deliveryTimeFunc = config.DeliveryTimeFunc
	}
	return d
}

// HandleDeliver handles POST /deliver requests to deliver pizza orders.
// It validates the request and starts the delivery simulation asynchronously.
func (d *Delivery) HandleDeliver(w http.ResponseWriter, r *http.Request) {
	var req DeliverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate that at least one item is provided
	if len(req.OrderItems) == 0 {
		http.Error(w, "Order must contain at least one item", http.StatusBadRequest)
		return
	}

	slog.Info("delivery request received", "orderId", req.OrderID, "items", len(req.OrderItems))

	// Start delivery in a goroutine (background; detach from request context)
	go d.deliverOrder(context.Background(), req.OrderID)

	// Return accepted response immediately
	resp := DeliverResponse{
		OrderID: req.OrderID,
		Status:  "delivering",
		Message: fmt.Sprintf("Started delivering %d item(s)", len(req.OrderItems)),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

// deliverOrder simulates delivering an order with a random delivery time between 5-20 seconds.
// It sends percentage-based progress updates every second and a final DELIVERED event.
func (d *Delivery) deliverOrder(ctx context.Context, orderID uuid.UUID) {
	deliveryTime := d.deliveryTimeFunc()
	startTime := time.Now()
	slog.Info("delivery started", "orderId", orderID, "deliveryTime", deliveryTime)

	for elapsed := 1; elapsed <= deliveryTime; elapsed++ {
		select {
		case <-ctx.Done():
			slog.Warn("delivery cancelled", "orderId", orderID, "error", ctx.Err())
			return
		default:
		}

		time.Sleep(1 * time.Second)

		// Calculate and send percentage update
		percent := (elapsed * 100) / deliveryTime
		d.sendEvent(ctx, orderID, fmt.Sprintf("delivering %d%%", percent))
	}

	duration := time.Since(startTime)
	slog.Info("delivery completed", "orderId", orderID, "duration", duration.Round(time.Second))

	// Send DELIVERED event
	d.sendEvent(ctx, orderID, "DELIVERED")
}

// sendEvent sends an event to the store service.
func (d *Delivery) sendEvent(ctx context.Context, orderID uuid.UUID, status string) {
	event := OrderEvent{
		OrderID: orderID,
		Status:  status,
		Source:  "delivery",
	}

	body, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal event", "orderId", orderID, "error", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.storeURL+"/events", bytes.NewReader(body))
	if err != nil {
		slog.Error("failed to create event request", "orderId", orderID, "error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		slog.Error("failed to send event to store", "orderId", orderID, "status", status, "error", err)
		return
	}
	defer resp.Body.Close()
}
