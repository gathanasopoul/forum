package database

import (
	"os"
	"testing"
)

func TestConnect(t *testing.T) {
	// Use a temporary database file instead of the actual forum.db
	// By temporarily changing the working directory or relying on Connect's default path
	// Connect() opens "./forum.db". If we create it and then delete it, it should be fine.
	// But it's better to just call Connect and test if it works, then close and clean up.

	db, err := Connect()
	if err != nil {
		t.Fatalf("Expected no error connecting to db, got %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("Expected ping to succeed, got %v", err)
	}

	// Clean up the created forum.db
	err = os.Remove("./forum.db")
	if err != nil && !os.IsNotExist(err) {
		t.Logf("Warning: failed to clean up test database file: %v", err)
	}
}
