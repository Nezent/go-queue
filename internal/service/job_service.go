package service

import (
	"context"

	"github.com/Nezent/go-queue/common"
	"github.com/Nezent/go-queue/internal/domain"
	"github.com/Nezent/go-queue/internal/middleware"
	"github.com/Nezent/go-queue/internal/repository"
	"github.com/Nezent/go-queue/internal/worker/task"
	"github.com/google/uuid"
)

type JobService interface {
	// CreateJob creates a new job in the database.
	CreateJob(context.Context, domain.JobCreateRequestDTO) (*domain.Job, *common.AppError)
	// GetJobPayload retrieves a job payload by its ID.
	GetJobPayload(context.Context, uuid.UUID) (*task.JobPayload, *common.AppError)
	// UpdateJobStatus updates an existing job status in the database.
	UpdateJobStatus(context.Context, uuid.UUID) (*domain.Job, *common.AppError)
	// GetJobStatus retrieves the status of a job by its ID.
	GetJobStatus(context.Context, uuid.UUID) (*domain.JobStatusResponseDTO, *common.AppError)
}

type jobService struct {
	jobRepo repository.JobRepository
}

func (js *jobService) CreateJob(ctx context.Context, job domain.JobCreateRequestDTO) (*domain.Job, *common.AppError) {

	// Validate the job type
	if job.Type == "" {
		return nil, common.NewBadRequestError("Job type is required")
	}
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, common.NewUnauthorizedError("User ID not found in context")
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, common.NewBadRequestError("Invalid User ID format")
	}

	jobEntity := domain.Job{
		UserID:   parsedUserID,
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

func (js *jobService) GetJobStatus(ctx context.Context, jobID uuid.UUID) (*domain.JobStatusResponseDTO, *common.AppError) {
	// Retrieve job status from the repository
	job, appErr := js.jobRepo.GetJobStatus(ctx, jobID)
	if appErr != nil {
		return nil, appErr
	}

	if job == nil {
		return nil, common.NewNotFoundError("Job not found")
	}

	return job, nil
}

func (js *jobService) GetJobPayload(ctx context.Context, jobID uuid.UUID) (*task.JobPayload, *common.AppError) {
	// Retrieve job payload from the repository
	payload, appErr := js.jobRepo.GetJobPayload(ctx, jobID)
	if appErr != nil {
		return nil, appErr
	}

	return payload, nil
}
func (js *jobService) UpdateJobStatus(ctx context.Context, jobID uuid.UUID) (*domain.Job, *common.AppError) {
	// Update job status in the repository
	job, appErr := js.jobRepo.UpdateJobStatus(ctx, jobID)
	if appErr != nil {
		return nil, appErr
	}

	return job, nil
}

func NewJobService(jobRepo repository.JobRepository) *jobService {
	return &jobService{
		jobRepo: jobRepo,
	}
}
