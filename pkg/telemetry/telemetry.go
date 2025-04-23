package telemetry

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Config holds configuration for telemetry components
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	// Tracing config
	TracingEnabled bool
	OTLPEndpoint   string
	JaegerEndpoint string
	// Metrics config
	MetricsEnabled bool
	MetricsPort    int
	PrometheusPath string // Path for Prometheus metrics endpoint (default: /metrics)
	// Logging config
	LogLevel string
}

// DefaultConfig returns a configuration with values from environment variables
func DefaultConfig() Config {
	return Config{
		ServiceName:    getEnv("SERVICE_NAME", "unknown-service"),
		ServiceVersion: getEnv("SERVICE_VERSION", "0.0.1"),
		Environment:    getEnv("ENV", "development"),
		TracingEnabled: getEnvAsBool("TELEMETRY_ENABLED", true),
		OTLPEndpoint:   getEnv("TELEMETRY_OTLP_ENDPOINT", "jaeger:4317"),
		JaegerEndpoint: getEnv("TELEMETRY_JAEGER_ENDPOINT", "http://jaeger:14268/api/traces"),
		MetricsEnabled: getEnvAsBool("TELEMETRY_METRICS_ENABLED", true),
		MetricsPort:    getEnvAsInt("TELEMETRY_METRICS_PORT", 9090),
		PrometheusPath: getEnv("TELEMETRY_PROMETHEUS_PATH", "/metrics"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}
}

// Helper functions to read environment variables
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}

func Setup(ctx context.Context, cfg Config) (shutdown func(context.Context) error, err error) {
	// Create propagator
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	// Initialize tracer provider
	tracerProvider, tracerShutdown, err := setupJaegerTracing(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to setup tracing: %w", err)
	}
	otel.SetTracerProvider(tracerProvider)

	// Initialize meter provider
	_, metricsShutdown, err := setupPrometheusMetrics(ctx, cfg)
	if err != nil {
		// Clean up tracer if metrics setup fails
		_ = tracerShutdown(ctx)
		return nil, fmt.Errorf("failed to setup metrics: %w", err)
	}

	// Initialize logger
	if err := setupLogging(cfg); err != nil {
		// Clean up tracer and metrics if logging setup fails
		_ = tracerShutdown(ctx)
		_ = metricsShutdown(ctx)
		return nil, fmt.Errorf("failed to setup logging: %w", err)
	}

	// Return a shutdown function that closes all telemetry components
	return func(ctx context.Context) error {
		// Use a timeout for shutdown
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var shutdownErr error
		if err := tracerShutdown(ctx); err != nil {
			shutdownErr = fmt.Errorf("failed to shutdown tracer provider: %w", err)
			log.Printf("Error shutting down tracer provider: %v", err)
		}

		if err := metricsShutdown(ctx); err != nil {
			shutdownErr = fmt.Errorf("failed to shutdown meter provider: %w", err)
			log.Printf("Error shutting down meter provider: %v", err)
		}

		return shutdownErr
	}, nil
}
