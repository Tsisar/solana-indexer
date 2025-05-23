#!/bin/bash

set -e

# Name of the Docker image
IMAGE_NAME="intothefathom/solana-indexer.vaults"

# Get the latest Git tag
TAG=$(git describe --tags --abbrev=0)

# Function to print header
function print_header() {
    echo "======================================"
    echo "$1"
    echo "======================================"
}

# Check if docker buildx is available
if ! docker buildx version > /dev/null 2>&1; then
    echo "Docker buildx is not installed or available. Please install Docker buildx."
    exit 1
fi

# Create and use a new buildx builder if not already exists
if ! docker buildx inspect mybuilder >/dev/null 2>&1; then
    print_header "Creating and using new buildx builder"
    docker buildx create --use --name mybuilder
else
    print_header "Using existing buildx builder"
    docker buildx use mybuilder
fi

# Inspect and bootstrap the builder
docker buildx inspect mybuilder --bootstrap

# Build and push the image for multiple platforms
print_header "Building and pushing Docker image ${IMAGE_NAME}:${TAG}-dev"
docker buildx build --platform linux/amd64,linux/arm64 -t ${IMAGE_NAME}:${TAG}-dev --push .

# Clean up builder
print_header "Cleaning up buildx builder"
docker buildx rm mybuilder

# Remove local container with the same name, if exists
CONTAINER_ID=$(docker ps -aqf "name=$(basename $IMAGE_NAME)")

if [ -n "$CONTAINER_ID" ]; then
    print_header "Removing existing local container $(basename $IMAGE_NAME)"
    docker rm -f ${CONTAINER_ID}
fi

print_header "Docker image ${IMAGE_NAME}:${TAG}-dev has been built and pushed successfully"