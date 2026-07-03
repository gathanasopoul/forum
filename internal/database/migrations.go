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

	// Migrate existing users table to have oauth columns
	_, _ = db.Exec("ALTER TABLE users ADD COLUMN oauth_provider TEXT;")
	_, _ = db.Exec("ALTER TABLE users ADD COLUMN oauth_id TEXT;")

	return nil
}
