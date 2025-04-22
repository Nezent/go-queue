package domain

import "github.com/google/uuid"

type User struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	EmailVerified     bool      `json:"email_verified"`
	VerificationToken string    `json:"verification_token"`
	LastLoginAt       string    `json:"last_login_at"`
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
	LastLoginAt       string    `json:"last_login_at"`
}

type UserLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
