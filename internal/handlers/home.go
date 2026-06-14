package handlers

import (
	"html/template"
	"net/http"
)

// Register page
func (h *Handlers) RegisterPage(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		"templates/register.html",
	)

	if err != nil {
		http.Error(
			w,
			"Template Error",
			http.StatusInternalServerError,
		)
		return
	}

	tmpl.Execute(w, nil)
}
