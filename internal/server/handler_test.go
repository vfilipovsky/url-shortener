package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/v1/ping", nil)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ping(w, r)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "\"pong\"", strings.TrimSpace(w.Body.String()))
}
