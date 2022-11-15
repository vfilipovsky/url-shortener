package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vfilipovsky/url-shortener/pkg/exception"
)

type request struct {
	Code string `validate:"required" schema:"foo"`
}

func TestValidateRequestReturnsNoErrorOnValidRequest(t *testing.T) {
	req := &request{
		Code: "code",
	}

	actual := ValidateRequestPayload(req)

	assert.Nil(t, actual)
}

func TestValidateRequestReturnsErrorOnInvalidRequest(t *testing.T) {
	actual := ValidateRequestPayload(&request{})

	expected := "Key: 'request.Code' Error:Field validation for 'Code' failed on the 'required' tag"
	assert.NotNil(t, actual)
	assert.EqualError(t, actual, expected)
}

func TestParseQueryParamsReturnsNilOnValidQueryParam(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test?foo=bar", nil)

	type parseForm struct {
		Foo string `schema:"foo"`
	}

	pf := &parseForm{}
	expected := "bar"
	actual := ParseQueryParams(r, pf)

	assert.Nil(t, actual)
	assert.Equal(t, expected, pf.Foo)
}

func TestParseQueryParamsReturnsErrorOnInvalidQueryParam(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test?bar=foo", nil)

	type parseForm struct {
		Foo string `schema:"foo"`
	}

	pf := &parseForm{}
	actual := ParseQueryParams(r, pf)

	assert.NotNil(t, actual)
	assert.EqualError(t, actual, "schema: invalid path \"bar\"")
}

func TestDecodeAndValidateReturnsNilOnValidBody(t *testing.T) {
	type payload struct {
		Foo string `json:"foo"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"foo":"bar"}`))
	pl := &payload{}

	expected := "bar"
	actual := DecodeAndValidate(r.Body, pl)

	assert.Nil(t, actual)
	assert.Equal(t, expected, pl.Foo)
}

func TestDecodeAndValidateReturnsEmptyBodyErrOnEmptyJson(t *testing.T) {
	type payload struct {
		Foo string `json:"foo"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	pl := &payload{}

	expected := ""
	expectedErr := exception.EmptyBody
	actual := DecodeAndValidate(r.Body, pl)

	assert.NotNil(t, actual)
	assert.ErrorIs(t, expectedErr, actual)
	assert.Equal(t, expectedErr.Code, actual.(*exception.StatusError).Code)
	assert.Equal(t, expectedErr.Message, actual.(*exception.StatusError).Message)
	assert.Equal(t, expected, pl.Foo)
}

func TestDecodeAndValidateReturnsJsonErrorOnInvalidJson(t *testing.T) {
	type payload struct {
		Foo string `json:"foo"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid json"))
	pl := &payload{}

	expected := ""
	expectedErr := "invalid character 'i' looking for beginning of value"
	actual := DecodeAndValidate(r.Body, pl)

	assert.NotNil(t, actual)
	assert.EqualError(t, actual, expectedErr)
	assert.Equal(t, expected, pl.Foo)
}

func TestDecodeAndValidateReturnsValidateErrorIfRequestIsInvalid(t *testing.T) {
	type payload struct {
		Foo string `json:"foo" validate:"required,gte=15"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"foo":"bar"}`))
	pl := &payload{}
	expectedErr := "Key: 'payload.Foo' Error:Field validation for 'Foo' failed on the 'gte' tag"
	actual := DecodeAndValidate(r.Body, pl)

	assert.NotNil(t, actual)
	assert.EqualError(t, actual, expectedErr)
}
