package user

import (
	"database/sql"
)

// The incoming email and password of a user; the `Password` should be unhashed
type UserCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// This will take the credentials for a login attempt and return the new
// session object or an error. This is the same as using `CheckCredentials`,
// except that function doesn't create a session.
//
// Potential errors include:
//	- `ErrEmailNotFound`
//	- `ErrIncorrectPassword`
//	- `ErrDatabase`
//	- `ErrServer`
func LogIn(db *sql.DB, c *UserCredentials) (*Session, error) {
	err := CheckCredentials(db, c.Email, c.Password)
	if err != nil {
		return nil, err
	}

	return CreateSession(db, c.Email)
}
