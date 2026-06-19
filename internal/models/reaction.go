package models

// Reaction represents a like or dislike.
type Reaction struct {
	ID     int
	UserID int
	PostID int
	Value  int
}
