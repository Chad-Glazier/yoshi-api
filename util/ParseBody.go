package util

import (
	"net/http"
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Parses a JSON body from a request, and then validates the struct
// with the `validator` package from https://github.com/go-playground/validator/v10
func ParseBody[T any](r *http.Request) (*T, error) {
	validate = validator.New(validator.WithRequiredStructEnabled())
	decoder := json.NewDecoder(r.Body)
	var body T
	err := decoder.Decode(&body)
	if err != nil {
		return nil, err
	}
	err = validate.Struct(body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}