package telemetry

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel/trace"
)

// setupLogging configures the logger to include trace information
func setupLogging(_ Config) error {
	// Set up a basic logger for now
	// In a real implementation, you might want to use a structured logger
	// that can include trace IDs and other context
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// Set output to stdout
	log.SetOutput(os.Stdout)

	return nil
}

// LogWithContext adds trace information to log messages
func LogWithContext(ctx context.Context, format string, args ...any) {
	// Extract trace information from context
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()

		// Prepend trace info to the format string
		format = "[trace_id=%s span_id=%s] " + format

		// Insert trace IDs at the beginning of args
		newArgs := make([]any, len(args)+2)
		newArgs[0] = traceID
		newArgs[1] = spanID
		copy(newArgs[2:], args)

		log.Printf(format, newArgs...)
	} else {
		// No span in context, log normally
		log.Printf(format, args...)
	}
}
