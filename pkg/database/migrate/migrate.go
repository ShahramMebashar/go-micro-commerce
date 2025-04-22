package migrate

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var ErrNoMigrations = migrate.ErrNoChange

type Migrator struct {
	migrate        *migrate.Migrate
	migrationsPath string
	dsn            string
}

type Config struct {
	DSN            string
	MigrationsPath string
}

func New(cfg Config) (*Migrator, error) {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.MigrationsPath),
		cfg.DSN,
	)
	if err != nil {
		return nil, err
	}

	migrator := &Migrator{
		migrate:        m,
		migrationsPath: cfg.MigrationsPath,
		dsn:            cfg.DSN,
	}
	return migrator, nil
}

func (m *Migrator) Up() error {
	err := m.migrate.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return ErrNoMigrations
		}
		return fmt.Errorf("migration up failed: %w", err)
	}
	return nil
}

func (m *Migrator) Down() error {
	return m.migrate.Down()
}

func (m *Migrator) Steps(n int) error {
	return m.migrate.Steps(n)
}

func (m *Migrator) Version() (uint, bool, error) {
	return m.migrate.Version()
}

func (m *Migrator) Close() (error, error) {
	return m.migrate.Close()
}
