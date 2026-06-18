package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"forum/internal/auth"
	"forum/internal/models"
	"forum/internal/services"
)

// CreateComment handles comment creation.
func (h *Handlers) CreateComment(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	userID, err := auth.GetCurrentUserID(
		h.DB,
		r,
	)
	if err != nil {
		http.Error(
			w,
			"Server error",
			http.StatusInternalServerError,
		)
		return
	}

	if userID == 0 {
		http.Redirect(
			w,
			r,
			"/login",
			http.StatusSeeOther,
		)
		return
	}

	postID, err := strconv.Atoi(
		r.FormValue("post_id"),
	)
	if err != nil {
		http.Error(
			w,
			"Invalid post",
			http.StatusBadRequest,
		)
		return
	}

	content := strings.TrimSpace(
		r.FormValue("content"),
	)

	if content == "" {
		http.Error(
			w,
			"Comment cannot be empty",
			http.StatusBadRequest,
		)
		return
	}

	comment := models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}

	err = services.CreateComment(
		h.DB,
		comment,
	)
	if err != nil {
		http.Error(
			w,
			"Server error",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/post?id="+strconv.Itoa(postID),
		http.StatusSeeOther,
	)
}
