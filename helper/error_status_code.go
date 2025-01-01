package helper

import (
    "errors"
    "net/http"
)

// Define some common errors
var (
    ErrNotFound       = errors.New("record not found")
    ErrUnauthorized   = errors.New("unauthorized")
    ErrBadRequest     = errors.New("bad request")
    ErrConflict       = errors.New("data conflict")
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
    default:
        return http.StatusInternalServerError
    }
}

