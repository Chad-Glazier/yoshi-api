package handler

import (
	"net/http"
	"yoshi-api/base_handlers"
	"yoshi-api/util"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	handler := util.HandlerConfig{
		HttpMethod:  http.MethodGet,
		BaseHandler: base_handlers.Ping,
	}

	handler.Execute(w, r)
}
