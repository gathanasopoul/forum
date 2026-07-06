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

func TestCreatePostPageGET(t *testing.T) {
	// Not authenticated, should redirect to login
	h := New(setupHandlerTestDB(t))

	req := httptest.NewRequest(http.MethodGet, "/create-post", nil)
	rec := httptest.NewRecorder()

	h.CreatePostPage(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
}

func TestViewPostPageInvalidID(t *testing.T) {
	h := New(setupHandlerTestDB(t))

	req := httptest.NewRequest(http.MethodGet, "/post?id=invalid", nil)
	rec := httptest.NewRecorder()

	h.ViewPostPage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestViewPostPageNotFound(t *testing.T) {
	h := New(setupHandlerTestDB(t))

	req := httptest.NewRequest(http.MethodGet, "/post?id=999", nil)
	rec := httptest.NewRecorder()

	h.ViewPostPage(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestCreatePostPagePOSTAuthenticatedSuccess(t *testing.T) {
	db := setupHandlerTestDB(t)
	h := New(db)

	services.CreateUser(db, models.User{Username: "testuser", Email: "test@test.com", Password: "pwd"})
	user, _ := services.GetUserByEmail(db, "test@test.com")
	token, _ := services.CreateSession(db, user.ID)

	form := url.Values{}
	form.Set("title", "My Test Post")
	form.Set("content", "This is the content")
	form.Add("categories", "Go")

	req := httptest.NewRequest(http.MethodPost, "/create-post", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_token", Value: token})
	
	rec := httptest.NewRecorder()
	h.CreatePostPage(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
}

func TestCreatePostPagePOSTEmptyData(t *testing.T) {
	db := setupHandlerTestDB(t)
	h := New(db)

	services.CreateUser(db, models.User{Username: "testuser2", Email: "test2@test.com", Password: "pwd"})
	user, _ := services.GetUserByEmail(db, "test2@test.com")
	token, _ := services.CreateSession(db, user.ID)

	form := url.Values{}
	form.Set("title", "")
	form.Set("content", "This is the content")
	form.Add("categories", "Go")

	req := httptest.NewRequest(http.MethodPost, "/create-post", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_token", Value: token})
	
	rec := httptest.NewRecorder()
	h.CreatePostPage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

