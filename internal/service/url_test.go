package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	mockRepository "github.com/vfilipovsky/url-shortener/internal/repository/mock"
)

func TestRemoveOverdueReturnsNilOnSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlRepository := mockRepository.NewMockUrl(ctrl)

	urlRepository.EXPECT().RemoveOverdue().Return(nil)

	urlService := NewUrl(urlRepository)

	actual := urlService.RemoveOverdue()

	assert.Nil(t, actual)
}

func TestRemoveOverdueReturnsAnErrorOnFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlRepository := mockRepository.NewMockUrl(ctrl)

	urlRepository.EXPECT().RemoveOverdue().Return(errors.New("repository error"))

	urlService := NewUrl(urlRepository)

	actual := urlService.RemoveOverdue()

	assert.NotNil(t, actual)
	assert.Error(t, actual, "repository error")

}

func TestGetByAccessIDReturnsNilOnSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlRepository := mockRepository.NewMockUrl(ctrl)

	id := uuid.Must(uuid.NewUUID())
	urls := []*entity.Url{
		{
			Code: "qwerty",
		},
		{
			Code: "1234",
		},
	}

	urlRepository.EXPECT().FindByAccessID(id, 0, 100).Return(urls, nil)

	urlService := NewUrl(urlRepository)

	actual, err := urlService.GetByAccessID(id, 0, 100)

	assert.Nil(t, err)
	assert.Equal(t, urls, actual)
	assert.Equal(t, urls[0].Code, actual[0].Code)
	assert.Equal(t, urls[1].Code, actual[1].Code)
}

func TestGetByAccessIDReturnsAnErrorOnFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlRepository := mockRepository.NewMockUrl(ctrl)

	id := uuid.Must(uuid.NewUUID())
	urlRepository.EXPECT().FindByAccessID(id, 0, 100).Return(nil, errors.New("repository error"))

	urlService := NewUrl(urlRepository)

	actual, err := urlService.GetByAccessID(id, 0, 100)

	assert.NotNil(t, err)
	assert.Nil(t, actual)
	assert.Error(t, err, "repository error")
}

func TestGetByCodeReturnsUrlOnSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlRepository := mockRepository.NewMockUrl(ctrl)

	code := "qwerty"
	url := &entity.Url{Code: "qwerty"}

	urlRepository.EXPECT().FindByCode(code).Return(url, nil)

	urlService := NewUrl(urlRepository)

	actual, err := urlService.GetByCode(code)

	assert.Nil(t, err)
	assert.Equal(t, url, actual)
	assert.Equal(t, url.Code, actual.Code)
	assert.Equal(t, url.Code, actual.Code)
}

func TestGetByCodeReturnsAnErrorOnFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlRepository := mockRepository.NewMockUrl(ctrl)

	code := "qwerty"
	urlRepository.EXPECT().FindByCode(code).Return(nil, errors.New("repository error"))

	urlService := NewUrl(urlRepository)

	actual, err := urlService.GetByCode(code)

	assert.NotNil(t, err)
	assert.Nil(t, actual)
	assert.Error(t, err, "repository error")
}

func TestSaveReturnNilOnSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlRepository := mockRepository.NewMockUrl(ctrl)

	id := uuid.Must(uuid.NewUUID())
	url := &entity.Url{ID: id, Code: "qwerty"}

	urlRepository.EXPECT().Create(url).Return(nil)
	urlService := NewUrl(urlRepository)

	actual := urlService.Create(url)

	assert.Nil(t, actual)
	assert.Equal(t, url.Code, "qwerty")
}

func TestSaveReturnAnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlRepository := mockRepository.NewMockUrl(ctrl)

	id := uuid.Must(uuid.NewUUID())
	url := &entity.Url{ID: id, Code: "qwerty"}

	urlRepository.EXPECT().Create(url).Return(errors.New("repository error"))
	urlService := NewUrl(urlRepository)

	actual := urlService.Create(url)

	assert.NotNil(t, actual)
	assert.Error(t, actual, "repository error")
}

func TestDeleteReturnNilOnSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlRepository := mockRepository.NewMockUrl(ctrl)

	id := uuid.Must(uuid.NewUUID())

	urlRepository.EXPECT().Delete(id).Return(nil)
	urlService := NewUrl(urlRepository)

	actual := urlService.Delete(id)

	assert.Nil(t, actual)
}

func TestDeleteReturnAnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	urlRepository := mockRepository.NewMockUrl(ctrl)

	id := uuid.Must(uuid.NewUUID())

	urlRepository.EXPECT().Delete(id).Return(errors.New("repository error"))
	urlService := NewUrl(urlRepository)

	actual := urlService.Delete(id)

	assert.NotNil(t, actual)
	assert.Error(t, actual, "repository error")
}
