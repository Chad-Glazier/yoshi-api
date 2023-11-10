package handlers

import (
	"net/http"
	"yoshi/db"
	"yoshi/db/user"
	"yoshi/util"
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	stop := util.AllowCors(w, r)
	if stop {
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	credentials, err := util.ParseBody[user.UserCredentials](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := db.Connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()

	session, err := user.LogIn(db, credentials)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.SetCookie(w)
}
