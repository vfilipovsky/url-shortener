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

func TestSuccessDeleteWithGivenValidAccessTokenAndUrlCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessService := mockService.NewMockAccess(ctrl)
	urlService := mockService.NewMockUrl(ctrl)

	accessToken := "zcgt2Q0s4U4F2MP0Zhxosq3xzqhvv4elZUaOF8wnWuqMwrmCWq"
	access := &entity.Access{
		Token:    accessToken,
		IsActive: true,
	}

	url := &entity.Url{
		ID: uuid.Must(uuid.NewUUID()),
	}

	accessService.EXPECT().GetByToken(accessToken).Return(access, nil)
	urlService.EXPECT().GetByCode("qwerty").Return(url, nil)
	urlService.EXPECT().Delete(url.ID).Return(nil)

	delete := &deleteAction{accessService: accessService, urlService: urlService}

	actual := delete.Run("qwerty", accessToken)

	assert.Nil(t, actual)
}

func TestAccessRestrictedIfAccessNotFoundByAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessService := mockService.NewMockAccess(ctrl)
	urlService := mockService.NewMockUrl(ctrl)

	accessToken := "zcgt2Q0s4U4F2MP0Zhxosq3xzqhvv4elZUaOF8wnWuqMwrmCWq"

	accessService.EXPECT().GetByToken(accessToken).Return(nil, nil)

	delete := &deleteAction{accessService: accessService, urlService: urlService}

	expected := exception.AccessRestricted
	actual := delete.Run("qwerty", accessToken)

	assert.ErrorIs(t, expected, actual)
}

func TestAccessRestrictedIfAccessIsNotActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessService := mockService.NewMockAccess(ctrl)
	urlService := mockService.NewMockUrl(ctrl)

	accessToken := "zcgt2Q0s4U4F2MP0Zhxosq3xzqhvv4elZUaOF8wnWuqMwrmCWq"
	access := &entity.Access{
		Token:    accessToken,
		IsActive: false,
	}

	accessService.EXPECT().GetByToken(accessToken).Return(access, nil)

	delete := &deleteAction{accessService: accessService, urlService: urlService}

	expected := exception.AccessRestricted
	actual := delete.Run("qwerty", accessToken)

	assert.ErrorIs(t, expected, actual)
}

func TestReturnNilIfUrlNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessService := mockService.NewMockAccess(ctrl)
	urlService := mockService.NewMockUrl(ctrl)

	accessToken := "zcgt2Q0s4U4F2MP0Zhxosq3xzqhvv4elZUaOF8wnWuqMwrmCWq"
	access := &entity.Access{
		Token:    accessToken,
		IsActive: true,
	}

	accessService.EXPECT().GetByToken(accessToken).Return(access, nil)
	urlService.EXPECT().GetByCode("qwerty").Return(nil, nil)

	delete := &deleteAction{accessService: accessService, urlService: urlService}

	actual := delete.Run("qwerty", accessToken)

	assert.Nil(t, actual)
}
