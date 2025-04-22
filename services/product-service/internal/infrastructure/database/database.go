package database

import (
	"context"
	"fmt"
	"log"
	"microservice/pkg/config"
	"microservice/pkg/database/migrate"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Initialize(cfg *config.Config) (*pgxpool.Pool, error) {
	log.Println("Initializing database...")
	dsn := cfg.DB.GetDSN()
	log.Printf("Using DSN: %s", dsn)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(cfg); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func runMigrations(cfg *config.Config) error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Try multiple possible migration paths to handle different working directories
	possiblePaths := []string{
		filepath.Join(cwd, "internal/migrations"),                          // When run from services/product-service
		filepath.Join(cwd, "services/product-service/internal/migrations"), // When run from root
	}

	// Find the first path that exists
	var migrationsPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			migrationsPath = path
			break
		}
	}

	if migrationsPath == "" {
		return fmt.Errorf("could not find migrations directory in any of the expected locations")
	}

	log.Printf("Using migrations path: %s", migrationsPath)

	migrator, err := migrate.New(migrate.Config{
		DSN:            cfg.DB.GetDSN(),
		MigrationsPath: migrationsPath,
	})
	if err != nil {
		return err
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil {
		if err == migrate.ErrNoMigrations {
			log.Println("No migrations to apply")
			return nil
		}
		return err
	}
	log.Println("Migrations completed successfully")
	return nil
}
