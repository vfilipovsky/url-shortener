package main

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	"github.com/vfilipovsky/url-shortener/internal/repository"
	"github.com/vfilipovsky/url-shortener/pkg/config"
	"github.com/vfilipovsky/url-shortener/pkg/database"
	"github.com/vfilipovsky/url-shortener/pkg/logger"
	"github.com/vfilipovsky/url-shortener/pkg/randomizer"
	"github.com/vfilipovsky/url-shortener/pkg/tokenizer"
)

func main() {
	var err error

	if err = godotenv.Load(".env"); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	cfg := config.Init()

	logger.NewInstance(cfg.Logger)

	db, err := database.NewInstance(cfg.DbCredentials)

	if err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			logger.Errorf("Failed to close database connection: %s", err)
			return
		}

		logger.Infof("Database connection closed")
	}(db)

	random := randomizer.New()
	accessRepository := repository.NewAccess(db)
	accessToken := random.Random(50, randomizer.Chars)
	id, err := tokenizer.NewUUID()

	if err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}

	if err = accessRepository.Create(&entity.Access{
		ID:       id,
		Token:    accessToken,
		IsActive: true,
	}); err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}

	logger.Infof("Access successfully created")
	logger.Infof(accessToken)
}
