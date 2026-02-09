// Package inventory provides the inventory service for the Pizza Vibe application.
// It manages pizza ingredient stock levels and provides REST endpoints for inventory operations.
package inventory

// ItemResponse represents the response for a single inventory item query.
type ItemResponse struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

// AcquireResponse represents the response when acquiring an item from inventory.
type AcquireResponse struct {
	Item              string `json:"item"`
	Status            string `json:"status"`
	RemainingQuantity int    `json:"remainingQuantity"`
}

// Status constants for inventory acquisition responses.
const (
	StatusAcquired = "ACQUIRED"
	StatusEmpty    = "EMPTY"
)

// DefaultInventory returns the default inventory stock levels.
func DefaultInventory() map[string]int {
	return map[string]int{
		"Pepperoni":  10,
		"Pineapple":  10,
		"PizzaDough": 10,
		"Mozzarella": 10,
		"Sauce":      10,
	}
}
