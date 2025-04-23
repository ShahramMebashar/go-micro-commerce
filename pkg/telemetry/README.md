# Telemetry Package

This package provides observability features for your microservices:

- **Distributed Tracing**: Track requests as they flow through your services
- **Metrics**: Measure request counts, durations, and active requests
- **Logging with Context**: Include trace IDs in your logs for correlation

## Setup

### 1. Initialize the Telemetry Package

In your `main.go` file:

```go
import (
    "context"
    "log"
    "microservice/pkg/telemetry"
)

func main() {
    // Create telemetry configuration
    telemetryCfg := telemetry.Config{
        ServiceName:    "product-service",
        ServiceVersion: "1.0.0",
        Environment:    "development",
        TracingEnabled: true,
        OTLPEndpoint:   "jaeger:4317",  // OTLP gRPC endpoint for Jaeger
        MetricsEnabled: true,
        MetricsPort:    8090,           // Port for metrics endpoint
        PrometheusPath: "/metrics",     // Path for Prometheus metrics
        LogLevel:       "info",
    }
    
    // Initialize telemetry
    ctx := context.Background()
    shutdown, err := telemetry.Setup(ctx, telemetryCfg)
    if err != nil {
        log.Fatalf("Failed to setup telemetry: %v", err)
    }
    defer shutdown(ctx)
    
    // Rest of your application...
}
```

### 2. Add Middleware to Your Router

```go
import (
    "microservice/pkg/telemetry"
    "github.com/go-chi/chi/v5"
)

func setupRouter() *chi.Mux {
    r := chi.NewRouter()
    
    // Add telemetry middleware
    r.Use(telemetry.Middleware)
    
    // Add your routes
    r.Get("/products", productHandler.ListProducts)
    
    return r
}
```

### 3. Start Jaeger and Prometheus

Use the provided Docker Compose file:

```bash
docker-compose -f docker-compose-telemetry.yml up -d
```

## Usage

### Tracing

```go
func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
    // Create a span for this operation
    ctx, span := telemetry.Tracer().Start(ctx, "ProductService.GetByID")
    defer span.End()
    
    // Add attributes to the span
    span.SetAttributes(attribute.String("product.id", id.String()))
    
    // Call the repository
    product, err := s.repo.GetByID(ctx, id)
    
    // Record error if any
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }
    
    return product, nil
}
```

### Logging

```go
// Instead of:
log.Printf("Getting product with ID: %s", id)

// Use:
telemetry.LogWithContext(ctx, "Getting product with ID: %s", id)
```

### Metrics

```go
// Record a count of business operations
telemetry.RequestCounter.Add(ctx, 1)

// Record the duration of an operation
startTime := time.Now()
// ... perform operation
duration := time.Since(startTime).Seconds()
telemetry.RequestDuration.Record(ctx, duration)

// Track concurrent operations
telemetry.ActiveRequests.Add(ctx, 1)
// ... perform operation
telemetry.ActiveRequests.Add(ctx, -1)
```

## Viewing Telemetry Data

- **Jaeger UI**: http://localhost:16686
- **Prometheus UI**: http://localhost:9090

## Best Practices

1. Always propagate context through your call chain
2. Use descriptive span names
3. Add relevant attributes to spans
4. Record errors in spans
5. Always defer span.End()
6. Use LogWithContext for logging
