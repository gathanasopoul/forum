package handlers

import (
	"database/sql"
	"html/template"
	"path/filepath"
)

// Handlers holds shared dependencies for all HTTP handlers.
type Handlers struct {
	DB *sql.DB
}

// New creates a Handlers instance with the database connection.
func New(db *sql.DB) *Handlers {
	return &Handlers{DB: db}
}

// parseTemplateWithBase loads the base template and a page template.
// It automatically falls back to test paths if execution happens from subdirectories.
func parseTemplateWithBase(pageTmplName string, funcs template.FuncMap) (*template.Template, error) {
	paths := []string{
		filepath.Join("templates", pageTmplName),
		filepath.Join("..", "..", "templates", pageTmplName),
	}

	var lastErr error
	for _, path := range paths {
		dir := filepath.Dir(path)
		basePath := filepath.Join(dir, "base.html")

		t := template.New("base")
		if funcs != nil {
			t = t.Funcs(funcs)
		}

		tmpl, err := t.ParseFiles(basePath, path)
		if err == nil {
			return tmpl, nil
		}
		lastErr = err
	}
	return nil, lastErr
}
