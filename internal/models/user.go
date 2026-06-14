package models

import "time"

// User entity
type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}
