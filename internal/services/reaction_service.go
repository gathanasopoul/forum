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

	var reactionID int

	err := db.QueryRow(
		`
		SELECT id
		FROM post_reactions
		WHERE user_id = ?
		AND post_id = ?
		`,
		userID,
		postID,
	).Scan(&reactionID)

	if err == sql.ErrNoRows {

		_, err = db.Exec(
			`
			INSERT INTO post_reactions (
				user_id,
				post_id,
				value
			)
			VALUES (?, ?, ?)
			`,
			userID,
			postID,
			value,
		)

		return err
	}

	if err != nil {
		return err
	}

	_, err = db.Exec(
		`
		UPDATE post_reactions
		SET value = ?
		WHERE id = ?
		`,
		value,
		reactionID,
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
