package db

import (
	"database/sql"
	"net/http"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id          int
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	password    string
}

func (u *User) Id() int {
	return u.id
}

func (u *User) CheckPassword(passwordAttempt string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(u.password),
		[]byte(passwordAttempt),
	)
	return err == nil
}

// Reads a User from the database, from a given `id`.
//
//	- `nil, error`: There was some error connecting to, 
//    or executing a query on, the database.
//	- `nil, nil`: There was no error, and the user was 
//    not found.
//	- `*User, nil`: The user was successfully found.
//
func ReadUser(db *sql.DB, userId int) (*User, error) {
	rows, err := db.Query(`
		SELECT
			id,
			email,
			display_name,
			first_name,
			last_name,
			password
		FROM
			users
		WHERE
			id = ?
	`, userId)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}
	u := User{}
	rows.Scan(
		&u.id,
		&u.Email,
		&u.DisplayName,
		&u.FirstName,
		&u.LastName,
		&u.password,
	)
	return &u, nil
}

func (u *User) Send(w http.ResponseWriter) {
	body, err := json.Marshal(*u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}