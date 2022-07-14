package httputil

import (
	"fmt"
	"net/http"
)

var (
	ErrBadRequest          = NewAPIError(http.StatusBadRequest, "bad_request", nil)
	ErrNotFound            = NewAPIError(http.StatusNotFound, "not_found", nil)
	ErrConflict            = NewAPIError(http.StatusConflict, "conflict", nil)
	ErrUnprocessableEntity = NewAPIError(http.StatusUnprocessableEntity, "entity_validation", nil)
	ErrInternalServerError = NewAPIError(http.StatusInternalServerError, "server_error", nil)
	ErrUnauthorized        = NewAPIError(http.StatusUnauthorized, "unauthorized", nil)
	ErrForbidden           = NewAPIError(http.StatusForbidden, "forbidden", nil)
)

var Errors = map[int]*APIError{
	http.StatusBadRequest:          ErrBadRequest,
	http.StatusNotFound:            ErrNotFound,
	http.StatusConflict:            ErrConflict,
	http.StatusUnprocessableEntity: ErrUnprocessableEntity,
	http.StatusInternalServerError: ErrInternalServerError,
	http.StatusUnauthorized:        ErrUnauthorized,
	http.StatusForbidden:           ErrForbidden,
}

type APIError struct {
	StatusCode int    `json:"status"`
	ErrorCode  string `json:"code"`
	Message    string `json:"message"`
	Errors     any    `json:"errors,omitempty"`
}

func NewAPIError(status int, code string, errors any) *APIError {
	return &APIError{
		StatusCode: status,
		Message:    http.StatusText(status),
		ErrorCode:  code,
		Errors:     errors,
	}
}

func (e *APIError) SetErrors(errors any) *APIError {
	e.Errors = errors
	return e
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%d][%s] %s %v", e.StatusCode, e.ErrorCode, e.Message, e.Errors)
}
