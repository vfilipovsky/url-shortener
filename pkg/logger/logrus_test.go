package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestErrorLog(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.Error("error")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "error", hook.LastEntry().Message)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}
