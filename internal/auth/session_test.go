package auth

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"forum/internal/models"
	"forum/internal/services"

	_ "github.com/mattn/go-sqlite3"
)

func TestSetSessionCookie(t *testing.T) {
	w := httptest.NewRecorder()
	token := "test-token"

	SetSessionCookie(w, token)

	resp := w.Result()
	cookies := resp.Cookies()

	if len(cookies) == 0 {
		t.Fatalf("Expected cookie to be set, got none")
	}

	cookie := cookies[0]
	if cookie.Name != CookieName {
		t.Errorf("Expected cookie name %s, got %s", CookieName, cookie.Name)
	}
	if cookie.Value != token {
		t.Errorf("Expected cookie value %s, got %s", token, cookie.Value)
	}
	if cookie.HttpOnly != true {
		t.Errorf("Expected HttpOnly to be true")
	}
	if cookie.Path != "/" {
		t.Errorf("Expected Path to be /, got %s", cookie.Path)
	}
}

func TestGetSessionToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	token := "test-token"

	cookie := &http.Cookie{
		Name:  CookieName,
		Value: token,
	}
	req.AddCookie(cookie)

	got := GetSessionToken(req)
	if got != token {
		t.Errorf("Expected token %s, got %s", token, got)
	}
}

func TestGetSessionTokenEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	got := GetSessionToken(req)
	if got != "" {
		t.Errorf("Expected empty token, got %s", got)
	}
}

func TestClearSessionCookie(t *testing.T) {
	w := httptest.NewRecorder()

	ClearSessionCookie(w)

	resp := w.Result()
	cookies := resp.Cookies()

	if len(cookies) == 0 {
		t.Fatalf("Expected cookie to be set, got none")
	}

	cookie := cookies[0]
	if cookie.Name != CookieName {
		t.Errorf("Expected cookie name %s, got %s", CookieName, cookie.Name)
	}
	if cookie.Value != "" {
		t.Errorf("Expected cookie value to be empty, got %s", cookie.Value)
	}
	if cookie.Expires.After(time.Now()) {
		t.Errorf("Expected cookie to be expired")
	}
}

func TestGetCurrentUserID(t *testing.T) {
	// Set up an in-memory database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Read schema
	schemaBytes, err := os.ReadFile("../database/schema.sql")
	if err != nil {
		t.Fatalf("Failed to read schema: %v", err)
	}

	_, err = db.Exec(string(schemaBytes))
	if err != nil {
		t.Fatalf("Failed to execute schema: %v", err)
	}

	// Insert a test user
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	err = services.CreateUser(db, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create a session for the user
	token, err := services.CreateSession(db, 1)
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}

	// Test case: Valid token
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  CookieName,
		Value: token,
	})

	userID, err := GetCurrentUserID(db, req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if userID != 1 {
		t.Errorf("Expected userID 1, got %d", userID)
	}

	// Test case: Invalid token
	reqInvalid := httptest.NewRequest(http.MethodGet, "/", nil)
	reqInvalid.AddCookie(&http.Cookie{
		Name:  CookieName,
		Value: "invalid-token",
	})

	userIDInvalid, err := GetCurrentUserID(db, reqInvalid)
	if err != nil {
		t.Errorf("Expected no error for invalid token, got %v", err)
	}
	if userIDInvalid != 0 {
		t.Errorf("Expected userID 0 for invalid token, got %d", userIDInvalid)
	}

	// Test case: No token
	reqNoToken := httptest.NewRequest(http.MethodGet, "/", nil)
	userIDNoToken, err := GetCurrentUserID(db, reqNoToken)
	if err != nil {
		t.Errorf("Expected no error for no token, got %v", err)
	}
	if userIDNoToken != 0 {
		t.Errorf("Expected userID 0 for no token, got %d", userIDNoToken)
	}
}
