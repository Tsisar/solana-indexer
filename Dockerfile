# Stage 1: Build the Go application with cross-compilation
FROM --platform=$BUILDPLATFORM golang:1.24.2-alpine AS builder

# Build-time arguments provided automatically by BuildKit
ARG TARGETOS
ARG TARGETARCH

# Install necessary packages for CGO and timezone
RUN apk add --no-cache gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Cross-compile the binary for the target platform
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o /indexer ./cmd/indexer

# Stage 2: Minimal runtime image
FROM alpine:latest

# Install CA certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy compiled binary
COPY --from=builder /indexer .

EXPOSE 9040

ENV METRICS_ENABLED="true" \
    METRICS_PORT="8040" \
    SAVE_LOGS="false" \
    LOG_LEVEL="debug"

ENTRYPOINT ["./indexer"]