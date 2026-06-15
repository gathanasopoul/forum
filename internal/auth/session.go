package auth

import (
	"database/sql"
	"net/http"
	"time"

	"forum/internal/services"
)

// CookieName is the name of the session cookie sent to the browser.
const CookieName = "session_token"

// SetSessionCookie sends the session token to the browser as an HttpOnly cookie.
// The cookie expires after 24 hours.
func SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
}

// GetSessionToken reads the session token from the request cookie.
// Returns an empty string if no cookie is present.
func GetSessionToken(r *http.Request) string {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// ClearSessionCookie tells the browser to delete the session cookie.
// Used during logout.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

// GetCurrentUserID returns the ID of the logged-in user from the session cookie.
// Returns 0 if the user is not logged in or the session is invalid/expired.
// Other handlers use this to check if a user is authenticated.
func GetCurrentUserID(db *sql.DB, r *http.Request) (int, error) {
	token := GetSessionToken(r)
	if token == "" {
		return 0, nil
	}

	session, err := services.GetSessionByToken(db, token)
	if err != nil {
		return 0, err
	}
	if session == nil {
		return 0, nil
	}
	return session.UserID, nil
}
