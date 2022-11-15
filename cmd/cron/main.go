package main

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/vfilipovsky/url-shortener/internal/cronjob"
	"github.com/vfilipovsky/url-shortener/internal/repository"
	"github.com/vfilipovsky/url-shortener/internal/service"
	"github.com/vfilipovsky/url-shortener/pkg/config"
	"github.com/vfilipovsky/url-shortener/pkg/cron"
	"github.com/vfilipovsky/url-shortener/pkg/database"
	"github.com/vfilipovsky/url-shortener/pkg/logger"
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

	s := cron.NewGocron(db)

	registerJobs(s, db, cfg)

	s.Start()
}

func registerJobs(runner cron.Cron, db *sql.DB, cfg *config.Config) {
	urlRepository := repository.NewUrl(db)
	urlService := service.NewUrl(urlRepository)

	runner.Schedule(cron.EVERY_MINUTE, cronjob.NewClearOverDueUrls(urlService))
}
