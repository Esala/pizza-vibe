package inventory

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// TestHandleGetAll tests the GET /inventory endpoint
func TestHandleGetAll(t *testing.T) {
	inv := NewInventory()

	r := chi.NewRouter()
	r.Get("/inventory", inv.HandleGetAll)

	req, err := http.NewRequest("GET", "/inventory", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var stock map[string]int
	if err := json.Unmarshal(rr.Body.Bytes(), &stock); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	expectedItems := map[string]int{
		"Pepperoni":  10,
		"Pineapple":  10,
		"PizzaDough": 10,
		"Mozzarella": 10,
		"Sauce":      10,
	}

	for item, expectedQty := range expectedItems {
		if qty, ok := stock[item]; !ok {
			t.Errorf("expected item %s not found in inventory", item)
		} else if qty != expectedQty {
			t.Errorf("expected quantity %d for %s, got %d", expectedQty, item, qty)
		}
	}
}

// TestHandleGetItem tests the GET /inventory/{item} endpoint
func TestHandleGetItem(t *testing.T) {
	inv := NewInventory()

	r := chi.NewRouter()
	r.Get("/inventory/{item}", inv.HandleGetItem)

	req, err := http.NewRequest("GET", "/inventory/Pepperoni", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response ItemResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Item != "Pepperoni" {
		t.Errorf("expected item Pepperoni, got %s", response.Item)
	}
	if response.Quantity != 10 {
		t.Errorf("expected quantity 10, got %d", response.Quantity)
	}
}

// TestHandleGetItemNotFound tests GET /inventory/{item} for non-existent item
func TestHandleGetItemNotFound(t *testing.T) {
	inv := NewInventory()

	r := chi.NewRouter()
	r.Get("/inventory/{item}", inv.HandleGetItem)

	req, err := http.NewRequest("GET", "/inventory/NonExistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

// TestHandleAcquireItem tests POST /inventory/{item} - successfully acquiring an item
func TestHandleAcquireItem(t *testing.T) {
	inv := NewInventory()

	r := chi.NewRouter()
	r.Post("/inventory/{item}", inv.HandleAcquireItem)

	req, err := http.NewRequest("POST", "/inventory/Pepperoni", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response AcquireResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Status != StatusAcquired {
		t.Errorf("expected status ACQUIRED, got %s", response.Status)
	}
	if response.Item != "Pepperoni" {
		t.Errorf("expected item Pepperoni, got %s", response.Item)
	}
	if response.RemainingQuantity != 9 {
		t.Errorf("expected remaining quantity 9, got %d", response.RemainingQuantity)
	}
}

// TestHandleAcquireItemEmpty tests POST /inventory/{item} when item is empty
func TestHandleAcquireItemEmpty(t *testing.T) {
	// Create inventory with empty Pepperoni
	inv := NewInventoryWithStock(map[string]int{
		"Pepperoni":  0,
		"Pineapple":  10,
		"PizzaDough": 10,
		"Mozzarella": 10,
		"Sauce":      10,
	})

	r := chi.NewRouter()
	r.Post("/inventory/{item}", inv.HandleAcquireItem)

	req, err := http.NewRequest("POST", "/inventory/Pepperoni", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response AcquireResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Status != StatusEmpty {
		t.Errorf("expected status EMPTY, got %s", response.Status)
	}
	if response.Item != "Pepperoni" {
		t.Errorf("expected item Pepperoni, got %s", response.Item)
	}
	if response.RemainingQuantity != 0 {
		t.Errorf("expected remaining quantity 0, got %d", response.RemainingQuantity)
	}
}

// TestHandleAcquireItemNotFound tests POST /inventory/{item} for non-existent item
func TestHandleAcquireItemNotFound(t *testing.T) {
	inv := NewInventory()

	r := chi.NewRouter()
	r.Post("/inventory/{item}", inv.HandleAcquireItem)

	req, err := http.NewRequest("POST", "/inventory/NonExistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

// TestHandleAddQuantity tests POST /inventory/{item}/add - adding quantity to an item
func TestHandleAddQuantity(t *testing.T) {
	inv := NewInventory()

	r := chi.NewRouter()
	r.Post("/inventory/{item}/add", inv.HandleAddQuantity)

	body := strings.NewReader(`{"quantity": 5}`)
	req, err := http.NewRequest("POST", "/inventory/Pepperoni/add", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response ItemResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Item != "Pepperoni" {
		t.Errorf("expected item Pepperoni, got %s", response.Item)
	}
	if response.Quantity != 15 {
		t.Errorf("expected quantity 15, got %d", response.Quantity)
	}
}

// TestHandleAddQuantityNotFound tests POST /inventory/{item}/add for non-existent item
func TestHandleAddQuantityNotFound(t *testing.T) {
	inv := NewInventory()

	r := chi.NewRouter()
	r.Post("/inventory/{item}/add", inv.HandleAddQuantity)

	body := strings.NewReader(`{"quantity": 5}`)
	req, err := http.NewRequest("POST", "/inventory/NonExistent/add", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

// TestHandleAddQuantityInvalidBody tests POST /inventory/{item}/add with invalid JSON
func TestHandleAddQuantityInvalidBody(t *testing.T) {
	inv := NewInventory()

	r := chi.NewRouter()
	r.Post("/inventory/{item}/add", inv.HandleAddQuantity)

	body := strings.NewReader(`{invalid json}`)
	req, err := http.NewRequest("POST", "/inventory/Pepperoni/add", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
