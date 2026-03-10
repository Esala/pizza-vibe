# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Pizza Vibe is agentic pizza store, which uses Langchain4j, Quarkus and Dapr Workflows to provide a seamless experience for customers.

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

The Pizza application is composed of: 
- Five services written in Go: bikes, drinks-stock, inventory, oven, store 
- Three agents written in Java with Quarkus: cooking-agent, delivery-agent, store-mgmt-agent
- MCP Server using Quarkus: pizza-mcp 

The front-end is written in React and the back-end is written in Go.



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
