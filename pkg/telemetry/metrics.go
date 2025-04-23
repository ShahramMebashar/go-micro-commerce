package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// Global metrics
var (
	RequestCounter  metric.Int64Counter
	RequestDuration metric.Float64Histogram
	ActiveRequests  metric.Int64UpDownCounter
)

// setupMetrics initializes the meter provider and global metrics
func setupMetrics(_ context.Context, cfg Config) (metric.MeterProvider, func(context.Context) error, error) {
	if !cfg.MetricsEnabled {
		// Return a no-op meter provider if metrics are disabled
		return otel.GetMeterProvider(), func(context.Context) error { return nil }, nil
	}

	// In a real implementation, you would create a proper meter provider with exporters
	// For simplicity, we're using the global no-op meter provider
	meterProvider := otel.GetMeterProvider()

	// Create global metrics
	meter := meterProvider.Meter("microservice")

	var err error
	RequestCounter, err = meter.Int64Counter(
		"request_count",
		metric.WithDescription("Number of requests received"),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request counter: %w", err)
	}

	RequestDuration, err = meter.Float64Histogram(
		"request_duration",
		metric.WithDescription("Duration of requests in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request duration: %w", err)
	}

	ActiveRequests, err = meter.Int64UpDownCounter(
		"active_requests",
		metric.WithDescription("Number of requests currently being processed"),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create active requests: %w", err)
	}

	// Return the provider and a no-op shutdown function
	return meterProvider, func(context.Context) error { return nil }, nil
}
