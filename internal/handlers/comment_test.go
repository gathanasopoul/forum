package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
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
