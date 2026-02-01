# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Pizza Vibe is a Go-based pizza application ("vibecoded").

## Build and Run Commands

```bash
# Build the application
go build -o pizza-vibe ./...

# Run the application
go run .

# Run tests
go test ./...

# Run a specific test
go test -run TestName ./path/to/package

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run docker-compose validation tests
./scripts/test-docker-compose.sh

# Build and run with docker-compose
docker-compose build
docker-compose up -d
docker-compose down
```

## Architecture

The Pizza application is composed by three services written in Go: 
- Store service which exposes the APIs that will be consumed by the front-end. This service acts as the orchestrator for 
    pizza orders between the Kitchen and Delivery Services.
- Kitchen service which will be responsible cooking the pizzas. 
- Delivery service which will be responsible for the delivery of the pizza to the customer.

The Store service must use Dapr Workflows to orchestrate the pizza order flow.
The Kitchen and delivery services must use Dapr Pub/Sub to provide updates to the Store service.

## Best practices

General: 
- Do not do more than what is asked for
Frontend:
- Everytime that you send a request to the store service validate the data types to make sure that the request is valid.
- Use the store service data types (@store/models.go) to create mock data for the jest tests.
- Always use Fetch to call other services using http.
- Do not add styles unless it is specified by the user.
- When creating content in pages, only add what is explicitly requested or ask if recommending additional content is needed.
- Never add styles unless specifically requested by the user.

Backend:
- Always keep update the docker-compose.yaml file with all the services of the application.
- Run `./scripts/test-docker-compose.sh` to validate docker-compose changes before committing.
- Always provide Kubernetes manifests for each service and infrastructure component.
- Always implement Dockerfile for each service
