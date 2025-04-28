package service

import (
	"context"

	"github.com/Nezent/go-queue/common"
	"github.com/Nezent/go-queue/internal/domain"
	"github.com/Nezent/go-queue/internal/repository"
	"github.com/google/uuid"
)

type JobService interface {
	// CreateJob creates a new job in the database.
	CreateJob(context.Context, domain.JobCreateRequestDTO) (*domain.Job, *common.AppError)
	// GetJob retrieves a job by its ID.
	// GetJob(ctx context.Context, jobID uuid.UUID) (*domain.Job, *common.AppError)
	// // UpdateJob updates an existing job in the database.
	// UpdateJob(ctx context.Context, job domain.Job) (*domain.Job, *common.AppError)
	// // DeleteJob deletes a job from the database.
	// DeleteJob(ctx context.Context, jobID uuid.UUID) *common.AppError
}

type jobService struct {
	jobRepo repository.JobRepository
}

func (js *jobService) CreateJob(ctx context.Context, job domain.JobCreateRequestDTO) (*domain.Job, *common.AppError) {

	// Validate the job type
	if job.Type == "" {
		return nil, common.NewBadRequestError("Job type is required")
	}
	if job.UserID == uuid.Nil {
		return nil, common.NewBadRequestError("User ID is required")
	}

	jobEntity := domain.Job{
		UserID:   job.UserID,
		Type:     job.Type,
		Payload:  job.Payload,
		Priority: job.Priority,
		RunAt:    job.RunAt,
	}

	createdJob, appErr := js.jobRepo.CreateJob(ctx, jobEntity)
	if appErr != nil {
		return nil, appErr
	}

	return createdJob, nil
}

func NewJobService(jobRepo repository.JobRepository) *jobService {
	return &jobService{
		jobRepo: jobRepo,
	}
}
