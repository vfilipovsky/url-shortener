package handler

import (
	"os"

	"github.com/vfilipovsky/url-shortener/pkg/logger"
)

const (
	V1ApiPrefix = "/api/v1"

	V1Url = V1ApiPrefix + "/url"
)

func GetPublicDir() string {
	dir, err := os.Getwd()

	if err != nil {
		logger.Errorf(err.Error())
		return ""
	}

	return dir + "/public"
}
