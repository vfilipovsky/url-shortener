package url

import (
	"github.com/vfilipovsky/url-shortener/internal/entity"
	"github.com/vfilipovsky/url-shortener/internal/service"
	"github.com/vfilipovsky/url-shortener/pkg/exception"
)

type GetAllByAccess interface {
	Run(accessToken string, offset, limit int) ([]*entity.Url, error)
}

type getAllUrlsByAccess struct {
	accessService service.Access
	urlService    service.Url
}

func NewGetAllUrlsByAccessAction(
	accessService service.Access,
	urlService service.Url,
) GetAllByAccess {
	return &getAllUrlsByAccess{
		accessService: accessService,
		urlService:    urlService,
	}
}

func (a *getAllUrlsByAccess) Run(accessToken string, offset, limit int) ([]*entity.Url, error) {
	access, err := a.accessService.GetByToken(accessToken)

	if err != nil {
		return nil, err
	}

	if access == nil || !access.IsActive {
		return nil, exception.AccessRestricted
	}

	urls, err := a.urlService.GetByAccessID(access.ID, offset, limit)

	if err != nil {
		return nil, err
	}

	return urls, nil
}
