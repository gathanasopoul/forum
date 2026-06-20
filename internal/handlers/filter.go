package handlers

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"forum/internal/auth"
	"forum/internal/models"
	"forum/internal/services"
)

const (
	filterAll    = "all"
	filterMine   = "mine"
	filterLiked  = "liked"
	filterCategory = "category"
)

// HomeFilter describes the active post filter on the home page.
type HomeFilter struct {
	ActiveType string
	Category   string
	Label      string
}

// homePageData holds template data for the home page.
type homePageData struct {
	Posts      []*models.Post
	UserID     int
	Categories []string
	Filter     HomeFilter
}

// parseHomeFilter reads the filter selection from the query string.
func parseHomeFilter(r *http.Request) HomeFilter {
	category := strings.TrimSpace(r.URL.Query().Get("category"))
	if category != "" {
		return HomeFilter{
			ActiveType: filterCategory,
			Category:   category,
			Label:      category,
		}
	}

	switch strings.TrimSpace(r.URL.Query().Get("filter")) {
	case filterMine:
		return HomeFilter{ActiveType: filterMine, Label: "My Posts"}
	case filterLiked:
		return HomeFilter{ActiveType: filterLiked, Label: "Liked Posts"}
	default:
		return HomeFilter{ActiveType: filterAll, Label: "All Posts"}
	}
}

// requiresAuth reports whether the filter is only available to logged-in users.
func (f HomeFilter) requiresAuth() bool {
	return f.ActiveType == filterMine || f.ActiveType == filterLiked
}

// loadPostsForFilter returns posts matching the selected filter.
func (h *Handlers) loadPostsForFilter(filter HomeFilter, userID int) ([]*models.Post, error) {
	switch filter.ActiveType {
	case filterCategory:
		return services.GetPostsByCategory(h.DB, filter.Category)
	case filterMine:
		return services.GetPostsByUserID(h.DB, userID)
	case filterLiked:
		return services.GetLikedPostsByUserID(h.DB, userID)
	default:
		return services.GetAllPosts(h.DB)
	}
}

// renderHomePage renders the home page with optional post filtering.
func (h *Handlers) renderHomePage(w http.ResponseWriter, r *http.Request) {
	filter := parseHomeFilter(r)

	userID, err := auth.GetCurrentUserID(h.DB, r)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if filter.requiresAuth() && userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	posts, err := h.loadPostsForFilter(filter, userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	categories, err := services.GetAllCategories(h.DB)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("home.html").Funcs(template.FuncMap{
		"categoryURL": func(name string) template.URL {
			return template.URL("/?category=" + url.QueryEscape(name))
		},
	}).ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}

	data := homePageData{
		Posts:      posts,
		UserID:     userID,
		Categories: categories,
		Filter:     filter,
	}

	tmpl.Execute(w, data)
}
