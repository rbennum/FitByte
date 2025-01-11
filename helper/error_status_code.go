package helper

import (
	"errors"
	"fmt"
	"net/http"
)

// Define some common errors
var (
	ErrNotFound     = errors.New("record not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrBadRequest   = errors.New("bad request")
	ErrConflict     = errors.New("data conflict")

	ErrInvalidDepartmentId    = errors.New("invalid department id")
	ErrConflictIdentityNumber = errors.New("identity number conflict")

	ErrInternalServer = errors.New("internal server error")
)

// ErrorResponse represents error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error response implements the Error interface
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

func NewErrorResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{Code: code, Message: message}
}

func GetErrorStatusCode(err error) int {
	// detect if the error instance if one of the ErrorResponse
	if httpErr, ok := err.(*ErrorResponse); ok {
		return httpErr.Code
	}
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrConflict:
		return http.StatusConflict
	case ErrInvalidDepartmentId:
		return http.StatusBadRequest
	case ErrConflictIdentityNumber:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func GetErrorMessage(err error) string {
	// detect if the error instance if one of the ErrorResponse
	if httpErr, ok := err.(*ErrorResponse); ok {
		return httpErr.Message
	}
	switch err {
	case ErrNotFound:
		return ErrNotFound.Error()
	case ErrUnauthorized:
		return ErrUnauthorized.Error()
	case ErrBadRequest:
		return ErrBadRequest.Error()
	case ErrConflict:
		return ErrConflict.Error()
	case ErrInvalidDepartmentId:
		return ErrInvalidDepartmentId.Error()
	case ErrConflictIdentityNumber:
		return ErrConflictIdentityNumber.Error()
	default:
		return ErrInternalServer.Error()
	}
}
