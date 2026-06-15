package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"forum/internal/models"
	"forum/internal/services"

	"golang.org/x/crypto/bcrypt"
)

// RegisterPage handles GET and POST requests to /register.
// GET shows the registration form, POST creates a new user account.
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

// handleRegisterPost processes the registration form submission.
// It validates input, checks for duplicate username/email, hashes the password,
// saves the user to the database, and redirects to the login page.
func (h *Handlers) handleRegisterPost(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	existing, err := services.GetUserByUsername(h.DB, username)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if existing != nil {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	existing, err = services.GetUserByEmail(h.DB, email)
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
