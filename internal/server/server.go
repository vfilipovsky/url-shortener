package server

import (
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/vfilipovsky/url-shortener/pkg/config"
	"github.com/vfilipovsky/url-shortener/pkg/logger"
)

type server struct {
	config *config.Config
	db     *sql.DB
	quit   chan os.Signal
}

func New(cfg *config.Config, db *sql.DB) server {
	return server{config: cfg, db: db}
}

func (s server) Run() {
	r := register(mux.NewRouter(), s.db, s.config)

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"Access-Control-Allow-Origin", "Content-Type", "Access-Control-Allow-Methods"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{http.MethodDelete, http.MethodPost, http.MethodGet}),
	)

	go func() {
		if err := http.ListenAndServe(":"+s.config.Http.Port, cors(r)); err != nil {
			logger.Fatalf(err.Error())
			os.Exit(1)
		}
	}()

	logger.Infof("Listening on port :%s", s.config.Http.Port)
}

func (s server) Shutdown() chan os.Signal {
	s.quit = make(chan os.Signal, 1)
	signal.Notify(s.quit, syscall.SIGTERM, syscall.SIGINT)

	return s.quit
}
