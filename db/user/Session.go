package user

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// represents a session as it appears in the database.
type Session struct {
	Email       string
	LastRenewed string
	Uuid        string
}

// this is the max amount of time that a session is allowed before it is
// considered to have expired. This constant represents 4 days in nanoseconds.
const MaxSessionAge = 3.456e+14

// Retrieves a `Session` record from the database that matches the `session_id`
// cookie in a request.
//
// Possible error values:
//   - `ErrNoAuthCookie`
//   - `ErrUnrecognizedSession`
//   - `ErrExpiredSession`
//   - `ErrServer`
//   - `ErrDatabase`
func ExistingSession(db *sql.DB, r *http.Request) (*Session, error) {
	sessionId := ""
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionId = cookie.Value
		}
	}
	if sessionId == "" {
		return nil, ErrNoAuthCookie
	}

	rows, err := db.Query(`
		SELECT email, uuid, last_renewed
		FROM user_sessions
		WHERE uuid = ?
		`,
		sessionId,
	)
	if err != nil {
		return nil, ErrDatabase
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, ErrUnrecognizedSession
	}
	session := &Session{}
	err = rows.Scan(
		&session.Email,
		&session.Uuid,
		&session.LastRenewed,
	)
	if err != nil {
		return nil, ErrServer
	}
	return session, nil
}

// Returns `true` if the session is expired, otherwise returning `false`.
// Expiration is determined by `MaxSessionAge`.
func (s *Session) IsExpired() bool {
	timestamp, _ := time.Parse("2006-01-02 15:04:05", s.LastRenewed)
	// ignoring the possible error since it should not happen,
	// unless the database is migrated to another that does not
	// have dates in this format.
	return time.Since(timestamp) > MaxSessionAge
}

// Removes the record of a session in the database.
//
// Possible error values:
// - `ErrDatabase`
func (s *Session) Terminate(db *sql.DB) error {
	_, err := db.Exec(`
		DELETE FROM user_sessions
		WHERE uuid = ?
		`,
		s.Uuid,
	)
	if err != nil {
		return ErrDatabase
	}
	return nil
}

// Renews a session in the database.
//
// Possible error values:
// - `ErrDatabase`
func (s *Session) Renew(db *sql.DB) error {
	_, err := db.Exec(`
		UPDATE user_sessions 
		SET last_renewed = CURRENT_TIMESTAMP() 
		WHERE uuid = ?
		`,
		s.Uuid,
	)
	return err
}

// Sets a cookie in the response that represents the session.
func (s *Session) SetCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    s.Uuid,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
}

// Unsets the `session_id` cookie.
func UnsetSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().AddDate(0, -1, 0),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, cookie)
}

// Creates a new session in the database for the specified user.
//
// Possible error values:
//   - `ErrDatabase`
//   - `ErrServer`
func CreateSession(db *sql.DB, email string) (*Session, error) {
	newSessionId, err := uuid.NewRandom()
	if err != nil {
		return nil, ErrServer
	}
	_, err = db.Exec(`
		insert into user_sessions (email, uuid)
		values (?, ?)
		`,
		email,
		newSessionId.String(),
	)
	if err != nil {
		return nil, ErrDatabase
	}
	rows, err := db.Query(`
		SELECT email, uuid, last_renewed
		FROM user_sessions
		WHERE uuid = ?
	`, newSessionId.String())
	if err != nil {
		return nil, ErrDatabase
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, ErrServer
	}
	newSession := Session{}
	rows.Scan(
		&newSession.Email,
		&newSession.Uuid,
		&newSession.LastRenewed,
	)
	return &newSession, nil
}
