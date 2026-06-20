package services

import (
	"database/sql"

	"forum/internal/models"
)

// GetPostsByCategory returns posts linked to the given category name.
func GetPostsByCategory(db *sql.DB, category string) ([]*models.Post, error) {
	query := `
		SELECT DISTINCT p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_categories pc ON p.id = pc.post_id
		JOIN categories c ON pc.category_id = c.id
		WHERE c.name = ?
		ORDER BY p.created_at DESC`

	return queryPostsWithCategories(db, query, category)
}

// GetPostsByUserID returns posts created by the given user.
func GetPostsByUserID(db *sql.DB, userID int) ([]*models.Post, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC`

	return queryPostsWithCategories(db, query, userID)
}

// GetLikedPostsByUserID returns posts the user has liked.
func GetLikedPostsByUserID(db *sql.DB, userID int) ([]*models.Post, error) {
	query := `
		SELECT DISTINCT p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_reactions pr ON p.id = pr.post_id
		WHERE pr.user_id = ? AND pr.value = 1
		ORDER BY p.created_at DESC`

	return queryPostsWithCategories(db, query, userID)
}
