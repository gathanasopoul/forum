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

// AddCommentReaction creates or updates a comment reaction.
func AddCommentReaction(
	db *sql.DB,
	userID int,
	commentID int,
	value int,
) error {

	var reactionID int

	err := db.QueryRow(
		`
		SELECT id
		FROM comment_reactions
		WHERE user_id = ?
		AND comment_id = ?
		`,
		userID,
		commentID,
	).Scan(&reactionID)

	if err == sql.ErrNoRows {

		_, err = db.Exec(
			`
			INSERT INTO comment_reactions (
				user_id,
				comment_id,
				value
			)
			VALUES (?, ?, ?)
			`,
			userID,
			commentID,
			value,
		)

		return err
	}

	if err != nil {
		return err
	}

	_, err = db.Exec(
		`
		UPDATE comment_reactions
		SET value = ?
		WHERE id = ?
		`,
		value,
		reactionID,
	)

	return err
}

// GetCommentReactionCounts returns likes and dislikes for a comment.
func GetCommentReactionCounts(
	db *sql.DB,
	commentID int,
) (int, int, error) {

	var likes int
	var dislikes int

	err := db.QueryRow(
		`
		SELECT COUNT(*)
		FROM comment_reactions
		WHERE comment_id = ?
		AND value = 1
		`,
		commentID,
	).Scan(&likes)

	if err != nil {
		return 0, 0, err
	}

	err = db.QueryRow(
		`
		SELECT COUNT(*)
		FROM comment_reactions
		WHERE comment_id = ?
		AND value = -1
		`,
		commentID,
	).Scan(&dislikes)

	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}
