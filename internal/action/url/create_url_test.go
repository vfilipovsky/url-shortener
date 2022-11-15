package url

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	"github.com/vfilipovsky/url-shortener/internal/payload"
	mockService "github.com/vfilipovsky/url-shortener/internal/service/mock"
	"github.com/vfilipovsky/url-shortener/pkg/exception"
	"github.com/vfilipovsky/url-shortener/pkg/randomizer"
	mockRandomizer "github.com/vfilipovsky/url-shortener/pkg/randomizer/mock"
	mockTokenizer "github.com/vfilipovsky/url-shortener/pkg/tokenizer/mock"
)

func TestSuccessCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessService := mockService.NewMockAccess(ctrl)
	urlService := mockService.NewMockUrl(ctrl)
	random := mockRandomizer.NewMockRandomizer(ctrl)
	tokenizer := mockTokenizer.NewMockTokenizer(ctrl)

	id := uuid.Must(uuid.NewUUID())
	accessToken := "zcgt2Q0s4U4F2MP0Zhxosq3xzqhvv4elZUaOF8wnWuqMwrmCWq"
	access := &entity.Access{
		Token:    accessToken,
		IsActive: true,
	}

	pl := &payload.CreateUrl{
		AccessToken: accessToken,
		Url:         "https://google.com",
		AliveUntil:  time.Now().Add(time.Duration(time.Now().Weekday())),
		IsSecured:   true,
	}

	newUrl := &entity.Url{
		Url:        pl.Url,
		AliveUntil: pl.AliveUntil,
		IsSecured:  pl.IsSecured,
		AccessID:   access.ID,
		Code:       "qwe1234",
		Pin:        "1234",
	}

	accessService.EXPECT().GetByToken(accessToken).Return(access, nil)
	tokenizer.EXPECT().Generate().Return(id, nil)
	random.EXPECT().Random(codeLength, randomizer.Chars).Return("qwe1234")
	random.EXPECT().Random(pinLength, randomizer.Numbers).Return("1234")

	newUrl.ID = id

	urlService.EXPECT().Create(newUrl).Return(nil)

	create := &createAction{
		accessService: accessService,
		urlService:    urlService,
		randomizer:    random,
		tokenizer:     tokenizer,
	}

	actual, err := create.Run(pl.AccessToken, pl.Url, pl.IsSecured, pl.AliveUntil)

	assert.Nil(t, err)
	assert.Equal(t, actual.AccessID, newUrl.AccessID)
}

func TestAccessRestrictedDueToInvalidAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessService := mockService.NewMockAccess(ctrl)

	accessToken := "zcgt2Q0s4U4F2MP0Zhxosq3xzqhvv4elZUaOF8wnWuqMwrmCWq"
	pl := &payload.CreateUrl{
		AccessToken: accessToken,
		Url:         "https://google.com",
		AliveUntil:  time.Now().Add(time.Duration(time.Now().Weekday())),
		IsSecured:   true,
	}

	accessService.EXPECT().GetByToken(accessToken).Return(nil, nil)

	create := &createAction{accessService: accessService}

	expectedErr := exception.AccessRestricted
	actual, err := create.Run(pl.AccessToken, pl.Url, pl.IsSecured, pl.AliveUntil)

	assert.Nil(t, actual)
	assert.ErrorIs(t, expectedErr, err)
}

func TestAccessRestrictedDueToAccessNotActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessService := mockService.NewMockAccess(ctrl)

	accessToken := "zcgt2Q0s4U4F2MP0Zhxosq3xzqhvv4elZUaOF8wnWuqMwrmCWq"

	access := &entity.Access{
		IsActive: false,
		Token:    accessToken,
	}

	pl := &payload.CreateUrl{
		AccessToken: accessToken,
		Url:         "https://google.com",
		AliveUntil:  time.Now().Add(time.Duration(time.Now().Weekday())),
		IsSecured:   true,
	}

	accessService.EXPECT().GetByToken(accessToken).Return(access, nil)

	create := &createAction{accessService: accessService}

	expectedErr := exception.AccessRestricted
	actual, err := create.Run(pl.AccessToken, pl.Url, pl.IsSecured, pl.AliveUntil)

	assert.Nil(t, actual)
	assert.ErrorIs(t, expectedErr, err)
}
