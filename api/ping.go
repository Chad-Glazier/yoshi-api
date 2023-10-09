package handler

import (
	"net/http"
	"yoshi-api/handlers"
	"yoshi-api/middleware"
	"yoshi-api/util"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	handler := util.HandlerConfig{
		HttpMethod:  "GET",
		BaseHandler: handlers.Ping,
	}

	handler.Use(middleware.Authenticated)

	handler.Execute(w, r)
}
