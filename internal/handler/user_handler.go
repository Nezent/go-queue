package handler

import (
	"encoding/json"
	"net/http"
	"time"

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

func (uh *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Parse the request body into a UserLoginRequestDTO
	var userDTO domain.UserLoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		common.RespondJSON(w, http.StatusBadRequest, common.ErrorResponse("Invalid request payload"))
		return
	}

	// Call the service to login the user
	userID, err := uh.Service.LoginUser(ctx, userDTO)
	if err != nil {
		common.RespondJSON(w, http.StatusUnauthorized, common.ErrorResponse("Invalid email or password"))
		return
	}

	// Generate JWT token
	accessToken, appError := common.GenerateJWT(userID.String(), "user", time.Minute*15)
	if appError != nil {
		common.RespondJSON(w, http.StatusInternalServerError, common.ErrorResponse("Failed to generate token"))
		return
	}

	// Set the access token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Minute * 15),
	})

	common.RespondJSON(w, http.StatusOK, common.SuccessResponse("Login successful", domain.UserLoginResponseDTO{AccessToken: accessToken}))

}

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expire := time.Unix(0, 0)
		cookies := []string{"access_token", "refresh_token"}
		for _, name := range cookies {
			http.SetCookie(w, &http.Cookie{
				Name:     name,
				Value:    "",
				Path:     "/",
				Expires:  expire,
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})
		}
		w.WriteHeader(http.StatusOK)
	}
}
