FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o product-service ./services/product-service/cmd/api

# Use a minimal alpine image for the final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/product-service .

# Copy migrations
COPY --from=builder /app/services/product-service/internal/migrations ./services/product-service/internal/migrations

# Expose the ports
EXPOSE 8080
EXPOSE 9091

# Run the application
CMD ["./product-service"]
