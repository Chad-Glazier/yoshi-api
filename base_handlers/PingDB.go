package base_handlers

import (
	"database/sql"
	"net/http"
	"os"
	"yoshi-api/types"

	_ "github.com/go-sql-driver/mysql"
)

func PingDB(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		(&types.ApiError{
			Status:  http.StatusInternalServerError,
			Message: "failed to connect to database",
		}).Send(w)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		(&types.ApiError{
			Status: http.StatusServiceUnavailable,
			Message: "database failed to respond",
		}).Send(w)
	}

	w.WriteHeader(http.StatusOK)
}
