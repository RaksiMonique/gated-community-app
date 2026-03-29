package domain

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int               `json:"-"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func NewValidationError(errors map[string]string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "Validation failed",
		Errors:  errors,
	}
}

var ErrInternal = &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}
