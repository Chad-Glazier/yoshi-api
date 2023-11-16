package handlers

import (
	"net/http"
	"yoshi/db/user"
	"yoshi/mw"
	"yoshi/util"
)

func Info(w http.ResponseWriter, r *http.Request) {
	mw.
		NewPipeline(info).
		Use(
			mw.Cors,
			mw.Method(http.MethodGet),
			mw.DB,
		).
		Run(w, r)
}

func info(res *mw.Resources, w http.ResponseWriter, r *http.Request) {
	displayNames, ok := r.URL.Query()["displayName"]
	if !ok {
		sendPersonalInfo(res, w, r)
		return
	}
	publicInfo, err := user.GetPublicInfo(res.DB, displayNames...)
	switch err {
	case nil:
		if len(publicInfo) == 1 {
			util.SendJSON(w, publicInfo[0])
		} else {
			util.SendJSON(w, publicInfo)
		}
	case user.ErrDisplayNameNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sendPersonalInfo(res *mw.Resources, w http.ResponseWriter, r *http.Request) {
	session, err := user.ExistingSession(res.DB, r)
	switch err {
	case nil:
		break
	case user.ErrNoAuthCookie, user.ErrUnrecognizedSession, user.ErrExpiredSession:
		w.WriteHeader(http.StatusUnauthorized)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	personal, err := user.GetPersonalInfo(res.DB, session.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SendJSON(w, personal)
}
