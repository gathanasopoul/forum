package services

import (
	"database/sql"

	"forum/internal/models"
)

// Create comment
func CreateComment(
	db *sql.DB,
	comment models.Comment,
) error {

	query := `
	INSERT INTO comments (
		post_id,
		user_id,
		content
	)
	VALUES (?, ?, ?)
	`

	_, err := db.Exec(
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	)

	return err
}

// Get comments for post
func GetCommentsByPostID(
	db *sql.DB,
	postID int,
) ([]*models.Comment, error) {

	query := `
		SELECT
			c.id,
			c.post_id,
			c.user_id,
			u.username,
			c.content,
			c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`

	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment

	for rows.Next() {
		var c models.Comment

		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Username,
			&c.Content,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		likes, dislikes, err := GetCommentReactionCounts(
			db,
			c.ID,
		)
		if err != nil {
			return nil, err
		}

		c.Likes = likes
		c.Dislikes = dislikes

		comments = append(comments, &c)
	}

	return comments, nil
}
