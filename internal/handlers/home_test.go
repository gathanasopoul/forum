package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"forum/internal/models"
	"forum/internal/services"
)

func TestHomePage(t *testing.T) {
	h := New(setupHandlerTestDB(t))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.HomePage(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestHomePageNotFound(t *testing.T) {
	h := New(setupHandlerTestDB(t))

	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	rec := httptest.NewRecorder()

	h.HomePage(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestHomePageFilterUnauthenticated(t *testing.T) {
	h := New(setupHandlerTestDB(t))

	req := httptest.NewRequest(http.MethodGet, "/?filter=mine", nil)
	rec := httptest.NewRecorder()

	h.HomePage(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
}

func TestHomePageFilterAuthenticated(t *testing.T) {
	db := setupHandlerTestDB(t)
	h := New(db)

	services.CreateUser(db, models.User{Username: "testuser", Email: "test@test.com", Password: "pwd"})
	user, _ := services.GetUserByEmail(db, "test@test.com")
	token, _ := services.CreateSession(db, user.ID)

	req := httptest.NewRequest(http.MethodGet, "/?filter=mine", nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: token})
	rec := httptest.NewRecorder()

	h.HomePage(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
