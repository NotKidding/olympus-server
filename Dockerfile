# --- Stage 1: Build the Go binary ---
FROM golang:1.26-alpine AS builder

# Install git and certs (standard practice)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy dependency files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a static binary (CGO_ENABLED=0 is key for Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o olympus ./cmd/olympus/main.go

# --- Stage 2: Final Runtime Image ---
FROM alpine:latest

WORKDIR /app

# Copy only the compiled binary from the builder
COPY --from=builder /app/olympus .

# Expose the C2 port and the gRPC management port
EXPOSE 8080 9090

# Run the Teamserver
CMD ["./olympus"]