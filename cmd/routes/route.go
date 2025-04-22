package routes

import (
	"github.com/Nezent/go-queue/internal/bootstrap"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, c *bootstrap.Container) {
	r.Route("/api/v1", func(r chi.Router) {
		// Group user routes
		r.Route("/users", func(r chi.Router) {
			r.Post("/", c.UserHandler.RegisterUser)
			// r.Get("/{user_id}", c.UserHandler.GetUserById)
			// r.Get("/", c.UserHandler.GetUsers)
		})

		// Continue with /auth, /jobs, etc.
	})

}
