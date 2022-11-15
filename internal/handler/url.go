package handler

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/vfilipovsky/url-shortener/internal/action/url"
	"github.com/vfilipovsky/url-shortener/internal/payload"
	"github.com/vfilipovsky/url-shortener/pkg/api"
	"github.com/vfilipovsky/url-shortener/pkg/config"
	"github.com/vfilipovsky/url-shortener/pkg/exception"
	"github.com/vfilipovsky/url-shortener/pkg/logger"
)

type urlHandler struct {
	cfg             *config.Http
	infoAction      url.Info
	createUrlAction url.Create
	deleteUrlAction url.Delete
	getUrlAction    url.Get
	getUrlsByAccess url.GetAllByAccess
}

func newUrl(
	cfg *config.Http,
	createUrlAction url.Create,
	deleteUrlAction url.Delete,
	getUrlAction url.Get,
	infoAction url.Info,
	getUrlsByAccessAction url.GetAllByAccess,
) *urlHandler {
	return &urlHandler{
		cfg:             cfg,
		createUrlAction: createUrlAction,
		deleteUrlAction: deleteUrlAction,
		getUrlAction:    getUrlAction,
		infoAction:      infoAction,
		getUrlsByAccess: getUrlsByAccessAction,
	}
}

func RegisterUrlRoutes(
	r *mux.Router,
	cfg *config.Http,
	createUrlAction url.Create,
	deleteUrlAction url.Delete,
	getUrlAction url.Get,
	infoAction url.Info,
	getUrlsByAccess url.GetAllByAccess,
) {
	h := newUrl(
		cfg,
		createUrlAction,
		deleteUrlAction,
		getUrlAction,
		infoAction,
		getUrlsByAccess,
	)

	s := http.StripPrefix("/public/", http.FileServer(http.Dir("./public/")))
	r.PathPrefix("/public/").Handler(s)

	r.HandleFunc(V1Url, h.Create).Methods(http.MethodPost)
	r.HandleFunc(V1Url+"/{code}", h.Delete).Methods(http.MethodDelete)
	r.HandleFunc(V1Url+"/{code}", h.GetUrl).Methods(http.MethodPost)
	r.HandleFunc(V1Url+"/access/list", h.GetUrlsByAccess).Methods(http.MethodPost)
	r.HandleFunc("/{code}", h.GetPage).Methods(http.MethodGet)
}

func somethingWentWrong(w http.ResponseWriter, r *http.Request, err error) {
	logger.Errorf(err.Error())
	http.ServeFile(w, r, GetPublicDir()+"/something_went_wrong.html")
}

func urlNotFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, GetPublicDir()+"/url_not_found.html")
}

func (h *urlHandler) GetUrlsByAccess(w http.ResponseWriter, r *http.Request) {
	pl := &payload.GetUrlsByAccess{}

	if err := api.DecodeAndValidate(r.Body, pl); err != nil {
		api.Respond(w, err)
		return
	}

	if err := api.ParseQueryParams(r, pl); err != nil {
		api.Respond(w, err)
		return
	}

	urls, err := h.getUrlsByAccess.Run(pl.AccessToken, pl.Offset, pl.Limit)

	if err != nil {
		api.Respond(w, err)
		return
	}

	api.Respond(w, urls)
}

func (h *urlHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	url, err := h.infoAction.Run(mux.Vars(r)["code"])

	if err != nil {
		if _, ok := err.(*exception.StatusError); ok && err.(*exception.StatusError).Code == 404 {
			urlNotFound(w, r)
			return
		}

		somethingWentWrong(w, r, err)
		return
	}

	if !url.IsSecured {
		http.Redirect(w, r, url.Url, http.StatusPermanentRedirect)
		return
	}

	tmpl, err := template.ParseFiles(GetPublicDir() + "/url.html")

	if err != nil {
		somethingWentWrong(w, r, err)
		return
	}

	if err := tmpl.Execute(w, struct {
		Code string
		Url  string
	}{
		Code: url.Code,
		Url:  h.cfg.Host + ":" + h.cfg.Port,
	}); err != nil {
		somethingWentWrong(w, r, err)
	}
}

func (h *urlHandler) GetUrl(w http.ResponseWriter, r *http.Request) {
	pl := &payload.GetUrl{Code: mux.Vars(r)["code"]}

	if err := api.DecodeAndValidate(r.Body, pl); err != nil {
		api.Respond(w, err)
		return
	}

	url, err := h.getUrlAction.Run(pl.Pin, pl.Code)

	if err != nil {
		api.Respond(w, err)
		return
	}

	api.Respond(w, url)
}

func (h *urlHandler) Create(w http.ResponseWriter, r *http.Request) {
	pl := &payload.CreateUrl{}

	if err := api.DecodeAndValidate(r.Body, pl); err != nil {
		api.Respond(w, err)
		return
	}

	url, err := h.createUrlAction.Run(pl.AccessToken, pl.Url, pl.IsSecured, pl.AliveUntil)

	if err != nil {
		api.Respond(w, err)
		return
	}

	api.Respond(w, url)
}

func (h *urlHandler) Delete(w http.ResponseWriter, r *http.Request) {
	pl := &payload.DeleteUrl{Code: mux.Vars(r)["code"]}

	if err := api.DecodeAndValidate(r.Body, pl); err != nil {
		api.Respond(w, err)
		return
	}

	api.Respond(w, h.deleteUrlAction.Run(pl.Code, pl.AccessToken))
}
