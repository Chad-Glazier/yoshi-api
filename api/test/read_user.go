package handler

import (
	"net/http"
	"yoshi-api/base_handlers/test"
	"yoshi-api/util"
)

func TestReadUser(w http.ResponseWriter, r *http.Request) {
	handler := util.HandlerConfig{
		HttpMethod: http.MethodGet,
		BaseHandler: test.ReadUser,
	}

	handler.Execute(w, r)
}