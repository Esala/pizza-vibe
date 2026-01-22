// Package kitchen provides the kitchen service for the Pizza Vibe application.
// It handles cooking pizza orders by processing order items with simulated cooking times.
package kitchen

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Kitchen manages pizza cooking operations and provides HTTP handlers for the kitchen service.
type Kitchen struct {
	rng *rand.Rand
}

// NewKitchen creates a new Kitchen instance with a seeded random number generator.
func NewKitchen() *Kitchen {
	return &Kitchen{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
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

	// Start cooking in a goroutine
	go k.cookItems(req.OrderID.String(), req.OrderItems)

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
// It prints the cooking progress to the terminal.
func (k *Kitchen) cookItems(orderID string, items []OrderItem) {
	for _, item := range items {
		for i := 0; i < item.Quantity; i++ {
			// Random cooking time between 1-10 seconds
			cookingTime := k.rng.Intn(10) + 1
			time.Sleep(time.Duration(cookingTime) * time.Second)
			fmt.Printf("[Order %s] Cooked %s (took %d seconds)\n", orderID, item.PizzaType, cookingTime)
		}
	}
	fmt.Printf("[Order %s] All items cooked!\n", orderID)
}
