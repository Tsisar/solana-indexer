# Application name and registry configuration
APP := solana-indexer.vaults
REGISTRY := intothefathom
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)-$(shell git rev-parse --short HEAD)

# Build configuration
TARGETOS ?= darwin
TARGETARCH ?= arm64
CGO_ENABLED ?= 0

# Phony targets
.PHONY: help format lint test get build image push clean dev release

# Help target
help:
	@echo "Available targets:"
	@echo "  help     - Show this help message"
	@echo "  format   - Format Go code"
	@echo "  lint     - Run golangci-lint"
	@echo "  test     - Run tests"
	@echo "  get      - Get dependencies"
	@echo "  build    - Build the application"
	@echo "  image    - Build Docker image (multi-platform)"
	@echo "  push     - Push Docker image to registry"
	@echo "  clean    - Clean build artifacts"
	@echo "  dev      - Development build (with debug info)"
	@echo "  release  - Production release build and push"
	@echo ""
	@echo "Configuration:"
	@echo "  TARGETOS   - Target OS (linux, darwin, windows) [$(TARGETOS)]"
	@echo "  TARGETARCH - Target architecture (amd64, arm64) [$(TARGETARCH)]"
	@echo "  CGO_ENABLED - Enable CGO (0 or 1) [$(CGO_ENABLED)]"

# Format Go code
format:
	@echo "Formatting Go code..."
	@gofmt -s -w ./

# Install golangci-lint if not present
.PHONY: install-lint
install-lint:
	@which golangci-lint >/dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)

# Run linter
lint: install-lint
	@echo "Running linter..."
	@golangci-lint run ./...

# Run tests
test:
	@echo "Running tests..."
	@go test -v -cover ./...

# Get dependencies
get:
	@echo "Getting dependencies..."
	@go mod tidy
	@go mod download

# Development build
dev: format get
	@echo "Building development binary..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) \
		go build -v -o indexer ./cmd/indexer

# Production build
build: format get
	@echo "Building production binary..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) \
		go build -v -o indexer -ldflags="-s -w" ./cmd/indexer

# Build multi-platform Docker image (no push)
image:
	@echo "Building and pushing multi-platform Docker image..."
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg TARGETARCH=$(TARGETARCH) \
		--build-arg TARGETOS=$(TARGETOS) \
		--tag $(REGISTRY)/$(APP):$(VERSION) \
		--push .

# Clean artifacts
clean:
	@echo "Cleaning artifacts..."
	@rm -f indexer
	@docker rmi $(REGISTRY)/$(APP):$(VERSION) || true

# Full release pipeline
release: clean test build image
	@echo "Release $(VERSION) completed successfully"