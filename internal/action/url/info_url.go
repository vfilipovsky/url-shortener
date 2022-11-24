package url

//go:generate mockgen -source=./info_url.go -destination=./mock/mock_info_url.go

import (
	"fmt"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	"github.com/vfilipovsky/url-shortener/internal/service"
	"github.com/vfilipovsky/url-shortener/pkg/exception"
)

type Info interface {
	Run(code string) (*entity.Url, error)
}

type infoAction struct {
	urlService service.Url
}

func NewInfoAction(urlService service.Url) Info {
	return &infoAction{urlService: urlService}
}

func (a *infoAction) Run(code string) (*entity.Url, error) {
	url, err := a.urlService.GetByCode(code)

	if err != nil {
		return nil, err
	}

	if url == nil {
		return nil, exception.NotFound(fmt.Sprintf("url not found by code '%s'", code))
	}

	return url, nil
}
