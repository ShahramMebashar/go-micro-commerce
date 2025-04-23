package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// setupJaegerTracing initializes the OTLP tracer provider for Jaeger
func setupJaegerTracing(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, func(context.Context) error, error) {
	if !cfg.TracingEnabled {
		// Return a no-op tracer if tracing is disabled
		return sdktrace.NewTracerProvider(), func(context.Context) error { return nil }, nil
	}

	// Create a resource describing the service
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create OTLP exporter for Jaeger
	var endpoint string
	if cfg.OTLPEndpoint != "" {
		endpoint = cfg.OTLPEndpoint
	} else {
		endpoint = "jaeger:4317" // Default OTLP gRPC endpoint for Jaeger
	}

	// Create gRPC connection to the collector
	conn, err := grpc.DialContext(ctx, endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Create OTLP exporter
	exporter, err := otlptrace.New(ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithGRPCConn(conn),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// Configure the trace provider
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// Return the provider and a shutdown function
	return tracerProvider, func(ctx context.Context) error {
		return tracerProvider.Shutdown(ctx)
	}, nil
}
