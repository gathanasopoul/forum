package handlers

import (
	"log"
	"net/http"
	"strings"

	"forum/internal/auth"
	"forum/internal/services"

	"golang.org/x/crypto/bcrypt"
)

// loginPageData holds form values and validation errors for the login template.
type loginPageData struct {
	Email string
	Error string
}

// renderLoginPage renders the login template with the given data and status code.
func (h *Handlers) renderLoginPage(w http.ResponseWriter, data loginPageData, status int) {
	tmpl, err := parseTemplateWithBase("login.html", nil)
	if err != nil {
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"UserID": 0,
		"Form":   data,
	})
}

// LoginPage handles GET and POST requests to /login.
// GET shows the login form, POST authenticates the user.
func (h *Handlers) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.handleLoginPost(w, r)
		return
	}

	h.renderLoginPage(w, loginPageData{}, http.StatusOK)
}

// handleLoginPost processes the login form submission.
// It validates input, checks credentials against the database,
// creates a session, sets a cookie, and redirects to the home page.
func (h *Handlers) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if email == "" || password == "" {
		h.renderLoginPage(w, loginPageData{
			Email: email,
			Error: "Email and Password are required",
		}, http.StatusBadRequest)
		return
	}

	user, err := services.GetUserByEmail(h.DB, email)
	if err != nil {
		log.Printf("Login GetUserByEmail error: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		h.renderLoginPage(w, loginPageData{
			Email: email,
			Error: "Invalid credentials",
		}, http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		h.renderLoginPage(w, loginPageData{
			Email: email,
			Error: "Invalid credentials",
		}, http.StatusUnauthorized)
		return
	}

	token, err := services.CreateSession(h.DB, user.ID)
	if err != nil {
		log.Printf("Login CreateSession error: %v", err)
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
