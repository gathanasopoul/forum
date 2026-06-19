package services

import (
	"database/sql"
)

// AddPostReaction creates a like or dislike for a post.
func AddPostReaction(
	db *sql.DB,
	userID int,
	postID int,
	value int,
) error {

	query := `
	INSERT INTO post_reactions (
		user_id,
		post_id,
		value
	)
	VALUES (?, ?, ?)
	`

	_, err := db.Exec(
		query,
		userID,
		postID,
		value,
	)

	return err
}
