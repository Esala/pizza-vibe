// Package delivery provides the delivery service for the Pizza Vibe application.
// It handles delivering pizza orders by calling the delivery agent and forwarding updates.
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

// AgentDeliverRequest represents the request body for the delivery-agent.
// It contains the order ID to be delivered.
type AgentDeliverRequest struct {
	OrderID string `json:"orderId"`
}

// DeliveryUpdate represents a streaming update from the delivery agent.
// These updates inform the client about the current action being performed.
type DeliveryUpdate struct {
	Type      string `json:"type"`      // Type of update: "action", "partial", "result"
	Action    string `json:"action"`    // The action being performed
	Message   string `json:"message"`   // Human-readable message describing the update
	ToolName  string `json:"toolName"`  // The name of the tool being executed (if applicable)
	ToolInput string `json:"toolInput"` // The input to the tool (if applicable)
}
