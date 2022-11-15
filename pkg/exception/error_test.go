package exception

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusError(t *testing.T) {
	expectedCode := http.StatusBadRequest
	expectedMessage := "validation error"

	err := &StatusError{
		Code:    http.StatusBadRequest,
		Message: expectedMessage,
	}

	assert.Equal(t, expectedCode, err.StatusCode())
	assert.Equal(t, expectedMessage, err.Error())
}
