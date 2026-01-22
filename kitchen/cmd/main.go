// Package main is the entry point for the Kitchen service.
// It sets up the HTTP server with the /cook endpoint.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/salaboy/pizza-vibe/kitchen"
)

func main() {
	// Get port from environment variable or default to 8081
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Create kitchen instance
	k := kitchen.NewKitchen()

	// Set up router with middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Register routes
	r.Post("/cook", k.HandleCook)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Kitchen service starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
