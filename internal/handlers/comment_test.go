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

func TestCreateCommentUnauthenticated(t *testing.T) {
	h := New(setupHandlerTestDB(t))

	form := url.Values{}
	form.Set("content", "This is a test comment")
	form.Set("post_id", "1")

	req := httptest.NewRequest(http.MethodPost, "/comment/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.CreateComment(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
}

func TestCreateCommentAuthenticatedSuccess(t *testing.T) {
	db := setupHandlerTestDB(t)
	h := New(db)

	services.CreateUser(db, models.User{Username: "testuser", Email: "test@test.com", Password: "pwd"})
	user, _ := services.GetUserByEmail(db, "test@test.com")
	token, _ := services.CreateSession(db, user.ID)
	services.CreatePost(db, models.Post{UserID: user.ID, Title: "Title", Content: "Content"}, []string{"Go"})

	form := url.Values{}
	form.Set("content", "This is a valid comment")
	form.Set("post_id", "1")

	req := httptest.NewRequest(http.MethodPost, "/comment/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_token", Value: token})
	rec := httptest.NewRecorder()

	h.CreateComment(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
}

func TestCreateCommentEmptyContent(t *testing.T) {
	db := setupHandlerTestDB(t)
	h := New(db)

	services.CreateUser(db, models.User{Username: "testuser2", Email: "test2@test.com", Password: "pwd"})
	user, _ := services.GetUserByEmail(db, "test2@test.com")
	token, _ := services.CreateSession(db, user.ID)

	form := url.Values{}
	form.Set("content", "  ")
	form.Set("post_id", "1")

	req := httptest.NewRequest(http.MethodPost, "/comment/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_token", Value: token})
	rec := httptest.NewRecorder()

	h.CreateComment(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
