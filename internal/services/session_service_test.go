package services

import (
	"database/sql"
	"forum/internal/models"
	"testing"
)

func createTestUser(t *testing.T, db *sql.DB) int {
	t.Helper()

	user := models.User{
		Username: "testuser",
		Email:    "test@gmail.com",
		Password: "hashed-password",
	}
	if err := CreateUser(db, user); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	found, err := GetUserByEmail(db, "test@gmail.com")
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}
	return found.ID
}

func TestCreateSessionAndGetByToken(t *testing.T) {
	db := setupTestDB(t)
	userID := createTestUser(t, db)

	token, err := CreateSession(db, userID)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	session, err := GetSessionByToken(db, token)
	if err != nil {
		t.Fatalf("GetSessionByToken: %v", err)
	}
	if session == nil {
		t.Fatal("expected session, got nil")
	}
	if session.UserID != userID {
		t.Errorf("userID = %d, want %d", session.UserID, userID)
	}
}

func TestGetSessionByTokenNotFound(t *testing.T) {
	db := setupTestDB(t)

	session, err := GetSessionByToken(db, "missing-token")
	if err != nil {
		t.Fatalf("GetSessionByToken: %v", err)
	}
	if session != nil {
		t.Fatal("expected nil session for missing token")
	}
}

func TestGetSessionByTokenExpired(t *testing.T) {
	db := setupTestDB(t)
	userID := createTestUser(t, db)

	_, err := db.Exec(
		`INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, datetime('now', '-1 hour'))`,
		userID, "expired-token",
	)
	if err != nil {
		t.Fatalf("insert expired session: %v", err)
	}

	session, err := GetSessionByToken(db, "expired-token")
	if err != nil {
		t.Fatalf("GetSessionByToken: %v", err)
	}
	if session != nil {
		t.Fatal("expected nil for expired session")
	}
}

func TestDeleteSession(t *testing.T) {
	db := setupTestDB(t)
	userID := createTestUser(t, db)

	token, err := CreateSession(db, userID)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	if err := DeleteSession(db, token); err != nil {
		t.Fatalf("DeleteSession: %v", err)
	}

	session, err := GetSessionByToken(db, token)
	if err != nil {
		t.Fatalf("GetSessionByToken: %v", err)
	}
	if session != nil {
		t.Fatal("expected session to be deleted")
	}
}
