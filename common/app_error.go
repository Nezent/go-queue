package common

import (
	"encoding/json"
	"net/http"
)

// AppError defines a structured error used throughout the app.
type AppError struct {
	StatusCode int    `json:"status_code,omitempty"` // HTTP status code
	Message    string `json:"message"`               // Client-facing message
	Err        error  `json:"-"`                     // Internal error (not exposed in JSON)
}

// Error implements the built-in error interface.
func (e *AppError) Error() string {
	return e.Message
}

// AsMessage returns a sanitized version of AppError for safe client responses.
func (e *AppError) AsMessage() *AppError {
	return &AppError{
		Message: e.Message,
	}
}

// WriteJSON sends the error as a JSON response to the client.
func (e *AppError) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)
	json.NewEncoder(w).Encode(e.Err.Error())
}

// Error Constructors

func NewUnexpectedServerError(message string, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		Err:        err,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func NewDuplicateError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		Message:    message,
	}
}
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

// Generic wrapper for any error
func WrapError(statusCode int, message string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}
