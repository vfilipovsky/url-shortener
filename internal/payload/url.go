package payload

import "time"

type CreateUrl struct {
	AccessToken string    `json:"access_token" validate:"required"`
	Url         string    `json:"url" validate:"required"`
	AliveUntil  time.Time `json:"alive_until" validate:"required"`
	IsSecured   *bool     `json:"is_secured" validate:"required"`
}

type DeleteUrl struct {
	AccessToken string `json:"access_token" validate:"required"`
	Code        string `json:"code" validate:"required"`
}

type GetUrl struct {
	Code string `json:"code" validate:"required"`
	Pin  string `json:"pin" validate:"omitempty"`
}

type GetUrlsByAccess struct {
	AccessToken string `json:"access_token" validate:"required"`
	Limit       int    `json:"limit" validate:"omitempty,gte=1" schema:"limit"`
	Offset      int    `json:"offset" validate:"omitempty,gte=0" schema:"offset"`
}
