package auth

import (
	"database/sql"
	"net/http"
	"time"

	"forum/internal/services"
)

const CookieName = "session_token"

func SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
}

func GetSessionToken(r *http.Request) string {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

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
