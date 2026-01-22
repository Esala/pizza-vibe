// Package store provides the store service for the Pizza Vibe application.
// It handles pizza orders, receives events from kitchen and delivery services,
// and sends real-time updates to frontend clients via WebSocket.
package store

import (
	"encoding/json"

	"github.com/google/uuid"
)

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

// MarshalJSON serializes the Order to JSON format.
func (o Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	return json.Marshal(&struct {
		Alias
	}{
		Alias: Alias(o),
	})
}

// UnmarshalJSON deserializes JSON data into an Order.
func (o *Order) UnmarshalJSON(data []byte) error {
	type Alias Order
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	return json.Unmarshal(data, aux)
}
