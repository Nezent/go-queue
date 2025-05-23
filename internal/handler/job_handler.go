package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nezent/go-queue/common"
	"github.com/Nezent/go-queue/internal/domain"
	"github.com/Nezent/go-queue/internal/service"
	"github.com/Nezent/go-queue/internal/worker/task"
	"github.com/google/uuid"
)

type JobHandler struct {
	Service service.JobService
}

func (jh *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var jobDTO domain.JobCreateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&jobDTO); err != nil {
		common.RespondJSON(w, http.StatusBadRequest, common.ErrorResponse("Invalid request payload"))
		return
	}

	jobResponse, err := jh.Service.CreateJob(ctx, jobDTO)
	if err != nil {
		common.RespondJSON(w, http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	common.RespondJSON(w, http.StatusOK, common.SuccessResponse("Job created successfully", jobResponse))
}

func (jh *JobHandler) GetJobPayload(ctx context.Context, jobID uuid.UUID) (*task.JobPayload, error) {
	if jobID == uuid.Nil {
		return nil, errors.New("invalid job id")
	}
	jobPayload, err := jh.Service.GetJobPayload(ctx, jobID)
	if err != nil {
		return nil, err
	}
	return jobPayload, nil
}

func (jh *JobHandler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobIDStr := r.URL.Query().Get("job_id")
	if jobIDStr == "" {
		common.RespondJSON(w, http.StatusBadRequest, common.ErrorResponse("Job ID is required"))
		return
	}

	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		common.RespondJSON(w, http.StatusBadRequest, common.ErrorResponse("Invalid Job ID format"))
		return
	}

	jobStatus, appErr := jh.Service.GetJobStatus(ctx, jobID)
	if appErr != nil {
		common.RespondJSON(w, http.StatusInternalServerError, common.ErrorResponse(appErr))
		return
	}

	common.RespondJSON(w, http.StatusOK, common.SuccessResponse("Job status retrieved successfully", jobStatus))
}

func (jh *JobHandler) UpdateJobStatus(ctx context.Context, jobID uuid.UUID, status string, attempts int) *common.AppError {

	_, appErr := jh.Service.UpdateJobStatus(ctx, jobID, status, attempts)
	if appErr != nil {
		return appErr
	}

	return nil
}
