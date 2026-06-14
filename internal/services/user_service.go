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
