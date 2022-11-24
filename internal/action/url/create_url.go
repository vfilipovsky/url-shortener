package url

//go:generate mockgen -source=./create_url.go -destination=./mock/mock_create_url.go

import (
	"time"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	"github.com/vfilipovsky/url-shortener/internal/service"
	"github.com/vfilipovsky/url-shortener/pkg/exception"
	"github.com/vfilipovsky/url-shortener/pkg/randomizer"
	"github.com/vfilipovsky/url-shortener/pkg/tokenizer"
)

const (
	codeLength = 12
	pinLength  = 4
)

type Create interface {
	Run(accessToken, url string, isSecured bool, aliveUntil time.Time) (*entity.Url, error)
}

type createAction struct {
	accessService service.Access
	urlService    service.Url
	randomizer    randomizer.Randomizer
	tokenizer     tokenizer.Tokenizer
}

func NewCreateAction(
	accessService service.Access,
	urlService service.Url,
	randomizer randomizer.Randomizer,
	tokenizer tokenizer.Tokenizer,
) Create {
	return &createAction{
		accessService: accessService,
		urlService:    urlService,
		randomizer:    randomizer,
		tokenizer:     tokenizer,
	}
}

func (a *createAction) Run(accessToken, url string, isSecured bool, aliveUntil time.Time) (*entity.Url, error) {
	access, err := a.accessService.GetByToken(accessToken)

	if err != nil {
		return nil, err
	}

	if access == nil || !access.IsActive {
		return nil, exception.AccessRestricted
	}

	id, err := a.tokenizer.Generate()

	if err != nil {
		return nil, err
	}

	newUrl := &entity.Url{
		ID:         id,
		Url:        url,
		Code:       a.randomizer.Random(codeLength, randomizer.Chars),
		AliveUntil: aliveUntil,
		AccessID:   access.ID,
	}

	if isSecured {
		newUrl.IsSecured = true
		newUrl.Pin = a.randomizer.Random(pinLength, randomizer.Numbers)
	}

	return newUrl, a.urlService.Create(newUrl)
}
