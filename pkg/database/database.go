package database

import (
	"context"
	"fmt"
	"log"
	"microservice/pkg/config"
	"microservice/pkg/database/migrate"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Initialize(cfg *config.Config) (*pgxpool.Pool, error) {
	log.Println("Initializing database...")

	dsn := cfg.DB.GetDSN()
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(cfg); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func RunMigrations(cfg *config.Config) error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	cwd = strings.TrimSuffix(cwd, cfg.DB.MigrationsPath)

	migrationsPath := filepath.Join(
		cwd,
		cfg.DB.MigrationsPath,
	)

	if _, err := os.Stat(migrationsPath); err != nil {
		return err
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
