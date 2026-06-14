package handlers

import "database/sql"

// Handlers holds shared dependencies for all HTTP handlers.
type Handlers struct {
	DB *sql.DB
}

// New creates a Handlers instance with the database connection.
func New(db *sql.DB) *Handlers {
	return &Handlers{DB: db}
}
