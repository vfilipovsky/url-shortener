package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/v1/ping", nil)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ping(w, r)

	var result string
	err := json.NewDecoder(w.Result().Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", result)
}
