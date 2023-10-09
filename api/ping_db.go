package handler

import (
	"net/http"
	"yoshi-api/handlers"
	"yoshi-api/util"
)

func PingDB(w http.ResponseWriter, r *http.Request) {
	handler := util.HandlerConfig{
		HttpMethod:  "GET",
		BaseHandler: handlers.PingDB,
	}

	handler.Execute(w, r)
}
