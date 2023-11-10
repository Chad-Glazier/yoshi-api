package handlers

import (
	"net/http"
	"yoshi/db/user"
	"yoshi/util"
)

func LogOut(w http.ResponseWriter, r *http.Request) {
	stop := util.AllowCors(w, r)
	if stop {
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user.UnsetSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}
