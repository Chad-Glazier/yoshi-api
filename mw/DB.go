package mw

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Grants access to:
//   - `DB`
//
// May send back the following error:
//   - `http.StatusInternalServerError` (text body)
func DB(res *Resources, w http.ResponseWriter, r *http.Request) (bool, CleanupFunc) {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to connect to database"))
		return false, nil
	}
	res.DB = db
	return true, cleanupDB
}

func cleanupDB(res *Resources) {
	res.DB.Close()
}
