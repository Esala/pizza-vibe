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
```

## Architecture

The Pizza application is composed by three services written in Go: 
- Store service which exposes the APIs that will be consumed by the front-end. This service acts as the orchestrator for 
    pizza orders between the Kitchen and Delivery Services.
- Kitchen service which will be responsible cooking the pizzas. 
- Delivery service which will be responsible for the delivery of the pizza to the customer.

The Store service must use Dapr Workflows to orchestrate the pizza order flow.
The Kitchen and delivery services must use Dapr Pub/Sub to provide updates to the Store service.