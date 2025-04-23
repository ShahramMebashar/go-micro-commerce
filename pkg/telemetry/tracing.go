package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// Tracer returns the global tracer
func Tracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer("microservice")
}

// setupTracing initializes the tracer provider
func setupTracing(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, func(context.Context) error, error) {
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

	// Create stdout exporter
	var exporter sdktrace.SpanExporter
	if cfg.OTLPEndpoint == "stdout" {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	} else {
		// In a real implementation, you would use OTLP exporter here
		// For simplicity, we're using stdout exporter for all cases
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create trace exporter: %w", err)
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
