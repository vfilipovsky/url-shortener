package url

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	mockService "github.com/vfilipovsky/url-shortener/internal/service/mock"
	"github.com/vfilipovsky/url-shortener/pkg/exception"
)

func TestSuccessGetUrlsByAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessService := mockService.NewMockAccess(ctrl)
	urlService := mockService.NewMockUrl(ctrl)

	accessToken := "zcgt2Q0s4U4F2MP0Zhxosq3xzqhvv4elZUaOF8wnWuqMwrmCWq"
	access := &entity.Access{
		Token:    accessToken,
		IsActive: true,
	}

	accessService.EXPECT().GetByToken(accessToken).Return(access, nil)

	urls := []*entity.Url{
		{
			ID: uuid.Must(uuid.NewUUID()),
		},
		{
			ID: uuid.Must(uuid.NewUUID()),
		},
	}

	urlService.EXPECT().GetByAccessID(access.ID, 0, 100).Return(urls, nil)

	getUrls := &getAllUrlsByAccess{accessService: accessService, urlService: urlService}

	expected := urls
	actual, err := getUrls.Run(accessToken, 0, 100)

	assert.Nil(t, err)
	assert.Equal(t, expected[0].ID, actual[0].ID)

}

func TestGetUrlsAccessRestrictedIfAccessNotFoundByAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessService := mockService.NewMockAccess(ctrl)
	urlService := mockService.NewMockUrl(ctrl)

	accessToken := "zcgt2Q0s4U4F2MP0Zhxosq3xzqhvv4elZUaOF8wnWuqMwrmCWq"

	accessService.EXPECT().GetByToken(accessToken).Return(nil, nil)

	getUrls := &getAllUrlsByAccess{accessService: accessService, urlService: urlService}

	expected := exception.AccessRestricted
	actual, err := getUrls.Run(accessToken, 0, 100)

	assert.Nil(t, actual)
	assert.ErrorIs(t, expected, err)
}

func TestGetUrlsAccessRestrictedIfAccessIsNotActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessService := mockService.NewMockAccess(ctrl)
	urlService := mockService.NewMockUrl(ctrl)

	accessToken := "zcgt2Q0s4U4F2MP0Zhxosq3xzqhvv4elZUaOF8wnWuqMwrmCWq"
	access := &entity.Access{
		Token:    accessToken,
		IsActive: false,
	}

	accessService.EXPECT().GetByToken(accessToken).Return(access, nil)

	getUrls := &getAllUrlsByAccess{accessService: accessService, urlService: urlService}

	expected := exception.AccessRestricted
	actual, err := getUrls.Run(accessToken, 0, 100)

	assert.Nil(t, actual)
	assert.ErrorIs(t, expected, err)
}
