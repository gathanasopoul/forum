package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
