# Stage 1: Build the Go application
FROM golang:1.24.2-alpine AS builder

# Install necessary packages for CGO and timezone
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build only the indexer binary from cmd/indexer
RUN go build -ldflags="-s -w" -o /indexer ./cmd/indexer

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Install CA certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /root/

# Copy the compiled Go binary
COPY --from=builder /indexer .

# Expose metrics port (optional)
EXPOSE 9040

# Default environment variables
ENV METRICS_ENABLED="true" \
    METRICS_PORT="9040" \
    SAVE_LOGS="false" \
    LOG_LEVEL="debug"

# Run the binary
ENTRYPOINT ["./indexer"]