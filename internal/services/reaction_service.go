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

// GetPostReactionCounts returns likes and dislikes for a post.
func GetPostReactionCounts(
	db *sql.DB,
	postID int,
) (int, int, error) {

	var likes int
	var dislikes int

	err := db.QueryRow(
		`
		SELECT COUNT(*)
		FROM post_reactions
		WHERE post_id = ?
		AND value = 1
		`,
		postID,
	).Scan(&likes)

	if err != nil {
		return 0, 0, err
	}

	err = db.QueryRow(
		`
		SELECT COUNT(*)
		FROM post_reactions
		WHERE post_id = ?
		AND value = -1
		`,
		postID,
	).Scan(&dislikes)

	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}
