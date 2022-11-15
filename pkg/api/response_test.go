package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vfilipovsky/url-shortener/pkg/exception"
)

func TestRespond(t *testing.T) {
	type result struct {
		Foo string `json:"foo"`
	}

	w := httptest.NewRecorder()
	Respond(w, &result{Foo: "bar"}, http.StatusOK)

	res := w.Result()
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.Equal(t, "{\"foo\":\"bar\"}\n", string(body))
}

func TestRespondAnError(t *testing.T) {
	err := exception.StatusError{
		Message: "test error",
		Code:    400,
	}

	w := httptest.NewRecorder()
	Respond(w, &err)

	res := w.Result()
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.Equal(t, "{\"errors\":[{\"message\":\"test error\"}]}\n", string(body))
}
