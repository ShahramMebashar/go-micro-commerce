package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	// Database config errors
	ErrDBHostIsRequired     = errors.New("database host is required")
	ErrDBUserIsRequired     = errors.New("database user is required")
	ErrDBPortInvalid        = errors.New("database port must be between 1 and 65535")
	ErrDBDatabaseIsRequired = errors.New("database name is required")

	// Server config errors
	ErrServerPortInvalid    = errors.New("server port must be a valid number")
	ErrServerTimeoutInvalid = errors.New("server timeout must be positive")
)

// Environment represents the application environment (development, testing, production)
type Environment string

const (
	Development Environment = "development"
	Testing     Environment = "testing"
	Production  Environment = "production"
)

// DBConfig holds database connection configuration
type DBConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Database       string
	SSLMode        string
	MigrationsPath string
}

func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
}

// Validate checks if the database configuration is valid
func (c *DBConfig) Validate() error {
	// Check required fields
	if c.Host == "" {
		return ErrDBHostIsRequired
	}

	port, err := strconv.Atoi(c.Port)
	if err != nil || port < 1 || port > 65535 {
		return ErrDBPortInvalid
	}

	if c.User == "" {
		return ErrDBUserIsRequired
	}

	if c.Database == "" {
		return ErrDBDatabaseIsRequired
	}

	if c.MigrationsPath == "" {
		return fmt.Errorf("migrations path is required")
	}

	// Validate SSLMode (should be one of: disable, require, verify-ca, verify-full)
	validSSLModes := map[string]bool{
		"disable":     true,
		"require":     true,
		"verify-ca":   true,
		"verify-full": true,
	}
	if !validSSLModes[c.SSLMode] {
		return fmt.Errorf("invalid SSL mode: %s", c.SSLMode)
	}

	return nil
}

type ServerConfig struct {
	Port             string
	Timeout          int
	LogLevel         string
	AllowedOrigins   string
	AllowedMethods   string
	AllowedHeaders   string
	AllowCredentials bool
	MaxAge           int
}

func (s ServerConfig) GetAddr() string {
	return ":" + s.Port
}

func (s ServerConfig) Validate() error {
	port, err := strconv.Atoi(s.Port)
	if err != nil || port < 1 || port > 65535 {
		return ErrServerPortInvalid
	}

	if s.Timeout <= 0 {
		return ErrServerTimeoutInvalid
	}

	// Validate LogLevel (e.g., check if it's one of: "debug", "info", "warn", "error")
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[strings.ToLower(s.LogLevel)] {
		return fmt.Errorf("log level must be one of: debug, info, warn, error")
	}

	return nil
}

type TelemetryConfig struct {
	Enabled        bool
	ServiceName    string
	ServiceVersion string
	OTLPEndpoint   string
	JaegerEndpoint string
	MetricsEnabled bool
	MetricsPort    int
	PrometheusPath string
	LogLevel       string
}

// Config holds all application configuration
type Config struct {
	Env       Environment
	DB        DBConfig
	Server    ServerConfig
	Telemetry TelemetryConfig
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate database configuration
	if err := c.DB.Validate(); err != nil {
		return fmt.Errorf("database config: %w", err)
	}

	// Validate server configuration
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	return nil
}

// InProduction returns true if the application is running in production
func (c *Config) inProduction() bool {
	return c.Env == "production"
}

// LoadConfig loads configuration from environment variables and .env files and sets default values if they are not provided
func LoadConfig(envPath string) (*Config, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	cwd = strings.TrimSuffix(cwd, envPath)

	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	envFile := filepath.Join(
		cwd,
		envPath,
		".env",
	)

	err = godotenv.Load(envFile)
	if err != nil {
		log.Printf("Failed to load .env file %s: %v", envFile, err)
		err = godotenv.Load()
	}

	// Try loading from the root directory
	if err != nil {
		log.Printf("Failed to load .env file %v", err)
	}

	// Continue with configuration regardless of whether .env was loaded
	cfg := &Config{
		DB: DBConfig{
			Host:           GetEnv("DB_HOST", "localhost"),
			Port:           GetEnv("DB_PORT", "5432"),
			User:           GetEnv("DB_USER", "postgres"),
			Password:       GetEnv("DB_PASSWORD", "postgres"),
			Database:       GetEnv("DB_DATABASE", "products"),
			SSLMode:        GetEnv("DB_SSLMODE", "disable"),
			MigrationsPath: GetEnv("MIGRATIONS_PATH", "migrations"),
		},
		Server: ServerConfig{
			Port:             GetEnv("SERVER_PORT", "8080"),
			Timeout:          getEnvAsInt("SERVER_TIMEOUT", 30),
			LogLevel:         GetEnv("LOG_LEVEL", "info"),
			AllowedOrigins:   GetEnv("ALLOWED_ORIGINS", "*"),
			AllowedMethods:   GetEnv("ALLOWED_METHODS", "GET, POST, PUT, DELETE, OPTIONS"),
			AllowedHeaders:   GetEnv("ALLOWED_HEADERS", "Content-Type, Authorization, X-Requested-With, X-Request-ID"),
			AllowCredentials: getEnvAsBool("ALLOW_CREDENTIALS", false),
			MaxAge:           getEnvAsInt("MAX_AGE", 86400),
		},
		Env: Environment(GetEnv("ENV", "development")),
		Telemetry: TelemetryConfig{
			Enabled:        getEnvAsBool("TELEMETRY_ENABLED", true),
			ServiceName:    GetEnv("SERVICE_NAME", "product-service"),
			ServiceVersion: GetEnv("SERVICE_VERSION", "0.0.1"),
			OTLPEndpoint:   GetEnv("TELEMETRY_OTLP_ENDPOINT", "jaeger:4317"),
			JaegerEndpoint: GetEnv("TELEMETRY_JAEGER_ENDPOINT", "http://jaeger:14268/api/traces"),
			MetricsEnabled: getEnvAsBool("TELEMETRY_METRICS_ENABLED", true),
			MetricsPort:    getEnvAsInt("TELEMETRY_METRICS_PORT", 9090),
			PrometheusPath: GetEnv("TELEMETRY_PROMETHEUS_PATH", "/metrics"),
			LogLevel:       GetEnv("LOG_LEVEL", "info"),
		},
	}

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
