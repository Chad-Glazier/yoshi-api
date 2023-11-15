package handlers

import (
	"net/http"
	"yoshi/db/user"
	"yoshi/mw"
	"yoshi/util"
)

type RegistrationConflict struct {
	Email       bool `json:"email"`
	DisplayName bool `json:"displayName"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	mw.NewPipeline(
		register,
		mw.Cors,
		mw.Method(http.MethodPost),
		mw.DB,
	).Run(w, r)
}

func register(res *mw.Resources, w http.ResponseWriter, r *http.Request) {
	u, err := util.ParseBody[user.UserRegistration](r)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	session, err := user.Register(res.DB, u)
	switch err {
	case nil:
		break
	case user.ErrEmailAndDisplayNameTaken:
		w.WriteHeader(http.StatusConflict)
		util.SendJSON(w, RegistrationConflict{
			Email: true,
			DisplayName: true,
		})
	case user.ErrEmailTaken:
		w.WriteHeader(http.StatusConflict)
		util.SendJSON(w, RegistrationConflict{
			Email: true,
		})
		return
	case user.ErrDisplayNameTaken:
		w.WriteHeader(http.StatusConflict)
		util.SendJSON(w, RegistrationConflict{
			DisplayName: true,
		})		
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	existingSession, _ := user.ExistingSession(res.DB, r)
	if existingSession != nil {
		existingSession.Terminate(res.DB)
	}

	session.SetCookie(w)
	w.WriteHeader(http.StatusCreated)
}
