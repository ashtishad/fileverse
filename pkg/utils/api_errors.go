package utils

import (
	"net/http"
)

// APIError represents an custom api error,
// Code returns the associated HTTP status code as an integer.
// Error returns a user-friendly error message, potentially hiding internal errors.
// Wrap encapsulates the given error within the APIError, providing additional context.
// The wrapped error is not included in the message returned by Error().
type APIError interface {
	Error() string
	Wrap(err error) APIError
	Code() int
}

// apiError is a concrete implementation of the APIError interface.
type apiError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
	Causes     string `json:"causes"`
}

// Code returns http status code
func (e *apiError) Code() int {
	return e.StatusCode
}

// Error returns error message and code. But hides internal server/db related errors
func (e *apiError) Error() string {
	return e.Message
}

// Wrap wraps internal causes of error into an APIError
// use-cases: adding internal causes into error.
func (e *apiError) Wrap(err error) APIError {
	if err != nil {
		e.Causes = err.Error()
	}

	return e
}

// InternalServerError creates a new APIError for internal server errors.
// returns http.StatusInternalServerError 500.
// Example usage:
//
//	err := InternalServerError("internal server error", innerErr)
func InternalServerError(message string, err error) APIError {
	result := &apiError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}

	return result.Wrap(err)
}

// BadRequestError creates a new APIError for bad requests.
// returns http.StatusBadRequest 400.
// Example usage:
//
//	err := BadRequestError("invalid input").Wrap(innerErr)
func BadRequestError(message string) APIError {
	return &apiError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NotFoundError creates a new APIError for not found errors.
// returns http.StatusNotFound 404.
// Example usage:
//
//	err := NotFoundError("resource not found")
func NotFoundError(message string) APIError {
	return &apiError{
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

// UnauthorizedError creates a new APIError for unauthorized requests.
// returns http.StatusUnauthorized 401.
// Example usage:
//
//	err := UnauthorizedError("unauthorized")
func UnauthorizedError(message string) APIError {
	return &apiError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// RateLimitError creates a new APIError for rate limit error.
// returns http.StatusTooManyRequests 429.
// Example usage:
//
//	err := RateLimitError("api request limited")
func RateLimitError(message string) APIError {
	return &apiError{
		Message:    message,
		StatusCode: http.StatusTooManyRequests,
	}
}

// ConflictError creates a new APIError for duplicate fields,
// returns http.StatusConflict 409.
// Example usage:
//
//	err := ConflictError("user name already exists")
func ConflictError(message string) APIError {
	return &apiError{
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}
