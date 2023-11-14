package handlers

import (
	"net/http"
	"yoshi/db"
	"yoshi/db/user"
	"yoshi/util"
)

type DeletionConfirmation struct {
	Password string `json:"password"`
}

func Unregister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
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
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	c, err := util.ParseBody[DeletionConfirmation](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = user.CheckCredentials(db, email, c.Password)
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

	user.Delete(db, email)
	user.UnsetSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}
