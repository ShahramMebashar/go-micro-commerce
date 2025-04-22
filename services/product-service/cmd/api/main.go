package main

import (
	"context"
	"log"
	"microservice/pkg/config"
	"microservice/pkg/logger"
	"microservice/services/product-service/internal/application"
	"microservice/services/product-service/internal/infrastructure/api"
	"microservice/services/product-service/internal/infrastructure/database"
	"microservice/services/product-service/internal/infrastructure/persistence/postgres"
	"microservice/services/product-service/internal/infrastructure/validator"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load configuration with multiple possible .env files
	cfg, err := config.LoadConfig("services/product-service")

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbpool, err := database.Initialize(cfg)

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer dbpool.Close()
	lg := logger.GetDefaultLogger()

	productRepo := postgres.NewProductRepository(dbpool)
	productService := application.NewProductService(productRepo, lg)
	productHandler := api.NewProductHandler(productService, validator.New())

	runSrever(cfg, productHandler, lg)
}

func runSrever(cfg *config.Config, productHandler *api.ProductHandler, logger logger.Logger) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(api.RequestID)
	r.Use(api.Logger(logger))
	r.Use(middleware.Timeout(time.Duration(cfg.Server.Timeout) * time.Second))
	r.Use(middleware.CleanPath)
	r.Use(middleware.Compress(5))
	r.Use(api.CORS(cfg))
	r.Use(api.ContentTypeJson)

	r.Route("/api", func(r chi.Router) {
		productHandler.RegisterRoutes(r)
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		api.RespondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "product-service",
		})
	})
	server := &http.Server{
		Addr:    cfg.Server.GetAddr(),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("Failed to start server: %v", err)
			}
		}
	}()

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down server...")

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}

	log.Println("Server stopped")
}
