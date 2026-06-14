package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Open SQLite connection
func Connect() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
