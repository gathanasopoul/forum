package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"forum/internal/models"
	"forum/internal/services"
)

func TestLikePostAuthenticated(t *testing.T) {
	db := setupHandlerTestDB(t)
	h := New(db)

	services.CreateUser(db, models.User{Username: "testuser", Email: "test@test.com", Password: "pwd"})
	user, _ := services.GetUserByEmail(db, "test@test.com")
	token, _ := services.CreateSession(db, user.ID)
	services.CreatePost(db, models.Post{UserID: user.ID, Title: "Title", Content: "Content"}, []string{"Go"})

	form := url.Values{}
	form.Set("post_id", "1")

	req := httptest.NewRequest(http.MethodPost, "/post/like", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_token", Value: token})
	rec := httptest.NewRecorder()

	h.LikePost(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
}

func TestLikePostUnauthenticated(t *testing.T) {
	db := setupHandlerTestDB(t)
	h := New(db)

	form := url.Values{}
	form.Set("post_id", "1")

	req := httptest.NewRequest(http.MethodPost, "/post/like", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.LikePost(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if !strings.Contains(rec.Header().Get("Location"), "/login") {
		t.Errorf("expected redirect to login")
	}
}

func TestLikePostInvalidMethod(t *testing.T) {
	h := New(setupHandlerTestDB(t))

	req := httptest.NewRequest(http.MethodGet, "/post/like", nil)
	rec := httptest.NewRecorder()

	h.LikePost(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
