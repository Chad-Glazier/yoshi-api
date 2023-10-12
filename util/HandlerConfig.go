package util

import (
	"net/http"
	"yoshi-api/types"
)

/*
	A struct that defines certain settings on an endpoint, including
	what HTTP method it accepts, the "base handler" (the function that
	should process the request once it's passed through the middleware),
	and any number of middleware functions that will be executed in
	sequence.
*/
type HandlerConfig struct {
	/* 
		The handler function that is ultimately responsible for the
		response. This will be executed if and only if all of the
		middleware execute without issue.
	 */
	BaseHandler http.HandlerFunc
	/*
		The middleware functions to be applied. Middleware follow the
		signature, `func (r *http.Request) (err *ErrorInfo)`; if the
		middleware executed normally, then you should return `nil`. If
		there was a problem and you want to end the request early, 
		without executing the rest of the middleware and the `BaseHandler`,
		then return a non-`nil` `ErrorInfo` object.
	*/
	Middleware  []types.MiddlewareFunc
	HttpMethod  string
}

/*
	Runs the handler function with the specified configuration.
*/
func (config *HandlerConfig) Execute(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if config.HttpMethod == "" {
		config.HttpMethod = http.MethodGet
	}

	if r.Method != config.HttpMethod {
		(&types.ApiError{
			Status:  405,
			Message: "This endpoint only accepts " + config.HttpMethod + " requests",
		}).Send(w)
		return
	}

	for _, middleware := range config.Middleware {
		err := middleware(r)
		if err != nil {
			err.Send(w)
			return
		}
	}

	config.BaseHandler(w, r)
}
