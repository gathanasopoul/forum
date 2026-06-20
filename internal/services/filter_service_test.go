package services

import (
	"forum/internal/models"
	"testing"
)

func TestGetPostsByCategory(t *testing.T) {
	db := setupTestDB(t)
	userID := createTestUser(t, db)

	if _, err := CreatePost(db, models.Post{UserID: userID, Title: "Go Post", Content: "Go content"}, []string{"Go"}); err != nil {
		t.Fatalf("CreatePost Go: %v", err)
	}
	if _, err := CreatePost(db, models.Post{UserID: userID, Title: "Web Post", Content: "Web content"}, []string{"Web Development"}); err != nil {
		t.Fatalf("CreatePost Web: %v", err)
	}

	posts, err := GetPostsByCategory(db, "Go")
	if err != nil {
		t.Fatalf("GetPostsByCategory: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("len(posts) = %d, want 1", len(posts))
	}
	if posts[0].Title != "Go Post" {
		t.Errorf("title = %q, want Go Post", posts[0].Title)
	}
}

func TestGetPostsByUserID(t *testing.T) {
	db := setupTestDB(t)
	userOne := createTestUser(t, db)

	userTwo := models.User{
		Username: "otheruser",
		Email:    "other@gmail.com",
		Password: "hashed-password",
	}
	if err := CreateUser(db, userTwo); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	found, err := GetUserByEmail(db, "other@gmail.com")
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}

	if _, err := CreatePost(db, models.Post{UserID: userOne, Title: "Mine", Content: "My post"}, []string{"General Talk"}); err != nil {
		t.Fatalf("CreatePost mine: %v", err)
	}
	if _, err := CreatePost(db, models.Post{UserID: found.ID, Title: "Theirs", Content: "Other post"}, []string{"General Talk"}); err != nil {
		t.Fatalf("CreatePost theirs: %v", err)
	}

	posts, err := GetPostsByUserID(db, userOne)
	if err != nil {
		t.Fatalf("GetPostsByUserID: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("len(posts) = %d, want 1", len(posts))
	}
	if posts[0].Title != "Mine" {
		t.Errorf("title = %q, want Mine", posts[0].Title)
	}
}

func TestGetLikedPostsByUserID(t *testing.T) {
	db := setupTestDB(t)
	userID := createTestUser(t, db)

	postOneID, err := CreatePost(db, models.Post{UserID: userID, Title: "Liked", Content: "liked post"}, []string{"Go"})
	if err != nil {
		t.Fatalf("CreatePost liked: %v", err)
	}
	postTwoID, err := CreatePost(db, models.Post{UserID: userID, Title: "Not Liked", Content: "other post"}, []string{"Go"})
	if err != nil {
		t.Fatalf("CreatePost other: %v", err)
	}

	if err := AddPostReaction(db, userID, postOneID, 1); err != nil {
		t.Fatalf("AddPostReaction: %v", err)
	}
	if err := AddPostReaction(db, userID, postTwoID, -1); err != nil {
		t.Fatalf("AddPostReaction dislike: %v", err)
	}

	posts, err := GetLikedPostsByUserID(db, userID)
	if err != nil {
		t.Fatalf("GetLikedPostsByUserID: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("len(posts) = %d, want 1", len(posts))
	}
	if posts[0].Title != "Liked" {
		t.Errorf("title = %q, want Liked", posts[0].Title)
	}
}
