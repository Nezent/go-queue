package service

import (
	"context"

	"github.com/Nezent/go-queue/common"
	"github.com/Nezent/go-queue/internal/domain"
	"github.com/Nezent/go-queue/internal/repository"
)

type UserService interface {
	RegisterUser(context.Context, domain.UserRegisterDTO) (*domain.UserResponseDTO, *common.AppError)
	// LoginUser(username, password string) (int, error)
	// GetUserByID(userID int) (string, error)
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
func (us *userService) RegisterUser(context context.Context, user domain.UserRegisterDTO) (*domain.UserResponseDTO, *common.AppError) {
	// Validate user data
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return nil, common.NewBadRequestError("Name, email, and password are required")
	}
	if !common.ValidateEmailWithRegex(user.Email) {
		return nil, common.NewBadRequestError("Invalid email format")
	}
	if len(user.Password) < 6 {
		return nil, common.NewBadRequestError("Password must be at least 6 characters long")
	}

	domainUser := domain.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	// Register user
	userResponse, err := us.repo.RegisterUser(context, domainUser)
	if err != nil {
		return nil, err
	}
	// Convert to response DTO
	responseDTO := &domain.UserResponseDTO{
		ID:                userResponse.ID,
		Name:              userResponse.Name,
		Email:             userResponse.Email,
		EmailVerified:     userResponse.EmailVerified,
		VerificationToken: userResponse.VerificationToken,
		LastLoginAt:       userResponse.LastLoginAt,
	}
	return responseDTO, nil
}
