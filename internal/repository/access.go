package repository

//go:generate mockgen -source=./access.go -destination=./mock/mock_access.go

import (
	"database/sql"

	"github.com/vfilipovsky/url-shortener/internal/entity"
)

type Access interface {
	FindByToken(token string) (*entity.Access, error)
	Create(access *entity.Access) error
}

type accessRepository struct {
	db *sql.DB
}

func NewAccess(db *sql.DB) Access {
	return &accessRepository{db: db}
}

func (r *accessRepository) Create(access *entity.Access) error {
	sql := `
	INSERT INTO accesses (id, token, is_active) 
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(
		sql,
		access.ID,
		access.Token,
		access.IsActive,
	)

	return err
}

func (r *accessRepository) FindByToken(token string) (*entity.Access, error) {
	access := entity.Access{}

	err := r.db.QueryRow(`
		SELECT 
			id, token, is_active, created_at, updated_at 
		FROM accesses 
		WHERE token = $1`, token).Scan(
		&access.ID,
		&access.Token,
		&access.IsActive,
		&access.CreatedAt,
		&access.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &access, nil
}
