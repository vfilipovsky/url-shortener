package service

//go:generate mockgen -source=./access.go -destination=./mock/mock_access.go

import (
	"github.com/vfilipovsky/url-shortener/internal/entity"
	"github.com/vfilipovsky/url-shortener/internal/repository"
)

type Access interface {
	GetByToken(token string) (*entity.Access, error)
}

type accessService struct {
	repo repository.Access
}

func NewAccess(repo repository.Access) Access {
	return &accessService{repo: repo}
}

func (s *accessService) GetByToken(token string) (*entity.Access, error) {
	return s.repo.FindByToken(token)
}
