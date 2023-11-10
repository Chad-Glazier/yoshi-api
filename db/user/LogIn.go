package user

import (
	"net/http"
	"yoshi/db"

	"golang.org/x/crypto/bcrypt"
)

// This is meant to take the credentials for a login attempt. Upon success,
// it will return the UUID of the new session and `nil`.
//
// Potential errors include:
// - `ErrEmailNotFound`
// - `ErrIncorrectPassword`
// - `ErrDatabaseError`
func LogIn(r *http.Request, c *UserCredentials) (string, error) {
	db, err := db.Connect()

	if err != nil {
		return "", ErrDatabaseError
	}

	defer db.Close()

	rows, err := db.Query(`
		SELECT email, password FROM user_credentials
		WHERE email = ?
	`, c.Email)

	if err != nil {
		return "", ErrDatabaseError
	}

	defer rows.Close()

	if !rows.Next() {
		return "", ErrEmailNotFound
	}

	storedCredentials := UserCredentials{}
	rows.Scan(
		&storedCredentials.Email,
		&storedCredentials.Password,
	)
	err = bcrypt.CompareHashAndPassword(
		[]byte(storedCredentials.Password),
		[]byte(c.Password),
	)
	if err != nil {
		return "", ErrIncorrectPassword
	}

	cookies := r.Cookies()
	sessionId := ""
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionId = cookie.Value
			break
		}
	}
	if sessionId != "" {
		go endSession(sessionId, db)
	}

	return createNewSession(c.Email, db)
}
