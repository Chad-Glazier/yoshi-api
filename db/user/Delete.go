package user

import (
	"database/sql"
)

// deletes a user from the database. This function is idempotent,
// and it should never err.
func Delete(db *sql.DB, email string) {
	db.Exec(`
		DELETE FROM user_data
		WHERE email = ?
		`,
		email,
	)

	db.Exec(`
		DELETE FROM user_preferences
		WHERE email = ?
		`,
		email,
	)

	db.Exec(`
		DELETE FROM user_credentials
		WHERE email = ?
		`,
		email,
	)

	db.Exec(`
		DELETE FROM user_sessions
		WHERE email = ?
		`,
		email,
	)
}
