package user

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

// Checks the credentials of a user. Will return `nil` if the 
// credentials are valid, otherwise will return some error.
//
// Potential errors include:
//	- `ErrDatabase`
//	- `ErrEmailNotFound`
//	- `ErrIncorrectPassword`
//	- `ErrServer`
func CheckCredentials(db *sql.DB, email, password string) error {
	rows, err := db.Query(`
		SELECT password
		FROM user_credentials
		WHERE email = ?
		`,
		email,
	)
	if err != nil {
		return ErrDatabase
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrEmailNotFound
	}

	var storedPasswrd string
	err = rows.Scan(
		&storedPasswrd,
	)
	if err != nil {
		return ErrServer
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(storedPasswrd),
		[]byte(password),
	)
	if err != nil {
		return ErrIncorrectPassword
	}

	return nil
}
