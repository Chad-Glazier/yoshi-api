package handlers

import (
	"net/http"
	"yoshi/db/user"
	"yoshi/util"
	"yoshi/mw"
)

func Login(w http.ResponseWriter, r *http.Request) {
	mw.
		NewPipeline(login).
		Use(
			mw.Cors,
			mw.Method(http.MethodPost),
			mw.DB,			
		).
		Run(w, r)
}

func login(res *mw.Resources, w http.ResponseWriter, r *http.Request) {
	credentials, err := util.ParseBody[user.UserCredentials](r)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	session, err := user.LogIn(res.DB, credentials)
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
	w.WriteHeader(http.StatusOK)
}
