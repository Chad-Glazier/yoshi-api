package handler

import (
	"net/http"
	"yoshi-api/handlers"
	"yoshi-api/util"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	handler := util.HandlerConfig{
		HttpMethod:  "GET",
		BaseHandler: handlers.Ping,
	}

	handler.Execute(w, r)
}
