package user

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

// The incoming email and password of a user; the `Password` should be unhashed
type UserCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// This will take the credentials for a login attempt and return the new
// session object or an error.
//
// Potential errors include:
//	- `ErrEmailNotFound`
//	- `ErrIncorrectPassword`
//	- `ErrDatabase`
//	- `ErrServer`
func LogIn(db *sql.DB, c *UserCredentials) (*Session, error) {
	rows, err := db.Query(`
		SELECT email, password FROM user_credentials
		WHERE email = ?
	`, c.Email)
	if err != nil {
		return nil, ErrDatabase
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrEmailNotFound
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
		return nil, ErrIncorrectPassword
	}

	return CreateSession(db, c.Email)
}
