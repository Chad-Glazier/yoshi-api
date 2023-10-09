package util

import (
	"net/http"
	"yoshi-api/types"
)

type HandlerConfig struct {
	BaseHandler http.HandlerFunc
	Middleware  []types.MiddlewareFunc
	HttpMethod  string
}

func (config *HandlerConfig) Execute(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if config.HttpMethod == "" {
		config.HttpMethod = "GET"
	}

	if r.Method != config.HttpMethod {
		SendError(w, types.ErrorInfo{
			Status:  405,
			Message: "This endpoint only accepts " + config.HttpMethod + " requests",
		})
		return
	}

	for _, middleware := range config.Middleware {
		err := middleware(r)
		if err != nil {
			SendError(w, *err)
			return
		}
	}

	config.BaseHandler(w, r)
}

func (config *HandlerConfig) Use(newMiddleware ...types.MiddlewareFunc) {
	config.Middleware = append(config.Middleware, newMiddleware...)
}
