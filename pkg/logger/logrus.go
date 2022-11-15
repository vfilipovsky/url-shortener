package logger

import (
	"github.com/sirupsen/logrus"

	"github.com/vfilipovsky/url-shortener/pkg/config"
)

var log *logrus.Logger

// NewInstance - creates a static instance of logger
func NewInstance(cfg *config.Logger) {
	log = logrus.New()
	log.SetLevel(cfg.Level)

	// setup sentry etc..
}

func Infof(format string, v ...any) {
	log.Infof(format, v...)
}

func Warnf(format string, v ...any) {
	log.Warnf(format, v...)
}

func Errorf(format string, v ...any) {
	log.Errorf(format, v...)
}

func Fatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}
