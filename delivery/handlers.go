// Package delivery provides the delivery service for the Pizza Vibe application.
// It handles delivering pizza orders by calling the delivery agent and forwarding updates.
package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// DeliveryConfig contains configuration options for the Delivery service.
type DeliveryConfig struct {
	StoreURL         string
	DeliveryAgentURL string
}

// OrderEvent represents an event sent to the store service.
type OrderEvent struct {
	OrderID   uuid.UUID `json:"orderId"`
	Status    string    `json:"status"`
	Source    string    `json:"source"`
	Message   string    `json:"message,omitempty"`
	ToolName  string    `json:"toolName,omitempty"`
	ToolInput string    `json:"toolInput,omitempty"`
}

// Delivery manages pizza delivery operations and provides HTTP handlers for the delivery service.
type Delivery struct {
	storeURL         string
	deliveryAgentURL string
	httpClient       *http.Client
	streamingClient  *http.Client // No timeout for SSE streaming connections
}

// NewDelivery creates a new Delivery instance.
func NewDelivery() *Delivery {
	return &Delivery{
		storeURL:         "http://store:8080",
		deliveryAgentURL: "http://delivery-agent:8089",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		streamingClient: &http.Client{}, // No timeout for SSE streaming
	}
}

// NewDeliveryWithConfig creates a new Delivery instance with the given configuration.
func NewDeliveryWithConfig(config DeliveryConfig) *Delivery {
	d := NewDelivery()
	if config.StoreURL != "" {
		d.storeURL = config.StoreURL
	}
	if config.DeliveryAgentURL != "" {
		d.deliveryAgentURL = config.DeliveryAgentURL
	}
	return d
}

// HandleDeliver handles POST /deliver requests to deliver pizza orders.
// It validates the request and starts the delivery via the delivery agent asynchronously.
func (d *Delivery) HandleDeliver(w http.ResponseWriter, r *http.Request) {
	var req DeliverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate that at least one item is provided
	if len(req.OrderItems) == 0 {
		http.Error(w, "Order must contain at least one item", http.StatusBadRequest)
		return
	}

	slog.Info("delivery request received", "orderId", req.OrderID, "items", len(req.OrderItems))

	// Start delivery in a goroutine (background; detach from request context)
	go d.deliverOrder(context.Background(), req.OrderID)

	// Return accepted response immediately
	resp := DeliverResponse{
		OrderID: req.OrderID,
		Status:  "delivering",
		Message: fmt.Sprintf("Started delivering %d item(s)", len(req.OrderItems)),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

// deliverOrder calls the delivery-agent to deliver the order and forwards streaming updates to the store.
func (d *Delivery) deliverOrder(ctx context.Context, orderID uuid.UUID) {
	// Call delivery-agent with streaming to get updates
	updates := make(chan DeliveryUpdate, 100)
	go func() {
		err := d.callDeliveryAgentStream(ctx, orderID, updates)
		if err != nil {
			slog.Error("failed to deliver order via agent", "orderId", orderID, "error", err)
			d.sendEvent(ctx, orderID, fmt.Sprintf("failed to deliver: %v", err))
		}
	}()

	// Forward each update to the store
	var result string
	for update := range updates {
		// Only forward action events (not partial response tokens)
		if update.Type == "action" && update.Action != "" {
			slog.Info("delivery update", "orderId", orderID, "action", update.Action, "message", update.Message)
			d.sendDeliveryEvent(ctx, orderID, update)
		}
		if update.Type == "result" {
			result = update.Message
		}
	}

	slog.Info("delivery completed by agent", "orderId", orderID, "result", result)

	// Send DELIVERED event
	d.sendEvent(ctx, orderID, "DELIVERED")
}

// callDeliveryAgentStream sends a request to the delivery-agent streaming endpoint
// and sends DeliveryUpdate events to the provided channel.
func (d *Delivery) callDeliveryAgentStream(ctx context.Context, orderID uuid.UUID, updates chan<- DeliveryUpdate) error {
	defer close(updates)

	agentReq := AgentDeliverRequest{
		OrderID: orderID.String(),
	}

	body, err := json.Marshal(agentReq)
	if err != nil {
		return fmt.Errorf("failed to marshal agent request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.deliveryAgentURL+"/deliver/stream", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create agent request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := d.streamingClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call delivery agent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delivery agent returned status %d", resp.StatusCode)
	}

	// Parse SSE stream
	return d.parseSSEStream(resp.Body, updates)
}

// parseSSEStream parses a Server-Sent Events stream and sends updates to the channel.
// It handles both "data: value" and "data:value" formats, and multi-line SSE events
// that may include id:, event:, or other fields before the data: field.
func (d *Delivery) parseSSEStream(body interface{ Read([]byte) (int, error) }, updates chan<- DeliveryUpdate) error {
	buf := make([]byte, 4096)
	var dataBuffer bytes.Buffer

	for {
		n, err := body.Read(buf)
		if n > 0 {
			dataBuffer.Write(buf[:n])

			// Process complete SSE events (delimited by blank line)
			for {
				data := dataBuffer.Bytes()
				idx := bytes.Index(data, []byte("\n\n"))
				if idx == -1 {
					break
				}

				eventData := data[:idx]
				dataBuffer.Next(idx + 2)

				// Parse each line of the SSE event looking for data: fields
				lines := bytes.Split(eventData, []byte("\n"))
				for _, line := range lines {
					line = bytes.TrimRight(line, "\r") // Handle \r\n line endings
					if !bytes.HasPrefix(line, []byte("data:")) {
						continue
					}
					jsonData := line[5:] // Skip "data:"
					// Skip optional space after colon per SSE spec
					if len(jsonData) > 0 && jsonData[0] == ' ' {
						jsonData = jsonData[1:]
					}
					if len(jsonData) == 0 {
						continue
					}
					var update DeliveryUpdate
					if err := json.Unmarshal(jsonData, &update); err != nil {
						slog.Warn("failed to parse SSE event", "error", err, "data", string(jsonData))
						continue
					}
					updates <- update
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// sendDeliveryEvent sends a rich delivery update event to the store service.
func (d *Delivery) sendDeliveryEvent(ctx context.Context, orderID uuid.UUID, update DeliveryUpdate) {
	event := OrderEvent{
		OrderID:   orderID,
		Status:    update.Action,
		Source:    "delivery",
		Message:   update.Message,
		ToolName:  update.ToolName,
		ToolInput: update.ToolInput,
	}
	d.postEvent(ctx, event)
}

// sendEvent sends a simple status event to the store service.
func (d *Delivery) sendEvent(ctx context.Context, orderID uuid.UUID, status string) {
	event := OrderEvent{
		OrderID: orderID,
		Status:  status,
		Source:  "delivery",
	}
	d.postEvent(ctx, event)
}

// postEvent posts an OrderEvent to the store service.
func (d *Delivery) postEvent(ctx context.Context, event OrderEvent) {
	body, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal event", "orderId", event.OrderID, "error", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.storeURL+"/events", bytes.NewReader(body))
	if err != nil {
		slog.Error("failed to create event request", "orderId", event.OrderID, "error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		slog.Error("failed to send event to store", "orderId", event.OrderID, "status", event.Status, "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("store returned unexpected status for event", "orderId", event.OrderID, "status", event.Status, "httpStatus", resp.StatusCode)
	}
}
