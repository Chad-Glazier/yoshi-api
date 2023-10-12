package handler

import (
	"net/http"
	"yoshi-api/base_handlers"
	"yoshi-api/util"
)

func PingDB(w http.ResponseWriter, r *http.Request) {
	handler := util.HandlerConfig{
		HttpMethod:  http.MethodGet,
		BaseHandler: base_handlers.PingDB,
	}

	handler.Execute(w, r)
}
