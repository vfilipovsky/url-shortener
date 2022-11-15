package api

import (
	"encoding/json"
	"net/http"

	"github.com/vfilipovsky/url-shortener/pkg/logger"
)

func Respond(w http.ResponseWriter, result any, statusCode ...int) {
	if _, ok := result.(error); !ok {
		code := http.StatusOK

		if len(statusCode) > 0 {
			code = statusCode[0]
		}

		send(w, result, code)
		return
	}

	message, code := handleError(result.(error))

	send(w, message, code)
}

func send(w http.ResponseWriter, result any, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(&result); err != nil {
		logger.Errorf("failed to send response: %s", err.Error())
	}
}
