package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/vfilipovsky/url-shortener/internal/server"
	"github.com/vfilipovsky/url-shortener/pkg/config"
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

	server := server.New(cfg, db)
	server.Run()
	<-server.Shutdown()

	fmt.Println()
}
