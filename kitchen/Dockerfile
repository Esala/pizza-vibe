# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the kitchen service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/kitchen ./kitchen/cmd/

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/kitchen .

# Expose port
EXPOSE 8081

# Run the service
CMD ["./kitchen"]
