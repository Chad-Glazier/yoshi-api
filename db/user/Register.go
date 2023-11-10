package user

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

// Registers a user and creates a new session, returning the session cookie.
//
// Potential errors include:
//	- `ErrServer`
//	- `ErrDatabase`
//	- `ErrEmailTaken`
//	- `ErrDisplayNameTaken`
//	- `ErrEmailAndDisplayNameTaken`
//	- `ErrPwdTooLong`
func Register(db *sql.DB, u *UserRegistration) (*Session, error) {
	rows, err := db.Query(`
		SELECT email 
		FROM user_credentials 
		WHERE email = ?`,
		u.Email,
	)
	if err != nil {
		return nil, ErrDatabase
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
		return nil, ErrDatabase
	}
	displayNameTaken := rows.Next()
	rows.Close()

	switch {
	case emailTaken && displayNameTaken:
		return nil, ErrEmailAndDisplayNameTaken
	case emailTaken:
		return nil, ErrEmailTaken
	case displayNameTaken:
		return nil, ErrDisplayNameTaken
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(u.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, ErrPwdTooLong
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
		return nil, ErrDatabase
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
		// if that query failed, then we should undo the last one as well.
		db.Exec(`
			DELETE FROM user_data 
			WHERE email = ?
			`,
			u.Email,
		)
		return nil, ErrDatabase
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
		// if that query failed, then we should undo the previous two as well.
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
		return nil, ErrDatabase
	}

	return CreateSession(db, u.Email)
}
