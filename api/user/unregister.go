package handlers

import (
	"net/http"
	"yoshi/db"
	"yoshi/db/user"
	"yoshi/util"
)

func Unregister(w http.ResponseWriter, r *http.Request) {
	stop := util.AllowCors(w, r)
	if stop {
		return
	}

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	db, err := db.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	email, err := user.Authorized(db, r)
	switch err {
	case nil:
		break
	case user.ErrExpiredSession, user.ErrUnrecognizedSession, user.ErrNoAuthCookie:
		w.WriteHeader(http.StatusUnauthorized)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	user.Delete(db, email)
	user.UnsetSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}
