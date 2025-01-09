package helper

import (
	"errors"
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

func GetErrorStatusCode(err error) int {
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
