PROJECT_NAME := "aeternum"

# List all available commands
_default:
    @just --list --unsorted

# Run debug server
run-local port="8080":
    go run . --port={{ port }} --dev=true

# Run production server
run-prod port="8080":
    go run . --port {{ port }} --dev=false

# Execute unit tests
test:
    @echo "Running unit tests!"
    go clean -testcache
    go test -cover ./...

# Sync Go modules
tidy:
    @go mod tidy
    @echo "Go modules synced successfully!"

# Build Docker image manually and push to K8s server
build-k8s-deployment tag="latest":
    #!/usr/bin/env bash
    eval $(minikube docker-env)
    echo "Using Minikube Docker environment."
    IMAGE_NAME="{{ PROJECT_NAME }}-api"
    docker build \
        --no-cache \
        -f Dockerfile \
        -t "${IMAGE_NAME}:{{ tag }}" .
    echo "Docker image built successfully!"
    docker images | grep "$IMAGE_NAME"

# Start Docker Compose with load-balancer
compose-up:
    docker compose -f compose.yaml up --build

# Stop all Docker Compose containers and delete images created
compose-down:
    docker compose -f compose.yaml down
    docker rmi $(docker images | grep "{{ PROJECT_NAME }}" | awk "{print \$3}")

# Run the docs server locally
docs:
    mkdocs build --strict --clean
    mkdocs serve --open
