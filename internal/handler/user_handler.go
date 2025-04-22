package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Nezent/go-queue/common"
	"github.com/Nezent/go-queue/internal/domain"
	"github.com/Nezent/go-queue/internal/service"
)

type UserHandler struct {
	// UserService is the service layer for user-related operations.
	Service service.UserService
}

func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Parse the request body into a UserRegisterDTO
	var userDTO domain.UserRegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		common.RespondJSON(w, http.StatusBadRequest, common.ErrorResponse("Invalid request payload"))
		return
	}

	// Call the service to register the user
	userResponse, err := uh.Service.RegisterUser(ctx, userDTO)
	if err != nil {
		common.RespondJSON(w, http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	// Respond with the user data
	common.RespondJSON(w, http.StatusOK, common.SuccessResponse("User registered successfully", userResponse))

}
