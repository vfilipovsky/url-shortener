package service

//go:generate mockgen -source=./url.go -destination=./mock/mock_url.go

import (
	"github.com/google/uuid"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	"github.com/vfilipovsky/url-shortener/internal/repository"
)

type Url interface {
	RemoveOverdue() error
	GetByCode(code string) (*entity.Url, error)
	GetByAccessID(accessID uuid.UUID, offset, limit int) ([]*entity.Url, error)
	Create(url *entity.Url) error
	Delete(ID uuid.UUID) error
}

type urlService struct {
	repo repository.Url
}

func NewUrl(repo repository.Url) Url {
	return &urlService{repo: repo}
}

func (s *urlService) RemoveOverdue() error {
	return s.repo.RemoveOverdue()
}

func (s *urlService) GetByAccessID(id uuid.UUID, offset, limit int) ([]*entity.Url, error) {
	return s.repo.FindByAccessID(id, offset, limit)
}

func (s *urlService) GetByCode(code string) (*entity.Url, error) {
	return s.repo.FindByCode(code)
}

func (s *urlService) Create(url *entity.Url) error {
	return s.repo.Create(url)
}

func (s *urlService) Delete(ID uuid.UUID) error {
	return s.repo.Delete(ID)
}
