#!/bin/bash

# Docker Compose Validation Tests
# This script validates that the docker-compose.yaml works correctly

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
COMPOSE_FILE="$PROJECT_DIR/docker-compose.yaml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0

# Cleanup function
cleanup() {
    echo -e "\n${YELLOW}Cleaning up...${NC}"
    cd "$PROJECT_DIR"
    docker-compose down --volumes --remove-orphans 2>/dev/null || true
}

# Set trap to cleanup on exit
trap cleanup EXIT

# Helper function to print test results
pass() {
    echo -e "${GREEN}PASS${NC}: $1"
    TESTS_PASSED=$((TESTS_PASSED + 1))
}

fail() {
    echo -e "${RED}FAIL${NC}: $1"
    TESTS_FAILED=$((TESTS_FAILED + 1))
}

# Test 1: Validate docker-compose.yaml syntax
test_syntax() {
    echo -e "\n${YELLOW}Test: Validating docker-compose.yaml syntax...${NC}"
    cd "$PROJECT_DIR"

    if docker-compose config > /dev/null 2>&1; then
        pass "docker-compose.yaml syntax is valid"
    else
        fail "docker-compose.yaml syntax is invalid"
        docker-compose config 2>&1
        return 1
    fi
}

# Test 2: Verify all required services are defined
test_required_services() {
    echo -e "\n${YELLOW}Test: Checking required services are defined...${NC}"
    cd "$PROJECT_DIR"

    local services=$(docker-compose config --services 2>/dev/null)
    local required_services=("store" "kitchen")

    for service in "${required_services[@]}"; do
        if echo "$services" | grep -q "^${service}$"; then
            pass "Service '$service' is defined"
        else
            fail "Service '$service' is not defined"
        fi
    done
}

# Test 3: Verify Dockerfiles exist for all services
test_dockerfiles_exist() {
    echo -e "\n${YELLOW}Test: Checking Dockerfiles exist...${NC}"
    cd "$PROJECT_DIR"

    local services=$(docker-compose config --services 2>/dev/null)

    for service in $services; do
        local dockerfile="$PROJECT_DIR/$service/Dockerfile"
        if [[ -f "$dockerfile" ]]; then
            pass "Dockerfile exists for '$service'"
        else
            fail "Dockerfile missing for '$service' (expected at $dockerfile)"
        fi
    done
}

# Test 4: Build all services
test_build_services() {
    echo -e "\n${YELLOW}Test: Building all services...${NC}"
    cd "$PROJECT_DIR"

    if docker-compose build --no-cache 2>&1; then
        pass "All services built successfully"
    else
        fail "Failed to build services"
        return 1
    fi
}

# Test 5: Start services and verify they're running
test_services_start() {
    echo -e "\n${YELLOW}Test: Starting services...${NC}"
    cd "$PROJECT_DIR"

    docker-compose up -d 2>&1

    # Wait for services to start
    sleep 5

    local services=$(docker-compose config --services 2>/dev/null)

    for service in $services; do
        local status=$(docker-compose ps --format json 2>/dev/null | jq -r ".[] | select(.Service == \"$service\") | .State" 2>/dev/null || docker-compose ps 2>/dev/null | grep "$service" | awk '{print $4}')
        if [[ "$status" == "running" ]] || [[ "$status" == "Up" ]] || docker-compose ps 2>/dev/null | grep "$service" | grep -q "Up"; then
            pass "Service '$service' is running"
        else
            fail "Service '$service' is not running (status: $status)"
        fi
    done
}

# Test 6: Verify health endpoints respond
test_health_endpoints() {
    echo -e "\n${YELLOW}Test: Checking health endpoints...${NC}"
    cd "$PROJECT_DIR"

    # Wait for health checks to pass
    local max_attempts=30
    local attempt=0

    echo "Waiting for services to become healthy..."

    # Test store health endpoint
    attempt=0
    while [[ $attempt -lt $max_attempts ]]; do
        if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
            pass "Store service health endpoint responding"
            break
        fi
        ((attempt++))
        sleep 1
    done
    if [[ $attempt -eq $max_attempts ]]; then
        fail "Store service health endpoint not responding after ${max_attempts}s"
    fi

    # Test kitchen health endpoint
    attempt=0
    while [[ $attempt -lt $max_attempts ]]; do
        if curl -sf http://localhost:8081/health > /dev/null 2>&1; then
            pass "Kitchen service health endpoint responding"
            break
        fi
        ((attempt++))
        sleep 1
    done
    if [[ $attempt -eq $max_attempts ]]; then
        fail "Kitchen service health endpoint not responding after ${max_attempts}s"
    fi
}

# Test 7: Verify services expose correct ports
test_port_mappings() {
    echo -e "\n${YELLOW}Test: Checking port mappings...${NC}"
    cd "$PROJECT_DIR"

    local expected_ports=(
        "store:8080"
        "kitchen:8081"
    )

    for mapping in "${expected_ports[@]}"; do
        local service="${mapping%%:*}"
        local port="${mapping##*:}"

        if docker-compose ps 2>/dev/null | grep "$service" | grep -q "0.0.0.0:$port->"; then
            pass "Service '$service' exposes port $port"
        else
            fail "Service '$service' does not expose port $port correctly"
        fi
    done
}

# Test 8: Check service logs for errors
test_no_startup_errors() {
    echo -e "\n${YELLOW}Test: Checking for startup errors in logs...${NC}"
    cd "$PROJECT_DIR"

    local services=$(docker-compose config --services 2>/dev/null)

    for service in $services; do
        local errors=$(docker-compose logs "$service" 2>&1 | grep -i "error\|panic\|fatal" | grep -v "level=error" || true)
        if [[ -z "$errors" ]]; then
            pass "No startup errors in '$service' logs"
        else
            fail "Found errors in '$service' logs:"
            echo "$errors"
        fi
    done
}

# Test 9: Verify network connectivity between services
test_network_connectivity() {
    echo -e "\n${YELLOW}Test: Checking network connectivity between services...${NC}"
    cd "$PROJECT_DIR"

    # Check if store can reach kitchen
    if docker-compose exec -T store wget -q --spider http://kitchen:8081/health 2>/dev/null; then
        pass "Store can reach Kitchen service"
    else
        fail "Store cannot reach Kitchen service"
    fi

    # Check if kitchen can reach store
    if docker-compose exec -T kitchen wget -q --spider http://store:8080/health 2>/dev/null; then
        pass "Kitchen can reach Store service"
    else
        fail "Kitchen cannot reach Store service"
    fi
}

# Main function
main() {
    echo "========================================"
    echo "Docker Compose Validation Tests"
    echo "========================================"
    echo "Project: $PROJECT_DIR"
    echo "Compose file: $COMPOSE_FILE"

    # Run tests
    test_syntax
    test_required_services
    test_dockerfiles_exist
    test_build_services
    test_services_start
    test_health_endpoints
    test_port_mappings
    test_no_startup_errors
    test_network_connectivity

    # Print summary
    echo -e "\n========================================"
    echo "Test Summary"
    echo "========================================"
    echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
    echo -e "${RED}Failed: $TESTS_FAILED${NC}"

    if [[ $TESTS_FAILED -gt 0 ]]; then
        echo -e "\n${RED}Some tests failed!${NC}"
        exit 1
    else
        echo -e "\n${GREEN}All tests passed!${NC}"
        exit 0
    fi
}

# Run main function
main "$@"
