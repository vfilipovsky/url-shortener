package server

import (
	"database/sql"
	"strings"

	"github.com/gorilla/mux"

	"github.com/vfilipovsky/url-shortener/internal/action/url"
	"github.com/vfilipovsky/url-shortener/internal/handler"
	"github.com/vfilipovsky/url-shortener/internal/repository"
	"github.com/vfilipovsky/url-shortener/internal/service"
	"github.com/vfilipovsky/url-shortener/pkg/config"
	"github.com/vfilipovsky/url-shortener/pkg/logger"
	"github.com/vfilipovsky/url-shortener/pkg/randomizer"
	"github.com/vfilipovsky/url-shortener/pkg/tokenizer"
)

func register(r *mux.Router, db *sql.DB, cfg *config.Config) *mux.Router {
	// utils
	randomizer := randomizer.New()
	tokenizer := tokenizer.New()

	// repositories
	urlRepo := repository.NewUrl(db)
	accessRepo := repository.NewAccess(db)

	// services
	urlService := service.NewUrl(urlRepo)
	accessService := service.NewAccess(accessRepo)

	// actions
	createUrlAction := url.NewCreateAction(accessService, urlService, randomizer, tokenizer)
	deleteUrlAction := url.NewDeleteAction(accessService, urlService)
	getUrlAction := url.NewGetAction(urlService)
	infoAction := url.NewInfoAction(urlService)
	getUrlsByAccessAction := url.NewGetAllUrlsByAccessAction(accessService, urlService)

	// routes
	RegisterDefaultRoutes(r)
	handler.RegisterUrlRoutes(
		r,
		cfg.Http,
		createUrlAction,
		deleteUrlAction,
		getUrlAction,
		infoAction,
		getUrlsByAccessAction,
	)

	logRoutes(r)
	logger.Infof("Routes registered")

	return r
}

func logRoutes(r *mux.Router) {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, err := route.GetMethods()

		if err != nil {
			return nil
		}

		path, err := route.GetPathTemplate()

		if err != nil {
			return nil
		}

		logger.Infof("[%s] %s", strings.Join(methods, "|"), path)
		return nil
	})
}
