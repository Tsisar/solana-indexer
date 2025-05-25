#!/bin/bash

set -euo pipefail

# Name of the Docker image
IMAGE_NAME="intothefathom/solana-indexer.vaults"

# Get the latest Git tag (fallback to "latest-dev" if not found)
TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "latest-dev")

# Function to print header
print_header() {
    echo "======================================"
    echo "$1"
    echo "======================================"
}

# Ensure buildx is available
if ! docker buildx version > /dev/null 2>&1; then
    echo "[ERROR] Docker buildx is not installed or available. Please install it."
    exit 1
fi

# Create and use buildx builder if it doesn't exist
if ! docker buildx inspect mybuilder >/dev/null 2>&1; then
    print_header "Creating and using new buildx builder"
    docker buildx create --use --name mybuilder
else
    print_header "Using existing buildx builder"
    docker buildx use mybuilder
fi

# Ensure builder is bootstrapped
docker buildx inspect mybuilder --bootstrap

# Build and push multi-arch image
print_header "Building and pushing Docker image ${IMAGE_NAME}:${TAG}-dev"
docker buildx build \
    --platform linux/amd64,linux/arm64 \
    -t ${IMAGE_NAME}:${TAG}-dev \
    --push .

# Optional: remove the buildx builder (comment out if you want to reuse)
print_header "Cleaning up buildx builder"
docker buildx rm mybuilder || true

# Optional: remove local container with the same name
CONTAINER_ID=$(docker ps -aqf "name=$(basename $IMAGE_NAME)")

if [ -n "$CONTAINER_ID" ]; then
    print_header "Removing existing local container $(basename $IMAGE_NAME)"
    docker rm -f ${CONTAINER_ID}
fi

print_header "Docker image ${IMAGE_NAME}:${TAG}-dev has been built and pushed successfully"