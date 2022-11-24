package url

//go:generate mockgen -source=./get_url.go -destination=./mock/mock_get_url.go

import (
	"github.com/vfilipovsky/url-shortener/internal/service"
	"github.com/vfilipovsky/url-shortener/pkg/exception"
)

type Get interface {
	Run(pin, code string) (string, error)
}

type getUrlAction struct {
	urlService service.Url
}

func NewGetAction(urlService service.Url) Get {
	return &getUrlAction{urlService: urlService}
}

func (a *getUrlAction) Run(pin, code string) (string, error) {
	url, err := a.urlService.GetByCode(code)

	if err != nil {
		return "", err
	}

	if url == nil {
		return "", exception.NotFound("url not found")
	}

	if !url.IsSecured {
		return url.Url, nil
	}

	if url.Pin == pin {
		return url.Url, nil
	}

	return "", exception.AccessRestricted
}
