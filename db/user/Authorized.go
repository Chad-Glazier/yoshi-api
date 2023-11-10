package user

import (
	"database/sql"
	"net/http"
)

// returns the email of the user if the request contained a valid session
// id cookie. Additionally, this will renew the session cookie if it is valid.
//
// Potential errors include:
// - `ErrNoAuthCookie`
// - `ErrDatabase`
// - `ErrUnrecognizedSession`
// - `ErrExpiredSession`
// - `ErrServer`
func Authorized(db *sql.DB, r *http.Request) (string, error) {
	session, err := ExistingSession(db, r)
	if err != nil {
		return "", err
	}

	if session.IsExpired() {
		err = session.Terminate(db)
		if err != nil {
			return "", err
		}
	}

	err = session.Renew(db)
	if err != nil {
		return "", err
	}
	
	return session.Email, nil
}
