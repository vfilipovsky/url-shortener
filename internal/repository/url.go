package repository

//go:generate mockgen -source=./url.go -destination=./mock/mock_url.go

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/vfilipovsky/url-shortener/internal/entity"
)

type Url interface {
	RemoveOverdue() error
	FindByCode(code string) (*entity.Url, error)
	FindByAccessID(accessID uuid.UUID, offset, limit int) ([]*entity.Url, error)
	Create(url *entity.Url) error
	Delete(ID uuid.UUID) error
}

type urlRepository struct {
	db *sql.DB
}

func NewUrl(db *sql.DB) Url {
	return &urlRepository{db: db}
}

func (r *urlRepository) RemoveOverdue() error {
	if _, err := r.db.Exec("DELETE FROM urls WHERE alive_until < NOW()"); err != nil {
		return err
	}

	return nil
}

func (r *urlRepository) FindByCode(code string) (*entity.Url, error) {
	var url entity.Url

	s := `
	SELECT 
		id,
		code,
		url,
		is_secured,
		pin,
		alive_until,
		access_id,
		created_at,
		updated_at
	FROM urls
	WHERE code = $1
`

	err := r.db.QueryRow(s, code).Scan(
		&url.ID,
		&url.Code,
		&url.Url,
		&url.IsSecured,
		&url.Pin,
		&url.AliveUntil,
		&url.AccessID,
		&url.CreatedAt,
		&url.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &url, nil
}

func (r *urlRepository) FindByAccessID(accessID uuid.UUID, offset, limit int) ([]*entity.Url, error) {
	s := `
		SELECT 
			id,
			code,
			url,
			is_secured,
			pin,
			alive_until,
			access_id,
			created_at,
			updated_at
		FROM urls
		WHERE access_id = $1
		LIMIT $2 OFFSET $3
	`

	if limit == 0 {
		limit = DefaultLimit
	}

	rows, err := r.db.Query(s, accessID.String(), limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var urls []*entity.Url

	for rows.Next() {
		var url entity.Url

		if err = rows.Scan(
			&url.ID,
			&url.Code,
			&url.Url,
			&url.IsSecured,
			&url.Pin,
			&url.AliveUntil,
			&url.AccessID,
			&url.CreatedAt,
			&url.UpdatedAt,
		); err != nil {
			if err == sql.ErrNoRows {
				return urls, nil
			}

			return urls, err
		}

		urls = append(urls, &url)
	}

	return urls, nil
}

func (r *urlRepository) Create(url *entity.Url) error {
	sql := `
		INSERT INTO urls (id, code, url, is_secured, pin, alive_until, access_id) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		sql,
		url.ID,
		url.Code,
		url.Url,
		url.IsSecured,
		url.Pin,
		url.AliveUntil,
		url.AccessID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *urlRepository) Delete(ID uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM urls WHERE id = $1", ID.String())

	if err != nil {
		return err
	}

	return nil
}
