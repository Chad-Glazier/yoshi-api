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

type Profanity struct {
	DisplayName bool `json:"displayName"`
	FirstName   bool `json:"firstName"`
	LastName    bool `json:"lastName"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	mw.
		NewPipeline(register).
		Use(
			mw.Cors,
			mw.Method(http.MethodPost),
			mw.DB,
		).
		Run(w, r)
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
			Email:       true,
			DisplayName: true,
		})
		return
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
	case user.ErrProfaneDisplayName:
		w.WriteHeader(http.StatusBadRequest)
		util.SendJSON(w, Profanity{
			DisplayName: true,
		})
		return
	case user.ErrProfaneFirstName:
		w.WriteHeader(http.StatusBadRequest)
		util.SendJSON(w, Profanity{
			FirstName: true,
		})
		return
	case user.ErrProfaneLastName:
		w.WriteHeader(http.StatusBadRequest)
		util.SendJSON(w, Profanity{
			LastName: true,
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
