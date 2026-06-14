package handlers

import (
	"html/template"
	"net/http"

	"forum/internal/models"
	"forum/internal/services"

	"golang.org/x/crypto/bcrypt"
)

// Register page
func (h *Handlers) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.handleRegisterPost(w, r)
		return
	}

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

func (h *Handlers) handleRegisterPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	existing, err := services.GetUserByEmail(h.DB, email)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if existing != nil {
		http.Error(w, "Email already taken", http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashed),
	}
	err = services.CreateUser(h.DB, user)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
