# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the store service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/store ./store/cmd/

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/store .

# Expose port
EXPOSE 8080

# Run the service
CMD ["./store"]
