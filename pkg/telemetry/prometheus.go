package telemetry

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	prom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// Define Prometheus metrics
var (
	// HTTP request metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTP request duration metrics
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Active requests gauge
	activeRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_requests",
			Help: "Number of requests currently being processed",
		},
	)

	// We don't need to define go_goroutines as it's already provided by the Prometheus client

	// Memory stats
	memoryAllocBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_memory_alloc_bytes",
			Help: "Number of bytes allocated and still in use",
		},
	)
)

// Initialize Prometheus metrics
func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(activeRequests)
	prometheus.MustRegister(memoryAllocBytes)

	// Start a goroutine to update runtime metrics
	go updateRuntimeMetrics()
}

// updateRuntimeMetrics periodically updates runtime metrics
func updateRuntimeMetrics() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var memStats runtime.MemStats

	for range ticker.C {
		// No need to update goroutines count as it's handled by the Prometheus client

		// Update memory stats
		runtime.ReadMemStats(&memStats)
		memoryAllocBytes.Set(float64(memStats.Alloc))
	}
}

// RecordHTTPRequest records metrics for an HTTP request
func RecordHTTPRequest(method, endpoint, status string, duration time.Duration) {
	httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// IncreaseActiveRequests increases the active requests counter
func IncreaseActiveRequests() {
	activeRequests.Inc()
}

// DecreaseActiveRequests decreases the active requests counter
func DecreaseActiveRequests() {
	activeRequests.Dec()
}

// setupPrometheusMetrics initializes the Prometheus meter provider
func setupPrometheusMetrics(ctx context.Context, cfg Config) (metric.MeterProvider, func(context.Context) error, error) {
	if !cfg.MetricsEnabled {
		// Return a no-op meter provider if metrics are disabled
		return otel.GetMeterProvider(), func(context.Context) error { return nil }, nil
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

	// Create Prometheus exporter
	exporter, err := prom.New()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	// Create meter provider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter),
	)

	// Set as global meter provider
	otel.SetMeterProvider(meterProvider)

	// Create global metrics
	meter := meterProvider.Meter("microservice")

	var metricErr error
	RequestCounter, metricErr = meter.Int64Counter(
		"request_count",
		metric.WithDescription("Number of requests received"),
	)
	if metricErr != nil {
		return nil, nil, fmt.Errorf("failed to create request counter: %w", metricErr)
	}

	RequestDuration, metricErr = meter.Float64Histogram(
		"request_duration",
		metric.WithDescription("Duration of requests in seconds"),
		metric.WithUnit("s"),
	)
	if metricErr != nil {
		return nil, nil, fmt.Errorf("failed to create request duration: %w", metricErr)
	}

	ActiveRequests, metricErr = meter.Int64UpDownCounter(
		"active_requests",
		metric.WithDescription("Number of requests currently being processed"),
	)
	if metricErr != nil {
		return nil, nil, fmt.Errorf("failed to create active requests: %w", metricErr)
	}

	// Start Prometheus HTTP server
	if cfg.MetricsPort > 0 {
		go func() {
			log.Printf("Starting Prometheus metrics server on :%d", cfg.MetricsPort)

			// Create a new ServeMux to avoid conflicts with the main HTTP server
			mux := http.NewServeMux()

			// Register the Prometheus HTTP handler
			mux.Handle(cfg.PrometheusPath, promhttp.Handler())

			// Start the HTTP server on the metrics port
			server := &http.Server{
				Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.MetricsPort),
				Handler: mux,
			}

			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("Error starting Prometheus metrics server: %v", err)
			}
		}()
	}

	// Return the provider and a shutdown function
	return meterProvider, func(context.Context) error {
		// Prometheus doesn't have a shutdown function
		return nil
	}, nil
}
