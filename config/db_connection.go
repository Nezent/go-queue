package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// ConnectDB establishes a pooled connection to PostgreSQL using pgxpool.Pool with min/max connections configured.
func ConnectDB() (*pgxpool.Pool, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("[INFO] .env file not found, using system environment variables")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "postgres"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Load config from DSN
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to parse database config: %v\n", err)
		return nil, err
	}

	// Set min/max connections from environment or use default
	cfg.MinConns = int32(getEnvAsInt("DB_MIN_CONNS", 2))
	cfg.MaxConns = int32(getEnvAsInt("DB_MAX_CONNS", 10))

	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Printf("[ERROR] Failed to create DB pool: %v\n", err)
		return nil, err
	}

	// Ping to verify connection
	if err := dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		log.Printf("[ERROR] Database ping failed: %v\n", err)
		return nil, err
	}

	log.Printf("[INFO] Connected to PostgreSQL — MinConns: %d, MaxConns: %d", cfg.MinConns, cfg.MaxConns)
	return dbpool, nil
}

// getEnv returns the value of an environment variable or fallback if not found
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt parses an environment variable as int or returns fallback
func getEnvAsInt(key string, fallback int) int {
	valStr := os.Getenv(key)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}
