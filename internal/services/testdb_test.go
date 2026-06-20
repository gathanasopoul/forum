package services

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates an in-memory SQLite database with the project schema.
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", "file:memdb1?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	db.SetMaxOpenConns(1)

	schema, err := os.ReadFile("../database/schema.sql")
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	if _, err = db.Exec(string(schema)); err != nil {
		t.Fatalf("run schema: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}
