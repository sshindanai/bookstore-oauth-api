package errors

import (
	"net/http"
)

// RestErr is a custom error handling struct for this project
type RestErr struct {
	Message string `json:"message"`
	Code    int    `json:"status_code"`
	Error   string `json:"error"`
}

// NewBadRequestError is a custom error handling for BadRequest error
func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Code:    http.StatusBadRequest,
		Error:   "bad_request",
	}
}

// NewNotFoundError is a custom error handling for NotFound error
func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Code:    http.StatusNotFound,
		Error:   "not_found",
	}
}

// NewUnauthorizedError is a custom error handling for Unauthorized error
func NewUnauthorizedError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Code:    http.StatusUnauthorized,
		Error:   "unauthorized",
	}
}

// NewInternalServerError is a custom error handling for InternalServerError error
func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Code:    http.StatusInternalServerError,
		Error:   "internal_server_error",
	}
}

// NewConflictError is a custom error handling for Conflict error
func NewConflictError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Code:    http.StatusConflict,
		Error:   "conflict",
	}
}

func NewRestErrorFromBytes(b []byte) *RestErr {
	return &RestErr{
		Message: string(b),
		Code:    http.StatusBadRequest,
		Error:   "bad_request",
	}
}
