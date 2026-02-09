// Package store provides the store service for the Pizza Vibe application.
// It handles pizza orders, receives events from kitchen and delivery services,
// and sends real-time updates to frontend clients via WebSocket.
package store

import "github.com/google/uuid"

// OrderItem represents a single item in an order, containing the pizza type
// and the quantity requested.
type OrderItem struct {
	PizzaType string `json:"pizzaType"`
	Quantity  int    `json:"quantity"`
}

// Order represents a pizza order with a unique identifier, items, additional data,
// and current status.
type Order struct {
	OrderID     uuid.UUID   `json:"orderId"`
	OrderItems  []OrderItem `json:"orderItems"`
	OrderData   string      `json:"orderData"`
	OrderStatus string      `json:"orderStatus"`
}
