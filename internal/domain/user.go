package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	EmailVerified     bool      `json:"email_verified"`
	VerificationToken string    `json:"verification_token"`
	LastLoginAt       time.Time `json:"last_login_at"`
}

type UserRegisterDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponseDTO struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	EmailVerified     bool      `json:"email_verified"`
	VerificationToken string    `json:"verification_token"`
	LastLoginAt       time.Time `json:"last_login_at"`
}

type UserLoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponseDTO struct {
	AccessToken string `json:"access_token"`
}
