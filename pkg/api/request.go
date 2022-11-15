package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/schema"

	"github.com/vfilipovsky/url-shortener/pkg/exception"
)

var (
	validate = validator.New()
	decoder  = schema.NewDecoder()
)

func ParseQueryParams(r *http.Request, payload any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := decoder.Decode(payload, r.Form); err != nil {
		return err
	}

	return nil
}

func DecodeAndValidate(requestBody io.Reader, payload any) error {
	body, err := io.ReadAll(requestBody)

	if err != nil {
		return err
	}

	if len(body) == 0 {
		return exception.EmptyBody
	}

	if err := json.Unmarshal(body, payload); err != nil {
		return err
	}

	if err := validate.Struct(payload); err != nil {
		return err
	}

	return nil
}

func ValidateRequestPayload(payload any) error {
	return validate.Struct(payload)
}
