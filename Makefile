# ==============================
# Build configuration
# ==============================

# Target platform (override via `make TARGETOS=linux`)
TARGETOS     ?= linux
TARGETARCH   ?= amd64
CGO_ENABLED  ?= 0

# Application name and registry
APP      := solana-indexer.vaults
REGISTRY := intothefathom
VERSION  := $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)-$(shell git rev-parse --short HEAD)

# ==============================
# Phony targets
# ==============================
.PHONY: help format lint test get build image push clean dev release install-lint

# ==============================
# Help
# ==============================
help:
	@echo "Available targets:"
	@echo "  help     - Show this help message"
	@echo "  format   - Format Go code"
	@echo "  lint     - Run golangci-lint"
	@echo "  test     - Run tests"
	@echo "  get      - Get dependencies"
	@echo "  build    - Build the application"
	@echo "  image    - Build Docker image"
	@echo "  push     - Push Docker image to registry"
	@echo "  clean    - Clean build artifacts"
	@echo "  dev      - Development build (with debug info)"
	@echo "  release  - Production release build and push"
	@echo ""
	@echo "Configuration:"
	@echo "  TARGETOS     = $(TARGETOS) (linux, windows, darwin)"
	@echo "  TARGETARCH   = $(TARGETARCH) (amd64, arm64)"
	@echo "  CGO_ENABLED  = $(CGO_ENABLED) (0 for no cgo, 1 for cgo enabled)"

# ==============================
# Go Tools
# ==============================
format:
	@echo "Formatting Go code..."
	@gofmt -s -w ./

install-lint:
	@which golangci-lint >/dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)

lint: install-lint
	@echo "Running linter..."
	@golangci-lint run ./...

test:
	@echo "Running tests..."
	@go test -v -cover ./...

get:
	@echo "Getting dependencies..."
	@go mod tidy
	@go mod download

# ==============================
# Build targets
# ==============================
dev: format get
	@echo "Building development binary..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) \
		go build -v -o indexer ./cmd/indexer

build: format get
	@echo "Building production binary..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) \
		go build -v -o indexer -ldflags="-s -w" ./cmd/indexer

# ==============================
# Docker
# ==============================
image:
	@echo "Building Docker image for $(TARGETOS)/$(TARGETARCH)..."
	@docker buildx build \
		--platform $(TARGETOS)/$(TARGETARCH) \
		--build-arg TARGETOS=$(TARGETOS) \
		--build-arg TARGETARCH=$(TARGETARCH) \
		--tag $(REGISTRY)/$(APP):$(VERSION) \
		--load .

push:
	@echo "Pushing image to registry..."
	@docker push $(REGISTRY)/$(APP):$(VERSION)

# ==============================
# Utilities
# ==============================
clean:
	@echo "Cleaning artifacts..."
	@rm -f indexer
	@docker rmi $(REGISTRY)/$(APP):$(VERSION) || true

release: clean test build image push
	@echo "Release $(VERSION) completed successfully"