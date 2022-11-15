package tokenizer

//go:generate mockgen -source=./uuid.go -destination=./mock/mock_uuid.go

import (
	"github.com/google/uuid"

	"github.com/vfilipovsky/url-shortener/pkg/exception"
)

type Tokenizer interface {
	Parse(token string) (uuid.UUID, error)
	Generate() (uuid.UUID, error)
	Random() (uuid.UUID, error)
}

type tokenizer struct{}

func New() Tokenizer {
	return &tokenizer{}
}

func NewUUID() (uuid.UUID, error) {
	return uuid.NewUUID()
}

func (t *tokenizer) Parse(token string) (uuid.UUID, error) {
	id, err := uuid.Parse(token)

	if err != nil {
		return uuid.Nil, exception.BadRequest("invalid ID")
	}

	return id, nil
}

func (t *tokenizer) Generate() (uuid.UUID, error) {
	return uuid.NewUUID()
}

func (t *tokenizer) Random() (uuid.UUID, error) {
	return uuid.NewRandom()
}
