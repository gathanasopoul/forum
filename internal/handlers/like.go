package handlers

import (
	"net/http"
	"strconv"

	"forum/internal/auth"
	"forum/internal/services"
)

// LikePost handles post likes.
func (h *Handlers) LikePost(
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

	err = services.AddPostReaction(
		h.DB,
		userID,
		postID,
		1,
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

// DislikePost handles post dislikes.
func (h *Handlers) DislikePost(
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

	err = services.AddPostReaction(
		h.DB,
		userID,
		postID,
		-1,
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

// LikeComment handles comment likes.
func (h *Handlers) LikeComment(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.GetCurrentUserID(h.DB, r)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentID, err := strconv.Atoi(
		r.FormValue("comment_id"),
	)
	if err != nil {
		http.Error(w, "Invalid comment", http.StatusBadRequest)
		return
	}

	postID := r.FormValue("post_id")

	err = services.AddCommentReaction(
		h.DB,
		userID,
		commentID,
		1,
	)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/post?id="+postID,
		http.StatusSeeOther,
	)
}

// DislikeComment handles comment dislikes.
func (h *Handlers) DislikeComment(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.GetCurrentUserID(h.DB, r)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentID, err := strconv.Atoi(
		r.FormValue("comment_id"),
	)
	if err != nil {
		http.Error(w, "Invalid comment", http.StatusBadRequest)
		return
	}

	postID := r.FormValue("post_id")

	err = services.AddCommentReaction(
		h.DB,
		userID,
		commentID,
		-1,
	)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/post?id="+postID,
		http.StatusSeeOther,
	)
}
