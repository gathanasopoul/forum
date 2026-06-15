package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"forum/internal/models"
	"forum/internal/services"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

func setupHandlerTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

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

func TestRegisterPostEmptyEmail(t *testing.T) {
	h := New(setupHandlerTestDB(t))

	form := url.Values{}
	form.Set("username", "salas")
	form.Set("email", "")
	form.Set("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.RegisterPage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "Email is required") {
		t.Errorf("body = %q, want email required message", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Register</h1>") {
		t.Error("expected register form to be re-rendered")
	}
}

func TestRegisterPostDuplicateEmail(t *testing.T) {
	db := setupHandlerTestDB(t)
	h := New(db)

	hashed, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if err := services.CreateUser(db, models.User{
		Username: "existing",
		Email:    "taken@gmail.com",
		Password: string(hashed),
	}); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	form := url.Values{}
	form.Set("username", "newuser")
	form.Set("email", "taken@gmail.com")
	form.Set("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.RegisterPage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "Email already taken") {
		t.Errorf("body = %q, want duplicate email message", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `value="newuser"`) {
		t.Error("expected username value to be preserved in form")
	}
}

func TestRegisterPostDuplicateUsernameAndEmail(t *testing.T) {
	db := setupHandlerTestDB(t)
	h := New(db)

	hashed, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if err := services.CreateUser(db, models.User{
		Username: "salas",
		Email:    "salas@gmail.com",
		Password: string(hashed),
	}); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	form := url.Values{}
	form.Set("username", "salas")
	form.Set("email", "salas@gmail.com")
	form.Set("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.RegisterPage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "Username already taken") {
		t.Error("expected username already taken message")
	}
	if !strings.Contains(rec.Body.String(), "Email already taken") {
		t.Error("expected email already taken message")
	}
}

func TestLoginPostEmptyPassword(t *testing.T) {
	h := New(setupHandlerTestDB(t))

	form := url.Values{}
	form.Set("email", "salas@gmail.com")
	form.Set("password", "")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.LoginPage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "Password is required") {
		t.Errorf("body = %q, want password required message", rec.Body.String())
	}
}

func TestLoginPostInvalidCredentials(t *testing.T) {
	db := setupHandlerTestDB(t)
	h := New(db)

	hashed, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if err := services.CreateUser(db, models.User{
		Username: "salas",
		Email:    "salas@gmail.com",
		Password: string(hashed),
	}); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	form := url.Values{}
	form.Set("email", "salas@gmail.com")
	form.Set("password", "wrong-password")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.LoginPage(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if !strings.Contains(rec.Body.String(), "Invalid credentials") {
		t.Errorf("body = %q, want invalid credentials message", rec.Body.String())
	}
}
