// Package kitchen provides the kitchen service for the Pizza Vibe application.
// It handles cooking pizza orders by processing order items with simulated cooking times.
package kitchen

import "github.com/google/uuid"

// OrderItem represents a single item in an order, containing the pizza type
// and the quantity requested.
type OrderItem struct {
	PizzaType string `json:"pizzaType"`
	Quantity  int    `json:"quantity"`
}

// CookRequest represents the request body for cooking an order.
// It contains the order ID and the items to be cooked.
type CookRequest struct {
	OrderID    uuid.UUID   `json:"orderId"`
	OrderItems []OrderItem `json:"orderItems"`
}

// CookResponse represents the response returned after accepting a cook request.
type CookResponse struct {
	OrderID uuid.UUID `json:"orderId"`
	Status  string    `json:"status"`
	Message string    `json:"message,omitempty"`
}

// CookedItem represents a single item that has been cooked, including the time
// it took to cook.
type CookedItem struct {
	PizzaType   string `json:"pizzaType"`
	Quantity    int    `json:"quantity"`
	CookingTime int    `json:"cookingTime"` // in seconds
}
