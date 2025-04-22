package common

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{Success: true, Message: message, Data: data}
}

func ErrorResponse(errMsg string) APIResponse {
	return APIResponse{Success: false, Error: errMsg}
}

func RespondJSON(w http.ResponseWriter, statusCode int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, `{"success":false,"error":"Internal Server Error"}`, http.StatusInternalServerError)
	}
}
