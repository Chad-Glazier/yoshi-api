package handlers

import (
	"net/http"
	"yoshi/db/user"
	"yoshi/mw"
	"yoshi/util"
)

func Unregister(w http.ResponseWriter, r *http.Request) {
	mw.NewPipeline(
		unregister,
		mw.Cors,
		mw.Method(http.MethodDelete),
		mw.DB,
		mw.Session,
	).Run(w, r)
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
	case user.ErrDatabase, user.ErrServer:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	case user.ErrEmailNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	case user.ErrIncorrectPassword:
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user.Delete(res.DB, res.Session.Email)
	user.UnsetSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}
