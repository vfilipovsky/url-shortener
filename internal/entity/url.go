package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/vfilipovsky/url-shortener/pkg/timestamp"
)

type Url struct {
	ID         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Url        string    `json:"url"`
	IsSecured  bool      `json:"is_secured"`
	Pin        string    `json:"pin"`
	AliveUntil time.Time `json:"alive_until"`

	AccessID uuid.UUID `json:"-"`
	Access   *Access   `json:"-"`

	timestamp.Timestamp
}
