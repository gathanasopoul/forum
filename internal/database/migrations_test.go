package database

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestRunMigrations(t *testing.T) {
	// Use an in-memory database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Since RunMigrations reads from "internal/database/schema.sql",
	// when running `go test ./internal/database`, the working directory is `forum/internal/database`.
	// So "internal/database/schema.sql" will resolve to `forum/internal/database/internal/database/schema.sql`, which doesn't exist.
	// Wait, let's fix RunMigrations or create a temporary schema file if needed.
	// We can temporarily create `internal/database/schema.sql` in the test directory to mimic the working dir from root.

	err = os.MkdirAll("internal/database", 0755)
	if err != nil {
		t.Fatalf("Failed to create temporary directory for test schema: %v", err)
	}
	defer os.RemoveAll("internal")

	schemaContent, err := os.ReadFile("schema.sql")
	if err != nil {
		t.Fatalf("Failed to read original schema.sql: %v", err)
	}

	err = os.WriteFile("internal/database/schema.sql", schemaContent, 0644)
	if err != nil {
		t.Fatalf("Failed to write temporary schema.sql: %v", err)
	}

	err = RunMigrations(db)
	if err != nil {
		t.Errorf("Expected RunMigrations to succeed, got %v", err)
	}
}
