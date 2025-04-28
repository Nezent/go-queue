package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Nezent/go-queue/common"
	"github.com/Nezent/go-queue/internal/domain"
	"github.com/Nezent/go-queue/internal/service"
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
