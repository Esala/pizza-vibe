# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the delivery service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/delivery ./delivery/cmd/

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/delivery .

# Expose port
EXPOSE 8082

# Run the service
CMD ["./delivery"]
