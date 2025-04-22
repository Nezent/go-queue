package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Nezent/go-queue/cmd/routes"
	"github.com/Nezent/go-queue/config"
	"github.com/Nezent/go-queue/internal/bootstrap"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Replace with your allowed origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           30,
	}))

	// Middleware
	// r.Use(chi.Logger)
	// r.Use(middleware.WithTransaction(db))
	// r.Use(middleware.AuthMiddleware)
	// r.Use(middleware.Recoverer)

	// Database connection
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close(context.Background())
	// Dependency injection
	container := bootstrap.Initialize(db)

	// Add all routes
	routes.RegisterRoutes(r, container)

	log.Println("🚀 Server running at :8080")
	http.ListenAndServe(":8080", r)
}
