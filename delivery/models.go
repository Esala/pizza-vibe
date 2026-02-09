// Package delivery provides the delivery service for the Pizza Vibe application.
// It handles delivering pizza orders by simulating delivery with progress updates.
package delivery

import "github.com/google/uuid"

// OrderItem represents a single item in an order, containing the pizza type
// and the quantity requested.
type OrderItem struct {
	PizzaType string `json:"pizzaType"`
	Quantity  int    `json:"quantity"`
}

// DeliverRequest represents the request body for delivering an order.
// It contains the order ID and the items to be delivered.
type DeliverRequest struct {
	OrderID    uuid.UUID   `json:"orderId"`
	OrderItems []OrderItem `json:"orderItems"`
}

// DeliverResponse represents the response returned after accepting a delivery request.
type DeliverResponse struct {
	OrderID uuid.UUID `json:"orderId"`
	Status  string    `json:"status"`
	Message string    `json:"message,omitempty"`
}
