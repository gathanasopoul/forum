package models

import "time"

type Post struct {
	// Post represents a forum post stored in the posts table.
	ID         int
	UserID     int
	Username   string // Resolved from users table for display
	Title      string
	Content    string
	CreatedAt  time.Time
	Categories []string // List of category names linked to this post
}
