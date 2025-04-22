package bootstrap

import (
	"github.com/Nezent/go-queue/internal/handler"
	"github.com/Nezent/go-queue/internal/repository"
	"github.com/Nezent/go-queue/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	UserHandler handler.UserHandler
}

func Initialize(db *pgxpool.Pool) *Container {

	return &Container{
		UserHandler: handler.UserHandler{
			Service: service.NewUserService(repository.NewUserRepository(db)),
		},
		// other handlers...
	}
}
