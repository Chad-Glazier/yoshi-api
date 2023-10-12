package db

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	return sql.Open("mysql", os.Getenv("DSN"))
} 

func Close(connection *sql.DB) (error) {
	return connection.Close()
}