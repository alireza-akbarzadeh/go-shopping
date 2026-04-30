# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/shopping-platform-api ./cmd/api

# Stage 2: Create a minimal runtime image
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS calls (if needed)
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/bin/shopping-platform-api /app/shopping-platform-api

# Copy .env.example (optional, but we'll use environment variables)
COPY .env.example .env

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["/app/shopping-platform-api"]