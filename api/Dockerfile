# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Final stage
FROM alpine:latest

# Install ca-certificates and postgresql-client for migrations
RUN apk --no-cache add ca-certificates postgresql-client

# Create non-root user
RUN addgroup -g 1001 -S moogie && \
    adduser -u 1001 -S moogie -G moogie

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bin/server .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy migration files and entrypoint
COPY migrations/ ./migrations/
COPY entrypoint.sh .

# Create directories and set permissions
RUN mkdir -p /app/logs && \
    chown -R moogie:moogie /app && \
    chmod +x entrypoint.sh

# Switch to non-root user
USER moogie

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Use entrypoint script to run migrations and start server
ENTRYPOINT ["./entrypoint.sh"]
