package url

//go:generate mockgen -source=./delete_url.go -destination=./mock/mock_delete_url.go

import (
	"github.com/vfilipovsky/url-shortener/internal/service"
	"github.com/vfilipovsky/url-shortener/pkg/exception"
)

type Delete interface {
	Run(code, accessToken string) error
}

type deleteAction struct {
	accessService service.Access
	urlService    service.Url
}

func NewDeleteAction(accessService service.Access, urlService service.Url) Delete {
	return &deleteAction{accessService: accessService, urlService: urlService}
}

func (a *deleteAction) Run(code, accessToken string) error {
	access, err := a.accessService.GetByToken(accessToken)

	if err != nil {
		return err
	}

	if access == nil || !access.IsActive {
		return exception.AccessRestricted
	}

	url, err := a.urlService.GetByCode(code)

	if err != nil {
		return err
	}

	if url == nil {
		return nil
	}

	if url.AccessID != access.ID {
		return exception.AccessRestricted
	}

	return a.urlService.Delete(url.ID)
}
