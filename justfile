PROJECT_NAME := "aeternum"

# Default command
default:
    @just --list --unsorted

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

# Sync Go modules
tidy:
    go mod tidy
    cd api && go mod tidy
    cd execution && go mod tidy
    go work sync

# CLI run wrapper
cli *args:
    @go run main.go {{ args }}

# Build CLI binary
build-cli:
    #!/usr/bin/env bash
    echo "Building {{ PROJECT_NAME }} binary..."
    go mod download all
    VERSION=$(jq -r .version specs.json)
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.version=${VERSION}" -o ./aeternum main.go
    echo "Built binary for {{ PROJECT_NAME }} ${VERSION} successfully!"

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

# Start Compose with load-balancer
compose-up:
    docker compose -f compose.yaml up --build

# Stop all Compose containers and delete images created
compose-down:
    docker compose -f compose.yaml down
    docker rmi $(docker images | grep "{{ PROJECT_NAME }}" | awk "{print \$3}")

# Run a sample execution
test-sample-request:
    #!/usr/bin/env bash
    FULL_URL="${API_BASE_URL}/v0/tests/run"
    curl -X POST "$FULL_URL" \
        --header "Content-Type: application/json" \
        -d @sample/basic_request.json

# Run the docs server locally
docs:
    mkdocs build --strict --clean
    mkdocs serve --open
