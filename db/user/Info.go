package user

import (
	"database/sql"
)

type Preferences struct {
	Theme string `json:"theme"`
	Kilos bool   `json:"kilos"`
}

type PersonalInfo struct {
	Email       string      `json:"email"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	DisplayName string      `json:"displayName"`
	Preferences Preferences `json:"preferences"`
}

type PublicInfo struct {
	DisplayName string `json:"displayName"`
}

// Retrieves the public info of a user from the database by their display name.
//
// Possible error values:
//   - `ErrDisplayNameNotFound`
//   - `ErrServer`
//   - `ErrDatabase`
func GetPublicInfo(db *sql.DB, displayNames ...string) ([]PublicInfo, error) {
	// this function is trivial right now because there is not yet much
	// info about users. In the future, this might include a bio, a profile
	// pic, etc.
	publicInfo := make([]PublicInfo, len(displayNames))
	for i, displayName := range displayNames {
		rows, err := db.Query(`
			SELECT display_name
			FROM user_data
			WHERE display_name = ?
			`,
			displayName,
		)
		if err != nil {
			return nil, ErrDatabase
		}
		if !rows.Next() {
			return nil, ErrDisplayNameNotFound
		}

		err = rows.Scan(
			&(publicInfo[i].DisplayName),
		)
		if err != nil {
			return nil, ErrServer
		}
	}
	return publicInfo, nil
}

// Retrieves the personal info of a user from the database by their email.
//
// Possible error values:
//   - `ErrEmailNotFound`
//   - `ErrServer`
//   - `ErrDatabase`
func GetPersonalInfo(db *sql.DB, email string) (*PersonalInfo, error) {
	rows, err := db.Query(`
		SELECT 
			d.first_name, 
			d.last_name, 
			d.display_name,
			p.theme,
			p.kilos
		FROM
			user_data d,
			user_preferences p
		WHERE
			d.email = ? AND
			p.email = d.email
		`,
		email,
	)
	if err != nil {
		return nil, ErrDatabase
	}
	if !rows.Next() {
		return nil, ErrEmailNotFound
	}

	personal := PersonalInfo{
		Email: email,
	}
	err = rows.Scan(
		&personal.FirstName,
		&personal.LastName,
		&personal.DisplayName,
		&personal.Preferences.Theme,
		&personal.Preferences.Kilos,
	)
	if err != nil {
		return nil, ErrServer
	}

	return &personal, nil
}
