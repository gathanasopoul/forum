package models

import "time"

// User represents a registered forum user stored in the users table.
type User struct {
	ID        int
	Username  string
	Email     string
	Password  string // bcrypt hash, never plain text
	CreatedAt time.Time
}
