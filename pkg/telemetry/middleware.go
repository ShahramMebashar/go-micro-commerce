package telemetry

import (
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ResponseWriter is a custom response writer that captures the status code and size
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Size       int
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the response size
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.Size += size
	return size, err
}

// Middleware returns an HTTP middleware that adds tracing, metrics, and logging
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract context from request
		ctx := r.Context()

		// Start a new span for this request
		ctx, span := Tracer().Start(
			ctx,
			r.URL.Path,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.user_agent", r.UserAgent()),
				attribute.String("http.remote_addr", r.RemoteAddr),
			),
		)
		defer span.End()

		// Create a custom response writer to capture the status code and size
		rw := &ResponseWriter{w, http.StatusOK, 0}

		// Increment active requests counter using Prometheus metrics
		activeRequests.Inc()

		// Record the start time
		startTime := time.Now()

		// Log the request
		LogWithContext(ctx, "Request started: %s %s", r.Method, r.URL.Path)

		// Call the next handler with the tracing context
		next.ServeHTTP(rw, r.WithContext(ctx))

		// Calculate request duration
		duration := time.Since(startTime)

		// Record metrics using Prometheus metrics
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", rw.StatusCode)).Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())

		// Also record using OpenTelemetry metrics for compatibility
		RequestCounter.Add(ctx, 1)
		RequestDuration.Record(ctx, duration.Seconds())

		// Decrement active requests counter
		activeRequests.Dec()
		ActiveRequests.Add(ctx, -1)

		// Add response details to span
		span.SetAttributes(
			attribute.Int("http.status_code", rw.StatusCode),
			attribute.Int("http.response_size", rw.Size),
		)

		// Set span status based on HTTP status code
		if rw.StatusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(rw.StatusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}

		// Log request completion
		LogWithContext(ctx, "Request completed: %s %s %d %d %.6fs",
			r.Method, r.URL.Path, rw.StatusCode, rw.Size, duration)
	})
}
