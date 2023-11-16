package handlers

import (
	"net/http"
	"yoshi/db/user"
	"yoshi/mw"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	mw.
		NewPipeline(logout).
		Use(
			mw.Cors,
			mw.Method(http.MethodDelete),
			mw.DB,
		).
		Run(w, r)
}

func logout(res *mw.Resources, w http.ResponseWriter, r *http.Request) {
	session, _ := user.ExistingSession(res.DB, r)
	if session != nil {
		session.Terminate(res.DB)
	}

	user.UnsetSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}
