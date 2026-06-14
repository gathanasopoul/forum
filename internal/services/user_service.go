package services

import (
	"database/sql"
	"forum/internal/models"
)

// Create user
func CreateUser(db *sql.DB, user models.User) error {

	query := `
	INSERT INTO users (
		username,
		email,
		password
	)
	VALUES (?, ?, ?)
	`

	_, err := db.Exec(
		query,
		user.Username,
		user.Email,
		user.Password,
	)

	return err
}

// GetUserByEmail finds a user by email. Returns nil if not found.
func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = ?`

	row := db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
