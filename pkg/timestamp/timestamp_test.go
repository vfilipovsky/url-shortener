package timestamp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp(t *testing.T) {
	createdAt := time.Now()
	updatedAt := time.Now()

	actual := Timestamp{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	assert.Equal(t, createdAt, actual.CreatedAt)
	assert.Equal(t, updatedAt, actual.UpdatedAt)
}
