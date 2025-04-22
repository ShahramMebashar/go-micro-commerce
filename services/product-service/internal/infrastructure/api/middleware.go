package api

import (
	"context"
	"fmt"
	"microservice/pkg/config"
	"microservice/pkg/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type MiddlewareFunc func(next http.Handler) http.Handler

// ContentTypeJson sets the Content-Type header to application/json
func ContentTypeJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// RequestIDKey is the context key for the request ID

// RequestIDHeader is the header name for the request ID
const RequestIDHeader = "X-Request-ID"

// RequestID adds a request ID to the context and response headers
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)

		if requestID == "" {
			uuid, err := uuid.NewUUID()
			if err == nil {
				requestID = uuid.String()
			}
		}

		w.Header().Set(RequestIDHeader, requestID)

		ctx := context.WithValue(r.Context(), middleware.RequestIDKey, requestID)

		// set context value
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// responseWriter is a custom response writer that captures the status code and response size
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Size       int
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the response size
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.Size += size
	return size, err
}

// Logger logs information about each request
func Logger(logger logger.Logger) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &ResponseWriter{w, http.StatusOK, 0}

			next.ServeHTTP(rw, r)

			requestID, _ := r.Context().Value(middleware.RequestIDKey).(string)

			logFmts := fmt.Sprintf(
				"request_id=%s method=%s status=%d(%s) path=%s query=%s remote_ip=%s user_agent=%s duration=%s size=%d",
				requestID,
				r.Method,
				rw.StatusCode,
				http.StatusText(rw.StatusCode),
				r.URL.Path,
				r.URL.RawQuery,
				r.RemoteAddr,
				r.UserAgent(),
				time.Since(start).String(),
				rw.Size,
			)

			if rw.StatusCode < 400 {
				logger.Info(logFmts)
			} else if rw.StatusCode < 500 {
				logger.Warn(logFmts)
			} else {
				logger.Error(logFmts)
			}
		}
		return http.HandlerFunc(fn)
	}
}

// CORS adds CORS headers to responses
func CORS(cfg *config.Config) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			origin := cfg.Server.AllowedOrigins
			if origin == "" {
				origin = "*" // Default to allow all origins
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", cfg.Server.AllowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", cfg.Server.AllowedHeaders)

			if cfg.Server.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
