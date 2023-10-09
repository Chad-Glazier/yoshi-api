package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"yoshi-api/types"
	"yoshi-api/util"

	_ "github.com/go-sql-driver/mysql"
)

func PingDB(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		util.SendError(w, types.ErrorInfo{
			Status: http.StatusInternalServerError,
			Message: "failed to connect to database",
		})
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		util.SendError(w, types.ErrorInfo{
			Status: http.StatusServiceUnavailable,
			Message: "database failed to respond",
		})
	}

	w.WriteHeader(200)
}
