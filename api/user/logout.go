package handlers

import (
	"net/http"
	"yoshi/db"
	"yoshi/db/user"
	"yoshi/util"
)

func LogOut(w http.ResponseWriter, r *http.Request) {
	stop := util.AllowCors(w, r)
	if stop {
		return
	}

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	db, err := db.Connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()

	session, _ := user.ExistingSession(db, r)
	if session != nil {
		session.Terminate(db)
	}

	user.UnsetSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}
