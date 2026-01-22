# pizza-vibe
Vibecoded Pizza App in Go

## Architecture

The Pizza Vibe application is composed of three services written in Go:

- **Store Service** (port 8080): Exposes the APIs consumed by the front-end. Acts as the orchestrator for pizza orders between the Kitchen and Delivery Services.
- **Kitchen Service** (port 8081): Responsible for cooking the pizzas. Each item takes a random time from 1-10 seconds to cook.
- **Delivery Service** (coming soon): Responsible for the delivery of the pizza to the customer.

## Running the Services

### Using Docker Compose (Recommended)

The easiest way to run all services is using Docker Compose:

```bash
# Build and start all services
docker-compose up --build

# Run in detached mode
docker-compose up --build -d

# Stop all services
docker-compose down
```

### Running Locally

#### Prerequisites
- Go 1.24 or later

#### Store Service
```bash
go run ./store/cmd/
```
The store service will start on port 8080.

#### Kitchen Service
```bash
go run ./kitchen/cmd/
```
The kitchen service will start on port 8081.

## API Endpoints

### Store Service (port 8080)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/order` | POST | Create a new pizza order |
| `/events` | POST | Receive events from kitchen/delivery |
| `/ws` | GET | WebSocket for real-time order updates |
| `/health` | GET | Health check endpoint |

### Kitchen Service (port 8081)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/cook` | POST | Cook order items |
| `/health` | GET | Health check endpoint |

#### Example: Cook Request
```bash
curl -X POST http://localhost:8081/cook \
  -H "Content-Type: application/json" \
  -d '{
    "orderId": "550e8400-e29b-41d4-a716-446655440000",
    "orderItems": [
      {"pizzaType": "Margherita", "quantity": 2},
      {"pizzaType": "Pepperoni", "quantity": 1}
    ]
  }'
```

## Development

### Build
```bash
go build -o pizza-vibe ./...
```

### Run Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test -coverprofile=coverage.out ./...
```

### Format Code
```bash
go fmt ./...
```

### Vet Code
```bash
go vet ./...
```
