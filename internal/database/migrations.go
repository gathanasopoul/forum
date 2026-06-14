package database

import (
	"database/sql"
	"os"
)

// Execute schema.sql
func RunMigrations(db *sql.DB) error {

	schema, err := os.ReadFile("internal/database/schema.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return err
	}

	return nil
}
