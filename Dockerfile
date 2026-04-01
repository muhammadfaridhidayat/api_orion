# ---- Build Stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install ca-certificates for HTTPS and tzdata for timezone support
RUN apk add --no-cache ca-certificates tzdata

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with optimizations
# CGO_ENABLED=0 for a fully static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/api_orion \
    .

# ---- Runtime Stage ----
FROM alpine:3.21

WORKDIR /app

# Install ca-certificates (needed for HTTPS/TLS connections to external services like Postgres with SSL)
RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the compiled binary from the builder stage
COPY --from=builder /app/api_orion .

# Set ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Set default environment variables (can be overridden by Dokploy)
ENV PORT=8080

# Run the binary
CMD ["./api_orion"]
