package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOAuthLogin(t *testing.T) {
	// Initialize oauth configs just in case, otherwise they might be nil
	InitOAuth()
	
	h := New(setupHandlerTestDB(t))

	req := httptest.NewRequest(http.MethodGet, "/auth/github/login", nil)
	rec := httptest.NewRecorder()

	h.OAuthLogin(rec, req)

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusTemporaryRedirect)
	}

	loc := rec.Header().Get("Location")
	if loc == "" {
		t.Error("expected Location header to be set for redirect")
	}

	cookies := rec.Result().Cookies()
	var foundState bool
	for _, cookie := range cookies {
		if cookie.Name == "oauthstate" {
			foundState = true
			break
		}
	}
	if !foundState {
		t.Error("expected oauthstate cookie to be set")
	}
}

func TestOAuthLoginInvalidProvider(t *testing.T) {
	InitOAuth()
	
	h := New(setupHandlerTestDB(t))

	req := httptest.NewRequest(http.MethodGet, "/auth/invalidprovider/login", nil)
	rec := httptest.NewRecorder()

	h.OAuthLogin(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestOAuthCallbackInvalidState(t *testing.T) {
	InitOAuth()
	h := New(setupHandlerTestDB(t))

	req := httptest.NewRequest(http.MethodGet, "/auth/github/callback?state=invalid", nil)
	rec := httptest.NewRecorder()

	h.OAuthCallback(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestOAuthCallbackInvalidProvider(t *testing.T) {
	InitOAuth()
	h := New(setupHandlerTestDB(t))

	// Provide a matching state to bypass the state check, so we can reach the provider check
	state := "teststate"
	req := httptest.NewRequest(http.MethodGet, "/auth/invalid/callback?state="+state, nil)
	req.AddCookie(&http.Cookie{Name: "oauthstate", Value: state})
	rec := httptest.NewRecorder()

	h.OAuthCallback(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
