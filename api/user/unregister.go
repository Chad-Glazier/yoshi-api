package handlers

import (
	"net/http"
	"yoshi/db/user"
	"yoshi/mw"
	"yoshi/util"
)

func Unregister(w http.ResponseWriter, r *http.Request) {
	mw.
		NewPipeline(unregister).
		Use(	
			mw.Cors,
			mw.Method(http.MethodDelete),
			mw.DB,
			mw.Session,
		).
		Run(w, r)
}

type DeletionConfirmation struct {
	Password string `json:"password"`
}

func unregister(res *mw.Resources, w http.ResponseWriter, r *http.Request) {
	c, err := util.ParseBody[DeletionConfirmation](r)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	err = user.CheckCredentials(res.DB, res.Session.Email, c.Password)
	switch err {
	case nil:
		break
	case user.ErrEmailNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case user.ErrIncorrectPassword:
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Delete(res.DB, res.Session.Email)
	user.UnsetSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}
