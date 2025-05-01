package routes

import (
	"net/http"

	"github.com/Nezent/go-queue/internal/bootstrap"
	"github.com/Nezent/go-queue/internal/handler"
	"github.com/Nezent/go-queue/internal/middleware"
	"github.com/Nezent/go-queue/internal/websocket"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, c *bootstrap.Container) {
	r.Route("/api/v1", func(api chi.Router) {

		// 🔐 Auth Routes (Public)
		api.Route("/auth", func(auth chi.Router) {
			auth.Post("/login", c.UserHandler.LoginHandler)
			auth.Post("/register", c.UserHandler.RegisterUser)
			auth.Post("/logout", handler.LogoutHandler)
			auth.Get("/verify", c.UserHandler.VerifyUser) // Verify user email
			// Optional: Uncomment if you have a refresh token endpoint
			// auth.Post("/refresh", c.UserHandler.RefreshTokenHandler)
		})

		// 👤 User Routes (Protected)
		api.Route("/users", func(users chi.Router) {
			users.Use(middleware.AuthMiddleware)

			// Add optional role-based access
			// users.Use(middleware.RequireRole("admin"))

			// users.Get("/", c.UserHandler.GetUsers)
			// users.Get("/{user_id}", c.UserHandler.GetUserById)
		})

		// 💼 Job Routes (Protected)
		api.Route("/jobs", func(jobs chi.Router) {
			jobs.Use(middleware.AuthMiddleware)

			// Optional: jobs.Use(middleware.RequireRole("admin", "hr"))

			// jobs.Get("/", c.JobHandler.GetJobs)
			jobs.Post("/", c.JobHandler.CreateJob)
			// jobs.Get("/{job_id}", c.JobHandler.GetJobById)
		})

		// 📦 WebSocket Routes
		api.Route("/ws", func(ws chi.Router) {
			ws.Use(middleware.AuthMiddleware)
			ws.Get("/jobs", func(w http.ResponseWriter, r *http.Request) {
				websocket.HandleWebSocket(c.WebSocketHub, w, r)
			})
		})
		// ✅ Add more domain-specific groups below (example: /tasks, /reports, /analytics)
	})
}
