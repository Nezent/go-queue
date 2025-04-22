package repository

import (
	"github.com/Nezent/go-queue/common"
	"github.com/Nezent/go-queue/internal/domain"
	"github.com/jackc/pgx"
)

type UserRepository interface {
	RegisterUser(domain.User) common.APIResponse
	LoginUser(username, password string) (int, error)
	GetUserByID(userID int) (string, error)
}

type userRepository struct {
	db *pgx.Conn
}

func (ur userRepository) RegisterUser(user domain.User) common.APIResponse {
	// Implement the logic to register a user in the database
	// This is a placeholder implementation
	return common.SuccessResponse("User registered successfully", nil)
}

func NewUserRepository(db *pgx.Conn) userRepository {
	return userRepository{db: db}
}
