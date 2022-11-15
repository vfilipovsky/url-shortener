package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/vfilipovsky/url-shortener/pkg/api"
)

const V1 = "/api/v1"

func RegisterDefaultRoutes(r *mux.Router) {
	r.HandleFunc(V1+"/ping", ping).Methods(http.MethodGet)
}

func ping(w http.ResponseWriter, _ *http.Request) {
	api.Respond(w, "pong", http.StatusOK)
}
