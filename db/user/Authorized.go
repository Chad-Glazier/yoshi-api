package user

import (
	"errors"
	"net/http"
	"time"
	"yoshi/db"
)

// returns the email of the user if the request contained a valid session
// id cookie. Additionally, this will renew the session cookie if it is valid.
//
// Potential errors include:
// - `ErrNoAuthCookie`
// - `ErrDatabaseError`
// - `ErrUnrecognizedSession`
// - `ErrExpiredSession`
func Authorized(r *http.Request) (string, error) {
	cookies := r.Cookies()
	sessionId := ""
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionId = cookie.Value
			break
		}
	}
	if sessionId == "" {
		return "", ErrNoAuthCookie
	}
	db, err := db.Connect()
	if err != nil {
		return "", ErrDatabaseError
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT email, last_renewed, uuid
		FROM user_sessions
		WHERE uuid = ?
		`,
		sessionId,
	)
	if err != nil {
		return "", ErrDatabaseError
	}
	defer rows.Close()
	if !rows.Next() {
		return "", ErrUnrecognizedSession
	}
	session := &Session{}
	err = rows.Scan(
		&session.Email,
		&session.LastRenewed,
		&session.Uuid,
	)
	if err != nil {
		return "", err
	}
	timestamp, err := time.Parse("2006-01-02 15:04:05", session.LastRenewed)
	if err != nil {
		return "", errors.New(session.LastRenewed)
	}
	if time.Since(timestamp) > maxSessionAge {
		db.Exec(`
			DELETE FROM user_sessions
			WHERE uuid = ?
			`,
			session.Uuid,
		)
		return "", ErrExpiredSession
	}
	go renewSession(session.Uuid, db)
	return session.Email, nil
}
