PROJECT_NAME := "aeternum"

# Default command
default:
    @just --list

# Run debug server
run-local port="8080":
    go run ./api/cmd/main.go --port={{port}} --dev=true

# Run production server
run-prod port:
    go run ./api/cmd/main.go --port {{port}} --dev=false

# Execute unit tests
test:
    @echo "Running unit tests!"
    go clean -testcache
    go test -cover ./api/...

# Build Docker image
build tag="latest": test
	@echo "Building Docker image (tag={{ tag }})..."
	docker build -t {{ PROJECT_NAME }}:{{ tag }} -f ./docker/server.Dockerfile .
	@echo "Docker image built successfully!"

# Sync Go modules
tidy:
    cd api && go mod tidy
    go work sync

# Start Compose with load-balancer
compose-up:
    docker compose -f compose.yaml up --build

# Stop all Compose containers and delete images created
compose-down:
    docker compose -f compose.yaml up down
    docker rmi $(docker images | grep "{{ PROJECT_NAME }}" | awk "{print \$3}")

# Run a sample execution
test-sample-request:
    #!/usr/bin/env bash
    API_HOST="localhost"
    PORT="8080"
    ENDPOINT="v0/tests/run"
    LOCAL_URL="http://${API_HOST}:${PORT}/${ENDPOINT}"
    curl -vX POST "$LOCAL_URL" \
        --header "Content-Type: application/json" \
        -d @sample/basic_request.json
