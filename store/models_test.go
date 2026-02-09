package store

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

// TestOrderItem verifies that OrderItem contains pizzaType and quantity fields.
func TestOrderItem(t *testing.T) {
	item := OrderItem{
		PizzaType: "Margherita",
		Quantity:  2,
	}

	if item.PizzaType != "Margherita" {
		t.Errorf("expected PizzaType 'Margherita', got '%s'", item.PizzaType)
	}
	if item.Quantity != 2 {
		t.Errorf("expected Quantity 2, got %d", item.Quantity)
	}
}

// TestOrder verifies that Order contains orderId(UUID), OrderItems, orderData, and orderStatus.
func TestOrder(t *testing.T) {
	orderID := uuid.New()
	items := []OrderItem{
		{PizzaType: "Pepperoni", Quantity: 1},
		{PizzaType: "Hawaiian", Quantity: 3},
	}

	order := Order{
		OrderID:     orderID,
		OrderItems:  items,
		OrderData:   "Extra cheese please",
		OrderStatus: "pending",
	}

	if order.OrderID != orderID {
		t.Errorf("expected OrderID %v, got %v", orderID, order.OrderID)
	}
	if len(order.OrderItems) != 2 {
		t.Errorf("expected 2 OrderItems, got %d", len(order.OrderItems))
	}
	if order.OrderItems[0].PizzaType != "Pepperoni" {
		t.Errorf("expected first item PizzaType 'Pepperoni', got '%s'", order.OrderItems[0].PizzaType)
	}
	if order.OrderData != "Extra cheese please" {
		t.Errorf("expected OrderData 'Extra cheese please', got '%s'", order.OrderData)
	}
	if order.OrderStatus != "pending" {
		t.Errorf("expected OrderStatus 'pending', got '%s'", order.OrderStatus)
	}
}

// TestOrderJSONSerialization verifies that Order can be serialized to and from JSON.
func TestOrderJSONSerialization(t *testing.T) {
	orderID := uuid.New()
	order := Order{
		OrderID: orderID,
		OrderItems: []OrderItem{
			{PizzaType: "Veggie", Quantity: 2},
		},
		OrderData:   "No onions",
		OrderStatus: "preparing",
	}

	// Test JSON marshaling
	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("failed to marshal order: %v", err)
	}

	// Test JSON unmarshaling
	var decoded Order
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal order: %v", err)
	}

	if decoded.OrderID != orderID {
		t.Errorf("expected OrderID %v after unmarshal, got %v", orderID, decoded.OrderID)
	}
	if decoded.OrderStatus != "preparing" {
		t.Errorf("expected OrderStatus 'preparing' after unmarshal, got '%s'", decoded.OrderStatus)
	}
}
