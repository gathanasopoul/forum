package models

import "time"

// Session represents an active login session stored in the sessions table.
type Session struct {
	ID        int
	UserID    int
	Token     string // UUID sent to the browser as a cookie
	ExpiresAt time.Time
}
