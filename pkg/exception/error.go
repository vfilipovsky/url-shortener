package exception

import "net/http"

var (
	AccessRestricted = Forbidden("access restricted")
	EmptyBody        = BadRequest("empty body")
)

type Error interface {
	Code() int
	Message() string
}

type StatusError struct {
	Message string
	Code    int
}

func (se *StatusError) StatusCode() int {
	return se.Code
}

func (se *StatusError) Error() string {
	return se.Message
}

func InternalServiceError(message string) *StatusError {
	return &StatusError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func BadRequest(message string) *StatusError {
	return &StatusError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NotFound(message string) *StatusError {
	return &StatusError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func Unauthorized(message string) *StatusError {
	return &StatusError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func Forbidden(message string) *StatusError {
	return &StatusError{
		Code:    http.StatusForbidden,
		Message: message,
	}
}
