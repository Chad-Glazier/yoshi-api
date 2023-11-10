package user

import (
	"golang.org/x/crypto/bcrypt"
	"yoshi/db"
)


// Registers a user and creates a new session, returning the session ID.
//
// Potential errors include:
// - `ErrDatabaseError`
// - `ErrEmailTaken`
// - `ErrDisplayNameTaken`
// - `ErrPwdTooLong`
func Register(u *UserRegistration) (string, error) {
	db, err := db.Connect()
	if err != nil {
		return "", ErrDatabaseError
	}
	defer db.Close()
	rows, err := db.Query(`
		SELECT email 
		FROM user_credentials 
		WHERE email = ?`,
		u.Email,
	)
	if err != nil {
		return "", ErrDatabaseError
	}
	emailTaken := rows.Next()
	rows.Close()
	rows, err = db.Query(`
		SELECT display_name 
		FROM user_data 
		WHERE display_name = ?`,
		u.DisplayName,
	)
	if err != nil {
		return "", ErrDatabaseError
	}
	displayNameTaken := rows.Next()
	rows.Close()
	switch {
	case emailTaken && displayNameTaken:
		return "", ErrEmailAndDisplayNameTaken
	case emailTaken:
		return "", ErrEmailTaken
	case displayNameTaken:
		return "", ErrDisplayNameTaken
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(u.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", ErrPwdTooLong
	}
	_, err = db.Exec(`
		INSERT INTO user_data
			(email, first_name, last_name, display_name)
		VALUES
			(?, ?, ?, ?)`,
		u.Email,
		u.FirstName,
		u.LastName,
		u.DisplayName,
	)
	if err != nil {
		return "", ErrDatabaseError
	}
	_, err = db.Exec(`
		INSERT INTO user_credentials
			(email, password)
		VALUES
			(?, ?)`,
		u.Email,
		encryptedPassword,
	)
	if err != nil {
		db.Exec(`
			DELETE FROM user_data 
			WHERE email = ?
			`,
			u.Email,
		)
		return "", ErrDatabaseError
	}
	_, err = db.Exec(`
		INSERT INTO user_preferences
			(email)
		VALUES
			(?)	
		`,
		u.Email,
	)
	if err != nil {
		db.Exec(`
			DELETE FROM user_data 
			WHERE email = ?
			`,
			u.Email,
		)
		db.Exec(`
			DELETE FROM user_credentials
			WHERE email = ?
			`,
			u.Email,
		)
		return "", ErrDatabaseError
	}
	return createNewSession(u.Email, db)
}
