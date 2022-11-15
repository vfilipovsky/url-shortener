package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/vfilipovsky/url-shortener/pkg/config"
	"github.com/vfilipovsky/url-shortener/pkg/logger"
)

func NewInstance(cfg *config.DatabaseCredentials) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SslMode)

	driver := "postgres"
	conn, err := sql.Open(driver, dsn)

	if err != nil {
		return nil, err
	}

	logger.Infof("Database connected: (%s) %s:%s", driver, cfg.Host, cfg.Port)

	return conn, err
}
