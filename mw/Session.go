package mw

import (
	"net/http"
	"yoshi/db/user"
)

// Expects access to:
//	- `DB`
//
// Grants access to:
//	- `Session`
//
// May send back the following errors:
//	- `http.StatusUnauthorized` (no body)
//	- `http.StatusInternalServerError` (text body)
func Session(res *Resources, w http.ResponseWriter, r *http.Request) (bool, CleanupFunc) {
	session, err := user.ExistingSession(res.DB, r)
	
	switch err {
	case nil: 
		res.Session = session
		return true, nil
	case user.ErrNoAuthCookie, user.ErrUnrecognizedSession, user.ErrExpiredSession:
		w.WriteHeader(http.StatusUnauthorized)
		return false, nil
	default: 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return false, nil
	}
}
