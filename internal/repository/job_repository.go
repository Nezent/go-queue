package repository

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Nezent/go-queue/common"
	"github.com/Nezent/go-queue/internal/domain"
	"github.com/Nezent/go-queue/internal/middleware"
	"github.com/Nezent/go-queue/internal/worker/task"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository interface {
	// CreateJob creates a new job in the database.
	CreateJob(context.Context, domain.Job) (*domain.Job, *common.AppError)
	// GetJobPayload retrieves a job payload by its ID.
	GetJobPayload(context.Context, uuid.UUID) (*task.JobPayload, *common.AppError)
	// // UpdateJobStatus updates an existing job status in the database.
	UpdateJobStatus(context.Context, uuid.UUID, string, int) (*domain.Job, *common.AppError)
	// GetJobStatus retrieves the status of a job by its ID.
	GetJobStatus(context.Context, uuid.UUID) (*domain.JobStatusResponseDTO, *common.AppError)
}
type jobRepository struct {
	db *pgxpool.Pool
}

func (jr jobRepository) CreateJob(ctx context.Context, job domain.Job) (*domain.Job, *common.AppError) {
	// Extract transaction from context
	tx, err := middleware.GetTxFromContext(ctx)
	if err != nil {
		return nil, common.NewUnexpectedServerError("Transaction context not found", err)
	}

	// Insert into database
	query := `
		INSERT INTO jobs (user_id, type, payload, status, priority, attempts, run_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`
	run_at := job.RunAt
	job.Status = "pending"
	job.CreatedAt = time.Now().In(common.DhakaTZ)
	job.UpdatedAt = time.Now().In(common.DhakaTZ)
	job.Attempts = 0
	var jobID uuid.UUID
	err = tx.QueryRow(ctx, query,
		job.UserID, job.Type, job.Payload,
		job.Status, job.Priority, job.Attempts,
		run_at,
	).Scan(&jobID)
	if err != nil {
		return nil, common.NewUnexpectedServerError("Failed to create job", err)
	}
	job.ID = jobID
	return &job, nil
}

func (jr jobRepository) GetJobPayload(ctx context.Context, jobID uuid.UUID) (*task.JobPayload, *common.AppError) {
	query := `
		SELECT type, payload, status, priority, attempts, run_at FROM jobs WHERE id = $1
	`

	var payload task.JobPayload
	var rawPayload []byte // payload column as JSON
	err := jr.db.QueryRow(ctx, query, jobID).Scan(
		&payload.JobType,
		&rawPayload,
		&payload.Status,
		&payload.Priority,
		&payload.Attempts,
		&payload.RunAt,
	)
	if err != nil {
		log.Printf("[DEBUG] QueryRow scan failed for jobID %s: %v", jobID, err)
		return nil, common.NewUnexpectedServerError("Failed to retrieve job payload", err)
	}

	// Unmarshal JSON payload
	var emailPayload task.EmailPayload
	if err := json.Unmarshal(rawPayload, &emailPayload); err != nil {
		return nil, common.NewUnexpectedServerError("Failed to parse job payload JSON", err)
	}

	payload.Payload = emailPayload
	return &payload, nil
}

func (jr jobRepository) UpdateJobStatus(ctx context.Context, jobID uuid.UUID, status string, attempts int) (*domain.Job, *common.AppError) {
	// Extract transaction from context
	tx, err := middleware.GetTxFromContext(ctx)
	if err != nil {
		return nil, common.NewUnexpectedServerError("Transaction context not found", err)
	}

	// Update job status in database
	query := `
		UPDATE jobs SET status = $1, attempts = $2, updated_at = $3 WHERE id = $4 RETURNING *
	`
	job := domain.Job{}
	err = tx.QueryRow(ctx, query, status, attempts, time.Now().In(common.DhakaTZ), jobID).Scan(&job.ID, &job.UserID, &job.Type, &job.Payload, &job.Status, &job.Priority, &job.Attempts, &job.RunAt, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, common.NewUnexpectedServerError("Failed to update job status", err)
	}
	job.UpdatedAt = time.Now().In(common.DhakaTZ)
	return &job, nil
}

func (jr jobRepository) GetJobStatus(ctx context.Context, jobID uuid.UUID) (*domain.JobStatusResponseDTO, *common.AppError) {
	// Retrieve job status from database
	query := `
		SELECT type, status, priority, attempts, run_at
		FROM jobs WHERE id = $1
	`
	job := domain.JobStatusResponseDTO{}
	err := jr.db.QueryRow(ctx, query, jobID).Scan(&job.Type, &job.Status, &job.Priority, &job.Attempts, &job.RunAt)
	if err != nil {
		return nil, common.NewUnexpectedServerError("Failed to retrieve job status", err)
	}
	return &job, nil
}

func NewJobRepository(db *pgxpool.Pool) jobRepository {
	return jobRepository{
		db: db,
	}
}
