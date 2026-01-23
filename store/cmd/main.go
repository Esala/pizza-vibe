// Package main provides the entry point for the store service.
// The store service exposes REST endpoints for pizza orders and
// a WebSocket endpoint for real-time order updates.
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/salaboy/pizza-vibe/store"
)

func main() {
	s := store.NewStore()
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// REST endpoints
	r.Post("/order", s.HandleCreateOrder) // Create a new pizza order
	r.Post("/events", s.HandleEvent)      // Receive events from kitchen/delivery

	// WebSocket endpoint
	r.Get("/ws", s.HandleWebSocket) // Real-time order updates

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Store service starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
