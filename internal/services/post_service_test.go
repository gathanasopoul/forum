package services

import (
	"forum/internal/models"
	"testing"
)

func TestCreatePostAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	userID := createTestUser(t, db)

	categories := []string{"Go", "Web"}
	post := models.Post{
		UserID:  userID,
		Title:   "First Post",
		Content: "Hello World",
	}

	postID, err := CreatePost(db, post, categories)
	if err != nil {
		t.Fatalf("CreatePost: %v", err)
	}
	if postID == 0 {
		t.Fatal("expected non-zero post ID")
	}

	found, err := GetPostByID(db, postID)
	if err != nil {
		t.Fatalf("GetPostByID: %v", err)
	}
	if found == nil {
		t.Fatal("expected post, got nil")
	}
	if found.Title != "First Post" {
		t.Errorf("title = %q, want First Post", found.Title)
	}
	if found.Content != "Hello World" {
		t.Errorf("content = %q, want Hello World", found.Content)
	}
	if found.Username != "testuser" {
		t.Errorf("username = %q, want testuser", found.Username)
	}
	if len(found.Categories) != 2 {
		t.Errorf("len(categories) = %d, want 2", len(found.Categories))
	}
}

func TestGetAllPosts(t *testing.T) {
	db := setupTestDB(t)
	userID := createTestUser(t, db)

	p1 := models.Post{UserID: userID, Title: "Title 1", Content: "Content 1"}
	p2 := models.Post{UserID: userID, Title: "Title 2", Content: "Content 2"}

	if _, err := CreatePost(db, p1, []string{"General"}); err != nil {
		t.Fatalf("CreatePost 1: %v", err)
	}
	if _, err := CreatePost(db, p2, []string{"Help"}); err != nil {
		t.Fatalf("CreatePost 2: %v", err)
	}

	posts, err := GetAllPosts(db)
	if err != nil {
		t.Fatalf("GetAllPosts: %v", err)
	}
	if len(posts) != 2 {
		t.Errorf("len(posts) = %d, want 2", len(posts))
	}
}
