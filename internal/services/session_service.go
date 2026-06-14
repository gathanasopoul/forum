package services

import (
	"database/sql"
	"forum/internal/models"
	"time"

	"github.com/google/uuid"
)

func CreateSession(db *sql.DB, userID int) (string, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	query := `INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)`
	_, err := db.Exec(query, userID, token, expiresAt)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetSessionByToken(db *sql.DB, token string) (*models.Session, error) {
	query := `SELECT id, user_id, token, expires_at FROM sessions WHERE token = ?`

	row := db.QueryRow(query, token)
	var s models.Session
	err := row.Scan(&s.ID, &s.UserID, &s.Token, &s.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if time.Now().After(s.ExpiresAt) {
		return nil, nil
	}
	return &s, nil
}

func DeleteSession(db *sql.DB, token string) error {
	query := `DELETE FROM sessions WHERE token = ?`
	_, err := db.Exec(query, token)
	return err
}
