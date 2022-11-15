package entity

import (
	"github.com/google/uuid"

	"github.com/vfilipovsky/url-shortener/pkg/timestamp"
)

type Access struct {
	ID uuid.UUID `json:"id"`

	Token    string `json:"token"`
	IsActive bool   `json:"is_active"`

	timestamp.Timestamp
}
