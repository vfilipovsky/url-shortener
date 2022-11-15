package tokenizer

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	id, err := New().Random()

	_, parseErr := uuid.Parse(id.String())

	assert.Nil(t, err)
	assert.Nil(t, parseErr)
}

func TestNewUUID(t *testing.T) {
	id, err := NewUUID()

	_, parseErr := uuid.Parse(id.String())

	assert.Nil(t, err)
	assert.Nil(t, parseErr)
}

func TestGenerate(t *testing.T) {
	id, err := New().Generate()

	_, parseErr := uuid.Parse(id.String())

	assert.Nil(t, err)
	assert.Nil(t, parseErr)
}

func TestParse(t *testing.T) {
	id := uuid.Must(uuid.NewUUID()).String()

	parsedUuid, err := New().Parse(id)

	assert.Nil(t, err)
	assert.Equal(t, id, parsedUuid.String())
}
