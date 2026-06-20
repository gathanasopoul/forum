package services

import (
	"database/sql"
	"forum/internal/models"
)

// CreatePost creates a new post and registers the linked categories.
func CreatePost(db *sql.DB, post models.Post, categories []string) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	query := `INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)`
	res, err := tx.Exec(query, post.UserID, post.Title, post.Content)
	if err != nil {
		return 0, err
	}

	postID64, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	postID := int(postID64)

	for _, catName := range categories {
		var catID int
		// Query category ID, insert if it doesn't exist
		err = tx.QueryRow(`SELECT id FROM categories WHERE name = ?`, catName).Scan(&catID)
		if err == sql.ErrNoRows {
			resCat, err := tx.Exec(`INSERT INTO categories (name) VALUES (?)`, catName)
			if err != nil {
				return 0, err
			}
			catID64, err := resCat.LastInsertId()
			if err != nil {
				return 0, err
			}
			catID = int(catID64)
		} else if err != nil {
			return 0, err
		}

		// Link post and category
		_, err = tx.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, catID)
		if err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return postID, nil
}

// GetPostByID retrieves a single post with its categories and user details.
func GetPostByID(db *sql.DB, id int) (*models.Post, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?`

	row := db.QueryRow(query, id)
	var p models.Post
	err := row.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Fetch associated categories
	categories, err := loadCategoriesForPost(db, p.ID)
	if err != nil {
		return nil, err
	}
	p.Categories = categories

	return &p, nil
}

// GetAllPosts lists all posts on the forum.
func GetAllPosts(db *sql.DB) ([]*models.Post, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC`

	return queryPostsWithCategories(db, query)
}

// queryPostsWithCategories runs a post SELECT and attaches category names to each row.
func queryPostsWithCategories(db *sql.DB, query string, args ...interface{}) ([]*models.Post, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var posts []*models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			rows.Close()
			return nil, err
		}
		posts = append(posts, &p)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	for _, p := range posts {
		categories, err := loadCategoriesForPost(db, p.ID)
		if err != nil {
			return nil, err
		}
		p.Categories = categories
	}

	return posts, nil
}

// loadCategoriesForPost returns the category names linked to a post.
func loadCategoriesForPost(db *sql.DB, postID int) ([]string, error) {
	rows, err := db.Query(`
		SELECT c.name
		FROM categories c
		JOIN post_categories pc ON c.id = pc.category_id
		WHERE pc.post_id = ?`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		categories = append(categories, name)
	}

	return categories, nil
}

// GetAllCategories returns all registered category names.
func GetAllCategories(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`SELECT name FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		categories = append(categories, name)
	}
	return categories, nil
}
