package handlers

import (
	"net/http"
)

func LogOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: "",
	})
	w.WriteHeader(http.StatusOK)
}
