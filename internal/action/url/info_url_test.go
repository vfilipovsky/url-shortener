package url

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	mockService "github.com/vfilipovsky/url-shortener/internal/service/mock"
	"github.com/vfilipovsky/url-shortener/pkg/exception"
)

func TestSuccessInfoUrlOnValidCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlService := mockService.NewMockUrl(ctrl)

	url := &entity.Url{
		Url: "https://google.com",
	}

	info := &infoAction{urlService: urlService}
	urlService.EXPECT().GetByCode("qwerty").Return(url, nil)

	expected := url
	actual, err := info.Run("qwerty")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestNotFoundExceptionIfUrlNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlService := mockService.NewMockUrl(ctrl)

	code := "qwerty"
	info := &infoAction{urlService: urlService}
	urlService.EXPECT().GetByCode(code).Return(nil, nil)

	expected := exception.NotFound(fmt.Sprintf("url not found by code '%s'", code))
	actual, err := info.Run("qwerty")

	assert.Nil(t, actual)
	assert.Equal(t, expected.Code, err.(*exception.StatusError).Code)
	assert.Equal(t, expected.Message, err.(*exception.StatusError).Message)
}
