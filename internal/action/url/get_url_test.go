package url

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	mockService "github.com/vfilipovsky/url-shortener/internal/service/mock"
	"github.com/vfilipovsky/url-shortener/pkg/exception"
)

func TestSuccessGetUrlByValidPinAndCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlService := mockService.NewMockUrl(ctrl)

	url := &entity.Url{
		Url: "https://google.com",
	}

	getUrl := &getUrlAction{urlService: urlService}
	urlService.EXPECT().GetByCode("qwerty").Return(url, nil)

	expected := url.Url
	actual, err := getUrl.Run("1234", "qwerty")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestSuccessUnsecuredUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlService := mockService.NewMockUrl(ctrl)

	getUrl := &getUrlAction{urlService: urlService}

	url := &entity.Url{
		IsSecured: false,
	}

	urlService.EXPECT().GetByCode("qwerty").Return(url, nil)

	expected := url.Url
	actual, err := getUrl.Run("", "qwerty")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestUrlNotFoundError(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlService := mockService.NewMockUrl(ctrl)

	getUrl := &getUrlAction{urlService: urlService}
	urlService.EXPECT().GetByCode("qwerty").Return(nil, nil)

	expected := ""
	expectedErr := exception.NotFound("url not found")
	actual, err := getUrl.Run("1234", "qwerty")

	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedErr.Code, err.(*exception.StatusError).Code)
	assert.EqualError(t, expectedErr, err.Error())
}

func TestAccessRestrictedOnWrongPin(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlService := mockService.NewMockUrl(ctrl)

	url := &entity.Url{
		IsSecured: true,
	}

	getUrl := &getUrlAction{urlService: urlService}
	urlService.EXPECT().GetByCode("qwerty").Return(url, nil)

	expected := ""
	expectedErr := exception.AccessRestricted
	actual, err := getUrl.Run("1235", "qwerty")

	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedErr.Code, err.(*exception.StatusError).Code)
	assert.EqualError(t, expectedErr, err.Error())
}
