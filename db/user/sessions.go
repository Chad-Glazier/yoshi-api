package user

import 	(
	"database/sql"
	"github.com/google/uuid"
)

// creates a new session associated with the provided email. Returns
// the UUID for the session.
func createNewSession(email string, db *sql.DB) (string, error) {
	sessionId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	_, err = db.Exec(`
		insert into user_sessions (email, uuid)
		values (?, ?)
	`, email, sessionId.String())
	if err != nil {
		return "", err
	}
	return sessionId.String(), nil
}

// renews a session (replaces the timestamp with the current time).
func renewSession(sessionId string, db *sql.DB) error {
	_, err := db.Exec(`
		UPDATE user_sessions 
		SET last_renewed = CURRENT_TIMESTAMP() 
		WHERE uuid = ?
		`,
		sessionId,
	)
	return err
}

// deletes a session from the database.
func endSession(sessionId string, db *sql.DB) error {
	_, err := db.Exec(`
		DELETE FROM user_sessions 
		WHERE uuid = ?
		`,
		sessionId,
	)
	return err
}