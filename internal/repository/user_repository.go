package repository

import (
	"context"
	"time"

	"github.com/Nezent/go-queue/common"
	"github.com/Nezent/go-queue/internal/domain"
	"github.com/Nezent/go-queue/internal/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	RegisterUser(context.Context, domain.User) (*domain.User, *common.AppError)
	LoginUser(context.Context, string, string) (*uuid.UUID, *common.AppError)
	VerifyUser(context.Context, string) *common.AppError
}

type userRepository struct {
	db *pgxpool.Pool
}

func (ur userRepository) RegisterUser(ctx context.Context, user domain.User) (*domain.User, *common.AppError) {
	// Extract transaction from context
	tx, err := middleware.GetTxFromContext(ctx)
	if err != nil {
		return nil, common.NewUnexpectedServerError("Transaction context not found", err)
	}

	// Hash password
	hashedPassword, err := common.GenerateHashPassword(user.Password)
	if err != nil {
		return nil, common.NewUnexpectedServerError("Failed to hash password", err)
	}

	// Generate verification token (basic version, in production use a secure random generator)
	verificationToken, err := common.GenerateHash(uuid.New().String())
	if err != nil {
		return nil, common.NewUnexpectedServerError("Failed to generate validation token", err)
	}
	// Prepare user data
	var userID uuid.UUID

	// Insert into database
	query := `
		INSERT INTO users (name, email, password_hash, email_verified, verification_token, last_login_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`
	err = tx.QueryRow(ctx, query,
		user.Name, user.Email, hashedPassword, false,
		verificationToken, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339),
	).Scan(&userID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, common.NewDuplicateError("Email already in use")
		}
		return nil, common.NewUnexpectedServerError("Failed to register user", err)
	}

	user.ID = userID
	verification_token := verificationToken
	user.VerificationToken = *verification_token
	user.LastLoginAt = time.Now().Format(time.RFC3339)

	return &user, nil
}

func (ur userRepository) LoginUser(ctx context.Context, email, password string) (*uuid.UUID, *common.AppError) {
	var userID uuid.UUID
	var passwordHash string

	query := `SELECT id, password_hash FROM users WHERE email = $1`
	err := ur.db.QueryRow(ctx, query, email).Scan(&userID, &passwordHash)
	if err != nil {
		return nil, common.NewUnauthorizedError("Invalid credentials")
	}
	if err := common.CompareHashPassword(passwordHash, password); err != nil {
		return nil, common.NewUnauthorizedError("Invalid credentials")
	}
	return &userID, nil
}

func (r userRepository) VerifyUser(ctx context.Context, token string) *common.AppError {
	tx, err := middleware.GetTxFromContext(ctx)
	if err != nil {
		return common.NewUnexpectedServerError("Transaction context not found", err)
	}
	var userID uuid.UUID
	err = tx.QueryRow(ctx,
		`SELECT id FROM users WHERE verification_token = $1`, token).
		Scan(&userID)
	if err != nil {
		return common.NewUnexpectedServerError("Failed to verify user", err)
	}
	_, err = r.db.Exec(ctx,
		`UPDATE users SET email_verified = TRUE WHERE id = $1`, userID)
	if err != nil {
		return common.NewUnexpectedServerError("Failed to verify user", err)
	}
	return nil
}

func NewUserRepository(db *pgxpool.Pool) userRepository {
	return userRepository{db: db}
}
