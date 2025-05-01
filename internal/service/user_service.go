package service

import (
	"context"

	"github.com/Nezent/go-queue/common"
	"github.com/Nezent/go-queue/internal/domain"
	"github.com/Nezent/go-queue/internal/repository"
	"github.com/Nezent/go-queue/internal/worker/enqueue"
	"github.com/Nezent/go-queue/internal/worker/task"
	"github.com/google/uuid"
)

type UserService interface {
	RegisterUser(context.Context, domain.UserRegisterDTO) (*domain.UserResponseDTO, *common.AppError)
	LoginUser(context.Context, domain.UserLoginRequestDTO) (*uuid.UUID, *common.AppError)
	VerifyUser(context.Context, string) *common.AppError
}
type userService struct {
	repo       repository.UserRepository
	dispatcher *enqueue.TaskDispatcher
}

func NewUserService(repo repository.UserRepository, dispatcher *enqueue.TaskDispatcher) UserService {
	return &userService{repo: repo, dispatcher: dispatcher}
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
	// Send verification email
	sendVerification(context, userResponse.Email, userResponse.VerificationToken, us.dispatcher)
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

func (us *userService) LoginUser(ctx context.Context, user domain.UserLoginRequestDTO) (*uuid.UUID, *common.AppError) {
	// Validate user data
	if user.Email == "" || user.Password == "" {
		return nil, common.NewBadRequestError("Email and password are required")
	}
	if !common.ValidateEmailWithRegex(user.Email) {
		return nil, common.NewBadRequestError("Invalid email format")
	}
	if len(user.Password) < 6 {
		return nil, common.NewBadRequestError("Password must be at least 6 characters long")
	}

	userID, err := us.repo.LoginUser(ctx, user.Email, user.Password)
	if err != nil {
		return nil, err
	}
	return userID, nil
}

func (us *userService) VerifyUser(ctx context.Context, token string) *common.AppError {
	// Validate token
	if token == "" {
		return common.NewBadRequestError("Token is required")
	}

	// Verify user
	err := us.repo.VerifyUser(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func sendVerification(context context.Context, email string, token string, dispatcher *enqueue.TaskDispatcher) {

	_ = dispatcher.EnqueueSendVerificationEmail(context, task.SendVerificationEmailPayload{
		Email: email,
		Token: token,
	})
}
