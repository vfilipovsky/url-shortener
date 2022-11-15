package api

import (
	"errors"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"

	"github.com/vfilipovsky/url-shortener/pkg/config"
	"github.com/vfilipovsky/url-shortener/pkg/exception"
	"github.com/vfilipovsky/url-shortener/pkg/logger"
)

func disableLogger() {
	logger.NewInstance(&config.Logger{
		Level: 0,
	})
	log.SetOutput(io.Discard)
}

func TestHandleErrorExceptionError(t *testing.T) {
	expectedErrorMessage := "user not found"
	expectedErrorCode := http.StatusNotFound

	statusError := &exception.StatusError{
		Code:    http.StatusNotFound,
		Message: expectedErrorMessage,
	}

	expectedErrorResp := ErrorResponse{Errors: []ErrorMessage{{
		Message: expectedErrorMessage,
	}}}

	actualResp, actualCode := handleError(statusError)

	assert.Equal(t, expectedErrorCode, actualCode)
	assert.Equal(t, expectedErrorResp, actualResp)
}

func TestHandlerErrorValidationError(t *testing.T) {
	type test struct {
		Name string `json:"name" validate:"required"`
	}

	expectedErrorMessage := "Field validation for 'Name' failed on the 'required' tag"
	expectedErrorCode := http.StatusBadRequest

	validate := validator.New()
	testStruct := &test{}
	err := validate.Struct(testStruct)

	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	expectedErrorResp := ErrorResponse{Errors: []ErrorMessage{{
		Message: expectedErrorMessage,
	}}}

	actualResp, actualCode := handleError(err)

	assert.Equal(t, expectedErrorCode, actualCode)
	assert.Equal(t, expectedErrorResp, actualResp)
}

func TestHandleErrorDefaultCase(t *testing.T) {
	disableLogger()

	expectedErrorMessage := "internal service error"
	expectedErrorCode := http.StatusInternalServerError

	err := errors.New("test name")

	expectedErrorResp := ErrorResponse{Errors: []ErrorMessage{{
		Message: expectedErrorMessage,
	}}}

	actualResp, actualCode := handleError(err)

	assert.Equal(t, expectedErrorCode, actualCode)
	assert.Equal(t, expectedErrorResp, actualResp)
}
