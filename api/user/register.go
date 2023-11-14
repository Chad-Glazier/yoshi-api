package handlers

import (
	"net/http"
	"yoshi/db"
	"yoshi/db/user"
	"yoshi/util"
)

type RegistrationConflict struct {
	Email       bool `json:"email"`
	DisplayName bool `json:"displayName"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	registrationDetails, err := util.ParseBody[user.UserRegistration](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	db, err := db.Connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()

	session, err := user.Register(db, registrationDetails)
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

	session.SetCookie(w)
	w.WriteHeader(http.StatusCreated)
}
