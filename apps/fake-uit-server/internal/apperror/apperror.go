package apperror

import "net/http"

type AppError struct {
	Status  int
	Message string
	Code    string
}

func (e AppError) Error() string {
	return e.Message
}

var (
	ErrBadRequest = AppError{
		Status:  http.StatusBadRequest,
		Message: "Bad request!",
		Code:    "BAD_REQUEST",
	}
	ErrUnauthorized = AppError{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized!",
		Code:    "UNAUTHORIZED",
	}
	ErrForbidden = AppError{
		Status:  http.StatusForbidden,
		Message: "Forbidden!",
		Code:    "FORBIDDEN",
	}
	ErrNotFound = AppError{
		Status:  http.StatusNotFound,
		Message: "Not found!",
		Code:    "NOT_FOUND",
	}
	ErrTooManyRequests = AppError{
		Status:  http.StatusTooManyRequests,
		Message: "Too many requests!",
		Code:    "TOO_MANY_REQUESTS",
	}
	ErrServiceUnavailable = AppError{
		Status:  http.StatusServiceUnavailable,
		Message: "Service unavailable!",
		Code:    "SERVICE_UNAVAILABLE",
	}
	ErrInternal = AppError{
		Status:  http.StatusInternalServerError,
		Message: "Internal error happened!",
		Code:    "INTERNAL_ERROR",
	}
)

func FromStatus(status int) AppError {
	switch status {
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	case http.StatusServiceUnavailable:
		return ErrServiceUnavailable
	default:
		return ErrInternal
	}
}
