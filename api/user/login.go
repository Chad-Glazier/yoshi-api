package handlers

import (
	"net/http"
	"yoshi/db/user"
	"yoshi/util"
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	stop := util.AllowCors(w, r)
	if stop {
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	credentials, err := util.ParseBody[user.UserCredentials](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionId, err := user.LogIn(r, credentials)
	switch err {
	case nil:
		break
	case user.ErrEmailNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	case user.ErrIncorrectPassword:
		w.WriteHeader(http.StatusUnauthorized)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	util.SetSessionCookie(sessionId, w)
}