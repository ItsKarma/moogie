# Multi-stage build for Moogie Runner
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the runner binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o runner .

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/runner /app/runner

# Run as non-root user for security
RUN addgroup -g 1000 runner && \
    adduser -D -u 1000 -G runner runner && \
    chown -R runner:runner /app

USER runner

ENTRYPOINT ["/app/runner"]
