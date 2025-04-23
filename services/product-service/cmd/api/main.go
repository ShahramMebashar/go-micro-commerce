package main

import (
	"context"
	"log"
	"microservice/pkg/config"
	"microservice/pkg/database"
	"microservice/pkg/logger"
	"microservice/pkg/telemetry"
	"microservice/services/product-service/internal/application"
	"microservice/services/product-service/internal/infrastructure/api"
	"microservice/services/product-service/internal/infrastructure/persistence/postgres"
	"microservice/services/product-service/internal/infrastructure/validator"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "microservice/docs" // This is important for swagger to find the docs

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Product Service API
// @version 1.0
// @description This is the product service API for the microservice architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	// Load configuration
	appCfg, err := config.LoadConfig("services/product-service")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize telemetry first
	telemetryCfg := telemetry.Config{
		ServiceName:    appCfg.Telemetry.ServiceName,
		ServiceVersion: appCfg.Telemetry.ServiceVersion,
		Environment:    string(appCfg.Env),
		TracingEnabled: appCfg.Telemetry.Enabled,
		OTLPEndpoint:   appCfg.Telemetry.OTLPEndpoint,
		JaegerEndpoint: appCfg.Telemetry.JaegerEndpoint,
		MetricsEnabled: appCfg.Telemetry.MetricsEnabled,
		MetricsPort:    appCfg.Telemetry.MetricsPort,
		PrometheusPath: appCfg.Telemetry.PrometheusPath,
		LogLevel:       appCfg.Server.LogLevel,
	}

	telShutdown, err := telemetry.Setup(context.Background(), telemetryCfg)
	if err != nil {
		log.Fatalf("Failed to setup telemetry: %v", err)
	}
	defer telShutdown(context.Background())

	// Initialize database
	dbpool, err := database.Initialize(appCfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbpool.Close()

	// Now use the tracer after it's been initialized
	tr := telemetry.Tracer()
	lg := logger.GetDefaultLogger()
	productRepo := postgres.NewProductRepository(dbpool, tr)
	productService := application.NewProductService(productRepo, tr)
	productHandler := api.NewProductHandler(productService, validator.New(), lg)

	runServer(appCfg, productHandler, lg)
}

func runServer(cfg *config.Config, productHandler *api.ProductHandler, logger logger.Logger) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(api.RequestID)
	r.Use(api.Logger(logger))
	r.Use(telemetry.Middleware)
	r.Use(middleware.Timeout(time.Duration(cfg.Server.Timeout) * time.Second))
	r.Use(middleware.CleanPath)
	r.Use(middleware.Compress(5))
	r.Use(api.CORS(cfg))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL pointing to API definition
	))

	r.Route("/api", func(r chi.Router) {
		r.Use(api.ContentTypeJson)
		productHandler.RegisterRoutes(r)
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Get hostname to identify which container is responding
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}
		api.RespondWithJSON(w, http.StatusOK, map[string]string{
			"status":   "ok",
			"service":  "product-service",
			"hostname": hostname,
		})
	})
	server := &http.Server{
		Addr:    cfg.Server.GetAddr(),
		Handler: r,
	}

	go func() {
		logger.Info("Server started on http://localhost" + cfg.Server.GetAddr())
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logger.Fatal("Failed to start server: %v", err)
			}
		}
	}()

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")

	err := server.Shutdown(ctx)
	if err != nil {
		logger.Fatal("Failed to shutdown server: %v", err)
	}

	logger.Info("Server stopped")
}
