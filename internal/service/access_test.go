package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	mockRepository "github.com/vfilipovsky/url-shortener/internal/repository/mock"
)

func TestGetByTokenSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessRepository := mockRepository.NewMockAccess(ctrl)

	token := "qwerty"
	access := &entity.Access{Token: token}

	accessService := NewAccess(accessRepository)

	accessRepository.EXPECT().FindByToken(token).Return(access, nil)

	actual, err := accessService.GetByToken(token)

	assert.Nil(t, err)
	assert.Equal(t, access, actual)
}

func TestGetByTokenReturnsAnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessRepository := mockRepository.NewMockAccess(ctrl)

	token := "qwerty"
	var nilAccess *entity.Access

	accessService := NewAccess(accessRepository)

	accessRepository.EXPECT().FindByToken(token).Return(nilAccess, errors.New("repository error"))

	actual, err := accessService.GetByToken(token)

	assert.NotNil(t, err)
	assert.Equal(t, nilAccess, actual)
}
