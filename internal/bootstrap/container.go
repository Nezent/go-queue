package bootstrap

import (
	"github.com/jackc/pgx/v5"
)

type Container struct {
	// UserHandler handler.UserHandler
	// ... other handlers
}

func Initialize(db *pgx.Conn) *Container {

	return &Container{
		// UserHandler: handler.UserHandler{
		// 	Service: service.NewUserService(repository.NewUserRepositoryDatabase(db)),
		// },
		// other handlers...
	}
}
