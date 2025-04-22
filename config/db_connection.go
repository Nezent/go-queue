package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// ConnectDB establishes a new connection to the PostgreSQL database using pgx.Conn.
// It loads environment variables and sets a connection timeout.
// It returns a connection and an error for proper error handling.
func ConnectDB() (*pgx.Conn, error) {
	// Load .env file (optional)
	if err := godotenv.Load(".env"); err != nil {
		log.Println("[INFO] .env file not found, falling back to system environment variables")
	}

	// Construct DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "postgres"),
	)

	// Connect with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to database: %v\n", err)
		return nil, err
	}

	// Ping to verify
	if err := conn.Ping(ctx); err != nil {
		_ = conn.Close(ctx)
		log.Printf("[ERROR] Database ping failed: %v\n", err)
		return nil, err
	}

	log.Println("[INFO] Successfully connected to the PostgreSQL database")
	return conn, nil
}

// getEnv fetches env var or returns a fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
