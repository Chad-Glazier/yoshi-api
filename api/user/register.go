package handlers

import (
	"net/http"
	"yoshi/db/user"
	"yoshi/util"
)

type RegistrationConflict struct {
	Email       bool `json:"email"`
	DisplayName bool `json:"displayName"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	stop := util.AllowCors(w, r)
	if stop {
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	registrationDetails, err := util.ParseBody[user.UserRegistration](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionId, err := user.Register(registrationDetails)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	util.SetSessionCookie(sessionId, w)
}
