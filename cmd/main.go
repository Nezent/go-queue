package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Nezent/go-queue/cmd/routes"
	"github.com/Nezent/go-queue/config"
	"github.com/Nezent/go-queue/internal/bootstrap"
	"github.com/Nezent/go-queue/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func main() {

	// Connect to DB
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db.Close()
		log.Println("Database pool closed")
	}()

	// Initialize Chi router
	r := chi.NewRouter()
	r.Use(middleware.WithTransaction(db))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           30,
	}))

	if err := godotenv.Load(".env"); err != nil {
		log.Println("[INFO] .env file not found, using system environment variables")
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: os.Getenv("REDIS_ADDR"),
	}

	dispatcher := bootstrap.InitializeDispatcher(redisOpt)

	// Dependency injection
	container := bootstrap.Initialize(db, dispatcher)

	// Register all routes
	routes.RegisterRoutes(r, container)

	log.Println("🚀 Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
