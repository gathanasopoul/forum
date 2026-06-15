package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"forum/internal/models"
	"forum/internal/services"

	"golang.org/x/crypto/bcrypt"
)

var registerTemplatePaths = []string{
	"templates/register.html",
	"../../templates/register.html",
}

// registerPageData holds form values and validation errors for the register template.
type registerPageData struct {
	Username      string
	Email         string
	UsernameError string
	EmailError    string
	PasswordError string
}

// RegisterPage handles GET and POST requests to /register.
// GET shows the registration form, POST creates a new user account.
func (h *Handlers) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.handleRegisterPost(w, r)
		return
	}

	h.renderRegisterPage(w, registerPageData{}, http.StatusOK)
}

// renderRegisterPage renders the register template with the given data and status code.
func (h *Handlers) renderRegisterPage(w http.ResponseWriter, data registerPageData, status int) {
	tmpl, err := parseRegisterTemplate()
	if err != nil {
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	tmpl.Execute(w, data)
}

// parseRegisterTemplate loads the register template from the project root or test path.
func parseRegisterTemplate() (*template.Template, error) {
	var lastErr error
	for _, path := range registerTemplatePaths {
		tmpl, err := template.ParseFiles(path)
		if err == nil {
			return tmpl, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

// handleRegisterPost processes the registration form submission.
// It validates input, checks for duplicate username/email, hashes the password,
// saves the user to the database, and redirects to the login page.
func (h *Handlers) handleRegisterPost(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	data := registerPageData{
		Username: username,
		Email:    email,
	}

	if username == "" {
		data.UsernameError = "Username is required"
	}
	if email == "" {
		data.EmailError = "Email is required"
	}
	if password == "" {
		data.PasswordError = "Password is required"
	}
	if data.UsernameError != "" || data.EmailError != "" || data.PasswordError != "" {
		h.renderRegisterPage(w, data, http.StatusBadRequest)
		return
	}

	existing, err := services.GetUserByUsername(h.DB, username)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if existing != nil {
		data.UsernameError = "Username already taken"
	}

	existing, err = services.GetUserByEmail(h.DB, email)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if existing != nil {
		data.EmailError = "Email already taken"
	}

	if data.UsernameError != "" || data.EmailError != "" {
		h.renderRegisterPage(w, data, http.StatusBadRequest)
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
