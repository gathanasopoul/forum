package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/auth"
	"forum/internal/models"
	"forum/internal/services"
)

// CreatePostPage handles displaying the creation form (GET) and creating the post (POST).
func (h *Handlers) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetCurrentUserID(h.DB, r)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		h.handleCreatePostPost(w, r, userID)
		return
	}

	// Categories predefined for post options
	categories := []string{"Go", "Web Development", "Databases", "General Talk"}

	tmpl, err := template.ParseFiles("templates/create-post.html")
	if err != nil {
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
		"Categories": categories,
	})
}

func (h *Handlers) handleCreatePostPost(w http.ResponseWriter, r *http.Request, userID int) {
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	
	r.ParseForm()
	selectedCategories := r.Form["categories"]

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}
	if len(selectedCategories) == 0 {
		http.Error(w, "Please select at least one category", http.StatusBadRequest)
		return
	}

	post := models.Post{
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	_, err := services.CreatePost(h.DB, post, selectedCategories)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ViewPostPage handles displaying a single post by ID.
func (h *Handlers) ViewPostPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	post, err := services.GetPostByID(h.DB, id)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.Error(w, "Post Not Found", http.StatusNotFound)
		return
	}

	userID, _ := auth.GetCurrentUserID(h.DB, r)

	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
		"Post":   post,
		"UserID": userID,
	})
}
