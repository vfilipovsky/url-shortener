package api

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"

	"github.com/vfilipovsky/url-shortener/pkg/exception"
	"github.com/vfilipovsky/url-shortener/pkg/logger"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Errors []ErrorMessage `json:"errors"`
}

func handleError(err error) (any, int) {
	errResp := ErrorResponse{Errors: []ErrorMessage{}}

	switch e := err.(type) {
	case *exception.StatusError:
		errResp.Errors = append(errResp.Errors, ErrorMessage{Message: e.Error()})

		return errResp, e.StatusCode()
	case validator.ValidationErrors:
		for _, fieldError := range e {
			msg := fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fieldError.Field(), fieldError.Tag())

			errResp.Errors = append(errResp.Errors, ErrorMessage{Message: msg})
		}

		return errResp, http.StatusBadRequest
	default:
		msg := "internal service error"

		logger.Errorf("%s -> %s", msg, e.Error())

		errResp.Errors = append(errResp.Errors, ErrorMessage{Message: msg})

		return errResp, http.StatusInternalServerError
	}
}
