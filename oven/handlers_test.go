package oven

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// TestHandleGetAll tests the GET /ovens/ endpoint
func TestHandleGetAll(t *testing.T) {
	svc := NewOvenService()

	r := chi.NewRouter()
	r.Get("/ovens/", svc.HandleGetAll)

	req, err := http.NewRequest("GET", "/ovens/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var ovens []Oven
	if err := json.Unmarshal(rr.Body.Bytes(), &ovens); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if len(ovens) != 4 {
		t.Errorf("expected 4 ovens, got %d", len(ovens))
	}

	for _, oven := range ovens {
		if oven.Status != StatusAvailable {
			t.Errorf("expected oven status AVAILABLE, got %s", oven.Status)
		}
	}
}

// TestHandleGetByID tests the GET /ovens/{ovenId} endpoint
func TestHandleGetByID(t *testing.T) {
	svc := NewOvenService()

	r := chi.NewRouter()
	r.Get("/ovens/{ovenId}", svc.HandleGetByID)

	req, err := http.NewRequest("GET", "/ovens/oven-1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var oven Oven
	if err := json.Unmarshal(rr.Body.Bytes(), &oven); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if oven.ID != "oven-1" {
		t.Errorf("expected oven ID oven-1, got %s", oven.ID)
	}
	if oven.Status != StatusAvailable {
		t.Errorf("expected oven status AVAILABLE, got %s", oven.Status)
	}
}

// TestHandleGetByIDNotFound tests GET /ovens/{ovenId} for non-existent oven
func TestHandleGetByIDNotFound(t *testing.T) {
	svc := NewOvenService()

	r := chi.NewRouter()
	r.Get("/ovens/{ovenId}", svc.HandleGetByID)

	req, err := http.NewRequest("GET", "/ovens/oven-99", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

// TestHandleReserve tests POST /ovens/{ovenId} - successfully reserving an oven
func TestHandleReserve(t *testing.T) {
	svc := NewOvenService()

	r := chi.NewRouter()
	r.Post("/ovens/{ovenId}", svc.HandleReserve)

	req, err := http.NewRequest("POST", "/ovens/oven-1?user=chef1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var oven Oven
	if err := json.Unmarshal(rr.Body.Bytes(), &oven); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if oven.ID != "oven-1" {
		t.Errorf("expected oven ID oven-1, got %s", oven.ID)
	}
	if oven.Status != StatusReserved {
		t.Errorf("expected oven status RESERVED, got %s", oven.Status)
	}
	if oven.User != "chef1" {
		t.Errorf("expected user chef1, got %s", oven.User)
	}
}

// TestHandleReserveAlreadyReserved tests POST /ovens/{ovenId} when oven is already reserved
func TestHandleReserveAlreadyReserved(t *testing.T) {
	svc := NewOvenService()

	r := chi.NewRouter()
	r.Post("/ovens/{ovenId}", svc.HandleReserve)

	// First reservation
	req1, _ := http.NewRequest("POST", "/ovens/oven-1?user=chef1", nil)
	rr1 := httptest.NewRecorder()
	r.ServeHTTP(rr1, req1)

	// Second reservation attempt
	req2, err := http.NewRequest("POST", "/ovens/oven-1?user=chef2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)

	if status := rr2.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}
}

// TestHandleReserveMissingUser tests POST /ovens/{ovenId} without user parameter
func TestHandleReserveMissingUser(t *testing.T) {
	svc := NewOvenService()

	r := chi.NewRouter()
	r.Post("/ovens/{ovenId}", svc.HandleReserve)

	req, err := http.NewRequest("POST", "/ovens/oven-1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// TestHandleReserveNotFound tests POST /ovens/{ovenId} for non-existent oven
func TestHandleReserveNotFound(t *testing.T) {
	svc := NewOvenService()

	r := chi.NewRouter()
	r.Post("/ovens/{ovenId}", svc.HandleReserve)

	req, err := http.NewRequest("POST", "/ovens/oven-99?user=chef1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

// TestHandleRelease tests DELETE /ovens/{ovenId} - successfully releasing an oven
func TestHandleRelease(t *testing.T) {
	svc := NewOvenService()

	r := chi.NewRouter()
	r.Post("/ovens/{ovenId}", svc.HandleReserve)
	r.Delete("/ovens/{ovenId}", svc.HandleRelease)

	// First reserve the oven
	req1, _ := http.NewRequest("POST", "/ovens/oven-1?user=chef1", nil)
	rr1 := httptest.NewRecorder()
	r.ServeHTTP(rr1, req1)

	// Now release it
	req2, err := http.NewRequest("DELETE", "/ovens/oven-1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)

	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var oven Oven
	if err := json.Unmarshal(rr2.Body.Bytes(), &oven); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if oven.ID != "oven-1" {
		t.Errorf("expected oven ID oven-1, got %s", oven.ID)
	}
	if oven.Status != StatusAvailable {
		t.Errorf("expected oven status AVAILABLE, got %s", oven.Status)
	}
	if oven.User != "" {
		t.Errorf("expected empty user, got %s", oven.User)
	}
}

// TestHandleReleaseAlreadyAvailable tests DELETE /ovens/{ovenId} when oven is already available
func TestHandleReleaseAlreadyAvailable(t *testing.T) {
	svc := NewOvenService()

	r := chi.NewRouter()
	r.Delete("/ovens/{ovenId}", svc.HandleRelease)

	req, err := http.NewRequest("DELETE", "/ovens/oven-1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}
}

// TestHandleReleaseNotFound tests DELETE /ovens/{ovenId} for non-existent oven
func TestHandleReleaseNotFound(t *testing.T) {
	svc := NewOvenService()

	r := chi.NewRouter()
	r.Delete("/ovens/{ovenId}", svc.HandleRelease)

	req, err := http.NewRequest("DELETE", "/ovens/oven-99", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
