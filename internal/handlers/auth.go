package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"forum/internal/auth"
	"forum/internal/services"

	"golang.org/x/crypto/bcrypt"
)

// LoginPage handles GET and POST requests to /login.
// GET shows the login form, POST authenticates the user.
func (h *Handlers) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.handleLoginPost(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// handleLoginPost processes the login form submission.
// It validates input, checks credentials against the database,
// creates a session, sets a cookie, and redirects to the home page.
func (h *Handlers) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	user, err := services.GetUserByEmail(h.DB, email)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := services.CreateSession(h.DB, user.ID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(w, token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout ends the current user session.
// It deletes the session from the database, clears the cookie,
// and redirects to the login page.
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	token := auth.GetSessionToken(r)
	if token != "" {
		services.DeleteSession(h.DB, token)
	}
	auth.ClearSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
