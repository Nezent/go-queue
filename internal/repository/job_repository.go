package repository

import (
	"context"
	"time"

	"github.com/Nezent/go-queue/common"
	"github.com/Nezent/go-queue/internal/domain"
	"github.com/Nezent/go-queue/internal/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository interface {
	// CreateJob creates a new job in the database.
	CreateJob(context.Context, domain.Job) (*domain.Job, *common.AppError)
	// GetJob retrieves a job by its ID.
	// GetJob(ctx context.Context, jobID uuid.UUID) (*domain.Job, *common.AppError)
	// // UpdateJob updates an existing job in the database.
	// UpdateJob(ctx context.Context, job domain.Job) (*domain.Job, *common.AppError)
	// // DeleteJob deletes a job from the database.
	// DeleteJob(ctx context.Context, jobID uuid.UUID) *common.AppError
	// // ListJobs retrieves a list of jobs with pagination.
	// ListJobs(ctx context.Context, page, limit int) ([]domain.Job, *common.AppError)

}
type jobRepository struct {
	db *pgxpool.Pool
}

func (jr *jobRepository) CreateJob(ctx context.Context, job domain.Job) (*domain.Job, *common.AppError) {
	// Extract transaction from context
	tx, err := middleware.GetTxFromContext(ctx)
	if err != nil {
		return nil, common.NewUnexpectedServerError("Transaction context not found", err)
	}

	// Insert into database
	query := `
		INSERT INTO jobs (user_id, type, payload, status, priority, attempts, run_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`
	run_at, _ := time.Parse(time.RFC3339, job.RunAt)
	job.Status = "pending"
	job.CreatedAt = time.Now().Format(time.RFC3339)
	job.UpdatedAt = time.Now().Format(time.RFC3339)
	job.Attempts = 0
	var jobID uuid.UUID
	err = tx.QueryRow(ctx, query,
		job.UserID, job.Type, job.Payload,
		job.Status, job.Priority, job.Attempts,
		run_at, time.Now().Format(time.RFC3339),
	).Scan(&jobID)
	if err != nil {
		return nil, common.NewUnexpectedServerError("Failed to create job", err)
	}
	job.ID = jobID
	return &job, nil
}

func NewJobRepository(db *pgxpool.Pool) *jobRepository {
	return &jobRepository{
		db: db,
	}
}
