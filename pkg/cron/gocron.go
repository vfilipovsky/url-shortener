package cron

import (
	"database/sql"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/vfilipovsky/url-shortener/pkg/logger"
)

type gocronService struct {
	db        *sql.DB
	scheduler *gocron.Scheduler
}

func NewGocron(db *sql.DB) Cron {
	s := gocron.NewScheduler(time.UTC)

	return &gocronService{
		db:        db,
		scheduler: s,
	}
}

func (g *gocronService) Schedule(cronExpression Expression, job Job) {
	_, err := g.scheduler.Cron(string(cronExpression)).Do(job.Run)
	if err != nil {
		logger.Errorf("Job '%s' register failed: %s", job.GetName(), err.Error())
		return
	}

	logger.Infof("Job '%s' registered", job.GetName())
}

func (g gocronService) Start() {
	g.scheduler.StartBlocking()
}
